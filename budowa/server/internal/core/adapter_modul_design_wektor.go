// Odpowiedzialność pliku: dziesięć czynności warsztatu wektorowego modułu
// Design — `design.vector.*`. Metody stoją na `*adapterDesignu`
// (`adapter_modul_design.go`); rachunek na ścieżkach leży
// w `adapter_modul_design_wektor_sciezki.go`, katalog krojów
// w `adapter_modul_design_kroje.go`, wpięcie w `adapter_modul_design_uchwyty.go`.
//
// ── Kształt powstaje OD RAZU jako węzły ścieżki ─────────────────────────────
// `design.vector.shape.add` nie zakłada „prostokąta" jako osobnego bytu do
// późniejszej zamiany w ścieżkę. Prostokąt, elipsa, wielokąt i gwiazda wchodzą
// do bazy jako komplet węzłów z uchwytami, więc Operator ciągnie je piórem od
// pierwszej chwili. Byt pośredni wymagałby komendy „zamień w ścieżkę", której
// kontrakt nie ma, i dawał na planszy dwa rodzaje kształtu różniące się tym,
// czego z nimi wolno zrobić.
//
// ── Zamiana tekstu w kontury jest nieodwracalna dla WYNIKU ──────────────────
// Kontury nie wiedzą, że były literami: po zamianie nie da się poprawić
// literówki. Dlatego tekst źródłowy ZOSTAJE — w nazwie ścieżki
// (`sciezka_wektorowa_design.nazwa`), skąd okno go odczytuje i pokazuje obok
// konturów. Odpowiedź mówi wprost polem `outlined`, co Operator dostał.
//
// ── Operacja logiczna oddaje węzły, nie obrazek ─────────────────────────────
// Suma, różnica, część wspólna i wykluczenie idą przez `tdewolff/canvas`, a
// wynik wraca do bazy znowu jako WĘZŁY — nie jako gotowy napis SVG. Ścieżka po
// operacji ma dać się ciągnąć piórem dalej, inaczej pierwsza suma dwóch kół
// kończyłaby edycję kształtu.
//
// ── Ścieżki źródłowe usunięte wracają w bilansie ────────────────────────────
// `keepSources` bez wskazania znaczy „usuń źródła" — tak działa operacja
// logiczna w każdym programie wektorowym. Usunięte wracają w `removedPathIds`:
// pole puste tam, gdzie ścieżki zniknęły, byłoby ciszą, a okno pokazywałoby
// kształty, których w bazie już nie ma.
package core

import (
	"bytes"
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/tdewolff/canvas"
	renderPdf "github.com/tdewolff/canvas/renderers/pdf"
	renderPs "github.com/tdewolff/canvas/renderers/ps"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

const (
	// przedrostekSciezkiDesign i przedrostekSymboluDesign znakują identyfikatory
	// zewnętrzne bytów warsztatu wektorowego — obok `plansza-` i `warstwa-`
	// z obszaru kompozycji.
	przedrostekSciezkiDesign = "sciezka-"
	przedrostekSymboluDesign = "symbol-"

	// najmniejWezlowSciezkiDesignu jest kresem sensu: ścieżka o jednym węźle
	// jest punktem, a punkt nie jest kształtem, który da się wyrysować.
	najmniejWezlowSciezkiDesignu = 2
)

// UstawSciezke zakłada ścieżkę albo nadpisuje zastaną — obsługuje
// `design.vector.path.set`.
func (a *adapterDesignu) UstawSciezke(ctx context.Context,
	z shared.DesignVectorPathSetRequest) (shared.DesignVectorPathSetResponse, error) {

	if strings.TrimSpace(z.BoardId) == "" {
		return shared.DesignVectorPathSetResponse{}, bladWskazaniaDesignu(
			"komenda design.vector.path.set bez wskazania kompozycji")
	}
	if err := sprawdzWezlySciezkiDesignu("design.vector.path.set", z.Nodes); err != nil {
		return shared.DesignVectorPathSetResponse{}, err
	}

	kompozycja, err := a.repozytorium.Kompozycja(ctx, strings.TrimSpace(z.BoardId))
	if err != nil {
		return shared.DesignVectorPathSetResponse{}, bladNieznanejKompozycjiDesignu(z.BoardId, err)
	}

	kod := nowyIdentyfikator(przedrostekSciezkiDesign)
	kolejnosc := 0
	if z.PathId != nil && strings.TrimSpace(*z.PathId) != "" {
		kod = strings.TrimSpace(*z.PathId)
		// Ścieżka wskazana a nieznana jest ODMOWĄ, nie cichym założeniem nowej:
		// Operator poprawiający kształt oczekuje, że poprawił ten jeden, a nie
		// że dostał drugi obok.
		zastana, err := a.repozytorium.SciezkaWektorowaDesignuPoKodzie(ctx, kod)
		if err != nil {
			return shared.DesignVectorPathSetResponse{}, bladNieznanejSciezkiDesignu(kod, err)
		}
		if zastana.KompozycjaID != kompozycja.ID {
			return shared.DesignVectorPathSetResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"ścieżka %s leży na innej kompozycji niż %s — komenda design.vector.path.set nie "+
					"przenosi kształtów między planszami", kod, kompozycja.Kod))
		}
		kolejnosc = zastana.Kolejnosc
	}

	wiersz, err := a.zapiszSciezkeDesignu(ctx, kompozycja.ID, kod, kolejnosc, z.LayerId, z.Name,
		z.Nodes, z.Closed != nil && *z.Closed, z.Fill, z.Stroke)
	if err != nil {
		return shared.DesignVectorPathSetResponse{}, err
	}
	sciezka, err := sciezkaKontraktuDesignu(kompozycja.Kod, wiersz)
	if err != nil {
		return shared.DesignVectorPathSetResponse{}, bladWydaniaDesignu(err.Error())
	}
	return shared.DesignVectorPathSetResponse{Path: sciezka}, nil
}

// zapiszSciezkeDesignu składa wiersz ścieżki i utrwala go — droga wspólna dla
// rysowania piórem, kształtu, tekstu i wyniku operacji logicznej.
func (a *adapterDesignu) zapiszSciezkeDesignu(ctx context.Context, kompozycjaID int64,
	kod string, kolejnosc int, warstwa, nazwa *string, wezly []shared.DesignVectorNode,
	zamknieta bool, wypelnienie *shared.DesignFill,
	obrys *shared.DesignStroke) (dane.SciezkaWektorowaDesignu, error) {

	zapis, err := zapisWezlowDesignu(wezly)
	if err != nil {
		return dane.SciezkaWektorowaDesignu{}, bladWydaniaDesignu(err.Error())
	}
	zapisWypelnienia, err := zapisWypelnieniaDesignu(wypelnienie)
	if err != nil {
		return dane.SciezkaWektorowaDesignu{}, bladWydaniaDesignu(err.Error())
	}
	zapisObrysu, err := zapisObrysuDesignu(obrys)
	if err != nil {
		return dane.SciezkaWektorowaDesignu{}, bladWydaniaDesignu(err.Error())
	}
	wiersz, err := a.repozytorium.ZapiszSciezkeWektorowaDesignu(ctx, dane.SciezkaWektorowaDesignu{
		Kod:             kod,
		KompozycjaID:    kompozycjaID,
		WarstwaKod:      wskaznikNiepustegoDesignu(warstwa),
		Nazwa:           wskaznikNiepustegoDesignu(nazwa),
		WezlyJSON:       zapis,
		Zamknieta:       zamknieta,
		WypelnienieJSON: zapisWypelnienia,
		ObrysJSON:       zapisObrysu,
		Kolejnosc:       kolejnosc,
	})
	if err != nil {
		return dane.SciezkaWektorowaDesignu{}, bladDesignu(err)
	}
	return wiersz, nil
}

