// Odpowiedzialność pliku: sesje bramki (tabela `sesja_bramki`) — druga połowa
// repozytorium bramki. Katalog metod wejścia leży w `auth.go`.
//
// Sesja bramki nie jest kartą sesji. Tabela `sesja` opisuje pracę: rozmowę,
// okna, kolejki, projekt. Wiersz poniżej opisuje wejście: kto przeszedł bramkę,
// kiedy i jaką metodą. Karta pracy istnieje niezależnie od tego, czy ktokolwiek
// się zalogował, a wygaśnięcie wejścia rozmowy nie zabiera.
//
// Struktura niesie wyłącznie skrót tokenu; token surowy zna rdzeń w chwili
// założenia sesji oraz klient, któremu go oddano. Repozytorium nie ma jak go
// odtworzyć — taki jest zamiar.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// SesjaBramki to wiersz tabeli `sesja_bramki`. Czasy idą w milisekundach epoki
// — w tej samej jednostce, w której niesie je kontrakt (`AuthSession.expiresAt`).
type SesjaBramki struct {
	ID            int64
	SkrotTokenu   string
	RodzajMetody  *string
	UrzadzenieKod *string
	Wygasa        int64
	Utworzono     int64
	Uniewazniono  *int64
	// Trwanie to długość życia sesji w milisekundach — ta sama, którą wybiera
	// przełącznik „nie wyloguj mnie" przy zakładaniu. Jest zapisana, bo wygasanie
	// jest przesuwne: odnowienie musi wiedzieć, o ile przesunąć koniec, a sama
	// chwila `Wygasa` tego nie mówi. Zero znaczy wiersz bez zapisanego trwania —
	// wołający bierze wtedy trwanie podstawowe.
	Trwanie int64
}

const (
	kolumnySesjiBramki = `id, token_skrot, metoda_rodzaj, urzadzenie_kod,
	                      wygasa, utworzono, uniewazniono, trwanie`

	wstawSesjeBramki = `INSERT INTO sesja_bramki
	                    (token_skrot, metoda_rodzaj, urzadzenie_kod, wygasa, utworzono, trwanie)
	                    VALUES (?, ?, ?, ?, ?, ?)`

	sesjaBramkiPoSkrocie = `SELECT ` + kolumnySesjiBramki +
		` FROM sesja_bramki WHERE token_skrot = ?`

	przedluzSesjeBramki = `UPDATE sesja_bramki SET wygasa = ?
	                       WHERE token_skrot = ? AND uniewazniono IS NULL`

	// Pusty napis w pierwszym argumencie zdejmuje wyłączenie — wtedy polecenie
	// unieważnia komplet sesji czynnych. Jedno zapytanie na oba warianty, żeby
	// „poza bieżącą" i „wszystkie" nie rozjechały się na dwie ścieżki.
	uniewaznijSesjeBramki = `UPDATE sesja_bramki SET uniewazniono = ?
	                         WHERE uniewazniono IS NULL AND (? = '' OR token_skrot <> ?)`

	// Wykaz urządzeń powstaje z sesji, nie z osobnej tabeli: urządzeniem konta
	// jest to, które kiedykolwiek weszło. Grupowanie po kodzie daje jeden wiersz
	// na urządzenie, a nie jeden na każde logowanie.
	//
	// `MAX(uniewazniono IS NULL AND wygasa > ?)` mówi, czy urządzenie ma DZIŚ
	// ważny token — czyli czy unieważnienie ma co odbierać.
	urzadzeniaSesjiBramki = `SELECT urzadzenie_kod,
	                                MAX(utworzono) AS ostatnio,
	                                MAX(CASE WHEN uniewazniono IS NULL AND wygasa > ?
	                                         THEN 1 ELSE 0 END) AS czynny
	                         FROM sesja_bramki
	                         WHERE urzadzenie_kod IS NOT NULL AND urzadzenie_kod <> ''
	                         GROUP BY urzadzenie_kod
	                         ORDER BY ostatnio DESC`

	uniewaznijSesjeUrzadzenia = `UPDATE sesja_bramki SET uniewazniono = ?
	                             WHERE uniewazniono IS NULL AND urzadzenie_kod = ?`
)

// ZalozSesjeBramki zakłada sesję wejścia i oddaje ją odczytaną z bazy.
func (r *repozytoriumUwierzytelnienia) ZalozSesjeBramki(ctx context.Context,
	sesja SesjaBramki) (SesjaBramki, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, wstawSesjeBramki)
	if err != nil {
		return SesjaBramki{}, err
	}
	_, err = polecenie.ExecContext(ctx, sesja.SkrotTokenu,
		tekstDoKolumny(sesja.RodzajMetody), tekstDoKolumny(sesja.UrzadzenieKod),
		sesja.Wygasa, sesja.Utworzono, sesja.Trwanie)
	if err != nil {
		return SesjaBramki{}, fmt.Errorf("dane: nie można założyć sesji bramki: %w", err)
	}
	return r.SesjaBramkiPoSkrocie(ctx, sesja.SkrotTokenu)
}

