
# requesto

[![Go Reference](https://pkg.go.dev/badge/github.com/Kaguya233qwq/requesto.svg)](https://pkg.go.dev/github.com/Kaguya233qwq/requesto)

[[English](../README.md)] | [[简体中文](/docs/README_zh-CN.md)] | [[繁體中文](/docs/README_zh-TW.md)] | [[日本語](/docs/README_ja.md)] | [[한국어](/docs/README_ko.md)] | [[Русский](/docs/README_ru.md)] | Español | [[Deutsch](/docs/README_de.md)]

> [!TIP]
> Este proyecto aún está en desarrollo activo. Si tienes buenas ideas o sugerencias, te invitamos a abrir un issue o enviar un pull request.

`requesto` es una biblioteca cliente HTTP para Go, diseñada para ser elegante y potente. Construida sobre la biblioteca estándar `net/http`, ofrece una API encadenada, un robusto sistema de middleware y una serie de funcionalidades convenientes, con el objetivo de realizar peticiones HTTP de la forma más intuitiva posible para mejorar la experiencia de desarrollo.

## 🧠 Sobre el nombre

- 1. Ya hay demasiadas librerías llamadas `requests`.
- 2. Tiene un aire vivo y apasionado, como las lenguas romances, lo que lo hace pegadizo y fácil de recordar.
- 3. También se puede entender como una combinación de `request` + `go`.

## ⭐ Características

*   📦️ Ofrece métodos de petición listos para usar y opciones avanzadas para satisfacer fácilmente diferentes escenarios.
*   ✨ Una API encadenada intuitiva y fluida que hace que las llamadas sean claras y ágiles.
*   🚀 Soporte para JSON, x-www-form-urlencoded, datos binarios y subida de archivos.
*   🍪 Gestión automática de sesiones con `CookieJar`.
*   ⏱️ Control de timeouts y cancelaciones a través de `Context`.
*   🔧 Cliente configurable (timeouts, políticas de redirección, etc.).
*   🧅 Un sistema de middleware (hooks) potente y fácil de escribir.

## 💿 Instalación

```bash
go get github.com/Kaguya233qwq/requesto
```

## ⚡ Inicio Rápido

Envía una petición y parsea una respuesta JSON en solo unas pocas líneas:

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
		log.Fatalf("Petición fallida: %v", err)
	}
	fmt.Printf("Código de estado: %d\n", resp.StatusCode())
	// Parsear el cuerpo de la respuesta a un mapa
	jsonData, _ := resp.Json()
	fmt.Printf("Args from server: %v\n", jsonData["args"])
}
```

## 📚 Uso

### 🚀 Peticiones Directas

#### Pasar parámetros de consulta (query)

```go
params := requesto.Params{"q": "1"}
resp, err = requesto.Get("https://example.com", params)
if err != nil {
	log.Fatalf("Petición fallida: %v", err)
}
fmt.Printf("Código de estado: %d\n", resp.StatusCode())
```
#### Establecer cabeceras (headers)

```go
params := requesto.Params{"q": "1"}
headers := requesto.Headers{"Accept": "text/html,application"}
resp, err = requesto.Get("https://example.com", params, headers)
if err != nil {
	log.Fatalf("Petición fallida: %v", err)
}
fmt.Printf("Código de estado: %d\n", resp.StatusCode())
```

#### Enviar datos de formulario (Form)

```go
form := requesto.AsForm(map[string]string{"arg1": "value1"})
resp, err = requesto.Post("https://example.com", form)
if err != nil {
	log.Fatalf("Petición fallida: %v", err)
}
fmt.Printf("Código de estado: %d\n", resp.StatusCode())
```

#### Enviar JSON

```go
json := requesto.AsJson(map[string]any{"arg1": "value1", "arg2": 0})
resp, err = requesto.Post("https://example.com", json)
if err != nil {
	log.Fatalf("Petición fallida: %v", err)
}
fmt.Printf("Código de estado: %d\n", resp.StatusCode())
```

### 🛠️ Peticiones con Cliente

El cliente (`Client`) es reutilizable y contiene un pool de conexiones, un `CookieJar` y configuraciones globales, funcionando como una sesión más persistente.

#### Cliente por defecto

```go
client := requesto.NewClient("https://example.com")
```

#### Cliente personalizado

`NewClient` soporta opciones funcionales para configurar timeouts, políticas de redirección, etc.

```go
client := requesto.NewClient(
    "https://example.com",
    requesto.WithTimeout(10*time.Second),      // Establecer un timeout global de 10 segundos
    requesto.WithFollowRedirects(false),     // Deshabilitar redirecciones automáticas
)
```

Muchos atributos del cliente pueden ser establecidos o modificados directamente antes de realizar una petición:

```go
client := requesto.NewClient("https://example.com")
// Establecer cabeceras
client.Headers.Set("Content-Type", "text/html; charset=utf-8")
// Establecer parámetros de consulta
client.Params = map[string]string{
    "page": "1",
    "size": "10",
}
// Establecer cookies desde un mapa
client.SetCookiesFromMap(
    map[string]string{
        "__token": "xyz",
    },
)
client.Get()
```

#### Crear una nueva petición

Usa `NewRequest()` para construir una nueva petición, lo que permite un control más complejo de la sesión:

```go
client := requesto.NewClient("https://api.example.com")
req1 := client.NewRequest()
req2 := client.NewRequest()
// Las peticiones del mismo cliente comparten el pool de conexiones, pero son independientes entre sí.
```

También puedes realizar una petición directamente:

```go
client := requesto.NewClient("https://example.com")
resp, err := client.Get()
```

#### Cuerpo de la Petición (Request Body)

En el modo cliente, `requesto` soporta varias formas de establecer el cuerpo de la petición.

#### Enviar JSON

Soporta tanto `structs` como `map[string]any`:

```go
type User struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

