// Odpowiedzialność pliku: katalog punktów dostępu (tabela `punkt_dostepu`) —
// struktura, kontrakt repozytorium i odczyt. Zapis leży w `dostep_punkty_zapis.go`.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"danacoconsole/shared"
)

// PunktDostepu to wiersz tabeli `punkt_dostepu` wraz z listą korzeni i słownictwem
// trybu mostu. `Kod` jest trwałym identyfikatorem punktu i odpowiada polu `id`
// struktury AccessPoint kontraktu.
type PunktDostepu struct {
	ID                     int64
	Kod                    string
	Nazwa                  string
	Opis                   string
	Rodzaj                 shared.AccessPointKind
	UrzadzenieID           *int64
	Host                   string
	Port                   int
	Uzytkownik             string
	SciezkaKlucza          string
	PolecenieStartu        string
	NazwaMostu             string
	PoswiadczenieOdwolanie *string
	TrybDomyslny           shared.AccessMode
	Stan                   shared.AccessPointStatus
	Sprawdzono             *string
	Aktywny                bool
	Kolejnosc              int
	Korzenie               []string
	// ArgumentyTrybu niosą słowo, jakim dany most nazywa tryb kontraktu przy uruchomieniu.
	ArgumentyTrybu map[shared.AccessMode]string
	Utworzono      string
	Zaktualizowano string
}

// RepozytoriumPunktowDostepu jest kontraktem katalogu punktów dostępu, określającym operacje dostępne na wykazie.
type RepozytoriumPunktowDostepu interface {
	Lista(ctx context.Context, tylkoAktywne bool) ([]PunktDostepu, error)
	Pobierz(ctx context.Context, id int64) (PunktDostepu, error)
	PoKodzie(ctx context.Context, kod string) (PunktDostepu, error)
	Dodaj(ctx context.Context, punkt PunktDostepu) (int64, error)
	Aktualizuj(ctx context.Context, punkt PunktDostepu) error
	Usun(ctx context.Context, id int64) error
	ZapiszWynikSprawdzenia(ctx context.Context, id int64,
		stan shared.AccessPointStatus, sprawdzono string) error
}

const (
	kolumnyPunktuDostepu = `id, kod, nazwa, opis, rodzaj, urzadzenie_id, host, port, uzytkownik,
	                        sciezka_klucza, polecenie_startu, nazwa_mostu, poswiadczenie_odwolanie,
	                        tryb_domyslny, stan, sprawdzono, aktywny, kolejnosc,
	                        utworzono, zaktualizowano`

	listaPunktowDostepu = `SELECT ` + kolumnyPunktuDostepu + ` FROM punkt_dostepu
	                       WHERE (? = 0 OR aktywny = 1) AND ` + WarunekKonta + `
	                       ORDER BY kolejnosc, kod`

	pobierzPunktDostepu = `SELECT ` + kolumnyPunktuDostepu + ` FROM punkt_dostepu
	                       WHERE id = ? AND ` + WarunekKonta

	// Kolumna `kod` jest unikalna w całej tabeli, nie w obrębie konta, więc kod
	// zajęty przez punkt innego konta wraca po zawężeniu jako brak wiersza.
	punktDostepuPoKodzie = `SELECT ` + kolumnyPunktuDostepu + ` FROM punkt_dostepu
	                        WHERE kod = ? AND ` + WarunekKonta

	// Nadanie wskazuje punkt identyfikatorem przyniesionym przez żądanie,
	// a punkt cudzy oddaje pusty wykaz korzeni — nie do odróżnienia od punktu
	// bez ograniczenia obszaru. Przynależność punktu rozstrzyga się osobno.
	punktDostepuKonta = `SELECT 1 FROM punkt_dostepu WHERE id = ? AND ` + WarunekKonta
)

type repozytoriumPunktowDostepu struct {
	zapytania *zapytania
	db        *sql.DB
}

// Zgodność implementacji z kontraktem sprawdzana jest przy kompilacji, a nie
// dopiero przy złożeniu zestawu repozytoriów.
var _ RepozytoriumPunktowDostepu = (*repozytoriumPunktowDostepu)(nil)