// SciezkiWektorowe zwraca ścieżki kompozycji — obsługuje
// `design.vector.path.list`.
func (a *adapterDesignu) SciezkiWektorowe(ctx context.Context,
	z shared.DesignVectorPathListRequest) (shared.DesignVectorPathListResponse, error) {

	if strings.TrimSpace(z.BoardId) == "" {
		return shared.DesignVectorPathListResponse{}, bladWskazaniaDesignu(
			"komenda design.vector.path.list bez wskazania kompozycji")
	}
	kompozycja, err := a.repozytorium.Kompozycja(ctx, strings.TrimSpace(z.BoardId))
	if err != nil {
		return shared.DesignVectorPathListResponse{}, bladNieznanejKompozycjiDesignu(z.BoardId, err)
	}
	wiersze, err := a.repozytorium.SciezkiWektoroweDesignu(ctx, kompozycja.ID,
		wskaznikNiepustegoDesignu(z.LayerId))
	if err != nil {
		return shared.DesignVectorPathListResponse{}, bladDesignu(err)
	}
	sciezki := make([]shared.DesignVectorPath, 0, len(wiersze))
	for _, wiersz := range wiersze {
		sciezka, err := sciezkaKontraktuDesignu(kompozycja.Kod, wiersz)
		if err != nil {
			return shared.DesignVectorPathListResponse{}, bladWydaniaDesignu(err.Error())
		}
		sciezki = append(sciezki, sciezka)
	}
	return shared.DesignVectorPathListResponse{Paths: sciezki, Total: len(sciezki)}, nil
}

// UsunSciezke usuwa ścieżkę — obsługuje `design.vector.path.remove`.
//
// Ścieżki, której nie ma, nie odmawiamy: kontrakt pyta polem `removed`, czy
// wiersz istniał, a nie czy polecenie SQL się udało — usunięcie czegoś, czego
// nie było, jest odpowiedzią, nie usterką.
func (a *adapterDesignu) UsunSciezke(ctx context.Context,
	z shared.DesignVectorPathRemoveRequest) (shared.DesignVectorPathRemoveResponse, error) {

	if strings.TrimSpace(z.PathId) == "" {
		return shared.DesignVectorPathRemoveResponse{}, bladWskazaniaDesignu(
			"komenda design.vector.path.remove bez wskazania ścieżki")
	}
	usunieta, err := a.repozytorium.UsunSciezkeWektorowaDesignu(ctx, strings.TrimSpace(z.PathId))
	if err != nil {
		return shared.DesignVectorPathRemoveResponse{}, bladDesignu(err)
	}
	return shared.DesignVectorPathRemoveResponse{Removed: usunieta}, nil
}

// DolozKsztalt zakłada kształt podstawowy jako ścieżkę o węzłach — obsługuje
// `design.vector.shape.add`.
func (a *adapterDesignu) DolozKsztalt(ctx context.Context,
	z shared.DesignVectorShapeAddRequest) (shared.DesignVectorShapeAddResponse, error) {

	if strings.TrimSpace(z.BoardId) == "" {
		return shared.DesignVectorShapeAddResponse{}, bladWskazaniaDesignu(
			"komenda design.vector.shape.add bez wskazania kompozycji")
	}
	if err := sprawdzWyliczenieDesignu("design.vector.shape.add", "kind", z.Kind,
		shared.WartosciDesignShapeKind()); err != nil {
		return shared.DesignVectorShapeAddResponse{}, err
	}
	// Odcinek jest jedynym kształtem, którego bok wolno mieć zerowy: linia
	// pionowa ma zerową szerokość i nadal jest linią. Pozostałe kształty o boku
	// niedodatnim nie istnieją.
	if z.Kind == shared.DesignShapeKindLine {
		if z.Width == 0 && z.Height == 0 {
			return shared.DesignVectorShapeAddResponse{}, bladWskazaniaDesignu(
				"komenda design.vector.shape.add z odcinkiem o zerowej długości: odcinek o obu " +
					"bokach zerowych jest punktem, nie kształtem")
		}
	} else if z.Width <= 0 || z.Height <= 0 {
		return shared.DesignVectorShapeAddResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
			"komenda design.vector.shape.add z kształtem %v×%v: kształt o niedodatnim boku nie istnieje",
			z.Width, z.Height))
	}

	kompozycja, err := a.repozytorium.Kompozycja(ctx, strings.TrimSpace(z.BoardId))
	if err != nil {
		return shared.DesignVectorShapeAddResponse{}, bladNieznanejKompozycjiDesignu(z.BoardId, err)
	}

	wezly, zamknieta := ksztaltWezlamiDesignu(z)
	if len(wezly) < najmniejWezlowSciezkiDesignu {
		return shared.DesignVectorShapeAddResponse{}, bladWydaniaDesignu(fmt.Sprintf(
			"rachunek kształtu %s nie dał ani jednego węzła", string(z.Kind)))
	}
	nazwa := string(z.Kind)
	wiersz, err := a.zapiszSciezkeDesignu(ctx, kompozycja.ID,
		nowyIdentyfikator(przedrostekSciezkiDesign), 0, z.LayerId, &nazwa,
		wezly, zamknieta, z.Fill, z.Stroke)
	if err != nil {
		return shared.DesignVectorShapeAddResponse{}, err
	}
	sciezka, err := sciezkaKontraktuDesignu(kompozycja.Kod, wiersz)
	if err != nil {
		return shared.DesignVectorShapeAddResponse{}, bladWydaniaDesignu(err.Error())
	}
	return shared.DesignVectorShapeAddResponse{Path: sciezka}, nil
}

