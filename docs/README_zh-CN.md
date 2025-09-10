# requesto

[![Go Reference](https://pkg.go.dev/badge/github.com/Kaguya233qwq/requesto.svg)](https://pkg.go.dev/github.com/Kaguya233qwq/requesto)

[[English](../README.md)] | 简体中文 | [[繁體中文](/docs/README_zh-TW.md)] | [[日本語](/docs/README_ja.md)] | [[한국어](/docs/README_ko.md)] | [[Русский](/docs/README_ru.md)] | [[Español](/docs/README_es.md)] | [[Deutsch](/docs/README_de.md)]

> [!TIP]
> 本项目还在积极完善中，如果你有好的想法或建议，欢迎提交issue或者pr。

`requesto` 是一个为 Go 语言设计的、优雅且强大的 HTTP 客户端库。基于标准库 `net/http` 构建，提供了链式 API、强大的中间件系统和一系列便利的功能，旨在以最符合开发者直觉的方式发起HTTP请求，提升开发体验。

## 🧠关于命名

- 1. 叫requests的实在太多了

- 2. 有一种拉丁语系的活泼与热情，好听又好记

- 3. 其他理解：request+go的结合


## ⭐功能特性

*   📦️ 提供开箱即用的请求方式和专业方式，轻松满足不同场景
*   ✨ 最符合直觉的链式 API调用，清晰又流畅
*   🚀 支持 JSON, x-www-form-urlencoded, 二进制及文件上传
*   🍪 自动化的 CookieJar 会话管理
*   ⏱️ 通过 `Context` 实现超时和取消控制
*   🔧 可配置的客户端（超时、重定向策略等）
*   🧅 强大且易于编写的中间件（钩子）系统

## 💿安装

```bash
go get github.com/Kaguya233qwq/requesto
```

## ⚡快速开始

在几行代码内发送请求并解析 JSON 响应：

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
		log.Fatalf("请求失败: %v", err)
	}
	fmt.Printf("状态码: %d\n", resp.StatusCode())
	// 将响应体解析为 map
	jsonData, _ := resp.Json()
	fmt.Printf("Args from server: %v\n", jsonData["args"])
}
```

## 📚使用方法

### 🚀直接请求

#### 传入query参数

```go
params := requesto.Params{"q": "1"}
resp, err = requesto.Get("https://example.com", params)
if err != nil {
	log.Fatalf("请求失败: %v", err)
}
fmt.Printf("状态码: %d\n", resp.StatusCode())
```
#### 设置请求头

```go
params := requesto.Params{"q": "1"}
headers := requesto.Headers{"Accept": "text/html,application"}
resp, err = requesto.Get("https://example.com", params, headers)
if err != nil {
	log.Fatalf("请求失败: %v", err)
}
fmt.Printf("状态码: %d\n", resp.StatusCode())
```

#### 设置Form表单参数

```go
form := requesto.AsForm(map[string]string{"arg1": "value1"})
resp, err = requesto.Post("https://example.com", form)
if err != nil {
	log.Fatalf("请求失败: %v", err)
}
fmt.Printf("状态码: %d\n", resp.StatusCode())
```

#### 设置Json

```go
json := requesto.AsJson(map[string]any{"arg1": "value1", "arg2": 0})
resp, err = requesto.Post("https://example.com", json)
if err != nil {
	log.Fatalf("请求失败: %v", err)
}
fmt.Printf("状态码: %d\n", resp.StatusCode())
```


### 🛠️客户端请求

客户端 (`Client`) 是可复用的，它包含连接池、CookieJar 和全局配置，相当于一个更加持久的会话

#### 默认客户端

```go
client := requesto.NewClient("https://example.com")
```

#### 自定义客户端

`NewClient` 支持函数式选项，可配置超时、重定向策略等。

```go
client := requesto.NewClient(
    "https://example.com",
    requesto.WithTimeout(10*time.Second),      // 设置全局超时为 10 秒
    requesto.WithFollowRedirects(false),     // 禁止自动重定向
)
```

client中的很多属性在发起请求前均可直接进行设置和修改：

```go
client := requesto.NewClient("https://example.com")
//设置请求头
client.Headers.Set("Content-Type", "text/html; charset=utf-8")
//设置查询参数
client.Params = map[string]string{
    "page": "1",
    "size": "10",
}
//以kv形式设置cookie
client.SetCookiesFromMap(
    map[string]string{
        "__token": "xyz",
    },
)
client.Get()
```

#### 发起新请求

使用`NewRequest()`来构建一个新的请求，以便于实现更加复杂的会话控制：

```go
client := requesto.NewClient("https://api.example.com")
req1 := client.NewRequest()
req2 := client.NewRequest()
// 同一个client下的Request共用一个连接池，但互不干扰
```

也可直接发起请求：

```go
client := requesto.NewClient("https://example.com")
resp, err := client.Get()
```

#### 请求体

在client模式下，`requesto` 支持多种方式设置请求体。

#### 发送 JSON

支持结构体类型和map[string]any类型：

```go
type User struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

