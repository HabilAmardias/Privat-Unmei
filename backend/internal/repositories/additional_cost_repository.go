package repositories

import (
	"context"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/db"
)

type AdditionalCostRepositoryImpl struct {
	DB *db.CustomDB
}

func CreateAdditionalCostRepository(db *db.CustomDB) *AdditionalCostRepositoryImpl {
	return &AdditionalCostRepositoryImpl{db}
}

func (acr *AdditionalCostRepositoryImpl) FindByName(ctx context.Context, name string, count *int64) error {
	var driver RepoDriver = acr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	SELECT count(*)
	FROM additional_costs
	WHERE LOWER(name) = LOWER($1) AND deleted_at IS NULL
	`
	if err := driver.QueryRow(query, name).Scan(
		count,
	); err != nil {
		return customerrors.NewError(
			"failed to get operational cost",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (acr *AdditionalCostRepositoryImpl) CreateOperationalCost(ctx context.Context, name string, amount float64, id *int) error {
	var driver RepoDriver = acr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	INSERT INTO additional_costs (name, amount)
	VALUES
	($1, $2)
	ON CONFLICT (name)
	DO UPDATE SET
		deleted_at = NULL,
		amount = EXCLUDED.amount,
		updated_at = CURRENT_TIMESTAMP
	RETURNING (id)
	`
	if err := driver.QueryRow(query, name, amount).Scan(id); err != nil {
		return customerrors.NewError(
			"failed to create new operational cost",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (acr *AdditionalCostRepositoryImpl) GetOperationalCost(ctx context.Context, totalCost *float64) error {
	var driver RepoDriver = acr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	SELECT
		COALESCE(SUM(amount), 0)
	FROM additional_costs
	WHERE deleted_at IS NULL
	`
	if err := driver.QueryRow(query).Scan(totalCost); err != nil {
		return customerrors.NewError(
			"failed to get operational cost",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}
