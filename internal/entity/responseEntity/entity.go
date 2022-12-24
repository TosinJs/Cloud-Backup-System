package responseEntity

import "tosinjs/cloud-backup/internal/entity/errorEntity"

type ResponseEntity struct {
	StatusCode int    `json:"status,omitempty"`
	Message    string `json:"message,omitempty"`
	Data       any    `json:"data,omitempty"`
	Error      any    `json:"error,omitempty"`
	Path       string `json:"path,omitempty"`
}

func BuildResponseObject(statusCode int, path string, body any) ResponseEntity {
	return ResponseEntity{
		StatusCode: statusCode,
		Message:    "success",
		Data:       body,
		Path:       path,
	}
}

func BuildErrorResponseObject(statusCode int, errorMessage, path string) ResponseEntity {
	return ResponseEntity{
		StatusCode: statusCode,
		Message:    "failed",
		Error:      errorMessage,
		Path:       path,
	}
}

func BuildServiceErrorResponseObject(err *errorEntity.ServiceError, path string) ResponseEntity {
	return ResponseEntity{
		StatusCode: err.StatusCode,
		Message:    "failed",
		Error:      err.Message,
		Path:       path,
	}
}
