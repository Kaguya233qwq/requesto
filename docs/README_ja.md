# requesto

[![Go Reference](https://pkg.go.dev/badge/github.com/Kaguya233qwq/requesto.svg)](https://pkg.go.dev/github.com/Kaguya233qwq/requesto)

[[English](../README.md)] | [[简体中文](/docs/README_zh-CN.md)] | [[繁體中文](/docs/README_zh-TW.md)] | 日本語 | [[한국어](/docs/README_ko.md)] | [[Русский](/docs/README_ru.md)] | [[Español](/docs/README_es.md)] | [[Deutsch](/docs/README_de.md)]

> [!TIP]
> このプロジェクトは現在も積極的に開発が進められています。素晴らしいアイデアや提案がありましたら、Issue の提出やプルリクエストをお待ちしています。

`requesto` は、Go言語向けに設計された、エレガントで強力なHTTPクライアントライブラリです。標準ライブラリ `net/http` をベースに構築されており、流れるようなチェーンAPI、強力なミドルウェアシステム、そして一連の便利な機能を提供します。開発者にとって最も直感的な方法でHTTPリクエストを送信し、開発体験を向上させることを目指しています。

## 🧠 命名について

- 1. `requests` という名前のライブラリは既に多すぎます。
- 2. ラテン語系の言語が持つ活発さと情熱を思わせ、覚えやすく親しみやすい響きがあります。
- 3. また、`request` + `go` の組み合わせとしても解釈できます。

## ⭐ 主な特徴

*   📦️ 手軽な直接リクエストと、詳細な設定が可能なクライアントの両方を提供し、様々なシナリオに柔軟に対応します。
*   ✨ 直感的で流れるようなチェーンAPIにより、明確でスムーズなコーディングが可能です。
*   🚀 JSON、x-www-form-urlencoded、バイナリデータ、ファイルアップロードをサポートします。
*   🍪 `CookieJar` による自動的なセッション管理を行います。
*   ⏱️ `Context` を利用したタイムアウトとキャンセル制御が可能です。
*   🔧 設定可能なクライアント（タイムアウト、リダイレクトポリシーなど）。
*   🧅 強力かつシンプルに記述できるミドルウェア（フック）システムを備えています。

## 💿 インストール

```bash
go get github.com/Kaguya233qwq/requesto
```

## ⚡ クイックスタート

数行のコードでリクエストを送信し、JSONレスポンスを解析します。

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
		log.Fatalf("リクエスト失敗: %v", err)
	}
	fmt.Printf("ステータスコード: %d\n", resp.StatusCode())
	// レスポンスボディをマップに解析
	jsonData, _ := resp.Json()
	fmt.Printf("Args from server: %v\n", jsonData["args"])
}
```

## 📚 使い方

### 🚀 直接リクエスト

#### クエリパラメータの渡し方

```go
params := requesto.Params{"q": "1"}
resp, err = requesto.Get("https://example.com", params)
if err != nil {
	log.Fatalf("リクエスト失敗: %v", err)
}
fmt.Printf("ステータスコード: %d\n", resp.StatusCode())
```
#### リクエストヘッダーの設定

```go
params := requesto.Params{"q": "1"}
headers := requesto.Headers{"Accept": "text/html,application"}
resp, err = requesto.Get("https://example.com", params, headers)
if err != nil {
	log.Fatalf("リクエスト失敗: %v", err)
}
fmt.Printf("ステータスコード: %d\n", resp.StatusCode())
```

#### フォームデータの設定

```go
form := requesto.AsForm(map[string]string{"arg1": "value1"})
resp, err = requesto.Post("https://example.com", form)
if err != nil {
	log.Fatalf("リクエスト失敗: %v", err)
}
fmt.Printf("ステータスコード: %d\n", resp.StatusCode())
```

#### JSONデータの設定

```go
json := requesto.AsJson(map[string]any{"arg1": "value1", "arg2": 0})
resp, err = requesto.Post("https://example.com", json)
if err != nil {
	log.Fatalf("リクエスト失敗: %v", err)
}
fmt.Printf("ステータスコード: %d\n", resp.StatusCode())
```


### 🛠️ クライアントを利用したリクエスト

クライアント (`Client`) は再利用可能で、接続プール、CookieJar、グローバル設定を保持し、永続的なセッションのように機能します。

#### デフォルトクライアント

```go
client := requesto.NewClient("https://example.com")
```

#### カスタムクライアント

`NewClient` は関数型オプションをサポートしており、タイムアウトやリダイレクトポリシーなどを設定できます。

```go
client := requesto.NewClient(
    "https://example.com",
    requesto.WithTimeout(10*time.Second),      // グローバルタイムアウトを10秒に設定
    requesto.WithFollowRedirects(false),     // 自動リダイレクトを無効化
)
```

クライアントが保持する各種設定は、リクエスト送信前に直接変更できます。

```go
client := requesto.NewClient("https://example.com")
// リクエストヘッダーを設定
client.Headers.Set("Content-Type", "text/html; charset=utf-8")
// クエリパラメータを設定
client.Params = map[string]string{
    "page": "1",
    "size": "10",
}
// キー・バリュー形式でCookieを設定
client.SetCookiesFromMap(
    map[string]string{
        "__token": "xyz",
    },
)
client.Get()
```

#### 新しいリクエストの作成

`NewRequest()` を使用して新しいリクエストを構築し、より複雑なセッション制御を実現します。

```go
client := requesto.NewClient("https://api.example.com")
req1 := client.NewRequest()
req2 := client.NewRequest()
// 同じクライアントから作成されたリクエストは接続プールを共有しますが、互いに干渉しません
```

直接リクエストを送信することも可能です。

```go
client := requesto.NewClient("https://example.com")
resp, err := client.Get()
```

#### リクエストボディ

クライアントモードでは、`requesto` は様々な方法でリクエストボディを設定できます。

#### JSONの送信

構造体と `map[string]any` の両方をサポートしています。

```go
type User struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

