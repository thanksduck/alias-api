package db

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pool *pgxpool.Pool
	once sync.Once
)

// Connect initializes the database connection pool
func Connect() {
	once.Do(func() {
		var err error
		pool, err = pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
		if err != nil {
			log.Fatalf("Unable to create connection pool: %v\n", err)
		}
		log.Println("Connected to database")
	})
}

// GetPool returns the database connection pool
func GetPool() *pgxpool.Pool {
	if pool == nil {
		Connect()
	}
	return pool
}
