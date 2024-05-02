package http

import (
	"errors"
	"net/http"

	"github.com/arasan1289/hexagonal-demo/internal/core/domain"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// response represents a response body format
type response struct {
	Success           bool              `json:"success"`
	Data              any               `json:"data,omitempty"`
	Errors            map[string]string `json:"errors,omitempty"`
	DescriptiveErrors []ValidationError `json:"descriptive_errors,omitempty"`
}

// newResponse is a helper function to create a response body
func newResponse(success bool, data any, errors map[string]string, descriptiveErrs []ValidationError) response {
	return response{
		Success:           success,
		Data:              data,
		Errors:            errors,
		DescriptiveErrors: descriptiveErrs,
	}
}

// errorStatusMap is a map of defined error messages and their corresponding http status codes
var errorStatusMap = map[error]int{
	domain.ErrDataNotFound:               http.StatusNotFound,
	domain.ErrConflictingData:            http.StatusConflict,
	domain.ErrInvalidCredentials:         http.StatusUnauthorized,
	domain.ErrUnauthorized:               http.StatusUnauthorized,
	domain.ErrEmptyAuthorizationHeader:   http.StatusUnauthorized,
	domain.ErrInvalidAuthorizationHeader: http.StatusUnauthorized,
	domain.ErrInvalidAuthorizationType:   http.StatusUnauthorized,
	domain.ErrInvalidToken:               http.StatusUnauthorized,
	domain.ErrExpiredToken:               http.StatusUnauthorized,
	domain.ErrForbidden:                  http.StatusForbidden,
	domain.ErrNoUpdatedData:              http.StatusBadRequest,
	domain.ErrInternal:                   http.StatusInternalServerError,
	domain.ErrRateLimitExceeded:          http.StatusTooManyRequests,
	domain.ErrOTPExpired:                 http.StatusBadRequest,
	domain.ErrOTPMismatch:                http.StatusBadRequest,
}

// parseError parses error messages from the error object and returns a slice of error messages
func parseError(ctx *gin.Context, err error) (map[string]string, []ValidationError) {
	var errMsgs = make(map[string]string)
	var descriptiveError []ValidationError

	// Check if the error is a validation error
	var validationErr validator.ValidationErrors
	if errors.As(err, &validationErr) {
		// Determine error type (simple or descriptive)
		errorType := ctx.Query("error_type")

		// Log each validation error
		for _, err := range validationErr {
			ctx.Error(err)
		}

		// Generate error messages based on error type
		jsonFormatter := NewJSONFormatter()
		if errorType == "descriptive" {
			descriptiveError = jsonFormatter.Descriptive(validationErr)
		} else {
			errMsgs = jsonFormatter.Simple(validationErr)
		}
	} else {
		// Log the error
		ctx.Error(err)

		// Add a generic error message
		errMsgs["message"] = err.Error()

		// If descriptive error is requested, append to list
		errorType := ctx.Query("error_type")
		if errorType == "descriptive" {
			descriptiveError = append(descriptiveError, ValidationError{Field: "message", Reason: err.Error()})
		}
	}

	return errMsgs, descriptiveError
}

// handleSuccess sends a success response with the specified status code and optional data
func handleSuccess(ctx *gin.Context, data any) {
	rsp := newResponse(true, data, nil, nil)
	ctx.JSON(http.StatusOK, rsp)
}

// handleError sends a error response with the specified status code and error message
func handleError(ctx *gin.Context, err error) {
	switch err {
	case gorm.ErrRecordNotFound:
		err = domain.ErrDataNotFound
	case gorm.ErrDuplicatedKey:
		err = domain.ErrConflictingData
	}
	statusCode, ok := errorStatusMap[err]
	if !ok {
		statusCode = http.StatusInternalServerError
	}
	errMsg, descriptiveErrs := parseError(ctx, err)
	errRsp := newResponse(false, nil, errMsg, descriptiveErrs)
	ctx.JSON(statusCode, errRsp)
}

// validationError sends a error response with the specified status code and error message
func validationError(ctx *gin.Context, err error) {
	errs, descriptiveErrs := parseError(ctx, err)
	errRsp := newResponse(false, nil, errs, descriptiveErrs)
	ctx.JSON(http.StatusBadRequest, errRsp)
}
