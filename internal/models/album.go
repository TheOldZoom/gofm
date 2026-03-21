package models

import "encoding/json"

type Tag struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type AlbumTags struct {
	Tag []Tag `json:"tag"`
}

func (t *AlbumTags) UnmarshalJSON(data []byte) error {
	var empty string
	if err := json.Unmarshal(data, &empty); err == nil {
		t.Tag = nil
		return nil
	}

	type alias AlbumTags
	var raw alias
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	t.Tag = raw.Tag
	return nil
}

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
	Tags AlbumTags `json:"tags"`
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
