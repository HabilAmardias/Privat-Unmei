package repositories

import (
	"context"
	"database/sql"
	"errors"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/db"
	"privat-unmei/internal/entity"
)

type DiscountRepositoryImpl struct {
	DB *db.CustomDB
}

func CreateDiscountRepository(db *db.CustomDB) *DiscountRepositoryImpl {
	return &DiscountRepositoryImpl{db}
}

func (dr *DiscountRepositoryImpl) GetAllDiscount(ctx context.Context, limit int, page int, totalRow *int64, discounts *[]entity.GetDiscountQuery) error {
	var driver RepoDriver = dr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	countQuery := `
	SELECT
		count(*)
	FROM discounts
	WHERE deleted_at IS NULL
	`
	if err := driver.QueryRow(countQuery).Scan(
		totalRow,
	); err != nil {
		return customerrors.NewError(
			"failed to get discount",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	query := `
	SELECT
		id,
		number_of_participant,
		amount
	FROM discounts
	WHERE deleted_at IS NULL
	LIMIT $1
	OFFSET $2
	`

	rows, err := driver.Query(query, limit, limit*(page-1))
	if err != nil {
		return customerrors.NewError(
			"failed to get discount",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	defer rows.Close()
	for rows.Next() {
		var item entity.GetDiscountQuery
		if err := rows.Scan(
			&item.ID,
			&item.NumberOfParticipant,
			&item.Amount,
		); err != nil {
			return customerrors.NewError(
				"failed to get discount",
				err,
				customerrors.DatabaseExecutionError,
			)
		}
		*discounts = append(*discounts, item)
	}
	return nil
}

func (dr *DiscountRepositoryImpl) DeleteDiscount(ctx context.Context, id int) error {
	var driver RepoDriver = dr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	UPDATE discounts
	SET
		deleted_at = CURRENT_TIMESTAMP,
		updated_at = CURRENT_TIMESTAMP
	WHERE id = $1 AND deleted_at IS NULL
	`
	_, err := driver.Exec(query, id)
	if err != nil {
		return customerrors.NewError(
			"failed to delete discount amount",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (dr *DiscountRepositoryImpl) UpdateAmount(ctx context.Context, id int, amount *float64) error {
	var driver RepoDriver = dr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	UPDATE discounts
	SET
		amount = COALESCE($1, amount),
		updated_at = CURRENT_TIMESTAMP
	WHERE id = $2 AND deleted_at IS NULL
	`
	_, err := driver.Exec(query, amount, id)
	if err != nil {
		return customerrors.NewError(
			"failed to update discount amount",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (dr *DiscountRepositoryImpl) FindByID(ctx context.Context, id int, discount *entity.Discount) error {
	var driver RepoDriver = dr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	SELECT
		id,
		number_of_participant,
		amount,
		created_at,
		updated_at,
		deleted_at
	FROM discounts
	WHERE id = $1 AND deleted_at IS NULL
	`
	if err := driver.QueryRow(query, id).Scan(
		&discount.ID,
		&discount.NumberOfParticipant,
		&discount.Amount,
		&discount.CreatedAt,
		&discount.UpdatedAt,
		&discount.DeletedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return customerrors.NewError(
				"discount not found",
				err,
				customerrors.ItemNotExist,
			)
		}
		return customerrors.NewError(
			"failed to get discount",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (dr *DiscountRepositoryImpl) CreateNewDiscount(ctx context.Context, numberOfParticipant int, amount float64, id *int) error {
	var driver RepoDriver = dr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	INSERT INTO discounts (number_of_participant, amount)
	VALUES
	($1, $2)
	ON CONFLICT (number_of_participant)
	DO UPDATE SET
		number_of_participant = EXCLUDED.number_of_participant,
		deleted_at = NULL,
		updated_at = CURRENT_TIMESTAMP
	RETURNING id
	`
	if err := driver.QueryRow(query, numberOfParticipant, amount).Scan(
		id,
	); err != nil {
		return customerrors.NewError(
			"failed to create new discount",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (dr *DiscountRepositoryImpl) GetMaxParticipant(ctx context.Context, maxParticipant *int) error {
	var driver RepoDriver = dr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	SELECT
		COALESCE(MAX(number_of_participant),1)
	FROM discounts
	WHERE deleted_at IS NULL
	`
	if err := driver.QueryRow(query).Scan(maxParticipant); err != nil {
		return customerrors.NewError(
			"failed to get discount",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (dr *DiscountRepositoryImpl) GetDiscountByNumberOfParticipant(ctx context.Context, numberOfParticipant int, discount *entity.Discount) error {
	var driver RepoDriver = dr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	SELECT
		id,
		number_of_participant,
		amount,
		created_at,
		updated_at,
		deleted_at
	FROM discounts
	WHERE number_of_participant = $1 AND deleted_at IS NULL
	`
	if err := driver.QueryRow(query, numberOfParticipant).Scan(
		&discount.ID,
		&discount.NumberOfParticipant,
		&discount.Amount,
		&discount.CreatedAt,
		&discount.UpdatedAt,
		&discount.DeletedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return customerrors.NewError(
				"discount not found",
				err,
				customerrors.ItemNotExist,
			)
		}
		return customerrors.NewError(
			"failed to get discount",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}
