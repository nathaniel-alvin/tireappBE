package db

import (
	"github.com/jmoiron/sqlx"
	"github.com/nathaniel-alvin/tireappBE/types"
)

type InventoryRepo interface {
}

type InventoryDatabase struct {
	db *sqlx.DB
}

func NewInventoryRepo(db *sqlx.DB) *InventoryDatabase {
	return &InventoryDatabase{
		db: db,
	}
}
