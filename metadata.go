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

// ModKVP represents a mod's KVP metadata
type ModKVP struct {
	Metakey   string `json:"metakey"`
	Metavalue string `json:"metavalue"`
}

// ModMetadata respesents multiple KVP metadata objects
type ModMetadata struct {
	Data         []ModKVP `json:"data"`
	ResultCount  int      `json:"result_count"`
	ResultLimit  int      `json:"result_limit"`
	ResultTotal  int      `json:"result_total"`
	ResultOffset int      `json:"result_offset"`
}

// DeleteModMetadata deletes a mod's metadata
func (u *User) DeleteModMetadata(metadata []string, modID, gameID int) (err error) {
	reqBody := url.Values{
		"metadata": {"[\"" + strings.Join(metadata, "\",\"") + "\"]"},
	}
	client := http.Client{Timeout: time.Duration(5 * time.Second)}
	req, err := http.NewRequest("DELETE", "https://api.mod.io/v1/games/"+strconv.Itoa(gameID)+"mods/"+strconv.Itoa(modID)+"/metadetakvp?api_key="+u.APIKey(), strings.NewReader(reqBody.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return err
	}
	if u.OAuth2Token() != "" || len(u.OAuth2Token()) > 1 {
		req.Header.Set("Authorization", "Bearer "+u.OAuth2Token())
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

// AddModMetadata adds metadata to a mod
func (u *User) AddModMetadata(metadata []string, modID, gameID int) (m *Message, err error) {
	reqBody := url.Values{
		"metadata": {"[\"" + strings.Join(metadata, "\",\"") + "\"]"},
	}
	client := http.Client{Timeout: time.Duration(5 * time.Second)}
	req, err := http.NewRequest("POST", "https://api.mod.io/v1/games/"+strconv.Itoa(gameID)+"mods/"+strconv.Itoa(modID)+"/metadetakvp?api_key="+u.APIKey(), strings.NewReader(reqBody.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return nil, err
	}
	if u.OAuth2Token() != "" || len(u.OAuth2Token()) > 1 {
		req.Header.Set("Authorization", "Bearer "+u.OAuth2Token())
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

// GetModMetadata gets a mod's metadata
func (u *User) GetModMetadata(modID, gameID int) (mm *ModMetadata, err error) {
	client := http.Client{Timeout: time.Duration(5 * time.Second)}
	req, err := http.NewRequest("GET", "https://api.mod.io/v1/games/"+strconv.Itoa(gameID)+"mods/"+strconv.Itoa(modID)+"/metadetakvp?api_key="+u.APIKey(), nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return nil, err
	}
	if u.OAuth2Token() != "" || len(u.OAuth2Token()) > 1 {
		req.Header.Set("Authorization", "Bearer "+u.OAuth2Token())
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
	err = json.Unmarshal(body, &mm)
	if err != nil {
		return nil, err
	}
	return mm, nil
}
