// Odpowiedzialność pliku: notatki powiązane ze źródłem przeglądania, wiersze tabeli notatka_przegladania, obsługujące komendę browser.note.add.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// NotatkaPrzegladania to wiersz tabeli `notatka_przegladania` — notatka
// Operatora zapisana w toku przeglądania, opcjonalnie powiązana ze źródłem
// (`ZrodloID`) i cytatem fragmentu strony (`Cytat`).
type NotatkaPrzegladania struct {
	ID       int64
	Kod      string
	Okno     string
	ZrodloID *string
	Tresc    string
	Cytat    *string
	// Trzy pola dołożone komendą update: rodzaj notatki, wątek tematyczny i przypięcie na początek wykazu.
	Klasyfikacja   *string
	Watek          *string
	Przypieta      bool
	Utworzono      string
	Zaktualizowano string
}

const (
	kolumnyNotatki = `id, identyfikator_zewnetrzny, okno, zrodlo_zewnetrzny_id,
	                  tresc, cytat, klasyfikacja, watek, przypieta,
	                  utworzono, zaktualizowano`

	zapiszNotatkePrzegladania = `INSERT INTO notatka_przegladania
	                 (identyfikator_zewnetrzny, okno, zrodlo_zewnetrzny_id, tresc, cytat,
	                  klasyfikacja, watek, przypieta, konto_id)
	                 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ` + WskazanieKonta + `)`

	pobierzNotatkePrzegladania = `SELECT ` + kolumnyNotatki + ` FROM notatka_przegladania
	                  WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	// Kolejność zgodna z indeksem idx_notatka_przegladania_okno (okno, utworzono DESC, id).
	// Puste zawężenie do źródła wyłącza warunek — jedno zapytanie na oba
	// warianty `browser.note.list`, bez sklejania SQL-a w locie.
	listaNotatekOkna = `SELECT ` + kolumnyNotatki + ` FROM notatka_przegladania
	                    WHERE okno = ? AND ` + WarunekKonta + `
	                      AND (? = '' OR zrodlo_zewnetrzny_id = ?)
	                      AND (? = '' OR watek = ?)
	                      AND (? = '' OR klasyfikacja = ?)
	                      AND (? = '' OR tresc LIKE ? OR IFNULL(cytat,'') LIKE ?)
	                    ORDER BY (? = 1 AND przypieta = 1) DESC, utworzono DESC, id DESC LIMIT ?`

	// Aktualizacja notatki zmienia wyłącznie pola, które żądanie naprawdę przyniosło; wzorzec COALESCE nie kasuje pól przy zmianie tylko jednego z nich.
	aktualizujNotatkePrzegladania = `UPDATE notatka_przegladania SET
	                                     tresc = COALESCE(?, tresc),
	                                     cytat = COALESCE(?, cytat),
	                                     zrodlo_zewnetrzny_id = COALESCE(?, zrodlo_zewnetrzny_id),
	                                     klasyfikacja = COALESCE(?, klasyfikacja),
	                                     watek = COALESCE(?, watek),
	                                     przypieta = COALESCE(?, przypieta),
	                                     zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                                 WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta
)

// FiltrNotatekPrzegladania zawęża wykaz notatek okna, obsługując pola żądania browser.note.list w całości.
type FiltrNotatekPrzegladania struct {
	// Okno operacyjne, którego wykaz dotyczy — pole wymagane kontraktem.
	Okno string
	// ZrodloID zawęża do notatek jednego źródła; pusty napis znaczy wszystkie notatki okna.
	ZrodloID string
	// Watek zawęża do jednego wątku tematycznego notatek.
	Watek string
	// Klasyfikacja zawęża do jednego rodzaju notatki.
	Klasyfikacja string
	// Szukaj przegląda treść i cytat naraz.
	Szukaj string
	// PrzypieteNaPoczatku wynosi notatki przypięte na czoło wykazu, zamiast
	// zostawiać je w porządku czasu.
	PrzypieteNaPoczatku bool
	// Limit 0 lub ujemny znaczy wykaz pełny, nie wykaz pusty.
	Limit int
}

// ZapiszNotatke wstawia nową notatkę i zwraca ją odczytaną z bazy, z nadanym
// identyfikatorem i czasami utworzenia/aktualizacji.
func (r *repozytoriumPrzegladania) ZapiszNotatke(ctx context.Context,
	notatka NotatkaPrzegladania) (NotatkaPrzegladania, error) {

	if notatka.Kod == "" {
		return NotatkaPrzegladania{}, fmt.Errorf("dane: notatka przeglądania bez identyfikatora")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszNotatkePrzegladania)
	if err != nil {
		return NotatkaPrzegladania{}, err
	}
	_, err = polecenie.ExecContext(ctx, notatka.Kod, notatka.Okno,
		tekstDoKolumny(notatka.ZrodloID), notatka.Tresc, tekstDoKolumny(notatka.Cytat),
		tekstDoKolumny(notatka.Klasyfikacja), tekstDoKolumny(notatka.Watek),
		liczbaLogiczna(notatka.Przypieta), KontoOperatora(ctx))
	if err != nil {
		return NotatkaPrzegladania{}, fmt.Errorf("dane: nie można zapisać notatki %q: %w", notatka.Kod, err)
	}
	return r.pobierzNotatke(ctx, notatka.Kod)
}