user := User{Name: "Alice", Age: 20}
client := requesto.NewClient("https://example.com")
resp, err := client.NewRequest().
    JoinPath("/api"). // 使用JoinPath方法可以动态拼接路径
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

#### 发送 Form 表单

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

#### 上传文件

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

### 处理响应

`Response` 对象提供了一些便捷方法来解析响应体。

```go
resp, _ := client.NewRequest().Get()

// 获取响应状态码和头信息
statusCode := resp.StatusCode()
headers := resp.Headers()

// 获取响应体
text, _ := resp.Text()
bytes, _ := resp.Bytes()

// 解析 JSON
jsonData, _ := resp.Json()

// 解析响应中的 Cookie
cookies := resp.Cookies()
cookiesMap := resp.CookiesMap()
```

### Cookie 管理

`Client` 内置了一个 `CookieJar`，可以自动处理会话。

```go
client.NewRequest().JoinPath("/cookies/set").SetParams(map[string]string{
    "session_id": "my_secret_session",
}).Get()

// 访问需要登录认证的接口，当响应头有set-cookie操作时，对应的cookie会被设置到当前的client
resp, _ := client.NewRequest().JoinPath("/login").Get()
cookies := resp.Cookies()
client.

// 也可以为客户端手动设置 Cookie
client.SetCookiesFromMap(map[string]string{
    "token": "123",
})
```

### 使用 Context 进行超时和取消

当你需要进行并发请求时，可以使用 `NewRequestWithContext`传入上下文来为单次请求设置超时或传递取消信号。

```go
// 创建一个 5 秒后会超时的 context
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// 这个请求会因为 context 超时而失败
_, err := client.NewRequestWithContext(ctx).
    JoinPath("/delay/3").
    Get()

if errors.Is(err, context.DeadlineExceeded) {
    fmt.Println("请求超时")
}
```

## 进阶功能：中间件

中间件（middleware）是 `requesto` 强大的进阶功能之一，它允许你在请求发送前后注入自定义逻辑。

### 使用中间件

`requesto` 提供了一些开箱即用的中间件，如重试和日志。

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

### 编写用户自定义中间件

推荐使用 `NewHook` 构建器来快速构建简单的中间件

- #### 使用 `NewHook`

下面一个在每个请求前添加认证头的中间件示例：

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

- #### 实现完整接口

对于需要控制请求流的复杂逻辑（如缓存），你可以实现完整的 `requesto.Middleware` 接口：

```go
func MyComplexMiddleware() requesto.Middleware {
    return func(req *requesto.Request, next requesto.Next) (*requesto.Response, error) {
        // 请求前逻辑...

        // 调用下一个中间件
        resp, err := next(req)

        // 响应后逻辑...

        return resp, err
    }
}
```

## 📜 License

本项目使用MIT开源协议

## 🙏 致谢

感谢以下项目为本项目提供开发灵感与代码参考，正是因为有你们的存在才使得开源社区更加美好。

*   [earthboundkid/requests](https://github.com/earthboundkid/requests)
*   [asmcos/requests](https://github.com/asmcos/requests)