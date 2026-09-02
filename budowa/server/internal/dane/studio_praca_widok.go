// Odpowiedzialność pliku: nastawy widoku okna pracy z dokumentem modułu Studio, czyli kolumny widoku wiersza nastaw pracy dokumentu.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// NastawaWidokuStudia to kolumny widoku wiersza nastawa_pracy_studio; pola nie są wskaźnikami, bo każda kolumna ma wartość domyślną.
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

// warunekDokumentuNastawy prowadzi do konta przez dokument_studio (konto_id od migracji 484);
// nastawa_pracy_studio nie ma własnego konto_id, a wiersz samego okna nie ma drogi do konta.
const warunekDokumentuNastawy = `EXISTS (SELECT 1 FROM dokument_studio
	WHERE dokument_studio.id = nastawa_pracy_studio.dokument_id AND ` + WarunekKonta + `)`

const (
	kolumnyWidokuStudia = `id, okno, dokument_id, tryb_powierzchni, kierunek_podzialu,
	                       granica_podzialu, tryb_widoku, skala_procent, skala_nastawa,
	                       linijki_widoczne, linijka_jednostka, granica_marginesu,
	                       znaki_formatowania, stron_w_rzedzie, widok_rozkladowki,
	                       przewijanie, podswietlenie_zmian_modelu, przybornik_widoczny`

	widokStudiaPobierzDokument = `SELECT ` + kolumnyWidokuStudia + `
	                              FROM nastawa_pracy_studio
	                              WHERE okno = ? AND dokument_id = ? AND ` + warunekDokumentuNastawy

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
	                     WHERE id = ? AND (dokument_id IS NULL OR ` + warunekDokumentuNastawy + `)`

	widokStudiaIstnieje = `SELECT 1 FROM nastawa_pracy_studio WHERE id = ?`
)

// NastawaWidoku oddaje kolumny widoku wiersza nastaw pracy dla pary okna i dokumentu albo dla samego okna, nie zakładając wiersza.
func (r *repozytoriumStudia) NastawaWidoku(ctx context.Context, okno string,
	dokumentID *int64) (NastawaWidokuStudia, error) {

	if okno == "" {
		return NastawaWidokuStudia{}, fmt.Errorf("dane: nastawy widoku studio bez okna")
	}
	polecenieSQL, argumenty := widokStudiaPobierzOkno, []any{okno}
	if dokumentID != nil {
		polecenieSQL, argumenty = widokStudiaPobierzDokument, []any{okno, *dokumentID, KontoOperatora(ctx)}
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

// ZapiszNastaweWidoku zapisuje kolumny widoku wskazanego wiersza nastaw pracy okna oraz jego dokumentu.
func (r *repozytoriumStudia) ZapiszNastaweWidoku(ctx context.Context,
	nastawa NastawaWidokuStudia) error {

	if nastawa.ID == 0 {
		return fmt.Errorf("dane: zapis nastaw widoku bez wiersza nastaw")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, widokStudiaZapisz)
	if err != nil {
		return err
	}
	wynik, err := polecenie.ExecContext(ctx, nastawa.TrybPowierzchni, nastawa.KierunekPodzialu,
		nastawa.GranicaPodzialu, nastawa.TrybWidoku, nastawa.SkalaProcent, nastawa.SkalaNastawa,
		liczbaLogiczna(nastawa.LinijkiWidoczne), nastawa.LinijkaJednostka,
		liczbaLogiczna(nastawa.GranicaMarginesu), liczbaLogiczna(nastawa.ZnakiFormatowania),
		nastawa.StronWRzedzie, liczbaLogiczna(nastawa.WidokRozkladowki), nastawa.Przewijanie,
		liczbaLogiczna(nastawa.PodswietlenieZmianModelu),
		liczbaLogiczna(nastawa.PrzybornikWidoczny), nastawa.ID, KontoOperatora(ctx))
	if err != nil {
		return fmt.Errorf("dane: nie można zapisać nastaw widoku %d: %w", nastawa.ID, err)
	}
	zmienione, err := wynik.RowsAffected()
	if err != nil {
		return fmt.Errorf("dane: nieznany skutek zapisu nastaw widoku %d: %w", nastawa.ID, err)
	}
	if zmienione > 0 {
		return nil
	}
	var jest int
	err = r.db.QueryRowContext(ctx, widokStudiaIstnieje, nastawa.ID).Scan(&jest)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrBrakWiersza
	}
	if err != nil {
		return fmt.Errorf("dane: nieczytelny wiersz nastaw widoku %d: %w", nastawa.ID, err)
	}
	return fmt.Errorf("dane: nastawa widoku %d należy do innego konta: %w", nastawa.ID, ErrKolizjaWiersza)
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
