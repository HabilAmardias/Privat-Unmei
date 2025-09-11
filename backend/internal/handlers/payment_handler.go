package handlers

import (
	"net/http"
	"privat-unmei/internal/constants"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/dtos"
	"privat-unmei/internal/entity"
	"privat-unmei/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

// TODO: Add Payment Method CRUD feature for admin
type PaymentHandlerImpl struct {
	ps *services.PaymentServiceImpl
}

func CreatePaymentHandler(ps *services.PaymentServiceImpl) *PaymentHandlerImpl {
	return &PaymentHandlerImpl{ps}
}

func (ph *PaymentHandlerImpl) GetAllPaymentMethod(ctx *gin.Context) {
	var req dtos.GetAllPaymentMethodReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(err)
		return
	}
	param := entity.GetAllPaymentMethodParam{
		SeekPaginatedParam: entity.SeekPaginatedParam{
			Limit:  req.Limit,
			LastID: req.LastID,
		},
		Search: req.Search,
	}
	if param.Limit <= 0 || param.Limit > constants.MaxLimit {
		param.Limit = constants.DefaultLimit
	}
	if param.LastID <= 0 {
		param.LastID = constants.DefaultLastID
	}
	methods, totalRow, err := ph.ps.GetAllPaymentMethod(ctx, param)
	if err != nil {
		ctx.Error(err)
		return
	}
	entries := []dtos.GetPaymentMethodRes{}
	for _, method := range *methods {
		entries = append(entries, dtos.GetPaymentMethodRes(method))
	}
	var filter []dtos.FilterInfo
	if req.Search != nil {
		filter = append(filter, dtos.FilterInfo{
			Name:  "Search",
			Value: *req.Search,
		})
	}
	lastID := constants.DefaultLastID
	if len(entries) > 0 {
		lastID = entries[len(entries)-1].ID
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.SeekPaginatedResponse[dtos.GetPaymentMethodRes]{
			Entries: entries,
			PageInfo: dtos.SeekPaginatedInfo{
				FilterBy: filter,
				TotalRow: *totalRow,
				Limit:    param.Limit,
				LastID:   lastID,
			},
		},
	})
}

func (ph *PaymentHandlerImpl) UpdatePaymentMethod(ctx *gin.Context) {
	claim, err := getAuthenticationPayload(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	methodID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Error(err)
		return
	}
	var req dtos.UpdatePaymentMethodReq
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.Error(err)
		return
	}
	param := entity.UpdatePaymentMethodParam{
		AdminID:       claim.Subject,
		MethodID:      methodID,
		MethodNewName: req.MethodNewName,
	}
	if err := ph.ps.UpdatePaymentMethod(ctx, param); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.UpdatePaymentMethodRes{
			ID: param.MethodID,
		},
	})
}

func (ph *PaymentHandlerImpl) DeletePaymentMethod(ctx *gin.Context) {
	claim, err := getAuthenticationPayload(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	paymentMethodID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Error(customerrors.NewError(
			"invalid payment method",
			err,
			customerrors.InvalidAction,
		))
		return
	}
	param := entity.DeletePaymentMethodParam{
		AdminID:  claim.Subject,
		MethodID: paymentMethodID,
	}
	if err := ph.ps.DeletePaymentMethod(ctx, param); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.DeletePaymentMethodRes{
			ID: param.MethodID,
		},
	})
}

func (ph *PaymentHandlerImpl) CreatePaymentMethod(ctx *gin.Context) {
	claim, err := getAuthenticationPayload(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	var req dtos.CreatePaymentMethodReq
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.Error(err)
		return
	}
	param := entity.CreatePaymentMethodParam{
		AdminID:    claim.Subject,
		MethodName: req.Name,
	}
	id, err := ph.ps.CreatePaymentMethod(ctx, param)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusCreated, dtos.Response{
		Success: true,
		Data: dtos.CreatePaymentMethodRes{
			ID: *id,
		},
	})
}
