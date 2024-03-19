package pg

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"time"
)

const (
	dbName    = "shortener"
	tableName = "shortener"
)

type PgRepo struct {
	DB *pgxpool.Pool
}

func NewConnect(dbDSN string) (PgRepo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(dbDSN)
	if err != nil {
		return PgRepo{}, err
	}

	db, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return PgRepo{}, err
	}

	var exists bool
	err = db.QueryRow(context.Background(), existDBQuery).Scan(&exists)
	if err != nil {
		return PgRepo{}, err
	}

	if !exists {
		_, err = db.Exec(ctx, createDBQuery)
		if err != nil {
			return PgRepo{}, err
		}
	}

	_, err = db.Exec(ctx, createTableQuery)
	if err != nil {
		return PgRepo{}, err
	}

	return PgRepo{db}, nil
}
