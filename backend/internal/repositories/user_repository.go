package repositories

import (
	"context"
	"database/sql"
	"errors"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/db"
	"privat-unmei/internal/entity"
)

type UserRepositoryImpl struct {
	DB *db.CustomDB
}

func CreateUserRepository(db *db.CustomDB) *UserRepositoryImpl {
	return &UserRepositoryImpl{db}
}

func (ur *UserRepositoryImpl) UpdateUserProfile(ctx context.Context, queryEntity *entity.UpdateUserQuery, id string) error {
	var driver RepoDriver
	driver = ur.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	UPDATE users
	SET
		name = COALESCE($1, name),
		bio = COALESCE($2, bio),
		profile_image = COALESCE($3, profile_image),
		updated_at = NOW()
	WHERE id = $4 AND deleted_at IS NULL
	`
	_, err := driver.Exec(
		query,
		queryEntity.Name,
		queryEntity.Bio,
		queryEntity.ProfileImage,
		id,
	)
	if err != nil {
		return customerrors.NewError(
			"failed to update profile",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (ur *UserRepositoryImpl) UpdateUserPassword(ctx context.Context, password string, id string) error {
	var driver RepoDriver
	driver = ur.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	UPDATE users SET password_hash = $1, updated_at = NOW() WHERE id = $2 AND deleted_at IS NULL
	`
	_, err := driver.Exec(query, password, id)
	if err != nil {
		return customerrors.NewError(
			"Failed to update user password",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (ur *UserRepositoryImpl) UpdateUserStatus(ctx context.Context, status string, id string) error {
	var driver RepoDriver
	driver = ur.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	UPDATE users SET status = $1, updated_at = NOW() WHERE id = $2 AND deleted_at IS NULL
	`
	_, err := driver.Exec(query, status, id)
	if err != nil {
		return customerrors.NewError(
			"Failed to update user status",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
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

	if err := row.Scan(&user.ID); err != nil {
		return customerrors.NewError(
			"failed to create account",
			err,
			customerrors.DatabaseExecutionError,
		)
	}

	return nil
}

func (ur *UserRepositoryImpl) FindByID(ctx context.Context, id string, user *entity.User) error {
	var driver RepoDriver
	driver = ur.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	SELECT id, name, email, password_hash, bio, profile_image, status, created_at, updated_at, deleted_at FROM users
	WHERE id = $1 AND deleted_at IS NULL
	`
	row := driver.QueryRow(query, id)
	if err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Bio,
		&user.ProfileImage,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return customerrors.NewError(
				customerrors.UserNotFound,
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

func (ur *UserRepositoryImpl) FindByEmail(ctx context.Context, email string, user *entity.User) error {
	var driver RepoDriver
	driver = ur.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	SELECT id, name, email, password_hash, bio, profile_image, status, created_at, updated_at, deleted_at FROM users
	WHERE email = $1 AND deleted_at IS NULL
	`
	row := driver.QueryRow(query, email)
	if err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Bio,
		&user.ProfileImage,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return customerrors.NewError(
				customerrors.UserNotFound,
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

func (ar *UserRepositoryImpl) DeleteUser(ctx context.Context, id string) error {
	var driver RepoDriver
	driver = ar.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	UPDATE users
	SET
		deleted_at = NOW(),
		updated_at = NOW()
	WHERE id = $1 AND deleted_at IS NULL
	`
	_, err := driver.Exec(query, id)
	if err != nil {
		return customerrors.NewError("failed to delete user", err, customerrors.DatabaseExecutionError)
	}
	return nil
}
