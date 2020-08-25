package gomodio_test

import (
	"fmt"
	"log"

	"github.com/M4cs/gomodio"
)

var (
	apiKey = "your-api-key"
	email  = "your-email"
)

// ExampleNewUserSimple is an example of simple authentication
func ExampleNewUserSimple() {
	user := gomodio.NewUser(apiKey, email)
	fmt.Println(user.APIKey())
	// Output: your-api-key
}

// ExampleNewUserOAuth is an example of OAuth2 Authentication
func ExampleNewUserOAuth() {
	user := gomodio.NewUser(apiKey, email)
	user.RequestSecurityCode()
	user.SetOAuth2Token("token from e-mail")
	fmt.Println(user.OAuth2Token())
	// Output: your-oauth2-token
}

// ExampleGetGames is an example of user.GetGames
func ExampleGetGames() {
	user := gomodio.NewUser(apiKey, email)
	games, err := user.GetGames(map[string]string{
		"name": "Skater XL",
		"q":    "Skater",
	})
	if err != nil {
		log.Fatalln(err.Error())
	}
	fmt.Println(games.ResultCount)
	// Output: 1
}
