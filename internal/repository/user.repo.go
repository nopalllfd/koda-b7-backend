package repository

import (
	errs "backend-golang/internal/err"
	"backend-golang/internal/model"
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (ur *UserRepository) Create(ctx context.Context, userID int, full_name, photo, phone *string) error {
	sql := "INSERT into profiles (user_id, full_name, photo, phone) VALUES ($1,$2,$3,$4)"
	if _, err := ur.db.Exec(ctx, sql, userID, full_name, photo, phone); err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) GetProfile(ctx context.Context, id int) (model.Profile, error) {
	sql := "SELECT user_id, COALESCE(full_name, '') as full_name, COALESCE(photo, '') as photo, COALESCE(phone, '') as phone FROM profiles WHERE user_id = $1"

	var user model.Profile
	if err := ur.db.QueryRow(ctx, sql, id).Scan(&user.User_id, &user.FullName, &user.Photo, &user.Phone); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Profile{}, errs.ErrProfileNotFound
		}
		return model.Profile{}, errs.ErrInternalServer
	}
	return user, nil
}

func (ur *UserRepository) Edit(ctx context.Context, full_name, photo, phone *string, userID int) (int64, error) {
	log.Println(userID, *full_name)
	sql := `
UPDATE profiles
SET
	full_name = $1,
	photo = $2,
	phone = $3
WHERE user_id = $4
`
	commandTag, err := ur.db.Exec(ctx, sql, *full_name, *photo, *phone, userID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505": // Kode Postgres untuk Unique Violation
				return 0, errs.ErrPhoneAlreadyUsed
			case "22001": // Kode Postgres untuk String Data Right Truncation (input terlalu panjang)
				return 0, errs.ErrInvalidInput
			}
		}
		// Error database lainnya (koneksi putus, dll)
		return 0, errs.ErrInternalServer
	}
	return commandTag.RowsAffected(), nil
}

// func (ur *UserRepository) GetDashboard(ctx context.Context, id int)