user := User{Name: "Alice", Age: 20}
client := requesto.NewClient("https://example.com")
resp, err := client.NewRequest().
    JoinPath("/api"). // JoinPathメソッドで動的にパスを連結できます
    SetJsonData(user). // structまたはmapが利用可能です
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

#### フォームデータの送信

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

#### ファイルのアップロード

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

### レスポンスの処理

`Response` オブジェクトは、レスポンスボディを解析するための便利なメソッドを提供します。

```go
resp, _ := client.NewRequest().Get()

// ステータスコードとヘッダー情報の取得
statusCode := resp.StatusCode()
headers := resp.Headers()

// レスポンスボディの取得
text, _ := resp.Text()
bytes, _ := resp.Bytes()

// JSONの解析
jsonData, _ := resp.Json()

// レスポンス内のCookieを解析
cookies := resp.Cookies()
cookiesMap := resp.CookiesMap()
```

### Cookie管理

`Client` は `CookieJar` を内蔵しており、セッションを自動で処理します。

```go
client.NewRequest().JoinPath("/cookies/set").SetParams(map[string]string{
    "session_id": "my_secret_session",
}).Get()

// レスポンスヘッダーに `set-cookie` がある場合、対応するCookieが現在のクライアントに自動で設定されます。
resp, _ := client.NewRequest().JoinPath("/login").Get()
cookies := resp.Cookies()

// クライアントに手動でCookieを設定することも可能です。
client.SetCookiesFromMap(map[string]string{
    "token": "123",
})
```

### Contextによるタイムアウトとキャンセル

並行リクエストを行う際、`NewRequestWithContext` を使用してコンテキストを渡し、個別のリクエストにタイムアウトを設定したり、キャンセルシグナルを送信したりできます。

```go
// 5秒でタイムアウトする `context` を作成
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// このリクエストは `context` のタイムアウトにより失敗します
_, err := client.NewRequestWithContext(ctx).
    JoinPath("/delay/3").
    Get()

if errors.Is(err, context.DeadlineExceeded) {
    fmt.Println("リクエストがタイムアウトしました")
}
```

## 高度な機能：ミドルウェア

ミドルウェア（middleware）は `requesto` の強力な機能の一つで、リクエストの送信前後で独自のロジックを注入することができます。

### ミドルウェアの使用

`requesto` は、リトライやロギングなど、すぐに使えるミドルウェアをいくつか提供しています。

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

### カスタムミドルウェアの作成

簡単なミドルウェアを素早く構築するには `NewHook` ビルダーの使用を推奨します。

- #### `NewHook` の使用

下記は、各リクエストの前に認証ヘッダーを追加するミドルウェアの例です。

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

- #### 完全なインターフェースの実装

キャッシュのような、リクエストフローの制御が必要な複雑なロジックには、`requesto.Middleware` インターフェースを完全に実装します。

```go
func MyComplexMiddleware() requesto.Middleware {
    return func(req *requesto.Request, next requesto.Next) (*requesto.Response, error) {
        // リクエスト前のロジック...

        // 次のミドルウェアを呼び出す
        resp, err := next(req)

        // レスポンス後のロジック...

        return resp, err
    }
}
```

## 📜 ライセンス

このプロジェクトはMITライセンスの下で公開されています。

## 🙏 謝辞

以下のプロジェクトには、開発におけるインスピレーションやコードの参考にさせていただきました。皆様の存在がオープンソースコミュニティをより素晴らしいものにしています。心より感謝申し上げます。

*   [earthboundkid/requests](https://github.com/earthboundkid/requests)
*   [asmcos/requests](https://github.com/asmcos/requests)