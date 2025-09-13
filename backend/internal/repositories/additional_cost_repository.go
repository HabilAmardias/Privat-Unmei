package repositories

import (
	"context"
	"database/sql"
	"errors"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/db"
	"privat-unmei/internal/entity"
)

type AdditionalCostRepositoryImpl struct {
	DB *db.CustomDB
}

func CreateAdditionalCostRepository(db *db.CustomDB) *AdditionalCostRepositoryImpl {
	return &AdditionalCostRepositoryImpl{db}
}

func (acr *AdditionalCostRepositoryImpl) GetAllAdditionalCost(ctx context.Context, limit int, page int, totalRow *int64, costs *[]entity.GetAdditionalCostQuery) error {
	var driver RepoDriver = acr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	countQuery := `
	SELECT
		count(*)
	FROM additional_costs
	WHERE deleted_at IS NULL
	`
	if err := driver.QueryRow(countQuery).Scan(
		totalRow,
	); err != nil {
		return customerrors.NewError(
			"failed to get cost",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	query := `
	SELECT
		id,
		name,
		amount
	FROM additional_costs
	WHERE deleted_at IS NULL
	LIMIT $1
	OFFSET $2
	`

	rows, err := driver.Query(query, limit, limit*(page-1))
	if err != nil {
		return customerrors.NewError(
			"failed to get cost",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	defer rows.Close()
	for rows.Next() {
		var item entity.GetAdditionalCostQuery
		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Amount,
		); err != nil {
			return customerrors.NewError(
				"failed to get cost",
				err,
				customerrors.DatabaseExecutionError,
			)
		}
		*costs = append(*costs, item)
	}
	return nil
}

func (acr *AdditionalCostRepositoryImpl) DeleteCost(ctx context.Context, id int) error {
	var driver RepoDriver = acr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	UPDATE additional_costs
	SET
		deleted_at = CURRENT_TIMESTAMP,
		updated_at = CURRENT_TIMESTAMP
	WHERE id = $1 AND deleted_at IS NULL
	`
	_, err := driver.Exec(query, id)
	if err != nil {
		return customerrors.NewError(
			"failed to delete cost",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (acr *AdditionalCostRepositoryImpl) UpdateCostAmount(ctx context.Context, id int, amount *float64) error {
	var driver RepoDriver = acr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	UPDATE additional_costs
	SET
		amount = COALESCE($1, amount),
		updated_at = CURRENT_TIMESTAMP
	WHERE id = $2 AND deleted_at IS NULL
	`
	_, err := driver.Exec(query, amount, id)
	if err != nil {
		return customerrors.NewError(
			"failed to update cost amount",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (acr *AdditionalCostRepositoryImpl) FindByID(ctx context.Context, id int, cost *entity.AdditionalCost) error {
	var driver RepoDriver = acr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	SELECT
		id,
		name,
		amount,
		created_at,
		updated_at,
		deleted_at
	FROM additional_costs
	WHERE id = $1 AND deleted_at IS NULL
	`
	if err := driver.QueryRow(query, id).Scan(
		&cost.ID,
		&cost.Name,
		&cost.Amount,
		&cost.CreatedAt,
		&cost.UpdatedAt,
		&cost.DeletedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return customerrors.NewError(
				"additional cost does not exist",
				err,
				customerrors.ItemNotExist,
			)
		}
		return customerrors.NewError(
			"failed to get additional cost",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
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
