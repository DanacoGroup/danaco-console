// Plik prowadzi katalog metod wejścia przez bramkę: hasło, PIN urządzenia albo klucz Windows Hello; sesje bramki
// leżą obok, w auth_sesje_bramki.go, a odwołanie do sejfu poświadczeń repozytorium traktuje jak nieprzezroczysty napis.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// MetodaUwierzytelnienia to wiersz tabeli `metoda_uwierzytelnienia`; kod jest identyfikatorem trwałym, kod urządzenia niesie napis nadany przez klienta.
type MetodaUwierzytelnienia struct {
	ID               int64
	Kod              string
	Rodzaj           string
	Etykieta         *string
	UrzadzenieKod    *string
	NazwaUrzadzenia  *string
	Kotwica          bool
	OdwolanieSekretu string
	Utworzono        int64
	OstatnioUzyto    *int64
}

// RepozytoriumUwierzytelnienia jest kontraktem bramki: katalog metod
// wejścia oraz sesje bramki w jednym bycie, bo obie tabele zmieniają się razem
// — zmiana hasła unieważnia sesje, a wejście odnotowuje użycie metody.
type RepozytoriumUwierzytelnienia interface {
	// Metody zwraca komplet metod wejścia w kolejności wyświetlania.
	Metody(ctx context.Context) ([]MetodaUwierzytelnienia, error)
	// MetodaPoKodzie zwraca metodę wskazaną identyfikatorem trwałym; brak wiersza daje ErrBrakWiersza.
	MetodaPoKodzie(ctx context.Context, kod string) (MetodaUwierzytelnienia, error)
	// Kotwica zwraca hasło bramki; brak wiersza znaczy bramkę nieustawioną, stan otwierający rejestrację.
	Kotwica(ctx context.Context) (MetodaUwierzytelnienia, error)
	// MetodaUrzadzenia zwraca metodę danego rodzaju założoną na urządzeniu.
	MetodaUrzadzenia(ctx context.Context, rodzaj, urzadzenieKod string) (MetodaUwierzytelnienia, error)
	// ZalozMetode wstawia wiersz metody i oddaje go w postaci zapisanej.
	ZalozMetode(ctx context.Context, metoda MetodaUwierzytelnienia) (MetodaUwierzytelnienia, error)
	// ZapiszOdwolanieSekretu podmienia odwołanie do sejfu przy zmianie hasła, bez zakładania kotwicy.
	ZapiszOdwolanieSekretu(ctx context.Context, kod, odwolanie string) error
	// OdnotujUzycie zapisuje czas ostatniego wejścia tą metodą.
	OdnotujUzycie(ctx context.Context, kod string, teraz int64) error
	// UsunMetode kasuje metodę; drugi wynik mówi, czy wiersz istniał, dla odróżnienia od braku wiersza.
	UsunMetode(ctx context.Context, kod string) (bool, error)

	// ZalozSesjeBramki zakłada sesję wejścia; token nie wchodzi do bazy, wchodzi wyłącznie jego skrót.
	ZalozSesjeBramki(ctx context.Context, sesja SesjaBramki) (SesjaBramki, error)
	// SesjaBramkiPoSkrocie zwraca sesję rozpoznaną skrótem tokenu.
	SesjaBramkiPoSkrocie(ctx context.Context, skrot string) (SesjaBramki, error)
	// PrzedluzSesjeBramki przesuwa wygaśnięcie sesji i oddaje ją po zmianie.
	PrzedluzSesjeBramki(ctx context.Context, skrot string, wygasa int64) (SesjaBramki, error)
	// UniewaznijSesjeBramkiPoza unieważnia wszystkie sesje czynne poza wskazaną i zwraca ich liczbę.
	UniewaznijSesjeBramkiPoza(ctx context.Context, skrotZachowany string, teraz int64) (int, error)
	// UrzadzeniaKonta zwraca urządzenia, które weszły przez bramkę, z ostatnią chwilą wejścia.
	UrzadzeniaKonta(ctx context.Context, teraz int64) ([]UrzadzenieKonta, error)
	// UniewaznijSesjeUrzadzenia zamyka sesje czynne wskazanego urządzenia
	// i zwraca ich liczbę.
	UniewaznijSesjeUrzadzenia(ctx context.Context, urzadzenie string, teraz int64) (int, error)
}

