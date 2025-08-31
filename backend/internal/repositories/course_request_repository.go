package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

func (cr *CourseRequestRepositoryImpl) StudentCourseRequestList(
	ctx context.Context,
	studentID string,
	status *string,
	search *string,
	lastID int,
	limit int,
	totalRow *int64,
	requests *[]entity.StudentCourseRequestQuery,
) error {
	var driver RepoDriver = cr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	args := []any{lastID, studentID}
	countArgs := []any{lastID, studentID}
	query := `
	SELECT
		cr.id,
		cr.student_id,
		cr.course_id,
		cr.total_price,
		cr.status,
		u.name,
		u.email,
		c.title
	FROM course_requests cr
	JOIN courses c on cr.course_id = c.id
	JOIN mentors m on c.mentor_id = m.id
	JOIN users u on u.id = m.id
	WHERE
		cr.deleted_at IS NULL
		AND c.deleted_at IS NULL
		AND m.deleted_at IS NULL
		AND u.deleted_at IS NULL
		AND cr.id < $1
		AND cr.student_id = $2
	`
	countQuery := `
	SELECT
		count(*)
	FROM course_requests cr
	JOIN courses c on cr.course_id = c.id
	JOIN mentors m on c.mentor_id = m.id
	JOIN users u on u.id = m.id
	WHERE
		cr.deleted_at IS NULL
		AND c.deleted_at IS NULL
		AND m.deleted_at IS NULL
		AND u.deleted_at IS NULL
		AND cr.id < $1
		AND cr.student_id = $2
	`
	if status != nil {
		args = append(args, *status)
		countArgs = append(countArgs, *status)
		query += fmt.Sprintf(`
			AND cr.status = $%d
		`, len(args))
		countQuery += fmt.Sprintf(`
			AND cr.status = $%d
		`, len(countArgs))
	}
	if search != nil {
		args = append(args, "%"+*search+"%")
		countArgs = append(countArgs, "%"+*search+"%")
		query += fmt.Sprintf(`
		AND c.title ILIKE $%d
		`, len(args))
		countQuery += fmt.Sprintf(`
		AND c.title ILIKE $%d
		`, len(countArgs))
	}

	query += `ORDER BY cr.id DESC `
	args = append(args, limit)
	query += fmt.Sprintf(`LIMIT $%d`, len(args))

	if err := driver.QueryRow(countQuery, countArgs...).Scan(totalRow); err != nil {
		return customerrors.NewError(
			"failed to get mentor course request list",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	log.Println(query)
	rows, err := driver.Query(query, args...)
	if err != nil {
		return customerrors.NewError(
			"failed to get mentor course request list",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	defer rows.Close()

	for rows.Next() {
		var item entity.StudentCourseRequestQuery
		if err := rows.Scan(
			&item.ID,
			&item.StudentID,
			&item.CourseID,
			&item.TotalPrice,
			&item.Status,
			&item.MentorName,
			&item.MentorEmail,
			&item.CourseName,
		); err != nil {
			return customerrors.NewError(
				"failed to get mentor course request list",
				err,
				customerrors.DatabaseExecutionError,
			)
		}
		*requests = append(*requests, item)
	}
	return nil
}

func (cr *CourseRequestRepositoryImpl) MentorCourseRequestList(
	ctx context.Context,
	mentorID string,
	status *string,
	lastID int,
	limit int,
	totalRow *int64,
	requests *[]entity.MentorCourseRequestQuery,
) error {
	var driver RepoDriver = cr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	args := []any{lastID, mentorID}
	countArgs := []any{lastID, mentorID}
	query := `
	SELECT
		cr.id,
		cr.student_id,
		cr.course_id,
		cr.total_price,
		cr.status,
		u.name,
		u.email,
		c.title
	FROM course_requests cr
	JOIN courses c on cr.course_id = c.id
	JOIN students s on cr.student_id = s.id
	JOIN users u on u.id = s.id
	WHERE
		cr.deleted_at IS NULL
		AND c.deleted_at IS NULL
		AND s.deleted_at IS NULL
		AND u.deleted_at IS NULL
		AND cr.id < $1
		AND c.mentor_id = $2
	`
	countQuery := `
	SELECT
		count(*)
	FROM course_requests cr
	JOIN courses c on cr.course_id = c.id
	JOIN students s on cr.student_id = s.id
	JOIN users u on u.id = s.id
	WHERE
		cr.deleted_at IS NULL
		AND c.deleted_at IS NULL
		AND s.deleted_at IS NULL
		AND u.deleted_at IS NULL
		AND cr.id < $1
		AND c.mentor_id = $2
	`
	if status != nil {
		args = append(args, *status)
		countArgs = append(countArgs, *status)
		query += fmt.Sprintf(`
			AND cr.status = $%d
		`, len(args))
		countQuery += fmt.Sprintf(`
			AND cr.status = $%d
		`, len(countArgs))
	}

	query += `ORDER BY cr.id DESC `
	args = append(args, limit)
	query += fmt.Sprintf(`LIMIT $%d`, len(args))

	if err := driver.QueryRow(countQuery, countArgs...).Scan(totalRow); err != nil {
		return customerrors.NewError(
			"failed to get mentor course request list",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	log.Println(query)
	rows, err := driver.Query(query, args...)
	if err != nil {
		return customerrors.NewError(
			"failed to get mentor course request list",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	defer rows.Close()

	for rows.Next() {
		var item entity.MentorCourseRequestQuery
		if err := rows.Scan(
			&item.ID,
			&item.StudentID,
			&item.CourseID,
			&item.TotalPrice,
			&item.Status,
			&item.Name,
			&item.Email,
			&item.CourseName,
		); err != nil {
			return customerrors.NewError(
				"failed to get mentor course request list",
				err,
				customerrors.DatabaseExecutionError,
			)
		}
		*requests = append(*requests, item)
	}
	return nil
}

func (cr *CourseRequestRepositoryImpl) GetPaymentDetail(ctx context.Context, id int, cre *entity.CourseRequest) error {
	var driver RepoDriver = cr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	SELECT
		student_id,
		course_id,
		status,
		subtotal,
		operational_cost,
		total_price,
		expired_at
	FROM course_requests
	WHERE id = $1 AND deleted_at IS NULL
	`
	if err := driver.QueryRow(query, id).Scan(
		&cre.StudentID,
		&cre.CourseID,
		&cre.Status,
		&cre.SubTotal,
		&cre.OperationalCost,
		&cre.TotalPrice,
		&cre.ExpiredAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return customerrors.NewError(
				"no course request found",
				err,
				customerrors.ItemNotExist,
			)
		}
		return customerrors.NewError(
			"failed to get course request",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (cr *CourseRequestRepositoryImpl) CompleteRequest(ctx context.Context) error {
	var driver RepoDriver = cr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	UPDATE course_requests cr
	SET
		status = 'completed',
		updated_at = NOW()
	WHERE cr.status = 'scheduled'
	AND cr.deleted_at IS NULL
	AND NOT EXISTS (
		SELECT 1
		FROM course_schedule cs
		WHERE cs.course_request_id = cr.id
		AND cs.status != 'completed'
		AND cs.deleted_at IS NULL
	)
	`
	log.Println(query)
	_, err := driver.Exec(query)
	return err
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
	WHERE NOW() >= expired_at 
	AND deleted_at IS NULL
	AND status IN ('reserved','pending payment')
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
		subtotal,
		operational_cost,
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
		&courseRequest.SubTotal,
		&courseRequest.OperationalCost,
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
		deleted_at IS NULL
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
	INSERT INTO course_requests (student_id, course_id, subtotal, operational_cost, total_price, number_of_sessions, expired_at)
	VALUES
	($1, $2, $3, $4, $5, $6, $7)
	RETURNING (id)
	`
	row := driver.QueryRow(
		query,
		courseRequest.StudentID,
		courseRequest.CourseID,
		courseRequest.SubTotal,
		courseRequest.OperationalCost,
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
	WHERE course_id = $1 AND status NOT IN ('completed', 'cancelled') AND deleted_at IS NULL
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
