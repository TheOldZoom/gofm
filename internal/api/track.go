package api

import (
	"github.com/theOldZoom/gofm/internal/models"
	"github.com/theOldZoom/gofm/internal/verbose"
)

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
