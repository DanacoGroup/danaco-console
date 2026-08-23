// Odpowiedzialność pliku: obszar modułu Roundtable — byty debaty, kontrakt
// całego obszaru i uczestnik (tabela `debata_uczestnik`
// z `migracja_044_roundtable.sql`). Tury i wypowiedzi leżą w
// `roundtable_tury.go`, stanowisko końcowe w `roundtable_stanowisko.go` — jedno
// repozytorium, trzy pliki wedle odpowiedzialności.
//
// Repozytorium nie rozmawia z modelem. Wywołanie kanałów uczestników prowadzi
// rdzeń przez rejestr kanałów; tutaj leży wyłącznie to, co
// po debacie zostaje: kto brał w niej udział, o co pytano i co odpowiedziano.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// UczestnikDebaty to wiersz tabeli `debata_uczestnik`. Ten sam kanał modelu
// może wystąpić dwukrotnie pod odrębnymi tożsamościami, więc jednoznaczny jest
// wyłącznie Kod — nie para (Okno, KanalModelu).
type UczestnikDebaty struct {
	Kod             string
	Okno            string
	KanalModelu     string
	NazwaTozsamosci *string
	PromptSystemowy *string
	Wyciszony       bool
	Kolejnosc       int
	// Pola z migracji 190 — tożsamość i pozycja uczestnika w naradzie.
	Kluczowy     bool
	Waga         float64
	Rola         string
	Agent        *string
	Awatar       *string
	OpisRoli     *string
	LiczbaProbek int

	Utworzono      string
	Zaktualizowano string
}

// TuraDebaty to wiersz tabeli `debata_tura`. Format i Stan niosą wartości
// kontraktu wprost (RoundtableFormat, RoundtableTurnStatus).
type TuraDebaty struct {
	Kod         string
	Okno        string
	Numer       int
	Zagadnienie *string
	Pytanie     string
	Format      string
	Stan        string
	GranicaTur  int
	// Pola z migracji 191 — wariant tury, wątek boczny i granice tury.
	TuraNadrzedna  string
	GranicaCzasuMs int
	GranicaZnakow  int
	Anonimowa      bool

	Rozpoczeto string
	Zamknieto  *string
}

// WypowiedzDebaty to wiersz tabeli `debata_wypowiedz`. Uczestnik jest kodem,
// nie kluczem obcym — moderator wypowiada się w turze, a uczestnikiem nie jest.
type WypowiedzDebaty struct {
	Kod       string
	TuraKod   string
	Uczestnik string
	Tresc     string
	// Pola z migracji 191. Pewność ujemna znaczy „nie deklarowano”, akt mowy
	// pusty — „jeszcze nieklasyfikowany”.
	OdpowiedzNa string
	AktMowy     string
	Pewnosc     float64
	Redakcja    int

	Utworzono string
}

// StanowiskoDebaty to wiersz tabeli `debata_stanowisko`. Pusta Tura znaczy
// stanowisko całej debaty.
type StanowiskoDebaty struct {
	Kod            string
	Okno           string
	Tura           string
	Tresc          *string
	Wersja         int
	Zaktualizowano string
	// Pola z migracji 198 — stanowisko redagowane przez Operatora wraz
	// z zapisem decyzji.
	Redagowane    bool
	Zaakceptowane bool
	Kontekst      *string
	Warianty      *string
	Konsekwencje  *string
	Tury          string
}

// RepozytoriumRoundtable jest kontraktem obszaru Roundtable.
//
// Obszar urósł ponad jeden plik, więc kontrakt składa się z części: rdzeń
// debaty stoi tutaj, a zdolności dobudowane migracjami 190–199 leżą we własnych
// plikach i wchodzą tu przez zanurzenie. Jedno repozytorium, jeden kontrakt,
// tyle plików, ile odpowiedzialności.
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
	// UstawKolejnosc zapisuje kolejność głosu w jednej transakcji: wskazanie
	// moderatora jest jednym ruchem, więc połowiczny zapis zostawiłby debatę
	// z porządkiem, którego nikt nie wybrał.
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

	zapiszUczestnikaDebaty = `INSERT INTO debata_uczestnik
	                          (identyfikator_zewnetrzny, okno, kanal_modelu, nazwa_tozsamosci,
	                           prompt_systemowy, wyciszony, kolejnosc)
	                          VALUES (?, ?, ?, ?, ?, ?, ?)
	                          ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                              nazwa_tozsamosci = excluded.nazwa_tozsamosci,
	                              prompt_systemowy = excluded.prompt_systemowy,
	                              zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')`

	pobierzUczestnika = `SELECT ` + kolumnyUczestnika + `
	                     FROM debata_uczestnik WHERE identyfikator_zewnetrzny = ?`

	pobierzUczestnikow = `SELECT ` + kolumnyUczestnika + `
	                      FROM debata_uczestnik WHERE okno = ?
	                      ORDER BY kolejnosc ASC, id ASC`

	ustawWyciszenieUczestnika = `UPDATE debata_uczestnik
	                             SET wyciszony = ?,
	                                 zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                             WHERE identyfikator_zewnetrzny = ?`

	ustawKolejnoscUczestnika = `UPDATE debata_uczestnik
	                            SET kolejnosc = ?,
	                                zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                            WHERE okno = ? AND identyfikator_zewnetrzny = ?`
)

