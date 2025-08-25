package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/entity"
	"time"
)

type CourseScheduleRepositoryImpl struct {
	DB *sql.DB
}

func CreateCourseScheduleRepository(db *sql.DB) *CourseScheduleRepositoryImpl {
	return &CourseScheduleRepositoryImpl{db}
}

func (csr *CourseScheduleRepositoryImpl) UpdateScheduleStatusByCourseRequestID(ctx context.Context, courseRequestID int, newStatus string) error {
	var driver RepoDriver
	driver = csr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	UPDATE course_schedule
	SET
		status = $1,
		updated_at = NOW()
	WHERE course_request_id = $2 AND deleted_at IS NULL
	`
	_, err := driver.Exec(query, newStatus, courseRequestID)
	if err != nil {
		return customerrors.NewError(
			"failed to update course schedule",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (csr *CourseScheduleRepositoryImpl) FindReservedScheduleByCourseRequestID(ctx context.Context, courseRequestID int, schedules *[]entity.CourseRequestSchedule) error {
	var driver RepoDriver
	driver = csr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	SELECT
		id,
		course_request_id,
		scheduled_date,
		start_time,
		end_time,
		status,
		created_at,
		updated_at,
		deleted_at
	FROM course_schedule
	WHERE course_request_id = $1 AND status = 'reserved' AND deleted_at IS NULL
	`
	rows, err := driver.Query(query, courseRequestID)
	if err != nil {
		return customerrors.NewError(
			"failed to get course request schedule",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	defer rows.Close()
	for rows.Next() {
		var item entity.CourseRequestSchedule
		if err := rows.Scan(
			&item.ID,
			&item.CourseRequestID,
			&item.ScheduledDate,
			&item.StartTime,
			&item.EndTime,
			&item.Status,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.DeletedAt,
		); err != nil {
			return customerrors.NewError(
				"failed to get couse request schedule",
				err,
				customerrors.DatabaseExecutionError,
			)
		}
		*schedules = append(*schedules, item)
	}
	return nil
}

func (csr *CourseScheduleRepositoryImpl) IsScheduleExist(ctx context.Context, mentorID string, date time.Time, startTime string, endTime string, existingSchedule *int64) error {
	var driver RepoDriver
	driver = csr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	SELECT COUNT(*) FROM course_schedule cs
	JOIN course_requests cr ON cs.course_request_id = cr.id
	JOIN courses c ON cr.course_id = c.id
	WHERE c.mentor_id = $1 AND cs.scheduled_date = $2 
	AND cs.deleted_at IS NULL AND cr.deleted_at IS NULL AND c.deleted_at IS NULL
	AND cs.status IN ('scheduled','reserved')
	AND NOT (cs.end_time <= $3 OR cs.start_time >= $4)
	`
	row := driver.QueryRow(
		query,
		mentorID,
		fmt.Sprintf("%d-%d-%d", date.Year(), date.Month(), date.Day()),
		startTime,
		endTime,
	)
	if err := row.Scan(existingSchedule); err != nil {
		return customerrors.NewError(
			"failed to get schedule",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (csr *CourseScheduleRepositoryImpl) CreateSchedule(ctx context.Context, courseRequestID int, slots *[]entity.CreateRequestSchedule) error {
	var driver RepoDriver
	driver = csr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	args := []any{courseRequestID}
	query := `
	INSERT INTO course_schedule (course_request_id, scheduled_date, start_time, end_time)
	VALUES
	`
	lastIndex := len(*slots) - 1
	sprintIndex := 2
	for i, slot := range *slots {
		if i < lastIndex {
			query += fmt.Sprintf(
				`($1, $%d, $%d, $%d),`,
				sprintIndex,
				sprintIndex+1,
				sprintIndex+2,
			)
		} else {
			query += fmt.Sprintf(
				`($1, $%d, $%d, $%d);`,
				sprintIndex,
				sprintIndex+1,
				sprintIndex+2,
			)
		}
		args = append(args, slot.Date, slot.StartTime, slot.EndTime)
		sprintIndex += 3
	}
	log.Println(query)
	_, err := driver.Exec(query, args...)
	if err != nil {
		return customerrors.NewError(
			"failed to create order",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}
