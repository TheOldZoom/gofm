package api

import (
	"sync"

	"github.com/theOldZoom/gofm/internal/models"
	"github.com/theOldZoom/gofm/internal/verbose"
)

func enrichTracksConcurrently(source string, tracks []models.Track, includeAlbum bool) {
	var wg sync.WaitGroup

	for i := range tracks {
		wg.Add(1)

		go func(track *models.Track) {
			defer wg.Done()

			if includeAlbum {
				if err := EnrichTrackAlbumFromAPI(track); err != nil {
					verbose.Printf("%s track album enrichment failed for %s - %s: %v", source, track.Artist.Name, track.Name, err)
				}
			}

			if err := EnrichTrackImageFromPage(track); err != nil {
				verbose.Printf("%s track image fallback failed for %s: %v", source, track.Name, err)
			}
		}(&tracks[i])
	}

	wg.Wait()
}

func enrichTrackPlayCountsConcurrently(source string, username string, tracks []models.Track) {
	var wg sync.WaitGroup

	for i := range tracks {
		wg.Add(1)

		go func(track *models.Track) {
			defer wg.Done()

			if err := EnrichTrackUserPlayCountFromAPI(track, username); err != nil {
				verbose.Printf("%s track playcount enrichment failed for %s - %s: %v", source, track.Artist.Name, track.Name, err)
			}
		}(&tracks[i])
	}

	wg.Wait()
}

func enrichArtistsConcurrently(source string, artists []models.Artist) {
	var wg sync.WaitGroup

	for i := range artists {
		wg.Add(1)

		go func(artist *models.Artist) {
			defer wg.Done()

			if err := EnrichArtistImageFromPage(artist); err != nil {
				verbose.Printf("%s artist image fallback failed for %s: %v", source, artist.Name, err)
			}
		}(&artists[i])
	}

	wg.Wait()
}
