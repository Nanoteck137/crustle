package api

import (
	"fmt"
	"net/http"
)

// func (c *Client) Login(body PostAuthSigninBody) (*PostAuthSignin, error) {
// 	url := c.addr + "/api/v1/auth/signin"
// 	return Request[PostAuthSignin](url, http.MethodPost, c.token, body)
// }
//
// func (c *Client) GetTracks(filter, sort string) (*GetTracks, error) {
// 	filter = url.QueryEscape(filter)
// 	sort = url.QueryEscape(sort)
//
// 	url := c.addr + fmt.Sprintf("/api/v1/tracks?filter=%s&sort=%s", filter, sort)
// 	return Request[GetTracks](url, http.MethodGet, c.token, nil)
// }
//
// func (c *Client) GetPlaylists() (*GetPlaylists, error) {
// 	url := c.addr + "/api/v1/playlists"
// 	return Request[GetPlaylists](url, http.MethodGet, c.token, nil)
// }


func (c *Client) GetPlaylistById(id string, options Options) (*GetPlaylistById, error) {
	path := fmt.Sprintf("/api/v1/playlists/%v", id)
	url, err := createUrl(c.addr, path, options.QueryParams)
	if err != nil {
		return nil, err
	}

	data := RequestData{
		Url: url,
		Method: http.MethodGet, 
		Token: c.token, 
		Body: nil,
	}
	return Request[GetPlaylistById](data)
}

