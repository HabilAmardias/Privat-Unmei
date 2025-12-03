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

type PaymentRepositoryImpl struct {
	DB *db.CustomDB
}

func CreatePaymentRepository(db *db.CustomDB) *PaymentRepositoryImpl {
	return &PaymentRepositoryImpl{db}
}

func (pr *PaymentRepositoryImpl) CreatePaymentDetail(ctx context.Context, id int, subtotal float64, operationalCost float64, totalPrice float64, methodName string, accountNumber string) error {
	var driver RepoDriver = pr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	INSERT INTO payments (
		course_request_id,
		subtotal,
		operational_cost,
		total_price,
		payment_method_name,
		account_number
	)
	VALUES
	($1, $2, $3, $4, $5, $6)
	`
	_, err := driver.Exec(query, id, subtotal, operationalCost, totalPrice, methodName, accountNumber)
	if err != nil {
		return customerrors.NewError(
			"failed to create payment",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (pr *PaymentRepositoryImpl) FindPaymentByRequestID(ctx context.Context, id int, payment *entity.Payment) error {
	var driver RepoDriver = pr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}

	query := `
	SELECT
		id,
		course_request_id,
		subtotal,
		operational_cost,
		total_price,
		payment_method_name,
		account_number,
		created_at,
		updated_at,
		deleted_at
	FROM payments
	WHERE course_request_id = $1 AND deleted_at IS NULL
	`
	if err := driver.QueryRow(query, id).Scan(
		&payment.ID,
		&payment.CourseRequestID,
		&payment.SubTotal,
		&payment.OperationalCost,
		&payment.TotalPrice,
		&payment.PaymentMethodName,
		&payment.AccountNumber,
		&payment.CreatedAt,
		&payment.UpdatedAt,
		&payment.DeletedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return customerrors.NewError(
				"payment does not found",
				err,
				customerrors.ItemNotExist,
			)
		}
		return customerrors.NewError(
			"failed to get payment",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (pr *PaymentRepositoryImpl) GetLeastPaymentMethodCount(ctx context.Context, id int, count *int) error {
	var driver RepoDriver = pr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}

	query := `
	with counts as (
		select mp.mentor_id, count(*) as payment_method_count
		from mentor_payments mp
		where mp.mentor_id in (
			select mentor_id
			from mentor_payments
			where payment_method_id = $1 and deleted_at is null
		) and mp.deleted_at is null
		group by mp.mentor_id
	)
	select coalesce(min(c.payment_method_count), 0) from counts c;
	`
	if err := driver.QueryRow(query, id).Scan(count); err != nil {
		return customerrors.NewError(
			"failed to get payment method",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (pr *PaymentRepositoryImpl) HardDeleteMentorPayment(ctx context.Context, mentorID string) error {
	var driver RepoDriver
	driver = pr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	DELETE FROM mentor_payments
	WHERE mentor_id = $1 AND deleted_at IS NULL
	`
	_, err := driver.Exec(query, mentorID)
	return err
}

