package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/nanoteck137/crustle/api"
	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use: "export",
	Run: func(cmd *cobra.Command, args []string) {
		output, err := cmd.Flags().GetString("output")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("output: %v\n", output)

		workDir, err := config.BootstrapDataDir()
		if err != nil {
			log.Fatal(err)
		}

		data, err := ReadDataFile(workDir)

		client := api.New("http://127.0.0.1:3000")
		client.SetToken(data.Token)

		res, err := client.GetPlaylists(api.Options{})
		if err != nil {
			log.Fatal(err)
		}

		type PlaylistItem struct {
			TrackName  string `json:"trackName"`
			AlbumName  string `json:"albumName"`
			ArtistName string `json:"artstiName"`
		}

		type Playlist struct {
			Name  string         `json:"name"`
			Items []PlaylistItem `json:"items"`
		}

		exportedPlaylists := make([]Playlist, 0, len(res.Playlists))
		for _, playlist := range res.Playlists {
			p, err := client.GetPlaylistById(playlist.Id, api.Options{})
			if err != nil {
				log.Fatal(err)
			}

			items := make([]PlaylistItem, 0, len(p.Items))
			for _, track := range p.Items {
				items = append(items, PlaylistItem{
					TrackName:  track.Name,
					AlbumName:  track.AlbumName,
					ArtistName: track.ArtistName,
				})
			}

			exportedPlaylists = append(exportedPlaylists, Playlist{
				Name:  playlist.Name,
				Items: items,
			})
		}

		b, err := json.MarshalIndent(exportedPlaylists, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		err = os.WriteFile(output, b, 0644)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	exportCmd.Flags().StringP("output", "o", "", "Output File")
	exportCmd.MarkFlagRequired("output")

	rootCmd.AddCommand(exportCmd)
}
