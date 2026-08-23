// Odpowiedzialność pliku: wersje dokumentu modułu Studio (tabela
// `wersja_dokumentu_studio`) — repozytorium sesji pokazywane w Repository Panel.
// Typ i interfejs deklaruje `dane/studio.go`; ten plik implementuje wyłącznie
// metody obszaru wersji na tym samym `*repozytoriumStudia`, jak w module
// Automations (`automations_przebiegi.go`).
//
// Treść wersji: `Tresc` niesie treść krótką wprost, `TrescOdwolanie` — odwołanie
// do pliku dla treści obszernej, tym samym sposobem co `dane/wiadomosci.go`.
//
// Przywrócenie jest zapisem dwutabelowym: `PrzywrocWersje` czyta wersję docelową
// i nadpisuje treść dokumentu w jednej transakcji (`dane/transakcja.go`) — bez
// niej odczyt wersji i zapis dokumentu mogłyby rozjechać się przy równoległym
// zapisie tego samego dokumentu z innego okna.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// WersjaDokumentu to wiersz tabeli `wersja_dokumentu_studio` wraz z kodem
// dokumentu nadrzędnego, doczytanym złączeniem — Repository Panel pokazuje
// wersję zawsze w kontekście jednego dokumentu, więc kod towarzyszy każdemu
// odczytowi tak samo jak `AutomatykaKod` w `Przebieg`.
type WersjaDokumentu struct {
	ID             int64
	Kod            string
	DokumentID     int64
	DokumentKod    string
	Etykieta       *string
	Podsumowanie   *string
	SkrotTresci    *string
	Tresc          *string
	TrescOdwolanie *string
	// Autor rozróżnia zmianę Operatora od zmiany modelu (rozdz. 3.6
	// opracowania). Wskaźnik, nie napis: wersje założone przed dobudową autora
	// nie niosą i mają czytać się jako autor NIEZNANY, a nie jako Operator,
	// którym mogły nie być.
	Autor         *string
	KamienMilowy  bool
	GalazKod      *string
	PropozycjaKod *string
	Utworzono     string
}

const (
	kolumnyWersjiDokumentu = `w.id, w.identyfikator_zewnetrzny, w.dokument_id, d.identyfikator_zewnetrzny,
	                          w.etykieta, w.podsumowanie, w.skrot_tresci, w.tresc, w.tresc_odwolanie,
	                          w.autor, w.kamien_milowy, w.galaz_id, w.propozycja_id,
	                          w.utworzono`

	zrodloWersjiDokumentu = ` FROM wersja_dokumentu_studio w
	                          JOIN dokument_studio d ON d.id = w.dokument_id`

	// Zapis idzie przez ON CONFLICT jak w `ZapiszDokument` — `document.save` z
	// `createVersion=true` bywa powtórzony przy ponowieniu żądania (retry
	// klienta po zerwanym połączeniu), a wtedy drugi zapis tej samej wersji ma
	// nadpisać, nie założyć duplikatu.
	zapiszWersjeDokumentu = `INSERT INTO wersja_dokumentu_studio
	                         (identyfikator_zewnetrzny, dokument_id, etykieta, podsumowanie,
	                          skrot_tresci, tresc, tresc_odwolanie,
	                          autor, kamien_milowy, galaz_id, propozycja_id)
	                         VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	                         ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                             etykieta = excluded.etykieta,
	                             podsumowanie = excluded.podsumowanie,
	                             skrot_tresci = excluded.skrot_tresci,
	                             tresc = excluded.tresc,
	                             tresc_odwolanie = excluded.tresc_odwolanie,
	                             autor = excluded.autor,
	                             kamien_milowy = excluded.kamien_milowy,
	                             galaz_id = excluded.galaz_id,
	                             propozycja_id = excluded.propozycja_id`

	pobierzWersjeDokumentu = `SELECT ` + kolumnyWersjiDokumentu + zrodloWersjiDokumentu +
		` WHERE w.identyfikator_zewnetrzny = ?`

	listaWersjiDokumentu = `SELECT ` + kolumnyWersjiDokumentu + zrodloWersjiDokumentu +
		` WHERE w.dokument_id = ? ORDER BY w.utworzono DESC, w.id DESC`

	// Odczyt wewnątrz transakcji `PrzywrocWersje` — wymusza przynależność
	// wersji do dokumentu (`dokument_id = ?`), więc przywrócenie wersji obcego
	// dokumentu wraca jako brak wiersza, nie jako cudzy zapis.
	pobierzIDDokumentuStudia = `SELECT id FROM dokument_studio WHERE identyfikator_zewnetrzny = ?`

	pobierzTrescWersjiWDokumencie = `SELECT tresc, tresc_odwolanie FROM wersja_dokumentu_studio
	                                 WHERE identyfikator_zewnetrzny = ? AND dokument_id = ?`

	przywrocDokumentStudia = `UPDATE dokument_studio
	                          SET tresc = ?, tresc_odwolanie = ?, wersja_biezaca_id = ?,
	                              zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                          WHERE id = ?`
)

