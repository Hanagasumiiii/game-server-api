package database

import (
	"database/sql"
	_ "github.com/lib/pq"
)

func NewConnection(connectionString string) *sql.DB {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic("failed to connect to database: " + err.Error())
	}
	return db
}
