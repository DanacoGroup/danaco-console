// Odpowiedzialność pliku: dwie drogi powstania makiety — z opisu
// (`design.mockup.generate`) i ze zrzutu ekranu (`design.mockup.import`).
// Ramki, więzy, komponenty i prototyp leżą w `adapter_modul_design_makiety.go`.
//
// ── Makieta jest UKŁADEM, nie obrazkiem ─────────────────────────────────────
// Obie komendy oddają ramkę i WARSTWY, nie plik graficzny. Dlatego kanał
// obrazowy — ten, którym jedzie `design.asset.generate` — nie jest tu drogą:
// obraz ekranu nie da się przesunąć piórem, a makieta ma być czymś, w czym
// Operator pracuje dalej. Rozkład na sekcje liczy więc rdzeń, a kanał modelu
// wchodzi wyłącznie tam, gdzie dodaje wiedzę: przy rozpisaniu OPISU na nazwy
// sekcji.
//
// ── Rozkład opisu: rachunek rdzenia, kanał gdy stoi ─────────────────────────
// Bez wskazanego kanału rdzeń rozkłada opis własną regułą — rozdziela go na
// części i rozpoznaje nazwy sekcji, które w makietach ekranów występują. Wynik
// jest zawsze, na każdej maszynie i bez konta u dostawcy. Ze wskazanym kanałem
// pyta model o rozpisanie tego samego opisu; odpowiedź nieczytelna albo błąd
// kanału NIE kończy komendy — liczy wtedy reguła rdzenia, bo makieta ma powstać.
//
// ── Zrzut: rdzeń wykrywa układ RACHUNKIEM, a litery czyta PROGRAMEM ─────────
// `design.mockup.import` rozkłada zrzut na obszary rachunkiem własnym na
// pikselach: tło z obwodu obrazu, progowanie odstępstwa od tła, spójne bloki,
// prostokąty otaczające. To daje prawdziwy układ warstw ze zrzutu i nie zależy
// od niczego spoza binarium.
//
// TREŚĆ napisów czyta program `tesseract` — składnik pakietu serwera
// zadeklarowany w `zaleznosci_zewnetrzne.go` i wołany przez `zewnetrzne.Wolaj`,
// tą samą drogą, co w Studiu i w Translate. Czytnika liter w czystym Go nie ma,
// więc jest to drugi przypadek zasady: nie ma biblioteki — program idzie
// z pakietem serwera. Ten jeden plik obszaru Design jest objęty nazwanym
// wyjątkiem zapory (`plikiDesignuZOdczytemPisma`); własny `exec` zostaje
// zabroniony i tutaj.
//
// ── Czego brak programu NIE psuje ───────────────────────────────────────────
// Rozstrzygnięcie jest dwustopniowe, bo pole `recognizeText` domyślnie ZNACZY
// „tak", a układ makiety nie ma prawa zależeć od odczytu liter:
//
//   - `recognizeText` podane WPROST jako `true` — Operator poprosił o odczyt.
//     Brak programu jest wtedy ODMOWĄ nazwaną, nie pustym odczytem udającym, że
//     napisów nie było.
//   - pole POMINIĘTE — obowiązuje domyślne „tak", ale odczyt jest dodatkiem do
//     układu. Brak programu daje wtedy makietę z układem, a każda linia tekstu
//     niesie w adnotacji zdanie o tym, że treści nie odczytano i dlaczego. Cisza
//     jest zakazana także tutaj — adnotacja jest jedynym polem, którym da się to
//     powiedzieć, bo kontrakt nie ma na to osobnego.
//
// Obszary, których rdzeń nie rozłożył na elementy, wracają
// w `unrecognizedRegions` — bilans zamiast ciszy.
package core

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/disintegration/imaging"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/models"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/session"
	"danacoconsole/server/internal/zewnetrzne"
	"danacoconsole/shared"
)

const (
	// bokBlokuZrzutuDesignu jest bokiem kratki, na której rdzeń szuka spójnych
	// obszarów. Piksel na piksel byłby rachunkiem na milionach komórek przy
	// każdym zrzucie; kratka ośmiopikselowa zostawia dokładność wystarczającą,
	// żeby rozdzielić przycisk od nagłówka.
	bokBlokuZrzutuDesignu = 8

	// progOdstepstwaOdTlaDesignu jest różnicą składowych, od której punkt uznaje
	// się za treść, a nie za tło. Wartość jest w skali 0–255.
	progOdstepstwaOdTlaDesignu = 24

	// najmniejszyObszarZrzutuDesignu jest najmniejszą liczbą kratek, jaką musi
	// mieć obszar, żeby wejść do makiety jako warstwa. Poniżej tego są pojedyncze
	// litery i drobiny kompresji, a warstwa na literę nie jest układem.
	najmniejszyObszarZrzutuDesignu = 4

	// granicaObszarowZrzutuDesignu chroni makietę przed dwoma tysiącami warstw:
	// zrzut fotografii rozsypuje się na tyle plam, a makieta z dwóch tysięcy
	// warstw nie jest makietą. Obszary powyżej granicy wracają w bilansie.
	granicaObszarowZrzutuDesignu = 64

	// granicaBokuZrzutuDesignu chroni rachunek przed zrzutem o boku
	// dwudziestotysięcznym: kratka liczy się po całej powierzchni.
	granicaBokuZrzutuDesignu = 20000

	// granicaCzasuOdczytuPismaZrzutuDesignu jest granicą czasu JEDNEGO odczytu
	// linii tekstu. Odczyt idzie po jednym obszarze, a obszar jest wycinkiem
	// wysokości wiersza — minuta na taki wycinek to zapas, nie miara.
	granicaCzasuOdczytuPismaZrzutuDesignu = 30 * time.Second

	// jezykiOdczytuPismaZrzutuDesignu to języki, którymi czyta się napisy zrzutu.
	// Oba, bo makiety powstają i z okien polskich, i z angielskich, a wskazania
	// języka kontrakt tej komendy nie ma. Oba pochodzą z pakietu serwera
	// (`tesseract-ocr` niesie angielski, `tesseract-ocr-pol` polski), więc żaden
	// nie jest założeniem o maszynie Operatora.
	jezykiOdczytuPismaZrzutuDesignu = "pol+eng"

	// powiekszenieOdczytuPismaZrzutuDesignu podnosi wycinek przed odczytem.
	// Wiersz interfejsu ma na zrzucie kilkanaście punktów wysokości, a czytnik
	// pisma pracuje na wysokościach kilkakrotnie większych — bez powiększenia
	// odczyt zwracałby znaki przypadkowe zamiast słów.
	powiekszenieOdczytuPismaZrzutuDesignu = 3

	// granicaObszarowOdczytuPismaZrzutuDesignu ogranicza liczbę odczytów w jednym
	// wciągnięciu zrzutu. Każdy odczyt jest osobnym uruchomieniem programu, więc
	// zrzut o sześćdziesięciu liniach tekstu trwałby minuty. Linie powyżej granicy
	// zostają liniami tekstu bez odczytanej treści.
	granicaObszarowOdczytuPismaZrzutuDesignu = 24

	// granicaZnakowOdczytuPismaZrzutuDesignu przycina odczyt wchodzący do
	// adnotacji warstwy. Adnotacja jest podpisem w oknie, nie miejscem na akapit.
	granicaZnakowOdczytuPismaZrzutuDesignu = 160
)

