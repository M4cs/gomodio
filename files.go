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
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Files struct which maps to the JSON of Get Files
type Files struct {
	Data         []File `json:"data"`
	ResultCount  int    `json:"result_count"`
	ResultLimit  int    `json:"result_limit"`
	ResultTotal  int    `json:"result_total"`
	ResultOffset int    `json:"result_offset"`
}

// File struct which maps to the JSON of Get/Add/Delete File
type File struct {
	ID             int    `json:"id"`
	ModID          int    `json:"mod_id"`
	DateAdded      int    `json:"date_added"`
	DateScanned    int    `json:"date_scanned"`
	VirusStatus    int    `json:"virus_status"`
	VirusPositive  int    `json:"virus_positive"`
	VirustotalHash string `json:"virustotal_hash"`
	Filesize       int    `json:"filesize"`
	Filehash       struct {
		Md5 string `json:"md5"`
	} `json:"filehash"`
	Filename     string `json:"filename"`
	Version      string `json:"version"`
	Changelog    string `json:"changelog"`
	MetadataBlob string `json:"metadata_blob"`
	Download     struct {
		BinaryURL   string `json:"binary_url"`
		DateExpires int    `json:"date_expires"`
	} `json:"download"`
}

// GetModfiles grabs modfiles and returns a Files struct
func GetModfiles(modID int, gameID int, options map[string]string, user *User) (f *Files, err error) {
	options["api_key"] = user.APIKey()
	queryStr := ParseArgsGet(options)
	client := http.Client{Timeout: time.Duration(5 * time.Second)}
	req, err := http.NewRequest("GET", "https://api.mod.io/v1/games/"+strconv.Itoa(gameID)+"/mods/"+strconv.Itoa(modID)+"?"+queryStr, nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return f, err
	}
	if user.OAuth2Token() != "" || len(user.OAuth2Token()) > 1 {
		req.Header.Set("Authorization", "Bearer "+user.OAuth2Token())
	}
	resp, err := client.Do(req)
	if err != nil {
		return f, err
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
		return f, err
	}
	err = json.Unmarshal(body, &f)
	if err != nil {
		return f, err
	}
	return f, nil
}

// GetModfile grabs a modfile and returns a File struct
func GetModfile(fileID int, modID int, gameID int, user *User) (f *File, err error) {
	client := http.Client{Timeout: time.Duration(5 * time.Second)}
	req, err := http.NewRequest("GET", "https://api.mod.io/v1/games/"+strconv.Itoa(gameID)+"/mods/"+strconv.Itoa(modID)+"?api_key="+user.APIKey(), nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return f, err
	}
	if user.OAuth2Token() != "" || len(user.OAuth2Token()) > 1 {
		req.Header.Set("Authorization", "Bearer "+user.OAuth2Token())
	}
	resp, err := client.Do(req)
	if err != nil {
		return f, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return f, err
	}
	if resp.StatusCode != 200 {
		var errObj ErrorCase
		e := json.Unmarshal(body, &errObj)
		if e != nil {
			log.Fatalln(e)
		}
		return nil, HandleResponseError(errObj)
	}
	err = json.Unmarshal(body, &f)
	if err != nil {
		return f, err
	}
	return f, nil
}

// EditModfile sends a PUT request to edit a mod file
func EditModfile(modID int, gameID int, options map[string]string, user *User) (f *File, err error) {
	if user.OAuth2Token() == "" {
		return f, errors.New("requires OAuth2 token")
	}
	client := http.Client{Timeout: time.Duration(5 * time.Second)}
	reqBody := ParseArgsBody(options)
	req, err := http.NewRequest("PUT", "https://api.mod.io/v1/games/"+strconv.Itoa(gameID)+"/mods/"+strconv.Itoa(modID)+"/files", strings.NewReader(reqBody.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return f, err
	}
	if user.OAuth2Token() != "" || len(user.OAuth2Token()) > 1 {
		req.Header.Set("Authorization", "Bearer "+user.OAuth2Token())
	}
	resp, err := client.Do(req)
	if err != nil {
		return f, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return f, err
	}
	if resp.StatusCode != 200 {
		var errObj ErrorCase
		e := json.Unmarshal(b, &errObj)
		if e != nil {
			log.Fatalln(e)
		}
		return nil, HandleResponseError(errObj)
	}

	err = json.Unmarshal(b, &f)
	if err != nil {
		return f, err
	}
	return f, nil
}

// DeleteModfile sends a DELETE request to delete a mod file
func DeleteModfile(fileID int, modID int, gameID int, user *User) (err error) {
	if user.OAuth2Token() == "" {
		return errors.New("requires OAuth2 token")
	}
	client := http.Client{Timeout: time.Duration(5 * time.Second)}
	req, err := http.NewRequest("DELETE", "https://api.mod.io/v1/games/"+strconv.Itoa(gameID)+"/mods/"+strconv.Itoa(modID)+"/files/"+strconv.Itoa(fileID), nil)
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

// AddModfile sends a POST request to upload a mod file
func AddModfile(modID int, gameID int, fp string, options map[string]string, user *User) (f *File, err error) {
	if user.OAuth2Token() == "" {
		return f, errors.New("requires OAuth2 token")
	}
	file, err := os.Open(fp)
	if err != nil {
		return f, err
	}
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("filedata", filepath.Base(fp))
	if err != nil {
		return f, err
	}
	_, err = io.Copy(part, file)
	if options != nil {
		for k, v := range options {
			_ = writer.WriteField(k, v)
		}
	}
	client := http.Client{Timeout: time.Duration(5 * time.Second)}
	req, err := http.NewRequest("POST", "https://api.mod.io/v1/games/"+strconv.Itoa(gameID)+"/mods/"+strconv.Itoa(modID)+"/files", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return f, err
	}
	if user.OAuth2Token() != "" || len(user.OAuth2Token()) > 1 {
		req.Header.Set("Authorization", "Bearer "+user.OAuth2Token())
	}
	resp, err := client.Do(req)
	if err != nil {
		return f, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return f, err
	}
	if resp.StatusCode != 201 {
		var errObj ErrorCase
		e := json.Unmarshal(b, &errObj)
		if e != nil {
			log.Fatalln(e)
		}
		return nil, HandleResponseError(errObj)
	}
	err = json.Unmarshal(b, &f)
	if err != nil {
		return f, err
	}
	return f, err
}