func (pr *PaymentRepositoryImpl) GetMentorPaymentMethod(
	ctx context.Context,
	mentorID string,
	methods *[]entity.GetMentorPaymentMethodQuery,
) error {
	var driver RepoDriver
	driver = pr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	SELECT
		mp.payment_method_id,
		pm.name,
		mp.account_number
	FROM mentor_payments mp
	JOIN payment_methods pm ON mp.payment_method_id = pm.id
	WHERE mp.mentor_id = $1 AND mp.deleted_at IS NULL AND pm.deleted_at IS NULL
	`
	rows, err := driver.Query(query, mentorID)
	if err != nil {
		return customerrors.NewError(
			"failed to get mentor payment method",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	defer rows.Close()
	for rows.Next() {
		var item entity.GetMentorPaymentMethodQuery
		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.AccountNumber,
		); err != nil {
			return customerrors.NewError(
				"failed to get mentor payment method",
				err,
				customerrors.DatabaseExecutionError,
			)
		}
		*methods = append(*methods, item)
	}
	return nil
}

func (pr *PaymentRepositoryImpl) GetAllPaymentMethod(
	ctx context.Context,
	search *string,
	limit int,
	page int,
	totalRow *int64,
	methods *[]entity.GetPaymentMethodQuery,
) error {
	var driver RepoDriver
	driver = pr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	args := []any{}
	countArgs := []any{}
	query := `
	SELECT
		id,
		name
	FROM payment_methods
	WHERE deleted_at IS NULL
	`
	countQuery := `
	SELECT count(*)
	FROM payment_methods
	WHERE deleted_at IS NULL
	`
	if search != nil {
		args = append(args, "%"+*search+"%")
		countArgs = append(countArgs, "%"+*search+"%")
		query += fmt.Sprintf(`
		AND name ILIKE $%d
		`, len(args))
		countQuery += fmt.Sprintf(`
		AND name ILIKE $%d
		`, len(countArgs))
	}

	args = append(args, limit)
	query += fmt.Sprintf(` LIMIT $%d`, len(args))

	args = append(args, limit*(page-1))
	query += fmt.Sprintf(" OFFSET $%d", len(args))
	if err := driver.QueryRow(countQuery, countArgs...).Scan(totalRow); err != nil {
		return customerrors.NewError(
			"failed to get payment method list",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	rows, err := driver.Query(query, args...)
	if err != nil {
		return customerrors.NewError(
			"failed to get payment method list",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	defer rows.Close()
	for rows.Next() {
		var item entity.GetPaymentMethodQuery
		if err := rows.Scan(
			&item.ID,
			&item.Name,
		); err != nil {
			return customerrors.NewError(
				"failed to get payment method list",
				err,
				customerrors.DatabaseExecutionError,
			)
		}
		*methods = append(*methods, item)
	}
	return nil
}

func (pr *PaymentRepositoryImpl) UpdatePaymentMethod(ctx context.Context, newName *string, id int) error {
	var driver RepoDriver
	driver = pr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	UPDATE payment_methods
	SET name = COALESCE($1, name)
	WHERE id = $2 AND deleted_at IS NULL
	`
	_, err := driver.Exec(query, newName, id)
	if err != nil {
		return customerrors.NewError(
			"failed to update payment method",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (pr *PaymentRepositoryImpl) DeletePaymentMethod(ctx context.Context, paymentMethodID int) error {
	var driver RepoDriver
	driver = pr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	UPDATE payment_methods
	SET deleted_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
	WHERE id = $1 AND deleted_at IS NULL
	`
	_, err := driver.Exec(query, paymentMethodID)
	if err != nil {
		return customerrors.NewError(
			"failed to delete payment method",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (pr *PaymentRepositoryImpl) UnassignPaymentMethodFromAllMentor(ctx context.Context, paymentMethodID int) error {
	var driver RepoDriver
	driver = pr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}

	query := `
	UPDATE mentor_payments
	SET deleted_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
	WHERE payment_method_id = $1 AND deleted_at IS NULL
	`
	_, err := driver.Exec(query, paymentMethodID)
	if err != nil {
		return customerrors.NewError(
			"failed to delete payment method",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (pr *PaymentRepositoryImpl) FindPaymentMethodByName(ctx context.Context, paymentName string, count *int64) error {
	var driver RepoDriver
	driver = pr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	SELECT COUNT(*)
	FROM payment_methods
	WHERE LOWER(name) = LOWER($1) AND deleted_at IS NULL
	`
	if err := driver.QueryRow(query, paymentName).Scan(count); err != nil {
		return customerrors.NewError(
			"failed to find payment method",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (pr *PaymentRepositoryImpl) CreatePaymentMethod(ctx context.Context, paymentName string, method *entity.PaymentMethod) error {
	var driver RepoDriver
	driver = pr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	INSERT INTO payment_methods (name)
	VALUES
	($1)
	ON CONFLICT (name) 
	DO UPDATE SET
		deleted_at = NULL,
		updated_at = CURRENT_TIMESTAMP,
		name = EXCLUDED.name
	RETURNING id, name, created_at, updated_at, deleted_at
	`
	if err := driver.QueryRow(query, paymentName).Scan(
		&method.ID,
		&method.Name,
		&method.CreatedAt,
		&method.UpdatedAt,
		&method.DeletedAt,
	); err != nil {
		return customerrors.NewError(
			"failed to create payment method",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (pr *PaymentRepositoryImpl) UnassignPaymentMethodFromMentor(ctx context.Context, mentorID string) error {
	var driver RepoDriver
	driver = pr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	UPDATE mentor_payments
	SET deleted_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
	WHERE mentor_id = $1 AND deleted_at IS NULL
	`
	_, err := driver.Exec(query, mentorID)
	if err != nil {
		return customerrors.NewError(
			"failed to unassign payment method",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (pr *PaymentRepositoryImpl) AssignPaymentMethodToMentor(ctx context.Context, mentorID string, paymentInfo []entity.AddMentorPaymentInfo) error {
	var driver RepoDriver
	driver = pr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	INSERT INTO mentor_payments (payment_method_id, mentor_id, account_number)
	VALUES
	`
	args := []any{}
	sprintIndex := 1
	for i, info := range paymentInfo {
		if i != len(paymentInfo)-1 {
			query += fmt.Sprintf(`
			($%d, $%d, $%d),
			`, sprintIndex, sprintIndex+1, sprintIndex+2)
		} else {
			query += fmt.Sprintf(`
			($%d, $%d, $%d)
			`, sprintIndex, sprintIndex+1, sprintIndex+2)
		}
		args = append(args, info.PaymentMethodID, mentorID, info.AccountNumber)
		sprintIndex += 3
	}
	query += `
	ON CONFLICT (mentor_id, payment_method_id)
	DO UPDATE SET
		mentor_id = EXCLUDED.mentor_id,
		payment_method_id = EXCLUDED.payment_method_id,
		account_number = EXCLUDED.account_number,
		deleted_at = NULL,
		updated_at = CURRENT_TIMESTAMP
	`
	_, err := driver.Exec(query, args...)
	if err != nil {
		return customerrors.NewError(
			"failed to create mentor payment method",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (pr *PaymentRepositoryImpl) MentorPaymentInfo(ctx context.Context, mentorID string, paymentInfo *[]entity.MentorPaymentInfo) error {
	var driver RepoDriver
	driver = pr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	SELECT
		pm.name,
		mp.account_number,
		pm.id
	FROM mentor_payments mp
	JOIN payment_methods pm ON pm.id = mp.payment_method_id
	WHERE mp.mentor_id = $1 AND pm.deleted_at IS NULL AND mp.deleted_at IS NULL
	`
	rows, err := driver.Query(query, mentorID)
	if err != nil {
		return customerrors.NewError(
			"failed to get mentor payment info",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	defer rows.Close()
	for rows.Next() {
		var item entity.MentorPaymentInfo
		if err := rows.Scan(
			&item.PaymentMethodName,
			&item.AccountNumber,
			&item.PaymentMethodID,
		); err != nil {
			return customerrors.NewError(
				"failed to get mentor payment info",
				err,
				customerrors.DatabaseExecutionError,
			)
		}
		*paymentInfo = append(*paymentInfo, item)
	}
	return nil
}

func (pr *PaymentRepositoryImpl) FindPaymentMethodByID(ctx context.Context, id int, method *entity.PaymentMethod) error {
	var driver RepoDriver
	driver = pr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	SELECT
		id,
		name,
		created_at,
		updated_at,
		deleted_at
	FROM payment_methods
	WHERE id = $1 AND deleted_at IS NULL
	`
	if err := driver.QueryRow(query, id).Scan(
		&method.ID,
		&method.Name,
		&method.CreatedAt,
		&method.UpdatedAt,
		&method.DeletedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return customerrors.NewError(
				"payment method not found",
				err,
				customerrors.ItemNotExist,
			)
		}
		return customerrors.NewError(
			"failed to get payment method detail",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (pr *PaymentRepositoryImpl) GetPaymentInfoByMentorAndMethodID(ctx context.Context, mentorID string, paymentMethodID int, mentorPayment *entity.MentorPayment) error {
	var driver RepoDriver
	driver = pr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	SELECT
		mp.payment_method_id,
		p.name,
		mp.mentor_id,
		mp.account_number
	FROM mentor_payments mp
	JOIN payment_methods p ON p.id = mp.payment_method_id
	WHERE mp.mentor_id = $1 AND mp.payment_method_id = $2 AND mp.deleted_at IS NULL AND p.deleted_at IS NULL
	`
	if err := driver.QueryRow(query, mentorID, paymentMethodID).Scan(
		&mentorPayment.PaymentMethodID,
		&mentorPayment.PaymentMethodName,
		&mentorPayment.MentorID,
		&mentorPayment.AccountNumber,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return customerrors.NewError(
				"mentor payment info not found",
				err,
				customerrors.ItemNotExist,
			)
		}
		return customerrors.NewError(
			"failed to get mentor payment info",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}
