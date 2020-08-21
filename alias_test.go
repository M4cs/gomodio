package gomodio_test

import (
	"testing"

	"github.com/M4cs/gomodio"
)

var user *gomodio.User

func init() {
	user = gomodio.NewUser("c1ac9d8941bb8d42bf4f43472f9ec940", "mabridgland@protonmail.com")
}

// TestNewUser tests the intialization of a user
func TestNewUser(t *testing.T) {
	if user.APIKey() == "" {
		t.Errorf("User APIKey is Empty")
	}
	if !user.RequestSecurityCode() {
		t.Errorf("Could not send security code!")
	}
}

func TestGetGames(t *testing.T) {
	games, err := gomodio.GetGames(map[string]string{
		"q": "Skater XL",
	}, user)
	if err != nil {
		t.Errorf("Hit Error: " + err.Error())
	}
	if games.ResultCount < 0 {
		t.Errorf("No Results!")
	}
}
