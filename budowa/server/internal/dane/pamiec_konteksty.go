// Odpowiedzialność pliku: nazwane konteksty pamięci, kontekst czynny karty sesji oraz zasady retencji
// (tabele `kontekst_pamieci`, `kontekst_pamieci_czynny`, `zasada_retencji_pamieci`, migracja 294).
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// KontekstPamieci to wiersz tabeli `kontekst_pamieci`, niosący nazwany zestaw wskazań pamięci profilu.
type KontekstPamieci struct {
	Kod             string
	Nazwa           string
	Opis            *string
	ProfilKod       string
	PoziomyJSON     string
	WpisyJSON       string
	PromptSystemowy *string
	Czynny          bool
	Utworzono       int64
	Zaktualizowano  int64
}

// ZasadaRetencjiPamieci to wiersz tabeli `zasada_retencji_pamieci`, niosący próg wygaszania wpisów pamięci.
type ZasadaRetencjiPamieci struct {
	Kod              string
	Zasieg           string
	ZasiegKod        string
	ProfilKod        string
	DniWygasania     int
	WrazliweDomyslne bool
	WzorceJSON       string
	Czynna           bool
	Zaktualizowano   int64
}

// RepozytoriumKontekstowPamieci jest kontraktem kontekstów pamięci i ich zasad retencji dla danego profilu.
type RepozytoriumKontekstowPamieci interface {
	ZapiszKontekstPamieci(ctx context.Context, kontekst KontekstPamieci) (KontekstPamieci, error)
	KontekstPamieciPoKodzie(ctx context.Context, kod string) (KontekstPamieci, error)
	KontekstyPamieci(ctx context.Context, profil string, zWylaczonymi bool) ([]KontekstPamieci, error)
	UsunKontekstPamieci(ctx context.Context, kod string) (bool, error)
	UaktywnijKontekstPamieci(ctx context.Context, sesja, kontekst string, chwila int64) error
	CzynnyKontekstPamieci(ctx context.Context, sesja string) (string, error)
	ZapiszZasadeRetencjiPamieci(ctx context.Context, zasada ZasadaRetencjiPamieci) (ZasadaRetencjiPamieci, error)
	ZasadyRetencjiPamieci(ctx context.Context, zasieg, zasiegKod, profil string) ([]ZasadaRetencjiPamieci, error)
	LiczbaWpisowPamieciProfilu(ctx context.Context, granica string) (int, error)
}

