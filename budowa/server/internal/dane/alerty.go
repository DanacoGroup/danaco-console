// Plik definiuje magazyn reguł wyzwalania i rejestru wyzwoleń rodziny alert.
// Repozytorium nie ewaluuje reguł; ewaluacja należy do adaptera.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// RegulaAlertu to wiersz tabeli regula_alertu: definicja progu, miary i
// kanałów powiadomień, wraz z licznikiem dotychczasowych wyzwoleń.
type RegulaAlertu struct {
	ID                 int64
	Kod                string
	Nazwa              string
	Opis               *string
	Rodzaj             string
	Miara              string
	Porownanie         *string
	Prog               *float64
	OknoMs             int64
	Waga               string
	KanalyJSON         string
	Zasieg             *string
	ZasiegKod          *string
	Czynna             bool
	WyciszonaDo        *int64
	EskalacjaPoMs      *int64
	AdresZwrotny       *string
	Utworzono          int64
	Zaktualizowano     *int64
	OstatnieWyzwolenie *int64
	LiczbaWyzwolen     int
}

// WyzwolenieAlertu to wiersz tabeli `wyzwolenie_alertu` — jedno zawołanie
// produktu wraz z wartością zmierzoną w chwili zawołania.
type WyzwolenieAlertu struct {
	Kod                   string
	RegulaKod             string
	RegulaNazwa           *string
	Stan                  string
	Waga                  string
	Miara                 string
	WartoscObserwowana    float64
	Prog                  *float64
	Komunikat             string
	Wyzwolono             int64
	Potwierdzono          *int64
	Rozwiazano            *int64
	Notatka               *string
	BladKod               *string
	SondaKod              *string
	WywolanieKod          *string
	KanalyDostarczoneJSON string
}

// SitoWyzwolenAlertu zawęża odczyt rejestru wyzwoleń po regule, stanie,
// wadze i przedziale czasu, z granicą liczby wierszy.
type SitoWyzwolenAlertu struct {
	RegulaKod string
	Stan      string
	Waga      string
	OdCzasu   int64
	DoCzasu   int64
	Granica   int
}

// RepozytoriumAlertow jest kontraktem magazynu reguł i wyzwoleń: zapis,
// odczyt, potwierdzenie oraz pomiary potrzebne ewaluacji reguł.
type RepozytoriumAlertow interface {
	ZapiszRegule(ctx context.Context, regula RegulaAlertu) (RegulaAlertu, bool, error)
	Regula(ctx context.Context, kod string) (RegulaAlertu, error)
	Reguly(ctx context.Context, rodzaj, miara string, tylkoCzynne bool, granica int) ([]RegulaAlertu, error)
	UsunRegule(ctx context.Context, kod string) (int, error)
	ZapiszWyzwolenie(ctx context.Context, wyzwolenie WyzwolenieAlertu) (WyzwolenieAlertu, error)
	Wyzwolenie(ctx context.Context, kod string) (WyzwolenieAlertu, error)
	Wyzwolenia(ctx context.Context, sito SitoWyzwolenAlertu) ([]WyzwolenieAlertu, int, error)
	PotwierdzWyzwolenie(ctx context.Context, kod string, chwila int64, notatka *string) (WyzwolenieAlertu, error)
	// Dwie miary, których nie da się wziąć z magazynu śladu wywołań ani z
	// dziennika błędów.
	LiczbaNieudanychPomiarowSond(ctx context.Context, od, do int64) (int, error)
	LiczbaNieudanychPozycjiKolejki(ctx context.Context, od, do int64) (int, error)
}

