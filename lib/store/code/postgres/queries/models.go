// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package queries

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Datum struct {
	ID          int32
	CodeBlock   pgtype.Text
	FilePath    pgtype.Text
	Line        pgtype.Int4
	Description pgtype.Text
}
