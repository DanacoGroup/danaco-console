// Odpowiedzialność pliku: szablony listów wbudowane w binarkę rdzenia
// i podane pakietowi jako gotowy komplet.
package mail

import "embed"

/*
Szablony idą w binarce, nie z dysku: rdzeń wysyła listy z maszyny wdrożenia,
na której katalogu `design/` nie ma. Układ katalogów html/ i text/ jest ten,
którego oczekuje Load.

Znaki marki nie należą do tego pakietu — Build przyjmuje je argumentem Logos.
*/

//go:embed html/*.html text/*.txt
var szablonyWbudowane embed.FS

// WbudowanyKomplet zwraca komplet szablonów wczytany z binarki. Wzór wywołania
// podaje komentarz do Load.
func WbudowanyKomplet() (*Set, error) {
	return Load(szablonyWbudowane)
}
