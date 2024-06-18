package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/nanoteck137/crustle/api"
	"github.com/pelletier/go-toml/v2"
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

type Filter struct {
	Name   string `toml:"name"`
	Filter string `toml:"filter"`
	Sort   string `toml:"sort"`
}

type FilterFile struct {
	Filters []Filter `toml:"filters"`
}

func ReadFilters(file string) ([]Filter, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var res FilterFile
	err = toml.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}

	return res.Filters, nil
}

var downloadFilterCmd = &cobra.Command{
	Use: "filter",
	Run: func(cmd *cobra.Command, args []string) {
		filterFile, err := cmd.Flags().GetString("filter-file")
		if err != nil {
			log.Fatal(err)
		}

		workDir, err := config.BootstrapDataDir()
		if err != nil {
			log.Fatal(err)
		}

		data, err := ReadDataFile(workDir)

		client := api.New(config.ServerAddress)
		client.SetToken(data.Token)

		filters, err := ReadFilters(filterFile)
		if err != nil {
			log.Fatal(err)
		}

		var options []huh.Option[Filter]

		for _, filter := range filters {
			options = append(options, huh.NewOption(filter.Name, filter))
		}

		var filter Filter
		s := huh.NewSelect[Filter]().
			Title("Select Filter").
			Options(options...).
			Value(&filter)
		err = s.Run()
		if err != nil {
			log.Fatal(err)
		}

		res, err := client.GetTracks(filter.Filter, filter.Sort)
		if err != nil {
			log.Fatal(err)
		}

		err = os.MkdirAll(config.DownloadDir, 0755)
		if err != nil {
			log.Fatal(err)
		}

		// // TODO(patrik): Sanitize the name
		name := "Filter - " + filter.Name
		dst := path.Join(config.DownloadDir, name)

		_, err = os.Stat(dst)
		if err == nil {
			val := false
			s := huh.NewConfirm().
				Title(fmt.Sprintf("'%s' is not empty\nThis operation will delete the folder", dst)).
				Value(&val)
			err = s.Run()
			if err != nil {
				log.Fatal(err)
			}

			err = os.RemoveAll(dst)
			if err != nil {
				log.Fatal(err)
			}
		} 

		err = os.MkdirAll(dst, 0755)
		if err != nil {
			log.Fatal(err)
		}

		tmpDir := path.Join(config.DownloadDir, "tmp")

		err = os.MkdirAll(tmpDir, 0755)
		if err != nil {
			log.Fatal(err)
		}

		for i, track := range res.Tracks {
			fmt.Printf("Downloading: %s\n", track.Name)
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
			args = append(args, "--album", name)
			args = append(args, "--number", strconv.FormatInt(int64(i + 1), 10))
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
	},
}

var downloadUpdateCmd = &cobra.Command{
	Use: "update",
	Run: func(cmd *cobra.Command, args []string) {},
}

func init() {
	downloadFilterCmd.Flags().StringP("filter-file", "f", "", "Filter file to use")
	downloadFilterCmd.MarkFlagRequired("filter-file")

	downloadCmd.AddCommand(downloadPlaylistCmd)
	downloadCmd.AddCommand(downloadFilterCmd)
	downloadCmd.AddCommand(downloadUpdateCmd)

	rootCmd.AddCommand(downloadCmd)
}
