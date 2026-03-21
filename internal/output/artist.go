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

const artistImageWidth = 14
const lastFMPlaceholderImageID = "2a96cbd8b46e442fc41c2b86b821562f"

func RenderArtist(artist models.Artist, title string) {
	infoLines := []string{
		title,
		fmt.Sprintf("\x1b[38;2;100;100;100m%s", artist.Name),
	}

	if artist.PlayCount != "" {
		playCount := artist.PlayCount
		if plays, err := strconv.Atoi(artist.PlayCount); err == nil {
			playCount = utils.Commas(plays)
		}

		infoLines = append(infoLines,
			"",
			"",
			fmt.Sprintf("%s plays", playCount),
		)
	}

	img := bestArtistImageURL(artist)
	if img == "" {
		verbose.Printf("rendering artist without image: %s", artist.Name)
		printLines(infoLines)
		return
	}

	imgLines, err := image.RenderANSILines(img, artistImageWidth)
	if err != nil {
		verbose.Printf("artist image render failed for %s: %v", artist.Name, err)
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