const (
	kolumnyKontekstuPamieci = `identyfikator_zewnetrzny, nazwa, opis, profil_kod, poziomy_json,
	                           wpisy_json, prompt_systemowy, czynny, utworzono, zaktualizowano`

	// Kod kontekstu jest unikalny w całej tabeli, więc kod cudzego konta trafia
	// w konflikt; warunek przy DO UPDATE zostawia wtedy wiersz nietknięty,
	// a zapis kończy się ErrKolizjaWiersza.
	zapiszKontekstPamieci = `INSERT INTO kontekst_pamieci (` + kolumnyKontekstuPamieci + `, konto_id)
	                         VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ` + WskazanieKonta + `)
	                         ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                             nazwa = excluded.nazwa,
	                             opis = excluded.opis,
	                             profil_kod = excluded.profil_kod,
	                             poziomy_json = excluded.poziomy_json,
	                             wpisy_json = excluded.wpisy_json,
	                             prompt_systemowy = excluded.prompt_systemowy,
	                             czynny = excluded.czynny,
	                             zaktualizowano = excluded.zaktualizowano
	                         WHERE ` + WarunekKonta

	pobierzKontekstPamieci = `SELECT ` + kolumnyKontekstuPamieci +
		` FROM kontekst_pamieci WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	usunKontekstPamieci = `DELETE FROM kontekst_pamieci
	                       WHERE identyfikator_zewnetrzny = ? AND ` + WarunekKonta

	// Kod karty sesji jest kluczem głównym całej tabeli, więc karta cudzego konta
	// trafia tu w konflikt. Nadpisywane jest wskazanie zastane, więc granica
	// dochodzi do korzenia przez kod kontekstu stojący w wierszu przed zapisem,
	// nie przez kod wskazywany żądaniem.
	uaktywnijKontekstPamieci = `INSERT INTO kontekst_pamieci_czynny (sesja_kod, kontekst_kod, uaktywniono)
	                            VALUES (?, ?, ?)
	                            ON CONFLICT(sesja_kod) DO UPDATE SET
	                                kontekst_kod = excluded.kontekst_kod,
	                                uaktywniono = excluded.uaktywniono
	                            WHERE EXISTS (SELECT 1 FROM kontekst_pamieci
	                                           WHERE identyfikator_zewnetrzny = kontekst_pamieci_czynny.kontekst_kod
	                                             AND ` + WarunekKonta + `)`

	// Kontekst czynny karty własnej kolumny konta nie ma; granica dochodzi do
	// niego złączeniem z kontekstem, w którym wskazanie konta stoi.
	pobierzCzynnyKontekstPamieci = `SELECT c.kontekst_kod FROM kontekst_pamieci_czynny c
	                                JOIN kontekst_pamieci k
	                                  ON k.identyfikator_zewnetrzny = c.kontekst_kod
	                                WHERE c.sesja_kod = ? AND ` + WarunekKonta

	kolumnyZasadyRetencji = `identyfikator_zewnetrzny, zasieg, zasieg_kod, profil_kod,
	                         dni_wygasania, wrazliwe_domyslnie, wzorce_json, czynna, zaktualizowano`

	zapiszZasadeRetencji = `INSERT INTO zasada_retencji_pamieci (` + kolumnyZasadyRetencji + `)
	                        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	                        ON CONFLICT(zasieg, zasieg_kod, profil_kod,
	                                    COALESCE(konto_id, 0)) DO UPDATE SET
	                            dni_wygasania = excluded.dni_wygasania,
	                            wrazliwe_domyslnie = excluded.wrazliwe_domyslnie,
	                            wzorce_json = excluded.wzorce_json,
	                            czynna = excluded.czynna,
	                            zaktualizowano = excluded.zaktualizowano`

	pobierzZasadeRetencji = `SELECT ` + kolumnyZasadyRetencji +
		` FROM zasada_retencji_pamieci WHERE zasieg = ? AND zasieg_kod = ? AND profil_kod = ?`
)

type repozytoriumKontekstowPamieci struct {
	zapytania *zapytania
	db        *sql.DB
}

// noweRepozytoriumKontekstowPamieci zakłada magazyn kontekstów pamięci nad bazą danych tego całego zestawu.
func noweRepozytoriumKontekstowPamieci(z *zapytania, db *sql.DB) *repozytoriumKontekstowPamieci {
	return &repozytoriumKontekstowPamieci{zapytania: z, db: db}
}

// ZapiszKontekstPamieci zakłada kontekst pamięci albo nadpisuje zastany wiersz po jego unikalnym kodzie.
func (r *repozytoriumKontekstowPamieci) ZapiszKontekstPamieci(ctx context.Context,
	kontekst KontekstPamieci) (KontekstPamieci, error) {

	if strings.TrimSpace(kontekst.Kod) == "" {
		return KontekstPamieci{}, fmt.Errorf("dane: kontekst pamięci bez identyfikatora")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszKontekstPamieci)
	if err != nil {
		return KontekstPamieci{}, err
	}
	wynik, err := polecenie.ExecContext(ctx, kontekst.Kod, kontekst.Nazwa, tekstDoKolumny(kontekst.Opis),
		kontekst.ProfilKod, kontekst.PoziomyJSON, kontekst.WpisyJSON,
		tekstDoKolumny(kontekst.PromptSystemowy), liczbaLogiczna(kontekst.Czynny),
		kontekst.Utworzono, kontekst.Zaktualizowano,
		KontoOperatora(ctx), KontoOperatora(ctx))
	if err != nil {
		return KontekstPamieci{}, fmt.Errorf("dane: nie można zapisać kontekstu pamięci %q: %w",
			kontekst.Kod, err)
	}
	if err := sprawdzTrafienieZapisu(wynik, "kontekst pamięci", kontekst.Kod); err != nil {
		return KontekstPamieci{}, err
	}
	return r.KontekstPamieciPoKodzie(ctx, kontekst.Kod)
}

// KontekstPamieciPoKodzie zwraca jeden kontekst pamięci wskazany jego unikalnym kodem tekstowym w bazie.
func (r *repozytoriumKontekstowPamieci) KontekstPamieciPoKodzie(ctx context.Context,
	kod string) (KontekstPamieci, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzKontekstPamieci)
	if err != nil {
		return KontekstPamieci{}, err
	}
	kontekst, err := odczytajKontekstPamieci(polecenie.QueryRowContext(ctx, kod, KontoOperatora(ctx)))
	if errors.Is(err, sql.ErrNoRows) {
		return KontekstPamieci{}, fmt.Errorf("dane: kontekst pamięci %q nie istnieje: %w",
			kod, ErrBrakWiersza)
	}
	return kontekst, err
}

// KontekstyPamieci zwraca konteksty pamięci danego profilu w kolejności ustalonej do ich wyświetlania.
func (r *repozytoriumKontekstowPamieci) KontekstyPamieci(ctx context.Context, profil string,
	zWylaczonymi bool) ([]KontekstPamieci, error) {

	warunki := []string{WarunekKonta}
	argumenty := []any{KontoOperatora(ctx)}
	if strings.TrimSpace(profil) != "" {
		warunki = append(warunki, "profil_kod = ?")
		argumenty = append(argumenty, profil)
	}
	if !zWylaczonymi {
		warunki = append(warunki, "czynny = 1")
	}
	tekst := `SELECT ` + kolumnyKontekstuPamieci + ` FROM kontekst_pamieci WHERE ` +
		strings.Join(warunki, " AND ") + ` ORDER BY nazwa, identyfikator_zewnetrzny`

	wiersze, err := r.db.QueryContext(ctx, tekst, argumenty...)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać kontekstów pamięci: %w", err)
	}
	defer wiersze.Close()

	lista := []KontekstPamieci{}
	for wiersze.Next() {
		kontekst, err := odczytajKontekstPamieci(wiersze)
		if err != nil {
			return nil, err
		}
		lista = append(lista, kontekst)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt kontekstów pamięci: %w", err)
	}
	return lista, nil
}

// UsunKontekstPamieci kasuje jedynie wskazanie kontekstu pamięci; wpisów samej pamięci w ogóle nie rusza.
func (r *repozytoriumKontekstowPamieci) UsunKontekstPamieci(ctx context.Context,
	kod string) (bool, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, usunKontekstPamieci)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, kod, KontoOperatora(ctx))
	if err != nil {
		return false, fmt.Errorf("dane: nie można usunąć kontekstu pamięci %q: %w", kod, err)
	}
	usuniete, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nie można policzyć usuniętych kontekstów: %w", err)
	}
	return usuniete > 0, nil
}

// UaktywnijKontekstPamieci wskazuje kontekst pamięci jako czynny dla wskazanej karty danej sesji rozmowy.
func (r *repozytoriumKontekstowPamieci) UaktywnijKontekstPamieci(ctx context.Context,
	sesja, kontekst string, chwila int64) error {

	polecenie, err := r.zapytania.przygotuj(ctx, uaktywnijKontekstPamieci)
	if err != nil {
		return err
	}
	wynik, err := polecenie.ExecContext(ctx, sesja, kontekst, chwila, KontoOperatora(ctx))
	if err != nil {
		return fmt.Errorf("dane: nie można uaktywnić kontekstu pamięci %q w karcie %q: %w",
			kontekst, sesja, err)
	}
	return sprawdzTrafienieZapisu(wynik, "kontekst czynny karty", sesja)
}

// CzynnyKontekstPamieci zwraca kontekst czynny karty sesji. Brak wskazania nie
// jest błędem — oddaje pusty napis.
func (r *repozytoriumKontekstowPamieci) CzynnyKontekstPamieci(ctx context.Context,
	sesja string) (string, error) {

	if strings.TrimSpace(sesja) == "" {
		return "", nil
	}
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzCzynnyKontekstPamieci)
	if err != nil {
		return "", err
	}
	var kod string
	err = polecenie.QueryRowContext(ctx, sesja, KontoOperatora(ctx)).Scan(&kod)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("dane: nie można odczytać kontekstu czynnego karty %q: %w", sesja, err)
	}
	return kod, nil
}

// ZapiszZasadeRetencjiPamieci zakłada zasadę retencji pamięci albo nadpisuje zastaną w tym samym zasięgu.
func (r *repozytoriumKontekstowPamieci) ZapiszZasadeRetencjiPamieci(ctx context.Context,
	zasada ZasadaRetencjiPamieci) (ZasadaRetencjiPamieci, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, zapiszZasadeRetencji)
	if err != nil {
		return ZasadaRetencjiPamieci{}, err
	}
	_, err = polecenie.ExecContext(ctx, zasada.Kod, zasada.Zasieg, zasada.ZasiegKod,
		zasada.ProfilKod, zasada.DniWygasania, liczbaLogiczna(zasada.WrazliweDomyslne),
		zasada.WzorceJSON, liczbaLogiczna(zasada.Czynna), zasada.Zaktualizowano)
	if err != nil {
		return ZasadaRetencjiPamieci{}, fmt.Errorf("dane: nie można zapisać zasady retencji pamięci: %w", err)
	}
	odczyt, err := r.zapytania.przygotuj(ctx, pobierzZasadeRetencji)
	if err != nil {
		return ZasadaRetencjiPamieci{}, err
	}
	return odczytajZasadeRetencjiPamieci(odczyt.QueryRowContext(ctx,
		zasada.Zasieg, zasada.ZasiegKod, zasada.ProfilKod))
}

// ZasadyRetencjiPamieci zwraca zasady obowiązujące w zasięgu, od najwęższej.
// Zasięg wskazany wprost stoi na czele, potem zasady szersze — okno dostaje
// materiał do rozstrzygnięcia, a nie jedną wartość wziętą z połowy przesłanek.
func (r *repozytoriumKontekstowPamieci) ZasadyRetencjiPamieci(ctx context.Context,
	zasieg, zasiegKod, profil string) ([]ZasadaRetencjiPamieci, error) {

	warunki := []string{"czynna = 1"}
	argumenty := []any{}
	if strings.TrimSpace(zasieg) != "" {
		warunki = append(warunki, "(zasieg = ? OR zasieg = 'global')")
		argumenty = append(argumenty, zasieg)
	}
	if strings.TrimSpace(zasiegKod) != "" {
		warunki = append(warunki, "(zasieg_kod = ? OR zasieg_kod = '')")
		argumenty = append(argumenty, zasiegKod)
	}
	if strings.TrimSpace(profil) != "" {
		warunki = append(warunki, "(profil_kod = ? OR profil_kod = '')")
		argumenty = append(argumenty, profil)
	}
	tekst := `SELECT ` + kolumnyZasadyRetencji + ` FROM zasada_retencji_pamieci WHERE ` +
		strings.Join(warunki, " AND ") +
		` ORDER BY CASE WHEN profil_kod <> '' THEN 0 ELSE 1 END,
		           CASE WHEN zasieg_kod <> '' THEN 0 ELSE 1 END,
		           CASE WHEN zasieg <> 'global' THEN 0 ELSE 1 END, id`

	wiersze, err := r.db.QueryContext(ctx, tekst, argumenty...)
	if err != nil {
		return nil, fmt.Errorf("dane: nie można odczytać zasad retencji pamięci: %w", err)
	}
	defer wiersze.Close()

	lista := []ZasadaRetencjiPamieci{}
	for wiersze.Next() {
		zasada, err := odczytajZasadeRetencjiPamieci(wiersze)
		if err != nil {
			return nil, err
		}
		lista = append(lista, zasada)
	}
	if err := wiersze.Err(); err != nil {
		return nil, fmt.Errorf("dane: przerwany odczyt zasad retencji pamięci: %w", err)
	}
	return lista, nil
}

// LiczbaWpisowPamieciProfilu liczy wpisy pamięci założone nie później niż wskazana chwila — te zastane,
// których zasada dotknie przy najbliższym wygaszaniu.
func (r *repozytoriumKontekstowPamieci) LiczbaWpisowPamieciProfilu(ctx context.Context,
	granica string) (int, error) {

	liczba := 0
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM wpis_pamieci_projektu WHERE utworzono <= ?`, granica).Scan(&liczba)
	if err != nil {
		return 0, fmt.Errorf("dane: nie można policzyć wpisów pamięci objętych zasadą: %w", err)
	}
	return liczba, nil
}

