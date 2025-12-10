package repositories

import (
	"context"
	"database/sql"
	"errors"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/db"
	"privat-unmei/internal/entity"
)

type StudentRepositoryImpl struct {
	DB *db.CustomDB
}

func CreateStudentRepository(db *db.CustomDB) *StudentRepositoryImpl {
	return &StudentRepositoryImpl{db}
}

func (sr *StudentRepositoryImpl) GetStudentList(ctx context.Context, totalRow *int64, limit int, page int, students *[]entity.ListStudentQuery) error {
	var driver RepoDriver
	driver = sr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	countQuery := `
	SELECT count(*)
	FROM users u
	JOIN students s ON s.id = u.id
	WHERE u.deleted_at IS NULL AND s.deleted_at IS NULL
	`
	query := `
	SELECT
		s.id, 
		u.name, 
		u.public_id, 
		u.bio, 
		u.profile_image, 
		u.status
	FROM users u
	JOIN students s ON s.id = u.id
	WHERE u.deleted_at IS NULL AND s.deleted_at IS NULL
	`
	row := driver.QueryRow(countQuery)
	if err := row.Scan(totalRow); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return customerrors.NewError(
				"user not found",
				err,
				customerrors.ItemNotExist,
			)
		}
		return customerrors.NewError(
			"failed to get student list",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	args := []any{}
	args = append(args, limit)
	args = append(args, limit*(page-1))
	query += ` 
		LIMIT $1
		OFFSET $2
	`

	rows, err := driver.Query(query, args...)
	if err != nil {
		return customerrors.NewError(
			"failed to get student list",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	defer rows.Close()
	for rows.Next() {
		var student entity.ListStudentQuery
		if err := rows.Scan(
			&student.ID,
			&student.Name,
			&student.PublicID,
			&student.Bio,
			&student.ProfileImage,
			&student.Status,
		); err != nil {
			return customerrors.NewError(
				"failed to get student list",
				err,
				customerrors.DatabaseExecutionError,
			)
		}
		*students = append(*students, student)
	}
	return nil
}

func (sr *StudentRepositoryImpl) UpdateResetToken(ctx context.Context, id string, token *string) error {
	var driver RepoDriver
	driver = sr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	UPDATE students SET reset_token = $1, updated_at = NOW() WHERE id = $2 AND deleted_at IS NULL
	`
	_, err := driver.Exec(query, token, id)
	if err != nil {
		return customerrors.NewError(
			"failed to update data",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (sr *StudentRepositoryImpl) UpdateVerifyToken(ctx context.Context, id string, token *string) error {
	var driver RepoDriver
	driver = sr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	UPDATE students SET verify_token = $1, updated_at = NOW() WHERE id = $2 AND deleted_at IS NULL
	`
	_, err := driver.Exec(query, token, id)
	if err != nil {
		return customerrors.NewError(
			"failed to update data",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (sr *StudentRepositoryImpl) FindByID(ctx context.Context, id string, student *entity.Student) error {
	var driver RepoDriver
	driver = sr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	SELECT id, verify_token, reset_token, created_at, updated_at, deleted_at
	FROM students
	WHERE id = $1 AND deleted_at IS NULL
	`
	row := driver.QueryRow(query, id)
	if err := row.Scan(
		&student.ID,
		&student.VerifyToken,
		&student.ResetToken,
		&student.CreatedAt,
		&student.UpdatedAt,
		&student.DeletedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return customerrors.NewError(
				"user not found",
				err,
				customerrors.ItemNotExist,
			)
		}
		return customerrors.NewError(
			"failed to find user",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (sr *StudentRepositoryImpl) AddNewStudent(ctx context.Context, student *entity.Student) error {
	var driver RepoDriver
	driver = sr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	INSERT INTO students (id, reset_token, verify_token)
	VALUES
	($1, $2, $3);
	`
	_, err := driver.Exec(query, student.ID, student.ResetToken, student.VerifyToken)
	if err != nil {
		return customerrors.NewError(
			"failed to create account",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}
