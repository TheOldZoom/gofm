package api

import "github.com/theOldZoom/gofm/internal/models"

func EnrichTrackImageFromPage(track *models.Track) error {
	if track == nil || track.Url == "" || !trackNeedsImageFallback(*track) {
		return nil
	}

	imageURL, err := GetPageImageURL(track.Url)
	if err != nil || imageURL == "" {
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
