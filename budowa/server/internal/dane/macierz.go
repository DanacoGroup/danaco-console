// Odpowiedzialność pliku: odczyt macierzy widoczności modułów w środowiskach
// (tabela `srodowisko_modul`). Macierz jest konfiguracją: wiersz mówi, czy
// moduł stoi w bocznej nawigacji danego środowiska i na którym miejscu.
// Wiersze wnoszą migracje schematu — repozytorium ich nie zakłada.
//
// Osobny byt obok `RepozytoriumModulow`: tamto repozytorium czyta moduły,
// a macierzy używa wyłącznie jako filtra (`ListaSrodowiska`). Tu bytem jest sam
// wiersz macierzy — z parą kodów, kolejnością i widocznością — którego tamten
// kształt nie umie oddać. To nie druga prawda o module: jedna tabela, dwa różne
// pytania.
package dane

import (
	"context"
	"fmt"
)

// WierszMacierzy to jeden wiersz tabeli `srodowisko_modul` wraz z kodami obu
// stron pary. Kody wchodzą tu razem z identyfikatorami, bo czytelnik macierzy
// prawie zawsze potrzebuje kodu, a nie numeru wiersza — a drugie zapytanie po
// słownik byłoby powrotem do pętli N+1, którą ten byt właśnie znosi.
type WierszMacierzy struct {
	SrodowiskoID  int64
	SrodowiskoKod string
	ModulID       int64
	ModulKod      string
	// Kolejnosc to pozycja modułu w bocznej nawigacji środowiska. Dla pary
	// niewidocznej nie niesie treści — pozycji poza wykazem nie ma.
	Kolejnosc int64
	Widoczny  bool
}

// RepozytoriumMacierzy jest kontraktem odczytu macierzy widoczności.
//
// Zapisu tu nie ma świadomie. Kontrakt platformy nie definiuje ani jednej
// komendy zmieniającej macierz, więc metoda zapisu nie miałaby drogi wywołania
// — a byt bez drogi wywołania jest atrapą. Gdy komenda powstanie, pisarz
// dopisze się do tego samego interfejsu.
type RepozytoriumMacierzy interface {
	// Pelna zwraca wszystkie wiersze macierzy — widoczne i niewidoczne —
	// w kolejności kart środowisk, a wewnątrz środowiska w kolejności nawigacji.
	Pelna(ctx context.Context) ([]WierszMacierzy, error)
	// KodySrodowiskModulow zwraca odwzorowanie `modul.id` → kody środowisk,
	// w których moduł jest WIDOCZNY, w kolejności kart środowisk. Moduł
	// nieobecny w wyniku nie ma okna modułowego w żadnym środowisku — jest
	// dostępny wyłącznie ze strony głównej.
	KodySrodowiskModulow(ctx context.Context) (map[int64][]string, error)
}

const (
	kolumnyMacierzy = `s.id, s.kod, m.id, m.kod, sm.kolejnosc, sm.widoczny`

	zlaczenieMacierzy = ` FROM srodowisko_modul sm
	                      JOIN srodowisko s ON s.id = sm.srodowisko_id
	                      JOIN modul m      ON m.id = sm.modul_id`

	// Porządek środowisk jest ten sam, co w `listaSrodowisk` (kolejnosc, kod),
	// a wewnątrz środowiska ten sam, co w `listaModulowSrodowiska`
	// (sm.kolejnosc, m.kod). Dzięki temu jedno złączenie oddaje dokładnie ten
	// wynik, który dawała pętla po środowiskach — co do zawartości i kolejności.
	porzadekMacierzy = ` ORDER BY s.kolejnosc, s.kod, sm.kolejnosc, m.kod`

	pelnaMacierz = `SELECT ` + kolumnyMacierzy + zlaczenieMacierzy + porzadekMacierzy

	widocznaMacierz = `SELECT ` + kolumnyMacierzy + zlaczenieMacierzy +
		` WHERE sm.widoczny = 1` + porzadekMacierzy
)

type repozytoriumMacierzy struct {
	zapytania *zapytania
}

func noweRepozytoriumMacierzy(z *zapytania) *repozytoriumMacierzy {
	return &repozytoriumMacierzy{zapytania: z}
}

// Pelna zwraca komplet wierszy macierzy.
func (r *repozytoriumMacierzy) Pelna(ctx context.Context) ([]WierszMacierzy, error) {
	return r.wykaz(ctx, pelnaMacierz, "macierzy widoczności")
}

// KodySrodowiskModulow składa odwzorowanie moduł → kody środowisk widocznych.
//
// Jedno zapytanie zamiast zapytania na każde środowisko: macierz jest mała,
// ale czyta ją każde `home.enter`, `environment.list`, `environment.enter`
// i `module.list`, więc N+1 płaciło się przy każdym wejściu Operatora.
func (r *repozytoriumMacierzy) KodySrodowiskModulow(ctx context.Context) (map[int64][]string, error) {
	wiersze, err := r.wykaz(ctx, widocznaMacierz, "widocznej macierzy")
	if err != nil {
		return nil, err
	}
	kody := map[int64][]string{}
	for _, wiersz := range wiersze {
		kody[wiersz.ModulID] = append(kody[wiersz.ModulID], wiersz.SrodowiskoKod)
	}
	return kody, nil
}

// wykaz wykonuje zapytanie zwracające wiele wierszy macierzy.
func (r *repozytoriumMacierzy) wykaz(ctx context.Context, zapytanie, opis string) ([]WierszMacierzy, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać %s: %w", opis, err)
	}
	defer wiersze.Close()

	lista := []WierszMacierzy{}
	for wiersze.Next() {
		wiersz, err := odczytajWierszMacierzy(wiersze)
		if err != nil {
			return nil, err
		}
		lista = append(lista, wiersz)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt %s: %w", opis, err)
	}
	return lista, nil
}

// odczytajWierszMacierzy składa strukturę z jednego wiersza wyniku.
func odczytajWierszMacierzy(wiersz skaner) (WierszMacierzy, error) {
	var para WierszMacierzy
	var widoczny int
	if err := wiersz.Scan(
		&para.SrodowiskoID,
		&para.SrodowiskoKod,
		&para.ModulID,
		&para.ModulKod,
		&para.Kolejnosc,
		&widoczny,
	); err != nil {
		return WierszMacierzy{}, err
	}
	para.Widoczny = widoczny != 0
	return para, nil
}
