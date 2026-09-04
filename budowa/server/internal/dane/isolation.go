// Plik obsługuje obszar profili izolacji oraz odczyt słownika poziomów zasięgu z tabeli
// poziom_zasiegu na potrzeby rodziny poleceń isolation. Uzasadnienie granicy wobec
// repozytorium konfiguracji niesie rozdział isolation.go dokumentacji architektury.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"danacoconsole/shared"
)

// PoziomZasiegu to wiersz tabeli `poziom_zasiegu`: poziom kontraktu, jego nazwa
// do wyświetlenia i pierwszeństwo (1 — najszerszy, 8 — najwęższy).
type PoziomZasiegu struct {
	Poziom        shared.ConfigScope
	Nazwa         string
	Pierwszenstwo int
}

// PrzelacznikProfiluIzolacji to jeden punkt izolacji zapisany w profilu. Klucz
// i wartość są dokładnie te, które czyta rozstrzygacz ustawień — profil nie ma
// własnego słownika nazw.
type PrzelacznikProfiluIzolacji struct {
	Klucz   string
	Wartosc string
}

// ProfilIzolacji odwzorowuje wiersz tabeli profil_izolacji wraz z kompletem jego
// przełączników — nazwany szablon ustawień izolacji, który operator przypisuje do poziomów.
type ProfilIzolacji struct {
	Kod            string
	Nazwa          string
	Opis           string
	Przelaczniki   []PrzelacznikProfiluIzolacji
	Utworzono      string
	Zaktualizowano string
}

// RepozytoriumIzolacji jest kontraktem obszaru profili izolacji: definiuje odczyt słownika
// poziomów, zapis i odczyt profili, przypisanie profilu do poziomu oraz wybór warstwy.
type RepozytoriumIzolacji interface {
	// PoziomyZasiegu zwraca osiem poziomów w kolejności pierwszeństwa, od najszerszego do najwęższego.
	PoziomyZasiegu(ctx context.Context) ([]PoziomZasiegu, error)
	// ZapiszProfil zakłada albo zmienia profil wraz z kompletem przełączników, zastępowanych w całości.
	ZapiszProfil(ctx context.Context, profil ProfilIzolacji) (ProfilIzolacji, error)
	// Profil zwraca profil po kodzie. Brak wiersza daje ErrBrakWiersza.
	Profil(ctx context.Context, kod string) (ProfilIzolacji, error)
	// Profile zwraca wszystkie profile w kolejności nazwy.
	Profile(ctx context.Context) ([]ProfilIzolacji, error)
	// UsunProfil kasuje profil. Drugi wynik mówi, czy wiersz istniał.
	UsunProfil(ctx context.Context, kod string) (bool, error)
	// PrzypiszProfil zapisuje, z którego profilu pochodzi polityka bytu poziomu.
	PrzypiszProfil(ctx context.Context, kod string, poziom shared.ConfigScope, kluczZasiegu string) error
	// ProfilPoziomu zwraca kod profilu przypisanego pod adresem. Drugi wynik mówi, czy przypisanie jest.
	ProfilPoziomu(ctx context.Context, poziom shared.ConfigScope, kluczZasiegu string) (string, bool, error)
	// ZapiszWarstwe utrwala warstwę wybraną przez Operatora dla bytu poziomu.
	ZapiszWarstwe(ctx context.Context, poziom shared.ConfigScope, kluczZasiegu, warstwa string) error
	// Warstwa zwraca warstwę wybraną dla bytu poziomu. Drugi wynik mówi, czy wybór w ogóle zapadł.
	Warstwa(ctx context.Context, poziom shared.ConfigScope, kluczZasiegu string) (string, bool, error)
}

const listaPoziomowZasiegu = `SELECT kod, nazwa, pierwszenstwo FROM poziom_zasiegu
                              ORDER BY pierwszenstwo`

