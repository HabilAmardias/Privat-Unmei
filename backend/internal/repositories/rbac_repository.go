package repositories

import (
	"context"
	"database/sql"
	"errors"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/db"
	"privat-unmei/internal/entity"
)

type RBACRepository struct {
	DB *db.CustomDB
}

func CreateRBACRepository(db *db.CustomDB) *RBACRepository {
	return &RBACRepository{db}
}

func (rr *RBACRepository) CheckRoleAccess(ctx context.Context, rbac *entity.Rbac, role int, permission int, resource int) error {
	var driver RepoDriver
	driver = rr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	SELECT id
	FROM rbac
	WHERE role_id = $1 AND permission_id = $2 AND resource_id = $3 AND deleted_at IS NULL
	`

	row := driver.QueryRow(query, role, permission, resource)
	if err := row.Scan(&rbac.ID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return customerrors.NewError("user does not have access", err, customerrors.Unauthenticate)
		}
		return customerrors.NewError("failed to authorize", err, customerrors.DatabaseExecutionError)
	}
	return nil
}
