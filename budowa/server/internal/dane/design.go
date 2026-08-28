// Plik definiuje obszar promptu strukturalnego modułu Design: tabelę
// prompt_design oraz kontrakt RepozytoriumDesignu, wspólny z
// design_zasoby.go i design_kompozycje.go.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// PromptDesignu to wiersz tabeli prompt_design. Pola opcjonalne kontraktu
// niosą wskaźnik: nil znaczy nieustawione, nie puste.
type PromptDesignu struct {
	ID          int64
	Kod         string
	Temat       string
	Styl        *string
	Kompozycja  *string
	Oswietlenie *string
	Paleta      *string
	Proporcje   *string
	Wykluczenia *string
	Ziarno      *int
	Warianty    *int
	Silnik      *string
	Kreatywnosc *float64
	// Prompt bez okna jest parametrem wywołania; prompt z oknem wchodzi do
	// historii promptów okna.
	Okno      string
	Kanal     *string
	Utworzono string
}

// RepozytoriumDesignu jest kontraktem obszaru Design, obejmującym prompty,
// zasoby, kompozycje, kolekcje, szablony i pozostałe zasoby modułu.
type RepozytoriumDesignu interface {
	// --- prompt ---
	ZapiszPrompt(ctx context.Context, prompt PromptDesignu) (PromptDesignu, error)
	Prompt(ctx context.Context, kod string) (PromptDesignu, error)

	// --- zasoby i etykiety ---
	ZapiszZasob(ctx context.Context, zasob ZasobDesignu) (ZasobDesignu, error)
	Zasoby(ctx context.Context, filtr FiltrZasobow) ([]ZasobDesignu, int, error)
	// Zasob domyka design.asset.tag.set: przekłada kod kontraktu na klucz
	// wiersza z etykietami.
	Zasob(ctx context.Context, kod string) (ZasobDesignu, error)
	UstawEtykietyZasobu(ctx context.Context, zasobID int64, etykiety []string) error
	EtykietyZasobu(ctx context.Context, zasobID int64) ([]string, error)
	// UstawUlubionyZasobu domyka favorite.set osobno, by nie tracić cudzej
	// zmiany przy wyścigu.
	UstawUlubionyZasobu(ctx context.Context, zasobID int64, ulubiony bool) error
	// UsunZasob domyka design.asset.remove i oddaje, czy wiersz istniał,
	// zgodnie z polem Removed.
	UsunZasob(ctx context.Context, kod string) (bool, error)

	// --- kompozycje ---
	ZapiszKompozycje(ctx context.Context, kompozycja KompozycjaDesignu,
		warstwy []WarstwaKompozycji) (KompozycjaDesignu, error)
	Kompozycja(ctx context.Context, kod string) (KompozycjaDesignu, error)
	// Kompozycje domyka design.board.list: odczyt kompozycji okna po kodzie
	// nieznanym klientowi.
	Kompozycje(ctx context.Context, okno string) ([]KompozycjaDesignu, error)
	Warstwy(ctx context.Context, kompozycjaID int64) ([]WarstwaKompozycji, error)

	// --- kolekcje zasobów ---
	ZapiszKolekcjeDesignu(ctx context.Context, kolekcja KolekcjaDesignu) (KolekcjaDesignu, error)
	KolekcjaDesignuPoKodzie(ctx context.Context, kod string) (KolekcjaDesignu, error)
	KolekcjeDesignu(ctx context.Context, okno string, zasob *string) ([]KolekcjaDesignu, error)
	// ZmienPrzypisaniaKolekcjiDesignu oddaje liczbę przypisań naprawdę
	// zmienionych, pole changed.
	ZmienPrzypisaniaKolekcjiDesignu(ctx context.Context, kolekcjaID int64,
		zasoby []string, zdejmij bool) (int, error)
	ZasobyKolekcjiDesignu(ctx context.Context, kolekcjaID int64) ([]string, error)

	// --- szablony promptu i historia promptów ---
	ZapiszSzablonPromptuDesignu(ctx context.Context, szablon SzablonPromptuDesignu) (SzablonPromptuDesignu, error)
	SzablonyPromptuDesignu(ctx context.Context, okno string) ([]SzablonPromptuDesignu, error)
	// PromptyOknaDesignu oddaje prompty wydane w oknie, od najświeższego,
	// wraz z liczbą wszystkich.
	PromptyOknaDesignu(ctx context.Context, okno string, limit int) ([]PromptDesignu, int, error)
	// ZasobyPromptuDesignu oddaje kody zasobów powstałych z promptu, nośnik
	// pola AssetIds kontraktu.
	ZasobyPromptuDesignu(ctx context.Context, promptID int64) ([]string, error)

	// --- wersje kompozycji ---
	ZapiszWersjeKompozycjiDesignu(ctx context.Context, wersja WersjaKompozycjiDesignu,
		warstwy []WarstwaKompozycji) (WersjaKompozycjiDesignu, error)
	WersjeKompozycjiDesignu(ctx context.Context, kompozycjaID int64, limit int) ([]WersjaKompozycjiDesignu, int, error)
	WersjaKompozycjiDesignuPoKodzie(ctx context.Context, kod string) (WersjaKompozycjiDesignu, error)
	WarstwyWersjiKompozycjiDesignu(ctx context.Context, wersjaID int64) ([]WarstwaKompozycji, error)
	// KompozycjaDesignuPoKluczu czyta kompozycję po kluczu wiersza, oddając
	// jej kod zewnętrzny.
	KompozycjaDesignuPoKluczu(ctx context.Context, id int64) (KompozycjaDesignu, error)

	// --- adnotacje kompozycji ---
	ZapiszAdnotacjeDesignu(ctx context.Context, adnotacja AdnotacjaDesignu) (AdnotacjaDesignu, error)
	AdnotacjaDesignuPoKodzie(ctx context.Context, kod string) (AdnotacjaDesignu, error)
	AdnotacjeKompozycjiDesignu(ctx context.Context, kompozycjaID int64, tylkoOtwarte bool) ([]AdnotacjaDesignu, error)

	// --- szablony materiału i profile druku ---
	ZapiszSzablonMaterialuDesignu(ctx context.Context, szablon SzablonMaterialuDesignu,
		warstwy []WarstwaKompozycji) (SzablonMaterialuDesignu, error)
	SzablonMaterialuDesignuPoKodzie(ctx context.Context, kod string) (SzablonMaterialuDesignu, error)
	SzablonyMaterialuDesignu(ctx context.Context, okno, rodzaj string) ([]SzablonMaterialuDesignu, error)
	WarstwySzablonuMaterialuDesignu(ctx context.Context, szablonID int64) ([]WarstwaKompozycji, error)
	// Strony szablonu czynią z niego publikację wielostronicową; bez nich
	// szablon jest jednostronicowy.
	ZapiszStronySzablonuMaterialuDesignu(ctx context.Context, szablonID int64,
		strony []StronaSzablonuMaterialuDesignu, warstwy map[string][]WarstwaKompozycji) error
	StronySzablonuMaterialuDesignu(ctx context.Context, szablonID int64) ([]StronaSzablonuMaterialuDesignu, error)
	WarstwyStronySzablonuMaterialuDesignu(ctx context.Context, stronaID int64) ([]WarstwaKompozycji, error)
	ZapiszProfilDrukuDesignu(ctx context.Context, profil ProfilDrukuDesignu) (ProfilDrukuDesignu, error)
	ProfilDrukuDesignuPoKodzie(ctx context.Context, kod string) (ProfilDrukuDesignu, error)
	ProfileDrukuDesignu(ctx context.Context, okno string) ([]ProfilDrukuDesignu, error)
	// ZapiszLicencjeZasobuDesignu i odczyt licencji domykają
	// design.stock.import dla zasobu zewnętrznego.
	ZapiszLicencjeZasobuDesignu(ctx context.Context, licencja LicencjaZasobuDesignu) error
	LicencjaZasobuDesignuPoZasobie(ctx context.Context, zasobID int64) (LicencjaZasobuDesignu, error)

	// --- zestawy żetonów ---
	ZapiszZestawZetonowDesignu(ctx context.Context, zestaw ZestawZetonowDesignu,
		zetony []ZetonDesignu) (ZestawZetonowDesignu, error)
	ZestawZetonowDesignuPoKodzie(ctx context.Context, kod string) (ZestawZetonowDesignu, error)
	ZestawyZetonowDesignu(ctx context.Context, okno string, kod *string) ([]ZestawZetonowDesignu, error)
	ZetonyZestawuDesignu(ctx context.Context, zestawID int64) ([]ZetonDesignu, error)

	// --- ikony własne ---
	ZapiszIkoneDesignu(ctx context.Context, ikona IkonaDesignu) (IkonaDesignu, error)
	IkonaDesignuPoKodzie(ctx context.Context, kod string) (IkonaDesignu, error)
	IkonyDesignuOkna(ctx context.Context, okno string) ([]IkonaDesignu, error)

	// --- warsztat fotografii ---
	ZapiszCzynnoscFotografiiDesignu(ctx context.Context, czynnosc CzynnoscFotografiiDesignu) error
	CzynnosciFotografiiDesignuZasobu(ctx context.Context, zasobID int64) ([]CzynnoscFotografiiDesignu, error)
	CzynnoscFotografiiDesignuWyniku(ctx context.Context, zasobID int64) (CzynnoscFotografiiDesignu, error)
	ZapiszNastaweFotografiiDesignu(ctx context.Context, nastawa NastawaFotografiiDesignu) (NastawaFotografiiDesignu, error)
	NastawaFotografiiDesignuPoKodzie(ctx context.Context, kod string) (NastawaFotografiiDesignu, error)
	NastawyFotografiiDesignu(ctx context.Context, okno string) ([]NastawaFotografiiDesignu, error)

	// --- ścieżki wektorowe i symbole ---
	ZapiszSciezkeWektorowaDesignu(ctx context.Context,
		sciezka SciezkaWektorowaDesignu) (SciezkaWektorowaDesignu, error)
	SciezkaWektorowaDesignuPoKodzie(ctx context.Context, kod string) (SciezkaWektorowaDesignu, error)
	SciezkiWektoroweDesignu(ctx context.Context, kompozycjaID int64,
		warstwaKod *string) ([]SciezkaWektorowaDesignu, error)
	UsunSciezkeWektorowaDesignu(ctx context.Context, kod string) (bool, error)
	ZapiszSymbolDesignu(ctx context.Context, symbol SymbolDesignu,
		czlonkowie []CzlonekSymbolyDesignu) (SymbolDesignu, error)
	SymbolDesignuPoKodzie(ctx context.Context, kod string) (SymbolDesignu, error)
	SymboleDesignu(ctx context.Context, kompozycjaID int64) ([]SymbolDesignu, error)
	CzlonkowieSymbolyDesignu(ctx context.Context, symbolID int64) ([]CzlonekSymbolyDesignu, error)

	// --- ramki makiety, więzy responsywne i siatki ---
	ZapiszRamkeDesignu(ctx context.Context, ramka RamkaDesignu) (RamkaDesignu, error)
	RamkaDesignuPoKodzie(ctx context.Context, kod string) (RamkaDesignu, error)
	RamkiDesignu(ctx context.Context, kompozycjaID int64) ([]RamkaDesignu, error)
	UsunRamkeDesignu(ctx context.Context, kod string) (bool, error)
	WarstwyRamkiDesignu(ctx context.Context, ramkaID int64) ([]string, error)
	PrzypiszWarstwyDoRamkiDesignu(ctx context.Context, ramkaID int64, warstwy []string) error
	ZwolnijWarstwyRamkiDesignu(ctx context.Context, ramkaID int64) ([]string, error)
	WiezyRamkiDesignu(ctx context.Context, ramkaID int64) ([]WiezRamkiDesignu, error)
	ZapiszWiezyRamkiDesignu(ctx context.Context, ramkaID int64, wiezy []WiezRamkiDesignu) error
	SiatkaKompozycjiDesignu(ctx context.Context, kompozycjaID int64) (string, error)
	ZapiszSiatkeKompozycjiDesignu(ctx context.Context, kompozycjaID int64, siatkaJSON string) error
	// Te trzy metody ruszają jedną warstwę, osobno od zapisu pełnego planszy,
	// by nie gubić cudzej zmiany.
	WarstwaKompozycjiDesignuPoKodzie(ctx context.Context, kod string) (WarstwaKompozycji, error)
	PrzestawWarstweKompozycjiDesignu(ctx context.Context, warstwa WarstwaKompozycji) error
	DolozWarstweKompozycjiDesignu(ctx context.Context, warstwa WarstwaKompozycji) (WarstwaKompozycji, error)

	// --- komponenty i ich instancje (design_komponenty.go) ---
	ZapiszKomponentDesignu(ctx context.Context, komponent KomponentDesignu) (KomponentDesignu, error)
	KomponentDesignuPoKodzie(ctx context.Context, kod string) (KomponentDesignu, error)
	KomponentyDesignu(ctx context.Context, okno string, kod *string) ([]KomponentDesignu, error)
	ZapiszInstancjeKomponentuDesignu(ctx context.Context, instancja InstancjaKomponentuDesignu) error
	InstancjeKomponentuDesignu(ctx context.Context, komponentID int64) ([]InstancjaKomponentuDesignu, error)

	// --- połączenia prototypu (design_prototyp.go) ---
	ZapiszPolaczeniePrototypuDesignu(ctx context.Context,
		polaczenie PolaczeniePrototypuDesignu) (PolaczeniePrototypuDesignu, error)
	PolaczeniePrototypuDesignuPoKodzie(ctx context.Context, kod string) (PolaczeniePrototypuDesignu, error)
	PolaczeniaPrototypuDesignu(ctx context.Context, kompozycjaID int64) ([]PolaczeniePrototypuDesignu, error)
	UsunPolaczeniePrototypuDesignu(ctx context.Context, kod string) (bool, error)

	// --- gradienty wypełnienia ---
	ZapiszGradientDesignu(ctx context.Context, gradient GradientDesignu) (GradientDesignu, error)
	GradientDesignuPoCelu(ctx context.Context, kompozycjaID int64,
		sciezka, warstwa string) (GradientDesignu, error)
}

