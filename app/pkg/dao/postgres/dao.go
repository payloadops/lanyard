package postgresdao

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type DAO struct {
	db *sql.DB
}

func NewDAO(dataSourceName string) *DAO {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	return &DAO{db: db}
}
