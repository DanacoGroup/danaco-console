// Odpowiedzialność pliku: materiał wniesiony do tłumaczenia — dokument
// z segmentami (`dokument_tlumaczenia`, `segment_dokumentu_tlumaczenia`,
// migracja 164), zasób lokalizacyjny z kluczami (`zasob_lokalizacji`,
// `klucz_lokalizacji`) i kwestie napisów (`kwestia_napisow`, migracja 165).
//
// Trzy rodzaje materiału, jeden plik, bo wszystkie trzy wchodzą do modułu tą
// samą drogą: plik na dysku → wiersze w bazie → panele przekładu. Różnią się
// wyłącznie tym, co w pliku jest jednostką: akapit, klucz, kwestia.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

// DokumentTlumaczenia to wiersz `dokument_tlumaczenia`.
type DokumentTlumaczenia struct {
	ID             int64
	Kod            string
	OknoID         int64
	OknoKod        string
	Sciezka        string
	Format         string
	LiczbaStron    *int64
	UzytoOcr       bool
	Utworzono      int64
	Zaktualizowano int64
}

// SegmentDokumentu to wiersz `segment_dokumentu_tlumaczenia` — kawałek treści
// razem z miejscem w strukturze pliku.
type SegmentDokumentu struct {
	Kolejnosc    int64
	Tresc        string
	SciezkaWezla *string
	Strona       *int64
	Styl         *string
}

// ZapiszDokument zakłada dokument albo nadpisuje zastany po kodzie i wymienia
// jego segmenty w całości: ponowne wczytanie tego samego pliku jest nowym
// odczytem materiału, nie przyrostem do poprzedniego.
func (r *repozytoriumTlumaczen) ZapiszDokument(ctx context.Context,
	dokument DokumentTlumaczenia, segmenty []SegmentDokumentu) (DokumentTlumaczenia, error) {

	if strings.TrimSpace(dokument.Kod) == "" {
		return DokumentTlumaczenia{}, fmt.Errorf("dane: dokument tłumaczenia bez identyfikatora")
	}
	teraz := time.Now().UnixMilli()
	utworzono := dokument.Utworzono
	if utworzono == 0 {
		utworzono = teraz
	}
	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		if _, err := transakcja.ExecContext(ctx, `INSERT INTO dokument_tlumaczenia
			(identyfikator_zewnetrzny, okno_id, sciezka, format, liczba_stron, uzyto_ocr,
			 utworzono, zaktualizowano)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
				sciezka = excluded.sciezka, format = excluded.format,
				liczba_stron = excluded.liczba_stron, uzyto_ocr = excluded.uzyto_ocr,
				zaktualizowano = excluded.zaktualizowano`,
			dokument.Kod, dokument.OknoID, dokument.Sciezka, dokument.Format,
			liczbaDoKolumny(dokument.LiczbaStron), wartoscLogicznaDoKolumny(dokument.UzytoOcr),
			utworzono, teraz); err != nil {
			return fmt.Errorf("dane: nie można zapisać dokumentu %q: %w", dokument.Kod, err)
		}
		var dokumentID int64
		if err := transakcja.QueryRowContext(ctx,
			`SELECT id FROM dokument_tlumaczenia WHERE identyfikator_zewnetrzny = ?`,
			dokument.Kod).Scan(&dokumentID); err != nil {
			return fmt.Errorf("dane: nie można odczytać dokumentu %q po zapisie: %w", dokument.Kod, err)
		}
		if _, err := transakcja.ExecContext(ctx,
			`DELETE FROM segment_dokumentu_tlumaczenia WHERE dokument_id = ?`, dokumentID); err != nil {
			return fmt.Errorf("dane: nie można zdjąć segmentów dokumentu %q: %w", dokument.Kod, err)
		}
		for numer, segment := range segmenty {
			if _, err := transakcja.ExecContext(ctx, `INSERT INTO segment_dokumentu_tlumaczenia
				(dokument_id, kolejnosc, tresc, sciezka_wezla, strona, styl)
				VALUES (?, ?, ?, ?, ?, ?)`,
				dokumentID, numer, segment.Tresc, tekstDoKolumny(segment.SciezkaWezla),
				liczbaDoKolumny(segment.Strona), tekstDoKolumny(segment.Styl)); err != nil {
				return fmt.Errorf("dane: nie można zapisać segmentu dokumentu %q: %w", dokument.Kod, err)
			}
		}
		return nil
	})
	if err != nil {
		return DokumentTlumaczenia{}, err
	}
	return r.Dokument(ctx, dokument.Kod)
}

