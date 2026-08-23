// Odpowiedzialność pliku: szablony materiału marketingowego —
// `design.template.save`, `design.template.list`, `design.template.apply`.
// Metody stoją na `*adapterDesignu` (`adapter_modul_design.go`); szablony
// promptu leżą w `adapter_modul_design_szablony.go` i są innym bytem.
//
// ── Szablon jest UKŁADEM, nie poleceniem ────────────────────────────────────
// Szablon materiału niesie rozmiar (baner 1200×628, wizytówka 90×50) i komplet
// warstw. Zastosowanie zakłada z niego kompozycję gotową do pracy — po to
// istnieje. Wykaz szablonów, z którego nie da się szablonu użyć, byłby spisem
// cudzej pracy.
//
// ── Podstawienie treści idzie po nazwie warstwy ─────────────────────────────
// Warstwa niesie adnotację (`note`) i to ona jest jej nazwą w szablonie:
// „logo", „nagłówek", „zdjęcie produktu". Podstawienie wskazuje zasób, który ma
// w tej warstwie stanąć. Nazwa, której szablon nie ma, wraca w `unmatchedNames`
// — bilans zamiast ciszy. Cicha zgoda oznaczałaby komplet kampanii złożony
// z szablonu, w którym połowa podstawień nie weszła, a Operator dowiedziałby
// się o tym dopiero z wydruku.
package core

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// przedrostekSzablonuMaterialuDesign znakuje identyfikatory zewnętrzne
// szablonów materiału — obok `szablon-promptu-` z obszaru Prompt Buildera.
const przedrostekSzablonuMaterialuDesign = "szablon-materialu-"

// przedrostekStronySzablonuDesign znakuje identyfikatory zewnętrzne stron
// publikacji wielostronicowej (migracja 339).
const przedrostekStronySzablonuDesign = "strona-szablonu-"

// ZapiszSzablonMaterialu utrwala szablon materiału wraz z warstwami —
// obsługuje `design.template.save`.
//
// Szablon wskazany a nieznany jest ODMOWĄ, nie cichym założeniem nowego —
// wzorem szablonu promptu: Operator, który nadpisuje, oczekuje że nadpisał ten
// jeden, a nie że dostał drugi obok.
func (a *adapterDesignu) ZapiszSzablonMaterialu(ctx context.Context,
	z shared.DesignTemplateSaveRequest) (shared.DesignTemplateSaveResponse, error) {

	if strings.TrimSpace(z.WindowId) == "" {
		return shared.DesignTemplateSaveResponse{}, bladWskazaniaDesignu(
			"komenda design.template.save bez wskazania okna")
	}
	if strings.TrimSpace(z.Name) == "" {
		return shared.DesignTemplateSaveResponse{}, bladWskazaniaDesignu(
			"komenda design.template.save bez nazwy szablonu")
	}
	if err := sprawdzRodzajSzablonuMaterialu(z.Kind); err != nil {
		return shared.DesignTemplateSaveResponse{}, err
	}
	if z.Width <= 0 || z.Height <= 0 {
		return shared.DesignTemplateSaveResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
			"komenda design.template.save z materiałem %v×%v: materiał o niedodatnim boku nie "+
				"istnieje, a kompozycja z niego założona nie miałaby czego pokazać", z.Width, z.Height))
	}

	kod := nowyIdentyfikator(przedrostekSzablonuMaterialuDesign)
	if z.TemplateId != nil && strings.TrimSpace(*z.TemplateId) != "" {
		kod = strings.TrimSpace(*z.TemplateId)
		zastany, err := a.repozytorium.SzablonMaterialuDesignuPoKodzie(ctx, kod)
		if err != nil {
			return shared.DesignTemplateSaveResponse{}, bladNieznanegoSzablonuMaterialu(kod, err)
		}
		if zastany.Okno != z.WindowId {
			return shared.DesignTemplateSaveResponse{}, bladWskazaniaDesignu(fmt.Sprintf(
				"szablon %s należy do okna %s, a komenda design.template.save przyszła z okna %s",
				kod, zastany.Okno, z.WindowId))
		}
	}

	zapisany, err := a.repozytorium.ZapiszSzablonMaterialuDesignu(ctx, dane.SzablonMaterialuDesignu{
		Kod:       kod,
		Okno:      z.WindowId,
		Nazwa:     strings.TrimSpace(z.Name),
		Rodzaj:    string(z.Kind),
		Szerokosc: z.Width,
		Wysokosc:  z.Height,
		Opis:      z.Description,
	}, przelozWarstwyDoZapisu(z.Layers))
	if err != nil {
		return shared.DesignTemplateSaveResponse{}, bladDesignu(err)
	}

	// Strony czynią z szablonu PUBLIKACJĘ. Zapis idzie po zapisie samego
	// szablonu, bo strona wskazuje jego klucz. Szablon bez stron zostaje
	// jednostronicowy i jego warstwy leżą tam, gdzie leżały — baner nie ma stron.
	if len(z.Pages) > 0 {
		strony, warstwyStron, err := stronySzablonuDoZapisuDesignu(z.Pages)
		if err != nil {
			return shared.DesignTemplateSaveResponse{}, err
		}
		if err := a.repozytorium.ZapiszStronySzablonuMaterialuDesignu(ctx, zapisany.ID,
			strony, warstwyStron); err != nil {
			return shared.DesignTemplateSaveResponse{}, bladDesignu(err)
		}
	}

	szablon, err := a.zlozSzablonMaterialu(ctx, zapisany)
	if err != nil {
		return shared.DesignTemplateSaveResponse{}, bladDesignu(err)
	}
	return shared.DesignTemplateSaveResponse{Template: szablon}, nil
}

