package types

import (
	"context"
	"database/sql"
)

type User struct {
	ID         int            `json:"id"`
	Email      string         `json:"email"`
	Password   string         `json:"password"`
	Username   string         `json:"userName" db:"display_name"`
	ProfileUrl sql.NullString `json:"profileUrl" db:"profile_url"`
	Active     bool           `json:"active"`
}

type UserRepo interface {
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	CreateUser(ctx context.Context, user *User) error
	GetUserById(ctx context.Context, id int) (*User, error)
}
