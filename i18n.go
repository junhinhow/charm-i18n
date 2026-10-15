// Package i18n fornece internacionalizacao para interfaces de terminal.
//
// Suporta catalogos de mensagens, deteccao de locale, cadeia de fallback
// (ex: "pt-BR" -> "pt" -> "en"), formatacao com argumentos, carregamento de
// traducoes a partir de arquivos JSON e suporte a direcao de escrita da
// direita para a esquerda (RTL).
package i18n

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
)

// Catalog armazena traducoes organizadas por locale e chave.
// Estrutura: locale -> chave -> valor traduzido.
type Catalog struct {
	mu           sync.RWMutex
	translations map[string]map[string]string
	locale       string
	fallback     string
}

// NewCatalog cria um novo catalogo de traducoes com locale padrao "en".
func NewCatalog() *Catalog {
	return &Catalog{
		translations: make(map[string]map[string]string),
		locale:       defaultLocale,
		fallback:     defaultLocale,
	}
}

// AddLocale adiciona ou substitui as mensagens de um locale no catalogo.
//
// O locale e normalizado para a forma canonica (ex: "PT_br" -> "pt-BR"),
// de modo que a chave procurada depois por SetLocale/T sempre encontra as
// mensagens.
func (c *Catalog) AddLocale(locale string, messages map[string]string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	copied := make(map[string]string, len(messages))
	for k, v := range messages {
		copied[k] = v
	}
	c.translations[localeKey(locale)] = copied
}

// SetLocale define o locale ativo para traducoes.
//
// O locale e normalizado (ex: "PT_br.UTF-8" -> "pt-BR") e a busca passa a
// seguir a cadeia de fallback: locale ativo, idioma base do locale ativo,
// locale de fallback e idioma base do fallback.
func (c *Catalog) SetLocale(locale string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.locale = localeKey(locale)
}

// SetFallback define o locale de fallback quando a chave nao existe no locale ativo.
func (c *Catalog) SetFallback(locale string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.fallback = localeKey(locale)
}

// Locale retorna o locale ativo atual.
func (c *Catalog) Locale() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.locale
}

// T traduz uma chave usando o locale ativo.
//
// A busca percorre a cadeia de fallback (locale ativo -> idioma base ->
// fallback -> idioma base do fallback). Se a chave nao existir em nenhum
// deles, a propria chave e devolvida.
func (c *Catalog) T(key string) string {
	value, ok := c.lookup(key)
	if !ok {
		return key
	}
	return value
}

// Tf traduz uma chave e formata o resultado usando fmt.Sprintf.
//
// Quando a chave nao existe no catalogo, a PROPRIA CHAVE e devolvida sem
// aplicar a formatacao (evitando artefatos como
// "missing.key%!(EXTRA string=Gopher)").
//
// Quando o template traduzido nao pode ser formatado com os argumentos
// informados (faltam argumentos, o "%" literal nao foi escapado, etc.), o
// template original e devolvido sem modificacoes. Em nenhum dos casos a
// saida contem os artefatos de erro do fmt ("%!d(MISSING)", "%!(EXTRA ...)").
func (c *Catalog) Tf(key string, args ...any) string {
	template, ok := c.lookup(key)
	if !ok {
		// Chave inexistente: devolve a chave, sem formatar.
		return key
	}

	formatted, ok := formatTemplate(template, args)
	if !ok {
		// Template invalido ou argumentos faltando: devolve o template cru.
		return template
	}
	return formatted
}

// DetectLocale detecta o locale do sistema a partir das variaveis de ambiente.
// Verifica LC_ALL, LC_MESSAGES e LANG nessa ordem de prioridade.
// Retorna "en" se nenhuma variavel estiver definida ou se elas apontarem
// para um pseudo-locale do POSIX ("C", "C.UTF-8", "POSIX.UTF-8").
func DetectLocale() string {
	// Ordem de prioridade: LC_ALL > LC_MESSAGES > LANG
	for _, env := range []string{"LC_ALL", "LC_MESSAGES", "LANG"} {
		if val := os.Getenv(env); val != "" {
			return parseLocale(val)
		}
	}
	return defaultLocale
}

// LoadJSON carrega traducoes de um arquivo JSON.
// O formato esperado do JSON e:
//
//	{
//	  "en": { "hello": "Hello", "bye": "Goodbye" },
//	  "pt-BR": { "hello": "Ola", "bye": "Tchau" }
//	}
func (c *Catalog) LoadJSON(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("erro ao ler arquivo de traducoes: %w", err)
	}

	var raw map[string]map[string]string
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("erro ao decodificar JSON de traducoes: %w", err)
	}

	for locale, messages := range raw {
		c.AddLocale(locale, messages)
	}

	return nil
}

// Locales retorna a lista de locales disponiveis no catalogo, em ordem
// alfabetica (a ordem e estavel entre chamadas).
func (c *Catalog) Locales() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	locales := make([]string, 0, len(c.translations))
	for locale := range c.translations {
		locales = append(locales, locale)
	}
	sort.Strings(locales)
	return locales
}

// Has verifica se uma chave existe.
//
// Tem o mesmo comportamento de T: a chave e procurada em toda a cadeia de
// fallback, e nao apenas no locale ativo.
func (c *Catalog) Has(key string) bool {
	_, ok := c.lookup(key)
	return ok
}

// lookup procura a chave na cadeia de fallback do locale ativo.
// Devolve o valor e true quando a chave existe.
func (c *Catalog) lookup(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, locale := range c.chainLocked() {
		if msgs, ok := c.translations[locale]; ok {
			if value, ok := msgs[key]; ok {
				return value, true
			}
		}
	}
	return "", false
}

// chainLocked devolve a cadeia de busca dos locales: locale ativo, idioma
// base do locale ativo, locale de fallback e idioma base do fallback.
//
// Exemplo: locale "pt-BR" com fallback "en" -> ["pt-BR", "pt", "en"].
func (c *Catalog) chainLocked() []string {
	chain := make([]string, 0, 4)

	appendLocale := func(locale string) {
		if locale == "" {
			return
		}
		for _, existing := range chain {
			if existing == locale {
				return
			}
		}
		chain = append(chain, locale)
	}

	appendLocale(c.locale)
	appendLocale(baseLanguage(c.locale))
	appendLocale(c.fallback)
	appendLocale(baseLanguage(c.fallback))

	return chain
}

// localeKey normaliza um locale para uso como chave de armazenamento.
// Quando a normalizacao nao produz nada (locale vazio), o valor original e
// preservado.
func localeKey(locale string) string {
	if key := normalizeLocale(locale); key != "" {
		return key
	}
	return strings.TrimSpace(locale)
}