// stronySzablonuDoZapisuDesignu przekłada strony kontraktu na wiersze wraz
// z warstwami przypisanymi do kodu strony.
//
// Numery stron są sprawdzane PRZED zapisem: numer niedodatni nie jest numerem
// strony, a numer powtórzony odbiłby się od unikatu schematu i wrócił jako
// awaria rdzenia oznaczona jako ponawialna — a to jest pomyłka wołającego.
func stronySzablonuDoZapisuDesignu(strony []shared.DesignTemplatePage) (
	[]dane.StronaSzablonuMaterialuDesignu, map[string][]dane.WarstwaKompozycji, error) {

	wiersze := make([]dane.StronaSzablonuMaterialuDesignu, 0, len(strony))
	warstwy := make(map[string][]dane.WarstwaKompozycji, len(strony))
	numery := map[int]bool{}
	for numer, strona := range strony {
		if strona.Number <= 0 {
			return nil, nil, bladWskazaniaDesignu(fmt.Sprintf(
				"strona numer %d w żądaniu ma numer %d: strony publikacji liczy się od jednego",
				numer+1, strona.Number))
		}
		if numery[strona.Number] {
			return nil, nil, bladWskazaniaDesignu(fmt.Sprintf(
				"numer strony %d występuje dwa razy — dwie strony o tym samym numerze dałyby "+
					"przy wydaniu kolejność zależną od porządku odczytu, czyli żadną", strona.Number))
		}
		numery[strona.Number] = true

		kod := nowyIdentyfikator(przedrostekStronySzablonuDesign)
		if strona.Id != nil && strings.TrimSpace(*strona.Id) != "" {
			kod = strings.TrimSpace(*strona.Id)
		}
		wiersze = append(wiersze, dane.StronaSzablonuMaterialuDesignu{
			Kod: kod, Numer: strona.Number, Nazwa: wskaznikNiepustegoDesignu(strona.Name),
		})
		warstwy[kod] = przelozWarstwyDoZapisu(strona.Layers)
	}
	return wiersze, warstwy, nil
}

