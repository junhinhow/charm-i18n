package i18n

import "testing"

func TestFormatTemplate(t *testing.T) {
	testes := []struct {
		nome      string
		template  string
		args      []any
		esperado  string
		deveFalha bool
	}{
		{
			nome:     "simples",
			template: "Hello, %s!",
			args:     []any{"Gopher"},
			esperado: "Hello, Gopher!",
		},
		{
			nome:     "multiplos verbos",
			template: "%s tem %d mensagens",
			args:     []any{"Ana", 3},
			esperado: "Ana tem 3 mensagens",
		},
		{
			nome:     "percentual escapado",
			template: "50%% off",
			args:     []any{"Gopher"},
			esperado: "50% off",
		},
		{
			nome:     "percentual literal",
			template: "50% off",
			args:     []any{"Gopher"},
			esperado: "50% off",
		},
		{
			nome:      "argumento faltando",
			template:  "%s tem %d mensagens",
			args:      []any{"Ana"},
			deveFalha: true,
		},
		{
			nome:     "argumentos extras sao ignorados",
			template: "%s",
			args:     []any{"Ana", 3, true},
			esperado: "Ana",
		},
		{
			nome:     "largura e precisao",
			template: "[%10.3f]",
			args:     []any{3.14159},
			esperado: "[     3.142]",
		},
		{
			nome:     "indice explicito",
			template: "%[2]s e %[1]s",
			args:     []any{"Ana", "Bruno"},
			esperado: "Bruno e Ana",
		},
		{
			nome:     "verbo invalido vira literal",
			template: "%z",
			args:     []any{"Ana"},
			esperado: "%z",
		},
		{
			nome:     "percentual no fim vira literal",
			template: "100%",
			args:     []any{"Ana"},
			esperado: "100%",
		},
		{
			nome:      "tipo incompativel",
			template:  "%d",
			args:      []any{"texto"},
			deveFalha: true,
		},
	}

	for _, tt := range testes {
		t.Run(tt.nome, func(t *testing.T) {
			got, ok := formatTemplate(tt.template, tt.args)
			if tt.deveFalha {
				if ok {
					t.Errorf("esperava falha, obteve %q", got)
				}
				return
			}
			if !ok {
				t.Fatalf("esperava sucesso, obtive falha")
			}
			if got != tt.esperado {
				t.Errorf("esperava %q, obteve %q", tt.esperado, got)
			}
		})
	}
}

func TestFormatTemplateNuncaEmiteArtefatos(t *testing.T) {
	templates := []string{
		"%s e %d",
		"50% off",
		"%z",
		"%d",
		"%w",
		"%",
		"%%",
		"%%s",
		"%[3]s",
	}

	for _, template := range templates {
		t.Run(template, func(t *testing.T) {
			got, ok := formatTemplate(template, []any{"x", 1})
			if !ok {
				return
			}
			for i := 0; i+1 < len(got); i++ {
				if got[i] == '%' && got[i+1] == '!' {
					t.Errorf("artefato do fmt em %q: %q", template, got)
					return
				}
			}
		})
	}
}
