package gomodio

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Games struct which maps to the JSON response of Games/Edit in Get Game/s
type Games struct {
	Data        []Game `json:"data"`
	ResultCount int    `json:"result_count"`
	ResultLimit int    `json:"result_limit"`
	ResultTotal int    `json:"result_total"`
}

// Game struct which maps to the JSON response of Get/Edit Game/s
type Game struct {
	ID          int `json:"id"`
	Status      int `json:"status"`
	SubmittedBy struct {
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
	} `json:"submitted_by"`
	DateAdded          int    `json:"date_added"`
	DateUpdated        int    `json:"date_updated"`
	DateLive           int    `json:"date_live"`
	PresentationOption int    `json:"presentation_option"`
	SubmissionOption   int    `json:"submission_option"`
	CurationOption     int    `json:"curation_option"`
	CommunityOptions   int    `json:"community_options"`
	RevenueOptions     int    `json:"revenue_options"`
	APIAccessOptions   int    `json:"api_access_options"`
	MaturityOptions    int    `json:"maturity_options"`
	UgcName            string `json:"ugc_name"`
	Icon               struct {
		Filename     string `json:"filename"`
		Original     string `json:"original"`
		Thumb64X64   string `json:"thumb_64x64"`
		Thumb128X128 string `json:"thumb_128x128"`
		Thumb256X256 string `json:"thumb_256x256"`
	} `json:"icon"`
	Logo struct {
		Filename      string `json:"filename"`
		Original      string `json:"original"`
		Thumb320X180  string `json:"thumb_320x180"`
		Thumb640X360  string `json:"thumb_640x360"`
		Thumb1280X720 string `json:"thumb_1280x720"`
	} `json:"logo"`
	Header struct {
		Filename string `json:"filename"`
		Original string `json:"original"`
	} `json:"header"`
	Name            string `json:"name"`
	NameID          string `json:"name_id"`
	Summary         string `json:"summary"`
	Instructions    string `json:"instructions"`
	InstructionsURL string `json:"instructions_url"`
	ProfileURL      string `json:"profile_url"`
	TagOptions      []struct {
		Name   string   `json:"name"`
		Type   string   `json:"type"`
		Tags   []string `json:"tags"`
		Hidden bool     `json:"hidden"`
	} `json:"tag_options"`
}

// GetGames from mod.io
func (user *User) GetGames(query map[string]string) (res *Games, err error) {
	query["api_key"] = user.APIKey()
	queryString := ParseArgsGet(query)
	client := http.Client{Timeout: time.Duration(5 * time.Second)}
	req, err := http.NewRequest("GET", "https://api.mod.io/v1/games?"+queryString, nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return res, err
	}
	if user.OAuth2Token() != "" || len(user.OAuth2Token()) > 1 {
		req.Header.Set("Authorization", "Bearer "+user.OAuth2Token())
	}
	resp, err := client.Do(req)
	if err != nil {
		return res, err
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
		return res, err
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return res, err
	}
	return res, nil

}

// EditGame function makes a PUT request and returns the updated Game Object
func (user *User) EditGame(gameID int, query map[string]string) (res *Game, err error) {
	if user.OAuth2Token() != "" || len(user.OAuth2Token()) > 1 {
		return nil, errors.New("requires OAuth2 token")
	}
	query["api_key"] = user.APIKey()
	reqBody := ParseArgsBody(query)
	client := http.Client{Timeout: time.Duration(5 * time.Second)}
	req, err := http.NewRequest("PUT", "https://api.mod.io/v1/games/"+strconv.Itoa(gameID), strings.NewReader(reqBody.Encode()))
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

// GetGame function returns a Game struct
func (user *User) GetGame(gameID int, query map[string]string) (res *Game, err error) {
	var queryString string
	if query != nil {
		query["api_key"] = user.APIKey()
		queryString = ParseArgsGet(query)
	} else {
		queryString = "api_key=" + user.APIKey()
	}
	client := http.Client{Timeout: time.Duration(5 * time.Second)}
	req, err := http.NewRequest("GET", "https://api.mod.io/v1/games/"+strconv.Itoa(gameID)+"?"+queryString, nil)
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

// ToJSON returns JSON string of Games struct
func (g *Games) ToJSON() (jsonStr string, err error) {
	jsonString, err := json.Marshal(g)
	if err != nil {
		return "", err
	}
	return string(jsonString), nil
}

// ToJSON returns JSON string of Game struct
func (g *Game) ToJSON() (jsonStr string, err error) {
	jsonString, err := json.Marshal(g)
	if err != nil {
		return "", err
	}
	return string(jsonString), nil
}
