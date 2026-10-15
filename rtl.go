package i18n

import "strings"

// rtlLanguages sao os codigos de idioma (ISO 639-1/639-3) escritos da direita
// para a esquerda.
var rtlLanguages = map[string]bool{
	"ar":  true, // arabe
	"arc": true, // aramaico imperial
	"ckb": true, // curdo central
	"dv":  true, // divehi
	"fa":  true, // persa
	"he":  true, // hebraico
	"iw":  true, // hebraico (codigo legado)
	"ji":  true, // iidiche (codigo legado)
	"khw": true, // khowar
	"ks":  true, // caxemira
	"mzn": true, // mazandarani
	"nqo": true, // n'ko
	"ps":  true, // pashto
	"sd":  true, // sindhi
	"syr": true, // siriaco
	"ug":  true, // uigur
	"ur":  true, // urdu
	"yi":  true, // iidiche
}

// rtlScripts sao os codigos de script (ISO 15924) da direita para a
// esquerda.
var rtlScripts = map[string]bool{
	"Adlm": true, // adlam
	"Arab": true, // arabe
	"Armi": true, // aramaico imperial
	"Avst": true, // avestano
	"Chrs": true, // corasmiano
	"Elym": true, // elimaico
	"Hatr": true, // hatran
	"Hung": true, // hungaro antigo
	"Mand": true, // mandaico
	"Mani": true, // maniqueismo
	"Mend": true, // mende kikakui
	"Narb": true, // arabe antigo do norte
	"Nbat": true, // nabateano
	"Nkoo": true, // n'ko
	"Orkh": true, // turquico antigo
	"Ougr": true, // uigur antigo
	"Phlp": true, // pahlavi dos salmos
	"Phnx": true, // fenicio
	"Rohg": true, // rohingya hanifi
	"Samr": true, // samaritano
	"Sogo": true, // sogdiano
	"Syrc": true, // siriaco
	"Thaa": true, // thaana
	"Yezi": true, // yezidi
}

// IsRTL informa se o locale usa um sistema de escrita da direita para a
// esquerda.
//
// A verificacao usa o subtag de idioma (ex: "ar", "he", "fa", "ur") e,
// quando presente, o subtag de script (ex: "az-Arab"). Quando o script esta
// presente, ele tem precedencia sobre o idioma.
//
//	IsRTL("ar-EG")   // true
//	IsRTL("he")      // true
//	IsRTL("az-Arab") // true
//	IsRTL("az-Latn") // false
//	IsRTL("pt-BR")   // false
func IsRTL(locale string) bool {
	locale = normalizeLocale(locale)
	if locale == "" {
		return false
	}

	parts := strings.Split(locale, "-")
	for _, part := range parts[1:] {
		if isScriptSubtag(part) {
			// O script tem precedencia sobre o idioma.
			return rtlScripts[part]
		}
	}

	return rtlLanguages[strings.ToLower(parts[0])]
}

// Reverse reordena a string na ordem visual (da esquerda para a direita),
// aplicando as regras bidirecionais do Unicode (UAX #9) sobre clusters de
// grapheme.
//
// O texto e dividido em clusters de grapheme, a direcao de cada caractere e
// classificada e os niveis bidirecionais sao resolvidos. O resultado e a
// ordem em que o texto DEVE ser desenhado por um terminal que nao faz
// reordenacao bidirecional.
//
// Exemplos:
//
//	Reverse("שלום")         // "םולש" (palavra RTL na ordem visual)
//	Reverse("Hello שלום")    // "Hello םולש"
//	Reverse("Hello")         // "Hello" (texto LTR mantem a ordem)
//
// A funcao devolve a string sem alteracoes quando ela nao contem texto da
// direita para a esquerda: em texto LTR a ordem visual e igual a ordem
// logica.
//
// Escopo: implementa um subconjunto da UAX #9 (resolucao de neutros e
// numeros e reordenacao por niveis). Embeddings explicitos (LRE/RLE/LRO/RLO)
// e isoladores nao sao processados.
func Reverse(s string) string {
	return reverseVisual(s)
}