// sekcjeMakietyDesignu to nazwy sekcji rozpoznawane w opisie ekranu wraz
// z wysokością, jaką rdzeń im nadaje. Wysokość nie jest zgadywana za Operatora
// — jest proporcją utrzymaną między sekcjami, żeby makieta była czytelna od
// pierwszego spojrzenia; Operator zmienia ją dalej sam.
var sekcjeMakietyDesignu = []struct {
	Nazwa    string
	Slowa    []string
	Wysokosc float64
}{
	{"pasek stanu", []string{"pasek stanu", "status"}, 44},
	{"nagłówek", []string{"nagłówek", "naglowek", "header", "tytuł", "tytul"}, 80},
	{"nawigacja", []string{"nawigacja", "menu", "zakładki", "zakladki", "nav"}, 56},
	{"bohater", []string{"bohater", "hero", "baner", "banner", "okładka", "okladka"}, 360},
	{"wyszukiwanie", []string{"wyszukiwanie", "szukaj", "search", "filtr"}, 64},
	{"formularz", []string{"formularz", "form", "logowanie", "rejestracja"}, 280},
	{"karty", []string{"karty", "kafle", "siatka", "lista", "grid", "cards"}, 320},
	{"tabela", []string{"tabela", "table", "wykaz"}, 300},
	{"wykres", []string{"wykres", "chart", "statystyki"}, 240},
	{"treść", []string{"treść", "tresc", "opis", "artykuł", "artykul", "content"}, 240},
	{"przyciski", []string{"przycisk", "przyciski", "button", "akcje"}, 64},
	{"stopka", []string{"stopka", "footer"}, 160},
}

// GenerujMakiete zakłada ramkę makiety wraz z warstwami sekcji — obsługuje
// `design.mockup.generate`.
func (a *adapterDesignu) GenerujMakiete(ctx context.Context,
	z shared.DesignMockupGenerateRequest) (shared.DesignMockupGenerateResponse, error) {

	if strings.TrimSpace(z.WindowId) == "" {
		return shared.DesignMockupGenerateResponse{}, bladWskazaniaDesignu(
			"komenda design.mockup.generate bez wskazania okna")
	}
	if strings.TrimSpace(z.BoardId) == "" {
		return shared.DesignMockupGenerateResponse{}, bladWskazaniaDesignu(
			"komenda design.mockup.generate bez wskazania kompozycji")
	}
	if strings.TrimSpace(z.Prompt) == "" {
		return shared.DesignMockupGenerateResponse{}, bladWskazaniaDesignu(
			"komenda design.mockup.generate bez opisu ekranu: rdzeń nie wymyśla, co ma być na " +
				"makiecie, bo makieta wymyślona nie jest makietą Operatora")
	}
	kompozycja, err := a.repozytorium.Kompozycja(ctx, strings.TrimSpace(z.BoardId))
	if err != nil {
		return shared.DesignMockupGenerateResponse{}, bladNieznanejKompozycjiDesignu(z.BoardId, err)
	}
	if z.TokenSetId != nil && strings.TrimSpace(*z.TokenSetId) != "" {
		if _, err := a.repozytorium.ZestawZetonowDesignuPoKodzie(ctx,
			strings.TrimSpace(*z.TokenSetId)); err != nil {
			if czyBrakZasobuDesignu(err) {
				return shared.DesignMockupGenerateResponse{}, bladNieznanegoBytuDesignu(
					"zestawu żetonów " + *z.TokenSetId + " nie ma w tym rdzeniu")
			}
			return shared.DesignMockupGenerateResponse{}, bladDesignu(err)
		}
	}

	szerokosc, wysokosc := domyslnaSzerokoscRamkiDesignu, domyslnaWysokoscRamkiDesignu
	nastawa := ""
	if z.DevicePreset != nil && strings.TrimSpace(*z.DevicePreset) != "" {
		nastawaSzerokosc, nastawaWysokosc, znana := nastawaUrzadzeniaDesignu(*z.DevicePreset)
		if !znana {
			return shared.DesignMockupGenerateResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"komenda design.mockup.generate z nastawą urządzenia %q, której rdzeń nie zna; "+
					"nastawy znane: %s", *z.DevicePreset,
				strings.Join(nazwyNastawUrzadzenDesignu(), ", ")))
		}
		szerokosc, wysokosc, nastawa = nastawaSzerokosc, nastawaWysokosc,
			strings.TrimSpace(*z.DevicePreset)
	}

	nazwySekcji := a.sekcjeOpisuMakietyDesignu(ctx, z)
	if len(nazwySekcji) == 0 {
		return shared.DesignMockupGenerateResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
			"z opisu %q rdzeń nie wyjął ani jednej sekcji ekranu — makiety z jednego słowa nie "+
				"złoży; wymień sekcje po przecinku albo w osobnych wierszach, np. „nagłówek, "+
				"nawigacja, bohater, karty, stopka\"", strings.TrimSpace(z.Prompt)))
	}

	// Ramka makiety powstaje jako nowa: `design.mockup.generate` zakłada ekran,
	// a nie przestawia zastany.
	nazwaRamki := nazwaMakietyDesignu(z.Prompt)
	odpowiedzRamki, err := a.UstawRamke(ctx, shared.DesignFrameSetRequest{
		BoardId: kompozycja.Kod, Name: nazwaRamki,
		Width: &szerokosc, Height: &wysokosc,
		DevicePreset: wskaznikNiepustegoDesignu(&nastawa),
	})
	if err != nil {
		return shared.DesignMockupGenerateResponse{}, err
	}
	ramka, err := a.repozytorium.RamkaDesignuPoKodzie(ctx, odpowiedzRamki.Frame.Id)
	if err != nil {
		return shared.DesignMockupGenerateResponse{}, bladDesignu(err)
	}

	// Komponenty okna, których nazwa pada w opisie, wchodzą do makiety: Operator,
	// który ma w oknie komponent „przycisk główny" i pisze o przyciskach, ma
	// dostać swój komponent, a nie prostokąt bez nazwy.
	komponenty, err := a.repozytorium.KomponentyDesignu(ctx, z.WindowId, nil)
	if err != nil {
		return shared.DesignMockupGenerateResponse{}, bladDesignu(err)
	}

	kodyWarstw := make([]string, 0, len(nazwySekcji))
	uzyteKomponenty := []string{}
	for _, sekcja := range nazwySekcji {
		wysokoscSekcji := sekcja.Wysokosc
		szerokoscSekcji := szerokosc
		adnotacja := sekcja.Nazwa
		warstwa, err := a.repozytorium.DolozWarstweKompozycjiDesignu(ctx, warstwaSekcjiMakietyDesignu(
			kompozycja.ID, adnotacja, szerokoscSekcji, wysokoscSekcji))
		if err != nil {
			return shared.DesignMockupGenerateResponse{}, bladDesignu(err)
		}
		kodyWarstw = append(kodyWarstw, warstwa.Kod)
		for _, komponent := range komponenty {
			if pasujeKomponentDoSekcjiDesignu(komponent.Nazwa, sekcja.Nazwa, z.Prompt) {
				uzyteKomponenty = append(uzyteKomponenty, komponent.Kod)
			}
		}
	}
	if err := a.repozytorium.PrzypiszWarstwyDoRamkiDesignu(ctx, ramka.ID, kodyWarstw); err != nil {
		return shared.DesignMockupGenerateResponse{}, bladDesignu(err)
	}

	// Sekcje układają się pionowo jedna pod drugą — to jest układ ekranu i po to
	// jest układ automatyczny; liczenie tego drugą drogą byłoby drugą prawdą.
	odstep := 0.0
	uklad := shared.DesignAutoLayout{
		Direction: shared.DesignLayoutDirectionVertical,
		Gap:       &odstep,
	}
	odpowiedzUkladu, err := a.UlozAutomatycznie(ctx, shared.DesignLayoutAutoRequest{
		FrameId: ramka.Kod, Layout: uklad, LayerIds: kodyWarstw,
	})
	if err != nil {
		return shared.DesignMockupGenerateResponse{}, err
	}

	swiezaRamka, err := a.repozytorium.RamkaDesignuPoKodzie(ctx, ramka.Kod)
	if err != nil {
		return shared.DesignMockupGenerateResponse{}, bladDesignu(err)
	}
	return shared.DesignMockupGenerateResponse{
		Frame:        ramkaKontraktuDesignu(kompozycja.Kod, swiezaRamka),
		Layers:       odpowiedzUkladu.Layers,
		ComponentIds: uporzadkujBilansDesignu(bezPowtorzenDesignu(uzyteKomponenty)),
	}, nil
}

