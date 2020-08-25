package gomodio

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Events struct represents the events object of mod.io's API
type Events struct {
	Data         []Event `json:"data"`
	ResultCount  int     `json:"result_count"`
	ResultOffset int     `json:"result_offset"`
	ResultLimit  int     `json:"result_limit"`
	ResultTotal  int     `json:"result_total"`
}

// Event struct represents the event object of mod.io's API
type Event struct {
	ID        int    `json:"id"`
	ModID     int    `json:"mod_id"`
	UserID    int    `json:"user_id"`
	DateAdded int    `json:"date_added"`
	EventType string `json:"event_type"`
}

// GetModsEvents gets all mods events
func (user *User) GetModsEvents(gameID int, options map[string]string) (e *Events, err error) {
	var queryStr string
	if options != nil {
		options["api_key"] = user.APIKey()
		queryStr = ParseArgsGet(options)
	} else {
		queryStr = "api_key=" + user.APIKey()
	}
	client := http.Client{Timeout: time.Duration(5 * time.Second)}
	req, err := http.NewRequest("GET", "https://api.mod.io/v1/games/"+strconv.Itoa(gameID)+"/mods/events?"+queryStr, nil)
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
	err = json.Unmarshal(body, &e)
	if err != nil {
		return nil, err
	}
	return e, nil
}

// GetModEvents gets a single mod's events
func (user *User) GetModEvents(gameID int, modID int) (e *Events, err error) {
	client := http.Client{Timeout: time.Duration(5 * time.Second)}
	req, err := http.NewRequest("GET", "https://api.mod.io/v1/games/"+strconv.Itoa(gameID)+"/mods/"+strconv.Itoa(modID)+"/events", nil)
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
	err = json.Unmarshal(body, &e)
	if err != nil {
		return nil, err
	}
	return e, nil
}
