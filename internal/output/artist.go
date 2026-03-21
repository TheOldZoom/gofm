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

const artistImageWidth = 14
const lastFMPlaceholderImageID = "2a96cbd8b46e442fc41c2b86b821562f"

var htmlTagPattern = regexp.MustCompile(`<[^>]+>`)

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

func RenderArtistInfo(artist models.Artist) {
	link := func(url string, label string) string {
		if url == "" || label == "" {
			return label
		}

		return fmt.Sprintf("\x1b]8;;%s\x1b\\%s\x1b]8;;\x1b\\", url, label)
	}

	infoLines := []string{
		artist.Name,
	}

	if artist.Url != "" {
		infoLines = append(infoLines, fmt.Sprintf("\x1b[38;2;100;100;100m%s", link(artist.Url, artist.Url)))
	}

	if artist.OnTour == "1" {
		infoLines = append(infoLines, "", "On tour")
	}

	stats := make([]string, 0, 3)
	if artist.Stats.Listeners != "" {
		listeners := artist.Stats.Listeners
		if n, err := strconv.Atoi(listeners); err == nil {
			listeners = utils.Commas(n)
		}
		stats = append(stats, fmt.Sprintf("%s listeners", listeners))
	}

	playCount := artist.Stats.PlayCount
	if playCount == "" {
		playCount = artist.PlayCount
	}
	if playCount != "" {
		if n, err := strconv.Atoi(playCount); err == nil {
			playCount = utils.Commas(n)
		}
		stats = append(stats, fmt.Sprintf("%s total plays", playCount))
	}

	if artist.Stats.UserPlayCount != "" {
		userPlayCount := artist.Stats.UserPlayCount
		if n, err := strconv.Atoi(userPlayCount); err == nil {
			userPlayCount = utils.Commas(n)
		}
		stats = append(stats, fmt.Sprintf("%s plays", userPlayCount))
	}

	if len(stats) > 0 {
		infoLines = append(infoLines, "")
		infoLines = append(infoLines, stats...)
	}

	if len(artist.Tags.Tag) > 0 {
		tags := make([]string, 0, len(artist.Tags.Tag))
		for _, tag := range artist.Tags.Tag {
			if tag.Name != "" {
				tags = append(tags, link(tag.Url, tag.Name))
			}
		}
		if len(tags) > 0 {
			infoLines = append(infoLines, "", fmt.Sprintf("Tags: %s", strings.Join(tags, ", ")))
		}
	}

	summary := cleanArtistBio(artist.Bio.Summary)
	if summary != "" {
		infoLines = append(infoLines, "")
		infoLines = append(infoLines, wrapText(summary, 72)...)
	}

	img := bestArtistImageURL(artist)
	if img == "" {
		verbose.Printf("rendering artist without image: %s", artist.Name)
		printLines(infoLines)
		return
	}

	imgLines, err := image.RenderANSILines(img, 40)
	if err != nil {
		verbose.Printf("artist image render failed for %s: %v", artist.Name, err)
		printLines(infoLines)
		return
	}

	image.RenderSideBySide(imgLines, infoLines, 40)
}

func bestArtistImageURL(artist models.Artist) string {
	for i := len(artist.Image) - 1; i >= 0; i-- {
		if artist.Image[i].Url != "" && !strings.Contains(artist.Image[i].Url, lastFMPlaceholderImageID) {
			return artist.Image[i].Url
		}
	}

	return ""
}

func cleanArtistBio(text string) string {
	if text == "" {
		return ""
	}

	text = htmlTagPattern.ReplaceAllString(text, "")
	text = html.UnescapeString(text)
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.Join(strings.Fields(text), " ")
	text = strings.TrimSpace(strings.TrimSuffix(text, "Read more on Last.fm"))
	return strings.TrimSpace(strings.TrimSuffix(text, ". User-contributed text is available under the Creative Commons By-SA License; additional terms may apply."))
}

func wrapText(text string, width int) []string {
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
