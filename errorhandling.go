package gomodio

import (
	"errors"
	"strconv"
)

// ErrorCase for gomodio
type ErrorCase struct {
	Error Error `json:"error"`
}

// Error for gomodio
type Error struct {
	Code    int    `json:"error_ref"`
	Message string `json:"message"`
}

// HandleResponseError checks for detailed codes and returns a detailed error response
func HandleResponseError(e ErrorCase) (err error) {
	return errors.New("code:" + strconv.Itoa(e.Error.Code) + " message:" + e.Error.Message)
}
