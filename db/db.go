package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/nathaniel-alvin/tireappBE/config"
)

func NewPostgresStore() (*sqlx.DB, error) {
	connStr := fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable", config.Envs.PublicHost, config.Envs.DBUser, config.Envs.DBPassword, config.Envs.DBName)
	// connStr := "user=postgres dbname=postgres password=tireapp sslmode=disable"
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
