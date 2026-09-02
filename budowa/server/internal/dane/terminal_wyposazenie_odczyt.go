// Repozytorium czyta wyposażenie modułu Terminal: książkę hostów, bibliotekę
// skryptów, wykaz kluczy, tuneli i obserwacji, zawężane parametrem zapytania.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

const (
	kolumnyHostaTerminala = `kod, nazwa, cel, port, grupa, katalog_roboczy, klucz_kod,
	                         host_posredni_kod, notatka, utworzono, zaktualizowano`

	pobierzHostaTerminala = `SELECT ` + kolumnyHostaTerminala + `
	                         FROM terminal_host WHERE kod = ? AND ` + WarunekKonta

	// Fraza szuka w nazwie, adresie celu i notatce. LIKE w SQLite nie rozróżnia
	// wielkości liter dla znaków ASCII; dla pozostałych schodzi do porównania
	// dosłownego — i to jest cała obietnica, jaką książka hostów składa.
	pobierzHostyTerminala = `SELECT ` + kolumnyHostaTerminala + `
	                         FROM terminal_host
	                         WHERE (? = '' OR grupa = ?)
	                           AND (? = '' OR nazwa LIKE ? OR cel LIKE ? OR notatka LIKE ?)
	                           AND ` + WarunekKonta + `
	                         ORDER BY grupa ASC, nazwa ASC, id ASC`

	kolumnySkryptuTerminala = `kod, nazwa, rodzaj, powloka, tresc, znaczniki, alias,
	                           wersja, uruchomiono, utworzono, zaktualizowano`

	pobierzSkryptTerminala = `SELECT ` + kolumnySkryptuTerminala + `
	                          FROM terminal_skrypt WHERE kod = ?`

	// Znacznik dopasowuje się do jednej pozycji wykazu rozdzielanego przecinkiem.
	// Otoczenie obu stron przecinkami sprawia, że `test` nie łapie `testowy`.
	pobierzSkryptyTerminala = `SELECT ` + kolumnySkryptuTerminala + `
	                           FROM terminal_skrypt
	                           WHERE (? = '' OR rodzaj = ?)
	                             AND (? = '' OR (',' || replace(znaczniki, ' ', '') || ',') LIKE ?)
	                             AND (? = '' OR nazwa LIKE ? OR tresc LIKE ?)
	                           ORDER BY nazwa ASC, id ASC`

	kolumnyKluczaTerminala = `kod, nazwa, rodzaj, odcisk, klucz_jawny, sciezka, haslo, utworzono`

	pobierzKluczTerminala = `SELECT ` + kolumnyKluczaTerminala + `
	                         FROM terminal_klucz WHERE kod = ? AND ` + WarunekKonta

	pobierzKluczeTerminala = `SELECT ` + kolumnyKluczaTerminala + `
	                          FROM terminal_klucz WHERE ` + WarunekKonta + `
	                          ORDER BY nazwa ASC, id ASC`

	kolumnyTuneluTerminala = `kod, okno_kod, rodzaj, host_kod, cel, port_lokalny,
	                          host_docelowy, port_docelowy, stan, powod, zalozono, zamknieto`

	pobierzTunelTerminala = `SELECT ` + kolumnyTuneluTerminala + `
	                         FROM terminal_tunel WHERE kod = ?`

	pobierzTuneleTerminala = `SELECT ` + kolumnyTuneluTerminala + `
	                          FROM terminal_tunel
	                          WHERE (? = '' OR okno_kod = ?)
	                            AND (? = '' OR stan = ?)
	                          ORDER BY zalozono DESC, id DESC`

	kolumnyObserwacjiTerminala = `kod, okno_kod, karta_kod, wzorzec, polecenie, tlumienie,
	                              rekurencyjnie, stan, licznik, wyzwolono, powod, zalozono`

	pobierzObserwacjeTerminala = `SELECT ` + kolumnyObserwacjiTerminala + `
	                              FROM terminal_obserwacja WHERE kod = ?`

	pobierzObserwacjeTerminalaWykaz = `SELECT ` + kolumnyObserwacjiTerminala + `
	                                   FROM terminal_obserwacja
	                                   WHERE (? = '' OR okno_kod = ?)
	                                     AND (? = '' OR stan = ?)
	                                   ORDER BY zalozono DESC, id DESC`
)