// SzablonyMaterialu zwraca szablony okna — obsługuje `design.template.list`.
func (a *adapterDesignu) SzablonyMaterialu(ctx context.Context,
	z shared.DesignTemplateListRequest) (shared.DesignTemplateListResponse, error) {

	if strings.TrimSpace(z.WindowId) == "" {
		return shared.DesignTemplateListResponse{}, bladWskazaniaDesignu(
			"komenda design.template.list bez wskazania okna")
	}
	rodzaj := ""
	if z.Kind != nil {
		if err := sprawdzRodzajSzablonuMaterialu(*z.Kind); err != nil {
			return shared.DesignTemplateListResponse{}, err
		}
		rodzaj = string(*z.Kind)
	}

	wiersze, err := a.repozytorium.SzablonyMaterialuDesignu(ctx, z.WindowId, rodzaj)
	if err != nil {
		return shared.DesignTemplateListResponse{}, bladDesignu(err)
	}
	szablony := make([]shared.DesignTemplate, 0, len(wiersze))
	for _, wiersz := range wiersze {
		szablon, err := a.zlozSzablonMaterialu(ctx, wiersz)
		if err != nil {
			return shared.DesignTemplateListResponse{}, bladDesignu(err)
		}
		szablony = append(szablony, szablon)
	}
	return shared.DesignTemplateListResponse{Templates: szablony, Total: len(szablony)}, nil
}

