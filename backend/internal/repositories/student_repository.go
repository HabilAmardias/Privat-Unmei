package repositories

import (
	"context"
	"database/sql"
	"errors"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/entity"
)

type StudentRepositoryImpl struct {
	DB *sql.DB
}

func CreateStudentRepository(db *sql.DB) *StudentRepositoryImpl {
	return &StudentRepositoryImpl{db}
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
	WHERE id = $1
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
				customerrors.UserNotFound,
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
	INSERT INTO students (id)
	VALUES
	($1);
	`
	_, err := driver.Exec(query, student.ID)
	if err != nil {
		return customerrors.NewError(
			"failed to create account",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}
