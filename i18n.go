// Package i18n fornece internacionalizacao para interfaces de terminal.
//
// Suporta catalogos de mensagens, deteccao de locale, formatacao com argumentos
// e carregamento de traducoes a partir de arquivos JSON.
package i18n

import (
	"encoding/json"
	"fmt"
	"os"
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
		locale:       "en",
		fallback:     "en",
	}
}

// AddLocale adiciona ou substitui as mensagens de um locale no catalogo.
func (c *Catalog) AddLocale(locale string, messages map[string]string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	copied := make(map[string]string, len(messages))
	for k, v := range messages {
		copied[k] = v
	}
	c.translations[locale] = copied
}

// SetLocale define o locale ativo para traducoes.
func (c *Catalog) SetLocale(locale string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.locale = locale
}

// SetFallback define o locale de fallback quando a chave nao existe no locale ativo.
func (c *Catalog) SetFallback(locale string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.fallback = locale
}

// Locale retorna o locale ativo atual.
func (c *Catalog) Locale() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.locale
}

// T traduz uma chave usando o locale ativo.
// Se a chave nao for encontrada no locale ativo, tenta o fallback.
// Se tambem nao encontrar no fallback, retorna a propria chave.
func (c *Catalog) T(key string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Tenta locale ativo
	if msgs, ok := c.translations[c.locale]; ok {
		if val, ok := msgs[key]; ok {
			return val
		}
	}

	// Tenta fallback
	if c.locale != c.fallback {
		if msgs, ok := c.translations[c.fallback]; ok {
			if val, ok := msgs[key]; ok {
				return val
			}
		}
	}

	// Retorna a chave como ultimo recurso
	return key
}

// Tf traduz uma chave e formata o resultado usando fmt.Sprintf.
// Os argumentos sao aplicados ao template traduzido.
func (c *Catalog) Tf(key string, args ...any) string {
	template := c.T(key)
	if len(args) == 0 {
		return template
	}
	return fmt.Sprintf(template, args...)
}

// DetectLocale detecta o locale do sistema a partir das variaveis de ambiente.
// Verifica LC_ALL, LC_MESSAGES e LANG nessa ordem de prioridade.
// Retorna "en" se nenhuma variavel estiver definida.
func DetectLocale() string {
	// Ordem de prioridade: LC_ALL > LC_MESSAGES > LANG
	for _, env := range []string{"LC_ALL", "LC_MESSAGES", "LANG"} {
		if val := os.Getenv(env); val != "" {
			return parseLocale(val)
		}
	}
	return "en"
}

// parseLocale extrai o codigo do locale de um valor de variavel de ambiente.
// Exemplos: "pt_BR.UTF-8" -> "pt-BR", "en_US" -> "en-US", "C" -> "en"
func parseLocale(raw string) string {
	if raw == "C" || raw == "POSIX" {
		return "en"
	}

	// Remove encoding (ex: .UTF-8)
	if idx := strings.IndexByte(raw, '.'); idx >= 0 {
		raw = raw[:idx]
	}

	// Remove modifier (ex: @euro)
	if idx := strings.IndexByte(raw, '@'); idx >= 0 {
		raw = raw[:idx]
	}

	// Converte underscore para hifen (pt_BR -> pt-BR)
	raw = strings.ReplaceAll(raw, "_", "-")

	if raw == "" {
		return "en"
	}

	return raw
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

// Locales retorna a lista de locales disponiveis no catalogo.
func (c *Catalog) Locales() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	locales := make([]string, 0, len(c.translations))
	for locale := range c.translations {
		locales = append(locales, locale)
	}
	return locales
}

// Has verifica se uma chave existe no locale ativo.
func (c *Catalog) Has(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if msgs, ok := c.translations[c.locale]; ok {
		_, exists := msgs[key]
		return exists
	}
	return false
}
