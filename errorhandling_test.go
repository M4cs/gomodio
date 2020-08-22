package gomodio_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/M4cs/gomodio"
)

func TestErrorhandling(t *testing.T) {
	fmt.Println("BUSTED")
	er := gomodio.Error{
		Code:    1,
		Message: "Test Msg",
	}
	var ec gomodio.ErrorCase
	ec = gomodio.ErrorCase{
		Error: er,
	}
	err := gomodio.HandleResponseError(ec)
	if err.Error() != ("code:" + strconv.Itoa(ec.Error.Code) + " message:" + ec.Error.Message) {
		t.Errorf("Error Msg is Incorrect!")
	}
}
