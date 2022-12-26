package errorEntity

import (
	"fmt"
	"net/http"
)

type ServiceError struct {
	StatusCode int
	message    string
	err        error
}

func (se ServiceError) Error() string {
	return fmt.Sprintf("Error: %s", se.message)
}

func New(statusCode int, message string, err error) *ServiceError {
	return &ServiceError{
		StatusCode: statusCode,
		message:    message,
		err:        err,
	}
}

func ConflictError(message string, err error) *ServiceError {
	return New(http.StatusConflict, message, err)
}

func InternalServerError(err error) *ServiceError {
	return New(http.StatusInternalServerError, "Internal Server Error", err)
}

func NotFoundError(message string, err error) *ServiceError {
	return New(http.StatusNotFound, message, err)
}

func BadRequestError(message string, err error) *ServiceError {
	return New(http.StatusBadRequest, message, err)
}