const (
	kolumnyPromptuDesign = `id, identyfikator_zewnetrzny, temat, styl, kompozycja,
	                        oswietlenie, paleta, proporcje_kadru, wykluczenia, ziarno,
	                        warianty, silnik, kreatywnosc, okno, kanal, utworzono`

	// Zapis zakłada prompt albo nadpisuje zastany po identyfikatorze
	// zewnętrznym — regeneracja wariantów tego samego promptu jest normalną
	// ścieżką, nie usterką.
	zapiszPromptDesign = `INSERT INTO prompt_design
	                      (identyfikator_zewnetrzny, temat, styl, kompozycja, oswietlenie,
	                       paleta, proporcje_kadru, wykluczenia, ziarno, warianty, silnik,
	                       kreatywnosc, okno, kanal)
	                      VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	                      ON CONFLICT(identyfikator_zewnetrzny) DO UPDATE SET
	                          okno = excluded.okno,
	                          kanal = excluded.kanal,
	                          temat = excluded.temat,
	                          styl = excluded.styl,
	                          kompozycja = excluded.kompozycja,
	                          oswietlenie = excluded.oswietlenie,
	                          paleta = excluded.paleta,
	                          proporcje_kadru = excluded.proporcje_kadru,
	                          wykluczenia = excluded.wykluczenia,
	                          ziarno = excluded.ziarno,
	                          warianty = excluded.warianty,
	                          silnik = excluded.silnik,
	                          kreatywnosc = excluded.kreatywnosc`

	pobierzPromptDesign = `SELECT ` + kolumnyPromptuDesign + ` FROM prompt_design
	                       WHERE identyfikator_zewnetrzny = ?`
)