user := User{Name: "Alice", Age: 20}
client := requesto.NewClient("https://example.com")
resp, err := client.NewRequest().
    JoinPath("/api"). // Usa JoinPath para añadir segmentos a la URL dinámicamente
    SetJsonData(user). // Puede ser un struct o un mapa
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

#### Enviar datos de Formulario

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

#### Subir archivos

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

### Manejo de la Respuesta

El objeto `Response` ofrece métodos convenientes para parsear el cuerpo de la respuesta.

```go
resp, _ := client.NewRequest().Get()

// Obtener código de estado y cabeceras
statusCode := resp.StatusCode()
headers := resp.Headers()

// Obtener el cuerpo de la respuesta
text, _ := resp.Text()
bytes, _ := resp.Bytes()

// Parsear JSON
jsonData, _ := resp.Json()

// Parsear Cookies de la respuesta
cookies := resp.Cookies()
cookiesMap := resp.CookiesMap()
```

### Gestión de Cookies

El `Client` tiene un `CookieJar` integrado que maneja las sesiones automáticamente.

```go
client.NewRequest().JoinPath("/cookies/set").SetParams(map[string]string{
    "session_id": "my_secret_session",
}).Get()

// Al acceder a un endpoint que requiere autenticación, si la respuesta incluye la cabecera `set-cookie`, la cookie se guardará automáticamente en el cliente.
resp, _ := client.NewRequest().JoinPath("/login").Get()
cookies := resp.Cookies()

// También puedes establecer cookies manualmente en el cliente
client.SetCookiesFromMap(map[string]string{
    "token": "123",
})
```

### Usar Context para Timeouts y Cancelaciones

Cuando necesites realizar peticiones concurrentes, puedes usar `NewRequestWithContext` para pasar un contexto y establecer un timeout o una señal de cancelación para una petición específica.

```go
// Crear un contexto que expirará en 5 segundos
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// Esta petición fallará debido al timeout del contexto
_, err := client.NewRequestWithContext(ctx).
    JoinPath("/delay/3").
    Get()

if errors.Is(err, context.DeadlineExceeded) {
    fmt.Println("La petición ha expirado (timeout)")
}
```

## Funcionalidades Avanzadas: Middleware

El middleware es una de las funcionalidades avanzadas más potentes de `requesto`, permitiéndote inyectar lógica personalizada antes de que se envíe una petición o después de que se reciba una respuesta.

### Uso de Middleware

`requesto` proporciona algunos middlewares listos para usar, como reintentos (retry) y logging.

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

### Escribir Middleware Personalizado

Se recomienda usar el constructor `NewHook` para crear rápidamente middlewares sencillos.

- #### Usando `NewHook`

A continuación, un ejemplo de un middleware que añade una cabecera de autenticación a cada petición:

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

- #### Implementar la interfaz completa

Para lógica más compleja que necesita controlar el flujo de la petición (como el caché), puedes implementar la interfaz `requesto.Middleware` completa:

```go
func MyComplexMiddleware() requesto.Middleware {
    return func(req *requesto.Request, next requesto.Next) (*requesto.Response, error) {
        // Lógica antes de la petición...

        // Llamar al siguiente middleware en la cadena
        resp, err := next(req)

        // Lógica después de la respuesta...

        return resp, err
    }
}
```

## 📜 Licencia

Este proyecto se distribuye bajo la licencia MIT.

## 🙏 Agradecimientos

Agradecemos a los siguientes proyectos por servir de inspiración y referencia de código. La comunidad de código abierto es un lugar mejor gracias a vuestra existencia.

*   [earthboundkid/requests](https://github.com/earthboundkid/requests)
*   [asmcos/requests](https://github.com/asmcos/requests)