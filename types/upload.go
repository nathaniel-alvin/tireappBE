package types

import "net/http"

type UploadStore interface {
	InsertFileFromRequest(r *http.Request, userID int) (int, error)
	CreateImageForInventory(r *http.Request, inventoryID int) error
	UpdateImageForInventory(r *http.Request, inventoryID int) error
}
