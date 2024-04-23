package dbClient

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pool    *pgxpool.Pool
	initErr error
	once    sync.Once
)

func Init() (*pgxpool.Pool, error) {
	once.Do(func() {
		var err error
		config, err := pgxpool.ParseConfig("")
		if err != nil {
			initErr = fmt.Errorf("error parsing DATABASE_URL: %v", err)
			return
		}

		pool, err = pgxpool.NewWithConfig(context.Background(), config)
		if err != nil {
			initErr = fmt.Errorf("unable to connect to database: %v", err)
			return
		}
	})

	return pool, initErr
}

func GetPGClient() *pgxpool.Pool {
	return pool
}

func CleanUp() {
	if pool != nil {
		pool.Close()
	}
}
