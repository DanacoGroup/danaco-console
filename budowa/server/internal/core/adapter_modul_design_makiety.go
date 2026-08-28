// Odpowiedzialność pliku: warsztat makiety modułu Design — ramki, układ
// automatyczny, więzy responsywne, siatki, komponenty i prototyp; metody
// stoją na `*adapterDesignu`.
package core

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

const (
	// Przedrostki identyfikatorów bytów warsztatu makiety: ramka, komponent
	// i połączenie prototypu dzielą wspólną przestrzeń kodów, odróżnioną tym
	// przedrostkiem od pozostałych bytów kompozycji.
	przedrostekRamkiDesign      = "ramka-"
	przedrostekKomponentuDesign = "komponent-"
	przedrostekPolaczeniaDesign = "przejscie-"

	// domyslnaSzerokoscRamkiDesignu i domyslnaWysokoscRamkiDesignu wchodzą
	// wyłącznie wtedy, gdy Operator nie podał ani wymiarów, ani nastawy
	// urządzenia — jest to jedyna para liczb, jaką da się podać, gdy żądanie
	// nie mówi o rozmiarze niczego.
	domyslnaSzerokoscRamkiDesignu = 1440.0
	domyslnaWysokoscRamkiDesignu  = 1024.0
)

// nastawyUrzadzenDesignu to rozmiary ekranów znane rdzeniowi, czytane przez
// zakładanie ramki i przez wykaz ramek, żeby oba miejsca zgadzały się co do
// wymiarów jednej nastawy.
var nastawyUrzadzenDesignu = []struct {
	Nazwa     string
	Szerokosc float64
	Wysokosc  float64
}{
	{"telefon-maly", 320, 568},
	{"telefon", 390, 844},
	{"telefon-duzy", 430, 932},
	{"tablet", 768, 1024},
	{"tablet-duzy", 1024, 1366},
	{"laptop", 1440, 900},
	{"pulpit", 1920, 1080},
	{"pulpit-szeroki", 2560, 1440},
	{"zegarek", 184, 224},
}

// nastawaUrzadzeniaDesignu odnajduje nastawę urządzenia po nazwie bez względu
// na wielkość liter i oddaje jej szerokość oraz wysokość w pikselach.
func nastawaUrzadzeniaDesignu(nazwa string) (float64, float64, bool) {
	szukana := strings.ToLower(strings.TrimSpace(nazwa))
	for _, nastawa := range nastawyUrzadzenDesignu {
		if strings.ToLower(nastawa.Nazwa) == szukana {
			return nastawa.Szerokosc, nastawa.Wysokosc, true
		}
	}
	return 0, 0, false
}

// nazwyNastawUrzadzenDesignu oddaje nazwy nastaw w kolejności wykazu — wchodzą
// do odpowiedzi `design.frame.list` i do treści odmów.
func nazwyNastawUrzadzenDesignu() []string {
	nazwy := make([]string, 0, len(nastawyUrzadzenDesignu))
	for _, nastawa := range nastawyUrzadzenDesignu {
		nazwy = append(nazwy, nastawa.Nazwa)
	}
	return nazwy
}

// UstawRamke zakłada ramkę na kompozycji albo nadpisuje ramkę zastaną tego
// samego kodu — obsługuje `design.frame.set`.
func (a *adapterDesignu) UstawRamke(ctx context.Context,
	z shared.DesignFrameSetRequest) (shared.DesignFrameSetResponse, error) {

	if strings.TrimSpace(z.BoardId) == "" {
		return shared.DesignFrameSetResponse{}, bladWskazaniaDesignu(
			"komenda design.frame.set bez wskazania kompozycji")
	}
	if strings.TrimSpace(z.Name) == "" {
		return shared.DesignFrameSetResponse{}, bladWskazaniaDesignu(
			"komenda design.frame.set bez nazwy ramki")
	}
	if z.Layout != nil {
		if err := sprawdzUkladDesignu("design.frame.set", *z.Layout); err != nil {
			return shared.DesignFrameSetResponse{}, err
		}
	}

	kompozycja, err := a.repozytorium.Kompozycja(ctx, strings.TrimSpace(z.BoardId))
	if err != nil {
		return shared.DesignFrameSetResponse{}, bladNieznanejKompozycjiDesignu(z.BoardId, err)
	}

	// Wskazanie wymiaru wprost bije nastawę urządzenia: jawna liczba wygrywa
	// z nazwą nastawy.
	szerokosc, wysokosc := 0.0, 0.0
	if z.DevicePreset != nil && strings.TrimSpace(*z.DevicePreset) != "" {
		nastawaSzerokosc, nastawaWysokosc, znana := nastawaUrzadzeniaDesignu(*z.DevicePreset)
		if !znana {
			return shared.DesignFrameSetResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"komenda design.frame.set z nastawą urządzenia %q, której rdzeń nie zna; nastawy znane: %s",
				*z.DevicePreset, strings.Join(nazwyNastawUrzadzenDesignu(), ", ")))
		}
		szerokosc, wysokosc = nastawaSzerokosc, nastawaWysokosc
	}
	if z.Width != nil && *z.Width > 0 {
		szerokosc = *z.Width
	}
	if z.Height != nil && *z.Height > 0 {
		wysokosc = *z.Height
	}

	kod := nowyIdentyfikator(przedrostekRamkiDesign)
	if z.FrameId != nil && strings.TrimSpace(*z.FrameId) != "" {
		kod = strings.TrimSpace(*z.FrameId)
		zastana, err := a.repozytorium.RamkaDesignuPoKodzie(ctx, kod)
		if err != nil {
			return shared.DesignFrameSetResponse{}, bladNieznanejRamkiDesignu(kod, err)
		}
		if zastana.KompozycjaID != kompozycja.ID {
			return shared.DesignFrameSetResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"ramka %s leży na innej kompozycji niż %s — komenda design.frame.set nie przenosi "+
					"ramek między planszami", kod, kompozycja.Kod))
		}
		// Zmiana ramki bez podanych wymiarów zostawia wymiary dotychczasowe,
		// nie domyślne.
		if szerokosc <= 0 {
			szerokosc = zastana.Szerokosc
		}
		if wysokosc <= 0 {
			wysokosc = zastana.Wysokosc
		}
	}
	if szerokosc <= 0 {
		szerokosc = domyslnaSzerokoscRamkiDesignu
	}
	if wysokosc <= 0 {
		wysokosc = domyslnaWysokoscRamkiDesignu
	}

	siatkaZapis, err := zapisJsonDesignu(z.Grid)
	if err != nil {
		return shared.DesignFrameSetResponse{}, bladWydaniaDesignu(err.Error())
	}
	ukladZapis, err := zapisJsonDesignu(z.Layout)
	if err != nil {
		return shared.DesignFrameSetResponse{}, bladWydaniaDesignu(err.Error())
	}

	zapisana, err := a.repozytorium.ZapiszRamkeDesignu(ctx, dane.RamkaDesignu{
		Kod:               kod,
		KompozycjaID:      kompozycja.ID,
		Nazwa:             strings.TrimSpace(z.Name),
		Szerokosc:         szerokosc,
		Wysokosc:          wysokosc,
		X:                 z.X,
		Y:                 z.Y,
		NastawaUrzadzenia: wskaznikNiepustegoDesignu(z.DevicePreset),
		SiatkaJSON:        siatkaZapis,
		UkladJSON:         ukladZapis,
	})
	if err != nil {
		return shared.DesignFrameSetResponse{}, bladDesignu(err)
	}
	return shared.DesignFrameSetResponse{Frame: ramkaKontraktuDesignu(kompozycja.Kod, zapisana)}, nil
}

