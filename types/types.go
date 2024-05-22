package types

import (
	"net/http"
	"time"
)

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	CreateUser(User) (int, error)
	GetUserById(id int) (*User, error)
}

type UploadStore interface {
	InsertFileFromRequest(r *http.Request, userID int) (ImageID, error)
}

type (
	ImageID  int
	Filename string
	Filetype string
)

type RegisterUserRequest struct {
	Username string `json:"userName" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3,max=130"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UploadPictureRequest struct {
	Filename string `json:"filename"`
}

type UpdateTireModelRequest struct {
	Brand string `json:"brand"`
	Type  string `json:"type"`
	Size  string `json:"size"`
	DOT   string `json:"dot"`
}
