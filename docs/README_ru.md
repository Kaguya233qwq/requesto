# requesto

[![Go Reference](https://pkg.go.dev/badge/github.com/Kaguya233qwq/requesto.svg)](https://pkg.go.dev/github.com/Kaguya233qwq/requesto)

[[English](../README.md)] | [[简体中文](/docs/README_zh-CN.md)] | [[繁體中文](/docs/README_zh-TW.md)] | [[日本語](/docs/README_ja.md)] | [[한국어](/docs/README_ko.md)] | Русский | [[Español](/docs/README_es.md)] | [[Deutsch](/docs/README_de.md)]

> [!TIP]
> Этот проект находится в активной разработке. Если у вас есть хорошие идеи или предложения, мы будем рады, если вы создадите Issue или Pull Request.

`requesto` — это элегантная и мощная HTTP-клиентская библиотека для языка Go. Она построена на основе стандартной библиотеки `net/http` и предоставляет API с цепочечными вызовами, мощную систему промежуточного ПО (middleware) и ряд удобных функций. Цель `requesto` — сделать отправку HTTP-запросов максимально интуитивной для разработчика и улучшить опыт разработки.

## 🧠 О названии

- 1. Существует слишком много библиотек с названием `requests`.
- 2. Название звучит живо и страстно, в духе романских языков, что делает его приятным и запоминающимся.
- 3. Его также можно интерпретировать как сочетание `request` + `go`.

## ⭐ Основные возможности

*   📦️ Готовые к использованию простые и профессиональные способы отправки запросов, подходящие для разных сценариев.
*   ✨ Интуитивно понятный API с цепочечными вызовами, обеспечивающий ясность и плавность кода.
*   🚀 Поддержка JSON, x-www-form-urlencoded, бинарных данных и загрузки файлов.
*   🍪 Автоматическое управление сессиями с помощью `CookieJar`.
*   ⏱️ Управление таймаутами и отменой через `Context`.
*   🔧 Настраиваемый клиент (таймауты, политика редиректов и т.д.).
*   🧅 Мощная и простая в написании система промежуточного ПО (хуков).

## 💿 Установка

```bash
go get github.com/Kaguya233qwq/requesto
```

## ⚡ Быстрый старт

Отправьте запрос и разберите JSON-ответ всего в нескольких строках кода:

```go
package main

import (
	"fmt"
	"log"

	"github.comcom/Kaguya233qwq/requesto"
)

func main() {
	resp, err := requesto.Get("https://httpbingo.org/get")
	if err != nil {
		log.Fatalf("Ошибка запроса: %v", err)
	}
	fmt.Printf("Код состояния: %d\n", resp.StatusCode())
	// Парсинг тела ответа в map
	jsonData, _ := resp.Json()
	fmt.Printf("Args from server: %v\n", jsonData["args"])
}
```

## 📚 Использование

### 🚀 Прямые запросы

#### Передача параметров запроса (Query)

```go
params := requesto.Params{"q": "1"}
resp, err = requesto.Get("https://example.com", params)
if err != nil {
	log.Fatalf("Ошибка запроса: %v", err)
}
fmt.Printf("Код состояния: %d\n", resp.StatusCode())
```
#### Установка заголовков (Headers)

```go
params := requesto.Params{"q": "1"}
headers := requesto.Headers{"Accept": "text/html,application"}
resp, err = requesto.Get("https://example.com", params, headers)
if err != nil {
	log.Fatalf("Ошибка запроса: %v", err)
}
fmt.Printf("Код состояния: %d\n", resp.StatusCode())
```

#### Установка данных формы (Form)

```go
form := requesto.AsForm(map[string]string{"arg1": "value1"})
resp, err = requesto.Post("https://example.com", form)
if err != nil {
	log.Fatalf("Ошибка запроса: %v", err)
}
fmt.Printf("Код состояния: %d\n", resp.StatusCode())
```

#### Установка JSON

```go
json := requesto.AsJson(map[string]any{"arg1": "value1", "arg2": 0})
resp, err = requesto.Post("https://example.com", json)
if err != nil {
	log.Fatalf("Ошибка запроса: %v", err)
}
fmt.Printf("Код состояния: %d\n", resp.StatusCode())
```


### 🛠️ Использование клиента

Клиент (`Client`) является многоразовым. Он содержит пул соединений, `CookieJar` и глобальные настройки, что делает его похожим на постоянную сессию.

#### Клиент по умолчанию

```go
client := requesto.NewClient("https://example.com")
```

#### Пользовательский клиент

`NewClient` поддерживает функциональные опции для настройки таймаутов, политики редиректов и т.д.

```go
client := requesto.NewClient(
    "https://example.com",
    requesto.WithTimeout(10*time.Second),      // Установить глобальный таймаут 10 секунд
    requesto.WithFollowRedirects(false),     // Запретить автоматические редиректы
)
```

Многие атрибуты клиента можно настраивать и изменять непосредственно перед отправкой запроса:

```go
client := requesto.NewClient("https://example.com")
// Установка заголовков
client.Headers.Set("Content-Type", "text/html; charset=utf-8")
// Установка параметров запроса
client.Params = map[string]string{
    "page": "1",
    "size": "10",
}
// Установка cookie в формате ключ-значение
client.SetCookiesFromMap(
    map[string]string{
        "__token": "xyz",
    },
)
client.Get()
```

#### Создание нового запроса

Используйте `NewRequest()` для создания нового запроса и реализации более сложного управления сессией:

```go
client := requesto.NewClient("https://api.example.com")
req1 := client.NewRequest()
req2 := client.NewRequest()
// Запросы от одного клиента используют общий пул соединений, но не мешают друг другу
```

Также можно отправлять запросы напрямую:

```go
client := requesto.NewClient("https://example.com")
resp, err := client.Get()
```

#### Тело запроса

В режиме клиента `requesto` поддерживает несколько способов установки тела запроса.

#### Отправка JSON

Поддерживаются структуры и `map[string]any`:

```go
type User struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

user := User{Name: "Alice", Age: 20}
client := requesto.NewClient("https://example.com")
resp, err := client.NewRequest().
    JoinPath("/api"). // Метод JoinPath для динамического добавления пути
    SetJsonData(user). // Может быть структурой или map
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

#### Отправка данных формы

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

#### Загрузка файлов

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

### Обработка ответа

Объект `Response` предоставляет удобные методы для парсинга тела ответа.

```go
resp, _ := client.NewRequest().Get()

// Получение кода состояния и заголовков ответа
statusCode := resp.StatusCode()
headers := resp.Headers()

// Получение тела ответа
text, _ := resp.Text()
bytes, _ := resp.Bytes()

// Парсинг JSON
jsonData, _ := resp.Json()

// Парсинг Cookie из ответа
cookies := resp.Cookies()
cookiesMap := resp.CookiesMap()
```

### Управление Cookie

В `Client` встроен `CookieJar`, который автоматически управляет сессиями.

```go
client.NewRequest().JoinPath("/cookies/set").SetParams(map[string]string{
    "session_id": "my_secret_session",
}).Get()

// При доступе к защищенным ресурсам cookie будут установлены автоматически, если сервер отправит заголовок set-cookie
resp, _ := client.NewRequest().JoinPath("/login").Get()
cookies := resp.Cookies()

// Также можно установить Cookie для клиента вручную
client.SetCookiesFromMap(map[string]string{
    "token": "123",
})
```

### Использование Context для таймаутов и отмены

При выполнении параллельных запросов вы можете использовать `NewRequestWithContext`, чтобы передать контекст для установки таймаута или сигнала отмены для одного запроса.

```go
// Создание контекста, который отменится через 5 секунд
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// Этот запрос завершится неудачно из-за таймаута контекста
_, err := client.NewRequestWithContext(ctx).
    JoinPath("/delay/3").
    Get()

if errors.Is(err, context.DeadlineExceeded) {
    fmt.Println("Таймаут запроса")
}
```

## Продвинутые возможности: Промежуточное ПО (Middleware)

Промежуточное ПО (middleware) — одна из мощных продвинутых функций `requesto`, позволяющая встраивать собственную логику до отправки запроса и после получения ответа.

### Использование промежуточного ПО

`requesto` предоставляет готовые middleware, такие как повторные попытки (retry) и логирование.

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

### Написание собственного промежуточного ПО

Для создания простого middleware рекомендуется использовать конструктор `NewHook`.

- #### Использование `NewHook`

Пример middleware, которое добавляет заголовок аутентификации перед каждым запросом:

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

- #### Реализация полного интерфейса

Для более сложной логики, требующей управления потоком запросов (например, кеширование), вы можете реализовать полный интерфейс `requesto.Middleware`:

```go
func MyComplexMiddleware() requesto.Middleware {
    return func(req *requesto.Request, next requesto.Next) (*requesto.Response, error) {
        // Логика перед запросом...

        // Вызов следующего middleware в цепочке
        resp, err := next(req)

        // Логика после ответа...

        return resp, err
    }
}
```

## 📜 Лицензия

Этот проект распространяется под лицензией MIT.

## 🙏 Благодарности

Мы благодарим следующие проекты за вдохновение и примеры кода. Благодаря вашему существованию сообщество открытого исходного кода становится лучше.

*   [earthboundkid/requests](https://github.com/earthboundkid/requests)
*   [asmcos/requests](https://github.com/asmcos/requests)