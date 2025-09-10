package requesto

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
	"strings"
)

// Next defines the next handler in the middleware chain.
// A middleware must call the next function to pass the request along to the next handler.
type Next func(req *Request) (*Response, error)

// Middleware defines the function signature for a requesto middleware.
type Middleware func(req *Request, next Next) (*Response, error)

// Request represents a single HTTP request that can be configured and sent.
type Request struct {
	client    *Client
	ctx       context.Context
	method    string
	url       *url.URL
	headers   http.Header
	params    map[string]string
	jsonData  any
	formData  map[string]string
	bodyBytes []byte
	files     map[string]File
	err       error
}

// JoinPath intelligently appends a path segment to the request's URL.
// If the provided path `p` is an absolute URL, it will replace the current request URL entirely.
// Otherwise, it joins the path to the existing URL's path.
func (r *Request) JoinPath(p string) *Request {
	if r.err != nil {
		return r
	}

	// Check if p is an absolute URL.
	parsedPath, err := url.Parse(p)
	if err == nil && parsedPath.IsAbs() {
		// If it's an absolute URL, replace the current one.
		r.url = parsedPath
		return r
	}

	// Otherwise, perform the join logic.
	if r.url == nil {
		baseURL, err := url.Parse(r.client.BaseURL)
		if err != nil {
			r.err = err
			return r
		}
		r.url = baseURL
	}
	r.url.Path = path.Join(r.url.Path, p)
	return r
}

// SetURL sets the raw URL for the request, replacing any existing URL.
func (r *Request) SetURL(rawURL string) *Request {
	if r.err != nil {
		return r
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		r.err = err
		return r
	}
	r.url = u
	return r
}

// SetParams sets the URL query parameters for the request.
func (r *Request) SetParams(params map[string]string) *Request {
	if r.err != nil {
		return r
	}
	r.params = params
	return r
}

// SetHeaders sets the request headers.
func (r *Request) SetHeaders(h map[string]string) *Request {
	if r.err != nil {
		return r
	}
	r.headers = make(http.Header)
	for k, v := range h {
		r.headers.Add(k, v)
	}
	return r
}

// SetJsonData sets the request body to be JSON-encoded from the provided data (struct or map).
// It also sets the Content-Type header to "application/json; charset=utf-8".
func (r *Request) SetJsonData(data any) *Request {
	if r.err != nil {
		return r
	}
	r.jsonData = data

	r.headers.Set("Content-Type", "application/json; charset=utf-8")
	return r
}

// SetFormData sets the request body to be form-urlencoded from the provided map.
// It also sets the Content-Type header to "application/x-www-form-urlencoded".
func (r *Request) SetFormData(data map[string]string) *Request {
	if r.err != nil {
		return r
	}

	r.formData = data
	r.headers.Set("Content-Type", "application/x-www-form-urlencoded")
	return r
}

// SetBinary sets the request body to the provided raw byte slice.
// If the Content-Type header is not already set, it defaults to "application/octet-stream".
func (r *Request) SetBinary(data []byte) *Request {
	if r.err != nil {
		return r
	}
	r.bodyBytes = data
	if r.headers.Get("Content-Type") == "" {
		r.headers.Set("Content-Type", "application/octet-stream")
	}
	return r
}

// SetFiles sets the files to be uploaded as part of a multipart/form-data request.
func (r *Request) SetFiles(files map[string]File) *Request {
	if r.err != nil {
		return r
	}
	r.files = files
	return r
}