// wzorzecFrazyTerminala otacza frazę znakami wieloznacznymi. Fraza pusta zostaje
// pusta, bo pusty parametr wyłącza całe zawężenie.
func wzorzecFrazyTerminala(fraza string) string {
	if fraza == "" {
		return ""
	}
	return "%" + fraza + "%"
}

// Host zwraca z tabeli `terminal_host` wpis książki hostów modułu Terminal,
// wskazany parametrem kodu identyfikującego wiersz.
func (r *repozytoriumTerminala) Host(ctx context.Context, kod string) (HostTerminala, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzHostaTerminala)
	if err != nil {
		return HostTerminala{}, err
	}
	host, err := odczytajHostaTerminala(polecenie.QueryRowContext(ctx, kod, KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return HostTerminala{}, ErrBrakWiersza
	}
	if err != nil {
		return HostTerminala{}, fmt.Errorf("dane: nieczytelny wiersz hosta terminala %q: %w", kod, err)
	}
	return host, nil
}

// Hosty zwraca całą książkę hostów modułu Terminal z tabeli `terminal_host`,
// zawężoną przekazanym filtrem wyszukiwania.
func (r *repozytoriumTerminala) Hosty(ctx context.Context,
	filtr FiltrHostowTerminala) ([]HostTerminala, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzHostyTerminala)
	if err != nil {
		return nil, err
	}
	wzorzec := wzorzecFrazyTerminala(filtr.Fraza)
	wiersze, err := polecenie.QueryContext(ctx,
		filtr.Grupa, filtr.Grupa,
		filtr.Fraza, wzorzec, wzorzec, wzorzec,
		KontoOperatora(ctx))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać książki hostów: %w", err)
	}
	defer wiersze.Close()

	hosty := make([]HostTerminala, 0, 8)
	for wiersze.Next() {
		host, err := odczytajHostaTerminala(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz hosta terminala: %w", err)
		}
		hosty = append(hosty, host)
	}
	return hosty, wiersze.Err()
}

// Skrypt zwraca z tabeli `terminal_skrypt` pozycję biblioteki skryptów modułu
// Terminal, wskazaną parametrem kodu.
func (r *repozytoriumTerminala) Skrypt(ctx context.Context, kod string) (SkryptTerminala, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzSkryptTerminala)
	if err != nil {
		return SkryptTerminala{}, err
	}
	skrypt, err := odczytajSkryptTerminala(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return SkryptTerminala{}, ErrBrakWiersza
	}
	if err != nil {
		return SkryptTerminala{}, fmt.Errorf("dane: nieczytelny wiersz skryptu %q: %w", kod, err)
	}
	return skrypt, nil
}

// Skrypty zwraca bibliotekę skryptów modułu Terminal zawężoną filtrem,
// w brzmieniu bieżącym każdej pozycji.
func (r *repozytoriumTerminala) Skrypty(ctx context.Context,
	filtr FiltrSkryptowTerminala) ([]SkryptTerminala, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzSkryptyTerminala)
	if err != nil {
		return nil, err
	}
	rodzaj := string(filtr.Rodzaj)
	wzorzecZnacznika := ""
	if filtr.Znacznik != "" {
		wzorzecZnacznika = "%," + filtr.Znacznik + ",%"
	}
	wzorzecFrazy := wzorzecFrazyTerminala(filtr.Fraza)
	wiersze, err := polecenie.QueryContext(ctx,
		rodzaj, rodzaj,
		filtr.Znacznik, wzorzecZnacznika,
		filtr.Fraza, wzorzecFrazy, wzorzecFrazy)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać biblioteki skryptów: %w", err)
	}
	defer wiersze.Close()

	skrypty := make([]SkryptTerminala, 0, 8)
	for wiersze.Next() {
		skrypt, err := odczytajSkryptTerminala(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz skryptu: %w", err)
		}
		skrypty = append(skrypty, skrypt)
	}
	return skrypty, wiersze.Err()
}

