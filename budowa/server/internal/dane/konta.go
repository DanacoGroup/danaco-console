// Odpowiedzialność pliku: odczyt katalogu kont modeli i kont programów code CLI
// (tabela `konto`). Katalog jest sterowany danymi — nowe konto to nowy wiersz, nie nowy typ w kodzie.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"danacoconsole/shared"
)

// Konto to wiersz katalogu kont, niosący pełny profil uwierzytelnienia modelu albo programu CLI rdzenia.
type Konto struct {
	ID                      int64
	Nazwa                   string
	Rodzaj                  shared.AccountKind
	Dostawca                string
	IdentyfikatorZewnetrzny *string
	ModelDomyslny           *string
	AdresBazowy             *string
	KatalogKonfiguracji     *string
	// MaPoswiadczenie mówi tylko, że odwołanie jest zapisane; treści odwołania odczyt katalogu nie niesie.
	MaPoswiadczenie bool
	Stan            StanKonta
	// WyczerpaneDo jest chwilą odnowienia limitu w zapisie ISO 8601; puste
	// znaczy „nieznana".
	WyczerpaneDo   *string
	Domyslne       bool
	Aktywne        bool
	Kolejnosc      int
	Utworzono      string
	Zaktualizowano string
}

// FiltrKont zawęża wykaz kont wynikowych; pusty filtr znaczy „komplet katalogu" bez żadnego zawężenia.
type FiltrKont struct {
	// Rodzaj ogranicza wykaz do jednego rodzaju kont; nil znaczy wszystkie.
	Rodzaj *shared.AccountKind
	// TylkoAktywne pomija konta wyłączone przez Operatora.
	TylkoAktywne bool
}

// RepozytoriumKont jest kontraktem katalogu kont, określającym operacje dostępne na całym wykazie kont.
type RepozytoriumKont interface {
	Dodaj(ctx context.Context, konto Konto, odwolaniePoswiadczenia *string) (int64, error)
	Aktualizuj(ctx context.Context, konto Konto) error
	UstawPoswiadczenie(ctx context.Context, id int64, odwolanie *string) error
	OdwolaniePoswiadczenia(ctx context.Context, id int64) (string, error)
	Usun(ctx context.Context, id int64) ([]int64, error)
	Lista(ctx context.Context, filtr FiltrKont) ([]Konto, error)
	Pobierz(ctx context.Context, id int64) (Konto, error)
	PobierzPoNazwie(ctx context.Context, nazwa string) (Konto, error)
	UstawDomyslne(ctx context.Context, id int64) (*int64, error)
	Domyslne(ctx context.Context, rodzaj shared.AccountKind) (Konto, error)
	KontaRotacji(ctx context.Context, rodzaj shared.AccountKind) ([]Konto, error)
	OznaczStan(ctx context.Context, id int64, stan StanKonta, doChwili *string) error
}

const (
	// znacznikPoswiadczenia zamienia odwołanie na samo „jest / nie ma" już
	// w zapytaniu — kolumna z odwołaniem nie opuszcza bazy.
	znacznikPoswiadczenia = `CASE WHEN poswiadczenie_odwolanie IS NULL
	                              OR poswiadczenie_odwolanie = '' THEN 0 ELSE 1 END`

	kolumnyKonta = `id, nazwa, rodzaj, dostawca, identyfikator_zewnetrzny, model_domyslny,
	                adres_bazowy, katalog_konfiguracji, ` + znacznikPoswiadczenia + `,
	                stan, wyczerpane_do, domyslne, aktywne, kolejnosc, utworzono, zaktualizowano`

	listaKont = `SELECT ` + kolumnyKonta + ` FROM konto
	             WHERE (? = '' OR rodzaj = ?) AND (? = 0 OR aktywne = 1)
	               AND ` + WarunekKonta + `
	             ORDER BY rodzaj, kolejnosc, id`

	pobierzKonto = `SELECT ` + kolumnyKonta + ` FROM konto
	                WHERE id = ? AND ` + WarunekKonta

	pobierzKontoPoNazwie = `SELECT ` + kolumnyKonta + ` FROM konto
	                        WHERE nazwa = ? AND ` + WarunekKonta

	pobierzKontoDomyslne = `SELECT ` + kolumnyKonta + ` FROM konto
	                        WHERE rodzaj = ? AND domyslne = 1 AND ` + WarunekKonta
)

