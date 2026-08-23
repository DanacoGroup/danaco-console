// Odpowiedzialność pliku: dziennik wykonanych akcji okna (tabela
// `log_akcji_okna`, migracja 052) — treść komendy `window.action`. Katalog
// akcji (`core/akcje.go`) mówi, jakie akcje istnieją; ten plik mówi, kiedy
// i z jakim skutkiem konkretne okno je wykonało, wzorem `log_akcji_kolejki`
// (migracja 003): przejrzystość zamiast bramy. Typ repozytorium i konstruktor
// deklaruje `przekazanie_okna.go`; ten plik dokłada wyłącznie metody dziennika
// akcji.
//
// Parametry i wynik są surowym zapisem, nierozbieranym. Kształt obu pól
// zależy od konkretnej akcji z katalogu, którego warstwa danych akcji nie
// zna — tak samo jak `komplet_kontekstu` w `przekazanie_okna.go` niesie
// `ContextBundle` bez rozbioru w SQL. Tu obie kolumny to zwykły TEXT,
// przenoszony i oddawany jako wskaźnik na napis.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// AkcjaOkna to wiersz tabeli `log_akcji_okna` — jedno wykonanie akcji z
// katalogu przez konkretne okno. `Okno` niesie identyfikator zewnętrzny
// (`okno_komunikacji.identyfikator_zewnetrzny`), ten sam, którym okno
// wychodzi kontraktem jako `Window.Id` / `WindowActionRequest.WindowId` —
// warstwa wyższa nie zna wewnętrznych kluczy liczbowych.
type AkcjaOkna struct {
	ID        int64
	Okno      string
	AkcjaID   string
	Parametry *string
	Wynik     *string
	Utworzono string
}

// limitDomyslnyAkcjiOkna obowiązuje, gdy wywołujący nie poda dodatniego
// limitu — dziennik rośnie z każdym wykonaniem akcji katalogu, więc odczyt
// bez granicy byłby pułapką wydajności przy oknie długo działającym.
const limitDomyslnyAkcjiOkna = 100

const (
	idOknaPoIdentyfikatorzeAkcji = `SELECT id FROM okno_komunikacji WHERE identyfikator_zewnetrzny = ?`

	zapiszAkcjeOkna = `INSERT INTO log_akcji_okna (okno_komunikacji_id, akcja_id, parametry, wynik)
	                   VALUES (?, ?, ?, ?)`

	pobierzAkcjeOkna = `SELECT log.id, ok.identyfikator_zewnetrzny, log.akcja_id, log.parametry,
	                            log.wynik, log.utworzono
	                    FROM log_akcji_okna log
	                    JOIN okno_komunikacji ok ON ok.id = log.okno_komunikacji_id
	                    WHERE ok.identyfikator_zewnetrzny = ?
	                    ORDER BY log.utworzono DESC, log.id DESC
	                    LIMIT ?`
)

// ZapiszAkcje odnotowuje wykonanie akcji katalogu przez wskazane okno i
// zwraca zapisany wiersz — obsługuje `window.action`. Okno nieznane wraca
// jako ErrBrakWiersza: dziennik nie zakłada wpisu dla okna, którego nie ma.
func (r *repozytoriumPrzekazan) ZapiszAkcje(ctx context.Context, akcja AkcjaOkna) (AkcjaOkna, error) {
	if akcja.Okno == "" {
		return AkcjaOkna{}, fmt.Errorf("dane: akcja okna bez identyfikatora okna")
	}
	if akcja.AkcjaID == "" {
		return AkcjaOkna{}, fmt.Errorf("dane: akcja okna bez identyfikatora akcji")
	}

	szukanieOkna, err := r.zapytania.przygotuj(ctx, idOknaPoIdentyfikatorzeAkcji)
	if err != nil {
		return AkcjaOkna{}, err
	}
	var oknoID int64
	err = szukanieOkna.QueryRowContext(ctx, akcja.Okno).Scan(&oknoID)
	if errors.Is(err, sql.ErrNoRows) {
		return AkcjaOkna{}, fmt.Errorf("dane: okno %q nie istnieje: %w", akcja.Okno, ErrBrakWiersza)
	}
	if err != nil {
		return AkcjaOkna{}, fmt.Errorf("dane: nie można odnaleźć okna %q: %w", akcja.Okno, err)
	}

	wstawianie, err := r.zapytania.przygotuj(ctx, zapiszAkcjeOkna)
	if err != nil {
		return AkcjaOkna{}, err
	}
	wynik, err := wstawianie.ExecContext(ctx, oknoID, akcja.AkcjaID,
		tekstDoKolumny(akcja.Parametry), tekstDoKolumny(akcja.Wynik))
	if err != nil {
		return AkcjaOkna{}, fmt.Errorf("dane: nie można zapisać akcji %q okna %q: %w",
			akcja.AkcjaID, akcja.Okno, err)
	}
	nowyID, err := wynik.LastInsertId()
	if err != nil {
		return AkcjaOkna{}, fmt.Errorf("dane: nie można odczytać id akcji okna: %w", err)
	}

	zapisana := akcja
	zapisana.ID = nowyID
	return zapisana, nil
}

// AkcjeOkna zwraca ostatnie akcje wskazanego okna, od najnowszej — zasila
// panel akcji okna i Mission Control. Limit niedodatni wraca do wartości
// domyślnej zamiast oznaczać „bez granicy": dziennik rośnie bezterminowo.
func (r *repozytoriumPrzekazan) AkcjeOkna(ctx context.Context, okno string, limit int) ([]AkcjaOkna, error) {
	if limit <= 0 {
		limit = limitDomyslnyAkcjiOkna
	}

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzAkcjeOkna)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno, limit)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać akcji okna %q: %w", okno, err)
	}
	defer wiersze.Close()

	lista := []AkcjaOkna{}
	for wiersze.Next() {
		akcja, err := odczytajAkcjeOkna(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz akcji okna %q: %w", okno, err)
		}
		lista = append(lista, akcja)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt akcji okna %q: %w", okno, err)
	}
	return lista, nil
}

// odczytajAkcjeOkna składa strukturę z jednego wiersza wyniku złączenia
// `log_akcji_okna` z `okno_komunikacji`.
func odczytajAkcjeOkna(wiersz skaner) (AkcjaOkna, error) {
	var akcja AkcjaOkna
	var parametry, wynikTekst sql.NullString
	err := wiersz.Scan(&akcja.ID, &akcja.Okno, &akcja.AkcjaID, &parametry, &wynikTekst, &akcja.Utworzono)
	if err != nil {
		return AkcjaOkna{}, err
	}
	akcja.Parametry = tekstZKolumny(parametry)
	akcja.Wynik = tekstZKolumny(wynikTekst)
	return akcja, nil
}
