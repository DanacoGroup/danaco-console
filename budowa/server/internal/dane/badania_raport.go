// Plik utrzymuje raport badania, jego sekcje, eksporty oraz jednowierszową
// przestrzeń badania, jako trzecią część repozytorium badań.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// RaportBadania to wiersz tabeli raport_badania: raport bez sekcji, oddawanych osobno metodą Sekcje raportu.
type RaportBadania struct {
	ID             int64
	Kod            string
	Okno           string
	Tytul          string
	Utworzono      string
	Zaktualizowano string
}

// SekcjaRaportu to wiersz `sekcja_raportu_badania`; treść niesie para
// `Tresc`/`TrescOdwolanie` — tekst wprost albo odwołanie do zasobu.
type SekcjaRaportu struct {
	ID             int64
	Kod            string
	RaportID       int64
	Tytul          string
	Tresc          *string
	TrescOdwolanie *string
	Kolejnosc      int
	UstalenieKody  []string
	Utworzono      string
}

// EksportRaportu to wiersz `eksport_raportu_badania` — trwały ślad wyniku
// eksportu. `RaportKod` jest kodem zewnętrznym raportu.
type EksportRaportu struct {
	ID               int64
	Kod              string
	RaportKod        string
	Format           string
	SciezkaDocelowa  *string
	PlikBibliotekiID *string
	SciezkaWyniku    *string
	RozmiarBajtow    *int64
	// Cel jest miejscem docelowym eksportu: pobranie, Library, Studio albo Roundtable.
	Cel       string
	Utworzono string
}

