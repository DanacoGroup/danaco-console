// Odpowiedzialność pliku: dostęp do rejestru centrum powiadomień (tabela
// `powiadomienie_centrum` z migracji 379).
//
// Rejestr jest trwały, nie ulotny: zdarzenie zapisane tu przeżywa zamknięcie
// okna i restart rdzenia, bo centrum jest — jak mówi karta komponentu
// (rozdz. 11.6) — „trwałym rejestrem tych samych zdarzeń", których ulotną
// postacią jest Toast.
//
// Repozytorium nie zna kontraktu ani zdarzeń rozgłaszanych: czyta i zapisuje
// wiersze. Rozgłoszenie i przekład na kształt kontraktu należą do rdzenia.
package dane

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ZdarzenieCentrum to wiersz rejestru centrum powiadomień.
type ZdarzenieCentrum struct {
	ID            int64
	Klasa         string
	Waga          string
	Tresc         string
	ZrodloTyp     string
	ZrodloID      string
	SrodowiskoKod string
	SesjaID       string
	Akcje         string
	Stan          string
	Kanal         string
	OdlozoneDo    time.Time
	ZnacznikCzasu time.Time
}

// FiltrCentrum zawęża odczyt rejestru. Wykaz pusty znaczy „bez zawężenia".
type FiltrCentrum struct {
	Klasy         []string
	Wagi          []string
	Stany         []string
	SrodowiskoKod string
	SesjaID       string
	Limit         int
}

// RepozytoriumCentrumPowiadomien jest kontraktem rejestru centrum.
type RepozytoriumCentrumPowiadomien interface {
	// Zapisz dopisuje zdarzenie i oddaje je z nadanym identyfikatorem.
	Zapisz(ctx context.Context, zdarzenie ZdarzenieCentrum) (ZdarzenieCentrum, error)
	// Wykaz oddaje zdarzenia od najnowszego, zawężone filtrem.
	Wykaz(ctx context.Context, filtr FiltrCentrum) ([]ZdarzenieCentrum, error)
	// Nowe liczy zdarzenia w stanie `nowe` w CAŁYM rejestrze, nie w widoku —
	// plakietka paska kontekstu liczy rejestr, a nie przefiltrowaną kolumnę.
	Nowe(ctx context.Context) (int, error)
	// Odczytaj przenosi wskazane zdarzenia do stanu `odczytane`; wykaz pusty
	// bierze wszystkie zdarzenia nowe (działanie zbiorcze centrum).
	Odczytaj(ctx context.Context, identyfikatory []int64) (int, error)
	// Zamknij przenosi jedno zdarzenie do stanu `obsluzone`.
	Zamknij(ctx context.Context, id int64) (bool, error)
	// Odloz przenosi jedno zdarzenie do stanu `odlozone` wraz z chwilą powrotu.
	Odloz(ctx context.Context, id int64, doKiedy time.Time) (bool, error)
	// Przywroc oddaje do stanu `nowe` zdarzenia, których chwila powrotu minęła.
	// Bez tego odłożenie byłoby cichym skasowaniem.
	Przywroc(ctx context.Context, teraz time.Time) (int, error)
}

const (
	kolumnyCentrum = `id, klasa, waga, tresc, COALESCE(zrodlo_typ,''), COALESCE(zrodlo_id,''),
	                  srodowisko_kod, sesja_id, akcje, stan, kanal_dostarczenia,
	                  COALESCE(odlozone_do,''), znacznik_czasu`

	zapiszZdarzenieCentrum = `INSERT INTO powiadomienie_centrum
	    (klasa, waga, tresc, zrodlo_typ, zrodlo_id, srodowisko_kod, sesja_id, akcje,
	     stan, kanal_dostarczenia, znacznik_czasu)
	    VALUES (?, ?, ?, NULLIF(?,''), NULLIF(?,''), ?, ?, ?, 'nowe', ?, ?)`

	policzNoweCentrum = `SELECT COUNT(*) FROM powiadomienie_centrum WHERE stan = 'nowe'`

	// Zbiorcze oznaczenie odczytania: wykaz pusty bierze wszystkie nowe.
	odczytajWszystkieCentrum = `UPDATE powiadomienie_centrum
	                               SET stan = 'odczytane', odlozone_do = NULL
	                             WHERE stan = 'nowe'`

	zamknijZdarzenieCentrum = `UPDATE powiadomienie_centrum
	                              SET stan = 'obsluzone', odlozone_do = NULL
	                            WHERE id = ? AND stan <> 'obsluzone'`

	odlozZdarzenieCentrum = `UPDATE powiadomienie_centrum
	                            SET stan = 'odlozone', odlozone_do = ?
	                          WHERE id = ? AND stan <> 'obsluzone'`

	przywrocOdlozoneCentrum = `UPDATE powiadomienie_centrum
	                              SET stan = 'nowe', odlozone_do = NULL
	                            WHERE stan = 'odlozone' AND odlozone_do <= ?`
)

