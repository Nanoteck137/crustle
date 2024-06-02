package cmd

import (
	"log"

	"github.com/kr/pretty"
	"github.com/nanoteck137/crustle/api"
	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Use: "download",
}

var downloadPlaylistCmd = &cobra.Command{
	Use: "playlist",
	Run: func(cmd *cobra.Command, args []string) {
		workDir, err := config.BootstrapDataDir()
		if err != nil {
			log.Fatal(err)
		}

		data, err := ReadDataFile(workDir)

		client := api.New("http://127.0.0.1:3000")
		client.SetToken(data.Token)

		res, err := client.GetPlaylists()
		if err != nil {
			log.Fatal(err)
		}

		pretty.Println(res.Playlists)

		playlist := res.Playlists[0]

		p, err := client.GetPlaylistById(playlist.Id)
		if err != nil {
			log.Fatal(err)
		}
		
		pretty.Println(p)


		// playlist.
	},
}

var downloadUpdateCmd = &cobra.Command{
	Use: "update",
	Run: func(cmd *cobra.Command, args []string) {},
}

func init() {
	downloadCmd.AddCommand(downloadPlaylistCmd)
	downloadCmd.AddCommand(downloadUpdateCmd)

	rootCmd.AddCommand(downloadCmd)
}