// warstwaSekcjiMakietyDesignu składa wiersz warstwy jednej sekcji makiety.
func warstwaSekcjiMakietyDesignu(kompozycjaID int64, adnotacja string,
	szerokosc, wysokosc float64) dane.WarstwaKompozycji {

	return dane.WarstwaKompozycji{
		Kod:          nowyIdentyfikator(przedrostekWarstwyDesign),
		KompozycjaID: kompozycjaID,
		Szerokosc:    &szerokosc,
		Wysokosc:     &wysokosc,
		Adnotacja:    &adnotacja,
	}
}

// pasujeKomponentDoSekcjiDesignu rozstrzyga, czy komponent okna należy do tej
// sekcji makiety. Dopasowanie idzie po nazwie: komponent „przycisk główny"
// wchodzi do sekcji „przyciski", a komponent wymieniony w opisie wprost — do
// sekcji, w której jego nazwa padła.
func pasujeKomponentDoSekcjiDesignu(nazwaKomponentu, nazwaSekcji, opis string) bool {
	komponent := strings.ToLower(strings.TrimSpace(nazwaKomponentu))
	if komponent == "" {
		return false
	}
	sekcja := strings.ToLower(nazwaSekcji)
	if strings.Contains(komponent, sekcja) || strings.Contains(sekcja, komponent) {
		return true
	}
	return strings.Contains(strings.ToLower(opis), komponent)
}

// sekcjeOpisuMakietyDesignu rozkłada opis ekranu na sekcje. Kanał modelu wchodzi
// wyłącznie wtedy, gdy Operator go wskazał, i wyłącznie jako ROZPISANIE opisu;
// jego niepowodzenie nie kończy komendy — liczy wtedy reguła rdzenia (nagłówek
// pliku).
func (a *adapterDesignu) sekcjeOpisuMakietyDesignu(ctx context.Context,
	z shared.DesignMockupGenerateRequest) []sekcjaMakietyDesignu {

	opis := z.Prompt
	if z.ChannelId != nil && strings.TrimSpace(*z.ChannelId) != "" {
		polecenie := "Rozpisz opis ekranu na sekcje interfejsu, każdą w osobnym wierszu, " +
			"bez numeracji i bez komentarza. Opis: " + strings.TrimSpace(z.Prompt)
		if rozpisane, err := a.zapytajModelDesignu(ctx, z.WindowId,
			strings.TrimSpace(*z.ChannelId), polecenie); err == nil &&
			strings.TrimSpace(rozpisane) != "" {

			opis = rozpisane
		}
	}
	return sekcjeZOpisuDesignu(opis)
}

// sekcjaMakietyDesignu to jedna sekcja makiety wraz z wysokością.
type sekcjaMakietyDesignu struct {
	Nazwa    string
	Wysokosc float64
}

