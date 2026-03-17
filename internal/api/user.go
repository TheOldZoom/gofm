package api

import (
	"fmt"

	"github.com/theOldZoom/gofm/internal/models"

	"github.com/spf13/viper"
)

func GetRecentTracks(username string, limit int) ([]models.Track, error) {
	client := &Client{
		ApiKey: viper.GetString("api_key"),
	}
	var resp models.RecentTracksResponse

	err := client.Get("user.getRecentTracks", map[string]string{
		"user":  username,
		"limit": fmt.Sprintf("%d", limit),
	}, &resp)
	if err != nil {
		return nil, err
	}

	tracks := resp.RecentTracks.Track
	if len(tracks) > limit {
		tracks = tracks[:limit]
	}

	return tracks, nil
}

func GetNowPlaying(username string) (*models.Track, error) {
	client := &Client{
		ApiKey: viper.GetString("api_key"),
	}
	var resp models.RecentTracksResponse

	err := client.Get("user.getRecentTracks", map[string]string{
		"user":  username,
		"limit": "1",
	}, &resp)
	if err != nil {
		return nil, err
	}

	tracks := resp.RecentTracks.Track
	if len(tracks) > 0 && tracks[0].Attr.NowPlaying == "true" {
		return &tracks[0], nil
	}

	return nil, nil
}

func GetInfo(username string) (*models.UserGetInfoResponse, error) {
	client := &Client{
		ApiKey: viper.GetString("api_key"),
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
	client := &Client{
		ApiKey: apiKey,
	}
	var resp models.UserGetInfoResponse

	err := client.Get("user.getInfo", map[string]string{
		"user": username,
	}, &resp)
	if err != nil {
		return err
	}

	return nil
}
