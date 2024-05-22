package types

import "time"

type Workshop struct {
	ID              int       `json:"id"`
	TireInventoryID int       `json:"tireInventoryId"`
	Name            string    `json:"name"`
	Address         string    `json:"address"`
	ContactNumber   string    `json:"contactNumber"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
	DeletedAt       time.Time `json:"deletedAt"`
}
