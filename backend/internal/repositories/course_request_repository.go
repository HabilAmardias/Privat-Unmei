package repositories

import (
	"context"
	"fmt"

	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/db"
	"privat-unmei/internal/entity"
	"time"
)

type CourseRequestRepositoryImpl struct {
	DB *db.CustomDB
}

func CreateCourseRequestRepository(db *db.CustomDB) *CourseRequestRepositoryImpl {
	return &CourseRequestRepositoryImpl{db}
}

func (cr *CourseRequestRepositoryImpl) GetMonthlyCostReport(ctx context.Context, reports *[]entity.MonthlyCostReportQuery) error {
	var driver RepoDriver = cr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	SELECT
		COALESCE(SUM(p.operational_cost), 0),
		EXTRACT(MONTH from p.created_at) as month_report
	FROM payments p
	JOIN course_requests cr ON p.course_request_id = cr.id
	WHERE
		p.deleted_at IS NULL AND
		EXTRACT(YEAR FROM p.created_at) = EXTRACT(YEAR FROM CURRENT_TIMESTAMP) AND
		cr.deleted_at IS NULL AND
		cr.status IN ('scheduled', 'completed')
	GROUP BY month_report
	ORDER BY month_report ASC
	`
	rows, err := driver.Query(query)
	if err != nil {
		return customerrors.NewError(
			"something went wrong",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	defer rows.Close()
	for rows.Next() {
		var item entity.MonthlyCostReportQuery
		if err := rows.Scan(
			&item.TotalCost,
			&item.Month,
		); err != nil {
			return customerrors.NewError(
				"something went wrong",
				err,
				customerrors.DatabaseExecutionError,
			)
		}
		*reports = append(*reports, item)
	}
	return nil
}

func (cr *CourseRequestRepositoryImpl) GetMonthlySessionReport(ctx context.Context, reports *[]entity.MonthlySessionReportQuery) error {
	var driver RepoDriver = cr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	SELECT 
		COALESCE(SUM(number_of_sessions), 0),
		EXTRACT(MONTH FROM created_at) as month_report
	FROM course_requests
	WHERE
		deleted_at IS NULL AND
		status IN ('scheduled', 'completed') AND
		EXTRACT(YEAR FROM created_at) = EXTRACT(YEAR FROM CURRENT_TIMESTAMP)
	GROUP BY
		month_report
	ORDER BY month_report ASC
	`
	rows, err := driver.Query(query)
	if err != nil {
		return customerrors.NewError(
			"something went wrong",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	defer rows.Close()
	for rows.Next() {
		var item entity.MonthlySessionReportQuery
		if err := rows.Scan(
			&item.TotalSession,
			&item.Month,
		); err != nil {
			return customerrors.NewError(
				"something went wrong",
				err,
				customerrors.DatabaseExecutionError,
			)
		}
		*reports = append(*reports, item)
	}
	return nil
}

func (cr *CourseRequestRepositoryImpl) GetThisMonthMentorReport(ctx context.Context, reports *[]entity.MonthlyMentorReportQuery) error {
	var driver RepoDriver = cr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	// TODO: Cache this with redis
	query := `
	SELECT
		u.name,
		u.email,
		COALESCE(SUM(cr.number_of_sessions), 0),
		COALESCE(SUM(p.operational_cost), 0)
	FROM payments p
	JOIN course_requests cr ON p.course_request_id = cr.id
	JOIN courses c ON cr.course_id = c.id
	JOIN mentors m ON c.mentor_id = m.id
	JOIN users u ON u.id = m.id
	WHERE
		p.deleted_at IS NULL AND
		EXTRACT(MONTH FROM p.created_at) = EXTRACT(MONTH FROM CURRENT_TIMESTAMP) AND
		EXTRACT(YEAR FROM p.created_at) = EXTRACT(YEAR FROM CURRENT_TIMESTAMP) AND
		cr.deleted_at IS NULL AND
		cr.status IN ('scheduled', 'completed') AND
		c.deleted_at IS NULL AND
		m.deleted_at IS NULL AND
		u.deleted_at IS NULL
	GROUP BY u.email, u.name
	`
	rows, err := driver.Query(query)
	if err != nil {
		return customerrors.NewError(
			"something went wrong",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	defer rows.Close()
	for rows.Next() {
		var item entity.MonthlyMentorReportQuery
		if err := rows.Scan(
			&item.MentorName,
			&item.MentorEmail,
			&item.TotalSession,
			&item.TotalCost,
		); err != nil {
			return customerrors.NewError(
				"something went wrong",
				err,
				customerrors.DatabaseExecutionError,
			)
		}
		*reports = append(*reports, item)
	}
	return nil
}

func (cr *CourseRequestRepositoryImpl) GetThisMonthSessions(ctx context.Context, totalSession *int64) error {
	var driver RepoDriver = cr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}

	query := `
	SELECT COALESCE(SUM(number_of_sessions), 0)
	FROM course_requests
	WHERE
		deleted_at IS NULL AND
		status IN ('scheduled', 'completed') AND
		EXTRACT(MONTH FROM created_at) = EXTRACT(MONTH FROM CURRENT_TIMESTAMP) AND
		EXTRACT(YEAR FROM created_at) = EXTRACT(YEAR FROM CURRENT_TIMESTAMP)
	`
	if err := driver.QueryRow(query).Scan(totalSession); err != nil {
		return customerrors.NewError(
			"something went wrong",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (cr *CourseRequestRepositoryImpl) GetThisMonthOperationalCost(ctx context.Context, totalCost *float64) error {
	var driver RepoDriver = cr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}

	query := `
	SELECT COALESCE(SUM(p.operational_cost), 0)
	FROM payments p
	JOIN course_requests cr ON p.course_request_id = cr.id
	WHERE
		p.deleted_at IS NULL AND
		cr.deleted_at IS NULL AND
		cr.status IN ('scheduled', 'completed') AND
		EXTRACT(MONTH FROM p.created_at) = EXTRACT(MONTH FROM CURRENT_TIMESTAMP) AND
		EXTRACT(YEAR FROM p.created_at) = EXTRACT(YEAR FROM CURRENT_TIMESTAMP)
	`
	if err := driver.QueryRow(query).Scan(totalCost); err != nil {
		return customerrors.NewError(
			"something went wrong",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (cr *CourseRequestRepositoryImpl) DeleteAllStudentOrders(ctx context.Context, id string) error {
	var driver RepoDriver = cr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	UPDATE course_requests
	SET deleted_at = CURRENT_TIMESTAMP
	WHERE student_id = $1 AND deleted_at IS NULL
	`
	_, err := driver.Exec(query, id)
	if err != nil {
		return customerrors.NewError(
			"something went wrong",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (cr *CourseRequestRepositoryImpl) DeleteAllMentorOrders(ctx context.Context, id string) error {
	var driver RepoDriver = cr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	UPDATE course_requests
	SET deleted_at = NOW(), updated_at = NOW()
	WHERE course_id IN (
		SELECT id
		FROM courses
		WHERE mentor_id = $1 AND deleted_at IS NULL
	) AND deleted_at IS NULL
	`
	_, err := driver.Exec(query, id)
	if err != nil {
		return customerrors.NewError(
			"failed to delete course requests",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (cr *CourseRequestRepositoryImpl) FindOngoingOrderByMentorID(ctx context.Context, mentorID string, count *int64) error {
	var driver RepoDriver = cr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	SELECT COUNT(*)
	FROM course_requests cr
	JOIN courses c ON cr.course_id = c.id
	JOIN mentors m ON c.mentor_id = m.id
	WHERE cr.status IN ('reserved','pending payment') AND
	cr.deleted_at IS NULL AND cr.expired_at > CURRENT_TIMESTAMP AND
	c.deleted_at IS NULL AND
	m.deleted_at IS NULL AND
	c.mentor_id = $1
	`
	if err := driver.QueryRow(query, mentorID).Scan(count); err != nil {
		return customerrors.NewError(
			"failed to get ongoing order",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (cr *CourseRequestRepositoryImpl) StudentCourseRequestList(
	ctx context.Context,
	studentID string,
	status *string,
	search *string,
	page int,
	limit int,
	totalRow *int64,
	requests *[]entity.StudentCourseRequestQuery,
) error {
	var driver RepoDriver = cr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	args := []any{studentID}
	countArgs := []any{studentID}
	// TODO: cache 1st page of this with redis
	query := `
	SELECT
		cr.id,
		cr.student_id,
		cr.course_id,
		py.total_price,
		cr.status,
		u.name,
		u.public_id,
		c.title
	FROM course_requests cr
	JOIN courses c on cr.course_id = c.id
	JOIN mentors m on c.mentor_id = m.id
	JOIN users u on u.id = m.id
	JOIN payments py on py.course_request_id = cr.id
	WHERE
		cr.deleted_at IS NULL
		AND c.deleted_at IS NULL
		AND m.deleted_at IS NULL
		AND u.deleted_at IS NULL
		AND cr.student_id = $1
	`
	countQuery := `
	SELECT
		count(*)
	FROM course_requests cr
	JOIN courses c on cr.course_id = c.id
	JOIN mentors m on c.mentor_id = m.id
	JOIN users u on u.id = m.id
	JOIN payments py on py.course_request_id = cr.id
	WHERE
		cr.deleted_at IS NULL
		AND c.deleted_at IS NULL
		AND m.deleted_at IS NULL
		AND u.deleted_at IS NULL
		AND cr.student_id = $1
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
	args = append(args, limit)
	query += fmt.Sprintf(`LIMIT $%d`, len(args))

	args = append(args, limit*(page-1))
	query += fmt.Sprintf(" OFFSET $%d", len(args))

	if err := driver.QueryRow(countQuery, countArgs...).Scan(totalRow); err != nil {
		return customerrors.NewError(
			"failed to get mentor course request list",
			err,
			customerrors.DatabaseExecutionError,
		)
	}

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
			&item.MentorPublicID,
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
	page int,
	limit int,
	totalRow *int64,
	requests *[]entity.MentorCourseRequestQuery,
) error {
	var driver RepoDriver = cr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	args := []any{mentorID}
	countArgs := []any{mentorID}
	// TODO: cache 1st page of this with redis
	query := `
	SELECT
		cr.id,
		cr.student_id,
		cr.course_id,
		p.total_price,
		cr.status,
		u.name,
		u.public_id,
		c.title
	FROM course_requests cr
	JOIN courses c on cr.course_id = c.id
	JOIN students s on cr.student_id = s.id
	JOIN users u on u.id = s.id
	JOIN payments p on p.course_request_id = cr.id
	WHERE
		cr.deleted_at IS NULL
		AND c.deleted_at IS NULL
		AND s.deleted_at IS NULL
		AND u.deleted_at IS NULL
		AND p.deleted_at IS NULL
		AND c.mentor_id = $1
	`
	countQuery := `
	SELECT
		count(*)
	FROM course_requests cr
	JOIN courses c on cr.course_id = c.id
	JOIN students s on cr.student_id = s.id
	JOIN users u on u.id = s.id
	JOIN payments p on p.course_request_id = cr.id
	WHERE
		cr.deleted_at IS NULL
		AND c.deleted_at IS NULL
		AND s.deleted_at IS NULL
		AND u.deleted_at IS NULL
		AND p.deleted_at IS NULL
		AND c.mentor_id = $1
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
	args = append(args, limit)
	query += fmt.Sprintf(`LIMIT $%d`, len(args))

	args = append(args, limit*(page-1))
	query += fmt.Sprintf(" OFFSET $%d", len(args))

	if err := driver.QueryRow(countQuery, countArgs...).Scan(totalRow); err != nil {
		return customerrors.NewError(
			"failed to get mentor course request list",
			err,
			customerrors.DatabaseExecutionError,
		)
	}

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
			&item.PublicID,
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

	_, err := driver.Exec(query)
	return err
}

func (cr *CourseRequestRepositoryImpl) CancelExpiredRequest(ctx context.Context, courseRequestIDs *[]string) error {
	var driver RepoDriver = cr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	UPDATE course_requests
	SET
		status = 'cancelled',
		expired_at = NULL,
		updated_at = NOW()
	WHERE NOW() >= expired_at 
	AND deleted_at IS NULL
	AND status IN ('reserved','pending payment')
	RETURNING (id)
	`

	rows, err := driver.Query(query)
	if err != nil {
		// not wrapping it on customerror because its just for cron
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var item string
		if err := rows.Scan(&item); err != nil {
			// not wrapping it on customerror because its just for cron
			return err
		}
		*courseRequestIDs = append(*courseRequestIDs, item)
	}
	return nil
}

func (cr *CourseRequestRepositoryImpl) ChangeRequestStatus(ctx context.Context, id string, newStatus string, eat *time.Time) error {
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

func (cr *CourseRequestRepositoryImpl) FindByID(ctx context.Context, id string, courseRequest *entity.CourseRequest) error {
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
		number_of_sessions,
		number_of_participant,
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
		&courseRequest.NumberOfSessions,
		&courseRequest.NumberOfParticipant,
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
		expired_at > CURRENT_TIMESTAMP AND
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
	INSERT INTO course_requests (student_id, course_id, number_of_sessions, number_of_participant, expired_at)
	VALUES
	($1, $2, $3, $4, $5)
	RETURNING (id)
	`
	row := driver.QueryRow(
		query,
		courseRequest.StudentID,
		courseRequest.CourseID,
		courseRequest.NumberOfSessions,
		courseRequest.NumberOfParticipant,
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

func (cr *CourseRequestRepositoryImpl) FindOngoingByCourseID(ctx context.Context, courseID int, count *int) error {
	var driver RepoDriver = cr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	SELECT COALESCE(COUNT(*), 0)
	FROM course_requests
	WHERE course_id = $1 AND status NOT IN ('scheduled', 'completed', 'cancelled') AND deleted_at IS NULL AND expired_at > CURRENT_TIMESTAMP
	`
	if err := driver.QueryRow(query, courseID).Scan(count); err != nil {
		return customerrors.NewError(
			"failed to get order count",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (cr *CourseRequestRepositoryImpl) FindCompletedByStudentIDAndCourseID(ctx context.Context, studentID string, courseID int, count *int) error {
	var driver RepoDriver = cr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	SELECT COALESCE(COUNT(*), 0)
	FROM course_requests
	WHERE student_id = $1 AND course_id = $2 AND status = 'completed' AND deleted_at IS NULL
	`
	if err := driver.QueryRow(query, studentID, courseID).Scan(count); err != nil {
		return customerrors.NewError(
			"failed to get order count",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}
