package db

import (
	"bytes"
	"log"
	"net/http"
	"strings"
	"time"

	tireapperror "github.com/nathaniel-alvin/tireappBE/error"
	"github.com/nathaniel-alvin/tireappBE/service/s3"

	"github.com/jmoiron/sqlx"
)

// MaxFilenameBytes is the maximum number of bytes allowed for uploaded files
// There's no technical reason on PicoShare's side for this limitation, but it's
// useful to have some upper bound to limit malicious inputs, and 255 is a
// common filename limit (in single-byte characters) across most filesystems.
const MaxFilenameBytes = 255

var (
	ErrFilenameEmpty             = tireapperror.Errorf(tireapperror.EINVALID, "filename must be non-empty")
	ErrFilenameTooLong           = tireapperror.Errorf(tireapperror.EINVALID, "filename too long")
	ErrFilenameHasDotPrefix      = tireapperror.Errorf(tireapperror.EINVALID, "filename cannot begin with dots")
	ErrFilenameIllegalCharacters = tireapperror.Errorf(tireapperror.EINVALID, "illegal characters in filename")
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
		return 0, tireapperror.Errorf(tireapperror.EINVALID, "%v", err)
	}
	defer func() {
		if err := r.MultipartForm.RemoveAll(); err != nil {
			log.Printf("failed to free multipart form resources: %v", err)
		}
	}()

	file, metadata, err := r.FormFile("file")
	if err != nil {
		return 0, tireapperror.Errorf(tireapperror.EINVALID, "%v", err)
	}

	if metadata.Size == 0 {
		return 0, tireapperror.Errorf(tireapperror.EINVALID, "file is empty")
	}

	filename, err := parse(metadata.Filename)
	if err != nil {
		return 0, err
	}

	contentType, err := parseContentType(metadata.Header.Get("Content-Type"))
	if err != nil {
		return 0, err
	}

	var buf bytes.Buffer
	_, err = buf.ReadFrom(file)
	if err != nil {
		return 0, tireapperror.Errorf(tireapperror.EINTERNAL, "failed to read file: %v", err)
	}

	// upload file to s3
	ImageURL, err := s3.UploadImageToS3(buf.Bytes(), filename)
	if err != nil {
		return 0, err
	}

	// begin transaction
	tx, err := s.db.Beginx()
	if err != nil {
		return 0, tireapperror.Errorf(tireapperror.EINTERNAL, "%v", err)
	}
	defer tx.Rollback()

	var InventoryID int
	var TireID int
	var ImageID int

	// create tire model -> tire_inventory -> image
	err = tx.QueryRowx("INSERT INTO tire_model (created_at) VALUES ($1) RETURNING id;", time.Now()).Scan(&TireID)
	if err != nil {
		// log.Printf("Failed to create tire model: %v", err)
		return 0, tireapperror.Errorf(tireapperror.EINTERNAL, "%v", err)
	}

	err = tx.QueryRowx("INSERT INTO tire_inventory (user_id, tire_id, is_saved, created_at) VALUES ($1, $2, $3, $4) RETURNING id;", userID, TireID, false, time.Now()).Scan(&InventoryID)
	if err != nil {
		// log.Printf("Failed to create scanned tire: %v", err)
		return 0, tireapperror.Errorf(tireapperror.EINTERNAL, "%v", err)
	}

	err = tx.QueryRowx("INSERT INTO image (inventory_id, data_url, type, size, created_at, filename) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;", InventoryID, ImageURL, contentType, metadata.Size, time.Now(), filename).Scan(&ImageID)
	if err != nil {
		// log.Printf("failed to save entry: %v", err)
		return 0, tireapperror.Errorf(tireapperror.EINTERNAL, "%v", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, tireapperror.Errorf(tireapperror.EINTERNAL, "%v", err)
	}

	return InventoryID, nil
}

func (s *UploadRepo) CreateImageForInventory(r *http.Request, inventoryID int) error {
	multipartMaxMemory := mibToBytes(1)
	if err := r.ParseMultipartForm(multipartMaxMemory); err != nil {
		return tireapperror.Errorf(tireapperror.EINVALID, "%v", err)
	}
	defer func() {
		if err := r.MultipartForm.RemoveAll(); err != nil {
			log.Printf("failed to free multipart form resources: %v", err)
		}
	}()

	file, metadata, err := r.FormFile("file")
	if err != nil {
		return tireapperror.Errorf(tireapperror.EINVALID, "%v", err)
	}

	if metadata.Size == 0 {
		return tireapperror.Errorf(tireapperror.EINVALID, "file is empty")
	}

	filename, err := parse(metadata.Filename)
	if err != nil {
		return err
	}

	contentType, err := parseContentType(metadata.Header.Get("Content-Type"))
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	_, err = buf.ReadFrom(file)
	if err != nil {
		return tireapperror.Errorf(tireapperror.EINTERNAL, "failed to read file: %v", err)
	}

	// upload file to s3
	ImageURL, err := s3.UploadImageToS3(buf.Bytes(), filename)
	if err != nil {
		return err
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return tireapperror.Errorf(tireapperror.EINTERNAL, "%v", err)
	}
	defer tx.Rollback()

	query := "INSERT INTO image (inventory_id, data_url, type, size, created_at, filename) VALUES ($1, $2, $3, $4, $5, $6);"
	_, err = tx.Exec(query, inventoryID, ImageURL, contentType, metadata.Size, time.Now(), filename)
	if err != nil {
		return tireapperror.Errorf(tireapperror.EINTERNAL, "%v", err)
	}

	if err := tx.Commit(); err != nil {
		return tireapperror.Errorf(tireapperror.EINTERNAL, "%v", err)
	}

	return nil
}

func (s *UploadRepo) UpdateImageForInventory(r *http.Request, inventoryID int) error {
	multipartMaxMemory := mibToBytes(1)
	if err := r.ParseMultipartForm(multipartMaxMemory); err != nil {
		return tireapperror.Errorf(tireapperror.EINVALID, "%v", err)
	}
	defer func() {
		if err := r.MultipartForm.RemoveAll(); err != nil {
			log.Printf("failed to free multipart form resources: %v", err)
		}
	}()

	file, metadata, err := r.FormFile("file")
	if err != nil {
		return tireapperror.Errorf(tireapperror.EINVALID, "%v", err)
	}

	if metadata.Size == 0 {
		return tireapperror.Errorf(tireapperror.EINVALID, "file is empty")
	}

	filename, err := parse(metadata.Filename)
	if err != nil {
		return err
	}

	contentType, err := parseContentType(metadata.Header.Get("Content-Type"))
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	_, err = buf.ReadFrom(file)
	if err != nil {
		return tireapperror.Errorf(tireapperror.EINTERNAL, "failed to read file: %v", err)
	}

	// upload file to s3
	ImageURL, err := s3.UploadImageToS3(buf.Bytes(), filename)
	if err != nil {
		return err
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return tireapperror.Errorf(tireapperror.EINTERNAL, "%v", err)
	}
	defer tx.Rollback()

	query := `
	UPDATE 
		image 
	SET 
		data_url = $2,
		type = $3,
		size = $4,
		updated_at = $5,
		filename = $6
	WHERE 
		inventory_id = $1`
	_, err = tx.Exec(query, inventoryID, ImageURL, contentType, metadata.Size, time.Now(), filename)
	if err != nil {
		return tireapperror.Errorf(tireapperror.EINTERNAL, "%v", err)
	}

	if err := tx.Commit(); err != nil {
		return tireapperror.Errorf(tireapperror.EINTERNAL, "%v", err)
	}

	return nil
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
