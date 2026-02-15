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

type DiscountHandlerImpl struct {
	ds *services.DiscountServiceImpl
}

func CreateDiscountHandler(ds *services.DiscountServiceImpl) *DiscountHandlerImpl {
	return &DiscountHandlerImpl{ds}
}

func (dh *DiscountHandlerImpl) GetDiscount(ctx *gin.Context) {
	participantStr := ctx.Param("participant")
	participant, err := strconv.Atoi(participantStr)
	if err != nil {
		ctx.Error(customerrors.NewError(
			"invalid data",
			err,
			customerrors.InvalidAction,
		))
		return
	}
	claim, err := getAuthenticationPayload(ctx, constants.CTX_AUTH_PAYLOAD_KEY)
	if err != nil {
		ctx.Error(err)
		return
	}
	param := entity.GetDiscountParam{
		Participant: participant,
		UserID:      claim.Subject,
	}
	discount, err := dh.ds.GetDiscount(ctx, param)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.GetDiscountRes{
			Amount: *discount,
		},
	})
}

func (dh *DiscountHandlerImpl) GetAllDiscount(ctx *gin.Context) {
	var req dtos.GetAllDiscountReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(err)
		return
	}
	claim, err := getAuthenticationPayload(ctx, constants.CTX_AUTH_PAYLOAD_KEY)
	if err != nil {
		ctx.Error(err)
		return
	}
	param := entity.GetAllDiscountParam{
		AdminID: claim.Subject,
		PaginatedParam: entity.PaginatedParam{
			Limit: req.Limit,
			Page:  req.Page,
		},
	}
	if param.Limit <= 0 || param.Limit > constants.MaxLimit {
		param.Limit = constants.DefaultLimit
	}
	if param.Page <= 0 {
		param.Page = constants.DefaultPage
	}
	discounts, totalRow, err := dh.ds.GetAllDiscount(ctx, param)
	if err != nil {
		ctx.Error(err)
		return
	}
	entries := []dtos.GetAllDiscountRes{}
	for _, d := range *discounts {
		entries = append(entries, dtos.GetAllDiscountRes(d))
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.PaginatedResponse[dtos.GetAllDiscountRes]{
			Entries: entries,
			PageInfo: dtos.PaginatedInfo{
				Limit:    param.Limit,
				Page:     param.Page,
				TotalRow: *totalRow,
			},
		},
	})
}

func (dh *DiscountHandlerImpl) DeleteDiscount(ctx *gin.Context) {
	discountID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Error(customerrors.NewError(
			"invalid discount",
			err,
			customerrors.InvalidAction,
		))
		return
	}
	claim, err := getAuthenticationPayload(ctx, constants.CTX_AUTH_PAYLOAD_KEY)
	if err != nil {
		ctx.Error(err)
		return
	}
	param := entity.DeleteDiscountParam{
		AdminID:    claim.Subject,
		DiscountID: discountID,
	}
	if err := dh.ds.DeleteDiscount(ctx, param); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.DiscountIDRes{
			ID: discountID,
		},
	})
}

func (dh *DiscountHandlerImpl) UpdateDiscountAmount(ctx *gin.Context) {
	discountID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Error(customerrors.NewError(
			"invalid discount",
			err,
			customerrors.InvalidAction,
		))
		return
	}
	claim, err := getAuthenticationPayload(ctx, constants.CTX_AUTH_PAYLOAD_KEY)
	if err != nil {
		ctx.Error(err)
		return
	}
	var req dtos.UpdateDiscountReq
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.Error(err)
		return
	}
	param := entity.UpdateDiscountParam{
		AdminID:    claim.Subject,
		DiscountID: discountID,
		Amount:     req.Amount,
	}
	if err := dh.ds.UpdateDiscountAmount(ctx, param); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.DiscountIDRes{
			ID: discountID,
		},
	})
}

func (dh *DiscountHandlerImpl) CreateNewDiscount(ctx *gin.Context) {
	var req dtos.CreateNewDiscountReq
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.Error(err)
		return
	}
	claim, err := getAuthenticationPayload(ctx, constants.CTX_AUTH_PAYLOAD_KEY)
	if err != nil {
		ctx.Error(err)
		return
	}
	param := entity.CreateNewDiscountParam{
		AdminID:             claim.Subject,
		NumberOfParticipant: req.NumberOfParticipant,
		Amount:              req.Amount,
	}
	id, err := dh.ds.CreateNewDiscount(ctx, param)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusCreated, dtos.Response{
		Success: true,
		Data: dtos.DiscountIDRes{
			ID: id,
		},
	})
}
