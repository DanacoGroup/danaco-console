// Odpowiedzialność pliku: karty przeglądania, grupy kart i przestrzenie robocze
// (tabele `karta_przegladania`, `grupa_kart_przegladania`,
// `przestrzen_przegladania`, migracja 170) — trwałość rzędu kart Browser Window
// i przełącznika przestrzeni.
//
// Karta zamknięta zostaje w tabeli, a odsiewa ją kolumna `zamknieta`. Skasowanie
// wiersza zabrałoby przestrzeni roboczej to, co zapamiętała: przestrzeń wskazuje
// karty, więc karta usunięta z bazy zamieniłaby zapisany zestaw w zestaw
// dziurawy, o którym nikt by się nie dowiedział aż do jego otwarcia.
//
// Przynależność do grupy i do przestrzeni stoi po stronie karty. Karta należy do
// jednej grupy i jednej przestrzeni naraz, więc skład jednej i drugiej jest
// zapytaniem po kolumnie — a nie drugą listą, którą trzeba by prostować przy
// każdym zamknięciu karty.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// KartaPrzegladania to wiersz tabeli `karta_przegladania` — jedna karta rzędu
// kart okna przeglądarki.
type KartaPrzegladania struct {
	ID               int64
	Kod              string
	Okno             string
	Url              *string
	Tytul            *string
	Stan             string
	Przypieta        bool
	Zamknieta        bool
	Grupa            *string
	Przestrzen       *string
	KartaOtwierajaca *string
	Kolejnosc        int64
	OstatnioCzynna   *string
	Utworzono        string
}

// GrupaKart to wiersz tabeli `grupa_kart_przegladania` — nazwany, kolorowany
// zestaw kart zwijany jednym kliknięciem.
type GrupaKart struct {
	ID        int64
	Kod       string
	Okno      string
	Nazwa     string
	Barwa     *string
	Zwinieta  bool
	Utworzono string
}

// PrzestrzenPrzegladania to wiersz tabeli `przestrzen_przegladania` — zapisany
// zestaw kart przełączany bez utraty stanu.
type PrzestrzenPrzegladania struct {
	ID             int64
	Kod            string
	Okno           *string
	Nazwa          string
	Profil         *string
	Utworzono      string
	Zaktualizowano string
}

// FiltrKartPrzegladania zawęża wykaz kart okna — obsługuje pola żądania
// `browser.tab.list`.
type FiltrKartPrzegladania struct {
	// Okno operacyjne, którego rząd kart jest odczytywany.
	Okno string
	// Przestrzen zawęża wykaz do kart jednej przestrzeni roboczej; pusta
	// wartość znaczy „wszystkie karty okna", nie „karty bez przestrzeni".
	Przestrzen string
	// ZZawieszonymi dopuszcza karty uśpione. Karta zawieszona nadal jest
	// otwarta, więc domyślnie wchodzi do wykazu; pole zostawione tu po to,
	// żeby żądanie mogło wykaz zawęzić do kart żywych.
	ZZawieszonymi bool
	// ZZamknietymi dopuszcza karty zamknięte. Rząd kart ich nie pokazuje, ale
	// przestrzeń robocza pamięta także te, które Operator zamknął — bez tego
	// przywrócenie zapisanego zestawu oddawałoby zestaw okrojony i nikt by się
	// o tym nie dowiedział.
	ZZamknietymi bool
}