// SetCookiesFromMap adds cookies to the client's underlying cookie jar from a map.
// This requires the client to have a valid BaseURL to determine the cookie domain.
func (r *Request) SetCookiesFromMap(cookies map[string]string) *Request {
	if r.err != nil {
		return r
	}

	baseURL := r.client.BaseURL
	if baseURL == "" {
		r.err = errors.New("requesto: cannot set cookies, client BaseURL is not configured")
		return r
	}

	originalURL, err := url.Parse(baseURL)
	if err != nil {
		r.err = fmt.Errorf("requesto: client BaseURL is invalid: %w", err)
		return r
	}

	rootURL := &url.URL{
		Scheme: originalURL.Scheme,
		Host:   originalURL.Host,
	}

	if rootURL.Scheme == "" || rootURL.Host == "" {
		r.err = fmt.Errorf("requesto: could not determine a valid scheme and host from BaseURL: %s", baseURL)
		return r
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

	r.client.CookieJar.SetCookies(rootURL, cookieObjects)

	// Return the request instance to allow for continued chaining.
	return r
}

// Get sets the method to GET and sends the request.
func (r *Request) Get() (*Response, error) {
	r.method = "GET"
	return r.send()
}

// Post sets the method to POST and sends the request.
func (r *Request) Post() (*Response, error) {
	r.method = "POST"
	return r.send()
}

// Put sets the method to PUT and sends the request.
func (r *Request) Put() (*Response, error) {
	r.method = "PUT"
	return r.send()
}

// Delete sets the method to DELETE and sends the request.
func (r *Request) Delete() (*Response, error) {
	r.method = "DELETE"
	return r.send()
}

// Cookies parses and returns any cookies set in the request headers.
func (r *Request) Cookies() []*http.Cookie {
	dummyReq := &http.Request{Header: r.headers}
	return dummyReq.Cookies()
}

// CookiesMap parses any cookies set in the request headers and returns them as a map.
func (r *Request) CookiesMap() map[string]string {
	result := make(map[string]string)
	for _, cookie := range r.Cookies() {
		result[cookie.Name] = cookie.Value
	}
	return result
}

// buildURL constructs the final URL for the request.
func (r *Request) buildURL() (*url.URL, error) {
	// If no URL is set on the request, use the client's BaseURL.
	if r.url == nil || r.url.String() == "" {
		if r.client.BaseURL == "" {
			return nil, errors.New("requesto: no URL specified for the request and no BaseURL in client")
		}
		baseURL, err := url.Parse(r.client.BaseURL)
		if err != nil {
			return nil, err
		}
		r.url = baseURL
	}

	finalURL := r.url

	finalQuery := finalURL.Query()
	for key, value := range r.client.Params {
		finalQuery.Set(key, value)
	}
	for key, value := range r.params {
		finalQuery.Set(key, value)
	}
	finalURL.RawQuery = finalQuery.Encode()

	return finalURL, nil
}

// buildBody constructs the request's io.Reader body and determines its Content-Type.
func (r *Request) buildBody() (body io.Reader, contentType string, err error) {
	finalFormData := make(map[string]string)
	maps.Copy(finalFormData, r.client.FormData)
	maps.Copy(finalFormData, r.formData)

	finalFiles := make(map[string]File)
	maps.Copy(finalFiles, r.client.Files)
	maps.Copy(finalFiles, r.files)

	// The body is built based on priority:
	// Files (multipart) > Binary > JSON > Form Data
	if len(finalFiles) > 0 {
		return r.buildMultipartBody(finalFiles, finalFormData)
	}

	// Merge JsonData and FormData
	finalBodyBytes := r.bodyBytes
	if len(r.bodyBytes) == 0 {
		finalBodyBytes = r.client.BodyBytes
	}

	clientJson, clientIsMap := r.client.JsonData.(map[string]any)
	reqJson, reqIsMap := r.jsonData.(map[string]any)
	var finalJsonData any
	if r.client.JsonData != nil && r.jsonData != nil && clientIsMap && reqIsMap {
		mergedJson := make(map[string]any)
		maps.Copy(mergedJson, clientJson)
		maps.Copy(mergedJson, reqJson)
		finalJsonData = mergedJson
	} else if r.jsonData != nil {
		finalJsonData = r.jsonData
	} else {
		finalJsonData = r.client.JsonData
	}

	if len(finalBodyBytes) > 0 {
		return bytes.NewReader(finalBodyBytes), "", nil
	}

	if finalJsonData != nil {
		jsonDataBytes, err := json.Marshal(finalJsonData)
		if err != nil {
			return nil, "", err
		}
		if r.headers.Get("Content-Type") == "" {
			contentType = "application/json; charset=utf-8"
		}
		return bytes.NewReader(jsonDataBytes), contentType, nil
	}

	if len(finalFormData) > 0 {
		formData := url.Values{}
		for key, value := range finalFormData {
			formData.Set(key, value)
		}
		if r.headers.Get("Content-Type") == "" {
			contentType = "application/x-www-form-urlencoded"
		}
		return strings.NewReader(formData.Encode()), contentType, nil
	}

	return nil, "", nil
}

// buildMultipartBody creates a multipart/form-data body for file uploads.
func (r *Request) buildMultipartBody(files map[string]File, formData map[string]string) (body io.Reader, contentType string, err error) {
	bodyBuf := &bytes.Buffer{}
	writer := multipart.NewWriter(bodyBuf)

	for key, value := range formData {
		if err := writer.WriteField(key, value); err != nil {
			return nil, "", err
		}
	}

	for fieldName, file := range files {
		// If the content is an io.Closer (like *os.File), ensure it gets closed.
		if closer, ok := file.Content.(io.Closer); ok {
			defer closer.Close()
		}

		part, err := writer.CreateFormFile(fieldName, file.Name)
		if err != nil {
			return nil, "", err
		}

		if _, err = io.Copy(part, file.Content); err != nil {
			return nil, "", err
		}
	}

	if err := writer.Close(); err != nil {
		return nil, "", err
	}

	return bodyBuf, writer.FormDataContentType(), nil
}

// buildHeaders merges headers from the client and the request,
// with request-level headers taking precedence.
func (r *Request) buildHeaders() http.Header {
	finalHeader := r.client.Headers.Clone()
	maps.Copy(finalHeader, r.headers)
	return finalHeader
}

// send executes the request by building and running the middleware chain.
// All HTTP method functions (Get, Post, etc.) call this method.
func (r *Request) send() (*Response, error) {
	var terminator Next = func(req *Request) (*Response, error) {
		return req.do()
	}

	// Build the middleware chain in reverse.
	chain := terminator
	for i := len(r.client.middlewares) - 1; i >= 0; i-- {
		m := r.client.middlewares[i]
		chain = func(currentChain Next, currentMiddleware Middleware) Next {
			return func(req *Request) (*Response, error) {
				return currentMiddleware(req, currentChain)
			}
		}(chain, m)
	}

	return chain(r)
}

// do is the final step in the middleware chain. It builds the http.Request,
// sends it using the client's http.Client, and wraps the response.
func (r *Request) do() (*Response, error) {
	if r.err != nil {
		return nil, r.err
	}

	// Build the final URL.
	finalURL, err := r.buildURL()
	if err != nil {
		return nil, err
	}

	// Build the request body.
	body, contentType, err := r.buildBody()
	if err != nil {
		return nil, err
	}

	if contentType != "" {
		r.headers.Set("Content-Type", contentType)
	}

	// Create the standard http.Request.
	req, err := http.NewRequestWithContext(r.ctx, r.method, finalURL.String(), body)
	if err != nil {
		return nil, err
	}

	// Merge headers.
	req.Header = r.buildHeaders()

	// Send the request.
	resp, err := r.client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return newResponse(resp), nil
}
