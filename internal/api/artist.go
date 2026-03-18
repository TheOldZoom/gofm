package api

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/theOldZoom/gofm/internal/models"
	"github.com/theOldZoom/gofm/internal/verbose"
)

const lastFMPlaceholderImageID = "2a96cbd8b46e442fc41c2b86b821562f"

var lastFMOGImagePattern = regexp.MustCompile(`<meta[^>]+property="og:image"[^>]+content="([^"]+)"`)

func GetPageImageURL(pageURL string) (string, error) {
	for attempt := 1; attempt <= 2; attempt++ {
		verbose.Printf("page image lookup attempt %d: %s", attempt, pageURL)

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
			verbose.Printf("page image read failed: %v", readErr)
			return "", readErr
		}
		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("artist page returned status code %d", resp.StatusCode)
		}

		matches := lastFMOGImagePattern.FindSubmatch(body)
		if len(matches) < 2 {
			verbose.Printf("page image not found in metadata: %s", pageURL)
			return "", nil
		}

		imageURL := string(matches[1])
		if isPlaceholderImageURL(imageURL) {
			verbose.Printf("page image was placeholder: %s", imageURL)
			return "", nil
		}

		ok, err := imageURLExists(imageURL)
		if err != nil {
			verbose.Printf("page image verification failed: %v", err)
			return "", err
		}
		if ok {
			verbose.Printf("page image verified: %s", imageURL)
			return imageURL, nil
		}

		verbose.Printf("page image missing, refetching metadata: %s", imageURL)
	}

	verbose.Printf("page image lookup exhausted retries: %s", pageURL)
	return "", nil
}

func EnrichArtistImageFromPage(artist *models.Artist) error {
	if artist == nil || artist.Url == "" {
		return nil
	}

	imageURL, err := GetPageImageURL(artist.Url)
	if err != nil || imageURL == "" {
		if err != nil {
			verbose.Printf("artist image enrichment failed for %s: %v", artist.Name, err)
		}
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

	verbose.Printf("artist image enriched from page: %s -> %s", artist.Name, imageURL)
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

	ok := resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices
	if !ok {
		verbose.Printf("image url check failed: %s status=%d", imageURL, resp.StatusCode)
	}

	return ok, nil
}