const (
	kolumnyKartyPrzegladania = `id, identyfikator_zewnetrzny, okno, url, tytul, stan,
	                            przypieta, zamknieta, grupa, przestrzen, karta_otwierajaca,
	                            kolejnosc, ostatnio_czynna, utworzono`

	zapiszKartePrzegladania = `INSERT INTO karta_przegladania
	                           (identyfikator_zewnetrzny, okno, url, tytul, stan, przypieta,
	                            zamknieta, grupa, przestrzen, karta_otwierajaca, kolejnosc,
	                            ostatnio_czynna)
	                           VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	                           ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                               okno = excluded.okno,
	                               url = excluded.url,
	                               tytul = excluded.tytul,
	                               stan = excluded.stan,
	                               przypieta = excluded.przypieta,
	                               zamknieta = excluded.zamknieta,
	                               grupa = excluded.grupa,
	                               przestrzen = excluded.przestrzen,
	                               karta_otwierajaca = excluded.karta_otwierajaca,
	                               kolejnosc = excluded.kolejnosc,
	                               ostatnio_czynna = excluded.ostatnio_czynna`

	pobierzKartePrzegladania = `SELECT ` + kolumnyKartyPrzegladania + `
	                            FROM karta_przegladania WHERE identyfikator_zewnetrzny = ?`

	// Karty zamknięte nie wchodzą do wykazu nigdy: rząd kart pokazuje to, co
	// otwarte. Zawężenie do przestrzeni działa tym samym idiomem co reszta
	// wykazów modułu — pusty tekst wyłącza warunek, więc plan zapytania jest
	// jeden.
	listaKartPrzegladania = `SELECT ` + kolumnyKartyPrzegladania + `
	                         FROM karta_przegladania
	                         WHERE okno = ? AND (? = 1 OR zamknieta = 0)
	                           AND (? = '' OR przestrzen = ?)
	                           AND (? = 1 OR stan <> 'suspended')
	                         ORDER BY kolejnosc, id`

	najwyzszaKolejnoscKart = `SELECT COALESCE(MAX(kolejnosc), -1) FROM karta_przegladania
	                          WHERE okno = ? AND zamknieta = 0`

	zamknijKartePrzegladania = `UPDATE karta_przegladania
	                            SET zamknieta = 1, stan = 'inactive'
	                            WHERE identyfikator_zewnetrzny = ? AND zamknieta = 0`

	odznaczCzynneKarty = `UPDATE karta_przegladania SET stan = 'inactive'
	                      WHERE okno = ? AND zamknieta = 0 AND stan = 'active'
	                        AND identyfikator_zewnetrzny <> ?`

	kolumnyGrupyKart = `id, identyfikator_zewnetrzny, okno, nazwa, barwa, zwinieta, utworzono`

	zapiszGrupeKart = `INSERT INTO grupa_kart_przegladania
	                   (identyfikator_zewnetrzny, okno, nazwa, barwa, zwinieta)
	                   VALUES (?, ?, ?, ?, ?)
	                   ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                       okno = excluded.okno,
	                       nazwa = excluded.nazwa,
	                       barwa = excluded.barwa,
	                       zwinieta = excluded.zwinieta`

	pobierzGrupeKart = `SELECT ` + kolumnyGrupyKart + `
	                    FROM grupa_kart_przegladania WHERE identyfikator_zewnetrzny = ?`

	listaGrupKart = `SELECT ` + kolumnyGrupyKart + `
	                 FROM grupa_kart_przegladania WHERE okno = ?
	                 ORDER BY utworzono DESC, id DESC LIMIT ?`

	usunGrupeKart = `DELETE FROM grupa_kart_przegladania WHERE identyfikator_zewnetrzny = ?`

	odepnijKartyGrupy = `UPDATE karta_przegladania SET grupa = NULL WHERE grupa = ?`

	przypiszKarteDoGrupy = `UPDATE karta_przegladania SET grupa = ?
	                        WHERE identyfikator_zewnetrzny = ? AND okno = ?`

	kolumnyPrzestrzeni = `id, identyfikator_zewnetrzny, okno, nazwa, profil, utworzono, zaktualizowano`

	zapiszPrzestrzen = `INSERT INTO przestrzen_przegladania
	                    (identyfikator_zewnetrzny, okno, nazwa, profil)
	                    VALUES (?, ?, ?, ?)
	                    ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                        okno = excluded.okno,
	                        nazwa = excluded.nazwa,
	                        profil = excluded.profil,
	                        zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')`

	pobierzPrzestrzen = `SELECT ` + kolumnyPrzestrzeni + `
	                     FROM przestrzen_przegladania WHERE identyfikator_zewnetrzny = ?`

	listaPrzestrzeni = `SELECT ` + kolumnyPrzestrzeni + `
	                    FROM przestrzen_przegladania
	                    WHERE (? = '' OR okno = ?)
	                    ORDER BY utworzono DESC, id DESC LIMIT ?`

	usunPrzestrzen = `DELETE FROM przestrzen_przegladania WHERE identyfikator_zewnetrzny = ?`

	przypiszKarteDoPrzestrzeni = `UPDATE karta_przegladania SET przestrzen = ?
	                              WHERE identyfikator_zewnetrzny = ? AND okno = ?`

	liczbaKartPrzestrzeni = `SELECT COUNT(*) FROM karta_przegladania
	                         WHERE przestrzen = ? AND zamknieta = 0`
)

