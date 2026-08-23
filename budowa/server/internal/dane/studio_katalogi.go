// Odpowiedzialność pliku: byty modułu Studio, które NIE należą do jednego
// dokumentu — operacje własne Tools Panel, łańcuchy operacji, profile wydania,
// szablony dokumentów, gałęzie oraz odwołania do wersji.
//
// Pierwsze cztery mieszkają w zasięgu konfiguracji (globalny, środowisko,
// projekt, sesja), a nie przy dokumencie: Operator zapisuje własny prompt raz
// i sięga po niego w każdym dokumencie. Wiązanie ich z dokumentem kazałoby
// przepisywać je przy każdym nowym pliku.
//
// Gałąź i odwołanie do wersji należą do dokumentu, ale stoją tutaj, bo są
// bytami rejestru — mają własny cykl życia i własny identyfikator zewnętrzny,
// inaczej niż zmiana śledzona, która bez dokumentu nie znaczy nic.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// WpisKatalogowyStudia to wiersz jednej z trzech tabel katalogowych modułu:
// operacji własnej, łańcucha operacji albo profilu wydania. Trzy tabele mają tę
// samą postać — kod, nazwa, ładunek, zasięg — więc jedna struktura opisuje
// wszystkie trzy, a rozróżnia je tabela, z której wiersz przyszedł.
type WpisKatalogowyStudia struct {
	ID    int64
	Kod   string
	Nazwa string
	// Dodatkowa nazwa kolumny drugiej: kategoria operacji albo format profilu.
	Rodzaj string
	// Ładunek: prompt operacji, kroki łańcucha albo ustawienia profilu.
	Ladunek   string
	Zasieg    string
	ZasiegID  *string
	Utworzono string
}

// SzablonStudia to wiersz tabeli `szablon_studio`.
type SzablonStudia struct {
	ID        int64
	Kod       string
	Nazwa     string
	Opis      *string
	Format    string
	Tresc     string
	PolaJSON  *string
	Fabryczny bool
	Utworzono string
}

// GalazStudia to wiersz tabeli `galaz_studio`.
type GalazStudia struct {
	ID               int64
	Kod              string
	DokumentKod      string
	Nazwa            string
	WersjaStartowaID string
	WersjaBiezacaID  *string
	Scalona          bool
	Utworzono        string
}

// OdwolanieWersji to wiersz tabeli `odwolanie_wersji_studio`.
type OdwolanieWersji struct {
	ID        int64
	Kod       string
	WersjaKod string
	Zasieg    string
	Wygasa    *string
	Utworzono string
}

// opisTabeliKatalogu wskazuje tabelę i nazwy jej dwóch kolumn zmiennych. Dzięki
// niemu trzy tabele obsługuje jeden komplet metod, a różnice nazw kolumn stoją
// w jednym miejscu, nie w sześciu zapytaniach przepisanych z drobną zmianą.
type opisTabeliKatalogu struct {
	tabela         string
	rodzaj         string
	ladunek        string
	domyslnyZasieg string
}

var (
	tabelaOperacjiStudia = opisTabeliKatalogu{
		tabela: "operacja_studio", rodzaj: "kategoria", ladunek: "prompt", domyslnyZasieg: "global",
	}
	tabelaLancuchowStudia = opisTabeliKatalogu{
		tabela: "lancuch_studio", rodzaj: "", ladunek: "kroki_json", domyslnyZasieg: "project",
	}
	tabelaProfiliStudia = opisTabeliKatalogu{
		tabela: "profil_wydania_studio", rodzaj: "format", ladunek: "ustawienia_json", domyslnyZasieg: "global",
	}
)

// kolumnyKatalogu składa listę kolumn odczytu. Tabela bez kolumny rodzaju
// oddaje w jej miejsce pusty napis, żeby odczyt miał zawsze tyle samo pól.
func (o opisTabeliKatalogu) kolumnyOdczytu() string {
	rodzaj := "''"
	if o.rodzaj != "" {
		rodzaj = o.rodzaj
	}
	ladunek := o.ladunek
	if o.tabela == "profil_wydania_studio" {
		// Ustawienia profilu dopuszczają NULL, a odczyt oczekuje napisu.
		ladunek = "COALESCE(" + o.ladunek + ", '')"
	}
	return "id, identyfikator_zewnetrzny, nazwa, " + rodzaj + ", " + ladunek +
		", zasieg, zasieg_id, utworzono"
}

