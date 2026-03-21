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

const trackImageWidth = 14
const lastFMTrackPlaceholderImageID = "2a96cbd8b46e442fc41c2b86b821562f"

var htmlTrackTagPattern = regexp.MustCompile(`<[^>]+>`)

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
	} else if track.UserPlayCount != "" {
		playCount := track.UserPlayCount
		if plays, err := strconv.Atoi(track.UserPlayCount); err == nil {
			playCount = utils.Commas(plays)
		}
		infoLines = append(infoLines,
			fmt.Sprintf("%s plays", playCount),
		)
	}

	if track.Album.Name != "" {
		if track.PlayCount == "" && track.UserPlayCount == "" {
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

func RenderTrackInfo(track models.Track) {
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

	infoLines := []string{track.Name}

	if track.Artist.Name != "" {
		label := track.Artist.Name
		if track.Artist.URL != "" {
			label = link(track.Artist.URL, track.Artist.Name)
		}
		infoLines = append(infoLines, fmt.Sprintf("\x1b[38;2;100;100;100m%s", label))
	}

	if track.Url != "" {
		infoLines = append(infoLines, fmt.Sprintf("\x1b[38;2;100;100;100m%s", link(track.Url, track.Url)))
	}

	stats := make([]string, 0, 4)
	if plays := formatCount(track.PlayCount, "total plays"); plays != "" {
		stats = append(stats, plays)
	}
	if listeners := formatCount(track.Listeners, "listeners"); listeners != "" {
		stats = append(stats, listeners)
	}
	if userPlays := formatCount(track.UserPlayCount, "plays"); userPlays != "" {
		stats = append(stats, fmt.Sprintf("%s", userPlays))
	}
	if duration := formatTrackDuration(string(track.Duration)); duration != "" {
		stats = append(stats, fmt.Sprintf("%s", duration))
	}
	if track.UserLoved == "1" {
		stats = append(stats, "Loved track")
	}
	if len(stats) > 0 {
		infoLines = append(infoLines, "")
		infoLines = append(infoLines, stats...)
	}

	if track.Album.Name != "" {
		albumLabel := track.Album.Name
		if track.Album.URL != "" {
			albumLabel = link(track.Album.URL, track.Album.Name)
		}
		infoLines = append(infoLines, "", fmt.Sprintf("%s", albumLabel))
	}

	if len(track.TopTags.Tag) > 0 {
		tags := make([]string, 0, len(track.TopTags.Tag))
		for _, tag := range track.TopTags.Tag {
			if tag.Name != "" {
				tags = append(tags, link(tag.Url, tag.Name))
			}
		}
		if len(tags) > 0 {
			infoLines = append(infoLines, "", fmt.Sprintf("Tags: %s", strings.Join(tags, ", ")))
		}
	}

	summary := cleanTrackWiki(track.Wiki.Summary)
	if summary != "" {
		infoLines = append(infoLines, "")
		infoLines = append(infoLines, wrapTrackText(summary, 72)...)
	}

	img := bestTrackImageURL(track)
	if img == "" {
		verbose.Printf("rendering track without image: %s", track.Name)
		printLines(infoLines)
		return
	}

	imgLines, err := image.RenderANSILines(img, 40)
	if err != nil {
		verbose.Printf("track image render failed for %s: %v", track.Name, err)
		printLines(infoLines)
		return
	}

	image.RenderSideBySide(imgLines, infoLines, 40)
}

func bestTrackImageURL(track models.Track) string {
	for i := len(track.Image) - 1; i >= 0; i-- {
		if track.Image[i].Url != "" && !strings.Contains(track.Image[i].Url, lastFMTrackPlaceholderImageID) {
			return track.Image[i].Url
		}
	}

	return ""
}

func cleanTrackWiki(text string) string {
	if text == "" {
		return ""
	}

	text = htmlTrackTagPattern.ReplaceAllString(text, "")
	text = html.UnescapeString(text)
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.Join(strings.Fields(text), " ")
	text = strings.TrimSpace(strings.TrimSuffix(text, "Read more on Last.fm."))
	return strings.TrimSpace(strings.TrimSuffix(text, "User-contributed text is available under the Creative Commons By-SA License; additional terms may apply."))
}

func wrapTrackText(text string, width int) []string {
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

func formatTrackDuration(duration string) string {
	if duration == "" {
		return ""
	}

	ms, err := strconv.Atoi(duration)
	if err != nil || ms <= 0 {
		return ""
	}

	totalSeconds := ms / 1000
	minutes := totalSeconds / 60
	seconds := totalSeconds % 60
	if minutes == 0 {
		return fmt.Sprintf("%d seconds", seconds)
	}
	return fmt.Sprintf("%d minutes %d seconds", minutes, seconds)
}

func printLines(lines []string) {
	for _, line := range lines {
		fmt.Println(line)
	}
}
