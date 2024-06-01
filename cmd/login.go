package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/kr/pretty"
	"github.com/nanoteck137/crustle/api"
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

type ApiError[E any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Errors  E      `json:"errors,omitempty"`
}

func (err *ApiError[E]) Error() string {
	return err.Message
}

type ApiResponse[D any, E any] struct {
	Status string       `json:"status"`
	Data   D            `json:"data,omitempty"`
	Error  *ApiError[E] `json:"error,omitempty"`
}

var loginCmd = &cobra.Command{
	Use: "login",
	Run: func(cmd *cobra.Command, args []string) {
		creds, err := ReadCreds()
		if err != nil {
			log.Fatal(err)
		}

		a := api.New("http://127.0.0.1:3000")

		body := api.PostAuthSigninBody{
			Username: creds.Username,
			Password: creds.Password,
		}

		buf := bytes.Buffer{}

		err = json.NewEncoder(&buf).Encode(&body)
		if err != nil {
			log.Fatal(err)
		}

		res, err := api.Request[ApiResponse[api.PostAuthSignin, any]](a, "/api/v1/auth/signin", "POST", &buf)
		if err != nil {
			log.Fatal(err)
		}

		if res.Status == "success" {
			fmt.Printf("res.Data.Token: %v\n", res.Data.Token)
		}

		pretty.Println(res)
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
