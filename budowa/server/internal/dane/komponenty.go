// Odpowiedzialność pliku: rejestr komponentów własnych Strefy 2 Strony głównej
// (tabela `komponent`) — trwałość rodziny `component.*`.
//
// Wiersz `komponent` jest kaflem Strefy 2, który wskazuje byt magazynu
// modułowego kolumną `byt_docelowy`. Nie powiela bytu modułowego i nie jest jego
// drugą prawdą: kroki automatyki, umiejętności eksperta i pamięć projektu
// zostają w swoich tabelach, a to repozytorium nie tyka żadnej z nich.
//
// Konstruktor bierze współdzieloną pamięć zapytań zestawu, tak jak pozostałe
// repozytoria pakietu — `Zestaw.Zamknij` zwalnia wyłącznie tę jedną pamięć
// poleceń. `*sql.DB` konstruktor nie bierze i brać nie musi: żaden zapis tego
// rejestru nie obejmuje drugiej tabeli, więc transakcji wielotabelowej tu nie ma.
package dane

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
)

// Komponent to wiersz tabeli `komponent`. `Kod` jest identyfikatorem trwałym
// i odpowiada polu `Component.id` kontraktu, a `BytDocelowy` — polu
// `Component.targetId`.
type Komponent struct {
	ID             int64
	Kod            string
	Rodzaj         string
	Nazwa          string
	Opis           *string
	BytDocelowy    *string
	Czynny         bool
	Konfiguracja   string
	PoziomZasiegu  *string
	KluczZasiegu   string
	Utworzono      int64
	Zaktualizowano int64
}

// FiltrKomponentow zawęża wykaz — obsługuje oba pola żądania `component.list`.
type FiltrKomponentow struct {
	// Rodzaj pusty znaczy „wszystkie rodzaje”.
	Rodzaj string
	// DolaczWylaczone otwiera wykaz na komponenty niczynne. Domyślnie zamknięty:
	// kontrakt oznacza `includeDisabled` jako niewymagane, a Strefa 2 pokazuje
	// kafle czynne.
	DolaczWylaczone bool
}

// ZmianaKomponentu niesie pola `component.update`. Wskaźnik pusty znaczy „bez
// zmiany”, zgodnie ze zdaniem kontraktu „pola pominięte zostają bez zmian”.
type ZmianaKomponentu struct {
	Nazwa        *string
	Opis         *string
	Czynny       *bool
	Konfiguracja *string
}

// RepozytoriumKomponentow jest kontraktem rejestru komponentów.
type RepozytoriumKomponentow interface {
	ZalozKomponent(ctx context.Context, komponent Komponent) (Komponent, error)
	Komponent(ctx context.Context, kod string) (Komponent, error)
	Komponenty(ctx context.Context, filtr FiltrKomponentow) ([]Komponent, error)
	// ZmienKomponent i PrzypiszKomponent biorą czas zmiany od warstwy wyższej,
	// w milisekundach epoki — tak samo jak ZalozKomponent. Baza nie wstawia
	// własnego „teraz”, bo kolumna niesie wartość kontraktu bez przekładu;
	// dwa zegary dla jednego pola byłyby dwiema prawdami.
	ZmienKomponent(ctx context.Context, kod string, zmiana ZmianaKomponentu, teraz int64) (Komponent, error)
	UsunKomponent(ctx context.Context, kod string) (bool, error)
	// PrzypiszKomponent zapisuje parę (poziom zasięgu, klucz zasięgu). Drugi
	// wynik mówi, czy przypisanie coś zmieniło — powtórzenie tego samego
	// przypisania nie dochodzi do skutku i `component.assign` oddaje wtedy
	// `assigned: false` zamiast udawać czynność.
	PrzypiszKomponent(ctx context.Context, kod, poziom, kluczZasiegu string, teraz int64) (Komponent, bool, error)
}

