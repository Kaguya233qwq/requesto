# requesto

[![Go Reference](https://pkg.go.dev/badge/github.com/Kaguya233qwq/requesto.svg)](https://pkg.go.dev/github.com/Kaguya233qwq/requesto)

English | [[简体中文](/docs/README_zh-CN.md)] | [[繁體中文](/docs/README_zh-TW.md)] | [[日本語](/docs/README_ja.md)] | [[한국어](/docs/README_ko.md)] | [[Русский](/docs/README_ru.md)] | [[Español](/docs/README_es.md)] | [[Deutsch](/docs/README_de.md)]

> [!TIP]
> This project is under active development. If you have any great ideas or suggestions, feel free to open an issue or submit a pull request!

`requesto` is an elegant and powerful HTTP client library for Go. Built on top of the standard `net/http` library, it offers a fluent, chainable API, a robust middleware system, and a suite of convenient features. It's designed to make sending HTTP requests intuitive and to enhance the developer experience.

## 🧠 About the Name

- 1. There are just too many libraries named `requests`.
- 2. It has a lively and passionate feel, reminiscent of Romance languages, making it catchy and memorable.
- 3. You can also see it as a blend of `request` + `go`.


## ⭐ Features

*   📦️ Out-of-the-box convenience for simple requests, with powerful options for complex scenarios.
*   ✨ An intuitive, fluent, and chainable API.
*   🚀 Built-in support for JSON, x-www-form-urlencoded, binary data, and file uploads.
*   🍪 Automatic cookie jar management for sessions.
*   ⏱️ Timeout and cancellation control via `context.Context`.
*   🔧 Configurable clients (timeouts, redirect policies, etc.).
*   🧅 A powerful yet simple middleware (hook) system.

## 💿 Installation

```bash
go get github.com/Kaguya233qwq/requesto
```

## ⚡ Quick Start

Send a request and parse a JSON response in just a few lines:

```go
package main

import (
	"fmt"
	"log"

	"github.com/Kaguya233qwq/requesto"
)

func main() {
	resp, err := requesto.Get("https://httpbingo.org/get")
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}
	fmt.Printf("Status Code: %d\n", resp.StatusCode())
	
	// Parse the response body into a map
	jsonData, _ := resp.Json()
	fmt.Printf("Args from server: %v\n", jsonData["args"])
}
```

## 📚 Usage

### 🚀 Direct Requests

#### Passing Query Parameters

```go
params := requesto.Params{"q": "1"}
resp, err := requesto.Get("https://example.com", params)
if err != nil {
	log.Fatalf("Request failed: %v", err)
}
fmt.Printf("Status Code: %d\n", resp.StatusCode())
```
#### Setting Headers

```go
params := requesto.Params{"q": "1"}
headers := requesto.Headers{"Accept": "text/html,application/json"}
resp, err := requesto.Get("https://example.com", params, headers)
if err != nil {
	log.Fatalf("Request failed: %v", err)
}
fmt.Printf("Status Code: %d\n", resp.StatusCode())
```

#### Sending Form Data

```go
form := requesto.AsForm(map[string]string{"arg1": "value1"})
resp, err := requesto.Post("https://httpbingo.org/post", form)
if err != nil {
	log.Fatalf("Request failed: %v", err)
}
fmt.Printf("Status Code: %d\n", resp.StatusCode())
```

#### Sending JSON Data

```go
jsonData := requesto.AsJson(map[string]any{"arg1": "value1", "arg2": 0})
resp, err := requesto.Post("https://httpbingo.org/post", jsonData)
if err != nil {
	log.Fatalf("Request failed: %v", err)
}
fmt.Printf("Status Code: %d\n", resp.StatusCode())
```


### 🛠️ Using a Client

A `Client` is reusable and contains a connection pool, cookie jar, and global settings, which acts like a persistent session.

#### Default Client

```go
client := requesto.NewClient("https://example.com")
```

#### Custom Client

`NewClient` supports functional options to configure timeouts, redirect policies, and more.

```go
client := requesto.NewClient(
    "https://example.com",
    requesto.WithTimeout(10*time.Second),      // Set a global timeout of 10 seconds
    requesto.WithFollowRedirects(false),     // Disable automatic redirects
)
```

You can set or modify many client properties before making a request:

```go
client := requesto.NewClient("https://example.com")
// Set headers
client.Headers.Set("Content-Type", "text/html; charset=utf-8")
// Set query parameters
client.Params = map[string]string{
    "page": "1",
    "size": "10",
}
// Set cookies from a key-value map
client.SetCookiesFromMap(
    map[string]string{
        "__token": "xyz",
    },
)
client.Get()
```

#### Creating a New Request

Use `NewRequest()` to build a new request for more complex session control:

```go
client := requesto.NewClient("https://api.example.com")
req1 := client.NewRequest()
req2 := client.NewRequest()
// Requests from the same client share a connection pool but are otherwise independent.
```

Or, you can make a request directly:

```go
client := requesto.NewClient("https://example.com")
resp, err := client.Get()
```

#### Request Body

In client mode, `requesto` supports various ways to set the request body.

##### Sending JSON

Both structs and `map[string]any` are supported:

```go
type User struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

user := User{Name: "Alice", Age: 20}
client := requesto.NewClient("https://httpbingo.org")
resp, err := client.NewRequest().
    JoinPath("/post"). // Use JoinPath to dynamically append to the base URL
    SetJsonData(user). // Can be a struct or a map
    Post()
```

```go
user := map[string]any{"name": "Alice", "age": 20}
client := requesto.NewClient("https://httpbingo.org")
resp, err := client.NewRequest().
    JoinPath("/post").
    SetJsonData(user).
    Post()
```

##### Sending Form Data

```go
formData := map[string]string{
    "username": "bob",
    "password": "123",
}

resp, err := client.NewRequest().
    JoinPath("/post").
    SetFormData(formData).
    Post()
```

##### Uploading Files

```go
// Assuming 'hello.txt' exists and contains 'hello world'
resp, err := client.NewRequest().
    JoinPath("/post").
    SetFormData(map[string]string{
        "user_id": "123",
    }).
    SetFiles(map[string]string{
        "upload_file": "hello.txt",
    }).
    Post()
```

### Handling Responses

The `Response` object provides convenient methods for parsing the response body.

```go
resp, _ := client.NewRequest().Get()

// Get status code and headers
statusCode := resp.StatusCode()
headers := resp.Headers()

// Get the response body
text, _ := resp.Text()
bytes, _ := resp.Bytes()

// Parse JSON
var jsonData map[string]any
err := resp.Json(&jsonData) // It's better to pass a pointer to unmarshal into

// Parse cookies from the response
cookies := resp.Cookies()
cookiesMap := resp.CookiesMap()
```

### Cookie Management

The `Client` has a built-in `CookieJar` that automatically handles session cookies.

```go
// The server sets a cookie, and the client's jar will automatically store it.
client.NewRequest().JoinPath("/cookies/set").SetParams(map[string]string{
    "session_id": "my_secret_session",
}).Get()

// Now, access an authenticated endpoint. The cookie will be sent automatically.
resp, _ := client.NewRequest().JoinPath("/cookies").Get()
// You can inspect the cookies on the client
cookies := client.Cookies()

// You can also manually set cookies on the client
client.SetCookiesFromMap(map[string]string{
    "token": "123",
})
```

### Using Context for Timeouts and Cancellation

When you need to manage concurrent requests, use `NewRequestWithContext` to pass a context for per-request timeouts or cancellation signals.

```go
// Create a context that will be canceled after 2 seconds
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()

// This request will fail due to the context timeout
_, err := client.NewRequestWithContext(ctx).
    JoinPath("/delay/3"). // This endpoint waits for 3 seconds
    Get()

if errors.Is(err, context.DeadlineExceeded) {
    fmt.Println("Request timed out")
}
```

## Advanced: Middleware

Middleware is one of `requesto`'s most powerful features, allowing you to inject custom logic before a request is sent or after a response is received.

### Using Middleware

`requesto` provides some built-in middleware, such as a logger and a retrier.

```go
import "github.com/Kaguya233qwq/requesto/middleware"

client := requesto.NewClient("http://example.com")

// Apply middleware to the client
client.Use(
    middleware.NewLogger(
        middleware.WithLevel(middleware.LevelDebug),
        middleware.WithHeaders(true),
    ),
    middleware.NewRetrier(
        middleware.WithRetryCount(3),
        middleware.WithRetryOnServerErrors(), // Retry on 5xx status codes
    ),
)
```

### Writing Custom Middleware

The recommended way to create simple middleware is with the `NewHook` builder.

- #### Using `NewHook`

Here's an example of a middleware that adds an authentication header to every request:

```go
func AuthMiddleware(token string) requesto.Middleware {
    return middleware.NewHook(
        middleware.WithBeforeRequest(func(req *requesto.Request) error {
            req.SetHeader("Authorization", "Bearer " + token)
            return nil
        }),
    )
}

client.Use(AuthMiddleware("my_secret_token"))
```

- #### Implementing the Full Interface

For more complex logic that needs to control the request flow (like caching), you can implement the full `requesto.Middleware` interface:

```go
func MyComplexMiddleware() requesto.Middleware {
    return func(req *requesto.Request, next requesto.Next) (*requesto.Response, error) {
        // Logic before the request...

        // Call the next middleware in the chain
        resp, err := next(req)

        // Logic after the response...

        return resp, err
    }
}
```

## 📜 License

This project is licensed under the MIT License. See the LICENSE file for more details.

## 🙏 Acknowledgements

Special thanks to the following projects for providing inspiration and code references. The open-source community is a better place because of you.

*   [earthboundkid/requests](https://github.com/earthboundkid/requests)
*   [asmcos/requests](https://github.com/asmcos/requests)