// ZastosujSzablonMaterialu zakłada kompozycję z szablonu wraz z podstawieniem
// treści — obsługuje `design.template.apply`.
func (a *adapterDesignu) ZastosujSzablonMaterialu(ctx context.Context,
	z shared.DesignTemplateApplyRequest) (shared.DesignTemplateApplyResponse, error) {

	if strings.TrimSpace(z.TemplateId) == "" {
		return shared.DesignTemplateApplyResponse{}, bladWskazaniaDesignu(
			"komenda design.template.apply bez wskazania szablonu")
	}
	if strings.TrimSpace(z.WindowId) == "" {
		return shared.DesignTemplateApplyResponse{}, bladWskazaniaDesignu(
			"komenda design.template.apply bez wskazania okna")
	}

	podstawienia, err := podstawieniaSzablonu(z.Replacements)
	if err != nil {
		return shared.DesignTemplateApplyResponse{}, err
	}

	szablon, err := a.repozytorium.SzablonMaterialuDesignuPoKodzie(ctx, strings.TrimSpace(z.TemplateId))
	if err != nil {
		return shared.DesignTemplateApplyResponse{}, bladNieznanegoSzablonuMaterialu(z.TemplateId, err)
	}
	// Publikacja wielostronicowa: na kompozycję wchodzi JEDNA strona — ta
	// wskazana albo pierwsza. Wszystkie strony naraz na jednym płótnie leżałyby
	// jedna na drugiej, bo kompozycja jest arkuszem, nie plikiem.
	strony, err := a.repozytorium.StronySzablonuMaterialuDesignu(ctx, szablon.ID)
	if err != nil {
		return shared.DesignTemplateApplyResponse{}, bladDesignu(err)
	}
	warstwySzablonu := []dane.WarstwaKompozycji{}
	if len(strony) > 0 {
		wybrana := strony[0]
		if z.PageNumber != nil {
			znaleziona := false
			for _, strona := range strony {
				if strona.Numer == *z.PageNumber {
					wybrana, znaleziona = strona, true
					break
				}
			}
			if !znaleziona {
				numery := make([]string, 0, len(strony))
				for _, strona := range strony {
					numery = append(numery, fmt.Sprintf("%d", strona.Numer))
				}
				return shared.DesignTemplateApplyResponse{}, bladNieznanegoBytuDesignu(fmt.Sprintf(
					"szablon %s nie ma strony numer %d; strony szablonu: %s",
					szablon.Kod, *z.PageNumber, strings.Join(numery, ", ")))
			}
		}
		warstwySzablonu, err = a.repozytorium.WarstwyStronySzablonuMaterialuDesignu(ctx, wybrana.ID)
		if err != nil {
			return shared.DesignTemplateApplyResponse{}, bladDesignu(err)
		}
	} else {
		warstwySzablonu, err = a.repozytorium.WarstwySzablonuMaterialuDesignu(ctx, szablon.ID)
		if err != nil {
			return shared.DesignTemplateApplyResponse{}, bladDesignu(err)
		}
	}

	kodKompozycji := nowyIdentyfikator(przedrostekKompozycjiDesign)
	if z.BoardId != nil && strings.TrimSpace(*z.BoardId) != "" {
		kodKompozycji = strings.TrimSpace(*z.BoardId)
		if _, err := a.repozytorium.Kompozycja(ctx, kodKompozycji); err != nil {
			return shared.DesignTemplateApplyResponse{}, bladNieznanejKompozycjiDesignu(kodKompozycji, err)
		}
	}

	// Warstwy szablonu dostają WŁASNE identyfikatory w kompozycji: kompozycja
	// jest odtąd bytem osobnym i jej zmiana nie ma prawa ruszyć szablonu, z
	// którego powstała.
	uzyte := map[string]bool{}
	warstwy := make([]dane.WarstwaKompozycji, 0, len(warstwySzablonu))
	for numer, warstwa := range warstwySzablonu {
		nowa := dane.WarstwaKompozycji{
			Kod:         nowyIdentyfikator(przedrostekWarstwyDesign),
			ZasobID:     warstwa.ZasobID,
			X:           warstwa.X,
			Y:           warstwa.Y,
			Szerokosc:   warstwa.Szerokosc,
			Wysokosc:    warstwa.Wysokosc,
			Kolejnosc:   numer + 1,
			Zablokowana: warstwa.Zablokowana,
			Adnotacja:   warstwa.Adnotacja,
		}
		if warstwa.Adnotacja != nil {
			nazwa := strings.TrimSpace(*warstwa.Adnotacja)
			if zasob, jest := podstawienia[nazwa]; jest {
				uzyte[nazwa] = true
				wartosc := zasob
				nowa.ZasobID = &wartosc
			}
		}
		warstwy = append(warstwy, nowa)
	}

	nazwaKompozycji := szablon.Nazwa
	zapisana, err := a.repozytorium.ZapiszKompozycje(ctx, dane.KompozycjaDesignu{
		Kod:   kodKompozycji,
		Okno:  z.WindowId,
		Nazwa: &nazwaKompozycji,
	}, warstwy)
	if err != nil {
		return shared.DesignTemplateApplyResponse{}, bladDesignu(err)
	}
	kompozycja, err := a.zlozBoard(ctx, zapisana)
	if err != nil {
		return shared.DesignTemplateApplyResponse{}, bladDesignu(err)
	}

	odpowiedz := shared.DesignTemplateApplyResponse{Board: kompozycja}
	// Strony wracają w odpowiedzi, żeby okno publikacji wiedziało, ile stron
	// zostało — bez tego Operator dostawałby jedną kompozycję i nie miałby po
	// czym poznać, że publikacja ma jeszcze dwadzieścia trzy.
	if len(strony) > 0 {
		odpowiedz.Pages, err = a.stronyKontraktuSzablonuDesignu(ctx, strony)
		if err != nil {
			return shared.DesignTemplateApplyResponse{}, bladDesignu(err)
		}
	}
	niedopasowane := []string{}
	for nazwa := range podstawienia {
		if !uzyte[nazwa] {
			niedopasowane = append(niedopasowane, nazwa)
		}
	}
	if len(niedopasowane) > 0 {
		sort.Strings(niedopasowane)
		odpowiedz.UnmatchedNames = niedopasowane
	}
	return odpowiedz, nil
}

