package output

import (
	"fmt"
	"html"
	"regexp"
	"strconv"
	"strings"

	"github.com/theOldZoom/gofm/internal/image"
	"github.com/theOldZoom/gofm/internal/models"
	"github.com/theOldZoom/gofm/internal/utils"
	"github.com/theOldZoom/gofm/internal/verbose"
)

const albumImageWidth = 14
const lastFMAlbumPlaceholderImageID = "2a96cbd8b46e442fc41c2b86b821562f"

var htmlAlbumTagPattern = regexp.MustCompile(`<[^>]+>`)

func RenderAlbum(album models.Album, title string) {
	infoLines := []string{
		title,
	}

	if album.Artist.Name != "" {
		infoLines = append(infoLines, fmt.Sprintf("\x1b[38;2;100;100;100m%s", album.Artist.Name))
	}

	if album.PlayCount != "" {
		playCount := string(album.PlayCount)
		if plays, err := strconv.Atoi(playCount); err == nil {
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

func RenderAlbumInfo(album models.Album) {
	link := func(url string, label string) string {
		if url == "" || label == "" {
			return label
		}

		return fmt.Sprintf("\x1b]8;;%s\x1b\\%s\x1b]8;;\x1b\\", url, label)
	}

	formatCount := func(value string, suffix string) string {
		if value == "" {
			return ""
		}
		if n, err := strconv.Atoi(value); err == nil {
			value = utils.Commas(n)
		}
		if suffix == "" {
			return value
		}
		return fmt.Sprintf("%s %s", value, suffix)
	}

	infoLines := []string{album.Name}

	if album.Artist.Name != "" {
		label := album.Artist.Name
		if album.Artist.URL != "" {
			label = link(album.Artist.URL, album.Artist.Name)
		}
		infoLines = append(infoLines, fmt.Sprintf("\x1b[38;2;100;100;100m%s", label))
	}

	if album.URL != "" {
		infoLines = append(infoLines, fmt.Sprintf("\x1b[38;2;100;100;100m%s", link(album.URL, album.URL)))
	}

	stats := make([]string, 0, 4)

	if listeners := formatCount(string(album.Listeners), "listeners"); listeners != "" {
		stats = append(stats, listeners)
	}
	if plays := formatCount(string(album.PlayCount), "total plays"); plays != "" {
		stats = append(stats, plays)
	}
	if userPlays := formatCount(string(album.UserPlayCount), "plays"); userPlays != "" {
		stats = append(stats, fmt.Sprintf("%s", userPlays))
	}
	if released := strings.TrimSpace(album.Releasedate); released != "" {
		stats = append(stats, fmt.Sprintf("%s", released))
	}
	if len(album.Tracks.Track) > 0 {
		stats = append(stats, fmt.Sprintf("%d tracks", len(album.Tracks.Track)))
	}
	if len(stats) > 0 {
		infoLines = append(infoLines, "")
		infoLines = append(infoLines, stats...)
	}

	if len(album.Tags.Tag) > 0 {
		tags := make([]string, 0, len(album.Tags.Tag))
		for _, tag := range album.Tags.Tag {
			if tag.Name != "" {
				tags = append(tags, link(tag.Url, tag.Name))
			}
		}
		if len(tags) > 0 {
			infoLines = append(infoLines, "", fmt.Sprintf("Tags: %s", strings.Join(tags, ", ")))
		}
	}

	if len(album.Tracks.Track) > 0 {
		infoLines = append(infoLines, "", "Tracks:")
		for i, track := range album.Tracks.Track {
			line := fmt.Sprintf("%d. %s", i+1, track.Name)
			if duration := formatAlbumTrackDuration(string(track.Duration)); duration != "" {
				line = fmt.Sprintf("%s (%s)", line, duration)
			}
			infoLines = append(infoLines, line)
		}
	}

	summary := cleanAlbumWiki(album.Wiki.Summary)
	if summary != "" {
		infoLines = append(infoLines, "")
		infoLines = append(infoLines, wrapAlbumText(summary, 72)...)
	}

	img := bestAlbumImageURL(album)
	if img == "" {
		verbose.Printf("rendering album without image: %s", album.Name)
		fmt.Println(strings.Join(infoLines, "\n"))
		return
	}

	imgLines, err := image.RenderANSILines(img, 40)
	if err != nil {
		verbose.Printf("album image render failed for %s: %v", album.Name, err)
		fmt.Println(strings.Join(infoLines, "\n"))
		return
	}

	image.RenderSideBySide(imgLines, infoLines, 40)
}

func bestAlbumImageURL(album models.Album) string {
	for i := len(album.Image) - 1; i >= 0; i-- {
		if album.Image[i].Url != "" && !strings.Contains(album.Image[i].Url, lastFMAlbumPlaceholderImageID) {
			return album.Image[i].Url
		}
	}

	return ""
}

func cleanAlbumWiki(text string) string {
	if text == "" {
		return ""
	}

	text = htmlAlbumTagPattern.ReplaceAllString(text, "")
	text = html.UnescapeString(text)
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.Join(strings.Fields(text), " ")
	text = strings.TrimSpace(strings.TrimSuffix(text, "Read more on Last.fm."))
	return strings.TrimSpace(strings.TrimSuffix(text, "User-contributed text is available under the Creative Commons By-SA License; additional terms may apply."))
}

func wrapAlbumText(text string, width int) []string {
	if text == "" {
		return nil
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return nil
	}

	lines := make([]string, 0, len(words)/8+1)
	current := words[0]
	for _, word := range words[1:] {
		if len(current)+1+len(word) <= width {
			current += " " + word
			continue
		}
		lines = append(lines, current)
		current = word
	}

	lines = append(lines, current)
	return lines
}

func formatAlbumTrackDuration(duration string) string {
	if duration == "" {
		return ""
	}

	seconds, err := strconv.Atoi(duration)
	if err != nil || seconds <= 0 {
		return ""
	}

	minutes := seconds / 60
	remainingSeconds := seconds % 60
	if minutes == 0 {
		return fmt.Sprintf("%ds", remainingSeconds)
	}

	return fmt.Sprintf("%dm%ds", minutes, remainingSeconds)
}
