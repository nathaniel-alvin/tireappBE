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
	req, err := utils.Decode[types.LoginUserRequest](r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate request
	if err := utils.Validate.Struct(req); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid request payload %v", errors))
		return
	}

	// check user exists
	u, err := h.store.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("incorrect email or password"))
		return
	}

	if !auth.ComparePassword(u.Password, []byte(req.Password)) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("incorrect email or password"))
		return
	}

	tokenStr, err := auth.CreateTokensAndSetCookies(w, u.ID, config.Envs.JWTExpirationInSeconds)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.Encode(w, http.StatusOK, map[string]string{"token": tokenStr})
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	req, err := utils.Decode[types.RegisterUserRequest](r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate request
	if err := utils.Validate.Struct(req); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid request payload %v", errors))
		return
	}

	// check if email has been used
	if _, err := h.store.GetUserByEmail(r.Context(), req.Email); err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with email %s already exixsts", req.Email))
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("fail to hash password: %w", err))
	}

	if err = h.store.CreateUser(r.Context(), &types.User{
		ID:         0,
		Email:      req.Email,
		Password:   hashedPassword,
		Username:   req.Username,
		Active:     true,
		ProfileUrl: sql.NullString{},
	}); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.Encode(w, http.StatusCreated, nil)
}

func (h *Handler) handleRefresh(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			utils.WriteError(w, http.StatusUnauthorized, err)
			return
		}
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	err = auth.TokenRefresher(w, c, config.Envs.JWTExpirationInSeconds)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.Encode(w, http.StatusOK, nil)
}
