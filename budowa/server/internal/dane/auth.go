// Plik prowadzi katalog metod wejścia przez bramkę: hasło, PIN urządzenia albo klucz Windows Hello, wraz z miarą tempa
// zgadywania dróg potwierdzenia; sesje bramki leżą obok, w auth_sesje_bramki.go, a odwołanie do sejfu poświadczeń
// repozytorium traktuje jak nieprzezroczysty napis.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// MetodaUwierzytelnienia to wiersz tabeli `metoda_uwierzytelnienia`; kod jest identyfikatorem trwałym, kod urządzenia niesie napis nadany przez klienta.
type MetodaUwierzytelnienia struct {
	ID              int64
	Kod             string
	Rodzaj          string
	Etykieta        *string
	UrzadzenieKod   *string
	NazwaUrzadzenia *string
	Kotwica         bool
	// KontoId wiąże metodę z kontem. Zero znaczy wiersz zastany, sprzed
	// migracji 406, gdy konto było jedno i wskazania nie potrzebowało.
	KontoId          int64
	OdwolanieSekretu string
	Utworzono        int64
	OstatnioUzyto    *int64
}

// RepozytoriumUwierzytelnienia jest kontraktem bramki: katalog metod
// wejścia oraz sesje bramki w jednym bycie, bo obie tabele zmieniają się razem
// — zmiana hasła unieważnia sesje, a wejście odnotowuje użycie metody.
//
// Każda czynność czyta konto wołającego z kontekstu żądania i zawęża się do
// niego. Konto wędruje kontekstem, a nie argumentem, bo te same zapytania
// wołają trzy warstwy wyżej i żadna z nich konta by nie miała skąd wziąć.
type RepozytoriumUwierzytelnienia interface {
	// Metody zwraca komplet metod wejścia konta wołającego w kolejności wyświetlania.
	Metody(ctx context.Context) ([]MetodaUwierzytelnienia, error)
	// MetodaPoKodzie zwraca metodę konta wołającego wskazaną identyfikatorem
	// trwałym; brak wiersza daje ErrBrakWiersza. Identyfikator sam nie
	// wystarcza: metoda cudzego konta ma być dla wołającego nieodróżnialna od
	// nieistniejącej.
	MetodaPoKodzie(ctx context.Context, kod string) (MetodaUwierzytelnienia, error)
	// KotwicaKonta zwraca hasło wskazanego konta. Wiersze zastane, bez wskazania
	// konta, należą do konta najstarszego — tak stały przed migracją 406.
	KotwicaKonta(ctx context.Context, kontoId int64) (MetodaUwierzytelnienia, error)
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
	// UniewaznijSesjeBramkiPoza unieważnia sesje czynne konta wołającego poza wskazaną i zwraca ich liczbę.
	UniewaznijSesjeBramkiPoza(ctx context.Context, skrotZachowany string, teraz int64) (int, error)
	// UrzadzeniaKonta zwraca urządzenia konta wołającego, które weszły przez bramkę, z ostatnią chwilą wejścia.
	UrzadzeniaKonta(ctx context.Context, teraz int64) ([]UrzadzenieKonta, error)
	// UniewaznijSesjeUrzadzenia zamyka sesje czynne wskazanego urządzenia
	// w obrębie konta wołającego i zwraca ich liczbę.
	UniewaznijSesjeUrzadzenia(ctx context.Context, urzadzenie string, teraz int64) (int, error)

	/* Trzy czynności niżej dotyczą dróg potwierdzenia (tabela
	   `potwierdzenie_tozsamosci`), a stoją przy bramce, bo wszystkie trzy są
	   miarą tempa zgadywania kodu, nie treścią konta: licznik pomyłek zamyka
	   drogę, zamknięcie dróg poprzednich zostawia w obiegu jeden kod naraz,
	   a chwila wydania ostatniej drogi wyznacza odstęp między listami. */

	// OdnotujPomylkeDrogi dolicza pomyłkę drogom czynnym danego celu i zamyka te,
	// które doszły do pułapu; zwraca liczbę zamkniętych.
	OdnotujPomylkeDrogi(ctx context.Context, cel string, pulapProb int, teraz int64) (int, error)
	// ZamknijDrogiKonta zamyka drogi czynne wskazanego konta i celu; zwraca ich liczbę.
	ZamknijDrogiKonta(ctx context.Context, kontoId int64, cel string) (int, error)
	// OstatniaDrogaKonta oddaje chwilę wydania ostatniej drogi tego konta i celu;
	// zero znaczy, że takiej drogi nie było.
	OstatniaDrogaKonta(ctx context.Context, kontoId int64, cel string) (int64, error)
}

