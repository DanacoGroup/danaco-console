// Odpowiedzialność pliku: podstawa modułu Translate — okno tłumaczenia (tekst
// źródłowy, tabela `okno_tlumaczenia`) wraz z deklaracją całego kontraktu
// modułu. Panel tłumaczenia leży w `tlumaczenie_panele.go`, zmiany treści panelu
// w `tlumaczenie_tresc.go`, słownik w `slownik*.go`, jakość i mowa
// w `jakosc*.go` — jedno repozytorium, dziewięć plików wedle odpowiedzialności.
//
// Interfejs deklaruje wyłącznie ten plik, w całości — wraz z metodami
// implementowanymi w pozostałych plikach. Interfejs rozdzielony na dziewięć
// plików byłby dziewięcioma prawdami o jednym kontrakcie.
//
// Rdzeń nie rozpoznaje języka i nie tłumaczy: `jezyk_zrodlowy` bywa NULL, dopóki
// Operator albo `source.set` go nie poda — nie dorabiamy tu wartości domyślnej
// udającej rozpoznanie. Treść panelu bywa pusta z tego samego powodu.
//
// Czas jest liczbą (ms epoki), wzorem `dane/asystent.go`, nie tekstem.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// OknoTlumaczenia to wiersz tabeli `okno_tlumaczenia` — tekst źródłowy i jego
// rozpoznany (albo nie) język. Treść źródłowa bywa obszerna → plik na dysku,
// baza trzyma odwołanie, wzór `dane/wiadomosci.go`: `TekstZrodlowy`
// niesie treść wprost, gdy jest krótka, `TekstZrodlowyOdwolanie` — odwołanie
// do pliku, gdy jest obszerna. Rdzeń nie rozstrzyga, które pole wypełnić;
// zapisuje to, co przyszło z wyższej warstwy.
type OknoTlumaczenia struct {
	ID                     int64
	Kod                    string
	TekstZrodlowy          *string
	TekstZrodlowyOdwolanie *string
	JezykZrodlowy          *string
	LiczbaSegmentow        *int64
	Utworzono              int64
	Zaktualizowano         int64
}

