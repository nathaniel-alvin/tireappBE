package types

import (
	"context"
	"database/sql"
)

type User struct {
	ID         int            `json:"id" db:"id"`
	Email      string         `json:"email" db:"email"`
	Password   string         `json:"password" db:"password"`
	Username   string         `json:"userName" db:"display_name"`
	ProfileUrl sql.NullString `json:"profileUrl" db:"profile_url"`
	Active     bool           `json:"active" db:"active"`
}

type UserRepo interface {
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	CreateUser(ctx context.Context, user *User) error
	GetUserById(ctx context.Context, id int) (*User, error)
}