type repozytoriumCentrumPowiadomien struct {
	zapytania *zapytania
}

func noweRepozytoriumCentrumPowiadomien(z *zapytania) *repozytoriumCentrumPowiadomien {
	return &repozytoriumCentrumPowiadomien{zapytania: z}
}

func (r *repozytoriumCentrumPowiadomien) Zapisz(ctx context.Context,
	zdarzenie ZdarzenieCentrum) (ZdarzenieCentrum, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, zapiszZdarzenieCentrum)
	if err != nil {
		return ZdarzenieCentrum{}, err
	}
	if zdarzenie.ZnacznikCzasu.IsZero() {
		zdarzenie.ZnacznikCzasu = time.Now().UTC()
	}
	if strings.TrimSpace(zdarzenie.Akcje) == "" {
		zdarzenie.Akcje = "[]"
	}
	wynik, err := polecenie.ExecContext(ctx, zdarzenie.Klasa, zdarzenie.Waga, zdarzenie.Tresc,
		zdarzenie.ZrodloTyp, zdarzenie.ZrodloID, zdarzenie.SrodowiskoKod, zdarzenie.SesjaID,
		zdarzenie.Akcje, zdarzenie.Kanal, chwilaTekstem(zdarzenie.ZnacznikCzasu))
	if err != nil {
		return ZdarzenieCentrum{}, fmt.Errorf("dane: nie można zapisać zdarzenia centrum: %w", err)
	}
	numer, err := wynik.LastInsertId()
	if err != nil {
		return ZdarzenieCentrum{}, fmt.Errorf("dane: zdarzenie centrum bez identyfikatora: %w", err)
	}
	zdarzenie.ID = numer
	zdarzenie.Stan = "nowe"
	return zdarzenie, nil
}

func (r *repozytoriumCentrumPowiadomien) Wykaz(ctx context.Context,
	filtr FiltrCentrum) ([]ZdarzenieCentrum, error) {

	warunki := []string{}
	parametry := []any{}

	// Stan pominięty nie znaczy „wszystko": zdarzenia obsłużone są zamknięte
	// i kolumna centrum ich nie pokazuje, dopóki Operator wprost o nie nie pyta.
	stany := filtr.Stany
	if len(stany) == 0 {
		stany = []string{"nowe", "odczytane", "odlozone"}
	}
	warunki = append(warunki, wKolumnie("stan", len(stany)))
	parametry = append(parametry, naArgumenty(stany)...)

	if len(filtr.Klasy) > 0 {
		warunki = append(warunki, wKolumnie("klasa", len(filtr.Klasy)))
		parametry = append(parametry, naArgumenty(filtr.Klasy)...)
	}
	if len(filtr.Wagi) > 0 {
		warunki = append(warunki, wKolumnie("waga", len(filtr.Wagi)))
		parametry = append(parametry, naArgumenty(filtr.Wagi)...)
	}
	if filtr.SrodowiskoKod != "" {
		warunki = append(warunki, "srodowisko_kod = ?")
		parametry = append(parametry, filtr.SrodowiskoKod)
	}
	if filtr.SesjaID != "" {
		warunki = append(warunki, "sesja_id = ?")
		parametry = append(parametry, filtr.SesjaID)
	}

	limit := filtr.Limit
	if limit <= 0 {
		limit = 200
	}
	polecenieSQL := `SELECT ` + kolumnyCentrum + ` FROM powiadomienie_centrum WHERE ` +
		strings.Join(warunki, " AND ") +
		` ORDER BY znacznik_czasu DESC, id DESC LIMIT ` + strconv.Itoa(limit)

	polecenie, err := r.zapytania.przygotuj(ctx, polecenieSQL)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, parametry...)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać rejestru centrum: %w", err)
	}
	defer wiersze.Close()

	lista := []ZdarzenieCentrum{}
	for wiersze.Next() {
		zdarzenie, err := odczytajZdarzenieCentrum(wiersze)
		if err != nil {
			return nil, err
		}
		lista = append(lista, zdarzenie)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt rejestru centrum: %w", err)
	}
	return lista, nil
}

func (r *repozytoriumCentrumPowiadomien) Nowe(ctx context.Context) (int, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, policzNoweCentrum)
	if err != nil {
		return 0, err
	}
	var ile int
	if err := polecenie.QueryRowContext(ctx).Scan(&ile); err != nil {
		return 0, fmt.Errorf("dane: nie można policzyć zdarzeń nowych: %w", err)
	}
	return ile, nil
}