// UrzadzenieKonta to jeden wiersz wykazu urządzeń powiązanych z kontem; urządzeniem konta jest to, które kiedykolwiek weszło przez bramkę.
type UrzadzenieKonta struct {
	Kod string
	// OstatnioWidziane to chwila ostatniego wejścia w milisekundach epoki.
	OstatnioWidziane int64
	// MaToken mówi, czy urządzenie ma dziś ważny, nieunieważniony token.
	MaToken bool
}

/*
warunekKontaBramki zawęża wiersz bramki do konta, które o niego pyta. Jest
jeden na wszystkie tabele bramki, bo rozstrzygnięcie musi wypaść tak samo dla
metody wejścia, sesji i drogi potwierdzenia — inaczej granica konta trzymałaby
w jednym zapytaniu, a w sąsiednim nie.

Obie strony porównania sprowadzają brak wskazania do konta najstarszego: wiersz
bez `konto_id` powstał przed migracją 406 i nie ma jak wskazać konta wstecz,
a żądanie bez rozpoznanego konta przychodzi z połączenia przed zalogowaniem.
Instalacja jednokontowa rozstrzyga oba przypadki jednoznacznie.
*/
const warunekKontaBramki = `COALESCE(konto_id, (SELECT id FROM konto_wlasciciela ORDER BY id LIMIT 1))
	                        = COALESCE(NULLIF(?, 0), (SELECT id FROM konto_wlasciciela ORDER BY id LIMIT 1))`

