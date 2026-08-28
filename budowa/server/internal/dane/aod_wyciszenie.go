// Plik zapisuje i odczytuje wyciszenia nakładki Always On Display oraz sygnały klas zdarzeń wyzwalających; wyciszenie
// wskazuje moduł, kartę sesji albo klasę zdarzeń i sięga wszystkich powłok Operatora.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"danacoconsole/shared"
)

// WyciszenieNakladki to wiersz tabeli `wyciszenie_nakladki` wraz z chwilą końca i rodzajem oraz zakresem zniesienia.
type WyciszenieNakladki struct {
	ID            int64
	Identyfikator string
	Rodzaj        shared.AodMuteKind
	Zakres        shared.AodMuteScope
	KluczZakresu  string
	NazwaZakresu  string
	KlasaZdarzen  shared.AodEventClass
	Urzadzenie    string
	// KonczySie jest chwilą końca w zapisie ISO-8601 UTC; pusty znaczy trwanie do zniesienia.
	KonczySie string
	Utworzono string
}

// SygnalNakladki to wiersz tabeli `sygnal_nakladki` niosący klasę zdarzenia wyzwalającego nakładkę Operatora.
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

// RepozytoriumWyciszenNakladki jest kontraktem wyciszeń i sygnałów nakładki: zapis, odczyt i zniesienie.
type RepozytoriumWyciszenNakladki interface {
	// ZapiszWyciszenieNakladki zakłada wyciszenie albo oddaje zastane, mówiąc, czy wiersz powstał.
	ZapiszWyciszenieNakladki(ctx context.Context,
		wyciszenie WyciszenieNakladki) (WyciszenieNakladki, bool, error)
	// ZniesWyciszenieNakladki usuwa wyciszenie wskazane identyfikatorem kontraktu.
	ZniesWyciszenieNakladki(ctx context.Context, identyfikator string) (bool, error)
	// WyciszenieNakladkiPoBycie odnajduje wyciszenie złożone z rodzaju i zakresu, bez identyfikatora.
	WyciszenieNakladkiPoBycie(ctx context.Context,
		wzor WyciszenieNakladki) (WyciszenieNakladki, bool, error)
	// WyciszeniaNakladki zwraca wyciszenia czynne o wskazanej chwili, usuwając po drodze przeterminowane.
	WyciszeniaNakladki(ctx context.Context, teraz string) ([]WyciszenieNakladki, error)

	// ZapiszSygnalNakladki odkłada sygnał klasy zdarzeń wyzwalających.
	ZapiszSygnalNakladki(ctx context.Context, sygnal SygnalNakladki) (SygnalNakladki, error)
	// SygnalyNakladki zwraca sygnały zawężone niepustymi polami wzoru, najświeższe na początku.
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

// noweRepozytoriumWyciszenNakladki zakłada magazyn wyciszeń i sygnałów nakładki nad zapytaniami zestawu.
func noweRepozytoriumWyciszenNakladki(z *zapytania) *repozytoriumWyciszenNakladki {
	return &repozytoriumWyciszenNakladki{zapytania: z}
}

// ZapiszWyciszenieNakladki zakłada wyciszenie; drugie żądanie tego samego bytu oddaje wiersz zastany, wyjąwszy czas trwania, który drugie żądanie zastępuje.
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

// ZniesWyciszenieNakladki usuwa z bazy danych wiersz wyciszenia nakładki po jego identyfikatorze trwałym.
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

// WyciszenieNakladkiPoBycie odnajduje wyciszenie złożone z rodzaju i zakresu bytu, bez jego identyfikatora.
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

// WyciszeniaNakladki zwraca wyciszenia czynne o wskazanej chwili, usuwając po drodze wpisy przeterminowane.
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

// ZapiszSygnalNakladki odkłada sygnał klasy zdarzeń wyzwalających nakładkę wraz z chwilą jego wystąpienia.
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

// SygnalyNakladki zwraca sygnały zawężone niepustymi polami wzoru wyszukiwania klasy zdarzeń nakładki.
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

// odczytajWyciszenieNakladki składa strukturę wyciszenia wprost z jednego wiersza wyniku zapytania SQL.
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

// odczytajSygnalNakladki składa strukturę sygnału nakładki wprost z jednego wiersza wyniku zapytania do bazy.
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
