// Dławik prób wejścia przez bramkę: rosnąca zwłoka po próbach nieudanych,
// zerowana pierwszym wejściem udanym.
//
// Dławik nie odmawia żadnej próby — nie ma tu progu, stanu „zablokowane” ani
// kodu błędu „za dużo prób”. Ogranicza wyłącznie prędkość zgadywania: sekret
// zgadywany po łączu lokalnym idzie tysiącami prób na sekundę, przy zwłoce
// sięgającej pięciu sekund schodzi do dwunastu prób na minutę.
//
// Zwłoka nakładana jest na wejściu czynności, przed sprawdzeniem sekretu.
// Czekanie dopiero po rozpoznaniu sekretu jako błędnego czyniłoby z czasu
// odpowiedzi wskaźnik poprawności sekretu.
//
// Licznik prób żyje w pamięci, nie w bazie: jest stanem biegu procesu, nie
// faktem o Operatorze. Zapisany w bazie przeżywałby restart i kazałby czekać
// komuś, kto dopiero zaczyna.
package core

import (
	"context"
	"sync"
	"time"
)

// zwlokaPierwszaDlawika jest zwłoką nałożoną na próbę NASTĘPUJĄCĄ PO pierwszej
// nieudanej. Ćwierć sekundy jest poniżej progu, na którym człowiek zauważa
// opóźnienie interfejsu, a maszynie odbiera już trzy czwarte prędkości.
const zwlokaPierwszaDlawika = 250 * time.Millisecond

// zwlokaGranicznaDlawika jest sufitem, powyżej którego zwłoka nie rośnie.
// Wzrost wykładniczy bez sufitu po kilkunastu próbach daje czekanie liczone
// w godzinach, czyli odmowę wykonaną zegarem. Pięć sekund zostawia Operatorowi
// wejście w każdej chwili, a zgadującemu wyznacza pułap dwunastu prób na minutę.
const zwlokaGranicznaDlawika = 5 * time.Second

// dlawikWejscia trzyma licznik prób nieudanych osobno dla każdej drogi wejścia.
//
// Kluczem jest droga wejścia, nie wołający: bramka jest jedna, a Operator
// bezimienny (wzorzec Danaco HUB), więc nie ma konta, po którym można by liczyć.
// Rozdzielenie po metodzie i urządzeniu sprawia, że seria chybionych PIN-ów na
// tablecie nie spowalnia wejścia hasłem na maszynie roboczej — to dwa różne
// sekrety.
type dlawikWejscia struct {
	mu    sync.Mutex
	proby map[string]int

	// czekaj podstawia własne czekanie w miejsce zegara, żeby test mógł zmierzyć
	// zwłokę bez odczekiwania jej. Wartość zerowa oznacza czekanie prawdziwe.
	czekaj func(ctx context.Context, ile time.Duration)
}

// nowyDlawikWejscia zakłada dławik z pustym licznikiem.
func nowyDlawikWejscia() *dlawikWejscia {
	return &dlawikWejscia{proby: map[string]int{}}
}

// Zaczekaj nakłada zwłokę należną drodze wejścia i zwraca jej długość.
//
// Zwłoka wynika z prób wcześniejszych, więc pierwsza próba nie czeka nigdy.
// Zerwanie kontekstu (rozłączony klient) kończy czekanie natychmiast.
func (d *dlawikWejscia) Zaczekaj(ctx context.Context, droga string) time.Duration {
	if d == nil {
		return 0
	}
	d.mu.Lock()
	nieudane := d.proby[droga]
	d.mu.Unlock()

	zwloka := zwlokaPoProbach(nieudane)
	if zwloka <= 0 {
		return 0
	}
	if d.czekaj != nil {
		d.czekaj(ctx, zwloka)
		return zwloka
	}
	zegar := time.NewTimer(zwloka)
	defer zegar.Stop()
	select {
	case <-zegar.C:
	case <-ctx.Done():
	}
	return zwloka
}

// Niepowodzenie dolicza próbę nieudaną i zwraca zwłokę, która obejmie próbę
// następną. Wartość zwrócona służy wyłącznie opisaniu stanu w odpowiedzi.
func (d *dlawikWejscia) Niepowodzenie(droga string) time.Duration {
	if d == nil {
		return 0
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	// Licznik przestaje rosnąć tam, gdzie zwłoka i tak stoi na suficie —
	// inaczej rósłby bez końca i przepełniłby się przy dość długiej serii.
	if zwlokaPoProbach(d.proby[droga]) < zwlokaGranicznaDlawika {
		d.proby[droga]++
	}
	return zwlokaPoProbach(d.proby[droga])
}

// Wyzeruj kasuje licznik drogi wejścia. Woła się to po wejściu udanym: sekret
// był dobry, więc następna próba nie czeka wcale.
func (d *dlawikWejscia) Wyzeruj(droga string) {
	if d == nil {
		return
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.proby, droga)
}

// zwlokaPoProbach przekłada liczbę prób nieudanych na zwłokę: 0, 250 ms, 500,
// 1000, 2000, 4000, dalej równo 5000 ms. Podwajanie daje szybki spadek prędkości
// zgadywania przy pierwszych kilku próbach, a sufit nie pozwala mu przerodzić
// się w odmowę.
func zwlokaPoProbach(nieudane int) time.Duration {
	if nieudane <= 0 {
		return 0
	}
	zwloka := zwlokaPierwszaDlawika
	for i := 1; i < nieudane; i++ {
		zwloka *= 2
		if zwloka >= zwlokaGranicznaDlawika {
			return zwlokaGranicznaDlawika
		}
	}
	return zwloka
}

// drogaWejscia składa klucz licznika z metody i urządzenia żądania. Urządzenie
// puste (wejście hasłem) daje klucz samej metody, bo hasło bramki jest jedno
// dla całej platformy.
func drogaWejscia(metoda, urzadzenie string) string {
	if urzadzenie == "" {
		return metoda
	}
	return metoda + "\x00" + urzadzenie
}
