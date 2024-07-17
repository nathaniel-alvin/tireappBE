package types

import (
	"context"
	"database/sql"
)

type UserRepo interface {
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	CreateUser(ctx context.Context, user *User) error
	GetUserById(ctx context.Context, id int) (*User, error)
	EditUser(ctx context.Context, userID int, user *User) error
}

type User struct {
	ID int `json:"id" db:"id"`
	// Email      string         `json:"email" db:"email"`
	Password   string         `json:"password" db:"password"`
	Username   string         `json:"userName" db:"display_name"`
	ProfileUrl sql.NullString `json:"profileUrl" db:"profile_url"`
	Active     bool           `json:"active" db:"active"`
}
