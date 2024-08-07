// THIS FILE IS GENERATED BY PYRIN GOGEN CODE GENERATOR
package api

type Artist struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Picture string `json:"picture"`
}

type GetArtists struct {
	Artists []Artist `json:"artists"`
}

type GetArtistById Artist

type Album struct {
	Id string `json:"id"`
	Name string `json:"name"`
	CoverArt string `json:"coverArt"`
	ArtistId string `json:"artistId"`
}

type GetArtistAlbumsById struct {
	Albums []Album `json:"albums"`
}

type GetAlbums struct {
	Albums []Album `json:"albums"`
}

type GetAlbumById Album

type Track struct {
	Id string `json:"id"`
	Number int `json:"number"`
	Name string `json:"name"`
	CoverArt string `json:"coverArt"`
	Duration int `json:"duration"`
	BestQualityFile string `json:"bestQualityFile"`
	MobileQualityFile string `json:"mobileQualityFile"`
	AlbumId string `json:"albumId"`
	ArtistId string `json:"artistId"`
	AlbumName string `json:"albumName"`
	ArtistName string `json:"artistName"`
	Tags []string `json:"tags"`
	Genres []string `json:"genres"`
}

type GetAlbumTracksById struct {
	Tracks []Track `json:"tracks"`
}

type GetTracks struct {
	Tracks []Track `json:"tracks"`
}

type GetTrackById Track

type GetSync struct {
	IsSyncing int `json:"isSyncing"`
}

type PostQueue struct {
	Tracks []Track `json:"tracks"`
}

type Tag struct {
	Id string `json:"id"`
	Name string `json:"name"`
}

type GetTags struct {
	Tags []Tag `json:"tags"`
}

type PostAuthSignup struct {
	Id string `json:"id"`
	Username string `json:"username"`
}

type PostAuthSignupBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
	PasswordConfirm string `json:"passwordConfirm"`
}

type PostAuthSignin struct {
	Token string `json:"token"`
}

type PostAuthSigninBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type GetAuthMe struct {
	Id string `json:"id"`
	Username string `json:"username"`
}

type Playlist struct {
	Id string `json:"id"`
	Name string `json:"name"`
}

type GetPlaylists struct {
	Playlists []Playlist `json:"playlists"`
}

type PostPlaylist Playlist

type PostPlaylistBody struct {
	Name string `json:"name"`
}

type GetPlaylistById struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Items []Track `json:"items"`
}

type PostPlaylistItemsByIdBody struct {
	Tracks []string `json:"tracks"`
}

type DeletePlaylistItemsByIdBody struct {
	TrackIndices []int `json:"trackIndices"`
}

type PostPlaylistsItemMoveByIdBody struct {
	TrackId string `json:"trackId"`
	ToIndex int `json:"toIndex"`
}

type GetSystemInfo struct {
	Version string `json:"version"`
	IsSetup int `json:"isSetup"`
}

type PostSystemSetupBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
	PasswordConfirm string `json:"passwordConfirm"`
}

type ExportTrack struct {
	Name string `json:"name"`
	Album string `json:"album"`
	Artist string `json:"artist"`
}

type ExportPlaylist struct {
	Name string `json:"name"`
	Tracks []ExportTrack `json:"tracks"`
}

type ExportUser struct {
	Username string `json:"username"`
	Playlists []ExportPlaylist `json:"playlists"`
}

type PostSystemExport struct {
	Users []ExportUser `json:"users"`
}

