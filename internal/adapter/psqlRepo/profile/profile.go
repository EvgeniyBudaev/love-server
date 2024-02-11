package profile

import (
	"context"
	"database/sql"
	"errors"
	"github.com/EvgeniyBudaev/love-server/internal/entity/pagination"
	"github.com/EvgeniyBudaev/love-server/internal/entity/profile"
	"github.com/EvgeniyBudaev/love-server/internal/logger"
	"go.uber.org/zap"
	"strconv"
	"time"
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
	birthday := p.Birthday.Format("2006-01-02")
	query := "INSERT INTO profiles (display_name, birthday, gender, search_gender, location, description, height," +
		" weight, looking_for, is_deleted, is_blocked, is_premium, is_show_distance, is_invisible, created_at," +
		" updated_at, last_online) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15," +
		" $16, $17) RETURNING id"
	err := r.db.QueryRowContext(ctx, query, p.DisplayName, birthday, p.Gender, p.SearchGender, p.Location,
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

func (r *RepositoryProfile) UpdateLastOnline(ctx context.Context, profileID uint64) error {
	tx, err := r.db.Begin()
	if err != nil {
		r.logger.Debug("error func UpdateLastOnline, method Begin by path"+
			" internal/adapter/psqlRepo/profile/profile.go", zap.Error(err))
		return err
	}
	defer tx.Rollback()
	query := "UPDATE profiles SET last_online=$1 WHERE id=$2"
	_, err = r.db.ExecContext(ctx, query, time.Now(), profileID)
	if err != nil {
		r.logger.Debug(
			"error func UpdateLastOnline, method ExecContext by path internal/adapter/psqlRepo/profile/profile.go",
			zap.Error(err))
		return err
	}
	tx.Commit()
	return nil
}

func (r *RepositoryProfile) Delete(ctx context.Context, p *profile.Profile) (*profile.Profile, error) {
	tx, err := r.db.Begin()
	if err != nil {
		r.logger.Debug(
			"error func Delete, method Begin by path internal/adapter/psqlRepo/profile/profile.go", zap.Error(err))
		return nil, err
	}
	defer tx.Rollback()
	query := "UPDATE profiles SET display_name=$1, birthday=$2, gender=$3, search_gender=$4, location=$5," +
		" description=$6, height=$7, weight=$8, looking_for=$9, is_deleted=$10, is_blocked=$11, is_premium=$12," +
		" is_show_distance=$13, is_invisible=$14, updated_at=$15, last_online=$16 WHERE id=$17"
	_, err = r.db.ExecContext(ctx, query, p.DisplayName, p.Birthday, p.Gender, p.SearchGender, p.Location,
		p.Description, p.Height, p.Weight, p.LookingFor, p.IsDeleted, p.IsBlocked, p.IsPremium, p.IsShowDistance,
		p.IsInvisible, p.UpdatedAt, p.LastOnline, p.ID)
	if err != nil {
		r.logger.Debug(
			"error func Delete, method ExecContext by path internal/adapter/psqlRepo/profile/profile.go",
			zap.Error(err))
		return nil, err
	}
	tx.Commit()
	return p, nil
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

func (r *RepositoryProfile) FindByTelegramId(ctx context.Context, telegramID uint64) (*profile.Profile, error) {
	p := profile.Profile{}
	query := `SELECT p.id, p.display_name, p.birthday, p.gender, p.search_gender, p.location, p.description, p.height,
       p.weight, p.looking_for, p.is_deleted, p.is_blocked, p.is_premium, p.is_show_distance, p.is_invisible,
       p.created_at, p.updated_at,  p.last_online
			  FROM profiles p
			  JOIN profile_telegram pt ON p.id = pt.profile_id
			  WHERE pt.telegram_id = $1`
	row := r.db.QueryRowContext(ctx, query, telegramID)
	if row == nil {
		err := errors.New("no rows found")
		r.logger.Debug(
			"error func FindByTelegramId, method QueryRowContext by path internal/adapter/psqlRepo/profile/profile.go",
			zap.Error(err))
		return nil, err
	}
	err := row.Scan(&p.ID, &p.DisplayName, &p.Birthday, &p.Gender, &p.SearchGender, &p.Location, &p.Description,
		&p.Height, &p.Weight, &p.LookingFor, &p.IsDeleted, &p.IsBlocked, &p.IsPremium, &p.IsShowDistance,
		&p.IsInvisible, &p.CreatedAt, &p.UpdatedAt, &p.LastOnline)
	if err != nil {
		r.logger.Debug("error func FindByTelegramId, method Scan by path internal/adapter/psqlRepo/profile/profile.go",
			zap.Error(err))
		return nil, err
	}
	return &p, nil
}

func (r *RepositoryProfile) SelectList(
	ctx context.Context, qp *profile.QueryParamsProfileList) (*profile.ResponseListProfile, error) {
	ageFromInt, err := strconv.Atoi(qp.AgeFrom)
	if err != nil {
		return nil, err
	}
	ageToInt, err := strconv.Atoi(qp.AgeTo)
	if err != nil {
		return nil, err
	}
	birthYearStart := time.Now().Year() - ageToInt - 1
	birthYearEnd := time.Now().Year() - ageFromInt
	birthdateFrom := time.Date(birthYearStart, time.January, 1, 0, 0, 0, 0, time.UTC)
	birthdateTo := time.Date(birthYearEnd, time.December, 31, 23, 59, 59, 999999999, time.UTC)
	query := "SELECT id, display_name, birthday, gender, search_gender, location, description, height, weight," +
		" looking_for, is_deleted, is_blocked, is_premium, is_show_distance, is_invisible, created_at, updated_at," +
		" last_online FROM profiles WHERE is_deleted=false AND is_blocked=false AND birthday BETWEEN $1 AND $2 AND ($3 = 'all' OR gender=$3)"
	countQuery := "SELECT COUNT(*) FROM profiles WHERE is_deleted=false AND is_blocked=false AND birthday BETWEEN $1 AND $2 AND ($3 = 'all' OR gender=$3)"
	limit := qp.Limit
	page := qp.Page
	// get totalItems
	totalItems, err := pagination.GetTotalItems(ctx, r.db, countQuery, birthdateFrom, birthdateTo, qp.SearchGender)
	if err != nil {
		r.logger.Debug(
			"error func SelectList, method GetTotalItems by path internal/adapter/psqlRepo/profile/profile.go",
			zap.Error(err))
		return nil, err
	}
	// pagination
	query = pagination.ApplyPagination(query, page, limit)
	countQuery = pagination.ApplyPagination(countQuery, page, limit)
	queryParams := []interface{}{birthdateFrom, birthdateTo, qp.SearchGender}
	rows, err := r.db.QueryContext(ctx, query, queryParams...)
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

func (r *RepositoryProfile) DeleteTelegram(
	ctx context.Context, p *profile.TelegramProfile) (*profile.TelegramProfile, error) {
	tx, err := r.db.Begin()
	if err != nil {
		r.logger.Debug("error func DeleteTelegram, method Begin by path internal/adapter/psqlRepo/profile/profile.go",
			zap.Error(err))
		return nil, err
	}
	defer tx.Rollback()
	query := "UPDATE profile_telegram SET telegram_id=$1, username=$2, first_name=$3, last_name=$4, language_code=$5," +
		" allows_write_to_pm=$6, query_id=$7 WHERE id=$8"
	_, err = r.db.ExecContext(ctx, query, p.TelegramID, p.UserName, p.Firstname, p.Lastname, p.LanguageCode,
		p.AllowsWriteToPm, p.QueryID, &p.ID)
	if err != nil {
		r.logger.Debug(
			"error func DeleteTelegram, method QueryRowContext by path internal/adapter/psqlRepo/profile/profile.go",
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
			"error func FindTelegramById, method QueryRowContext by path internal/adapter/psqlRepo/profile/profile.go",
			zap.Error(err))
		return nil, err
	}
	err := row.Scan(&p.ID, &p.ProfileID, &p.TelegramID, &p.UserName, &p.Firstname, &p.Lastname, &p.LanguageCode,
		&p.AllowsWriteToPm, &p.QueryID)
	if err != nil {
		r.logger.Debug("error func FindTelegramById, method Scan by path internal/adapter/psqlRepo/profile/profile.go",
			zap.Error(err))
		return nil, err
	}
	return &p, nil
}

func (r *RepositoryProfile) AddNavigator(
	ctx context.Context, p *profile.NavigatorProfile) (*profile.NavigatorProfile, error) {
	query := "INSERT INTO profile_navigators (profile_id, latitude, longitude) VALUES ($1, $2, $3) RETURNING id"
	err := r.db.QueryRowContext(ctx, query, p.ProfileID, p.Latitude, p.Longitude).Scan(&p.ID)
	if err != nil {
		r.logger.Debug(
			"error func AddNavigator, method QueryRowContext by path internal/adapter/psqlRepo/profile/profile.go",
			zap.Error(err))
		return nil, err
	}
	return p, nil
}

func (r *RepositoryProfile) UpdateNavigator(
	ctx context.Context, p *profile.NavigatorProfile) (*profile.NavigatorProfile, error) {
	tx, err := r.db.Begin()
	if err != nil {
		r.logger.Debug("error func UpdateNavigator, method Begin by path internal/adapter/psqlRepo/profile/profile.go",
			zap.Error(err))
		return nil, err
	}
	defer tx.Rollback()
	query := "UPDATE profile_navigators SET latitude=$1, longitude=$2 WHERE profile_id=$3"
	_, err = r.db.ExecContext(ctx, query, p.Latitude, p.Longitude, &p.ProfileID)
	if err != nil {
		r.logger.Debug(
			"error func UpdateNavigator, method QueryRowContext by path internal/adapter/psqlRepo/profile/profile.go",
			zap.Error(err))
		return nil, err
	}
	tx.Commit()
	return p, nil
}

func (r *RepositoryProfile) FindNavigatorById(
	ctx context.Context, profileID uint64) (*profile.NavigatorProfile, error) {
	p := profile.NavigatorProfile{}
	query := `SELECT id, profile_id, latitude, longitude
			  FROM profile_navigators
			  WHERE profile_id = $1`
	row := r.db.QueryRowContext(ctx, query, profileID)
	if row == nil {
		err := errors.New("no rows found")
		r.logger.Debug(
			"error func FindNavigatorById, method QueryRowContext by path internal/adapter/psqlRepo/profile/profile.go",
			zap.Error(err))
		return nil, err
	}
	err := row.Scan(&p.ID, &p.ProfileID, &p.Latitude, &p.Longitude)
	if err != nil {
		r.logger.Debug("error func FindNavigatorById, method Scan by path internal/adapter/psqlRepo/profile/profile.go",
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

func (r *RepositoryProfile) DeleteImage(ctx context.Context, p *profile.ImageProfile) (*profile.ImageProfile, error) {
	tx, err := r.db.Begin()
	if err != nil {
		r.logger.Debug("error func DeleteImage, method Begin by path internal/adapter/psqlRepo/profile/profile.go",
			zap.Error(err))
		return nil, err
	}
	defer tx.Rollback()
	query := "UPDATE profile_images SET is_deleted=$1 WHERE id=$2"
	_, err = r.db.ExecContext(ctx, query, p.IsDeleted, &p.ID)
	if err != nil {
		r.logger.Debug(
			"error func DeleteImage method QueryRowContext by path internal/adapter/psqlRepo/profile/profile.go",
			zap.Error(err))
		return nil, err
	}
	tx.Commit()
	return p, nil
}

func (r *RepositoryProfile) FindImageById(ctx context.Context, imageID uint64) (*profile.ImageProfile, error) {
	p := profile.ImageProfile{}
	query := `SELECT id, profile_id, name, url, size, created_at, updated_at, is_deleted, is_blocked, is_primary,
       is_private
			  FROM profile_images
			  WHERE id=$1 AND is_deleted=false AND is_blocked=false`
	row := r.db.QueryRowContext(ctx, query, imageID)
	if row == nil {
		err := errors.New("no rows found")
		r.logger.Debug(
			"error func FindImageById, method QueryRowContext by path internal/adapter/psqlRepo/profile/profile.go",
			zap.Error(err))
		return nil, err
	}
	err := row.Scan(&p.ID, &p.ProfileID, &p.Name, &p.Url, &p.Size, &p.CreatedAt, &p.UpdatedAt,
		&p.IsDeleted, &p.IsBlocked, &p.IsPrimary, &p.IsPrivate)
	if err != nil {
		r.logger.Debug("error func FindImageById, method Scan by path internal/adapter/psqlRepo/profile/profile.go",
			zap.Error(err))
		return nil, err
	}
	return &p, nil
}

func (r *RepositoryProfile) SelectListPublicImage(
	ctx context.Context, profileID uint64) ([]*profile.ImageProfile, error) {
	query := `SELECT id, profile_id, name, url, size, created_at, updated_at, is_deleted, is_blocked, is_primary,
       is_private
	FROM profile_images
	WHERE profile_id=$1 AND is_deleted=false AND is_blocked=false AND is_private=false`
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

func (r *RepositoryProfile) SelectListImage(
	ctx context.Context, profileID uint64) ([]*profile.ImageProfile, error) {
	query := `SELECT id, profile_id, name, url, size, created_at, updated_at, is_deleted, is_blocked, is_primary,
       is_private
	FROM profile_images
	WHERE profile_id=$1`
	rows, err := r.db.QueryContext(ctx, query, profileID)
	if err != nil {
		r.logger.Debug("error func SelectListImage,"+
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
			r.logger.Debug("error func SelectListImage,"+
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
