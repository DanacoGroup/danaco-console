package core

import (
	"context"
	"io"
	"log"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"danacoconsole/server/internal/konfiguracja"
	"danacoconsole/server/internal/store"
)

// Uprząż sprawdzianów rdzenia.
//
// Rdzeń nie ma odbiornika, który dałoby się złożyć w pamięci: montaż otwiera
// repozytoria nad bazą, odtwarza stan sesji z wierszy i startuje budzik
// harmonogramu. Zaślepianie tego łańcucha dałoby sprawdzian zaślepki, nie
// sprawdzian produktu — dlatego uprząż montuje rdzeń prawdziwy, tylko nad bazą
// jednorazową.
//
// Baza jest plikiem w katalogu tymczasowym, nie `:memory:`. Pula połączeń
// rdzenia trzyma cztery połączenia (store.maksPolaczen), a każde połączenie do
// `:memory:` dostaje w SQLite własną, osobną bazę — migracje wykonałyby się na
// jednej, a zapytania trafiły na trzy puste. Plik w `t.TempDir()` znosi ten
// problem bez kosztu: sterownik jest czystym Go, więc pełny przejazd migracji
// idzie w milisekundach, a katalog znika po sprawdzianie sam.

// zmontujDoSprawdzenia składa rdzeń nad świeżą bazą i pilnuje zwolnienia
// zasobów po zakończeniu sprawdzianu. Zwraca zmontowany rdzeń wraz z jego
// kontekstem życia — kontekst przydaje się sprawdzianom wywołującym komendy.
func zmontujDoSprawdzenia(t *testing.T) (*Zmontowany, context.Context) {
	t.Helper()

	return zmontujNadDziennikiem(t, io.Discard)
}

// zmontujZDziennikiem składa ten sam rdzeń, ale nad dziennikiem, który da się
// przeczytać. Potrzebny tam, gdzie przedmiotem pomiaru jest sam zapis — rdzeń
// odnotowuje w dzienniku rzeczy, których nie ma jak oddać w odpowiedzi, więc
// dziennik niemy zamieniłby taki sprawdzian w sprawdzenie niczego.
func zmontujZDziennikiem(t *testing.T) (*Zmontowany, context.Context, *dziennikDoOdczytu) {
	t.Helper()

	dziennik := &dziennikDoOdczytu{}
	zmontowany, zycie := zmontujNadDziennikiem(t, dziennik)
	return zmontowany, zycie, dziennik
}

// zmontujNadDziennikiem jest wspólnym montażem obu uprzęży. Osobna, bo dwa
// montaże rozjechałyby się przy pierwszej zmianie nastaw sprawdzianu.
func zmontujNadDziennikiem(t *testing.T, zapis io.Writer) (*Zmontowany, context.Context) {
	t.Helper()

	katalog := t.TempDir()
	baza, err := store.Otworz(filepath.Join(katalog, "dane.sqlite"))
	if err != nil {
		t.Fatalf("nie można otworzyć bazy sprawdzianu: %v", err)
	}
	t.Cleanup(func() { _ = baza.Zamknij() })

	// Kontekst życia zamyka się przed zwolnieniem zasobów: budzik harmonogramu
	// i pętle adapterów wiszą na nim, a nie na Zamknij.
	zycie, zakoncz := context.WithCancel(context.Background())
	t.Cleanup(zakoncz)

	ustawienia := konfiguracja.Domyslna()
	ustawienia.KatalogDanych = katalog
	// Pakietu klienta i profili kanału głównego sprawdzian nie ma i mieć nie
	// musi: montaż opisuje ich brak jako dopuszczalny, a sprawdzian zgodności
	// z kontraktem dotyczy rejestru komend, nie plików obok gniazda.
	ustawienia.KatalogKlienta = ""
	ustawienia.KatalogProfili = ""
	// Port zerowy: uprząż nie wystawia nasłuchu, więc żaden port nie jest
	// zajmowany. Sprawdziany transportu podnoszą własny serwer osobno.
	ustawienia.Port = 0

	zmontowany, err := Zmontuj(zycie, Montaz{
		Konfiguracja: ustawienia,
		Baza:         baza,
		Dziennik:     log.New(zapis, "", 0),
	})
	if err != nil {
		t.Fatalf("montaż rdzenia nie powiódł się: %v", err)
	}
	t.Cleanup(zmontowany.Zamknij)

	return zmontowany, zycie
}

// dziennikNiemy zwraca dziennik zapisujący donikąd. Rdzeń pisze do dziennika
// przy starcie i przy każdej degradacji; w sprawdzianie te zapisy zaśmiecałyby
// wyjście, a ich brak niczego nie zmienia — dziennik jest opcjonalny w każdym
// miejscu montażu.
func dziennikNiemy() *log.Logger {
	return log.New(io.Discard, "", 0)
}

// dziennikDoOdczytu zbiera zapisy dziennika rdzenia, żeby sprawdzian mógł je
// przeczytać.
//
// Pod zamkiem, bo rdzeń pisze do dziennika także z własnych gorutyn — budzika
// harmonogramu i pętli adapterów — a sprawdzian czyta z gorutyny swojej.
type dziennikDoOdczytu struct {
	zamek sync.Mutex
	tresc strings.Builder
}

func (d *dziennikDoOdczytu) Write(bajty []byte) (int, error) {
	d.zamek.Lock()
	defer d.zamek.Unlock()
	return d.tresc.Write(bajty)
}

// Tresc oddaje wszystko, co dotąd trafiło do dziennika.
func (d *dziennikDoOdczytu) Tresc() string {
	d.zamek.Lock()
	defer d.zamek.Unlock()
	return d.tresc.String()
}
