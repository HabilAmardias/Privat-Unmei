package repositories

import (
	"context"
	"database/sql"
	"privat-unmei/internal/entity"
)

type MentorRepository struct {
	DB *sql.DB
}

func (mr *MentorRepository) AddNewMentor(ctx context.Context, mentor *entity.Mentor) error {
	var driver RepoDriver
	driver = mr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	INSERT INTO mentors(resume_url, years_of_experience, whatsapp_number, degree, major, campus)
	VALUES ($1, $2, $3, $4, $5, $6);
	`
	_, err := driver.Exec(
		query,
		mentor.Resume,
		mentor.YearsOfExperience,
		mentor.WhatsappNumber,
		mentor.Degree,
		mentor.Major,
		mentor.Campus,
	)
	if err != nil {
		return err
	}
	return nil
}
