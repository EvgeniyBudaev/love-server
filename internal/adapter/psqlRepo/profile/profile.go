package profile

import (
	"context"
	"database/sql"
	"github.com/EvgeniyBudaev/love-server/internal/entity/profile"
	"github.com/EvgeniyBudaev/love-server/internal/logger"
	"go.uber.org/zap"
)

type RepositoryProfile struct {
	logger logger.Logger
	db     *sql.DB
}

func NewRepositoryProfile(logger logger.Logger, db *sql.DB) *RepositoryProfile {
	return &RepositoryProfile{
		logger: logger,
		db:     db,
	}
}

func (r *RepositoryProfile) Create(ctx context.Context, p *profile.Profile) (*profile.Profile, error) {
	query := "INSERT INTO profiles (display_name) VALUES ($1) RETURNING id"
	err := r.db.QueryRowContext(ctx, query, p.DisplayName).Scan(&p.ID)
	if err != nil {
		r.logger.Debug(
			"error func Create, method QueryRowContext by path internal/adapter/psqlRepo/profile/profile.go",
			zap.Error(err))
		return nil, err
	}
	return p, nil
}
