package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"privat-unmei/internal/customerrors"
	"privat-unmei/internal/dtos"
	"privat-unmei/internal/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ErrorMiddleware(logger logger.CustomLogger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		if len(ctx.Errors) == 0 {
			return
		}
		err := ctx.Errors[0]
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			fieldErrors := make([]dtos.DetailsError, 0)

			for _, fe := range ve {
				logger.Warnln(err.Error())

				fieldErrors = append(fieldErrors, dtos.DetailsError{
					Title:   fe.Field(),
					Message: fmt.Sprintf("invalid input on field %s", fe.Field()),
				})
			}
			ctx.JSON(http.StatusBadRequest, dtos.Response{
				Success: false,
				Data:    fieldErrors,
			})
			return
		}
		var ce *customerrors.CustomError
		if errors.As(err, &ce) {
			logger.Warnln(ce.ErrLog.Error())
			ctx.JSON(ce.GetStatusCode(), dtos.Response{
				Success: false,
				Data: dtos.MessageResponse{
					Message: ce.ErrUser,
				},
			})
			return
		}
		logger.Warnln(err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, dtos.Response{
			Success: false,
			Data: dtos.MessageResponse{
				Message: "Something went wrong",
			},
		})
	}
}
