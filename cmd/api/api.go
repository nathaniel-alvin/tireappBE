package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/nathaniel-alvin/tireappBE/db"
	"github.com/nathaniel-alvin/tireappBE/service/inventory"
	"github.com/nathaniel-alvin/tireappBE/service/leaderboard"
	"github.com/nathaniel-alvin/tireappBE/service/user"
)

type APIServer struct {
	addr string
	db   *sqlx.DB
}

func NewAPIServer(addr string, db *sqlx.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userStore := db.NewUserRepo(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)

	inventoryStore := db.NewInventoryRepo(s.db)
	inventoryHandler := inventory.NewHandler(inventoryStore, userStore)
	inventoryHandler.RegisterRoutes(subrouter)

	leaderboardStore := db.NewLeaderboardRepo(s.db)
	leaderboardHandler := leaderboard.NewHandler(leaderboardStore, userStore)
	leaderboardHandler.RegisterRoutes(subrouter)

	log.Println("Listening on", s.addr)

	return http.ListenAndServe(s.addr, router)
}
