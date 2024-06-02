package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/kr/pretty"
	"github.com/nanoteck137/crustle/api"
	"github.com/nanoteck137/crustle/types"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

type Creds struct {
	Username string
	Password string
}

func ReadCreds() (Creds, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Username: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		return Creds{}, err
	}

	fmt.Print("Enter Password: ")
	bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return Creds{}, err
	}

	fmt.Println()

	password := string(bytePassword)

	return Creds{
		Username: strings.TrimSpace(username),
		Password: strings.TrimSpace(password),
	}, nil
}

func ReadDataFile(workDir types.WorkDir) (types.DataFile, error) {
	d, err := os.ReadFile(workDir.DataFile())
	if err != nil {
		if !os.IsNotExist(err) {
			return types.DataFile{}, err
		}

	}

	var data types.DataFile
	err = json.Unmarshal(d, &data)
	if err != nil {
		return types.DataFile{}, err
	}

	return data, nil

}

func WriteDataFile(workDir types.WorkDir, data types.DataFile) error {
	d, err := json.Marshal(&data)
	if err != nil {
		return err
	}

	err = os.WriteFile(workDir.DataFile(), d, 0644)
	if err != nil {
		return err
	}

	return nil
}

var loginCmd = &cobra.Command{
	Use: "login",
	Run: func(cmd *cobra.Command, args []string) {
		workDir, err := config.BootstrapDataDir()
		if err != nil {
			log.Fatal(err)
		}

		data, err := ReadDataFile(workDir)

		// TODO(patrik): Check if token is valid
		if data.Token != "" {
			log.Println("Warning: Already logged in")
		}

		creds, err := ReadCreds()
		if err != nil {
			log.Fatal(err)
		}

		client := api.New("http://127.0.0.1:3000")

		res, err := client.Login(api.PostAuthSigninBody{
			Username: creds.Username,
			Password: creds.Password,
		})
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("res.Token: %v\n", res.Token)

		data.Token = res.Token

		err = WriteDataFile(workDir, data)
		if err != nil {
			log.Fatal(err)
		}

		pretty.Println(res)
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
