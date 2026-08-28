// Odpowiedzialność pliku: obszar dokumentu modułu Studio wraz z kontraktem całego obszaru, w tym wersji, propozycji i postaci dokumentu.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"danacoconsole/shared"
)

// DokumentStudia to wiersz tabeli dokument_studio, niosący treść, tytuł i kod dokumentu w module Studio.
type DokumentStudia struct {
	ID                 int64
	Kod                string
	Okno               string
	Tytul              *string
	Format             shared.StudioDocumentFormat
	Tresc              *string
	TrescOdwolanie     *string
	PlikRepozytoriumID *string
	WersjaBiezacaKod   *string
	Utworzono          string
	Zaktualizowano     string
}

// RepozytoriumStudia jest kontraktem obszaru Studio; postać dokumentu wchodzi tu zagnieżdżonym interfejsem osobnego repozytorium postaci.
type RepozytoriumStudia interface {
	RepozytoriumPostaciStudia

	// --- dokument ---
	ZapiszDokument(ctx context.Context, dokument DokumentStudia) (DokumentStudia, error)
	Dokument(ctx context.Context, kod string) (DokumentStudia, error)
	Dokumenty(ctx context.Context, okno string) ([]DokumentStudia, error)

	// --- wersje ---
	ZapiszWersje(ctx context.Context, dokumentID int64, wersja WersjaDokumentu) (WersjaDokumentu, error)
	Wersje(ctx context.Context, dokumentID int64) ([]WersjaDokumentu, error)
	Wersja(ctx context.Context, kodWersji string) (WersjaDokumentu, error)
	PrzywrocWersje(ctx context.Context, kodDokumentu, kodWersji string) (DokumentStudia, error)

	// --- propozycje ---
	ZapiszPropozycje(ctx context.Context, dokumentID int64, propozycja PropozycjaZmiany) (PropozycjaZmiany, error)
	Propozycja(ctx context.Context, kod string) (PropozycjaZmiany, error)
	Propozycje(ctx context.Context, dokumentID int64) ([]PropozycjaZmiany, error)

	// --- komentarze redakcyjne i adnotacje różnic (`studio_adnotacje.go`) ---
	ZapiszKomentarz(ctx context.Context, dokumentID int64, komentarz KomentarzStudia) (KomentarzStudia, error)
	Komentarz(ctx context.Context, kod string) (KomentarzStudia, error)
	Komentarze(ctx context.Context, dokumentID int64, rodzaj string) ([]KomentarzStudia, error)
	RozstrzygnijKomentarz(ctx context.Context, kod string, rozwiazany bool) (KomentarzStudia, error)

	// --- śledzenie zmian (`studio_adnotacje.go`) ---
	ZapiszZmianeSledzona(ctx context.Context, dokumentID int64, zmiana ZmianaSledzona) (ZmianaSledzona, error)
	ZmianySledzone(ctx context.Context, dokumentID int64) ([]ZmianaSledzona, error)
	RozstrzygnijZmianeSledzona(ctx context.Context, kod, decyzja string) (bool, error)
	UstawSledzenie(ctx context.Context, kodDokumentu string, czynne bool) error
	Sledzenie(ctx context.Context, kodDokumentu string) (bool, error)

	// --- kolejka cyfryzacji (`studio_cyfryzacja.go`) ---
	ZapiszPozycjeWczytywania(ctx context.Context, pozycja PozycjaWczytywania) (PozycjaWczytywania, error)
	PozycjaWczytywania(ctx context.Context, kod string) (PozycjaWczytywania, error)
	PozycjeWczytywania(ctx context.Context, okno string, tylkoNieprzetworzone bool) ([]PozycjaWczytywania, error)

	// --- katalogi zasięgu: operacje, łańcuchy, profile (`studio_katalogi.go`) ---
	ZapiszOperacje(ctx context.Context, wpis WpisKatalogowyStudia) (WpisKatalogowyStudia, error)
	Operacje(ctx context.Context, zasieg string, zasiegID *string) ([]WpisKatalogowyStudia, error)
	UsunOperacje(ctx context.Context, kod string) (bool, error)
	ZapiszLancuch(ctx context.Context, wpis WpisKatalogowyStudia) (WpisKatalogowyStudia, error)
	Lancuch(ctx context.Context, kod string) (WpisKatalogowyStudia, error)
	Lancuchy(ctx context.Context, zasieg string, zasiegID *string) ([]WpisKatalogowyStudia, error)
	ZapiszProfilWydania(ctx context.Context, wpis WpisKatalogowyStudia) (WpisKatalogowyStudia, error)
	ProfilWydania(ctx context.Context, kod string) (WpisKatalogowyStudia, error)
	ProfileWydania(ctx context.Context, zasieg string, zasiegID *string) ([]WpisKatalogowyStudia, error)

	// --- szablony, gałęzie, odwołania (`studio_katalogi.go`) ---
	ZapiszSzablon(ctx context.Context, szablon SzablonStudia) (SzablonStudia, error)
	Szablon(ctx context.Context, kod string) (SzablonStudia, error)
	Szablony(ctx context.Context) ([]SzablonStudia, error)
	ZapiszGalaz(ctx context.Context, dokumentID int64, galaz GalazStudia) (GalazStudia, error)
	Galaz(ctx context.Context, kod string) (GalazStudia, error)
	Galezie(ctx context.Context, dokumentID int64) ([]GalazStudia, error)
	ZapiszOdwolanieWersji(ctx context.Context, odwolanie OdwolanieWersji) (OdwolanieWersji, error)

	// --- przebieg wsadu (`studio_wsad.go`) ---
	ZapiszPrzebiegWsadu(ctx context.Context, przebieg PrzebiegWsaduStudia,
		pozycje []PozycjaWsaduStudia) (PrzebiegWsaduStudia, error)
	PrzebiegWsadu(ctx context.Context, kod string) (PrzebiegWsaduStudia, error)
	PozycjeWsadu(ctx context.Context, kodPrzebiegu string) ([]PozycjaWsaduStudia, error)
}

