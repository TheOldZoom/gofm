package cache

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/theOldZoom/gofm/internal/verbose"
)

const pageImageCacheFileName = "page-images.json"

const (
	pageImageCacheKindHit         = "hit"
	pageImageCacheKindMiss        = "miss"
	pageImageCacheKindRateLimited = "rate_limited"
)

type pageImageCacheEntry struct {
	Kind      string    `json:"kind"`
	ImageURL  string    `json:"image_url,omitempty"`
	Error     string    `json:"error,omitempty"`
	ExpiresAt time.Time `json:"expires_at"`
}

type pageImageCacheFile struct {
	Version int                            `json:"version"`
	Entries map[string]pageImageCacheEntry `json:"entries"`
}

type pageImageCacheStore struct {
	once    sync.Once
	mu      sync.Mutex
	enabled bool
	path    string
	entries map[string]pageImageCacheEntry
}

var pageImageCache pageImageCacheStore

func LookupPageImageURL(pageURL string) (string, error, bool) {
	pageImageCache.load()

	pageImageCache.mu.Lock()
	defer pageImageCache.mu.Unlock()

	if !pageImageCache.enabled {
		return "", nil, false
	}

	entry, ok := pageImageCache.entries[pageURL]
	if !ok {
		return "", nil, false
	}

	if time.Now().After(entry.ExpiresAt) {
		delete(pageImageCache.entries, pageURL)
		return "", nil, false
	}

	switch entry.Kind {
	case pageImageCacheKindHit:
		return entry.ImageURL, nil, true
	case pageImageCacheKindMiss:
		return "", nil, true
	case pageImageCacheKindRateLimited:
		return "", errors.New(entry.Error), true
	default:
		return "", nil, false
	}
}

func StorePageImageURL(pageURL string, imageURL string, ttl time.Duration) {
	pageImageCache.store(pageURL, pageImageCacheEntry{
		Kind:      pageImageCacheKindHit,
		ImageURL:  imageURL,
		ExpiresAt: time.Now().Add(ttl),
	})
}

func StorePageImageMiss(pageURL string, ttl time.Duration) {
	pageImageCache.store(pageURL, pageImageCacheEntry{
		Kind:      pageImageCacheKindMiss,
		ExpiresAt: time.Now().Add(ttl),
	})
}

func StorePageImageRateLimit(pageURL string, err error, ttl time.Duration) {
	if err == nil {
		return
	}

	pageImageCache.store(pageURL, pageImageCacheEntry{
		Kind:      pageImageCacheKindRateLimited,
		Error:     err.Error(),
		ExpiresAt: time.Now().Add(ttl),
	})
}

func (s *pageImageCacheStore) load() {
	s.once.Do(func() {
		s.entries = make(map[string]pageImageCacheEntry)

		cacheDir, err := os.UserCacheDir()
		if err != nil {
			verbose.Printf("page image cache disabled: %v", err)
			return
		}

		dir := filepath.Join(cacheDir, "gofm")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			verbose.Printf("page image cache disabled: %v", err)
			return
		}

		s.path = filepath.Join(dir, pageImageCacheFileName)
		s.enabled = true

		data, err := os.ReadFile(s.path)
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				verbose.Printf("page image cache read failed: %v", err)
			}
			return
		}

		var cacheFile pageImageCacheFile
		if err := json.Unmarshal(data, &cacheFile); err != nil {
			verbose.Printf("page image cache parse failed: %v", err)
			return
		}

		if cacheFile.Entries != nil {
			s.entries = cacheFile.Entries
		}
	})
}

func (s *pageImageCacheStore) store(pageURL string, entry pageImageCacheEntry) {
	s.load()

	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.enabled {
		return
	}

	s.entries[pageURL] = entry
	if err := s.persistLocked(); err != nil {
		verbose.Printf("page image cache write failed: %v", err)
	}
}

func (s *pageImageCacheStore) persistLocked() error {
	now := time.Now()
	for pageURL, entry := range s.entries {
		if now.After(entry.ExpiresAt) {
			delete(s.entries, pageURL)
		}
	}

	data, err := json.Marshal(pageImageCacheFile{
		Version: 1,
		Entries: s.entries,
	})
	if err != nil {
		return err
	}

	tmpPath := s.path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0o644); err != nil {
		return err
	}

	return os.Rename(tmpPath, s.path)
}