// Ramki zwraca ramki kompozycji wraz z nastawami urządzeń znanymi rdzeniowi —
// obsługuje `design.frame.list`.
func (a *adapterDesignu) Ramki(ctx context.Context,
	z shared.DesignFrameListRequest) (shared.DesignFrameListResponse, error) {

	if strings.TrimSpace(z.BoardId) == "" {
		return shared.DesignFrameListResponse{}, bladWskazaniaDesignu(
			"komenda design.frame.list bez wskazania kompozycji")
	}
	kompozycja, err := a.repozytorium.Kompozycja(ctx, strings.TrimSpace(z.BoardId))
	if err != nil {
		return shared.DesignFrameListResponse{}, bladNieznanejKompozycjiDesignu(z.BoardId, err)
	}
	wiersze, err := a.repozytorium.RamkiDesignu(ctx, kompozycja.ID)
	if err != nil {
		return shared.DesignFrameListResponse{}, bladDesignu(err)
	}
	ramki := make([]shared.DesignFrame, 0, len(wiersze))
	for _, wiersz := range wiersze {
		ramki = append(ramki, ramkaKontraktuDesignu(kompozycja.Kod, wiersz))
	}
	return shared.DesignFrameListResponse{
		Frames: ramki, Total: len(ramki), DevicePresets: nazwyNastawUrzadzenDesignu(),
	}, nil
}

// UsunRamke usuwa ramkę i zwalnia jej warstwy do `releasedLayerIds`, nie
// usuwając ich z kompozycji — obsługuje `design.frame.remove`.
func (a *adapterDesignu) UsunRamke(ctx context.Context,
	z shared.DesignFrameRemoveRequest) (shared.DesignFrameRemoveResponse, error) {

	if strings.TrimSpace(z.FrameId) == "" {
		return shared.DesignFrameRemoveResponse{}, bladWskazaniaDesignu(
			"komenda design.frame.remove bez wskazania ramki")
	}
	kod := strings.TrimSpace(z.FrameId)
	ramka, err := a.repozytorium.RamkaDesignuPoKodzie(ctx, kod)
	if err != nil {
		if czyBrakZasobuDesignu(err) {
			// Ramki, której nie ma, nie odmawia się — pole `removed` mówi, czy
			// wiersz istniał.
			return shared.DesignFrameRemoveResponse{Removed: false}, nil
		}
		return shared.DesignFrameRemoveResponse{}, bladDesignu(err)
	}
	zwolnione, err := a.repozytorium.ZwolnijWarstwyRamkiDesignu(ctx, ramka.ID)
	if err != nil {
		return shared.DesignFrameRemoveResponse{}, bladDesignu(err)
	}
	usunieta, err := a.repozytorium.UsunRamkeDesignu(ctx, kod)
	if err != nil {
		return shared.DesignFrameRemoveResponse{}, bladDesignu(err)
	}
	odpowiedz := shared.DesignFrameRemoveResponse{Removed: usunieta}
	if len(zwolnione) > 0 {
		odpowiedz.ReleasedLayerIds = zwolnione
	}
	return odpowiedz, nil
}

// UlozAutomatycznie przelicza i zapisuje położenia warstw ramki układem
// automatycznym, biorąc kolejność wskazaną żądaniem albo kolejność
// przypisania do ramki — obsługuje `design.layout.auto`.
func (a *adapterDesignu) UlozAutomatycznie(ctx context.Context,
	z shared.DesignLayoutAutoRequest) (shared.DesignLayoutAutoResponse, error) {

	if strings.TrimSpace(z.FrameId) == "" {
		return shared.DesignLayoutAutoResponse{}, bladWskazaniaDesignu(
			"komenda design.layout.auto bez wskazania ramki")
	}
	if err := sprawdzUkladDesignu("design.layout.auto", z.Layout); err != nil {
		return shared.DesignLayoutAutoResponse{}, err
	}
	ramka, err := a.repozytorium.RamkaDesignuPoKodzie(ctx, strings.TrimSpace(z.FrameId))
	if err != nil {
		return shared.DesignLayoutAutoResponse{}, bladNieznanejRamkiDesignu(z.FrameId, err)
	}

	kody := z.LayerIds
	if len(kody) == 0 {
		kody, err = a.repozytorium.WarstwyRamkiDesignu(ctx, ramka.ID)
		if err != nil {
			return shared.DesignLayoutAutoResponse{}, bladDesignu(err)
		}
	}
	if len(kody) == 0 {
		return shared.DesignLayoutAutoResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
			"ramka %s nie ma ani jednej warstwy do ułożenia, a żądanie żadnej nie wskazało — "+
				"układ automatyczny nie ma czego przeliczyć", ramka.Kod))
	}

	warstwy := make([]dane.WarstwaKompozycji, 0, len(kody))
	for _, kod := range kody {
		warstwa, err := a.repozytorium.WarstwaKompozycjiDesignuPoKodzie(ctx, strings.TrimSpace(kod))
		if err != nil {
			return shared.DesignLayoutAutoResponse{}, bladNieznanejWarstwyDesignu(kod, err)
		}
		if warstwa.KompozycjaID != ramka.KompozycjaID {
			return shared.DesignLayoutAutoResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"warstwa %s leży na innej kompozycji niż ramka %s", warstwa.Kod, ramka.Kod))
		}
		warstwy = append(warstwy, warstwa)
	}

	szerokoscTresci, wysokoscTresci := ulozWarstwyDesignu(ramka, z.Layout, warstwy)
	for _, warstwa := range warstwy {
		if err := a.repozytorium.PrzestawWarstweKompozycjiDesignu(ctx, warstwa); err != nil {
			return shared.DesignLayoutAutoResponse{}, bladDesignu(err)
		}
	}
	// Nastawy układu zapisuje się przy ramce, żeby zmiana rozmiaru mogła je wziąć
	// bez powtórnego żądania.
	ukladZapis, err := zapisJsonDesignu(&z.Layout)
	if err != nil {
		return shared.DesignLayoutAutoResponse{}, bladWydaniaDesignu(err.Error())
	}
	ramka.UkladJSON = ukladZapis
	if _, err := a.repozytorium.ZapiszRamkeDesignu(ctx, ramka); err != nil {
		return shared.DesignLayoutAutoResponse{}, bladDesignu(err)
	}

	// Warstwy oddaje się odczytane z bazy, nie policzone w pamięci, żeby mówić
	// o stanie po zapisie.
	odczytane := make([]shared.DesignBoardLayer, 0, len(warstwy))
	for _, warstwa := range warstwy {
		swieza, err := a.repozytorium.WarstwaKompozycjiDesignuPoKodzie(ctx, warstwa.Kod)
		if err != nil {
			return shared.DesignLayoutAutoResponse{}, bladDesignu(err)
		}
		odczytane = append(odczytane, przelozWarstwyKontraktu([]dane.WarstwaKompozycji{swieza})...)
	}
	return shared.DesignLayoutAutoResponse{
		Layers: odczytane, ContentWidth: &szerokoscTresci, ContentHeight: &wysokoscTresci,
	}, nil
}