// RepozytoriumTlumaczen jest kontraktem całego modułu Translate — jedno
// repozytorium rozłożone na dziewięć plików wedle odpowiedzialności.
type RepozytoriumTlumaczen interface {
	// --- okno źródłowe ---
	ZapiszOkno(ctx context.Context, okno OknoTlumaczenia) (OknoTlumaczenia, error)
	Okno(ctx context.Context, kod string) (OknoTlumaczenia, error)

	// --- panele ---
	ZapiszPanel(ctx context.Context, oknoID int64, panel PanelTlumaczenia) (PanelTlumaczenia, error)
	Panel(ctx context.Context, kod string) (PanelTlumaczenia, error)
	Panele(ctx context.Context, oknoID int64) ([]PanelTlumaczenia, error)
	// WszystkiePanele daje zakres komendom, które go nie dostają w żądaniu:
	// `glossary.apply` z pustym panelem i `glossary.occurrences`.
	WszystkiePanele(ctx context.Context) ([]PanelTlumaczenia, error)

	// --- zmiany treści panelu ---
	UstawTlumaczenie(ctx context.Context, kodPanelu string, tresc, odwolanie *string) (PanelTlumaczenia, error)
	UstawTlumaczenieZwrotne(ctx context.Context, kodPanelu, tresc string) (PanelTlumaczenia, error)
	UstawTon(ctx context.Context, kodPanelu, ton string) (PanelTlumaczenia, error)
	UstawStanPanelu(ctx context.Context, kodPanelu, stan string) (PanelTlumaczenia, error)

	// --- słownik: terminy ---
	ZapiszTerminy(ctx context.Context, terminy []TerminSlownika) ([]TerminSlownika, error)
	Terminy(ctx context.Context) ([]TerminSlownika, error)
	Termin(ctx context.Context, kod string) (TerminSlownika, error)

	// --- słownik: ślady wymiany (dwie tabele, nie jedna) ---
	ZapiszImportSlownika(ctx context.Context, slad SladImportuSlownika) (SladImportuSlownika, error)
	ZapiszEksportSlownika(ctx context.Context, slad SladEksportuSlownika) (SladEksportuSlownika, error)
	ImportySlownika(ctx context.Context, limit int) ([]SladImportuSlownika, error)
	EksportySlownika(ctx context.Context, limit int) ([]SladEksportuSlownika, error)

	// --- słownik: pamięć tłumaczeń ---
	ZapiszPamiec(ctx context.Context, panelID int64, wpis WpisPamieciTlumaczen) (WpisPamieciTlumaczen, error)
	Podpowiedzi(ctx context.Context, jezyk, fraza string, limit int) ([]WpisPamieciTlumaczen, error)

	// --- jakość: niezgodności ---
	ZapiszNiezgodnosci(ctx context.Context, panelID int64, niezgodnosci []NiezgodnoscTlumaczenia) error
	Niezgodnosci(ctx context.Context, panelID int64) ([]NiezgodnoscTlumaczenia, error)

	// --- jakość: synteza mowy ---
	ZapiszSyntezeMowy(ctx context.Context, synteza SyntezaMowy) (SyntezaMowy, error)
	Syntezy(ctx context.Context, panelID int64) ([]SyntezaMowy, error)

	// --- jakość: eksport panelu ---
	ZapiszEksportPanelu(ctx context.Context, eksport EksportPanelu) (EksportPanelu, error)
	EksportyPanelu(ctx context.Context, panelID int64) ([]EksportPanelu, error)

	// --- pamięć tłumaczeń jako byt Operatora (`tlumaczenie_pamiec.go`) ---
	WpisyPamieci(ctx context.Context, filtr FiltrPamieciTlumaczen) ([]WpisPamieciTlumaczenPelny, int, error)
	WpisPamieci(ctx context.Context, kod string) (WpisPamieciTlumaczenPelny, error)
	ZapiszWpisPamieci(ctx context.Context, wpis WpisPamieciTlumaczenPelny) (WpisPamieciTlumaczenPelny, error)
	UsunWpisPamieci(ctx context.Context, kod string) (bool, error)
	UsunWpisyPamieci(ctx context.Context, kody []string) (int, error)
	PolitykaPamieci(ctx context.Context, oknoID int64) (PolitykaPamieciOkna, error)
	ZapiszPolitykePamieci(ctx context.Context, polityka PolitykaPamieciOkna) (PolitykaPamieciOkna, error)

	// --- segmentacja (`tlumaczenie_segmentacja.go`) ---
	ZestawyRegulSegmentacji(ctx context.Context, jezyk string) ([]ZestawRegulSegmentacji, error)
	ZestawRegulSegmentacji(ctx context.Context, kod string) (ZestawRegulSegmentacji, error)
	ZapiszZestawRegulSegmentacji(ctx context.Context, zestaw ZestawRegulSegmentacji) (ZestawRegulSegmentacji, error)
	SegmentyOkna(ctx context.Context, oknoID int64) ([]string, error)
	UstawSegmentyOkna(ctx context.Context, oknoID int64, segmenty []string) error

	// --- kontrola: terminy zawężone, profile QA, obieg, korekta
	//     (`tlumaczenie_kontrola.go`) ---
	TerminyZawezone(ctx context.Context, filtr FiltrTerminow) ([]TerminSlownika, int, error)
	ProfileQa(ctx context.Context, zasieg, zasiegID string) ([]ProfilQa, error)
	ProfilQaPoKodzie(ctx context.Context, kod string) (ProfilQa, error)
	ZapiszProfilQa(ctx context.Context, profil ProfilQa) (ProfilQa, error)
	UsunProfilQa(ctx context.Context, kod string) (bool, error)
	ZapiszZatwierdzenie(ctx context.Context, zapis ZatwierdzeniePanelu) (ZatwierdzeniePanelu, error)
	Zatwierdzenia(ctx context.Context, panelID, oknoID int64) ([]ZatwierdzeniePanelu, error)
	ZapiszUstaleniaKorekty(ctx context.Context, panelID int64, ustalenia []UstalenieKorekty) error
	UstaleniaKorekty(ctx context.Context, panelID int64) ([]UstalenieKorekty, error)
	UstalenieKorektyPoKodzie(ctx context.Context, kod string) (UstalenieKorekty, error)
	RozstrzygnijUstalenieKorekty(ctx context.Context, kod string, odrzucone bool) error

	// --- materiał: dokument, zasób lokalizacyjny, napisy
	//     (`tlumaczenie_dokument.go`) ---
	ZapiszDokument(ctx context.Context, dokument DokumentTlumaczenia, segmenty []SegmentDokumentu) (DokumentTlumaczenia, error)
	Dokument(ctx context.Context, kod string) (DokumentTlumaczenia, error)
	SegmentyDokumentu(ctx context.Context, dokumentID int64) ([]SegmentDokumentu, error)
	ZapiszZasobLokalizacji(ctx context.Context, zasob ZasobLokalizacji, klucze []KluczLokalizacji) (ZasobLokalizacji, error)
	ZasobLokalizacji(ctx context.Context, kod string) (ZasobLokalizacji, error)
	KluczeLokalizacji(ctx context.Context, zasobID int64) ([]KluczLokalizacji, error)
	ZapiszKluczLokalizacji(ctx context.Context, zasobID int64, klucz KluczLokalizacji) error
	ZapiszKwestieNapisow(ctx context.Context, oknoID, panelID int64, kwestie []KwestiaNapisow) error
	KwestieNapisow(ctx context.Context, oknoID, panelID int64) ([]KwestiaNapisow, error)

	// --- silniki i wymiana zewnętrzna (`tlumaczenie_wymiana_zewnetrzna.go`) ---
	ProfileSilnikow(ctx context.Context, zasieg, zasiegID string) ([]ProfilSilnika, error)
	ProfilSilnikaPoKodzie(ctx context.Context, kod string) (ProfilSilnika, error)
	ZapiszProfilSilnika(ctx context.Context, profil ProfilSilnika) (ProfilSilnika, error)
	PolitykaPivotaZasiegu(ctx context.Context, zasieg, zasiegID string) (PolitykaPivota, error)
	ZapiszPolitykePivota(ctx context.Context, polityka PolitykaPivota) (PolitykaPivota, error)
	ZapiszPakietPrzekazania(ctx context.Context, pakiet PakietPrzekazania) (PakietPrzekazania, error)
	PakietPrzekazania(ctx context.Context, kod string) (PakietPrzekazania, error)
	ZalozZleceniePakietu(ctx context.Context, kod string, oknoID int64, pozycje []PozycjaPakietuTlumaczenia) error
	PozycjePakietu(ctx context.Context, kod string) ([]PozycjaPakietuTlumaczenia, error)
	ZapiszMost(ctx context.Context, most MostTlumaczenia) error
	Most(ctx context.Context, oknoID int64) (MostTlumaczenia, error)
}

