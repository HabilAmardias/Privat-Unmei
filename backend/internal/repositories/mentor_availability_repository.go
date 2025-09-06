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

type MentorAvailabilityRepositoryImpl struct {
	DB *db.CustomDB
}

func CreateCourseAvailabilityRepository(db *db.CustomDB) *MentorAvailabilityRepositoryImpl {
	return &MentorAvailabilityRepositoryImpl{db}
}

func (car *MentorAvailabilityRepositoryImpl) GetDOWAvailability(ctx context.Context, mentorID string, dows *[]int) error {
	var driver RepoDriver
	driver = car.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	SELECT
		day_of_week
	FROM mentor_availability
	WHERE mentor_id = $1 AND deleted_at IS NULL
	`
	rows, err := driver.Query(query, mentorID)
	if err != nil {
		return customerrors.NewError(
			"failed to get availability",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	defer rows.Close()
	for rows.Next() {
		var item int
		if err := rows.Scan(&item); err != nil {
			return customerrors.NewError(
				"failed to get availability",
				err,
				customerrors.DatabaseExecutionError,
			)
		}
		*dows = append(*dows, item)
	}
	return nil
}

func (car *MentorAvailabilityRepositoryImpl) GetAvailabilityByMentorID(ctx context.Context, mentorID string, scheds *[]entity.MentorAvailability) error {
	var driver RepoDriver
	driver = car.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	SELECT
		id,
		mentor_id,
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
	WHERE mentor_id = $1 AND deleted_at IS NULL
	`
	rows, err := driver.Query(query, mentorID)
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
			&item.MentorID,
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

func (car *MentorAvailabilityRepositoryImpl) DeleteAvailability(ctx context.Context, mentorID string) error {
	var driver RepoDriver
	driver = car.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	UPDATE mentor_availability
	SET deleted_at = NOW(), updated_at = NOW()
	WHERE mentor_id = $1 AND deleted_at IS NULL
	`
	_, err := driver.Exec(query, mentorID)
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
	INSERT INTO mentor_availability (mentor_id, day_of_week, start_time, end_time)
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
		($%d, $%d, $%d, $%d)
		`, sprintIndex, sprintIndex+1, sprintIndex+2, sprintIndex+3)
		}
		args = append(args, schedule.MentorID)
		args = append(args, schedule.DayOfWeek)
		args = append(args, schedule.StartTime.ToString())
		args = append(args, schedule.EndTime.ToString())
		sprintIndex += 4
	}
	query += `
	ON CONFLICT(mentor_id, day_of_week, start_time, end_time)
	DO UPDATE SET
		mentor_id = EXCLUDED.mentor_id,
		day_of_week = EXCLUDED.day_of_week,
		start_time = EXCLUDED.start_time,
		end_time = EXCLUDED.end_time,
		deleted_at = NULL,
		updated_at = NOW();
	`

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

func (car *MentorAvailabilityRepositoryImpl) CheckMentorAvailability(
	ctx context.Context,
	mentorID string,
	dates []time.Time,
	startTimes []string,
	endTimes []string,
	availabilityRes *entity.AvailabilityResult,
) error {
	var driver RepoDriver
	driver = car.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	WITH requested_slots AS (
		SELECT 
			unnest($1::date[]) as requested_date,
			unnest($2::text[])::time as requested_start_time,
			unnest($3::text[])::time as requested_end_time
	),
	mentor_available_slots AS (
		SELECT 
			rs.requested_date,
			rs.requested_start_time,
			rs.requested_end_time,
			CASE 
				WHEN ma.id IS NOT NULL THEN 1 
				ELSE 0 
			END as is_available
		FROM requested_slots rs
		LEFT JOIN mentor_availability ma ON (
			ma.mentor_id = $4
			AND ma.day_of_week = EXTRACT(DOW FROM rs.requested_date)
			AND ma.start_time <= rs.requested_start_time
			AND ma.end_time >= rs.requested_end_time
			AND ma.deleted_at IS NULL
		)
	)
	SELECT 
		COUNT(*) as total_requested,
		COUNT(CASE WHEN is_available = 1 THEN 1 END) as available_slots,
		array_agg(
			CASE WHEN is_available = 0 THEN 
				requested_date || ' ' || requested_start_time 
			END
		) FILTER (WHERE is_available = 0) as unavailable_slots
	FROM mentor_available_slots;
	`
	if err := driver.QueryRow(query, pq.Array(dates), pq.Array(startTimes), pq.Array(endTimes), mentorID).Scan(
		&availabilityRes.TotalRequested,
		&availabilityRes.AvailableSlots,
		(*pq.StringArray)(&availabilityRes.UnavailableSlots),
	); err != nil {
		return customerrors.NewError(
			"failed to get availability",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}
