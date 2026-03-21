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

type StringValue string

func (s *StringValue) UnmarshalJSON(data []byte) error {
	var text string
	if err := json.Unmarshal(data, &text); err == nil {
		*s = StringValue(text)
		return nil
	}

	var number json.Number
	if err := json.Unmarshal(data, &number); err == nil {
		*s = StringValue(number.String())
		return nil
	}

	return json.Unmarshal(data, (*string)(s))
}

type TrackStreamable struct {
	Text      string
	FullTrack string
}

func (s *TrackStreamable) UnmarshalJSON(data []byte) error {
	var text string
	if err := json.Unmarshal(data, &text); err == nil {
		s.Text = text
		return nil
	}

	var raw struct {
		Text      string `json:"#text"`
		FullTrack string `json:"fulltrack"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	s.Text = raw.Text
	s.FullTrack = raw.FullTrack
	return nil
}

type TrackArtist struct {
	Name string
	URL  string
	MBID string
}

func (a *TrackArtist) UnmarshalJSON(data []byte) error {
	var text string
	if err := json.Unmarshal(data, &text); err == nil {
		a.Name = text
		return nil
	}

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
	Name   string
	URL    string
	MBID   string
	Artist string
	Image  []struct {
		Size string `json:"size"`
		Url  string `json:"#text"`
	}
}

func (a *TrackAlbum) UnmarshalJSON(data []byte) error {
	var raw struct {
		Text   string `json:"#text"`
		Title  string `json:"title"`
		URL    string `json:"url"`
		MBID   string `json:"mbid"`
		Artist string `json:"artist"`
		Image  []struct {
			Size string `json:"size"`
			Url  string `json:"#text"`
		} `json:"image"`
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
	a.Artist = raw.Artist
	a.Image = raw.Image
	return nil
}

func (a TrackAlbum) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name   string `json:"name,omitempty"`
		URL    string `json:"url,omitempty"`
		MBID   string `json:"mbid,omitempty"`
		Artist string `json:"artist,omitempty"`
		Image  []struct {
			Size string `json:"size"`
			Url  string `json:"#text"`
		} `json:"image,omitempty"`
	}{
		Name:   a.Name,
		URL:    a.URL,
		MBID:   a.MBID,
		Artist: a.Artist,
		Image:  a.Image,
	})
}

type Track struct {
	Name       string          `json:"name"`
	MBID       string          `json:"mbid"`
	Url        string          `json:"url"`
	Duration   StringValue     `json:"duration"`
	Streamable TrackStreamable `json:"streamable"`
	Listeners  string          `json:"listeners"`
	Artist     TrackArtist     `json:"artist"`
	Album      TrackAlbum      `json:"album"`
	Image      []struct {
		Size string `json:"size"`
		Url  string `json:"#text"`
	} `json:"image"`
	Attr struct {
		NowPlaying string `json:"nowplaying"`
	} `json:"@attr,omitempty"`
	PlayCount     string `json:"playcount"`
	UserPlayCount string `json:"userplaycount"`
	UserLoved     string `json:"userloved"`
	TopTags       struct {
		Tag []struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"tag"`
	} `json:"toptags"`
	Wiki struct {
		Published string `json:"published"`
		Summary   string `json:"summary"`
		Content   string `json:"content"`
	} `json:"wiki"`
}

type TrackGetInfoResponse struct {
	Track Track `json:"track"`
}
