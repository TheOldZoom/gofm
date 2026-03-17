package output

import (
	"fmt"
	"strings"

	"github.com/theOldZoom/gofm/internal/image"
	"github.com/theOldZoom/gofm/internal/models"
)

const artistImageWidth = 14
const lastFMPlaceholderImageID = "2a96cbd8b46e442fc41c2b86b821562f"

func RenderArtist(artist models.Artist, title string) {
	infoLines := []string{
		title,
		fmt.Sprintf("\x1b[38;2;100;100;100m%s", artist.Name),
	}
	if artist.PlayCount != "" {
		infoLines = append(infoLines,
			"",
			"",
			fmt.Sprintf("Play Count: %s", artist.PlayCount),
		)
	}

	img := bestArtistImageURL(artist)
	if img == "" {
		printLines(infoLines)
		return
	}

	imgLines, err := image.RenderANSILines(img, artistImageWidth)
	if err != nil {
		printLines(infoLines)
		return
	}

	image.RenderSideBySide(imgLines, infoLines, artistImageWidth)
}

func bestArtistImageURL(artist models.Artist) string {
	for i := len(artist.Image) - 1; i >= 0; i-- {
		if artist.Image[i].Url != "" && !strings.Contains(artist.Image[i].Url, lastFMPlaceholderImageID) {
			return artist.Image[i].Url
		}
	}

	return ""
}
