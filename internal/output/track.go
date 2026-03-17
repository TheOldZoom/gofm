package output

import (
	"fmt"

	"github.com/theOldZoom/gofm/internal/image"
	"github.com/theOldZoom/gofm/internal/models"
)

const trackImageWidth = 14

func RenderTrack(track models.Track, title string) {
	infoLines := []string{
		title,
		fmt.Sprintf("\x1b[38;2;100;100;100m%s", track.Artist.Name),
	}
	if track.Album.Name != "" {
		infoLines = append(infoLines,
			"",
			"",
			fmt.Sprintf("Album: %s", track.Album.Name),
		)
	}

	img := bestTrackImageURL(track)
	if img == "" {
		printLines(infoLines)
		return
	}

	imgLines, err := image.RenderANSILines(img, trackImageWidth)
	if err != nil {
		printLines(infoLines)
		return
	}

	image.RenderSideBySide(imgLines, infoLines, trackImageWidth)
}

func bestTrackImageURL(track models.Track) string {
	for i := len(track.Image) - 1; i >= 0; i-- {
		if track.Image[i].Url != "" {
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
