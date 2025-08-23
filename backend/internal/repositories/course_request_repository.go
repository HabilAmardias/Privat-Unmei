package repositories

import (
	"context"
	"database/sql"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/entity"
)

type CourseRequestRepositoryImpl struct {
	DB *sql.DB
}

func CreateCourseRequestRepository(db *sql.DB) *CourseRequestRepositoryImpl {
	return &CourseRequestRepositoryImpl{db}
}

func (cr *CourseRequestRepositoryImpl) FindOngoingByCourseIDAndStudentID(ctx context.Context, courseID int, studentID string, count *int64) error {
	var driver RepoDriver = cr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	SELECT COUNT(*)
	FROM course_requests
	WHERE
		student_id = $1 AND
		course_id = $2 AND
		status IN ('reserved','pending payment', 'scheduled') AND
		deleted_at IS NULL AND
		expired_at > NOW()
	`
	if err := driver.QueryRow(query, studentID, courseID).Scan(count); err != nil {
		return customerrors.NewError(
			"failed to get order",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (cr *CourseRequestRepositoryImpl) CreateOrder(
	ctx context.Context,
	courseRequest *entity.CourseRequest,
) error {
	var driver RepoDriver = cr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	INSERT INTO course_requests (student_id, course_id, total_price, number_of_sessions, expired_at)
	VALUES
	($1, $2, $3, $4, $5)
	RETURNING (id)
	`
	row := driver.QueryRow(
		query,
		courseRequest.StudentID,
		courseRequest.CourseID,
		courseRequest.TotalPrice,
		courseRequest.NumberOfSessions,
		courseRequest.ExpiredAt,
	)
	if err := row.Scan(&courseRequest.ID); err != nil {
		return customerrors.NewError(
			"failed to create order",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (cr *CourseRequestRepositoryImpl) FindOngoingByCourseID(ctx context.Context, courseID int, orders *[]entity.CourseRequest) error {
	var driver RepoDriver = cr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	SELECT
		id,
		student_id,
		course_id,
		status,
		total_price,
		number_of_sessions,
		accepted_at,
		payment_due,
		created_at,
		updated_at,
		deleted_at
	FROM course_requests
	WHERE course_id = $1 AND status NOT IN ('completed', 'cancelled') AND deleted_at IS NULL AND expired_at > NOW()
	`
	rows, err := driver.Query(query, courseID)
	if err != nil {
		return customerrors.NewError(
			"failed to get order",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	defer rows.Close()
	for rows.Next() {
		var item entity.CourseRequest
		if err := rows.Scan(
			&item.ID,
			&item.StudentID,
			&item.CourseID,
			&item.Status,
			&item.TotalPrice,
			&item.NumberOfSessions,
			&item.AcceptedAt,
			&item.PaymentDue,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.DeletedAt,
		); err != nil {
			return customerrors.NewError(
				"failed to get orders",
				err,
				customerrors.DatabaseExecutionError,
			)
		}
		*orders = append(*orders, item)
	}
	return nil
}

func (cr *CourseRequestRepositoryImpl) FindCompletedByStudentIDAndCourseID(ctx context.Context, studentID string, courseID int, orders *[]entity.CourseRequest) error {
	var driver RepoDriver = cr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	SELECT
		id,
		student_id,
		course_id,
		status,
		total_price,
		number_of_sessions,
		accepted_at,
		payment_due,
		created_at,
		updated_at,
		deleted_at
	FROM course_requests
	WHERE student_id = $1 AND course_id = $2 AND status = 'completed' AND deleted_at IS NULL
	`
	rows, err := driver.Query(query, studentID, courseID)
	if err != nil {
		return customerrors.NewError(
			"failed to get order",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	defer rows.Close()
	for rows.Next() {
		var item entity.CourseRequest
		if err := rows.Scan(
			&item.ID,
			&item.StudentID,
			&item.CourseID,
			&item.Status,
			&item.TotalPrice,
			&item.NumberOfSessions,
			&item.AcceptedAt,
			&item.PaymentDue,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.DeletedAt,
		); err != nil {
			return customerrors.NewError(
				"failed to get orders",
				err,
				customerrors.DatabaseExecutionError,
			)
		}
		*orders = append(*orders, item)
	}
	return nil
}
