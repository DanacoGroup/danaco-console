// Odpowiedzialność pliku: układ sekcji panelu okna, czyli kolejność, zwinięcie
// i zdjęcie z widoku, jest zapisywany i odczytywany całościowo w tabeli
// sekcja_panelu dla pary okno i panel.
package dane

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"
)

// SekcjaPanelu to jedna sekcja panelu widziana od strony układu: gdzie stoi,
// czy jest zwinięta i czy w ogóle jest na widoku. Nazwy ani treści sekcji ten
// byt nie zna — należą do warstwy widoku, która sekcje rysuje.
type SekcjaPanelu struct {
	// Id jest identyfikatorem sekcji w obrębie panelu; nadaje go widok.
	Id string
	// Kolejnosc to miejsce na widoku, liczone od 1.
	Kolejnosc int
	// Zwinieta znaczy „sekcja jest na widoku, zawinięta do nagłówka”.
	Zwinieta bool
	// Zdjeta znaczy „sekcji na widoku nie ma wcale”, niezależnie od zwinięcia.
	Zdjeta bool
}

// RepozytoriumSekcjiPaneli jest kontraktem odczytu i zapisu układu sekcji
// jednego panelu jednego okna.
type RepozytoriumSekcjiPaneli interface {
	// SekcjePanelu oddaje układ w kolejności widoku; panel nigdy nieustawiany
	// oddaje wykaz pusty.
	SekcjePanelu(ctx context.Context, okno, panel string) ([]SekcjaPanelu, error)
	// ZapiszSekcje zastępuje układ panelu w całości i oddaje układ obowiązujący
	// z numeracją 1..N.
	ZapiszSekcje(ctx context.Context, okno, panel string, sekcje []SekcjaPanelu) ([]SekcjaPanelu, error)
}

const (
	sekcjePaneluWykaz = `SELECT sekcja_id, kolejnosc, zwinieta, zdjeta
	                     FROM sekcja_panelu
	                     WHERE okno_id = ? AND panel_id = ?
	                     ORDER BY kolejnosc, sekcja_id`

	sekcjePaneluCzyszczenie = `DELETE FROM sekcja_panelu
	                           WHERE okno_id = ? AND panel_id = ?`

	sekcjaPaneluZapis = `INSERT INTO sekcja_panelu
	                         (okno_id, panel_id, sekcja_id, kolejnosc, zwinieta, zdjeta, zaktualizowano)
	                     VALUES (?, ?, ?, ?, ?, ?, strftime('%Y-%m-%dT%H:%M:%fZ','now'))`
)

type repozytoriumSekcjiPaneli struct {
	zapytania *zapytania
	db        *sql.DB
}

func noweRepozytoriumSekcjiPaneli(z *zapytania, db *sql.DB) *repozytoriumSekcjiPaneli {
	return &repozytoriumSekcjiPaneli{zapytania: z, db: db}
}

// SekcjePaneli oddaje repozytorium układów paneli nad tą samą bazą, co reszta
// zestawu. Jest metodą, a nie polem struktury — wzorem `RoleOkien` — bo
// repozytorium nie trzyma stanu poza wskaźnikiem na wspólną pamięć zapytań.
func (z *Zestaw) SekcjePaneli() RepozytoriumSekcjiPaneli {
	if z == nil || z.zapytania == nil || z.zapytania.db == nil {
		return nil
	}
	return noweRepozytoriumSekcjiPaneli(z.zapytania, z.zapytania.db)
}

// SekcjePanelu czyta bieżący układ sekcji jednego panelu wskazanego okna z tabeli sekcja_panelu, w kolejności widoku.
func (r *repozytoriumSekcjiPaneli) SekcjePanelu(ctx context.Context,
	okno, panel string) ([]SekcjaPanelu, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, sekcjePaneluWykaz)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno, panel)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać układu panelu %q okna %q: %w", panel, okno, err)
	}
	defer wiersze.Close()

	uklad := make([]SekcjaPanelu, 0, 8)
	for wiersze.Next() {
		var sekcja SekcjaPanelu
		if err := wiersze.Scan(&sekcja.Id, &sekcja.Kolejnosc, &sekcja.Zwinieta, &sekcja.Zdjeta); err != nil {
			return nil, fmt.Errorf("dane: nie można odczytać sekcji panelu %q okna %q: %w", panel, okno, err)
		}
		uklad = append(uklad, sekcja)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt układu panelu %q okna %q: %w", panel, okno, err)
	}
	return uklad, nil
}

