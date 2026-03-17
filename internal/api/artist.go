package api

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/theOldZoom/gofm/internal/models"
)

const lastFMPlaceholderImageID = "2a96cbd8b46e442fc41c2b86b821562f"

var lastFMOGImagePattern = regexp.MustCompile(`<meta[^>]+property="og:image"[^>]+content="([^"]+)"`)

func GetPageImageURL(pageURL string) (string, error) {
	for range 2 {
		req, err := http.NewRequest(http.MethodGet, pageURL, nil)
		if err != nil {
			return "", err
		}
		setBrowserHeaders(req, "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return "", err
		}

		body, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			return "", readErr
		}
		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("artist page returned status code %d", resp.StatusCode)
		}

		matches := lastFMOGImagePattern.FindSubmatch(body)
		if len(matches) < 2 {
			return "", nil
		}

		imageURL := string(matches[1])
		if isPlaceholderImageURL(imageURL) {
			return "", nil
		}

		ok, err := imageURLExists(imageURL)
		if err != nil {
			return "", err
		}
		if ok {
			return imageURL, nil
		}
	}

	return "", nil
}

func EnrichArtistImageFromPage(artist *models.Artist) error {
	if artist == nil || artist.Url == "" {
		return nil
	}

	imageURL, err := GetPageImageURL(artist.Url)
	if err != nil || imageURL == "" {
		return err
	}

	artist.Image = []struct {
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

func isPlaceholderImageURL(imageURL string) bool {
	return strings.Contains(imageURL, lastFMPlaceholderImageID)
}

func imageURLExists(imageURL string) (bool, error) {
	req, err := http.NewRequest(http.MethodGet, imageURL, nil)
	if err != nil {
		return false, err
	}
	setBrowserHeaders(req, "image/avif,image/webp,image/apng,image/*,*/*;q=0.8")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	return resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices, nil
}