// ulozWarstwyDesignu liczy położenia warstw w układzie automatycznym i oddaje
// wymiary treści po ułożeniu, zmieniając warstwy w miejscu; warstwa bez
// wymiarów zostaje zerowa, nie dostaje wymiaru zmyślonego.
func ulozWarstwyDesignu(ramka dane.RamkaDesignu, uklad shared.DesignAutoLayout,
	warstwy []dane.WarstwaKompozycji) (float64, float64) {

	odstep := 0.0
	if uklad.Gap != nil {
		odstep = *uklad.Gap
	}
	gora, prawo, dol, lewo := odstepyWewnetrzneDesignu(uklad)
	poziomo := uklad.Direction == shared.DesignLayoutDirectionHorizontal

	// Rozciągnięcie w poprzek kierunku układa warstwy na wymiar ramki
	// pomniejszony o odstępy wewnętrzne.
	poprzeczna := ramka.Wysokosc - gora - dol
	if poziomo {
		poprzeczna = ramka.Wysokosc - gora - dol
	} else {
		poprzeczna = ramka.Szerokosc - lewo - prawo
	}
	rozciagac := uklad.Align != nil && *uklad.Align == shared.DesignLayoutAlignStretch

	kursor := lewo
	if !poziomo {
		kursor = gora
	}
	najwiekszaPoprzeczna := 0.0
	for numer := range warstwy {
		warstwa := &warstwy[numer]
		szerokosc, wysokosc := 0.0, 0.0
		if warstwa.Szerokosc != nil {
			szerokosc = *warstwa.Szerokosc
		}
		if warstwa.Wysokosc != nil {
			wysokosc = *warstwa.Wysokosc
		}
		if rozciagac && poprzeczna > 0 {
			if poziomo {
				wysokosc = poprzeczna
				warstwa.Wysokosc = &wysokosc
			} else {
				szerokosc = poprzeczna
				warstwa.Szerokosc = &szerokosc
			}
		}

		if poziomo {
			x := kursor
			y := gora + wyrownaniePoprzeczneDesignu(uklad.Align, poprzeczna, wysokosc)
			warstwa.X, warstwa.Y = &x, &y
			kursor += szerokosc + odstep
			najwiekszaPoprzeczna = wiekszaDesignu(najwiekszaPoprzeczna, wysokosc)
			continue
		}
		x := lewo + wyrownaniePoprzeczneDesignu(uklad.Align, poprzeczna, szerokosc)
		y := kursor
		warstwa.X, warstwa.Y = &x, &y
		kursor += wysokosc + odstep
		najwiekszaPoprzeczna = wiekszaDesignu(najwiekszaPoprzeczna, szerokosc)
	}
	// Ostatni odstęp nie należy do treści — jest za ostatnią warstwą.
	if len(warstwy) > 0 {
		kursor -= odstep
	}
	if poziomo {
		return kursor - lewo, najwiekszaPoprzeczna
	}
	return najwiekszaPoprzeczna, kursor - gora
}

// odstepyWewnetrzneDesignu oddaje cztery odstępy układu; brak wskazania znaczy
// zero, a nie wartość „ładną" dobraną przez rdzeń.
func odstepyWewnetrzneDesignu(uklad shared.DesignAutoLayout) (float64, float64, float64, float64) {
	wartosc := func(pole *float64) float64 {
		if pole == nil {
			return 0
		}
		return *pole
	}
	return wartosc(uklad.PaddingTop), wartosc(uklad.PaddingRight),
		wartosc(uklad.PaddingBottom), wartosc(uklad.PaddingLeft)
}

// wyrownaniePoprzeczneDesignu liczy odsunięcie warstwy w poprzek kierunku
// układania, według wyrównania początku, środka albo końca.
func wyrownaniePoprzeczneDesignu(wyrownanie *shared.DesignLayoutAlign,
	dostepne, wlasne float64) float64 {

	if wyrownanie == nil || dostepne <= 0 {
		return 0
	}
	switch *wyrownanie {
	case shared.DesignLayoutAlignCenter:
		return (dostepne - wlasne) / 2
	case shared.DesignLayoutAlignEnd:
		return dostepne - wlasne
	}
	return 0
}

// UstawWiezy zapisuje więzy responsywne warstw ramki i oddaje liczbę więzi
// rzeczywiście zmienionych — obsługuje `design.constraint.set`.
func (a *adapterDesignu) UstawWiezy(ctx context.Context,
	z shared.DesignConstraintSetRequest) (shared.DesignConstraintSetResponse, error) {

	if strings.TrimSpace(z.FrameId) == "" {
		return shared.DesignConstraintSetResponse{}, bladWskazaniaDesignu(
			"komenda design.constraint.set bez wskazania ramki")
	}
	if len(z.Constraints) == 0 {
		return shared.DesignConstraintSetResponse{}, bladWskazaniaDesignu(
			"komenda design.constraint.set bez ani jednego więzu: żądanie puste nie jest " +
				"zdjęciem więzi, bo kontrakt nie ma pola, którym to powiedzieć")
	}
	for numer, wiez := range z.Constraints {
		if strings.TrimSpace(wiez.LayerId) == "" {
			return shared.DesignConstraintSetResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"więz numer %d bez wskazania warstwy", numer+1))
		}
		if err := sprawdzWyliczenieDesignu("design.constraint.set",
			fmt.Sprintf("constraints[%d].horizontal", numer), wiez.Horizontal,
			shared.WartosciDesignConstraintAnchor()); err != nil {
			return shared.DesignConstraintSetResponse{}, err
		}
		if err := sprawdzWyliczenieDesignu("design.constraint.set",
			fmt.Sprintf("constraints[%d].vertical", numer), wiez.Vertical,
			shared.WartosciDesignConstraintAnchor()); err != nil {
			return shared.DesignConstraintSetResponse{}, err
		}
	}

	ramka, err := a.repozytorium.RamkaDesignuPoKodzie(ctx, strings.TrimSpace(z.FrameId))
	if err != nil {
		return shared.DesignConstraintSetResponse{}, bladNieznanejRamkiDesignu(z.FrameId, err)
	}

	// Stan zastany odczytuje się przed zapisem: `changed` liczy zmiany
	// rzeczywiste, nie wiersze żądania.
	zastane, err := a.repozytorium.WiezyRamkiDesignu(ctx, ramka.ID)
	if err != nil {
		return shared.DesignConstraintSetResponse{}, bladDesignu(err)
	}
	poprzednie := make(map[string]dane.WiezRamkiDesignu, len(zastane))
	for _, wiez := range zastane {
		poprzednie[wiez.WarstwaKod] = wiez
	}

	wiezy := make([]dane.WiezRamkiDesignu, 0, len(z.Constraints))
	zmienionych := 0
	for _, wiez := range z.Constraints {
		nowy := dane.WiezRamkiDesignu{
			WarstwaKod: strings.TrimSpace(wiez.LayerId),
			Poziomo:    string(wiez.Horizontal),
			Pionowo:    string(wiez.Vertical),
		}
		stary, byl := poprzednie[nowy.WarstwaKod]
		if !byl || stary.Poziomo != nowy.Poziomo || stary.Pionowo != nowy.Pionowo {
			zmienionych++
		}
		wiezy = append(wiezy, nowy)
	}
	if err := a.repozytorium.ZapiszWiezyRamkiDesignu(ctx, ramka.ID, wiezy); err != nil {
		return shared.DesignConstraintSetResponse{}, bladDesignu(err)
	}

	// Odpowiedź niesie więzy odczytane z bazy, komplet obowiązujący ramkę, nie
	// echo żądania częściowego.
	po, err := a.repozytorium.WiezyRamkiDesignu(ctx, ramka.ID)
	if err != nil {
		return shared.DesignConstraintSetResponse{}, bladDesignu(err)
	}
	return shared.DesignConstraintSetResponse{
		Constraints: wiezyKontraktuDesignu(po), Changed: zmienionych,
	}, nil
}

