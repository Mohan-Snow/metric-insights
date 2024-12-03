package repo

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Postgres struct {
	db *pgxpool.Pool
}

func New(ctx context.Context) (*Postgres, error) {
	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		"postgres-db", "5432", "postgres", "postgres", "postgres")

	pool, err := pgxpool.Connect(ctx, connString)
	if err != nil {
		log.Println("Failed creating connection pool")
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		log.Println("Failed to ping db")
		return nil, err
	}

	return &Postgres{db: pool}, nil
}

func (p *Postgres) Close() {
	p.db.Close()
}