// sekcjeZOpisuDesignu rozkłada tekst na sekcje regułą rdzenia.
//
// Reguła jest dwustopniowa. Najpierw tekst dzieli się na części po wierszach,
// przecinkach, średnikach i myślnikach — tak Operator wylicza sekcje. Potem
// każda część szuka swojej nazwy w wykazie sekcji; część, która nie trafia
// w żadną, zostaje sekcją o własnej nazwie i wysokości sekcji treści, bo to
// wciąż jest coś, o co Operator prosił.
//
// Sekcje w kolejności ich wystąpienia w opisie: to Operator ustala, co jest
// u góry ekranu, a nie kolejność w wykazie rdzenia.
func sekcjeZOpisuDesignu(opis string) []sekcjaMakietyDesignu {
	czesci := strings.FieldsFunc(opis, func(znak rune) bool {
		return znak == '\n' || znak == ',' || znak == ';' || znak == '•' || znak == '|'
	})
	sekcje := make([]sekcjaMakietyDesignu, 0, len(czesci))
	widziane := map[string]bool{}
	for _, czesc := range czesci {
		czesc = strings.TrimSpace(strings.Trim(czesc, "-–—*·. \t"))
		if czesc == "" {
			continue
		}
		maly := strings.ToLower(czesc)
		nazwa := czesc
		wysokosc := 0.0
		for _, znana := range sekcjeMakietyDesignu {
			trafiona := false
			for _, slowo := range znana.Slowa {
				if strings.Contains(maly, slowo) {
					trafiona = true
					break
				}
			}
			if trafiona {
				nazwa, wysokosc = znana.Nazwa, znana.Wysokosc
				break
			}
		}
		if wysokosc == 0 {
			// Część nierozpoznana zostaje sekcją własną. Wysokość bierze proporcję
			// sekcji treści — jedynej, o której da się coś powiedzieć, gdy nie
			// wiadomo, co to za sekcja.
			wysokosc = 240
			if len(nazwa) > 60 {
				nazwa = strings.TrimSpace(string([]rune(nazwa)[:60])) + "…"
			}
		}
		if widziane[nazwa] {
			continue
		}
		widziane[nazwa] = true
		sekcje = append(sekcje, sekcjaMakietyDesignu{Nazwa: nazwa, Wysokosc: wysokosc})
	}
	return sekcje
}

// nazwaMakietyDesignu składa nazwę ramki z opisu. Nazwa jest podpisem, nie
// wynikiem: ekran opisuje układ warstw, a nazwa go tylko oznacza w wykazie ramek.
func nazwaMakietyDesignu(opis string) string {
	nazwa := strings.TrimSpace(strings.ReplaceAll(opis, "\n", " "))
	const granica = 80
	if znaki := []rune(nazwa); len(znaki) > granica {
		nazwa = strings.TrimSpace(string(znaki[:granica])) + "…"
	}
	if nazwa == "" {
		return "makieta"
	}
	return nazwa
}

