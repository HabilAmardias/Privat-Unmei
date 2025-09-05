package entity

import (
	"privat-unmei/internal/logger"
	"time"

	"github.com/gorilla/websocket"
)

type (
	Chatroom struct {
		ID        int
		StudentID string
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
	MessageDetailQuery struct {
		ID          int
		SenderID    string
		SenderName  string
		SenderEmail string
		ChatroomID  int
		Content     string
	}
	ChatroomDetailQuery struct {
		ID               int
		UserID           string
		Username         string
		UserEmail        string
		UserProfileImage string
	}
	CreateChatroomParam struct {
		StudentID string
		MentorID  string
	}
	GetChatroomParam struct {
		ChatroomID int
		UserID     string
		Role       int
	}
	GetUserChatroomsParam struct {
		PaginatedParam
		UserID string
		Role   int
	}
	SendMessageParam struct {
		ChatroomID int
		UserID     string
		Content    string
		Role       int
	}
	GetMessagesParam struct {
		SeekPaginatedParam
		UserID     string
		ChatroomID int
		Role       int
	}
	ChatClient struct {
		Conn       *websocket.Conn
		Send       chan []byte
		UserID     string
		ChatroomID int
		Logger     logger.CustomLogger
	}
	ChatHub struct {
		Clients    map[*ChatClient]bool
		Unregister chan *ChatClient
		Register   chan *ChatClient
		Broadcast  chan []byte
		Chatrooms  map[int]map[*ChatClient]bool
	}
)

func CreateChatHub() *ChatHub {
	return &ChatHub{
		Clients:    make(map[*ChatClient]bool),
		Unregister: make(chan *ChatClient),
		Register:   make(chan *ChatClient),
		Broadcast:  make(chan []byte),
		Chatrooms:  make(map[int]map[*ChatClient]bool),
	}
}

func (c *ChatClient) Write() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				c.Logger.Warnln(err.Error())
			}
		case <-ticker.C:
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.Logger.Errorln(err.Error())
				return
			}

		}
	}
}

func (c *ChatClient) Read(hub *ChatHub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()
	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			c.Logger.Errorln(err.Error())
			break
		}
	}
}

func (h *ChatHub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
			if h.Chatrooms[client.ChatroomID] == nil {
				h.Chatrooms[client.ChatroomID] = make(map[*ChatClient]bool)
			}
			h.Chatrooms[client.ChatroomID][client] = true
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				delete(h.Chatrooms[client.ChatroomID], client)
				if len(h.Chatrooms[client.ChatroomID]) == 0 {
					delete(h.Chatrooms, client.ChatroomID)
				}
				close(client.Send)
			}
		}
	}
}

func (h *ChatHub) BroadcastMessage(message []byte, chatroomID int) {
	if clients, exist := h.Chatrooms[chatroomID]; exist {
		for client := range clients {
			select {
			case client.Send <- message:
			default:
				close(client.Send)
				delete(h.Clients, client)
				delete(clients, client)
			}
		}
	}
}
