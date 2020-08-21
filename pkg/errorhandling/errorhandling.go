package errorhandling

import (
	"errors"
	"strconv"
)

type ErrorCase struct {
	Error struct {
		Code    int
		Message string
	}
}

// HandleResponseError checks for detailed codes and returns a detailed error response
func HandleResponseError(e *ErrorCase) (err error) {
	return errors.New("code:" + strconv.Itoa(e.Error.Code) + " message:" + e.Error.Message)
}
