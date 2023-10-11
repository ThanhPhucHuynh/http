package httpclient

import (
	"time"

	avhttp "github.com/ThanhPhucHuynh/http"
)

// Option represents the client options
type Option func(*Client)

// WithHTTPTimeout sets hystrix timeout
func WithHTTPTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.timeout = timeout
	}
}

// WithRetryCount sets the retry count for the hystrixHTTPClient
func WithRetryCount(retryCount int) Option {
	return func(c *Client) {
		c.retryCount = retryCount
	}
}

// WithRetrier sets the strategy for retrying
func WithRetrier(retrier avhttp.Retriable) Option {
	return func(c *Client) {
		c.retrier = retrier
	}
}

// WithHTTPClient sets a custom http client
func WithHTTPClient(client avhttp.Doer) Option {
	return func(c *Client) {
		c.client = client
	}
}