// ZapiszWersje zakłada wersję dokumentu w repozytorium sesji albo nadpisuje
// zastaną i zwraca stan po zapisie. Nie przestawia `wersja_biezaca_id`
// dokumentu — o tym decyduje wywołujący `document.save`, wołając osobno
// `ZapiszDokument`, bo dokument istnieje i bez żadnej wersji.
func (r *repozytoriumStudia) ZapiszWersje(ctx context.Context, dokumentID int64,
	wersja WersjaDokumentu) (WersjaDokumentu, error) {

	if wersja.Kod == "" || dokumentID == 0 {
		return WersjaDokumentu{}, fmt.Errorf("dane: wersja dokumentu studio bez identyfikatora albo bez dokumentu")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszWersjeDokumentu)
	if err != nil {
		return WersjaDokumentu{}, err
	}
	_, err = polecenie.ExecContext(ctx, wersja.Kod, dokumentID, tekstDoKolumny(wersja.Etykieta),
		tekstDoKolumny(wersja.Podsumowanie), tekstDoKolumny(wersja.SkrotTresci),
		tekstDoKolumny(wersja.Tresc), tekstDoKolumny(wersja.TrescOdwolanie),
		tekstDoKolumny(wersja.Autor), liczbaLogiczna(wersja.KamienMilowy),
		tekstDoKolumny(wersja.GalazKod), tekstDoKolumny(wersja.PropozycjaKod))
	if err != nil {
		return WersjaDokumentu{}, fmt.Errorf("dane: nie można zapisać wersji %q dokumentu studio %d: %w",
			wersja.Kod, dokumentID, err)
	}
	return r.Wersja(ctx, wersja.Kod)
}

// Wersje zwraca historię wersji dokumentu od najnowszej — kolejność, w jakiej
// Repository Panel pokazuje repozytorium sesji.
func (r *repozytoriumStudia) Wersje(ctx context.Context, dokumentID int64) ([]WersjaDokumentu, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaWersjiDokumentu)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, dokumentID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać wersji dokumentu studio %d: %w", dokumentID, err)
	}
	defer wiersze.Close()

	lista := []WersjaDokumentu{}
	for wiersze.Next() {
		wersja, err := odczytajWersjeDokumentu(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz wersji dokumentu studio: %w", err)
		}
		lista = append(lista, wersja)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt wersji dokumentu studio %d: %w", dokumentID, err)
	}
	return lista, nil
}

// Wersja zwraca wersję dokumentu o wskazanym kodzie. Brak wiersza wraca jako
// ErrBrakWiersza.
func (r *repozytoriumStudia) Wersja(ctx context.Context, kodWersji string) (WersjaDokumentu, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzWersjeDokumentu)
	if err != nil {
		return WersjaDokumentu{}, err
	}
	wersja, err := odczytajWersjeDokumentu(polecenie.QueryRowContext(ctx, kodWersji))
	if errors.Is(err, sql.ErrNoRows) {
		return WersjaDokumentu{}, ErrBrakWiersza
	}
	if err != nil {
		return WersjaDokumentu{}, fmt.Errorf("dane: nieczytelny wiersz wersji dokumentu studio %q: %w", kodWersji, err)
	}
	return wersja, nil
}

