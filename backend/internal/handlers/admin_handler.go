package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"privat-unmei/internal/constants"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/dtos"
	"privat-unmei/internal/entity"
	"privat-unmei/internal/services"

	"github.com/gin-gonic/gin"
)

type AdminHandlerImpl struct {
	as *services.AdminServiceImpl
	ms *services.MentorServiceImpl
}

func CreateAdminHandler(as *services.AdminServiceImpl, ms *services.MentorServiceImpl) *AdminHandlerImpl {
	return &AdminHandlerImpl{as, ms}
}

func (ah *AdminHandlerImpl) GetMentorList(ctx *gin.Context) {
	var req dtos.ListMentorReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(err)
		return
	}
	param := entity.ListMentorParam{
		PaginatedParam: entity.PaginatedParam{
			Limit: req.Limit,
			Page:  req.Page,
		},
		Search:               req.Search,
		SortYearOfExperience: req.SortYearOfExperience,
	}
	if req.Limit < 0 {
		param.Limit = constants.DefaultLimit

	}
	if req.Page <= 0 {
		param.Page = constants.DefaultPage
	}
	res, totalRow, err := ah.ms.GetMentorList(ctx, param)
	if err != nil {
		ctx.Error(err)
		return
	}
	entries := []dtos.ListMentorRes{}
	for _, mentor := range *res {
		entries = append(entries, dtos.ListMentorRes(mentor))
	}
	var filters []dtos.FilterInfo
	if req.Search != nil {
		filter := dtos.FilterInfo{
			Name:  "Search",
			Value: *req.Search,
		}
		filters = append(filters, filter)
	}
	var sorts []dtos.SortInfo
	if req.SortYearOfExperience != nil {
		sortInfo := dtos.SortInfo{
			Name: "years_of_experience",
			ASC:  *req.SortYearOfExperience,
		}
		sorts = append(sorts, sortInfo)
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.PaginatedResponse[dtos.ListMentorRes]{
			Entries: entries,
			PageInfo: dtos.PaginatedInfo{
				Page:     param.Page,
				Limit:    param.Limit,
				TotalRow: *totalRow,
				FilterBy: filters,
				SortBy:   sorts,
			},
		},
	})

}
func (ah *AdminHandlerImpl) DeleteMentor(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.Error(customerrors.NewError(
			"no mentor id given",
			errors.New("no mentor id given"),
			customerrors.InvalidAction,
		))
		return
	}
	param := entity.DeleteMentorParam{
		ID: id,
	}
	if err := ah.ms.DeleteMentor(ctx, param); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.UpdateMentorForAdminRes{
			ID: id,
		},
	})
}
func (ah *AdminHandlerImpl) UpdateMentor(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.Error(customerrors.NewError(
			"no mentor id given",
			errors.New("no mentor id given"),
			customerrors.InvalidAction,
		))
		return
	}
	var req dtos.UpdateMentorForAdminReq
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.Error(err)
		return
	}
	if req.WhatsappNumber != nil && !ValidatePhoneNumber(*req.WhatsappNumber) {
		ctx.Error(customerrors.NewError(
			"invalid whatsapp number",
			errors.New("invalid whatsapp number"),
			customerrors.InvalidAction,
		))
		return
	}
	param := entity.UpdateMentorParam{
		ID: id,
		UpdateMentorQuery: entity.UpdateMentorQuery{
			WhatsappNumber:    req.WhatsappNumber,
			YearsOfExperience: req.YearsOfExperience,
		},
	}
	if err := ah.ms.UpdateMentorForAdmin(ctx, param); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.UpdateMentorForAdminRes{
			ID: param.ID,
		},
	})
}
func (ah *AdminHandlerImpl) GenerateRandomPassword(ctx *gin.Context) {
	pass, err := ah.as.GenerateRandomPassword()
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusCreated, dtos.Response{
		Success: true,
		Data: dtos.GeneratePasswordRes{
			Password: pass,
		},
	})
}