const (
	kolumnyReguluAlertu = `r.id, r.identyfikator_zewnetrzny, r.nazwa, r.opis, r.rodzaj, r.miara,
	                       r.porownanie, r.prog, r.okno_ms, r.waga, r.kanaly_json, r.zasieg,
	                       r.zasieg_kod, r.czynna, r.wyciszona_do, r.eskalacja_po_ms,
	                       r.adres_zwrotny, r.utworzono, r.zaktualizowano, r.ostatnie_wyzwolenie,
	                       (SELECT COUNT(*) FROM wyzwolenie_alertu w
	                         WHERE w.regula_kod = r.identyfikator_zewnetrzny)`

	wstawReguleAlertu = `INSERT INTO regula_alertu
	                     (identyfikator_zewnetrzny, nazwa, opis, rodzaj, miara, porownanie, prog,
	                      okno_ms, waga, kanaly_json, zasieg, zasieg_kod, czynna, wyciszona_do,
	                      eskalacja_po_ms, adres_zwrotny, utworzono, zaktualizowano, konto_id)
	                     VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ` +
		WskazanieKonta + `)`

	aktualizujReguleAlertu = `UPDATE regula_alertu SET
	                              nazwa = ?, opis = ?, rodzaj = ?, miara = ?, porownanie = ?, prog = ?,
	                              okno_ms = ?, waga = ?, kanaly_json = ?, zasieg = ?, zasieg_kod = ?,
	                              czynna = ?, wyciszona_do = ?, eskalacja_po_ms = ?, adres_zwrotny = ?,
	                              zaktualizowano = ?
	                          WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	usunReguleAlertu = `DELETE FROM regula_alertu
	                    WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	policzWyzwoleniaReguly = `SELECT COUNT(*) FROM wyzwolenie_alertu WHERE regula_kod = ?`

	kolumnyWyzwoleniaAlertu = `w.identyfikator_zewnetrzny, w.regula_kod,
	                           (SELECT r.nazwa FROM regula_alertu r
	                             WHERE r.identyfikator_zewnetrzny = w.regula_kod),
	                           w.stan, w.waga, w.miara, w.wartosc_obserwowana, w.prog, w.komunikat,
	                           w.wyzwolono, w.potwierdzono, w.rozwiazano, w.notatka, w.blad_kod,
	                           w.sonda_kod, w.wywolanie_kod, w.kanaly_dostarczone_json`

	wstawWyzwolenieAlertu = `INSERT INTO wyzwolenie_alertu
	                         (identyfikator_zewnetrzny, regula_kod, stan, waga, miara,
	                          wartosc_obserwowana, prog, komunikat, wyzwolono, potwierdzono,
	                          rozwiazano, notatka, blad_kod, sonda_kod, wywolanie_kod,
	                          kanaly_dostarczone_json)
	                         VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	odbijWyzwolenieReguly = `UPDATE regula_alertu SET ostatnie_wyzwolenie = ?
	                         WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta
)

// Wyzwolenie nie ma kolumny konto_id: konto sięga przez regułę wskazaną kodem
// zewnętrznym, jednoznacznym w całej tabeli reguł (migracja 484).
var (
	warunekKontaReguly = strings.ReplaceAll(WarunekKonta, "konto_id", "r.konto_id")

	kontoWyzwoleniaAlertu = ` regula_kod IN (SELECT k.identyfikator_zewnetrzny FROM regula_alertu k
	                                         WHERE ` +
		strings.ReplaceAll(WarunekKonta, "konto_id", "k.konto_id") + `)`

	pobierzReguleAlertu = `SELECT ` + kolumnyReguluAlertu +
		` FROM regula_alertu r WHERE r.identyfikator_zewnetrzny = ? AND ` + warunekKontaReguly

	pobierzWyzwolenieAlertu = `SELECT ` + kolumnyWyzwoleniaAlertu +
		` FROM wyzwolenie_alertu w WHERE w.identyfikator_zewnetrzny = ? AND ` + kontoWyzwoleniaAlertu

	potwierdzWyzwolenieAlertu = `UPDATE wyzwolenie_alertu
	                             SET stan = 'acknowledged', potwierdzono = ?,
	                                 notatka = COALESCE(?, notatka)
	                             WHERE identyfikator_zewnetrzny = ? AND ` + kontoWyzwoleniaAlertu
)

type repozytoriumAlertow struct {
	zapytania *zapytania
	db        *sql.DB
}

// noweRepozytoriumAlertow zakłada magazyn reguł i wyzwoleń nad bazą
// wskazanego zestawu repozytoriów danych.
func noweRepozytoriumAlertow(z *zapytania, db *sql.DB) *repozytoriumAlertow {
	return &repozytoriumAlertow{zapytania: z, db: db}
}