func (r *repozytoriumCentrumPowiadomien) Odczytaj(ctx context.Context,
	identyfikatory []int64) (int, error) {

	polecenieSQL := odczytajWszystkieCentrum
	parametry := []any{}
	if len(identyfikatory) > 0 {
		polecenieSQL = `UPDATE powiadomienie_centrum
		                   SET stan = 'odczytane', odlozone_do = NULL
		                 WHERE stan IN ('nowe','odlozone') AND ` +
			wKolumnie("id", len(identyfikatory))
		for _, numer := range identyfikatory {
			parametry = append(parametry, numer)
		}
	}
	polecenie, err := r.zapytania.przygotuj(ctx, polecenieSQL)
	if err != nil {
		return 0, err
	}
	wynik, err := polecenie.ExecContext(ctx, parametry...)
	if err != nil {
		return 0, fmt.Errorf("dane: nie można oznaczyć zdarzeń jako odczytane: %w", err)
	}
	return liczbaZmian(wynik)
}

func (r *repozytoriumCentrumPowiadomien) Zamknij(ctx context.Context, id int64) (bool, error) {
	ile, err := r.wykonajNaZdarzeniu(ctx, zamknijZdarzenieCentrum, id)
	return ile > 0, err
}

func (r *repozytoriumCentrumPowiadomien) Odloz(ctx context.Context, id int64,
	doKiedy time.Time) (bool, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, odlozZdarzenieCentrum)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, chwilaTekstem(doKiedy.UTC()), id)
	if err != nil {
		return false, fmt.Errorf("dane: nie można odłożyć zdarzenia centrum: %w", err)
	}
	ile, err := liczbaZmian(wynik)
	return ile > 0, err
}

func (r *repozytoriumCentrumPowiadomien) Przywroc(ctx context.Context,
	teraz time.Time) (int, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, przywrocOdlozoneCentrum)
	if err != nil {
		return 0, err
	}
	wynik, err := polecenie.ExecContext(ctx, chwilaTekstem(teraz.UTC()))
	if err != nil {
		return 0, fmt.Errorf("dane: nie można przywrócić zdarzeń odłożonych: %w", err)
	}
	return liczbaZmian(wynik)
}

func (r *repozytoriumCentrumPowiadomien) wykonajNaZdarzeniu(ctx context.Context,
	polecenieSQL string, id int64) (int, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, polecenieSQL)
	if err != nil {
		return 0, err
	}
	wynik, err := polecenie.ExecContext(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("dane: nie można zmienić stanu zdarzenia centrum: %w", err)
	}
	return liczbaZmian(wynik)
}

// odczytajZdarzenieCentrum składa wiersz rejestru.
func odczytajZdarzenieCentrum(wiersz skaner) (ZdarzenieCentrum, error) {
	var zdarzenie ZdarzenieCentrum
	var odlozone, znacznik string
	err := wiersz.Scan(&zdarzenie.ID, &zdarzenie.Klasa, &zdarzenie.Waga, &zdarzenie.Tresc,
		&zdarzenie.ZrodloTyp, &zdarzenie.ZrodloID, &zdarzenie.SrodowiskoKod, &zdarzenie.SesjaID,
		&zdarzenie.Akcje, &zdarzenie.Stan, &zdarzenie.Kanal, &odlozone, &znacznik)
	if err != nil {
		return ZdarzenieCentrum{}, fmt.Errorf("dane: nieczytelny wiersz rejestru centrum: %w", err)
	}
	zdarzenie.OdlozoneDo = chwilaZTekstu(odlozone)
	zdarzenie.ZnacznikCzasu = chwilaZTekstu(znacznik)
	return zdarzenie, nil
}

// liczbaZmian oddaje liczbę wierszy zmienionych poleceniem.
func liczbaZmian(wynik sql.Result) (int, error) {
	ile, err := wynik.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("dane: nie można policzyć zmienionych wierszy: %w", err)
	}
	return int(ile), nil
}

// wKolumnie składa warunek `kolumna IN (?, ?, …)` o zadanej liczbie miejsc.
func wKolumnie(kolumna string, ile int) string {
	return kolumna + " IN (" + strings.TrimSuffix(strings.Repeat("?,", ile), ",") + ")"
}

// naArgumenty przenosi wykaz napisów do argumentów zapytania.
func naArgumenty(wartosci []string) []any {
	argumenty := make([]any, 0, len(wartosci))
	for _, wartosc := range wartosci {
		argumenty = append(argumenty, wartosc)
	}
	return argumenty
}

// chwilaTekstem zapisuje chwilę w postaci używanej przez kolumny tej bazy.
func chwilaTekstem(chwila time.Time) string {
	return chwila.UTC().Format("2006-01-02T15:04:05.000Z")
}

// chwilaZTekstu odczytuje chwilę; zapis pusty albo nieczytelny daje chwilę zerową.
func chwilaZTekstu(zapis string) time.Time {
	if strings.TrimSpace(zapis) == "" {
		return time.Time{}
	}
	for _, uklad := range []string{"2006-01-02T15:04:05.000Z", time.RFC3339Nano, time.RFC3339} {
		if chwila, err := time.Parse(uklad, zapis); err == nil {
			return chwila.UTC()
		}
	}
	return time.Time{}
}
