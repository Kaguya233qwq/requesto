# requesto

[![Go Reference](https://pkg.go.dev/badge/github.com/Kaguya233qwq/requesto.svg)](https://pkg.go.dev/github.com/Kaguya233qwq/requesto)

[[English](../README.md)] | [[简体中文](/docs/README_zh-CN.md)] | 繁體中文 | [[日本語](/docs/README_ja.md)] | [[한국어](/docs/README_ko.md)] | [[Русский](/docs/README_ru.md)] | [[Español](/docs/README_es.md)] | [[Deutsch](/docs/README_de.md)]

> [!TIP]
> 本專案還在積極完善中，如果您有好的想法或建議，歡迎提交 Issue 或 PR。

`requesto` 是一個為 Go 語言設計的、優雅且強大的 HTTP 客戶端函式庫。它基於標準函式庫 `net/http` 建構，提供了鏈式 API、強大的中介軟體系統和一系列便利的功能，旨在以最符合開發者直覺的方式發起 HTTP 請求，提升開發體驗。

## 🧠 關於命名

- 1. 叫做 requests 的實在太多了
- 2. 有一種拉丁語系的活潑與熱情，好聽又好記
- 3. 其他理解：request+go 的結合

## ⭐ 功能特性

*   📦️ 提供開箱即用的請求方式與專業方式，輕鬆滿足不同場景
*   ✨ 最符合直覺的鏈式 API 呼叫，清晰又流暢
*   🚀 支援 JSON, x-www-form-urlencoded, 二進位及檔案上傳
*   🍪 自動化的 CookieJar 會話管理
*   ⏱️ 透過 `Context` 實作逾時和取消控制
*   🔧 可設定的客戶端（逾時、重新導向策略等）
*   🧅 強大且易於編寫的中介軟體（掛鉤）系統

## 💿 安裝

```bash
go get github.com/Kaguya233qwq/requesto
```

## ⚡ 快速開始

在幾行程式碼內傳送請求並解析 JSON 回應：

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
		log.Fatalf("請求失敗: %v", err)
	}
	fmt.Printf("狀態碼: %d\n", resp.StatusCode())
	// 將回應主體解析為 map
	jsonData, _ := resp.Json()
	fmt.Printf("Args from server: %v\n", jsonData["args"])
}
```

## 📚 使用方法

### 🚀 直接請求

#### 傳入查詢參數 (Query)

```go
params := requesto.Params{"q": "1"}
resp, err = requesto.Get("https://example.com", params)
if err != nil {
	log.Fatalf("請求失敗: %v", err)
}
fmt.Printf("狀態碼: %d\n", resp.StatusCode())
```
#### 設定請求標頭 (Header)

```go
params := requesto.Params{"q": "1"}
headers := requesto.Headers{"Accept": "text/html,application"}
resp, err = requesto.Get("https://example.com", params, headers)
if err != nil {
	log.Fatalf("請求失敗: %v", err)
}
fmt.Printf("狀態碼: %d\n", resp.StatusCode())
```

#### 設定 Form 表單參數

```go
form := requesto.AsForm(map[string]string{"arg1": "value1"})
resp, err = requesto.Post("https://example.com", form)
if err != nil {
	log.Fatalf("請求失敗: %v", err)
}
fmt.Printf("狀態碼: %d\n", resp.StatusCode())
```

#### 設定 JSON

```go
json := requesto.AsJson(map[string]any{"arg1": "value1", "arg2": 0})
resp, err = requesto.Post("https://example.com", json)
if err != nil {
	log.Fatalf("請求失敗: %v", err)
}
fmt.Printf("狀態碼: %d\n", resp.StatusCode())
```

### 🛠️ 客戶端請求

客戶端 (`Client`) 是可複用的，它包含連線池、CookieJar 和全域設定，相當於一個更持久的會話。

#### 預設客戶端

```go
client := requesto.NewClient("https://example.com")
```

#### 自訂客戶端

`NewClient` 支援函式式選項，可設定逾時、重新導向策略等。

```go
client := requesto.NewClient(
    "https://example.com",
    requesto.WithTimeout(10*time.Second),      // 設定全域逾時為 10 秒
    requesto.WithFollowRedirects(false),     // 禁止自動重新導向
)
```

Client 中的許多屬性在發起請求前均可直接進行設定和修改：

```go
client := requesto.NewClient("https://example.com")
//設定請求標頭
client.Headers.Set("Content-Type", "text/html; charset=utf-8")
//設定查詢參數
client.Params = map[string]string{
    "page": "1",
    "size": "10",
}
//以 kv 形式設定 cookie
client.SetCookiesFromMap(
    map[string]string{
        "__token": "xyz",
    },
)
client.Get()
```

#### 發起新請求

使用`NewRequest()`來建構一個新的請求，以便於實作更複雜的會話控制：

```go
client := requesto.NewClient("https://api.example.com")
req1 := client.NewRequest()
req2 := client.NewRequest()
// 同一個 client 下的 Request 共用一個連線池，但互不干擾
```

也可直接發起請求：

```go
client := requesto.NewClient("https://example.com")
resp, err := client.Get()
```

#### 請求主體 (Request Body)

在客戶端模式下，`requesto` 支援多種方式設定請求主體。

#### 傳送 JSON

支援結構體 (struct) 類型和 map[string]any 類型：

```go
type User struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