// ZapiszWpisKatalogowy zakłada albo nadpisuje wiersz jednej z trzech tabel
// katalogowych modułu.
func (r *repozytoriumStudia) ZapiszWpisKatalogowy(ctx context.Context,
	tabela opisTabeliKatalogu, wpis WpisKatalogowyStudia) (WpisKatalogowyStudia, error) {

	if wpis.Kod == "" {
		return WpisKatalogowyStudia{}, fmt.Errorf("dane: wpis katalogu studio bez identyfikatora")
	}
	if wpis.Zasieg == "" {
		wpis.Zasieg = tabela.domyslnyZasieg
	}

	kolumny := "identyfikator_zewnetrzny, nazwa, " + tabela.ladunek + ", zasieg, zasieg_id"
	miejsca := "?, ?, ?, ?, ?"
	nadpisanie := "nazwa = excluded.nazwa, " + tabela.ladunek + " = excluded." + tabela.ladunek +
		", zasieg = excluded.zasieg, zasieg_id = excluded.zasieg_id"
	argumenty := []any{wpis.Kod, wpis.Nazwa, wpis.Ladunek, wpis.Zasieg, tekstDoKolumny(wpis.ZasiegID)}
	if tabela.rodzaj != "" {
		kolumny = "identyfikator_zewnetrzny, nazwa, " + tabela.rodzaj + ", " + tabela.ladunek +
			", zasieg, zasieg_id"
		miejsca = "?, ?, ?, ?, ?, ?"
		nadpisanie = "nazwa = excluded.nazwa, " + tabela.rodzaj + " = excluded." + tabela.rodzaj +
			", " + tabela.ladunek + " = excluded." + tabela.ladunek +
			", zasieg = excluded.zasieg, zasieg_id = excluded.zasieg_id"
		argumenty = []any{wpis.Kod, wpis.Nazwa, wpis.Rodzaj, wpis.Ladunek, wpis.Zasieg,
			tekstDoKolumny(wpis.ZasiegID)}
	}

	zapytanie := "INSERT INTO " + tabela.tabela + " (" + kolumny + ") VALUES (" + miejsca + ")" +
		" ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET " + nadpisanie
	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return WpisKatalogowyStudia{}, err
	}
	if _, err := polecenie.ExecContext(ctx, argumenty...); err != nil {
		return WpisKatalogowyStudia{}, fmt.Errorf("dane: nie można zapisać wpisu %q w %s: %w",
			wpis.Kod, tabela.tabela, err)
	}
	return r.WpisKatalogowy(ctx, tabela, wpis.Kod)
}

// WpisKatalogowy zwraca jeden wiersz tabeli katalogowej.
func (r *repozytoriumStudia) WpisKatalogowy(ctx context.Context,
	tabela opisTabeliKatalogu, kod string) (WpisKatalogowyStudia, error) {

	zapytanie := "SELECT " + tabela.kolumnyOdczytu() + " FROM " + tabela.tabela +
		" WHERE identyfikator_zewnetrzny = ?"
	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return WpisKatalogowyStudia{}, err
	}
	wpis, err := odczytajWpisKatalogowy(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return WpisKatalogowyStudia{}, ErrBrakWiersza
	}
	if err != nil {
		return WpisKatalogowyStudia{}, fmt.Errorf("dane: nieczytelny wpis %q z %s: %w",
			kod, tabela.tabela, err)
	}
	return wpis, nil
}

// WpisyKatalogowe zwraca wiersze tabeli katalogowej dla wskazanego zasięgu.
func (r *repozytoriumStudia) WpisyKatalogowe(ctx context.Context,
	tabela opisTabeliKatalogu, zasieg string, zasiegID *string) ([]WpisKatalogowyStudia, error) {

	if zasieg == "" {
		zasieg = tabela.domyslnyZasieg
	}
	// Zasięg pusty w kolumnie i zasięg podany to dwa różne wiersze, więc
	// porównanie identyfikatora zasięgu musi znosić NULL po obu stronach.
	zapytanie := "SELECT " + tabela.kolumnyOdczytu() + " FROM " + tabela.tabela +
		" WHERE zasieg = ? AND (zasieg_id IS ? OR zasieg_id = ?) ORDER BY nazwa, id"
	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, zasieg, tekstDoKolumny(zasiegID), tekstDoKolumny(zasiegID))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać wykazu z %s: %w", tabela.tabela, err)
	}
	defer wiersze.Close()

	lista := []WpisKatalogowyStudia{}
	for wiersze.Next() {
		wpis, err := odczytajWpisKatalogowy(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz z %s: %w", tabela.tabela, err)
		}
		lista = append(lista, wpis)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt z %s: %w", tabela.tabela, err)
	}
	return lista, nil
}

