package db

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	tireappbe "github.com/nathaniel-alvin/tireappBE"
	"github.com/nathaniel-alvin/tireappBE/types"

	tireapperror "github.com/nathaniel-alvin/tireappBE/error"
)

type InventoryRepo struct {
	db *sqlx.DB
}

func NewInventoryRepo(db *sqlx.DB) *InventoryRepo {
	return &InventoryRepo{
		db: db,
	}
}

// TODO: check each row if deleted date == nil
func (s *InventoryRepo) GetInventories(ctx context.Context, userID int) (*[]types.TireInventory, error) {
	tireInventories := []types.TireInventory{}
	query := `
    SELECT 
        ti.id AS tire_inventory_id,
		ti.user_id,
		ti.is_saved,
		ti.note,
        ti.created_at AS tire_inventory_created_at,
        ti.updated_at AS tire_inventory_updated_at,
        ti.deleted_at AS tire_inventory_deleted_at,
        tm.id AS tire_model_id,
        tm.brand,
        tm.type,
        tm.size,
        tm.dot,
        tm.created_at AS tire_model_created_at,
        v.id AS car_detail_id,
		v.inventory_id as car_detail_inventory_id,
        v.make AS car_make,
        v.model,
        v.year,
        v.license_plate,
        v.created_at AS car_detail_created_at,
		i.id AS image_id,
		i.data_url AS image_url,
		i.type AS image_type,
		i.size AS image_size,
		i.created_at AS image_created_at
    FROM 
        tire_inventory ti
    JOIN 
        tire_model tm ON ti.tire_id = tm.id
    LEFT JOIN 
        vehicle v ON v.inventory_id = ti.id
    LEFT JOIN 
		image i ON i.inventory_id = ti.id
    WHERE 
        ti.user_id = $1 AND ti.is_saved = true
    ORDER BY 
        ti.created_at DESC;
    `
	// iseng ngececk aja
	if userID != tireappbe.UserIDFromContext(ctx) {
		return nil, tireapperror.Errorf(tireapperror.EINTERNAL, "incorrect user ID")
	}

	err := s.db.Select(&tireInventories, query, userID)
	if err != nil {
		return nil, tireapperror.Errorf(tireapperror.EINTERNAL, "%v", err)
	}

	return &tireInventories, nil
}

func (s *InventoryRepo) GetInventoryByID(ctx context.Context, userID int, inventoryID int) (*types.TireInventory, error) {
	tireInventories := []types.TireInventory{}
	query := `
    SELECT 
        ti.id AS tire_inventory_id,
		ti.user_id,
		ti.is_saved,
		ti.note,
        ti.created_at AS tire_inventory_created_at,
        ti.updated_at AS tire_inventory_updated_at,
        ti.deleted_at AS tire_inventory_deleted_at,
        tm.id AS tire_model_id,
        tm.brand,
        tm.type,
        tm.size,
        tm.dot,
        tm.created_at AS tire_model_created_at,
        v.id AS car_detail_id,
		v.inventory_id as car_detail_inventory_id,
        v.make AS car_make,
        v.model,
        v.year,
        v.license_plate,
        v.created_at AS car_detail_created_at,
		i.id AS image_id,
		i.data_url AS image_url,
		i.type AS image_type,
		i.size AS image_size,
		i.created_at AS image_created_at
    FROM 
        tire_inventory ti
    JOIN 
        tire_model tm ON ti.tire_id = tm.id
    LEFT JOIN 
        vehicle v ON v.inventory_id = ti.id
    LEFT JOIN 
		image i ON i.inventory_id = ti.id
    WHERE 
        ti.user_id = $1 AND ti.is_saved = true AND ti.id = $2
    ORDER BY 
        ti.created_at DESC;
    `
	// iseng ngececk aja
	if userID != tireappbe.UserIDFromContext(ctx) {
		return nil, tireapperror.Errorf(tireapperror.EINTERNAL, "incorrect user ID")
	}

	err := s.db.Select(&tireInventories, query, userID, inventoryID)
	if err != nil {
		return nil, tireapperror.Errorf(tireapperror.EINTERNAL, "%v", err)
	}

	if len(tireInventories) == 0 {
		return nil, tireapperror.Errorf(tireapperror.ENOTFOUND, "")
	}

	if len(tireInventories) > 1 {
		return nil, tireapperror.Errorf(tireapperror.EINTERNAL, "get inventory by ID: two results returned")
	}

	return &tireInventories[0], nil
}

