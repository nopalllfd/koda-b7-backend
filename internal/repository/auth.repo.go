package repository

import (
	"backend-golang/internal/model"
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type AuthRepository struct {
	dbtx DBTX
	rc   *redis.Client
}

func NewAuthRepo(db *pgxpool.Pool, rc *redis.Client) *AuthRepository {
	return &AuthRepository{
		dbtx: db,
		rc:   rc,
	}
}

func (ar *AuthRepository) GetUserByEmail(
	ctx context.Context,
	email string,
) (*model.User, error) {

	sql := `
	SELECT 
		u.id,
		u.email,
		u.password,
		COALESCE(u.pin, '') AS pin,
		COALESCE(p.full_name, ''),
		COALESCE(p.photo, '')
	FROM users u
	JOIN profiles p ON p.user_id = u.id
	WHERE u.email = $1
	`

	var user model.User

	err := ar.dbtx.QueryRow(
		ctx,
		sql,
		email,
	).Scan(
		&user.Id,
		&user.Email,
		&user.Password,
		&user.Pin,
		&user.FullName,
		&user.Photo,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ar *AuthRepository) FindByEmail(
	ctx context.Context,
	email string,
) (bool, error) {

	sql := `
	SELECT EXISTS(
		SELECT 1 FROM users WHERE email = $1
	)
	`

	var exists bool

	if err := ar.dbtx.QueryRow(
		ctx,
		sql,
		email,
	).Scan(&exists); err != nil {

		return false, err
	}

	return exists, nil
}

func (ar *AuthRepository) Create(
	ctx context.Context,
	email string,
	password string,
) (int, error) {

	sql := `
	INSERT INTO users (
		email,
		password
	)
	VALUES ($1, $2)
	RETURNING id
	`

	var userID int

	err := ar.dbtx.QueryRow(
		ctx,
		sql,
		email,
		password,
	).Scan(&userID)

	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (ar *AuthRepository) SetPin(
	ctx context.Context,
	pin string,
	id int,
) error {

	sql := `
	UPDATE users
	SET pin = $1
	WHERE id = $2
	`

	if _, err := ar.dbtx.Exec(
		ctx,
		sql,
		pin,
		id,
	); err != nil {

		return err
	}

	return nil
}

func (ar *AuthRepository) GetUserPassword(
	ctx context.Context,
	id int,
) (string, error) {

	sql := `
	SELECT password
	FROM users
	WHERE id = $1
	`

	var userPassword string

	if err := ar.dbtx.QueryRow(
		ctx,
		sql,
		id,
	).Scan(&userPassword); err != nil {

		return "", err
	}

	return userPassword, nil
}

func (ar *AuthRepository) SetPassword(
	ctx context.Context,
	newPass string,
	id int,
) error {

	sql := `
	UPDATE users
	SET password = $1
	WHERE id = $2
	`

	if _, err := ar.dbtx.Exec(
		ctx,
		sql,
		newPass,
		id,
	); err != nil {

		return err
	}

	return nil
}

func (ar *AuthRepository) BlacklistToken(
	ctx context.Context,
	token string,
	expiredAt time.Time,
) error {

	ttl := time.Until(expiredAt)

	if ttl < 0 {
		ttl = 0
	}

	return ar.rc.Set(
		ctx,
		"bl:"+token,
		"true",
		ttl,
	).Err()
}

func (ar *AuthRepository) SaveTokenForgotPass(
	ctx context.Context,
	token string,
	userID int,
) error {

	sql := `
	INSERT INTO password_reset_tokens (
		user_id,
		token,
		expires_at
	)
	VALUES (
		$1,
		$2,
		NOW() + INTERVAL '15 minutes'
	)
	`

	if _, err := ar.dbtx.Exec(
		ctx,
		sql,
		userID,
		token,
	); err != nil {

		return err
	}

	return nil
}

func (ar *AuthRepository) TokenValidCheck(
	ctx context.Context,
	token string,
) (int, error) {

	sql := `
	SELECT user_id
	FROM password_reset_tokens
	WHERE token = $1
	AND expires_at > NOW()
	AND used_at IS NULL
	`

	var userID int

	if err := ar.dbtx.QueryRow(
		ctx,
		sql,
		token,
	).Scan(&userID); err != nil {

		return 0, err
	}

	return userID, nil
}

func (ar *AuthRepository) ValidateChangedPassword(
	ctx context.Context,
	token string,
) error {

	sql := `
	UPDATE password_reset_tokens
	SET used_at = NOW()
	WHERE token = $1
	`

	if _, err := ar.dbtx.Exec(
		ctx,
		sql,
		token,
	); err != nil {

		return err
	}

	return nil
}
