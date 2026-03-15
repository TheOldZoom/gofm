package models

type UserGetInfoResponse struct {
	User struct {
		Name     string `json:"name"`
		RealName string `json:"realname"`
		Url      string `json:"url"`
		Image    []struct {
			Size string `json:"size"`
			Url  string `json:"#text"`
		} `json:"image"`
	} `json:"user"`
}
