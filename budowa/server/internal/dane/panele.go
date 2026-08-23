// Odpowiedzialność pliku: układ sekcji panelu okna — kolejność, zwinięcie
// i zdjęcie z widoku — zapisany i odczytany po parze (okno, panel). Nośnikiem
// jest tabela `sekcja_panelu`; drugiego miejsca ten układ nie ma.
//
// Zapis jest całościowy, nie różnicowy: podane sekcje wyznaczają układ panelu,
// a czego w żądaniu nie ma, tego po zapisie nie ma w bazie. Wynika to wprost
// z kontraktu — `panel.sections.set` niesie samo pole `sections`, bez znacznika
// czynności, więc jedynym czytelnym znaczeniem listy jest układ docelowy. Zapis
// różnicowy wymagałby, żeby klient wiedział, co w bazie leży teraz; wtedy dwa
// okna przestawiające ten sam panel rozjechałyby układ, bo każde dopisywałoby
// swoje do cudzego stanu.
//
// Skasowanie starego układu i wpisanie nowego idzie jedną transakcją: przerwane
// w połowie zostawiłyby panel bez sekcji, nieodróżnialny od panelu nigdy
// nieustawianego.
//
// Kolejność nadaje ten plik, nie wołający. Baza pilnuje wyłącznie dolnej granicy
// (CHECK kolejnosc >= 1), bo liczba sekcji panelu jest znana dopiero w chwili
// zapisu. Numery 1..N nanosi `ZapiszSekcje`, żeby wykaz czytany po `kolejnosc`
// był tym samym, co wykaz czytany po miejscu na liście.
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
	// Zdjeta znaczy „sekcji na widoku nie ma wcale”. Stan niezależny od zwinięcia:
	// sekcja zdjęta rozwinięta wraca na widok rozwinięta.
	Zdjeta bool
}

// RepozytoriumSekcjiPaneli jest kontraktem odczytu i zapisu układu sekcji
// jednego panelu jednego okna.
type RepozytoriumSekcjiPaneli interface {
	// SekcjePanelu oddaje układ w kolejności widoku. Panel nigdy nieustawiany
	// oddaje wykaz pusty; układ domyślny należy do widoku, a nie do bazy.
	SekcjePanelu(ctx context.Context, okno, panel string) ([]SekcjaPanelu, error)
	// ZapiszSekcje zastępuje układ panelu w całości i oddaje układ obowiązujący
	// wraz z nadaną numeracją 1..N.
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

// SekcjePanelu czyta układ jednego panelu jednego okna.
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

	// Układ oddawany czyta się z bazy, nie z żądania. Tablica `uklad` opisuje
	// wyłącznie treść żądania, a klient bierze odpowiedź za dowód skutku
	// (`client/src/powloka/zrodlo-sekcji-paneli.ts`), więc dowód musi pochodzić
	// z nośnika.
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
	// Odczyt idzie po zamknięciu transakcji i po opadnięciu zapisów zbiegłych
	// w czasie: układ obowiązujący to stan, który zastanie następny czytelnik.
	// Dwa okna przestawiające ten sam panel dostają dzięki temu tę samą treść,
	// więc odpowiedź niepodobna do żądania znaczy tyle, że cudzy zapis wszedł
	// po naszym.
	zapisyPaneli.poczekaj(adres)
	return r.SekcjePanelu(ctx, okno, panel)
}

// ── zbieg zapisów jednego panelu ─────────────────────────────────────────────

// zapisyPaneli liczy zapisy będące w toku dla każdego adresu (okno, panel).
//
// Transakcja pilnuje całościowości zapisu, ale nie tego, by odczyt kontrolny
// zastał stan już ustalony. Licznik pozwala odczytać układ dopiero wtedy, gdy
// zapisy zbiegłe w czasie opadły. Czekanie ma kres `kresCzekaniaPaneli`: panel
// przestawiany bez ustanku nie może wstrzymać odpowiedzi w nieskończoność,
// a odczyt po upływie kresu jest nadal odczytem z nośnika.
var zapisyPaneli = licznikZapisowPaneli{wToku: map[string]int{}}

// kresCzekaniaPaneli ogranicza czekanie na opadnięcie zapisów zbiegłych.
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

// poczekaj wstrzymuje odczyt, dopóki trwa cudzy zapis tego samego panelu.
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
