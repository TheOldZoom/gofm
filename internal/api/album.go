package api

import (
	"github.com/theOldZoom/gofm/internal/models"
	"github.com/theOldZoom/gofm/internal/verbose"
)

func EnrichAlbumImageFromPage(album *models.Album) error {
	if album == nil || album.URL == "" {
		return nil
	}

	imageURL, err := GetPageImageURL(album.URL)
	if err != nil || imageURL == "" {
		if err != nil {
			verbose.Printf("album image enrichment failed for %s: %v", album.Name, err)
		}
		return err
	}

	album.Image = []struct {
		Size string `json:"size"`
		Url  string `json:"#text"`
	}{
		{
			Size: "page",
			Url:  imageURL,
		},
	}

	verbose.Printf("album image enriched from page: %s -> %s", album.Name, imageURL)
	return nil
}