const (
	kolumnyRaportuBadania = `id, identyfikator_zewnetrzny, okno, tytul, utworzono, zaktualizowano`
	zapiszRaportBadania   = `INSERT INTO raport_badania
	                       (identyfikator_zewnetrzny, okno, tytul, zaktualizowano)
	                       VALUES (?, ?, ?, strftime('%Y-%m-%dT%H:%M:%fZ','now'))
	                       ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                           tytul = excluded.tytul, zaktualizowano = excluded.zaktualizowano`
	pobierzRaportBadania = `SELECT ` + kolumnyRaportuBadania + ` FROM raport_badania
	                        WHERE identyfikator_zewnetrzny = ?`
	usunSekcjeRaportuBadania  = `DELETE FROM sekcja_raportu_badania WHERE raport_id = ?`
	wstawSekcjeRaportuBadania = `INSERT INTO sekcja_raportu_badania
	                             (identyfikator_zewnetrzny, raport_id, tytul, tresc,
	                              tresc_odwolanie, kolejnosc) VALUES (?, ?, ?, ?, ?, ?)`
	znajdzIDUstaleniaPoKodzie   = `SELECT id FROM ustalenie_badania WHERE identyfikator_zewnetrzny = ?`
	wstawUstalenieSekcjiRaportu = `INSERT INTO ustalenie_sekcji_raportu_badania
	                               (sekcja_id, ustalenie_id) VALUES (?, ?)`
	listaSekcjiRaportuBadania = `SELECT id, identyfikator_zewnetrzny, raport_id, tytul, tresc,
	                                     tresc_odwolanie, kolejnosc, utworzono
	                             FROM sekcja_raportu_badania WHERE raport_id = ? ORDER BY kolejnosc, id`
	listaKodowUstalenSekcji = `SELECT u.identyfikator_zewnetrzny
	                           FROM ustalenie_sekcji_raportu_badania s
	                           JOIN ustalenie_badania u ON u.id = s.ustalenie_id
	                           WHERE s.sekcja_id = ? ORDER BY u.id`
	znajdzIDRaportuPoKodzie    = `SELECT id FROM raport_badania WHERE identyfikator_zewnetrzny = ?`
	wstawEksportRaportuBadania = `INSERT INTO eksport_raportu_badania
	                              (identyfikator_zewnetrzny, raport_id, format, cel, sciezka_docelowa,
	                               plik_biblioteki_id, sciezka_wyniku, rozmiar_bajtow)
	                              VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	pobierzEksportRaportuBadania = `SELECT id, identyfikator_zewnetrzny, raport_id, format, cel,
	                                        sciezka_docelowa, plik_biblioteki_id, sciezka_wyniku,
	                                        rozmiar_bajtow, utworzono
	                                FROM eksport_raportu_badania WHERE identyfikator_zewnetrzny = ?`
	ustawPrzestrzenBadania = `INSERT INTO przestrzen_badania (id, zakres, zaktualizowano)
	                          VALUES (1, ?, strftime('%Y-%m-%dT%H:%M:%fZ','now'))
	                          ON CONFLICT(id) DO UPDATE SET
	                              zakres = excluded.zakres, zaktualizowano = excluded.zaktualizowano`
	usunEtapyPrzestrzeniBadania     = `DELETE FROM etap_przestrzeni_badania`
	wstawEtapPrzestrzeniBadania     = `INSERT INTO etap_przestrzeni_badania (kolejnosc, etap) VALUES (?, ?)`
	pobierzZakresPrzestrzeniBadania = `SELECT zakres FROM przestrzen_badania WHERE id = 1`
	listaEtapowPrzestrzeniBadania   = `SELECT etap FROM etap_przestrzeni_badania ORDER BY kolejnosc`
)

// ZapiszRaport zakłada albo nadpisuje raport po kodzie i podmienia komplet jego sekcji w jednej transakcji.
func (r *repozytoriumBadan) ZapiszRaport(ctx context.Context, raport RaportBadania, sekcje []SekcjaRaportu) (RaportBadania, error) {
	if raport.Kod == "" {
		return RaportBadania{}, fmt.Errorf("dane: raport badania bez identyfikatora")
	}
	if raport.Okno == "" {
		return RaportBadania{}, fmt.Errorf("dane: raport badania %q bez okna", raport.Kod)
	}
	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		zapis, err := r.zapytania.wTransakcji(ctx, transakcja, zapiszRaportBadania)
		if err != nil {
			return err
		}
		if _, err := zapis.ExecContext(ctx, raport.Kod, raport.Okno, raport.Tytul); err != nil {
			return fmt.Errorf("dane: nie można zapisać raportu badania %q: %w", raport.Kod, err)
		}
		odczyt, err := r.zapytania.wTransakcji(ctx, transakcja, pobierzRaportBadania)
		if err != nil {
			return err
		}
		zapisany, err := odczytajRaportBadania(odczyt.QueryRowContext(ctx, raport.Kod))
		if err != nil {
			return fmt.Errorf("dane: nie można odczytać zapisanego raportu badania %q: %w", raport.Kod, err)
		}
		czyszczenie, err := r.zapytania.wTransakcji(ctx, transakcja, usunSekcjeRaportuBadania)
		if err != nil {
			return err
		}
		if _, err := czyszczenie.ExecContext(ctx, zapisany.ID); err != nil {
			return fmt.Errorf("dane: nie można wyczyścić sekcji raportu badania %q: %w", raport.Kod, err)
		}
		return zapiszSekcjeRaportu(ctx, r.zapytania, transakcja, zapisany.ID, raport.Kod, sekcje)
	})
	if err != nil {
		return RaportBadania{}, err
	}
	return r.Raport(ctx, raport.Kod)
}

// zapiszSekcjeRaportu wstawia sekcje raportu i ich powiązania z ustaleniami badania w jednej transakcji.
func zapiszSekcjeRaportu(ctx context.Context, z *zapytania, transakcja *sql.Tx, raportID int64, kodRaportu string, sekcje []SekcjaRaportu) error {
	wstawienie, err := z.wTransakcji(ctx, transakcja, wstawSekcjeRaportuBadania)
	if err != nil {
		return err
	}
	szukanieUstalenia, err := z.wTransakcji(ctx, transakcja, znajdzIDUstaleniaPoKodzie)
	if err != nil {
		return err
	}
	wstawieniePowiazania, err := z.wTransakcji(ctx, transakcja, wstawUstalenieSekcjiRaportu)
	if err != nil {
		return err
	}
	for numer, sekcja := range sekcje {
		if sekcja.Kod == "" {
			return fmt.Errorf("dane: sekcja numer %d raportu badania %q bez identyfikatora", numer, kodRaportu)
		}
		kolejnosc := sekcja.Kolejnosc
		if kolejnosc == 0 {
			kolejnosc = numer + 1
		}
		wynik, err := wstawienie.ExecContext(ctx, sekcja.Kod, raportID, sekcja.Tytul, tekstDoKolumny(sekcja.Tresc), tekstDoKolumny(sekcja.TrescOdwolanie), kolejnosc)
		if err != nil {
			return fmt.Errorf("dane: nie można zapisać sekcji %q raportu badania %q: %w", sekcja.Kod, kodRaportu, err)
		}
		sekcjaID, err := wynik.LastInsertId()
		if err != nil {
			return fmt.Errorf("dane: nie można ustalić identyfikatora sekcji %q raportu badania %q: %w", sekcja.Kod, kodRaportu, err)
		}
		for _, kodUstalenia := range sekcja.UstalenieKody {
			var ustalenieID int64
			err := szukanieUstalenia.QueryRowContext(ctx, kodUstalenia).Scan(&ustalenieID)
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Errorf("dane: sekcja %q odwołuje się do nieznanego ustalenia %q: %w", sekcja.Kod, kodUstalenia, ErrBrakWiersza)
			}
			if err != nil {
				return fmt.Errorf("dane: nie można znaleźć ustalenia %q dla sekcji %q: %w", kodUstalenia, sekcja.Kod, err)
			}
			if _, err := wstawieniePowiazania.ExecContext(ctx, sekcjaID, ustalenieID); err != nil {
				return fmt.Errorf("dane: nie można powiązać sekcji %q z ustaleniem %q: %w", sekcja.Kod, kodUstalenia, err)
			}
		}
	}
	return nil
}

// Raport zwraca raport badania o wskazanym kodzie zewnętrznym; brak wiersza wraca jako ErrBrakWiersza.
func (r *repozytoriumBadan) Raport(ctx context.Context, kod string) (RaportBadania, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzRaportBadania)
	if err != nil {
		return RaportBadania{}, err
	}
	raport, err := odczytajRaportBadania(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return RaportBadania{}, ErrBrakWiersza
	}
	if err != nil {
		return RaportBadania{}, fmt.Errorf("dane: nieczytelny wiersz raportu badania %q: %w", kod, err)
	}
	return raport, nil
}

// odczytajRaportBadania składa całą strukturę raportu badania z jednego wiersza wyniku danego zapytania.
func odczytajRaportBadania(wiersz skaner) (RaportBadania, error) {
	var raport RaportBadania
	err := wiersz.Scan(&raport.ID, &raport.Kod, &raport.Okno, &raport.Tytul,
		&raport.Utworzono, &raport.Zaktualizowano)
	if err != nil {
		return RaportBadania{}, err
	}
	return raport, nil
}

// Sekcje zwraca sekcje raportu w zapisanej kolejności wraz z kodami zasilających je ustaleń badania wprost.
func (r *repozytoriumBadan) Sekcje(ctx context.Context, raportID int64) ([]SekcjaRaportu, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaSekcjiRaportuBadania)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, raportID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać sekcji raportu badania %d: %w", raportID, err)
	}
	defer wiersze.Close()

	lista := []SekcjaRaportu{}
	for wiersze.Next() {
		var sekcja SekcjaRaportu
		var tresc, odwolanie sql.NullString
		err := wiersze.Scan(&sekcja.ID, &sekcja.Kod, &sekcja.RaportID, &sekcja.Tytul,
			&tresc, &odwolanie, &sekcja.Kolejnosc, &sekcja.Utworzono)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz sekcji raportu badania %d: %w", raportID, err)
		}
		sekcja.Tresc = tekstZKolumny(tresc)
		sekcja.TrescOdwolanie = tekstZKolumny(odwolanie)
		lista = append(lista, sekcja)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt sekcji raportu badania %d: %w", raportID, err)
	}

	// Kody ustaleń dociągane osobnym zapytaniem na sekcję: liczba sekcji jest mała, N+1 nie waży wiele.
	powiazania, err := r.zapytania.przygotuj(ctx, listaKodowUstalenSekcji)
	if err != nil {
		return nil, err
	}
	for indeks := range lista {
		wierszeKodow, err := powiazania.QueryContext(ctx, lista[indeks].ID)
		if err != nil {
			return nil, fmt.Errorf("dane: nie można odczytać ustaleń sekcji %d: %w", lista[indeks].ID, err)
		}
		kody := []string{}
		for wierszeKodow.Next() {
			var kod string
			if err := wierszeKodow.Scan(&kod); err != nil {
				wierszeKodow.Close()
				return nil, fmt.Errorf("dane: nieczytelny kod ustalenia sekcji %d: %w", lista[indeks].ID, err)
			}
			kody = append(kody, kod)
		}
		blad := wierszeKodow.Err()
		wierszeKodow.Close()
		if blad != nil {
			return nil, fmt.Errorf("dane: przerwany odczyt ustaleń sekcji %d: %w", lista[indeks].ID, blad)
		}
		lista[indeks].UstalenieKody = kody
	}
	return lista, nil
}

// ZapiszEksport zapisuje ślad eksportu raportu badania: dokłada wiersz historii, nic nigdy nie nadpisuje.
func (r *repozytoriumBadan) ZapiszEksport(ctx context.Context, eksport EksportRaportu) (EksportRaportu, error) {
	if eksport.Kod == "" {
		return EksportRaportu{}, fmt.Errorf("dane: eksport raportu badania bez identyfikatora")
	}
	if eksport.RaportKod == "" {
		return EksportRaportu{}, fmt.Errorf("dane: eksport %q bez kodu raportu", eksport.Kod)
	}
	szukanie, err := r.zapytania.przygotuj(ctx, znajdzIDRaportuPoKodzie)
	if err != nil {
		return EksportRaportu{}, err
	}
	var raportID int64
	err = szukanie.QueryRowContext(ctx, eksport.RaportKod).Scan(&raportID)
	if errors.Is(err, sql.ErrNoRows) {
		return EksportRaportu{}, fmt.Errorf("dane: raport badania %q nie istnieje: %w", eksport.RaportKod, ErrBrakWiersza)
	}
	if err != nil {
		return EksportRaportu{}, fmt.Errorf("dane: nie można znaleźć raportu badania %q: %w", eksport.RaportKod, err)
	}
	wstawienie, err := r.zapytania.przygotuj(ctx, wstawEksportRaportuBadania)
	if err != nil {
		return EksportRaportu{}, err
	}
	cel := eksport.Cel
	if cel == "" {
		cel = "download"
	}
	if _, err := wstawienie.ExecContext(ctx, eksport.Kod, raportID, eksport.Format, cel,
		tekstDoKolumny(eksport.SciezkaDocelowa),
		tekstDoKolumny(eksport.PlikBibliotekiID), tekstDoKolumny(eksport.SciezkaWyniku), liczbaDoKolumny(eksport.RozmiarBajtow)); err != nil {
		return EksportRaportu{}, fmt.Errorf("dane: nie można zapisać eksportu %q raportu badania %q: %w", eksport.Kod, eksport.RaportKod, err)
	}
	odczyt, err := r.zapytania.przygotuj(ctx, pobierzEksportRaportuBadania)
	if err != nil {
		return EksportRaportu{}, err
	}
	zapisany, err := odczytajEksportRaportu(odczyt.QueryRowContext(ctx, eksport.Kod))
	if err != nil {
		return EksportRaportu{}, fmt.Errorf("dane: nie można odczytać zapisanego eksportu %q: %w", eksport.Kod, err)
	}
	zapisany.RaportKod = eksport.RaportKod
	return zapisany, nil
}

// odczytajEksportRaportu składa strukturę eksportu z jednego wiersza; RaportKod wywołujący uzupełnia sam.
func odczytajEksportRaportu(wiersz skaner) (EksportRaportu, error) {
	var eksport EksportRaportu
	var raportID int64
	var sciezkaDocelowa, plikBiblioteki, sciezkaWyniku sql.NullString
	var rozmiar sql.NullInt64
	err := wiersz.Scan(&eksport.ID, &eksport.Kod, &raportID, &eksport.Format, &eksport.Cel,
		&sciezkaDocelowa, &plikBiblioteki, &sciezkaWyniku, &rozmiar, &eksport.Utworzono)
	if err != nil {
		return EksportRaportu{}, err
	}
	eksport.SciezkaDocelowa = tekstZKolumny(sciezkaDocelowa)
	eksport.PlikBibliotekiID = tekstZKolumny(plikBiblioteki)
	eksport.SciezkaWyniku = tekstZKolumny(sciezkaWyniku)
	eksport.RozmiarBajtow = liczbaZKolumny(rozmiar)
	return eksport, nil
}

// UstawPrzestrzen nadpisuje jedyny wiersz przestrzeni badania i wymienia jej wszystkie etapy w całości.
func (r *repozytoriumBadan) UstawPrzestrzen(ctx context.Context, zakres string, etapy []string) (string, []string, error) {
	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		zapis, err := r.zapytania.wTransakcji(ctx, transakcja, ustawPrzestrzenBadania)
		if err != nil {
			return err
		}
		if _, err := zapis.ExecContext(ctx, zakres); err != nil {
			return fmt.Errorf("dane: nie można zapisać przestrzeni badania: %w", err)
		}
		czyszczenie, err := r.zapytania.wTransakcji(ctx, transakcja, usunEtapyPrzestrzeniBadania)
		if err != nil {
			return err
		}
		if _, err := czyszczenie.ExecContext(ctx); err != nil {
			return fmt.Errorf("dane: nie można wyczyścić etapów przestrzeni badania: %w", err)
		}
		wstawienie, err := r.zapytania.wTransakcji(ctx, transakcja, wstawEtapPrzestrzeniBadania)
		if err != nil {
			return err
		}
		for indeks, etap := range etapy {
			if _, err := wstawienie.ExecContext(ctx, indeks+1, etap); err != nil {
				return fmt.Errorf("dane: nie można zapisać etapu %q przestrzeni badania: %w", etap, err)
			}
		}
		return nil
	})
	if err != nil {
		return "", nil, err
	}
	return r.Przestrzen(ctx)
}

// Przestrzen zwraca zakres i etapy badania; tabela pusta wraca jako para pusta,
// nie jako błąd — to stan startowy.
func (r *repozytoriumBadan) Przestrzen(ctx context.Context) (string, []string, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzZakresPrzestrzeniBadania)
	if err != nil {
		return "", nil, err
	}
	var zakres string
	err = polecenie.QueryRowContext(ctx).Scan(&zakres)
	if errors.Is(err, sql.ErrNoRows) {
		return "", []string{}, nil
	}
	if err != nil {
		return "", nil, fmt.Errorf("dane: nie można odczytać zakresu przestrzeni badania: %w", err)
	}
	listaEtapow, err := r.zapytania.przygotuj(ctx, listaEtapowPrzestrzeniBadania)
	if err != nil {
		return "", nil, err
	}
	wiersze, err := listaEtapow.QueryContext(ctx)
	if err != nil {
		return "", nil, fmt.Errorf("dane: nie można odczytać etapów przestrzeni badania: %w", err)
	}
	defer wiersze.Close()
	etapy := []string{}
	for wiersze.Next() {
		var etap string
		if err := wiersze.Scan(&etap); err != nil {
			return "", nil, fmt.Errorf("dane: nieczytelny etap przestrzeni badania: %w", err)
		}
		etapy = append(etapy, etap)
	}
	if err := wiersze.Err(); err != nil {
		return "", nil, fmt.Errorf("dane: przerwany odczyt etapów przestrzeni badania: %w", err)
	}
	return zakres, etapy, nil
}
