package db

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/nathaniel-alvin/tireappBE/types"

	tireapperror "github.com/nathaniel-alvin/tireappBE/error"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (s *UserRepo) GetUserByUsername(ctx context.Context, username string) (*types.User, error) {
	users := []types.User{}
	err := s.db.Select(&users, "SELECT * FROM Account WHERE display_name = $1", username)
	if err != nil {
		return nil, tireapperror.Errorf(tireapperror.EINTERNAL, "query error")
	}

	if len(users) > 1 {
		return nil, tireapperror.Errorf(tireapperror.EINVALID, "data error. more than one user found")
	}

	if len(users) == 0 {
		return nil, tireapperror.Errorf(tireapperror.EINVALID, "data error. no user found")
	}

	return &users[0], nil
}

func (s *UserRepo) GetUserById(ctx context.Context, id int) (*types.User, error) {
	user := types.User{}
	err := s.db.Get(&user, "SELECT * FROM account where id = $1", id)
	if err != nil {
		return nil, tireapperror.Errorf(tireapperror.EINTERNAL, "db error: %v", err)
	}
	return &user, nil
}

func (s *UserRepo) CreateUser(ctx context.Context, u *types.User) error {
	// begin transaction
	tx, err := s.db.Beginx()
	if err != nil {
		return tireapperror.Errorf(tireapperror.EINTERNAL, "db error: %v", err)
	}
	defer tx.Rollback()

	// insert into table and return the id
	tx.MustExec("INSERT INTO Account (password, display_name) VALUES ($1, $2);", u.Password, u.Username)

	// commit if all operation are successful
	if err := tx.Commit(); err != nil {
		return tireapperror.Errorf(tireapperror.EINTERNAL, "db error: %v", err)
	}

	return nil
}

func (s *UserRepo) EditUser(ctx context.Context, userID int, u *types.User) error {
	query := `UPDATE account
        SET
			display_name = $2,
			password = $3,
			active = $4
        WHERE id = $1`

	_, err := s.db.Exec(query, userID, u.Username, u.Password, u.Active)
	if err != nil {
		return tireapperror.Errorf(tireapperror.EINTERNAL, "db error: %v", err)
	}

	return nil
}