const (
	kolumnyMetodyUwierzytelnienia = `id, identyfikator_zewnetrzny, rodzaj, etykieta,
	                                 urzadzenie_kod, nazwa_urzadzenia, kotwica,
	                                 COALESCE(konto_id, 0), sekret_odwolanie,
	                                 utworzono, ostatnio_uzyto`

	listaMetodUwierzytelnienia = `SELECT ` + kolumnyMetodyUwierzytelnienia +
		` FROM metoda_uwierzytelnienia WHERE ` + warunekKontaBramki +
		` ORDER BY kotwica DESC, rodzaj, utworzono, id`

	metodaUwierzytelnieniaPoKodzie = `SELECT ` + kolumnyMetodyUwierzytelnienia +
		` FROM metoda_uwierzytelnienia
		  WHERE identyfikator_zewnetrzny = ? AND ` + warunekKontaBramki

	// Hasło wskazanego konta.
	metodaUwierzytelnieniaKotwicaKonta = `SELECT ` + kolumnyMetodyUwierzytelnienia +
		` FROM metoda_uwierzytelnienia WHERE kotwica = 1 AND ` + warunekKontaBramki

	metodaUwierzytelnieniaUrzadzenia = `SELECT ` + kolumnyMetodyUwierzytelnienia +
		` FROM metoda_uwierzytelnienia WHERE rodzaj = ? AND urzadzenie_kod = ?`

	// Zero w kolumnie konta znaczyłoby konto o identyfikatorze zero, czyli
	// żadne; wiersz z zerem byłby niewidoczny dla własnego właściciela, więc
	// wskazanie puste wchodzi jako NULL, tak jak wiersz zastany.
	wstawMetodeUwierzytelnienia = `INSERT INTO metoda_uwierzytelnienia
	                               (identyfikator_zewnetrzny, rodzaj, etykieta, urzadzenie_kod,
	                                nazwa_urzadzenia, kotwica, konto_id, sekret_odwolanie, utworzono)
	                               VALUES (?, ?, ?, ?, ?, ?, NULLIF(?, 0), ?, ?)`

	zapiszOdwolanieMetody = `UPDATE metoda_uwierzytelnienia SET sekret_odwolanie = ?
	                         WHERE identyfikator_zewnetrzny = ?`

	odnotujUzycieMetody = `UPDATE metoda_uwierzytelnienia SET ostatnio_uzyto = ?
	                       WHERE identyfikator_zewnetrzny = ?`

	usunMetodeUwierzytelnienia = `DELETE FROM metoda_uwierzytelnienia
	                              WHERE identyfikator_zewnetrzny = ?`

	// Pomyłka nie trafia w żaden wiersz — kod nietrafiony nie nazywa drogi,
	// której dotyczył. Liczy się więc każdej drodze czynnej tego celu: zgadujący
	// mierzy w którąś z nich, a pułap zamyka je, zanim dojdzie do wartości.
	doliczPomylkeDrogom = `UPDATE potwierdzenie_tozsamosci SET proby = proby + 1
	                       WHERE cel = ? AND uzyte = 0 AND wygasa > ?`

	zamknijDrogiPrzepalone = `UPDATE potwierdzenie_tozsamosci SET uzyte = 1
	                          WHERE cel = ? AND uzyte = 0 AND proby >= ?`

	zamknijDrogiKonta = `UPDATE potwierdzenie_tozsamosci SET uzyte = 1
	                     WHERE cel = ? AND uzyte = 0 AND ` + warunekKontaBramki

	ostatniaDrogaKonta = `SELECT COALESCE(MAX(utworzono), 0) FROM potwierdzenie_tozsamosci
	                      WHERE cel = ? AND ` + warunekKontaBramki
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

// Metody zwraca komplet metod wejścia konta wołającego. Kotwica idzie pierwsza,
// bo sekcja Uwierzytelnianie Okna Ustawień pokazuje hasło na czele, a metody
// szybkiego wejścia pod nim. Katalog pusty nie jest błędem — znaczy bramkę
// nieustawioną albo konto bez własnej metody.
func (r *repozytoriumUwierzytelnienia) Metody(ctx context.Context) ([]MetodaUwierzytelnienia, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaMetodUwierzytelnienia)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, KontoOperatora(ctx))
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

// MetodaPoKodzie zwraca jedną metodę konta wołającego, wskazaną jej
// identyfikatorem trwałym. Identyfikator cudzej metody daje ErrBrakWiersza:
// sam identyfikator nie jest poświadczeniem, a bez tego warunku wystarczał, by
// zdjąć metodę wejścia obcego konta.
func (r *repozytoriumUwierzytelnienia) MetodaPoKodzie(ctx context.Context,
	kod string) (MetodaUwierzytelnienia, error) {

	return r.jedna(ctx, metodaUwierzytelnieniaPoKodzie, fmt.Sprintf("%q", kod), kod, KontoOperatora(ctx))
}

// KotwicaKonta zwraca hasło wskazanego konta. Warunek obejmuje też wiersze
// zastane (`konto_id IS NULL`) należące do konta najstarszego — inaczej
// instalacja sprzed migracji 406 przestałaby wpuszczać własnego Operatora.
func (r *repozytoriumUwierzytelnienia) KotwicaKonta(ctx context.Context,
	kontoId int64) (MetodaUwierzytelnienia, error) {

	return r.jedna(ctx, metodaUwierzytelnieniaKotwicaKonta, "kotwica konta", kontoId)
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
		metoda.KontoId, metoda.OdwolanieSekretu, metoda.Utworzono)
	if czyKolizja(err) {
		// Kolizja z indeksem nie jest awarią zapisu: jest odpowiedzią, że taki wiersz już istnieje w bazie.
		return MetodaUwierzytelnienia{}, fmt.Errorf(
			"dane: metoda wejścia %q koliduje z istniejącą: %w", metoda.Kod, ErrKolizjaWiersza)
	}
	if err != nil {
		return MetodaUwierzytelnienia{}, fmt.Errorf(
			"dane: nie można założyć metody wejścia %q: %w", metoda.Kod, err)
	}
	// Odczyt idzie kontem wiersza, nie kontem żądania: rejestracja zakłada
	// hasło konta, którego jeszcze nikt nie rozpoznał w kontekście, a wiersz ma
	// wrócić do wołającego mimo to.
	return r.MetodaPoKodzie(ZKontemOperatora(ctx, metoda.KontoId), metoda.Kod)
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

