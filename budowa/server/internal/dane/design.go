// Odpowiedzialność pliku: obszar promptu strukturalnego modułu Design
// (tabela `prompt_design`) wraz z kontraktem całego obszaru. Zasoby i ich
// etykiety (Assets Panel) leżą w `design_zasoby.go`, kompozycje i warstwy
// (Design Board) w `design_kompozycje.go` — jedno repozytorium, trzy pliki
// wedle odpowiedzialności, jak w modułach Automations i Studio. Interfejs
// deklaruje wyłącznie ten plik, w całości — także metody obszarów zasobów
// i kompozycji — żeby kontrakt stał w jednym miejscu.
//
// Prompt ma własną tabelę, bo wiele zasobów i wariantów powstaje z tego samego
// promptu: `prompt_design` jest jedną prawdą o promptcie, a
// `zasob_design.prompt_id` jedynym odwołaniem.
package dane

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// PromptDesignu to wiersz tabeli `prompt_design`. Pola opcjonalne kontraktu
// (`DesignPrompt`) niosą wskaźnik — nil znaczy „nieustawione”, nie „puste”;
// tekstDoKolumny/tekstZKolumny i odpowiedniki liczbowe rozróżniają te dwa
// stany przy przejściu przez kolumnę dopuszczającą NULL.
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
	// Okno i Kanal doszły migracją 232. Prompt bez okna jest parametrem
	// wywołania; prompt z oknem jest zapisem zdarzenia w tym oknie i dopiero
	// taki wchodzi do `design.prompt.history.list`.
	Okno      string
	Kanal     *string
	Utworzono string
}

