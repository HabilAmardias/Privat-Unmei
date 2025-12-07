package entity

import (
	"context"
	"privat-unmei/internal/logger"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
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
	GetChatroomParam struct {
		UserID       string
		SecondUserID string
		Role         int
	}
	GetChatroomInfoParam struct {
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
		cancelFunc context.CancelFunc
		cancelCtx  context.Context
		conn       *websocket.Conn
		lg         logger.CustomLogger
		sub        *redis.PubSub
		once       sync.Once
	}
)

func CreateChatClient(conn *websocket.Conn, sub *redis.PubSub, lg logger.CustomLogger) *ChatClient {
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})
	wsCtx, cancel := context.WithCancel(context.Background())

	return &ChatClient{
		cancelFunc: cancel,
		cancelCtx:  wsCtx,
		conn:       conn,
		lg:         lg,
		sub:        sub,
	}
}

func (cc *ChatClient) CloseConn() {
	cc.once.Do(func() {
		cc.cancelFunc()
		cc.sub.Close()
		cc.conn.Close()
	})
}

func (cc *ChatClient) Read(wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
		cc.CloseConn()
	}()
	for {
		_, _, err := cc.conn.ReadMessage()
		if err != nil {
			cc.lg.Errorln(err.Error())
			return
		}
	}
}

func (cc *ChatClient) Write(wg *sync.WaitGroup) {
	ticker := time.NewTicker(54 * time.Second)
	ch := cc.sub.Channel()
	defer func() {
		wg.Done()
		ticker.Stop()
		cc.CloseConn()
	}()
	for {
		select {
		case <-cc.cancelCtx.Done():
			cc.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		case message, ok := <-ch:
			if !ok {
				cc.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := cc.conn.WriteMessage(websocket.TextMessage, []byte(message.Payload)); err != nil {
				cc.lg.Errorln(err.Error())
				return
			}
		case <-ticker.C:
			if err := cc.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				cc.lg.Errorln(err.Error())
				return
			}
		}
	}
}
