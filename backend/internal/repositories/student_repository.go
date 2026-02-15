package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/db"
	"privat-unmei/internal/entity"
	"time"
)

type StudentRepositoryImpl struct {
	DB *db.CustomDB
}

func CreateStudentRepository(db *db.CustomDB) *StudentRepositoryImpl {
	return &StudentRepositoryImpl{db}
}

func (sr *StudentRepositoryImpl) UpdateLoginToken(ctx context.Context, studentID string, loginToken *string) error {
	var driver RepoDriver = sr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	UPDATE students
	SET login_token = $1, updated_at = NOW()
	WHERE id = $2 AND deleted_at IS NULL
	`

	_, err := driver.Exec(query, loginToken, studentID)
	if err != nil {
		return customerrors.NewError(
			"something went wrong",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (sr *StudentRepositoryImpl) UpdateOTP(ctx context.Context, studentID string, lastUpdateOTP *time.Time, otp *int64) error {
	var driver RepoDriver = sr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	UPDATE students
	SET otp = $1, otp_last_updated_at = $2, updated_at = NOW()
	WHERE id = $3 AND deleted_at IS NULL
	`

	_, err := driver.Exec(query, otp, lastUpdateOTP, studentID)
	if err != nil {
		return customerrors.NewError(
			"something went wrong",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (sr *StudentRepositoryImpl) DeleteStudent(ctx context.Context, studentID string) error {
	var driver RepoDriver = sr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	UPDATE students
	SET deleted_at = CURRENT_TIMESTAMP
	WHERE id = $1 AND deleted_at IS NULL
	`
	_, err := driver.Exec(query, studentID)
	if err != nil {
		return customerrors.NewError(
			"something went wrong",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (sr *StudentRepositoryImpl) GetStudentList(ctx context.Context, totalRow *int64, limit int, page int, students *[]entity.ListStudentQuery, search *string) error {
	var driver RepoDriver
	driver = sr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	args := []any{}
	countArgs := []any{}
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
		u.status
	FROM users u
	JOIN students s ON s.id = u.id
	WHERE u.deleted_at IS NULL AND s.deleted_at IS NULL
	`
	if search != nil {
		args = append(args, "%"+*search+"%")
		countArgs = append(countArgs, "%"+*search+"%")
		query += fmt.Sprintf(" AND (u.name ILIKE $%d OR u.public_id ILIKE $%d)", len(args), len(args))
		countQuery += fmt.Sprintf(" AND (u.name ILIKE $%d OR u.public_id ILIKE $%d)", len(countArgs), len(countArgs))
	}
	if err := driver.QueryRow(countQuery, countArgs...).Scan(totalRow); err != nil {
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
	args = append(args, limit)
	query += fmt.Sprintf(` LIMIT $%d`, len(args))
	args = append(args, limit*(page-1))
	query += fmt.Sprintf(` OFFSET $%d`, len(args))

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
	SELECT id, verify_token, reset_token, login_token, otp, otp_last_updated_at, created_at, updated_at, deleted_at
	FROM students
	WHERE id = $1 AND deleted_at IS NULL
	`
	row := driver.QueryRow(query, id)
	if err := row.Scan(
		&student.ID,
		&student.VerifyToken,
		&student.ResetToken,
		&student.LoginToken,
		&student.OTP,
		&student.OTPLastUpdatedAt,
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
