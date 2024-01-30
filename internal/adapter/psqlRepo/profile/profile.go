package profile

import (
	"context"
	"database/sql"
	"errors"
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

func (r *RepositoryProfile) Add(ctx context.Context, p *profile.Profile) (*profile.Profile, error) {
	query := "INSERT INTO profiles (display_name, age, gender, location, description, is_deleted," +
		" is_blocked, is_premium, created_at, updated_at, last_online) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9," +
		" $10, $11) RETURNING id"
	err := r.db.QueryRowContext(ctx, query, p.DisplayName, p.Age, p.Gender, p.Location, p.Description, p.IsDeleted,
		p.IsBlocked, p.IsPremium, p.CreatedAt, p.UpdatedAt, p.LastOnline).Scan(&p.ID)
	if err != nil {
		r.logger.Debug(
			"error func Add, method QueryRowContext by path internal/adapter/psqlRepo/profile/profile.go",
			zap.Error(err))
		return nil, err
	}
	return p, nil
}

func (r *RepositoryProfile) FindById(ctx context.Context, id uint64) (*profile.Profile, error) {
	p := profile.Profile{}
	query := `SELECT id, display_name, age, gender, location, description, is_deleted, is_blocked, is_premium,
       created_at, updated_at, last_online
			  FROM profiles
			  WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)
	if row == nil {
		err := errors.New("no rows found")
		r.logger.Debug(
			"error func FindById, method QueryRowContext by path internal/adapter/psqlRepo/profile/profile.go",
			zap.Error(err))
		return nil, err
	}
	err := row.Scan(&p.ID, &p.DisplayName, &p.Age, &p.Gender, &p.Location, &p.Description, &p.IsDeleted, &p.IsBlocked,
		&p.IsPremium, &p.CreatedAt, &p.UpdatedAt, &p.LastOnline)
	if err != nil {
		r.logger.Debug("error func FindById, method Scan by path internal/adapter/psqlRepo/profile/profile.go",
			zap.Error(err))
		return nil, err
	}
	return &p, nil
}

func (r *RepositoryProfile) AddTelegram(
	ctx context.Context, p *profile.TelegramProfile) (*profile.TelegramProfile, error) {
	query := "INSERT INTO profile_telegram (profile_id, telegram_id, username, first_name, last_name, language_code," +
		" allows_write_to_pm, query_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id"
	err := r.db.QueryRowContext(ctx, query, p.ProfileID, p.TelegramID, p.UserName, p.Firstname, p.Lastname,
		p.LanguageCode, p.AllowsWriteToPm, p.QueryID).Scan(&p.ID)
	if err != nil {
		r.logger.Debug(
			"error func AddTelegram, method QueryRowContext by path internal/adapter/psqlRepo/profile/profile.go",
			zap.Error(err))
		return nil, err
	}
	return p, nil
}

func (r *RepositoryProfile) FindTelegramById(ctx context.Context, profileID uint64) (*profile.TelegramProfile, error) {
	p := profile.TelegramProfile{}
	query := `SELECT id, profile_id, telegram_id, username, first_name, last_name, language_code, allows_write_to_pm,
       query_id
			  FROM profile_telegram
			  WHERE profile_id = $1`
	row := r.db.QueryRowContext(ctx, query, profileID)
	if row == nil {
		err := errors.New("no rows found")
		r.logger.Debug(
			"error func FindById, method QueryRowContext by path internal/adapter/psqlRepo/profile/profile.go",
			zap.Error(err))
		return nil, err
	}
	err := row.Scan(&p.ID, &p.ProfileID, &p.TelegramID, &p.Firstname, &p.Lastname, &p.UserName, &p.LanguageCode,
		&p.AllowsWriteToPm, &p.QueryID)
	if err != nil {
		r.logger.Debug("error func FindById, method Scan by path internal/adapter/psqlRepo/profile/profile.go",
			zap.Error(err))
		return nil, err
	}
	return &p, nil
}

func (r *RepositoryProfile) AddImage(ctx context.Context, p *profile.ImageProfile) (*profile.ImageProfile, error) {
	query := "INSERT INTO profile_images (profile_id, name, url, size, created_at, updated_at, is_deleted," +
		" is_blocked, is_primary, is_private) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id"
	err := r.db.QueryRowContext(ctx, query, p.ProfileID, p.Name, p.Url, p.Size, p.CreatedAt, p.UpdatedAt, p.IsDeleted,
		p.IsBlocked, p.IsPrimary, p.IsPrivate).Scan(&p.ID)
	if err != nil {
		r.logger.Debug(
			"error func AddImage, method QueryRowContext by path internal/adapter/psqlRepo/profile/profile.go",
			zap.Error(err))
		return nil, err
	}
	return p, nil
}

func (r *RepositoryProfile) SelectListPublicImage(
	ctx context.Context, profileID uint64) ([]*profile.ImageProfile, error) {
	query := `SELECT id, profile_id, name, url, size, created_at, updated_at, is_deleted, is_blocked, is_primary,
       is_private
	FROM profile_images
	WHERE profile_id = $1`
	rows, err := r.db.QueryContext(ctx, query, profileID)
	if err != nil {
		r.logger.Debug("error func SelectListPublicImage,"+
			" method QueryContext by path internal/adapter/psqlRepo/profile/profile.go", zap.Error(err))
		return nil, err
	}
	defer rows.Close()
	list := make([]*profile.ImageProfile, 0)
	for rows.Next() {
		p := profile.ImageProfile{}
		err := rows.Scan(&p.ID, &p.ProfileID, &p.Name, &p.Url, &p.Size, &p.CreatedAt, &p.UpdatedAt, &p.IsDeleted,
			&p.IsBlocked, &p.IsPrimary, &p.IsPrivate)
		if err != nil {
			r.logger.Debug("error func SelectListPublicImage,"+
				" method Scan by path internal/adapter/psqlRepo/profile/profile.go", zap.Error(err))
			continue
		}
		list = append(list, &p)
	}
	return list, nil
}
