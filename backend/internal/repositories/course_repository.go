package repositories

import (
	"context"
	"database/sql"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/entity"
)

type CourseRepositoryImpl struct {
	DB *sql.DB
}

func CreateCourseRepository(db *sql.DB) *CourseRepositoryImpl {
	return &CourseRepositoryImpl{db}
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
	course *entity.Course,
) error {
	var driver RepoDriver
	driver = cr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	INSERT INTO courses (mentor_id, title, description, domicile, min_price, max_price, min_duration_days, max_duration_days)
	VALUES
	($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING id
	`
	row := driver.QueryRow(query, mentorID, title, description, domicile, minPrice, maxPrice, minDuration, maxDuration)
	if err := row.Scan(&course.ID); err != nil {
		return customerrors.NewError(
			"failed to create course",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}
