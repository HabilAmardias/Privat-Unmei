package repositories

import (
	"context"
	"database/sql"
	"errors"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/entity"
)

type UserRepositoryImpl struct {
	DB *sql.DB
}

func CreateUserRepository(db *sql.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{db}
}

func (ur *UserRepositoryImpl) AddNewUser(ctx context.Context, user *entity.User) error {
	var driver RepoDriver
	driver = ur.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	INSERT INTO users (name, email, password_hash, bio, profile_image, status)
	VALUES
	($1, $2, $3, $4, $5, $6)
	RETURNING id;
	`
	row := driver.QueryRow(query,
		user.Name,
		user.Email,
		user.Password,
		user.Bio,
		user.ProfileImage,
		user.Status)

	if err := row.Scan(user.ID); err != nil {
		return customerrors.NewError(
			"failed to create account",
			err,
			customerrors.DatabaseExecutionError,
		)
	}

	return nil
}

func (ur *UserRepositoryImpl) FindByEmail(ctx context.Context, email string, user *entity.User) error {
	var driver RepoDriver
	driver = ur.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	SELECT id, name, email, password_hash, bio, profile_image, status, created_at, updated_at, deleted_at FROM users
	WHERE email = $1
	`
	row := driver.QueryRow(query, email)
	if err := row.Scan(
		user.ID,
		user.Name,
		user.Email,
		user.Password,
		user.Bio,
		user.ProfileImage,
		user.Status,
		user.CreatedAt,
		user.UpdatedAt,
		user.DeletedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return customerrors.NewError(
				"user not found",
				err,
				customerrors.ItemNotExist,
			)
		}
		return customerrors.NewError(
			"failed to get user data",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}
