package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/db"
	"privat-unmei/internal/entity"

	"github.com/lib/pq"
)

type CourseCategoryRepositoryImpl struct {
	DB *db.CustomDB
}

func CreateCourseCategoryRepository(db *db.CustomDB) *CourseCategoryRepositoryImpl {
	return &CourseCategoryRepositoryImpl{db}
}

func (ccr *CourseCategoryRepositoryImpl) DeleteCategory(ctx context.Context, categoryID int) error {
	var driver RepoDriver
	driver = ccr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	UPDATE course_categories
	SET
		deleted_at = CURRENT_TIMESTAMP,
		updated_at = CURRENT_TIMESTAMP
	WHERE id = $1 AND deleted_at IS NULL
	`
	_, err := driver.Exec(query, categoryID)
	if err != nil {
		return customerrors.NewError(
			"failed to delete category",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (ccr *CourseCategoryRepositoryImpl) UnassignCategoriesByCategoryID(ctx context.Context, categoryID int) error {
	var driver RepoDriver
	driver = ccr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	UPDATE course_category_assignments
	SET
		deleted_at = CURRENT_TIMESTAMP,
		updated_at = CURRENT_TIMESTAMP
	WHERE category_id = $1 AND deleted_at IS NULL
	`
	_, err := driver.Exec(query, categoryID)
	if err != nil {
		return customerrors.NewError(
			"failed to delete category",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (ccr *CourseCategoryRepositoryImpl) GetCategoriesByCourseID(ctx context.Context, courseID int, categories *[]entity.GetCategoriesQuery) error {
	var driver RepoDriver
	driver = ccr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}

	query := `
	SELECT
		cc.id,
		cc.name
	FROM course_categories cc
	JOIN course_category_assignments cca on cca.category_id = cc.id
	WHERE cca.course_id = $1 AND cca.deleted_at IS NULL AND cc.deleted_at IS NULL
	`

	rows, err := driver.Query(query, courseID)
	if err != nil {
		return customerrors.NewError(
			"failed to get course categories",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	defer rows.Close()

	for rows.Next() {
		var item entity.GetCategoriesQuery
		if err := rows.Scan(&item.CategoryID, &item.CategoryName); err != nil {
			return customerrors.NewError(
				"failed to get course categories",
				err,
				customerrors.DatabaseExecutionError,
			)
		}
		*categories = append(*categories, item)
	}

	return nil
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

func (ccr *CourseCategoryRepositoryImpl) UnassignCategoriesMultipleCourse(ctx context.Context, courseIDs []int) error {
	var driver RepoDriver
	driver = ccr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	UPDATE course_category_assignments
	SET updated_at = NOW(), deleted_at = NOW()
	WHERE course_id = ANY($1) AND deleted_at IS NULL
	`
	_, err := driver.Exec(query, pq.Array(courseIDs))
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
			($1, $%d)
			`, len(args)+1)
		}
		args = append(args, cat)
	}
	query += `
	ON CONFLICT(course_id, category_id)
	DO UPDATE SET
		course_id = EXCLUDED.course_id,
		category_id = EXCLUDED.category_id,
		deleted_at = NULL,
		updated_at = NOW();
	`

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
	page int,
	limit int,
	search *string,
) error {
	var driver RepoDriver
	driver = ccr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	args := []any{}
	countArgs := []any{}
	query := `
	SELECT id, name
	FROM course_categories
	WHERE deleted_at IS NULL
	`
	countQuery := `
	SELECT count(*)
	FROM course_categories
	WHERE deleted_at IS NULL
	`
	if search != nil {
		query += fmt.Sprintf(" AND name ILIKE $%d ", len(args)+1)
		countQuery += fmt.Sprintf(" AND name ILIKE $%d ", len(countArgs)+1)
		args = append(args, "%"+*search+"%")
		countArgs = append(countArgs, "%"+*search+"%")
	}
	args = append(args, limit)
	query += fmt.Sprintf(" LIMIT $%d ", len(args))

	args = append(args, limit*(page-1))
	query += fmt.Sprintf(" OFFSET $%d ", len(args))

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
