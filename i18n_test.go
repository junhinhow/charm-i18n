package i18n

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewCatalog(t *testing.T) {
	c := NewCatalog()
	if c.Locale() != "en" {
		t.Errorf("locale padrao esperado 'en', obteve '%s'", c.Locale())
	}
}

func TestAddLocaleETraducao(t *testing.T) {
	c := NewCatalog()
	c.AddLocale("pt-BR", map[string]string{
		"hello": "Ola",
		"bye":   "Tchau",
	})
	c.SetLocale("pt-BR")

	if got := c.T("hello"); got != "Ola" {
		t.Errorf("esperava 'Ola', obteve '%s'", got)
	}
	if got := c.T("bye"); got != "Tchau" {
		t.Errorf("esperava 'Tchau', obteve '%s'", got)
	}
}

func TestChaveInexistenteRetornaChave(t *testing.T) {
	c := NewCatalog()
	c.AddLocale("en", map[string]string{})
	if got := c.T("missing.key"); got != "missing.key" {
		t.Errorf("esperava 'missing.key', obteve '%s'", got)
	}
}

func TestFallback(t *testing.T) {
	c := NewCatalog()
	c.AddLocale("en", map[string]string{
		"hello": "Hello",
	})
	c.SetLocale("fr")

	// Deve usar fallback (en)
	if got := c.T("hello"); got != "Hello" {
		t.Errorf("esperava fallback 'Hello', obteve '%s'", got)
	}
}

func TestTfFormatacao(t *testing.T) {
	c := NewCatalog()
	c.AddLocale("en", map[string]string{
		"greeting": "Hello, %s! You have %d messages.",
	})

	got := c.Tf("greeting", "Alice", 5)
	expected := "Hello, Alice! You have 5 messages."
	if got != expected {
		t.Errorf("esperava '%s', obteve '%s'", expected, got)
	}
}

func TestTfSemArgumentos(t *testing.T) {
	c := NewCatalog()
	c.AddLocale("en", map[string]string{
		"hello": "Hello",
	})

	if got := c.Tf("hello"); got != "Hello" {
		t.Errorf("esperava 'Hello', obteve '%s'", got)
	}
}

func TestDetectLocaleComLANG(t *testing.T) {
	// Limpa variaveis que podem interferir
	os.Unsetenv("LC_ALL")
	os.Unsetenv("LC_MESSAGES")

	os.Setenv("LANG", "pt_BR.UTF-8")
	defer os.Unsetenv("LANG")

	got := DetectLocale()
	if got != "pt-BR" {
		t.Errorf("esperava 'pt-BR', obteve '%s'", got)
	}
}

func TestDetectLocaleComLCALL(t *testing.T) {
	os.Setenv("LC_ALL", "fr_FR.UTF-8")
	defer os.Unsetenv("LC_ALL")

	got := DetectLocale()
	if got != "fr-FR" {
		t.Errorf("esperava 'fr-FR', obteve '%s'", got)
	}
}

func TestDetectLocaleSemVariaveis(t *testing.T) {
	os.Unsetenv("LC_ALL")
	os.Unsetenv("LC_MESSAGES")
	os.Unsetenv("LANG")

	got := DetectLocale()
	if got != "en" {
		t.Errorf("esperava 'en', obteve '%s'", got)
	}
}

func TestParseLocale(t *testing.T) {
	testes := []struct {
		entrada  string
		esperado string
	}{
		{"pt_BR.UTF-8", "pt-BR"},
		{"en_US", "en-US"},
		{"C", "en"},
		{"POSIX", "en"},
		{"de_DE.UTF-8@euro", "de-DE"},
		{"", "en"},
	}

	for _, tt := range testes {
		t.Run(tt.entrada, func(t *testing.T) {
			got := parseLocale(tt.entrada)
			if got != tt.esperado {
				t.Errorf("parseLocale(%q): esperava %q, obteve %q", tt.entrada, tt.esperado, got)
			}
		})
	}
}

func TestLoadJSON(t *testing.T) {
	conteudo := `{
		"en": {"hello": "Hello", "bye": "Goodbye"},
		"pt-BR": {"hello": "Ola", "bye": "Tchau"}
	}`

	tmpDir := t.TempDir()
	caminho := filepath.Join(tmpDir, "translations.json")
	if err := os.WriteFile(caminho, []byte(conteudo), 0644); err != nil {
		t.Fatal(err)
	}

	c := NewCatalog()
	if err := c.LoadJSON(caminho); err != nil {
		t.Fatalf("LoadJSON falhou: %v", err)
	}

	c.SetLocale("pt-BR")
	if got := c.T("hello"); got != "Ola" {
		t.Errorf("esperava 'Ola', obteve '%s'", got)
	}

	c.SetLocale("en")
	if got := c.T("bye"); got != "Goodbye" {
		t.Errorf("esperava 'Goodbye', obteve '%s'", got)
	}
}

func TestLoadJSONArquivoInexistente(t *testing.T) {
	c := NewCatalog()
	err := c.LoadJSON("/caminho/inexistente.json")
	if err == nil {
		t.Error("esperava erro para arquivo inexistente")
	}
}

func TestLoadJSONInvalido(t *testing.T) {
	tmpDir := t.TempDir()
	caminho := filepath.Join(tmpDir, "invalid.json")
	if err := os.WriteFile(caminho, []byte("invalido"), 0644); err != nil {
		t.Fatal(err)
	}

	c := NewCatalog()
	err := c.LoadJSON(caminho)
	if err == nil {
		t.Error("esperava erro para JSON invalido")
	}
}

func TestLocales(t *testing.T) {
	c := NewCatalog()
	c.AddLocale("en", map[string]string{"a": "b"})
	c.AddLocale("pt-BR", map[string]string{"a": "b"})

	locales := c.Locales()
	if len(locales) != 2 {
		t.Errorf("esperava 2 locales, obteve %d", len(locales))
	}
}

func TestHas(t *testing.T) {
	c := NewCatalog()
	c.AddLocale("en", map[string]string{"hello": "Hello"})

	if !c.Has("hello") {
		t.Error("esperava Has('hello') = true")
	}
	if c.Has("missing") {
		t.Error("esperava Has('missing') = false")
	}
}
