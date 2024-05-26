package inventory

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nathaniel-alvin/tireappBE/db"
)

type Handler struct {
	store db.InventoryRepo
}

func NewHandler(store db.InventoryRepo) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/inventories", h.handleInventoryIndex).Methods("GET")
	router.HandleFunc("/inventories/{id}", h.handleInventoryView).Methods("GET")

	router.HandleFunc("/inventories", h.handleInventoryCreate).Methods("POST")
	router.HandleFunc("/inventories/{id}", h.handleInventoryEdit).Methods("PUT")
	router.HandleFunc("/inventories/{id}", h.handleInventoryDelete).Methods("DELETE")

	router.HandleFunc("/inventories/{id}/car-details", h.handleInventorySetCarDetail).Methods("POST")
	router.HandleFunc("/inventories/{id}/workshops", h.handleInventorySetWorkshop).Methods("POST")
}

func (h *Handler) handleInventoryIndex(w http.ResponseWriter, r *http.Request) {
	h.store.GetInventories(r.Context())
}
func (h *Handler) handleInventoryView(w http.ResponseWriter, r *http.Request)         {}
func (h *Handler) handleInventoryCreate(w http.ResponseWriter, r *http.Request)       {}
func (h *Handler) handleInventoryEdit(w http.ResponseWriter, r *http.Request)         {}
func (h *Handler) handleInventoryDelete(w http.ResponseWriter, r *http.Request)       {}
func (h *Handler) handleInventorySetCarDetail(w http.ResponseWriter, r *http.Request) {}
func (h *Handler) handleInventorySetWorkshop(w http.ResponseWriter, r *http.Request)  {}
