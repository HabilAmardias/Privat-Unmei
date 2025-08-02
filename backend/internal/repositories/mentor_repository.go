package repositories

import (
	"context"
	"database/sql"
	"errors"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/entity"
)

type MentorRepositoryImpl struct {
	DB *sql.DB
}

func CreateMentorRepositoryImpl(db *sql.DB) *MentorRepositoryImpl {
	return &MentorRepositoryImpl{db}
}

func (mr *MentorRepositoryImpl) FindByID(ctx context.Context, id string, mentor *entity.Mentor) error {
	var driver RepoDriver
	driver = mr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	SELECT 
		id, 
		total_rating, 
		rating_count, 
		resume_url, 
		years_of_experience, 
		whatsapp_number, 
		degree, 
		major, 
		campus, 
		created_at,
		updated_at,
		deleted_at
	FROM mentors
	WHERE id = $1 and deleted_at IS NULL
	`
	row := driver.QueryRow(
		query,
		id,
	)
	if err := row.Scan(
		&mentor.ID,
		&mentor.TotalRating,
		&mentor.RatingCount,
		&mentor.Resume,
		&mentor.YearsOfExperience,
		&mentor.WhatsappNumber,
		&mentor.Degree,
		&mentor.Major,
		&mentor.Campus,
		&mentor.CreatedAt,
		&mentor.UpdatedAt,
		&mentor.DeletedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return customerrors.NewError(
				customerrors.UserNotFound,
				err,
				customerrors.ItemNotExist,
			)
		}
		return customerrors.NewError(
			"failed to get user",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (mr *MentorRepositoryImpl) AddNewMentor(ctx context.Context, mentor *entity.Mentor) error {
	var driver RepoDriver
	driver = mr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	INSERT INTO mentors(id, resume_url, years_of_experience, whatsapp_number, degree, major, campus)
	VALUES ($1, $2, $3, $4, $5, $6, $7);
	`
	_, err := driver.Exec(
		query,
		mentor.ID,
		mentor.Resume,
		mentor.YearsOfExperience,
		mentor.WhatsappNumber,
		mentor.Degree,
		mentor.Major,
		mentor.Campus,
	)
	if err != nil {
		return customerrors.NewError(
			"failed to create account",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (mr *MentorRepositoryImpl) UpdateMentor(ctx context.Context, id string, queryEntity *entity.UpdateMentorQuery) error {
	var driver RepoDriver
	driver = mr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	UPDATE mentors
	SET
		resume_url = COALESCE($1, resume_url),
		years_of_experience = COALESCE($2, years_of_experience),
		whatsapp_number = COALESCE($3, whatsapp_number),
		degree = COALESCE($4, degree),
		major = COALESCE($5, major),
		campus = COALESCE($6, campus),
		updated_at = NOW()
	WHERE id = $7 AND deleted_at IS NULL;
	`
	_, err := driver.Exec(query, queryEntity.Resume, queryEntity.YearsOfExperience, queryEntity.WhatsappNumber, queryEntity.Degree, queryEntity.Major, queryEntity.Campus, id)
	if err != nil {
		return customerrors.NewError("failed to update mentor", err, customerrors.DatabaseExecutionError)
	}
	return nil
}

func (mr *MentorRepositoryImpl) DeleteMentor(ctx context.Context, id string) error {
	var driver RepoDriver
	driver = mr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	UPDATE mentors
	SET
		deleted_at = NOW(),
		updated_at = NOW()
	WHERE id = $1 AND deleted_at IS NULL
	`
	_, err := driver.Exec(query, id)
	if err != nil {
		return customerrors.NewError("failed to delete mentor", err, customerrors.DatabaseExecutionError)
	}
	return nil
}
