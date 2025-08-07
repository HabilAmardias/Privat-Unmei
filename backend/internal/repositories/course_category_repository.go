package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/entity"

	"github.com/lib/pq"
)

type CourseCategoryRepositoryImpl struct {
	DB *sql.DB
}

func CreateCourseCategoryRepository(db *sql.DB) *CourseCategoryRepositoryImpl {
	return &CourseCategoryRepositoryImpl{db}
}

func (ccr *CourseCategoryRepositoryImpl) UnassignCategories(ctx context.Context, courseID int) error {
	var driver RepoDriver
	driver = ccr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	UPDATE course_category_assignments
	SET updated_at = NOW(), deleted_at = NOW()
	WHERE course_id = $1 AND deleted_at IS NULL
	`
	_, err := driver.Exec(query, courseID)
	if err != nil {
		return customerrors.NewError(
			"failed to unassign course categories",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (ccr *CourseCategoryRepositoryImpl) FindByMultipleIDs(ctx context.Context, ids []int, categories *[]entity.CourseCategory) error {
	var driver RepoDriver
	driver = ccr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	args := []any{pq.Array(ids)}
	query := `
	SELECT id, name, created_at, updated_at, deleted_at
	FROM course_categories
	WHERE id = ANY($1) AND deleted_at IS NULL
	`
	log.Println(query)
	rows, err := driver.Query(query, args...)
	if err != nil {
		return customerrors.NewError(
			"failed to find course category",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	defer rows.Close()
	for rows.Next() {
		var item entity.CourseCategory
		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.DeletedAt,
		); err != nil {
			return customerrors.NewError(
				"failed to find course category",
				err,
				customerrors.DatabaseExecutionError,
			)
		}
		*categories = append(*categories, item)
	}
	return nil
}

func (ccr *CourseCategoryRepositoryImpl) AssignCategories(ctx context.Context, courseID int, categories []int) error {
	var driver RepoDriver
	driver = ccr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	args := []any{courseID}
	query := `
	INSERT INTO course_category_assignments (course_id, category_id)
	VALUES
	`
	for i, cat := range categories {
		if i != len(categories)-1 {
			query += fmt.Sprintf(`
			($1, $%d),
			`, len(args)+1)
		} else {
			query += fmt.Sprintf(`
			($1, $%d);
			`, len(args)+1)
		}
		args = append(args, cat)
	}
	log.Println(query)
	_, err := driver.Exec(query, args...)
	if err != nil {
		return customerrors.NewError(
			"failed to assign course category",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (ccr *CourseCategoryRepositoryImpl) FindByID(ctx context.Context, id int, category *entity.CourseCategory) error {
	var driver RepoDriver
	driver = ccr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	SELECT id, name, created_at, updated_at, deleted_at
	FROM course_categories
	WHERE id = $1 AND deleted_at IS NULL
	`
	row := driver.QueryRow(query, id)
	if err := row.Scan(
		&category.ID,
		&category.Name,
		&category.CreatedAt,
		&category.UpdatedAt,
		&category.DeletedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return customerrors.NewError(
				"category does not exist",
				err,
				customerrors.ItemNotExist,
			)
		}
		return customerrors.NewError(
			"failed to get category",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (ccr *CourseCategoryRepositoryImpl) FindByName(ctx context.Context, name string, category *entity.CourseCategory) error {
	var driver RepoDriver
	driver = ccr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	SELECT id, name, created_at, updated_at, deleted_at
	FROM course_categories
	WHERE LOWER(name) = LOWER($1) AND deleted_at IS NULL
	`
	row := driver.QueryRow(query, name)
	if err := row.Scan(
		&category.ID,
		&category.Name,
		&category.CreatedAt,
		&category.UpdatedAt,
		&category.DeletedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return customerrors.NewError(
				"category does not exist",
				err,
				customerrors.ItemNotExist,
			)
		}
		return customerrors.NewError(
			"failed to get category",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (ccr *CourseCategoryRepositoryImpl) UpdateCategory(ctx context.Context, param entity.UpdateCategoryParam) error {
	var driver RepoDriver
	driver = ccr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	UPDATE course_categories
	SET name = COALESCE($1, name), updated_at = NOW()
	WHERE id = $2 AND deleted_at IS NULL
	`
	_, err := driver.Exec(query, param.Name, param.ID)
	if err != nil {
		return customerrors.NewError(
			"failed to update category",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (ccr *CourseCategoryRepositoryImpl) CreateCategory(ctx context.Context, category *entity.CreateCategoryQuery) error {
	var driver RepoDriver
	driver = ccr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	INSERT INTO course_categories (name)
	VALUES
	($1)
	RETURNING id, name
	`
	row := driver.QueryRow(query, category.Name)
	if err := row.Scan(
		&category.ID,
		&category.Name,
	); err != nil {
		return customerrors.NewError(
			"failed to create categories",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
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
