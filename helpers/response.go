package helpers

import (
	"fmt"
	"strings"
)

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Error   interface{}      `json:"error"`
	Data    interface{}      `json:"data"`
}

type EmptyResponse struct {}

func BuildResponse(status int, message string, data interface{}) Response {
	return Response{
		Status:  status,
		Message: message,
		Data:    data,
	}
}

func BuildErrorResponse(status int, message string, error interface{}) Response {
	splittedError := strings.Split(fmt.Sprint(error), "\n")
	return Response{
		Status:  status,
		Message: message,
		Error:   splittedError,
	}
}