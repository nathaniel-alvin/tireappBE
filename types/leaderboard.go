package types

import "context"

type LeaderboardRepo interface {
	GetTireModelLeaderboard(ctx context.Context) (*[]TireModelLeaderboard, error)
}

type TireModelLeaderboard struct {
	Brand string `json:"brand" db:"brand"`
	Count int    `json:"count" db:"count"`
}
