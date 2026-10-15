package i18n

import "testing"

func TestNormalizeLocale(t *testing.T) {
	testes := []struct {
		entrada  string
		esperado string
	}{
		{"pt_BR", "pt-BR"},
		{"PT_br", "pt-BR"},
		{"PT-BR", "pt-BR"},
		{"pt", "pt"},
		{"EN", "en"},
		{"zh_Hant_TW", "zh-Hant-TW"},
		{"zh-hant-tw", "zh-Hant-TW"},
		{"es_419", "es-419"},
		{"pt_BR.UTF-8", "pt-BR"},
		{"de_DE.UTF-8@euro", "de-DE"},
		{"  pt-BR  ", "pt-BR"},
		{"", ""},
	}

	for _, tt := range testes {
		t.Run(tt.entrada, func(t *testing.T) {
			if got := normalizeLocale(tt.entrada); got != tt.esperado {
				t.Errorf("normalizeLocale(%q): esperava %q, obteve %q", tt.entrada, tt.esperado, got)
			}
		})
	}
}

func TestBaseLanguage(t *testing.T) {
	testes := []struct {
		entrada  string
		esperado string
	}{
		{"pt-BR", "pt"},
		{"pt_BR", "pt"},
		{"pt", "pt"},
		{"zh-Hant-TW", "zh"},
		{"", ""},
	}

	for _, tt := range testes {
		t.Run(tt.entrada, func(t *testing.T) {
			if got := baseLanguage(tt.entrada); got != tt.esperado {
				t.Errorf("baseLanguage(%q): esperava %q, obteve %q", tt.entrada, tt.esperado, got)
			}
		})
	}
}

func TestLocaleKey(t *testing.T) {
	if got := localeKey("PT_br"); got != "pt-BR" {
		t.Errorf("esperava 'pt-BR', obteve '%s'", got)
	}
	if got := localeKey("  "); got != "" {
		t.Errorf("esperava string vazia, obteve '%s'", got)
	}
}

func TestSetLocaleNormaliza(t *testing.T) {
	c := NewCatalog()
	c.SetLocale("PT_br.UTF-8")
	if got := c.Locale(); got != "pt-BR" {
		t.Errorf("esperava 'pt-BR', obteve '%s'", got)
	}
}

func TestSetFallbackNormaliza(t *testing.T) {
	c := NewCatalog()
	c.AddLocale("EN", map[string]string{"hello": "Hello"})
	c.SetLocale("pt-BR")
	c.SetFallback("en")

	if got := c.T("hello"); got != "Hello" {
		t.Errorf("esperava 'Hello', obteve '%s'", got)
	}
}
