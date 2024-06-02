package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

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

type Client struct {
	addr string
}

func New(addr string) *Client {
	return &Client{
		addr: addr,
	}
}

func (c *Client) Login(body PostAuthSigninBody) (*PostAuthSignin, error) {
	url := c.addr + "/api/v1/auth/signin"
	return Request[PostAuthSignin](url, http.MethodPost, body)
}

func Request[D any](url, method string, body any) (*D, error) {

	var r io.Reader

	if body != nil {
		buf := bytes.Buffer{}

		err := json.NewEncoder(&buf).Encode(body)
		if err != nil {
			return nil, err
		}

		r = &buf
	}

	req, err := http.NewRequest(method, url, r)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res ApiResponse[D, any]
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	if res.Status == "error" {
		return nil, res.Error
	}

	return &res.Data, nil
}
