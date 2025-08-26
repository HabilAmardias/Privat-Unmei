package handlers

import (
	"errors"
	"net/http"
	"privat-unmei/internal/constants"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/dtos"
	"privat-unmei/internal/entity"
	"privat-unmei/internal/services"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CourseRequestHandlerImpl struct {
	cos *services.CourseRequestServiceImpl
}

func CreateCourseRequestHandler(cos *services.CourseRequestServiceImpl) *CourseRequestHandlerImpl {
	return &CourseRequestHandlerImpl{cos}
}

func (crh *CourseRequestHandlerImpl) ConfirmPayment(ctx *gin.Context) {
	claim, err := getAuthenticationPayload(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	courseRequestID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Error(customerrors.NewError(
			"invalid course request",
			err,
			customerrors.InvalidAction,
		))
		return
	}
	param := entity.ConfirmPaymentParam{
		MentorID:        claim.Subject,
		CourseRequestID: courseRequestID,
	}
	if err := crh.cos.MentorConfirmPayment(ctx, param); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.CreateCourseRequestRes{
			CourseRequestID: param.CourseRequestID,
		},
	})
}

func (crh *CourseRequestHandlerImpl) RejectCourseRequest(ctx *gin.Context) {
	claim, err := getAuthenticationPayload(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	courseRequestID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Error(customerrors.NewError(
			"invalid course request",
			err,
			customerrors.InvalidAction,
		))
		return
	}
	param := entity.HandleCourseRequestParam{
		MentorID:        claim.Subject,
		CourseRequestID: courseRequestID,
		Accept:          false,
	}
	if err := crh.cos.HandleCourseRequest(ctx, param); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.CreateCourseRequestRes{
			CourseRequestID: param.CourseRequestID,
		},
	})
}

func (crh *CourseRequestHandlerImpl) AcceptCourseRequest(ctx *gin.Context) {
	claim, err := getAuthenticationPayload(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	courseRequestID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Error(customerrors.NewError(
			"invalid course request",
			err,
			customerrors.InvalidAction,
		))
		return
	}
	param := entity.HandleCourseRequestParam{
		MentorID:        claim.Subject,
		CourseRequestID: courseRequestID,
		Accept:          true,
	}
	if err := crh.cos.HandleCourseRequest(ctx, param); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.CreateCourseRequestRes{
			CourseRequestID: param.CourseRequestID,
		},
	})
}

func (crh *CourseRequestHandlerImpl) CreateReservation(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Error(customerrors.NewError(
			"invalid course",
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
	var req dtos.CreateCourseRequstReq
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.Error(err)
		return
	}
	if len(req.PreferredSlots) <= 0 {
		ctx.Error(customerrors.NewError(
			"need to enter preferred schedule to create order",
			errors.New("no preferred slots"),
			customerrors.InvalidAction,
		))
		return
	}
	if len(req.PreferredSlots) > constants.MaxRequestSlot {
		ctx.Error(customerrors.NewError(
			"can only reserve up to seven session in one order",
			errors.New("there are more than 7 requested slots"),
			customerrors.InvalidAction,
		))
		return
	}
	if err := CheckDateUniqueness(req.PreferredSlots); err != nil {
		ctx.Error(err)
		return
	}
	param := entity.CreateCourseRequestParam{
		CourseID:       id,
		StudentID:      claim.Subject,
		PreferredSlots: []entity.PreferredSlot{},
	}
	for _, slot := range req.PreferredSlots {
		parsedDate, err := time.Parse("2006-01-02", slot.Date)
		if err != nil {
			ctx.Error(customerrors.NewError(
				"invalid date",
				err,
				customerrors.InvalidAction,
			))
			return
		}
		if !ValidateDate(parsedDate) {
			ctx.Error(customerrors.NewError(
				"invalid date",
				errors.New("invalid date"),
				customerrors.InvalidAction,
			))
			return
		}
		param.PreferredSlots = append(param.PreferredSlots, entity.PreferredSlot{
			Date:      parsedDate,
			StartTime: entity.TimeOnly(slot.StartTime),
		})
	}
	courseRequestID, err := crh.cos.CreateReservation(ctx, param)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusCreated, dtos.Response{
		Success: true,
		Data: dtos.CreateCourseRequestRes{
			CourseRequestID: courseRequestID,
		},
	})
}
