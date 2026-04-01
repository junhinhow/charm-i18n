# charm-i18n

Internationalization for terminal UIs — message catalogs, locale detection, RTL support.

> [Leia em Portugues](README.pt-br.md)

## Features

- **Message catalogs** with locale -> key -> value mapping
- **Locale detection** from `LC_ALL`, `LC_MESSAGES`, `LANG` environment variables
- **Formatted translations** with `fmt.Sprintf`-style arguments
- **JSON loading** for translation files
- **Fallback locale** when a key is missing in the active locale
- **Thread-safe** catalog operations
- Zero external dependencies (stdlib only)

## Install

```bash
go get github.com/junhinhow/charm-i18n@latest
```

## Usage

```go
package main

import (
    "fmt"

    i18n "github.com/junhinhow/charm-i18n"
)

func main() {
    cat := i18n.NewCatalog()

    cat.AddLocale("en", map[string]string{
        "hello":    "Hello, World!",
        "greeting": "Hello, %s! You have %d messages.",
    })

    cat.AddLocale("pt-BR", map[string]string{
        "hello":    "Ola, Mundo!",
        "greeting": "Ola, %s! Voce tem %d mensagens.",
    })

    // Auto-detect locale from environment
    cat.SetLocale(i18n.DetectLocale())

    // Simple translation
    fmt.Println(cat.T("hello"))

    // Formatted translation
    fmt.Println(cat.Tf("greeting", "Gopher", 42))
}
```

### Load from JSON

```go
// translations.json:
// {
//   "en": { "hello": "Hello" },
//   "pt-BR": { "hello": "Ola" }
// }

cat := i18n.NewCatalog()
if err := cat.LoadJSON("translations.json"); err != nil {
    log.Fatal(err)
}
```

## API

| Function / Method | Description |
|-------------------|-------------|
| `NewCatalog() *Catalog` | Create a new empty catalog (default locale: `"en"`) |
| `AddLocale(locale, messages)` | Add translations for a locale |
| `SetLocale(locale)` | Set the active locale |
| `SetFallback(locale)` | Set the fallback locale |
| `Locale() string` | Get the current active locale |
| `T(key) string` | Translate a key |
| `Tf(key, args...) string` | Translate with `fmt.Sprintf` formatting |
| `DetectLocale() string` | Detect locale from env vars |
| `LoadJSON(path) error` | Load translations from a JSON file |
| `Locales() []string` | List available locales |
| `Has(key) bool` | Check if a key exists in the active locale |

## License

[MIT](LICENSE)
