package requesto

import (
	"fmt"
	"maps"
	"net/http"
)

// bodyArgument is a marker interface used to identify different types of
// request arguments, such as JSON bodies, form data, or query parameters.
type bodyArgument interface {
	isArgument()
}

// requestHeaders is an internal type that wraps http.Header.
type requestHeaders struct{ data http.Header }

func (r requestHeaders) isArgument() {}

// queryParams is an internal type that wraps a map for URL query parameters.
type queryParams struct{ data map[string]string }

func (q queryParams) isArgument() {}

// jsonBody is an internal type that wraps any data to be sent as a JSON body.
type jsonBody struct{ data any }

func (j jsonBody) isArgument() {}

// formBody is an internal type that wraps a map for a form-urlencoded body.
type formBody struct{ data map[string]string }

func (f formBody) isArgument() {}

// filesBody is an internal type that wraps a map for file uploads.
type filesBody struct{ data map[string]File }

func (f filesBody) isArgument() {}

// AsHeaders marks the provided map to be used as HTTP request headers.
func AsHeaders(data map[string]string) bodyArgument {
	var headers = http.Header{}
	for key, value := range data {
		headers.Set(key, value)
	}
	return requestHeaders{data: headers}
}

// AsParams marks the provided map to be used as URL query parameters.
func AsParams(data map[string]string) bodyArgument {
	return queryParams{data: data}
}

// AsJson marks the provided data to be sent as a JSON request body.
func AsJson(data any) bodyArgument {
	return jsonBody{data: data}
}

// AsForm marks the provided map to be sent as a form-urlencoded request body.
func AsForm(data map[string]string) bodyArgument {
	return formBody{data: data}
}

// AsFiles marks the provided map to be used for file uploads.
func AsFiles(data map[string]File) bodyArgument {
	return filesBody{data: data}
}

// Get sends a GET request to the specified URL.
// It accepts optional arguments, which can be of type Params or Headers.
// Body-related arguments are ignored. An error is returned for any unsupported argument type.
func Get(URL string, args ...any) (*Response, error) {

	client := NewClient(URL)
	for _, arg := range args {
		switch v := arg.(type) {
		case Headers:
			for key, values := range v {
				client.Headers.Set(key, values)
			}
		case requestHeaders:
			client.Headers = v.data
		case queryParams:
			client.Params = v.data
		case Params:
			maps.Copy(client.Params, v)
		case JsonData, FormData, Files, jsonBody, formBody, filesBody:
			// Ignore body-related arguments for a GET request.
			continue
		case nil:
			continue
		default:
			return nil, fmt.Errorf("unsupported argument type for Get: %T", v)
		}
	}
	return client.Get()
}

// Post sends a convenient POST request.
//
// It accepts a URL and one or more optional arguments, which can be:
// - Headers: to set request headers.
// - Params: to set URL query parameters.
// - AsJson, AsForm, or AsFiles: to set the request body.
func Post(URL string, args ...any) (*Response, error) {
	client := NewClient(URL)

	for _, arg := range args {
		switch v := arg.(type) {
		case Headers:
			for key, value := range v {
				client.Headers.Set(key, value)
			}
		case Params:
			maps.Copy(client.Params, v)
		case bodyArgument:
			// Handle different body argument types.
			switch body := v.(type) {
			case jsonBody:
				client.JsonData = body.data
			case formBody:
				maps.Copy(client.FormData, body.data)
			case filesBody:
				maps.Copy(client.Files, body.data)
			case queryParams:
				client.Params = body.data
			case requestHeaders:
				client.Headers = body.data
			}
		case nil:
			continue
		default:
			return nil, fmt.Errorf("unsupported argument type for Post: %T", v)
		}
	}

	return client.Post()
}