// ZapiszRegule zakłada regułę albo nadpisuje zastaną po kodzie zewnętrznym,
// zwracając stan reguły i informację, czy powstała.
func (r *repozytoriumAlertow) ZapiszRegule(ctx context.Context,
	regula RegulaAlertu) (RegulaAlertu, bool, error) {

	if strings.TrimSpace(regula.Kod) == "" {
		return RegulaAlertu{}, false, fmt.Errorf("dane: reguła alertu bez identyfikatora")
	}
	powstala := false
	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		odczyt, err := r.zapytania.wTransakcji(ctx, transakcja, pobierzReguleAlertu)
		if err != nil {
			return err
		}
		zastana, err := odczytajReguleAlertu(odczyt.QueryRowContext(ctx, regula.Kod, KontoOperatora(ctx)))
		switch {
		case errors.Is(err, sql.ErrNoRows):
			powstala = true
		case err != nil:
			return fmt.Errorf("dane: nie można odczytać reguły alertu %q: %w", regula.Kod, err)
		default:
			regula.Utworzono = zastana.Utworzono
			regula.OstatnieWyzwolenie = zastana.OstatnieWyzwolenie
		}

		tekst := aktualizujReguleAlertu
		if powstala {
			tekst = wstawReguleAlertu
		}
		zapis, err := r.zapytania.wTransakcji(ctx, transakcja, tekst)
		if err != nil {
			return err
		}
		wspolne := []any{regula.Nazwa, tekstDoKolumny(regula.Opis), regula.Rodzaj, regula.Miara,
			tekstDoKolumny(regula.Porownanie), liczbaRzeczywistaDoKolumny(regula.Prog),
			regula.OknoMs, regula.Waga, regula.KanalyJSON, tekstDoKolumny(regula.Zasieg),
			tekstDoKolumny(regula.ZasiegKod), liczbaLogiczna(regula.Czynna),
			liczbaDoKolumny(regula.WyciszonaDo), liczbaDoKolumny(regula.EskalacjaPoMs),
			tekstDoKolumny(regula.AdresZwrotny)}
		var argumenty []any
		if powstala {
			argumenty = append([]any{regula.Kod}, wspolne...)
			argumenty = append(argumenty, regula.Utworzono,
				liczbaDoKolumny(regula.Zaktualizowano), KontoOperatora(ctx))
		} else {
			argumenty = append(wspolne, liczbaDoKolumny(regula.Zaktualizowano),
				regula.Kod, KontoOperatora(ctx))
		}
		wynik, err := zapis.ExecContext(ctx, argumenty...)
		// Odczyt zawężony nie widzi reguły cudzego konta, więc zakładanie trafia
		// na jednoznaczność kodu pilnowaną przez bazę.
		if czyKolizja(err) {
			return fmt.Errorf("dane: reguła alertu %q należy do innego konta: %w",
				regula.Kod, ErrKolizjaWiersza)
		}
		if err != nil {
			return fmt.Errorf("dane: nie można zapisać reguły alertu %q: %w", regula.Kod, err)
		}
		if zmienione, err := wynik.RowsAffected(); !powstala && err == nil && zmienione == 0 {
			return fmt.Errorf("dane: reguła alertu %q nie istnieje: %w", regula.Kod, ErrBrakWiersza)
		}
		return nil
	})
	if err != nil {
		return RegulaAlertu{}, false, err
	}
	zapisana, err := r.Regula(ctx, regula.Kod)
	return zapisana, powstala, err
}

// Regula zwraca jedną regułę wraz z licznikiem jej dotychczasowych
// wyzwoleń, liczonym osobnym zapytaniem.
func (r *repozytoriumAlertow) Regula(ctx context.Context, kod string) (RegulaAlertu, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzReguleAlertu)
	if err != nil {
		return RegulaAlertu{}, err
	}
	regula, err := odczytajReguleAlertu(polecenie.QueryRowContext(ctx, kod, KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return RegulaAlertu{}, fmt.Errorf("dane: reguła alertu %q nie istnieje: %w", kod, ErrBrakWiersza)
	}
	return regula, err
}

// Reguly zwraca reguły spełniające zawężenie rodzaju i miary wraz z
// licznikiem wyzwoleń każdej z nich.
func (r *repozytoriumAlertow) Reguly(ctx context.Context, rodzaj, miara string,
	tylkoCzynne bool, granica int) ([]RegulaAlertu, error) {

	warunki := []string{"1 = 1"}
	argumenty := []any{}
	if strings.TrimSpace(rodzaj) != "" {
		warunki = append(warunki, "r.rodzaj = ?")
		argumenty = append(argumenty, rodzaj)
	}
	if strings.TrimSpace(miara) != "" {
		warunki = append(warunki, "r.miara = ?")
		argumenty = append(argumenty, miara)
	}
	if tylkoCzynne {
		warunki = append(warunki, "r.czynna = 1")
	}
	warunki = append(warunki, warunekKontaReguly)
	argumenty = append(argumenty, KontoOperatora(ctx))
	tekst := `SELECT ` + kolumnyReguluAlertu + ` FROM regula_alertu r WHERE ` +
		strings.Join(warunki, " AND ") + ` ORDER BY r.id`
	if granica > 0 {
		tekst += fmt.Sprintf(" LIMIT %d", granica)
	}
	wiersze, err := r.db.QueryContext(ctx, tekst, argumenty...)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać reguł alertu: %w", err)
	}
	defer wiersze.Close()

	lista := []RegulaAlertu{}
	for wiersze.Next() {
		regula, err := odczytajReguleAlertu(wiersze)
		if err != nil {
			return nil, err
		}
		lista = append(lista, regula)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt reguł alertu: %w", err)
	}
	return lista, nil
}

