package upload

import (
	"net/http"

	"github.com/gorilla/mux"
	tireappbe "github.com/nathaniel-alvin/tireappBE"
	"github.com/nathaniel-alvin/tireappBE/service/auth"
	"github.com/nathaniel-alvin/tireappBE/types"
	"github.com/nathaniel-alvin/tireappBE/utils"
)

type Handler struct {
	store          types.UploadStore
	userStore      types.UserRepo
	inventoryStore types.InventoryRepo
}

func NewHandler(store types.UploadStore, userStore types.UserRepo, inventoryStore types.InventoryRepo) *Handler {
	return &Handler{
		store:          store,
		userStore:      userStore,
		inventoryStore: inventoryStore,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	// Upload picture (ngambil data dari AI)
	router.HandleFunc("/upload", auth.WithJWTAuth(h.handleUpload, h.userStore)).Methods("POST")
}

func (h *Handler) handleUpload(w http.ResponseWriter, r *http.Request) {
	type UploadPictureRequest struct {
		Filename string `json:"filename"`
		// TireModel types.TireModel `json:"tireModel"`
	}

	userID := tireappbe.UserIDFromContext(r.Context())
	// req, err := utils.Decode[UploadPictureRequest](r)
	// if err != nil {
	// 	utils.WriteError(w, http.StatusBadRequest, err)
	// 	return
	// }

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
