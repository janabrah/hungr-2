package storage

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

func Init(connString string) error {
	var err error
	db, err = pgxpool.New(context.Background(), connString)
	return err
}

// Tx wraps a pgx transaction for use in handlers
type Tx struct {
	tx pgx.Tx
}

// BeginTx starts a new transaction
func BeginTx(ctx context.Context) (*Tx, error) {
	tx, err := db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return &Tx{tx: tx}, nil
}

// Commit commits the transaction
func (t *Tx) Commit(ctx context.Context) error {
	return t.tx.Commit(ctx)
}

// Rollback rolls back the transaction
func (t *Tx) Rollback(ctx context.Context) error {
	return t.tx.Rollback(ctx)
}
