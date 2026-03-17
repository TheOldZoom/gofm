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

func GetArtistPageImageURL(artistURL string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, artistURL, nil)
	if err != nil {
		return "", err
	}
	setBrowserHeaders(req, "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("artist page returned status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	matches := lastFMOGImagePattern.FindSubmatch(body)
	if len(matches) < 2 {
		return "", nil
	}

	imageURL := string(matches[1])
	if isPlaceholderArtistImageURL(imageURL) {
		return "", nil
	}

	return imageURL, nil
}

func EnrichArtistImageFromPage(artist *models.Artist) error {
	if artist == nil || artist.Url == "" {
		return nil
	}

	imageURL, err := GetArtistPageImageURL(artist.Url)
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

func isPlaceholderArtistImageURL(imageURL string) bool {
	return strings.Contains(imageURL, lastFMPlaceholderImageID)
}
