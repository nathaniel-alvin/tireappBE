package db

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/nathaniel-alvin/tireappBE/types"
)

type UserDatabase struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserDatabase {
	return &UserDatabase{
		db: db,
	}
}

func (s *UserDatabase) GetUserByEmail(email string) (*types.User, error) {
	users := []types.User{}
	err := s.db.Select(&users, "SELECT * FROM account WHERE email = $1", email)
	if err != nil {
		return nil, err
	}

	if len(users) > 1 {
		return nil, fmt.Errorf("GetUserByEmail: data error. more than one user found")
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("GetUserByEmail: data error. no user found")
	}

	return &users[0], nil
}

func (s *UserDatabase) GetUserById(id int) (*types.User, error) {
	user := types.User{}
	err := s.db.Select(&user, "SELECT * FROM account where id = $1", id)
	if err != nil {
		return nil, err
	}
	return &user, err
}

func (s *UserDatabase) CreateUser(u types.User) (int, error) {
	// begin transaction
	tx, err := s.db.Beginx()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var userId int

	// insert into table and return the id
	err = tx.QueryRowx("INSERT INTO account (email, password, display_name) VALUES ($1, $2, $3) RETURNING id;", u.Email, u.Password, u.Username).Scan(&userId)
	if err != nil {
		return 0, err
	}

	// commit if all operation are successful
	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return userId, nil
}

func getUserByEmail(ctx context.Context)
