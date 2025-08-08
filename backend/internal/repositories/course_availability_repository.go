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

type CourseAvailabilityRepositoryImpl struct {
	DB *sql.DB
}

func CreateCourseAvailabilityRepository(db *sql.DB) *CourseAvailabilityRepositoryImpl {
	return &CourseAvailabilityRepositoryImpl{db}
}

func (car *CourseAvailabilityRepositoryImpl) DeleteAvailabilityMultipleCourses(ctx context.Context, courseIDs []int) error {
	var driver RepoDriver
	driver = car.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	UPDATE course_availability
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

func (car *CourseAvailabilityRepositoryImpl) DeleteAvailability(ctx context.Context, courseID int) error {
	var driver RepoDriver
	driver = car.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	UPDATE course_availability
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

func (car *CourseAvailabilityRepositoryImpl) CreateAvailability(ctx context.Context, schedules *[]entity.CourseAvailability) error {
	var driver RepoDriver
	driver = car.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	INSERT INTO course_availability (course_id, day_of_week, start_time, end_time)
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
