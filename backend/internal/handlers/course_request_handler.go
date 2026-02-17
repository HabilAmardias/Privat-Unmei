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
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type CourseRequestHandlerImpl struct {
	cos *services.CourseRequestServiceImpl
}

func CreateCourseRequestHandler(cos *services.CourseRequestServiceImpl) *CourseRequestHandlerImpl {
	return &CourseRequestHandlerImpl{cos}
}

func (crh *CourseRequestHandlerImpl) GetIncomeHistoryReport(ctx *gin.Context) {
	claim, err := getAuthenticationPayload(ctx, constants.CTX_AUTH_PAYLOAD_KEY)
	if err != nil {
		ctx.Error(err)
		return
	}
	param := entity.GetAdminMonthlyReportParam{
		AdminID: claim.Subject,
	}
	costReport, sessionReport, err := crh.cos.GetIncomeHistoryReport(ctx, param)
	if err != nil {
		ctx.Error(err)
		return
	}
	costRes := new([]dtos.MonthlyCostReportRes)
	sessionRes := new([]dtos.MonthlySessionReportRes)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		for _, e := range *costReport {
			*costRes = append(*costRes, dtos.MonthlyCostReportRes(e))
		}
		wg.Done()
	}()
	go func() {
		for _, e := range *sessionReport {
			*sessionRes = append(*sessionRes, dtos.MonthlySessionReportRes(e))
		}
		wg.Done()
	}()
	wg.Wait()

	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.HistoryIncomeReportRes{
			SessionReport: *sessionRes,
			CostReport:    *costRes,
		},
	})
}

func (crh *CourseRequestHandlerImpl) GetThisMonthIncomeReport(ctx *gin.Context) {
	claim, err := getAuthenticationPayload(ctx, constants.CTX_AUTH_PAYLOAD_KEY)
	if err != nil {
		ctx.Error(err)
		return
	}
	param := entity.GetAdminMonthlyReportParam{
		AdminID: claim.Subject,
	}
	totalSession, totalCost, err := crh.cos.GetThisMonthIncomeReport(ctx, param)
	res := dtos.MonthlyIncomeReportRes{
		TotalSession: *totalSession,
		TotalCost:    *totalCost,
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data:    res,
	})
}

func (crh *CourseRequestHandlerImpl) GetThisMonthMentorReport(ctx *gin.Context) {
	claim, err := getAuthenticationPayload(ctx, constants.CTX_AUTH_PAYLOAD_KEY)
	if err != nil {
		ctx.Error(err)
		return
	}
	param := entity.GetAdminMonthlyReportParam{
		AdminID: claim.Subject,
	}
	reports, err := crh.cos.GetThisMonthMentorReport(ctx, param)
	if err != nil {
		ctx.Error(err)
		return
	}
	var res []dtos.MonthlyMentorReportRes
	for _, e := range *reports {
		res = append(res, dtos.MonthlyMentorReportRes(e))
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data:    res,
	})
}

func (crh *CourseRequestHandlerImpl) StudentCourseRequestDetail(ctx *gin.Context) {
	id := ctx.Param("id")
	claim, err := getAuthenticationPayload(ctx, constants.CTX_AUTH_PAYLOAD_KEY)
	if err != nil {
		ctx.Error(err)
		return
	}
	param := entity.StudentCourseRequestDetailParam{
		CourseRequestID: id,
		StudentID:       claim.Subject,
	}
	detail, err := crh.cos.StudentCourseRequestDetail(ctx, param)
	if err != nil {
		ctx.Error(err)
		return
	}
	res := dtos.StudentCourseRequestDetailRes{
		CourseRequestID:     param.CourseRequestID,
		CourseName:          detail.CourseName,
		CourseID:            detail.CourseID,
		MentorName:          detail.MentorName,
		MentorID:            detail.MentorID,
		MentorPublicID:      detail.MentorPublicID,
		TotalPrice:          detail.TotalPrice,
		Subtotal:            detail.Subtotal,
		OperationalCost:     detail.OperationalCost,
		NumberOfSessions:    detail.NumberOfSessions,
		Status:              detail.Status,
		ExpiredAt:           detail.ExpiredAt,
		NumberOfParticipant: detail.NumberOfParticipant,
		AccountNumber:       detail.AccountNumber,
		PaymentMethodName:   detail.PaymentMethodName,
		Schedules:           []dtos.CourseScheduleRes{},
	}
	for _, sc := range detail.Schedules {
		res.Schedules = append(res.Schedules, dtos.CourseScheduleRes{
			ScheduledDate: sc.ScheduledDate,
			StartTime:     sc.StartTime,
			EndTime:       sc.EndTime,
		})
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data:    res,
	})
}

