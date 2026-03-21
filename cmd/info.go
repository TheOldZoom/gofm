/*
Copyright © 2026 Zoom theoldzoom@proton.me

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/theOldZoom/gofm/internal/api"
	"github.com/theOldZoom/gofm/internal/output"
	"github.com/theOldZoom/gofm/internal/verbose"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show information about an artist, track, or album",

	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("info called")
	// },
}

var infoArtistCmd = &cobra.Command{
	Use:   "artist",
	Short: "Show information about an artist",
	Run: func(cmd *cobra.Command, args []string) {
		username, _ := cmd.Flags().GetString("username")
		if username == "" {
			username = viper.GetString("username")
		}
		artistName := strings.Join(args, " ")
		if artistName == "" {
			tracks, err := api.GetRecentTracks(username, 1)
			if err != nil {
				verbose.Printf("command info artist failed: %v", err)
				fmt.Println("Failed to get recent tracks:", err)
				return
			}
			artistName = tracks[0].Artist.Name
			if artistName == "" {
				verbose.Printf("command info artist: artist not found")
				fmt.Println("Artist not found.")
				return
			}

		}
		verbose.Printf("command info artist: artist=%q args=%v", artistName, args)
		artist, err := api.GetArtistInfo(artistName)
		if err != nil {
			verbose.Printf("command info artist failed: %v", err)
			fmt.Println("Failed to get artist:", err)
			return
		}
		if artist == nil {
			verbose.Printf("command info artist: artist not found")
			fmt.Println("Artist not found.")
			return
		}
		verbose.Printf("command info artist rendering artist: %s", artist.Name)
		output.RenderArtistInfo(*artist)
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
	infoCmd.AddCommand(infoArtistCmd)
	infoArtistCmd.Flags().StringP("username", "u", "", "Username")

}