// WczytajMakiete rozkłada zrzut ekranu na ramkę i warstwy — obsługuje
// `design.mockup.import`.
func (a *adapterDesignu) WczytajMakiete(ctx context.Context,
	z shared.DesignMockupImportRequest) (shared.DesignMockupImportResponse, error) {

	if strings.TrimSpace(z.WindowId) == "" {
		return shared.DesignMockupImportResponse{}, bladWskazaniaDesignu(
			"komenda design.mockup.import bez wskazania okna")
	}
	if strings.TrimSpace(z.BoardId) == "" {
		return shared.DesignMockupImportResponse{}, bladWskazaniaDesignu(
			"komenda design.mockup.import bez wskazania kompozycji")
	}
	if strings.TrimSpace(z.AssetId) == "" {
		return shared.DesignMockupImportResponse{}, bladWskazaniaDesignu(
			"komenda design.mockup.import bez wskazania zrzutu ekranu")
	}
	kompozycja, err := a.repozytorium.Kompozycja(ctx, strings.TrimSpace(z.BoardId))
	if err != nil {
		return shared.DesignMockupImportResponse{}, bladNieznanejKompozycjiDesignu(z.BoardId, err)
	}
	zasob, err := a.repozytorium.Zasob(ctx, strings.TrimSpace(z.AssetId))
	if err != nil {
		return shared.DesignMockupImportResponse{}, bladNieznanegoZasobuDesignu(z.AssetId, err)
	}
	if zasob.URI == nil || strings.TrimSpace(*zasob.URI) == "" {
		return shared.DesignMockupImportResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
			"zasób %s nie ma treści w magazynie — ze zrzutu bez bajtów nie da się odczytać układu",
			zasob.Kod))
	}
	obraz, err := obrazZasobuDesignu(*zasob.URI)
	if err != nil {
		return shared.DesignMockupImportResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
			"zasób %s nie jest obrazem, który rdzeń potrafi rozłożyć: %s", zasob.Kod, err.Error()))
	}
	granice := obraz.Bounds()
	if granice.Dx() > granicaBokuZrzutuDesignu || granice.Dy() > granicaBokuZrzutuDesignu {
		return shared.DesignMockupImportResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
			"zrzut %d×%d przekracza granicę boku %d — rachunek obszarów idzie po całej "+
				"powierzchni i przy takim boku zamawiałby pamięć, której rdzeń nie dostanie",
			granice.Dx(), granice.Dy(), granicaBokuZrzutuDesignu))
	}

	// Rozmiar ramki bierze wymiary zrzutu, a nastawa urządzenia — gdy wskazana —
	// tylko go NAZYWA. Nadpisanie wymiarów nastawą przeskalowałoby obszary
	// względem zrzutu, z którego je odczytano.
	nastawa := ""
	if z.DevicePreset != nil && strings.TrimSpace(*z.DevicePreset) != "" {
		if _, _, znana := nastawaUrzadzeniaDesignu(*z.DevicePreset); !znana {
			return shared.DesignMockupImportResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"komenda design.mockup.import z nastawą urządzenia %q, której rdzeń nie zna; "+
					"nastawy znane: %s", *z.DevicePreset,
				strings.Join(nazwyNastawUrzadzenDesignu(), ", ")))
		}
		nastawa = strings.TrimSpace(*z.DevicePreset)
	}

	obszary, nierozlozonych := obszaryZrzutuDesignu(obraz)
	if len(obszary) == 0 {
		return shared.DesignMockupImportResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
			"ze zrzutu %s rdzeń nie wyodrębnił ani jednego obszaru — obraz jednolity nie ma "+
				"układu do odtworzenia (obszarów odrzuconych jako zbyt drobne: %d)",
			zasob.Kod, nierozlozonych))
	}

	nazwaRamki := "zrzut " + zasob.Kod
	if zasob.Nazwa != nil && strings.TrimSpace(*zasob.Nazwa) != "" {
		nazwaRamki = strings.TrimSpace(*zasob.Nazwa)
	}
	szerokosc, wysokosc := float64(granice.Dx()), float64(granice.Dy())
	odpowiedzRamki, err := a.UstawRamke(ctx, shared.DesignFrameSetRequest{
		BoardId: kompozycja.Kod, Name: nazwaRamki,
		Width: &szerokosc, Height: &wysokosc,
		DevicePreset: wskaznikNiepustegoDesignu(&nastawa),
	})
	if err != nil {
		return shared.DesignMockupImportResponse{}, err
	}
	ramka, err := a.repozytorium.RamkaDesignuPoKodzie(ctx, odpowiedzRamki.Frame.Id)
	if err != nil {
		return shared.DesignMockupImportResponse{}, bladDesignu(err)
	}

	// Rozpoznanie tekstu jest domyślnie włączone (kontrakt: „brak znaczy tak").
	// Rdzeń wykrywa rachunkiem, KTÓRE obszary są liniami tekstu, i czyta ich
	// treść programem pakietu serwera — nagłówek pliku mówi, którą drogą i co się
	// dzieje, gdy programu nie ma.
	znaczycTekst := z.RecognizeText == nil || *z.RecognizeText
	odczyt := map[int]string{}
	powodBrakuOdczytu := ""
	if znaczycTekst {
		odczytany, err := a.odczytajPismoObszarowDesignu(ctx, obraz, obszary)
		switch {
		case err == nil:
			odczyt = odczytany
		case z.RecognizeText != nil && *z.RecognizeText:
			// Operator poprosił o odczyt WPROST. Makieta bez odczytu byłaby
			// odpowiedzią na inne żądanie, niż zadał.
			return shared.DesignMockupImportResponse{}, bladOdczytuPismaZrzutuDesignu(err)
		default:
			// Pole pominięte: układ jest tym, o co Operator prosił, a odczyt
			// dodatkiem. Brak wchodzi w adnotację każdej linii tekstu, żeby nie
			// wyglądało to na zrzut bez napisów.
			powodBrakuOdczytu = err.Error()
		}
	}

	kody := make([]string, 0, len(obszary))
	warstwy := make([]shared.DesignBoardLayer, 0, len(obszary))
	for numer, obszar := range obszary {
		adnotacja := fmt.Sprintf("obszar %d", numer+1)
		if znaczycTekst && obszar.LiniaTekstu {
			adnotacja = fmt.Sprintf("linia tekstu %d", numer+1)
			switch {
			case odczyt[numer] != "":
				adnotacja = fmt.Sprintf("linia tekstu %d: %s", numer+1, odczyt[numer])
			case powodBrakuOdczytu != "":
				adnotacja = fmt.Sprintf("linia tekstu %d (treści nie odczytano: %s)",
					numer+1, powodBrakuOdczytu)
			}
		}
		x, y := float64(obszar.X), float64(obszar.Y)
		szerokoscObszaru, wysokoscObszaru := float64(obszar.Szerokosc), float64(obszar.Wysokosc)
		warstwa, err := a.repozytorium.DolozWarstweKompozycjiDesignu(ctx, dane.WarstwaKompozycji{
			Kod:          nowyIdentyfikator(przedrostekWarstwyDesign),
			KompozycjaID: kompozycja.ID,
			X:            &x,
			Y:            &y,
			Szerokosc:    &szerokoscObszaru,
			Wysokosc:     &wysokoscObszaru,
			Adnotacja:    &adnotacja,
		})
		if err != nil {
			return shared.DesignMockupImportResponse{}, bladDesignu(err)
		}
		kody = append(kody, warstwa.Kod)
		warstwy = append(warstwy, przelozWarstwyKontraktu(
			[]dane.WarstwaKompozycji{warstwa})...)
	}
	if err := a.repozytorium.PrzypiszWarstwyDoRamkiDesignu(ctx, ramka.ID, kody); err != nil {
		return shared.DesignMockupImportResponse{}, bladDesignu(err)
	}

	odpowiedz := shared.DesignMockupImportResponse{
		Frame:  ramkaKontraktuDesignu(kompozycja.Kod, ramka),
		Layers: warstwy,
	}
	// Bilans wchodzi zawsze, także zerowy: pole mówi, ile obszarów rdzeń
	// odrzucił, a zero jest tu odpowiedzią, nie brakiem odpowiedzi.
	odpowiedz.UnrecognizedRegions = &nierozlozonych
	return odpowiedz, nil
}

// bladOdczytuPismaZrzutuDesignu znakuje odmowę odczytu pisma kodem kontraktu.
//
// Rozstrzygnięcie jest to samo, co w warsztacie dokumentu, bo Operator naprawia
// te trzy rzeczy różnymi ruchami: brak programu jest ZAPLECZEM NIEDOSTĘPNYM
// (kod ponawialny — po dołożeniu pakietu to samo żądanie przejdzie), naruszenie
// izolacji jest odmową punktu Operatora, a reszta usterką rdzenia.
func bladOdczytuPismaZrzutuDesignu(err error) error {
	var brak *zewnetrzne.BrakNarzedzia
	if errors.As(err, &brak) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
			"moduł Design: komenda design.mockup.import z wskazanym odczytem tekstu: "+
				brak.Error()+"; naprawa: dołożyć program do pakietu serwera albo wywołać "+
				"komendę z recognizeText=false — układ obszarów powstanie bez odczytu napisów"))
	}
	if errors.Is(err, session.ErrIzolacja) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodePermissionDenied,
			"moduł Design: odczyt pisma ze zrzutu zatrzymała izolacja: "+err.Error()))
	}
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"moduł Design: odczyt pisma ze zrzutu nie powiódł się: "+err.Error()))
}

