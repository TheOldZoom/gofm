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
)

var topCmd = &cobra.Command{
	Use:   "top",
	Short: "Show a user's top artists or tracks",

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
		artists, err := api.GetUserTopArtists(username, limit)
		if err != nil {
			fmt.Println("Failed to get top artists:", err)
			return
		}
		for i, artist := range artists {
			output.RenderArtist(artist, fmt.Sprintf("%d. %s", i+1, artist.Name))
			fmt.Println()
		}
	},
}

func init() {
	rootCmd.AddCommand(topCmd)
	topCmd.AddCommand(topArtistsCmd)
	topArtistsCmd.Flags().IntP("limit", "l", 10, "Number of artists to show")
}
