package repositories

import (
	"context"
	"fmt"
	"time"

	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/db"
	"privat-unmei/internal/entity"

	"github.com/lib/pq"
)

type CourseScheduleRepositoryImpl struct {
	DB *db.CustomDB
}

func CreateCourseScheduleRepository(db *db.CustomDB) *CourseScheduleRepositoryImpl {
	return &CourseScheduleRepositoryImpl{db}
}

func (csr *CourseScheduleRepositoryImpl) FindScheduleByCourseRequestID(ctx context.Context, courseRequestID int, schedules *[]entity.CourseRequestSchedule) error {
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
	WHERE course_request_id = $1 AND deleted_at IS NULL
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

func (csr *CourseScheduleRepositoryImpl) CompleteSchedule(ctx context.Context) error {
	var driver RepoDriver
	driver = csr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	UPDATE course_schedule
	SET
		status = 'completed',
		updated_at = NOW()
	WHERE
		end_time <= CURRENT_TIME AND
		scheduled_date <= CURRENT_DATE AND
		status = 'scheduled' AND
		deleted_at IS NULL
	`

	_, err := driver.Exec(query)
	return err
}

func (csr *CourseScheduleRepositoryImpl) CancelExpiredSchedule(ctx context.Context, ids []int) error {
	var driver RepoDriver
	driver = csr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	UPDATE course_schedule
	SET
		status = 'cancelled',
		updated_at = NOW()
	WHERE 
		course_request_id = ANY($1) 
		AND deleted_at IS NULL
	`

	_, err := driver.Exec(query, pq.Array(ids))
	// not wrapping it on customerror because its just for cron
	return err
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

func (csr *CourseScheduleRepositoryImpl) CheckScheduleConflicts(
	ctx context.Context,
	mentorID string,
	dates []time.Time,
	startTimes []string,
	endTimes []string,
	conflictingScheds *[]entity.ConflictingSchedule,
) error {
	var driver RepoDriver
	driver = csr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	WITH requested_slots AS (
		SELECT 
			unnest($1::date[]) as requested_date,
			unnest($2::text[])::time as requested_start_time,
			unnest($3::text[])::time as requested_end_time
	)
	SELECT 
		rs.requested_date,
		rs.requested_start_time,
		rs.requested_end_time,
		cs.id as conflicting_schedule_id
	FROM requested_slots rs
	LEFT JOIN course_schedule cs ON (
		cs.scheduled_date = rs.requested_date
		AND cs.start_time < rs.requested_end_time
		AND cs.end_time > rs.requested_start_time
		AND cs.deleted_at IS NULL
		AND cs.status IN ('reserved', 'scheduled')
	)
	JOIN course_requests cr ON (cs.course_request_id = cr.id)
	JOIN courses c ON (cr.course_id = c.id)
	WHERE c.mentor_id = $4::uuid;
	`

	rows, err := driver.Query(query, pq.Array(dates), pq.Array(startTimes), pq.Array(endTimes), mentorID)
	if err != nil {
		return customerrors.NewError(
			"failed to get conflicting schedules",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	defer rows.Close()
	for rows.Next() {
		var item entity.ConflictingSchedule
		if err := rows.Scan(
			&item.Date,
			&item.StartTime,
			&item.EndTime,
			&item.ScheduleID,
		); err != nil {
			return customerrors.NewError(
				"failed to get conflicting schedules",
				err,
				customerrors.DatabaseExecutionError,
			)
		}
		*conflictingScheds = append(*conflictingScheds, item)
	}
	return nil
}
