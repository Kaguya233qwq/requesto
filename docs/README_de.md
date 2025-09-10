# requesto

[![Go Reference](https://pkg.go.dev/badge/github.com/Kaguya233qwq/requesto.svg)](https://pkg.go.dev/github.com/Kaguya233qwq/requesto)

[[English](../README.md)] | [[简体中文](/docs/README_zh-CN.md)] | [[繁體中文](/docs/README_zh-TW.md)] | [[日本語](/docs/README_ja.md)] | [[한국어](/docs/README_ko.md)] | [[Русский](/docs/README_ru.md)] | [[Español](/docs/README_es.md)] | Deutsch

> [!TIP]
> Dieses Projekt wird aktiv weiterentwickelt. Wenn du gute Ideen oder Vorschläge hast, freuen wir uns über ein neues Issue oder einen Pull Request.

`requesto` ist eine elegante und leistungsstarke HTTP-Client-Bibliothek für Go. Sie basiert auf der Standardbibliothek `net/http` und bietet eine verkettbare API (Fluent API), ein mächtiges Middleware-System sowie eine Reihe praktischer Funktionen. Das Ziel ist es, HTTP-Anfragen auf die für Entwickler intuitivste Weise zu ermöglichen und die Entwicklungserfahrung zu verbessern.

## 🧠 Über den Namen

- 1. Es gibt einfach zu viele Bibliotheken, die `requests` heißen.
- 2. Er hat eine lebhafte und leidenschaftliche Ausstrahlung, ähnlich wie romanische Sprachen, was ihn eingängig und leicht zu merken macht.
- 3. Man kann ihn auch als eine Kombination aus `request` + `go` verstehen.

## ⭐ Features

*   📦️ Bietet sofort einsatzbereite, einfache Request-Methoden sowie erweiterte Optionen für komplexe Szenarien.
*   ✨ Eine äußerst intuitive, verkettbare API für klare und flüssige Aufrufe.
*   🚀 Unterstützung für JSON, x-www-form-urlencoded, Binärdaten und Datei-Uploads.
*   🍪 Automatische Sitzungsverwaltung durch `CookieJar`.
*   ⏱️ Steuerung von Timeouts und Abbrüchen über `Context`.
*   🔧 Konfigurierbarer Client (Timeouts, Weiterleitungs-Strategien etc.).
*   🧅 Ein leistungsstarkes und einfach zu schreibendes Middleware-System (Hooks).

## 💿 Installation

```bash
go get github.com/Kaguya233qwq/requesto
```

## ⚡ Schnellstart

Senden Sie eine Anfrage und parsen Sie eine JSON-Antwort in nur wenigen Zeilen:

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
		log.Fatalf("Anfrage fehlgeschlagen: %v", err)
	}
	fmt.Printf("Statuscode: %d\n", resp.StatusCode())
	// Den Antwort-Body in eine Map parsen
	jsonData, _ := resp.Json()
	fmt.Printf("Args from server: %v\n", jsonData["args"])
}
```

## 📚 Verwendung

### 🚀 Direkte Anfragen

#### Query-Parameter übergeben

```go
params := requesto.Params{"q": "1"}
resp, err = requesto.Get("https://example.com", params)
if err != nil {
	log.Fatalf("Anfrage fehlgeschlagen: %v", err)
}
fmt.Printf("Statuscode: %d\n", resp.StatusCode())
```
#### Request-Header setzen

```go
params := requesto.Params{"q": "1"}
headers := requesto.Headers{"Accept": "text/html,application"}
resp, err = requesto.Get("https://example.com", params, headers)
if err != nil {
	log.Fatalf("Anfrage fehlgeschlagen: %v", err)
}
fmt.Printf("Statuscode: %d\n", resp.StatusCode())
```

#### Formulardaten (Form) senden

```go
form := requesto.AsForm(map[string]string{"arg1": "value1"})
resp, err = requesto.Post("https://example.com", form)
if err != nil {
	log.Fatalf("Anfrage fehlgeschlagen: %v", err)
}
fmt.Printf("Statuscode: %d\n", resp.StatusCode())
```

#### JSON senden

```go
json := requesto.AsJson(map[string]any{"arg1": "value1", "arg2": 0})
resp, err = requesto.Post("https://example.com", json)
if err != nil {
	log.Fatalf("Anfrage fehlgeschlagen: %v", err)
}
fmt.Printf("Statuscode: %d\n", resp.StatusCode())
```


### 🛠️ Anfragen mit einem Client

Der Client (`Client`) ist wiederverwendbar. Er enthält einen Connection Pool, einen `CookieJar` und globale Konfigurationen, was ihn zu einer Art persistenten Sitzung macht.

#### Standard-Client

```go
client := requesto.NewClient("https://example.com")
```

#### Benutzerdefinierter Client

`NewClient` unterstützt funktionale Optionen zur Konfiguration von Timeouts, Weiterleitungs-Strategien und mehr.

```go
client := requesto.NewClient(
    "https://example.com",
    requesto.WithTimeout(10*time.Second),      // Globales Timeout auf 10 Sekunden setzen
    requesto.WithFollowRedirects(false),     // Automatische Weiterleitungen deaktivieren
)
```

Viele Eigenschaften des Clients können vor dem Senden einer Anfrage direkt gesetzt oder geändert werden:

```go
client := requesto.NewClient("https://example.com")
// Header setzen
client.Headers.Set("Content-Type", "text/html; charset=utf-8")
// Query-Parameter setzen
client.Params = map[string]string{
    "page": "1",
    "size": "10",
}
// Cookies aus einer Map setzen
client.SetCookiesFromMap(
    map[string]string{
        "__token": "xyz",
    },
)
client.Get()
```

#### Eine neue Anfrage erstellen

Verwenden Sie `NewRequest()`, um eine neue Anfrage zu erstellen, was eine komplexere Sitzungssteuerung ermöglicht:

```go
client := requesto.NewClient("https://api.example.com")
req1 := client.NewRequest()
req2 := client.NewRequest()
// Anfragen vom selben Client teilen sich einen Connection Pool, beeinflussen sich aber ansonsten nicht.
```

Sie können Anfragen auch direkt senden:

```go
client := requesto.NewClient("https://example.com")
resp, err := client.Get()
```

#### Request Body

Im Client-Modus unterstützt `requesto` verschiedene Möglichkeiten, den Request Body zu setzen.

#### JSON senden

Es werden sowohl Structs als auch `map[string]any` unterstützt:

```go
type User struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

