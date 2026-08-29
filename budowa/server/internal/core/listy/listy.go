// Odpowiedzialność pliku: postać graficzna listów systemowych — papeteria marki
// wraz ze znakiem, w którą rdzeń wpisuje treść potwierdzenia albo odzyskania.
package listy

import (
	_ "embed"
	"strings"
)

/*
Szablon i znak marki idą w binarce, nie z dysku: rdzeń wysyła listy z maszyny
wdrożenia, na której katalogu `design/` nie ma. Plik odczytywany w chwili
wysyłki byłby drogą, na której brak jednego pliku zamienia list w pustą kartkę.
*/

//go:embed potwierdzenie.html
var szablonPotwierdzenia string

//go:embed znak-marki.txt
var znakMarki string

// TrescListu niesie wartości wpisywane w papeterię. Wszystkie są tekstem
// gotowym do pokazania — składanie zdań należy do wołającego, nie do szablonu.
type TrescListu struct {
	// Naglowek nazywa sprawę listu, jednym zdaniem.
	Naglowek string
	// Wstep mówi, skąd ten list się wziął.
	Wstep string
	// EtykietaKodu stoi nad ramką z kodem, wersalikami.
	EtykietaKodu string
	// Kod jest tym, co Operator przepisuje do okna.
	Kod string
	// Polecenie mówi, co z kodem zrobić i jak długo jest ważny.
	Polecenie string
	// Nota niesie ostrzeżenie albo zdanie zamykające, drobnym drukiem.
	Nota string
}

/*
Zloz wpisuje treść w papeterię i oddaje gotową postać graficzną listu.

Podstawienie idzie przez `strings.NewReplacer`, nie przez szablon języka Go:
wartości pochodzą z rdzenia, nie od Operatora, a żadna z nich nie niesie
znaczników — pojedyncze przejście po tekście jest tu wystarczające i nie wnosi
zależności od pakietu szablonów.
*/
func Zloz(t TrescListu) string {
	return strings.NewReplacer(
		"{znak}", strings.TrimSpace(znakMarki),
		"{naglowek}", t.Naglowek,
		"{wstep}", t.Wstep,
		"{etykietaKodu}", t.EtykietaKodu,
		"{kod}", t.Kod,
		"{polecenie}", t.Polecenie,
		"{nota}", t.Nota,
	).Replace(szablonPotwierdzenia)
}
