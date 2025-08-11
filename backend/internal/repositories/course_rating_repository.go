package repositories

import (
	"context"
	"database/sql"
	"errors"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/entity"
)

type CourseRatingRepositoryImpl struct {
	DB *sql.DB
}

func CreateCourseRatingRepository(db *sql.DB) *CourseRatingRepositoryImpl {
	return &CourseRatingRepositoryImpl{db}
}

func (cr *CourseRatingRepositoryImpl) FindByCourseIDAndStudentID(
	ctx context.Context,
	courseID int,
	studentID string,
	review *entity.CourseRating,
) error {
	var driver RepoDriver = cr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	SELECT
		id,
		course_id,
		student_id,
		rating,
		feedback,
		created_at,
		updated_at,
		deleted_at
	FROM course_ratings
	WHERE course_id = $1 AND student_id = $2 AND deleted_at IS NULL
	`
	row := driver.QueryRow(query, courseID, studentID)
	if err := row.Scan(
		&review.ID,
		&review.CourseID,
		&review.StudentID,
		&review.Rating,
		&review.Feedback,
		&review.CreatedAt,
		&review.UpdatedAt,
		&review.DeletedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return customerrors.NewError(
				"rating does not exist",
				err,
				customerrors.ItemNotExist,
			)
		}
		return customerrors.NewError(
			"failed to get course review",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (cr *CourseRatingRepositoryImpl) CreateReview(
	ctx context.Context,
	courseID int,
	studentID string,
	rating int,
	feedback *string,
	review *entity.CourseRating,
) error {
	var driver RepoDriver = cr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}

	query := `
	INSERT INTO course_ratings(course_id, student_id, rating, feedback)
	VALUES
	($1, $2, $3, $4)
	RETURNING id
	`
	row := driver.QueryRow(
		query,
		courseID,
		studentID,
		rating,
		feedback,
	)
	if err := row.Scan(&review.ID); err != nil {
		return customerrors.NewError(
			"failed to create review",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}
