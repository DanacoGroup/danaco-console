// Odpowiedzialność pliku: warstwa danych tablicy znaków specjalnych modułu
// Studio — zasady autozamiany skrótu na znak i znaki ostatnio użyte
// (migracja 371).
//
// Tabela `autozamiana_znaku_studio` stoi w migracji 368 wraz z kolumną
// `fabryczna` i zasadami fabrycznymi wpisanymi wierszami; migracja 371 dokłada
// do niej znaki prawnicze i ułamki oraz zakłada `znak_ostatnio_uzyty_studio`.
//
// ── Czego tu NIE MA i dlaczego ───────────────────────────────────────────────
// Nie ma tablicy znaków. Nazwy znaków, ich punkty kodowe i grupy są wiedzą
// rdzenia — tak samo jak arkusz stylów fabryczny i wykaz nośników druku. Do bazy
// schodzi wyłącznie to, co Operator zmienił albo czym się posłużył. Drugi wykaz
// znaków w tabeli rozjechałby się z wykazem rdzenia przy pierwszym uzupełnieniu.
//
// Wykaz zasad autozamiany jest odwrotnie: on stoi W BAZIE, także w części
// fabrycznej, bo migracja 368 tak go założyła i bo Operator ma prawo zasadę
// fabryczną WYŁĄCZYĆ. Powtórzenie wykazu fabrycznego w kodzie rdzenia dałoby
// dwie prawdy o tym, co wchodzi w miejsce „(c)".
//
// ── Przedrostek nazw pomocniczych ────────────────────────────────────────────
// Nazwy pomocnicze tego pliku niosą przedrostek `symbol` — przestrzeń nazw
// pakietu `dane` jest dzielona z innymi wykonawcami.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// ZasadaAutozamianyStudia to wiersz tabeli `autozamiana_znaku_studio`.
//
// Zasada fabryczna ma tu swój wiersz z kolumną `Fabryczna` — wykaz fabryczny
// przyszedł migracją, nie z kodu. Zmiana zasady fabrycznej NIE zdejmuje jej
// oznaczenia: Operator, który wyłączył „--", nadal ma przed sobą zasadę
// fabryczną wyłączoną, a nie zasadę własną, którą wolno usunąć.
type ZasadaAutozamianyStudia struct {
	ID             int64
	Skrot          string
	Zamiennik      string
	Czynna         bool
	Fabryczna      bool
	Utworzono      string
	Zaktualizowano string
}

// ZnakOstatnioUzytyStudia to wiersz tabeli `znak_ostatnio_uzyty_studio`.
type ZnakOstatnioUzytyStudia struct {
	ID      int64
	Kod     string
	Znak    string
	IleUzyc int64
	Uzyto   string
}

const (
	symbolKolumnyZasady = `id, skrot, zamiennik, czynna, fabryczna, utworzono, zaktualizowano`

	// Zapis zasady jest nadpisaniem wiersza, nie założeniem drugiego: skrót jest
	// tożsamością zasady, a dwa wiersze o tym samym skrócie znaczyłyby dwie
	// prawdy o tym, co wchodzi w miejsce „(c)".
	// Zapis NIE rusza kolumny `fabryczna`: oznaczenie zasady przyszło migracją
	// i zmiana jej zamiennika ani wyłączenie nie czynią z niej zasady własnej.
	// Bez tego Operator mógłby usunąć zasadę fabryczną, wpisując ją najpierw
	// jako własną — czyli obejściem.
	symbolZapiszZasade = `INSERT INTO autozamiana_znaku_studio (skrot, zamiennik, czynna, fabryczna)
	                      VALUES (?, ?, ?, 0)
	                      ON CONFLICT(skrot) DO UPDATE SET
	                          zamiennik = excluded.zamiennik,
	                          czynna = excluded.czynna,
	                          zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')`

	symbolPobierzZasade = `SELECT ` + symbolKolumnyZasady + ` FROM autozamiana_znaku_studio
	                       WHERE skrot = ?`

	symbolListaZasad = `SELECT ` + symbolKolumnyZasady + ` FROM autozamiana_znaku_studio
	                    ORDER BY skrot`

	// Usunięcie obejmuje wyłącznie zasadę własną. Zasada fabryczna zostaje
	// i metoda oddaje fałsz — powód odmowy nazywa rdzeń. Warunek stoi w SQL,
	// a nie w rdzeniu, żeby żadna droga wołania go nie ominęła.
	symbolUsunZasade = `DELETE FROM autozamiana_znaku_studio
	                    WHERE skrot = ? AND fabryczna = 0`

	// Odnotowanie użycia podnosi licznik i przestawia czas. Licznik liczy baza,
	// nie wywołujący: dwa okna wstawiające ten sam znak naraz zgubiłyby jedno
	// z użyć, gdyby każde odczytało licznik i zapisało własną sumę.
	symbolOdnotujUzycie = `INSERT INTO znak_ostatnio_uzyty_studio (kod, znak)
	                       VALUES (?, ?)
	                       ON CONFLICT(kod) DO UPDATE SET
	                           znak = excluded.znak,
	                           ile_uzyc = znak_ostatnio_uzyty_studio.ile_uzyc + 1,
	                           uzyto = strftime('%Y-%m-%dT%H:%M:%fZ','now')`

	symbolListaOstatnich = `SELECT id, kod, znak, ile_uzyc, uzyto
	                        FROM znak_ostatnio_uzyty_studio
	                        ORDER BY uzyto DESC, ile_uzyc DESC, id DESC
	                        LIMIT ?`
)

