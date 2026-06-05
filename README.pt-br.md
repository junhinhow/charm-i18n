# charm-i18n

Internacionalizacao para interfaces de terminal — catalogos de mensagens, deteccao de locale, suporte RTL.

> [Read in English](README.md)

## Funcionalidades

- **Catalogos de mensagens** com mapeamento locale -> chave -> valor
- **Deteccao de locale** a partir das variaveis `LC_ALL`, `LC_MESSAGES`, `LANG`
- **Traducoes formatadas** com argumentos no estilo `fmt.Sprintf`
- **Carregamento JSON** para arquivos de traducao
- **Locale de fallback** quando a chave nao existe no locale ativo
- **Thread-safe** em todas as operacoes do catalogo
- Zero dependencias externas (apenas stdlib)

## Instalacao

```bash
go get github.com/junhinhow/charm-i18n@latest
```

## Uso

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

    // Detectar locale automaticamente do ambiente
    cat.SetLocale(i18n.DetectLocale())

    // Traducao simples
    fmt.Println(cat.T("hello"))

    // Traducao formatada
    fmt.Println(cat.Tf("greeting", "Gopher", 42))
}
```

### Carregar de JSON

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

| Funcao / Metodo | Descricao |
|-----------------|-----------|
| `NewCatalog() *Catalog` | Cria um novo catalogo vazio (locale padrao: `"en"`) |
| `AddLocale(locale, messages)` | Adiciona traducoes para um locale |
| `SetLocale(locale)` | Define o locale ativo |
| `SetFallback(locale)` | Define o locale de fallback |
| `Locale() string` | Retorna o locale ativo atual |
| `T(key) string` | Traduz uma chave |
| `Tf(key, args...) string` | Traduz com formatacao `fmt.Sprintf` |
| `DetectLocale() string` | Detecta locale das variaveis de ambiente |
| `LoadJSON(path) error` | Carrega traducoes de um arquivo JSON |
| `Locales() []string` | Lista locales disponiveis |
| `Has(key) bool` | Verifica se uma chave existe no locale ativo |

## Licenca

[MIT](LICENSE)