func (s *InventoryRepo) CreateInventory(ctx context.Context, userID int, i *types.TireInventory, m *types.TireModel) (int, error) {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, tireapperror.Errorf(tireapperror.EINTERNAL, "%v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	modelID, err := createTireModel(ctx, tx, m)
	if err != nil {
		return 0, err
	}

	i.TireModel.ID = modelID
	i.CreatedAt = time.Now()
	inventoryID, err := createTireInventory(ctx, tx, userID, i)
	if err != nil {
		return 0, err
	}

	// create empty car details so we can just update it
	_, err = createCarDetail(ctx, tx, inventoryID, nil)
	if err != nil {
		return 0, err
	}

	// create empty image to update later
	_, err = createEmptyTireImage(ctx, tx, inventoryID)

	return inventoryID, nil
}

func (s *InventoryRepo) UpdateTireModel(ctx context.Context, inventoryID int, tm types.TireModel) error {
	query := `
        UPDATE tire_model
        SET 
            brand = $2,
            type = $3,
            size = $4,
            dot = $5,
			updated_at = $6
        FROM tire_inventory
        WHERE tire_model.id = tire_inventory.tire_id
          AND tire_inventory.id = $1
    `
	_, err := s.db.Exec(query, inventoryID, tm.Brand.String, tm.Type.String, tm.Size.String, tm.DOT.String, time.Now())
	if err != nil {
		return tireapperror.Errorf(tireapperror.EINTERNAL, "%v", err)
	}
	return nil
}

func (s *InventoryRepo) UpdateCarDetail(ctx context.Context, inventoryID int, cd types.CarDetail) error {
	query := `
		UPDATE vehicle
		SET 
			make = $2,
			model = $3,
			year = $4,
			license_plate = $5,
			color = $6,
			updated_at = $7
		FROM tire_inventory ti
		WHERE vehicle.inventory_id = ti.id AND ti.id = $1
	`
	_, err := s.db.Exec(query, inventoryID, cd.Make.String, cd.Model.String, cd.Year.String, cd.LicensePlate.String, cd.Color.String, time.Now())
	if err != nil {
		return tireapperror.Errorf(tireapperror.EINTERNAL, "%v", err)
	}
	return nil
}

func (s *InventoryRepo) DeleteInventory(ctx context.Context, inventoryID int) error {
	query := `
		UPDATE tire_inventory
		SET
			is_saved = false,
			deleted_at = NOW()
		WHERE id = $1;
	`
	_, err := s.db.Exec(query, inventoryID)
	if err != nil {
		return tireapperror.Errorf(tireapperror.EINTERNAL, "%v", err)
	}
	return nil
}

func (s *InventoryRepo) GetInventoryHistory(ctx context.Context, userID int) (*[]types.TireInventory, error) {
	tireInventories := []types.TireInventory{}
	query := `
    SELECT 
        ti.id AS tire_inventory_id,
		ti.user_id,
		ti.is_saved,
		ti.note,
        ti.created_at AS tire_inventory_created_at,
        ti.updated_at AS tire_inventory_updated_at,
        ti.deleted_at AS tire_inventory_deleted_at,
        tm.id AS tire_model_id,
        tm.brand,
        tm.type,
        tm.size,
        tm.dot,
        tm.created_at AS tire_model_created_at,
        v.id AS car_detail_id,
		v.inventory_id as car_detail_inventory_id,
        v.make AS car_make,
        v.model,
        v.year,
        v.license_plate,
        v.created_at AS car_detail_created_at,
		i.id AS image_id,
		i.data_url AS image_url,
		i.type AS image_type,
		i.size AS image_size,
		i.created_at AS image_created_at
    FROM 
        tire_inventory ti
    JOIN 
        tire_model tm ON ti.tire_id = tm.id
    LEFT JOIN 
        vehicle v ON v.inventory_id = ti.id
    LEFT JOIN 
		image i ON i.inventory_id = ti.id
    WHERE 
        ti.user_id = $1 
    ORDER BY 
        ti.created_at DESC;
    `
	// iseng ngececk aja
	if userID != tireappbe.UserIDFromContext(ctx) {
		return nil, tireapperror.Errorf(tireapperror.EINTERNAL, "incorrect user ID")
	}

	err := s.db.Select(&tireInventories, query, userID)
	if err != nil {
		return nil, tireapperror.Errorf(tireapperror.EINTERNAL, "%v", err)
	}

	return &tireInventories, nil
}

func createTireModel(ctx context.Context, tx *sqlx.Tx, tm *types.TireModel) (int, error) {
	var TireID int
	query := "INSERT INTO tire_model (brand, type, size, dot, created_at) VALUES (:brand, :type, :size, :dot, :created_at) RETURNING id;"
	params := map[string]interface{}{
		"brand":      nil,
		"type":       nil,
		"size":       nil,
		"dot":        nil,
		"created_at": time.Now(),
	}

	if tm != nil {
		if tm.Brand.Valid {
			params["brand"] = tm.Brand.String
		}
		if tm.Type.Valid {
			params["type"] = tm.Type.String
		}
		if tm.Size.Valid {
			params["size"] = tm.Size.String
		}
		if tm.DOT.Valid {
			params["dot"] = tm.DOT.String
		}
	}

	rows, err := tx.NamedQuery(query, params)
	if err != nil {
		return 0, tireapperror.Errorf(tireapperror.EINTERNAL, "%v", err)
	}

	for rows.Next() {
		err := rows.Scan(&TireID)
		if err != nil {
			return 0, tireapperror.Errorf(tireapperror.EINTERNAL, "%v", err)
		}
	}
	return TireID, nil
}

func createTireInventory(ctx context.Context, tx *sqlx.Tx, userID int, ti *types.TireInventory) (int, error) {
	// authenticate
	// if ok := userID == tireappbe.UserIDFromContext(ctx); !ok {
	// 	return 0, fmt.Errorf("authentication failed when trying to create inventory")
	// }

	// TEMP: use userID from context
	userID = tireappbe.UserIDFromContext(ctx)

	var InventoryID int
	query := "INSERT INTO tire_inventory (user_id, tire_id, is_saved, note, created_at) VALUES (:user_id, :tire_model_id, :is_saved, :note, :created_at) RETURNING id;"

	params := map[string]interface{}{
		"user_id":       userID,
		"tire_model_id": ti.TireModel.ID,
		"is_saved":      ti.IsSaved,
		"note":          ti.Note.String,
		"created_at":    time.Now(),
	}

	rows, err := tx.NamedQuery(query, params)
	if err != nil {
		return 0, tireapperror.Errorf(tireapperror.EINTERNAL, "%v", err)
	}

	for rows.Next() {
		err := rows.Scan(&InventoryID)
		if err != nil {
			return 0, tireapperror.Errorf(tireapperror.EINTERNAL, "%v", err)
		}
	}

	return InventoryID, nil
}

func createCarDetail(ctx context.Context, tx *sqlx.Tx, inventoryID int, cd *types.CarDetail) (int, error) {
	var CarDetailID int
	query := "INSERT INTO vehicle (inventory_id, license_plate, color, make, model, year, created_at) VALUES (:inventoryID, :license_plate, :color, :make, :model, :year, :created_at) RETURNING id;"
	params := map[string]interface{}{
		"inventoryID":   inventoryID,
		"license_plate": nil,
		"color":         nil,
		"make":          nil,
		"model":         nil,
		"year":          nil,
		"created_at":    time.Now(),
	}
	if cd != nil {
		if cd.LicensePlate.Valid {
			params["license_plate"] = cd.LicensePlate.String
		}
		if cd.Color.Valid {
			params["color"] = cd.Color.String
		}
		if cd.Make.Valid {
			params["make"] = cd.Make.String
		}
		if cd.Model.Valid {
			params["model"] = cd.Model.String
		}
		if cd.Year.Valid {
			params["year"] = cd.Year.String
		}
	}
	rows, err := tx.NamedQuery(query, params)
	if err != nil {
		return 0, tireapperror.Errorf(tireapperror.EINTERNAL, "%v", err)
	}

	for rows.Next() {
		err := rows.Scan(&CarDetailID)
		if err != nil {
			return 0, tireapperror.Errorf(tireapperror.EINTERNAL, "%v", err)
		}
	}

	return CarDetailID, nil
}

func createEmptyTireImage(ctx context.Context, tx *sqlx.Tx, inventoryID int) (int, error) {
	var ImageID int
	query := "INSERT INTO image (inventory_id) VALUES ($1) RETURNING id;"

	rows, err := tx.Queryx(query, inventoryID)
	if err != nil {
		return 0, tireapperror.Errorf(tireapperror.EINTERNAL, "%v", err)
	}

	for rows.Next() {
		err := rows.Scan(&ImageID)
		if err != nil {
			return 0, tireapperror.Errorf(tireapperror.EINTERNAL, "%v", err)
		}
	}
	return ImageID, nil
}

func (s *InventoryRepo) GetTireNotes(ctx context.Context, inventoryID int) (string, error) {
	var tireNote []string

	query := `
	SELECT ti.note
	FROM tire_inventory ti
	WHERE ti.id = $1;
	`

	err := s.db.Select(&tireNote, query, inventoryID)
	if err != nil {
		return "", tireapperror.Errorf(tireapperror.EINTERNAL, "%v", err)
	}

	if len(tireNote) != 1 {
		return "", tireapperror.Errorf(tireapperror.EINTERNAL, "more than one inventory: %v", err)
	}

	return tireNote[0], nil
}

func (s *InventoryRepo) GetCarDetails(ctx context.Context, inventoryID int) (*types.CarDetail, error) {
	carDetail := []types.CarDetail{}

	query := `
	SELECT 
		vehicle.id as car_detail_id,
		vehicle.inventory_id as car_detail_inventory_id,
		vehicle.make as car_make,
		vehicle.model,
		vehicle.year,
		vehicle.license_plate,
		vehicle.color,
		vehicle.created_at as car_detail_created_at,
		vehicle.updated_at as car_detail_updated_at,
		vehicle.deleted_at as car_detail_deleted_at
	FROM 
		vehicle
	WHERE
		vehicle.inventory_id = $1;
	`

	err := s.db.Select(&carDetail, query, inventoryID)
	if err != nil {
		return nil, tireapperror.Errorf(tireapperror.EINTERNAL, "%v", err)
	}

	if len(carDetail) != 1 {
		return nil, tireapperror.Errorf(tireapperror.EINTERNAL, "more than one vehicle: %v", err)
	}

	return &carDetail[0], nil
}

func (s *InventoryRepo) UpdateInventoryNote(ctx context.Context, inventoryID int, note string) error {
	query := `
	UPDATE tire_inventory
	SET 
		note = $2,
		updated_at = $3
	WHERE id = $1;
	`
	_, err := s.db.Exec(query, inventoryID, note, time.Now())
	if err != nil {
		return tireapperror.Errorf(tireapperror.EINTERNAL, "%v", err)
	}

	return nil
}
