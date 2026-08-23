// Odpowiedzialność pliku: warstwa danych odcinka kontroli pracy modułu Studio,
// część trzecia — nastawy WIDOKU okna pracy z dokumentem, czyli kolumny widoku
// tabeli `nastawa_pracy_studio` z migracji 366.
//
// ── Dlaczego widok czyta ten plik, a nie `studio_kontrola_pracy.go` ─────────
// Tamten plik zna ten sam wiersz od strony AUTOZAPISU: odstęp, zdarzenia okna,
// wygasanie kopii, skutek ostatniego zapisu. Kolumny widoku (tryb powierzchni,
// skala, linijki, układ stron, przewijanie, podświetlenie zmian wykonawcy,
// przybornik) pytane są przez zupełnie inną parę komend — `studio.view.get`
// i `studio.view.set` — i nigdy razem z nastawami autozapisu. Osobny odczyt
// tych samych wierszy nie zakłada drugiego pojęcia nastawy: tabela jest jedna,
// wiersz zakłada się jedną drogą (`NastawaPracy`), a każda z dwóch grup kolumn
// ma własne polecenie zapisu. Jedno polecenie na obie grupy kazałoby widokowi
// przepisywać nastawy autozapisu, których nie zmieniał.
//
// Pochodzenie fragmentów leży poza tym plikiem (`studio_wejscie_pochodzenie.go`,
// odcinek wejścia): wiersz zakłada ten, kto fragment wnosi, a wykaz czyta ten
// odcinek — tamten odczyt jest pełny i drugiego się nie zakłada.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// NastawaWidokuStudia to kolumny WIDOKU wiersza `nastawa_pracy_studio`.
//
// Wiersz jest ten sam, co dla autozapisu; struktura jest osobna, bo osobne są
// pytania. Pola nie są wskaźnikami: każda kolumna ma w tabeli wartość domyślną
// i nigdy nie jest pusta, więc „brak" tu nie występuje — a wskaźnik sugerowałby,
// że występuje.
type NastawaWidokuStudia struct {
	ID                       int64
	Okno                     string
	DokumentID               *int64
	TrybPowierzchni          string
	KierunekPodzialu         string
	GranicaPodzialu          float64
	TrybWidoku               string
	SkalaProcent             int64
	SkalaNastawa             string
	LinijkiWidoczne          bool
	LinijkaJednostka         string
	GranicaMarginesu         bool
	ZnakiFormatowania        bool
	StronWRzedzie            int64
	WidokRozkladowki         bool
	Przewijanie              string
	PodswietlenieZmianModelu bool
	PrzybornikWidoczny       bool
}

const (
	kolumnyWidokuStudia = `id, okno, dokument_id, tryb_powierzchni, kierunek_podzialu,
	                       granica_podzialu, tryb_widoku, skala_procent, skala_nastawa,
	                       linijki_widoczne, linijka_jednostka, granica_marginesu,
	                       znaki_formatowania, stron_w_rzedzie, widok_rozkladowki,
	                       przewijanie, podswietlenie_zmian_modelu, przybornik_widoczny`

	widokStudiaPobierzDokument = `SELECT ` + kolumnyWidokuStudia + `
	                              FROM nastawa_pracy_studio
	                              WHERE okno = ? AND dokument_id = ?`

	widokStudiaPobierzOkno = `SELECT ` + kolumnyWidokuStudia + `
	                          FROM nastawa_pracy_studio
	                          WHERE okno = ? AND dokument_id IS NULL`

	widokStudiaZapisz = `UPDATE nastawa_pracy_studio SET
	                         tryb_powierzchni = ?, kierunek_podzialu = ?, granica_podzialu = ?,
	                         tryb_widoku = ?, skala_procent = ?, skala_nastawa = ?,
	                         linijki_widoczne = ?, linijka_jednostka = ?, granica_marginesu = ?,
	                         znaki_formatowania = ?, stron_w_rzedzie = ?, widok_rozkladowki = ?,
	                         przewijanie = ?, podswietlenie_zmian_modelu = ?,
	                         przybornik_widoczny = ?,
	                         zaktualizowano = strftime('%Y-%m-%dT%H:%M:%fZ','now')
	                     WHERE id = ?`
)

