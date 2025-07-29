package repositories

import (
	"context"
	"database/sql"
	"privat-unmei/internal/entity"
)

type UserRepositoryImpl struct {
	DB *sql.DB
}

func (ur *UserRepositoryImpl) AddNewUser(ctx context.Context, user *entity.User) error {
	var driver RepoDriver
	driver = ur.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	INSERT INTO users (name, email, password_hash, bio, profile_image)
	VALUES
	($1, $2, $3, $4, $5)
	RETURNING id;
	`
	row := driver.QueryRow(query,
		user.Name,
		user.Email,
		user.Password,
		user.Bio,
		user.ProfileImage)
	if err := row.Scan(user.ID); err != nil {
		return err
	}

	return nil
}
