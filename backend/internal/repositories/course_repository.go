package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/entity"
)

type CourseRepositoryImpl struct {
	DB *sql.DB
}

func CreateCourseRepository(db *sql.DB) *CourseRepositoryImpl {
	return &CourseRepositoryImpl{db}
}

func (cr *CourseRepositoryImpl) UpdateCourseTransactionCount(ctx context.Context) error {
	var driver RepoDriver
	driver = cr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}

	query := `
	WITH requests_to_cancel AS (
		SELECT id, course_id
		FROM course_requests 
		WHERE expired_at < CURRENT_TIMESTAMP 
		AND status IN ('reserved','pending payment')
		AND deleted_at IS NULL
	),
	course_cancel_counts AS (
		SELECT course_id, COUNT(*) as cancel_count
		FROM requests_to_cancel
		GROUP BY course_id
	)
	UPDATE courses 
	SET transaction_count = transaction_count - course_cancel_counts.cancel_count,
		updated_at = NOW()
	FROM course_cancel_counts
	WHERE courses.id = course_cancel_counts.course_id AND courses.deleted_at IS NULL;
	`
	log.Println(query)
	_, err := driver.Exec(query)
	// not wrapping it on customerror because its just for cron
	return err
}

func (cr *CourseRepositoryImpl) UpdateCourse(ctx context.Context, courseID int, updateQuery *entity.UpdateCourseQuery) error {
	var driver RepoDriver
	driver = cr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	UPDATE courses
	SET
		title = COALESCE($1, title),
		description = COALESCE($2, description),
		domicile = COALESCE($3, domicile),
		method = COALESCE($4, method),
		price = COALESCE($5, price),
		session_duration_minutes = COALESCE($6, session_duration_minutes),
		max_total_session = COALESCE($7, max_total_session),
		transaction_count = COALESCE($8, transaction_count),
		updated_at = NOW()
	WHERE id = $8 AND deleted_at IS NULL
	`
	_, err := driver.Exec(
		query,
		updateQuery.Title,
		updateQuery.Description,
		updateQuery.Domicile,
		updateQuery.Method,
		updateQuery.Price,
		updateQuery.SessionDuration,
		updateQuery.MaxSession,
		courseID,
	)
	if err != nil {
		return customerrors.NewError(
			"failed to update course",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (cr *CourseRepositoryImpl) CourseDetail(ctx context.Context, query *entity.CourseDetailQuery, courseID int) error {
	var driver RepoDriver
	driver = cr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	sqlQuery := `
		SELECT
		c.id,
		c.mentor_id,
		u.name,
		u.email,
		c.title,
		c.description,
		c.domicile,
		c.price,
		c.method,
		c.session_duration_minutes,
		c.max_total_session,
		COALESCE(
			(SELECT STRING_AGG(cc_all.name, ',') 
			FROM course_category_assignments cca_all 
			JOIN course_categories cc_all ON cca_all.category_id = cc_all.id 
			WHERE cca_all.course_id = c.id 
			AND cca_all.deleted_at IS NULL 
			AND cc_all.deleted_at IS NULL), 
			''
		) AS categories
	FROM courses c
	JOIN users u ON u.id = c.mentor_id
	WHERE c.id = $1 AND c.deleted_at IS NULL AND u.deleted_at IS NULL
	`
	log.Println(sqlQuery)
	row := driver.QueryRow(sqlQuery, courseID)
	if err := row.Scan(
		&query.ID,
		&query.MentorID,
		&query.MentorName,
		&query.MentorEmail,
		&query.Title,
		&query.Description,
		&query.Domicile,
		&query.Price,
		&query.Method,
		&query.SessionDuration,
		&query.MaxSession,
		&query.CourseCategories,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return customerrors.NewError(
				"course not found",
				err,
				customerrors.ItemNotExist,
			)
		}
		return customerrors.NewError(
			"failed to get course detail",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (cr *CourseRepositoryImpl) ListCourse(
	ctx context.Context,
	query *[]entity.CourseListQuery,
	totalRow *int64,
	param entity.ListCourseParam,
) error {
	var driver RepoDriver
	driver = cr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	countArgs := []any{param.LastID}
	args := []any{param.LastID}
	sqlQuery := `
		SELECT
		c.id,
		c.mentor_id,
		u.name,
		u.email,
		c.title,
		c.domicile,
		c.price,
		c.method,
		c.session_duration_minutes,
		c.max_total_session,
		COALESCE(
			(SELECT STRING_AGG(cc_all.name, ',') 
			FROM course_category_assignments cca_all 
			JOIN course_categories cc_all ON cca_all.category_id = cc_all.id 
			WHERE cca_all.course_id = c.id 
			AND cca_all.deleted_at IS NULL 
			AND cc_all.deleted_at IS NULL), 
			''
		) AS categories
	FROM courses c
	JOIN users u ON u.id = c.mentor_id
	WHERE c.deleted_at IS NULL AND u.deleted_at IS NULL AND c.id < $1
	`
	countQuery := `
	SELECT count(*)
	FROM courses c
	JOIN users u ON u.id = c.mentor_id
	WHERE c.deleted_at IS NULL AND u.deleted_at IS NULL AND c.id < $1
	`
	if param.Search != nil {
		countQuery += fmt.Sprintf(" AND (c.title ILIKE $%d OR u.name ILIKE $%d) ", len(countArgs)+1, len(countArgs)+1)
		sqlQuery += fmt.Sprintf(" AND (c.title ILIKE $%d OR u.name ILIKE $%d) ", len(args)+1, len(args)+1)
		args = append(args, "%"+*param.Search+"%")
		countArgs = append(countArgs, "%"+*param.Search+"%")
	}
	if param.Method != nil {
		countQuery += fmt.Sprintf(" AND c.method = $%d", len(countArgs)+1)
		sqlQuery += fmt.Sprintf(" AND c.method = $%d", len(args)+1)
		args = append(args, *param.Method)
		countArgs = append(countArgs, *param.Method)
	}
	if param.CourseCategory != nil {
		countQuery += fmt.Sprintf(` AND EXISTS (
			SELECT 1 
			FROM course_category_assignments cca_filter 
			JOIN course_categories cc_filter ON cca_filter.category_id = cc_filter.id
			WHERE cca_filter.course_id = c.id 
			AND cca_filter.deleted_at IS NULL
			AND cc_filter.deleted_at IS NULL
			AND cc_filter.id = $%d
		)`, len(args)+1)
		sqlQuery += fmt.Sprintf(` AND EXISTS (
			SELECT 1 
			FROM course_category_assignments cca_filter 
			JOIN course_categories cc_filter ON cca_filter.category_id = cc_filter.id
			WHERE cca_filter.course_id = c.id 
			AND cca_filter.deleted_at IS NULL
			AND cc_filter.deleted_at IS NULL
			AND cc_filter.id = $%d
		)`, len(args)+1)
		countArgs = append(countArgs, *param.CourseCategory)
		args = append(args, *param.CourseCategory)
	}
	sqlQuery += " ORDER BY c.id DESC "
	sqlQuery += fmt.Sprintf(" LIMIT $%d ", len(args)+1)
	args = append(args, param.Limit)
	row := driver.QueryRow(countQuery, countArgs...)
	if err := row.Scan(totalRow); err != nil {
		return customerrors.NewError(
			"failed to get course list",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	log.Println(sqlQuery)
	rows, err := driver.Query(sqlQuery, args...)
	if err != nil {
		return customerrors.NewError(
			"failed to get courses list",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	defer rows.Close()

	for rows.Next() {
		var item entity.CourseListQuery
		if err := rows.Scan(
			&item.ID,
			&item.MentorID,
			&item.MentorName,
			&item.MentorEmail,
			&item.Title,
			&item.Domicile,
			&item.Price,
			&item.Method,
			&item.SessionDuration,
			&item.MaxSession,
			&item.CourseCategories,
		); err != nil {
			return customerrors.NewError(
				"failed to get course list",
				err,
				customerrors.DatabaseExecutionError,
			)
		}
		*query = append(*query, item)
	}

	return nil
}
func (cr *CourseRepositoryImpl) GetMostBoughtCourses(ctx context.Context, courses *[]entity.CourseListQuery) error {
	var driver RepoDriver
	driver = cr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	sqlQuery := `
		SELECT
		c.id,
		c.mentor_id,
		u.name,
		u.email,
		c.title,
		c.domicile,
		c.price,
		c.method,
		c.session_duration_minutes,
		c.max_total_session,
		COALESCE(
			(SELECT STRING_AGG(cc_all.name, ',') 
			FROM course_category_assignments cca_all 
			JOIN course_categories cc_all ON cca_all.category_id = cc_all.id 
			WHERE cca_all.course_id = c.id 
			AND cca_all.deleted_at IS NULL 
			AND cc_all.deleted_at IS NULL), 
			''
		) AS categories
	FROM courses c
	JOIN users u ON u.id = c.mentor_id
	WHERE c.deleted_at IS NULL AND u.deleted_at IS NULL
	ORDER BY c.transaction_count DESC
	LIMIT 10
	`
	rows, err := driver.Query(sqlQuery)
	if err != nil {
		return customerrors.NewError(
			"failed to get courses",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	defer rows.Close()

	for rows.Next() {
		var item entity.CourseListQuery
		if err := rows.Scan(
			&item.ID,
			&item.MentorID,
			&item.MentorName,
			&item.MentorEmail,
			&item.Title,
			&item.Domicile,
			&item.Price,
			&item.Method,
			&item.SessionDuration,
			&item.MaxSession,
			&item.CourseCategories,
		); err != nil {
			return customerrors.NewError(
				"failed to get courses",
				err,
				customerrors.DatabaseExecutionError,
			)
		}
		*courses = append(*courses, item)
	}
	return nil
}

func (cr *CourseRepositoryImpl) GetMaximumTransactionCount(ctx context.Context, transactionCount *int64, mentorID string) error {
	var driver RepoDriver
	driver = cr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	SELECT MAX(transaction_count)
	FROM courses
	WHERE mentor_id = $1 AND deleted_at IS NULL
	`
	row := driver.QueryRow(query, mentorID)
	if err := row.Scan(transactionCount); err != nil {
		return customerrors.NewError(
			"failed to get course transaction count",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (cr *CourseRepositoryImpl) FindByMentorID(ctx context.Context, mentorID string, courseIDs *[]int) error {
	var driver RepoDriver
	driver = cr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	SELECT
		id
	FROM courses
	WHERE mentor_id = $1 and deleted_at IS NULL
	`
	rows, err := driver.Query(query, mentorID)
	if err != nil {
		return customerrors.NewError(
			"failed to get mentor courses",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return customerrors.NewError(
				"failed to get courses",
				err,
				customerrors.DatabaseExecutionError,
			)
		}
		*courseIDs = append(*courseIDs, id)
	}
	return nil
}

func (cr *CourseRepositoryImpl) MentorListCourse(
	ctx context.Context,
	query *[]entity.MentorListCourseQuery,
	totalRow *int64,
	param entity.MentorListCourseParam,
) error {
	var driver RepoDriver
	driver = cr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	countArgs := []any{param.MentorID, param.LastID}
	args := []any{param.MentorID, param.LastID}
	sqlQuery := `
		SELECT
		c.id,
		c.title,
		c.domicile,
		c.price,
		c.method,
		c.session_duration_minutes,
		c.max_total_session,
		COALESCE(
			(SELECT STRING_AGG(cc_all.name, ',') 
			FROM course_category_assignments cca_all 
			JOIN course_categories cc_all ON cca_all.category_id = cc_all.id 
			WHERE cca_all.course_id = c.id 
			AND cca_all.deleted_at IS NULL 
			AND cc_all.deleted_at IS NULL), 
			''
		) AS categories
	FROM courses c
	WHERE c.mentor_id = $1 AND c.deleted_at IS NULL AND c.id < $2
	`
	countQuery := `
	SELECT count(*)
	FROM courses c
	WHERE c.mentor_id = $1 AND c.deleted_at IS NULL AND c.id < $2
	`
	if param.Search != nil {
		countQuery += fmt.Sprintf(" AND c.title ILIKE $%d ", len(countArgs)+1)
		sqlQuery += fmt.Sprintf(" AND c.title ILIKE $%d ", len(args)+1)
		args = append(args, "%"+*param.Search+"%")
		countArgs = append(countArgs, "%"+*param.Search+"%")
	}
	if param.CourseCategory != nil {
		countQuery += fmt.Sprintf(` AND EXISTS (
			SELECT 1 
			FROM course_category_assignments cca_filter 
			JOIN course_categories cc_filter ON cca_filter.category_id = cc_filter.id
			WHERE cca_filter.course_id = c.id 
			AND cca_filter.deleted_at IS NULL
			AND cc_filter.deleted_at IS NULL
			AND cc_filter.id = $%d
		)`, len(args)+1)
		sqlQuery += fmt.Sprintf(` AND EXISTS (
			SELECT 1 
			FROM course_category_assignments cca_filter 
			JOIN course_categories cc_filter ON cca_filter.category_id = cc_filter.id
			WHERE cca_filter.course_id = c.id 
			AND cca_filter.deleted_at IS NULL
			AND cc_filter.deleted_at IS NULL
			AND cc_filter.id = $%d
		)`, len(args)+1)
		countArgs = append(countArgs, *param.CourseCategory)
		args = append(args, *param.CourseCategory)
	}
	sqlQuery += " ORDER BY c.id DESC "
	sqlQuery += fmt.Sprintf(" LIMIT $%d ", len(args)+1)
	args = append(args, param.Limit)
	row := driver.QueryRow(countQuery, countArgs...)
	if err := row.Scan(totalRow); err != nil {
		return customerrors.NewError(
			"failed to get course list",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	rows, err := driver.Query(sqlQuery, args...)
	if err != nil {
		return customerrors.NewError(
			"failed to get courses list",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	defer rows.Close()

	for rows.Next() {
		var item entity.MentorListCourseQuery
		if err := rows.Scan(
			&item.ID,
			&item.Title,
			&item.Domicile,
			&item.Price,
			&item.Method,
			&item.SessionDuration,
			&item.MaxSession,
			&item.CourseCategories,
		); err != nil {
			return customerrors.NewError(
				"failed to get course list",
				err,
				customerrors.DatabaseExecutionError,
			)
		}
		*query = append(*query, item)
	}

	return nil
}

func (cr *CourseRepositoryImpl) FindByID(ctx context.Context, id int, course *entity.Course, lock bool) error {
	var driver RepoDriver
	driver = cr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	SELECT
		id,
		mentor_id,
		title,
		description,
		domicile,
		price,
		method,
		session_duration_minutes,
		max_total_session,
		transaction_count,
		created_at,
		updated_at,
		deleted_at
	FROM courses
	WHERE id = $1 and deleted_at IS NULL
	`
	if lock {
		query += `
		FOR UPDATE
		`
	}
	row := driver.QueryRow(query, id)
	if err := row.Scan(
		&course.ID,
		&course.MentorID,
		&course.Title,
		&course.Description,
		&course.Domicile,
		&course.Price,
		&course.Method,
		&course.SessionDuration,
		&course.MaxSession,
		&course.TransactionCount,
		&course.CreatedAt,
		&course.UpdatedAt,
		&course.DeletedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return customerrors.NewError(
				"course does not exist",
				err,
				customerrors.ItemNotExist,
			)
		}
		return customerrors.NewError(
			"failed to get course",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (cr *CourseRepositoryImpl) DeleteCourse(ctx context.Context, id int) error {
	var driver RepoDriver
	driver = cr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	UPDATE courses
	SET deleted_at = NOW(), updated_at = NOW()
	WHERE id = $1 AND deleted_at IS NULL
	`
	_, err := driver.Exec(query, id)
	if err != nil {
		return customerrors.NewError(
			"failed to delete course",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (cr *CourseRepositoryImpl) DeleteMentorCourse(ctx context.Context, mentorID string) error {
	var driver RepoDriver
	driver = cr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	UPDATE courses
	SET deleted_at = NOW(), updated_at = NOW()
	WHERE mentor_id = $1 AND deleted_at IS NULL
	`
	_, err := driver.Exec(query, mentorID)
	if err != nil {
		return customerrors.NewError(
			"failed to delete course",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (cr *CourseRepositoryImpl) CreateCourse(
	ctx context.Context,
	mentorID string,
	title string,
	description string,
	domicile string,
	price float64,
	sessionDuration int,
	maxSession int,
	method string,
	course *entity.Course,
) error {
	var driver RepoDriver
	driver = cr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	INSERT INTO courses (mentor_id, title, description, domicile, price, session_duration_minutes, max_total_session, method)
	VALUES
	($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING id
	`
	row := driver.QueryRow(query, mentorID, title, description, domicile, price, sessionDuration, maxSession, method)
	if err := row.Scan(&course.ID); err != nil {
		return customerrors.NewError(
			"failed to create course",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}
