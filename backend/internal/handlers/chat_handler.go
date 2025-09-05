package handlers

import (
	"encoding/json"
	"net/http"
	"privat-unmei/internal/constants"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/dtos"
	"privat-unmei/internal/entity"
	"privat-unmei/internal/logger"
	"privat-unmei/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ChatHandlerImpl struct {
	chs *services.ChatServiceImpl
	upg *websocket.Upgrader
	hub *entity.ChatHub
	lg  logger.CustomLogger
}

func CreateChatHandler(chs *services.ChatServiceImpl, upg *websocket.Upgrader, hub *entity.ChatHub, lg logger.CustomLogger) *ChatHandlerImpl {
	return &ChatHandlerImpl{chs, upg, hub, lg}
}

func (chh *ChatHandlerImpl) ConnectChatChannel(ctx *gin.Context) {
	claim, err := getAuthenticationPayload(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	chatroomID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Error(customerrors.NewError(
			"invalid chatroom",
			err,
			customerrors.InvalidAction,
		))
		return
	}
	param := entity.GetChatroomParam{
		ChatroomID: chatroomID,
		UserID:     claim.Subject,
		Role:       claim.Role,
	}
	if _, err := chh.chs.GetChatroom(ctx, param); err != nil {
		ctx.Error(err)
		return
	}

	conn, err := chh.upg.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.Error(customerrors.NewError(
			"failed to connect",
			err,
			customerrors.CommonErr,
		))
		return
	}
	client := &entity.ChatClient{
		Conn:       conn,
		Send:       make(chan []byte),
		UserID:     claim.Subject,
		ChatroomID: chatroomID,
		Logger:     chh.lg,
	}
	chh.hub.Register <- client
	go func() {
		client.Read(chh.hub)
	}()
	go client.Write()
}

func (chh *ChatHandlerImpl) CreateChatroom(ctx *gin.Context) {
	var req dtos.CreateChatroomReq
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.Error(err)
		return
	}
	claim, err := getAuthenticationPayload(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	param := entity.CreateChatroomParam{
		UserID:   claim.Subject,
		MentorID: req.MentorID,
	}
	chatroom, err := chh.chs.CreateChatroom(ctx, param)
	if err != nil {
		ctx.Error(err)
		return
	}
	res := dtos.ChatroomRes{
		ID:       chatroom.ID,
		UserID:   chatroom.UserID,
		MentorID: chatroom.MentorID,
	}
	ctx.JSON(http.StatusCreated, dtos.Response{
		Success: true,
		Data:    res,
	})
}

func (chh *ChatHandlerImpl) GetChatroom(ctx *gin.Context) {
	chatroomID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Error(customerrors.NewError(
			"invalid chatroom",
			err,
			customerrors.InvalidAction,
		))
		return
	}
	claim, err := getAuthenticationPayload(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	param := entity.GetChatroomParam{
		ChatroomID: chatroomID,
		UserID:     claim.Subject,
		Role:       claim.Role,
	}
	chatroom, err := chh.chs.GetChatroom(ctx, param)
	if err != nil {
		ctx.Error(err)
		return
	}
	res := dtos.ChatroomRes{
		UserID:   chatroom.UserID,
		MentorID: chatroom.MentorID,
		ID:       chatroom.ID,
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data:    res,
	})
}

func (chh *ChatHandlerImpl) GetUserChatrooms(ctx *gin.Context) {
	var req dtos.GetUserChatroomsReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(err)
		return
	}
	claim, err := getAuthenticationPayload(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	param := entity.GetUserChatroomsParam{
		PaginatedParam: entity.PaginatedParam{
			Limit: req.Limit,
			Page:  req.Page,
		},
		UserID: claim.Subject,
		Role:   claim.Role,
	}
	if param.Limit <= 0 {
		param.Limit = constants.DefaultLimit
	}
	if param.Page <= 0 {
		param.Page = constants.DefaultPage
	}
	chatrooms, totalRow, err := chh.chs.GetUserChatrooms(ctx, param)
	if err != nil {
		ctx.Error(err)
		return
	}
	res := []dtos.ChatroomRes{}
	for _, ch := range *chatrooms {
		res = append(res, dtos.ChatroomRes{
			ID:       ch.ID,
			UserID:   ch.UserID,
			MentorID: ch.MentorID,
		})
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.PaginatedResponse[dtos.ChatroomRes]{
			Entries: res,
			PageInfo: dtos.PaginatedInfo{
				Limit:    param.Limit,
				Page:     param.Page,
				TotalRow: *totalRow,
			},
		},
	})
}

func (chh *ChatHandlerImpl) SendMessage(ctx *gin.Context) {
	var req dtos.SendMessageReq
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.Error(err)
		return
	}
	chatroomID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Error(customerrors.NewError(
			"invalid chatroom",
			err,
			customerrors.InvalidAction,
		))
	}
	claim, err := getAuthenticationPayload(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	param := entity.SendMessageParam{
		ChatroomID: chatroomID,
		UserID:     claim.Subject,
		Content:    req.Content,
	}
	message, err := chh.chs.SendMessage(ctx, param)
	if err != nil {
		ctx.Error(err)
		return
	}
	messagePayload, err := json.Marshal(message)
	if err != nil {
		ctx.Error(customerrors.NewError(
			"failed to send message",
			err,
			customerrors.CommonErr,
		))
		return
	}
	chh.hub.BroadcastMessage(messagePayload, chatroomID)
	ctx.JSON(http.StatusCreated, dtos.Response{
		Success: true,
		Data: dtos.MessageRes{
			ID:         message.ID,
			SenderID:   message.SenderID,
			ChatroomID: message.ChatroomID,
			Content:    message.Content,
		},
	})
}

func (chh *ChatHandlerImpl) GetMessages(ctx *gin.Context) {
	var req dtos.GetMessagesReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(err)
		return
	}
	claim, err := getAuthenticationPayload(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	chatroomID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Error(customerrors.NewError(
			"invalid chatroom",
			err,
			customerrors.InvalidAction,
		))
		return
	}
	param := entity.GetMessagesParam{
		SeekPaginatedParam: entity.SeekPaginatedParam{
			Limit:  req.Limit,
			LastID: req.LastID,
		},
		UserID:     claim.Subject,
		ChatroomID: chatroomID,
		Role:       claim.Role,
	}
	if param.Limit <= 0 {
		param.Limit = constants.DefaultLimit
	}
	if param.LastID <= 0 {
		param.LastID = constants.DefaultLastID
	}
	messages, totalRow, err := chh.chs.GetMessages(ctx, param)
	if err != nil {
		ctx.Error(err)
		return
	}
	res := []dtos.MessageRes{}
	for _, msg := range *messages {
		res = append(res, dtos.MessageRes{
			ID:         msg.ID,
			SenderID:   msg.SenderID,
			ChatroomID: msg.ChatroomID,
			Content:    msg.Content,
		})
	}
	lastID := 0
	if len(res) > 0 {
		lastIndex := len(res) - 1
		lastID = res[lastIndex].ID
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.SeekPaginatedResponse[dtos.MessageRes]{
			Entries: res,
			PageInfo: dtos.SeekPaginatedInfo{
				LastID:   lastID,
				Limit:    param.Limit,
				TotalRow: *totalRow,
			},
		},
	})
}
