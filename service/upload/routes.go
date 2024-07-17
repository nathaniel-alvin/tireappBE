package upload

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	tireappbe "github.com/nathaniel-alvin/tireappBE"
	"github.com/nathaniel-alvin/tireappBE/service/auth"
	"github.com/nathaniel-alvin/tireappBE/types"
	"github.com/nathaniel-alvin/tireappBE/utils"

	tireapperror "github.com/nathaniel-alvin/tireappBE/error"
)

type Handler struct {
	store     types.UploadStore
	userStore types.UserRepo
}

func NewHandler(store types.UploadStore, userStore types.UserRepo) *Handler {
	return &Handler{
		store:     store,
		userStore: userStore,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/upload/{id}", auth.WithJWTAuth(h.handleUploadForInventory, h.userStore)).Methods("POST")
	router.HandleFunc("/upload/{id}", auth.WithJWTAuth(h.handleUpdateForInventory, h.userStore)).Methods("PUT")
	router.HandleFunc("/upload", auth.WithJWTAuth(h.handleUpload, h.userStore)).Methods("POST")
}

func (h *Handler) handleUpload(w http.ResponseWriter, r *http.Request) {
	userID := tireappbe.UserIDFromContext(r.Context())

	_, err := h.store.InsertFileFromRequest(r, userID)
	if err != nil {
		utils.WriteError(w, err)
		return
	}
	// TODO: add tire model details pake method yang dari inventoryStore
	// CANCEL: gk jadi karena update ntar di button beda aja
	// h.inventoryStore.UpdateTireModel(r.Context(), inventoryID, req.TireModel)
	utils.Encode(w, http.StatusOK, nil)
}

func (h *Handler) handleUploadForInventory(w http.ResponseWriter, r *http.Request) {
	// Get inventory ID
	vars := mux.Vars(r)
	str, ok := vars["id"]
	if !ok {
		utils.WriteError(w, tireapperror.Errorf(tireapperror.EINVALID, "invalid ID format"))
		return
	}

	invID, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, tireapperror.Errorf(tireapperror.EINVALID, "invalid ID format"))
		return
	}

	err = h.store.CreateImageForInventory(r, invID)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.Encode(w, http.StatusCreated, nil)
}

func (h *Handler) handleUpdateForInventory(w http.ResponseWriter, r *http.Request) {
	// Get inventory ID
	vars := mux.Vars(r)
	str, ok := vars["id"]
	if !ok {
		utils.WriteError(w, tireapperror.Errorf(tireapperror.EINVALID, "invalid ID format"))
		return
	}

	invID, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, tireapperror.Errorf(tireapperror.EINVALID, "invalid ID format"))
		return
	}

	err = h.store.UpdateImageForInventory(r, invID)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.Encode(w, http.StatusOK, nil)
}
