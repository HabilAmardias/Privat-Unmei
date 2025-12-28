package repositories

import (
	"context"
	"database/sql"
	"errors"
	"privat-unmei/internal/constants"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/db"
	"privat-unmei/internal/entity"
)

type ChatRepositoryImpl struct {
	DB *db.CustomDB
}

func CreateChatRepository(db *db.CustomDB) *ChatRepositoryImpl {
	return &ChatRepositoryImpl{db}
}

func (chr *ChatRepositoryImpl) UpdateMentorLastRead(ctx context.Context, chatroomID string) error {
	var driver RepoDriver = chr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	UPDATE chatrooms
	SET
		last_read_by_mentor = CURRENT_TIMESTAMP
	WHERE id = $1 AND deleted_at IS NULL
	`
	_, err := driver.Exec(query, chatroomID)
	if err != nil {
		return customerrors.NewError(
			"something went wrong",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (chr *ChatRepositoryImpl) UpdateStudentLastRead(ctx context.Context, chatroomID string) error {
	var driver RepoDriver = chr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	UPDATE chatrooms
	SET
		last_read_by_student = CURRENT_TIMESTAMP
	WHERE id = $1 AND deleted_at IS NULL
	`
	_, err := driver.Exec(query, chatroomID)
	if err != nil {
		return customerrors.NewError(
			"something went wrong",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (chr *ChatRepositoryImpl) CreateChatroom(ctx context.Context, mentorID string, userID string, chatroom *entity.Chatroom) error {
	var driver RepoDriver
	driver = chr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}

	query := `
	INSERT INTO chatrooms (student_id, mentor_id)
	VALUES
	($1, $2)
	RETURNING id, student_id, mentor_id, created_at, updated_at, deleted_at
	`
	if err := driver.QueryRow(query, userID, mentorID).Scan(
		&chatroom.ID,
		&chatroom.StudentID,
		&chatroom.MentorID,
		&chatroom.CreatedAt,
		&chatroom.UpdatedAt,
		&chatroom.DeletedAt,
	); err != nil {
		return customerrors.NewError(
			"failed to create chatroom",
			err,
			customerrors.DatabaseExecutionError,
		)
	}

	return nil
}

func (chr *ChatRepositoryImpl) GetChatroom(ctx context.Context, mentorID string, userID string, chatroom *entity.Chatroom) error {
	var driver RepoDriver
	driver = chr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	SELECT
		id,
		student_id,
		mentor_id,
		created_at,
		updated_at,
		deleted_at
	FROM chatrooms
	WHERE student_id = $1 AND mentor_id = $2 AND deleted_at IS NULL
	`
	if err := driver.QueryRow(query, userID, mentorID).Scan(
		&chatroom.ID,
		&chatroom.StudentID,
		&chatroom.MentorID,
		&chatroom.CreatedAt,
		&chatroom.UpdatedAt,
		&chatroom.DeletedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return customerrors.NewError(
				"chatroom not found",
				err,
				customerrors.ItemNotExist,
			)
		}
		return customerrors.NewError(
			"failed to get chatroom",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (chr *ChatRepositoryImpl) GetUserChatrooms(ctx context.Context, userID string, role int, limit int, page int, totalRow *int64, chatrooms *[]entity.ChatroomDetailQuery) error {
	var driver RepoDriver
	driver = chr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	countQuery := `
	SELECT
		count(*)
	FROM chatrooms ch
	JOIN users u ON u.id = ch.mentor_id
	WHERE ch.student_id = $1 AND ch.deleted_at IS NULL AND u.deleted_at IS NULL
	`
	if role == constants.MentorRole {
		countQuery = `
		SELECT
			count(*)
		FROM chatrooms ch
		JOIN users u ON u.id = ch.student_id
		WHERE ch.mentor_id = $1 AND ch.deleted_at IS NULL AND u.deleted_at IS NULL
	`
	}
	if err := driver.QueryRow(countQuery, userID).Scan(totalRow); err != nil {
		return customerrors.NewError(
			"failed to get chatrooms",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	query := `
	SELECT
		ch.id,
		ch.mentor_id,
		u.name,
		u.public_id,
		u.profile_image,
		COALESCE(COUNT(m.id), 0) as unread_count
	FROM chatrooms ch
	JOIN users u ON u.id = ch.mentor_id
	LEFT JOIN messages m ON 
		m.chatroom_id = ch.id 
		AND m.created_at > ch.last_read_by_student 
		AND m.deleted_at IS NULL
	WHERE 
		ch.student_id = $1 
		AND ch.deleted_at IS NULL 
		AND u.deleted_at IS NULL 
	GROUP BY ch.id, ch.mentor_id, u.name, u.public_id, u.profile_image
	ORDER BY ch.updated_at DESC
	LIMIT $2
	OFFSET $3
	`
	if role == constants.MentorRole {
		query = `
		SELECT
			ch.id,
			ch.student_id,
			u.name,
			u.public_id,
			u.profile_image,
			COALESCE(COUNT(m.id), 0) as unread_count
		FROM chatrooms ch
		JOIN users u ON u.id = ch.student_id
		LEFT JOIN messages m ON 
			m.chatroom_id = ch.id 
			AND m.created_at > ch.last_read_by_mentor 
			AND m.deleted_at IS NULL
		WHERE 
			ch.mentor_id = $1 
			AND ch.deleted_at IS NULL 
			AND u.deleted_at IS NULL 
			
		GROUP BY ch.id, ch.student_id, u.name, u.public_id, u.profile_image
		ORDER BY ch.updated_at DESC
		LIMIT $2
		OFFSET $3
		`
	}
	rows, err := driver.Query(query, userID, limit, limit*(page-1))
	if err != nil {
		return customerrors.NewError(
			"failed to get chatrooms",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	defer rows.Close()
	for rows.Next() {
		var item entity.ChatroomDetailQuery
		if err := rows.Scan(
			&item.ID,
			&item.UserID,
			&item.Username,
			&item.UserPublicID,
			&item.UserProfileImage,
			&item.UnreadCount,
		); err != nil {
			return customerrors.NewError(
				"failed to get chatrooms",
				err,
				customerrors.DatabaseExecutionError,
			)
		}
		*chatrooms = append(*chatrooms, item)
	}
	return nil
}

func (chr *ChatRepositoryImpl) SendMessage(ctx context.Context, senderID string, chatroomID string, content string, message *entity.Message) error {
	var driver RepoDriver
	driver = chr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	INSERT into messages (sender_id, chatroom_id, content)
	VALUES
	($1, $2, $3)
	RETURNING id, sender_id, chatroom_id, content, created_at, updated_at, deleted_at
	`

	if err := driver.QueryRow(query, senderID, chatroomID, content).Scan(
		&message.ID,
		&message.SenderID,
		&message.ChatroomID,
		&message.Content,
		&message.CreatedAt,
		&message.UpdatedAt,
		&message.DeletedAt,
	); err != nil {
		return customerrors.NewError(
			"failed to send message",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (chr *ChatRepositoryImpl) UpdateChatroom(ctx context.Context, chatroomID string) error {
	var driver RepoDriver
	driver = chr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	UPDATE chatrooms 
	SET updated_at = CURRENT_TIMESTAMP 
	WHERE id = $1 AND deleted_at IS NULL
	`
	_, err := driver.Exec(query, chatroomID)
	if err != nil {
		return customerrors.NewError(
			"failed to send messages",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}

func (chr *ChatRepositoryImpl) GetMessages(ctx context.Context, chatroomID string, limit int, lastID *int, totalRow *int64, messages *[]entity.MessageDetailQuery) error {
	var driver RepoDriver
	driver = chr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	countQuery := `
	SELECT count(*)
	FROM messages m
	JOIN users u ON u.id = m.sender_id
	WHERE m.chatroom_id = $1 AND m.deleted_at IS NULL AND u.deleted_at IS NULL
	`
	if err := driver.QueryRow(countQuery, chatroomID).Scan(totalRow); err != nil {
		return customerrors.NewError(
			"failed to retrieve messages",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	query := `
	with chats as (
		SELECT
			m.id as id,
			m.sender_id as sender_id,
			u.name as sender_name,
			u.public_id as sender_public_id,
			m.chatroom_id as chatroom_id,
			m.content as content
		FROM messages m
		JOIN users u ON u.id = m.sender_id
		WHERE m.chatroom_id = $1 AND m.id < COALESCE($2, 9223372036854775807) AND m.deleted_at IS NULL AND u.deleted_at IS NULL
		ORDER BY m.id DESC
		LIMIT $3
	)
	select * from chats
	order by chats.id asc;
	`
	rows, err := driver.Query(query, chatroomID, lastID, limit)
	if err != nil {
		return customerrors.NewError(
			"failed to retrieve messages",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	defer rows.Close()
	for rows.Next() {
		var item entity.MessageDetailQuery
		if err := rows.Scan(
			&item.ID,
			&item.SenderID,
			&item.SenderName,
			&item.SenderPublicID,
			&item.ChatroomID,
			&item.Content,
		); err != nil {
			return customerrors.NewError(
				"failed to retrieve messages",
				err,
				customerrors.DatabaseExecutionError,
			)
		}
		*messages = append(*messages, item)
	}
	return nil
}

func (chr *ChatRepositoryImpl) FindByID(ctx context.Context, chatroomID string, chatroom *entity.Chatroom) error {
	var driver RepoDriver
	driver = chr.DB
	if tx := GetTransactionFromContext(ctx); tx != nil {
		driver = tx
	}
	query := `
	SELECT
		id,
		student_id,
		mentor_id,
		created_at,
		updated_at,
		deleted_at
	FROM chatrooms
	WHERE id = $1 AND deleted_at IS NULL
	`
	if err := driver.QueryRow(query, chatroomID).Scan(
		&chatroom.ID,
		&chatroom.StudentID,
		&chatroom.MentorID,
		&chatroom.CreatedAt,
		&chatroom.UpdatedAt,
		&chatroom.DeletedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return customerrors.NewError(
				"chatroom not found",
				err,
				customerrors.ItemNotExist,
			)
		}
		return customerrors.NewError(
			"failed to get chatroom",
			err,
			customerrors.DatabaseExecutionError,
		)
	}
	return nil
}
