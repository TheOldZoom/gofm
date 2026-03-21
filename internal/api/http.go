package api

import (
	"io"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/theOldZoom/gofm/internal/verbose"
)

const (
	httpMaxAttempts            = 3
	httpInitialRetryDelay      = 500 * time.Millisecond
	httpMaxRetryDelay          = 5 * time.Second
	httpConcurrentRequestLimit = 4
)

var httpRequestLimiter = make(chan struct{}, httpConcurrentRequestLimit)

func doRequestWithRetries(requestName string, buildRequest func() (*http.Request, error)) ([]byte, *http.Response, error) {
	var lastErr error

	for attempt := 1; attempt <= httpMaxAttempts; attempt++ {
		req, err := buildRequest()
		if err != nil {
			return nil, nil, err
		}

		resp, err := doLimitedRequest(req)
		if err != nil {
			lastErr = err
			if attempt == httpMaxAttempts || !isRetriableRequestError(err) {
				return nil, nil, err
			}

			delay := retryDelay(nil, attempt)
			verbose.Printf("%s request retrying after error (attempt %d/%d, wait=%s): %v", requestName, attempt, httpMaxAttempts, delay, err)
			time.Sleep(delay)
			continue
		}

		body, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			lastErr = readErr
			if attempt == httpMaxAttempts {
				return nil, resp, readErr
			}

			delay := retryDelay(resp, attempt)
			verbose.Printf("%s response read retry (attempt %d/%d, wait=%s): %v", requestName, attempt, httpMaxAttempts, delay, readErr)
			time.Sleep(delay)
			continue
		}

		if shouldRetryStatus(resp.StatusCode) && attempt < httpMaxAttempts {
			delay := retryDelay(resp, attempt)
			verbose.Printf("%s request retrying after status %d (attempt %d/%d, wait=%s)", requestName, resp.StatusCode, attempt, httpMaxAttempts, delay)
			time.Sleep(delay)
			continue
		}

		return body, resp, nil
	}

	return nil, nil, lastErr
}

func doLimitedRequest(req *http.Request) (*http.Response, error) {
	httpRequestLimiter <- struct{}{}
	defer func() { <-httpRequestLimiter }()

	return http.DefaultClient.Do(req)
}

func shouldRetryStatus(statusCode int) bool {
	return statusCode == http.StatusTooManyRequests ||
		statusCode == http.StatusRequestTimeout ||
		statusCode == http.StatusBadGateway ||
		statusCode == http.StatusServiceUnavailable ||
		statusCode == http.StatusGatewayTimeout
}

func isRetriableRequestError(err error) bool {
	if netErr, ok := err.(net.Error); ok {
		return netErr.Timeout() || netErr.Temporary()
	}

	return true
}

func retryDelay(resp *http.Response, attempt int) time.Duration {
	if resp != nil {
		if delay, ok := retryAfterDelay(resp); ok {
			return delay
		}
	}

	delay := httpInitialRetryDelay * time.Duration(1<<(attempt-1))
	if delay > httpMaxRetryDelay {
		return httpMaxRetryDelay
	}

	return delay
}

func retryAfterDelay(resp *http.Response) (time.Duration, bool) {
	if resp == nil {
		return 0, false
	}

	value := resp.Header.Get("Retry-After")
	if value == "" {
		return 0, false
	}

	if seconds, err := strconv.Atoi(value); err == nil {
		if seconds <= 0 {
			return 0, false
		}

		return time.Duration(seconds) * time.Second, true
	}

	if retryAt, err := http.ParseTime(value); err == nil {
		delay := time.Until(retryAt)
		if delay > 0 {
			return delay, true
		}
	}

	return 0, false
}
