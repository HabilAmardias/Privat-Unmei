package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/db"
	"privat-unmei/internal/entity"
)

type MentorRepositoryImpl struct {
	DB *db.CustomDB
}

func CreateMentorRepositoryImpl(db *db.CustomDB) *MentorRepositoryImpl {
	return &MentorRepositoryImpl{db}
}

func (mr *MentorRepositoryImpl) HardDeleteMentor(ctx context.Context, id string) error {
	var driver RepoDriver
	driver = mr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	DELETE FROM mentors
	WHERE id = $1 AND deleted_at IS NULL
	`
	_, err := driver.Exec(query, id)
	return err
}

func (mr *MentorRepositoryImpl) GetMentorList(ctx context.Context, mentors *[]entity.ListMentorQuery, totalRow *int64, param entity.ListMentorParam) error {
	var driver RepoDriver
	driver = mr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	args := []any{}
	countArgs := []any{}
	query := `
	SELECT
		m.id,
		u.name,
		u.email,
		m.years_of_experience
	FROM users u
	JOIN mentors m ON m.id = u.id
	WHERE
		m.deleted_at IS NULL AND
		u.deleted_at IS NULL
	`
	countQuery := `
	SELECT
		count(*)
	FROM users u
	JOIN mentors m ON m.id = u.id
	WHERE
		m.deleted_at IS NULL AND
		u.deleted_at IS NULL
	`
	if param.Search != nil {
		query += " AND "
		countQuery += " AND "
		query += fmt.Sprintf("u.name ILIKE $1 AND u.email ILIKE $%d", len(args)+1)
		countQuery += fmt.Sprintf("u.name ILIKE $1 AND u.email ILIKE $%d", len(countArgs)+1)
		args = append(args, "%"+*param.Search+"%")
		countArgs = append(countArgs, "%"+*param.Search+"%")
	}
	if param.SortYearOfExperience != nil {
		query += " ORDER BY years_of_experience"
		if *param.SortYearOfExperience {
			query += " ASC"
		} else {
			query += " DESC"
		}
	}
	query += fmt.Sprintf(` 
		LIMIT $%d
		OFFSET $%d
		`, len(args)+1, len(args)+2)
	args = append(args, param.Limit)
	args = append(args, param.Limit*(param.Page-1))

	row := driver.QueryRow(countQuery, countArgs...)
	if err := row.Scan(totalRow); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return customerrors.NewError(
				customerrors.UserNotFound,
				err,
				customerrors.ItemNotExist,
			)
		}
		return customerrors.NewError(
			"failed to get mentor list",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	rows, err := driver.Query(query, args...)
	if err != nil {
		return customerrors.NewError(
			"failed to get mentor list",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	defer rows.Close()
	for rows.Next() {
		var mentor entity.ListMentorQuery
		if err := rows.Scan(
			&mentor.ID,
			&mentor.Name,
			&mentor.Email,
			&mentor.YearsOfExperience,
		); err != nil {
			return customerrors.NewError(
				"failed to get mentor list",
				err,
				customerrors.DatabaseExecutionError,
			)
		}
		*mentors = append(*mentors, mentor)
	}
	return nil
}

func (mr *MentorRepositoryImpl) FindByID(ctx context.Context, id string, mentor *entity.Mentor, lock bool) error {
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
		degree, 
		major, 
		campus, 
		created_at,
		updated_at,
		deleted_at
	FROM mentors
	WHERE id = $1 and deleted_at IS NULL
	`
	if lock {
		query += `
		FOR UPDATE
		`
	}
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
	INSERT INTO mentors(id, resume_url, years_of_experience, degree, major, campus)
	VALUES ($1, $2, $3, $4, $5, $6);
	`
	_, err := driver.Exec(
		query,
		mentor.ID,
		mentor.Resume,
		mentor.YearsOfExperience,
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
		degree = COALESCE($3, degree),
		major = COALESCE($4, major),
		campus = COALESCE($5, campus),
		total_rating = COALESCE($6, total_rating),
		rating_count = COALESCE($7, rating_count),
		updated_at = NOW()
	WHERE id = $8 AND deleted_at IS NULL;
	`
	_, err := driver.Exec(
		query,
		queryEntity.Resume,
		queryEntity.YearsOfExperience,
		queryEntity.Degree,
		queryEntity.Major,
		queryEntity.Campus,
		queryEntity.TotalRating,
		queryEntity.RatingCount,
		id,
	)
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
