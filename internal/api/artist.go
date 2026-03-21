package api

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/theOldZoom/gofm/internal/cache"
	"github.com/theOldZoom/gofm/internal/models"
	"github.com/theOldZoom/gofm/internal/verbose"
)

const lastFMPlaceholderImageID = "2a96cbd8b46e442fc41c2b86b821562f"

const (
	pageImageSuccessTTL   = 7 * 24 * time.Hour
	pageImageMissTTL      = 24 * time.Hour
	pageImageRateLimitTTL = 15 * time.Minute
)

var lastFMOGImagePattern = regexp.MustCompile(`<meta[^>]+property="og:image"[^>]+content="([^"]+)"`)

func GetPageImageURL(pageURL string) (string, error) {
	if imageURL, err, ok := cache.LookupPageImageURL(pageURL); ok {
		switch {
		case err != nil:
			verbose.Printf("page image cache hit (error): %s -> %v", pageURL, err)
			return "", err
		case imageURL == "":
			verbose.Printf("page image cache hit (miss): %s", pageURL)
			return "", nil
		default:
			verbose.Printf("page image cache hit: %s -> %s", pageURL, imageURL)
			return imageURL, nil
		}
	}

	body, resp, err := doRequestWithRetries("page image lookup", func() (*http.Request, error) {
		req, err := http.NewRequest(http.MethodGet, pageURL, nil)
		if err != nil {
			return nil, err
		}
		setPageHeaders(req)
		return req, nil
	})
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("page returned status code %d", resp.StatusCode)
		if resp.StatusCode == http.StatusTooManyRequests {
			ttl := pageImageRateLimitTTL
			if retryAfter, ok := retryAfterDelay(resp); ok {
				ttl = retryAfter
			}
			cache.StorePageImageRateLimit(pageURL, err, ttl)
		} else if isPageImageMissStatus(resp.StatusCode) {
			verbose.Printf("page image lookup unavailable, caching miss: %s status=%d", pageURL, resp.StatusCode)
			cache.StorePageImageMiss(pageURL, pageImageMissTTL)
			return "", nil
		}
		return "", err
	}

	matches := lastFMOGImagePattern.FindSubmatch(body)
	if len(matches) < 2 {
		verbose.Printf("page image not found in metadata: %s", pageURL)
		cache.StorePageImageMiss(pageURL, pageImageMissTTL)
		return "", nil
	}

	imageURL := string(matches[1])
	if isPlaceholderImageURL(imageURL) {
		verbose.Printf("page image was placeholder: %s", imageURL)
		cache.StorePageImageMiss(pageURL, pageImageMissTTL)
		return "", nil
	}

	ok, err := imageURLExists(imageURL)
	if err != nil {
		verbose.Printf("page image verification failed: %v", err)
		return "", err
	}
	if ok {
		verbose.Printf("page image verified: %s", imageURL)
		cache.StorePageImageURL(pageURL, imageURL, pageImageSuccessTTL)
		return imageURL, nil
	}

	verbose.Printf("page image lookup exhausted retries: %s", pageURL)
	cache.StorePageImageMiss(pageURL, pageImageMissTTL)
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

func isPageImageMissStatus(statusCode int) bool {
	return statusCode == http.StatusForbidden ||
		statusCode == http.StatusNotFound ||
		statusCode == http.StatusNotAcceptable ||
		statusCode == http.StatusGone
}

func imageURLExists(imageURL string) (bool, error) {
	_, resp, err := doRequestWithRetries("page image verify", func() (*http.Request, error) {
		req, err := http.NewRequest(http.MethodGet, imageURL, nil)
		if err != nil {
			return nil, err
		}
		setImageHeaders(req)
		return req, nil
	})
	if err != nil {
		return false, err
	}

	ok := resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices
	if !ok {
		verbose.Printf("image url check failed: %s status=%d", imageURL, resp.StatusCode)
	}

	return ok, nil
}
