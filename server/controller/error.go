package controller

import "fmt"

type errorResponse struct {
	Message string `json:"message"`
}

func newErrorResponse(message string, args ...interface{}) *errorResponse {
	return &errorResponse{
		Message: fmt.Sprintf(message, args...),
	}
}
