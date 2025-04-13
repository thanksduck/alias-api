package db

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	q "github.com/thanksduck/alias-api/internal/db"
)

var (
	pool *pgxpool.Pool
	once sync.Once
)

var SQL *q.Queries

// DB represents the database client for transaction support
var DB *DatabaseClient

// DatabaseClient wraps the pgxpool to provide transaction methods
type DatabaseClient struct {
	Pool *pgxpool.Pool
}

// Begin starts a new transaction
func (d *DatabaseClient) Begin(ctx context.Context) (pgx.Tx, error) {
	return d.Pool.Begin(ctx)
}

// Connect initializes the database connection pool
func Connect() {
	once.Do(func() {
		var err error
		pool, err = pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
		if err != nil {
			log.Fatalf("Unable to create connection pool: %v\n", err)
		}
		log.Println("Connected to database")

		// Initialize the DB client
		DB = &DatabaseClient{
			Pool: pool,
		}
	})
}

// GetPool returns the database connection pool
func GetPool() *pgxpool.Pool {
	if pool == nil {
		Connect()
	}
	return pool
}

// InitQueries initializes the SQL queries with the connection pool
func InitQueries() {
	SQL = q.New(GetPool())
}
