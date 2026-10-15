package i18n

import "unicode"

// zwj e o "zero width joiner" (U+200D).
const zwj = '\u200D'

// graphemes divide a string em clusters de grapheme.
//
// Implementa um subconjunto de UAX #29: marcas combinantes, variation
// selectors, emoji modifiers, sequencias com ZWJ, tags e indicadores
// regionais sao mantidos juntos com o caractere base. Sem isso, reverter
// "e" acentuado produziria caracteres invalidos.
func graphemes(s string) []string {
	runes := []rune(s)
	if len(runes) == 0 {
		return nil
	}

	clusters := make([]string, 0, len(runes))
	for i := 0; i < len(runes); {
		start := i
		i++

		// Caracteres Prepend (ex: sinais de numero arabigo) sao unidos ao
		// caractere seguinte.
		if isPrepend(runes[start]) && i < len(runes) {
			i++
		}

	cluster:
		for i < len(runes) {
			r := runes[i]
			switch {
			case isExtend(r):
				// Marca combinante, variation selector, emoji modifier, tag.
				i++
			case runes[i-1] == zwj:
				// Sequencia de emoji com ZWJ.
				i++
			case isRegionalIndicator(r) && isRegionalIndicator(runes[i-1]) &&
				countRegionalIndicators(runes[start:i])%2 == 1:
				// Indicadores regionais formam pares (bandeiras).
				i++
			default:
				break cluster
			}
		}

		clusters = append(clusters, string(runes[start:i]))
	}

	return clusters
}

// isExtend informa se o caractere deve ficar junto do cluster anterior.
func isExtend(r rune) bool {
	switch {
	case r == zwj:
		return true
	case r == '\u200C': // ZWNJ
		return true
	case r >= 0xFE00 && r <= 0xFE0F: // variation selectors
		return true
	case r >= 0xE0100 && r <= 0xE01EF: // variation selectors supplement
		return true
	case r >= 0x1F3FB && r <= 0x1F3FF: // emoji modifiers
		return true
	case r >= 0xE0020 && r <= 0xE007F: // tags
		return true
	}
	return unicode.Is(unicode.Mn, r) || unicode.Is(unicode.Me, r) || unicode.Is(unicode.Mc, r)
}

// isPrepend informa se o caractere deve ser unido ao cluster seguinte.
func isPrepend(r rune) bool {
	switch {
	case r >= 0x0600 && r <= 0x0605: // sinais de numero arabigo
		return true
	case r == 0x06DD: // sinal de fim de ayah
		return true
	case r == 0x070F: // marca de numero syriaco
		return true
	case r == 0x0890 || r == 0x0891: // sinais de numero arabigo
		return true
	case r == 0x08E2: // marca de numero arabe
		return true
	case r == 0x110BD || r == 0x110CD: // numero kaithi
		return true
	}
	return false
}

// isRegionalIndicator informa se o caractere e um indicador regional
// (bandeiras/emoji de pais).
func isRegionalIndicator(r rune) bool {
	return r >= 0x1F1E6 && r <= 0x1F1FF
}

// countRegionalIndicators conta os indicadores regionais de um trecho.
func countRegionalIndicators(runes []rune) int {
	total := 0
	for _, r := range runes {
		if isRegionalIndicator(r) {
			total++
		}
	}
	return total
}