// RepozytoriumDesignu jest kontraktem obszaru Design.
type RepozytoriumDesignu interface {
	// --- prompt ---
	ZapiszPrompt(ctx context.Context, prompt PromptDesignu) (PromptDesignu, error)
	Prompt(ctx context.Context, kod string) (PromptDesignu, error)

	// --- zasoby i etykiety ---
	ZapiszZasob(ctx context.Context, zasob ZasobDesignu) (ZasobDesignu, error)
	Zasoby(ctx context.Context, filtr FiltrZasobow) ([]ZasobDesignu, int, error)
	// Zasob domyka `design.asset.tag.set`: komenda niesie identyfikator
	// kontraktu, a etykiety wiszą na kluczu wiersza — bez przekładu kod → wiersz
	// nadanie etykiet nie miałoby czego szukać ani czym odmówić.
	Zasob(ctx context.Context, kod string) (ZasobDesignu, error)
	UstawEtykietyZasobu(ctx context.Context, zasobID int64, etykiety []string) error
	EtykietyZasobu(ctx context.Context, zasobID int64) ([]string, error)
	// UstawUlubionyZasobu domyka `design.asset.favorite.set`. Osobna metoda,
	// a nie `ZapiszZasob` z podmienionym polem: pełny zapis wymaga kompletu
	// pól zasobu, więc oznaczenie ulubionego szłoby przez odczyt-zmień-zapisz
	// i przy wyścigu dwóch komend nadpisywałoby cudzą zmianę nazwy czy
	// wymiarów wartością sprzed odczytu.
	UstawUlubionyZasobu(ctx context.Context, zasobID int64, ulubiony bool) error
	// UsunZasob domyka `design.asset.remove` i oddaje informację, czy wiersz
	// istniał — kontrakt (`DesignAssetRemoveResponse.Removed`) pyta wprost
	// o to, a nie o powodzenie polecenia SQL, które usuwa zero wierszy równie
	// pomyślnie jak jeden.
	UsunZasob(ctx context.Context, kod string) (bool, error)

	// --- kompozycje ---
	ZapiszKompozycje(ctx context.Context, kompozycja KompozycjaDesignu,
		warstwy []WarstwaKompozycji) (KompozycjaDesignu, error)
	Kompozycja(ctx context.Context, kod string) (KompozycjaDesignu, error)
	// Kompozycje domyka `design.board.list` — odczyt kompozycji danego okna.
	// `Kompozycja` czyta po kodzie, którego klient po odświeżeniu okna nie zna:
	// kod nadaje adapter przy pierwszym `design.board.update` i poza samą
	// kompozycją nikt go nie przechowuje, więc bez wykazu po oknie plansza nie
	// ma jak wrócić.
	Kompozycje(ctx context.Context, okno string) ([]KompozycjaDesignu, error)
	Warstwy(ctx context.Context, kompozycjaID int64) ([]WarstwaKompozycji, error)

	// --- kolekcje zasobów (design_kolekcje.go) ---
	ZapiszKolekcjeDesignu(ctx context.Context, kolekcja KolekcjaDesignu) (KolekcjaDesignu, error)
	KolekcjaDesignuPoKodzie(ctx context.Context, kod string) (KolekcjaDesignu, error)
	KolekcjeDesignu(ctx context.Context, okno string, zasob *string) ([]KolekcjaDesignu, error)
	// ZmienPrzypisaniaKolekcjiDesignu dokłada zasoby do kolekcji albo je z niej
	// zdejmuje i oddaje liczbę przypisań, które NAPRAWDĘ się zmieniły —
	// kontrakt (`changed`) pyta o zmianę, nie o powodzenie polecenia SQL.
	ZmienPrzypisaniaKolekcjiDesignu(ctx context.Context, kolekcjaID int64,
		zasoby []string, zdejmij bool) (int, error)
	ZasobyKolekcjiDesignu(ctx context.Context, kolekcjaID int64) ([]string, error)

	// --- szablony promptu i historia promptów (design_szablony.go) ---
	ZapiszSzablonPromptuDesignu(ctx context.Context, szablon SzablonPromptuDesignu) (SzablonPromptuDesignu, error)
	SzablonyPromptuDesignu(ctx context.Context, okno string) ([]SzablonPromptuDesignu, error)
	// PromptyOknaDesignu oddaje prompty wydane w oknie, od najświeższego,
	// wraz z liczbą wszystkich — wykaz bywa przycięty limitem.
	PromptyOknaDesignu(ctx context.Context, okno string, limit int) ([]PromptDesignu, int, error)
	// ZasobyPromptuDesignu oddaje kody zasobów powstałych z promptu — nośnik
	// prowenancji `DesignPromptRecord.AssetIds`.
	ZasobyPromptuDesignu(ctx context.Context, promptID int64) ([]string, error)

	// --- wersje kompozycji (design_wersje.go) ---
	ZapiszWersjeKompozycjiDesignu(ctx context.Context, wersja WersjaKompozycjiDesignu,
		warstwy []WarstwaKompozycji) (WersjaKompozycjiDesignu, error)
	WersjeKompozycjiDesignu(ctx context.Context, kompozycjaID int64, limit int) ([]WersjaKompozycjiDesignu, int, error)
	WersjaKompozycjiDesignuPoKodzie(ctx context.Context, kod string) (WersjaKompozycjiDesignu, error)
	WarstwyWersjiKompozycjiDesignu(ctx context.Context, wersjaID int64) ([]WarstwaKompozycji, error)
	// KompozycjaDesignuPoKluczu czyta kompozycję po kluczu wiersza. Wersja
	// wskazuje kompozycję kluczem obcym, a przywrócenie musi oddać kompozycję
	// z jej identyfikatorem zewnętrznym — bez tego przekładu odpowiedź
	// `design.board.version.restore` nie miałaby czego nieść.
	KompozycjaDesignuPoKluczu(ctx context.Context, id int64) (KompozycjaDesignu, error)

	// --- adnotacje kompozycji (design_adnotacje.go) ---
	ZapiszAdnotacjeDesignu(ctx context.Context, adnotacja AdnotacjaDesignu) (AdnotacjaDesignu, error)
	AdnotacjaDesignuPoKodzie(ctx context.Context, kod string) (AdnotacjaDesignu, error)
	AdnotacjeKompozycjiDesignu(ctx context.Context, kompozycjaID int64, tylkoOtwarte bool) ([]AdnotacjaDesignu, error)

	// --- szablony materiału i profile druku (design_szablony_materialu.go, design_druk.go) ---
	ZapiszSzablonMaterialuDesignu(ctx context.Context, szablon SzablonMaterialuDesignu,
		warstwy []WarstwaKompozycji) (SzablonMaterialuDesignu, error)
	SzablonMaterialuDesignuPoKodzie(ctx context.Context, kod string) (SzablonMaterialuDesignu, error)
	SzablonyMaterialuDesignu(ctx context.Context, okno, rodzaj string) ([]SzablonMaterialuDesignu, error)
	WarstwySzablonuMaterialuDesignu(ctx context.Context, szablonID int64) ([]WarstwaKompozycji, error)
	// Strony szablonu (migracja 339) czynią z niego PUBLIKACJĘ wielostronicową.
	// Szablon bez stron zostaje jednostronicowy i czyta warstwy metodą wyżej —
	// baner nie ma stron.
	ZapiszStronySzablonuMaterialuDesignu(ctx context.Context, szablonID int64,
		strony []StronaSzablonuMaterialuDesignu, warstwy map[string][]WarstwaKompozycji) error
	StronySzablonuMaterialuDesignu(ctx context.Context, szablonID int64) ([]StronaSzablonuMaterialuDesignu, error)
	WarstwyStronySzablonuMaterialuDesignu(ctx context.Context, stronaID int64) ([]WarstwaKompozycji, error)
	ZapiszProfilDrukuDesignu(ctx context.Context, profil ProfilDrukuDesignu) (ProfilDrukuDesignu, error)
	ProfilDrukuDesignuPoKodzie(ctx context.Context, kod string) (ProfilDrukuDesignu, error)
	ProfileDrukuDesignu(ctx context.Context, okno string) ([]ProfilDrukuDesignu, error)
	// ZapiszLicencjeZasobuDesignu i odczyt licencji domykają `design.stock.import`:
	// zasób z katalogu zewnętrznego bez zapisanej licencji jest materiałem,
	// o którym nikt później nie powie, czy wolno go było użyć.
	ZapiszLicencjeZasobuDesignu(ctx context.Context, licencja LicencjaZasobuDesignu) error
	LicencjaZasobuDesignuPoZasobie(ctx context.Context, zasobID int64) (LicencjaZasobuDesignu, error)

	// --- zestawy żetonów (design_zetony.go) ---
	ZapiszZestawZetonowDesignu(ctx context.Context, zestaw ZestawZetonowDesignu,
		zetony []ZetonDesignu) (ZestawZetonowDesignu, error)
	ZestawZetonowDesignuPoKodzie(ctx context.Context, kod string) (ZestawZetonowDesignu, error)
	ZestawyZetonowDesignu(ctx context.Context, okno string, kod *string) ([]ZestawZetonowDesignu, error)
	ZetonyZestawuDesignu(ctx context.Context, zestawID int64) ([]ZetonDesignu, error)

	// --- ikony własne (design_ikony.go) ---
	// W bazie leżą wyłącznie ikony WŁASNE. Katalog ikon otwartoźródłowych jest
	// wkompilowany w binarium rdzenia, więc repozytorium go nie zna i nie ma
	// go czym zasiać ani zgubić.
	ZapiszIkoneDesignu(ctx context.Context, ikona IkonaDesignu) (IkonaDesignu, error)
	IkonaDesignuPoKodzie(ctx context.Context, kod string) (IkonaDesignu, error)
	IkonyDesignuOkna(ctx context.Context, okno string) ([]IkonaDesignu, error)

	// --- warsztat fotografii (design_fotografia.go) ---
	// Łańcuch edycji nie dubluje wariantów zasobu: że wariant powstał ze źródła,
	// mówi `zasob_design.wariant_zasobu_id`. Te metody trzymają to, czego kolumna
	// wariantu nie niesie — CZYNNOŚĆ, jej NASTAWY i DROGĘ rachunku (migracja 346).
	ZapiszCzynnoscFotografiiDesignu(ctx context.Context, czynnosc CzynnoscFotografiiDesignu) error
	CzynnosciFotografiiDesignuZasobu(ctx context.Context, zasobID int64) ([]CzynnoscFotografiiDesignu, error)
	CzynnoscFotografiiDesignuWyniku(ctx context.Context, zasobID int64) (CzynnoscFotografiiDesignu, error)
	ZapiszNastaweFotografiiDesignu(ctx context.Context, nastawa NastawaFotografiiDesignu) (NastawaFotografiiDesignu, error)
	NastawaFotografiiDesignuPoKodzie(ctx context.Context, kod string) (NastawaFotografiiDesignu, error)
	NastawyFotografiiDesignu(ctx context.Context, okno string) ([]NastawaFotografiiDesignu, error)

	// --- ścieżki wektorowe i symbole (design_wektor.go) ---
	// Ścieżka jest bytem osobnym od warstwy, bo krzywa nie mieści się
	// w prostokącie warstwy, a operacja logiczna potrzebuje obu krzywych
	// z osobna (migracja 316).
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

	// --- ramki makiety, więzy responsywne i siatki (design_makiety.go) ---
	// Przynależność warstwy do ramki wiąże się identyfikatorem ZEWNĘTRZNYM
	// warstwy, bo `design.board.update` przepisuje komplet warstw od nowa
	// i klucz wiersza warstwy tego zapisu nie przeżywa (migracja 317).
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
	// WarstwaKompozycjiDesignuPoKodzie, PrzestawWarstweKompozycjiDesignu
	// i DolozWarstweKompozycjiDesignu ruszają JEDNĄ warstwę. `ZapiszKompozycje`
	// przepisuje komplet, więc układ automatyczny i przeliczenie więzi szłyby
	// przez odczyt-zmień-zapisz całej planszy i gubiłyby cudzą zmianę.
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

	// --- gradienty wypełnienia (design_gradienty.go) ---
	// Gradient rozstrzyga się CELEM (kompozycja, ścieżka, warstwa), nie
	// identyfikatorem: `design.color.gradient.set` zakłada albo zmienia
	// gradient stojący na wskazanym bycie, a nie dokłada drugi obok.
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
	// ścieżką (kontrakt: `DesignAsset.VariantOfAssetId`), nie usterką.
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
	// Ziarno i Warianty niesie kontrakt jako *int, a pomocniki drzewa
	// (liczbaDoKolumny) przyjmują *int64 — przekład wprost, bez pomocnika
	// pośredniego, bo to jedyne miejsce w module, gdzie *int trafia do kolumny.
	var ziarno, warianty any
	if prompt.Ziarno != nil {
		ziarno = int64(*prompt.Ziarno)
	}
	if prompt.Warianty != nil {
		warianty = int64(*prompt.Warianty)
	}
	// Kreatywnosc jest *float64 kontraktu; drzewo nie ma pomocnika zapisu dla
	// kolumn REAL (liczbaRzeczywistaZKolumny czyta, ale nie pisze) — zapis wprost.
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

// odczytajPromptDesign składa strukturę z jednego wiersza wyniku.
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
