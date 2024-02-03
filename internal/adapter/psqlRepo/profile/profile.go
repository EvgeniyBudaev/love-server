package profile

import (
	"context"
	"database/sql"
	"errors"
	"github.com/EvgeniyBudaev/love-server/internal/entity/pagination"
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
	query := "INSERT INTO profiles (display_name, birthday, gender, search_gender, location, description, height," +
		" weight, looking_for, is_deleted, is_blocked, is_premium, is_show_distance, is_invisible, created_at," +
		" updated_at, last_online) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15," +
		" $16, $17) RETURNING id"
	err := r.db.QueryRowContext(ctx, query, p.DisplayName, p.Birthday, p.Gender, p.SearchGender, p.Location,
		p.Description, p.Height, p.Weight, p.LookingFor, p.IsDeleted, p.IsBlocked, p.IsPremium, p.IsShowDistance,
		p.IsInvisible, p.CreatedAt, p.UpdatedAt, p.LastOnline).Scan(&p.ID)
	if err != nil {
		r.logger.Debug(
			"error func Add, method QueryRowContext by path internal/adapter/psqlRepo/profile/profile.go",
			zap.Error(err))
		return nil, err
	}
	return p, nil
}

func (r *RepositoryProfile) Update(ctx context.Context, p *profile.Profile) (*profile.Profile, error) {
	tx, err := r.db.Begin()
	if err != nil {
		r.logger.Debug(
			"error func Update, method Begin by path internal/adapter/psqlRepo/profile/profile.go", zap.Error(err))
		return nil, err
	}
	defer tx.Rollback()
	query := "UPDATE profiles SET display_name=$1, birthday=$2, gender=$3, search_gender=$4, location=$5," +
		" description=$6, height=$7, weight=$8, looking_for=$9, is_blocked=$10, is_premium=$11, is_show_distance=$12," +
		" is_invisible=$13, updated_at=$14, last_online=$15 WHERE id=$16"
	_, err = r.db.ExecContext(ctx, query, p.DisplayName, p.Birthday, p.Gender, p.SearchGender, p.Location,
		p.Description, p.Height, p.Weight, p.LookingFor, p.IsBlocked, p.IsPremium, p.IsShowDistance,
		p.IsInvisible, p.UpdatedAt, p.LastOnline, p.ID)
	if err != nil {
		r.logger.Debug(
			"error func Update, method ExecContext by path internal/adapter/psqlRepo/profile/profile.go",
			zap.Error(err))
		return nil, err
	}
	tx.Commit()
	return p, nil
}

func (r *RepositoryProfile) SelectList(
	ctx context.Context, qp *profile.QueryParamsProfileList) (*profile.ResponseListProfile, error) {
	query := "SELECT id, display_name, birthday, gender, search_gender, location, description, height, weight," +
		" looking_for, is_deleted, is_blocked, is_premium, is_show_distance, is_invisible, created_at, updated_at," +
		" last_online FROM profiles WHERE is_deleted = false"
	countQuery := "SELECT COUNT(*) FROM profiles WHERE is_deleted = false"
	limit := qp.Limit
	page := qp.Page
	// get totalItems
	totalItems, err := pagination.GetTotalItems(ctx, r.db, countQuery)
	if err != nil {
		r.logger.Debug(
			"error func SelectList, method GetTotalItems by path internal/adapter/psqlRepo/profile/profile.go",
			zap.Error(err))
		return nil, err
	}
	// pagination
	query = pagination.ApplyPagination(query, page, limit)
	countQuery = pagination.ApplyPagination(countQuery, page, limit)
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		r.logger.Debug(
			"error func SelectList, method QueryContext by path internal/adapter/psqlRepo/profile/profile.go",
			zap.Error(err))
		return nil, err
	}
	defer rows.Close()
	list := make([]*profile.ContentListProfile, 0)
	for rows.Next() {
		p := profile.Profile{}
		err := rows.Scan(&p.ID, &p.DisplayName, &p.Birthday, &p.Gender, &p.SearchGender, &p.Location, &p.Description,
			&p.Height, &p.Weight, &p.LookingFor, &p.IsDeleted, &p.IsBlocked, &p.IsPremium, &p.IsShowDistance,
			&p.IsInvisible, &p.CreatedAt, &p.UpdatedAt, &p.LastOnline)
		if err != nil {
			r.logger.Debug("error func SelectList, method Scan by path internal/adapter/psqlRepo/profile/profile.go",
				zap.Error(err))
			continue
		}
		images, err := r.SelectListPublicImage(ctx, p.ID)
		if err != nil {
			r.logger.Debug("error func SelectList, method SelectListPublicImage by path"+
				" internal/adapter/psqlRepo/profile/profile.go", zap.Error(err))
			continue
		}
		lp := profile.ContentListProfile{
			ID:         p.ID,
			LastOnline: p.LastOnline,
			Image:      nil,
		}
		if len(images) > 0 {
			i := profile.ResponseImageProfile{
				Url: images[0].Url,
			}
			lp.Image = &i
		}
		list = append(list, &lp)
	}
	paging := pagination.GetPagination(limit, page, totalItems)
	response := profile.ResponseListProfile{
		Pagination: paging,
		Content:    list,
	}
	return &response, nil
}