// ZapiszKarte zakłada kartę albo nadpisuje zastaną (dopasowaną po kodzie
// zewnętrznym) i oddaje stan po zapisie.
func (r *repozytoriumPrzegladania) ZapiszKarte(ctx context.Context,
	karta KartaPrzegladania) (KartaPrzegladania, error) {

	if karta.Kod == "" || karta.Okno == "" {
		return KartaPrzegladania{}, fmt.Errorf("dane: karta przeglądania bez identyfikatora albo okna")
	}
	if karta.Stan == "" {
		karta.Stan = "inactive"
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszKartePrzegladania)
	if err != nil {
		return KartaPrzegladania{}, err
	}
	_, err = polecenie.ExecContext(ctx, karta.Kod, karta.Okno, tekstDoKolumny(karta.Url),
		tekstDoKolumny(karta.Tytul), karta.Stan, liczbaLogiczna(karta.Przypieta),
		liczbaLogiczna(karta.Zamknieta), tekstDoKolumny(karta.Grupa),
		tekstDoKolumny(karta.Przestrzen), tekstDoKolumny(karta.KartaOtwierajaca),
		karta.Kolejnosc, tekstDoKolumny(karta.OstatnioCzynna))
	if err != nil {
		return KartaPrzegladania{}, fmt.Errorf("dane: nie można zapisać karty przeglądania %q: %w", karta.Kod, err)
	}
	return r.Karta(ctx, karta.Kod)
}

// Karta oddaje kartę o wskazanym kodzie; brak wiersza wraca jako ErrBrakWiersza.
func (r *repozytoriumPrzegladania) Karta(ctx context.Context, kod string) (KartaPrzegladania, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzKartePrzegladania)
	if err != nil {
		return KartaPrzegladania{}, err
	}
	karta, err := odczytajKartePrzegladania(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return KartaPrzegladania{}, ErrBrakWiersza
	}
	if err != nil {
		return KartaPrzegladania{}, fmt.Errorf("dane: nieczytelna karta przeglądania %q: %w", kod, err)
	}
	return karta, nil
}

// Karty oddaje otwarte karty okna w kolejności rzędu kart.
func (r *repozytoriumPrzegladania) Karty(ctx context.Context,
	filtr FiltrKartPrzegladania) ([]KartaPrzegladania, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaKartPrzegladania)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, filtr.Okno, liczbaLogiczna(filtr.ZZamknietymi),
		filtr.Przestrzen, filtr.Przestrzen, liczbaLogiczna(filtr.ZZawieszonymi))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać kart okna %q: %w", filtr.Okno, err)
	}
	defer wiersze.Close()

	lista := []KartaPrzegladania{}
	for wiersze.Next() {
		karta, err := odczytajKartePrzegladania(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz kart okna %q: %w", filtr.Okno, err)
		}
		lista = append(lista, karta)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt kart okna %q: %w", filtr.Okno, err)
	}
	return lista, nil
}

// NastepnaKolejnoscKarty oddaje pozycję, na którą wchodzi nowa karta okna.
// Liczenie po stronie bazy, a nie po długości wykazu: dwie karty otwarte
// równocześnie dostałyby wtedy tę samą pozycję.
func (r *repozytoriumPrzegladania) NastepnaKolejnoscKarty(ctx context.Context, okno string) (int64, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, najwyzszaKolejnoscKart)
	if err != nil {
		return 0, err
	}
	var najwyzsza int64
	if err := polecenie.QueryRowContext(ctx, okno).Scan(&najwyzsza); err != nil {
		return 0, fmt.Errorf("dane: nie można ustalić kolejności kart okna %q: %w", okno, err)
	}
	return najwyzsza + 1, nil
}