// ZlozSciezkiLogicznie liczy operację logiczną na ścieżkach — obsługuje
// `design.vector.boolean`.
func (a *adapterDesignu) ZlozSciezkiLogicznie(ctx context.Context,
	z shared.DesignVectorBooleanRequest) (shared.DesignVectorBooleanResponse, error) {

	if len(z.PathIds) < 2 {
		return shared.DesignVectorBooleanResponse{}, bladWskazaniaDesignu(
			"komenda design.vector.boolean z mniej niż dwiema ścieżkami: operacja logiczna " +
				"potrzebuje dwóch kształtów, żeby mieć co z czym złożyć")
	}
	if err := sprawdzWyliczenieDesignu("design.vector.boolean", "operation", z.Operation,
		shared.WartosciDesignBooleanOp()); err != nil {
		return shared.DesignVectorBooleanResponse{}, err
	}

	wiersze := make([]dane.SciezkaWektorowaDesignu, 0, len(z.PathIds))
	sciezkiBiblioteki := make([]*canvas.Path, 0, len(z.PathIds))
	for _, kod := range z.PathIds {
		kod = strings.TrimSpace(kod)
		wiersz, err := a.repozytorium.SciezkaWektorowaDesignuPoKodzie(ctx, kod)
		if err != nil {
			return shared.DesignVectorBooleanResponse{}, bladNieznanejSciezkiDesignu(kod, err)
		}
		if len(wiersze) > 0 && wiersz.KompozycjaID != wiersze[0].KompozycjaID {
			return shared.DesignVectorBooleanResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"ścieżka %s leży na innej kompozycji niż %s — operacja logiczna nie łączy kształtów "+
					"z dwóch plansz, bo wynik nie miałby gdzie stanąć", kod, wiersze[0].Kod))
		}
		wezly, err := wezlyZeZapisuDesignu(wiersz.WezlyJSON)
		if err != nil {
			return shared.DesignVectorBooleanResponse{}, bladWydaniaDesignu(err.Error())
		}
		wiersze = append(wiersze, wiersz)
		sciezkiBiblioteki = append(sciezkiBiblioteki, sciezkaBibliotekiDesignu(wezly, wiersz.Zamknieta))
	}

	wynik := sciezkiBiblioteki[0]
	for _, kolejna := range sciezkiBiblioteki[1:] {
		switch z.Operation {
		case shared.DesignBooleanOpUnion:
			wynik = wynik.Or(kolejna)
		case shared.DesignBooleanOpSubtract:
			wynik = wynik.Not(kolejna)
		case shared.DesignBooleanOpIntersect:
			wynik = wynik.And(kolejna)
		case shared.DesignBooleanOpExclude:
			wynik = wynik.Xor(kolejna)
		}
	}
	wezlyWyniku, zamknieta := wezlyZeSciezkiBibliotekiDesignu(wynik)
	if len(wezlyWyniku) < najmniejWezlowSciezkiDesignu {
		// Wynik pusty jest PRAWDĄ o kształtach, nie usterką rachunku: część
		// wspólna dwóch rozłącznych kół jest pusta. Odmowa nazywa to wprost,
		// zamiast zakładać ścieżkę bez węzłów, której nic nie pokaże.
		return shared.DesignVectorBooleanResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
			"operacja %s na wskazanych ścieżkach dała kształt pusty — rdzeń nie zakłada ścieżki "+
				"bez węzłów; sprawdź, czy kształty się w ogóle nakładają", string(z.Operation)))
	}

	// Wynik dziedziczy wypełnienie i obrys pierwszej ścieżki: to ona jest
	// kształtem wiodącym (kolejność rozstrzyga przy różnicy), więc jej wygląd
	// jest tym, którego Operator się spodziewa.
	nazwa := string(z.Operation)
	nowa, err := a.zapiszSciezkeDesignu(ctx, wiersze[0].KompozycjaID,
		nowyIdentyfikator(przedrostekSciezkiDesign), 0, wiersze[0].WarstwaKod, &nazwa,
		wezlyWyniku, zamknieta, wypelnienieZeZapisuDesignu(wiersze[0].WypelnienieJSON),
		obrysZeZapisuDesignu(wiersze[0].ObrysJSON))
	if err != nil {
		return shared.DesignVectorBooleanResponse{}, err
	}

	usuniete := []string{}
	if z.KeepSources == nil || !*z.KeepSources {
		for _, wiersz := range wiersze {
			zniknela, err := a.repozytorium.UsunSciezkeWektorowaDesignu(ctx, wiersz.Kod)
			if err != nil {
				return shared.DesignVectorBooleanResponse{}, bladDesignu(err)
			}
			if zniknela {
				usuniete = append(usuniete, wiersz.Kod)
			}
		}
	}

	kompozycja, err := a.repozytorium.KompozycjaDesignuPoKluczu(ctx, wiersze[0].KompozycjaID)
	if err != nil {
		return shared.DesignVectorBooleanResponse{}, bladDesignu(err)
	}
	sciezka, err := sciezkaKontraktuDesignu(kompozycja.Kod, nowa)
	if err != nil {
		return shared.DesignVectorBooleanResponse{}, bladWydaniaDesignu(err.Error())
	}
	odpowiedz := shared.DesignVectorBooleanResponse{Path: sciezka}
	if len(usuniete) > 0 {
		odpowiedz.RemovedPathIds = usuniete
	}
	return odpowiedz, nil
}

