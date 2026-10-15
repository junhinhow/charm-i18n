# charm-i18n

Internacionalização para interfaces de terminal — catálogos de mensagens, detecção de locale, suporte a RTL.

> [Read in English](README.md)

## Funcionalidades

- **Catálogos de mensagens** com mapeamento locale -> chave -> valor
- **Detecção de locale** a partir das variáveis `LC_ALL`, `LC_MESSAGES`, `LANG`
- **Normalização de locale** (`PT_br.UTF-8` → `pt-BR`, `C.UTF-8` → `en`)
- **Cadeia de fallback** (`pt-BR` → `pt` → locale de fallback)
- **Traduções formatadas** com argumentos no estilo `fmt.Sprintf` (sem artefatos `%!(EXTRA ...)`)
- **Carregamento JSON** para arquivos de tradução
- **Suporte a RTL**: detecção de script e reordenação visual de texto bidirecional
- **Thread-safe** em todas as operações do catálogo
- Zero dependências externas (apenas stdlib)

## Instalação

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

    // Tradução simples
    fmt.Println(cat.T("hello"))

    // Tradução formatada
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

## Resolução de locale

`DetectLocale` lê `LC_ALL`, `LC_MESSAGES` e `LANG` (nessa ordem) e devolve um locale
normalizado:

| Valor no ambiente | Resultado |
|---|---|
| `pt_BR.UTF-8` | `pt-BR` |
| `PT_br.UTF-8` | `pt-BR` |
| `de_DE.UTF-8@euro` | `de-DE` |
| `C`, `c`, `POSIX`, `posix` | `en` |
| `C.UTF-8`, `POSIX.UTF-8` | `en` |

Os pseudo-locales do POSIX viram `en` porque `C.UTF-8` é o default de muitas distribuições
Linux modernas (systemd/Debian) e não descreve nenhum idioma real.

### Cadeia de fallback

A busca percorre a cadeia de locales até encontrar a chave:

```
locale ativo  ->  idioma base  ->  locale de fallback  ->  idioma base do fallback
    pt-BR      ->       pt       ->          en           ->            en
```

```go
cat.SetLocale("pt-BR")
cat.SetFallback("en")

cat.AddLocale("en",   map[string]string{"only.en": "Only EN"})
cat.AddLocale("pt",   map[string]string{"only.pt": "Only PT"})
cat.AddLocale("pt-BR", map[string]string{"only.br": "Only BR"})

cat.T("only.br") // "Only BR"  (locale exato)
cat.T("only.pt") // "Only PT"  (idioma base)
cat.T("only.en") // "Only EN"  (fallback)
```

O `Has` percorre exatamente a mesma cadeia do `T`, então os dois nunca se contradizem.

## Traduções formatadas

O `Tf` nunca deixa vazar os artefatos de erro do `fmt.Sprintf`:

| Chamada | Resultado |
|---|---|
| `Tf("greeting", "Alice", 5)` com `"Hello, %s! You have %d messages."` | `Hello, Alice! You have 5 messages.` |
| `Tf("missing.key", "Gopher")` | `missing.key` (a própria chave, sem formatação) |
| `Tf("greeting", "Alice")` com dois verbos | o template original (sem `%!d(MISSING)`) |
| `Tf("discount", "Gopher")` com `"50% off"` | `50% off` (um `%` solto é mantido literal) |
| `Tf("discount", "Gopher")` com `"50%% off"` | `50% off` |

Regras:

- Chave inexistente devolve a própria chave, sem nenhuma formatação.
- Template que não pode ser formatado com os argumentos informados (argumentos
  faltando, tipos incompatíveis) é devolvido sem alterações — nada é formatado pela metade.
- Argumentos extras são ignorados.
- `%%` produz um `%` literal; um `%` que não inicia uma diretiva válida permanece literal.

## Suporte a RTL (direita para esquerda)

O pacote tem suporte a RTL de verdade, e não apenas uma promessa:

```go
i18n.IsRTL("ar-EG")   // true
i18n.IsRTL("he")      // true
i18n.IsRTL("az-Arab") // true  (o script tem precedência sobre o idioma)
i18n.IsRTL("pt-BR")   // false

// Ordem visual de texto bidirecional (para terminais sem suporte a bidi)
i18n.Reverse("Hello שלום") // "Hello םולש"
i18n.Reverse("שלום")        // "םולש"
```

**Escopo (descrição honesta):**

- `IsRTL` decide a direção pelo subtag de idioma do locale (árabe, hebraico, persa,
  urdu, ...) ou, quando presente, pelo subtag de script (ISO 15924).
- `Reverse` reordena a string para a **ordem visual** (da esquerda para a direita) usando
  um subconjunto do algoritmo bidirecional do Unicode (UAX #9): segmentação em clusters
  de grapheme, resolução de neutros e números (regras W1–W7, N1–N2), níveis implícitos
  (I1–I2) e reordenação (L1–L2).
- **Não implementado**: controles de embedding explícitos (`LRE`/`RLE`/`LRO`/`RLO`),
  controles de isolamento e a base completa de caracteres do Unicode para as classes
  bidirecionais (a classificação cobre latim, grego, cirílico, árabe, hebraico e os
  demais scripts RTL comuns).
- Em texto puramente LTR a ordem visual é igual à ordem lógica, então `Reverse` devolve
  a entrada sem alterações.

## API

| Função / Método | Descrição |
|-----------------|-----------|
| `NewCatalog() *Catalog` | Cria um novo catálogo vazio (locale padrão: `"en"`) |
| `AddLocale(locale, messages)` | Adiciona traduções para um locale (a chave é normalizada) |
| `SetLocale(locale)` | Define o locale ativo (normalizado) |
| `SetFallback(locale)` | Define o locale de fallback (normalizado) |
| `Locale() string` | Retorna o locale ativo atual |
| `T(key) string` | Traduz uma chave, percorrendo a cadeia de fallback |
| `Tf(key, args...) string` | Traduz com formatação `fmt.Sprintf` (sem artefatos de erro) |
| `Has(key) bool` | Verifica se a chave é resolvida pela cadeia de fallback |
| `DetectLocale() string` | Detecta locale das variáveis de ambiente |
| `LoadJSON(path) error` | Carrega traduções de um arquivo JSON |
| `Locales() []string` | Lista os locales disponíveis (em ordem alfabética) |
| `IsRTL(locale) bool` | Verifica se o locale é da direita para a esquerda |
| `Reverse(s) string` | Reordena a string para a ordem visual (texto bidirecional) |

## Licença

[MIT](LICENSE)
