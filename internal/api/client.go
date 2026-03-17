package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const BaseURL = "https://ws.audioscrobbler.com/2.0/"
const browserUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/146.0.0.0 Safari/537.36"

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

func setBrowserHeaders(req *http.Request, accept string) {
	req.Header.Set("User-Agent", browserUserAgent)
	req.Header.Set("Accept", accept)
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("DNT", "1")
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
	setBrowserHeaders(req, "application/json,text/plain,*/*")

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
