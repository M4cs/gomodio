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

// Mods struct which maps to the JSON response of Get Mods
type Mods struct {
	Data         []Mod `json:"data"`
	ResultCount  int   `json:"result_count"`
	ResultLimit  int   `json:"result_limit"`
	ResultTotal  int   `json:"result_total"`
	ResultOffset int   `json:"result_offset"`
}

// Mod struct which maps to the JSON response of Get/Edit/Add/Delete Mod/s
type Mod struct {
	ID          int `json:"id"`
	GameID      int `json:"game_id"`
	Status      int `json:"status"`
	Visible     int `json:"visible"`
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
	DateAdded      int `json:"date_added"`
	DateUpdated    int `json:"date_updated"`
	DateLive       int `json:"date_live"`
	MaturityOption int `json:"maturity_option"`
	Logo           struct {
		Filename      string `json:"filename"`
		Original      string `json:"original"`
		Thumb320X180  string `json:"thumb_320x180"`
		Thumb640X360  string `json:"thumb_640x360"`
		Thumb1280X720 string `json:"thumb_1280x720"`
	} `json:"logo"`
	HomepageURL          string `json:"homepage_url"`
	Name                 string `json:"name"`
	NameID               string `json:"name_id"`
	Summary              string `json:"summary"`
	Description          string `json:"description"`
	DescriptionPlaintext string `json:"description_plaintext"`
	MetadataBlob         string `json:"metadata_blob"`
	ProfileURL           string `json:"profile_url"`
	Media                struct {
		Youtube   []string `json:"youtube"`
		Sketchfab []string `json:"sketchfab"`
		Images    []struct {
			Filename     string `json:"filename"`
			Original     string `json:"original"`
			Thumb320X180 string `json:"thumb_320x180"`
		} `json:"images"`
	} `json:"media"`
	Modfile struct {
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
	} `json:"modfile"`
	MetadataKvp []struct {
		Metakey   string `json:"metakey"`
		Metavalue string `json:"metavalue"`
	} `json:"metadata_kvp"`
	Tags []struct {
		Name      string `json:"name"`
		DateAdded int    `json:"date_added"`
	} `json:"tags"`
	Stats struct {
		ModID                     int     `json:"mod_id"`
		PopularityRankPosition    int     `json:"popularity_rank_position"`
		PopularityRankTotalMods   int     `json:"popularity_rank_total_mods"`
		DownloadsTotal            int     `json:"downloads_total"`
		SubscribersTotal          int     `json:"subscribers_total"`
		RatingsTotal              int     `json:"ratings_total"`
		RatingsPositive           int     `json:"ratings_positive"`
		RatingsNegative           int     `json:"ratings_negative"`
		RatingsPercentagePositive int     `json:"ratings_percentage_positive"`
		RatingsWeightedAggregate  float64 `json:"ratings_weighted_aggregate"`
		RatingsDisplayText        string  `json:"ratings_display_text"`
		DateExpires               int     `json:"date_expires"`
	} `json:"stats"`
}

// GetMods searches for mods and returns a Mods object
func GetMods(gameID int, query map[string]string, user *User) (res *Mods, err error) {
	query["api_key"] = user.APIKey()
	queryString := ParseArgsGet(query)
	client := http.Client{Timeout: time.Duration(5 * time.Second)}
	req, err := http.NewRequest("GET", "https://api.mod.io/v1/games/"+strconv.Itoa(gameID)+"/mods?"+queryString, nil)
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
	if err != nil {
		return res, err
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return res, err
	}
	return res, nil
}

// EditMod edits a mod
func EditMod(modID int, gameID int, options map[string]string, user *User) (res *Mod, err error) {
	if user.OAuth2Token() == "" {
		return res, errors.New("requires OAuth2 token")
	}
	reqBody := ParseArgsBody(options)
	client := http.Client{Timeout: time.Duration(5 * time.Second)}
	req, err := http.NewRequest("PUT", "https://api.mod.io/v1/games/"+strconv.Itoa(gameID)+"/mods/"+strconv.Itoa(modID), strings.NewReader(reqBody.Encode()))
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
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}
	err = json.Unmarshal(b, &res)
	if err != nil {
		return res, err
	}
	return res, nil
}

// DeleteMod sends a request to delete a mod
func DeleteMod(modID int, gameID int, user *User) (err error) {
	client := http.Client{Timeout: time.Duration(5 * time.Second)}
	req, err := http.NewRequest("DELETE", "https://api.mod.io/v1/games/"+strconv.Itoa(gameID)+"/mods/"+strconv.Itoa(modID), nil)
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
	if err != nil {
		return err
	}
	if resp.StatusCode != 204 {
		var errObj ErrorCase
		e := json.Unmarshal(body, &errObj)
		if e != nil {
			log.Fatalln(e)
		}
		return HandleResponseError(errObj)
	}
	return nil
}

// AddMod adds a mod taking bytes for files and returns Mod object
func AddMod(logo string, modName string, summary string, options map[string]string, gameID int, user *User) (res *Mod, err error) {
	if user.OAuth2Token() == "" {
		return res, errors.New("requires OAuth2 token")
	}
	file, err := os.Open(logo)
	if err != nil {
		return res, err
	}
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("logo", filepath.Base(logo))
	if err != nil {
		return res, err
	}
	_, err = io.Copy(part, file)
	_ = writer.WriteField("name", modName)
	_ = writer.WriteField("summary", summary)
	if options != nil {
		for k, v := range options {
			_ = writer.WriteField(k, v)
		}
	}
	client := http.Client{Timeout: time.Duration(5 * time.Second)}
	req, err := http.NewRequest("POST", "https://api.mod.io/v1/games/"+strconv.Itoa(gameID)+"/mods", body)
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
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}
	if resp.StatusCode != 200 {
		var errObj ErrorCase
		e := json.Unmarshal(b, &errObj)
		if e != nil {
			log.Fatalln(e)
		}
		return nil, HandleResponseError(errObj)
	}
	err = json.Unmarshal(b, &res)
	if err != nil {
		return res, err
	}
	return res, nil
}

// GetMod searches for a mod and returns a Mods object
func GetMod(modID int, gameID int, query map[string]string, user *User) (res *Mods, err error) {
	var queryString string
	if query != nil {
		query["api_key"] = user.APIKey()
		queryString = ParseArgsGet(query)
	} else {
		queryString = "api_key=" + user.APIKey()
	}
	client := http.Client{Timeout: time.Duration(5 * time.Second)}
	req, err := http.NewRequest("GET", "https://api.mod.io/v1/games/"+strconv.Itoa(gameID)+"/mods/"+strconv.Itoa(gameID)+"?"+queryString, nil)
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
	if err != nil {
		return res, err
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return res, err
	}
	return res, err
}