user := User{Name: "Alice", Age: 20}
client := requesto.NewClient("https://example.com")
resp, err := client.NewRequest().
    JoinPath("/api"). // 使用 JoinPath 方法可以動態拼接路徑
    SetJsonData(user). // 可以是 struct 或 map
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

#### 傳送 Form 表單

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

#### 上傳檔案

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

### 處理回應

`Response` 物件提供了一些便捷方法來解析回應主體。

```go
resp, _ := client.NewRequest().Get()

// 獲取回應狀態碼和標頭資訊
statusCode := resp.StatusCode()
headers := resp.Headers()

// 獲取回應主體
text, _ := resp.Text()
bytes, _ := resp.Bytes()

// 解析 JSON
jsonData, _ := resp.Json()

// 解析回應中的 Cookie
cookies := resp.Cookies()
cookiesMap := resp.CookiesMap()
```

### Cookie 管理

`Client` 內建了一個 `CookieJar`，可以自動處理會話。

```go
client.NewRequest().JoinPath("/cookies/set").SetParams(map[string]string{
    "session_id": "my_secret_session",
}).Get()

// 當回應標頭有 set-cookie 操作時，對應的 cookie 會被設定到當前的 client
resp, _ := client.NewRequest().JoinPath("/login").Get()
cookies := resp.Cookies()

// 也可以為客戶端手動設定 Cookie
client.SetCookiesFromMap(map[string]string{
    "token": "123",
})
```

### 使用 Context 進行逾時和取消

當您需要進行並行請求時，可以使用 `NewRequestWithContext` 傳入上下文來為單次請求設定逾時或傳遞取消訊號。

```go
// 建立一個 5 秒後會逾時的 context
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// 這個請求會因為 context 逾時而失敗
_, err := client.NewRequestWithContext(ctx).
    JoinPath("/delay/3").
    Get()

if errors.Is(err, context.DeadlineExceeded) {
    fmt.Println("請求逾時")
}
```

## 進階功能：中介軟體 (Middleware)

中介軟體（middleware）是 `requesto` 強大的進階功能之一，它允許您在請求傳送前後注入自訂邏輯。

### 使用中介軟體

`requesto` 提供了一些開箱即用的中介軟體，如重試和日誌。

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

### 編寫使用者自訂中介軟體

推薦使用 `NewHook` 建構器來快速建構簡單的中介軟體。

- #### 使用 `NewHook`

下面一個在每個請求前新增認證標頭的中介軟體範例：

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

- #### 實作完整介面

對於需要控制請求流程的複雜邏輯（如快取），您可以實作完整的 `requesto.Middleware` 介面：

```go
func MyComplexMiddleware() requesto.Middleware {
    return func(req *requesto.Request, next requesto.Next) (*requesto.Response, error) {
        // 請求前邏輯...

        // 呼叫下一個中介軟體
        resp, err := next(req)

        // 回應後邏輯...

        return resp, err
    }
}
```

## 📜 授權條款

本專案採用 MIT 開源授權條款。

## 🙏 致謝

感謝以下專案為本專案提供開發靈感與程式碼參考，正是因為有你們的存在才使得開源社群更加美好。

*   [earthboundkid/requests](https://github.com/earthboundkid/requests)
*   [asmcos/requests](https://github.com/asmcos/requests)