// Profil i warstwa izolacji niosą wskazanie konta; zapytania biorą je stąd.
var (

	// Człon WHERE przy DO UPDATE zatrzymuje nadpisanie profilu konta cudzego:
	// więz jednoznaczności kodu obejmuje całą tabelę, nie konto.
	zapiszProfilIzolacji = `INSERT INTO profil_izolacji (kod, nazwa, opis, konto_id)
	                        VALUES (?, ?, ?, ` + WskazanieKonta + `)
	                        ON CONFLICT(kod) DO UPDATE SET
	                            nazwa = excluded.nazwa,
	                            opis = excluded.opis,
	                            zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                        WHERE ` + WarunekKonta

	usunPrzelacznikiProfilu = `DELETE FROM przelacznik_profilu_izolacji
	                           WHERE profil_id = (SELECT id FROM profil_izolacji
	                                              WHERE kod = ? AND ` + WarunekKonta + `)`

	zapiszPrzelacznikProfilu = `INSERT INTO przelacznik_profilu_izolacji (profil_id, klucz, wartosc)
	                            VALUES ((SELECT id FROM profil_izolacji
	                                     WHERE kod = ? AND ` + WarunekKonta + `), ?, ?)
	                            ON CONFLICT(profil_id, klucz) DO UPDATE SET wartosc = excluded.wartosc`

	pobierzProfilIzolacji = `SELECT kod, nazwa, opis, utworzono, zaktualizowano
	                         FROM profil_izolacji WHERE kod = ? AND ` + WarunekKonta

	listaProfiliIzolacji = `SELECT kod, nazwa, opis, utworzono, zaktualizowano
	                        FROM profil_izolacji WHERE ` + WarunekKonta + `
	                        ORDER BY nazwa, kod`

	pobierzPrzelacznikiProfilu = `SELECT klucz, wartosc FROM przelacznik_profilu_izolacji
	                              WHERE profil_id = (SELECT id FROM profil_izolacji
	                                                 WHERE kod = ? AND ` + WarunekKonta + `)
	                              ORDER BY klucz`

	usunProfilIzolacji = `DELETE FROM profil_izolacji WHERE kod = ? AND ` + WarunekKonta

	zapiszPrzypisanieProfilu = `INSERT INTO przypisanie_profilu_izolacji
	                            (profil_id, poziom_zasiegu_id, klucz_zasiegu)
	                            VALUES ((SELECT id FROM profil_izolacji
	                                     WHERE kod = ? AND ` + WarunekKonta + `),
	                                    (SELECT id FROM poziom_zasiegu WHERE kod = ?), ?)
	                            ON CONFLICT(poziom_zasiegu_id, klucz_zasiegu) DO UPDATE SET
	                                profil_id = excluded.profil_id,
	                                przypisano = strftime('%Y-%m-%dT%H:%M:%fZ','now')`

	pobierzPrzypisanieProfilu = `SELECT p.kod FROM przypisanie_profilu_izolacji a
	                             JOIN profil_izolacji p ON p.id = a.profil_id
	                             JOIN poziom_zasiegu z ON z.id = a.poziom_zasiegu_id
	                             WHERE z.kod = ? AND a.klucz_zasiegu = ?
	                               AND ` + strings.ReplaceAll(WarunekKonta, "konto_id", "p.konto_id")

	zapiszWarstweIzolacji = `INSERT INTO warstwa_izolacji
	                         (poziom_zasiegu_id, klucz_zasiegu, warstwa, konto_id)
	                         VALUES ((SELECT id FROM poziom_zasiegu WHERE kod = ?), ?, ?,
	                                 ` + WskazanieKonta + `)
	                         ON CONFLICT(poziom_zasiegu_id, klucz_zasiegu,
	                                     COALESCE(konto_id, 0)) DO UPDATE SET
	                             warstwa = excluded.warstwa,
	                             zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                         WHERE ` + WarunekKonta

	pobierzWarstweIzolacji = `SELECT w.warstwa FROM warstwa_izolacji w
	                          JOIN poziom_zasiegu z ON z.id = w.poziom_zasiegu_id
	                          WHERE z.kod = ? AND w.klucz_zasiegu = ?
	                            AND ` + strings.ReplaceAll(WarunekKonta, "konto_id", "w.konto_id")
)

type repozytoriumIzolacji struct {
	zapytania *zapytania
	db        *sql.DB
}

// ProfileIzolacji oddaje repozytorium profili izolacji nad tą samą bazą, co pozostałe
// obszary zestawu. Jest metodą, nie polem struktury, ponieważ repozytorium jest bezstanowe
// i kolejne wywołania są równoważne.
func (z *Zestaw) ProfileIzolacji() RepozytoriumIzolacji {
	if z == nil || z.zapytania == nil {
		return nil
	}
	return &repozytoriumIzolacji{zapytania: z.zapytania, db: z.zapytania.db}
}

// PoziomyZasiegu czyta słownik poziomów zasięgu w kolejności pierwszeństwa rosnąco, od
// poziomu najszerszego do najwęższego, wraz z nazwą każdego poziomu do wyświetlenia.
func (r *repozytoriumIzolacji) PoziomyZasiegu(ctx context.Context) ([]PoziomZasiegu, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaPoziomowZasiegu)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać poziomów zasięgu: %w", err)
	}
	defer wiersze.Close()

	poziomy := make([]PoziomZasiegu, 0, 8)
	for wiersze.Next() {
		var kod string
		var poziom PoziomZasiegu
		if err := wiersze.Scan(&kod, &poziom.Nazwa, &poziom.Pierwszenstwo); err != nil {
			return nil, err
		}
		if poziom.Poziom, err = poziomZasieguZBazy(kod); err != nil {
			return nil, err
		}
		poziomy = append(poziomy, poziom)
	}
	return poziomy, wiersze.Err()
}

