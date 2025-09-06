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
	statuslist = []string{"reserved", "pending payment", "scheduled", "completed", "cancelled"}
)

func ValidateCategories(categories []int) error {
	catMap := make(map[int]bool)
	for _, cat := range categories {
		if _, exist := catMap[cat]; exist {
			return customerrors.NewError(
				"there are duplicate categories",
				errors.New("there are duplicate categories"),
				customerrors.InvalidAction,
			)
		} else {
			catMap[cat] = true
		}
	}
	return nil
}

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

func getAuthenticationToken(ctx *gin.Context) (string, error) {
	val, ok := ctx.Get(constants.CTX_AUTH_TOKEN_KEY)
	if !ok {
		return "", customerrors.NewError(
			"user credential identification failed",
			errors.New("cannot find authentication token"),
			customerrors.CommonErr,
		)
	}
	token, ok := val.(string)
	if !ok {
		return "", customerrors.NewError(
			"user credential identification failed",
			errors.New("cannot parse authentication token"),
			customerrors.CommonErr,
		)
	}
	return token, nil
}

func ValidateDegree(degree string) error {
	for _, item := range degreelist {
		if degree == item {
			return nil
		}
	}
	return customerrors.NewError(
		"invalid degree",
		errors.New("invalid degree"),
		customerrors.InvalidAction,
	)
}

func ValidatePhoneNumber(phoneNumber string) error {

	pattern := `^0\d{9,12}$`

	regex := regexp.MustCompile(pattern)

	if regex.MatchString(phoneNumber) {
		return nil
	}
	return customerrors.NewError(
		"invalid phone number",
		errors.New("invalid phone number"),
		customerrors.InvalidAction,
	)
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

func ValidateMethod(method string) error {
	if method == "offline" || method == "online" || method == "hybrid" {
		return nil
	}
	return customerrors.NewError(
		"invalid method",
		errors.New("invalid method"),
		customerrors.InvalidAction,
	)
}

func ValidateDate(date time.Time) error {
	now := time.Now()
	if date.After(now) {
		return nil
	}
	return customerrors.NewError(
		"invalid date",
		errors.New("invalid date"),
		customerrors.InvalidAction,
	)
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

func ValidateRequestStatus(status string) error {
	for _, item := range statuslist {
		if status == item {
			return nil
		}
	}
	return customerrors.NewError(
		"invalid status",
		errors.New("invalid status"),
		customerrors.InvalidAction,
	)
}