// TekstNaSciezce układa tekst wzdłuż ścieżki albo zamienia go w kontury —
// obsługuje `design.vector.text.path`.
//
// Tekst źródłowy zostaje w nazwie ścieżki także po zamianie w kontury (nagłówek
// pliku): kontury nie wiedzą, że były literami, a Operator ma po czym poznać, co
// tam napisał.
func (a *adapterDesignu) TekstNaSciezce(ctx context.Context,
	z shared.DesignVectorTextPathRequest) (shared.DesignVectorTextPathResponse, error) {

	if strings.TrimSpace(z.BoardId) == "" {
		return shared.DesignVectorTextPathResponse{}, bladWskazaniaDesignu(
			"komenda design.vector.text.path bez wskazania kompozycji")
	}
	if strings.TrimSpace(z.Text) == "" {
		return shared.DesignVectorTextPathResponse{}, bladWskazaniaDesignu(
			"komenda design.vector.text.path bez treści tekstu")
	}
	if z.FontSize <= 0 {
		return shared.DesignVectorTextPathResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
			"komenda design.vector.text.path z rozmiarem pisma %v: pismo o niedodatnim rozmiarze "+
				"nie ma konturu", z.FontSize))
	}
	krojWczytany, nazwaKroju, _, err := krojDesignu(z.FontFamily)
	if err != nil {
		return shared.DesignVectorTextPathResponse{}, bladWskazaniaDesignu(
			"komenda design.vector.text.path: " + err.Error())
	}

	kompozycja, err := a.repozytorium.Kompozycja(ctx, strings.TrimSpace(z.BoardId))
	if err != nil {
		return shared.DesignVectorTextPathResponse{}, bladNieznanejKompozycjiDesignu(z.BoardId, err)
	}

	// Ścieżka nośna: wskazana przez Operatora albo linia pisma od punktu (x, y).
	// Wskazana ścieżka rozstrzyga też o początku tekstu — pierwszy jej węzeł.
	poczatekX, poczatekY := 0.0, z.FontSize
	if z.X != nil {
		poczatekX = *z.X
	}
	if z.Y != nil {
		poczatekY = *z.Y
	}
	var nosna dane.SciezkaWektorowaDesignu
	maNosna := false
	if z.PathId != nil && strings.TrimSpace(*z.PathId) != "" {
		nosna, err = a.repozytorium.SciezkaWektorowaDesignuPoKodzie(ctx, strings.TrimSpace(*z.PathId))
		if err != nil {
			return shared.DesignVectorTextPathResponse{}, bladNieznanejSciezkiDesignu(*z.PathId, err)
		}
		if nosna.KompozycjaID != kompozycja.ID {
			return shared.DesignVectorTextPathResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"ścieżka nośna %s leży na innej kompozycji niż %s", nosna.Kod, kompozycja.Kod))
		}
		wezlyNosnej, err := wezlyZeZapisuDesignu(nosna.WezlyJSON)
		if err != nil {
			return shared.DesignVectorTextPathResponse{}, bladWydaniaDesignu(err.Error())
		}
		if len(wezlyNosnej) == 0 {
			return shared.DesignVectorTextPathResponse{}, bladWskazaniaDesignu(
				"ścieżka nośna " + nosna.Kod + " nie ma ani jednego węzła — nie ma wzdłuż czego układać tekstu")
		}
		poczatekX, poczatekY = wezlyNosnej[0].X, wezlyNosnej[0].Y
		maNosna = true
	}

	konturowac := z.Outline != nil && *z.Outline
	nazwa := strings.TrimSpace(z.Text)

	if !konturowac {
		// Bez konturowania ścieżka niesie LINIĘ PISMA: albo wskazaną ścieżkę
		// nośną (wtedy zapisujemy nową ścieżkę o tych samych węzłach, żeby tekst
		// dostał własny byt i nie nadpisał kształtu nośnego), albo odcinek
		// o zmierzonej długości tekstu. Odpowiedź mówi `outlined: false` — tekst
		// nadal jest tekstem i literówkę da się poprawić.
		wezly := []shared.DesignVectorNode{}
		zamknieta := false
		if maNosna {
			wezly, err = wezlyZeZapisuDesignu(nosna.WezlyJSON)
			if err != nil {
				return shared.DesignVectorTextPathResponse{}, bladWydaniaDesignu(err.Error())
			}
			zamknieta = nosna.Zamknieta
		} else {
			szerokosc, brakujacych := szerokoscTekstuDesignu(krojWczytany, nazwa, z.FontSize)
			if brakujacych > 0 && szerokosc == 0 {
				return shared.DesignVectorTextPathResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
					"krój %s nie ma ani jednego glifu z tego tekstu — linii pisma nie da się zmierzyć",
					nazwaKroju))
			}
			wezly = []shared.DesignVectorNode{
				{X: poczatekX, Y: poczatekY, Kind: shared.DesignVectorNodeKindCorner},
				{X: poczatekX + szerokosc, Y: poczatekY, Kind: shared.DesignVectorNodeKindCorner},
			}
		}
		wiersz, err := a.zapiszSciezkeDesignu(ctx, kompozycja.ID,
			nowyIdentyfikator(przedrostekSciezkiDesign), 0, nil, &nazwa,
			wezly, zamknieta, z.Fill, z.Stroke)
		if err != nil {
			return shared.DesignVectorTextPathResponse{}, err
		}
		sciezka, err := sciezkaKontraktuDesignu(kompozycja.Kod, wiersz)
		if err != nil {
			return shared.DesignVectorTextPathResponse{}, bladWydaniaDesignu(err.Error())
		}
		return shared.DesignVectorTextPathResponse{Path: sciezka, Outlined: false}, nil
	}

	kontury, err := sciezkaTekstuDesignu(krojWczytany, nazwa, z.FontSize, poczatekX, poczatekY)
	if err != nil {
		return shared.DesignVectorTextPathResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
			"komenda design.vector.text.path krojem %s: %s", nazwaKroju, err.Error()))
	}
	wezly, zamknieta := wezlyZeSciezkiBibliotekiDesignu(kontury)
	if len(wezly) < najmniejWezlowSciezkiDesignu {
		return shared.DesignVectorTextPathResponse{}, bladWydaniaDesignu(
			"kontury tekstu nie dały ani jednego węzła")
	}
	wiersz, err := a.zapiszSciezkeDesignu(ctx, kompozycja.ID,
		nowyIdentyfikator(przedrostekSciezkiDesign), 0, nil, &nazwa,
		wezly, zamknieta, z.Fill, z.Stroke)
	if err != nil {
		return shared.DesignVectorTextPathResponse{}, err
	}
	sciezka, err := sciezkaKontraktuDesignu(kompozycja.Kod, wiersz)
	if err != nil {
		return shared.DesignVectorTextPathResponse{}, bladWydaniaDesignu(err.Error())
	}
	return shared.DesignVectorTextPathResponse{Path: sciezka, Outlined: true}, nil
}