// wiezyKontraktuDesignu przekłada wykaz więzów warstwy z postaci danych na
// wykaz więzów kontraktu API modułu Design.
func wiezyKontraktuDesignu(wiezy []dane.WiezRamkiDesignu) []shared.DesignConstraint {
	lista := make([]shared.DesignConstraint, 0, len(wiezy))
	for _, wiez := range wiezy {
		lista = append(lista, shared.DesignConstraint{
			LayerId:    wiez.WarstwaKod,
			Horizontal: shared.DesignConstraintAnchor(wiez.Poziomo),
			Vertical:   shared.DesignConstraintAnchor(wiez.Pionowo),
		})
	}
	return lista
}

// ZmienRozmiarRamki zmienia rozmiar ramki i przelicza warstwy z ich więzów —
// obsługuje `design.frame.resize.apply`.
func (a *adapterDesignu) ZmienRozmiarRamki(ctx context.Context,
	z shared.DesignFrameResizeApplyRequest) (shared.DesignFrameResizeApplyResponse, error) {

	if strings.TrimSpace(z.FrameId) == "" {
		return shared.DesignFrameResizeApplyResponse{}, bladWskazaniaDesignu(
			"komenda design.frame.resize.apply bez wskazania ramki")
	}
	if z.Width <= 0 || z.Height <= 0 {
		return shared.DesignFrameResizeApplyResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
			"komenda design.frame.resize.apply z rozmiarem %v×%v: ramka o niedodatnim boku nie "+
				"jest ekranem", z.Width, z.Height))
	}
	ramka, err := a.repozytorium.RamkaDesignuPoKodzie(ctx, strings.TrimSpace(z.FrameId))
	if err != nil {
		return shared.DesignFrameResizeApplyResponse{}, bladNieznanejRamkiDesignu(z.FrameId, err)
	}

	wiezy, err := a.repozytorium.WiezyRamkiDesignu(ctx, ramka.ID)
	if err != nil {
		return shared.DesignFrameResizeApplyResponse{}, bladDesignu(err)
	}
	kotwice := make(map[string]dane.WiezRamkiDesignu, len(wiezy))
	for _, wiez := range wiezy {
		kotwice[wiez.WarstwaKod] = wiez
	}
	kody, err := a.repozytorium.WarstwyRamkiDesignu(ctx, ramka.ID)
	if err != nil {
		return shared.DesignFrameResizeApplyResponse{}, bladDesignu(err)
	}

	staraSzerokosc, staraWysokosc := ramka.Szerokosc, ramka.Wysokosc
	warstwy := make([]shared.DesignBoardLayer, 0, len(kody))
	for _, kod := range kody {
		warstwa, err := a.repozytorium.WarstwaKompozycjiDesignuPoKodzie(ctx, kod)
		if err != nil {
			if czyBrakZasobuDesignu(err) {
				// Warstwa zniknęła z kompozycji: pomija się ją, żeby nie blokować
				// przeliczenia pozostałych warstw.
				continue
			}
			return shared.DesignFrameResizeApplyResponse{}, bladDesignu(err)
		}
		wiez, maWiez := kotwice[kod]
		if maWiez {
			przeliczWarstweWiezemDesignu(&warstwa, wiez, staraSzerokosc, staraWysokosc,
				z.Width, z.Height)
			if err := a.repozytorium.PrzestawWarstweKompozycjiDesignu(ctx, warstwa); err != nil {
				return shared.DesignFrameResizeApplyResponse{}, bladDesignu(err)
			}
			warstwa, err = a.repozytorium.WarstwaKompozycjiDesignuPoKodzie(ctx, kod)
			if err != nil {
				return shared.DesignFrameResizeApplyResponse{}, bladDesignu(err)
			}
		}
		warstwy = append(warstwy, przelozWarstwyKontraktu([]dane.WarstwaKompozycji{warstwa})...)
	}

	ramka.Szerokosc, ramka.Wysokosc = z.Width, z.Height
	// Nastawa urządzenia przestaje obowiązywać po ręcznej zmianie rozmiaru,
	// by ramka nie kłamała.
	if ramka.NastawaUrzadzenia != nil {
		if szerokosc, wysokosc, znana := nastawaUrzadzeniaDesignu(*ramka.NastawaUrzadzenia); znana &&
			(szerokosc != z.Width || wysokosc != z.Height) {
			ramka.NastawaUrzadzenia = nil
		}
	}
	zapisana, err := a.repozytorium.ZapiszRamkeDesignu(ctx, ramka)
	if err != nil {
		return shared.DesignFrameResizeApplyResponse{}, bladDesignu(err)
	}
	kompozycja, err := a.repozytorium.KompozycjaDesignuPoKluczu(ctx, ramka.KompozycjaID)
	if err != nil {
		return shared.DesignFrameResizeApplyResponse{}, bladDesignu(err)
	}
	return shared.DesignFrameResizeApplyResponse{
		Frame: ramkaKontraktuDesignu(kompozycja.Kod, zapisana), Layers: warstwy,
	}, nil
}

// przeliczWarstweWiezemDesignu liczy nowe położenie i rozmiar warstwy po
// zmianie rozmiaru ramki, według jej kotwicy poziomej i pionowej.
func przeliczWarstweWiezemDesignu(warstwa *dane.WarstwaKompozycji, wiez dane.WiezRamkiDesignu,
	staraSzerokosc, staraWysokosc, nowaSzerokosc, nowaWysokosc float64) {

	x, szerokosc := przeliczOsWiezemDesignu(warstwa.X, warstwa.Szerokosc,
		shared.DesignConstraintAnchor(wiez.Poziomo), staraSzerokosc, nowaSzerokosc)
	y, wysokosc := przeliczOsWiezemDesignu(warstwa.Y, warstwa.Wysokosc,
		shared.DesignConstraintAnchor(wiez.Pionowo), staraWysokosc, nowaWysokosc)
	warstwa.X, warstwa.Szerokosc = x, szerokosc
	warstwa.Y, warstwa.Wysokosc = y, wysokosc
}

// przeliczOsWiezemDesignu liczy jedną oś. Warstwa bez położenia zostaje bez
// położenia: kotwica nie ma od czego liczyć, a wpisanie zera przesunęłoby ją
// w naroże ramki.
func przeliczOsWiezemDesignu(polozenie, rozmiar *float64, kotwica shared.DesignConstraintAnchor,
	stary, nowy float64) (*float64, *float64) {

	if polozenie == nil || stary <= 0 {
		return polozenie, rozmiar
	}
	wlasny := 0.0
	if rozmiar != nil {
		wlasny = *rozmiar
	}
	poczatek := *polozenie
	koniec := stary - (poczatek + wlasny)

	switch kotwica {
	case shared.DesignConstraintAnchorStart:
		return polozenie, rozmiar
	case shared.DesignConstraintAnchorEnd:
		nowePolozenie := nowy - koniec - wlasny
		return &nowePolozenie, rozmiar
	case shared.DesignConstraintAnchorCenter:
		nowePolozenie := (nowy - wlasny) / 2
		return &nowePolozenie, rozmiar
	case shared.DesignConstraintAnchorStretch:
		nowyRozmiar := nowy - poczatek - koniec
		if nowyRozmiar < 0 {
			nowyRozmiar = 0
		}
		return polozenie, &nowyRozmiar
	case shared.DesignConstraintAnchorScale:
		wspolczynnik := nowy / stary
		nowePolozenie := poczatek * wspolczynnik
		if rozmiar == nil {
			return &nowePolozenie, nil
		}
		nowyRozmiar := wlasny * wspolczynnik
		return &nowePolozenie, &nowyRozmiar
	}
	return polozenie, rozmiar
}

