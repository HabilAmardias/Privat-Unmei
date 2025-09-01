package repositories

import (
	"context"
	"database/sql"
	"errors"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/db"
	"privat-unmei/internal/entity"
)

type AdminRepositoryImpl struct {
	DB *db.CustomDB
}

func CreateAdminRepository(db *db.CustomDB) *AdminRepositoryImpl {
	return &AdminRepositoryImpl{db}
}

func (ar *AdminRepositoryImpl) FindByID(ctx context.Context, id string, admin *entity.Admin) error {
	var driver RepoDriver
	driver = ar.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}

	query := `
	SELECT id, created_at, updated_at, deleted_at
	FROM admins
	WHERE id = $1 AND deleted_at IS NULL
	`

	row := driver.QueryRow(query, id)
	if err := row.Scan(
		&admin.ID,
		&admin.CreatedAt,
		&admin.UpdatedAt,
		&admin.DeletedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return customerrors.NewError(
				customerrors.UserNotFound,
				err,
				customerrors.ItemNotExist,
			)
		}
		return customerrors.NewError(
			"failed to find user",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}
