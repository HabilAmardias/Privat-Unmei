package handlers

import (
	"errors"
	"net/http"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/dtos"
	"privat-unmei/internal/entity"
	"privat-unmei/internal/services"

	"github.com/gin-gonic/gin"
)

type CourseHandlerImpl struct {
	cs *services.CourseServiceImpl
}

func CreateCourseHandler(cs *services.CourseServiceImpl) *CourseHandlerImpl {
	return &CourseHandlerImpl{cs}
}

func (ch *CourseHandlerImpl) AddNewCourse(ctx *gin.Context) {
	claim, err := getAuthenticationPayload(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	var req dtos.CreateCourseReq
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.Error(err)
		return
	}
	if len(req.CourseAvailability) <= 0 {
		ctx.Error(customerrors.NewError(
			"mentor should insert their schedule",
			errors.New("mentor should insert their schedule"),
			customerrors.InvalidAction,
		))
		return
	}
	if len(req.Topics) <= 0 {
		ctx.Error(customerrors.NewError(
			"mentor should enter course topic at least one",
			errors.New("mentor should enter course topic at least one"),
			customerrors.InvalidAction,
		))
		return
	}
	param := entity.CreateCourseParam{
		MentorID:           claim.Subject,
		Title:              req.Title,
		Description:        req.Description,
		Domicile:           req.Domicile,
		MinPrice:           req.MinPrice,
		MaxPrice:           req.MaxPrice,
		Method:             req.Method,
		MinDuration:        req.MinDuration,
		MaxDuration:        req.MaxDuration,
		CourseAvailability: []entity.CreateSchedule{},
		Topics:             []entity.CreateTopic{},
	}
	for _, topic := range req.Topics {
		param.Topics = append(param.Topics, entity.CreateTopic(topic))
	}
	for _, schedule := range req.CourseAvailability {
		param.CourseAvailability = append(
			param.CourseAvailability,
			entity.CreateSchedule{
				DayOfWeek: schedule.DayOfWeek,
				StartTime: entity.TimeOnly(schedule.StartTime),
				EndTime:   entity.TimeOnly(schedule.EndTime),
			},
		)
	}
	courseID, err := ch.cs.CreateCourse(ctx, param)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusCreated, dtos.Response{
		Success: true,
		Data: dtos.CreateCourseRes{
			ID: courseID,
		},
	})
}