// UstawSiatke zapisuje siatkę układu, ramki albo całej kompozycji, jedną
// nastawą — obsługuje `design.grid.set`.
func (a *adapterDesignu) UstawSiatke(ctx context.Context,
	z shared.DesignGridSetRequest) (shared.DesignGridSetResponse, error) {

	if err := sprawdzSiatkeDesignu(z.Grid); err != nil {
		return shared.DesignGridSetResponse{}, err
	}
	zapis, err := zapisJsonDesignu(&z.Grid)
	if err != nil {
		return shared.DesignGridSetResponse{}, bladWydaniaDesignu(err.Error())
	}
	if zapis == nil {
		return shared.DesignGridSetResponse{}, bladWskazaniaDesignu(
			"komenda design.grid.set bez ani jednej nastawy siatki: siatka bez kolumn, rynny, " +
				"marginesu i linii bazowej nie jest siatką")
	}

	switch {
	case z.FrameId != nil && strings.TrimSpace(*z.FrameId) != "":
		ramka, err := a.repozytorium.RamkaDesignuPoKodzie(ctx, strings.TrimSpace(*z.FrameId))
		if err != nil {
			return shared.DesignGridSetResponse{}, bladNieznanejRamkiDesignu(*z.FrameId, err)
		}
		ramka.SiatkaJSON = zapis
		if _, err := a.repozytorium.ZapiszRamkeDesignu(ctx, ramka); err != nil {
			return shared.DesignGridSetResponse{}, bladDesignu(err)
		}
	case z.BoardId != nil && strings.TrimSpace(*z.BoardId) != "":
		kompozycja, err := a.repozytorium.Kompozycja(ctx, strings.TrimSpace(*z.BoardId))
		if err != nil {
			return shared.DesignGridSetResponse{}, bladNieznanejKompozycjiDesignu(*z.BoardId, err)
		}
		if err := a.repozytorium.ZapiszSiatkeKompozycjiDesignu(ctx, kompozycja.ID, *zapis); err != nil {
			return shared.DesignGridSetResponse{}, bladDesignu(err)
		}
	default:
		return shared.DesignGridSetResponse{}, bladWskazaniaDesignu(
			"komenda design.grid.set bez wskazania ramki ani kompozycji: siatka musi mieć do " +
				"czego przylgnąć")
	}
	return shared.DesignGridSetResponse{Grid: z.Grid}, nil
}

// sprawdzSiatkeDesignu odrzuca nastawy bezsensowne PRZED zapisem: siatka
// o zerowej liczbie kolumn nie dzieli niczego, a rynna ujemna nakładałaby
// kolumny na siebie.
func sprawdzSiatkeDesignu(siatka shared.DesignGrid) error {
	if siatka.Columns != nil && *siatka.Columns <= 0 {
		return bladWskazaniaDesignu(fmt.Sprintf(
			"komenda design.grid.set z %d kolumnami: siatka o niedodatniej liczbie kolumn nie "+
				"dzieli niczego", *siatka.Columns))
	}
	for nazwa, wartosc := range map[string]*float64{
		"gutter": siatka.Gutter, "margin": siatka.Margin, "baseline": siatka.Baseline,
	} {
		if wartosc != nil && *wartosc < 0 {
			return bladWskazaniaDesignu(fmt.Sprintf(
				"komenda design.grid.set z ujemną wartością pola %s (%v)", nazwa, *wartosc))
		}
	}
	return nil
}

// sprawdzUkladDesignu odrzuca nastawy układu automatycznego spoza kontraktu,
// zanim zapisze się je do ramki.
func sprawdzUkladDesignu(komenda string, uklad shared.DesignAutoLayout) error {
	if err := sprawdzWyliczenieDesignu(komenda, "layout.direction", uklad.Direction,
		shared.WartosciDesignLayoutDirection()); err != nil {
		return err
	}
	if uklad.Align != nil {
		if err := sprawdzWyliczenieDesignu(komenda, "layout.align", *uklad.Align,
			shared.WartosciDesignLayoutAlign()); err != nil {
			return err
		}
	}
	if uklad.Gap != nil && *uklad.Gap < 0 {
		return bladWskazaniaDesignu(fmt.Sprintf(
			"komenda %s z ujemnym odstępem %v: warstwy nakładałyby się na siebie", komenda, *uklad.Gap))
	}
	return nil
}

// ZapiszKomponent utrwala komponent wraz z wariantami i oddaje w
// `propagatedTo` liczbę instancji odczytaną z bazy po zapisie — obsługuje
// `design.component.save`.
func (a *adapterDesignu) ZapiszKomponent(ctx context.Context,
	z shared.DesignComponentSaveRequest) (shared.DesignComponentSaveResponse, error) {

	if strings.TrimSpace(z.WindowId) == "" {
		return shared.DesignComponentSaveResponse{}, bladWskazaniaDesignu(
			"komenda design.component.save bez wskazania okna")
	}
	if strings.TrimSpace(z.Name) == "" {
		return shared.DesignComponentSaveResponse{}, bladWskazaniaDesignu(
			"komenda design.component.save bez nazwy komponentu")
	}
	for numer, wariant := range z.Variants {
		if strings.TrimSpace(wariant.Name) == "" {
			return shared.DesignComponentSaveResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"wariant numer %d bez nazwy: wariant bez nazwy nie da się wskazać przy zakładaniu "+
					"instancji", numer+1))
		}
	}

	kod := nowyIdentyfikator(przedrostekKomponentuDesign)
	if z.ComponentId != nil && strings.TrimSpace(*z.ComponentId) != "" {
		kod = strings.TrimSpace(*z.ComponentId)
		zastany, err := a.repozytorium.KomponentDesignuPoKodzie(ctx, kod)
		if err != nil {
			return shared.DesignComponentSaveResponse{}, bladNieznanegoKomponentuDesignu(kod, err)
		}
		if zastany.Okno != z.WindowId {
			return shared.DesignComponentSaveResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"komponent %s należy do okna %s, a komenda design.component.save przyszła z okna %s",
				kod, zastany.Okno, z.WindowId))
		}
	}

	// Zestaw żetonów wskazany a nieznany jest odmową, żeby komponent nie
	// czerpał barw znikąd.
	if z.TokenSetId != nil && strings.TrimSpace(*z.TokenSetId) != "" {
		if _, err := a.repozytorium.ZestawZetonowDesignuPoKodzie(ctx,
			strings.TrimSpace(*z.TokenSetId)); err != nil {
			if czyBrakZasobuDesignu(err) {
				return shared.DesignComponentSaveResponse{}, bladNieznanegoBytuDesignu(
					"zestawu żetonów " + *z.TokenSetId + " nie ma w tym rdzeniu")
			}
			return shared.DesignComponentSaveResponse{}, bladDesignu(err)
		}
	}

	wariantyZapis, err := zapisJsonDesignu(z.Variants)
	if err != nil {
		return shared.DesignComponentSaveResponse{}, bladWydaniaDesignu(err.Error())
	}
	zapisany, err := a.repozytorium.ZapiszKomponentDesignu(ctx, dane.KomponentDesignu{
		Kod:              kod,
		Okno:             z.WindowId,
		Nazwa:            strings.TrimSpace(z.Name),
		WariantyJSON:     wariantyZapis,
		ZestawZetonowKod: wskaznikNiepustegoDesignu(z.TokenSetId),
	})
	if err != nil {
		return shared.DesignComponentSaveResponse{}, bladDesignu(err)
	}
	liczba := zapisany.Liczba
	return shared.DesignComponentSaveResponse{
		Component: komponentKontraktuDesignu(zapisany), PropagatedTo: &liczba,
	}, nil
}

