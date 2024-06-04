package types

import (
	"context"
	"database/sql"
	"time"
)

type InventoryRepo interface {
	GetInventories(ctx context.Context, userID int) (*[]TireInventory, error)
	GetInventoryByID(ctx context.Context, userID int, inventoryID int) (*TireInventory, error)
	CreateInventory(ctx context.Context, userId int, inventory *TireInventory) error
}

type Image struct {
	ID        ImageID      `json:"id"`
	Type      string       `json:"type"`
	Size      uint64       `json:"size"`
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt sql.NullTime `json:"updatedAt"`
	DeletedAt sql.NullTime `json:"deletedAt"`
}

type TireModel struct {
	ID        int            `json:"id" db:"tire_model_id"`
	Brand     sql.NullString `json:"brand" db:"tire_brand"`
	Type      sql.NullString `json:"type"`
	Size      sql.NullString `json:"size"`
	DOT       sql.NullString `json:"dot" db:"dot"`
	CreatedAt time.Time      `json:"createdAt" db:"tire_model_created_at"`
}

type TireInventory struct {
	ID      int  `json:"id" db:"tire_inventory_id"`
	UserID  int  `json:"userId" db:"user_id"`
	IsSaved bool `json:"isSaved" db:"is_saved"`

	TireModel `json:"tireModel"`
	CarDetail `json:"carDetail"`

	CreatedAt time.Time    `json:"createdAt" db:"tire_inventory_created_at"`
	UpdatedAt sql.NullTime `json:"updatedAt" db:"tire_inventory_updated_at"`
	DeletedAt sql.NullTime `json:"deletedAt" db:"tire_inventory_deleted_at"`
}

type CarDetail struct {
	ID              sql.NullInt32  `json:"id" db:"car_detail_id"`
	TireInventoryID sql.NullInt32  `json:"tireInventoryId"`
	Brand           sql.NullString `json:"brand" db:"car_brand"`
	Model           sql.NullString `json:"model"`
	Year            sql.NullString `json:"year"`
	LicensePlate    sql.NullString `json:"licensePlate"`
	Color           sql.NullString `json:"color"`

	CreatedAt sql.NullTime `json:"createdAt" db:"car_detail_created_at"`
	UpdatedAt sql.NullTime `json:"updatedAt"`
	DeletedAt sql.NullTime `json:"deletedAt"`
}
