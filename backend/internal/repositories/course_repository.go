package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/entity"
)

type CourseRepositoryImpl struct {
	DB *sql.DB
}

func CreateCourseRepository(db *sql.DB) *CourseRepositoryImpl {
	return &CourseRepositoryImpl{db}
}

func (cr *CourseRepositoryImpl) MentorListCourse(
	ctx context.Context,
	query *[]entity.MentorListCourseQuery,
	totalRow *int64,
	param entity.MentorListCourseParam,
) error {
	var driver RepoDriver
	driver = cr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	countArgs := []any{param.MentorID, param.LastID}
	args := []any{param.MentorID, param.LastID}
	sqlQuery := `
		SELECT
		c.id,
		c.title,
		c.domicile,
		c.min_price,
		c.max_price,
		c.method,
		c.min_duration_days,
		c.max_duration_days,
		COALESCE(
			(SELECT STRING_AGG(cc_all.name, ',') 
			FROM course_category_assignments cca_all 
			JOIN course_categories cc_all ON cca_all.category_id = cc_all.id 
			WHERE cca_all.course_id = c.id 
			AND cca_all.deleted_at IS NULL 
			AND cc_all.deleted_at IS NULL), 
			''
		) AS categories
	FROM courses c
	WHERE c.mentor_id = $1 AND c.deleted_at IS NULL AND c.id < $2
	`
	countQuery := `
	SELECT count(*)
	FROM courses c
	WHERE c.mentor_id = $1 AND c.deleted_at IS NULL AND c.id < $2
	`
	if param.Search != nil {
		countQuery += fmt.Sprintf(" AND c.title ILIKE $%d ", len(countArgs)+1)
		sqlQuery += fmt.Sprintf(" AND c.title ILIKE $%d ", len(args)+1)
		args = append(args, "%"+*param.Search+"%")
		countArgs = append(countArgs, "%"+*param.Search+"%")
	}
	if param.CourseCategory != nil {
		countQuery += fmt.Sprintf(` AND EXISTS (
			SELECT 1 
			FROM course_category_assignments cca_filter 
			JOIN course_categories cc_filter ON cca_filter.category_id = cc_filter.id
			WHERE cca_filter.course_id = c.id 
			AND cca_filter.deleted_at IS NULL
			AND cc_filter.deleted_at IS NULL
			AND cc_filter.id = $%d
		)`, len(args)+1)
		sqlQuery += fmt.Sprintf(` AND EXISTS (
			SELECT 1 
			FROM course_category_assignments cca_filter 
			JOIN course_categories cc_filter ON cca_filter.category_id = cc_filter.id
			WHERE cca_filter.course_id = c.id 
			AND cca_filter.deleted_at IS NULL
			AND cc_filter.deleted_at IS NULL
			AND cc_filter.id = $%d
		)`, len(args)+1)
		countArgs = append(countArgs, *param.CourseCategory)
		args = append(args, *param.CourseCategory)
	}
	sqlQuery += " ORDER BY c.id DESC "
	sqlQuery += fmt.Sprintf(" LIMIT $%d ", len(args)+1)
	args = append(args, param.Limit)
	row := driver.QueryRow(countQuery, countArgs...)
	if err := row.Scan(totalRow); err != nil {
		return customerrors.NewError(
			"failed to get course list",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	rows, err := driver.Query(sqlQuery, args...)
	if err != nil {
		return customerrors.NewError(
			"failed to get courses list",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	defer rows.Close()

	for rows.Next() {
		var item entity.MentorListCourseQuery
		if err := rows.Scan(
			&item.ID,
			&item.Title,
			&item.Domicile,
			&item.MinPrice,
			&item.MaxPrice,
			&item.Method,
			&item.MinDurationDays,
			&item.MaxDurationDays,
			&item.CourseCategories,
		); err != nil {
			return customerrors.NewError(
				"failed to get course list",
				err,
				customerrors.DatabaseExecutionError,
			)
		}
		*query = append(*query, item)
	}

	return nil
}

func (cr *CourseRepositoryImpl) FindByID(ctx context.Context, id int, course *entity.Course) error {
	var driver RepoDriver
	driver = cr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	SELECT
		id,
		mentor_id,
		title,
		description,
		domicile,
		min_price,
		method,
		max_price,
		min_duration_days,
		max_duration_days,
		transaction_count,
		created_at,
		updated_at,
		deleted_at
	FROM courses
	WHERE id = $1 and deleted_at IS NULL
	`
	row := driver.QueryRow(query, id)
	if err := row.Scan(
		&course.ID,
		&course.MentorID,
		&course.Title,
		&course.Description,
		&course.Domicile,
		&course.MinPrice,
		&course.Method,
		&course.MaxPrice,
		&course.MinDuration,
		&course.MaxDuration,
		&course.TransactionCount,
		&course.CreatedAt,
		&course.UpdatedAt,
		&course.DeletedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return customerrors.NewError(
				"course does not exist",
				err,
				customerrors.ItemNotExist,
			)
		}
		return customerrors.NewError(
			"failed to get course",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (cr *CourseRepositoryImpl) DeleteCourse(ctx context.Context, id int) error {
	var driver RepoDriver
	driver = cr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	UPDATE courses
	SET deleted_at = NOW(), updated_at = NOW()
	WHERE id = $1 AND deleted_at IS NULL
	`
	_, err := driver.Exec(query, id)
	if err != nil {
		return customerrors.NewError(
			"failed to delete course",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (cr *CourseRepositoryImpl) CreateCourse(
	ctx context.Context,
	mentorID string,
	title string,
	description string,
	domicile string,
	minPrice float64,
	maxPrice float64,
	minDuration int,
	maxDuration int,
	method string,
	course *entity.Course,
) error {
	var driver RepoDriver
	driver = cr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	INSERT INTO courses (mentor_id, title, description, domicile, min_price, max_price, min_duration_days, max_duration_days, method)
	VALUES
	($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING id
	`
	row := driver.QueryRow(query, mentorID, title, description, domicile, minPrice, maxPrice, minDuration, maxDuration, method)
	if err := row.Scan(&course.ID); err != nil {
		return customerrors.NewError(
			"failed to create course",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}