// odczytajKontekstPamieci przekłada wiersz wyniku zapytania na pełną strukturę kontekstu tej samej pamięci.
func odczytajKontekstPamieci(s skaner) (KontekstPamieci, error) {
	var kontekst KontekstPamieci
	var opis, prompt sql.NullString
	var czynny int
	err := s.Scan(&kontekst.Kod, &kontekst.Nazwa, &opis, &kontekst.ProfilKod,
		&kontekst.PoziomyJSON, &kontekst.WpisyJSON, &prompt, &czynny,
		&kontekst.Utworzono, &kontekst.Zaktualizowano)
	if err != nil {
		return KontekstPamieci{}, err
	}
	kontekst.Opis = tekstZKolumny(opis)
	kontekst.PromptSystemowy = tekstZKolumny(prompt)
	kontekst.Czynny = czynny == 1
	return kontekst, nil
}

// odczytajZasadeRetencjiPamieci przekłada wiersz wyniku zapytania na pełną strukturę tej zasady retencji.
func odczytajZasadeRetencjiPamieci(s skaner) (ZasadaRetencjiPamieci, error) {
	var zasada ZasadaRetencjiPamieci
	var wrazliwe, czynna int
	err := s.Scan(&zasada.Kod, &zasada.Zasieg, &zasada.ZasiegKod, &zasada.ProfilKod,
		&zasada.DniWygasania, &wrazliwe, &zasada.WzorceJSON, &czynna, &zasada.Zaktualizowano)
	if err != nil {
		return ZasadaRetencjiPamieci{}, err
	}
	zasada.WrazliweDomyslne = wrazliwe == 1
	zasada.Czynna = czynna == 1
	return zasada, nil
}
