package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/entity"
)

type CourseCategoryRepositoryImpl struct {
	DB *sql.DB
}

func CreateCourseCategoryRepository(db *sql.DB) *CourseCategoryRepositoryImpl {
	return &CourseCategoryRepositoryImpl{db}
}

func (ccr *CourseCategoryRepositoryImpl) GetCourseCategoryList(
	ctx context.Context,
	categories *[]entity.ListCourseCategoryQuery,
	totalRow *int64,
	param entity.ListCourseCategoryParam,
) error {
	var driver RepoDriver
	driver = ccr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	args := []any{param.LastID}
	countArgs := []any{param.LastID}
	query := `
	SELECT id, name
	FROM course_categories
	WHERE id < $1 AND deleted_at IS NULL
	`
	countQuery := `
	SELECT count(*)
	FROM course_categories
	WHERE id < $1 AND deleted_at IS NULL
	`
	if param.Search != nil {
		query += fmt.Sprintf(" AND name ILIKE $%d ", len(args)+1)
		countQuery += fmt.Sprintf(" AND name ILIKE $%d ", len(countArgs)+1)
		args = append(args, "%"+*param.Search+"%")
		countArgs = append(countArgs, "%"+*param.Search+"%")
	}
	query += " ORDER BY id DESC "
	query += fmt.Sprintf(" LIMIT $%d ", len(args)+1)
	args = append(args, param.Limit)

	row := driver.QueryRow(countQuery, countArgs...)
	if err := row.Scan(totalRow); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return customerrors.NewError(
				"no categories found",
				err,
				customerrors.ItemNotExist,
			)
		}
		return customerrors.NewError(
			"failed to get categories",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	log.Println(query)
	rows, err := driver.Query(query, args...)
	if err != nil {
		return customerrors.NewError(
			"failed to get categories",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	defer rows.Close()

	for rows.Next() {
		var cat entity.ListCourseCategoryQuery
		if err := rows.Scan(
			&cat.ID,
			&cat.Name,
		); err != nil {
			return customerrors.NewError(
				"failed to get categories",
				err,
				customerrors.DatabaseExecutionError,
			)
		}
		*categories = append(*categories, cat)
	}
	return nil
}
