package i18n

import (
	"fmt"
	"strconv"
	"strings"
)

// validVerbs sao os verbos aceitos pelo fmt.Sprintf.
//
// O verbo "%w" nao entra na lista: ele so faz sentido em fmt.Errorf e no
// fmt.Sprintf produziria o artefato "%!w(...)".
const validVerbs = "vTtbcdoOqxXUeEfFgGsp"

// formatTemplate aplica os argumentos a um template no estilo fmt.Sprintf.
//
// A diferenca para o fmt.Sprintf e que a funcao NUNCA devolve os artefatos de
// erro do fmt ("%!d(MISSING)", "%!(EXTRA ...)", "%!z(...)"). Quando o
// template nao pode ser formatado com os argumentos informados, a funcao
// devolve false e o chamador usa o template original.
//
// Regras aplicadas:
//
//   - "%%" produz um "%" literal;
//   - um "%" que nao inicia uma diretiva valida e tratado como "%" literal
//     (por exemplo "50% off" continua "50% off");
//   - argumentos extras sao ignorados (nenhum "%!(EXTRA ...)" e gerado);
//   - argumentos faltando invalidam o template inteiro: nada e formatado
//     pela metade, o template original e devolvido.
func formatTemplate(template string, args []any) (string, bool) {
	var b strings.Builder
	next := 0 // indice do proximo argumento implicito

	for i := 0; i < len(template); {
		if template[i] != '%' {
			b.WriteByte(template[i])
			i++
			continue
		}

		// "%%" -> "%" literal.
		if i+1 < len(template) && template[i+1] == '%' {
			b.WriteByte('%')
			i += 2
			continue
		}

		d, parsed := parseDirective(template[i:])
		if !parsed {
			// O "%" nao inicia uma diretiva valida: e um caractere literal.
			b.WriteByte('%')
			i++
			continue
		}

		indices, proximo, completo := d.resolveIndices(next, len(args))
		if !completo {
			// Faltam argumentos: o template inteiro e invalido.
			return "", false
		}
		next = proximo

		values := make([]any, 0, len(indices))
		for _, idx := range indices {
			values = append(values, args[idx])
		}

		formatted := fmt.Sprintf(d.spec, values...)
		if strings.Contains(formatted, "%!") {
			// O proprio fmt nao conseguiu formatar (ex: "%p" com string):
			// melhor devolver o template original do que um artefato.
			return "", false
		}

		b.WriteString(formatted)
		i += d.size
	}

	return b.String(), true
}

// directive e uma diretiva de formatacao ("%s", "%-10.3f", "%[2]d", ...).
type directive struct {
	spec      string // spec sem o indice explicito do argumento
	size      int    // quantidade de bytes consumidos no template
	explicit  int    // indice (base 0) do argumento explicito, ou -1
	starWidth bool   // a largura usa "*" (consome um argumento)
	starPrec  bool   // a precisao usa "*" (consome um argumento)
}

// resolveIndices devolve os indices dos argumentos consumidos pela
// diretiva e o proximo indice implicito.
//
// Retorna false quando algum argumento necessario nao existe.
func (d directive) resolveIndices(next, total int) ([]int, int, bool) {
	if d.explicit >= 0 {
		if d.explicit >= total {
			return nil, next, false
		}
		return []int{d.explicit}, d.explicit + 1, true
	}

	indices := make([]int, 0, 3)
	if d.starWidth {
		indices = append(indices, next)
		next++
	}
	if d.starPrec {
		indices = append(indices, next)
		next++
	}
	indices = append(indices, next)
	next++

	for _, idx := range indices {
		if idx >= total {
			return nil, next, false
		}
	}
	return indices, next, true
}

// parseDirective analisa uma diretiva de formatacao no inicio de format.
//
// Devolve false quando o texto nao e uma diretiva valida.
func parseDirective(format string) (directive, bool) {
	var d directive

	i := 1 // consome o "%"

	// Flags aceitos pelo fmt (o flag de espaco NAO e aceito: "% off" deve
	// ser tratado como "%" literal seguido de texto).
	flagsStart := i
	for i < len(format) {
		switch format[i] {
		case '+', '-', '#', '0':
			i++
			continue
		}
		break
	}
	flags := format[flagsStart:i]

	// Indice explicito do argumento: "[n]".
	d.explicit = -1
	if i < len(format) && format[i] == '[' {
		j := i + 1
		digitsStart := j
		for j < len(format) && format[j] >= '0' && format[j] <= '9' {
			j++
		}
		if j == digitsStart || j >= len(format) || format[j] != ']' {
			return d, false
		}
		n, err := strconv.Atoi(format[digitsStart:j])
		if err != nil || n <= 0 {
			return d, false
		}
		d.explicit = n - 1
		i = j + 1
	}

	// Largura.
	width := ""
	if i < len(format) && format[i] == '*' {
		d.starWidth = true
		width = "*"
		i++
	} else {
		widthStart := i
		for i < len(format) && format[i] >= '0' && format[i] <= '9' {
			i++
		}
		width = format[widthStart:i]
	}

	// Precisao.
	precision := ""
	if i < len(format) && format[i] == '.' {
		i++
		precStart := i
		for i < len(format) && format[i] >= '0' && format[i] <= '9' {
			i++
		}
		precision = "." + format[precStart:i]
		if i < len(format) && format[i] == '*' {
			d.starPrec = true
			precision = ".*"
			i++
		}
	}

	// Verbo.
	if i >= len(format) {
		return d, false
	}
	verb := rune(format[i])
	if !strings.ContainsRune(validVerbs, verb) {
		return d, false
	}
	i++

	// Um indice explicito combinado com "*" tem semantica ambigua no fmt
	// e nao e suportado aqui.
	if d.explicit >= 0 && (d.starWidth || d.starPrec) {
		return d, false
	}

	d.spec = "%" + flags + width + precision + string(verb)
	d.size = i
	return d, true
}
