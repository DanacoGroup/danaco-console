// Odpowiedzialność pliku: obszar modułu Roundtable, byty debaty, kontrakt całego obszaru i uczestnik, wraz z kanałem jego głosu.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// Ten sam kanał modelu może wystąpić dwukrotnie pod odrębnymi tożsamościami: jednoznaczny jest wyłącznie Kod.
type UczestnikDebaty struct {
	Kod             string
	Okno            string
	KanalModelu     string
	NazwaTozsamosci *string
	PromptSystemowy *string
	Wyciszony       bool
	Kolejnosc       int
	Kluczowy        bool
	Waga            float64
	Rola            string
	Agent           *string
	Awatar          *string
	OpisRoli        *string
	LiczbaProbek    int

	Utworzono      string
	Zaktualizowano string
}

// Format i Stan niosą wartości kontraktu wprost (RoundtableFormat, RoundtableTurnStatus).
type TuraDebaty struct {
	Kod            string
	Okno           string
	Numer          int
	Zagadnienie    *string
	Pytanie        string
	Format         string
	Stan           string
	GranicaTur     int
	TuraNadrzedna  string
	GranicaCzasuMs int
	GranicaZnakow  int
	Anonimowa      bool

	Rozpoczeto string
	Zamknieto  *string
}

// Uczestnik jest kodem, nie kluczem obcym: moderator wypowiada się w turze, a uczestnikiem nie jest.
type WypowiedzDebaty struct {
	Kod         string
	TuraKod     string
	Uczestnik   string
	Tresc       string
	OdpowiedzNa string
	AktMowy     string
	Pewnosc     float64
	Redakcja    int

	Utworzono string
}

// Pusta Tura znaczy stanowisko całej debaty, nie jednej tury.
type StanowiskoDebaty struct {
	Kod            string
	Okno           string
	Tura           string
	Tresc          *string
	Wersja         int
	Zaktualizowano string
	Redagowane     bool
	Zaakceptowane  bool
	Kontekst       *string
	Warianty       *string
	Konsekwencje   *string
	Tury           string
}

type RepozytoriumRoundtable interface {
	RepozytoriumDebatySkladu
	RepozytoriumDebatyGrafu
	RepozytoriumDebatyAnalizy
	RepozytoriumDebatyGlosowan
	RepozytoriumDebatyOceny
	RepozytoriumDebatyDecyzji
	RepozytoriumDebatyKonsensusu
	RepozytoriumDebatyWydania

	ZapiszUczestnika(ctx context.Context, uczestnik UczestnikDebaty) (UczestnikDebaty, error)
	Uczestnik(ctx context.Context, kod string) (UczestnikDebaty, error)
	Uczestnicy(ctx context.Context, okno string) ([]UczestnikDebaty, error)
	UstawWyciszenie(ctx context.Context, kod string, wyciszony bool) error
	UstawKolejnosc(ctx context.Context, okno string, kody []string) error

	ZalozTure(ctx context.Context, tura TuraDebaty) (TuraDebaty, error)
	ZmienTure(ctx context.Context, tura TuraDebaty) error
	Tura(ctx context.Context, kod string) (TuraDebaty, error)
	Tury(ctx context.Context, okno string, limit int) ([]TuraDebaty, error)

	ZapiszWypowiedz(ctx context.Context, wypowiedz WypowiedzDebaty) (WypowiedzDebaty, error)
	UzupelnijWypowiedz(ctx context.Context, kod, tresc string) error
	Wypowiedzi(ctx context.Context, turaKod string) ([]WypowiedzDebaty, error)
	Wypowiedz(ctx context.Context, kod string) (WypowiedzDebaty, error)
	WypowiedziOkna(ctx context.Context, okno string) ([]WypowiedzDebaty, error)
	ZastapWypowiedz(ctx context.Context, kod, tresc string) error
	OznaczWypowiedz(ctx context.Context, kod, aktMowy string, pewnosc float64) error

	ZapiszStanowisko(ctx context.Context, stanowisko StanowiskoDebaty) (StanowiskoDebaty, error)
	Stanowisko(ctx context.Context, okno, tura string) (StanowiskoDebaty, error)
}