// podstawieniaSzablonu rozkłada `replacements` żądania na mapę nazwa warstwy →
// identyfikator zasobu.
//
// Wartość musi być napisem: podstawienie wskazuje ZASÓB, który ma stanąć
// w warstwie. Liczba ani obiekt nie jest identyfikatorem zasobu, a przepuszczone
// po cichu dałyby warstwę pustą przy powodzeniu komendy.
func podstawieniaSzablonu(surowe json.RawMessage) (map[string]string, error) {
	if len(surowe) == 0 {
		return map[string]string{}, nil
	}
	var odczytane map[string]any
	if err := json.Unmarshal(surowe, &odczytane); err != nil {
		return nil, bladWskazaniaDesignu(
			"komenda design.template.apply z podstawieniami, które nie są obiektem JSON " +
				"nazwa warstwy → identyfikator zasobu: " + err.Error())
	}
	podstawienia := make(map[string]string, len(odczytane))
	for nazwa, wartosc := range odczytane {
		tekst, jest := wartosc.(string)
		if !jest || strings.TrimSpace(tekst) == "" {
			return nil, bladWskazaniaDesignu(fmt.Sprintf(
				"podstawienie warstwy %q nie wskazuje zasobu — wartością podstawienia jest "+
					"identyfikator zasobu Assets Panelu, a nie %T", nazwa, wartosc))
		}
		podstawienia[strings.TrimSpace(nazwa)] = strings.TrimSpace(tekst)
	}
	return podstawienia, nil
}

// zlozSzablonMaterialu składa `DesignTemplate` kontraktu z wiersza szablonu
// i jego warstw odczytanych osobno — tak, jak dzieli je schemat.
func (a *adapterDesignu) zlozSzablonMaterialu(ctx context.Context,
	szablon dane.SzablonMaterialuDesignu) (shared.DesignTemplate, error) {

	warstwy, err := a.repozytorium.WarstwySzablonuMaterialuDesignu(ctx, szablon.ID)
	if err != nil {
		return shared.DesignTemplate{}, err
	}
	return shared.DesignTemplate{
		Id:          szablon.Kod,
		WindowId:    szablon.Okno,
		Name:        szablon.Nazwa,
		Kind:        shared.DesignTemplateKind(szablon.Rodzaj),
		Width:       szablon.Szerokosc,
		Height:      szablon.Wysokosc,
		Layers:      przelozWarstwyKontraktu(warstwy),
		Description: szablon.Opis,
	}, nil
}

// sprawdzRodzajSzablonuMaterialu odrzuca rodzaj spoza kontraktu PRZED zapisem.
// Wykaz pochodzi z kontraktu i nie jest tu przepisywany — wartość dołożona do
// kontraktu wchodzi do sprawdzenia sama.
func sprawdzRodzajSzablonuMaterialu(rodzaj shared.DesignTemplateKind) error {
	dopuszczalne := shared.WartosciDesignTemplateKind()
	for _, znany := range dopuszczalne {
		if rodzaj == znany {
			return nil
		}
	}
	nazwy := make([]string, 0, len(dopuszczalne))
	for _, znany := range dopuszczalne {
		nazwy = append(nazwy, string(znany))
	}
	return bladWskazaniaDesignu(fmt.Sprintf(
		"rodzaj szablonu %q, którego kontrakt nie zna; rodzaje dopuszczalne: %s",
		string(rodzaj), strings.Join(nazwy, ", ")))
}

// bladNieznanegoSzablonuMaterialu nazywa szablon, którego rdzeń nie zna.
func bladNieznanegoSzablonuMaterialu(kod string, err error) error {
	if czyBrakZasobuDesignu(err) {
		return bladNieznanegoBytuDesignu("szablonu materiału " + kod + " nie ma w tym rdzeniu")
	}
	return bladDesignu(err)
}

// stronyKontraktuSzablonuDesignu składa strony kontraktu wraz z ich warstwami
// odczytanymi osobno — tak, jak dzieli je schemat (migracja 339).
func (a *adapterDesignu) stronyKontraktuSzablonuDesignu(ctx context.Context,
	strony []dane.StronaSzablonuMaterialuDesignu) ([]shared.DesignTemplatePage, error) {

	wykaz := make([]shared.DesignTemplatePage, 0, len(strony))
	for _, strona := range strony {
		warstwy, err := a.repozytorium.WarstwyStronySzablonuMaterialuDesignu(ctx, strona.ID)
		if err != nil {
			return nil, err
		}
		kod := strona.Kod
		wykaz = append(wykaz, shared.DesignTemplatePage{
			Id: &kod, Number: strona.Numer, Name: strona.Nazwa,
			Layers: przelozWarstwyKontraktu(warstwy),
		})
	}
	return wykaz, nil
}
