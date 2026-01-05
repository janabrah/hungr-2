package storage

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

func Init(connString string) error {
	var err error
	db, err = pgxpool.New(context.Background(), connString)
	return err
}
