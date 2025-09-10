# requesto

[![Go Reference](https://pkg.go.dev/badge/github.com/Kaguya233qwq/requesto.svg)](https://pkg.go.dev/github.com/Kaguya233qwq/requesto)

[[English](../README.md)] | [[简体中文](/docs/README_zh-CN.md)] | [[繁體中文](/docs/README_zh-TW.md)] | [[日本語](/docs/README_ja.md)] | 한국어 | [[Русский](/docs/README_ru.md)] | [[Español](/docs/README_es.md)] | [[Deutsch](/docs/README_de.md)]

> [!TIP]
> 이 프로젝트는 아직 활발히 개선 중입니다. 좋은 아이디어나 제안이 있다면 언제든지 이슈(issue)나 PR(pull request)을 남겨주세요.

`requesto`는 Go 언어를 위해 설계된 우아하고 강력한 HTTP 클라이언트 라이브러리입니다. 표준 라이브러리인 `net/http`를 기반으로 구축되었으며, 체인 API, 강력한 미들웨어 시스템 및 다양한 편의 기능을 제공합니다. 개발자가 가장 직관적인 방식으로 HTTP 요청을 보낼 수 있도록 하여 개발 경험을 향상시키는 것을 목표로 합니다.

## 🧠 이름에 대하여

- 1. `requests`라는 이름의 라이브러리가 너무 많습니다.
- 2. 라틴 계열 언어처럼 활기차고 열정적인 느낌을 주며, 듣기 좋고 기억하기 쉽습니다.
- 3. `request` + `go`의 조합으로도 이해할 수 있습니다.

## ⭐ 기능 및 특징

*   📦️ 간단한 요청 방식과 전문적인 방식을 모두 제공하여 다양한 시나리오를 손쉽게 만족시킵니다.
*   ✨ 가장 직관적인 체인 API 호출로 명확하고 부드러운 코딩이 가능합니다.
*   🚀 JSON, x-www-form-urlencoded, 바이너리 데이터 및 파일 업로드를 지원합니다.
*   🍪 `CookieJar`를 통한 자동화된 세션 관리를 제공합니다.
*   ⏱️ `Context`를 사용하여 타임아웃 및 취소를 제어합니다.
*   🔧 클라이언트 설정 가능 (타임아웃, 리다이렉션 정책 등).
*   🧅 강력하고 작성하기 쉬운 미들웨어(훅) 시스템을 갖추고 있습니다.

## 💿 설치

```bash
go get github.com/Kaguya233qwq/requesto
```

## ⚡ 빠른 시작

몇 줄의 코드로 요청을 보내고 JSON 응답을 파싱할 수 있습니다.

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
		log.Fatalf("요청 실패: %v", err)
	}
	fmt.Printf("상태 코드: %d\n", resp.StatusCode())
	// 응답 본문을 맵(map)으로 파싱
	jsonData, _ := resp.Json()
	fmt.Printf("Args from server: %v\n", jsonData["args"])
}
```

## 📚 사용법

### 🚀 직접 요청

#### 쿼리 파라미터 전달

```go
params := requesto.Params{"q": "1"}
resp, err = requesto.Get("https://example.com", params)
if err != nil {
	log.Fatalf("요청 실패: %v", err)
}
fmt.Printf("상태 코드: %d\n", resp.StatusCode())
```
#### 요청 헤더 설정

```go
params := requesto.Params{"q": "1"}
headers := requesto.Headers{"Accept": "text/html,application"}
resp, err = requesto.Get("https://example.com", params, headers)
if err != nil {
	log.Fatalf("요청 실패: %v", err)
}
fmt.Printf("상태 코드: %d\n", resp.StatusCode())
```

#### 폼(Form) 데이터 설정

```go
form := requesto.AsForm(map[string]string{"arg1": "value1"})
resp, err = requesto.Post("https://example.com", form)
if err != nil {
	log.Fatalf("요청 실패: %v", err)
}
fmt.Printf("상태 코드: %d\n", resp.StatusCode())
```

#### JSON 설정

```go
json := requesto.AsJson(map[string]any{"arg1": "value1", "arg2": 0})
resp, err = requesto.Post("https://example.com", json)
if err != nil {
	log.Fatalf("요청 실패: %v", err)
}
fmt.Printf("상태 코드: %d\n", resp.StatusCode())
```


### 🛠️ 클라이언트 요청

클라이언트(`Client`)는 재사용 가능한 객체로, 커넥션 풀, CookieJar 및 전역 설정을 포함하여 영구적인 세션처럼 작동합니다.

#### 기본 클라이언트

```go
client := requesto.NewClient("https://example.com")
```

#### 사용자 정의 클라이언트

`NewClient`는 함수형 옵션을 지원하여 타임아웃, 리다이렉션 정책 등을 설정할 수 있습니다.

```go
client := requesto.NewClient(
    "https://example.com",
    requesto.WithTimeout(10*time.Second),      // 전역 타임아웃 10초로 설정
    requesto.WithFollowRedirects(false),     // 자동 리다이렉션 비활성화
)
```

클라이언트의 많은 속성은 요청을 보내기 전에 직접 설정하거나 수정할 수 있습니다.

```go
client := requesto.NewClient("https://example.com")
// 요청 헤더 설정
client.Headers.Set("Content-Type", "text/html; charset=utf-8")
// 쿼리 파라미터 설정
client.Params = map[string]string{
    "page": "1",
    "size": "10",
}
// key-value 형태로 쿠키 설정
client.SetCookiesFromMap(
    map[string]string{
        "__token": "xyz",
    },
)
client.Get()
```

#### 새 요청 생성

`NewRequest()`를 사용하여 새 요청을 생성하면 더 복잡한 세션 제어가 가능합니다.

```go
client := requesto.NewClient("https://api.example.com")
req1 := client.NewRequest()
req2 := client.NewRequest()
// 같은 클라이언트에서 생성된 요청들은 커넥션 풀을 공유하지만 서로 간섭하지 않습니다.
```

요청을 직접 보낼 수도 있습니다.

```go
client := requesto.NewClient("https://example.com")
resp, err := client.Get()
```

#### 요청 본문 (Request Body)

클라이언트 모드에서 `requesto`는 다양한 방법으로 요청 본문을 설정할 수 있습니다.

#### JSON 전송

구조체(struct)와 `map[string]any` 타입을 모두 지원합니다.

```go
type User struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

