package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
	"net/http"
	"os"
	"privat-unmei/internal/constants"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/dtos"
	"privat-unmei/internal/entity"
	"privat-unmei/internal/services"
	"time"

	"github.com/gin-gonic/gin"
)

type StudentHandlerImpl struct {
	ss *services.StudentServiceImpl
}

func CreateStudentHandler(ss *services.StudentServiceImpl) *StudentHandlerImpl {
	return &StudentHandlerImpl{ss}
}

func (sh *StudentHandlerImpl) GetStudentProfile(ctx *gin.Context) {
	claim, err := getAuthenticationPayload(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	param := entity.StudentProfileParam{
		ID: claim.Subject,
	}
	profile, err := sh.ss.GetStudentProfile(ctx, param)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data:    dtos.StudentProfileRes(*profile),
	})
}

func (sh *StudentHandlerImpl) ChangePassword(ctx *gin.Context) {
	var req dtos.MentorChangePasswordReq
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.Error(err)
		return
	}
	claim, err := getAuthenticationPayload(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	param := entity.StudentChangePasswordParam{
		ID:          claim.Subject,
		NewPassword: req.NewPassword,
	}
	if err := sh.ss.ChangePassword(ctx, param); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.MessageResponse{
			Message: "Succesfully change password",
		},
	})
}

func (sh *StudentHandlerImpl) GoogleLogin(ctx *gin.Context) {
	expTime := time.Now().Add(30 * time.Minute)
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expTime}
	http.SetCookie(ctx.Writer, &cookie)

	url := sh.ss.GoogleLogin(state)

	log.Println("Redirecting to: ", url)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (sh *StudentHandlerImpl) GoogleLoginCallback(ctx *gin.Context) {
	state, err := ctx.Cookie("oauthstate")
	if err != nil {
		ctx.Error(customerrors.NewError(
			"failed to login",
			err,
			customerrors.CommonErr,
		))
		return
	}
	if ctx.Query("state") != state {
		ctx.Error(customerrors.NewError(
			"invalid credential",
			err,
			customerrors.InvalidAction,
		))
		return
	}
	code := ctx.Query("code")
	if code == "" {
		ctx.Error(customerrors.NewError(
			"failed to login",
			errors.New("missing code in query param"),
			customerrors.InvalidAction,
		))
		return
	}
	authToken, _, err := sh.ss.GoogleLoginCallback(ctx, code)
	if err != nil {
		ctx.Error(err)
		return
	}
	clientURL := os.Getenv("CLIENT_DOMAIN")
	ctx.Redirect(http.StatusTemporaryRedirect, clientURL+"/google-callback/"+authToken)
}

func (sh *StudentHandlerImpl) UpdateStudentProfile(ctx *gin.Context) {
	claim, err := getAuthenticationPayload(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	var req dtos.UpdateStudentReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.Error(err)
		return
	}

	headerFile, _ := ctx.FormFile("file")
	file, err := ValidateFile(headerFile, constants.FileSizeThreshold, constants.PNGType)
	if err != nil {
		ctx.Error(err)
		return
	}
	param := entity.UpdateStudentParam{
		ID:           claim.Subject,
		Name:         req.Name,
		Bio:          req.Bio,
		ProfileImage: file,
	}
	if err := sh.ss.UpdateStudentProfile(ctx, param); err != nil {
		ctx.Error(err)
		return
	}
	if file != nil {
		file.Close()
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.UpdateStudentRes{
			ID: claim.Subject,
		},
	})
}

func (sh *StudentHandlerImpl) GetStudentList(ctx *gin.Context) {
	var req dtos.GetStudentListReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.Error(err)
		return
	}
	claim, err := getAuthenticationPayload(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	param := entity.ListStudentParam{
		PaginatedParam: entity.PaginatedParam{
			Limit: req.Limit,
			Page:  req.Page,
		},
		AdminID: claim.Subject,
	}
	if param.Limit <= 0 || param.Limit > constants.MaxLimit {
		param.Limit = constants.DefaultLimit
	}
	if param.Page <= 0 {
		param.Page = constants.DefaultPage
	}
	students, totalRow, err := sh.ss.GetStudentList(ctx, param)
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

func (sh *StudentHandlerImpl) SendVerificationEmail(ctx *gin.Context) {
	claims, err := getAuthenticationPayload(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	if err := sh.ss.SendVerificationEmail(ctx, claims.Subject); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.MessageResponse{
			Message: "Successfully sent verification email",
		},
	})
}

func (sh *StudentHandlerImpl) ResetPassword(ctx *gin.Context) {
	var req dtos.ResetPasswordReq
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.Error(err)
		return
	}
	token, err := getAuthenticationToken(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	claim, err := getAuthenticationPayload(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	param := entity.ResetPasswordParam{
		ID:          claim.Subject,
		Token:       token,
		NewPassword: req.NewPassword,
	}
	if err := sh.ss.ResetPassword(ctx, param); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.MessageResponse{
			Message: "Sucessfully reset password",
		},
	})
}

func (sh *StudentHandlerImpl) SendResetTokenEmail(ctx *gin.Context) {
	var req dtos.SendResetTokenEmailReq
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.Error(err)
		return
	}
	if err := sh.ss.SendResetTokenEmail(ctx, req.Email); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.MessageResponse{
			Message: "Succesfully send reset password email",
		},
	})
}

func (sh *StudentHandlerImpl) Verify(ctx *gin.Context) {
	token, err := getAuthenticationToken(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	claim, err := getAuthenticationPayload(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	param := entity.VerifyStudentParam{
		Token: token,
		ID:    claim.Subject,
	}
	if err := sh.ss.Verify(ctx, param); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.MessageResponse{
			Message: "Successfully Verified",
		},
	})
}

func (sh *StudentHandlerImpl) Login(ctx *gin.Context) {
	domain := os.Getenv("COOKIE_DOMAIN")
	var req dtos.LoginStudentReq
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.Error(err)
		return
	}
	param := entity.StudentLoginParam{
		Email:    req.Email,
		Password: req.Password,
	}
	authToken, refreshToken, status, err := sh.ss.Login(ctx, param)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.SetSameSite(http.SameSiteNoneMode)
	ctx.SetCookie(constants.AUTH_COOKIE_KEY, *authToken, int(constants.AUTH_AGE), "/", domain, false, true)
	ctx.SetCookie(constants.REFRESH_COOKIE_KEY, *refreshToken, int(constants.REFRESH_AGE), "/", domain, false, true)
	ctx.JSON(http.StatusOK, dtos.Response{
		Success: true,
		Data: dtos.LoginStudentRes{
			Status: *status,
		},
	})
}

func (sh *StudentHandlerImpl) Register(ctx *gin.Context) {
	reg, exists := ctx.Get("validated_request")
	if !exists {
		ctx.Error(customerrors.NewError(
			"something went wrong",
			errors.New("request body not found"),
			customerrors.CommonErr,
		))
		return
	}
	req, ok := reg.(dtos.RegisterStudentReq)
	if !ok {
		ctx.Error(customerrors.NewError(
			"invalid request",
			errors.New("invalid request body"),
			customerrors.CommonErr,
		))
	}
	param := entity.StudentRegisterParam{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Status:   constants.UnverifiedStatus,
	}
	if err := sh.ss.Register(ctx, param); err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusCreated, dtos.Response{
		Success: true,
		Data: dtos.MessageResponse{
			Message: "Successfully Registered",
		},
	})
}