type repozytoriumKont struct {
	zapytania *zapytania
	db        *sql.DB
}

// noweRepozytoriumKont zakłada repozytorium katalogu kont. Uchwyt bazy jest
// potrzebny obok pamięci zapytań, bo wskazanie konta domyślnego, usunięcie konta
// i nadanie kolejności obejmują więcej niż jedno polecenie i muszą być jedną
// transakcją.
func noweRepozytoriumKont(z *zapytania, db *sql.DB) *repozytoriumKont {
	return &repozytoriumKont{zapytania: z, db: db}
}

// Lista zwraca wykaz kont z katalogu, zawężony przekazanym filtrem wyszukiwania kont modeli i programów.
func (r *repozytoriumKont) Lista(ctx context.Context, filtr FiltrKont) ([]Konto, error) {
	rodzaj := ""
	if filtr.Rodzaj != nil {
		wartosc, err := rodzajKontaNaBaze(*filtr.Rodzaj)
		if err != nil {
			return nil, err
		}
		rodzaj = wartosc
	}
	polecenie, err := r.zapytania.przygotuj(ctx, listaKont)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, rodzaj, rodzaj,
		liczbaLogiczna(filtr.TylkoAktywne), KontoOperatora(ctx))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać katalogu kont: %w", err)
	}
	defer wiersze.Close()
	return zbierzKonta(wiersze)
}

// Pobierz zwraca jedno konto z katalogu kont wskazane kluczem głównym wiersza tabeli kont w bazie danych.
func (r *repozytoriumKont) Pobierz(ctx context.Context, id int64) (Konto, error) {
	return r.jednoKonto(ctx, pobierzKonto, fmt.Sprintf("konto %d", id), id, KontoOperatora(ctx))
}

// PobierzPoNazwie zwraca konto po jego nazwie — nazwa jest w tabeli unikalna
// i to ona jest kodem konta w profilach rotacji.
func (r *repozytoriumKont) PobierzPoNazwie(ctx context.Context, nazwa string) (Konto, error) {
	return r.jednoKonto(ctx, pobierzKontoPoNazwie, fmt.Sprintf("konto %q", nazwa), nazwa,
		KontoOperatora(ctx))
}

// Domyslne zwraca konto domyślne wskazanego rodzaju. Brak wskazania jest stanem
// normalnym, nie awarią — sygnalizuje go ErrBrakWiersza.
func (r *repozytoriumKont) Domyslne(ctx context.Context, rodzaj shared.AccountKind) (Konto, error) {
	wartosc, err := rodzajKontaNaBaze(rodzaj)
	if err != nil {
		return Konto{}, err
	}
	return r.jednoKonto(ctx, pobierzKontoDomyslne,
		fmt.Sprintf("konto domyślne rodzaju %q", wartosc), wartosc, KontoOperatora(ctx))
}

// jednoKonto wykonuje odczyt jednego wiersza i przekłada brak wiersza wyniku na zgłoszony błąd braku wiersza.
func (r *repozytoriumKont) jednoKonto(ctx context.Context, zapytanie, opis string,
	argumenty ...any) (Konto, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return Konto{}, err
	}
	konto, err := odczytajKonto(polecenie.QueryRowContext(ctx, argumenty...))
	if errors.Is(err, sql.ErrNoRows) {
		return Konto{}, fmt.Errorf("%w: %s", ErrBrakWiersza, opis)
	}
	return konto, err
}
