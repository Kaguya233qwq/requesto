package requesto

import (
	"net/http"
	"time"
)

// clientConfig holds the configuration for a Client. It's used internally
// by ClientOption functions to modify the client's settings.
type clientConfig struct {
	httpClient *http.Client
}

// ClientOption is a function that configures a Client.
type ClientOption func(*clientConfig)

// WithTimeout sets the global request timeout for the http.Client.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *clientConfig) {
		c.httpClient.Timeout = timeout
	}
}

// WithFollowRedirects controls whether the client should automatically follow
// HTTP redirects. The default behavior is to follow redirects (true).
func WithFollowRedirects(follow bool) ClientOption {
	return func(c *clientConfig) {
		if !follow {
			c.httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			}
		} else {
			// Set to nil to restore the default redirect behavior.
			c.httpClient.CheckRedirect = nil
		}
	}
}

// WithTransport allows setting a custom http.RoundTripper (transport) for the client.
func WithTransport(transport http.RoundTripper) ClientOption {
	return func(c *clientConfig) {
		if transport != nil {
			c.httpClient.Transport = transport
		}
	}
}