// ZapiszSekcje zastępuje układ panelu w całości. Wykaz pusty jest poleceniem
// poprawnym i znaczy „panel wraca do układu domyślnego”: stary układ znika,
// a nowego nie ma.
func (r *repozytoriumSekcjiPaneli) ZapiszSekcje(ctx context.Context, okno, panel string,
	sekcje []SekcjaPanelu) ([]SekcjaPanelu, error) {

	uklad := make([]SekcjaPanelu, 0, len(sekcje))
	for numer, sekcja := range sekcje {
		sekcja.Kolejnosc = numer + 1
		uklad = append(uklad, sekcja)
	}

	adres := okno + "\x00" + panel
	zapisyPaneli.wejdz(adres)

	// Układ oddawany czyta się z bazy, nie z żądania, bo klient bierze odpowiedź za dowód skutku zapisu.
	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		czyszczenie, err := r.zapytania.wTransakcji(ctx, transakcja, sekcjePaneluCzyszczenie)
		if err != nil {
			return err
		}
		if _, err := czyszczenie.ExecContext(ctx, okno, panel); err != nil {
			return fmt.Errorf("dane: nie można zdjąć układu panelu %q okna %q: %w", panel, okno, err)
		}
		if len(uklad) == 0 {
			return nil
		}
		zapis, err := r.zapytania.wTransakcji(ctx, transakcja, sekcjaPaneluZapis)
		if err != nil {
			return err
		}
		for _, sekcja := range uklad {
			_, err := zapis.ExecContext(ctx, okno, panel, sekcja.Id,
				sekcja.Kolejnosc, sekcja.Zwinieta, sekcja.Zdjeta)
			if err != nil {
				return fmt.Errorf("dane: nie można zapisać sekcji %q panelu %q okna %q: %w",
					sekcja.Id, panel, okno, err)
			}
		}
		return nil
	})
	zapisyPaneli.wyjdz(adres)
	if err != nil {
		return nil, err
	}
	// Odczyt idzie po zamknięciu transakcji i opadnięciu zapisów zbiegłych, by oddać stan już ustalony.
	zapisyPaneli.poczekaj(adres)
	return r.SekcjePanelu(ctx, okno, panel)
}

// ── zbieg zapisów jednego panelu ─────────────────────────────────────────────

// zapisyPaneli liczy zapisy w toku dla każdego adresu okno i panel, aby odczyt
// zaczekał na opadnięcie zapisów zbieżnych w czasie, do kresu
// kresCzekaniaPaneli, i oddał układ już ustalony.
var zapisyPaneli = licznikZapisowPaneli{wToku: map[string]int{}}

// kresCzekaniaPaneli ogranicza czas oczekiwania odczytu na opadnięcie zapisów zbiegłych w czasie dla tego samego panelu.
const kresCzekaniaPaneli = 250 * time.Millisecond

type licznikZapisowPaneli struct {
	zamek sync.Mutex
	wToku map[string]int
}

func (l *licznikZapisowPaneli) wejdz(adres string) {
	l.zamek.Lock()
	l.wToku[adres]++
	l.zamek.Unlock()
}

func (l *licznikZapisowPaneli) wyjdz(adres string) {
	l.zamek.Lock()
	if l.wToku[adres] <= 1 {
		delete(l.wToku, adres)
	} else {
		l.wToku[adres]--
	}
	l.zamek.Unlock()
}

// poczekaj wstrzymuje odczyt, dopóki trwa cudzy zapis tego samego panelu, najwyżej do upływu kresu kresCzekaniaPaneli.
func (l *licznikZapisowPaneli) poczekaj(adres string) {
	koniec := time.Now().Add(kresCzekaniaPaneli)
	for {
		l.zamek.Lock()
		wolne := l.wToku[adres] == 0
		l.zamek.Unlock()
		if wolne || time.Now().After(koniec) {
			return
		}
		time.Sleep(time.Millisecond)
	}
}