const (
	kolumnyOknaTlumaczenia = `id, identyfikator_zewnetrzny, tekst_zrodlowy, tekst_zrodlowy_odwolanie,
	                          jezyk_zrodlowy, liczba_segmentow, utworzono, zaktualizowano`

	zapiszOknoTlumaczenia = `INSERT INTO okno_tlumaczenia
	                         (identyfikator_zewnetrzny, tekst_zrodlowy, tekst_zrodlowy_odwolanie,
	                          jezyk_zrodlowy, liczba_segmentow, utworzono, zaktualizowano)
	                         VALUES (?, ?, ?, ?, ?, ?, ?)
	                         ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                             tekst_zrodlowy = excluded.tekst_zrodlowy,
	                             tekst_zrodlowy_odwolanie = excluded.tekst_zrodlowy_odwolanie,
	                             jezyk_zrodlowy = excluded.jezyk_zrodlowy,
	                             liczba_segmentow = excluded.liczba_segmentow,
	                             zaktualizowano = excluded.zaktualizowano`

	pobierzOknoTlumaczenia = `SELECT ` + kolumnyOknaTlumaczenia + ` FROM okno_tlumaczenia
	                          WHERE identyfikator_zewnetrzny = ?`
)

type repozytoriumTlumaczen struct {
	zapytania *zapytania
	db        *sql.DB
}

