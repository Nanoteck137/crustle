package cmd

import (
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/kr/pretty"
	"github.com/nanoteck137/crustle/api"
	"github.com/spf13/cobra"
)

func DownloadFile(url, dst string) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = os.WriteFile(dst, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func DownloadTrack(track *api.Track, dst string) (string, error) {
	split := strings.Split(track.MobileQualityFile, ".")
	ext := split[len(split)-1]

	p := path.Join(dst, track.Id+"."+ext)

	err := DownloadFile(track.MobileQualityFile, p)
	if err != nil {
		return "", err
	}

	return p, nil
}

func DownloadTrackCover(track *api.Track, dst string) (string, error) {
	split := strings.Split(track.CoverArt, ".")
	ext := split[len(split)-1]

	p := path.Join(dst, "cover."+ext)

	err := DownloadFile(track.CoverArt, p)
	if err != nil {
		return "", err
	}

	return p, nil
}

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

		var options []huh.Option[api.Playlist]

		for _, playlist := range res.Playlists {
			options = append(options, huh.NewOption(playlist.Name, playlist))
		}

		var playlist api.Playlist
		s := huh.NewSelect[api.Playlist]().
			Title("Testing").
			Options(options...).
			Value(&playlist)
		err = s.Run()
		if err != nil {
			log.Fatal(err)
		}

		p, err := client.GetPlaylistById(playlist.Id)
		if err != nil {
			log.Fatal(err)
		}

		pretty.Println(p)

		err = os.MkdirAll(config.DownloadDir, 0755)
		if err != nil {
			log.Fatal(err)
		}

		// TODO(patrik): Sanitize the name
		name := "Playlist - " + playlist.Name
		dst := path.Join(config.DownloadDir, name)

		err = os.MkdirAll(dst, 0755)
		if err != nil {
			log.Fatal(err)
		}

		tmpDir := path.Join(config.DownloadDir, "tmp")

		err = os.MkdirAll(tmpDir, 0755)
		if err != nil {
			log.Fatal(err)
		}

		for _, track := range p.Items {
			p, err := DownloadTrack(&track, dst)
			if err != nil {
				log.Fatal(err)
			}

			c, err := DownloadTrackCover(&track, tmpDir)
			if err != nil {
				log.Fatal(err)
			}

			var args []string

			args = append(args, "--title", track.Name)
			args = append(args, "--artist", track.ArtistName)
			args = append(args, "--album", playlist.Name)
			args = append(args, "--number", strconv.Itoa(track.Number))
			args = append(args, "--image", c)
			args = append(args, "--remove")

			args = append(args, p)

			cmd := exec.Command("tagopus", args...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err = cmd.Run()
			if err != nil {
				log.Fatal(err)
			}
		}

		err = os.RemoveAll(tmpDir)
		if err != nil {
			log.Fatal(err)
		}

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