func (crh *CourseRequestHandlerImpl) StudentCourseRequestList(ctx *gin.Context) {
	var req dtos.StudentCourseRequestListReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(err)
		return
	}
	claim, err := getAuthenticationPayload(ctx, constants.CTX_AUTH_PAYLOAD_KEY)
	if err != nil {
		ctx.Error(err)
		return
	}
	if req.Status != nil {
		if err := ValidateRequestStatus(*req.Status); err != nil {
			ctx.Error(err)
			return
		}
	}
	param := entity.StudentCourseRequestListParam{
		PaginatedParam: entity.PaginatedParam{
			Limit: req.Limit,
			Page:  req.Page,
		},
		Status:    req.Status,
		StudentID: claim.Subject,
		Search:    req.Search,
	}
	if param.Limit <= 0 || param.Limit > constants.MaxLimit {
		param.Limit = constants.DefaultLimit
	}
	if param.Page <= 0 {
		param.Page = constants.DefaultPage
	}
	res, totalRow, err := crh.cos.StudentCourseRequestList(ctx, param)
	if err != nil {
		ctx.Error(err)
		return
	}
	entries := []dtos.StudentCourseRequestRes{}
	for _, req := range *res {
		entries = append(entries, dtos.StudentCourseRequestRes(req))
	}
	var filters []dtos.FilterInfo
	if req.Status != nil {
		filters = append(filters, dtos.FilterInfo{
			Name:  "Status",
			Value: *req.Status,
		})
	}
	if req.Search != nil {
		filters = append(filters, dtos.FilterInfo{
			Name:  "Search",
			Value: *req.Search,
		})
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.PaginatedResponse[dtos.StudentCourseRequestRes]{
			Entries: entries,
			PageInfo: dtos.PaginatedInfo{
				Page:     param.Page,
				FilterBy: filters,
				Limit:    param.Limit,
				TotalRow: *totalRow,
			},
		},
	})
}

func (crh *CourseRequestHandlerImpl) MentorCourseRequestDetail(ctx *gin.Context) {
	id := ctx.Param("id")
	claim, err := getAuthenticationPayload(ctx, constants.CTX_AUTH_PAYLOAD_KEY)
	if err != nil {
		ctx.Error(err)
		return
	}
	param := entity.MentorCourseRequestDetailParam{
		CourseRequestID: id,
		MentorID:        claim.Subject,
	}
	detail, err := crh.cos.MentorCourseRequestDetail(ctx, param)
	if err != nil {
		ctx.Error(err)
		return
	}
	res := dtos.MentorCourseRequestDetailRes{
		CourseRequestID:     param.CourseRequestID,
		CourseName:          detail.CourseName,
		StudentName:         detail.StudentName,
		StudentPublicID:     detail.StudentPublicID,
		TotalPrice:          detail.TotalPrice,
		Subtotal:            detail.Subtotal,
		OperationalCost:     detail.OperationalCost,
		NumberOfSessions:    detail.NumberOfSessions,
		Status:              detail.Status,
		ExpiredAt:           detail.ExpiredAt,
		PaymentMethod:       detail.PaymentMethod,
		AccountNumber:       detail.AccountNumber,
		NumberOfParticipant: detail.NumberOfParticipant,
		StudentID:           detail.StudentID,
		Schedules:           []dtos.CourseScheduleRes{},
	}
	for _, sc := range detail.Schedules {
		res.Schedules = append(res.Schedules, dtos.CourseScheduleRes{
			ScheduledDate: sc.ScheduledDate,
			StartTime:     sc.StartTime,
			EndTime:       sc.EndTime,
		})
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data:    res,
	})
}