user := User{Name: "Alice", Age: 20}
client := requesto.NewClient("https://example.com")
resp, err := client.NewRequest().
    JoinPath("/api"). // JoinPath 메서드로 동적으로 경로를 조합할 수 있습니다.
    SetJsonData(user). // 구조체나 맵을 사용할 수 있습니다.
    Post()
```

```go
user := map[string]any{"name": "Alice","age": 20}
client := requesto.NewClient("https://example.com")
resp, err := client.NewRequest().
    JoinPath("/api").
    SetJsonData(user).
    Post()
```

#### 폼 데이터 전송

```go
formData := map[string]string{
    "username": "bob",
    "password": "123",
}

resp, err := client.NewRequest().
    JoinPath("/login").
    SetFormData(formData).
    Post()
```

#### 파일 업로드

```go
resp, err := client.NewRequest().
    JoinPath("/api").
    SetFormData(map[string]string{
        "user_id": "123",
    }).
    SetFiles(map[string]string{
        "upload_file": "hello.txt",
    }).
    Post()
```

### 응답 처리

`Response` 객체는 응답 본문을 파싱하기 위한 편리한 메서드들을 제공합니다.

```go
resp, _ := client.NewRequest().Get()

// 응답 상태 코드 및 헤더 정보 가져오기
statusCode := resp.StatusCode()
headers := resp.Headers()

// 응답 본문 가져오기
text, _ := resp.Text()
bytes, _ := resp.Bytes()