// UsunRegule wykreśla regułę wraz z jej rejestrem wyzwoleń, oddając liczbę
// usuniętych wpisów rejestru.
func (r *repozytoriumAlertow) UsunRegule(ctx context.Context, kod string) (int, error) {
	usunietych := 0
	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		licznik, err := r.zapytania.wTransakcji(ctx, transakcja, policzWyzwoleniaReguly)
		if err != nil {
			return err
		}
		if err := licznik.QueryRowContext(ctx, kod).Scan(&usunietych); err != nil {
			return fmt.Errorf("dane: nie można policzyć wyzwoleń reguły %q: %w", kod, err)
		}
		kasowanie, err := r.zapytania.wTransakcji(ctx, transakcja, usunReguleAlertu)
		if err != nil {
			return err
		}
		wynik, err := kasowanie.ExecContext(ctx, kod, KontoOperatora(ctx))
		if err != nil {
			return fmt.Errorf("dane: nie można usunąć reguły alertu %q: %w", kod, err)
		}
		zmienione, err := wynik.RowsAffected()
		if err == nil && zmienione == 0 {
			return fmt.Errorf("dane: reguła alertu %q nie istnieje: %w", kod, ErrBrakWiersza)
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return usunietych, nil
}

// ZapiszWyzwolenie dopisuje wyzwolenie i odnotowuje jego chwilę w regule, w
// jednej transakcji bazy danych.
func (r *repozytoriumAlertow) ZapiszWyzwolenie(ctx context.Context,
	w WyzwolenieAlertu) (WyzwolenieAlertu, error) {

	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		zapis, err := r.zapytania.wTransakcji(ctx, transakcja, wstawWyzwolenieAlertu)
		if err != nil {
			return err
		}
		if _, err := zapis.ExecContext(ctx, w.Kod, w.RegulaKod, w.Stan, w.Waga, w.Miara,
			w.WartoscObserwowana, liczbaRzeczywistaDoKolumny(w.Prog), w.Komunikat, w.Wyzwolono,
			liczbaDoKolumny(w.Potwierdzono), liczbaDoKolumny(w.Rozwiazano),
			tekstDoKolumny(w.Notatka), tekstDoKolumny(w.BladKod), tekstDoKolumny(w.SondaKod),
			tekstDoKolumny(w.WywolanieKod), w.KanalyDostarczoneJSON); err != nil {
			return fmt.Errorf("dane: nie można zapisać wyzwolenia alertu %q: %w", w.Kod, err)
		}
		odbicie, err := r.zapytania.wTransakcji(ctx, transakcja, odbijWyzwolenieReguly)
		if err != nil {
			return err
		}
		if _, err := odbicie.ExecContext(ctx, w.Wyzwolono, w.RegulaKod, KontoOperatora(ctx)); err != nil {
			return fmt.Errorf("dane: nie można odnotować wyzwolenia reguły %q: %w", w.RegulaKod, err)
		}
		return nil
	})
	if err != nil {
		return WyzwolenieAlertu{}, err
	}
	return r.Wyzwolenie(ctx, w.Kod)
}

// Wyzwolenie zwraca jeden wpis rejestru wyzwoleń o wskazanym kodzie,
// zwracając błąd, gdy nie istnieje.
func (r *repozytoriumAlertow) Wyzwolenie(ctx context.Context, kod string) (WyzwolenieAlertu, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzWyzwolenieAlertu)
	if err != nil {
		return WyzwolenieAlertu{}, err
	}
	wyzwolenie, err := odczytajWyzwolenieAlertu(polecenie.QueryRowContext(ctx, kod, KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return WyzwolenieAlertu{}, fmt.Errorf("dane: wyzwolenie alertu %q nie istnieje: %w",
			kod, ErrBrakWiersza)
	}
	return wyzwolenie, err
}

