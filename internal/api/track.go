package api

import (
	"github.com/spf13/viper"
	"github.com/theOldZoom/gofm/internal/models"
	"github.com/theOldZoom/gofm/internal/verbose"
)

func EnrichTrackAlbumFromAPI(track *models.Track) error {
	if track == nil || track.Name == "" || track.Artist.Name == "" || track.Album.Name != "" {
		return nil
	}

	client := &Client{
		ApiKey: viper.GetString("api_key"),
	}
	var resp models.TrackGetInfoResponse

	err := client.Get("track.getInfo", map[string]string{
		"artist":      track.Artist.Name,
		"track":       track.Name,
		"autocorrect": "1",
	}, &resp)
	if err != nil {
		return err
	}

	if resp.Track.Album.Name == "" {
		return nil
	}

	track.Album = resp.Track.Album
	verbose.Printf("track album enriched from api: %s - %s -> %s", track.Artist.Name, track.Name, track.Album.Name)
	return nil
}

func EnrichTrackImageFromPage(track *models.Track) error {
	if track == nil || track.Url == "" || !trackNeedsImageFallback(*track) {
		return nil
	}

	imageURL, err := GetPageImageURL(track.Url)
	if err != nil || imageURL == "" {
		if err != nil {
			verbose.Printf("track image enrichment failed for %s: %v", track.Name, err)
		}
		return err
	}

	track.Image = []struct {
		Size string `json:"size"`
		Url  string `json:"#text"`
	}{
		{
			Size: "page",
			Url:  imageURL,
		},
	}

	verbose.Printf("track image enriched from page: %s -> %s", track.Name, imageURL)
	return nil
}

func trackNeedsImageFallback(track models.Track) bool {
	if len(track.Image) == 0 {
		return true
	}

	for _, image := range track.Image {
		if image.Url == "" {
			continue
		}

		return isPlaceholderImageURL(image.Url)
	}

	return true
}
