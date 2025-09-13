package handlers

import (
	"net/http"
	"privat-unmei/internal/dtos"
	"privat-unmei/internal/entity"
	"privat-unmei/internal/services"

	"github.com/gin-gonic/gin"
)

type AdditionalCostHandlerImpl struct {
	acs *services.AdditionalCostServiceImpl
}

func CreateAdditionalCostHandler(acs *services.AdditionalCostServiceImpl) *AdditionalCostHandlerImpl {
	return &AdditionalCostHandlerImpl{acs}
}

func (ach *AdditionalCostHandlerImpl) CreateNewAdditionalCost(ctx *gin.Context) {
	claim, err := getAuthenticationPayload(ctx)
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
