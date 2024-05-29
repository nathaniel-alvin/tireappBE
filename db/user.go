package db

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/nathaniel-alvin/tireappBE/types"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (s *UserRepo) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {
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

func (s *UserRepo) GetUserById(ctx context.Context, id int) (*types.User, error) {
	user := types.User{}
	err := s.db.Get(&user, "SELECT * FROM account where id = $1", id)
	if err != nil {
		return nil, err
	}
	return &user, err
}

func (s *UserRepo) CreateUser(ctx context.Context, u *types.User) error {
	// begin transaction
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// insert into table and return the id
	tx.MustExec("INSERT INTO account (email, password, display_name) VALUES ($1, $2, $3);", u.Email, u.Password, u.Username)

	// commit if all operation are successful
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func getUserByEmail(ctx context.Context, email string) (*types.User, error) {
	return nil, nil
}