const (
	kolumnyUczestnika = `identyfikator_zewnetrzny, okno, kanal_modelu, nazwa_tozsamosci,
	                     prompt_systemowy, wyciszony, kolejnosc, kluczowy, waga, rola,
	                     agent, awatar, opis_roli, liczba_probek, utworzono, zaktualizowano`

	// Identyfikator zewnętrzny jest jednoznaczny w całej tabeli: gałąź konfliktu bez warunku konta sięgałaby uczestnika konta cudzego.
	zapiszUczestnikaDebaty = `INSERT INTO debata_uczestnik
	                          (identyfikator_zewnetrzny, okno, kanal_modelu, nazwa_tozsamosci,
	                           prompt_systemowy, wyciszony, kolejnosc, konto_id)
	                          VALUES (?, ?, ?, ?, ?, ?, ?, ` + WskazanieKonta + `)
	                          ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                              nazwa_tozsamosci = excluded.nazwa_tozsamosci,
	                              prompt_systemowy = excluded.prompt_systemowy,
	                              zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                          WHERE ` + WarunekKonta

	pobierzUczestnika = `SELECT ` + kolumnyUczestnika + `
	                     FROM debata_uczestnik WHERE identyfikator_zewnetrzny = ?
	                       AND ` + WarunekKonta

	pobierzUczestnikow = `SELECT ` + kolumnyUczestnika + `
	                      FROM debata_uczestnik WHERE okno = ? AND ` + WarunekKonta + `
	                      ORDER BY kolejnosc ASC, id ASC`

	ustawWyciszenieUczestnika = `UPDATE debata_uczestnik
	                             SET wyciszony = ?,
	                                 zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                             WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	ustawKolejnoscUczestnika = `UPDATE debata_uczestnik
	                            SET kolejnosc = ?,
	                                zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                            WHERE okno = ? AND identyfikator_zewnetrzny = ?
	                              AND ` + WarunekKonta
)

type repozytoriumRoundtable struct {
	zapytania *zapytania
	db        *sql.DB
}

func noweRepozytoriumRoundtable(z *zapytania, db *sql.DB) *repozytoriumRoundtable {
	return &repozytoriumRoundtable{zapytania: z, db: db}
}

// Kolejność i wyciszenie zostają nietknięte: ustawia je moderator, a powtórny zapis tożsamości ich nie cofa.
func (r *repozytoriumRoundtable) ZapiszUczestnika(ctx context.Context,
	uczestnik UczestnikDebaty) (UczestnikDebaty, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, zapiszUczestnikaDebaty)
	if err != nil {
		return UczestnikDebaty{}, err
	}
	wynik, err := polecenie.ExecContext(ctx, uczestnik.Kod, uczestnik.Okno, uczestnik.KanalModelu,
		uczestnik.NazwaTozsamosci, uczestnik.PromptSystemowy, uczestnik.Wyciszony,
		uczestnik.Kolejnosc, KontoOperatora(ctx), KontoOperatora(ctx))
	if err != nil {
		return UczestnikDebaty{}, fmt.Errorf("dane: nie można zapisać uczestnika debaty %q: %w",
			uczestnik.Kod, err)
	}
	if err := sprawdzTrafienieZapisu(wynik, "uczestnik debaty", uczestnik.Kod); err != nil {
		return UczestnikDebaty{}, err
	}
	return r.Uczestnik(ctx, uczestnik.Kod)
}

func (r *repozytoriumRoundtable) Uczestnik(ctx context.Context, kod string) (UczestnikDebaty, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzUczestnika)
	if err != nil {
		return UczestnikDebaty{}, err
	}
	uczestnik, err := odczytajUczestnika(polecenie.QueryRowContext(ctx, kod, KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return UczestnikDebaty{}, ErrBrakWiersza
	}
	if err != nil {
		return UczestnikDebaty{}, fmt.Errorf("dane: nieczytelny wiersz uczestnika debaty %q: %w", kod, err)
	}
	return uczestnik, nil
}

func (r *repozytoriumRoundtable) Uczestnicy(ctx context.Context, okno string) ([]UczestnikDebaty, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzUczestnikow)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno, KontoOperatora(ctx))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać uczestników debaty okna %q: %w", okno, err)
	}
	defer wiersze.Close()

	uczestnicy := make([]UczestnikDebaty, 0, 8)
	for wiersze.Next() {
		uczestnik, err := odczytajUczestnika(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz uczestnika debaty: %w", err)
		}
		uczestnicy = append(uczestnicy, uczestnik)
	}
	return uczestnicy, wiersze.Err()
}

func (r *repozytoriumRoundtable) UstawWyciszenie(ctx context.Context, kod string, wyciszony bool) error {
	polecenie, err := r.zapytania.przygotuj(ctx, ustawWyciszenieUczestnika)
	if err != nil {
		return err
	}
	wynik, err := polecenie.ExecContext(ctx, wyciszony, kod, KontoOperatora(ctx))
	if err != nil {
		return fmt.Errorf("dane: nie można przestawić wyciszenia uczestnika %q: %w", kod, err)
	}
	return trafienieDebaty(wynik)
}

// Uczestnik spoza wskazania zachowuje pozycję: moderator ustawia porządek części składu równie dobrze jak całego.
func (r *repozytoriumRoundtable) UstawKolejnosc(ctx context.Context, okno string, kody []string) error {
	return wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		polecenie, err := r.zapytania.wTransakcji(ctx, transakcja, ustawKolejnoscUczestnika)
		if err != nil {
			return err
		}
		for pozycja, kod := range kody {
			wynik, err := polecenie.ExecContext(ctx, pozycja+1, okno, kod, KontoOperatora(ctx))
			if err != nil {
				return fmt.Errorf("dane: nie można ustawić kolejności uczestnika %q: %w", kod, err)
			}
			if err := trafienieDebaty(wynik); err != nil {
				return fmt.Errorf("dane: uczestnik %q nie stoi w oknie %q: %w", kod, okno, err)
			}
		}
		return nil
	})
}

func odczytajUczestnika(wiersz interface{ Scan(...any) error }) (UczestnikDebaty, error) {
	var uczestnik UczestnikDebaty
	err := wiersz.Scan(&uczestnik.Kod, &uczestnik.Okno, &uczestnik.KanalModelu,
		&uczestnik.NazwaTozsamosci, &uczestnik.PromptSystemowy, &uczestnik.Wyciszony,
		&uczestnik.Kolejnosc, &uczestnik.Kluczowy, &uczestnik.Waga, &uczestnik.Rola,
		&uczestnik.Agent, &uczestnik.Awatar, &uczestnik.OpisRoli, &uczestnik.LiczbaProbek,
		&uczestnik.Utworzono, &uczestnik.Zaktualizowano)
	return uczestnik, err
}

func trafienieDebaty(wynik sql.Result) error {
	zmienione, err := wynik.RowsAffected()
	if err != nil {
		return nil
	}
	if zmienione == 0 {
		return ErrBrakWiersza
	}
	return nil
}
