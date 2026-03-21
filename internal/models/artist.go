package models

type ArtistGetInfoResponse struct {
	Artist Artist `json:"artist"`
}

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
	OnTour     string `json:"ontour"`
	Image      []struct {
		Size string `json:"size"`
		Url  string `json:"#text"`
	} `json:"image"`
	Stats struct {
		Listeners     string `json:"listeners"`
		PlayCount     string `json:"playcount"`
		UserPlayCount string `json:"userplaycount"`
	} `json:"stats"`
	Similar struct {
		Artist []struct {
			Name  string `json:"name"`
			MBID  string `json:"mbid"`
			Match string `json:"match"`
			Url   string `json:"url"`
			Image []struct {
				Size string `json:"size"`
				Url  string `json:"#text"`
			} `json:"image"`
		} `json:"artist"`
	} `json:"similar"`
	Tags struct {
		Tag []struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"tag"`
	} `json:"tags"`
	Bio struct {
		Links struct {
			Link struct {
				Text string `json:"#text"`
				Rel  string `json:"rel"`
				Href string `json:"href"`
			} `json:"link"`
		} `json:"links"`
		Published string `json:"published"`
		Summary   string `json:"summary"`
		Content   string `json:"content"`
	} `json:"bio"`
	Attr struct {
		Rank string `json:"rank"`
	} `json:"@attr,omitempty"`
}
