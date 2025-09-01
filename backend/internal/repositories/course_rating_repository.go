package repositories

import (
	"context"
	"database/sql"
	"errors"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/db"
	"privat-unmei/internal/entity"
)

type CourseRatingRepositoryImpl struct {
	DB *db.CustomDB
}

func CreateCourseRatingRepository(db *db.CustomDB) *CourseRatingRepositoryImpl {
	return &CourseRatingRepositoryImpl{db}
}

func (cr *CourseRatingRepositoryImpl) GetCourseReviews(ctx context.Context, courseID int, lastID int, limit int, totalRow *int64, review *[]entity.CourseRatingQuery) error {
	var driver RepoDriver = cr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	SELECT
		cr.id,
		cr.course_id,
		cr.student_id,
		u.name,
		cr.rating,
		cr.feedback,
		cr.created_at
	FROM course_ratings cr
	JOIN users u ON u.id = cr.student_id
	WHERE
		u.deleted_at IS NULL AND 
		cr.course_id = $1 AND 
		cr.deleted_at IS NULL AND
		cr.id < $2
	ORDER BY cr.id DESC
	LIMIT $3
	`
	countQuery := `
	SELECT count(*)
	FROM course_ratings cr
	JOIN users u ON u.id = cr.student_id
	WHERE
		u.deleted_at IS NULL AND 
		cr.course_id = $1 AND 
		cr.deleted_at IS NULL AND
		cr.id < $2
	`
	row := driver.QueryRow(countQuery, courseID, lastID)
	if err := row.Scan(totalRow); err != nil {
		return customerrors.NewError(
			"failed to get course reviews",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	rows, err := driver.Query(query, courseID, lastID, limit)
	if err != nil {
		return customerrors.NewError(
			"failed to get course reviews",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	defer rows.Close()
	for rows.Next() {
		var item entity.CourseRatingQuery
		if err := rows.Scan(
			&item.ID,
			&item.CourseID,
			&item.StudentID,
			&item.StudentName,
			&item.Rating,
			&item.Feedback,
			&item.CreatedAt,
		); err != nil {
			return customerrors.NewError(
				"failed to get course reviews",
				err,
				customerrors.DatabaseExecutionError,
			)
		}
		*review = append(*review, item)
	}
	return nil
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
