package api

type Artist struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type Album struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	CoverArt string `json:"coverArt"`
	ArtistId string `json:"artistId"`
}

type Track struct {
	Id                string `json:"id"`
	Number            int    `json:"number"`
	Name              string `json:"name"`
	CoverArt          string `json:"coverArt"`
	Duration          int    `json:"duration"`
	BestQualityFile   string `json:"bestQualityFile"`
	MobileQualityFile string `json:"mobileQualityFile"`
	AlbumId           string `json:"albumId"`
	ArtistId          string `json:"artistId"`
	AlbumName         string `json:"albumName"`
	ArtistName        string `json:"artistName"`
}

type Tag struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Playlist struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type GetArtists struct {
	Artists []Artist `json:"artists"`
}

type GetArtistById struct {
	Artist
}

type GetArtistAlbumsById struct {
	Albums []Album `json:"albums"`
}

type GetAlbums struct {
	Albums []Album `json:"albums"`
}

type GetAlbumById struct {
	Album
}

type GetAlbumTracksById struct {
	Tracks []Track `json:"tracks"`
}

type GetTracks struct {
	Tracks []Track `json:"tracks"`
}

type GetTrackById struct {
	Track
}

type GetSync struct {
	IsSyncing bool `json:"isSyncing"`
}

type PostQueue struct {
	Tracks []Track `json:"tracks"`
}

type GetTags struct {
	Tags []Tag `json:"tags"`
}

type PostAuthSignupBody struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"passwordConfirm"`
}

type PostAuthSignup struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

type PostAuthSigninBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type PostAuthSignin struct {
	Token string `json:"token"`
}

type GetAuthMe struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

type PostPlaylist Playlist

type PostPlaylistBody struct {
	Name string `json:"name"`
}

type PostPlaylistItemsByIdBody struct {
	Tracks []string `json:"tracks"`
}

type GetPlaylistById struct {
	Playlist

	Items []Track `json:"items"`
}

type GetPlaylists struct {
	Playlists []Playlist `json:"playlists"`
}

type DeletePlaylistItemsByIdBody struct {
	TrackIndices []int `json:"trackIndices"`
}

type PostPlaylistsItemMoveByIdBody struct {
	ItemIndex int `json:"itemIndex"`
	ToIndex   int `json:"toIndex"`
}
