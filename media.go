package gomodio

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Message struct represents a message object in JSON
type Message struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// DeleteModMedia deletes mod media
func (user *User) DeleteModMedia(modID, gameID int, options map[string]string) (err error) {
	var reqBody url.Values
	if options == nil {
		return errors.New("must provide options. cannot be nil")
	}
	for k, v := range options {
		reqBody.Add(k, v)
	}
	client := http.Client{Timeout: time.Duration(5 * time.Second)}
	req, err := http.NewRequest("POST", "https://api.mod.io/v1/games/"+strconv.Itoa(gameID)+"/mods/"+strconv.Itoa(modID)+"/media", strings.NewReader(reqBody.Encode()))
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
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode != 204 {
		var errObj ErrorCase
		e := json.Unmarshal(b, &errObj)
		if e != nil {
			log.Fatalln(e)
		}
		return HandleResponseError(errObj)
	}
	return nil
}

// AddModMedia adds mod media
func (user *User) AddModMedia(modID, gameID int, options map[string]string) (msg *Message, err error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	if options != nil {
		for k, v := range options {
			if k == "logo" || k == "images" {
				file, err := os.Open(v)
				if err != nil {
					return nil, err
				}
				defer file.Close()
				part, err := writer.CreateFormFile(k, filepath.Base(v))
				if err != nil {
					return nil, err
				}
				_, err = io.Copy(part, file)
			} else {
				_ = writer.WriteField(k, v)
			}
		}
	}
	client := http.Client{Timeout: time.Duration(5 * time.Second)}
	req, err := http.NewRequest("POST", "https://api.mod.io/v1/games/"+strconv.Itoa(gameID)+"/mods/"+strconv.Itoa(modID)+"/media", body)
	req.Header.Set("Content-Type", "multipart/form-data")
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
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		var errObj ErrorCase
		e := json.Unmarshal(b, &errObj)
		if e != nil {
			log.Fatalln(e)
		}
		return nil, HandleResponseError(errObj)
	}
	err = json.Unmarshal(b, &msg)
	if err != nil {
		return nil, err
	}
	return msg, err
}

// AddGameMedia adds game media
func (user *User) AddGameMedia(logo, icon, header string, gameID int) (msg *Message, err error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	file, err := os.Open(logo)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	part, err := writer.CreateFormFile("logo", filepath.Base(logo))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	file1, err := os.Open(icon)
	if err != nil {
		return nil, err
	}
	defer file1.Close()
	part1, err := writer.CreateFormFile("icon", filepath.Base(icon))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part1, file1)
	file2, err := os.Open(header)
	if err != nil {
		return nil, err
	}
	defer file2.Close()
	part2, err := writer.CreateFormFile("header", filepath.Base(header))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part2, file2)
	client := http.Client{Timeout: time.Duration(5 * time.Second)}
	req, err := http.NewRequest("POST", "https://api.mod.io/v1/games/"+strconv.Itoa(gameID)+"/media", body)
	req.Header.Set("Content-Type", "multipart/form-data")
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
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		var errObj ErrorCase
		e := json.Unmarshal(b, &errObj)
		if e != nil {
			log.Fatalln(e)
		}
		return nil, HandleResponseError(errObj)
	}
	err = json.Unmarshal(b, &msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}
