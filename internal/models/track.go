package models

import "encoding/json"

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

type TrackArtist struct {
	Name string
	URL  string
	MBID string
}

func (a *TrackArtist) UnmarshalJSON(data []byte) error {
	var raw struct {
		Text string `json:"#text"`
		Name string `json:"name"`
		URL  string `json:"url"`
		MBID string `json:"mbid"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	a.Name = raw.Name
	if a.Name == "" {
		a.Name = raw.Text
	}
	a.URL = raw.URL
	a.MBID = raw.MBID
	return nil
}

func (a TrackArtist) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name string `json:"name,omitempty"`
		URL  string `json:"url,omitempty"`
		MBID string `json:"mbid,omitempty"`
	}{
		Name: a.Name,
		URL:  a.URL,
		MBID: a.MBID,
	})
}

type TrackAlbum struct {
	Name string
	URL  string
	MBID string
}

func (a *TrackAlbum) UnmarshalJSON(data []byte) error {
	var raw struct {
		Text  string `json:"#text"`
		Title string `json:"title"`
		URL   string `json:"url"`
		MBID  string `json:"mbid"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	a.Name = raw.Title
	if a.Name == "" {
		a.Name = raw.Text
	}
	a.URL = raw.URL
	a.MBID = raw.MBID
	return nil
}

func (a TrackAlbum) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name string `json:"name,omitempty"`
		URL  string `json:"url,omitempty"`
		MBID string `json:"mbid,omitempty"`
	}{
		Name: a.Name,
		URL:  a.URL,
		MBID: a.MBID,
	})
}

type Track struct {
	Name   string `json:"name"`
	Url    string `json:"url"`
	Artist TrackArtist `json:"artist"`
	Album  TrackAlbum `json:"album"`
	Image []struct {
		Size string `json:"size"`
		Url  string `json:"#text"`
	} `json:"image"`
	Attr struct {
		NowPlaying string `json:"nowplaying"`
	} `json:"@attr,omitempty"`
}

type TrackGetInfoResponse struct {
	Track Track `json:"track"`
}