func (crh *CourseRequestHandlerImpl) MentorCourseRequestList(ctx *gin.Context) {
	var req dtos.MentorCourseRequestListReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(err)
		return
	}
	claim, err := getAuthenticationPayload(ctx, constants.CTX_AUTH_PAYLOAD_KEY)
	if err != nil {
		ctx.Error(err)
		return
	}
	if req.Status != nil {
		if err := ValidateRequestStatus(*req.Status); err != nil {
			ctx.Error(err)
			return
		}
	}
	param := entity.MentorCourseRequestListParam{
		PaginatedParam: entity.PaginatedParam{
			Limit: req.Limit,
			Page:  req.Page,
		},
		Status:   req.Status,
		MentorID: claim.Subject,
	}
	if param.Limit <= 0 || param.Limit > constants.MaxLimit {
		param.Limit = constants.DefaultLimit
	}
	if param.Page <= 0 {
		param.Page = constants.DefaultPage
	}
	res, totalRow, err := crh.cos.MentorCourseRequestList(ctx, param)
	if err != nil {
		ctx.Error(err)
		return
	}
	entries := []dtos.MentorCourseRequestRes{}
	for _, req := range *res {
		entries = append(entries, dtos.MentorCourseRequestRes(req))
	}
	var filters []dtos.FilterInfo
	if req.Status != nil {
		filters = append(filters, dtos.FilterInfo{
			Name:  "Status",
			Value: *req.Status,
		})
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.PaginatedResponse[dtos.MentorCourseRequestRes]{
			Entries: entries,
			PageInfo: dtos.PaginatedInfo{
				Page:     param.Page,
				FilterBy: filters,
				Limit:    param.Limit,
				TotalRow: *totalRow,
			},
		},
	})
}

func (crh *CourseRequestHandlerImpl) ConfirmPayment(ctx *gin.Context) {
	claim, err := getAuthenticationPayload(ctx, constants.CTX_AUTH_PAYLOAD_KEY)
	if err != nil {
		ctx.Error(err)
		return
	}
	courseRequestID := ctx.Param("id")
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
	claim, err := getAuthenticationPayload(ctx, constants.CTX_AUTH_PAYLOAD_KEY)
	if err != nil {
		ctx.Error(err)
		return
	}
	courseRequestID := ctx.Param("id")
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
	claim, err := getAuthenticationPayload(ctx, constants.CTX_AUTH_PAYLOAD_KEY)
	if err != nil {
		ctx.Error(err)
		return
	}
	courseRequestID := ctx.Param("id")
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
	claim, err := getAuthenticationPayload(ctx, constants.CTX_AUTH_PAYLOAD_KEY)
	if err != nil {
		ctx.Error(err)
		return
	}
	var req dtos.CreateCourseRequestReq
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
	param := entity.CreateCourseRequestParam{
		CourseID:            id,
		StudentID:           claim.Subject,
		PreferredSlots:      []entity.PreferredSlot{},
		PaymentMethodID:     req.PaymentMethodID,
		NumberOfParticipant: req.NumberOfParticipant,
	}
	dateMap := make(map[time.Time]bool)
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
		if err := ValidateDate(parsedDate); err != nil {
			ctx.Error(err)
			return
		}
		if _, exist := dateMap[parsedDate]; exist {
			ctx.Error(customerrors.NewError(
				"cannot have 2 same request date",
				errors.New("there are duplicate date"),
				customerrors.InvalidAction,
			))
			return
		} else {
			dateMap[parsedDate] = true
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
