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

// Comments struct representing the JSON response of Get Comments
type Comments struct {
	Data         []Comment `json:"data"`
	ResultCount  int       `json:"result_count"`
	ResultOffset int       `json:"result_offset"`
	ResultLimit  int       `json:"result_limit"`
	ResultTotal  int       `json:"result_total"`
}

// Comment struct representing single comment objects
type Comment struct {
	ID    int `json:"id"`
	ModID int `json:"mod_id"`
	User  struct {
		ID         int    `json:"id"`
		NameID     string `json:"name_id"`
		Username   string `json:"username"`
		DateOnline int    `json:"date_online"`
		Avatar     struct {
			Filename     string `json:"filename"`
			Original     string `json:"original"`
			Thumb50X50   string `json:"thumb_50x50"`
			Thumb100X100 string `json:"thumb_100x100"`
		} `json:"avatar"`
		Timezone   string `json:"timezone"`
		Language   string `json:"language"`
		ProfileURL string `json:"profile_url"`
	} `json:"user"`
	DateAdded      int    `json:"date_added"`
	ReplyID        int    `json:"reply_id"`
	ThreadPosition string `json:"thread_position"`
	Karma          int    `json:"karma"`
	KarmaGuest     int    `json:"karma_guest"`
	Content        string `json:"content"`
}

// DeleteModComment deletes an existing mod comment
func DeleteModComment(commentID, modID, gameID int, user *User) (err error) {
	client := http.Client{Timeout: time.Duration(5 * time.Second)}
	req, err := http.NewRequest("DELETE", "https://api.mod.io/v1/games/"+strconv.Itoa(gameID)+"/mods/"+strconv.Itoa(modID)+"/comments", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return err
	}
	if user.OAuth2Token() != "" || len(user.OAuth2Token()) > 1 {
		req.Header.Set("Authorization", "Bearer "+user.OAuth2Token())
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 201 {
		var errObj ErrorCase
		e := json.Unmarshal(body, &errObj)
		if e != nil {
			log.Fatalln(e)
		}
		return HandleResponseError(errObj)
	}
	if err != nil {
		return err
	}
	return nil
}

// UpdateModComment updates an existing mod comment
func UpdateModComment(content string, commentID, modID, gameID int, user *User) (res *Comment, err error) {
	queryBody := url.Values{
		"content": {content},
	}
	client := http.Client{Timeout: time.Duration(5 * time.Second)}
	req, err := http.NewRequest("PUT", "https://api.mod.io/v1/games/"+strconv.Itoa(gameID)+"/mods/"+strconv.Itoa(modID)+"/comments", strings.NewReader(queryBody.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return nil, err
	}
	if user.OAuth2Token() != "" || len(user.OAuth2Token()) > 1 {
		req.Header.Set("Authorization", "Bearer "+user.OAuth2Token())
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		var errObj ErrorCase
		e := json.Unmarshal(body, &errObj)
		if e != nil {
			log.Fatalln(e)
		}
		return nil, HandleResponseError(errObj)
	}
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// AddModComment adds a mod comment
func (user *User) AddModComment(content string, modID, gameID int, options map[string]string) (res *Comment, err error) {
	var queryBody url.Values
	if options != nil {
		queryBody = ParseArgsBody(options)
	} else {
		queryBody = url.Values{
			"api_key": {user.APIKey()},
			"content": {content},
		}
	}
	client := http.Client{Timeout: time.Duration(5 * time.Second)}
	req, err := http.NewRequest("POST", "https://api.mod.io/v1/games/"+strconv.Itoa(gameID)+"/mods/"+strconv.Itoa(modID)+"/comments", strings.NewReader(queryBody.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return nil, err
	}
	if user.OAuth2Token() != "" || len(user.OAuth2Token()) > 1 {
		req.Header.Set("Authorization", "Bearer "+user.OAuth2Token())
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		var errObj ErrorCase
		e := json.Unmarshal(body, &errObj)
		if e != nil {
			log.Fatalln(e)
		}
		return nil, HandleResponseError(errObj)
	}
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetModComment searches for a mod comment specifically
func (user *User) GetModComment(commentID int, modID int, gameID int) (res *Comment, err error) {
	client := http.Client{Timeout: time.Duration(5 * time.Second)}
	req, err := http.NewRequest("GET", "https://api.mod.io/v1/games/"+strconv.Itoa(gameID)+"/mods/"+strconv.Itoa(modID)+"/comments/"+strconv.Itoa(commentID)+"?api_key="+user.APIKey(), nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return nil, err
	}
	if user.OAuth2Token() != "" || len(user.OAuth2Token()) > 1 {
		req.Header.Set("Authorization", "Bearer "+user.OAuth2Token())
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		var errObj ErrorCase
		e := json.Unmarshal(body, &errObj)
		if e != nil {
			log.Fatalln(e)
		}
		return nil, HandleResponseError(errObj)
	}
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetModComments searches for mod comments
func (user *User) GetModComments(modID int, gameID int, options map[string]string) (res *Comments, err error) {
	var queryStr string
	if options == nil {
		queryStr = "api_key=" + user.APIKey()
	} else {
		options["api_key"] = user.APIKey()
		queryStr = ParseArgsGet(options)
	}
	client := http.Client{Timeout: time.Duration(5 * time.Second)}
	req, err := http.NewRequest("GET", "https://api.mod.io/v1/games/"+strconv.Itoa(gameID)+"/mods/"+strconv.Itoa(modID)+"/comments"+"?"+queryStr, nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return nil, err
	}
	if user.OAuth2Token() != "" || len(user.OAuth2Token()) > 1 {
		req.Header.Set("Authorization", "Bearer "+user.OAuth2Token())
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 201 {
		var errObj ErrorCase
		e := json.Unmarshal(body, &errObj)
		if e != nil {
			log.Fatalln(e)
		}
		return nil, HandleResponseError(errObj)
	}
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