// OdnotujPomylkeDrogi dolicza pomyłkę drogom czynnym danego celu i zamyka te,
// które doszły do pułapu prób. Dwa polecenia, bo SQLite nie zna aktualizacji
// warunkowanej własnym wynikiem, a kolejność jest wiążąca: najpierw przyrost,
// potem zamknięcie, inaczej pułap zamykałby drogę o jedną pomyłkę za późno.
func (r *repozytoriumUwierzytelnienia) OdnotujPomylkeDrogi(ctx context.Context,
	cel string, pulapProb int, teraz int64) (int, error) {

	dolicz, err := r.zapytania.przygotuj(ctx, doliczPomylkeDrogom)
	if err != nil {
		return 0, err
	}
	if _, err := dolicz.ExecContext(ctx, cel, teraz); err != nil {
		return 0, fmt.Errorf("dane: nie można doliczyć pomyłki drogom celu %q: %w", cel, err)
	}
	zamknij, err := r.zapytania.przygotuj(ctx, zamknijDrogiPrzepalone)
	if err != nil {
		return 0, err
	}
	wynik, err := zamknij.ExecContext(ctx, cel, pulapProb)
	if err != nil {
		return 0, fmt.Errorf("dane: nie można zamknąć dróg celu %q po pułapie prób: %w", cel, err)
	}
	zamkniete, err := wynik.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("dane: nieczytelny wynik zamknięcia dróg celu %q: %w", cel, err)
	}
	return int(zamkniete), nil
}

// ZamknijDrogiKonta zamyka drogi czynne wskazanego konta i celu. Woła się przy
// wydaniu drogi nowej: w skrzynce zostaje wtedy jeden kod ważny, a nie tyle
// kodów, ile razy Operator poprosił o list.
func (r *repozytoriumUwierzytelnienia) ZamknijDrogiKonta(ctx context.Context,
	kontoId int64, cel string) (int, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, zamknijDrogiKonta)
	if err != nil {
		return 0, err
	}
	wynik, err := polecenie.ExecContext(ctx, cel, kontoId)
	if err != nil {
		return 0, fmt.Errorf("dane: nie można zamknąć dróg konta: %w", err)
	}
	zamkniete, err := wynik.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("dane: nieczytelny wynik zamknięcia dróg konta: %w", err)
	}
	return int(zamkniete), nil
}

// OstatniaDrogaKonta oddaje chwilę wydania ostatniej drogi tego konta i celu.
// Zero znaczy brak takiej drogi, nie brak odpowiedzi — wiersza po prostu nie ma.
func (r *repozytoriumUwierzytelnienia) OstatniaDrogaKonta(ctx context.Context,
	kontoId int64, cel string) (int64, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, ostatniaDrogaKonta)
	if err != nil {
		return 0, err
	}
	var utworzono int64
	if err := polecenie.QueryRowContext(ctx, cel, kontoId).Scan(&utworzono); err != nil {
		return 0, fmt.Errorf("dane: nie można odczytać chwili ostatniej drogi konta: %w", err)
	}
	return utworzono, nil
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
		&urzadzenie, &nazwaUrzadzenia, &kotwica, &metoda.KontoId,
		&metoda.OdwolanieSekretu, &metoda.Utworzono, &ostatnio)
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
