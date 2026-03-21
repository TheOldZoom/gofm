package output

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/theOldZoom/gofm/internal/image"
	"github.com/theOldZoom/gofm/internal/models"
	"github.com/theOldZoom/gofm/internal/utils"
	"github.com/theOldZoom/gofm/internal/verbose"
)

const trackImageWidth = 14
const lastFMTrackPlaceholderImageID = "2a96cbd8b46e442fc41c2b86b821562f"

func RenderTrack(track models.Track, title string) {
	infoLines := []string{
		title,
		fmt.Sprintf("\x1b[38;2;100;100;100m%s", track.Artist.Name),
	}

	if track.PlayCount != "" {
		playCount := track.PlayCount
		if plays, err := strconv.Atoi(track.PlayCount); err == nil {
			playCount = utils.Commas(plays)
		}
		infoLines = append(infoLines,
			fmt.Sprintf("%s plays", playCount),
		)
	}

	if track.Album.Name != "" {
		if track.PlayCount == "" {
			infoLines = append(infoLines,
				"",
			)
		}
		infoLines = append(infoLines,
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
