// Zlecenia kolejki widziane przez Queue Managera (tabela `zlecenie_kolejki`)
// oraz polityka kolejki (kolumny tabeli `kolejka`). To nie jest drugi silnik
// kolejek ani druga tabela pozycji obok `pozycja_kolejki`.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// Zlecenie to wiersz tabeli `zlecenie_kolejki`. Stan jest słownikiem bazy
// (kolumna `baza` wyliczenia `QueueItemStatus`).
type Zlecenie struct {
	ID                 int64
	Kod                string
	KolejkaID          int64
	Stan               string
	Ladunek            *string
	Priorytet          int
	Proby              int
	Warunek            *string
	KluczIdempotencji  *string
	PrzebiegID         *int64
	PrzebiegKod        *string
	EkspertDocelowy    *string
	ZlecenieZrodloweID *int64
	Termin             *string
	Utworzono          string
	Zaktualizowano     string
}

// PolitykaKolejki to komplet kolumn polityki tabeli `kolejka`: limity,
// wycofanie, rozproszenie i zadania martwe.
type PolitykaKolejki struct {
	Zasieg            string
	ZasiegID          *string
	LimitRownoleglych int
	TempoNaMinute     int
	LimitProb         int
	Wycofanie         string
	WycofanieSekundy  int
	Rozproszenie      bool
	ZadaniaMartwe     bool
	IdempotencjaZycie int
}

// PunktGlebokosci to jeden odcinek wykresu głębokości kolejki: chwila,
// liczba oczekujących i pracujących.
type PunktGlebokosci struct {
	Chwila       string
	Oczekujacych int
	Pracujacych  int
}

