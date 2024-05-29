package db

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	tireappbe "github.com/nathaniel-alvin/tireappBE"
	"github.com/nathaniel-alvin/tireappBE/types"
)

type InventoryRepo struct {
	db *sqlx.DB
}

func NewInventoryRepo(db *sqlx.DB) *InventoryRepo {
	return &InventoryRepo{
		db: db,
	}
}

func (s *InventoryRepo) GetInventories(ctx context.Context, userID int) (*[]types.TireInventory, error) {
	tireInventories := []types.TireInventory{}
	query := `
    SELECT 
        ti.id AS tire_inventory_id,
		ti.user_id,
		ti.is_saved,
        ti.created_at AS tire_inventory_created_at,
        ti.updated_at AS tire_inventory_updated_at,
        ti.deleted_at AS tire_inventory_deleted_at,
        tm.id AS tire_model_id,
        tm.brand AS tire_brand,
        tm.type,
        tm.size,
        tm.dot,
        tm.created_at AS tire_model_created_at,
        v.id AS car_detail_id,
        v.brand AS car_brand,
        v.model,
        v.year,
        v.created_at AS car_detail_created_at
    FROM 
        tire_inventory ti
    JOIN 
        tire_model tm ON ti.tire_id = tm.id
    LEFT JOIN 
        vehicle v ON v.scan_id = ti.id
    LEFT JOIN 
        workshop w ON w.scan_id = ti.id
    WHERE 
        ti.user_id = $1 AND ti.is_saved = true
    ORDER BY 
        ti.created_at DESC;
    `
	// iseng ngececk aja
	if userID != tireappbe.UserIDFromContext(ctx) {
		return nil, fmt.Errorf("different user id when checking with context")
	}

	err := s.db.Select(&tireInventories, query, userID)
	if err != nil {
		return nil, err
	}

	return &tireInventories, nil
}

func (s *InventoryRepo) GetInventoriesByID(ctx context.Context, userID int, inventoryID int) (*types.TireInventory, error) {
	tireInventories := types.TireInventory{}
	query := `
    SELECT 
        ti.id AS tire_inventory_id,
		ti.user_id,
		ti.is_saved,
        ti.created_at AS tire_inventory_created_at,
        ti.updated_at AS tire_inventory_updated_at,
        ti.deleted_at AS tire_inventory_deleted_at,
        tm.id AS tire_model_id,
        tm.brand AS tire_brand,
        tm.type,
        tm.size,
        tm.dot,
        tm.created_at AS tire_model_created_at,
        v.id AS car_detail_id,
        v.brand AS car_brand,
        v.model,
        v.year,
        v.created_at AS car_detail_created_at
    FROM 
        tire_inventory ti
    JOIN 
        tire_model tm ON ti.tire_id = tm.id
    LEFT JOIN 
        vehicle v ON v.scan_id = ti.id
    LEFT JOIN 
        workshop w ON w.scan_id = ti.id
    WHERE 
        ti.user_id = $1 AND ti.is_saved = true AND ti.id = $2
    ORDER BY 
        ti.created_at DESC;
    `
	// iseng ngececk aja
	if userID != tireappbe.UserIDFromContext(ctx) {
		return nil, fmt.Errorf("different user id when checking with context")
	}

	err := s.db.Select(&tireInventories, query, userID, inventoryID)
	if err != nil {
		return nil, err
	}

	return &tireInventories, nil
}