// UsunWpisKatalogowy usuwa wiersz i mówi, czy było co usuwać.
func (r *repozytoriumStudia) UsunWpisKatalogowy(ctx context.Context,
	tabela opisTabeliKatalogu, kod string) (bool, error) {

	zapytanie := "DELETE FROM " + tabela.tabela + " WHERE identyfikator_zewnetrzny = ?"
	polecenie, err := r.zapytania.przygotuj(ctx, zapytanie)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, kod)
	if err != nil {
		return false, fmt.Errorf("dane: nie można usunąć wpisu %q z %s: %w", kod, tabela.tabela, err)
	}
	usuniete, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nieznany skutek usunięcia wpisu %q: %w", kod, err)
	}
	return usuniete > 0, nil
}

func odczytajWpisKatalogowy(wiersz skaner) (WpisKatalogowyStudia, error) {
	var wpis WpisKatalogowyStudia
	var zasiegID sql.NullString
	err := wiersz.Scan(&wpis.ID, &wpis.Kod, &wpis.Nazwa, &wpis.Rodzaj, &wpis.Ladunek,
		&wpis.Zasieg, &zasiegID, &wpis.Utworzono)
	if err != nil {
		return WpisKatalogowyStudia{}, err
	}
	wpis.ZasiegID = tekstZKolumny(zasiegID)
	return wpis, nil
}

// ── Szablony dokumentów ─────────────────────────────────────────────────────

const (
	kolumnySzablonuStudia = `id, identyfikator_zewnetrzny, nazwa, opis, format, tresc,
	                         pola_json, fabryczny, utworzono`

	zapiszSzablonStudia = `INSERT INTO szablon_studio
	                       (identyfikator_zewnetrzny, nazwa, opis, format, tresc, pola_json, fabryczny)
	                       VALUES (?, ?, ?, ?, ?, ?, ?)
	                       ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                           nazwa = excluded.nazwa, opis = excluded.opis,
	                           format = excluded.format, tresc = excluded.tresc,
	                           pola_json = excluded.pola_json`

	pobierzSzablonStudia = `SELECT ` + kolumnySzablonuStudia + ` FROM szablon_studio
	                        WHERE identyfikator_zewnetrzny = ?`

	listaSzablonowStudia = `SELECT ` + kolumnySzablonuStudia + ` FROM szablon_studio
	                        ORDER BY fabryczny DESC, nazwa`
)

// ZapiszSzablon zakłada szablon dokumentu albo nadpisuje zastany.
func (r *repozytoriumStudia) ZapiszSzablon(ctx context.Context,
	szablon SzablonStudia) (SzablonStudia, error) {

	if szablon.Kod == "" {
		return SzablonStudia{}, fmt.Errorf("dane: szablon studio bez identyfikatora")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszSzablonStudia)
	if err != nil {
		return SzablonStudia{}, err
	}
	_, err = polecenie.ExecContext(ctx, szablon.Kod, szablon.Nazwa, tekstDoKolumny(szablon.Opis),
		szablon.Format, szablon.Tresc, tekstDoKolumny(szablon.PolaJSON),
		liczbaLogiczna(szablon.Fabryczny))
	if err != nil {
		return SzablonStudia{}, fmt.Errorf("dane: nie można zapisać szablonu %q: %w", szablon.Kod, err)
	}
	return r.Szablon(ctx, szablon.Kod)
}

// Szablon zwraca szablon o wskazanym kodzie.
func (r *repozytoriumStudia) Szablon(ctx context.Context, kod string) (SzablonStudia, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzSzablonStudia)
	if err != nil {
		return SzablonStudia{}, err
	}
	szablon, err := odczytajSzablonStudia(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return SzablonStudia{}, ErrBrakWiersza
	}
	if err != nil {
		return SzablonStudia{}, fmt.Errorf("dane: nieczytelny szablon %q: %w", kod, err)
	}
	return szablon, nil
}

// Szablony zwraca komplet szablonów, fabryczne na początku.
func (r *repozytoriumStudia) Szablony(ctx context.Context) ([]SzablonStudia, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaSzablonowStudia)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać szablonów studio: %w", err)
	}
	defer wiersze.Close()

	lista := []SzablonStudia{}
	for wiersze.Next() {
		szablon, err := odczytajSzablonStudia(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz szablonu studio: %w", err)
		}
		lista = append(lista, szablon)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt szablonów studio: %w", err)
	}
	return lista, nil
}

