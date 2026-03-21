package api

import (
	"fmt"

	"github.com/theOldZoom/gofm/internal/models"
	"github.com/theOldZoom/gofm/internal/verbose"

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

	enrichTracksConcurrently("recent", tracks, false)

	verbose.Printf("fetched %d recent tracks for %s", len(tracks), username)
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
		if err := EnrichTrackImageFromPage(&tracks[0]); err != nil {
			verbose.Printf("now playing image fallback failed for %s: %v", tracks[0].Name, err)
		}
		verbose.Printf("now playing track found for %s: %s", username, tracks[0].Name)
		return &tracks[0], nil
	}

	verbose.Printf("no now playing track for %s", username)
	return nil, nil
}

func GetUserTopAlbums(username string, limit int) ([]models.Album, error) {
	client := &Client{
		ApiKey: viper.GetString("api_key"),
	}
	var resp models.TopAlbumsResponse

	err := client.Get("user.getTopAlbums", map[string]string{
		"user":  username,
		"limit": fmt.Sprintf("%d", limit),
	}, &resp)
	if err != nil {
		return nil, err
	}
	albums := resp.TopAlbums.Album
	if len(albums) > limit {
		albums = albums[:limit]
	}

	enrichAlbumsConcurrently("top", albums)

	verbose.Printf("fetched %d top albums for %s", len(albums), username)
	return albums, nil
}

func GetUserTopArtists(username string, limit int) ([]models.Artist, error) {
	client := &Client{
		ApiKey: viper.GetString("api_key"),
	}
	var resp models.TopArtistsResponse

	err := client.Get("user.getTopArtists", map[string]string{
		"user":  username,
		"limit": fmt.Sprintf("%d", limit),
	}, &resp)
	if err != nil {
		return nil, err
	}

	artists := resp.TopArtists.Artist
	if len(artists) > limit {
		artists = artists[:limit]
	}

	enrichArtistsConcurrently("top", artists)

	verbose.Printf("fetched %d top artists for %s", len(artists), username)
	return artists, nil
}

func GetUserTopTracks(username string, limit int) ([]models.Track, error) {
	client := &Client{
		ApiKey: viper.GetString("api_key"),
	}
	var resp models.TopTracksResponse

	err := client.Get("user.getTopTracks", map[string]string{
		"user":  username,
		"limit": fmt.Sprintf("%d", limit),
	}, &resp)
	if err != nil {
		return nil, err
	}

	tracks := resp.TopTracks.Track
	if len(tracks) > limit {
		tracks = tracks[:limit]
	}

	enrichTracksConcurrently("top", tracks, true)

	verbose.Printf("fetched %d top tracks for %s", len(tracks), username)
	return tracks, nil
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

	verbose.Printf("fetched user info for %s", username)
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
