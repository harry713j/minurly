package errors

import (
	"net/http"
	"strings"
)

type ApiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func (e *ApiError) Error() string {
	return e.Message
}

func (e *ApiError) Is(target error) bool {
	_, ok := target.(*ApiError)

	return ok
}

func makeUpperCaseWithUnderscores(str string) string {
	return strings.ToUpper(strings.ReplaceAll(str, " ", "_"))
}

func NewBadRequestErr(msg string) *ApiError {
	return &ApiError{
		Code:    makeUpperCaseWithUnderscores(http.StatusText(http.StatusBadRequest)),
		Message: msg,
		Status:  http.StatusBadRequest,
	}
}

func NewNotFoundErr(msg string) *ApiError {
	return &ApiError{
		Code:    makeUpperCaseWithUnderscores(http.StatusText(http.StatusNotFound)),
		Message: msg,
		Status:  http.StatusNotFound,
	}
}

func NewUnauthorizedErr(msg string) *ApiError {
	return &ApiError{
		Code:    makeUpperCaseWithUnderscores(http.StatusText(http.StatusUnauthorized)),
		Message: msg,
		Status:  http.StatusUnauthorized,
	}
}

func NewInternalServerErr() *ApiError {
	return &ApiError{
		Code:    makeUpperCaseWithUnderscores(http.StatusText(http.StatusInternalServerError)),
		Message: http.StatusText(http.StatusInternalServerError),
		Status:  http.StatusInternalServerError,
	}
}
