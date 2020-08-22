package gomodio

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Subscribe Struct Maps to JSON Response for Subscribing
type Subscribe struct {
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

// SubscribeToMod sends a request to subscribe to a mod
func SubscribeToMod(modID, gameID int, user *User) (s *Subscribe, err error) {
	if user.OAuth2Token() == "" {
		return s, errors.New("requires OAuth2 token")
	}
	client := http.Client{Timeout: time.Duration(5 * time.Second)}

	req, err := http.NewRequest("POST", "https://api.mod.io/v1/games/"+strconv.Itoa(gameID)+"/mods/"+strconv.Itoa(modID)+"/subscribe", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+user.OAuth2Token())
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 201 {
		var errObj ErrorCase
		e := json.Unmarshal(body, &errObj)
		if e != nil {
			log.Fatalln(e)
		}
		return nil, HandleResponseError(errObj)
	}
	err = json.Unmarshal(body, &s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// UnsubscribeToMod sends a request to subscribe to a mod
func UnsubscribeToMod(modID, gameID int, user *User) (err error) {
	if user.OAuth2Token() == "" {
		return errors.New("requires OAuth2 token")
	}
	client := http.Client{Timeout: time.Duration(5 * time.Second)}

	req, err := http.NewRequest("DELETE", "https://api.mod.io/v1/games/"+strconv.Itoa(gameID)+"/mods/"+strconv.Itoa(modID)+"/subscribe", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+user.OAuth2Token())
	if err != nil {
		return err
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
	if resp.StatusCode != 201 {
		var errObj ErrorCase
		e := json.Unmarshal(body, &errObj)
		if e != nil {
			log.Fatalln(e)
		}
		return HandleResponseError(errObj)
	}
	return nil
}