// UrzadzenieKonta to jeden wiersz wykazu urządzeń powiązanych z kontem; urządzeniem konta jest to, które kiedykolwiek weszło przez bramkę.
type UrzadzenieKonta struct {
	Kod string
	// OstatnioWidziane to chwila ostatniego wejścia w milisekundach epoki.
	OstatnioWidziane int64
	// MaToken mówi, czy urządzenie ma dziś ważny, nieunieważniony token.
	MaToken bool
}

const (
	kolumnyMetodyUwierzytelnienia = `id, identyfikator_zewnetrzny, rodzaj, etykieta,
	                                 urzadzenie_kod, nazwa_urzadzenia, kotwica,
	                                 sekret_odwolanie, utworzono, ostatnio_uzyto`

	listaMetodUwierzytelnienia = `SELECT ` + kolumnyMetodyUwierzytelnienia +
		` FROM metoda_uwierzytelnienia ORDER BY kotwica DESC, rodzaj, utworzono, id`

	metodaUwierzytelnieniaPoKodzie = `SELECT ` + kolumnyMetodyUwierzytelnienia +
		` FROM metoda_uwierzytelnienia WHERE identyfikator_zewnetrzny = ?`

	metodaUwierzytelnieniaKotwica = `SELECT ` + kolumnyMetodyUwierzytelnienia +
		` FROM metoda_uwierzytelnienia WHERE kotwica = 1`

	metodaUwierzytelnieniaUrzadzenia = `SELECT ` + kolumnyMetodyUwierzytelnienia +
		` FROM metoda_uwierzytelnienia WHERE rodzaj = ? AND urzadzenie_kod = ?`

	wstawMetodeUwierzytelnienia = `INSERT INTO metoda_uwierzytelnienia
	                               (identyfikator_zewnetrzny, rodzaj, etykieta, urzadzenie_kod,
	                                nazwa_urzadzenia, kotwica, sekret_odwolanie, utworzono)
	                               VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	zapiszOdwolanieMetody = `UPDATE metoda_uwierzytelnienia SET sekret_odwolanie = ?
	                         WHERE identyfikator_zewnetrzny = ?`

	odnotujUzycieMetody = `UPDATE metoda_uwierzytelnienia SET ostatnio_uzyto = ?
	                       WHERE identyfikator_zewnetrzny = ?`

	usunMetodeUwierzytelnienia = `DELETE FROM metoda_uwierzytelnienia
	                              WHERE identyfikator_zewnetrzny = ?`
)

type repozytoriumUwierzytelnienia struct {
	zapytania *zapytania
}

// Zgodność implementacji repozytorium z kontraktem interfejsu sprawdzana jest przy kompilacji pakietu.
var _ RepozytoriumUwierzytelnienia = (*repozytoriumUwierzytelnienia)(nil)

// noweRepozytoriumUwierzytelnienia zakłada repozytorium bramki. `*sql.DB` nie
// jest potrzebne: żaden zapis bramki nie obejmuje dwóch tabel naraz, więc
// transakcji wielotabelowej tu nie ma.
func noweRepozytoriumUwierzytelnienia(z *zapytania) *repozytoriumUwierzytelnienia {
	return &repozytoriumUwierzytelnienia{zapytania: z}
}

// Metody zwraca komplet metod wejścia. Kotwica idzie pierwsza, bo sekcja
// Uwierzytelnianie Okna Ustawień pokazuje hasło na czele, a metody szybkiego
// wejścia pod nim. Katalog pusty nie jest błędem — znaczy bramkę nieustawioną.
func (r *repozytoriumUwierzytelnienia) Metody(ctx context.Context) ([]MetodaUwierzytelnienia, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaMetodUwierzytelnienia)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać metod wejścia: %w", err)
	}
	defer wiersze.Close()

	lista := []MetodaUwierzytelnienia{}
	for wiersze.Next() {
		metoda, err := odczytajMetodeUwierzytelnienia(wiersze)
		if err != nil {
			return nil, err
		}
		lista = append(lista, metoda)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt metod wejścia: %w", err)
	}
	return lista, nil
}

// MetodaPoKodzie zwraca jedną metodę wskazaną jej identyfikatorem trwałym wprost z bazy danych repozytorium.
func (r *repozytoriumUwierzytelnienia) MetodaPoKodzie(ctx context.Context,
	kod string) (MetodaUwierzytelnienia, error) {

	return r.jedna(ctx, metodaUwierzytelnieniaPoKodzie, fmt.Sprintf("%q", kod), kod)
}

// Kotwica zwraca hasło bramki, czyli jedyną metodę uwierzytelnienia, której zdjąć się w bramce nie da.
func (r *repozytoriumUwierzytelnienia) Kotwica(ctx context.Context) (MetodaUwierzytelnienia, error) {
	return r.jedna(ctx, metodaUwierzytelnieniaKotwica, "kotwica bramki")
}

// MetodaUrzadzenia zwraca metodę danego rodzaju założoną na wskazanym
// urządzeniu. Para (urządzenie, rodzaj) jest w bazie jednoznaczna.
func (r *repozytoriumUwierzytelnienia) MetodaUrzadzenia(ctx context.Context,
	rodzaj, urzadzenieKod string) (MetodaUwierzytelnienia, error) {

	return r.jedna(ctx, metodaUwierzytelnieniaUrzadzenia,
		fmt.Sprintf("%s urządzenia %q", rodzaj, urzadzenieKod), rodzaj, urzadzenieKod)
}

// ZalozMetode wstawia wiersz metody i oddaje go odczytany z bazy — dzięki temu
// warstwa wyżej dostaje dokładnie to, co zostało zapisane, a nie to, co wysłała.
func (r *repozytoriumUwierzytelnienia) ZalozMetode(ctx context.Context,
	metoda MetodaUwierzytelnienia) (MetodaUwierzytelnienia, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, wstawMetodeUwierzytelnienia)
	if err != nil {
		return MetodaUwierzytelnienia{}, err
	}
	_, err = polecenie.ExecContext(ctx, metoda.Kod, metoda.Rodzaj,
		tekstDoKolumny(metoda.Etykieta), tekstDoKolumny(metoda.UrzadzenieKod),
		tekstDoKolumny(metoda.NazwaUrzadzenia), liczbaLogiczna(metoda.Kotwica),
		metoda.OdwolanieSekretu, metoda.Utworzono)
	if czyKolizja(err) {
		// Kolizja z indeksem nie jest awarią zapisu: jest odpowiedzią, że taki wiersz już istnieje w bazie.
		return MetodaUwierzytelnienia{}, fmt.Errorf(
			"dane: metoda wejścia %q koliduje z istniejącą: %w", metoda.Kod, ErrKolizjaWiersza)
	}
	if err != nil {
		return MetodaUwierzytelnienia{}, fmt.Errorf(
			"dane: nie można założyć metody wejścia %q: %w", metoda.Kod, err)
	}
	return r.MetodaPoKodzie(ctx, metoda.Kod)
}

// ZapiszOdwolanieSekretu podmienia odwołanie do sejfu przy istniejącej metodzie.
// Brak wiersza jest błędem wołającego, nie ciszą: zmiana hasła bez kotwicy
// znaczyłaby, że rdzeń zgubił bramkę między odczytem a zapisem.
func (r *repozytoriumUwierzytelnienia) ZapiszOdwolanieSekretu(ctx context.Context,
	kod, odwolanie string) error {

	return r.zmien(ctx, zapiszOdwolanieMetody, kod, odwolanie, kod)
}

// OdnotujUzycie zapisuje w bazie danych repozytorium czas ostatniego udanego wejścia tą metodą uwierzytelnienia.
func (r *repozytoriumUwierzytelnienia) OdnotujUzycie(ctx context.Context,
	kod string, teraz int64) error {

	return r.zmien(ctx, odnotujUzycieMetody, kod, teraz, kod)
}

// UsunMetode kasuje wskazaną metodę uwierzytelnienia i mówi wołającemu, czy jej wiersz istniał w bazie danych.
func (r *repozytoriumUwierzytelnienia) UsunMetode(ctx context.Context, kod string) (bool, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, usunMetodeUwierzytelnienia)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, kod)
	if err != nil {
		return false, fmt.Errorf("dane: nie można zdjąć metody wejścia %q: %w", kod, err)
	}
	zmienione, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nieczytelny wynik zdjęcia metody wejścia %q: %w", kod, err)
	}
	return zmienione > 0, nil
}

// zmien wykonuje polecenie zmieniające jeden wiersz metody uwierzytelnienia i pilnuje, żeby wiersz naprawdę istniał.
func (r *repozytoriumUwierzytelnienia) zmien(ctx context.Context, zapytanie, kod string,
	argumenty ...any) error {

	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return err
	}
	wynik, err := polecenie.ExecContext(ctx, argumenty...)
	if err != nil {
		return fmt.Errorf("dane: nie można zmienić metody wejścia %q: %w", kod, err)
	}
	zmienione, err := wynik.RowsAffected()
	if err != nil {
		return fmt.Errorf("dane: nieczytelny wynik zmiany metody wejścia %q: %w", kod, err)
	}
	if zmienione == 0 {
		return fmt.Errorf("dane: metoda wejścia %q nie istnieje: %w", kod, ErrBrakWiersza)
	}
	return nil
}

// jedna wykonuje zapytanie zwracające najwyżej jeden wiersz metody uwierzytelnienia z bazy danych repozytorium.
func (r *repozytoriumUwierzytelnienia) jedna(ctx context.Context, zapytanie, opis string,
	argumenty ...any) (MetodaUwierzytelnienia, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return MetodaUwierzytelnienia{}, err
	}
	metoda, err := odczytajMetodeUwierzytelnienia(polecenie.QueryRowContext(ctx, argumenty...))
	if errors.Is(err, sql.ErrNoRows) {
		return MetodaUwierzytelnienia{}, fmt.Errorf(
			"dane: metoda wejścia %s nie istnieje: %w", opis, ErrBrakWiersza)
	}
	if err != nil {
		return MetodaUwierzytelnienia{}, fmt.Errorf(
			"dane: nieczytelny wiersz metody wejścia %s: %w", opis, err)
	}
	return metoda, nil
}

// odczytajMetodeUwierzytelnienia składa strukturę metody wprost z jednego wiersza wyniku zapytania SQL.
func odczytajMetodeUwierzytelnienia(wiersz skaner) (MetodaUwierzytelnienia, error) {
	var metoda MetodaUwierzytelnienia
	var etykieta, urzadzenie, nazwaUrzadzenia sql.NullString
	var ostatnio sql.NullInt64
	var kotwica int
	err := wiersz.Scan(&metoda.ID, &metoda.Kod, &metoda.Rodzaj, &etykieta,
		&urzadzenie, &nazwaUrzadzenia, &kotwica, &metoda.OdwolanieSekretu,
		&metoda.Utworzono, &ostatnio)
	if err != nil {
		return MetodaUwierzytelnienia{}, err
	}
	metoda.Etykieta = tekstZKolumny(etykieta)
	metoda.UrzadzenieKod = tekstZKolumny(urzadzenie)
	metoda.NazwaUrzadzenia = tekstZKolumny(nazwaUrzadzenia)
	metoda.Kotwica = kotwica != 0
	metoda.OstatnioUzyto = liczbaZKolumny(ostatnio)
	return metoda, nil
}