// Wyzwolenia zwraca rejestr spełniający zawężenie, od najnowszego, wraz
// z liczbą wierszy bez granicy limitu.
func (r *repozytoriumAlertow) Wyzwolenia(ctx context.Context,
	sito SitoWyzwolenAlertu) ([]WyzwolenieAlertu, int, error) {

	warunki := []string{"1 = 1"}
	argumenty := []any{}
	if strings.TrimSpace(sito.RegulaKod) != "" {
		warunki = append(warunki, "w.regula_kod = ?")
		argumenty = append(argumenty, sito.RegulaKod)
	}
	if strings.TrimSpace(sito.Stan) != "" {
		warunki = append(warunki, "w.stan = ?")
		argumenty = append(argumenty, sito.Stan)
	}
	if strings.TrimSpace(sito.Waga) != "" {
		warunki = append(warunki, "w.waga = ?")
		argumenty = append(argumenty, sito.Waga)
	}
	if sito.OdCzasu > 0 {
		warunki = append(warunki, "w.wyzwolono >= ?")
		argumenty = append(argumenty, sito.OdCzasu)
	}
	if sito.DoCzasu > 0 {
		warunki = append(warunki, "w.wyzwolono <= ?")
		argumenty = append(argumenty, sito.DoCzasu)
	}
	warunki = append(warunki, kontoWyzwoleniaAlertu)
	argumenty = append(argumenty, KontoOperatora(ctx))
	gdzie := strings.Join(warunki, " AND ")

	wszystkich := 0
	if err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM wyzwolenie_alertu w WHERE `+gdzie, argumenty...).
		Scan(&wszystkich); err != nil {
		return nil, 0, fmt.Errorf("dane: nie można policzyć wyzwoleń alertu: %w", err)
	}

	tekst := `SELECT ` + kolumnyWyzwoleniaAlertu + ` FROM wyzwolenie_alertu w WHERE ` + gdzie +
		` ORDER BY w.wyzwolono DESC, w.id DESC`
	if sito.Granica > 0 {
		tekst += fmt.Sprintf(" LIMIT %d", sito.Granica)
	}
	wiersze, err := r.db.QueryContext(ctx, tekst, argumenty...)
	if err != nil {
		return nil, 0, fmt.Errorf("dane: nie można odczytać wyzwoleń alertu: %w", err)
	}
	defer wiersze.Close()

	lista := []WyzwolenieAlertu{}
	for wiersze.Next() {
		wyzwolenie, err := odczytajWyzwolenieAlertu(wiersze)
		if err != nil {
			return nil, 0, err
		}
		lista = append(lista, wyzwolenie)
	}
	if err := wiersze.Err(); err != nil {
		return nil, 0, fmt.Errorf("dane: przerwany odczyt wyzwoleń alertu: %w", err)
	}
	return lista, wszystkich, nil
}

// PotwierdzWyzwolenie odnotowuje obsłużenie wyzwolenia. Notatka pusta zostawia
// notatkę zastaną — potwierdzenie bez słowa nie ma prawa skasować powodu
// zapisanego wcześniej.
func (r *repozytoriumAlertow) PotwierdzWyzwolenie(ctx context.Context, kod string,
	chwila int64, notatka *string) (WyzwolenieAlertu, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, potwierdzWyzwolenieAlertu)
	if err != nil {
		return WyzwolenieAlertu{}, err
	}
	wynik, err := polecenie.ExecContext(ctx, chwila, tekstDoKolumny(notatka), kod, KontoOperatora(ctx))
	if err != nil {
		return WyzwolenieAlertu{}, fmt.Errorf("dane: nie można potwierdzić wyzwolenia %q: %w", kod, err)
	}
	if zmienione, err := wynik.RowsAffected(); err == nil && zmienione == 0 {
		return WyzwolenieAlertu{}, fmt.Errorf("dane: wyzwolenie alertu %q nie istnieje: %w",
			kod, ErrBrakWiersza)
	}
	return r.Wyzwolenie(ctx, kod)
}

// odczytajReguleAlertu przekłada wiersz zapytania na regułę, zamieniając
// kolumny nullowalne na wskaźniki.
func odczytajReguleAlertu(s skaner) (RegulaAlertu, error) {
	var regula RegulaAlertu
	var opis, porownanie, zasieg, zasiegKod, adres sql.NullString
	var prog sql.NullFloat64
	var wyciszona, eskalacja, zaktualizowano, ostatnie sql.NullInt64
	var czynna int
	err := s.Scan(&regula.ID, &regula.Kod, &regula.Nazwa, &opis, &regula.Rodzaj, &regula.Miara,
		&porownanie, &prog, &regula.OknoMs, &regula.Waga, &regula.KanalyJSON, &zasieg, &zasiegKod,
		&czynna, &wyciszona, &eskalacja, &adres, &regula.Utworzono, &zaktualizowano, &ostatnie,
		&regula.LiczbaWyzwolen)
	if err != nil {
		return RegulaAlertu{}, err
	}
	regula.Opis = tekstZKolumny(opis)
	regula.Porownanie = tekstZKolumny(porownanie)
	regula.Prog = liczbaRzeczywistaZKolumny(prog)
	regula.Zasieg = tekstZKolumny(zasieg)
	regula.ZasiegKod = tekstZKolumny(zasiegKod)
	regula.Czynna = czynna == 1
	regula.WyciszonaDo = liczbaZKolumny(wyciszona)
	regula.EskalacjaPoMs = liczbaZKolumny(eskalacja)
	regula.AdresZwrotny = tekstZKolumny(adres)
	regula.Zaktualizowano = liczbaZKolumny(zaktualizowano)
	regula.OstatnieWyzwolenie = liczbaZKolumny(ostatnie)
	return regula, nil
}

// odczytajWyzwolenieAlertu przekłada wiersz zapytania na wyzwolenie,
// zamieniając kolumny nullowalne na wskaźniki.
func odczytajWyzwolenieAlertu(s skaner) (WyzwolenieAlertu, error) {
	var w WyzwolenieAlertu
	var nazwa, notatka, blad, sonda, wywolanie sql.NullString
	var prog sql.NullFloat64
	var potwierdzono, rozwiazano sql.NullInt64
	err := s.Scan(&w.Kod, &w.RegulaKod, &nazwa, &w.Stan, &w.Waga, &w.Miara, &w.WartoscObserwowana,
		&prog, &w.Komunikat, &w.Wyzwolono, &potwierdzono, &rozwiazano, &notatka, &blad, &sonda,
		&wywolanie, &w.KanalyDostarczoneJSON)
	if err != nil {
		return WyzwolenieAlertu{}, err
	}
	w.RegulaNazwa = tekstZKolumny(nazwa)
	w.Prog = liczbaRzeczywistaZKolumny(prog)
	w.Potwierdzono = liczbaZKolumny(potwierdzono)
	w.Rozwiazano = liczbaZKolumny(rozwiazano)
	w.Notatka = tekstZKolumny(notatka)
	w.BladKod = tekstZKolumny(blad)
	w.SondaKod = tekstZKolumny(sonda)
	w.WywolanieKod = tekstZKolumny(wywolanie)
	return w, nil
}

// LiczbaNieudanychPomiarowSond liczy pomiary kondycji zakończone stanem innym
// niż sprawny w zadanym oknie czasu. To materiał miary `probeFailure`.
func (r *repozytoriumAlertow) LiczbaNieudanychPomiarowSond(ctx context.Context,
	od, do int64) (int, error) {

	liczba := 0
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM wynik_sondy_kondycji
		  WHERE stan IN ('down','degraded') AND wykonano >= ? AND wykonano <= ?`,
		od, do).Scan(&liczba)
	if err != nil {
		return 0, fmt.Errorf("dane: nie można policzyć nieudanych pomiarów sond: %w", err)
	}
	return liczba, nil
}

// LiczbaNieudanychPozycjiKolejki liczy pozycje kolejek zakończone stanem
// błędnym w zadanym oknie czasu. Porównanie napisów jest tu poprawne, bo
// znacznik ISO-8601 rośnie leksykalnie razem z czasem.
func (r *repozytoriumAlertow) LiczbaNieudanychPozycjiKolejki(ctx context.Context,
	od, do int64) (int, error) {

	liczba := 0
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM pozycja_kolejki
		  WHERE stan = 'bledna'
		    AND zaktualizowano >= strftime('%Y-%m-%dT%H:%M:%fZ', ? / 1000.0, 'unixepoch')
		    AND zaktualizowano <= strftime('%Y-%m-%dT%H:%M:%fZ', ? / 1000.0, 'unixepoch')`,
		od, do).Scan(&liczba)
	if err != nil {
		return 0, fmt.Errorf("dane: nie można policzyć nieudanych pozycji kolejki: %w", err)
	}
	return liczba, nil
}