// odczytajPismoObszarowDesignu czyta TREŚĆ napisów z obszarów uznanych za linie
// tekstu i oddaje odczyt pod numerem obszaru.
//
// Odczyt idzie obszar po obszarze, nie całym zrzutem naraz. Powód jest
// rozstrzygnięciem: makieta wiąże napis z WARSTWĄ, więc odczyt musi wiedzieć,
// z którego prostokąta pochodzi. Odczyt całości oddałby jeden blok tekstu bez
// przypisania do obszarów i trzeba by go dzielić z powrotem po współrzędnych,
// których czytnik w tej postaci wyjścia nie podaje.
//
// Wycinek idzie do programu jako plik PNG w katalogu tymczasowym. Katalog
// roboczy wywołania zostaje PUSTY, a ścieżka argumentu jest bezwzględna — tak
// samo, jak w warsztacie dokumentu: własny katalog tymczasowy podany bramie
// izolacji zostałby odrzucony jako wyjście poza obszar okna.
func (a *adapterDesignu) odczytajPismoObszarowDesignu(ctx context.Context, obraz image.Image,
	obszary []obszarZrzutuDesignu) (map[int]string, error) {

	numery := []int{}
	for numer, obszar := range obszary {
		if !obszar.LiniaTekstu {
			continue
		}
		if len(numery) >= granicaObszarowOdczytuPismaZrzutuDesignu {
			break
		}
		numery = append(numery, numer)
	}
	if len(numery) == 0 {
		return map[int]string{}, nil
	}
	if a.uruchamiacz == nil {
		return nil, fmt.Errorf(
			"rdzeń nie ma uruchamiacza procesów, więc program rozpoznający pismo nie ma czym " +
				"wystartować; naprawa: podpiąć warstwę kanału przy składaniu rdzenia")
	}

	katalog, err := os.MkdirTemp("", "danaco-zrzut-pismo-")
	if err != nil {
		return nil, fmt.Errorf("nie można założyć katalogu tymczasowego na wycinki zrzutu: %w", err)
	}
	defer func() { _ = os.RemoveAll(katalog) }()

	okno, zasady, obszarZasiegu := a.zasiegOdczytuPismaDesignu()
	odczyt := map[int]string{}
	for _, numer := range numery {
		wycinek := wycinekObszaruZrzutuDesignu(obraz, obszary[numer])
		if wycinek == nil {
			continue
		}
		sciezka := filepath.Join(katalog, "obszar-"+strconv.Itoa(numer+1)+".png")
		var bufor bytes.Buffer
		if err := png.Encode(&bufor, wycinek); err != nil {
			return nil, fmt.Errorf("nie można złożyć wycinka obszaru %d: %w", numer+1, err)
		}
		if err := os.WriteFile(sciezka, bufor.Bytes(), 0o600); err != nil {
			return nil, fmt.Errorf("nie można zapisać wycinka obszaru %d: %w", numer+1, err)
		}

		// `--psm 7` mówi czytnikowi, że obraz jest JEDNĄ linią tekstu — a obszar
		// właśnie nią jest. Bez tego czytnik szuka układu stron w wycinku
		// o wysokości wiersza i oddaje pojedyncze znaki albo nic.
		wynik, err := zewnetrzne.Wolaj(ctx, a.uruchamiacz, okno, zasady, obszarZasiegu,
			narzedzieTesseract,
			[]string{sciezka, "stdout", "-l", jezykiOdczytuPismaZrzutuDesignu, "--psm", "7"},
			"", granicaCzasuOdczytuPismaZrzutuDesignu)
		if err != nil {
			// Brak programu dotyczy WSZYSTKICH obszarów, więc kończy odczyt i wraca
			// do wołającego — ten rozstrzyga, czy jest odmową całej komendy.
			var brak *zewnetrzne.BrakNarzedzia
			if errors.As(err, &brak) || errors.Is(err, session.ErrIzolacja) {
				return nil, err
			}
			// Potknięcie na jednym wycinku nie zabiera makiety: obszar zostaje
			// linią tekstu bez odczytanej treści i adnotacja niczego nie obiecuje.
			continue
		}
		if tresc := jednaLiniaOdczytuDesignu(string(wynik.Wyjscie)); tresc != "" {
			odczyt[numer] = tresc
		}
	}
	return odczyt, nil
}

// zasiegOdczytuPismaDesignu składa trójkę okno–zasady–obszar dla uruchomienia
// programu rozpoznającego pismo.
//
// Rozstrzygnięcie jest to samo, co w warsztacie dokumentu i w rodzinie mowy:
// odczyt pisma jest zdolnością PLATFORMY, nie czynnością okna Operatora, więc
// zasady i obszar składają się dla kontekstu najszerszego poziomu.
// `session.Zasady{}` podstawione z ręki znaczyłoby „izolacja wyłączona"
// niezależnie od tego, co Operator ustawił, a `SrodowiskoWykonania` wpisane jest
// jawnie, bo zrzut leży na hoście rdzenia i program musi ruszyć tam samo.
func (a *adapterDesignu) zasiegOdczytuPismaDesignu() (session.Okno, session.Zasady,
	session.Obszar) {

	okno := session.Okno{Ustawienia: session.Ustawienia{
		SrodowiskoWykonania: shared.ExecutionEnvCore,
	}}
	zasady := session.Zasady{}
	if a.rozstrzygacz != nil {
		zasady = ZasadyIzolacji(a.rozstrzygacz, konfig.Kontekst{})
	}
	obszar := session.Obszar{}
	if a.katalogRoboczy != nil {
		obszar = ObszarOkna(a.katalogRoboczy.Ustal(konfig.Kontekst{}, ""), "")
	}
	return okno, zasady, obszar
}

// wycinekObszaruZrzutuDesignu wycina obszar ze zrzutu i podnosi go przed
// odczytem. Wycinek pusty (obszar poza obrazem) daje `nil` — wołający pomija go
// wtedy zamiast wysyłać do programu plik bez pikseli.
func wycinekObszaruZrzutuDesignu(obraz image.Image, obszar obszarZrzutuDesignu) image.Image {
	granice := obraz.Bounds()
	kadr := image.Rect(
		granice.Min.X+obszar.X, granice.Min.Y+obszar.Y,
		granice.Min.X+obszar.X+obszar.Szerokosc, granice.Min.Y+obszar.Y+obszar.Wysokosc,
	).Intersect(granice)
	if kadr.Dx() < 1 || kadr.Dy() < 1 {
		return nil
	}
	wycinek := imaging.Crop(obraz, kadr)
	// Powiększenie dwuliniowe, nie najbliższym sąsiadem: schodki na literach
	// mylą czytnik bardziej niż lekkie rozmycie krawędzi.
	return imaging.Resize(wycinek,
		kadr.Dx()*powiekszenieOdczytuPismaZrzutuDesignu,
		kadr.Dy()*powiekszenieOdczytuPismaZrzutuDesignu, imaging.Linear)
}

