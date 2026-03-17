package image

import (
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
)

func LoadImage(source string) (image.Image, error) {
	reader, err := openImageSource(source)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func openImageSource(source string) (io.ReadCloser, error) {
	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		resp, err := http.Get(source)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
			resp.Body.Close()
			return nil, fmt.Errorf("failed to fetch image: %s", resp.Status)
		}
		return resp.Body, nil
	}

	return os.Open(source)
}

func RenderANSI(source string, width int) {
	lines, err := RenderANSILines(source, width)
	if err != nil {
		log.Fatal(err)
	}
	for _, line := range lines {
		fmt.Println(line)
	}
}

func RenderANSILines(source string, width int) ([]string, error) {
	img, err := LoadImage(source)
	if err != nil {
		return nil, err
	}
	bounds := img.Bounds()
	srcW := bounds.Dx()
	srcH := bounds.Dy()
	if srcW == 0 || srcH == 0 {
		return nil, fmt.Errorf("image has no size")
	}

	height := int(math.Max(1, math.Round(float64(srcH)*float64(width)/float64(srcW)/2)))
	lines := make([]string, 0, height)

	for y := 0; y < height; y++ {
		var line strings.Builder
		for x := 0; x < width; x++ {
			top := Sample(img, x, y*2, width, height*2)
			bottom := Sample(img, x, y*2+1, width, height*2)
			tr, tg, tb := Rgb(top)
			br, bg, bb := Rgb(bottom)
			fmt.Fprintf(&line, "\x1b[38;2;%d;%d;%dm\x1b[48;2;%d;%d;%dm▀", tr, tg, tb, br, bg, bb)
		}
		line.WriteString("\x1b[0m")
		lines = append(lines, line.String())
	}

	return lines, nil
}

func RenderASCII(img image.Image, width int) {
	const ramp = " .:-=+*#%@"

	bounds := img.Bounds()
	srcW := bounds.Dx()
	srcH := bounds.Dy()
	if srcW == 0 || srcH == 0 {
		log.Fatal("image has no size")
	}

	height := int(math.Max(1, math.Round(float64(srcH)*float64(width)/float64(srcW)*0.5)))

	for y := 0; y < height; y++ {
		var line strings.Builder
		for x := 0; x < width; x++ {
			c := Sample(img, x, y, width, height)
			r, g, b := Rgb(c)
			luma := 0.2126*float64(r) + 0.7152*float64(g) + 0.0722*float64(b)
			idx := int(math.Round((luma / 255) * float64(len(ramp)-1)))
			line.WriteByte(ramp[idx])
		}
		line.WriteByte('\n')
		fmt.Print(line.String())
	}
}

func Sample(img image.Image, x, y, dstW, dstH int) color.Color {
	b := img.Bounds()
	srcX := b.Min.X + x*b.Dx()/dstW
	srcY := b.Min.Y + y*b.Dy()/dstH
	return img.At(srcX, srcY)
}

func Rgb(c color.Color) (uint8, uint8, uint8) {
	r, g, b, _ := c.RGBA()
	return uint8(r >> 8), uint8(g >> 8), uint8(b >> 8)
}

func RenderSideBySide(leftLines, rightLines []string, leftWidth int) {
	rows := len(leftLines)
	if len(rightLines) > rows {
		rows = len(rightLines)
	}

	leftStart := 0
	rightStart := 1

	for i := 0; i < rows; i++ {
		left := strings.Repeat(" ", leftWidth)
		if i >= leftStart && i < leftStart+len(leftLines) {
			left = leftLines[i-leftStart]
		}

		right := ""
		if i >= rightStart && i < rightStart+len(rightLines) {
			right = rightLines[i-rightStart]
		}

		fmt.Printf("%s  %s\n", left, right)
	}
}