// OczyscSciezki skraca zapis współrzędnych i oddaje ZMIERZONY ubytek bajtów —
// obsługuje `design.vector.optimize`.
//
// Ubytek jest pomiarem, nie oszacowaniem: rdzeń mierzy długość zapisu przed
// i po, więc `savedBytes` ujemne (zapis dłuższy, bo precyzja wyższa niż zastana)
// jest tu prawdą, a nie usterką rachunku. Kontrakt pole opisuje właśnie tak.
func (a *adapterDesignu) OczyscSciezki(ctx context.Context,
	z shared.DesignVectorOptimizeRequest) (shared.DesignVectorOptimizeResponse, error) {

	miejsca := domyslnaPrecyzjaSciezkiDesignu
	if z.Precision != nil {
		if *z.Precision < 0 || *z.Precision > granicaPrecyzjiSciezkiDesignu {
			return shared.DesignVectorOptimizeResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"komenda design.vector.optimize z precyzją %d: rdzeń przycina do zakresu 0–%d miejsc "+
					"po przecinku — powyżej zapis liczby przestaje nieść informację",
				*z.Precision, granicaPrecyzjiSciezkiDesignu))
		}
		miejsca = *z.Precision
	}

	wiersze := []dane.SciezkaWektorowaDesignu{}
	switch {
	case len(z.PathIds) > 0:
		for _, kod := range z.PathIds {
			wiersz, err := a.repozytorium.SciezkaWektorowaDesignuPoKodzie(ctx, strings.TrimSpace(kod))
			if err != nil {
				return shared.DesignVectorOptimizeResponse{}, bladNieznanejSciezkiDesignu(kod, err)
			}
			wiersze = append(wiersze, wiersz)
		}
	case z.BoardId != nil && strings.TrimSpace(*z.BoardId) != "":
		kompozycja, err := a.repozytorium.Kompozycja(ctx, strings.TrimSpace(*z.BoardId))
		if err != nil {
			return shared.DesignVectorOptimizeResponse{}, bladNieznanejKompozycjiDesignu(*z.BoardId, err)
		}
		wiersze, err = a.repozytorium.SciezkiWektoroweDesignu(ctx, kompozycja.ID, nil)
		if err != nil {
			return shared.DesignVectorOptimizeResponse{}, bladDesignu(err)
		}
	default:
		return shared.DesignVectorOptimizeResponse{}, bladWskazaniaDesignu(
			"komenda design.vector.optimize bez wskazania ścieżek ani kompozycji: rdzeń nie " +
				"przepisuje wszystkiego, co ma w bazie, na podstawie żądania bez zakresu")
	}
	if len(wiersze) == 0 {
		return shared.DesignVectorOptimizeResponse{}, bladWskazaniaDesignu(
			"wskazany zakres nie ma ani jednej ścieżki do oczyszczenia")
	}

	przed, po := 0, 0
	for _, wiersz := range wiersze {
		wezly, err := wezlyZeZapisuDesignu(wiersz.WezlyJSON)
		if err != nil {
			return shared.DesignVectorOptimizeResponse{}, bladWydaniaDesignu(err.Error())
		}
		przed += len(wiersz.WezlyJSON)
		oczyszczone := przytnijPrecyzjeWezlowDesignu(wezly, miejsca)
		zapis, err := zapisWezlowDesignu(oczyszczone)
		if err != nil {
			return shared.DesignVectorOptimizeResponse{}, bladWydaniaDesignu(err.Error())
		}
		po += len(zapis)
		if _, err := a.repozytorium.ZapiszSciezkeWektorowaDesignu(ctx, dane.SciezkaWektorowaDesignu{
			Kod:             wiersz.Kod,
			KompozycjaID:    wiersz.KompozycjaID,
			WarstwaKod:      wiersz.WarstwaKod,
			Nazwa:           wiersz.Nazwa,
			WezlyJSON:       zapis,
			Zamknieta:       wiersz.Zamknieta,
			WypelnienieJSON: wiersz.WypelnienieJSON,
			ObrysJSON:       wiersz.ObrysJSON,
			Kolejnosc:       wiersz.Kolejnosc,
		}); err != nil {
			return shared.DesignVectorOptimizeResponse{}, bladDesignu(err)
		}
	}
	return shared.DesignVectorOptimizeResponse{
		SizeBytes: po, SavedBytes: przed - po, PathsAffected: len(wiersze),
	}, nil
}

// UstawSymbol zakłada symbol albo nadpisuje zastany — obsługuje
// `design.vector.symbol.set`.
//
// `propagatedTo` niesie liczbę członków NAPRAWDĘ zapisanych, odczytaną z bazy po
// zapisie. Liczba z żądania byłaby obietnicą, a nie skutkiem.
func (a *adapterDesignu) UstawSymbol(ctx context.Context,
	z shared.DesignVectorSymbolSetRequest) (shared.DesignVectorSymbolSetResponse, error) {

	if strings.TrimSpace(z.BoardId) == "" {
		return shared.DesignVectorSymbolSetResponse{}, bladWskazaniaDesignu(
			"komenda design.vector.symbol.set bez wskazania kompozycji")
	}
	if strings.TrimSpace(z.Name) == "" {
		return shared.DesignVectorSymbolSetResponse{}, bladWskazaniaDesignu(
			"komenda design.vector.symbol.set bez nazwy symbolu")
	}
	if len(z.PathIds) == 0 && len(z.LayerIds) == 0 {
		return shared.DesignVectorSymbolSetResponse{}, bladWskazaniaDesignu(
			"komenda design.vector.symbol.set bez ani jednej ścieżki i bez ani jednej warstwy: " +
				"symbol pusty nie jest definicją do wielokrotnego użycia")
	}
	kompozycja, err := a.repozytorium.Kompozycja(ctx, strings.TrimSpace(z.BoardId))
	if err != nil {
		return shared.DesignVectorSymbolSetResponse{}, bladNieznanejKompozycjiDesignu(z.BoardId, err)
	}

	kod := nowyIdentyfikator(przedrostekSymboluDesign)
	if z.SymbolId != nil && strings.TrimSpace(*z.SymbolId) != "" {
		kod = strings.TrimSpace(*z.SymbolId)
		zastany, err := a.repozytorium.SymbolDesignuPoKodzie(ctx, kod)
		if err != nil {
			return shared.DesignVectorSymbolSetResponse{}, bladNieznanegoSymboluDesignu(kod, err)
		}
		if zastany.KompozycjaID != kompozycja.ID {
			return shared.DesignVectorSymbolSetResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"symbol %s leży na innej kompozycji niż %s", kod, kompozycja.Kod))
		}
	}

	// Ścieżki wchodzące w skład symbolu sprawdzamy PRZED zapisem: symbol
	// wskazujący ścieżkę, której nie ma, byłby definicją bez kształtu.
	czlonkowie := make([]dane.CzlonekSymbolyDesignu, 0, len(z.PathIds)+len(z.LayerIds))
	for numer, kodSciezki := range z.PathIds {
		wiersz, err := a.repozytorium.SciezkaWektorowaDesignuPoKodzie(ctx, strings.TrimSpace(kodSciezki))
		if err != nil {
			return shared.DesignVectorSymbolSetResponse{}, bladNieznanejSciezkiDesignu(kodSciezki, err)
		}
		if wiersz.KompozycjaID != kompozycja.ID {
			return shared.DesignVectorSymbolSetResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"ścieżka %s leży na innej kompozycji niż symbol", wiersz.Kod))
		}
		czlonkowie = append(czlonkowie, dane.CzlonekSymbolyDesignu{
			Rodzaj: dane.CzlonekSymboluSciezka, Kod: wiersz.Kod, Kolejnosc: numer + 1,
		})
	}
	for numer, kodWarstwy := range z.LayerIds {
		czlonkowie = append(czlonkowie, dane.CzlonekSymbolyDesignu{
			Rodzaj:    dane.CzlonekSymboluWarstwa,
			Kod:       strings.TrimSpace(kodWarstwy),
			Kolejnosc: len(z.PathIds) + numer + 1,
		})
	}

	zapisany, err := a.repozytorium.ZapiszSymbolDesignu(ctx, dane.SymbolDesignu{
		Kod: kod, KompozycjaID: kompozycja.ID, Nazwa: strings.TrimSpace(z.Name),
	}, czlonkowie)
	if err != nil {
		return shared.DesignVectorSymbolSetResponse{}, bladDesignu(err)
	}
	symbol, err := a.zlozSymbolDesignu(ctx, kompozycja.Kod, zapisany)
	if err != nil {
		return shared.DesignVectorSymbolSetResponse{}, bladDesignu(err)
	}
	return shared.DesignVectorSymbolSetResponse{Symbol: symbol, PropagatedTo: zapisany.Liczba}, nil
}

