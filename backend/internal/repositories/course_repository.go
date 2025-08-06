package repositories

import (
	"context"
	"database/sql"
	"errors"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/entity"
)

type CourseRepositoryImpl struct {
	DB *sql.DB
}

func CreateCourseRepository(db *sql.DB) *CourseRepositoryImpl {
	return &CourseRepositoryImpl{db}
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