func (ah *AdminHandlerImpl) AddNewMentor(ctx *gin.Context) {
	var req dtos.AddNewMentorReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(err)
		return
	}
	if !ValidateDegree(req.Degree) {
		ctx.Error(customerrors.NewError(
			"invalid degree",
			errors.New("degree given is invalid"),
			customerrors.InvalidAction,
		))
		return
	}
	if !ValidatePhoneNumber(req.WhatsappNumber) {
		ctx.Error(customerrors.NewError(
			"invalid whatsapp number",
			errors.New("whatsapp number given is invalid"),
			customerrors.InvalidAction,
		))
		return
	}
	headerFile, err := ctx.FormFile("file")
	if err != nil {
		ctx.Error(
			customerrors.NewError(
				"Failed to upload pharmacy logo",
				err,
				customerrors.InvalidAction,
			),
		)
		return
	}
	if headerFile.Size > constants.FileSizeThreshold {
		ctx.Error(customerrors.NewError(
			"File size is too large",
			fmt.Errorf("file size is too large"),
			customerrors.InvalidAction,
		))
		return
	}
	file, err := headerFile.Open()
	if err != nil {
		ctx.Error(
			customerrors.NewError(
				"Failed to upload file",
				err,
				customerrors.InvalidAction,
			),
		)
		return
	}
	defer file.Close()

	buff := make([]byte, 512)
	if _, err := file.Read(buff); err != nil {
		ctx.Error(
			customerrors.NewError(
				"Failed to upload file",
				err,
				customerrors.InvalidAction,
			),
		)
		return
	}

	file.Seek(0, io.SeekStart)
	if fileType := http.DetectContentType(buff); fileType != constants.PDFType {
		ctx.Error(
			customerrors.NewError(
				"Uploaded file must be .pdf",
				err,
				customerrors.InvalidAction,
			),
		)
		return
	}
	param := entity.AddNewMentorParam{
		Name:              req.Name,
		Email:             req.Email,
		Password:          req.Password,
		ResumeFile:        file,
		Bio:               req.Bio,
		YearsOfExperience: req.YearsOfExperience,
		WhatsappNumber:    req.WhatsappNumber,
		Degree:            req.Degree,
		Major:             req.Major,
		Campus:            req.Campus,
	}
	if err := ah.ms.AddNewMentor(ctx, param); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusCreated, dtos.Response{
		Success: true,
		Data: dtos.MessageResponse{
			Message: "Succesfully create mentor account",
		},
	})
}

func (ah *AdminHandlerImpl) GetStudentList(ctx *gin.Context) {
	var req dtos.GetStudentListReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.Error(err)
		return
	}
	param := entity.ListStudentParam{
		PaginatedParam: entity.PaginatedParam{
			Limit: req.Limit,
			Page:  req.Page,
		},
	}
	if req.Limit < 0 {
		param.Limit = constants.DefaultLimit
	}
	if req.Page <= 0 {
		param.Page = constants.DefaultPage
	}
	students, totalRow, err := ah.as.GetStudentList(ctx, param)
	if err != nil {
		ctx.Error(err)
		return
	}
	entries := []dtos.ListStudentRes{}
	for _, student := range *students {
		entries = append(entries, dtos.ListStudentRes(student))
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.PaginatedResponse[dtos.ListStudentRes]{
			Entries: entries,
			PageInfo: dtos.PaginatedInfo{
				Page:     param.Page,
				Limit:    param.Limit,
				TotalRow: *totalRow,
			},
		},
	})
}

func (ah *AdminHandlerImpl) Login(ctx *gin.Context) {
	var req dtos.AdminLoginReq
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.Error(err)
		return
	}
	param := entity.AdminLoginParam(req)
	token, err := ah.as.Login(ctx, param)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.AdminLoginRes{
			Token: token,
		},
	})
}