// noweRepozytoriumPunktowDostepu zakłada repozytorium katalogu punktów dostępu na przekazanym połączeniu z bazą.
func noweRepozytoriumPunktowDostepu(z *zapytania, db *sql.DB) *repozytoriumPunktowDostepu {
	return &repozytoriumPunktowDostepu{zapytania: z, db: db}
}

// Lista zwraca katalog punktów — komplet albo same czynne. Katalog pusty nie jest
// błędem: okno bez nadań pracuje dalej, tylko niczego nie widzi.
func (r *repozytoriumPunktowDostepu) Lista(ctx context.Context, tylkoAktywne bool) ([]PunktDostepu, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaPunktowDostepu)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, liczbaLogiczna(tylkoAktywne), KontoOperatora(ctx))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać katalogu punktów dostępu: %w", err)
	}
	defer wiersze.Close()

	lista := []PunktDostepu{}
	for wiersze.Next() {
		punkt, err := odczytajPunktDostepu(wiersze)
		if err != nil {
			return nil, err
		}
		lista = append(lista, punkt)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt katalogu punktów dostępu: %w", err)
	}
	return r.uzupelnijListy(ctx, lista)
}

// Pobierz zwraca punkt dostępu wskazany kluczem głównym wiersza tabeli katalogu punktów dostępu w bazie.
func (r *repozytoriumPunktowDostepu) Pobierz(ctx context.Context, id int64) (PunktDostepu, error) {
	return r.jeden(ctx, pobierzPunktDostepu, fmt.Sprintf("%d", id), id)
}

// PoKodzie zwraca punkt wskazany trwałym kodem — tym samym, który wychodzi
// kontraktem jako `AccessPoint.id`.
func (r *repozytoriumPunktowDostepu) PoKodzie(ctx context.Context, kod string) (PunktDostepu, error) {
	return r.jeden(ctx, punktDostepuPoKodzie, fmt.Sprintf("%q", kod), kod)
}

// jeden odczytuje pojedynczy punkt dostępu wraz z jego listami podrzędnymi — korzeniami i słownictwem trybu.
func (r *repozytoriumPunktowDostepu) jeden(ctx context.Context, zapytanie, opis string,
	argument any) (PunktDostepu, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return PunktDostepu{}, err
	}
	punkt, err := odczytajPunktDostepu(polecenie.QueryRowContext(ctx, argument, KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return PunktDostepu{}, fmt.Errorf("dane: punkt dostępu %s nie istnieje: %w",
			opis, ErrBrakWiersza)
	}
	if err != nil {
		return PunktDostepu{}, err
	}
	return r.uzupelnijPunkt(ctx, punkt)
}

// uzupelnijPunkt dokłada listy podrzędne punktu dostępu — korzenie oraz słownictwo trybu używane przez most.
func (r *repozytoriumPunktowDostepu) uzupelnijPunkt(ctx context.Context,
	punkt PunktDostepu) (PunktDostepu, error) {

	korzenie, err := wczytajKorzenie(ctx, r.zapytania, listaKorzeniPunktu, punkt.ID,
		"punktu dostępu", KontoOperatora(ctx))
	if err != nil {
		return PunktDostepu{}, err
	}
	argumenty, err := argumentyTrybuMostu(ctx, r.zapytania, punkt.ID)
	if err != nil {
		return PunktDostepu{}, err
	}
	punkt.Korzenie = korzenie
	punkt.ArgumentyTrybu = argumenty
	return punkt, nil
}

// uzupelnijListy dokłada listy podrzędne każdemu punktowi dostępu z przekazanego wykazu wyników zapytania.
func (r *repozytoriumPunktowDostepu) uzupelnijListy(ctx context.Context,
	lista []PunktDostepu) ([]PunktDostepu, error) {

	for i, punkt := range lista {
		uzupelniony, err := r.uzupelnijPunkt(ctx, punkt)
		if err != nil {
			return nil, err
		}
		lista[i] = uzupelniony
	}
	return lista, nil
}
