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
