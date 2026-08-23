// Odpowiedzialność pliku: rola okna komunikacji widziana od strony identyfikatora
// zewnętrznego okna — droga potrzebna rodzinie komend `role.*`.
//
// Rola okna mieszka w kolumnie `okno_komunikacji.rola_okna`, a więź
// koordynator–wykonawca w `okno_komunikacji.okno_koordynatora_id` — tam, gdzie
// kontrakt widzi `Window.windowRole` i `Window.coordinatorWindowId`. Ten plik nie
// zakłada żadnego bytu; dokłada wyłącznie zapis i odczyt tych dwóch kolumn po
// identyfikatorze, którym posługuje się kontrakt.
//
// Osobno od `okna.go`, bo tamten plik czyta i pisze kolumnę `rola_okna` wyłącznie
// jako część pełnego wiersza okna (`Pobierz`/`Aktualizuj`, po identyfikatorze
// wewnętrznym `int64`), więc zapis samej roli musiałby wpierw wczytać całe okno
// wraz z katalogami roboczymi i zapisać je z powrotem — nadpisując po drodze
// pola, o które komenda `role.assign` nie prosi.
//
// Wcielenia roli (`persona`) tu nie ma: nie jest kolumną tego wiersza, tylko
// wpisem w tabeli `ustawienie` na poziomie zasięgu `window`, bo tam zapisuje je
// klient (`multitasking/wcielenia-analizy.ts`, klucz `multitasking.wcielenie`).
// Rdzeń sięga po ten sam adres przez RepozytoriumKonfiguracji, zamiast zakładać
// drugą prawdę o wcieleniu.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"danacoconsole/shared"
)

// RolaOkna to rola okna wraz z więzią, którą rola niesie. Koordynator pusty
// znaczy „okno nie podlega żadnemu koordynatorowi” i jest stanem poprawnym —
// tak samo poprawnym jak więź wypełniona.
type RolaOkna struct {
	Rola shared.WindowRole
	// Koordynator jest identyfikatorem zewnętrznym okna koordynatora, bo tylko
	// takim posługuje się kontrakt. Okno koordynatora bez identyfikatora
	// zewnętrznego oddaje pusty napis — wiersz istnieje, ale kontrakt nie ma jak
	// go nazwać.
	Koordynator string
}

// RepozytoriumRolOkien jest kontraktem zapisu i odczytu roli okna po
// identyfikatorze zewnętrznym.
type RepozytoriumRolOkien interface {
	// RolaOkna zwraca rolę okna wraz z więzią. Brak wiersza wraca jako
	// ErrBrakWiersza — okno żyjące wyłącznie w pamięci rdzenia nie ma tu nic.
	RolaOkna(ctx context.Context, okno string) (RolaOkna, error)
	// ZapiszRoleOkna zapisuje samą rolę. Rola inna niż wykonawca zdejmuje więź
	// koordynatora, bo koordynatora niesie wyłącznie okno wykonawcy.
	ZapiszRoleOkna(ctx context.Context, okno string, rola shared.WindowRole) error
}

const (
	rolaOknaZewnetrznego = `SELECT wykonawca.rola_okna, koordynator.identyfikator_zewnetrzny
	                        FROM okno_komunikacji wykonawca
	                        LEFT JOIN okno_komunikacji koordynator
	                            ON koordynator.id = wykonawca.okno_koordynatora_id
	                        WHERE wykonawca.identyfikator_zewnetrzny = ?`

	zapiszRoleOknaZewnetrznego = `UPDATE okno_komunikacji
	                              SET rola_okna = ?,
	                                  zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                              WHERE identyfikator_zewnetrzny = ?`

	// Zapis roli innej niż wykonawca zdejmuje więź tym samym poleceniem, a nie
	// drugim: rola i więź muszą się zgadzać w każdej chwili, także wtedy,
	// gdy drugie polecenie by nie doszło. Odpowiada to normalizacji roli
	// w pakiecie `session` (rola_okna.go), żeby pamięć i wiersz mówiły to samo.
	zapiszRoleSamodzielnaOkna = `UPDATE okno_komunikacji
	                             SET rola_okna = ?,
	                                 okno_koordynatora_id = NULL,
	                                 zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                             WHERE identyfikator_zewnetrzny = ?`
)

type repozytoriumRolOkien struct {
	zapytania *zapytania
}

func noweRepozytoriumRolOkien(z *zapytania) *repozytoriumRolOkien {
	return &repozytoriumRolOkien{zapytania: z}
}

// RoleOkien oddaje repozytorium ról okien nad tą samą bazą, co reszta zestawu.
//
// Metoda, a nie pole struktury: rola okna nie jest osobnym obszarem danych, tylko
// widokiem na dwie kolumny obszaru okien, który zestaw już niesie polem `Okna`.
// Pole dołożone obok tamtego zapowiadałoby drugie repozytorium tego samego bytu.
func (z *Zestaw) RoleOkien() RepozytoriumRolOkien {
	if z == nil || z.zapytania == nil {
		return nil
	}
	return noweRepozytoriumRolOkien(z.zapytania)
}

// RolaOkna zwraca rolę okna wraz z identyfikatorem zewnętrznym koordynatora.
func (r *repozytoriumRolOkien) RolaOkna(ctx context.Context, okno string) (RolaOkna, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, rolaOknaZewnetrznego)
	if err != nil {
		return RolaOkna{}, err
	}
	var rola string
	var koordynator sql.NullString
	err = polecenie.QueryRowContext(ctx, okno).Scan(&rola, &koordynator)
	if errors.Is(err, sql.ErrNoRows) {
		return RolaOkna{}, fmt.Errorf("%w: okno %q", ErrBrakWiersza, okno)
	}
	if err != nil {
		return RolaOkna{}, fmt.Errorf("dane: nie można odczytać roli okna %q: %w", okno, err)
	}
	wartosc, err := rolaOknaZBazy(rola)
	if err != nil {
		return RolaOkna{}, err
	}
	return RolaOkna{Rola: wartosc, Koordynator: koordynator.String}, nil
}

// ZapiszRoleOkna zapisuje rolę okna wskazanego identyfikatorem zewnętrznym.
//
// Okno bez wiersza wraca jako ErrBrakWiersza, a nie jako cichy brak skutku:
// wołający ma prawo wiedzieć, czy zapis się odbył. Sam rozstrzyga, czy brak
// wiersza jest dla niego awarią, czy stanem normalnym — okno komunikacji
// dostaje wiersz leniwie, przy pierwszej wiadomości.
func (r *repozytoriumRolOkien) ZapiszRoleOkna(ctx context.Context, okno string, rola shared.WindowRole) error {
	kolumna, err := rolaOknaNaBaze(rola)
	if err != nil {
		return err
	}
	zapytanie := zapiszRoleSamodzielnaOkna
	if rola == shared.WindowRoleExecutor {
		zapytanie = zapiszRoleOknaZewnetrznego
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return err
	}
	wynik, err := polecenie.ExecContext(ctx, kolumna, okno)
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać roli okna %q: %w", okno, err)
	}
	trafione, err := wynik.RowsAffected()
	if err != nil {
		return fmt.Errorf("dane: nieznany skutek zapisu roli okna %q: %w", okno, err)
	}
	if trafione == 0 {
		return fmt.Errorf("%w: okno %q", ErrBrakWiersza, okno)
	}
	return nil
}