// Symbole zwraca symbole kompozycji — obsługuje `design.vector.symbol.list`.
func (a *adapterDesignu) Symbole(ctx context.Context,
	z shared.DesignVectorSymbolListRequest) (shared.DesignVectorSymbolListResponse, error) {

	if strings.TrimSpace(z.BoardId) == "" {
		return shared.DesignVectorSymbolListResponse{}, bladWskazaniaDesignu(
			"komenda design.vector.symbol.list bez wskazania kompozycji")
	}
	kompozycja, err := a.repozytorium.Kompozycja(ctx, strings.TrimSpace(z.BoardId))
	if err != nil {
		return shared.DesignVectorSymbolListResponse{}, bladNieznanejKompozycjiDesignu(z.BoardId, err)
	}
	wiersze, err := a.repozytorium.SymboleDesignu(ctx, kompozycja.ID)
	if err != nil {
		return shared.DesignVectorSymbolListResponse{}, bladDesignu(err)
	}
	symbole := make([]shared.DesignSymbol, 0, len(wiersze))
	for _, wiersz := range wiersze {
		symbol, err := a.zlozSymbolDesignu(ctx, kompozycja.Kod, wiersz)
		if err != nil {
			return shared.DesignVectorSymbolListResponse{}, bladDesignu(err)
		}
		symbole = append(symbole, symbol)
	}
	return shared.DesignVectorSymbolListResponse{Symbols: symbole, Total: len(symbole)}, nil
}

// zlozSymbolDesignu składa `DesignSymbol` kontraktu z wiersza symbolu i jego
// członków odczytanych osobno — tak, jak dzieli je schemat.
func (a *adapterDesignu) zlozSymbolDesignu(ctx context.Context, kompozycja string,
	symbol dane.SymbolDesignu) (shared.DesignSymbol, error) {

	czlonkowie, err := a.repozytorium.CzlonkowieSymbolyDesignu(ctx, symbol.ID)
	if err != nil {
		return shared.DesignSymbol{}, err
	}
	wynik := shared.DesignSymbol{
		Id: symbol.Kod, BoardId: kompozycja, Name: symbol.Nazwa,
		InstanceCount: &symbol.Liczba,
	}
	for _, czlonek := range czlonkowie {
		switch czlonek.Rodzaj {
		case dane.CzlonekSymboluSciezka:
			wynik.PathIds = append(wynik.PathIds, czlonek.Kod)
		case dane.CzlonekSymboluWarstwa:
			wynik.LayerIds = append(wynik.LayerIds, czlonek.Kod)
		}
	}
	return wynik, nil
}

// WydajWektor wydaje ścieżki kompozycji jako SVG, PDF albo EPS — obsługuje
// `design.vector.export`.
//
// Wydanie wektorowe zostaje wektorem: rasteryzacja odebrałaby mu jedyną
// własność, dla której jest wektorem. PDF i EPS składa `tdewolff/canvas`
// wkompilowany w binarium, SVG — sklejenie dokumentu tutaj, tą samą drogą, co
// wyrys kompozycji (`adapter_modul_design_wyrys.go`).
func (a *adapterDesignu) WydajWektor(ctx context.Context,
	z shared.DesignVectorExportRequest) (shared.DesignVectorExportResponse, error) {

	if strings.TrimSpace(z.BoardId) == "" {
		return shared.DesignVectorExportResponse{}, bladWskazaniaDesignu(
			"komenda design.vector.export bez wskazania kompozycji")
	}
	if err := sprawdzWyliczenieDesignu("design.vector.export", "target", z.Target,
		shared.WartosciDesignVectorExportTarget()); err != nil {
		return shared.DesignVectorExportResponse{}, err
	}
	kompozycja, err := a.repozytorium.Kompozycja(ctx, strings.TrimSpace(z.BoardId))
	if err != nil {
		return shared.DesignVectorExportResponse{}, bladNieznanejKompozycjiDesignu(z.BoardId, err)
	}

	wiersze, err := a.sciezkiDoWydaniaDesignu(ctx, kompozycja.ID, z.PathIds)
	if err != nil {
		return shared.DesignVectorExportResponse{}, err
	}
	if len(wiersze) == 0 {
		return shared.DesignVectorExportResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
			"kompozycja %s nie ma ani jednej ścieżki wektorowej — rdzeń odmawia zamiast oddać "+
				"pusty dokument, bo plik pusty wygląda tak samo jak plik uszkodzony", kompozycja.Kod))
	}

	// Ramka zawężająca wydanie: jej prostokąt jest kadrem. Brak ramki bierze
	// prostokąt obejmujący wszystkie ścieżki.
	kadr, err := a.kadrWydaniaWektoraDesignu(ctx, kompozycja.ID, z.FrameId, wiersze)
	if err != nil {
		return shared.DesignVectorExportResponse{}, err
	}

	nazwa := oczyscNazwePlikuDesignu(nazwaKompozycjiDoPlikuDesignu(kompozycja)) + "." +
		string(z.Target)
	bajty, typTresci, err := zlozWydanieWektoraDesignu(z.Target, wiersze, kadr)
	if err != nil {
		return shared.DesignVectorExportResponse{}, bladWydaniaDesignu(err.Error())
	}
	return shared.DesignVectorExportResponse{
		ContentBase64: wBaza64Designu(bajty),
		FileName:      nazwa,
		MediaType:     typTresci,
		SizeBytes:     len(bajty),
	}, nil
}

// sciezkiDoWydaniaDesignu wybiera ścieżki objęte wydaniem: wskazane albo
// wszystkie ścieżki kompozycji.
func (a *adapterDesignu) sciezkiDoWydaniaDesignu(ctx context.Context, kompozycjaID int64,
	wskazane []string) ([]dane.SciezkaWektorowaDesignu, error) {

	if len(wskazane) == 0 {
		wiersze, err := a.repozytorium.SciezkiWektoroweDesignu(ctx, kompozycjaID, nil)
		if err != nil {
			return nil, bladDesignu(err)
		}
		return wiersze, nil
	}
	wiersze := make([]dane.SciezkaWektorowaDesignu, 0, len(wskazane))
	for _, kod := range wskazane {
		wiersz, err := a.repozytorium.SciezkaWektorowaDesignuPoKodzie(ctx, strings.TrimSpace(kod))
		if err != nil {
			return nil, bladNieznanejSciezkiDesignu(kod, err)
		}
		if wiersz.KompozycjaID != kompozycjaID {
			return nil, bladWskazaniaDesignu(fmt.Sprintf(
				"ścieżka %s nie leży na wydawanej kompozycji", wiersz.Kod))
		}
		wiersze = append(wiersze, wiersz)
	}
	return wiersze, nil
}

