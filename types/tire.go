package types

import "time"

type Image struct {
	ID        ImageID   `json:"id"`
	Type      string    `json:"type"`
	Size      uint64    `json:"size"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt time.Time `json:"deletedAt"`
}

type TireModel struct {
	ID        int       `json:"id" db:"tire_model_id"`
	Brand     string    `json:"brand" db:"tire_brand"`
	Type      string    `json:"type"`
	Size      string    `json:"size"`
	DOT       string    `json:"dot" db:"dot"`
	CreatedAt time.Time `json:"createdAt" db:"tire_model_created_at"`
}

type TireInventory struct {
	ID      int  `json:"id" db:"tire_inventory_id"`
	UserID  int  `json:"userId" db:"user_id"`
	IsSaved bool `json:"isSaved" db:"is_saved"`

	TireModel TireModel `json:"tireModel"`
	CarDetail CarDetail `json:"carDetail"`

	CreatedAt time.Time `json:"createdAt" db:"tire_inventory_created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"tire_inventory_updated_at"`
	DeletedAt time.Time `json:"deletedAt" db:"tire_inventory_deleted_at"`
}

type CarDetail struct {
	ID              int    `json:"id" db:"car_detail_id"`
	TireInventoryID int    `json:"tireInventoryId"`
	Brand           string `json:"brand" db:"car_brand"`
	Model           string `json:"model"`
	Year            string `json:"year"`
	LicensePlate    string `json:"licensePlate"`
	Color           string `json:"color"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt time.Time `json:"deletedAt"`
}
