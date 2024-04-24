package dbClient

import (
	"os"
	"sync"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

var (
	client  *bun.DB
	initErr error
	once    sync.Once
)

func Init() (*bun.DB, error) {
	once.Do(func() {
		config, err := pgx.ParseConfig(os.Getenv("DATABASE_URL"))
		if err != nil {
			panic(err)
		}
		config.PreferSimpleProtocol = true

		sqldb := stdlib.OpenDB(*config)
		client = bun.NewDB(sqldb, pgdialect.New())
	})

	return client, initErr
}

func GetClient() *bun.DB {
	return client
}

func CleanUp() {
	if client != nil {
		client.Close()
	}
}