// Dokument oddaje dokument po kodzie zewnętrznym.
func (r *repozytoriumTlumaczen) Dokument(ctx context.Context, kod string) (DokumentTlumaczenia, error) {
	var dokument DokumentTlumaczenia
	var strony sql.NullInt64
	var ocr int64
	err := r.db.QueryRowContext(ctx,
		`SELECT d.id, d.identyfikator_zewnetrzny, d.okno_id, o.identyfikator_zewnetrzny,
		        d.sciezka, d.format, d.liczba_stron, d.uzyto_ocr, d.utworzono, d.zaktualizowano
		   FROM dokument_tlumaczenia d
		   JOIN okno_tlumaczenia o ON o.id = d.okno_id
		  WHERE d.identyfikator_zewnetrzny = ?`, kod).
		Scan(&dokument.ID, &dokument.Kod, &dokument.OknoID, &dokument.OknoKod,
			&dokument.Sciezka, &dokument.Format, &strony, &ocr,
			&dokument.Utworzono, &dokument.Zaktualizowano)
	if errors.Is(err, sql.ErrNoRows) {
		return DokumentTlumaczenia{}, ErrBrakWiersza
	}
	if err != nil {
		return DokumentTlumaczenia{}, fmt.Errorf("dane: nieczytelny dokument %q: %w", kod, err)
	}
	dokument.LiczbaStron = liczbaZKolumny(strony)
	dokument.UzytoOcr = ocr == 1
	return dokument, nil
}

// SegmentyDokumentu oddaje segmenty dokumentu w kolejności zapisu.
func (r *repozytoriumTlumaczen) SegmentyDokumentu(ctx context.Context,
	dokumentID int64) ([]SegmentDokumentu, error) {

	wiersze, err := r.db.QueryContext(ctx,
		`SELECT kolejnosc, tresc, sciezka_wezla, strona, styl
		   FROM segment_dokumentu_tlumaczenia WHERE dokument_id = ? ORDER BY kolejnosc`, dokumentID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać segmentów dokumentu %d: %w", dokumentID, err)
	}
	defer wiersze.Close()

	segmenty := []SegmentDokumentu{}
	for wiersze.Next() {
		var segment SegmentDokumentu
		var sciezka, styl sql.NullString
		var strona sql.NullInt64
		if err := wiersze.Scan(&segment.Kolejnosc, &segment.Tresc, &sciezka, &strona, &styl); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny segment dokumentu %d: %w", dokumentID, err)
		}
		segment.SciezkaWezla = tekstZKolumny(sciezka)
		segment.Strona = liczbaZKolumny(strona)
		segment.Styl = tekstZKolumny(styl)
		segmenty = append(segmenty, segment)
	}
	return segmenty, wiersze.Err()
}

// ZasobLokalizacji to wiersz `zasob_lokalizacji`.
type ZasobLokalizacji struct {
	ID             int64
	Kod            string
	OknoID         int64
	OknoKod        string
	Sciezka        string
	Format         string
	JezykZrodlowy  *string
	Utworzono      int64
	Zaktualizowano int64
}

// KluczLokalizacji to wiersz `klucz_lokalizacji`. `Znaczniki` wchodzi i wychodzi
// jako wykaz — postać kolumny (linie) jest sprawą schematu, nie wywołującego.
type KluczLokalizacji struct {
	Klucz        string
	Tresc        string
	Znaczniki    []string
	Kontekst     *string
	ZrzutZasobID *string
	FormyMnogie  *string
	Kolejnosc    int64
}

