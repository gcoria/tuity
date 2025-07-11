package middleware

import (
	"net/http"
	"tuity/internal/adapters/driving/http/dto"
	"tuity/pkg/errors"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			if domainErr, ok := err.(*errors.DomainError); ok {
				statusCode := getStatusCodeFromErrorType(domainErr.Type)
				errorResp := dto.NewErrorResponse(string(domainErr.Type), domainErr.Message, domainErr.Details)
				c.JSON(statusCode, errorResp)
			} else {
				errorResp := dto.NewErrorResponse("internal_error", "An unexpected error occurred", "")
				c.JSON(http.StatusInternalServerError, errorResp)
			}

			c.Abort()
		}
	}
}

func getStatusCodeFromErrorType(errorType errors.ErrorType) int {
	switch errorType {
	case errors.ValidationError:
		return http.StatusBadRequest
	case errors.NotFoundError:
		return http.StatusNotFound
	case errors.ConflictError:
		return http.StatusConflict
	case errors.InternalError:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
