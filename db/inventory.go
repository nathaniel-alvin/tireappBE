package db

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	tireappbe "github.com/nathaniel-alvin/tireappBE"
	"github.com/nathaniel-alvin/tireappBE/types"
)

type InventoryRepo interface {
	GetInventories(ctx context.Context) ([]*types.TireInventory, error)
}

type InventoryDatabase struct {
	db *sqlx.DB
}

func NewInventoryRepo(db *sqlx.DB) *InventoryDatabase {
	return &InventoryDatabase{
		db: db,
	}
}

func (s *InventoryDatabase) GetInventories(ctx context.Context, id int) (*[]types.TireInventory, error) {
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
        v.brand,
        v.model,
        v.year,
        v.created_at AS car_detail_created_at,
        w.id AS workshop_id,
        w.name AS workshop_name,
        w.address AS workshop_address,
        w.created_at AS workshop_created_at
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
	if id != tireappbe.UserIDFromContext(ctx) {
		return nil, fmt.Errorf("different user id when checking with context")
	}

	err := s.db.Select(&tireInventories, query, id)
	if err != nil {
		return nil, err
	}

	return &tireInventories, nil
}
