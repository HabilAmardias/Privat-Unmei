package repositories

import (
	"context"
	"database/sql"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/entity"
)

type StudentRepositoryImpl struct {
	DB *sql.DB
}

func CreateStudentRepository(db *sql.DB) *StudentRepositoryImpl {
	return &StudentRepositoryImpl{db}
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
