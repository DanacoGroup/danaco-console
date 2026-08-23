// Odpowiedzialność pliku: nazwane konteksty pamięci, kontekst czynny karty
// sesji oraz zasady retencji (tabele `kontekst_pamieci`,
// `kontekst_pamieci_czynny`, `zasada_retencji_pamieci`, migracja 294) —
// dopełnienie rodziny `memory.*` obok `workspace_pamiec.go`.
//
// Kontekst jest zestawem WSKAZAŃ: usunięcie kontekstu nie kasuje ani jednego
// wpisu pamięci, bo kontekst nie jest właścicielem treści. Dlatego wpisy idą
// zapisem strukturalnym w kolumnie, a nie kluczem obcym z kaskadą.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// KontekstPamieci to wiersz tabeli `kontekst_pamieci`.
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

// ZasadaRetencjiPamieci to wiersz tabeli `zasada_retencji_pamieci`.
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

// RepozytoriumKontekstowPamieci jest kontraktem kontekstów i zasad retencji.
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

	zapiszKontekstPamieci = `INSERT INTO kontekst_pamieci (` + kolumnyKontekstuPamieci + `)
	                         VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	                         ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                             nazwa = excluded.nazwa,
	                             opis = excluded.opis,
	                             profil_kod = excluded.profil_kod,
	                             poziomy_json = excluded.poziomy_json,
	                             wpisy_json = excluded.wpisy_json,
	                             prompt_systemowy = excluded.prompt_systemowy,
	                             czynny = excluded.czynny,
	                             zaktualizowano = excluded.zaktualizowano`

	pobierzKontekstPamieci = `SELECT ` + kolumnyKontekstuPamieci +
		` FROM kontekst_pamieci WHERE identyfikator_zewnetrzny = ?`

	usunKontekstPamieci = `DELETE FROM kontekst_pamieci WHERE identyfikator_zewnetrzny = ?`

	uaktywnijKontekstPamieci = `INSERT INTO kontekst_pamieci_czynny (sesja_kod, kontekst_kod, uaktywniono)
	                            VALUES (?, ?, ?)
	                            ON CONFLICT(sesja_kod) DO UPDATE SET
	                                kontekst_kod = excluded.kontekst_kod,
	                                uaktywniono = excluded.uaktywniono`

	pobierzCzynnyKontekstPamieci = `SELECT kontekst_kod FROM kontekst_pamieci_czynny WHERE sesja_kod = ?`

	kolumnyZasadyRetencji = `identyfikator_zewnetrzny, zasieg, zasieg_kod, profil_kod,
	                         dni_wygasania, wrazliwe_domyslnie, wzorce_json, czynna, zaktualizowano`

	zapiszZasadeRetencji = `INSERT INTO zasada_retencji_pamieci (` + kolumnyZasadyRetencji + `)
	                        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	                        ON CONFLICT(zasieg, zasieg_kod, profil_kod) DO UPDATE SET
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

// noweRepozytoriumKontekstowPamieci zakłada magazyn kontekstów nad bazą zestawu.
func noweRepozytoriumKontekstowPamieci(z *zapytania, db *sql.DB) *repozytoriumKontekstowPamieci {
	return &repozytoriumKontekstowPamieci{zapytania: z, db: db}
}

// ZapiszKontekstPamieci zakłada kontekst albo nadpisuje zastany po kodzie.
func (r *repozytoriumKontekstowPamieci) ZapiszKontekstPamieci(ctx context.Context,
	kontekst KontekstPamieci) (KontekstPamieci, error) {

	if strings.TrimSpace(kontekst.Kod) == "" {
		return KontekstPamieci{}, fmt.Errorf("dane: kontekst pamięci bez identyfikatora")
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszKontekstPamieci)
	if err != nil {
		return KontekstPamieci{}, err
	}
	_, err = polecenie.ExecContext(ctx, kontekst.Kod, kontekst.Nazwa, tekstDoKolumny(kontekst.Opis),
		kontekst.ProfilKod, kontekst.PoziomyJSON, kontekst.WpisyJSON,
		tekstDoKolumny(kontekst.PromptSystemowy), liczbaLogiczna(kontekst.Czynny),
		kontekst.Utworzono, kontekst.Zaktualizowano)
	if err != nil {
		return KontekstPamieci{}, fmt.Errorf("dane: nie można zapisać kontekstu pamięci %q: %w",
			kontekst.Kod, err)
	}
	return r.KontekstPamieciPoKodzie(ctx, kontekst.Kod)
}

// KontekstPamieciPoKodzie zwraca jeden kontekst.
func (r *repozytoriumKontekstowPamieci) KontekstPamieciPoKodzie(ctx context.Context,
	kod string) (KontekstPamieci, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, pobierzKontekstPamieci)
	if err != nil {
		return KontekstPamieci{}, err
	}
	kontekst, err := odczytajKontekstPamieci(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return KontekstPamieci{}, fmt.Errorf("dane: kontekst pamięci %q nie istnieje: %w",
			kod, ErrBrakWiersza)
	}
	return kontekst, err
}

// KontekstyPamieci zwraca konteksty w kolejności wyświetlania.
func (r *repozytoriumKontekstowPamieci) KontekstyPamieci(ctx context.Context, profil string,
	zWylaczonymi bool) ([]KontekstPamieci, error) {

	warunki := []string{"1 = 1"}
	argumenty := []any{}
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

// UsunKontekstPamieci kasuje wskazanie. Wpisów pamięci nie rusza.
func (r *repozytoriumKontekstowPamieci) UsunKontekstPamieci(ctx context.Context,
	kod string) (bool, error) {

	polecenie, err := r.zapytania.przygotuj(ctx, usunKontekstPamieci)
	if err != nil {
		return false, err
	}
	wynik, err := polecenie.ExecContext(ctx, kod)
	if err != nil {
		return false, fmt.Errorf("dane: nie można usunąć kontekstu pamięci %q: %w", kod, err)
	}
	usuniete, err := wynik.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("dane: nie można policzyć usuniętych kontekstów: %w", err)
	}
	return usuniete > 0, nil
}

// UaktywnijKontekstPamieci wskazuje kontekst czynny karty sesji.
func (r *repozytoriumKontekstowPamieci) UaktywnijKontekstPamieci(ctx context.Context,
	sesja, kontekst string, chwila int64) error {

	polecenie, err := r.zapytania.przygotuj(ctx, uaktywnijKontekstPamieci)
	if err != nil {
		return err
	}
	if _, err := polecenie.ExecContext(ctx, sesja, kontekst, chwila); err != nil {
		return fmt.Errorf("dane: nie można uaktywnić kontekstu pamięci %q w karcie %q: %w",
			kontekst, sesja, err)
	}
	return nil
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
	err = polecenie.QueryRowContext(ctx, sesja).Scan(&kod)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("dane: nie można odczytać kontekstu czynnego karty %q: %w", sesja, err)
	}
	return kod, nil
}

// ZapiszZasadeRetencjiPamieci zakłada zasadę albo nadpisuje zastaną w tym samym
// zasięgu.
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

// LiczbaWpisowPamieciProfilu liczy wpisy pamięci założone nie później niż
// wskazana chwila — czyli te ZASTANE, których zasada dotknie przy najbliższym
// wygaszaniu. Kontrakt `memory.retention.set` oddaje tę liczbę wprost i nie ma
// prawa jej zgadywać: pochodzi z policzenia wierszy, nie z oszacowania.
//
// Granica jest znacznikiem czasu w zapisie kolumny `utworzono`
// (`wpis_pamieci_projektu`, migracja 035) — czyli ISO-8601 w UTC. Porównanie
// napisów jest tu poprawne, bo ten zapis rośnie leksykalnie razem z czasem.
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

// odczytajKontekstPamieci przekłada wiersz na kontekst.
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

// odczytajZasadeRetencjiPamieci przekłada wiersz na zasadę retencji.
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
