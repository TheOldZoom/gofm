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
	Use:   "artist \"[artist name]\"",
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

var infoTrackCmd = &cobra.Command{
	Use:   "track \"[track name]\" \"[artist name]\"",
	Short: "Show information about a track",
	Run: func(cmd *cobra.Command, args []string) {
		username, _ := cmd.Flags().GetString("username")
		if username == "" {
			username = viper.GetString("username")
		}
		trackName := ""
		artistName := ""
		switch {
		case len(args) == 0:
		case artistName != "":
			trackName = strings.Join(args, " ")
		case len(args) == 1:
			trackName = args[0]
		default:
			trackName = strings.Join(args[:len(args)-1], " ")
			artistName = args[len(args)-1]
		}
		if trackName == "" {
			if username == "" {
				verbose.Printf("command info track failed: missing track and username")
				fmt.Println("Track name is required unless a username is configured.")
				return
			}

			tracks, err := api.GetRecentTracks(username, 1)
			if err != nil {
				verbose.Printf("command info track failed: %v", err)
				fmt.Println("Failed to get recent tracks:", err)
				return
			}
			if len(tracks) == 0 {
				verbose.Printf("command info track: no recent tracks for %s", username)
				fmt.Println("No recent tracks found.")
				return
			}

			trackName = tracks[0].Name
			artistName = tracks[0].Artist.Name
		}

		if artistName == "" {
			verbose.Printf("command info track failed: missing artist for %s", trackName)
			fmt.Println("Artist name is required. Pass it as the second argument or with --artist.")
			return
		}

		track, err := api.GetTrackInfo(artistName, trackName, username)
		if err != nil {
			verbose.Printf("command info track failed: %v", err)
			fmt.Println("Failed to get track:", err)
			return
		}
		if track == nil {
			if username == "" {
				verbose.Printf("command info track: track not found and no username fallback")
				fmt.Println("Track not found.")
				return
			}

			verbose.Printf("command info track: track not found, falling back to recent track for %s", username)
			tracks, err := api.GetRecentTracks(username, 1)
			if err != nil {
				verbose.Printf("command info track fallback failed: %v", err)
				fmt.Println("Track not found, and failed to get recent tracks:", err)
				return
			}
			if len(tracks) == 0 {
				verbose.Printf("command info track fallback: no recent tracks for %s", username)
				fmt.Println("Track not found, and no recent tracks found.")
				return
			}

			track = &tracks[0]
		}
		verbose.Printf("command info track rendering track: %s", track.Name)
		output.RenderTrackInfo(*track)
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
	infoCmd.AddCommand(infoArtistCmd)
	infoArtistCmd.Flags().StringP("username", "u", "", "Username")
	infoCmd.AddCommand(infoTrackCmd)
	infoTrackCmd.Flags().StringP("username", "u", "", "Username")
}
