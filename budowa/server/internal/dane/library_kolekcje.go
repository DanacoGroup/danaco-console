// Odpowiedzialność pliku: kolekcje zasobów, przypisania plików do kolekcji i etykiety pliku (tabele
// `kolekcja_biblioteki`, `przypisanie_kolekcji_biblioteki`, `etykieta_pliku_biblioteki`) — obszar Tags
// & Collections modułu Library.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// KolekcjaBiblioteki to wiersz tabeli `kolekcja_biblioteki`. Kod jest
// identyfikatorem, którym kolekcja wychodzi kontraktem
// (`LibraryCollectionCreateResponse.collectionId`).
type KolekcjaBiblioteki struct {
	ID    int64
	Kod   string
	Nazwa string
	Opis  *string
	// RodzicKod i RegulaKod wypełnia wyłącznie odczyt pełny kolekcji swobodnej i korzeniowej domyślnie.
	RodzicKod *string
	RegulaKod *string
	// LiczbaPlikow jest wyliczeniem odczytu, nie kolumną; liczy przypisania w chwili pytania.
	LiczbaPlikow   int
	Utworzono      string
	Zaktualizowano string
}

const (
	kolumnyKolekcjiBiblioteki = `id, identyfikator_zewnetrzny, nazwa, opis, utworzono, zaktualizowano`

	// Kod jest UNIQUE w całej tabeli — DO UPDATE bez warunku konta sięgnąłby cudzego wiersza.
	zapiszKolekcjeBiblioteki = `INSERT INTO kolekcja_biblioteki
	                            (identyfikator_zewnetrzny, nazwa, opis, konto_id)
	                            VALUES (?, ?, ?, ` + WskazanieKonta + `)
	                            ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                                nazwa = excluded.nazwa,
	                                opis = excluded.opis,
	                                zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                            WHERE ` + WarunekKonta

	pobierzKolekcjeBiblioteki = `SELECT ` + kolumnyKolekcjiBiblioteki + ` FROM kolekcja_biblioteki
	                             WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	listaKolekcjiBiblioteki = `SELECT ` + kolumnyKolekcjiBiblioteki + ` FROM kolekcja_biblioteki
	                           WHERE ` + WarunekKonta + `
	                           ORDER BY nazwa, id`

	idPlikuBibliotekiPoKodzie = `SELECT id FROM plik_biblioteki WHERE identyfikator_zewnetrzny = ?`

	// Kod pliku idzie wprost z żądania — bez konta sięgnąłby cudzego wiersza.
	idPlikuBibliotekiPoKodzieWKoncie = `SELECT id FROM plik_biblioteki
	                                    WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	wstawPrzypisanieKolekcji = `INSERT INTO przypisanie_kolekcji_biblioteki (kolekcja_id, plik_id)
	                            VALUES (?, ?)
	                            ON CONFLICT(kolekcja_id, plik_id) DO NOTHING`

	usunEtykietyPliku = `DELETE FROM etykieta_pliku_biblioteki WHERE plik_id = ?`

	wstawEtykietePliku = `INSERT INTO etykieta_pliku_biblioteki (plik_id, etykieta)
	                      VALUES (?, ?)
	                      ON CONFLICT(plik_id, etykieta) DO NOTHING`

	listaEtykietPliku = `SELECT etykieta FROM etykieta_pliku_biblioteki
	                     WHERE plik_id = ? ORDER BY etykieta`
)

// UtworzKolekcje zakłada kolekcję albo nadpisuje zastaną i zwraca stan po
// zapisie — obsługuje `library.collection.create`.
func (r *repozytoriumBiblioteki) UtworzKolekcje(ctx context.Context,
	kolekcja KolekcjaBiblioteki) (KolekcjaBiblioteki, error) {

	if kolekcja.Kod == "" {
		return KolekcjaBiblioteki{}, fmt.Errorf("dane: kolekcja biblioteki bez identyfikatora")
	}
	if kolekcja.Nazwa == "" {
		return KolekcjaBiblioteki{}, fmt.Errorf("dane: kolekcja biblioteki %q bez nazwy", kolekcja.Kod)
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszKolekcjeBiblioteki)
	if err != nil {
		return KolekcjaBiblioteki{}, err
	}
	wynik, err := polecenie.ExecContext(ctx, kolekcja.Kod, kolekcja.Nazwa,
		tekstDoKolumny(kolekcja.Opis), KontoOperatora(ctx), KontoOperatora(ctx))
	if err != nil {
		return KolekcjaBiblioteki{}, fmt.Errorf("dane: nie można zapisać kolekcji biblioteki %q: %w",
			kolekcja.Kod, err)
	}
	if err := sprawdzTrafienieZapisu(wynik, "kolekcja biblioteki", kolekcja.Kod); err != nil {
		return KolekcjaBiblioteki{}, err
	}
	return r.kolekcjaPoKodzie(ctx, kolekcja.Kod)
}

// Kolekcje zwraca wszystkie kolekcje biblioteki, uporządkowane alfabetycznie rosnąco po ich pełnej nazwie.
func (r *repozytoriumBiblioteki) Kolekcje(ctx context.Context) ([]KolekcjaBiblioteki, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaKolekcjiBiblioteki)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, KontoOperatora(ctx))
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać kolekcji biblioteki: %w", err)
	}
	defer wiersze.Close()

	lista := []KolekcjaBiblioteki{}
	for wiersze.Next() {
		kolekcja, err := odczytajKolekcjeBiblioteki(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz kolekcji biblioteki: %w", err)
		}
		lista = append(lista, kolekcja)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt kolekcji biblioteki: %w", err)
	}
	return lista, nil
}

