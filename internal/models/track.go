package models

type RecentTracksResponse struct {
	RecentTracks struct {
		Track []Track `json:"track"`
	} `json:"recenttracks"`
}

type TopTracksResponse struct {
	TopTracks struct {
		Track []Track `json:"track"`
	} `json:"toptracks"`
}

type Track struct {
	Name   string `json:"name"`
	Url    string `json:"url"`
	Artist struct {
		Name string `json:"#text"`
	} `json:"artist"`
	Album struct {
		Name string `json:"#text"`
	} `json:"album"`
	Image []struct {
		Size string `json:"size"`
		Url  string `json:"#text"`
	} `json:"image"`
	Attr struct {
		NowPlaying string `json:"nowplaying"`
	} `json:"@attr,omitempty"`
}