// ZapiszZasadeAutozamiany zakłada zasadę autozamiany albo nadpisuje zastaną
// i oddaje stan po zapisie.
func (r *repozytoriumStudia) ZapiszZasadeAutozamiany(ctx context.Context,
	zasada ZasadaAutozamianyStudia) (ZasadaAutozamianyStudia, error) {

	skrot := strings.TrimSpace(zasada.Skrot)
	if skrot == "" {
		return ZasadaAutozamianyStudia{}, fmt.Errorf("dane: zasada autozamiany bez skrótu")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, symbolZapiszZasade)
	if err != nil {
		return ZasadaAutozamianyStudia{}, err
	}
	czynna := int64(0)
	if zasada.Czynna {
		czynna = 1
	}
	if _, err := polecenie.ExecContext(ctx, skrot, zasada.Zamiennik, czynna); err != nil {
		return ZasadaAutozamianyStudia{}, fmt.Errorf(
			"dane: nie można zapisać zasady autozamiany %q: %w", skrot, err)
	}
	return r.ZasadaAutozamiany(ctx, skrot)
}

// ZasadaAutozamiany oddaje zasadę o wskazanym skrócie.
func (r *repozytoriumStudia) ZasadaAutozamiany(ctx context.Context,
	skrot string) (ZasadaAutozamianyStudia, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, symbolPobierzZasade)
	if err != nil {
		return ZasadaAutozamianyStudia{}, err
	}
	zasada, err := symbolOdczytajZasade(polecenie.QueryRowContext(ctx, strings.TrimSpace(skrot)))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ZasadaAutozamianyStudia{}, ErrBrakWiersza
		}
		return ZasadaAutozamianyStudia{}, fmt.Errorf(
			"dane: nieczytelny wiersz zasady autozamiany %q: %w", skrot, err)
	}
	return zasada, nil
}

// ZasadyAutozamiany oddaje wszystkie zasady — fabryczne i własne Operatora.
func (r *repozytoriumStudia) ZasadyAutozamiany(ctx context.Context) ([]ZasadaAutozamianyStudia, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, symbolListaZasad)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać zasad autozamiany: %w", err)
	}
	defer wiersze.Close()

	lista := []ZasadaAutozamianyStudia{}
	for wiersze.Next() {
		zasada, err := symbolOdczytajZasade(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz zasady autozamiany: %w", err)
		}
		lista = append(lista, zasada)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt zasad autozamiany: %w", err)
	}
	return lista, nil
}

// UsunZasadeAutozamiany usuwa zasadę własną Operatora i oddaje, czy wiersz
// został usunięty. Zasada fabryczna zostaje — patrz zapytanie wyżej.
func (r *repozytoriumStudia) UsunZasadeAutozamiany(ctx context.Context, skrot string) (bool, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, symbolUsunZasade)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, strings.TrimSpace(skrot))
	if err != nil {
		return false, fmt.Errorf("dane: nie można usunąć zasady autozamiany %q: %w", skrot, err)
	}
	usuniete, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nieznany skutek usunięcia zasady autozamiany %q: %w", skrot, err)
	}
	return usuniete > 0, nil
}

// OdnotujUzycieZnaku zapisuje, że Operator posłużył się znakiem — po tym wykaz
// „znaki ostatnio użyte" ma z czego powstać.
func (r *repozytoriumStudia) OdnotujUzycieZnaku(ctx context.Context, kod, znak string) error {
	kod = strings.TrimSpace(strings.ToUpper(kod))
	if kod == "" || znak == "" {
		return fmt.Errorf("dane: odnotowanie użycia znaku bez kodu albo bez znaku")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, symbolOdnotujUzycie)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, kod, znak); err != nil {
		return fmt.Errorf("dane: nie można odnotować użycia znaku %q: %w", kod, err)
	}
	return nil
}

// ZnakiOstatnioUzyte oddaje znaki, którymi Operator posłużył się niedawno,
// od najbliższego ręce.
func (r *repozytoriumStudia) ZnakiOstatnioUzyte(ctx context.Context,
	ile int) ([]ZnakOstatnioUzytyStudia, error) {

	if ile <= 0 {
		ile = 24
	}
	polecenie, err := r.zapytania.przygotuj(ctx, symbolListaOstatnich)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, ile)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać znaków ostatnio użytych: %w", err)
	}
	defer wiersze.Close()

	lista := []ZnakOstatnioUzytyStudia{}
	for wiersze.Next() {
		var znak ZnakOstatnioUzytyStudia
		if err := wiersze.Scan(&znak.ID, &znak.Kod, &znak.Znak, &znak.IleUzyc, &znak.Uzyto); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz znaku ostatnio użytego: %w", err)
		}
		lista = append(lista, znak)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt znaków ostatnio użytych: %w", err)
	}
	return lista, nil
}

// symbolOdczytajZasade składa zasadę autozamiany z jednego wiersza wyniku.
func symbolOdczytajZasade(wiersz skaner) (ZasadaAutozamianyStudia, error) {
	var zasada ZasadaAutozamianyStudia
	var czynna, fabryczna int64
	err := wiersz.Scan(&zasada.ID, &zasada.Skrot, &zasada.Zamiennik, &czynna, &fabryczna,
		&zasada.Utworzono, &zasada.Zaktualizowano)
	if err != nil {
		return ZasadaAutozamianyStudia{}, err
	}
	zasada.Czynna = czynna != 0
	zasada.Fabryczna = fabryczna != 0
	return zasada, nil
}