// ZapiszZasobLokalizacji zakłada zasób albo nadpisuje zastany i wymienia jego
// klucze w całości.
func (r *repozytoriumTlumaczen) ZapiszZasobLokalizacji(ctx context.Context,
	zasob ZasobLokalizacji, klucze []KluczLokalizacji) (ZasobLokalizacji, error) {

	if strings.TrimSpace(zasob.Kod) == "" {
		return ZasobLokalizacji{}, fmt.Errorf("dane: zasób lokalizacyjny bez identyfikatora")
	}
	teraz := time.Now().UnixMilli()
	utworzono := zasob.Utworzono
	if utworzono == 0 {
		utworzono = teraz
	}
	err := wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		if _, err := transakcja.ExecContext(ctx, `INSERT INTO zasob_lokalizacji
			(identyfikator_zewnetrzny, okno_id, sciezka, format, jezyk_zrodlowy, utworzono, zaktualizowano)
			VALUES (?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
				sciezka = excluded.sciezka, format = excluded.format,
				jezyk_zrodlowy = excluded.jezyk_zrodlowy, zaktualizowano = excluded.zaktualizowano`,
			zasob.Kod, zasob.OknoID, zasob.Sciezka, zasob.Format,
			tekstDoKolumny(zasob.JezykZrodlowy), utworzono, teraz); err != nil {
			return fmt.Errorf("dane: nie można zapisać zasobu lokalizacyjnego %q: %w", zasob.Kod, err)
		}
		var zasobID int64
		if err := transakcja.QueryRowContext(ctx,
			`SELECT id FROM zasob_lokalizacji WHERE identyfikator_zewnetrzny = ?`,
			zasob.Kod).Scan(&zasobID); err != nil {
			return fmt.Errorf("dane: nie można odczytać zasobu %q po zapisie: %w", zasob.Kod, err)
		}
		if _, err := transakcja.ExecContext(ctx,
			`DELETE FROM klucz_lokalizacji WHERE zasob_id = ?`, zasobID); err != nil {
			return fmt.Errorf("dane: nie można zdjąć kluczy zasobu %q: %w", zasob.Kod, err)
		}
		for numer, klucz := range klucze {
			if _, err := transakcja.ExecContext(ctx, `INSERT INTO klucz_lokalizacji
				(zasob_id, klucz, tresc, znaczniki, kontekst, zrzut_zasob_id, formy_mnogie, kolejnosc)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
				zasobID, klucz.Klucz, klucz.Tresc, znacznikiDoKolumny(klucz.Znaczniki),
				tekstDoKolumny(klucz.Kontekst), tekstDoKolumny(klucz.ZrzutZasobID),
				tekstDoKolumny(klucz.FormyMnogie), numer); err != nil {
				return fmt.Errorf("dane: nie można zapisać klucza %q: %w", klucz.Klucz, err)
			}
		}
		return nil
	})
	if err != nil {
		return ZasobLokalizacji{}, err
	}
	return r.ZasobLokalizacji(ctx, zasob.Kod)
}

// znacznikiDoKolumny składa wykaz znaczników w jedną kolumnę tekstową.
func znacznikiDoKolumny(znaczniki []string) any {
	if len(znaczniki) == 0 {
		return nil
	}
	return strings.Join(znaczniki, "\n")
}

// ZasobLokalizacji oddaje zasób po kodzie zewnętrznym.
func (r *repozytoriumTlumaczen) ZasobLokalizacji(ctx context.Context,
	kod string) (ZasobLokalizacji, error) {

	var zasob ZasobLokalizacji
	var jezyk sql.NullString
	err := r.db.QueryRowContext(ctx,
		`SELECT z.id, z.identyfikator_zewnetrzny, z.okno_id, o.identyfikator_zewnetrzny,
		        z.sciezka, z.format, z.jezyk_zrodlowy, z.utworzono, z.zaktualizowano
		   FROM zasob_lokalizacji z
		   JOIN okno_tlumaczenia o ON o.id = z.okno_id
		  WHERE z.identyfikator_zewnetrzny = ?`, kod).
		Scan(&zasob.ID, &zasob.Kod, &zasob.OknoID, &zasob.OknoKod, &zasob.Sciezka,
			&zasob.Format, &jezyk, &zasob.Utworzono, &zasob.Zaktualizowano)
	if errors.Is(err, sql.ErrNoRows) {
		return ZasobLokalizacji{}, ErrBrakWiersza
	}
	if err != nil {
		return ZasobLokalizacji{}, fmt.Errorf("dane: nieczytelny zasób lokalizacyjny %q: %w", kod, err)
	}
	zasob.JezykZrodlowy = tekstZKolumny(jezyk)
	return zasob, nil
}

// KluczeLokalizacji oddaje klucze zasobu w kolejności zapisu.
func (r *repozytoriumTlumaczen) KluczeLokalizacji(ctx context.Context,
	zasobID int64) ([]KluczLokalizacji, error) {

	wiersze, err := r.db.QueryContext(ctx,
		`SELECT klucz, tresc, znaczniki, kontekst, zrzut_zasob_id, formy_mnogie, kolejnosc
		   FROM klucz_lokalizacji WHERE zasob_id = ? ORDER BY kolejnosc`, zasobID)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać kluczy zasobu %d: %w", zasobID, err)
	}
	defer wiersze.Close()

	klucze := []KluczLokalizacji{}
	for wiersze.Next() {
		var klucz KluczLokalizacji
		var znaczniki, kontekst, zrzut, formy sql.NullString
		if err := wiersze.Scan(&klucz.Klucz, &klucz.Tresc, &znaczniki, &kontekst,
			&zrzut, &formy, &klucz.Kolejnosc); err != nil {
			return nil, fmt.Errorf("dane: nieczytelny klucz zasobu %d: %w", zasobID, err)
		}
		if znaczniki.Valid && znaczniki.String != "" {
			klucz.Znaczniki = strings.Split(znaczniki.String, "\n")
		}
		klucz.Kontekst = tekstZKolumny(kontekst)
		klucz.ZrzutZasobID = tekstZKolumny(zrzut)
		klucz.FormyMnogie = tekstZKolumny(formy)
		klucze = append(klucze, klucz)
	}
	return klucze, wiersze.Err()
}

// ZapiszKluczLokalizacji nadpisuje jeden klucz zasobu — droga dla
// `resource.key.context.set` i `resource.plural.apply`, które ruszają klucze
// pojedynczo, a nie cały zasób.
func (r *repozytoriumTlumaczen) ZapiszKluczLokalizacji(ctx context.Context,
	zasobID int64, klucz KluczLokalizacji) error {

	_, err := r.db.ExecContext(ctx, `INSERT INTO klucz_lokalizacji
		(zasob_id, klucz, tresc, znaczniki, kontekst, zrzut_zasob_id, formy_mnogie, kolejnosc)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(zasob_id, klucz) DO UPDATE SET
			tresc = excluded.tresc, znaczniki = excluded.znaczniki,
			kontekst = excluded.kontekst, zrzut_zasob_id = excluded.zrzut_zasob_id,
			formy_mnogie = excluded.formy_mnogie`,
		zasobID, klucz.Klucz, klucz.Tresc, znacznikiDoKolumny(klucz.Znaczniki),
		tekstDoKolumny(klucz.Kontekst), tekstDoKolumny(klucz.ZrzutZasobID),
		tekstDoKolumny(klucz.FormyMnogie), klucz.Kolejnosc)
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać klucza %q: %w", klucz.Klucz, err)
	}
	return nil
}

// KwestiaNapisow to wiersz `kwestia_napisow`.
type KwestiaNapisow struct {
	Kolejnosc  int64
	PoczatekMs int64
	KoniecMs   int64
	Tresc      string
	Mowca      *string
}

// ZapiszKwestieNapisow wymienia kwestie okna albo panelu w całości. `panelID`
// równy zeru znaczy kwestie materiału źródłowego, nie przekładu.
func (r *repozytoriumTlumaczen) ZapiszKwestieNapisow(ctx context.Context,
	oknoID, panelID int64, kwestie []KwestiaNapisow) error {

	return wTransakcji(ctx, r.db, func(transakcja *sql.Tx) error {
		warunek := "okno_id = ? AND panel_id IS NULL"
		argumenty := []any{oknoID}
		if panelID > 0 {
			warunek = "panel_id = ?"
			argumenty = []any{panelID}
		}
		if _, err := transakcja.ExecContext(ctx,
			`DELETE FROM kwestia_napisow WHERE `+warunek, argumenty...); err != nil {
			return fmt.Errorf("dane: nie można zdjąć kwestii napisów: %w", err)
		}
		for numer, kwestia := range kwestie {
			var panel any
			if panelID > 0 {
				panel = panelID
			}
			if _, err := transakcja.ExecContext(ctx, `INSERT INTO kwestia_napisow
				(okno_id, panel_id, kolejnosc, poczatek_ms, koniec_ms, tresc, mowca)
				VALUES (?, ?, ?, ?, ?, ?, ?)`,
				oknoID, panel, numer, kwestia.PoczatekMs, kwestia.KoniecMs,
				kwestia.Tresc, tekstDoKolumny(kwestia.Mowca)); err != nil {
				return fmt.Errorf("dane: nie można zapisać kwestii napisów: %w", err)
			}
		}
		return nil
	})
}

// KwestieNapisow oddaje kwestie panelu, a przy `panelID` równym zeru — kwestie
// materiału źródłowego okna.
func (r *repozytoriumTlumaczen) KwestieNapisow(ctx context.Context,
	oknoID, panelID int64) ([]KwestiaNapisow, error) {

	warunek := "okno_id = ? AND panel_id IS NULL"
	argumenty := []any{oknoID}
	if panelID > 0 {
		warunek = "panel_id = ?"
		argumenty = []any{panelID}
	}
	wiersze, err := r.db.QueryContext(ctx,
		`SELECT kolejnosc, poczatek_ms, koniec_ms, tresc, mowca
		   FROM kwestia_napisow WHERE `+warunek+` ORDER BY kolejnosc`, argumenty...)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać kwestii napisów: %w", err)
	}
	defer wiersze.Close()

	kwestie := []KwestiaNapisow{}
	for wiersze.Next() {
		var kwestia KwestiaNapisow
		var mowca sql.NullString
		if err := wiersze.Scan(&kwestia.Kolejnosc, &kwestia.PoczatekMs, &kwestia.KoniecMs,
			&kwestia.Tresc, &mowca); err != nil {
			return nil, fmt.Errorf("dane: nieczytelna kwestia napisów: %w", err)
		}
		kwestia.Mowca = tekstZKolumny(mowca)
		kwestie = append(kwestie, kwestia)
	}
	return kwestie, wiersze.Err()
}
