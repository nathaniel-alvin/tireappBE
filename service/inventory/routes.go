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
	router.HandleFunc("/inventories/history", auth.WithJWTAuth(h.handleInventoryHistory, h.userStore)).Methods("GET")

	router.HandleFunc("/inventories/{id}/car-detail", auth.WithJWTAuth(h.handleInventorySetCarDetail, h.userStore)).Methods("POST")
	router.HandleFunc("/inventories/{id}/tire-model", auth.WithJWTAuth(h.handleInventorySetTireModel, h.userStore)).Methods("POST")
	// router.HandleFunc("/inventories/{id}/workshops", h.handleInventorySetWorkshop).Methods("POST")

	router.HandleFunc("/inventories/{id}/upload", auth.WithJWTAuth(h.handleInventoryAddImage, h.userStore)).Methods("POST")

	router.HandleFunc("/inventories/{id}", auth.WithJWTAuth(h.handleInventoryDelete, h.userStore)).Methods("DELETE")
	router.HandleFunc("/inventories/{id}", auth.WithJWTAuth(h.handleInventoryView, h.userStore)).Methods("GET")

	router.HandleFunc("/inventories", auth.WithJWTAuth(h.handleInventoryIndex, h.userStore)).Methods("GET")
	router.HandleFunc("/inventories", auth.WithJWTAuth(h.handleInventoryCreate, h.userStore)).Methods("POST")
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

	// Get inventory ID
	vars := mux.Vars(r)
	str, ok := vars["id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing inventory ID"))
		return
	}

	invID, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid inventory ID"))
		return
	}

	tire, err := h.store.GetInventoryByID(r.Context(), userID, invID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.Encode(w, http.StatusOK, tire)
}

func (h *Handler) handleInventoryCreate(w http.ResponseWriter, r *http.Request) {
	type CreateInventoryRequest struct {
		TireModel     types.TireModel     `json:"tireModel"`
		TireInventory types.TireInventory `json:"tireInventory"`
	}
	userID := tireappbe.UserIDFromContext(r.Context())

	req, err := utils.Decode[CreateInventoryRequest](r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	err = h.store.CreateInventory(r.Context(), userID, &req.TireInventory, &req.TireModel)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.Encode(w, http.StatusCreated, nil)
}

func (h *Handler) handleInventoryDelete(w http.ResponseWriter, r *http.Request) {
	// Get inventory ID
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing inventory ID"))
		return
	}

	inventoryID, err := strconv.Atoi(id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid inventory ID"))
		return
	}

	err = h.store.DeleteInventory(r.Context(), inventoryID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.Encode(w, http.StatusNoContent, nil)
}

func (h *Handler) handleInventorySetCarDetail(w http.ResponseWriter, r *http.Request) {
	type UpdateCarDetailRequest struct {
		Make         string `json:"make"`
		Model        string `json:"model"`
		Year         string `json:"year"`
		LicensePlate string `json:"licensePlate"`
		Color        string `json:"color"`
	}

	// Get inventory ID
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing inventory ID"))
		return
	}

	inventoryID, err := strconv.Atoi(id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid inventory ID"))
		return
	}

	// Decode req
	req, err := utils.Decode[UpdateCarDetailRequest](r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Map request to TireModel struct
	carDetail := types.CarDetail{
		Make:         utils.NewNullString(req.Make),
		Model:        utils.NewNullString(req.Model),
		Year:         utils.NewNullString(req.Year),
		LicensePlate: utils.NewNullString(req.LicensePlate),
		Color:        utils.NewNullString(req.Color),
	}

	err = h.store.UpdateCarDetail(r.Context(), inventoryID, carDetail)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.Encode(w, http.StatusOK, nil)
}
func (h *Handler) handleInventorySetWorkshop(w http.ResponseWriter, r *http.Request) {}

func (h *Handler) handleInventorySetTireModel(w http.ResponseWriter, r *http.Request) {
	type UpdateTireModelRequest struct {
		Brand string `json:"brand"`
		Type  string `json:"type"`
		Size  string `json:"size"`
		DOT   string `json:"dot"`
	}

	// Get inventory ID
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing inventory ID"))
		return
	}

	inventoryID, err := strconv.Atoi(id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid inventory ID"))
		return
	}

	// Decode req
	req, err := utils.Decode[UpdateTireModelRequest](r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Map request to TireModel struct
	tireModel := types.TireModel{
		Brand: utils.NewNullString(req.Brand),
		Type:  utils.NewNullString(req.Type),
		Size:  utils.NewNullString(req.Size),
		DOT:   utils.NewNullString(req.DOT),
	}

	err = h.store.UpdateTireModel(r.Context(), inventoryID, tireModel)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.Encode(w, http.StatusOK, nil)
}

// TODO
func (h *Handler) handleInventoryAddImage(w http.ResponseWriter, r *http.Request) {
	type AddImageRequest struct {
		Filename string `json:"filename"`
	}
}

func (h *Handler) handleInventoryHistory(w http.ResponseWriter, r *http.Request) {
	userID := tireappbe.UserIDFromContext(r.Context())

	tires, err := h.store.GetInventoryHistory(r.Context(), userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.Encode(w, http.StatusOK, tires)
}
