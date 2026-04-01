package i18n_test

import (
	"fmt"

	i18n "github.com/junhinhow/charm-i18n"
)

func Example() {
	cat := i18n.NewCatalog()
	cat.AddLocale("en", map[string]string{
		"hello": "Hello, World!",
	})
	cat.AddLocale("pt-BR", map[string]string{
		"hello": "Ola, Mundo!",
	})

	cat.SetLocale("pt-BR")
	fmt.Println(cat.T("hello"))
	// Output: Ola, Mundo!
}

func ExampleCatalog_Tf() {
	cat := i18n.NewCatalog()
	cat.AddLocale("en", map[string]string{
		"greeting": "Hello, %s!",
	})

	fmt.Println(cat.Tf("greeting", "Gopher"))
	// Output: Hello, Gopher!
}