// Klucz zwraca z tabeli `terminal_klucz` wpis wykazu kluczy SSH modułu
// Terminal, wskazany parametrem kodu.
func (r *repozytoriumTerminala) Klucz(ctx context.Context, kod string) (KluczTerminala, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzKluczTerminala)
	if err != nil {
		return KluczTerminala{}, err
	}
	klucz, err := odczytajKluczTerminala(polecenie.QueryRowContext(ctx, kod, KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return KluczTerminala{}, ErrBrakWiersza
	}
	if err != nil {
		return KluczTerminala{}, fmt.Errorf("dane: nieczytelny wiersz klucza %q: %w", kod, err)
	}
	return klucz, nil
}

// Klucze zwraca z tabeli `terminal_klucz` cały wykaz kluczy SSH modułu
// Terminal, uporządkowany rosnąco według nazwy.
func (r *repozytoriumTerminala) Klucze(ctx context.Context) ([]KluczTerminala, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzKluczeTerminala)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, KontoOperatora(ctx))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać wykazu kluczy: %w", err)
	}
	defer wiersze.Close()

	klucze := make([]KluczTerminala, 0, 8)
	for wiersze.Next() {
		klucz, err := odczytajKluczTerminala(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz klucza: %w", err)
		}
		klucze = append(klucze, klucz)
	}
	return klucze, wiersze.Err()
}

// Tunel zwraca z tabeli `terminal_tunel` przekierowanie portu modułu Terminal,
// wskazane parametrem kodu.
func (r *repozytoriumTerminala) Tunel(ctx context.Context, kod string) (TunelTerminala, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzTunelTerminala)
	if err != nil {
		return TunelTerminala{}, err
	}
	tunel, err := odczytajTunelTerminala(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return TunelTerminala{}, ErrBrakWiersza
	}
	if err != nil {
		return TunelTerminala{}, fmt.Errorf("dane: nieczytelny wiersz tunelu %q: %w", kod, err)
	}
	return tunel, nil
}

// Tunele zwraca z tabeli `terminal_tunel` przekierowania portów modułu
// Terminal zawężone filtrem, od najnowszego.
func (r *repozytoriumTerminala) Tunele(ctx context.Context,
	filtr FiltrTuneliTerminala) ([]TunelTerminala, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzTuneleTerminala)
	if err != nil {
		return nil, err
	}
	stan := string(filtr.Stan)
	wiersze, err := polecenie.QueryContext(ctx, filtr.OknoKod, filtr.OknoKod, stan, stan)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać wykazu tuneli: %w", err)
	}
	defer wiersze.Close()

	tunele := make([]TunelTerminala, 0, 8)
	for wiersze.Next() {
		tunel, err := odczytajTunelTerminala(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz tunelu: %w", err)
		}
		tunele = append(tunele, tunel)
	}
	return tunele, wiersze.Err()
}

// Obserwacja zwraca obserwację plików modułu Terminal o wskazanym kodzie,
// z tabeli `terminal_obserwacja`.
func (r *repozytoriumTerminala) Obserwacja(ctx context.Context, kod string) (ObserwacjaTerminala, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzObserwacjeTerminala)
	if err != nil {
		return ObserwacjaTerminala{}, err
	}
	obserwacja, err := odczytajObserwacjeTerminala(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return ObserwacjaTerminala{}, ErrBrakWiersza
	}
	if err != nil {
		return ObserwacjaTerminala{}, fmt.Errorf("dane: nieczytelny wiersz obserwacji %q: %w", kod, err)
	}
	return obserwacja, nil
}