const (
	kolumnyDokumentuStudia = `id, identyfikator_zewnetrzny, okno, tytul, format, tresc,
	                          tresc_odwolanie, plik_repozytorium_id, wersja_biezaca_id,
	                          utworzono, zaktualizowano`

	// Zapis zakłada dokument albo nadpisuje zastany po identyfikatorze
	// zewnętrznym — `document.open` wznawia sesję po `documentId`, więc drugi
	// zapis tego samego dokumentu jest normalną ścieżką, nie usterką.
	zapiszDokumentStudia = `INSERT INTO dokument_studio
	                        (identyfikator_zewnetrzny, okno, tytul, format, tresc,
	                         tresc_odwolanie, plik_repozytorium_id, wersja_biezaca_id)
	                        VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	                        ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                            okno = excluded.okno,
	                            tytul = excluded.tytul,
	                            format = excluded.format,
	                            tresc = excluded.tresc,
	                            tresc_odwolanie = excluded.tresc_odwolanie,
	                            plik_repozytorium_id = excluded.plik_repozytorium_id,
	                            wersja_biezaca_id = excluded.wersja_biezaca_id,
	                            zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')`

	pobierzDokumentStudia = `SELECT ` + kolumnyDokumentuStudia + ` FROM dokument_studio
	                         WHERE identyfikator_zewnetrzny = ?`

	listaDokumentowStudia = `SELECT ` + kolumnyDokumentuStudia + ` FROM dokument_studio
	                         WHERE okno = ? ORDER BY zaktualizowano DESC, id DESC`
)

type repozytoriumStudia struct {
	zapytania *zapytania
	db        *sql.DB
}

// noweRepozytoriumStudia zakłada repozytorium obszaru Studio. Parametr `db`
// służy obszarom wersji i propozycji, które zapisują w transakcji — sam
// dokument go nie potrzebuje, bo każdy jego zapis mieści się w jednym poleceniu.
func noweRepozytoriumStudia(z *zapytania, db *sql.DB) *repozytoriumStudia {
	return &repozytoriumStudia{zapytania: z, db: db}
}

