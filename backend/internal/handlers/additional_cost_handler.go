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

type AdditionalCostHandlerImpl struct {
	acs *services.AdditionalCostServiceImpl
}

func CreateAdditionalCostHandler(acs *services.AdditionalCostServiceImpl) *AdditionalCostHandlerImpl {
	return &AdditionalCostHandlerImpl{acs}
}

func (ach *AdditionalCostHandlerImpl) GetOperationalCost(ctx *gin.Context) {
	claim, err := getAuthenticationPayload(ctx, constants.CTX_AUTH_PAYLOAD_KEY)
	if err != nil {
		ctx.Error(err)
		return
	}
	param := entity.GetOperationalCostParam{
		UserID: claim.Subject,
	}
	cost, err := ach.acs.GetOperationalCost(ctx, param)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.OperationalCostRes{
			Cost: *cost,
		},
	})
}

func (ach *AdditionalCostHandlerImpl) GetAllAdditionalCost(ctx *gin.Context) {
	var req dtos.GetAllAdditionalCostReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(err)
		return
	}
	claim, err := getAuthenticationPayload(ctx, constants.CTX_AUTH_PAYLOAD_KEY)
	if err != nil {
		ctx.Error(err)
		return
	}
	param := entity.GetAllAdditionalCostParam{
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
	discounts, totalRow, err := ach.acs.GetAllAdditionalCost(ctx, param)
	if err != nil {
		ctx.Error(err)
		return
	}
	entries := []dtos.GetAdditionalCostRes{}
	for _, d := range *discounts {
		entries = append(entries, dtos.GetAdditionalCostRes(d))
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.PaginatedResponse[dtos.GetAdditionalCostRes]{
			Entries: entries,
			PageInfo: dtos.PaginatedInfo{
				Limit:    param.Limit,
				Page:     param.Page,
				TotalRow: *totalRow,
			},
		},
	})
}

func (ach *AdditionalCostHandlerImpl) DeleteCost(ctx *gin.Context) {
	claim, err := getAuthenticationPayload(ctx, constants.CTX_AUTH_PAYLOAD_KEY)
	if err != nil {
		ctx.Error(err)
		return
	}
	costID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Error(customerrors.NewError(
			"invalid cost",
			err,
			customerrors.InvalidAction,
		))
		return
	}
	param := entity.DeleteAdditionalCostParam{
		CostID:  costID,
		AdminID: claim.Subject,
	}
	if err := ach.acs.DeleteCost(ctx, param); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.AdditionalCostIDRes{
			ID: costID,
		},
	})
}

func (ach *AdditionalCostHandlerImpl) UpdateCostAmount(ctx *gin.Context) {
	claim, err := getAuthenticationPayload(ctx, constants.CTX_AUTH_PAYLOAD_KEY)
	if err != nil {
		ctx.Error(err)
		return
	}
	costID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Error(customerrors.NewError(
			"invalid cost",
			err,
			customerrors.InvalidAction,
		))
		return
	}
	var req dtos.UpdateAdditionalCostReq
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.Error(err)
		return
	}
	param := entity.UpdateAdditonalCostParam{
		AdminID: claim.Subject,
		CostID:  costID,
		Amount:  req.Amount,
	}
	if err := ach.acs.UpdateCostAmount(ctx, param); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.AdditionalCostIDRes{
			ID: param.CostID,
		},
	})
}

func (ach *AdditionalCostHandlerImpl) CreateNewAdditionalCost(ctx *gin.Context) {
	claim, err := getAuthenticationPayload(ctx, constants.CTX_AUTH_PAYLOAD_KEY)
	if err != nil {
		ctx.Error(err)
		return
	}
	var req dtos.CreateAdditionalCostReq
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.Error(err)
		return
	}
	param := entity.CreateAdditionalCostParam{
		AdminID: claim.Subject,
		Name:    req.Name,
		Amount:  req.Amount,
	}
	id, err := ach.acs.CreateNewAdditionalCost(ctx, param)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusCreated, dtos.Response{
		Success: true,
		Data: dtos.AdditionalCostIDRes{
			ID: id,
		},
	})
}
