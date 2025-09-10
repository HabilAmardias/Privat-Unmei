package handlers

import (
	"net/http"
	"privat-unmei/internal/dtos"
	"privat-unmei/internal/entity"
	"privat-unmei/internal/services"

	"github.com/gin-gonic/gin"
)

// TODO: Add Payment Method CRUD feature for admin
type PaymentHandlerImpl struct {
	ps *services.PaymentServiceImpl
}

func CreatePaymentHandler(ps *services.PaymentServiceImpl) *PaymentHandlerImpl {
	return &PaymentHandlerImpl{ps}
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