// Komponenty zwraca komponenty okna, wszystkie albo jeden wskazany kodem —
// obsługuje `design.component.list`.
func (a *adapterDesignu) Komponenty(ctx context.Context,
	z shared.DesignComponentListRequest) (shared.DesignComponentListResponse, error) {

	if strings.TrimSpace(z.WindowId) == "" {
		return shared.DesignComponentListResponse{}, bladWskazaniaDesignu(
			"komenda design.component.list bez wskazania okna")
	}
	wiersze, err := a.repozytorium.KomponentyDesignu(ctx, z.WindowId,
		wskaznikNiepustegoDesignu(z.ComponentId))
	if err != nil {
		return shared.DesignComponentListResponse{}, bladDesignu(err)
	}
	komponenty := make([]shared.DesignComponent, 0, len(wiersze))
	for _, wiersz := range wiersze {
		komponenty = append(komponenty, komponentKontraktuDesignu(wiersz))
	}
	return shared.DesignComponentListResponse{
		Components: komponenty, Total: len(komponenty),
	}, nil
}

// DolozInstancjeKomponentu zakłada warstwę instancji komponentu — obsługuje
// `design.component.instance.add`.
func (a *adapterDesignu) DolozInstancjeKomponentu(ctx context.Context,
	z shared.DesignComponentInstanceAddRequest) (shared.DesignComponentInstanceAddResponse, error) {

	if strings.TrimSpace(z.ComponentId) == "" {
		return shared.DesignComponentInstanceAddResponse{}, bladWskazaniaDesignu(
			"komenda design.component.instance.add bez wskazania komponentu")
	}
	if strings.TrimSpace(z.BoardId) == "" {
		return shared.DesignComponentInstanceAddResponse{}, bladWskazaniaDesignu(
			"komenda design.component.instance.add bez wskazania kompozycji")
	}
	komponent, err := a.repozytorium.KomponentDesignuPoKodzie(ctx, strings.TrimSpace(z.ComponentId))
	if err != nil {
		return shared.DesignComponentInstanceAddResponse{}, bladNieznanegoKomponentuDesignu(
			z.ComponentId, err)
	}
	kompozycja, err := a.repozytorium.Kompozycja(ctx, strings.TrimSpace(z.BoardId))
	if err != nil {
		return shared.DesignComponentInstanceAddResponse{}, bladNieznanejKompozycjiDesignu(
			z.BoardId, err)
	}

	warianty := wariantyKomponentuDesignu(komponent)
	wariant := ""
	if z.Variant != nil && strings.TrimSpace(*z.Variant) != "" {
		wariant = strings.TrimSpace(*z.Variant)
		znany := false
		for _, znanyWariant := range warianty {
			if znanyWariant.Name == wariant {
				znany = true
				break
			}
		}
		if !znany {
			nazwy := make([]string, 0, len(warianty))
			for _, znanyWariant := range warianty {
				nazwy = append(nazwy, znanyWariant.Name)
			}
			return shared.DesignComponentInstanceAddResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"komponent %s nie ma wariantu %q; warianty komponentu: %s",
				komponent.Kod, wariant, strings.Join(nazwy, ", ")))
		}
	} else if len(warianty) > 0 {
		wariant = warianty[0].Name
	}

	// Ramka wskazana a nieznana jest odmową, tak samo ramka z innej kompozycji.
	var ramka *dane.RamkaDesignu
	if z.FrameId != nil && strings.TrimSpace(*z.FrameId) != "" {
		wiersz, err := a.repozytorium.RamkaDesignuPoKodzie(ctx, strings.TrimSpace(*z.FrameId))
		if err != nil {
			return shared.DesignComponentInstanceAddResponse{}, bladNieznanejRamkiDesignu(*z.FrameId, err)
		}
		if wiersz.KompozycjaID != kompozycja.ID {
			return shared.DesignComponentInstanceAddResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"ramka %s leży na innej kompozycji niż %s", wiersz.Kod, kompozycja.Kod))
		}
		ramka = &wiersz
	}

	adnotacja := komponent.Nazwa
	if wariant != "" {
		adnotacja = komponent.Nazwa + " / " + wariant
	}
	warstwa, err := a.repozytorium.DolozWarstweKompozycjiDesignu(ctx, dane.WarstwaKompozycji{
		Kod:          nowyIdentyfikator(przedrostekWarstwyDesign),
		KompozycjaID: kompozycja.ID,
		X:            z.X,
		Y:            z.Y,
		Adnotacja:    &adnotacja,
	})
	if err != nil {
		return shared.DesignComponentInstanceAddResponse{}, bladDesignu(err)
	}

	var wariantWiersza *string
	if wariant != "" {
		wariantWiersza = &wariant
	}
	if err := a.repozytorium.ZapiszInstancjeKomponentuDesignu(ctx, dane.InstancjaKomponentuDesignu{
		KomponentID:  komponent.ID,
		KompozycjaID: kompozycja.ID,
		WarstwaKod:   warstwa.Kod,
		Wariant:      wariantWiersza,
	}); err != nil {
		return shared.DesignComponentInstanceAddResponse{}, bladDesignu(err)
	}
	if ramka != nil {
		if err := a.repozytorium.PrzypiszWarstwyDoRamkiDesignu(ctx, ramka.ID,
			[]string{warstwa.Kod}); err != nil {
			return shared.DesignComponentInstanceAddResponse{}, bladDesignu(err)
		}
	}

	// Komponent odczytuje się po zapisie instancji, żeby `instanceCount` niósł
	// liczbę po dołożeniu.
	po, err := a.repozytorium.KomponentDesignuPoKodzie(ctx, komponent.Kod)
	if err != nil {
		return shared.DesignComponentInstanceAddResponse{}, bladDesignu(err)
	}
	warstwy := przelozWarstwyKontraktu([]dane.WarstwaKompozycji{warstwa})
	return shared.DesignComponentInstanceAddResponse{
		Layer: warstwy[0], Component: komponentKontraktuDesignu(po),
	}, nil
}