const (
	kolumnyZlecenia = `z.id, z.identyfikator_zewnetrzny, z.kolejka_id, z.stan, z.ladunek,
	                   z.priorytet, z.proby, z.warunek, z.klucz_idempotencji, z.przebieg_id,
	                   p.identyfikator_zewnetrzny, z.ekspert_docelowy, z.zlecenie_zrodlowe_id,
	                   z.termin, z.utworzono, z.zaktualizowano`

	zrodloZlecenia = ` FROM zlecenie_kolejki z
	                   LEFT JOIN przebieg_automatyki p ON p.id = z.przebieg_id`

	wstawZlecenieKolejki = `INSERT INTO zlecenie_kolejki
	                        (identyfikator_zewnetrzny, kolejka_id, stan, ladunek, priorytet,
	                         warunek, klucz_idempotencji, przebieg_id, ekspert_docelowy,
	                         zlecenie_zrodlowe_id, termin)
	                        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	pobierzZlecenieKolejki = `SELECT ` + kolumnyZlecenia + zrodloZlecenia +
		` WHERE z.identyfikator_zewnetrzny = ?`

	pobierzZlecenieKluczem = `SELECT ` + kolumnyZlecenia + zrodloZlecenia +
		` WHERE z.kolejka_id = ? AND z.klucz_idempotencji = ?`

	// Stan pusty znaczy „wszystkie stany”. Zadania martwe i zdjęte wychodzą
	// wtedy razem z resztą — Queue Manager pokazuje je w kolumnie stanu,
	// a nie ukrywa przed Operatorem.
	listaZlecenKolejki = `SELECT ` + kolumnyZlecenia + zrodloZlecenia +
		` WHERE z.kolejka_id = ? AND (? = '' OR z.stan = ?)
		  ORDER BY z.priorytet, z.id LIMIT ?`

	liczbaZlecenKolejki = `SELECT COUNT(*) FROM zlecenie_kolejki WHERE kolejka_id = ?`

	// Zero w miejscu kolejki znaczy zadania martwe wszystkich kolejek, nie
	// kolejkę o identyfikatorze zero.
	listaZlecenMartwych = `SELECT ` + kolumnyZlecenia + zrodloZlecenia +
		` WHERE z.stan = 'martwe' AND (? = 0 OR z.kolejka_id = ?)
		  ORDER BY z.zaktualizowano DESC, z.id DESC LIMIT ?`

	zmienStanZlecenia = `UPDATE zlecenie_kolejki
	                     SET stan = ?, zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                     WHERE id = ?`

	zmienTerminZlecenia = `UPDATE zlecenie_kolejki
	                       SET stan = 'odlozone', termin = ?,
	                           zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                       WHERE id = ?`

	zmienWarunekZlecenia = `UPDATE zlecenie_kolejki
	                        SET warunek = ?, zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                        WHERE id = ?`

	skierujZlecenie = `UPDATE zlecenie_kolejki
	                   SET kolejka_id = ?, ekspert_docelowy = ?,
	                       zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                   WHERE id = ?`

	// Głębokość liczy się z bazy, nie z pętli w Go: odcinek wyznacza podłoga
	// z liczby sekund epoki podzielonej przez długość odcinka, a zlecenie
	// liczy się do odcinka swojego założenia.
	glebokoscKolejki = `SELECT CAST(strftime('%s', utworzono) AS INTEGER) / ? AS odcinek,
	                           SUM(CASE WHEN stan IN ('oczekuje','odlozone') THEN 1 ELSE 0 END),
	                           SUM(CASE WHEN stan = 'przetwarzane' THEN 1 ELSE 0 END)
	                    FROM zlecenie_kolejki
	                    WHERE (? = 0 OR kolejka_id = ?) AND utworzono >= ?
	                    GROUP BY odcinek ORDER BY odcinek`

	kolumnyPolitykiKolejki = `zasieg, zasieg_id, limit_rownoleglych, tempo_na_minute, limit_prob,
	                          wycofanie, wycofanie_sekundy, rozproszenie, zadania_martwe,
	                          idempotencja_zycie_sekundy`

	pobierzPolitykeKolejki = `SELECT ` + kolumnyPolitykiKolejki + ` FROM kolejka
	                          WHERE id = ? AND ` + WarunekKonta

	zapiszPolitykeKolejki = `UPDATE kolejka
	                         SET zasieg = ?, zasieg_id = ?, limit_rownoleglych = ?,
	                             tempo_na_minute = ?, limit_prob = ?, wycofanie = ?,
	                             wycofanie_sekundy = ?, rozproszenie = ?, zadania_martwe = ?,
	                             idempotencja_zycie_sekundy = ?,
	                             zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                         WHERE id = ? AND ` + WarunekKonta
)

// DodajZlecenie zakłada nowe zlecenie w kolejce, nadając mu stan oczekuje,
// gdy stan nie został wskazany.
func (r *repozytoriumKolejek) DodajZlecenie(ctx context.Context, zlecenie Zlecenie) (Zlecenie, error) {
	if zlecenie.Kod == "" {
		return Zlecenie{}, fmt.Errorf("dane: zlecenie kolejki bez identyfikatora")
	}
	stan := zlecenie.Stan
	if stan == "" {
		stan = "oczekuje"
	}
	polecenie, err := r.zapytania.przygotuj(ctx, wstawZlecenieKolejki)
	if err != nil {
		return Zlecenie{}, err
	}
	_, err = polecenie.ExecContext(ctx, zlecenie.Kod, zlecenie.KolejkaID, stan,
		tekstDoKolumny(zlecenie.Ladunek), zlecenie.Priorytet, tekstDoKolumny(zlecenie.Warunek),
		tekstDoKolumny(zlecenie.KluczIdempotencji), liczbaDoKolumny(zlecenie.PrzebiegID),
		tekstDoKolumny(zlecenie.EkspertDocelowy), liczbaDoKolumny(zlecenie.ZlecenieZrodloweID),
		tekstDoKolumny(zlecenie.Termin))
	if err != nil {
		return Zlecenie{}, fmt.Errorf("dane: nie można dodać zlecenia do kolejki %d: %w",
			zlecenie.KolejkaID, err)
	}
	return r.Zlecenie(ctx, zlecenie.Kod)
}

// Zlecenie zwraca jedno zlecenie kolejki po kodzie zewnętrznym. Brak wiersza
// wraca jako ErrBrakWiersza.
func (r *repozytoriumKolejek) Zlecenie(ctx context.Context, kod string) (Zlecenie, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzZlecenieKolejki)
	if err != nil {
		return Zlecenie{}, err
	}
	zlecenie, err := odczytajZlecenieKolejki(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return Zlecenie{}, ErrBrakWiersza
	}
	if err != nil {
		return Zlecenie{}, fmt.Errorf("dane: nieczytelne zlecenie kolejki %q: %w", kod, err)
	}
	return zlecenie, nil
}

// ZlecenieKluczem odnajduje zlecenie po kluczu idempotencji w obrębie kolejki.
// Ten sam klucz w dwóch kolejkach opisuje dwa różne zlecenia dwóch różnych
// torów, więc odczyt zawsze pyta o parę.
func (r *repozytoriumKolejek) ZlecenieKluczem(ctx context.Context, kolejkaID int64,
	klucz string) (Zlecenie, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzZlecenieKluczem)
	if err != nil {
		return Zlecenie{}, err
	}
	zlecenie, err := odczytajZlecenieKolejki(polecenie.QueryRowContext(ctx, kolejkaID, klucz))
	if errors.Is(err, sql.ErrNoRows) {
		return Zlecenie{}, ErrBrakWiersza
	}
	if err != nil {
		return Zlecenie{}, fmt.Errorf("dane: nieczytelne zlecenie o kluczu %q: %w", klucz, err)
	}
	return zlecenie, nil
}

// ZleceniaKolejki zwraca zlecenia w porządku przetwarzania wraz z liczbą
// wszystkich zleceń kolejki — wykaz bywa przycięty granicą, a licznik nie.
func (r *repozytoriumKolejek) ZleceniaKolejki(ctx context.Context, kolejkaID int64,
	stan string, limit int) ([]Zlecenie, int, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaZlecenKolejki)
	if err != nil {
		return nil, 0, err
	}
	wiersze, err := polecenie.QueryContext(ctx, kolejkaID, stan, stan, granicaWykazu(limit))
	if err != nil {
		return nil, 0, fmt.Errorf("dane: nie można odczytać zleceń kolejki %d: %w", kolejkaID, err)
	}
	lista, err := zbierzZleceniaKolejki(wiersze)
	if err != nil {
		return nil, 0, err
	}
	licznik, err := r.zapytania.przygotuj(ctx, liczbaZlecenKolejki)
	if err != nil {
		return nil, 0, err
	}
	var wszystkich int
	if err := licznik.QueryRowContext(ctx, kolejkaID).Scan(&wszystkich); err != nil {
		return nil, 0, fmt.Errorf("dane: nie można policzyć zleceń kolejki %d: %w", kolejkaID, err)
	}
	return lista, wszystkich, nil
}

// ZleceniaMartwe zwraca zlecenia trwale nieudane, jednej wskazanej kolejki
// albo wszystkich kolejek naraz.
func (r *repozytoriumKolejek) ZleceniaMartwe(ctx context.Context, kolejkaID int64,
	limit int) ([]Zlecenie, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, listaZlecenMartwych)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, kolejkaID, kolejkaID, granicaWykazu(limit))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać zadań martwych: %w", err)
	}
	return zbierzZleceniaKolejki(wiersze)
}

// ZmienStanZlecenia zapisuje nowy stan wskazanego zlecenia i odświeża
// znacznik ostatniej aktualizacji.
func (r *repozytoriumKolejek) ZmienStanZlecenia(ctx context.Context, zlecenieID int64,
	stan string) error {

	return wykonajZapisZlecenia(ctx, r.zapytania, zmienStanZlecenia, stan, zlecenieID)
}

// OdlozZlecenie przenosi wskazane zlecenie w stan odłożony wraz z terminem
// wykonania, do którego czeka.
func (r *repozytoriumKolejek) OdlozZlecenie(ctx context.Context, zlecenieID int64,
	termin string) error {

	return wykonajZapisZlecenia(ctx, r.zapytania, zmienTerminZlecenia, termin, zlecenieID)
}

// UstawWarunekZlecenia zapisuje warunek przetworzenia zlecenia; pusty
// warunek zdejmuje go całkowicie z wiersza.
func (r *repozytoriumKolejek) UstawWarunekZlecenia(ctx context.Context, zlecenieID int64,
	warunek *string) error {

	return wykonajZapisZlecenia(ctx, r.zapytania, zmienWarunekZlecenia,
		tekstDoKolumny(warunek), zlecenieID)
}

// SkierujZlecenie zmienia kolejkę zlecenia albo jego eksperta docelowego,
// nie ruszając reszty wiersza.
func (r *repozytoriumKolejek) SkierujZlecenie(ctx context.Context, zlecenieID, kolejkaID int64,
	ekspert *string) error {

	polecenie, err := r.zapytania.przygotuj(ctx, skierujZlecenie)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, kolejkaID, tekstDoKolumny(ekspert), zlecenieID); err != nil {
		return fmt.Errorf("dane: nie można skierować zlecenia %d: %w", zlecenieID, err)
	}
	return nil
}

// GlebokoscKolejki liczy zlecenia oczekujące i przetwarzane w kolejnych
// odcinkach czasu zadanej długości.
func (r *repozytoriumKolejek) GlebokoscKolejki(ctx context.Context, kolejkaID int64,
	odcinekSekundy int, od string) ([]PunktGlebokosci, error) {

	if odcinekSekundy <= 0 {
		odcinekSekundy = 3600
	}
	polecenie, err := r.zapytania.przygotuj(ctx, glebokoscKolejki)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, odcinekSekundy, kolejkaID, kolejkaID, od)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można policzyć głębokości kolejki: %w", err)
	}
	defer wiersze.Close()

	lista := []PunktGlebokosci{}
	for wiersze.Next() {
		var odcinek int64
		var punkt PunktGlebokosci
		if err := wiersze.Scan(&odcinek, &punkt.Oczekujacych, &punkt.Pracujacych); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny odcinek głębokości kolejki: %w", err)
		}
		// Odcinek wraca jako numer, więc chwilę odtwarza się mnożeniem —
		// początek odcinka w sekundach epoki.
		punkt.Chwila = fmt.Sprintf("%d", odcinek*int64(odcinekSekundy))
		lista = append(lista, punkt)
	}
	return lista, wiersze.Err()
}

// PolitykaKolejki zwraca politykę kolejki. Kolejka bez zapisanej polityki nie
// jest kolejką bez polityki: kolumny mają wartości domyślne modelu konfiguracji.
func (r *repozytoriumKolejek) PolitykaKolejki(ctx context.Context,
	kolejkaID int64) (PolitykaKolejki, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzPolitykeKolejki)
	if err != nil {
		return PolitykaKolejki{}, err
	}
	var polityka PolitykaKolejki
	var zasiegID sql.NullString
	var rozproszenie, martwe int
	err = polecenie.QueryRowContext(ctx, kolejkaID, KontoOperatora(ctx)).Scan(&polityka.Zasieg, &zasiegID,
		&polityka.LimitRownoleglych, &polityka.TempoNaMinute, &polityka.LimitProb,
		&polityka.Wycofanie, &polityka.WycofanieSekundy, &rozproszenie, &martwe,
		&polityka.IdempotencjaZycie)
	if errors.Is(err, sql.ErrNoRows) {
		return PolitykaKolejki{}, ErrBrakWiersza
	}
	if err != nil {
		return PolitykaKolejki{}, fmt.Errorf("dane: nieczytelna polityka kolejki %d: %w", kolejkaID, err)
	}
	polityka.ZasiegID = tekstZKolumny(zasiegID)
	polityka.Rozproszenie = rozproszenie == 1
	polityka.ZadaniaMartwe = martwe == 1
	return polityka, nil
}

// ZapiszPolitykeKolejki zapisuje politykę kolejki: limity, zasady
// wycofania, rozproszenie i zadania martwe.
func (r *repozytoriumKolejek) ZapiszPolitykeKolejki(ctx context.Context, kolejkaID int64,
	polityka PolitykaKolejki) error {

	polecenie, err := r.zapytania.przygotuj(ctx, zapiszPolitykeKolejki)
	if err != nil {
		return err
	}
	_, err = polecenie.ExecContext(ctx, polityka.Zasieg, tekstDoKolumny(polityka.ZasiegID),
		polityka.LimitRownoleglych, polityka.TempoNaMinute, polityka.LimitProb,
		polityka.Wycofanie, polityka.WycofanieSekundy, liczbaLogiczna(polityka.Rozproszenie),
		liczbaLogiczna(polityka.ZadaniaMartwe), polityka.IdempotencjaZycie, kolejkaID,
		KontoOperatora(ctx))
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać polityki kolejki %d: %w", kolejkaID, err)
	}
	return nil
}

// wykonajZapisZlecenia wykonuje jedno polecenie zmieniające zlecenie,
// wspólne dla kilku funkcji zapisu.
func wykonajZapisZlecenia(ctx context.Context, z *zapytania, tekst string,
	wartosc any, zlecenieID int64) error {

	polecenie, err := z.przygotuj(ctx, tekst)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, wartosc, zlecenieID); err != nil {
		return fmt.Errorf("dane: nie można zmienić zlecenia %d: %w", zlecenieID, err)
	}
	return nil
}

// zbierzZleceniaKolejki składa wykaz zleceń kolejki z wyniku zapytania,
// zamykając wiersze po odczycie.
func zbierzZleceniaKolejki(wiersze *sql.Rows) ([]Zlecenie, error) {
	defer wiersze.Close()

	lista := []Zlecenie{}
	for wiersze.Next() {
		zlecenie, err := odczytajZlecenieKolejki(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz zlecenia kolejki: %w", err)
		}
		lista = append(lista, zlecenie)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt zleceń kolejki: %w", err)
	}
	return lista, nil
}

// odczytajZlecenieKolejki składa strukturę Zlecenie z jednego wiersza
// wyniku zapytania SQL bazy danych.
func odczytajZlecenieKolejki(wiersz skaner) (Zlecenie, error) {
	var zlecenie Zlecenie
	var ladunek, warunek, klucz, kodPrzebiegu, ekspert, termin sql.NullString
	var przebieg, zrodlowe sql.NullInt64
	err := wiersz.Scan(&zlecenie.ID, &zlecenie.Kod, &zlecenie.KolejkaID, &zlecenie.Stan,
		&ladunek, &zlecenie.Priorytet, &zlecenie.Proby, &warunek, &klucz, &przebieg,
		&kodPrzebiegu, &ekspert, &zrodlowe, &termin, &zlecenie.Utworzono, &zlecenie.Zaktualizowano)
	if err != nil {
		return Zlecenie{}, err
	}
	zlecenie.Ladunek = tekstZKolumny(ladunek)
	zlecenie.Warunek = tekstZKolumny(warunek)
	zlecenie.KluczIdempotencji = tekstZKolumny(klucz)
	zlecenie.PrzebiegID = liczbaZKolumny(przebieg)
	zlecenie.PrzebiegKod = tekstZKolumny(kodPrzebiegu)
	zlecenie.EkspertDocelowy = tekstZKolumny(ekspert)
	zlecenie.ZlecenieZrodloweID = liczbaZKolumny(zrodlowe)
	zlecenie.Termin = tekstZKolumny(termin)
	return zlecenie, nil
}