type repozytoriumDesignu struct {
	zapytania *zapytania
	db        *sql.DB
}

// noweRepozytoriumDesignu zakłada repozytorium obszaru Design. Parametr `db`
// służy obszarowi kompozycji — zapis pełnej listy warstw prowadzi transakcję
// (usuń i wstaw od nowa), której prompt nie potrzebuje, bo mieści się
// w jednym poleceniu.
func noweRepozytoriumDesignu(z *zapytania, db *sql.DB) *repozytoriumDesignu {
	return &repozytoriumDesignu{zapytania: z, db: db}
}

// ZapiszPrompt zakłada prompt strukturalny albo nadpisuje zastany po
// identyfikatorze zewnętrznym i zwraca stan po zapisie.
func (r *repozytoriumDesignu) ZapiszPrompt(ctx context.Context,
	prompt PromptDesignu) (PromptDesignu, error) {

	if prompt.Kod == "" {
		return PromptDesignu{}, fmt.Errorf("dane: prompt design bez identyfikatora")
	}
	if prompt.Temat == "" {
		return PromptDesignu{}, fmt.Errorf("dane: prompt design %q bez tematu", prompt.Kod)
	}
	polecenie, err := r.zapytania.przygotuj(ctx, zapiszPromptDesign)
	if err != nil {
		return PromptDesignu{}, err
	}
	// Ziarno i Warianty są *int kontraktu; przekład na *int64 pomocnika
	// drzewa zapisany jest tu wprost.
	var ziarno, warianty any
	if prompt.Ziarno != nil {
		ziarno = int64(*prompt.Ziarno)
	}
	if prompt.Warianty != nil {
		warianty = int64(*prompt.Warianty)
	}
	// Kreatywnosc jest *float64 kontraktu; drzewo nie ma pomocnika zapisu dla
	// kolumn REAL, zapis wprost.
	var kreatywnosc any
	if prompt.Kreatywnosc != nil {
		kreatywnosc = *prompt.Kreatywnosc
	}
	_, err = polecenie.ExecContext(ctx, prompt.Kod, prompt.Temat,
		tekstDoKolumny(prompt.Styl), tekstDoKolumny(prompt.Kompozycja),
		tekstDoKolumny(prompt.Oswietlenie), tekstDoKolumny(prompt.Paleta),
		tekstDoKolumny(prompt.Proporcje), tekstDoKolumny(prompt.Wykluczenia),
		ziarno, warianty, tekstDoKolumny(prompt.Silnik), kreatywnosc,
		prompt.Okno, tekstDoKolumny(prompt.Kanal))
	if err != nil {
		return PromptDesignu{}, fmt.Errorf("dane: nie można zapisać promptu design %q: %w",
			prompt.Kod, err)
	}
	return r.Prompt(ctx, prompt.Kod)
}