// kadrWydaniaWektoraDesignu rozstrzyga prostokąt wydania.
func (a *adapterDesignu) kadrWydaniaWektoraDesignu(ctx context.Context, kompozycjaID int64,
	ramka *string, wiersze []dane.SciezkaWektorowaDesignu) (shared.DesignBoardRegion, error) {

	if ramka != nil && strings.TrimSpace(*ramka) != "" {
		wiersz, err := a.repozytorium.RamkaDesignuPoKodzie(ctx, strings.TrimSpace(*ramka))
		if err != nil {
			return shared.DesignBoardRegion{}, bladNieznanejRamkiDesignu(*ramka, err)
		}
		if wiersz.KompozycjaID != kompozycjaID {
			return shared.DesignBoardRegion{}, bladWskazaniaDesignu(fmt.Sprintf(
				"ramka %s leży na innej kompozycji niż wydawana", wiersz.Kod))
		}
		x, y := 0.0, 0.0
		if wiersz.X != nil {
			x = *wiersz.X
		}
		if wiersz.Y != nil {
			y = *wiersz.Y
		}
		return shared.DesignBoardRegion{
			X: x, Y: y, Width: wiersz.Szerokosc, Height: wiersz.Wysokosc,
		}, nil
	}

	pierwszy := true
	var lewa, gora, prawa, dol float64
	for _, wiersz := range wiersze {
		wezly, err := wezlyZeZapisuDesignu(wiersz.WezlyJSON)
		if err != nil {
			return shared.DesignBoardRegion{}, bladWydaniaDesignu(err.Error())
		}
		granice := sciezkaBibliotekiDesignu(wezly, wiersz.Zamknieta).Bounds()
		if pierwszy {
			lewa, gora = granice.X0, granice.Y0
			prawa, dol = granice.X1, granice.Y1
			pierwszy = false
			continue
		}
		lewa = mniejszaDesignu(lewa, granice.X0)
		gora = mniejszaDesignu(gora, granice.Y0)
		prawa = wiekszaDesignu(prawa, granice.X1)
		dol = wiekszaDesignu(dol, granice.Y1)
	}
	szerokosc, wysokosc := prawa-lewa, dol-gora
	if szerokosc <= 0 {
		szerokosc = domyslnyBokWyrysuDesignu
	}
	if wysokosc <= 0 {
		wysokosc = domyslnyBokWyrysuDesignu
	}
	return shared.DesignBoardRegion{X: lewa, Y: gora, Width: szerokosc, Height: wysokosc}, nil
}

// zlozWydanieWektoraDesignu składa bajty wydania w żądanej postaci.
func zlozWydanieWektoraDesignu(postac shared.DesignVectorExportTarget,
	wiersze []dane.SciezkaWektorowaDesignu,
	kadr shared.DesignBoardRegion) ([]byte, string, error) {

	if postac == shared.DesignVectorExportTargetSvg {
		tresc, err := zlozDokumentSvgWektoraDesignu(wiersze, kadr)
		if err != nil {
			return nil, "", err
		}
		return []byte(tresc), typTresciWydaniaDesignu("svg"), nil
	}

	// PDF i EPS mierzą stronę w punktach typograficznych i liczą oś Y od dołu.
	// Kompozycja liczy Y od góry, więc kształty jadą przez odbicie względem
	// wysokości kadru — bez tego wydanie byłoby lustrzanym odbiciem tego, co
	// Operator widzi na ekranie.
	plotno := canvas.New(kadr.Width, kadr.Height)
	kontekst := canvas.NewContext(plotno)
	odbicie := canvas.Identity.Translate(-kadr.X, kadr.Height+kadr.Y).Scale(1, -1)
	for _, wiersz := range wiersze {
		wezly, err := wezlyZeZapisuDesignu(wiersz.WezlyJSON)
		if err != nil {
			return nil, "", err
		}
		sciezka := sciezkaBibliotekiDesignu(wezly, wiersz.Zamknieta).Transform(odbicie)
		styl := stylWydaniaWektoraDesignu(wypelnienieZeZapisuDesignu(wiersz.WypelnienieJSON),
			obrysZeZapisuDesignu(wiersz.ObrysJSON))
		kontekst.RenderPath(sciezka, styl, canvas.Identity)
	}

	var bufor bytes.Buffer
	if postac == shared.DesignVectorExportTargetPdf {
		wydawca := renderPdf.New(&bufor, kadr.Width, kadr.Height, nil)
		plotno.RenderTo(wydawca)
		if err := wydawca.Close(); err != nil {
			return nil, "", fmt.Errorf("nie można zapisać wydania wektorowego jako pdf: %w", err)
		}
		return bufor.Bytes(), typTresciWydaniaDesignu("pdf"), nil
	}
	wydawca := renderPs.New(&bufor, kadr.Width, kadr.Height,
		&renderPs.Options{Format: renderPs.EncapsulatedPostScript})
	plotno.RenderTo(wydawca)
	if err := wydawca.Close(); err != nil {
		return nil, "", fmt.Errorf("nie można zapisać wydania wektorowego jako eps: %w", err)
	}
	return bufor.Bytes(), "application/postscript", nil
}

// stylWydaniaWektoraDesignu przekłada wypełnienie i obrys kontraktu na styl
// biblioteki.
//
// Ścieżka bez wypełnienia i bez obrysu dostaje obrys włoskowy: kształt bez
// żadnego z dwóch byłby w wydaniu niewidoczny, a plik, w którym nie widać
// niczego, wygląda identycznie jak plik uszkodzony.
func stylWydaniaWektoraDesignu(wypelnienie *shared.DesignFill,
	obrys *shared.DesignStroke) canvas.Style {

	styl := canvas.Style{FillRule: canvas.NonZero, StrokeCapper: canvas.ButtCap,
		StrokeJoiner: canvas.MiterJoin}
	maCokolwiek := false
	if wypelnienie != nil && wypelnienie.Color != nil {
		if barwa, err := rozpoznajBarweDesignu(*wypelnienie.Color); err == nil {
			styl.Fill = canvas.Paint{Color: barwaRgbaDesignu(barwa, wypelnienie.Opacity)}
			maCokolwiek = true
		}
	}
	if obrys != nil && obrys.Width > 0 {
		if barwa, err := rozpoznajBarweDesignu(obrys.Color); err == nil {
			styl.Stroke = canvas.Paint{Color: barwaRgbaDesignu(barwa, obrys.Opacity)}
			styl.StrokeWidth = obrys.Width
			styl.Dashes = obrys.Dash
			maCokolwiek = true
		}
	}
	if !maCokolwiek {
		styl.Stroke = canvas.Paint{Color: canvas.Black}
		styl.StrokeWidth = 1
	}
	return styl
}

// zlozDokumentSvgWektoraDesignu składa dokument SVG ze ścieżek kompozycji.
//
// Dokument jest samowystarczalny: nie odsyła do niczego na tej maszynie, bo ma
// być plikiem, który Operator wysyła dalej.
func zlozDokumentSvgWektoraDesignu(wiersze []dane.SciezkaWektorowaDesignu,
	kadr shared.DesignBoardRegion) (string, error) {

	var dokument strings.Builder
	dokument.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	fmt.Fprintf(&dokument,
		`<svg xmlns="http://www.w3.org/2000/svg" width="%g" height="%g" viewBox="%g %g %g %g">`+"\n",
		kadr.Width, kadr.Height, kadr.X, kadr.Y, kadr.Width, kadr.Height)
	for _, wiersz := range wiersze {
		wezly, err := wezlyZeZapisuDesignu(wiersz.WezlyJSON)
		if err != nil {
			return "", err
		}
		zapis := zapisSvgSciezkiDesignu(wezly, wiersz.Zamknieta)
		if zapis == "" {
			continue
		}
		fmt.Fprintf(&dokument, `  <path d="%s" %s/>`+"\n", zapis,
			atrybutyWygladuSvgDesignu(wypelnienieZeZapisuDesignu(wiersz.WypelnienieJSON),
				obrysZeZapisuDesignu(wiersz.ObrysJSON)))
	}
	dokument.WriteString("</svg>\n")
	return dokument.String(), nil
}