// ZapiszProfil utrwala profil i zastępuje komplet jego przełączników. Zapis idzie
// jedną transakcją: profil bez przełączników albo przełączniki bez profilu byłyby
// stanem, którego nikt nie zamawiał.
func (r *repozytoriumIzolacji) ZapiszProfil(ctx context.Context, profil ProfilIzolacji) (ProfilIzolacji, error) {
	if profil.Kod == "" || profil.Nazwa == "" {
		return ProfilIzolacji{}, fmt.Errorf("dane: profil izolacji bez kodu albo nazwy")
	}
	transakcja, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return ProfilIzolacji{}, fmt.Errorf("dane: nie można otworzyć transakcji profilu izolacji: %w", err)
	}
	defer func() { _ = transakcja.Rollback() }()

	konto := KontoOperatora(ctx)
	if _, err := transakcja.ExecContext(ctx, zapiszProfilIzolacji, profil.Kod, profil.Nazwa, profil.Opis,
		konto, konto); err != nil {
		return ProfilIzolacji{}, fmt.Errorf("dane: nie można zapisać profilu izolacji %q: %w", profil.Kod, err)
	}
	if _, err := transakcja.ExecContext(ctx, usunPrzelacznikiProfilu, profil.Kod, konto); err != nil {
		return ProfilIzolacji{}, fmt.Errorf("dane: nie można wyczyścić przełączników profilu %q: %w", profil.Kod, err)
	}
	for _, przelacznik := range profil.Przelaczniki {
		if przelacznik.Klucz == "" {
			continue
		}
		if _, err := transakcja.ExecContext(ctx, zapiszPrzelacznikProfilu,
			profil.Kod, konto, przelacznik.Klucz, przelacznik.Wartosc); err != nil {
			return ProfilIzolacji{}, fmt.Errorf("dane: nie można zapisać przełącznika %q profilu %q: %w",
				przelacznik.Klucz, profil.Kod, err)
		}
	}
	if err := transakcja.Commit(); err != nil {
		return ProfilIzolacji{}, fmt.Errorf("dane: nie można domknąć zapisu profilu izolacji %q: %w", profil.Kod, err)
	}
	return r.Profil(ctx, profil.Kod)
}

// Profil zwraca profil wraz z przełącznikami. Brak wiersza daje ErrBrakWiersza —
// warstwa wyższa rozpoznaje go przez errors.Is i odmawia z kodem `not_found`.
func (r *repozytoriumIzolacji) Profil(ctx context.Context, kod string) (ProfilIzolacji, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzProfilIzolacji)
	if err != nil {
		return ProfilIzolacji{}, err
	}
	profil, err := odczytajProfilIzolacji(polecenie.QueryRowContext(ctx, kod, KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return ProfilIzolacji{}, ErrBrakWiersza
	}
	if err != nil {
		return ProfilIzolacji{}, err
	}
	if profil.Przelaczniki, err = r.przelaczniki(ctx, kod); err != nil {
		return ProfilIzolacji{}, err
	}
	return profil, nil
}

// Profile zwraca wszystkie zapisane profile izolacji wraz z ich przełącznikami,
// uporządkowane według nazwy profilu.
func (r *repozytoriumIzolacji) Profile(ctx context.Context) ([]ProfilIzolacji, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaProfiliIzolacji)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, KontoOperatora(ctx))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać profili izolacji: %w", err)
	}
	defer wiersze.Close()

	profile := make([]ProfilIzolacji, 0)
	for wiersze.Next() {
		profil, err := odczytajProfilIzolacji(wiersze)
		if err != nil {
			return nil, err
		}
		profile = append(profile, profil)
	}
	if err := wiersze.Err(); err != nil {
		return nil, err
	}
	for indeks := range profile {
		if profile[indeks].Przelaczniki, err = r.przelaczniki(ctx, profile[indeks].Kod); err != nil {
			return nil, err
		}
	}
	return profile, nil
}

// UsunProfil kasuje profil. Drugi wynik mówi, czy wiersz istniał — odpowiedź
// „usunięto" o bycie, którego nie było, byłaby nieprawdą.
func (r *repozytoriumIzolacji) UsunProfil(ctx context.Context, kod string) (bool, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, usunProfilIzolacji)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, kod, KontoOperatora(ctx))
	if err != nil {
		return false, fmt.Errorf("dane: nie można usunąć profilu izolacji %q: %w", kod, err)
	}
	zmienione, err := wynik.RowsAffected()
	if err != nil {
		return false, err
	}
	return zmienione > 0, nil
}

