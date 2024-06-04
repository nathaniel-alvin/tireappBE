package inventory

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	tireappbe "github.com/nathaniel-alvin/tireappBE"
	"github.com/nathaniel-alvin/tireappBE/service/auth"
	"github.com/nathaniel-alvin/tireappBE/types"
	"github.com/nathaniel-alvin/tireappBE/utils"
)

type Handler struct {
	store     types.InventoryRepo
	userStore types.UserRepo
}

func NewHandler(store types.InventoryRepo, userStore types.UserRepo) *Handler {
	return &Handler{
		store:     store,
		userStore: userStore,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/inventories", auth.WithJWTAuth(h.handleInventoryIndex, h.userStore)).Methods("GET")
	router.HandleFunc("/inventories/{id}", h.handleInventoryView).Methods("GET")

	router.HandleFunc("/inventories", h.handleInventoryCreate).Methods("POST")
	router.HandleFunc("/inventories/{id}", h.handleInventoryEdit).Methods("PUT")
	router.HandleFunc("/inventories/{id}", h.handleInventoryDelete).Methods("DELETE")

	router.HandleFunc("/inventories/{id}/car-details", h.handleInventorySetCarDetail).Methods("POST")
	router.HandleFunc("/inventories/{id}/workshops", h.handleInventorySetWorkshop).Methods("POST")
}

func (h *Handler) handleInventoryIndex(w http.ResponseWriter, r *http.Request) {
	userID := tireappbe.UserIDFromContext(r.Context())

	tires, err := h.store.GetInventories(r.Context(), userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.Encode(w, http.StatusOK, tires)
}

func (h *Handler) handleInventoryView(w http.ResponseWriter, r *http.Request) {
	userID := tireappbe.UserIDFromContext(r.Context())

	vars := mux.Vars(r)
	str, ok := vars["id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing inventory ID"))
	}

	invID, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid inventory ID"))
	}

	tire, err := h.store.GetInventoryByID(r.Context(), userID, invID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.Encode(w, http.StatusOK, tire)
}

func (h *Handler) handleInventoryCreate(w http.ResponseWriter, r *http.Request) {
	userID := tireappbe.UserIDFromContext(r.Context())

	req, err := utils.Decode[types.TireInventory](r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err = h.store.CreateInventory(r.Context(), userID, &req)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.Encode(w, http.StatusOK, nil)
}
func (h *Handler) handleInventoryEdit(w http.ResponseWriter, r *http.Request)         {}
func (h *Handler) handleInventoryDelete(w http.ResponseWriter, r *http.Request)       {}
func (h *Handler) handleInventorySetCarDetail(w http.ResponseWriter, r *http.Request) {}
func (h *Handler) handleInventorySetWorkshop(w http.ResponseWriter, r *http.Request)  {}
