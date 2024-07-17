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

	tireappbe "github.com/nathaniel-alvin/tireappBE"
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
	router.HandleFunc("/user/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/user/register", h.handleRegister).Methods("POST")
	router.HandleFunc("/user/edituser", auth.WithJWTAuth(h.handleEditUser, h.store)).Methods("POST")
	router.HandleFunc("/user/info", auth.WithJWTAuth(h.handleGetUserDetails, h.store)).Methods("GET")

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
		utils.WriteError(w, tireapperror.Errorf(tireapperror.EINVALID, "Authentication Failed: %v", errors.Error()))
		return
	}

	// check user exists
	u, err := h.store.GetUserByUsername(r.Context(), req.Username)
	if err != nil {
		utils.WriteError(w, tireapperror.Errorf(tireapperror.EUNAUTHORIZED, "incorrect username or password"))
		return
	}

	if !auth.ComparePassword(u.Password, []byte(req.Password)) {
		utils.WriteError(w, tireapperror.Errorf(tireapperror.EUNAUTHORIZED, "incorrect username or password"))
		return
	}

	tokenStr, err := auth.CreateTokensAndSetCookies(w, u.ID, config.Envs.JWTExpirationInSeconds)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	type loginResponse struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
		IsActive bool   `json:"active"`
	}

	utils.Encode(w, http.StatusOK, map[string]any{
		"token": tokenStr,
		"user": loginResponse{
			ID:       u.ID,
			Username: u.Username,
			IsActive: u.Active,
		},
	})
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

	user, err := h.store.GetUserByUsername(r.Context(), req.Username)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	type createUserResponse struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
		IsActive bool   `json:"active"`
	}

	utils.Encode(w, http.StatusCreated, createUserResponse{
		ID:       user.ID,
		Username: user.Username,
		IsActive: user.Active,
	})
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
	newToken, err := auth.TokenRefresher(w, c, config.Envs.JWTExpirationInSeconds)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.Encode(w, http.StatusOK, map[string]string{"token": newToken})
}

func (h *Handler) handleEditUser(w http.ResponseWriter, r *http.Request) {
	type EditUserRequest struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required,min=3,max=130"`
		Active   bool   `json:"isActive" validate:"required"`
	}

	userID := tireappbe.UserIDFromContext(r.Context())

	req, err := utils.Decode[EditUserRequest](r)
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
	// if _, err = h.store.GetUserByUsername(r.Context(), req.Username); err == nil {
	// 	utils.WriteError(w, tireapperror.Errorf(tireapperror.EINVALID, fmt.Sprintf("user with username %s already exixsts", req.Username)))
	// 	return
	// }

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		utils.WriteError(w, tireapperror.Errorf(tireapperror.EINVALID, fmt.Sprintf("fail to hash password: %w", err)))
		return
	}

	err = h.store.EditUser(r.Context(), userID, &types.User{
		Username: req.Username,
		Password: hashedPassword,
		Active:   req.Active,
	})
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	user, err := h.store.GetUserById(r.Context(), userID)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	type getUserResponse struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
		IsActive bool   `json:"active"`
	}

	utils.Encode(w, http.StatusOK, getUserResponse{
		ID:       user.ID,
		Username: user.Username,
		IsActive: user.Active,
	})
}

func (h *Handler) handleGetUserDetails(w http.ResponseWriter, r *http.Request) {
	userID := tireappbe.UserIDFromContext(r.Context())

	user, err := h.store.GetUserById(r.Context(), userID)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	type getUserResponse struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
		IsActive bool   `json:"active"`
	}

	utils.Encode(w, http.StatusOK, getUserResponse{
		ID:       user.ID,
		Username: user.Username,
		IsActive: user.Active,
	})
}
