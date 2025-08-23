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
	AND cr.status IN ('scheduled', 'pending payment', 'reserved')
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

func (csr *CourseScheduleRepositoryImpl) CreateSchedule(ctx context.Context, courseRequestID int64, slots *[]entity.CreateRequestSchedule) error {
	var driver RepoDriver
	driver = csr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	args := []any{courseRequestID}
	query := `
	INSERT INTO course_schedule (course_request_id, session_number, scheduled_date, start_time, end_time)
	VALUES
	`
	lastIndex := len(*slots) - 1
	sprintIndex := 2
	for i, slot := range *slots {
		if i < lastIndex {
			query += fmt.Sprintf(
				`($1, $%d, $%d, $%d, $%d),`,
				sprintIndex,
				sprintIndex+1,
				sprintIndex+2,
				sprintIndex+3,
			)
		} else {
			query += fmt.Sprintf(
				`($1, $%d, $%d, $%d, $%d);`,
				sprintIndex,
				sprintIndex+1,
				sprintIndex+2,
				sprintIndex+3,
			)
		}
		args = append(args, i+1, slot.Date, slot.StartTime, slot.EndTime)
		sprintIndex += 4
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
