// Odpowiedzialność pliku: wyciszenia nakładki Always On Display (tabela
// `wyciszenie_nakladki`, migracja 374) oraz sygnały klas zdarzeń wyzwalających
// (tabela `sygnal_nakladki`, migracja 375).
//
// Dlaczego wyciszenie ma wiersz, a przypięcia obserwacji nie mają. Przypięcie
// wskazuje proces telemetrii, który ginie razem z rdzeniem, więc wiersz
// przeżyłby byt, na który wskazuje (`core/adapter_modul_aod.go`). Wyciszenie
// wskazuje moduł, kartę sesji albo klasę zdarzeń — byty, które restart rdzenia
// przeżywają — i ma sięgać WSZYSTKICH powłok Operatora. Bez wiersza Operator
// wyciszał w jednej powłoce, a w drugiej sugestie wchodziły dalej.
//
// Wyciszenie przeterminowane nie wchodzi do wykazu i jest z niego usuwane przy
// odczycie: Operator nie ma odklikiwać ciszy, która sama się skończyła. Usunięcie
// idzie przy odczycie, a nie zegarem w tle, bo wykaz czyta się przed każdym
// ujawnieniem sugestii — a proces budzony co minutę po to, żeby zwykle nie zrobić
// nic, jest kosztem bez skutku.
//
// Sygnał wyciszony ODKŁADA SIĘ NADAL: wyciszenie wstrzymuje ujawnienie, nie
// zapis (rozdz. 3.1 i 3.5 opracowania). Sito wyciszeń stoi po stronie rdzenia,
// nie w tym zapytaniu — tu leży wyłącznie zapis i odczyt.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"danacoconsole/shared"
)

// WyciszenieNakladki to wiersz tabeli `wyciszenie_nakladki`.
type WyciszenieNakladki struct {
	ID            int64
	Identyfikator string
	Rodzaj        shared.AodMuteKind
	Zakres        shared.AodMuteScope
	KluczZakresu  string
	NazwaZakresu  string
	KlasaZdarzen  shared.AodEventClass
	Urzadzenie    string
	// KonczySie jest chwilą końca w zapisie ISO-8601 UTC; pusty znaczy wyciszenie
	// trwające do zniesienia ręką Operatora.
	KonczySie string
	Utworzono string
}

// SygnalNakladki to wiersz tabeli `sygnal_nakladki`.
type SygnalNakladki struct {
	ID              int64
	Identyfikator   string
	KlasaZdarzen    shared.AodEventClass
	Tresc           string
	ModulKod        string
	SesjaKod        string
	LiczbaWystapien int
	ZdarzyloSie     string
}

// RepozytoriumWyciszenNakladki jest kontraktem wyciszeń i sygnałów nakładki.
type RepozytoriumWyciszenNakladki interface {
	// ZapiszWyciszenieNakladki zakłada wyciszenie albo oddaje zastane. Drugi wynik
	// mówi, czy wiersz naprawdę powstał: wyciszenie powtórzone nie jest zmianą
	// i nie ma czego rozgłaszać pozostałym powłokom.
	ZapiszWyciszenieNakladki(ctx context.Context,
		wyciszenie WyciszenieNakladki) (WyciszenieNakladki, bool, error)
	// ZniesWyciszenieNakladki usuwa wyciszenie wskazane identyfikatorem kontraktu.
	ZniesWyciszenieNakladki(ctx context.Context, identyfikator string) (bool, error)
	// WyciszenieNakladkiPoBycie odnajduje wyciszenie złożone z rodzaju i zakresu —
	// drogę zniesienia bez identyfikatora, którą idzie okno znoszące to, co samo
	// wcześniej założyło.
	WyciszenieNakladkiPoBycie(ctx context.Context,
		wzor WyciszenieNakladki) (WyciszenieNakladki, bool, error)
	// WyciszeniaNakladki zwraca wyciszenia czynne o wskazanej chwili, usuwając po
	// drodze te przeterminowane.
	WyciszeniaNakladki(ctx context.Context, teraz string) ([]WyciszenieNakladki, error)

	// ZapiszSygnalNakladki odkłada sygnał klasy zdarzeń wyzwalających.
	ZapiszSygnalNakladki(ctx context.Context, sygnal SygnalNakladki) (SygnalNakladki, error)
	// SygnalyNakladki zwraca sygnały zawężone niepustymi polami wzoru, najświeższe
	// na początku. Granica nieustawiona znaczy wykaz pełny.
	SygnalyNakladki(ctx context.Context, wzor SygnalNakladki, granica int) ([]SygnalNakladki, error)
}

