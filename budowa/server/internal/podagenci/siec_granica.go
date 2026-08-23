// Granica sieci podagentów: ilu podagentów pracuje jednocześnie pod jednym oknem
// wykonawcy. Granica dotyczy stanu sieci, nie kształtu pojedynczego żądania —
// przycinanie samego wywołania przepuszcza dwa wywołania po piętnaście.
//
// Miejsce w sieci zajmują wyłącznie podagenci czynni, czyli ci, którzy pracę mają
// przed sobą albo w toku. Podagent zakończony, błędny i zatrzymany miejsca nie
// trzyma — inaczej okno zużyłoby piętnaście miejsc raz na całą swoją historię.
//
// Żądanie większe niż liczba wolnych miejsc wykonuje się na tylu, ile jest wolnych.
// Odmowa zostaje na jeden przypadek: sieć pełna, gdzie przyciąć można wyłącznie do
// zera, a powołanie zerowe udawałoby wykonanie. Odmowa niesie wtedy zdanie mówiące,
// co zrobić — zebrać wyniki albo zatrzymać któregoś.
package podagenci

import (
	"context"
	"fmt"

	"danacoconsole/server/internal/dane"
)

// GranicaSieci to liczba podagentów, którzy mogą jednocześnie pracować pod
// jednym oknem wykonawcy.
const GranicaSieci = 15

// StanKoncowy mówi, czy podagent w tym stanie nie wróci już do pracy sam.
//
// Orzeczenie jest wspólne dla granicy sieci (kto zajmuje miejsce), zbierania
// wyników (czy komplet gotowy) i domknięcia zatrzymania (czy wiersz wciąż kłamie).
// Trzy osobne odpowiedzi rozjechałyby się przy pierwszej zmianie słownika stanów.
func StanKoncowy(stan string) bool {
	switch stan {
	case dane.StanPodagentaUkonczony, dane.StanPodagentaBledny, dane.StanPodagentaZatrzymany:
		return true
	default:
		return false
	}
}

// Czynni liczy podagentów zajmujących miejsce w sieci.
func Czynni(wiersze []dane.Podagent) int {
	ile := 0
	for _, podagent := range wiersze {
		if !StanKoncowy(podagent.Stan) {
			ile++
		}
	}
	return ile
}

// Przydzial jest odpowiedzią granicy na żądanie powołania.
type Przydzial struct {
	// Ile podagentów wolno powołać teraz. Zero znaczy sieć pełną.
	Ile int
	// Zajete to liczba miejsc zajętych w chwili pytania.
	Zajete int
	// Zadane to liczba wyprowadzona z żądania, po odczytaniu pustego wskazania.
	Zadane int
	// Powod jest niepusty wtedy, gdy przydział jest mniejszy niż liczba wskazana
	// w żądaniu — zdanie do dziennika albo do odmowy. Pusty powód znaczy żądanie
	// spełnione w całości.
	Powod string
}

// Pelna mówi, że sieć nie ma ani jednego wolnego miejsca. Tylko wtedy powołanie
// odmawia zamiast przyciąć.
func (p Przydzial) Pelna() bool { return p.Ile == 0 }

// LiczbaZadana odczytuje wskazanie liczby z żądania kontraktu.
//
// Brak wskazania znaczy jednego, wskazanie poniżej jedynki również jednego:
// powołanie zerowe nie jest żądaniem, tylko pomyłką klienta. Wskazanie ponad
// granicę przycina się do granicy.
func LiczbaZadana(wskazanie *int) int {
	if wskazanie == nil || *wskazanie < 1 {
		return 1
	}
	if *wskazanie > GranicaSieci {
		return GranicaSieci
	}
	return *wskazanie
}

// MiejscaWSieci rozstrzyga, ilu podagentów wolno powołać pod oknem, którego
// zastanych podagentów podano.
//
// Wiersze muszą pochodzić z jednego okna: granica jest granicą sieci wykonawcy,
// a nie granicą platformy — liczenie kompletu bazy zamknęłoby powołanie
// u wszystkich, gdy jedno okno zapełni swoją sieć.
func MiejscaWSieci(zastani []dane.Podagent, wskazanie *int) Przydzial {
	zadane := LiczbaZadana(wskazanie)
	zajete := Czynni(zastani)
	wolne := GranicaSieci - zajete
	if wolne < 0 {
		wolne = 0
	}
	przydzial := Przydzial{Ile: zadane, Zajete: zajete, Zadane: zadane}
	if wskazanie != nil && *wskazanie > GranicaSieci {
		przydzial.Powod = fmt.Sprintf(
			"żądano %d podagentów, a granica sieci wynosi %d — powołuję %d",
			*wskazanie, GranicaSieci, zadane)
	}
	if zadane <= wolne {
		return przydzial
	}
	przydzial.Ile = wolne
	if wolne == 0 {
		przydzial.Powod = fmt.Sprintf(
			"sieć podagentów tego okna jest pełna: %d z %d czynnych. "+
				"Zbierz wyniki zakończonych (subagent.result.collect) albo "+
				"zatrzymaj któregoś (subagent.stop), zanim powołasz kolejnego",
			zajete, GranicaSieci)
		return przydzial
	}
	przydzial.Powod = fmt.Sprintf(
		"żądano %d podagentów, a w sieci tego okna wolnych miejsc jest %d "+
			"(%d z %d czynnych) — powołuję %d",
		zadane, wolne, zajete, GranicaSieci, wolne)
	return przydzial
}

// WykazSieci jest wąskim kontraktem odczytu potrzebnym granicy — tyle
// z repozytorium podagentów, ile granica naprawdę czyta.
// `dane.RepozytoriumPodagentow` wypełnia go bez przeróbek.
type WykazSieci interface {
	Podagenci(ctx context.Context, filtr dane.FiltrPodagentow) ([]dane.Podagent, error)
}

// PrzydzialOkna pyta trwałość o zastaną sieć okna i rozstrzyga przydział.
//
// Bez repozytorium albo bez kodu okna granica nie ma czego liczyć i oddaje
// przydział samego wywołania. Błąd odczytu wraca do wołającego — przydział
// wyliczony z nieudanego odczytu byłby granicą zmyśloną.
func PrzydzialOkna(ctx context.Context, wykaz WykazSieci, oknoKod string,
	wskazanie *int) (Przydzial, error) {

	if wykaz == nil || oknoKod == "" {
		return Przydzial{Ile: LiczbaZadana(wskazanie), Zadane: LiczbaZadana(wskazanie)}, nil
	}
	zastani, err := wykaz.Podagenci(ctx, dane.FiltrPodagentow{OknoKod: oknoKod})
	if err != nil {
		return Przydzial{}, fmt.Errorf(
			"podagenci: nie można odczytać sieci okna %q, więc granicy nie da się zmierzyć: %w",
			oknoKod, err)
	}
	return MiejscaWSieci(zastani, wskazanie), nil
}