// UstawPolaczeniePrototypu zakłada przejście między ramkami albo nadpisuje
// zastane — obsługuje `design.prototype.link.set`.
func (a *adapterDesignu) UstawPolaczeniePrototypu(ctx context.Context,
	z shared.DesignPrototypeLinkSetRequest) (shared.DesignPrototypeLinkSetResponse, error) {

	if strings.TrimSpace(z.BoardId) == "" {
		return shared.DesignPrototypeLinkSetResponse{}, bladWskazaniaDesignu(
			"komenda design.prototype.link.set bez wskazania kompozycji")
	}
	if err := sprawdzWyliczenieDesignu("design.prototype.link.set", "trigger", z.Trigger,
		shared.WartosciDesignPrototypeTrigger()); err != nil {
		return shared.DesignPrototypeLinkSetResponse{}, err
	}
	if err := sprawdzWyliczenieDesignu("design.prototype.link.set", "transition", z.Transition,
		shared.WartosciDesignPrototypeTransition()); err != nil {
		return shared.DesignPrototypeLinkSetResponse{}, err
	}
	if z.DurationMs != nil && *z.DurationMs < 0 {
		return shared.DesignPrototypeLinkSetResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
			"komenda design.prototype.link.set z czasem przejścia %d ms: czas ujemny nie istnieje",
			*z.DurationMs))
	}

	kompozycja, err := a.repozytorium.Kompozycja(ctx, strings.TrimSpace(z.BoardId))
	if err != nil {
		return shared.DesignPrototypeLinkSetResponse{}, bladNieznanejKompozycjiDesignu(z.BoardId, err)
	}
	// Obie ramki sprawdza się przed zapisem, żeby przejście do brakującej nie
	// wyglądało jak gotowe.
	for pole, kod := range map[string]string{
		"fromFrameId": strings.TrimSpace(z.FromFrameId),
		"toFrameId":   strings.TrimSpace(z.ToFrameId),
	} {
		if kod == "" {
			return shared.DesignPrototypeLinkSetResponse{}, bladWskazaniaDesignu(
				"komenda design.prototype.link.set bez wskazania ramki w polu " + pole)
		}
		ramka, err := a.repozytorium.RamkaDesignuPoKodzie(ctx, kod)
		if err != nil {
			return shared.DesignPrototypeLinkSetResponse{}, bladNieznanejRamkiDesignu(kod, err)
		}
		if ramka.KompozycjaID != kompozycja.ID {
			return shared.DesignPrototypeLinkSetResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"ramka %s wskazana w polu %s leży na innej kompozycji niż %s",
				kod, pole, kompozycja.Kod))
		}
	}

	kod := nowyIdentyfikator(przedrostekPolaczeniaDesign)
	if z.LinkId != nil && strings.TrimSpace(*z.LinkId) != "" {
		kod = strings.TrimSpace(*z.LinkId)
		zastane, err := a.repozytorium.PolaczeniePrototypuDesignuPoKodzie(ctx, kod)
		if err != nil {
			return shared.DesignPrototypeLinkSetResponse{}, bladNieznanegoPolaczeniaDesignu(kod, err)
		}
		if zastane.KompozycjaID != kompozycja.ID {
			return shared.DesignPrototypeLinkSetResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"połączenie %s leży na innej kompozycji niż %s", kod, kompozycja.Kod))
		}
	}

	var czas *int64
	if z.DurationMs != nil {
		wartosc := int64(*z.DurationMs)
		czas = &wartosc
	}
	zapisane, err := a.repozytorium.ZapiszPolaczeniePrototypuDesignu(ctx,
		dane.PolaczeniePrototypuDesignu{
			Kod:          kod,
			KompozycjaID: kompozycja.ID,
			RamkaOdKod:   strings.TrimSpace(z.FromFrameId),
			RamkaDoKod:   strings.TrimSpace(z.ToFrameId),
			Wyzwalacz:    string(z.Trigger),
			Przejscie:    string(z.Transition),
			CzasMs:       czas,
			WarstwaKod:   wskaznikNiepustegoDesignu(z.LayerId),
		})
	if err != nil {
		return shared.DesignPrototypeLinkSetResponse{}, bladDesignu(err)
	}
	return shared.DesignPrototypeLinkSetResponse{
		Link: polaczenieKontraktuDesignu(zapisane),
	}, nil
}

// Prototyp zwraca ramki i przejścia kompozycji wraz z bilansem ramek
// osieroconych — obsługuje `design.prototype.get`.
func (a *adapterDesignu) Prototyp(ctx context.Context,
	z shared.DesignPrototypeGetRequest) (shared.DesignPrototypeGetResponse, error) {

	if strings.TrimSpace(z.BoardId) == "" {
		return shared.DesignPrototypeGetResponse{}, bladWskazaniaDesignu(
			"komenda design.prototype.get bez wskazania kompozycji")
	}
	kompozycja, err := a.repozytorium.Kompozycja(ctx, strings.TrimSpace(z.BoardId))
	if err != nil {
		return shared.DesignPrototypeGetResponse{}, bladNieznanejKompozycjiDesignu(z.BoardId, err)
	}
	wiersze, err := a.repozytorium.RamkiDesignu(ctx, kompozycja.ID)
	if err != nil {
		return shared.DesignPrototypeGetResponse{}, bladDesignu(err)
	}
	polaczenia, err := a.repozytorium.PolaczeniaPrototypuDesignu(ctx, kompozycja.ID)
	if err != nil {
		return shared.DesignPrototypeGetResponse{}, bladDesignu(err)
	}

	poczatkowa := ""
	if z.StartFrameId != nil && strings.TrimSpace(*z.StartFrameId) != "" {
		poczatkowa = strings.TrimSpace(*z.StartFrameId)
		znana := false
		for _, wiersz := range wiersze {
			if wiersz.Kod == poczatkowa {
				znana = true
				break
			}
		}
		if !znana {
			return shared.DesignPrototypeGetResponse{}, bladNieznanegoBytuDesignu(
				"ramki początkowej " + poczatkowa + " nie ma na kompozycji " + kompozycja.Kod)
		}
	}

	ramki := make([]shared.DesignFrame, 0, len(wiersze))
	for _, wiersz := range wiersze {
		ramki = append(ramki, ramkaKontraktuDesignu(kompozycja.Kod, wiersz))
	}
	przejscia := make([]shared.DesignPrototypeLink, 0, len(polaczenia))
	for _, polaczenie := range polaczenia {
		przejscia = append(przejscia, polaczenieKontraktuDesignu(polaczenie))
	}

	odpowiedz := shared.DesignPrototypeGetResponse{Frames: ramki, Links: przejscia}
	odpowiedz.UnreachableFrameIds = uporzadkujBilansDesignu(
		ramkiNieosiagalneDesignu(wiersze, polaczenia, poczatkowa))
	return odpowiedz, nil
}

// ramkiNieosiagalneDesignu oddaje ramki, do których nie da się dojść:
// przejściem po grafie od ramki początkowej, albo, bez niej, ramki bez
// połączenia wchodzącego.
func ramkiNieosiagalneDesignu(ramki []dane.RamkaDesignu,
	polaczenia []dane.PolaczeniePrototypuDesignu, poczatkowa string) []string {

	if len(ramki) == 0 {
		return nil
	}
	if poczatkowa == "" {
		wchodzace := map[string]bool{}
		for _, polaczenie := range polaczenia {
			wchodzace[polaczenie.RamkaDoKod] = true
		}
		osierocone := []string{}
		for _, ramka := range ramki {
			if !wchodzace[ramka.Kod] {
				osierocone = append(osierocone, ramka.Kod)
			}
		}
		// Wszystkie ramki bez połączeń to prototyp jeszcze niezbudowany, nie
		// zepsuty.
		if len(osierocone) == len(ramki) {
			return nil
		}
		return osierocone
	}

	nastepnicy := map[string][]string{}
	for _, polaczenie := range polaczenia {
		nastepnicy[polaczenie.RamkaOdKod] = append(nastepnicy[polaczenie.RamkaOdKod],
			polaczenie.RamkaDoKod)
	}
	osiagniete := map[string]bool{poczatkowa: true}
	kolejka := []string{poczatkowa}
	for len(kolejka) > 0 {
		biezaca := kolejka[0]
		kolejka = kolejka[1:]
		for _, nastepna := range nastepnicy[biezaca] {
			if osiagniete[nastepna] {
				continue
			}
			osiagniete[nastepna] = true
			kolejka = append(kolejka, nastepna)
		}
	}
	nieosiagalne := []string{}
	for _, ramka := range ramki {
		if !osiagniete[ramka.Kod] {
			nieosiagalne = append(nieosiagalne, ramka.Kod)
		}
	}
	return nieosiagalne
}

