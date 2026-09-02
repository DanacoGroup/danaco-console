// Układ sekcji panelu okna (kolejność, zwinięcie, zdjęcie) zapisywany całościowo w sekcja_panelu dla pary okno+panel.
package dane

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"
)

// SekcjaPanelu to jedna sekcja panelu od strony układu; nazw ani treści nie zna — należą do warstwy widoku.
type SekcjaPanelu struct {
	Id        string
	Kolejnosc int
	Zwinieta  bool
	Zdjeta    bool
}

// RepozytoriumSekcjiPaneli jest kontraktem odczytu i zapisu układu sekcji jednego panelu jednego okna.
type RepozytoriumSekcjiPaneli interface {
	SekcjePanelu(ctx context.Context, okno, panel string) ([]SekcjaPanelu, error)
	ZapiszSekcje(ctx context.Context, okno, panel string, sekcje []SekcjaPanelu) ([]SekcjaPanelu, error)
}

const (
	sekcjePaneluWykaz = `SELECT sekcja_id, kolejnosc, zwinieta, zdjeta
	                     FROM sekcja_panelu
	                     WHERE okno_id = ? AND panel_id = ? AND ` + WarunekKonta + `
	                     ORDER BY kolejnosc, sekcja_id`

	sekcjePaneluCzyszczenie = `DELETE FROM sekcja_panelu
	                           WHERE okno_id = ? AND panel_id = ? AND ` + WarunekKonta

	sekcjaPaneluZapis = `INSERT INTO sekcja_panelu
	                         (okno_id, panel_id, sekcja_id, kolejnosc, zwinieta, zdjeta, zaktualizowano, konto_id)
	                     VALUES (?, ?, ?, ?, ?, ?, strftime('%Y-%m-%dT%H:%M:%fZ','now'), ` + WskazanieKonta + `)`
)

type repozytoriumSekcjiPaneli struct {
	zapytania *zapytania
	db        *sql.DB
}

func noweRepozytoriumSekcjiPaneli(z *zapytania, db *sql.DB) *repozytoriumSekcjiPaneli {
	return &repozytoriumSekcjiPaneli{zapytania: z, db: db}
}

func (z *Zestaw) SekcjePaneli() RepozytoriumSekcjiPaneli {
	if z == nil || z.zapytania == nil || z.zapytania.db == nil {
		return nil
	}
	return noweRepozytoriumSekcjiPaneli(z.zapytania, z.zapytania.db)
}

func (r *repozytoriumSekcjiPaneli) SekcjePanelu(ctx context.Context,
	okno, panel string) ([]SekcjaPanelu, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, sekcjePaneluWykaz)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno, panel, KontoOperatora(ctx))
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

// ZapiszSekcje zastępuje układ panelu w całości; wykaz pusty znaczy „powrót do układu domyślnego".
func (r *repozytoriumSekcjiPaneli) ZapiszSekcje(ctx context.Context, okno, panel string,
	sekcje []SekcjaPanelu) ([]SekcjaPanelu, error) {

	uklad := make([]SekcjaPanelu, 0, len(sekcje))
	for numer, sekcja := range sekcje {
		sekcja.Kolejnosc = numer + 1
		uklad = append(uklad, sekcja)
	}

	adres := okno + "\x00" + panel
	zapisyPaneli.wejdz(adres)

	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		czyszczenie, err := r.zapytania.wTransakcji(ctx, transakcja, sekcjePaneluCzyszczenie)
		if err != nil {
			return err
		}
		if _, err := czyszczenie.ExecContext(ctx, okno, panel, KontoOperatora(ctx)); err != nil {
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
				sekcja.Kolejnosc, sekcja.Zwinieta, sekcja.Zdjeta, KontoOperatora(ctx))
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
	zapisyPaneli.poczekaj(adres)
	return r.SekcjePanelu(ctx, okno, panel)
}

// zapisyPaneli liczy zapisy w toku na adres okno+panel, by odczyt zaczekał na opadnięcie zapisów zbieżnych.
var zapisyPaneli = licznikZapisowPaneli{wToku: map[string]int{}}

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