type repozytoriumRoundtable struct {
	zapytania *zapytania
	db        *sql.DB
}

func noweRepozytoriumRoundtable(z *zapytania, db *sql.DB) *repozytoriumRoundtable {
	return &repozytoriumRoundtable{zapytania: z, db: db}
}

// ZapiszUczestnika dopisuje uczestnika albo odświeża jego tożsamość. Kolejność
// i wyciszenie zostają nietknięte: ustawia je moderator, a powtórny zapis
// tożsamości nie ma prawa cofnąć jego decyzji.
func (r *repozytoriumRoundtable) ZapiszUczestnika(ctx context.Context,
	uczestnik UczestnikDebaty) (UczestnikDebaty, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, zapiszUczestnikaDebaty)
	if err != nil {
		return UczestnikDebaty{}, err
	}
	if _, err := polecenie.ExecContext(ctx, uczestnik.Kod, uczestnik.Okno, uczestnik.KanalModelu,
		uczestnik.NazwaTozsamosci, uczestnik.PromptSystemowy, uczestnik.Wyciszony,
		uczestnik.Kolejnosc); err != nil {
		return UczestnikDebaty{}, fmt.Errorf("dane: nie można zapisać uczestnika debaty %q: %w",
			uczestnik.Kod, err)
	}
	return r.Uczestnik(ctx, uczestnik.Kod)
}

// Uczestnik zwraca uczestnika po kodzie; brak wiersza wraca jako ErrBrakWiersza.
func (r *repozytoriumRoundtable) Uczestnik(ctx context.Context, kod string) (UczestnikDebaty, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzUczestnika)
	if err != nil {
		return UczestnikDebaty{}, err
	}
	uczestnik, err := odczytajUczestnika(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return UczestnikDebaty{}, ErrBrakWiersza
	}
	if err != nil {
		return UczestnikDebaty{}, fmt.Errorf("dane: nieczytelny wiersz uczestnika debaty %q: %w", kod, err)
	}
	return uczestnik, nil
}

// Uczestnicy zwraca skład debaty okna w kolejności głosu.
func (r *repozytoriumRoundtable) Uczestnicy(ctx context.Context, okno string) ([]UczestnikDebaty, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzUczestnikow)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno)
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

// UstawWyciszenie przestawia wyciszenie uczestnika w turze.
func (r *repozytoriumRoundtable) UstawWyciszenie(ctx context.Context, kod string, wyciszony bool) error {
	polecenie, err := r.zapytania.przygotuj(ctx, ustawWyciszenieUczestnika)
	if err != nil {
		return err
	}
	wynik, err := polecenie.ExecContext(ctx, wyciszony, kod)
	if err != nil {
		return fmt.Errorf("dane: nie można przestawić wyciszenia uczestnika %q: %w", kod, err)
	}
	return trafienieDebaty(wynik)
}

// UstawKolejnosc zapisuje kolejność głosu wskazaną przez moderatora. Uczestnik
// spoza wskazania zachowuje swoją pozycję — moderator ustawia porządek części
// składu równie dobrze jak całego.
func (r *repozytoriumRoundtable) UstawKolejnosc(ctx context.Context, okno string, kody []string) error {
	return wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		polecenie, err := r.zapytania.wTransakcji(ctx, transakcja, ustawKolejnoscUczestnika)
		if err != nil {
			return err
		}
		for pozycja, kod := range kody {
			if _, err := polecenie.ExecContext(ctx, pozycja+1, okno, kod); err != nil {
				return fmt.Errorf("dane: nie można ustawić kolejności uczestnika %q: %w", kod, err)
			}
		}
		return nil
	})
}

// odczytajUczestnika składa uczestnika z jednego wiersza wyniku.
func odczytajUczestnika(wiersz interface{ Scan(...any) error }) (UczestnikDebaty, error) {
	var uczestnik UczestnikDebaty
	err := wiersz.Scan(&uczestnik.Kod, &uczestnik.Okno, &uczestnik.KanalModelu,
		&uczestnik.NazwaTozsamosci, &uczestnik.PromptSystemowy, &uczestnik.Wyciszony,
		&uczestnik.Kolejnosc, &uczestnik.Kluczowy, &uczestnik.Waga, &uczestnik.Rola,
		&uczestnik.Agent, &uczestnik.Awatar, &uczestnik.OpisRoli, &uczestnik.LiczbaProbek,
		&uczestnik.Utworzono, &uczestnik.Zaktualizowano)
	return uczestnik, err
}

// trafienieDebaty odróżnia zapis, który nic nie zmienił, od zapisu udanego.
// Bez tego przestawienie stanu bytu, którego nie ma, kończyłoby się cicho.
//
// Osobno od `sprawdzTrafienie` z `sesje.go`, bo tamten opisuje wiersz numerem
// klucza głównego i zwraca błąd opisowy, a obszar Roundtable rozpoznaje brak
// bytu przez `errors.Is(err, ErrBrakWiersza)` — rdzeń oddaje wtedy kod `not_found`
// zamiast usterki wewnętrznej.
func trafienieDebaty(wynik sql.Result) error {
	zmienione, err := wynik.RowsAffected()
	if err != nil {
		return nil // sterownik bez licznika — brak liczby nie jest błędem zapisu
	}
	if zmienione == 0 {
		return ErrBrakWiersza
	}
	return nil
}