func noweRepozytoriumTlumaczen(z *zapytania, db *sql.DB) *repozytoriumTlumaczen {
	return &repozytoriumTlumaczen{zapytania: z, db: db}
}

// ZapiszOkno zakłada wiersz okna tłumaczenia albo nadpisuje zastany po kodzie
// zewnętrznym. `Utworzono` nie wchodzi do klauzuli UPDATE — zapis powtórny nie
// ma prawa przesunąć chwili założenia okna, tylko chwilę ostatniej zmiany.
// Nie ustalamy tu języka źródłowego, gdy przyszedł pusty — rdzeń nie rozpoznaje
// języka; NULL zostaje NULL.
func (r *repozytoriumTlumaczen) ZapiszOkno(ctx context.Context, okno OknoTlumaczenia) (OknoTlumaczenia, error) {
	if okno.Kod == "" {
		return OknoTlumaczenia{}, fmt.Errorf("dane: okno tłumaczenia bez identyfikatora")
	}
	teraz := time.Now().UnixMilli()
	utworzono := okno.Utworzono
	if utworzono == 0 {
		utworzono = teraz
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszOknoTlumaczenia)
	if err != nil {
		return OknoTlumaczenia{}, err
	}
	_, err = polecenie.ExecContext(ctx, okno.Kod, tekstDoKolumny(okno.TekstZrodlowy),
		tekstDoKolumny(okno.TekstZrodlowyOdwolanie), tekstDoKolumny(okno.JezykZrodlowy),
		liczbaDoKolumny(okno.LiczbaSegmentow), utworzono, teraz)
	if err != nil {
		return OknoTlumaczenia{}, fmt.Errorf("dane: nie można zapisać okna tłumaczenia %q: %w", okno.Kod, err)
	}
	return r.Okno(ctx, okno.Kod)
}

// Okno zwraca okno tłumaczenia o wskazanym kodzie. Brak wiersza wraca jako
// ErrBrakWiersza — warstwa wyższa odróżnia „nie ma” od „odczyt się nie powiódł”.
func (r *repozytoriumTlumaczen) Okno(ctx context.Context, kod string) (OknoTlumaczenia, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzOknoTlumaczenia)
	if err != nil {
		return OknoTlumaczenia{}, err
	}
	okno, err := odczytajOknoTlumaczenia(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return OknoTlumaczenia{}, ErrBrakWiersza
	}
	if err != nil {
		return OknoTlumaczenia{}, fmt.Errorf("dane: nieczytelny wiersz okna tłumaczenia %q: %w", kod, err)
	}
	return okno, nil
}

// odczytajOknoTlumaczenia składa strukturę z jednego wiersza wyniku.
func odczytajOknoTlumaczenia(wiersz skaner) (OknoTlumaczenia, error) {
	var okno OknoTlumaczenia
	var tekstZrodlowy, tekstZrodlowyOdwolanie, jezykZrodlowy sql.NullString
	var liczbaSegmentow sql.NullInt64
	err := wiersz.Scan(&okno.ID, &okno.Kod, &tekstZrodlowy, &tekstZrodlowyOdwolanie,
		&jezykZrodlowy, &liczbaSegmentow, &okno.Utworzono, &okno.Zaktualizowano)
	if err != nil {
		return OknoTlumaczenia{}, err
	}
	okno.TekstZrodlowy = tekstZKolumny(tekstZrodlowy)
	okno.TekstZrodlowyOdwolanie = tekstZKolumny(tekstZrodlowyOdwolanie)
	okno.JezykZrodlowy = tekstZKolumny(jezykZrodlowy)
	okno.LiczbaSegmentow = liczbaZKolumny(liczbaSegmentow)
	return okno, nil
}
