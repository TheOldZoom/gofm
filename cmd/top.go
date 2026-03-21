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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/theOldZoom/gofm/internal/api"
	"github.com/theOldZoom/gofm/internal/output"
	"github.com/theOldZoom/gofm/internal/verbose"
)

var topCmd = &cobra.Command{
	Use:   "top",
	Short: "Show a user's top artists, tracks, or albums",

	// Run: func(cmd *cobra.Command, args []string) {
	// },
}

var topArtistsCmd = &cobra.Command{
	Use:   "artists",
	Short: "Show a user's top artists",
	Run: func(cmd *cobra.Command, args []string) {
		username := viper.GetString("username")
		if len(args) == 1 {
			username = args[0]
		}
		if username == "" {
			verbose.Printf("command top artists aborted: missing username")
			fmt.Println("No username provided. Pass one explicitly or run setup first.")
			return
		}
		limit, err := cmd.Flags().GetInt("limit")
		if err != nil {
			limit = 10
		}
		if limit < 1 {
			limit = 10
		}
		verbose.Printf("command top artists: username=%q limit=%d args=%v", username, limit, args)
		artists, err := api.GetUserTopArtists(username, limit)
		if err != nil {
			verbose.Printf("command top artists failed: %v", err)
			fmt.Println("Failed to get top artists:", err)
			return
		}
		verbose.Printf("command top artists rendering %d artists", len(artists))
		for i, artist := range artists {
			output.RenderArtist(artist, fmt.Sprintf("%d. %s", i+1, artist.Name))
			fmt.Println()
		}
	},
}

var topTracksCmd = &cobra.Command{
	Use:   "tracks",
	Short: "Show a user's top tracks",
	Run: func(cmd *cobra.Command, args []string) {
		username := viper.GetString("username")
		if len(args) == 1 {
			username = args[0]
		}
		if username == "" {
			verbose.Printf("command top tracks aborted: missing username")
			fmt.Println("No username provided. Pass one explicitly or run setup first.")
			return
		}
		limit, err := cmd.Flags().GetInt("limit")
		if err != nil {
			limit = 10
		}
		if limit < 1 {
			limit = 10
		}
		verbose.Printf("command top tracks: username=%q limit=%d args=%v", username, limit, args)
		tracks, err := api.GetUserTopTracks(username, limit)
		if err != nil {
			verbose.Printf("command top tracks failed: %v", err)
			fmt.Println("Failed to get top tracks:", err)
			return
		}
		verbose.Printf("command top tracks rendering %d tracks", len(tracks))
		for i, track := range tracks {
			output.RenderTrack(track, fmt.Sprintf("%d. %s", i+1, track.Name))
			fmt.Println()
		}
	},
}

var topAlbumsCmd = &cobra.Command{
	Use:   "albums",
	Short: "Show a user's top albums",
	Run: func(cmd *cobra.Command, args []string) {
		username := viper.GetString("username")
		if len(args) == 1 {
			username = args[0]
		}
		if username == "" {
			verbose.Printf("command top albums aborted: missing username")
			fmt.Println("No username provided. Pass one explicitly or run setup first.")
			return
		}
		limit, err := cmd.Flags().GetInt("limit")
		if err != nil {
			limit = 10
		}
		if limit < 1 {
			limit = 10
		}
		verbose.Printf("command top albums: username=%q limit=%d args=%v", username, limit, args)
		albums, err := api.GetUserTopAlbums(username, limit)
		if err != nil {
			verbose.Printf("command top albums failed: %v", err)
			fmt.Println("Failed to get top albums:", err)
			return
		}
		verbose.Printf("command top albums rendering %d albums", len(albums))
		for i, album := range albums {
			output.RenderAlbum(album, fmt.Sprintf("%d. %s", i+1, album.Name))
			fmt.Println()
		}
	},
}

func init() {
	rootCmd.AddCommand(topCmd)
	topCmd.AddCommand(topArtistsCmd)
	topArtistsCmd.Flags().IntP("limit", "l", 10, "Number of artists to show")
	topCmd.AddCommand(topTracksCmd)
	topTracksCmd.Flags().IntP("limit", "l", 10, "Number of tracks to show")
	topCmd.AddCommand(topAlbumsCmd)
	topAlbumsCmd.Flags().IntP("limit", "l", 10, "Number of albums to show")
}