// ZamknijKarte oznacza kartę jako zamkniętą i mówi, czy naprawdę była otwarta.
// Fałsz przy istniejącym wierszu znaczy „już zamknięta" — to nie to samo co
// brak karty i wołający ma prawo rozróżnić te dwie odpowiedzi.
func (r *repozytoriumPrzegladania) ZamknijKarte(ctx context.Context, kod string) (bool, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, zamknijKartePrzegladania)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, kod)
	if err != nil {
		return false, fmt.Errorf("dane: nie można zamknąć karty %q: %w", kod, err)
	}
	zmienione, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nie można ustalić skutku zamknięcia karty %q: %w", kod, err)
	}
	return zmienione > 0, nil
}

// OdznaczPozostaleKarty gasi czynność pozostałych kart okna. Karta czynna jest
// w oknie jedna — dwie naraz byłyby dwoma odpowiedziami na pytanie, co Operator
// i Wykonawca widzą.
func (r *repozytoriumPrzegladania) OdznaczPozostaleKarty(ctx context.Context, okno, kod string) error {
	polecenie, err := r.zapytania.przygotuj(ctx, odznaczCzynneKarty)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, okno, kod); err != nil {
		return fmt.Errorf("dane: nie można odznaczyć kart okna %q: %w", okno, err)
	}
	return nil
}

// PrzypiszKarteDoGrupy wiąże kartę z grupą albo zdejmuje wiązanie (kod pusty).
func (r *repozytoriumPrzegladania) PrzypiszKarteDoGrupy(ctx context.Context, okno, karta, grupa string) error {
	polecenie, err := r.zapytania.przygotuj(ctx, przypiszKarteDoGrupy)
	if err != nil {
		return err
	}
	var wartosc any
	if grupa != "" {
		wartosc = grupa
	}
	if _, err := polecenie.ExecContext(ctx, wartosc, karta, okno); err != nil {
		return fmt.Errorf("dane: nie można przypisać karty %q do grupy: %w", karta, err)
	}
	return nil
}

// PrzypiszKarteDoPrzestrzeni wiąże kartę z przestrzenią roboczą.
func (r *repozytoriumPrzegladania) PrzypiszKarteDoPrzestrzeni(ctx context.Context, okno, karta, przestrzen string) error {
	polecenie, err := r.zapytania.przygotuj(ctx, przypiszKarteDoPrzestrzeni)
	if err != nil {
		return err
	}
	var wartosc any
	if przestrzen != "" {
		wartosc = przestrzen
	}
	if _, err := polecenie.ExecContext(ctx, wartosc, karta, okno); err != nil {
		return fmt.Errorf("dane: nie można przypisać karty %q do przestrzeni: %w", karta, err)
	}
	return nil
}

// ZapiszGrupeKart zakłada grupę albo nadpisuje zastaną.
func (r *repozytoriumPrzegladania) ZapiszGrupeKart(ctx context.Context, grupa GrupaKart) (GrupaKart, error) {
	if grupa.Kod == "" || grupa.Okno == "" || grupa.Nazwa == "" {
		return GrupaKart{}, fmt.Errorf("dane: grupa kart bez identyfikatora, okna albo nazwy")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszGrupeKart)
	if err != nil {
		return GrupaKart{}, err
	}
	_, err = polecenie.ExecContext(ctx, grupa.Kod, grupa.Okno, grupa.Nazwa,
		tekstDoKolumny(grupa.Barwa), liczbaLogiczna(grupa.Zwinieta))
	if err != nil {
		return GrupaKart{}, fmt.Errorf("dane: nie można zapisać grupy kart %q: %w", grupa.Kod, err)
	}
	return r.GrupaKart(ctx, grupa.Kod)
}

// GrupaKart oddaje grupę o wskazanym kodzie.
func (r *repozytoriumPrzegladania) GrupaKart(ctx context.Context, kod string) (GrupaKart, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzGrupeKart)
	if err != nil {
		return GrupaKart{}, err
	}
	grupa, err := odczytajGrupeKart(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return GrupaKart{}, ErrBrakWiersza
	}
	if err != nil {
		return GrupaKart{}, fmt.Errorf("dane: nieczytelna grupa kart %q: %w", kod, err)
	}
	return grupa, nil
}