const (
	kolumnyKomponentu = `k.id, k.identyfikator_zewnetrzny, k.rodzaj, k.nazwa, k.opis,
	                     k.byt_docelowy, k.czynny, k.konfiguracja, p.kod, k.klucz_zasiegu,
	                     k.utworzono, k.zaktualizowano`

	zrodloKomponentu = ` FROM komponent k LEFT JOIN poziom_zasiegu p ON p.id = k.poziom_zasiegu_id`

	wstawKomponent = `INSERT INTO komponent
	                  (identyfikator_zewnetrzny, rodzaj, nazwa, opis, byt_docelowy,
	                   czynny, konfiguracja, utworzono, zaktualizowano)
	                  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	pobierzKomponent = `SELECT ` + kolumnyKomponentu + zrodloKomponentu +
		` WHERE k.identyfikator_zewnetrzny = ?`

	// Jedno zapytanie na cztery warianty żądania: puste zawężenie rodzaju
	// wyłącza pierwszy warunek, a otwarty `includeDisabled` — drugi. Porządek
	// biegnie indeksem idx_komponent_wykaz (rodzaj, nazwa, id), więc „kolejność
	// wyświetlania” kontraktu jest stała między wywołaniami.
	listaKomponentow = `SELECT ` + kolumnyKomponentu + zrodloKomponentu +
		` WHERE (? = '' OR k.rodzaj = ?) AND (? = 1 OR k.czynny = 1)
		  ORDER BY k.rodzaj, k.nazwa, k.id`

	// Zmiana idzie jednym poleceniem: NULL w argumencie zostawia kolumnę bez
	// zmiany. Dzięki temu „pola pominięte zostają bez zmian” jest własnością
	// zapytania, a nie kolejnością gałęzi w Go.
	zmienKomponent = `UPDATE komponent SET
	                     nazwa          = COALESCE(?, nazwa),
	                     opis           = CASE WHEN ? = 1 THEN ? ELSE opis END,
	                     czynny         = COALESCE(?, czynny),
	                     konfiguracja   = COALESCE(?, konfiguracja),
	                     zaktualizowano = ?
	                  WHERE identyfikator_zewnetrzny = ?`

	przypiszKomponent = `UPDATE komponent SET
	                        poziom_zasiegu_id = (SELECT id FROM poziom_zasiegu WHERE kod = ?),
	                        klucz_zasiegu     = ?,
	                        zaktualizowano    = ?
	                     WHERE identyfikator_zewnetrzny = ?`

	usunKomponent = `DELETE FROM komponent WHERE identyfikator_zewnetrzny = ?`
)

// repozytoriumKomponentow nie trzyma `*sql.DB` obok pamięci zapytań, bo żaden
// zapis tego rejestru nie obejmuje drugiej tabeli.
type repozytoriumKomponentow struct {
	zapytania *zapytania
}

// noweRepozytoriumKomponentow zakłada rejestr nad pamięcią zapytań zestawu.
func noweRepozytoriumKomponentow(zapytania *zapytania) RepozytoriumKomponentow {
	if zapytania == nil {
		return nil
	}
	return &repozytoriumKomponentow{zapytania: zapytania}
}

// ZalozKomponent wstawia kafel Strefy 2. Czas utworzenia i zmiany podaje
// warstwa wyższa w milisekundach epoki — kolumna niesie wartość kontraktu bez
// przekładu, więc baza nie wstawia własnego „teraz”.
func (r *repozytoriumKomponentow) ZalozKomponent(ctx context.Context,
	komponent Komponent) (Komponent, error) {

	if komponent.Kod == "" {
		return Komponent{}, fmt.Errorf("dane: komponent bez identyfikatora")
	}
	if komponent.Nazwa == "" {
		return Komponent{}, fmt.Errorf("dane: komponent %q bez nazwy", komponent.Kod)
	}
	konfiguracja, err := konfiguracjaKomponentu(komponent.Konfiguracja, komponent.Kod)
	if err != nil {
		return Komponent{}, err
	}
	polecenie, err := r.zapytania.przygotuj(ctx, wstawKomponent)
	if err != nil {
		return Komponent{}, err
	}
	if _, err := polecenie.ExecContext(ctx, komponent.Kod, komponent.Rodzaj, komponent.Nazwa,
		tekstDoKolumny(komponent.Opis), tekstDoKolumny(komponent.BytDocelowy),
		liczbaLogiczna(komponent.Czynny), konfiguracja,
		komponent.Utworzono, komponent.Zaktualizowano); err != nil {

		return Komponent{}, fmt.Errorf("dane: nie można założyć komponentu %q: %w", komponent.Kod, err)
	}
	return r.Komponent(ctx, komponent.Kod)
}

// Komponent zwraca kafel o wskazanym identyfikatorze. Brak wiersza wraca jako
// ErrBrakWiersza — warstwa wyższa odróżnia „nie ma” od „odczyt padł”.
func (r *repozytoriumKomponentow) Komponent(ctx context.Context, kod string) (Komponent, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzKomponent)
	if err != nil {
		return Komponent{}, err
	}
	komponent, err := odczytajKomponent(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return Komponent{}, fmt.Errorf("dane: komponent %q nie istnieje: %w", kod, ErrBrakWiersza)
	}
	if err != nil {
		return Komponent{}, fmt.Errorf("dane: nie można odczytać komponentu %q: %w", kod, err)
	}
	return komponent, nil
}

// Komponenty zwraca wykaz w kolejności wyświetlania.
func (r *repozytoriumKomponentow) Komponenty(ctx context.Context,
	filtr FiltrKomponentow) ([]Komponent, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaKomponentow)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, filtr.Rodzaj, filtr.Rodzaj,
		liczbaLogiczna(filtr.DolaczWylaczone))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać wykazu komponentów: %w", err)
	}
	defer wiersze.Close()

	komponenty := make([]Komponent, 0, 16)
	for wiersze.Next() {
		komponent, err := odczytajKomponent(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: uszkodzony wiersz komponentu: %w", err)
		}
		komponenty = append(komponenty, komponent)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać wykazu komponentów: %w", err)
	}
	return komponenty, nil
}

// ZmienKomponent zmienia wyłącznie pola wskazane w żądaniu. Opis ma osobny
// przełącznik, bo jego zmianą może być wyczyszczenie do NULL — COALESCE sam
// nie odróżniłby „bez zmiany” od „wyczyść”.
func (r *repozytoriumKomponentow) ZmienKomponent(ctx context.Context, kod string,
	zmiana ZmianaKomponentu, teraz int64) (Komponent, error) {

	// Odczyt przed zapisem, żeby zmiana komponentu nieistniejącego wróciła jako
	// ErrBrakWiersza, a nie jako UPDATE bez skutku odmeldowany jako sukces.
	if _, err := r.Komponent(ctx, kod); err != nil {
		return Komponent{}, err
	}
	konfiguracja, err := zmienionaKonfiguracja(zmiana.Konfiguracja, kod)
	if err != nil {
		return Komponent{}, err
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zmienKomponent)
	if err != nil {
		return Komponent{}, err
	}
	var czynny any
	if zmiana.Czynny != nil {
		czynny = liczbaLogiczna(*zmiana.Czynny)
	}
	zmianaOpisu := liczbaLogiczna(zmiana.Opis != nil)
	if _, err := polecenie.ExecContext(ctx, tekstDoKolumny(zmiana.Nazwa),
		zmianaOpisu, tekstDoKolumny(zmiana.Opis), czynny, konfiguracja,
		teraz, kod); err != nil {

		return Komponent{}, fmt.Errorf("dane: nie można zmienić komponentu %q: %w", kod, err)
	}
	return r.Komponent(ctx, kod)
}

// UsunKomponent zdejmuje kafel Strefy 2. Bytu magazynu modułowego nie tyka —
// uzasadnienie stoi przy uchwycie `component.delete` (`handlers_komponenty.go`).
func (r *repozytoriumKomponentow) UsunKomponent(ctx context.Context, kod string) (bool, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, usunKomponent)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, kod)
	if err != nil {
		return false, fmt.Errorf("dane: nie można usunąć komponentu %q: %w", kod, err)
	}
	usuniete, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nie można potwierdzić usunięcia komponentu %q: %w", kod, err)
	}
	return usuniete > 0, nil
}

// PrzypiszKomponent zapisuje parę (poziom zasięgu, klucz zasięgu). Powtórzenie
// tego samego przypisania niczego nie zmienia i wraca jako `false` — to jedyny
// przypadek, w którym `component.assign` może uczciwie oddać `assigned: false`
// bez odmowy.
func (r *repozytoriumKomponentow) PrzypiszKomponent(ctx context.Context,
	kod, poziom, kluczZasiegu string, teraz int64) (Komponent, bool, error) {

	zastany, err := r.Komponent(ctx, kod)
	if err != nil {
		return Komponent{}, false, err
	}
	if zastany.PoziomZasiegu != nil && *zastany.PoziomZasiegu == poziom &&
		zastany.KluczZasiegu == kluczZasiegu {

		return zastany, false, nil
	}
	polecenie, err := r.zapytania.przygotuj(ctx, przypiszKomponent)
	if err != nil {
		return Komponent{}, false, err
	}
	if _, err := polecenie.ExecContext(ctx, poziom, kluczZasiegu,
		teraz, kod); err != nil {

		return Komponent{}, false, fmt.Errorf("dane: nie można przypisać komponentu %q: %w", kod, err)
	}
	przypisany, err := r.Komponent(ctx, kod)
	if err != nil {
		return Komponent{}, false, err
	}
	// Poziom spoza słownika `poziom_zasiegu` dałby podzapytanie puste, więc
	// kolumna zostałaby NULL, a komenda odmeldowałaby sukces bez skutku.
	// Warstwa wyższa sprawdza wartość kontraktu, ta sprawdza skutek.
	if przypisany.PoziomZasiegu == nil {
		return Komponent{}, false,
			fmt.Errorf("dane: poziom zasięgu %q nie istnieje w słowniku", poziom)
	}
	return przypisany, true, nil
}

// konfiguracjaKomponentu sprawdza, że konfiguracja jest poprawnym JSON-em, i
// zamienia brak na pusty obiekt — kolumna jest NOT NULL DEFAULT '{}'.
func konfiguracjaKomponentu(tresc, kod string) (string, error) {
	if tresc == "" {
		return "{}", nil
	}
	if !json.Valid([]byte(tresc)) {
		return "", fmt.Errorf("dane: konfiguracja komponentu %q nie jest poprawnym JSON-em", kod)
	}
	return tresc, nil
}

// zmienionaKonfiguracja przekłada wskaźnik zmiany na argument zapytania: nil
// znaczy „bez zmiany” i zostawia kolumnę nietkniętą przez COALESCE.
func zmienionaKonfiguracja(tresc *string, kod string) (any, error) {
	if tresc == nil {
		return nil, nil
	}
	sprawdzona, err := konfiguracjaKomponentu(*tresc, kod)
	if err != nil {
		return nil, err
	}
	return sprawdzona, nil
}

func odczytajKomponent(wiersz skaner) (Komponent, error) {
	var komponent Komponent
	var opis, bytDocelowy, poziom sql.NullString
	var czynny int
	if err := wiersz.Scan(&komponent.ID, &komponent.Kod, &komponent.Rodzaj, &komponent.Nazwa,
		&opis, &bytDocelowy, &czynny, &komponent.Konfiguracja, &poziom, &komponent.KluczZasiegu,
		&komponent.Utworzono, &komponent.Zaktualizowano); err != nil {

		return Komponent{}, err
	}
	komponent.Opis = tekstZKolumny(opis)
	komponent.BytDocelowy = tekstZKolumny(bytDocelowy)
	komponent.PoziomZasiegu = tekstZKolumny(poziom)
	komponent.Czynny = czynny == 1
	return komponent, nil
}
