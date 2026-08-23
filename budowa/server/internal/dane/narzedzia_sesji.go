// Odpowiedzialność pliku: dołożenia narzędzi żyjące w stanie sesji — zapis,
// odczyt i zdjęcie (tabela `narzedzie_sesji`).
//
// Narzędzie dołożone komendą po ukośniku trwa do końca sesji: nie wchodzi na
// stałe do definicji eksperta i nie znika po jednej turze. Czas życia niesie
// klucz obcy z kasowaniem kaskadowym, nie kod — dlatego nie ma tu metody
// sprzątającej wygasłe wiersze; koniec życia dołożenia to koniec życia wiersza
// `sesja`.
//
// Plik prowadzi dołożenia jednej sesji, a nie katalog, z którego się je wybiera.
// Wykaz pozycji po ukośniku nie jest tabelą: składa go rdzeń na bieżąco z komend
// kontraktu, katalogu akcji (tabela `akcja`) i katalogu rozszerzeń (tabela
// `rozszerzenie`) — `core/adapter_narzedzia_sesji.go`. Odpisanie go do trzeciej
// tabeli byłoby drugą prawdą o tym, co platforma umie, i rozjechałoby się
// z pierwszą przy pierwszej instalacji rozszerzenia.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// NarzedzieSesji to jedno dołożenie: pozycja wykazu odpisana w chwili dołożenia.
//
// Wiersz niesie odpis, nie odwołanie: trzy pola nazewnicze i grupę, a nie klucz
// obcy do pozycji wykazu — wykaz nie jest tabelą i nie ma czego wskazać. Dzięki
// temu odinstalowanie rozszerzenia nie zamienia dołożenia w nazwę bez opisu
// i do końca sesji widać, co model dostał.
type NarzedzieSesji struct {
	ID int64
	// NazwaPelna niesie przedrostek źródła, np. `anthropic-skills:skill-creator`.
	// To ona jest tożsamością dołożenia w obrębie sesji.
	NazwaPelna string
	// NazwaSkrocona to nazwa bez przedrostka źródła — ta, którą się wpisuje.
	NazwaSkrocona string
	// Opis pełnym zdaniem: po co to jest i kiedy użyć.
	Opis string
	// Rodzaj niesie wartość kontraktu wprost (shared.SlashEntryKind).
	Rodzaj string
	// Grupa po przeznaczeniu — ta sama, którą niesie wykaz.
	Grupa string
	// ZrodloPozycji to przedrostek nazwy pełnej: skąd pozycja pochodzi.
	ZrodloPozycji string
	// Zrodlo niesie wartość kontraktu wprost (shared.SessionToolSource): czyja
	// ręka dołożyła.
	Zrodlo string
	// Dolozono to milisekundy epoki; podaje je warstwa wyższa, nie baza.
	Dolozono int64
}

// RepozytoriumNarzedziSesji jest kontraktem dołożeń jednej sesji.
//
// Identyfikator sesji jest tu kluczem wiersza `sesja`, nie identyfikatorem
// kontraktowym: przekład jednego na drugi należy do adaptera rdzenia, tak samo
// jak przy każdym innym repozytorium tego pakietu.
type RepozytoriumNarzedziSesji interface {
	// Narzedzia oddaje dołożenia sesji w kolejności dokładania. Sesja bez
	// dołożeń oddaje wykaz pusty i jest to stan poprawny, nie brak wiersza —
	// zestaw narzędzi tury jest wtedy samą definicją eksperta.
	Narzedzia(ctx context.Context, sesjaID int64) ([]NarzedzieSesji, error)
	// Doloz zapisuje dołożenie i oddaje je wraz z nadanym identyfikatorem oraz
	// znacznikiem, czy narzędzie było już dołożone wcześniej. Powtórzenie nie
	// jest błędem i nie mnoży wierszy.
	Doloz(ctx context.Context, sesjaID int64, narzedzie NarzedzieSesji) (NarzedzieSesji, bool, error)
	// Zdejmij kasuje dołożenie po nazwie pełnej. Falsz znaczy „nie było czego
	// zdejmować" i też nie jest błędem.
	Zdejmij(ctx context.Context, sesjaID int64, nazwaPelna string) (bool, error)
	// ZdejmijWszystkie kasuje wszystkie dołożenia sesji i oddaje ich liczbę.
	// Usunięcie sesji zdejmuje dołożenia kaskadą klucza obcego; ta czynność
	// obsługuje przypadek, w którym sesja zostaje, a zestaw ma wrócić do
	// podstawy.
	ZdejmijWszystkie(ctx context.Context, sesjaID int64) (int, error)
}