// Obserwacje zwraca z tabeli `terminal_obserwacja` obserwacje plików modułu
// Terminal zawężone filtrem, od najnowszej.
func (r *repozytoriumTerminala) Obserwacje(ctx context.Context,
	filtr FiltrObserwacjiTerminala) ([]ObserwacjaTerminala, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzObserwacjeTerminalaWykaz)
	if err != nil {
		return nil, err
	}
	stan := string(filtr.Stan)
	wiersze, err := polecenie.QueryContext(ctx, filtr.OknoKod, filtr.OknoKod, stan, stan)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać wykazu obserwacji: %w", err)
	}
	defer wiersze.Close()

	obserwacje := make([]ObserwacjaTerminala, 0, 8)
	for wiersze.Next() {
		obserwacja, err := odczytajObserwacjeTerminala(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz obserwacji: %w", err)
		}
		obserwacje = append(obserwacje, obserwacja)
	}
	return obserwacje, wiersze.Err()
}

// odczytajHostaTerminala składa strukturę HostTerminala z jednego wiersza
// wyniku zapytania do tabeli `terminal_host`.
func odczytajHostaTerminala(wiersz interface{ Scan(...any) error }) (HostTerminala, error) {
	var host HostTerminala
	err := wiersz.Scan(&host.Kod, &host.Nazwa, &host.Cel, &host.Port, &host.Grupa,
		&host.KatalogRoboczy, &host.KluczKod, &host.HostPosredniKod, &host.Notatka,
		&host.Utworzono, &host.Zaktualizowano)
	return host, err
}

// odczytajSkryptTerminala składa pozycję struktury SkryptTerminala z jednego
// wiersza wyniku zapytania SQL.
func odczytajSkryptTerminala(wiersz interface{ Scan(...any) error }) (SkryptTerminala, error) {
	var skrypt SkryptTerminala
	err := wiersz.Scan(&skrypt.Kod, &skrypt.Nazwa, &skrypt.Rodzaj, &skrypt.Powloka,
		&skrypt.Tresc, &skrypt.Znaczniki, &skrypt.Alias, &skrypt.Wersja,
		&skrypt.Uruchomiono, &skrypt.Utworzono, &skrypt.Zaktualizowano)
	return skrypt, err
}

// odczytajKluczTerminala składa strukturę KluczTerminala z jednego wiersza
// wyniku zapytania do tabeli `terminal_klucz`.
func odczytajKluczTerminala(wiersz interface{ Scan(...any) error }) (KluczTerminala, error) {
	var klucz KluczTerminala
	err := wiersz.Scan(&klucz.Kod, &klucz.Nazwa, &klucz.Rodzaj, &klucz.Odcisk,
		&klucz.KluczJawny, &klucz.Sciezka, &klucz.Haslo, &klucz.Utworzono)
	return klucz, err
}

// odczytajTunelTerminala składa przekierowanie struktury TunelTerminala
// z jednego wiersza wyniku zapytania SQL.
func odczytajTunelTerminala(wiersz interface{ Scan(...any) error }) (TunelTerminala, error) {
	var tunel TunelTerminala
	err := wiersz.Scan(&tunel.Kod, &tunel.OknoKod, &tunel.Rodzaj, &tunel.HostKod,
		&tunel.Cel, &tunel.PortLokalny, &tunel.HostDocelowy, &tunel.PortDocelowy,
		&tunel.Stan, &tunel.Powod, &tunel.Zalozono, &tunel.Zamknieto)
	return tunel, err
}

// odczytajObserwacjeTerminala składa obserwację ObserwacjaTerminala z jednego
// wiersza wyniku zapytania SQL.
func odczytajObserwacjeTerminala(wiersz interface{ Scan(...any) error }) (ObserwacjaTerminala, error) {
	var obserwacja ObserwacjaTerminala
	err := wiersz.Scan(&obserwacja.Kod, &obserwacja.OknoKod, &obserwacja.KartaKod,
		&obserwacja.Wzorzec, &obserwacja.Polecenie, &obserwacja.Tlumienie,
		&obserwacja.Rekurencyjnie, &obserwacja.Stan, &obserwacja.Licznik,
		&obserwacja.Wyzwolono, &obserwacja.Powod, &obserwacja.Zalozono)
	return obserwacja, err
}