// pobierzNotatke zwraca notatkę o wskazanym kodzie. Brak wiersza wraca jako
// ErrBrakWiersza — warstwa wyższa odróżnia „nie ma" od „odczyt się nie powiódł".
func (r *repozytoriumPrzegladania) pobierzNotatke(ctx context.Context, kod string) (NotatkaPrzegladania, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzNotatkePrzegladania)
	if err != nil {
		return NotatkaPrzegladania{}, err
	}
	notatka, err := odczytajNotatke(polecenie.QueryRowContext(ctx, kod, KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return NotatkaPrzegladania{}, ErrBrakWiersza
	}
	if err != nil {
		return NotatkaPrzegladania{}, fmt.Errorf("dane: nieczytelny wiersz notatki %q: %w", kod, err)
	}
	return notatka, nil
}

// Notatki zwraca notatki okna od najświeższej, w całości albo zawężone do
// jednego źródła. Brak notatek to wykaz pusty, nie brak wiersza — okno bez
// notatek nie jest oknem nieznanym (rozróżnienie robi `OknoZnane`).
func (r *repozytoriumPrzegladania) Notatki(ctx context.Context,
	filtr FiltrNotatekPrzegladania) ([]NotatkaPrzegladania, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaNotatekOkna)
	if err != nil {
		return nil, err
	}
	okno := filtr.Okno
	wzorzec := ""
	if filtr.Szukaj != "" {
		wzorzec = "%" + filtr.Szukaj + "%"
	}
	wiersze, err := polecenie.QueryContext(ctx, okno, KontoOperatora(ctx), filtr.ZrodloID, filtr.ZrodloID,
		filtr.Watek, filtr.Watek, filtr.Klasyfikacja, filtr.Klasyfikacja,
		filtr.Szukaj, wzorzec, wzorzec,
		liczbaLogiczna(filtr.PrzypieteNaPoczatku), granicaWykazu(filtr.Limit))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać notatek okna %q: %w", okno, err)
	}
	defer wiersze.Close()

	lista := []NotatkaPrzegladania{}
	for wiersze.Next() {
		notatka, err := odczytajNotatke(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz notatek okna %q: %w", okno, err)
		}
		lista = append(lista, notatka)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt notatek okna %q: %w", okno, err)
	}
	return lista, nil
}

// odczytajNotatke składa strukturę notatki przeglądania z jednego wiersza wyniku zapytania, kolumna po kolumnie.
func odczytajNotatke(wiersz skaner) (NotatkaPrzegladania, error) {
	var notatka NotatkaPrzegladania
	var zrodloID, cytat, klasyfikacja, watek sql.NullString
	var przypieta int64
	err := wiersz.Scan(&notatka.ID, &notatka.Kod, &notatka.Okno, &zrodloID,
		&notatka.Tresc, &cytat, &klasyfikacja, &watek, &przypieta,
		&notatka.Utworzono, &notatka.Zaktualizowano)
	if err != nil {
		return NotatkaPrzegladania{}, err
	}
	notatka.ZrodloID = tekstZKolumny(zrodloID)
	notatka.Cytat = tekstZKolumny(cytat)
	notatka.Klasyfikacja = tekstZKolumny(klasyfikacja)
	notatka.Watek = tekstZKolumny(watek)
	notatka.Przypieta = przypieta == 1
	return notatka, nil
}

// ZmianaNotatki niesie pola, które `browser.note.update` naprawdę przyniosła.
// Wskaźnik pusty znaczy „nie ruszaj tej kolumny" — dlatego typ jest osobny od
// `NotatkaPrzegladania`, w której puste pole znaczy „notatka tego nie ma".
type ZmianaNotatki struct {
	Kod          string
	Tresc        *string
	Cytat        *string
	ZrodloID     *string
	Klasyfikacja *string
	Watek        *string
	Przypieta    *bool
}

// AktualizujNotatke zmienia notatkę zapisaną wcześniej i oddaje jej stan po
// zmianie. Brak wiersza wraca jako ErrBrakWiersza, żeby komenda mogła odmówić
// zamiast meldować zmianę czegoś, czego nie ma.
func (r *repozytoriumPrzegladania) AktualizujNotatke(ctx context.Context,
	zmiana ZmianaNotatki) (NotatkaPrzegladania, error) {

	if zmiana.Kod == "" {
		return NotatkaPrzegladania{}, fmt.Errorf("dane: aktualizacja notatki bez identyfikatora")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, aktualizujNotatkePrzegladania)
	if err != nil {
		return NotatkaPrzegladania{}, err
	}
	var przypieta any
	if zmiana.Przypieta != nil {
		przypieta = liczbaLogiczna(*zmiana.Przypieta)
	}
	wynik, err := polecenie.ExecContext(ctx, tekstDoKolumny(zmiana.Tresc),
		tekstDoKolumny(zmiana.Cytat), tekstDoKolumny(zmiana.ZrodloID),
		tekstDoKolumny(zmiana.Klasyfikacja), tekstDoKolumny(zmiana.Watek),
		przypieta, zmiana.Kod, KontoOperatora(ctx))
	if err != nil {
		return NotatkaPrzegladania{}, fmt.Errorf("dane: nie można zmienić notatki %q: %w", zmiana.Kod, err)
	}
	zmienione, err := wynik.RowsAffected()
	if err != nil {
		return NotatkaPrzegladania{}, fmt.Errorf("dane: nie można ustalić skutku zmiany notatki %q: %w", zmiana.Kod, err)
	}
	if zmienione == 0 {
		return NotatkaPrzegladania{}, ErrBrakWiersza
	}
	return r.pobierzNotatke(ctx, zmiana.Kod)
}

// Notatka oddaje pojedynczą notatkę po kodzie zewnętrznym — potrzebna wykazom
// wątków i komendzie zmiany.
func (r *repozytoriumPrzegladania) Notatka(ctx context.Context, kod string) (NotatkaPrzegladania, error) {
	return r.pobierzNotatke(ctx, kod)
}