// UsunPolaczeniePrototypu usuwa przejście prototypu po jego kodzie —
// obsługuje `design.prototype.link.remove`.
func (a *adapterDesignu) UsunPolaczeniePrototypu(ctx context.Context,
	z shared.DesignPrototypeLinkRemoveRequest) (shared.DesignPrototypeLinkRemoveResponse, error) {

	if strings.TrimSpace(z.LinkId) == "" {
		return shared.DesignPrototypeLinkRemoveResponse{}, bladWskazaniaDesignu(
			"komenda design.prototype.link.remove bez wskazania połączenia")
	}
	usuniete, err := a.repozytorium.UsunPolaczeniePrototypuDesignu(ctx, strings.TrimSpace(z.LinkId))
	if err != nil {
		return shared.DesignPrototypeLinkRemoveResponse{}, bladDesignu(err)
	}
	return shared.DesignPrototypeLinkRemoveResponse{Removed: usuniete}, nil
}

// ramkaKontraktuDesignu składa `DesignFrame` kontraktu z wiersza ramki, wraz
// z jej siatką i układem, jeśli je ma.
func ramkaKontraktuDesignu(kompozycja string, wiersz dane.RamkaDesignu) shared.DesignFrame {
	ramka := shared.DesignFrame{
		Id: wiersz.Kod, BoardId: kompozycja, Name: wiersz.Nazwa,
		Width: wiersz.Szerokosc, Height: wiersz.Wysokosc,
		X: wiersz.X, Y: wiersz.Y, DevicePreset: wiersz.NastawaUrzadzenia,
	}
	if wiersz.SiatkaJSON != nil {
		var siatka shared.DesignGrid
		if err := json.Unmarshal([]byte(*wiersz.SiatkaJSON), &siatka); err == nil {
			ramka.Grid = &siatka
		}
	}
	if wiersz.UkladJSON != nil {
		var uklad shared.DesignAutoLayout
		if err := json.Unmarshal([]byte(*wiersz.UkladJSON), &uklad); err == nil {
			ramka.Layout = &uklad
		}
	}
	return ramka
}

// komponentKontraktuDesignu składa `DesignComponent` kontraktu z wiersza
// komponentu wraz z jego wariantami.
func komponentKontraktuDesignu(wiersz dane.KomponentDesignu) shared.DesignComponent {
	liczba := wiersz.Liczba
	return shared.DesignComponent{
		Id: wiersz.Kod, WindowId: wiersz.Okno, Name: wiersz.Nazwa,
		Variants: wariantyKomponentuDesignu(wiersz), TokenSetId: wiersz.ZestawZetonowKod,
		InstanceCount: &liczba,
	}
}

// wariantyKomponentuDesignu rozkłada zapis wariantów z kolumny. Zapis
// nieczytelny daje wykaz pusty, nie odmowę: warianty są cechą wtórną, a odmowa
// odczytu komponentu z ich powodu odebrałaby Operatorowi komponent.
func wariantyKomponentuDesignu(wiersz dane.KomponentDesignu) []shared.DesignComponentVariant {
	if wiersz.WariantyJSON == nil || strings.TrimSpace(*wiersz.WariantyJSON) == "" {
		return nil
	}
	var warianty []shared.DesignComponentVariant
	if err := json.Unmarshal([]byte(*wiersz.WariantyJSON), &warianty); err != nil {
		return nil
	}
	return warianty
}

// polaczenieKontraktuDesignu składa `DesignPrototypeLink` kontraktu z wiersza
// połączenia prototypu ramek.
func polaczenieKontraktuDesignu(wiersz dane.PolaczeniePrototypuDesignu) shared.DesignPrototypeLink {
	polaczenie := shared.DesignPrototypeLink{
		Id: wiersz.Kod, FromFrameId: wiersz.RamkaOdKod, ToFrameId: wiersz.RamkaDoKod,
		Trigger:    shared.DesignPrototypeTrigger(wiersz.Wyzwalacz),
		Transition: shared.DesignPrototypeTransition(wiersz.Przejscie),
		LayerId:    wiersz.WarstwaKod,
	}
	if wiersz.CzasMs != nil {
		czas := int(*wiersz.CzasMs)
		polaczenie.DurationMs = &czas
	}
	return polaczenie
}

// zapisJsonDesignu składa zapis JSON pola opcjonalnego do kolumny. Wartość nil
// i wykaz pusty dają brak zapisu — kolumna zostaje wtedy pusta, a nie wypełniona
// napisem „null", którego odczyt musiałby potem odróżniać od wartości.
func zapisJsonDesignu(wartosc any) (*string, error) {
	if wartosc == nil {
		return nil, nil
	}
	bajty, err := json.Marshal(wartosc)
	if err != nil {
		return nil, fmt.Errorf("nie można złożyć zapisu pola: %w", err)
	}
	zapis := string(bajty)
	if zapis == "null" || zapis == "{}" || zapis == "[]" {
		return nil, nil
	}
	return &zapis, nil
}

// bladNieznanejWarstwyDesignu, bladNieznanegoKomponentuDesignu
// i bladNieznanegoPolaczeniaDesignu nazywają byt makiety, którego rdzeń nie zna.
func bladNieznanejWarstwyDesignu(kod string, err error) error {
	if czyBrakZasobuDesignu(err) {
		return bladNieznanegoBytuDesignu("warstwy " + kod + " nie ma w tym rdzeniu")
	}
	return bladDesignu(err)
}

func bladNieznanegoKomponentuDesignu(kod string, err error) error {
	if czyBrakZasobuDesignu(err) {
		return bladNieznanegoBytuDesignu("komponentu " + kod + " nie ma w tym rdzeniu")
	}
	return bladDesignu(err)
}

func bladNieznanegoPolaczeniaDesignu(kod string, err error) error {
	if czyBrakZasobuDesignu(err) {
		return bladNieznanegoBytuDesignu("połączenia prototypu " + kod + " nie ma w tym rdzeniu")
	}
	return bladDesignu(err)
}

// KompozycjaZdarzenia odczytuje kompozycję po kodzie na potrzeby szyny zdarzeń.
// Powód istnienia tej metody stoi przy jej deklaracji w porcie
// (`adapter_modul_design_uchwyty.go`).
func (a *adapterDesignu) KompozycjaZdarzenia(ctx context.Context,
	kod string) (shared.DesignBoard, bool) {

	wiersz, err := a.repozytorium.Kompozycja(ctx, strings.TrimSpace(kod))
	if err != nil {
		return shared.DesignBoard{}, false
	}
	kompozycja, err := a.zlozBoard(ctx, wiersz)
	if err != nil {
		return shared.DesignBoard{}, false
	}
	return kompozycja, true
}

// KompozycjaRamkiZdarzenia odczytuje kompozycję, do której należy wskazana
// ramka, na potrzeby szyny zdarzeń.
func (a *adapterDesignu) KompozycjaRamkiZdarzenia(ctx context.Context,
	ramka string) (shared.DesignBoard, bool) {

	wiersz, err := a.repozytorium.RamkaDesignuPoKodzie(ctx, strings.TrimSpace(ramka))
	if err != nil {
		return shared.DesignBoard{}, false
	}
	kompozycja, err := a.repozytorium.KompozycjaDesignuPoKluczu(ctx, wiersz.KompozycjaID)
	if err != nil {
		return shared.DesignBoard{}, false
	}
	zlozona, err := a.zlozBoard(ctx, kompozycja)
	if err != nil {
		return shared.DesignBoard{}, false
	}
	return zlozona, true
}
