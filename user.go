package gomodio

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// User Struct for mod.io
type User struct {
	apikey      string
	email       string
	oauth2token string
}

// ExchangeResponse Struct for Response of Email Exchange
type ExchangeResponse struct {
	OAuthToken string `json:"access_token"`
	Code       int    `json:"code"`
}

// NewUser - Initializes a new User
func NewUser(apikey string, email string) *User {
	return &User{apikey, email, ""}
}

// APIKey returns the User's API key
func (u *User) APIKey() string {
	return u.apikey
}

// Email returns the User's Email
func (u *User) Email() string {
	return u.email
}

// OAuth2Token returns the User's OAuth2Token
func (u *User) OAuth2Token() string {
	return u.oauth2token
}

// SetOAuth2Token sets a User's OAuth2Token
// token is the OAuth2 Token from mod.io
func (u *User) SetOAuth2Token(token string) {
	u.oauth2token = token
}

// RequestSecurityCode Authenticate with mod.io Using API Key Only
func (u *User) RequestSecurityCode() bool {
	requestBody := url.Values{
		"api_key": {u.APIKey()},
		"email":   {u.Email()},
	}
	resp, err := http.Post("https://api.mod.io/v1/oauth/emailrequest", "application/x-www-form-urlencoded", strings.NewReader(requestBody.Encode()))
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Fatalln("Failed To Send Email Request. Check E-Mail/API Key")
	} else {
		return true
	}
	return false
}

// ExchangeSecurityCode Function
func (u *User) ExchangeSecurityCode(securitycode string) *User {
	reqBody := url.Values{
		"api_key":       {u.APIKey()},
		"security_code": {securitycode},
		"date_expires":  {strconv.FormatInt(time.Now().Unix()+31536000, 10)},
	}
	resp, err := http.Post("https://api.mod.io/v1/oauth/emailexchange", "application/x-www-form-urlencoded", strings.NewReader(reqBody.Encode()))
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	var res ExchangeResponse
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Fatalln(err)
	}
	if res.Code != 200 {
		log.Fatalln("Incorrect Security Code Provided")
	} else {
		u.SetOAuth2Token(res.OAuthToken)
	}
	return u
}
