package leaderboard

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nathaniel-alvin/tireappBE/service/auth"
	"github.com/nathaniel-alvin/tireappBE/types"
	"github.com/nathaniel-alvin/tireappBE/utils"
)

type Handler struct {
	store     types.LeaderboardRepo
	userStore types.UserRepo
}

func NewHandler(store types.LeaderboardRepo, userStore types.UserRepo) *Handler {
	return &Handler{
		store:     store,
		userStore: userStore,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/leaderboard", auth.WithJWTAuth(h.handleInventoryLeaderboard, h.userStore)).Methods("GET")
}

func (h *Handler) handleInventoryLeaderboard(w http.ResponseWriter, r *http.Request) {
	tires, err := h.store.GetTireModelLeaderboard(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.Encode(w, http.StatusOK, tires)
}
