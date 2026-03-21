package output

import (
	"fmt"
	"strings"

	"github.com/theOldZoom/gofm/internal/image"
	"github.com/theOldZoom/gofm/internal/models"
	"github.com/theOldZoom/gofm/internal/verbose"
)

const trackImageWidth = 14
const lastFMTrackPlaceholderImageID = "2a96cbd8b46e442fc41c2b86b821562f"

func RenderTrack(track models.Track, title string) {
	infoLines := []string{
		title,
		fmt.Sprintf("\x1b[38;2;100;100;100m%s", track.Artist.Name),
	}
	if track.Album.Name != "" {
		infoLines = append(infoLines,
			"",
			"",
			fmt.Sprintf("%s", track.Album.Name),
		)
	}

	img := bestTrackImageURL(track)
	if img == "" {
		verbose.Printf("rendering track without image: %s", track.Name)
		printLines(infoLines)
		return
	}

	imgLines, err := image.RenderANSILines(img, trackImageWidth)
	if err != nil {
		verbose.Printf("track image render failed for %s: %v", track.Name, err)
		printLines(infoLines)
		return
	}

	image.RenderSideBySide(imgLines, infoLines, trackImageWidth)
}

func bestTrackImageURL(track models.Track) string {
	for i := len(track.Image) - 1; i >= 0; i-- {
		if track.Image[i].Url != "" && !strings.Contains(track.Image[i].Url, lastFMTrackPlaceholderImageID) {
			return track.Image[i].Url
		}
	}

	return ""
}

func printLines(lines []string) {
	for _, line := range lines {
		fmt.Println(line)
	}
}
