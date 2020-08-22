package gomodio

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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

// GetModComments searches for mod comments
func GetModComments(modID int, gameID int, options map[string]string, user *User) (res *Comments, err error) {
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