// GrupyKart oddaje grupy okna od najnowszej.
func (r *repozytoriumPrzegladania) GrupyKart(ctx context.Context, okno string, limit int) ([]GrupaKart, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaGrupKart)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno, granicaWykazu(limit))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać grup kart okna %q: %w", okno, err)
	}
	defer wiersze.Close()

	lista := []GrupaKart{}
	for wiersze.Next() {
		grupa, err := odczytajGrupeKart(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz grup kart okna %q: %w", okno, err)
		}
		lista = append(lista, grupa)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt grup kart okna %q: %w", okno, err)
	}
	return lista, nil
}

// UsunGrupeKart zdejmuje grupę i odpina od niej karty. Karty zostają otwarte —
// grupa jest zestawem, nie właścicielem kart.
func (r *repozytoriumPrzegladania) UsunGrupeKart(ctx context.Context, kod string) (bool, error) {
	odpiecie, err := r.zapytania.przygotuj(ctx, odepnijKartyGrupy)
	if err != nil {
		return false, err
	}
	if _, err := odpiecie.ExecContext(ctx, kod); err != nil {
		return false, fmt.Errorf("dane: nie można odpiąć kart grupy %q: %w", kod, err)
	}
	polecenie, err := r.zapytania.przygotuj(ctx, usunGrupeKart)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, kod)
	if err != nil {
		return false, fmt.Errorf("dane: nie można usunąć grupy kart %q: %w", kod, err)
	}
	zmienione, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nie można ustalić skutku usunięcia grupy %q: %w", kod, err)
	}
	return zmienione > 0, nil
}

// ZapiszPrzestrzen zakłada przestrzeń roboczą albo nadpisuje zastaną.
func (r *repozytoriumPrzegladania) ZapiszPrzestrzen(ctx context.Context,
	przestrzen PrzestrzenPrzegladania) (PrzestrzenPrzegladania, error) {

	if przestrzen.Kod == "" || przestrzen.Nazwa == "" {
		return PrzestrzenPrzegladania{}, fmt.Errorf("dane: przestrzeń przeglądania bez identyfikatora albo nazwy")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszPrzestrzen)
	if err != nil {
		return PrzestrzenPrzegladania{}, err
	}
	_, err = polecenie.ExecContext(ctx, przestrzen.Kod, tekstDoKolumny(przestrzen.Okno),
		przestrzen.Nazwa, tekstDoKolumny(przestrzen.Profil))
	if err != nil {
		return PrzestrzenPrzegladania{}, fmt.Errorf("dane: nie można zapisać przestrzeni %q: %w", przestrzen.Kod, err)
	}
	return r.Przestrzen(ctx, przestrzen.Kod)
}

// Przestrzen oddaje przestrzeń roboczą o wskazanym kodzie.
func (r *repozytoriumPrzegladania) Przestrzen(ctx context.Context, kod string) (PrzestrzenPrzegladania, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzPrzestrzen)
	if err != nil {
		return PrzestrzenPrzegladania{}, err
	}
	przestrzen, err := odczytajPrzestrzen(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return PrzestrzenPrzegladania{}, ErrBrakWiersza
	}
	if err != nil {
		return PrzestrzenPrzegladania{}, fmt.Errorf("dane: nieczytelna przestrzeń przeglądania %q: %w", kod, err)
	}
	return przestrzen, nil
}

// Przestrzenie oddaje przestrzenie robocze, opcjonalnie zawężone do okna.
func (r *repozytoriumPrzegladania) Przestrzenie(ctx context.Context, okno string, limit int) ([]PrzestrzenPrzegladania, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaPrzestrzeni)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno, okno, granicaWykazu(limit))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać przestrzeni przeglądania: %w", err)
	}
	defer wiersze.Close()

	lista := []PrzestrzenPrzegladania{}
	for wiersze.Next() {
		przestrzen, err := odczytajPrzestrzen(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz przestrzeni przeglądania: %w", err)
		}
		lista = append(lista, przestrzen)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt przestrzeni przeglądania: %w", err)
	}
	return lista, nil
}

