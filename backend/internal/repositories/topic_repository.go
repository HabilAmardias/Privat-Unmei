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

type TopicRepositoryImpl struct {
	DB *sql.DB
}

func CreateTopicRepository(db *sql.DB) *TopicRepositoryImpl {
	return &TopicRepositoryImpl{db}
}

func (tr *TopicRepositoryImpl) GetTopicsByCourseID(ctx context.Context, courseID int, topics *[]entity.CourseTopic) error {
	var driver RepoDriver
	driver = tr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	SELECT
		id,
		course_id,
		title,
		description,
		created_at,
		updated_at,
		deleted_at
	FROM topics
	WHERE course_id = $1 AND deleted_at IS NULL
	`
	rows, err := driver.Query(query, courseID)
	if err != nil {
		return customerrors.NewError(
			"failed to get course topic",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	defer rows.Close()
	for rows.Next() {
		var item entity.CourseTopic
		if err := rows.Scan(
			&item.ID,
			&item.CourseID,
			&item.Title,
			&item.Description,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.DeletedAt,
		); err != nil {
			return customerrors.NewError(
				"failed to get course topic",
				err,
				customerrors.DatabaseExecutionError,
			)
		}
		*topics = append(*topics, item)
	}
	return nil
}

func (tr *TopicRepositoryImpl) DeleteTopicsMultipleCourse(ctx context.Context, courseIDs []int) error {
	var driver RepoDriver
	driver = tr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	UPDATE topics
	SET deleted_at = NOW(), updated_at = NOW()
	WHERE course_id = ANY($1) and deleted_at IS NULL
	`
	_, err := driver.Exec(query, pq.Array(courseIDs))
	if err != nil {
		return customerrors.NewError(
			"failed to delete course topics",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (tr *TopicRepositoryImpl) DeleteTopics(ctx context.Context, courseID int) error {
	var driver RepoDriver
	driver = tr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	UPDATE topics
	SET deleted_at = NOW(), updated_at = NOW()
	WHERE course_id = $1 and deleted_at IS NULL
	`
	_, err := driver.Exec(query, courseID)
	if err != nil {
		return customerrors.NewError(
			"failed to delete course topics",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (tr *TopicRepositoryImpl) CreateTopics(ctx context.Context, topics *[]entity.CourseTopic) error {
	var driver RepoDriver
	driver = tr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	INSERT INTO topics (course_id, title, description)
	VALUES
	`
	args := []any{}
	sprintIndex := 1
	for i, topic := range *topics {
		if i != len(*topics)-1 {
			query += fmt.Sprintf(`
		($%d, $%d, $%d),
		`, sprintIndex, sprintIndex+1, sprintIndex+2)
		} else {
			query += fmt.Sprintf(`
		($%d, $%d, $%d);
		`, sprintIndex, sprintIndex+1, sprintIndex+2)
		}
		args = append(args, topic.CourseID)
		args = append(args, topic.Title)
		args = append(args, topic.Description)
		sprintIndex += 3
	}
	log.Println(query)
	_, err := driver.Exec(query, args...)
	if err != nil {
		return customerrors.NewError(
			"failed to create course topic",
			err,
			customerrors.DatabaseExecutionError,
		)
	}

	return nil
}
