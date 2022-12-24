package errorEntity

import (
	"fmt"
	"net/http"
)

type ServiceError struct {
	StatusCode int
	Message    string
	Err        error
}

func (se ServiceError) Error() string {
	return fmt.Sprintf("Error: %s /n %v", se.Message, se.Err)
}

func New(statusCode int, message string, err error) *ServiceError {
	return &ServiceError{
		StatusCode: statusCode,
		Message:    message,
		Err:        err,
	}
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
