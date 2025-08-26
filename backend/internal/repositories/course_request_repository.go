package repositories

import (
	"context"
	"database/sql"
	"log"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/entity"
	"time"
)

type CourseRequestRepositoryImpl struct {
	DB *sql.DB
}

func CreateCourseRequestRepository(db *sql.DB) *CourseRequestRepositoryImpl {
	return &CourseRequestRepositoryImpl{db}
}

func (cr *CourseRequestRepositoryImpl) CancelExpiredRequest(ctx context.Context, courseRequestIDs *[]int) error {
	var driver RepoDriver = cr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	UPDATE course_requests
	SET
		status = 'cancelled',
		updated_at = NOW()
	WHERE NOW() >= expired_at AND deleted_at IS NULL
	RETURNING (id)
	`
	log.Println(query)
	rows, err := driver.Query(query)
	if err != nil {
		// not wrapping it on customerror because its just for cron
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var item int
		if err := rows.Scan(&item); err != nil {
			// not wrapping it on customerror because its just for cron
			return err
		}
		*courseRequestIDs = append(*courseRequestIDs, item)
	}
	return nil
}

func (cr *CourseRequestRepositoryImpl) ChangeRequestStatus(ctx context.Context, id int, newStatus string, eat *time.Time) error {
	var driver RepoDriver = cr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	UPDATE course_requests
	SET
		status = $1,
		updated_at = NOW(),
		expired_at = $2
	WHERE id = $3 AND deleted_at IS NULL
	`
	_, err := driver.Exec(query, newStatus, eat, id)
	if err != nil {
		return customerrors.NewError(
			"failed to update request status",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (cr *CourseRequestRepositoryImpl) FindByID(ctx context.Context, id int, courseRequest *entity.CourseRequest) error {
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
		expired_at,
		created_at,
		updated_at,
		deleted_at
	FROM course_requests
	WHERE id = $1 AND deleted_at IS NULL
	`
	if err := driver.QueryRow(query, id).Scan(
		&courseRequest.ID,
		&courseRequest.StudentID,
		&courseRequest.CourseID,
		&courseRequest.Status,
		&courseRequest.TotalPrice,
		&courseRequest.NumberOfSessions,
		&courseRequest.ExpiredAt,
		&courseRequest.CreatedAt,
		&courseRequest.UpdatedAt,
		&courseRequest.DeletedAt,
	); err != nil {
		return customerrors.NewError(
			"failed to get request",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
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
		expired_at,
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
			&item.ExpiredAt,
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
		expired_at,
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
			&item.ExpiredAt,
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
