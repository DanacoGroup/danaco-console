// Odpowiedzialność pliku: znaki marki dołączane częściami do listów
// transakcyjnych. Same listy składa pakiet `internal/mail`.
package listy

import _ "embed"

/*
Znaki obu wariantów. Nagłówek listu niesie dwa znaczniki obrazu: jasny widoczny
domyślnie, ciemny odsłaniany zapytaniem medialnym (opracowanie, rozdz. 4.3).
Lockup jasny ma atrament `#181818` i na tle ciemnym znika.
*/

// Znaki idą w binarce, nie z dysku: rdzeń wysyła listy z maszyny wdrożenia,
// na której katalogu `design/` nie ma.

//go:embed znak-marki.png
var ZnakJasny []byte

//go:embed znak-marki-ciemny.png
var ZnakCiemny []byte