// UsunPrzestrzen zdejmuje przestrzeń roboczą. Karty zostają nietknięte —
// kontrakt mówi o tym wprost przy `browser.workspace.remove`.
func (r *repozytoriumPrzegladania) UsunPrzestrzen(ctx context.Context, kod string) (bool, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, usunPrzestrzen)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, kod)
	if err != nil {
		return false, fmt.Errorf("dane: nie można usunąć przestrzeni %q: %w", kod, err)
	}
	zmienione, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nie można ustalić skutku usunięcia przestrzeni %q: %w", kod, err)
	}
	return zmienione > 0, nil
}

// LiczbaKartPrzestrzeni liczy karty zapisane w przestrzeni roboczej — kontrakt
// oddaje ją polem `tabCount`, a liczba wzięta z bazy nie rozjedzie się
// z rzeczywistym składem tak, jak rozjechałaby się liczba zapamiętana przy
// zapisie.
func (r *repozytoriumPrzegladania) LiczbaKartPrzestrzeni(ctx context.Context, kod string) (int64, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, liczbaKartPrzestrzeni)
	if err != nil {
		return 0, err
	}
	var liczba int64
	if err := polecenie.QueryRowContext(ctx, kod).Scan(&liczba); err != nil {
		return 0, fmt.Errorf("dane: nie można policzyć kart przestrzeni %q: %w", kod, err)
	}
	return liczba, nil
}

// odczytajKartePrzegladania składa kartę z jednego wiersza wyniku.
func odczytajKartePrzegladania(wiersz skaner) (KartaPrzegladania, error) {
	var karta KartaPrzegladania
	var url, tytul, grupa, przestrzen, otwierajaca, czynna sql.NullString
	var przypieta, zamknieta int64
	err := wiersz.Scan(&karta.ID, &karta.Kod, &karta.Okno, &url, &tytul, &karta.Stan,
		&przypieta, &zamknieta, &grupa, &przestrzen, &otwierajaca, &karta.Kolejnosc,
		&czynna, &karta.Utworzono)
	if err != nil {
		return KartaPrzegladania{}, err
	}
	karta.Url = tekstZKolumny(url)
	karta.Tytul = tekstZKolumny(tytul)
	karta.Przypieta = przypieta == 1
	karta.Zamknieta = zamknieta == 1
	karta.Grupa = tekstZKolumny(grupa)
	karta.Przestrzen = tekstZKolumny(przestrzen)
	karta.KartaOtwierajaca = tekstZKolumny(otwierajaca)
	karta.OstatnioCzynna = tekstZKolumny(czynna)
	return karta, nil
}

// odczytajGrupeKart składa grupę z jednego wiersza wyniku.
func odczytajGrupeKart(wiersz skaner) (GrupaKart, error) {
	var grupa GrupaKart
	var barwa sql.NullString
	var zwinieta int64
	if err := wiersz.Scan(&grupa.ID, &grupa.Kod, &grupa.Okno, &grupa.Nazwa, &barwa,
		&zwinieta, &grupa.Utworzono); err != nil {
		return GrupaKart{}, err
	}
	grupa.Barwa = tekstZKolumny(barwa)
	grupa.Zwinieta = zwinieta == 1
	return grupa, nil
}

// odczytajPrzestrzen składa przestrzeń roboczą z jednego wiersza wyniku.
func odczytajPrzestrzen(wiersz skaner) (PrzestrzenPrzegladania, error) {
	var przestrzen PrzestrzenPrzegladania
	var okno, profil sql.NullString
	if err := wiersz.Scan(&przestrzen.ID, &przestrzen.Kod, &okno, &przestrzen.Nazwa,
		&profil, &przestrzen.Utworzono, &przestrzen.Zaktualizowano); err != nil {
		return PrzestrzenPrzegladania{}, err
	}
	przestrzen.Okno = tekstZKolumny(okno)
	przestrzen.Profil = tekstZKolumny(profil)
	return przestrzen, nil
}
