package api

import (
	"encoding/json"
	"io"
	"net/http"
)

type Api struct {
	addr string
}

func New(addr string) *Api {
	return &Api{
		addr: addr,
	}
}

func Request[T any](api *Api, endpoint string, method string, body io.Reader) (*T, error) {
	url := api.addr + endpoint
	req, err := http.NewRequest(method, url, body)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res T
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
