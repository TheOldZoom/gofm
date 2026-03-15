package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const BaseURL = "https://ws.audioscrobbler.com/2.0/"

type Client struct {
	ApiKey string
}

type APIError struct {
	Code       int
	Message    string
	StatusCode int
}

func (e *APIError) Error() string {
	if e == nil {
		return ""
	}
	if e.Message == "" {
		return fmt.Sprintf("last.fm API error %d", e.Code)
	}
	return fmt.Sprintf("last.fm API error %d: %s", e.Code, e.Message)
}

type apiErrorResponse struct {
	Error   int    `json:"error"`
	Message string `json:"message"`
}

func (c *Client) Get(method string, params map[string]string, out any) error {
	q := url.Values{}
	q.Set("method", method)
	q.Set("api_key", c.ApiKey)
	q.Set("format", "json")

	for k, v := range params {
		q.Set(k, v)
	}

	req, err := http.NewRequest(http.MethodGet, BaseURL+"?"+q.Encode(), nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "gofm")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var apiErr apiErrorResponse
	if err := json.Unmarshal(body, &apiErr); err == nil && apiErr.Error != 0 {
		return &APIError{
			Code:       apiErr.Error,
			Message:    apiErr.Message,
			StatusCode: resp.StatusCode,
		}
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status code %d", resp.StatusCode)
	}

	return json.Unmarshal(body, out)
}