// PrzywrocWersje nadpisuje treść dokumentu treścią wskazanej wersji i
// przestawia `wersja_biezaca_id` na nią, bez usuwania wersji nowszych —
// przywrócenie samo staje się bieżącym stanem, a historia repozytorium sesji
// zostaje nietknięta (Operator może wrócić do dowolnej wersji ponownie).
func (r *repozytoriumStudia) PrzywrocWersje(ctx context.Context, kodDokumentu,
	kodWersji string) (DokumentStudia, error) {

	if kodDokumentu == "" || kodWersji == "" {
		return DokumentStudia{}, fmt.Errorf("dane: przywrócenie wersji wymaga dokumentu i wersji")
	}
	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		poleceniePoDokumencie, err := r.zapytania.wTransakcji(ctx, transakcja, pobierzIDDokumentuStudia)
		if err != nil {
			return err
		}
		var dokumentID int64
		err = poleceniePoDokumencie.QueryRowContext(ctx, kodDokumentu).Scan(&dokumentID)
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%w: dokument studio %q", ErrBrakWiersza, kodDokumentu)
		}
		if err != nil {
			return fmt.Errorf("dane: nieczytelny dokument studio %q: %w", kodDokumentu, err)
		}

		polecenieWersji, err := r.zapytania.wTransakcji(ctx, transakcja, pobierzTrescWersjiWDokumencie)
		if err != nil {
			return err
		}
		var tresc, odwolanie sql.NullString
		err = polecenieWersji.QueryRowContext(ctx, kodWersji, dokumentID).Scan(&tresc, &odwolanie)
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%w: wersja %q dokumentu studio %q", ErrBrakWiersza, kodWersji, kodDokumentu)
		}
		if err != nil {
			return fmt.Errorf("dane: nieczytelna wersja %q dokumentu studio %q: %w", kodWersji, kodDokumentu, err)
		}

		polecenieAktualizacji, err := r.zapytania.wTransakcji(ctx, transakcja, przywrocDokumentStudia)
		if err != nil {
			return err
		}
		_, err = polecenieAktualizacji.ExecContext(ctx, tresc, odwolanie, kodWersji, dokumentID)
		if err != nil {
			return fmt.Errorf("dane: nie można przywrócić wersji %q dokumentu studio %q: %w",
				kodWersji, kodDokumentu, err)
		}
		return nil
	})
	if err != nil {
		return DokumentStudia{}, err
	}
	return r.Dokument(ctx, kodDokumentu)
}

// odczytajWersjeDokumentu składa strukturę z jednego wiersza wyniku.
func odczytajWersjeDokumentu(wiersz skaner) (WersjaDokumentu, error) {
	var wersja WersjaDokumentu
	var etykieta, podsumowanie, skrot, tresc, odwolanie sql.NullString
	var autor, galaz, propozycja sql.NullString
	var kamienMilowy int
	err := wiersz.Scan(&wersja.ID, &wersja.Kod, &wersja.DokumentID, &wersja.DokumentKod,
		&etykieta, &podsumowanie, &skrot, &tresc, &odwolanie,
		&autor, &kamienMilowy, &galaz, &propozycja, &wersja.Utworzono)
	if err != nil {
		return WersjaDokumentu{}, err
	}
	wersja.Etykieta = tekstZKolumny(etykieta)
	wersja.Podsumowanie = tekstZKolumny(podsumowanie)
	wersja.SkrotTresci = tekstZKolumny(skrot)
	wersja.Tresc = tekstZKolumny(tresc)
	wersja.TrescOdwolanie = tekstZKolumny(odwolanie)
	wersja.Autor = tekstZKolumny(autor)
	wersja.KamienMilowy = kamienMilowy == 1
	wersja.GalazKod = tekstZKolumny(galaz)
	wersja.PropozycjaKod = tekstZKolumny(propozycja)
	return wersja, nil
}
