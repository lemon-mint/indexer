// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: queries.sql

package queries

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createData = `-- name: CreateData :one
INSERT INTO data (code_block, file_path, line, description) VALUES ($1, $2, $3, $4) RETURNING id, code_block, file_path, line, description
`

type CreateDataParams struct {
	CodeBlock   pgtype.Text
	FilePath    pgtype.Text
	Line        pgtype.Int4
	Description pgtype.Text
}

func (q *Queries) CreateData(ctx context.Context, arg CreateDataParams) (Datum, error) {
	row := q.db.QueryRow(ctx, createData,
		arg.CodeBlock,
		arg.FilePath,
		arg.Line,
		arg.Description,
	)
	var i Datum
	err := row.Scan(
		&i.ID,
		&i.CodeBlock,
		&i.FilePath,
		&i.Line,
		&i.Description,
	)
	return i, err
}

const deleteData = `-- name: DeleteData :one
DELETE FROM data WHERE id = $1 RETURNING id, code_block, file_path, line, description
`

func (q *Queries) DeleteData(ctx context.Context, id int32) (Datum, error) {
	row := q.db.QueryRow(ctx, deleteData, id)
	var i Datum
	err := row.Scan(
		&i.ID,
		&i.CodeBlock,
		&i.FilePath,
		&i.Line,
		&i.Description,
	)
	return i, err
}

const deleteDataList = `-- name: DeleteDataList :many
DELETE FROM data WHERE id = ANY($1::int[]) RETURNING id, code_block, file_path, line, description
`

func (q *Queries) DeleteDataList(ctx context.Context, dollar_1 []int32) ([]Datum, error) {
	rows, err := q.db.Query(ctx, deleteDataList, dollar_1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Datum
	for rows.Next() {
		var i Datum
		if err := rows.Scan(
			&i.ID,
			&i.CodeBlock,
			&i.FilePath,
			&i.Line,
			&i.Description,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getData = `-- name: GetData :one
SELECT id, code_block, file_path, line, description FROM data WHERE id = $1
`

func (q *Queries) GetData(ctx context.Context, id int32) (Datum, error) {
	row := q.db.QueryRow(ctx, getData, id)
	var i Datum
	err := row.Scan(
		&i.ID,
		&i.CodeBlock,
		&i.FilePath,
		&i.Line,
		&i.Description,
	)
	return i, err
}

const getDataList = `-- name: GetDataList :many
SELECT id, code_block, file_path, line, description FROM data WHERE id = ANY($1::int[])
`

func (q *Queries) GetDataList(ctx context.Context, dollar_1 []int32) ([]Datum, error) {
	rows, err := q.db.Query(ctx, getDataList, dollar_1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Datum
	for rows.Next() {
		var i Datum
		if err := rows.Scan(
			&i.ID,
			&i.CodeBlock,
			&i.FilePath,
			&i.Line,
			&i.Description,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateData = `-- name: UpdateData :one
UPDATE data SET code_block = $2, file_path = $3, line = $4, description = $5 WHERE id = $1 RETURNING id, code_block, file_path, line, description
`

type UpdateDataParams struct {
	ID          int32
	CodeBlock   pgtype.Text
	FilePath    pgtype.Text
	Line        pgtype.Int4
	Description pgtype.Text
}

func (q *Queries) UpdateData(ctx context.Context, arg UpdateDataParams) (Datum, error) {
	row := q.db.QueryRow(ctx, updateData,
		arg.ID,
		arg.CodeBlock,
		arg.FilePath,
		arg.Line,
		arg.Description,
	)
	var i Datum
	err := row.Scan(
		&i.ID,
		&i.CodeBlock,
		&i.FilePath,
		&i.Line,
		&i.Description,
	)
	return i, err
}
