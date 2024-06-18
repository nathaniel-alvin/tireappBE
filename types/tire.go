package types

import (
	"context"
	"time"

	"github.com/nathaniel-alvin/tireappBE/utils"
)

type InventoryRepo interface {
	GetInventories(ctx context.Context, userID int) (*[]TireInventory, error)
	GetInventoryByID(ctx context.Context, userID int, inventoryID int) (*TireInventory, error)
	CreateInventory(ctx context.Context, userID int, i *TireInventory, m *TireModel) error
	UpdateTireModel(ctx context.Context, inventoryID int, tm TireModel) error
	UpdateCarDetail(ctx context.Context, inventoryID int, cd CarDetail) error
	DeleteInventory(ctx context.Context, inventoryID int) error
	GetInventoryHistory(ctx context.Context, userID int) (*[]TireInventory, error)
	GetTireNotes(ctx context.Context, inventoryID int) (string, error)
	GetCarDetails(ctx context.Context, inventoryID int) (*CarDetail, error)
	UpdateInventoryNote(ctx context.Context, userID int, note string) error
}

type Image struct {
	ID        int            `json:"id"`
	Type      string         `json:"type"`
	Size      uint64         `json:"size"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt utils.NullTime `json:"updatedAt"`
	DeletedAt utils.NullTime `json:"deletedAt"`
}

type TireModel struct {
	ID        int              `json:"id" db:"tire_model_id"`
	Brand     utils.NullString `json:"brand"`
	Type      utils.NullString `json:"type"`
	Size      utils.NullString `json:"size"`
	DOT       utils.NullString `json:"dot" db:"dot"`
	CreatedAt time.Time        `json:"createdAt" db:"tire_model_created_at"`
}

type TireInventory struct {
	ID      int              `json:"id" db:"tire_inventory_id"`
	UserID  int              `json:"userId" db:"user_id"`
	IsSaved bool             `json:"isSaved" db:"is_saved"`
	Note    utils.NullString `json:"note" db:"note"`

	TireModel `json:"tireModel"`
	CarDetail `json:"carDetail"`

	CreatedAt time.Time      `json:"createdAt" db:"tire_inventory_created_at"`
	UpdatedAt utils.NullTime `json:"updatedAt" db:"tire_inventory_updated_at"`
	DeletedAt utils.NullTime `json:"deletedAt" db:"tire_inventory_deleted_at"`
}

type CarDetail struct {
	ID              utils.NullInt32  `json:"id" db:"car_detail_id"`
	TireInventoryID utils.NullInt32  `json:"tireInventoryId"`
	Make            utils.NullString `json:"make" db:"car_make"`
	Model           utils.NullString `json:"model"`
	Year            utils.NullString `json:"year"`
	LicensePlate    utils.NullString `json:"licensePlate" db:"license_plate"`
	Color           utils.NullString `json:"color"`

	CreatedAt utils.NullTime `json:"createdAt" db:"car_detail_created_at"`
	UpdatedAt utils.NullTime `json:"updatedAt"`
	DeletedAt utils.NullTime `json:"deletedAt"`
}
