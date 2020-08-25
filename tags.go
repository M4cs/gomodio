package gomodio

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Tag struct represents the tag object from mod.io
type Tag struct {
	Name      string `json:"name"`
	DateAdded int    `json:"date_added"`
}

// Tags struct is a collection of Tags
type Tags struct {
	Data         []Tag `json:"data"`
	ResultCount  int   `json:"result_count"`
	ResultLimit  int   `json:"result_limit"`
	ResultTotal  int   `json:"result_total"`
	ResultOffset int   `json:"result_offset"`
}

// GameTags struct is a game's tags object
type GameTags struct {
	Name   string   `json:"name"`
	Type   string   `json:"type"`
	Tags   []string `json:"string"`
	Hidden bool     `json:"bool"`
}

// DeleteGameTagOption deletes a game tag option
func (user *User) DeleteGameTagOption(tagGroupName string, tags []string, gameID int) (err error) {
	if user.OAuth2Token() == "" {
		return errors.New("requires oauth2 authentication")
	}
	var queryBody url.Values
	options := map[string]string{
		"tags": "[\"" + strings.Join(tags, "\",\"") + "\"]",
		"name": tagGroupName,
	}
	queryBody = ParseArgsBody(options)
	client := http.Client{Timeout: time.Duration(5 * time.Second)}
	req, err := http.NewRequest("DELETE", "https://api.mod.io/v1/games/"+strconv.Itoa(gameID)+"/tags", strings.NewReader(queryBody.Encode()))
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
	if resp.StatusCode != 204 {
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

// AddGameTagOption adds a single option to game tags
func (user *User) AddGameTagOption(tagGroupName, tagGroupType string, tags []string, gameID int, options map[string]string) (m *Message, err error) {
	if user.OAuth2Token() == "" {
		return nil, errors.New("requires oauth2 authentication")
	}
	var queryBody url.Values
	if tagGroupType != "dropdown" {
		return nil, errors.New("tagGroupType must be: dropdown or checkboxes")
	} else if tagGroupType != "checkboxes" {
		return nil, errors.New("tagGroupType must be: dropdown or checkboxes")
	}
	if options != nil {
		options["name"] = tagGroupName
		options["tagGroupType"] = tagGroupType
		queryBody = ParseArgsBody(options)
	} else {
		options = map[string]string{
			"name": tagGroupName,
			"type": tagGroupType,
		}
		queryBody = ParseArgsBody(options)
	}
	if tags != nil {
		queryBody.Add("tags", "[\""+strings.Join(tags, "\",\"")+"\"]")
	}
	client := http.Client{Timeout: time.Duration(5 * time.Second)}
	req, err := http.NewRequest("POST", "https://api.mod.io/v1/games/"+strconv.Itoa(gameID)+"/tags", nil)
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
	err = json.Unmarshal(body, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// GetGameTagOptions gets a game's tag options
func (user *User) GetGameTagOptions(gameID int) (t *Tags, err error) {
	client := http.Client{Timeout: time.Duration(5 * time.Second)}
	req, err := http.NewRequest("GET", "https://api.mod.io/v1/games/"+strconv.Itoa(gameID)+"/tags", nil)
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
	err = json.Unmarshal(body, &t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// DeleteModTags deletes a tag from a mod. Requires OAuth2
func (user *User) DeleteModTags(tags []string, modID, gameID int) (err error) {
	var queryBody url.Values
	if user.OAuth2Token() == "" {
		return errors.New("requires oauth2 authentication")
	}
	queryBody.Add("tags", "[\""+strings.Join(tags, "\",\"")+"\"]")
	client := http.Client{Timeout: time.Duration(5 * time.Second)}
	req, err := http.NewRequest("DELETE", "https://api.mod.io/v1/games/"+strconv.Itoa(gameID)+"/mods/"+strconv.Itoa(modID)+"/tags", strings.NewReader(queryBody.Encode()))
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
	if resp.StatusCode != 204 {
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

// AddModTags adds a tag to a mod. Requires OAuth2
func (user *User) AddModTags(tags []string, modID, gameID int) (t *Message, err error) {
	var queryBody url.Values
	if user.OAuth2Token() == "" {
		return nil, errors.New("requires oauth2 authentication")
	}
	for _, t := range tags {
		queryBody.Add("tags", t)
	}
	client := http.Client{Timeout: time.Duration(5 * time.Second)}
	req, err := http.NewRequest("POST", "https://api.mod.io/v1/games/"+strconv.Itoa(gameID)+"/mods/"+strconv.Itoa(modID)+"/tags", strings.NewReader(queryBody.Encode()))
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
	err = json.Unmarshal(body, &t)
	if err != nil {
		return nil, err
	}
	return t, nil

}

// GetModTags grabs tags from a mod
func (user *User) GetModTags(modID, gameID int, options map[string]string) (t *Tags, err error) {
	var queryStr string
	if options != nil {
		options["api_key"] = user.APIKey()
		queryStr = ParseArgsGet(options)
	} else {
		queryStr = "api_key=" + user.APIKey()
	}
	client := http.Client{Timeout: time.Duration(5 * time.Second)}
	req, err := http.NewRequest("GET", "https://api.mod.io/v1/games/"+strconv.Itoa(gameID)+"/mods/"+strconv.Itoa(modID)+"/tags?"+queryStr, nil)
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
	err = json.Unmarshal(body, &t)
	if err != nil {
		return nil, err
	}
	return t, nil
}