// Prompt zwraca prompt strukturalny o wskazanym kodzie. Brak wiersza wraca
// jako ErrBrakWiersza — warstwa wyższa odróżnia „nie ma” od „odczyt się nie
// powiódł”.
func (r *repozytoriumDesignu) Prompt(ctx context.Context, kod string) (PromptDesignu, error) {
	polecenie, err := r.zapytania.przygotuj(ctx, pobierzPromptDesign)
	if err != nil {
		return PromptDesignu{}, err
	}
	prompt, err := odczytajPromptDesign(polecenie.QueryRowContext(ctx, kod))
	if errors.Is(err, sql.ErrNoRows) {
		return PromptDesignu{}, ErrBrakWiersza
	}
	if err != nil {
		return PromptDesignu{}, fmt.Errorf("dane: nieczytelny wiersz promptu design %q: %w", kod, err)
	}
	return prompt, nil
}

// odczytajPromptDesign składa strukturę PromptDesignu z jednego wiersza
// wyniku zapytania, zamieniając kolumny nullowalne na wskaźniki.
func odczytajPromptDesign(wiersz skaner) (PromptDesignu, error) {
	var prompt PromptDesignu
	var styl, kompozycja, oswietlenie, paleta, proporcje, wykluczenia, silnik sql.NullString
	var kanal sql.NullString
	var ziarno, warianty sql.NullInt64
	var kreatywnosc sql.NullFloat64
	err := wiersz.Scan(&prompt.ID, &prompt.Kod, &prompt.Temat, &styl, &kompozycja,
		&oswietlenie, &paleta, &proporcje, &wykluczenia, &ziarno, &warianty, &silnik,
		&kreatywnosc, &prompt.Okno, &kanal, &prompt.Utworzono)
	if err != nil {
		return PromptDesignu{}, err
	}
	prompt.Styl = tekstZKolumny(styl)
	prompt.Kompozycja = tekstZKolumny(kompozycja)
	prompt.Oswietlenie = tekstZKolumny(oswietlenie)
	prompt.Paleta = tekstZKolumny(paleta)
	prompt.Proporcje = tekstZKolumny(proporcje)
	prompt.Wykluczenia = tekstZKolumny(wykluczenia)
	prompt.Silnik = tekstZKolumny(silnik)
	prompt.Kanal = tekstZKolumny(kanal)
	if ziarno.Valid {
		wartosc := int(ziarno.Int64)
		prompt.Ziarno = &wartosc
	}
	if warianty.Valid {
		wartosc := int(warianty.Int64)
		prompt.Warianty = &wartosc
	}
	prompt.Kreatywnosc = liczbaRzeczywistaZKolumny(kreatywnosc)
	return prompt, nil
}
