package requesto

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

// Headers is a type alias for a map of request headers.
type Headers map[string]string

// Params is a type alias for a map of URL query parameters.
type Params map[string]string

// JsonData is a type alias for any data that can be marshaled to JSON.
type JsonData any

// FormData is a type alias for a map of form data.
type FormData map[string]string

// Files is a type alias for a map of files to be uploaded.
type Files map[string]File

// Client is a reusable HTTP client that manages a connection pool, cookies,
// and default settings for requests.
type Client struct {
	httpClient  *http.Client
	middlewares []Middleware
	CookieJar   http.CookieJar
	BaseURL     string
	Headers     http.Header
	Params      map[string]string
	FormData    map[string]string
	JsonData    any
	BodyBytes   []byte
	Files       map[string]File
}

// NewClient creates and returns a new Client instance.
// It initializes a default http.Client with a cookie jar and applies any
// provided ClientOption functions.
func NewClient(baseUrl string, opts ...ClientOption) *Client {

	jar, err := cookiejar.New(nil)
	if err != nil {
		// This should not happen with a nil PublicSuffixList, but we handle it gracefully.
		panic(err)
	}
	defaultTransport := http.DefaultTransport.(*http.Transport).Clone()
	defaultTransport.MaxIdleConnsPerHost = 100

	defaultHttpClient := &http.Client{
		Transport: defaultTransport,
		Timeout:   30 * time.Second,
		Jar:       jar,
	}

	// Create and apply client configuration from options.
	config := &clientConfig{
		httpClient: defaultHttpClient,
	}
	for _, opt := range opts {
		opt(config)
	}

	return &Client{
		httpClient:  config.httpClient,
		middlewares: make([]Middleware, 0),
		CookieJar:   jar,
		BaseURL:     baseUrl,
		Headers:     make(http.Header),
		Params:      make(map[string]string),
		JsonData:    nil,
		FormData:    make(map[string]string),
		BodyBytes:   []byte{},
		Files:       make(map[string]File),
	}
}

// NewRequestWithContext creates a new Request instance with the provided context.
// If the provided context is nil, it defaults to context.Background().
func (c *Client) NewRequestWithContext(ctx context.Context) *Request {
	if ctx == nil {
		ctx = context.Background()
	}

	return &Request{
		client:    c,
		ctx:       ctx, // Use the provided context.
		headers:   make(http.Header),
		params:    make(map[string]string),
		jsonData:  nil,
		formData:  make(map[string]string),
		bodyBytes: []byte{},
		files:     make(map[string]File),
	}
}

// NewRequest creates a new Request instance with a default background context.
func (c *Client) NewRequest() *Request {
	return c.NewRequestWithContext(context.Background())
}

// SetHeaders sets default headers that will be sent with every request from this client.
func (c *Client) SetHeaders(h map[string]string) {
	c.Headers = make(http.Header)
	for k, v := range h {
		c.Headers.Add(k, v)
	}
}

// SetParams sets default query parameters that will be added to every request from this client.
func (c *Client) SetParams(params map[string]string) {
	c.Params = params
}

// SetJsonData sets a default JSON body for requests made by this client.
func (c *Client) SetJsonData(data any) {
	c.JsonData = data
}

// SetFormData sets a default form-urlencoded body for requests made by this client.
func (c *Client) SetFormData(data map[string]string) {
	c.FormData = data
}

// SetBinary sets a default raw binary body for requests made by this client.
func (c *Client) SetBinary(data []byte) {
	c.BodyBytes = data
}

// SetFiles sets default files to be uploaded with requests made by this client.
func (c *Client) SetFiles(files map[string]File) {
	c.Files = files
}

// SetCookiesFromMap adds cookies to the client's cookie jar from a map.
// It requires the client to have a valid BaseURL and returns an error if it's missing or invalid.
func (c *Client) SetCookiesFromMap(cookies map[string]string) error {
	if c.BaseURL == "" {
		return errors.New("requesto: cannot set cookies, client BaseURL is not configured")
	}

	originalURL, err := url.Parse(c.BaseURL)
	if err != nil {
		return fmt.Errorf("requesto: client BaseURL is invalid: %w", err)
	}

	rootURL := &url.URL{
		Scheme: originalURL.Scheme,
		Host:   originalURL.Host,
	}

	if rootURL.Scheme == "" || rootURL.Host == "" {
		return fmt.Errorf("requesto: could not determine a valid scheme and host from BaseURL: %s", c.BaseURL)
	}

	var cookieObjects []*http.Cookie
	for name, value := range cookies {
		cookie := &http.Cookie{
			Name:  name,
			Value: value,
			Path:  "/",
		}
		cookieObjects = append(cookieObjects, cookie)
	}

	c.CookieJar.SetCookies(rootURL, cookieObjects)
	return nil
}

// Get creates and sends a new GET request using the client's default settings.
func (c *Client) Get() (*Response, error) {
	return c.NewRequest().Get()
}

// Post creates and sends a new POST request using the client's default settings.
func (c *Client) Post() (*Response, error) {
	return c.NewRequest().Post()
}

// Use adds one or more middleware handlers to the client's middleware chain.
func (c *Client) Use(middlewares ...Middleware) {
	c.middlewares = append(c.middlewares, middlewares...)
}