func odczytajSzablonStudia(wiersz skaner) (SzablonStudia, error) {
	var szablon SzablonStudia
	var opis, pola sql.NullString
	var fabryczny int
	err := wiersz.Scan(&szablon.ID, &szablon.Kod, &szablon.Nazwa, &opis, &szablon.Format,
		&szablon.Tresc, &pola, &fabryczny, &szablon.Utworzono)
	if err != nil {
		return SzablonStudia{}, err
	}
	szablon.Opis = tekstZKolumny(opis)
	szablon.PolaJSON = tekstZKolumny(pola)
	szablon.Fabryczny = fabryczny == 1
	return szablon, nil
}

// ── Gałęzie dokumentu ───────────────────────────────────────────────────────

const (
	kolumnyGaleziStudia = `g.id, g.identyfikator_zewnetrzny, d.identyfikator_zewnetrzny, g.nazwa,
	                       g.wersja_startowa_id, g.wersja_biezaca_id, g.scalona, g.utworzono`

	zapiszGalazStudia = `INSERT INTO galaz_studio
	                     (identyfikator_zewnetrzny, dokument_id, nazwa, wersja_startowa_id,
	                      wersja_biezaca_id, scalona)
	                     VALUES (?, ?, ?, ?, ?, ?)
	                     ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                         nazwa = excluded.nazwa,
	                         wersja_biezaca_id = excluded.wersja_biezaca_id,
	                         scalona = excluded.scalona`

	pobierzGalazStudia = `SELECT ` + kolumnyGaleziStudia + ` FROM galaz_studio g
	                      JOIN dokument_studio d ON d.id = g.dokument_id
	                      WHERE g.identyfikator_zewnetrzny = ?`

	listaGaleziStudia = `SELECT ` + kolumnyGaleziStudia + ` FROM galaz_studio g
	                     JOIN dokument_studio d ON d.id = g.dokument_id
	                     WHERE g.dokument_id = ? ORDER BY g.id`
)

// ZapiszGalaz zakłada gałąź dokumentu albo nadpisuje jej stan.
func (r *repozytoriumStudia) ZapiszGalaz(ctx context.Context,
	dokumentID int64, galaz GalazStudia) (GalazStudia, error) {

	if galaz.Kod == "" {
		return GalazStudia{}, fmt.Errorf("dane: gałąź studio bez identyfikatora")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszGalazStudia)
	if err != nil {
		return GalazStudia{}, err
	}
	_, err = polecenie.ExecContext(ctx, galaz.Kod, dokumentID, galaz.Nazwa, galaz.WersjaStartowaID,
		tekstDoKolumny(galaz.WersjaBiezacaID), liczbaLogiczna(galaz.Scalona))
	if err != nil {
		return GalazStudia{}, fmt.Errorf("dane: nie można zapisać gałęzi %q: %w", galaz.Kod, err)
	}
	return r.Galaz(ctx, galaz.Kod)
}

// Galaz zwraca gałąź o wskazanym kodzie.
func (r *repozytoriumStudia) Galaz(ctx context.Context, kod string) (GalazStudia, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzGalazStudia)
	if err != nil {
		return GalazStudia{}, err
	}
	galaz, err := odczytajGalazStudia(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return GalazStudia{}, ErrBrakWiersza
	}
	if err != nil {
		return GalazStudia{}, fmt.Errorf("dane: nieczytelna gałąź %q: %w", kod, err)
	}
	return galaz, nil
}

// Galezie zwraca gałęzie dokumentu w kolejności założenia.
func (r *repozytoriumStudia) Galezie(ctx context.Context, dokumentID int64) ([]GalazStudia, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaGaleziStudia)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, dokumentID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać gałęzi studio: %w", err)
	}
	defer wiersze.Close()

	lista := []GalazStudia{}
	for wiersze.Next() {
		galaz, err := odczytajGalazStudia(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz gałęzi studio: %w", err)
		}
		lista = append(lista, galaz)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt gałęzi studio: %w", err)
	}
	return lista, nil
}

func odczytajGalazStudia(wiersz skaner) (GalazStudia, error) {
	var galaz GalazStudia
	var biezaca sql.NullString
	var scalona int
	err := wiersz.Scan(&galaz.ID, &galaz.Kod, &galaz.DokumentKod, &galaz.Nazwa,
		&galaz.WersjaStartowaID, &biezaca, &scalona, &galaz.Utworzono)
	if err != nil {
		return GalazStudia{}, err
	}
	galaz.WersjaBiezacaID = tekstZKolumny(biezaca)
	galaz.Scalona = scalona == 1
	return galaz, nil
}

