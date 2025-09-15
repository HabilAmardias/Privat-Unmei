package middlewares

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/dtos"

	"github.com/gin-gonic/gin"
)

type RecaptchaResponse struct {
	Success     bool     `json:"success"`
	ChallengeTS string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	ErrorCodes  []string `json:"error-codes"`
	Score       float64  `json:"score"`
	Action      string   `json:"action"`
}

func verifyCaptcha(token string, clientIP string) (bool, error) {
	RECAPTCHA_SECRET_KEY := os.Getenv("RECAPTCHA_SECRET_KEY")
	RECAPTCHA_VERIFY_URL := os.Getenv("RECAPTCHA_VERIFY_SITE")
	data := url.Values{}
	data.Set("secret", RECAPTCHA_SECRET_KEY)
	data.Set("response", token)
	if clientIP != "" {
		data.Set("remoteip", clientIP)
	}

	resp, err := http.PostForm(RECAPTCHA_VERIFY_URL, data)
	if err != nil {
		return false, customerrors.NewError(
			"failed to verify captcha",
			err,
			customerrors.CommonErr,
		)
	}
	defer resp.Body.Close()

	var recaptchaResp RecaptchaResponse
	if err := json.NewDecoder(resp.Body).Decode(&recaptchaResp); err != nil {
		return false, customerrors.NewError(
			"failed to verify captcha",
			err,
			customerrors.CommonErr,
		)
	}

	if !recaptchaResp.Success {
		return false, customerrors.NewError(
			"failed to verify captcha",
			errors.New("failed to verify captcha"),
			customerrors.InvalidAction,
		)
	}
	if recaptchaResp.Score < 0.5 {
		return false, customerrors.NewError(
			"invalid request",
			errors.New("captcha score is too low"),
			customerrors.InvalidAction,
		)
	}

	return true, nil
}

func CaptchaMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dtos.RegisterStudentReq
		if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
			ctx.Error(err)
			ctx.Abort()
			return
		}
		clientIP := ctx.ClientIP()
		valid, err := verifyCaptcha(req.CaptchaToken, clientIP)
		if err != nil {
			ctx.Error(err)
			ctx.Abort()
			return
		}
		if !valid {
			ctx.Error(customerrors.NewError(
				"invalid request",
				errors.New("invalid request"),
				customerrors.InvalidAction,
			))
			ctx.Abort()
			return
		}
		ctx.Set("validated_request", req)
		ctx.Next()
	}
}
