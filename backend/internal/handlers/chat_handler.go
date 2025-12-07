package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"privat-unmei/internal/constants"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/dtos"
	"privat-unmei/internal/entity"
	"privat-unmei/internal/logger"
	"privat-unmei/internal/services"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

type ChatHandlerImpl struct {
	chs *services.ChatServiceImpl
	upg *websocket.Upgrader
	rc  *redis.Client
	lg  logger.CustomLogger
}

func CreateChatHandler(chs *services.ChatServiceImpl, upg *websocket.Upgrader, rc *redis.Client, lg logger.CustomLogger) *ChatHandlerImpl {
	return &ChatHandlerImpl{chs, upg, rc, lg}
}

func (chh *ChatHandlerImpl) GetChatroom(ctx *gin.Context) {
	claim, err := getAuthenticationPayload(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	secondUserID := ctx.Param("id")
	param := entity.GetChatroomParam{
		UserID:       claim.Subject,
		SecondUserID: secondUserID,
		Role:         claim.Role,
	}
	chatroomID, err := chh.chs.GetChatroom(ctx, param)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.CreateChatroomRes{
			ID: chatroomID,
		},
	})

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
	param := entity.GetChatroomInfoParam{
		ChatroomID: chatroomID,
		UserID:     claim.Subject,
		Role:       claim.Role,
	}
	if _, err := chh.chs.GetChatroomInfo(ctx, param); err != nil {
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
	sub := chh.rc.Subscribe(context.Background(), fmt.Sprintf("chatroom:%d", chatroomID))

	var wg sync.WaitGroup
	cc := entity.CreateChatClient(conn, sub, chh.lg)
	defer cc.CloseConn()

	wg.Add(2)
	go cc.Read(&wg)
	go cc.Write(&wg)

	wg.Wait()
}

func (chh *ChatHandlerImpl) GetChatroomInfo(ctx *gin.Context) {
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
	param := entity.GetChatroomInfoParam{
		ChatroomID: chatroomID,
		UserID:     claim.Subject,
		Role:       claim.Role,
	}
	chatroom, err := chh.chs.GetChatroomInfo(ctx, param)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data:    dtos.ChatroomRes(*chatroom),
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
	if param.Limit <= 0 || param.Limit > constants.MaxLimit {
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
		res = append(res, dtos.ChatroomRes(ch))
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
		return
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
		Role:       claim.Role,
	}
	message, err := chh.chs.SendMessage(ctx, param)
	if err != nil {
		ctx.Error(err)
		return
	}
	res := dtos.MessageRes(*message)
	messagePayload, err := json.Marshal(res)
	if err != nil {
		ctx.Error(customerrors.NewError(
			"failed to send message",
			err,
			customerrors.CommonErr,
		))
		return
	}
	if err := chh.rc.Publish(ctx, fmt.Sprintf("chatroom:%d", chatroomID), messagePayload).Err(); err != nil {
		ctx.Error(customerrors.NewError(
			"something went wrong",
			err,
			customerrors.CommonErr,
		))
		return
	}
	ctx.JSON(http.StatusCreated, dtos.Response{
		Success: true,
		Data:    res,
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
	if param.Limit <= 0 || param.Limit > constants.MaxLimit {
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
		res = append(res, dtos.MessageRes(msg))
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
