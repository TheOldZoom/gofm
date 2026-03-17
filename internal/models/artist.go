package models

type TopArtistsResponse struct {
	TopArtists struct {
		Artist []Artist `json:"artist"`
	} `json:"topartists"`
}

type Artist struct {
	Name       string `json:"name"`
	PlayCount  string `json:"playcount"`
	MBID       string `json:"mbid"`
	Url        string `json:"url"`
	Streamable string `json:"streamable"`
	Image      []struct {
		Size string `json:"size"`
		Url  string `json:"#text"`
	} `json:"image"`
	Attr struct {
		Rank string `json:"rank"`
	} `json:"@attr,omitempty"`
}
