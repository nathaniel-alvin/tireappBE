package db

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

// MaxFilenameBytes is the maximum number of bytes allowed for uploaded files
// There's no technical reason on PicoShare's side for this limitation, but it's
// useful to have some upper bound to limit malicious inputs, and 255 is a
// common filename limit (in single-byte characters) across most filesystems.
const MaxFilenameBytes = 255

var (
	ErrFilenameEmpty             = errors.New("filename must be non-empty")
	ErrFilenameTooLong           = errors.New("filename too long")
	ErrFilenameHasDotPrefix      = errors.New("filename cannot begin with dots")
	ErrFilenameIllegalCharacters = errors.New("illegal characters in filename")
)

type UploadRepo struct {
	db *sqlx.DB
}

func NewUploadRepo(db *sqlx.DB) *UploadRepo {
	return &UploadRepo{
		db: db,
	}
}

func (s UploadRepo) InsertFileFromRequest(r *http.Request, userID int) (int, error) {
	multipartMaxMemory := mibToBytes(1)
	if err := r.ParseMultipartForm(multipartMaxMemory); err != nil {
		return 0, err
	}
	defer func() {
		if err := r.MultipartForm.RemoveAll(); err != nil {
			log.Printf("failed to free multipart form resources: %v", err)
		}
	}()

	_, metadata, err := r.FormFile("file")
	if err != nil {
		return 0, err
	}

	if metadata.Size == 0 {
		return 0, fmt.Errorf("file is empty")
	}

	filename, err := parse(metadata.Filename)
	if err != nil {
		return 0, err
	}

	contentType, err := parseContentType(metadata.Header.Get("Content-Type"))
	if err != nil {
		return 0, err
	}

	// begin transaction
	tx, err := s.db.Beginx()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var InventoryID int
	var TireID int
	var ImageID int

	// create tire model -> tire_inventory -> image
	err = tx.QueryRowx("INSERT INTO tire_model (created_at) VALUES ($1) RETURNING id;", time.Now()).Scan(&TireID)
	if err != nil {
		log.Printf("Failed to create tire model: %v", err)
		return 0, err
	}

	err = tx.QueryRowx("INSERT INTO tire_inventory (user_id, tire_id, is_saved, created_at) VALUES ($1, $2, $3, $4) RETURNING id;", userID, TireID, false, time.Now()).Scan(&InventoryID)
	if err != nil {
		log.Printf("Failed to create scanned tire: %v", err)
		return 0, err
	}

	// TODO: mikirin cara upload gambar ke s3
	var ImageURL string

	err = tx.QueryRowx("INSERT INTO image (scan_id, data_url, type, size, created_at, filename) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;", InventoryID, ImageURL, contentType, metadata.Size, time.Now(), filename).Scan(&ImageID)
	if err != nil {
		log.Printf("failed to save entry: %v", err)
		return 0, err
	}
	return InventoryID, nil
}

// mibToBytes converts an amount in MiB to an amount in bytes.
func mibToBytes(i int64) int64 {
	return i << 20
}

func parse(s string) (string, error) {
	if s == "" {
		return "", ErrFilenameEmpty
	}
	if len(s) > MaxFilenameBytes {
		return "", ErrFilenameTooLong
	}
	if s == "." || strings.HasPrefix(s, "..") {
		return "", ErrFilenameHasDotPrefix
	}
	if strings.ContainsAny(s, "\\/\a\b\t\n\v\f\r\n") {
		return "", ErrFilenameIllegalCharacters
	}
	return s, nil
}

func parseContentType(s string) (string, error) {
	// The content type header is fairly open-ended, so we're liberal in what
	// values we accept.
	return s, nil
}
