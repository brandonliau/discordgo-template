package httpx

import (
	"net/http"
)

type ClientOption func(*Client)

func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

func WithRetryPolicy(retryPolicy RetryPolicy) ClientOption {
	return func(c *Client) {
		c.retryPolicy = retryPolicy
	}
}

func WithBackoffPolicy(backoffPolicy BackoffPolicy) ClientOption {
	return func(c *Client) {
		c.backoffPolicy = backoffPolicy
	}
}