// ZapiszDokument zakłada dokument Studio albo nadpisuje zastany i zwraca jego pełny stan po samym zapisie.
func (r *repozytoriumStudia) ZapiszDokument(ctx context.Context,
	dokument DokumentStudia) (DokumentStudia, error) {

	if dokument.Kod == "" {
		return DokumentStudia{}, fmt.Errorf("dane: dokument studio bez identyfikatora")
	}
	if dokument.Format == "" {
		dokument.Format = shared.StudioDocumentFormatTxt
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszDokumentStudia)
	if err != nil {
		return DokumentStudia{}, err
	}
	_, err = polecenie.ExecContext(ctx, dokument.Kod, dokument.Okno,
		tekstDoKolumny(dokument.Tytul), string(dokument.Format),
		tekstDoKolumny(dokument.Tresc), tekstDoKolumny(dokument.TrescOdwolanie),
		tekstDoKolumny(dokument.PlikRepozytoriumID), tekstDoKolumny(dokument.WersjaBiezacaKod))
	if err != nil {
		return DokumentStudia{}, fmt.Errorf("dane: nie można zapisać dokumentu studio %q: %w",
			dokument.Kod, err)
	}
	return r.Dokument(ctx, dokument.Kod)
}

// Dokument zwraca dokument Studio o wskazanym kodzie zewnętrznym; brak wiersza wraca jako ErrBrakWiersza.
func (r *repozytoriumStudia) Dokument(ctx context.Context, kod string) (DokumentStudia, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzDokumentStudia)
	if err != nil {
		return DokumentStudia{}, err
	}
	dokument, err := odczytajDokumentStudia(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return DokumentStudia{}, ErrBrakWiersza
	}
	if err != nil {
		return DokumentStudia{}, fmt.Errorf("dane: nieczytelny wiersz dokumentu studio %q: %w", kod, err)
	}
	return dokument, nil
}

// Dokumenty zwraca dokumenty otwarte w danym oknie operacyjnym, od najświeżej zmienionego z nich dokumentu.
func (r *repozytoriumStudia) Dokumenty(ctx context.Context, okno string) ([]DokumentStudia, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, listaDokumentowStudia)
	if err != nil {
		return nil, err
	}
	wiersze, err := polecenie.QueryContext(ctx, okno)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać dokumentów studio okna %q: %w", okno, err)
	}
	defer wiersze.Close()

	lista := []DokumentStudia{}
	for wiersze.Next() {
		dokument, err := odczytajDokumentStudia(wiersze)
		if err != nil {
			return nil, fmt.Errorf("dane: nieczytelny wiersz dokumentu studio: %w", err)
		}
		lista = append(lista, dokument)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt dokumentów studio okna %q: %w", okno, err)
	}
	return lista, nil
}

// odczytajDokumentStudia składa strukturę dokumentu z jednego wiersza wyniku zapytania, kolumna po kolumnie.
func odczytajDokumentStudia(wiersz skaner) (DokumentStudia, error) {
	var dokument DokumentStudia
	var tytul, tresc, odwolanie, plikRepo, wersjaBiezaca sql.NullString
	var format string
	err := wiersz.Scan(&dokument.ID, &dokument.Kod, &dokument.Okno, &tytul, &format,
		&tresc, &odwolanie, &plikRepo, &wersjaBiezaca, &dokument.Utworzono, &dokument.Zaktualizowano)
	if err != nil {
		return DokumentStudia{}, err
	}
	dokument.Format = shared.StudioDocumentFormat(format)
	dokument.Tytul = tekstZKolumny(tytul)
	dokument.Tresc = tekstZKolumny(tresc)
	dokument.TrescOdwolanie = tekstZKolumny(odwolanie)
	dokument.PlikRepozytoriumID = tekstZKolumny(plikRepo)
	dokument.WersjaBiezacaKod = tekstZKolumny(wersjaBiezaca)
	return dokument, nil
}
