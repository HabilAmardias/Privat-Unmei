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
