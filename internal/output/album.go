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

const albumImageWidth = 14
const lastFMAlbumPlaceholderImageID = "2a96cbd8b46e442fc41c2b86b821562f"

func RenderAlbum(album models.Album, title string) {
	infoLines := []string{
		title,
	}

	if album.Artist.Name != "" {
		infoLines = append(infoLines, fmt.Sprintf("\x1b[38;2;100;100;100m%s", album.Artist.Name))
	}

	if album.PlayCount != "" {
		playCount := album.PlayCount
		if plays, err := strconv.Atoi(album.PlayCount); err == nil {
			playCount = utils.Commas(plays)
		}
		infoLines = append(infoLines,
			fmt.Sprintf("%s plays", playCount),
		)
	}

	img := bestAlbumImageURL(album)
	if img == "" {
		verbose.Printf("rendering album without image: %s", album.Name)
		fmt.Println(strings.Join(infoLines, "\n"))
		return
	}

	imgLines, err := image.RenderANSILines(img, albumImageWidth)
	if err != nil {
		verbose.Printf("album image render failed for %s: %v", album.Name, err)
		fmt.Println(strings.Join(infoLines, "\n"))
		return
	}

	image.RenderSideBySide(imgLines, infoLines, albumImageWidth)
}

func bestAlbumImageURL(album models.Album) string {
	for i := len(album.Image) - 1; i >= 0; i-- {
		if album.Image[i].Url != "" && !strings.Contains(album.Image[i].Url, lastFMAlbumPlaceholderImageID) {
			return album.Image[i].Url
		}
	}

	return ""
}
