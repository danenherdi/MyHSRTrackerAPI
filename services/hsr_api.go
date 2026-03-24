package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

type GachaLogResponse struct {
	Retcode int    `json:"retcode"`
	Message string `json:"message"`
	Data    *struct {
		Page   string `json:"page"`
		Size   string `json:"size"`
		List   []struct {
			UID       string `json:"uid"`
			GachaType string `json:"gacha_type"`
			ItemID    string `json:"item_id"`
			Count     string `json:"count"`
			Time      string `json:"time"`
			Name      string `json:"name"`
			Lang      string `json:"lang"`
			ItemType  string `json:"item_type"`
			RankType  string `json:"rank_type"`
			ID        string `json:"id"`
		} `json:"list"`
		Region   string `json:"region"`
		RegionTZ int    `json:"region_time_zone"`
	} `json:"data"`
}

// ExtractQueryParams parses the user's gacha URL and returns the query parameters
func ExtractQueryParams(rawURL string) (url.Values, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, errors.New("invalid url format")
	}
	
	q := u.Query()
	if q.Get("authkey") == "" {
		return nil, errors.New("authkey not found in url")
	}

	return q, nil
}

// FetchGachaLog makes a request to HoYoverse API
func FetchGachaLog(q url.Values, gachaType string, endId string) (*GachaLogResponse, error) {
	// Read Base URL from env
	baseURL := os.Getenv("HSR_API_URL")

	// Set specific pagination params
	q.Set("gacha_type", gachaType)
	q.Set("size", "20")
	q.Set("end_id", endId)

	// Build request URL
	reqURL := fmt.Sprintf("%s?%s", baseURL, q.Encode())

	// Create HTTP request
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}

	// HoYoverse APIs often require some generic headers, or at least a standard User-Agent
	req.Header.Set("User-Agent", "Mozilla/5.0")

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code: %d", res.StatusCode)
	}

	var gachaRes GachaLogResponse
	if err := json.NewDecoder(res.Body).Decode(&gachaRes); err != nil {
		return nil, errors.New("failed to decode response")
	}

	return &gachaRes, nil
}
