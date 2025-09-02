package dtos

type (
	CreateChatroomReq struct {
		MentorID string `json:"mentor_id" binding:"required"`
	}
	GetUserChatroomsReq struct {
		PaginatedReq
	}
	SendMessageReq struct {
		Content string `json:"content" binding:"required,max=50"`
	}
	GetMessagesReq struct {
		SeekPaginatedReq
	}
	ChatroomRes struct {
		ID       int    `json:"id"`
		UserID   string `json:"user_id"`
		MentorID string `json:"mentor_id"`
	}
	MessageRes struct {
		ID         int    `json:"id"`
		SenderID   string `json:"sender_id"`
		ChatroomID int    `json:"chatroom_id"`
		Content    string `json:"content"`
	}
)
