package handlers

import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"privat-unmei/internal/constants"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/dtos"
	"privat-unmei/internal/entity"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	degreelist = []string{"bachelor", "diploma", "high school", "master", "professor"}
)

func getAuthenticationPayload(ctx *gin.Context) (*entity.CustomClaim, error) {

	claims, ok := ctx.Get(constants.CTX_AUTH_PAYLOAD_KEY)
	if !ok {
		return nil, customerrors.NewError(
			"user credential identification failed",
			errors.New("cannot find authentication claim"),
			customerrors.CommonErr,
		)
	}

	customClaims, ok := claims.(*entity.CustomClaim)
	if !ok {
		return nil, customerrors.NewError(
			"user credential identification failed",
			errors.New("cannot parse authentication claim"),
			customerrors.CommonErr,
		)
	}
	return customClaims, nil
}

func ValidateDegree(degree string) bool {
	for _, item := range degreelist {
		if degree == item {
			return true
		}
	}
	return false
}

func ValidatePhoneNumber(phoneNumber string) bool {

	pattern := `^0\d{9,12}$`

	regex := regexp.MustCompile(pattern)

	return regex.MatchString(phoneNumber)
}

func ValidateFile(headerFile *multipart.FileHeader, fileSizeThresh int64, fileType string) (multipart.File, error) {
	if headerFile == nil {
		return nil, nil
	}
	if headerFile.Size > fileSizeThresh {
		return nil, customerrors.NewError(
			"file size is too large",
			errors.New("file size is too large"),
			customerrors.InvalidAction,
		)
	}
	file, err := headerFile.Open()
	if err != nil {
		return nil, customerrors.NewError(
			"failed to upload file",
			err,
			customerrors.CommonErr,
		)
	}
	defer file.Close()
	buff := make([]byte, 512)
	if _, err := file.Read(buff); err != nil {
		return nil, customerrors.NewError(
			"failed to upload file",
			err,
			customerrors.CommonErr,
		)
	}
	file.Seek(0, io.SeekStart)
	if fileExt := http.DetectContentType(buff); fileExt != fileType {
		return nil, customerrors.NewError(
			"invalid file format",
			errors.New("invalid file format"),
			customerrors.InvalidAction,
		)
	}
	return file, nil
}

func ValidateMethod(method string) bool {
	return method == "offline" || method == "online" || method == "hybrid"
}

func ValidateDate(date time.Time) bool {
	now := time.Now()
	return date.After(now)
}

func CheckDateUniqueness(slots []dtos.PreferredSlot) error {
	dateMap := make(map[string]bool)
	for _, slot := range slots {
		if _, exist := dateMap[slot.Date]; exist {
			return customerrors.NewError(
				"cannot reserve multiple session on same date",
				errors.New("there are duplicate date"),
				customerrors.InvalidAction,
			)
		} else {
			dateMap[slot.Date] = true
		}
	}
	return nil
}