const (
	kolumnyWyciszenia = `id, identyfikator_zewnetrzny, rodzaj, zakres, klucz_zakresu,
	                     nazwa_zakresu, klasa_zdarzen, urzadzenie_id, konczy_sie, utworzono`

	zapiszWyciszenieNakladki = `INSERT INTO wyciszenie_nakladki
	                            (identyfikator_zewnetrzny, rodzaj, zakres, klucz_zakresu,
	                             nazwa_zakresu, klasa_zdarzen, urzadzenie_id, konczy_sie)
	                            VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	pobierzWyciszenieNakladki = `SELECT ` + kolumnyWyciszenia +
		` FROM wyciszenie_nakladki WHERE identyfikator_zewnetrzny = ?`

	pobierzWyciszeniePoBycie = `SELECT ` + kolumnyWyciszenia +
		` FROM wyciszenie_nakladki
		  WHERE rodzaj = ? AND zakres = ? AND klucz_zakresu = ? AND klasa_zdarzen = ?`

	zniesWyciszenieNakladki = `DELETE FROM wyciszenie_nakladki WHERE identyfikator_zewnetrzny = ?`

	// Wyciszenie czasowe, którego chwila końca minęła, przestaje być stanem
	// platformy — znika przy pierwszym odczycie.
	usunWyciszeniaPrzeterminowane = `DELETE FROM wyciszenie_nakladki
	                                 WHERE konczy_sie <> '' AND konczy_sie <= ?`

	listaWyciszenNakladki = `SELECT ` + kolumnyWyciszenia +
		` FROM wyciszenie_nakladki ORDER BY utworzono ASC, id ASC`

	kolumnySygnalu = `id, identyfikator_zewnetrzny, klasa_zdarzen, tresc, modul_kod,
	                  sesja_kod, liczba_wystapien, zdarzylo_sie`

	zapiszSygnalNakladki = `INSERT INTO sygnal_nakladki
	                        (identyfikator_zewnetrzny, klasa_zdarzen, tresc, modul_kod,
	                         sesja_kod, liczba_wystapien)
	                        VALUES (?, ?, ?, ?, ?, ?)`

	pobierzSygnalNakladki = `SELECT ` + kolumnySygnalu +
		` FROM sygnal_nakladki WHERE identyfikator_zewnetrzny = ?`

	listaSygnalowNakladki = `SELECT ` + kolumnySygnalu +
		` FROM sygnal_nakladki
		  WHERE (? = '' OR klasa_zdarzen = ?)
		    AND (? = '' OR modul_kod = ?)
		    AND (? = '' OR sesja_kod = ?)
		  ORDER BY zdarzylo_sie DESC, id DESC
		  LIMIT ?`
)

type repozytoriumWyciszenNakladki struct {
	zapytania *zapytania
}

// noweRepozytoriumWyciszenNakladki zakłada magazyn wyciszeń nad zapytaniami zestawu.
func noweRepozytoriumWyciszenNakladki(z *zapytania) *repozytoriumWyciszenNakladki {
	return &repozytoriumWyciszenNakladki{zapytania: z}
}

// ZapiszWyciszenieNakladki zakłada wyciszenie. Wyciszenie tego samego bytu jest
// już zapisane, więc drugie żądanie oddaje wiersz zastany i mówi, że zmiany nie
// było.
//
// Wyjątek dotyczy wyciszenia czasowego: drugi czas ZASTĘPUJE poprzedni, bo dwa
// czasy naraz nie dałyby Operatorowi jednej odpowiedzi na pytanie „do kiedy".
func (r *repozytoriumWyciszenNakladki) ZapiszWyciszenieNakladki(ctx context.Context,
	wyciszenie WyciszenieNakladki) (WyciszenieNakladki, bool, error) {

	if strings.TrimSpace(wyciszenie.Identyfikator) == "" {
		return WyciszenieNakladki{}, false, fmt.Errorf("dane: wyciszenie nakładki bez identyfikatora")
	}
	zastane, jest, err := r.WyciszenieNakladkiPoBycie(ctx, wyciszenie)
	if err != nil {
		return WyciszenieNakladki{}, false, err
	}
	if jest {
		if wyciszenie.Rodzaj != shared.AodMuteKindTimed || zastane.KonczySie == wyciszenie.KonczySie {
			return zastane, false, nil
		}
		if _, err := r.ZniesWyciszenieNakladki(ctx, zastane.Identyfikator); err != nil {
			return WyciszenieNakladki{}, false, err
		}
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszWyciszenieNakladki)
	if err != nil {
		return WyciszenieNakladki{}, false, err
	}
	_, err = polecenie.ExecContext(ctx, wyciszenie.Identyfikator, string(wyciszenie.Rodzaj),
		string(wyciszenie.Zakres), wyciszenie.KluczZakresu, wyciszenie.NazwaZakresu,
		string(wyciszenie.KlasaZdarzen), wyciszenie.Urzadzenie, wyciszenie.KonczySie)
	if err != nil {
		return WyciszenieNakladki{}, false,
			fmt.Errorf("dane: nie można zapisać wyciszenia nakładki %q: %w",
				wyciszenie.Identyfikator, err)
	}
	polecenie, err = r.zapytania.przygotuj(ctx, pobierzWyciszenieNakladki)
	if err != nil {
		return WyciszenieNakladki{}, false, err
	}
	zapisane, err := odczytajWyciszenieNakladki(
		polecenie.QueryRowContext(ctx, wyciszenie.Identyfikator))
	if err != nil {
		return WyciszenieNakladki{}, false,
			fmt.Errorf("dane: nieczytelne wyciszenie nakładki %q: %w", wyciszenie.Identyfikator, err)
	}
	return zapisane, true, nil
}

// ZniesWyciszenieNakladki usuwa wiersz wyciszenia.
func (r *repozytoriumWyciszenNakladki) ZniesWyciszenieNakladki(ctx context.Context,
	identyfikator string) (bool, error) {

	if strings.TrimSpace(identyfikator) == "" {
		return false, fmt.Errorf("dane: zniesienie wyciszenia nakładki bez identyfikatora")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zniesWyciszenieNakladki)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, identyfikator)
	if err != nil {
		return false, fmt.Errorf("dane: nie można znieść wyciszenia nakładki %q: %w",
			identyfikator, err)
	}
	zmienione, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nieznany skutek zniesienia wyciszenia nakładki %q: %w",
			identyfikator, err)
	}
	return zmienione > 0, nil
}

// WyciszenieNakladkiPoBycie odnajduje wyciszenie po rodzaju i zakresie.
func (r *repozytoriumWyciszenNakladki) WyciszenieNakladkiPoBycie(ctx context.Context,
	wzor WyciszenieNakladki) (WyciszenieNakladki, bool, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzWyciszeniePoBycie)
	if err != nil {
		return WyciszenieNakladki{}, false, err
	}
	wyciszenie, err := odczytajWyciszenieNakladki(polecenie.QueryRowContext(ctx,
		string(wzor.Rodzaj), string(wzor.Zakres), wzor.KluczZakresu, string(wzor.KlasaZdarzen)))
	if errors.Is(err, sql.ErrNoRows) {
		return WyciszenieNakladki{}, false, nil
	}
	if err != nil {
		return WyciszenieNakladki{}, false,
			fmt.Errorf("dane: nieczytelne wyciszenie nakładki: %w", err)
	}
	return wyciszenie, true, nil
}

// WyciszeniaNakladki zwraca wyciszenia czynne, usuwając po drodze przeterminowane.
func (r *repozytoriumWyciszenNakladki) WyciszeniaNakladki(ctx context.Context,
	teraz string) ([]WyciszenieNakladki, error) {

	if teraz != "" {
		polecenie, err := r.zapytania.przygotuj(ctx, usunWyciszeniaPrzeterminowane)
		if err != nil {
			return nil, err
		}
		if _, err := polecenie.ExecContext(ctx, teraz); err != nil {
			return nil, fmt.Errorf("dane: nie można zdjąć przeterminowanych wyciszeń nakładki: %w", err)
		}
	}
	polecenie, err := r.zapytania.przygotuj(ctx, listaWyciszenNakladki)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać wyciszeń nakładki: %w", err)
	}
	defer wiersze.Close()

	wykaz := []WyciszenieNakladki{}
	for wiersze.Next() {
		wyciszenie, err := odczytajWyciszenieNakladki(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz wyciszeń nakładki: %w", err)
		}
		wykaz = append(wykaz, wyciszenie)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt wyciszeń nakładki: %w", err)
	}
	return wykaz, nil
}

// ZapiszSygnalNakladki odkłada sygnał klasy zdarzeń wyzwalających.
func (r *repozytoriumWyciszenNakladki) ZapiszSygnalNakladki(ctx context.Context,
	sygnal SygnalNakladki) (SygnalNakladki, error) {

	if strings.TrimSpace(sygnal.Identyfikator) == "" {
		return SygnalNakladki{}, fmt.Errorf("dane: sygnał nakładki bez identyfikatora")
	}
	wystapienia := sygnal.LiczbaWystapien
	if wystapienia < 1 {
		wystapienia = 1
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszSygnalNakladki)
	if err != nil {
		return SygnalNakladki{}, err
	}
	_, err = polecenie.ExecContext(ctx, sygnal.Identyfikator, string(sygnal.KlasaZdarzen),
		sygnal.Tresc, sygnal.ModulKod, sygnal.SesjaKod, wystapienia)
	if err != nil {
		return SygnalNakladki{}, fmt.Errorf("dane: nie można zapisać sygnału nakładki %q: %w",
			sygnal.Identyfikator, err)
	}
	polecenie, err = r.zapytania.przygotuj(ctx, pobierzSygnalNakladki)
	if err != nil {
		return SygnalNakladki{}, err
	}
	zapisany, err := odczytajSygnalNakladki(polecenie.QueryRowContext(ctx, sygnal.Identyfikator))
	if err != nil {
		return SygnalNakladki{}, fmt.Errorf("dane: nieczytelny sygnał nakładki %q: %w",
			sygnal.Identyfikator, err)
	}
	return zapisany, nil
}

// SygnalyNakladki zwraca sygnały zawężone niepustymi polami wzoru.
func (r *repozytoriumWyciszenNakladki) SygnalyNakladki(ctx context.Context,
	wzor SygnalNakladki, granica int) ([]SygnalNakladki, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaSygnalowNakladki)
	if err != nil {
		return nil, err
	}
	klasa := string(wzor.KlasaZdarzen)
	wiersze, err := polecenie.QueryContext(ctx, klasa, klasa, wzor.ModulKod, wzor.ModulKod,
		wzor.SesjaKod, wzor.SesjaKod, granicaWykazu(granica))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać sygnałów nakładki: %w", err)
	}
	defer wiersze.Close()

	wykaz := []SygnalNakladki{}
	for wiersze.Next() {
		sygnal, err := odczytajSygnalNakladki(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz sygnałów nakładki: %w", err)
		}
		wykaz = append(wykaz, sygnal)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt sygnałów nakładki: %w", err)
	}
	return wykaz, nil
}

// odczytajWyciszenieNakladki składa strukturę z jednego wiersza wyniku.
func odczytajWyciszenieNakladki(wiersz skaner) (WyciszenieNakladki, error) {
	var wyciszenie WyciszenieNakladki
	var rodzaj, zakres, klasa string
	err := wiersz.Scan(&wyciszenie.ID, &wyciszenie.Identyfikator, &rodzaj, &zakres,
		&wyciszenie.KluczZakresu, &wyciszenie.NazwaZakresu, &klasa, &wyciszenie.Urzadzenie,
		&wyciszenie.KonczySie, &wyciszenie.Utworzono)
	if err != nil {
		return WyciszenieNakladki{}, err
	}
	wyciszenie.Rodzaj = shared.AodMuteKind(rodzaj)
	wyciszenie.Zakres = shared.AodMuteScope(zakres)
	wyciszenie.KlasaZdarzen = shared.AodEventClass(klasa)
	return wyciszenie, nil
}

// odczytajSygnalNakladki składa strukturę z jednego wiersza wyniku.
func odczytajSygnalNakladki(wiersz skaner) (SygnalNakladki, error) {
	var sygnal SygnalNakladki
	var klasa string
	err := wiersz.Scan(&sygnal.ID, &sygnal.Identyfikator, &klasa, &sygnal.Tresc,
		&sygnal.ModulKod, &sygnal.SesjaKod, &sygnal.LiczbaWystapien, &sygnal.ZdarzyloSie)
	if err != nil {
		return SygnalNakladki{}, err
	}
	sygnal.KlasaZdarzen = shared.AodEventClass(klasa)
	return sygnal, nil
}