// PrzypiszProfil zapisuje, który profil izolacji obowiązuje dla wskazanego bytu poziomu
// zasięgu, nadpisując zastane przypisanie pod tym samym adresem.
func (r *repozytoriumIzolacji) PrzypiszProfil(ctx context.Context, kod string,
	poziom shared.ConfigScope, kluczZasiegu string) error {

	kodPoziomu, err := poziomZasieguNaBaze(poziom)
	if err != nil {
		return err
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszPrzypisanieProfilu)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, kod, KontoOperatora(ctx), kodPoziomu, kluczZasiegu); err != nil {
		return fmt.Errorf("dane: nie można przypisać profilu izolacji %q do poziomu %q: %w",
			kod, kodPoziomu, err)
	}
	return nil
}

// ProfilPoziomu zwraca profil przypisany pod adresem. Brak przypisania nie jest
// błędem — poziom bez profilu ma politykę złożoną z pojedynczych ustawień.
func (r *repozytoriumIzolacji) ProfilPoziomu(ctx context.Context, poziom shared.ConfigScope,
	kluczZasiegu string) (string, bool, error) {

	return r.jednaWartosc(ctx, pobierzPrzypisanieProfilu, poziom, kluczZasiegu)
}

// ZapiszWarstwe utrwala warstwę wybraną przez operatora dla wskazanego bytu poziomu,
// nadpisując zastany wybór pod tym samym adresem.
func (r *repozytoriumIzolacji) ZapiszWarstwe(ctx context.Context, poziom shared.ConfigScope,
	kluczZasiegu, warstwa string) error {

	kodPoziomu, err := poziomZasieguNaBaze(poziom)
	if err != nil {
		return err
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszWarstweIzolacji)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, kodPoziomu, kluczZasiegu, warstwa,
		KontoOperatora(ctx), KontoOperatora(ctx)); err != nil {
		return fmt.Errorf("dane: nie można zapisać warstwy izolacji %q na poziomie %q: %w",
			warstwa, kodPoziomu, err)
	}
	return nil
}

// Warstwa zwraca warstwę wybraną przez operatora dla wskazanego bytu poziomu. Drugi wynik
// mówi, czy taki wybór w ogóle zapadł.
func (r *repozytoriumIzolacji) Warstwa(ctx context.Context, poziom shared.ConfigScope,
	kluczZasiegu string) (string, bool, error) {

	return r.jednaWartosc(ctx, pobierzWarstweIzolacji, poziom, kluczZasiegu)
}

// jednaWartosc czyta jedną kolumnę tekstową spod adresu poziom + byt. Brak
// wiersza znaczy brak zapisu, nie błąd odczytu.
func (r *repozytoriumIzolacji) jednaWartosc(ctx context.Context, zapytanie string,
	poziom shared.ConfigScope, kluczZasiegu string) (string, bool, error) {

	kodPoziomu, err := poziomZasieguNaBaze(poziom)
	if err != nil {
		return "", false, err
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return "", false, err
	}
	var wartosc string
	err = polecenie.QueryRowContext(ctx, kodPoziomu, kluczZasiegu, KontoOperatora(ctx)).Scan(&wartosc)
	if errors.Is(err, sql.ErrNoRows) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return wartosc, true, nil
}

// przelaczniki czyta wszystkie punkty izolacji zapisane w profilu wskazanym kodem,
// uporządkowane według klucza przełącznika.
func (r *repozytoriumIzolacji) przelaczniki(ctx context.Context, kod string) ([]PrzelacznikProfiluIzolacji, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzPrzelacznikiProfilu)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, kod, KontoOperatora(ctx))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać przełączników profilu %q: %w", kod, err)
	}
	defer wiersze.Close()

	przelaczniki := make([]PrzelacznikProfiluIzolacji, 0)
	for wiersze.Next() {
		var przelacznik PrzelacznikProfiluIzolacji
		if err := wiersze.Scan(&przelacznik.Klucz, &przelacznik.Wartosc); err != nil {
			return nil, err
		}
		przelaczniki = append(przelaczniki, przelacznik)
	}
	return przelaczniki, wiersze.Err()
}

// odczytajProfilIzolacji składa strukturę ProfilIzolacji z jednego wiersza wyniku zapytania,
// bez jego przełączników, które czyta osobna funkcja.
func odczytajProfilIzolacji(wiersz skaner) (ProfilIzolacji, error) {
	var profil ProfilIzolacji
	err := wiersz.Scan(&profil.Kod, &profil.Nazwa, &profil.Opis, &profil.Utworzono, &profil.Zaktualizowano)
	if err != nil {
		return ProfilIzolacji{}, err
	}
	return profil, nil
}
