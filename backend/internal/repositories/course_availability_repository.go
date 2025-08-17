package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/entity"

	"github.com/lib/pq"
)

type MentorAvailabilityRepositoryImpl struct {
	DB *sql.DB
}

func CreateCourseAvailabilityRepository(db *sql.DB) *MentorAvailabilityRepositoryImpl {
	return &MentorAvailabilityRepositoryImpl{db}
}

func (car *MentorAvailabilityRepositoryImpl) GetAvailabilityByCourseID(ctx context.Context, courseID int, scheds *[]entity.MentorAvailability) error {
	var driver RepoDriver
	driver = car.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	SELECT
		id,
		course_id,
		day_of_week,
		CAST(EXTRACT(HOUR from start_time) AS INT),
		CAST(EXTRACT(MINUTE from start_time) AS INT),
		CAST(EXTRACT(SECOND from start_time) AS INT),
		CAST(EXTRACT(HOUR from end_time) AS INT),
		CAST(EXTRACT(MINUTE from end_time) AS INT),
		CAST(EXTRACT(SECOND from end_time) AS INT),
		created_at,
		updated_at,
		deleted_at
	FROM mentor_availability
	WHERE course_id = $1 AND deleted_at IS NULL
	`
	rows, err := driver.Query(query, courseID)
	if err != nil {
		return customerrors.NewError(
			"failed to get course schedules",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	defer rows.Close()
	for rows.Next() {
		var item entity.MentorAvailability
		if err := rows.Scan(
			&item.ID,
			&item.CourseID,
			&item.DayOfWeek,
			&item.StartTime.Hour,
			&item.StartTime.Minute,
			&item.StartTime.Second,
			&item.EndTime.Hour,
			&item.EndTime.Minute,
			&item.EndTime.Second,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.DeletedAt,
		); err != nil {
			return customerrors.NewError(
				"failed to get course schedules",
				err,
				customerrors.DatabaseExecutionError,
			)
		}
		*scheds = append(*scheds, item)
	}
	return nil
}

func (car *MentorAvailabilityRepositoryImpl) DeleteAvailabilityMultipleCourses(ctx context.Context, courseIDs []int) error {
	var driver RepoDriver
	driver = car.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	UPDATE mentor_availability
	SET deleted_at = NOW(), updated_at = NOW()
	WHERE course_id = ANY($1) AND deleted_at IS NULL
	`
	_, err := driver.Exec(query, pq.Array(courseIDs))
	if err != nil {
		return customerrors.NewError(
			"failed to delete course schedules",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (car *MentorAvailabilityRepositoryImpl) DeleteAvailability(ctx context.Context, courseID int) error {
	var driver RepoDriver
	driver = car.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	UPDATE mentor_availability
	SET deleted_at = NOW(), updated_at = NOW()
	WHERE course_id = $1 AND deleted_at IS NULL
	`
	_, err := driver.Exec(query, courseID)
	if err != nil {
		return customerrors.NewError(
			"failed to delete course schedules",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (car *MentorAvailabilityRepositoryImpl) CreateAvailability(ctx context.Context, schedules *[]entity.MentorAvailability) error {
	var driver RepoDriver
	driver = car.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	INSERT INTO mentor_availability (course_id, day_of_week, start_time, end_time)
	VALUES
	`
	args := []any{}
	sprintIndex := 1
	for i, schedule := range *schedules {
		if i != len(*schedules)-1 {
			query += fmt.Sprintf(`
		($%d, $%d, $%d, $%d),
		`, sprintIndex, sprintIndex+1, sprintIndex+2, sprintIndex+3)
		} else {
			query += fmt.Sprintf(`
		($%d, $%d, $%d, $%d);
		`, sprintIndex, sprintIndex+1, sprintIndex+2, sprintIndex+3)
		}
		args = append(args, schedule.CourseID)
		args = append(args, schedule.DayOfWeek)
		args = append(args, fmt.Sprintf("%d:%d:%d", schedule.StartTime.Hour, schedule.StartTime.Minute, schedule.StartTime.Second))
		args = append(args, fmt.Sprintf("%d:%d:%d", schedule.EndTime.Hour, schedule.EndTime.Minute, schedule.EndTime.Second))
		sprintIndex += 4
	}
	log.Println(query)
	_, err := driver.Exec(query, args...)
	if err != nil {
		return customerrors.NewError(
			"failed to create available schedule",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}
