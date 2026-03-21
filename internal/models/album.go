package models

type TopAlbumsResponse struct {
	TopAlbums struct {
		Album []Album `json:"album"`
	} `json:"topalbums"`
}

type Album struct {
	Name      string `json:"name"`
	PlayCount string `json:"playcount"`
	MBID      string `json:"mbid"`
	URL       string `json:"url"`
	Artist    TrackArtist `json:"artist"`
	Image     []struct {
		Size string `json:"size"`
		Url  string `json:"#text"`
	} `json:"image"`
	Attr struct {
		Rank string `json:"rank"`
	} `json:"@attr,omitempty"`
}

type AlbumGetInfoResponse struct {
	Album Album `json:"album"`
}