// NastawaWidoku oddaje kolumny widoku wiersza nastaw pracy dla pary
// (okno, dokument) albo dla samego okna.
//
// Wiersza NIE zakłada: zakłada go `NastawaPracy`, bo to jeden wiersz i dwie
// drogi zakładania rozjechałyby się przy pierwszej zmianie wartości domyślnej.
// Wołający sięga najpierw po `NastawaPracy`, potem po ten odczyt.
func (r *repozytoriumStudia) NastawaWidoku(ctx context.Context, okno string,
	dokumentID *int64) (NastawaWidokuStudia, error) {

	if okno == "" {
		return NastawaWidokuStudia{}, fmt.Errorf("dane: nastawy widoku studio bez okna")
	}
	polecenieSQL, argumenty := widokStudiaPobierzOkno, []any{okno}
	if dokumentID != nil {
		polecenieSQL, argumenty = widokStudiaPobierzDokument, []any{okno, *dokumentID}
	}
	polecenie, err := r.zapytania.przygotuj(ctx, polecenieSQL)
	if err != nil {
		return NastawaWidokuStudia{}, err
	}
	nastawa, err := odczytajNastaweWidokuStudia(polecenie.QueryRowContext(ctx, argumenty...))
	if errors.Is(err, sql.ErrNoRows) {
		return NastawaWidokuStudia{}, ErrBrakWiersza
	}
	if err != nil {
		return NastawaWidokuStudia{}, fmt.Errorf(
			"dane: nieczytelny wiersz nastaw widoku studio okna %q: %w", okno, err)
	}
	return nastawa, nil
}

// ZapiszNastaweWidoku zapisuje kolumny widoku wskazanego wiersza nastaw.
func (r *repozytoriumStudia) ZapiszNastaweWidoku(ctx context.Context,
	nastawa NastawaWidokuStudia) error {

	if nastawa.ID == 0 {
		return fmt.Errorf("dane: zapis nastaw widoku bez wiersza nastaw")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, widokStudiaZapisz)
	if err != nil {
		return err
	}
	_, err = polecenie.ExecContext(ctx, nastawa.TrybPowierzchni, nastawa.KierunekPodzialu,
		nastawa.GranicaPodzialu, nastawa.TrybWidoku, nastawa.SkalaProcent, nastawa.SkalaNastawa,
		liczbaLogiczna(nastawa.LinijkiWidoczne), nastawa.LinijkaJednostka,
		liczbaLogiczna(nastawa.GranicaMarginesu), liczbaLogiczna(nastawa.ZnakiFormatowania),
		nastawa.StronWRzedzie, liczbaLogiczna(nastawa.WidokRozkladowki), nastawa.Przewijanie,
		liczbaLogiczna(nastawa.PodswietlenieZmianModelu),
		liczbaLogiczna(nastawa.PrzybornikWidoczny), nastawa.ID)
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać nastaw widoku %d: %w", nastawa.ID, err)
	}
	return nil
}

func odczytajNastaweWidokuStudia(wiersz skaner) (NastawaWidokuStudia, error) {
	var nastawa NastawaWidokuStudia
	var dokument sql.NullInt64
	var linijki, margines, znaki, rozkladowka, podswietlenie, przybornik int
	err := wiersz.Scan(&nastawa.ID, &nastawa.Okno, &dokument, &nastawa.TrybPowierzchni,
		&nastawa.KierunekPodzialu, &nastawa.GranicaPodzialu, &nastawa.TrybWidoku,
		&nastawa.SkalaProcent, &nastawa.SkalaNastawa, &linijki, &nastawa.LinijkaJednostka,
		&margines, &znaki, &nastawa.StronWRzedzie, &rozkladowka, &nastawa.Przewijanie,
		&podswietlenie, &przybornik)
	if err != nil {
		return NastawaWidokuStudia{}, err
	}
	nastawa.DokumentID = liczbaZKolumny(dokument)
	nastawa.LinijkiWidoczne = linijki == 1
	nastawa.GranicaMarginesu = margines == 1
	nastawa.ZnakiFormatowania = znaki == 1
	nastawa.WidokRozkladowki = rozkladowka == 1
	nastawa.PodswietlenieZmianModelu = podswietlenie == 1
	nastawa.PrzybornikWidoczny = przybornik == 1
	return nastawa, nil
}