const (
	kolumnyNarzedziaSesji = `id, nazwa_pelna, nazwa_skrocona, opis, rodzaj, grupa,
	                         zrodlo_pozycji, zrodlo, dolozono`

	wykazNarzedziSesji = `SELECT ` + kolumnyNarzedziaSesji + `
	                      FROM narzedzie_sesji
	                      WHERE sesja_id = ?
	                      ORDER BY dolozono, id`

	narzedzieSesjiPoNazwie = `SELECT ` + kolumnyNarzedziaSesji + `
	                          FROM narzedzie_sesji
	                          WHERE sesja_id = ? AND nazwa_pelna = ?`

	wstawNarzedzieSesji = `INSERT INTO narzedzie_sesji
	                           (sesja_id, nazwa_pelna, nazwa_skrocona, opis, rodzaj,
	                            grupa, zrodlo_pozycji, zrodlo, dolozono)
	                       VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	skasujNarzedzieSesji = `DELETE FROM narzedzie_sesji
	                        WHERE sesja_id = ? AND nazwa_pelna = ?`

	skasujNarzedziaSesji = `DELETE FROM narzedzie_sesji WHERE sesja_id = ?`
)

type repozytoriumNarzedziSesji struct {
	zapytania *zapytania
	db        *sql.DB
}

func noweRepozytoriumNarzedziSesji(z *zapytania, db *sql.DB) *repozytoriumNarzedziSesji {
	return &repozytoriumNarzedziSesji{zapytania: z, db: db}
}

// NarzedziaSesji oddaje repozytorium dołożeń nad tą samą bazą, co reszta
// zestawu.
//
// Metoda, a nie pole struktury — z tego samego powodu, co `SekcjePaneli`
// (`panele.go`) i `Rozszerzenia` (`extension.go`): repozytorium nie trzyma stanu
// poza wskaźnikiem na wspólną pamięć zapytań.
func (z *Zestaw) NarzedziaSesji() RepozytoriumNarzedziSesji {
	if z == nil || z.zapytania == nil || z.zapytania.db == nil {
		return nil
	}
	return noweRepozytoriumNarzedziSesji(z.zapytania, z.zapytania.db)
}

// Narzedzia czyta dołożenia jednej sesji w kolejności dokładania.
func (r *repozytoriumNarzedziSesji) Narzedzia(ctx context.Context,
	sesjaID int64) ([]NarzedzieSesji, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, wykazNarzedziSesji)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, sesjaID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać dołożeń sesji %d: %w", sesjaID, err)
	}
	defer wiersze.Close()

	wykaz := make([]NarzedzieSesji, 0, 8)
	for wiersze.Next() {
		narzedzie, err := odczytajNarzedzieSesji(wiersze.Scan)
		if err != nil {
			return nil, fmt.Errorf("dane: uszkodzone dołożenie sesji %d: %w", sesjaID, err)
		}
		wykaz = append(wykaz, narzedzie)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt dołożeń sesji %d: %w", sesjaID, err)
	}
	return wykaz, nil
}

// Doloz zapisuje dołożenie. Odczyt i zapis idą jedną transakcją, bo razem
// odpowiadają na jedno pytanie: „czy to już jest, a jeśli nie — wpisz". Bez
// wspólnej transakcji dwie komendy po ukośniku wydane w tej samej chwili obie
// zastałyby pustą tabelę i obie próbowałyby wpisać wiersz; druga odbiłaby się
// o warunek UNIQUE i dostałaby odmowę za czynność, która była poprawna.
func (r *repozytoriumNarzedziSesji) Doloz(ctx context.Context, sesjaID int64,
	narzedzie NarzedzieSesji) (NarzedzieSesji, bool, error) {

	var zapisane NarzedzieSesji
	juzBylo := false
	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		odczyt, err := r.zapytania.wTransakcji(ctx, transakcja, narzedzieSesjiPoNazwie)
		if err != nil {
			return err
		}
		zastane, err := odczytajNarzedzieSesji(
			odczyt.QueryRowContext(ctx, sesjaID, narzedzie.NazwaPelna).Scan)
		if err == nil {
			// Powtórzenie nie jest błędem i nie nadpisuje zastanego wiersza.
			// Czas dołożenia ma mówić, kiedy model dostał narzędzie — odświeżony
			// przy każdym powtórzeniu kłamałby o chwili, od której je ma.
			zapisane, juzBylo = zastane, true
			return nil
		}
		if !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("dane: nie można sprawdzić dołożenia %q sesji %d: %w",
				narzedzie.NazwaPelna, sesjaID, err)
		}
		zapis, err := r.zapytania.wTransakcji(ctx, transakcja, wstawNarzedzieSesji)
		if err != nil {
			return err
		}
		wynik, err := zapis.ExecContext(ctx, sesjaID, narzedzie.NazwaPelna,
			narzedzie.NazwaSkrocona, narzedzie.Opis, narzedzie.Rodzaj, narzedzie.Grupa,
			narzedzie.ZrodloPozycji, narzedzie.Zrodlo, narzedzie.Dolozono)
		if err != nil {
			return fmt.Errorf("dane: nie można dołożyć %q do sesji %d: %w",
				narzedzie.NazwaPelna, sesjaID, err)
		}
		id, err := wynik.LastInsertId()
		if err != nil {
			return fmt.Errorf("dane: nieznany identyfikator dołożenia %q sesji %d: %w",
				narzedzie.NazwaPelna, sesjaID, err)
		}
		zapisane = narzedzie
		zapisane.ID = id
		return nil
	})
	if err != nil {
		return NarzedzieSesji{}, false, err
	}
	return zapisane, juzBylo, nil
}

// Zdejmij kasuje dołożenie po nazwie pełnej.
func (r *repozytoriumNarzedziSesji) Zdejmij(ctx context.Context, sesjaID int64,
	nazwaPelna string) (bool, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, skasujNarzedzieSesji)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, sesjaID, nazwaPelna)
	if err != nil {
		return false, fmt.Errorf("dane: nie można zdjąć dołożenia %q sesji %d: %w",
			nazwaPelna, sesjaID, err)
	}
	zmienione, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nieznany skutek zdjęcia dołożenia %q sesji %d: %w",
			nazwaPelna, sesjaID, err)
	}
	return zmienione > 0, nil
}

// ZdejmijWszystkie kasuje wszystkie dołożenia sesji.
func (r *repozytoriumNarzedziSesji) ZdejmijWszystkie(ctx context.Context,
	sesjaID int64) (int, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, skasujNarzedziaSesji)
	if err != nil {
		return 0, err
	}
	wynik, err := polecenie.ExecContext(ctx, sesjaID)
	if err != nil {
		return 0, fmt.Errorf("dane: nie można zdjąć dołożeń sesji %d: %w", sesjaID, err)
	}
	zmienione, err := wynik.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("dane: nieznany skutek zdjęcia dołożeń sesji %d: %w", sesjaID, err)
	}
	return int(zmienione), nil
}

// odczytajNarzedzieSesji składa wiersz z jednego skanu. Kolejność kolumn jest
// jedna (`kolumnyNarzedziaSesji`) i obsługuje oba odczyty — wykaz i pojedynczy
// wiersz — więc dopisanie kolumny nie rozjedzie jednego z nich.
func odczytajNarzedzieSesji(skan func(...any) error) (NarzedzieSesji, error) {
	var n NarzedzieSesji
	err := skan(&n.ID, &n.NazwaPelna, &n.NazwaSkrocona, &n.Opis, &n.Rodzaj, &n.Grupa,
		&n.ZrodloPozycji, &n.Zrodlo, &n.Dolozono)
	if err != nil {
		return NarzedzieSesji{}, err
	}
	return n, nil
}
