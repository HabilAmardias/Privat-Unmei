package repositories

import (
	"context"
	"database/sql"
	"privat-unmei/internal/constants"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/entity"
)

type CourseRequestRepositoryImpl struct {
	DB *sql.DB
}

func CreateCourseRequestRepository(db *sql.DB) *CourseRequestRepositoryImpl {
	return &CourseRequestRepositoryImpl{db}
}

func (cr *CourseRequestRepositoryImpl) FindCompletedByStudentIDAndCourseID(ctx context.Context, studentID string, courseID int, orders *[]entity.CourseOrder) error {
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
		price,
		duration_days,
		accepted_at,
		payment_due,
		start_date,
		end_date,
		created_at,
		updated_at,
		deleted_at
	FROM course_requests
	WHERE student_id = $1 AND course_id = $2 AND status = $3 AND deleted_at IS NULL
	`
	rows, err := driver.Query(query, studentID, courseID, constants.ConfirmedStatus)
	if err != nil {
		return customerrors.NewError(
			"failed to get order",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	defer rows.Close()
	for rows.Next() {
		var item entity.CourseOrder
		if err := rows.Scan(
			&item.ID,
			&item.StudentID,
			&item.CourseID,
			&item.Status,
			&item.Price,
			&item.DurationDays,
			&item.AcceptedAt,
			&item.PaymentDue,
			&item.StartDate,
			&item.EndDate,
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