// atrybutyWygladuSvgDesignu składa atrybuty wypełnienia i obrysu elementu SVG.
func atrybutyWygladuSvgDesignu(wypelnienie *shared.DesignFill, obrys *shared.DesignStroke) string {
	czesci := []string{}
	if wypelnienie != nil && wypelnienie.Color != nil {
		if barwa, err := rozpoznajBarweDesignu(*wypelnienie.Color); err == nil {
			czesci = append(czesci, fmt.Sprintf(`fill="%s"`, barwa.Clamped().Hex()))
			if wypelnienie.Opacity != nil {
				czesci = append(czesci, fmt.Sprintf(`fill-opacity="%g"`, *wypelnienie.Opacity))
			}
		}
	} else {
		czesci = append(czesci, `fill="none"`)
	}
	if obrys != nil && obrys.Width > 0 {
		if barwa, err := rozpoznajBarweDesignu(obrys.Color); err == nil {
			czesci = append(czesci, fmt.Sprintf(`stroke="%s" stroke-width="%g"`,
				barwa.Clamped().Hex(), obrys.Width))
			if obrys.Opacity != nil {
				czesci = append(czesci, fmt.Sprintf(`stroke-opacity="%g"`, *obrys.Opacity))
			}
			if len(obrys.Dash) > 0 {
				kreski := make([]string, 0, len(obrys.Dash))
				for _, kreska := range obrys.Dash {
					kreski = append(kreski, fmt.Sprintf("%g", kreska))
				}
				czesci = append(czesci, fmt.Sprintf(`stroke-dasharray="%s"`, strings.Join(kreski, " ")))
			}
		}
	}
	if len(czesci) == 0 {
		return `fill="none" stroke="#000000" stroke-width="1"`
	}
	return strings.Join(czesci, " ")
}

// sciezkaKontraktuDesignu składa `DesignVectorPath` kontraktu z wiersza.
func sciezkaKontraktuDesignu(kompozycja string,
	wiersz dane.SciezkaWektorowaDesignu) (shared.DesignVectorPath, error) {

	wezly, err := wezlyZeZapisuDesignu(wiersz.WezlyJSON)
	if err != nil {
		return shared.DesignVectorPath{}, err
	}
	return shared.DesignVectorPath{
		Id: wiersz.Kod, BoardId: kompozycja, LayerId: wiersz.WarstwaKod,
		Nodes: wezly, Closed: wiersz.Zamknieta,
		Fill:   wypelnienieZeZapisuDesignu(wiersz.WypelnienieJSON),
		Stroke: obrysZeZapisuDesignu(wiersz.ObrysJSON),
		Name:   wiersz.Nazwa,
	}, nil
}

// sprawdzWezlySciezkiDesignu odrzuca węzły bezsensowne PRZED zapisem: jeden
// węzeł nie jest kształtem, a rodzaj spoza kontraktu odbiłby się od warunku
// schematu i wrócił jako awaria rdzenia oznaczona jako ponawialna.
func sprawdzWezlySciezkiDesignu(komenda string, wezly []shared.DesignVectorNode) error {
	if len(wezly) < najmniejWezlowSciezkiDesignu {
		return bladWskazaniaDesignu(fmt.Sprintf(
			"komenda %s z %d węzłami: ścieżka o mniej niż dwóch węzłach jest punktem, nie kształtem",
			komenda, len(wezly)))
	}
	for numer, wezel := range wezly {
		if err := sprawdzWyliczenieDesignu(komenda, fmt.Sprintf("nodes[%d].kind", numer),
			wezel.Kind, shared.WartosciDesignVectorNodeKind()); err != nil {
			return err
		}
	}
	return nil
}

// wskaznikNiepustegoDesignu oddaje wskaźnik wyłącznie dla wartości niepustej.
// Napis z samych odstępów jest tu brakiem, nie wartością: `layerId: "  "`
// wpisany w kolumnę dawałby warstwę o nazwie z odstępów, której nikt nie zna.
func wskaznikNiepustegoDesignu(pole *string) *string {
	if pole == nil {
		return nil
	}
	wartosc := strings.TrimSpace(*pole)
	if wartosc == "" {
		return nil
	}
	return &wartosc
}

// nazwaKompozycjiDoPlikuDesignu oddaje rdzeń nazwy pliku wydania — z nazwy
// kompozycji, gdy Operator ją nadał, albo z jej identyfikatora.
func nazwaKompozycjiDoPlikuDesignu(kompozycja dane.KompozycjaDesignu) string {
	if kompozycja.Nazwa != nil && strings.TrimSpace(*kompozycja.Nazwa) != "" {
		return strings.TrimSpace(*kompozycja.Nazwa)
	}
	return kompozycja.Kod
}

// uporzadkujBilansDesignu porządkuje wykaz bilansu i oddaje nil dla pustego —
// pole niewymagane kontraktu ma wtedy nie wejść do odpowiedzi wcale, a wykaz
// niepusty ma stałą kolejność, żeby dwa wywołania tej samej komendy nie
// różniły się porządkiem zastrzeżeń.
func uporzadkujBilansDesignu(wykaz []string) []string {
	if len(wykaz) == 0 {
		return nil
	}
	sort.Strings(wykaz)
	return wykaz
}

// bladNieznanejSciezkiDesignu, bladNieznanegoSymboluDesignu
// i bladNieznanejRamkiDesignu nazywają byt warsztatu, którego rdzeń nie zna.
// Kod `not_found`, nie `internal_error`: okno ma po tym poznać, że wskazanie
// w jego wykazie jest nieaktualne.
func bladNieznanejSciezkiDesignu(kod string, err error) error {
	if czyBrakZasobuDesignu(err) {
		return bladNieznanegoBytuDesignu("ścieżki wektorowej " + kod + " nie ma w tym rdzeniu")
	}
	return bladDesignu(err)
}

func bladNieznanegoSymboluDesignu(kod string, err error) error {
	if czyBrakZasobuDesignu(err) {
		return bladNieznanegoBytuDesignu("symbolu " + kod + " nie ma w tym rdzeniu")
	}
	return bladDesignu(err)
}

func bladNieznanejRamkiDesignu(kod string, err error) error {
	if czyBrakZasobuDesignu(err) {
		return bladNieznanegoBytuDesignu("ramki " + kod + " nie ma w tym rdzeniu")
	}
	return bladDesignu(err)
}
