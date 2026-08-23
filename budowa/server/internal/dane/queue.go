// Odpowiedzialność pliku: dwie czynności obszaru kolejek spoza zakresu komend
// `queue.create` i `queue.action` obsługiwanych w `kolejki.go` — wykaz kolejek
// (komenda `queue.list`) oraz powiązania kolejki z oknami, ekspertem, projektem
// i automatyką (komenda `queue.link`, tabela `powiazanie_kolejki`).
//
// Metody wiszą na tym samym typie `repozytoriumKolejek`, nad tą samą bazą
// i z tym samym dziennikiem akcji (`log_akcji_kolejki`) — obszar kolejek ma
// jedno repozytorium, rozłożone na dwa pliki wyłącznie dla objętości.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"danacoconsole/shared"
)

// Rodzaje powiązań kolejki — wartości kolumny `powiazanie_kolejki.rodzaj`.
// Słownik jest zamknięty więzem CHECK na tej kolumnie; te stałe są jego jedynym
// odwzorowaniem w Go, żeby literał nie rozjechał się z więzem.
const (
	RodzajPowiazaniaOkno       = "okno"
	RodzajPowiazaniaEkspert    = "ekspert"
	RodzajPowiazaniaProjekt    = "projekt"
	RodzajPowiazaniaAutomatyka = "automatyka"
)

// PowiazanieKolejki to wiersz tabeli `powiazanie_kolejki`: jeden byt jednego
// rodzaju, związany z jedną kolejką. Pole `Byt` niesie identyfikator w kształcie
// kontraktu (identyfikator okna, kod eksperta, kod projektu, identyfikator
// automatyki), a nie klucz wewnętrzny.
type PowiazanieKolejki struct {
	KolejkaID int64
	Rodzaj    string
	Byt       string
	Utworzono string
}

const (
	listaKolejek = `SELECT ` + kolumnyKolejki + ` FROM kolejka ORDER BY id`

	listaKolejekStanu = `SELECT ` + kolumnyKolejki + ` FROM kolejka WHERE stan = ? ORDER BY id`

	// INSERT OR IGNORE, bo powtórzone powiązanie jest TYM SAMYM faktem, a nie
	// drugim. Kolizja z więzem UNIQUE nie jest tu błędem — jest brakiem zmiany.
	wstawPowiazanieKolejki = `INSERT OR IGNORE INTO powiazanie_kolejki (kolejka_id, rodzaj, byt)
	                          VALUES (?, ?, ?)`

	listaPowiazanKolejki = `SELECT kolejka_id, rodzaj, byt, utworzono
	                        FROM powiazanie_kolejki WHERE kolejka_id = ?
	                        ORDER BY rodzaj, id`

	czyKolejkaIstnieje = `SELECT 1 FROM kolejka WHERE id = ?`
)

// CzyKolejkaIstnieje mówi, czy wiersz kolejki jest w bazie. Osobna czynność,
// bo `PobierzKolejke` zwraca brak wiersza i awarię odczytu tym samym zwykłym
// błędem — a to dwie różne odpowiedzi kontraktu: `not_found` i `internal_error`.
func (r *repozytoriumKolejek) CzyKolejkaIstnieje(ctx context.Context, id int64) (bool, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, czyKolejkaIstnieje)
	if err != nil {
		return false, err
	}
	var jeden int
	err = polecenie.QueryRowContext(ctx, id).Scan(&jeden)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("dane: nie można sprawdzić istnienia kolejki %d: %w", id, err)
	}
	return true, nil
}

// ListaKolejek zwraca kolejki zawężone stanem. Pusty stan znaczy „wszystkie" —
// nie „żadne": komenda `queue.list` ma stan jako pole opcjonalne.
//
// Zawężenia po sesji i po oknie tu nie ma celowo. Sesja kolejki bywa znana
// wyłącznie z pamięci powiązań rdzenia (kolejka założona przed pierwszą
// utrwaloną wiadomością sesji nie ma czym wypełnić `kolejka.sesja_id`), a okna
// obsługiwane przez kolejkę leżą w `powiazanie_kolejki` albo w tej samej
// pamięci. Sito po obu tych bytach składa się więc w rdzeniu, na kolejce
// kontraktu — sito w SQL milczałoby o kolejkach, które warunek spełniają.
func (r *repozytoriumKolejek) ListaKolejek(ctx context.Context,
	stan *shared.QueueStatus) ([]Kolejka, error) {

	zapytanie := listaKolejek
	argumenty := []any{}
	if stan != nil {
		kolumna, err := stanKolejkiNaBaze(*stan)
		if err != nil {
			return nil, err
		}
		zapytanie = listaKolejekStanu
		argumenty = append(argumenty, kolumna)
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, argumenty...)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać wykazu kolejek: %w", err)
	}
	defer wiersze.Close()

	lista := []Kolejka{}
	for wiersze.Next() {
		kolejka, err := odczytajWierszKolejki(wiersze)
		if err != nil {
			return nil, err
		}
		lista = append(lista, kolejka)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt wykazu kolejek: %w", err)
	}
	return lista, nil
}

