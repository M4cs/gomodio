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

// GameStats struct represents a game's stats
type GameStats struct {
	GameID                    int `json:"game_id"`
	ModsCountTotal            int `json:"mods_count_total"`
	ModsDownloadsToday        int `json:"mods_downloads_today"`
	ModsDownloadsTotal        int `json:"mods_downloads_total"`
	ModsDownloadsDailyAverage int `json:"mods_downloads_daily_average"`
	ModsSubscribersTotal      int `json:"mods_subscribers_total"`
	DateExpires               int `json:"date_expires"`
}

// ModStats struct represents a group of stats of a mod
type ModStats struct {
	Data         []Stats `json:"data"`
	ResultCount  int     `json:"result_count"`
	ResultLimit  int     `json:"result_limit"`
	ResultTotal  int     `json:"result_total"`
	ResultOffset int     `json:"result_offset"`
}

// Stats struct represents a stats object
type Stats struct {
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
}

// GetGameStats gets a game's stats
func (u *User) GetGameStats(gameID int) (gs *GameStats, err error) {
	client := http.Client{Timeout: time.Duration(5 * time.Second)}
	req, err := http.NewRequest("GET", "https://api.mod.io/v1/games/"+strconv.Itoa(gameID)+"mods/stats?api_key="+u.APIKey(), nil)
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
	err = json.Unmarshal(body, &gs)
	if err != nil {
		return nil, err
	}
	return gs, nil
}

// GetModStats gets a mod's stats
func (u *User) GetModStats(modID, gameID int) (s *Stats, err error) {
	client := http.Client{Timeout: time.Duration(5 * time.Second)}
	req, err := http.NewRequest("GET", "https://api.mod.io/v1/games/"+strconv.Itoa(gameID)+"mods/"+strconv.Itoa(modID)+"/stats?api_key="+u.APIKey(), nil)
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
	err = json.Unmarshal(body, &s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// GetModsStats gets a game's mod's stats
func (u *User) GetModsStats(gameID int, options map[string]int) (ms *ModStats, err error) {
	var queryBody url.Values
	if options != nil {
		for k, v := range options {
			queryBody.Add(k, strconv.Itoa(v))
		}
	}
	queryBody.Add("api_key", u.APIKey())
	client := http.Client{Timeout: time.Duration(5 * time.Second)}
	req, err := http.NewRequest("GET", "https://api.mod.io/v1/games/"+strconv.Itoa(gameID)+"mods/stats", strings.NewReader(queryBody.Encode()))
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
	err = json.Unmarshal(body, &ms)
	if err != nil {
		return nil, err
	}
	return ms, nil
}
