package requesto

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Response wraps the standard http.Response to provide convenient access to the
// response body and other common properties. It also stores any error that
// occurred during the reading of the response body.
type Response struct {
	Resp      *http.Response
	bodyBytes []byte
	err       error
}

// newResponse creates a new Response instance. It reads the entire response
// body into memory and closes it, making the body accessible for multiple reads.
func newResponse(resp *http.Response) *Response {
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()

	return &Response{
		Resp:      resp,
		bodyBytes: body,
		err:       err,
	}
}

// StatusCode returns the HTTP status code of the response.
// It returns -1 if an error occurred while reading the response body.
func (r *Response) StatusCode() int {
	if r.err != nil {
		return -1
	}
	return r.Resp.StatusCode
}

// Header returns the response headers.
func (r *Response) Header() http.Header {
	if r.err != nil {
		return nil
	}
	return r.Resp.Header
}

// Text returns the response body as a string.
func (r *Response) Text() (string, error) {
	if r.err != nil {
		return "", r.err
	}
	return string(r.bodyBytes), nil
}

// Bytes returns the response body as a byte slice.
func (r *Response) Bytes() ([]byte, error) {
	if r.err != nil {
		return nil, r.err
	}
	return r.bodyBytes, nil
}

// Json unmarshals the response body into a map[string]any.
func (r *Response) Json() (json_results map[string]any, err error) {
	if r.err != nil {
		return json_results, r.err
	}
	if len(r.bodyBytes) == 0 {
		return json_results, nil
	}
	err = json.Unmarshal(r.bodyBytes, &json_results)
	if err != nil {
		return json_results, ErrUnmarshallingJSON
	}
	return json_results, nil
}

// Unmarshal unmarshals the response body into the provided value `v`,
// which should be a pointer.
func (r *Response) Unmarshal(v any) error {
	if r.err != nil {
		return r.err
	}
	if len(r.bodyBytes) == 0 {
		return nil
	}
	if err := json.Unmarshal(r.bodyBytes, v); err != nil {
		return fmt.Errorf("%w: %w", ErrUnmarshallingStruct, err)
	}
	return nil
}

// Cookies returns a slice of *http.Cookie containing all cookies set by the
// server in the response.
func (r *Response) Cookies() []*http.Cookie {
	if r.err != nil || r.Resp == nil {
		// Returns an empty, non-nil slice to allow for safe iteration.
		return []*http.Cookie{}
	}
	return r.Resp.Cookies()
}

// CookiesMap parses all response cookies into a map of name-value pairs.
func (r *Response) CookiesMap() map[string]string {
	result := make(map[string]string)

	// Reuses the Cookies() method to adhere to the DRY principle.
	for _, cookie := range r.Cookies() {
		result[cookie.Name] = cookie.Value
	}

	return result
}

// ToStruct unmarshals the response body into a new instance of the specified
// generic type T.
func ToStruct[T any](r *Response) (T, error) {
	var result T

	if r.err != nil {
		return result, r.err
	}
	if len(r.bodyBytes) == 0 {
		return result, nil
	}
	target := new(T)
	if err := json.Unmarshal(r.bodyBytes, target); err != nil {
		return result, fmt.Errorf("%w: %w", ErrUnmarshallingStruct, err)
	}
	return *target, nil
}
