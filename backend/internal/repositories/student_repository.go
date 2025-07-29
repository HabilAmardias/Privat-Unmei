package repositories

import (
	"context"
	"database/sql"
	"privat-unmei/internal/entity"
)

type StudentRepository struct {
	DB *sql.DB
}

func (sr *StudentRepository) AddNewStudent(ctx context.Context, student *entity.Student) error {
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
		return err
	}
	return nil
}