// SesjaBramkiPoSkrocie zwraca sesję rozpoznaną skrótem tokenu. Brak wiersza daje
// ErrBrakWiersza — komunikat nie niesie ani tokenu, ani jego skrótu, bo trafia
// do odpowiedzi i do dziennika.
func (r *repozytoriumUwierzytelnienia) SesjaBramkiPoSkrocie(ctx context.Context,
	skrot string) (SesjaBramki, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, sesjaBramkiPoSkrocie)
	if err != nil {
		return SesjaBramki{}, err
	}
	sesja, err := odczytajSesjeBramki(polecenie.QueryRowContext(ctx, skrot))
	if errors.Is(err, sql.ErrNoRows) {
		return SesjaBramki{}, fmt.Errorf("dane: sesja bramki nie istnieje: %w", ErrBrakWiersza)
	}
	if err != nil {
		return SesjaBramki{}, fmt.Errorf("dane: nieczytelny wiersz sesji bramki: %w", err)
	}
	return sesja, nil
}

// PrzedluzSesjeBramki przesuwa wygaśnięcie sesji czynnej. Sesja unieważniona
// przedłużeniu nie podlega: warunek stoi w zapytaniu, więc odmowa nie zależy
// od kolejności gałęzi w Go.
func (r *repozytoriumUwierzytelnienia) PrzedluzSesjeBramki(ctx context.Context,
	skrot string, wygasa int64) (SesjaBramki, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, przedluzSesjeBramki)
	if err != nil {
		return SesjaBramki{}, err
	}
	wynik, err := polecenie.ExecContext(ctx, wygasa, skrot)
	if err != nil {
		return SesjaBramki{}, fmt.Errorf("dane: nie można przedłużyć sesji bramki: %w", err)
	}
	zmienione, err := wynik.RowsAffected()
	if err != nil {
		return SesjaBramki{}, fmt.Errorf("dane: nieczytelny wynik przedłużenia sesji bramki: %w", err)
	}
	if zmienione == 0 {
		return SesjaBramki{}, fmt.Errorf("dane: czynna sesja bramki nie istnieje: %w", ErrBrakWiersza)
	}
	return r.SesjaBramkiPoSkrocie(ctx, skrot)
}

// UniewaznijSesjeBramkiPoza unieważnia sesje czynne poza wskazaną i zwraca ich
// liczbę — tej liczby żąda `auth.password.reset` polem `revokedSessions`.
func (r *repozytoriumUwierzytelnienia) UniewaznijSesjeBramkiPoza(ctx context.Context,
	skrotZachowany string, teraz int64) (int, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, uniewaznijSesjeBramki)
	if err != nil {
		return 0, err
	}
	wynik, err := polecenie.ExecContext(ctx, teraz, skrotZachowany, skrotZachowany)
	if err != nil {
		return 0, fmt.Errorf("dane: nie można unieważnić sesji bramki: %w", err)
	}
	zmienione, err := wynik.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("dane: nieczytelny wynik unieważnienia sesji bramki: %w", err)
	}
	return int(zmienione), nil
}

// UrzadzeniaKonta zwraca urządzenia, które kiedykolwiek weszły przez bramkę.
func (r *repozytoriumUwierzytelnienia) UrzadzeniaKonta(ctx context.Context,
	teraz int64) ([]UrzadzenieKonta, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, urzadzeniaSesjiBramki)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, teraz)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać urządzeń konta: %w", err)
	}
	defer wiersze.Close()

	lista := []UrzadzenieKonta{}
	for wiersze.Next() {
		var urzadzenie UrzadzenieKonta
		var czynny int64
		if err := wiersze.Scan(&urzadzenie.Kod, &urzadzenie.OstatnioWidziane, &czynny); err != nil {
			return nil, fmt.Errorf("dane: uszkodzony wiersz wykazu urządzeń: %w", err)
		}
		urzadzenie.MaToken = czynny != 0
		lista = append(lista, urzadzenie)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt urządzeń konta: %w", err)
	}
	return lista, nil
}

// UniewaznijSesjeUrzadzenia zamyka sesje czynne wskazanego urządzenia.
func (r *repozytoriumUwierzytelnienia) UniewaznijSesjeUrzadzenia(ctx context.Context,
	urzadzenie string, teraz int64) (int, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, uniewaznijSesjeUrzadzenia)
	if err != nil {
		return 0, err
	}
	wynik, err := polecenie.ExecContext(ctx, teraz, urzadzenie)
	if err != nil {
		return 0, fmt.Errorf("dane: nie można unieważnić sesji urządzenia %q: %w", urzadzenie, err)
	}
	zmienione, err := wynik.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("dane: nieczytelny wynik unieważnienia sesji urządzenia: %w", err)
	}
	return int(zmienione), nil
}

// odczytajSesjeBramki składa strukturę z jednego wiersza wyniku.
func odczytajSesjeBramki(wiersz skaner) (SesjaBramki, error) {
	var sesja SesjaBramki
	var rodzaj, urzadzenie sql.NullString
	var uniewazniono sql.NullInt64
	err := wiersz.Scan(&sesja.ID, &sesja.SkrotTokenu, &rodzaj, &urzadzenie,
		&sesja.Wygasa, &sesja.Utworzono, &uniewazniono, &sesja.Trwanie)
	if err != nil {
		return SesjaBramki{}, err
	}
	sesja.RodzajMetody = tekstZKolumny(rodzaj)
	sesja.UrzadzenieKod = tekstZKolumny(urzadzenie)
	sesja.Uniewazniono = liczbaZKolumny(uniewazniono)
	return sesja, nil
}