func (r *RepositoryProfile) FindById(ctx context.Context, id uint64) (*profile.Profile, error) {
	p := profile.Profile{}
	query := `SELECT id, display_name, birthday, gender, search_gender, location, description, height, weight,
       looking_for, is_deleted, is_blocked, is_premium, is_show_distance, is_invisible, created_at, updated_at,
       last_online
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
	err := row.Scan(&p.ID, &p.DisplayName, &p.Birthday, &p.Gender, &p.SearchGender, &p.Location, &p.Description,
		&p.Height, &p.Weight, &p.LookingFor, &p.IsDeleted, &p.IsBlocked, &p.IsPremium, &p.IsShowDistance,
		&p.IsInvisible, &p.CreatedAt, &p.UpdatedAt, &p.LastOnline)
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

func (r *RepositoryProfile) UpdateTelegram(
	ctx context.Context, p *profile.TelegramProfile) (*profile.TelegramProfile, error) {
	tx, err := r.db.Begin()
	if err != nil {
		r.logger.Debug("error func UpdateTelegram, method Begin by path internal/adapter/psqlRepo/profile/profile.go",
			zap.Error(err))
		return nil, err
	}
	defer tx.Rollback()
	query := "UPDATE profile_telegram SET username=$1, first_name=$2, last_name=$3, language_code=$4," +
		" allows_write_to_pm=$5 WHERE id=$6"
	_, err = r.db.ExecContext(ctx, query, p.UserName, p.Firstname, p.Lastname, p.LanguageCode, p.AllowsWriteToPm, &p.ID)
	if err != nil {
		r.logger.Debug(
			"error func UpdateTelegram, method QueryRowContext by path internal/adapter/psqlRepo/profile/profile.go",
			zap.Error(err))
		return nil, err
	}
	tx.Commit()
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
	err := row.Scan(&p.ID, &p.ProfileID, &p.TelegramID, &p.UserName, &p.Firstname, &p.Lastname, &p.LanguageCode,
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

func (r *RepositoryProfile) UpdateImage(ctx context.Context, p *profile.ImageProfile) (*profile.ImageProfile, error) {
	tx, err := r.db.Begin()
	if err != nil {
		r.logger.Debug("error func UpdateImage, method Begin by path internal/adapter/psqlRepo/profile/profile.go",
			zap.Error(err))
		return nil, err
	}
	defer tx.Rollback()
	query := "UPDATE profile_images SET name=$1, url=$2, size=$3, updated_at=$4, is_deleted=$5, is_blocked=$6," +
		" is_primary=$7, is_private=$8 WHERE id=$9"
	_, err = r.db.ExecContext(ctx, query, p.Name, p.Url, p.Size, p.UpdatedAt, p.IsDeleted, p.IsBlocked, p.IsPrimary,
		p.IsPrivate, &p.ID)
	if err != nil {
		r.logger.Debug(
			"error func UpdateImage method QueryRowContext by path internal/adapter/psqlRepo/profile/profile.go",
			zap.Error(err))
		return nil, err
	}
	tx.Commit()
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

func (r *RepositoryProfile) CheckIfCommonImageExists(
	ctx context.Context, profileID uint64, fileName string) (bool, uint64, error) {
	var imageID uint64
	query := "SELECT id" +
		" FROM profile_images WHERE profile_id=$1 AND name=$2"
	err := r.db.QueryRowContext(ctx, query, profileID, fileName).Scan(&imageID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, 0, nil
		}
		return false, 0, err
	}
	return true, imageID, nil
}
