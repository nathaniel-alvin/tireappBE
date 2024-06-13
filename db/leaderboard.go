package db

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/nathaniel-alvin/tireappBE/types"
)

type LeaderboardRepo struct {
	db *sqlx.DB
}

func NewLeaderboardRepo(db *sqlx.DB) *LeaderboardRepo {
	return &LeaderboardRepo{
		db: db,
	}
}

func (s *LeaderboardRepo) GetTireModelLeaderboard(ctx context.Context) (*[]types.TireModelLeaderboard, error) {
	leaderboard := []types.TireModelLeaderboard{}
	query := `
	SELECT brand, COUNT(*) as count
	FROM tire_model
	WHERE brand IS NOT NULL
	GROUP BY brand
	ORDER BY count DESC
	LIMIT 3;
	`
	err := s.db.Select(&leaderboard, query)
	if err != nil {
		return nil, err
	}
	return &leaderboard, nil
}