user := User{Name: "Alice", Age: 20}
client := requesto.NewClient("https://example.com")
resp, err := client.NewRequest().
    JoinPath("/api"). // JoinPath verwenden, um den Pfad dynamisch zusammenzufügen
    SetJsonData(user). // Kann ein Struct oder eine Map sein
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

#### Formulardaten senden

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

#### Dateien hochladen

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

### Antwort verarbeiten

Das `Response`-Objekt bietet praktische Methoden zum Parsen des Antwort-Bodys.

```go
resp, _ := client.NewRequest().Get()

// Statuscode und Header abrufen
statusCode := resp.StatusCode()
headers := resp.Headers()

// Antwort-Body abrufen
text, _ := resp.Text()
bytes, _ := resp.Bytes()

// JSON parsen
jsonData, _ := resp.Json()

// Cookies aus der Antwort parsen
cookies := resp.Cookies()
cookiesMap := resp.CookiesMap()
```

### Cookie-Verwaltung

Der `Client` verfügt über einen integrierten `CookieJar`, der Sitzungen automatisch verwaltet.

```go
client.NewRequest().JoinPath("/cookies/set").SetParams(map[string]string{
    "session_id": "my_secret_session",
}).Get()

// Wenn der Server einen `set-cookie`-Header sendet, wird das Cookie automatisch im Client gespeichert.
resp, _ := client.NewRequest().JoinPath("/login").Get()
cookies := resp.Cookies()

// Sie können Cookies auch manuell für den Client setzen
client.SetCookiesFromMap(map[string]string{
    "token": "123",
})
```

### Context für Timeouts und Abbrüche verwenden

Wenn Sie nebenläufige Anfragen durchführen, können Sie `NewRequestWithContext` verwenden, um einen Context zu übergeben und so ein Timeout oder ein Abbruchsignal für eine einzelne Anfrage zu setzen.

```go
// Einen Context erstellen, der nach 5 Sekunden abläuft
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// Diese Anfrage wird wegen des Context-Timeouts fehlschlagen
_, err := client.NewRequestWithContext(ctx).
    JoinPath("/delay/3").
    Get()

if errors.Is(err, context.DeadlineExceeded) {
    fmt.Println("Timeout bei der Anfrage")
}
```

## Erweiterte Funktionen: Middleware

Middleware ist eine der leistungsstärksten erweiterten Funktionen von `requesto`. Sie ermöglicht es Ihnen, benutzerdefinierte Logik vor dem Senden einer Anfrage oder nach dem Empfangen einer Antwort einzufügen.

### Middleware verwenden

`requesto` bietet einige sofort einsatzbereite Middlewares, wie z.B. für Wiederholungsversuche (Retry) und Logging.

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

### Benutzerdefinierte Middleware schreiben

Es wird empfohlen, den `NewHook`-Builder zu verwenden, um einfache Middlewares schnell zu erstellen.

- #### Verwendung von `NewHook`

Hier ist ein Beispiel für ein Middleware, das jeder Anfrage einen Authentifizierungs-Header hinzufügt:

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

- #### Vollständiges Interface implementieren

Für komplexere Logik, die den Anfragefluss steuern muss (z. B. Caching), können Sie das vollständige `requesto.Middleware`-Interface implementieren:

```go
func MyComplexMiddleware() requesto.Middleware {
    return func(req *requesto.Request, next requesto.Next) (*requesto.Response, error) {
        // Logik vor der Anfrage...

        // Das nächste Middleware in der Kette aufrufen
        resp, err := next(req)

        // Logik nach der Antwort...

        return resp, err
    }
}
```

## 📜 Lizenz

Dieses Projekt wird unter der MIT-Lizenz vertrieben.

## 🙏 Danksagung

Wir danken den folgenden Projekten für die Inspiration und als Code-Referenz. Die Open-Source-Community ist dank eurer Existenz ein besserer Ort.

*   [earthboundkid/requests](https://github.com/earthboundkid/requests)
*   [asmcos/requests](https://github.com/asmcos/requests)