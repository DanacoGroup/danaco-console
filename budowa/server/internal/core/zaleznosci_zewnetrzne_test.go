package core

import (
	"strings"
	"testing"

	"danacoconsole/server/internal/zewnetrzne"
)

// Sonda zależności ma mówić prawdę o maszynie, a nie o zamiarze.
//
// Wykaz zależności jest jedyną drogą, którą Operator dowiaduje się przed
// czynnością, czego rdzeń nie wykona. Szkoda, której te sprawdziany pilnują, ma
// dwie postacie: wykaz niekompletny (funkcja wymaga programu, którego nikt nie
// zapowiedział) oraz wykaz kłamiący o obecności (pozycja meldowana jako obecna,
// gdy programu nie ma). Pierwsza zostawia Operatora z niespodzianką po
// naciśnięciu, druga jest gorsza — odbiera sens całej sondzie.
//
// Sprawdziany nie mierzą, ile programów stoi na tej maszynie: ta liczba zależy
// od tego, co doinstalowano ręcznie, i na maszynie deweloperskiej jest zawyżona.
// Mierzą mechanizm.

// TestWykazZaleznosciNiesieKompletOpisu pilnuje, żeby każda pozycja mówiła
// Operatorowi trzy rzeczy: co to za program, czym go dociągnąć i co bez niego
// nie zadziała. Pozycja bez zakresu jest wpisem do dziennika, nie informacją.
func TestWykazZaleznosciNiesieKompletOpisu(t *testing.T) {
	wykaz := ZaleznosciZewnetrzne()
	if len(wykaz) == 0 {
		t.Fatal("wykaz zależności jest pusty — rdzeń woła programy spoza instalki i żadnego nie zapowiada")
	}
	for _, pozycja := range wykaz {
		if strings.TrimSpace(pozycja.Narzedzie.Nazwa) == "" {
			t.Errorf("pozycja %q nie ma nazwy czytelnej", pozycja.Narzedzie.Program)
		}
		if strings.TrimSpace(pozycja.Narzedzie.Program) == "" {
			t.Errorf("pozycja %q nie wskazuje programu", pozycja.Narzedzie.Nazwa)
		}
		if strings.TrimSpace(pozycja.Narzedzie.Pakiet) == "" {
			t.Errorf("pozycja %q nie podpowiada, czym ją dociągnąć", pozycja.Narzedzie.Nazwa)
		}
		if strings.TrimSpace(pozycja.Zakres) == "" {
			t.Errorf("pozycja %q nie mówi, co bez niej nie zadziała", pozycja.Narzedzie.Nazwa)
		}
	}
}

// TestSondaOdrozniaProgramObecnyOdNieobecnego wykazuje, że sonda mierzy stan
// maszyny, a nie oddaje wartość stałą. Program nieistniejący ma wracać jako
// brak — inaczej wykaz meldowałby komplet niezależnie od tego, co zastał.
func TestSondaOdrozniaProgramObecnyOdNieobecnego(t *testing.T) {
	if zewnetrzne.Stoi(zewnetrzne.Narzedzie{Program: "program-ktorego-nie-ma-na-zadnej-maszynie"}) {
		t.Fatal("sonda zameldowała obecność programu, którego nie ma — wykaz byłby bezwartościowy")
	}
	// Powłoka jest na każdej maszynie, na której ten sprawdzian w ogóle się
	// uruchomi, więc jej obecność jest miarą tego, że sonda potrafi też
	// potwierdzić — a nie tylko zaprzeczyć.
	if !zewnetrzne.Stoi(zewnetrzne.Narzedzie{Program: "sh"}) {
		t.Fatal("sonda nie widzi powłoki systemowej — mierzy coś innego niż ścieżkę wyszukiwania")
	}
}

// TestBrakujaceZaleznosciSaPodzbioremWykazu pilnuje spójności obu wejść:
// wykaz brakujących nie może zawierać pozycji spoza wykazu pełnego ani
// meldować jako brakującej pozycji, którą sonda uznała za obecną.
func TestBrakujaceZaleznosciSaPodzbioremWykazu(t *testing.T) {
	stan := make(map[string]bool)
	for _, pozycja := range ZaleznosciZewnetrzne() {
		stan[pozycja.Narzedzie.Program] = pozycja.Stoi
	}
	for _, brak := range BrakujaceZaleznosci() {
		obecny, znany := stan[brak.Narzedzie.Program]
		if !znany {
			t.Errorf("program %q melduje się jako brakujący, a nie ma go w wykazie zależności",
				brak.Narzedzie.Program)
		}
		if obecny {
			t.Errorf("program %q stoi na maszynie, a wykaz brakujących go wymienia",
				brak.Narzedzie.Program)
		}
	}
}