// jednaLiniaOdczytuDesignu sprowadza wyjście programu do jednej linii i przycina
// ją do długości adnotacji. Czytnik oddaje odczyt wraz ze znakami końca linii
// i podwójnymi odstępami, a adnotacja jest podpisem warstwy.
func jednaLiniaOdczytuDesignu(wyjscie string) string {
	tresc := strings.Join(strings.Fields(wyjscie), " ")
	if tresc == "" {
		return ""
	}
	if len(tresc) > granicaZnakowOdczytuPismaZrzutuDesignu {
		// Przycięcie idzie po RUNACH, nie po bajtach: cięcie w środku znaku
		// wielobajtowego dałoby w adnotacji bajt nieczytelny.
		runy := []rune(tresc)
		if len(runy) > granicaZnakowOdczytuPismaZrzutuDesignu {
			tresc = string(runy[:granicaZnakowOdczytuPismaZrzutuDesignu]) + "…"
		}
	}
	return tresc
}

// obszarZrzutuDesignu to jeden prostokąt odczytany ze zrzutu.
type obszarZrzutuDesignu struct {
	X           int
	Y           int
	Szerokosc   int
	Wysokosc    int
	LiniaTekstu bool
}

// obszaryZrzutuDesignu rozkłada zrzut na obszary i oddaje liczbę obszarów, których
// rdzeń nie wziął do makiety.
//
// Rachunek jest własny i wkompilowany, w czterech krokach:
//
//  1. barwa tła bierze się z OBWODU obrazu (najczęstsza wśród punktów brzegu) —
//     tło ekranu prawie zawsze dochodzi do krawędzi, a środek zrzutu bywa
//     wypełniony treścią;
//  2. obraz schodzi na kratkę bloków; blok jest „treścią", gdy jego średnia
//     odbiega od tła powyżej progu;
//  3. bloki treści łączą się w obszary spójne (przejście po kratce);
//  4. obszary za drobne wypadają, a nadliczbowe (powyżej granicy) też — jedne
//     i drugie wracają w liczbie odrzuconych.
func obszaryZrzutuDesignu(obraz image.Image) ([]obszarZrzutuDesignu, int) {
	granice := obraz.Bounds()
	kolumn := (granice.Dx() + bokBlokuZrzutuDesignu - 1) / bokBlokuZrzutuDesignu
	wierszy := (granice.Dy() + bokBlokuZrzutuDesignu - 1) / bokBlokuZrzutuDesignu
	if kolumn == 0 || wierszy == 0 {
		return nil, 0
	}

	tloR, tloG, tloB := barwaTlaZrzutuDesignu(obraz)
	tresc := make([]bool, kolumn*wierszy)
	for wiersz := 0; wiersz < wierszy; wiersz++ {
		for kolumna := 0; kolumna < kolumn; kolumna++ {
			r, g, b := sredniaBlokuZrzutuDesignu(obraz, granice, kolumna, wiersz)
			odstepstwo := absRoznicaDesignu(r, tloR) + absRoznicaDesignu(g, tloG) +
				absRoznicaDesignu(b, tloB)
			tresc[wiersz*kolumn+kolumna] = odstepstwo > progOdstepstwaOdTlaDesignu
		}
	}

	odwiedzone := make([]bool, len(tresc))
	surowe := []obszarZrzutuDesignu{}
	for wiersz := 0; wiersz < wierszy; wiersz++ {
		for kolumna := 0; kolumna < kolumn; kolumna++ {
			indeks := wiersz*kolumn + kolumna
			if !tresc[indeks] || odwiedzone[indeks] {
				continue
			}
			obszar, blokow := obszarSpojnyZrzutuDesignu(tresc, odwiedzone, kolumn, wierszy,
				kolumna, wiersz)
			if blokow < najmniejszyObszarZrzutuDesignu {
				continue
			}
			surowe = append(surowe, obszar)
		}
	}

	// Obszary większe idą pierwsze: przy granicy liczby warstw ma zostać układ,
	// a nie drobiny. Kolejność wynikowa wraca potem do porządku czytania (od
	// góry, potem od lewej), bo tak Operator patrzy na ekran.
	sort.SliceStable(surowe, func(i, j int) bool {
		return surowe[i].Szerokosc*surowe[i].Wysokosc > surowe[j].Szerokosc*surowe[j].Wysokosc
	})
	nierozlozonych := 0
	wybrane := surowe
	if len(surowe) > granicaObszarowZrzutuDesignu {
		nierozlozonych = len(surowe) - granicaObszarowZrzutuDesignu
		wybrane = surowe[:granicaObszarowZrzutuDesignu]
	}
	sort.SliceStable(wybrane, func(i, j int) bool {
		if wybrane[i].Y != wybrane[j].Y {
			return wybrane[i].Y < wybrane[j].Y
		}
		return wybrane[i].X < wybrane[j].X
	})
	for numer := range wybrane {
		wybrane[numer].LiniaTekstu = czyLiniaTekstuDesignu(wybrane[numer])
	}
	return wybrane, nierozlozonych
}

// obszarSpojnyZrzutuDesignu przechodzi po kratce od wskazanego bloku i oddaje
// prostokąt otaczający obszar spójny wraz z liczbą jego bloków.
func obszarSpojnyZrzutuDesignu(tresc, odwiedzone []bool, kolumn, wierszy,
	kolumnaStart, wierszStart int) (obszarZrzutuDesignu, int) {

	kolejka := [][2]int{{kolumnaStart, wierszStart}}
	odwiedzone[wierszStart*kolumn+kolumnaStart] = true
	lewa, prawa := kolumnaStart, kolumnaStart
	gora, dol := wierszStart, wierszStart
	blokow := 0
	for len(kolejka) > 0 {
		biezacy := kolejka[0]
		kolejka = kolejka[1:]
		blokow++
		kolumna, wiersz := biezacy[0], biezacy[1]
		if kolumna < lewa {
			lewa = kolumna
		}
		if kolumna > prawa {
			prawa = kolumna
		}
		if wiersz < gora {
			gora = wiersz
		}
		if wiersz > dol {
			dol = wiersz
		}
		for _, krok := range [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
			sasiadKolumna, sasiadWiersz := kolumna+krok[0], wiersz+krok[1]
			if sasiadKolumna < 0 || sasiadKolumna >= kolumn ||
				sasiadWiersz < 0 || sasiadWiersz >= wierszy {
				continue
			}
			indeks := sasiadWiersz*kolumn + sasiadKolumna
			if !tresc[indeks] || odwiedzone[indeks] {
				continue
			}
			odwiedzone[indeks] = true
			kolejka = append(kolejka, [2]int{sasiadKolumna, sasiadWiersz})
		}
	}
	return obszarZrzutuDesignu{
		X:         lewa * bokBlokuZrzutuDesignu,
		Y:         gora * bokBlokuZrzutuDesignu,
		Szerokosc: (prawa - lewa + 1) * bokBlokuZrzutuDesignu,
		Wysokosc:  (dol - gora + 1) * bokBlokuZrzutuDesignu,
	}, blokow
}

