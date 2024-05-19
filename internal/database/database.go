package database

import (
	"database/sql"
	"fmt"
	"game-server-api/internal/config"
	_ "github.com/lib/pq"
)

func NewConnection(cfg config.Config) *sql.DB {
	connectionString := fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=disable", cfg.Postgres.Host, cfg.Postgres.Database, cfg.Postgres.User, cfg.Postgres.Password)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic("failed to connect to database: " + err.Error())
	}
	return db
}
