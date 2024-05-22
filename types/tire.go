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
	ID        int       `json:"id"`
	Brand     string    `json:"brand"`
	Type      string    `json:"type"`
	Size      string    `json:"size"`
	DOT       string    `json:"dot"`
	CreatedAt time.Time `json:"createdAt"`
}

type TireInventory struct {
	ID      int  `json:"id"`
	UserID  int  `json:"userId"`
	TireID  int  `json:"tireId"`
	IsSaved bool `json:"isSaved"`

	CarDetail []CarDetail `json:"carDetail"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt time.Time `json:"deletedAt"`
}

type CarDetail struct {
	ID              int    `json:"id"`
	TireInventoryID int    `json:"tireInventoryId"`
	Brand           string `json:"brand"`
	Model           string `json:"model"`
	Year            string `json:"year"`
	LicensePlate    string `json:"licensePlate"`
	Color           string `json:"color"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt time.Time `json:"deletedAt"`
}