// czyLiniaTekstuDesignu rozstrzyga, czy obszar wygląda na linię pisma: niski
// i wyraźnie szerszy niż wysoki. To rozpoznanie UKŁADU, nie treści — rdzeń mówi
// „tu jest linia tekstu", a nie „tu jest napisane X" (nagłówek pliku).
func czyLiniaTekstuDesignu(obszar obszarZrzutuDesignu) bool {
	if obszar.Wysokosc <= 0 || obszar.Wysokosc > 6*bokBlokuZrzutuDesignu {
		return false
	}
	return obszar.Szerokosc >= 3*obszar.Wysokosc
}

// barwaTlaZrzutuDesignu odczytuje barwę tła z obwodu obrazu.
func barwaTlaZrzutuDesignu(obraz image.Image) (int, int, int) {
	granice := obraz.Bounds()
	licznik := map[[3]int]int{}
	dodaj := func(x, y int) {
		r, g, b, _ := obraz.At(x, y).RGBA()
		// Kwantyzacja po 16: tło z gradientem albo szumem kompresji nie ma
		// jednej dokładnej wartości, a bez kwantyzacji każdy punkt byłby własną
		// barwą i najczęstsza nie znaczyłaby nic.
		klucz := [3]int{int(r>>8) / 16, int(g>>8) / 16, int(b>>8) / 16}
		licznik[klucz]++
	}
	for x := granice.Min.X; x < granice.Max.X; x++ {
		dodaj(x, granice.Min.Y)
		dodaj(x, granice.Max.Y-1)
	}
	for y := granice.Min.Y; y < granice.Max.Y; y++ {
		dodaj(granice.Min.X, y)
		dodaj(granice.Max.X-1, y)
	}
	najczestsza := [3]int{}
	najwiecej := -1
	for klucz, ile := range licznik {
		if ile > najwiecej {
			najczestsza, najwiecej = klucz, ile
		}
	}
	return najczestsza[0]*16 + 8, najczestsza[1]*16 + 8, najczestsza[2]*16 + 8
}

// sredniaBlokuZrzutuDesignu liczy średnią składową jednego bloku kratki.
func sredniaBlokuZrzutuDesignu(obraz image.Image, granice image.Rectangle,
	kolumna, wiersz int) (int, int, int) {

	lewa := granice.Min.X + kolumna*bokBlokuZrzutuDesignu
	gora := granice.Min.Y + wiersz*bokBlokuZrzutuDesignu
	prawa := lewa + bokBlokuZrzutuDesignu
	dol := gora + bokBlokuZrzutuDesignu
	if prawa > granice.Max.X {
		prawa = granice.Max.X
	}
	if dol > granice.Max.Y {
		dol = granice.Max.Y
	}
	sumaR, sumaG, sumaB, punktow := 0, 0, 0, 0
	for y := gora; y < dol; y++ {
		for x := lewa; x < prawa; x++ {
			r, g, b, _ := obraz.At(x, y).RGBA()
			sumaR += int(r >> 8)
			sumaG += int(g >> 8)
			sumaB += int(b >> 8)
			punktow++
		}
	}
	if punktow == 0 {
		return 0, 0, 0
	}
	return sumaR / punktow, sumaG / punktow, sumaB / punktow
}

// absRoznicaDesignu oddaje wartość bezwzględną różnicy dwóch liczb całkowitych.
func absRoznicaDesignu(pierwsza, druga int) int {
	if pierwsza < druga {
		return druga - pierwsza
	}
	return pierwsza - druga
}

// bezPowtorzenDesignu usuwa powtórzenia z wykazu, zachowując kolejność
// pierwszego wystąpienia.
func bezPowtorzenDesignu(wykaz []string) []string {
	widziane := make(map[string]bool, len(wykaz))
	wynik := make([]string, 0, len(wykaz))
	for _, wartosc := range wykaz {
		if widziane[wartosc] {
			continue
		}
		widziane[wartosc] = true
		wynik = append(wynik, wartosc)
	}
	return wynik
}

// zapytajModelDesignu woła wskazany kanał modelu i zbiera całą odpowiedź
// tekstową. Most jest wspólny dla czynności słownych modułu Design (rozpisanie
// opisu makiety, propozycje zestawień krojów).
//
// Niepowodzenie wraca błędem, a wołający rozstrzyga, czy je przeżyć: przy
// makiecie i przy zestawieniach krojów rdzeń liczy wtedy własną regułą, bo wynik
// ma być zawsze. Odmowa całej komendy dlatego, że cudzy dostawca nie odpowiedział,
// byłaby brakiem funkcji u Operatora.
func (a *adapterDesignu) zapytajModelDesignu(ctx context.Context, okno, kanal,
	tresc string) (string, error) {

	if a.kanaly == nil {
		return "", fmt.Errorf("rejestr kanałów modelu nie jest wpięty")
	}
	if strings.TrimSpace(kanal) == "" {
		return "", fmt.Errorf("żądanie bez wskazania kanału modelu")
	}
	var zebrane strings.Builder
	ujscie := models.UjscieFunkcji(func(_ context.Context, f models.Fragment) error {
		if f.Kind == shared.ChunkKindText {
			zebrane.WriteString(models.TrescFragmentu(f))
		}
		return nil
	})
	zapytanie := models.Zapytanie{
		Zasiegi:   models.Zasiegi{Okno: okno},
		Wiadomosc: okno,
		Tresc:     tresc,
		Kanal:     kanal,
	}
	if err := a.kanaly.Wyslij(ctx, zapytanie, ujscie); err != nil {
		return "", err
	}
	return zebrane.String(), nil
}
