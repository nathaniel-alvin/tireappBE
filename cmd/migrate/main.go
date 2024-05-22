package main

import (
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/nathaniel-alvin/tireappBE/db"
)

func main() {
	// connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", config.Envs.DBUser, config.Envs.DBPassword, config.Envs.PublicHost, config.Envs.Port, config.Envs.DBName)
	// m, err := migrate.New(
	// 	"file://cmd/migrate/migrations",
	// 	connStr,
	// )
	db, err := db.NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://cmd/migrate/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatal(err)
	}

	cmd := os.Args[(len(os.Args) - 1)]
	if cmd == "up" {
		if err := m.Up(); err != nil {
			log.Fatal(err)
		}
	}
	if cmd == "down" {
		if err := m.Down(); err != nil {
			log.Fatal(err)
		}
	}
}
