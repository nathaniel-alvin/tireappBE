package main

import (
	"log"

	"github.com/nathaniel-alvin/tireappBE/cmd/api"
	"github.com/nathaniel-alvin/tireappBE/db"
)

func main() {
	db, err := db.NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	server := api.NewAPIServer(":8080", db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