// JSON 파싱
jsonData, _ := resp.Json()

// 응답에 포함된 쿠키 파싱
cookies := resp.Cookies()
cookiesMap := resp.CookiesMap()
```

### 쿠키 관리

`Client`에는 `CookieJar`가 내장되어 있어 세션을 자동으로 처리합니다.

```go
client.NewRequest().JoinPath("/cookies/set").SetParams(map[string]string{
    "session_id": "my_secret_session",
}).Get()

// 응답 헤더에 `set-cookie`가 있으면 해당 쿠키가 현재 클라이언트에 자동으로 설정됩니다.
resp, _ := client.NewRequest().JoinPath("/login").Get()
cookies := resp.Cookies()

// 클라이언트에 수동으로 쿠키를 설정할 수도 있습니다.
client.SetCookiesFromMap(map[string]string{
    "token": "123",
})
```

### Context를 사용한 타임아웃 및 취소

동시 요청을 처리해야 할 때, `NewRequestWithContext`를 사용하여 컨텍스트를 전달함으로써 개별 요청에 대한 타임아웃을 설정하거나 취소 신호를 보낼 수 있습니다.

```go
// 5초 후에 타임아웃되는 컨텍스트 생성
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// 이 요청은 컨텍스트 타임아웃으로 인해 실패합니다.
_, err := client.NewRequestWithContext(ctx).
    JoinPath("/delay/3").
    Get()

if errors.Is(err, context.DeadlineExceeded) {
    fmt.Println("요청 타임아웃")
}
```

## 고급 기능: 미들웨어

미들웨어(middleware)는 `requesto`의 강력한 고급 기능 중 하나로, 요청 전송 전후에 사용자 정의 로직을 주입할 수 있게 해줍니다.

### 미들웨어 사용하기

`requesto`는 재시도(retry), 로깅(logging) 등 바로 사용할 수 있는 몇 가지 미들웨어를 제공합니다.

```go
import "github.com/Kaguya233qwq/requesto/middleware"

client := requesto.NewClient("http://example.com")

client.Use(
    middleware.NewLogger(
        middleware.WithLevel(middleware.LevelDebug),
        middleware.WithHeaders(true),
    ),
    middleware.NewRetrier(
        middleware.WithRetryCount(3),
        middleware.WithRetryOnServerErrors(),
    ),
)
```

### 사용자 정의 미들웨어 작성하기

간단한 미들웨어는 `NewHook` 빌더를 사용하여 빠르게 구성하는 것을 추천합니다.

- #### `NewHook` 사용하기

아래는 모든 요청 전에 인증 헤더를 추가하는 미들웨어 예시입니다.

```go
func AuthMiddleware(token string) requesto.Middleware {
    return middleware.NewHook(
        middleware.WithBeforeRequest(func(req *requesto.Request) error {
            req.SetHeader(map[string]string{
                "Authorization": "Bearer " + token,
            })
            return nil
        }),
    )
}

client.Use(AuthMiddleware("my_secret_token"))
```

- #### 전체 인터페이스 구현하기

캐싱과 같이 요청 흐름을 제어해야 하는 복잡한 로직의 경우, `requesto.Middleware` 인터페이스를 직접 구현할 수 있습니다.

```go
func MyComplexMiddleware() requesto.Middleware {
    return func(req *requesto.Request, next requesto.Next) (*requesto.Response, error) {
        // 요청 전 로직...

        // 다음 미들웨어 호출
        resp, err := next(req)

        // 응답 후 로직...

        return resp, err
    }
}
```

## 📜 라이선스

이 프로젝트는 MIT 라이선스를 따릅니다.

## 🙏 감사 인사

다음 프로젝트들로부터 개발에 대한 영감과 코드 참조를 얻었습니다. 여러분 덕분에 오픈 소스 커뮤니티가 더욱 멋진 곳이 되었습니다.

*   [earthboundkid/requests](https://github.com/earthboundkid/requests)
*   [asmcos/requests](https://github.com/asmcos/requests)