// odczytajWierszKolejki składa kolejkę z jednego wiersza wyniku.
func odczytajWierszKolejki(wiersz skaner) (Kolejka, error) {
	var kolejka Kolejka
	var sesja, okno sql.NullInt64
	var stan string
	err := wiersz.Scan(&kolejka.ID, &kolejka.Nazwa, &kolejka.Rodzaj, &sesja, &okno, &stan,
		&kolejka.Utworzono, &kolejka.Zaktualizowano)
	if err != nil {
		return Kolejka{}, fmt.Errorf("dane: nieczytelny wiersz kolejki: %w", err)
	}
	kolejka.SesjaID = liczbaZKolumny(sesja)
	kolejka.OknoKoordynatoraID = liczbaZKolumny(okno)
	if kolejka.Stan, err = stanKolejkiZBazy(stan); err != nil {
		return Kolejka{}, err
	}
	return kolejka, nil
}

// PowiazaniaKolejki zwraca byty związane ze wskazaną kolejką.
func (r *repozytoriumKolejek) PowiazaniaKolejki(ctx context.Context,
	kolejkaID int64) ([]PowiazanieKolejki, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaPowiazanKolejki)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, kolejkaID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać powiązań kolejki %d: %w", kolejkaID, err)
	}
	defer wiersze.Close()

	lista := []PowiazanieKolejki{}
	for wiersze.Next() {
		var powiazanie PowiazanieKolejki
		if err := wiersze.Scan(&powiazanie.KolejkaID, &powiazanie.Rodzaj, &powiazanie.Byt,
			&powiazanie.Utworzono); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz powiązania kolejki %d: %w",
				kolejkaID, err)
		}
		lista = append(lista, powiazanie)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt powiązań kolejki %d: %w", kolejkaID, err)
	}
	return lista, nil
}

// ZwiazKolejke zapisuje powiązania kolejki i zostawia po nich ślad w dzienniku
// akcji — tym samym, w którym stoją zmiany stanu. Wszystko idzie w jednej
// transakcji: powiązanie zapisane bez wpisu w dzienniku byłoby zmianą bez
// przyczyny widocznej w przeglądzie kolejki.
//
// Powiązania dokładają się, nie zastępują. Kontrakt `queue.link` zna wyłącznie
// wiązanie; komendy rozwiązującej nie ma, więc zapis, który cicho zdejmuje
// wcześniejsze powiązania, robiłby czynność, o którą nikt nie prosił.
func (r *repozytoriumKolejek) ZwiazKolejke(ctx context.Context, kolejkaID int64,
	powiazania []PowiazanieKolejki) error {

	if len(powiazania) == 0 {
		return nil
	}
	return wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		polecenie, err := r.zapytania.wTransakcji(ctx, transakcja, wstawPowiazanieKolejki)
		if err != nil {
			return err
		}
		opisy := make([]string, 0, len(powiazania))
		for _, powiazanie := range powiazania {
			if _, err := polecenie.ExecContext(ctx, kolejkaID,
				powiazanie.Rodzaj, powiazanie.Byt); err != nil {
				return fmt.Errorf("dane: nie można związać kolejki %d z bytem %s %q: %w",
					kolejkaID, powiazanie.Rodzaj, powiazanie.Byt, err)
			}
			opisy = append(opisy, powiazanie.Rodzaj+"="+powiazanie.Byt)
		}
		szczegoly := strings.Join(opisy, " ")
		return dopiszWpisDziennika(ctx, r.zapytania, transakcja, WpisDziennika{
			KolejkaID: kolejkaID, Akcja: "powiazanie", Szczegoly: &szczegoly,
		})
	})
}
