package qdrant

import (
	"context"
	"fmt"

	"github.com/qdrant/go-client/qdrant"

	"github.com/snowmerak/indexer/lib/index/vector"
	"github.com/snowmerak/indexer/pkg/box"
)

type Config struct {
	host       string
	port       int
	name       string
	apiKey     string
	useTLS     bool
	volatility bool
}

func NewConfig(host string, port int, name string) *Config {
	return &Config{
		host: host,
		port: port,
		name: name,
	}
}

func (c *Config) WithAPIKey(apiKey string) *Config {
	c.apiKey = apiKey
	return c
}

func (c *Config) WithTLS() *Config {
	c.useTLS = true
	return c
}

func (c *Config) WithVolatility() *Config {
	c.volatility = true
	return c
}

var _ vector.Vector = (*Vector)(nil)

type Vector struct {
	client *qdrant.Client
	config *Config
}

func New(ctx context.Context, cfg *Config) (*Vector, error) {
	qc := &qdrant.Config{
		Host: cfg.host,
		Port: cfg.port,
	}

	if cfg.apiKey != "" {
		qc.APIKey = cfg.apiKey
	}

	if cfg.useTLS {
		qc.UseTLS = true
	}

	client, err := qdrant.NewClient(qc)
	if err != nil {
		return nil, err
	}

	context.AfterFunc(ctx, func() {
		client.Close()
	})

	if err := client.CreateCollection(ctx, &qdrant.CreateCollection{
		CollectionName: cfg.name,
	}); err != nil {
		if cfg.volatility {
			return nil, fmt.Errorf("failed to create collection: %w", err)
		}
	}

	if cfg.volatility {
		context.AfterFunc(ctx, func() {
			count := 0
		loop:
			for {
				if err := client.DeleteCollection(context.Background(), cfg.name); err != nil {
					count++
					if count > 3 {
						break loop
					}
					continue loop
				}

				break loop
			}
		})
	}

	return &Vector{
		client: client,
		config: cfg,
	}, nil
}

func (v *Vector) Store(ctx context.Context, id int, payload *vector.Payload) error {
	nv := make([]float32, len(payload.Vector))
	for i, v := range payload.Vector {
		nv[i] = float32(v)
	}

	if _, err := v.client.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: v.config.name,
		Wait:           box.Wrap(true),
		Points: []*qdrant.PointStruct{
			{
				Id:      qdrant.NewIDNum(uint64(id)),
				Vectors: qdrant.NewVectors(nv...),
				Payload: map[string]*qdrant.Value{
					vector.PayloadRelatedId: qdrant.NewValueInt(int64(payload.RelatedId)),
				},
			},
		},
	}); err != nil {
		return fmt.Errorf("failed to upsert: %w", err)
	}

	return nil
}

func (v *Vector) Get(ctx context.Context, id int) (*vector.Payload, error) {
	pt, err := v.client.Get(ctx, &qdrant.GetPoints{
		CollectionName: v.config.name,
		Ids:            []*qdrant.PointId{qdrant.NewIDNum(uint64(id))},
		WithVectors:    qdrant.NewWithVectors(true),
		WithPayload:    qdrant.NewWithPayload(true),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get: %w", err)
	}

	if len(pt) == 0 {
		return nil, fmt.Errorf("point not found")
	}

	p := &vector.Payload{
		Id: id,
	}

	fpt := pt[0]
	vt := fpt.GetVectors().GetVector().Data
	nv := make([]float64, len(vt))

	p.Vector = nv

	for i, v := range vt {
		nv[i] = float64(v)
	}

	ri := fpt.GetPayload()[vector.PayloadRelatedId]
	if ri != nil {
		p.RelatedId = int(ri.GetIntegerValue())
	}

	return p, nil
}

func (v *Vector) Search(ctx context.Context, vt []float64, limit int) ([]*vector.Payload, error) {
	nv := make([]float32, len(vt))
	for i, v := range vt {
		nv[i] = float32(v)
	}

	resp, err := v.client.Query(ctx, &qdrant.QueryPoints{
		CollectionName: v.config.name,
		Query:          qdrant.NewQuery(nv...),
		Limit:          box.Wrap(uint64(limit)),
		WithPayload:    qdrant.NewWithPayload(true),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}

	ps := make([]*vector.Payload, len(resp))
	for i, pt := range resp {
		ri := pt.GetPayload()[vector.PayloadRelatedId]
		p := &vector.Payload{
			Id: int(pt.GetId().GetNum()),
		}

		if ri != nil {
			p.RelatedId = int(ri.GetIntegerValue())
		}

		ps[i] = p
	}

	return ps, nil
}

func (v *Vector) Delete(ctx context.Context, id int) error {
	if _, err := v.client.Delete(ctx, &qdrant.DeletePoints{
		CollectionName: v.config.name,
		Points:         qdrant.NewPointsSelector(qdrant.NewIDNum(uint64(id))),
	}); err != nil {
		return fmt.Errorf("failed to delete: %w", err)
	}

	return nil
}

func (v *Vector) Drop(ctx context.Context) error {
	if err := v.client.DeleteCollection(ctx, v.config.name); err != nil {
		return fmt.Errorf("failed to drop: %w", err)
	}

	return nil
}
