package dtos

type (
	CreateChatroomReq struct {
		MentorID string `json:"mentor_id" binding:"required"`
	}
	GetUserChatroomsReq struct {
		PaginatedReq
	}
	SendMessageReq struct {
		Content string `json:"content" binding:"required,max=180"`
	}
	GetMessagesReq struct {
		SeekPaginatedReq
	}
	CreateChatroomRes struct {
		ID string `json:"id"`
	}
	ChatroomRes struct {
		ID               string `json:"id"`
		UserID           string `json:"user_id"`
		Username         string `json:"username"`
		UserPublicID     string `json:"public_id"`
		UserProfileImage string `json:"profile_image"`
	}
	MessageRes struct {
		ID             int    `json:"id"`
		SenderID       string `json:"sender_id"`
		SenderName     string `json:"sender_name"`
		SenderPublicID string `json:"sender_public_id"`
		ChatroomID     string `json:"chatroom_id"`
		Content        string `json:"content"`
	}
)
