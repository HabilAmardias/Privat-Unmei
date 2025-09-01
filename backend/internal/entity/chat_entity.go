package entity

import "time"

type (
	Chatroom struct {
		ID        int
		UserID    string
		MentorID  string
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt *time.Time
	}
	Message struct {
		ID         int
		SenderID   string
		ChatroomID int
		Content    string
		CreatedAt  time.Time
		UpdatedAt  time.Time
		DeletedAt  *time.Time
	}
	CreateChatroomParam struct {
		UserID   string
		MentorID string
	}
	GetChatroomParam struct {
		ChatroomID int
		UserID     string
	}
	GetUserChatroomsParam struct {
		PaginatedParam
		UserID string
	}
	SendMessageParam struct {
		ChatroomID int
		UserID     string
		Content    string
	}
	GetMessagesParam struct {
		SeekPaginatedParam
		UserID     string
		ChatroomID int
	}
)
