package api

import (
	"gofm/internal/models"
	"strconv"

	"github.com/spf13/viper"
)

func GetRecentTracks(username string, limit int) ([]models.Track, error) {
	client := &Client{
		ApiKey: viper.GetString("api_key"),
	}
	var resp models.RecentTracksResponse

	err := client.Get("user.getRecentTracks", map[string]string{
		"user":  username,
		"limit": strconv.Itoa(limit),
	}, &resp)
	if err != nil {
		return nil, err
	}
	return resp.RecentTracks.Track, nil
}

func GetInfo(username string, apiKey string) (*models.UserGetInfoResponse, error) {
	client := &Client{
		ApiKey: apiKey,
	}
	var resp models.UserGetInfoResponse

	err := client.Get("user.getInfo", map[string]string{
		"user": username,
	}, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func ValidateAPIKey(apiKey string) error {
	client := &Client{
		ApiKey: apiKey,
	}

	var resp struct {
		Tracks struct {
			Track []struct {
				Name string `json:"name"`
			} `json:"track"`
		} `json:"tracks"`
	}

	return client.Get("chart.getTopTracks", map[string]string{
		"limit": "1",
	}, &resp)
}

func ValidateUsername(username string, apiKey string) error {
	_, err := GetInfo(username, apiKey)
	return err
}
