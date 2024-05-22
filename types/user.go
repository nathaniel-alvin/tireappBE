package types

import "database/sql"

type User struct {
	ID         int            `json:"id"`
	Email      string         `json:"email"`
	Password   string         `json:"password"`
	Username   string         `json:"userName" db:"display_name"`
	ProfileUrl sql.NullString `json:"profileUrl" db:"profile_url"`
	Active     bool           `json:"active"`
}
