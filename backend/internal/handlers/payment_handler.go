package handlers

import (
	"net/http"
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