// PrzypiszDoKolekcji przypisuje pliki do kolekcji i oddaje kody plików faktycznie przypisanych, pomijając
// plik, którego nie ma.
func (r *repozytoriumBiblioteki) PrzypiszDoKolekcji(ctx context.Context,
	kodKolekcji string, kodyPlikow []string) ([]string, error) {

	kolekcja, err := r.kolekcjaPoKodzie(ctx, kodKolekcji)
	if err != nil {
		return nil, err
	}

	poszukiwaniePliku, err := r.zapytania.przygotuj(ctx, idPlikuBibliotekiPoKodzieWKoncie)
	if err != nil {
		return nil, err
	}
	wstawianie, err := r.zapytania.przygotuj(ctx, wstawPrzypisanieKolekcji)
	if err != nil {
		return nil, err
	}

	przypisane := []string{}
	for _, kodPliku := range kodyPlikow {
		var plikID int64
		err := poszukiwaniePliku.QueryRowContext(ctx, kodPliku, KontoOperatora(ctx)).Scan(&plikID)
		if errors.Is(err, sql.ErrNoRows) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("dane: nie można odnaleźć pliku biblioteki %q: %w", kodPliku, err)
		}
		if _, err := wstawianie.ExecContext(ctx, kolekcja.ID, plikID); err != nil {
			return nil, fmt.Errorf("dane: nie można przypisać pliku %q do kolekcji %q: %w",
				kodPliku, kodKolekcji, err)
		}
		przypisane = append(przypisane, kodPliku)
	}
	return przypisane, nil
}

// UstawEtykiety podmienia komplet etykiet pliku, bo `library.tag.set` nadsyła
// zawsze pełny zestaw, nie różnicę. Usunięcie i wstawienie zachodzi w jednej
// transakcji — plik nie zostaje przejściowo bez etykiet przy błędzie w trakcie.
func (r *repozytoriumBiblioteki) UstawEtykiety(ctx context.Context,
	kodPliku string, etykiety []string) ([]string, error) {

	plik, err := r.Plik(ctx, kodPliku)
	if err != nil {
		return nil, err
	}

	err = wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		czyszczenie, err := r.zapytania.wTransakcji(ctx, transakcja, usunEtykietyPliku)
		if err != nil {
			return err
		}
		if _, err := czyszczenie.ExecContext(ctx, plik.ID); err != nil {
			return fmt.Errorf("dane: nie można wyczyścić etykiet pliku %q: %w", kodPliku, err)
		}
		wstawienie, err := r.zapytania.wTransakcji(ctx, transakcja, wstawEtykietePliku)
		if err != nil {
			return err
		}
		for _, etykieta := range etykiety {
			if etykieta == "" {
				continue
			}
			if _, err := wstawienie.ExecContext(ctx, plik.ID, etykieta); err != nil {
				return fmt.Errorf("dane: nie można zapisać etykiety %q pliku %q: %w",
					etykieta, kodPliku, err)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return r.Etykiety(ctx, plik.ID)
}

// Etykiety zwraca etykiety przypisane danemu plikowi biblioteki w porządku alfabetycznym rosnącym po treści.
func (r *repozytoriumBiblioteki) Etykiety(ctx context.Context, plikID int64) ([]string, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaEtykietPliku)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, plikID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać etykiet pliku %d: %w", plikID, err)
	}
	defer wiersze.Close()

	lista := []string{}
	for wiersze.Next() {
		var etykieta string
		if err := wiersze.Scan(&etykieta); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz etykiety pliku %d: %w", plikID, err)
		}
		lista = append(lista, etykieta)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt etykiet pliku %d: %w", plikID, err)
	}
	return lista, nil
}

// kolekcjaPoKodzie zwraca kolekcję o wskazanym kodzie. Brak wiersza wraca jako
// ErrBrakWiersza — wspólne dla `UtworzKolekcje` i `PrzypiszDoKolekcji`.
func (r *repozytoriumBiblioteki) kolekcjaPoKodzie(ctx context.Context, kod string) (KolekcjaBiblioteki, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzKolekcjeBiblioteki)
	if err != nil {
		return KolekcjaBiblioteki{}, err
	}
	kolekcja, err := odczytajKolekcjeBiblioteki(polecenie.QueryRowContext(ctx, kod, KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return KolekcjaBiblioteki{}, ErrBrakWiersza
	}
	if err != nil {
		return KolekcjaBiblioteki{}, fmt.Errorf("dane: nieczytelny wiersz kolekcji biblioteki %q: %w", kod, err)
	}
	return kolekcja, nil
}

// odczytajKolekcjeBiblioteki składa pełną strukturę kolekcji z jednego wiersza wyniku zapytania do bazy.
func odczytajKolekcjeBiblioteki(wiersz skaner) (KolekcjaBiblioteki, error) {
	var kolekcja KolekcjaBiblioteki
	var opis sql.NullString
	err := wiersz.Scan(&kolekcja.ID, &kolekcja.Kod, &kolekcja.Nazwa, &opis,
		&kolekcja.Utworzono, &kolekcja.Zaktualizowano)
	if err != nil {
		return KolekcjaBiblioteki{}, err
	}
	kolekcja.Opis = tekstZKolumny(opis)
	return kolekcja, nil
}
