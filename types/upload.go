package types

import "net/http"

type UploadStore interface {
	InsertFileFromRequest(r *http.Request, userID int) (int, error)
}
