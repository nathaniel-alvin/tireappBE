package user

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/nathaniel-alvin/tireappBE/config"
	"github.com/nathaniel-alvin/tireappBE/service/auth"
	"github.com/nathaniel-alvin/tireappBE/types"
	"github.com/nathaniel-alvin/tireappBE/utils"

	tireapperror "github.com/nathaniel-alvin/tireappBE/error"
)

type Handler struct {
	store types.UserRepo
}

func NewHandler(store types.UserRepo) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/register", h.handleRegister).Methods("POST")

	router.HandleFunc("/refresh", auth.WithJWTAuth(h.handleRefresh, h.store)).Methods("GET")
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	type LoginUserRequest struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	req, err := utils.Decode[LoginUserRequest](r)
	if err != nil {
		utils.WriteError(w, tireapperror.Errorf(tireapperror.EINVALID, "invalid JSON body"))
		return
	}

	// validate request
	if err := utils.Validate.Struct(req); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, tireapperror.Errorf(tireapperror.EINVALID, errors.Error()))
		return
	}

	// check user exists
	u, err := h.store.GetUserByUsername(r.Context(), req.Username)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	if !auth.ComparePassword(u.Password, []byte(req.Password)) {
		utils.WriteError(w, tireapperror.Errorf(tireapperror.EINVALID, "incorrect username or password"))
		return
	}

	tokenStr, err := auth.CreateTokensAndSetCookies(w, u.ID, config.Envs.JWTExpirationInSeconds)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.Encode(w, http.StatusOK, map[string]string{"token": tokenStr})
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	type RegisterUserRequest struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required,min=3,max=130"`
	}

	req, err := utils.Decode[RegisterUserRequest](r)
	if err != nil {
		utils.WriteError(w, tireapperror.Errorf(tireapperror.EINVALID, "invalid JSON body"))
		return
	}

	// validate request
	if err := utils.Validate.Struct(req); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, tireapperror.Errorf(tireapperror.EINVALID, fmt.Sprintf("invalid request payload %v", errors.Error())))
		return
	}

	// check if username has been used
	if _, err := h.store.GetUserByUsername(r.Context(), req.Username); err == nil {
		utils.WriteError(w, tireapperror.Errorf(tireapperror.EINVALID, fmt.Sprintf("user with username %s already exixsts", req.Username)))
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		utils.WriteError(w, tireapperror.Errorf(tireapperror.EINVALID, fmt.Sprintf("fail to hash password: %w", err)))
	}

	if err = h.store.CreateUser(r.Context(), &types.User{
		ID:         0,
		Password:   hashedPassword,
		Username:   req.Username,
		Active:     true,
		ProfileUrl: sql.NullString{},
	}); err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.Encode(w, http.StatusCreated, nil)
}

func (h *Handler) handleRefresh(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			utils.WriteError(w, tireapperror.Errorf(tireapperror.EUNAUTHORIZED, "unauthorized"))
			return
		}
		utils.WriteError(w, tireapperror.Errorf(tireapperror.EINVALID, err.Error()))
		return
	}
	err = auth.TokenRefresher(w, c, config.Envs.JWTExpirationInSeconds)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.Encode(w, http.StatusOK, nil)
}
