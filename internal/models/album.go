package models

type TopAlbumsResponse struct {
	TopAlbums struct {
		Album []Album `json:"album"`
	} `json:"topalbums"`
}

type Album struct {
	Name          string      `json:"name"`
	PlayCount     StringValue `json:"playcount"`
	Listeners     StringValue `json:"listeners"`
	UserPlayCount StringValue `json:"userplaycount"`
	MBID          string      `json:"mbid"`
	URL           string      `json:"url"`
	Artist        TrackArtist `json:"artist"`
	Releasedate   string      `json:"releasedate"`
	Image         []struct {
		Size string `json:"size"`
		Url  string `json:"#text"`
	} `json:"image"`
	Tracks struct {
		Track []Track `json:"track"`
	} `json:"tracks"`
	Tags struct {
		Tag []struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"tag"`
	} `json:"tags"`
	Wiki struct {
		Published string `json:"published"`
		Summary   string `json:"summary"`
		Content   string `json:"content"`
	} `json:"wiki"`
	Attr struct {
		Rank string `json:"rank"`
	} `json:"@attr,omitempty"`
}

type AlbumGetInfoResponse struct {
	Album Album `json:"album"`
}