// ── Odwołania do wersji ─────────────────────────────────────────────────────

// ZapiszOdwolanieWersji zakłada odwołanie do wersji dokumentu.
func (r *repozytoriumStudia) ZapiszOdwolanieWersji(ctx context.Context,
	odwolanie OdwolanieWersji) (OdwolanieWersji, error) {

	if odwolanie.Kod == "" {
		return OdwolanieWersji{}, fmt.Errorf("dane: odwołanie do wersji bez identyfikatora")
	}
	if odwolanie.Zasieg == "" {
		odwolanie.Zasieg = "session"
	}
	polecenie, err := r.zapytania.przygotuj(ctx,
		`INSERT INTO odwolanie_wersji_studio (identyfikator_zewnetrzny, wersja_id, zasieg, wygasa)
		 VALUES (?, ?, ?, ?)`)
	if err != nil {
		return OdwolanieWersji{}, err
	}
	_, err = polecenie.ExecContext(ctx, odwolanie.Kod, odwolanie.WersjaKod, odwolanie.Zasieg,
		tekstDoKolumny(odwolanie.Wygasa))
	if err != nil {
		return OdwolanieWersji{}, fmt.Errorf("dane: nie można zapisać odwołania %q: %w",
			odwolanie.Kod, err)
	}
	return odwolanie, nil
}

// ── Nazwane wejścia do trzech tabel katalogowych ────────────────────────────
//
// Kontrakt repozytorium mówi nazwami dziedziny („zapisz operację"), a nie
// nazwą tabeli podanej parametrem: wołający nie ma rozstrzygać, w której
// tabeli byt mieszka. Opis tabeli zostaje szczegółem tego pliku.

// ZapiszOperacje zapisuje operację własną Tools Panel.
func (r *repozytoriumStudia) ZapiszOperacje(ctx context.Context,
	wpis WpisKatalogowyStudia) (WpisKatalogowyStudia, error) {
	return r.ZapiszWpisKatalogowy(ctx, tabelaOperacjiStudia, wpis)
}

// Operacje zwraca operacje własne zasięgu.
func (r *repozytoriumStudia) Operacje(ctx context.Context,
	zasieg string, zasiegID *string) ([]WpisKatalogowyStudia, error) {
	return r.WpisyKatalogowe(ctx, tabelaOperacjiStudia, zasieg, zasiegID)
}

// UsunOperacje usuwa operację własną.
func (r *repozytoriumStudia) UsunOperacje(ctx context.Context, kod string) (bool, error) {
	return r.UsunWpisKatalogowy(ctx, tabelaOperacjiStudia, kod)
}

// ZapiszLancuch zapisuje łańcuch operacji.
func (r *repozytoriumStudia) ZapiszLancuch(ctx context.Context,
	wpis WpisKatalogowyStudia) (WpisKatalogowyStudia, error) {
	return r.ZapiszWpisKatalogowy(ctx, tabelaLancuchowStudia, wpis)
}

// Lancuch zwraca jeden łańcuch operacji.
func (r *repozytoriumStudia) Lancuch(ctx context.Context, kod string) (WpisKatalogowyStudia, error) {
	return r.WpisKatalogowy(ctx, tabelaLancuchowStudia, kod)
}

// Lancuchy zwraca łańcuchy operacji zasięgu.
func (r *repozytoriumStudia) Lancuchy(ctx context.Context,
	zasieg string, zasiegID *string) ([]WpisKatalogowyStudia, error) {
	return r.WpisyKatalogowe(ctx, tabelaLancuchowStudia, zasieg, zasiegID)
}

// ZapiszProfilWydania zapisuje profil wydania dokumentu.
func (r *repozytoriumStudia) ZapiszProfilWydania(ctx context.Context,
	wpis WpisKatalogowyStudia) (WpisKatalogowyStudia, error) {
	return r.ZapiszWpisKatalogowy(ctx, tabelaProfiliStudia, wpis)
}

// ProfilWydania zwraca jeden profil wydania.
func (r *repozytoriumStudia) ProfilWydania(ctx context.Context, kod string) (WpisKatalogowyStudia, error) {
	return r.WpisKatalogowy(ctx, tabelaProfiliStudia, kod)
}

// ProfileWydania zwraca profile wydania zasięgu.
func (r *repozytoriumStudia) ProfileWydania(ctx context.Context,
	zasieg string, zasiegID *string) ([]WpisKatalogowyStudia, error) {
	return r.WpisyKatalogowe(ctx, tabelaProfiliStudia, zasieg, zasiegID)
}
