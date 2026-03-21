package api

import (
	"github.com/spf13/viper"
	"github.com/theOldZoom/gofm/internal/models"
)

func GetAlbumInfo(artistName string, albumName string, username string) (*models.Album, error) {
	client := &Client{
		ApiKey: viper.GetString("api_key"),
	}
	var resp models.AlbumGetInfoResponse

	params := map[string]string{
		"artist": artistName,
		"album":  albumName,
	}
	if username != "" {
		params["username"] = username
	}
	err := client.Get("album.getInfo", params, &resp)
	if err != nil {
		return nil, err
	}

	if resp.Album.Name == "" {
		return nil, nil
	}

	return &resp.Album, nil
}
