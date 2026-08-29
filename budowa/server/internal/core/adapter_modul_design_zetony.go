// Odpowiedzialność pliku: zestawy żetonów systemu projektowego — obsługuje
// `design.tokenset.save`, `design.tokenset.list`, `design.tokenset.import`
// oraz rozpoznanie ról znanych systemowi produktu.
package core

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// przedrostekZestawuZetonowDesign znakuje identyfikatory zewnętrzne zestawów
// żetonów. Nowy identyfikator dostaje go przy zapisie bez wskazanego kodu
// zestawu, a przy zapisie z kodem podanym przez wołającego przedrostek nie
// wchodzi w grę.
const przedrostekZestawuZetonowDesign = "zestaw-zetonow-"

// rolePrzedrostekDesignu jest przedrostkiem własności niestandardowych systemu
// wizualnego produktu. Ten sam, którego używa panel żetonów po stronie klienta
// (`client/src/moduly/design/zetony-systemu.ts`).
const rolePrzedrostekDesignu = "--dn-"

// roleSystemuWizualnegoDesignu wylicza role SEMANTYCZNE systemu wizualnego
// produktu — te, po które wolno sięgać komponentom.
var roleSystemuWizualnegoDesignu = map[string]bool{
	// Powierzchnie
	"tlo": true, "powierzchnia": true, "powierzchnia-2": true, "panel": true, "nakladka": true,
	// Tekst
	"tekst": true, "tekst-2": true, "tekst-3": true, "tekst-inv": true,
	// Obrysy i wskazanie
	"obrys": true, "obrys-mocny": true, "obrys-subtelny": true,
	"hover": true, "wcisniecie": true, "fokus": true,
	// Sygnał — jedyna barwa akcentu
	"sygnal": true, "sygnal-mocny": true, "sygnal-tlo": true, "sygnal-obrys": true,
	"sygnal-wypelnienie": true, "sygnal-wypelnienie-hover": true, "kropka": true,
	// Stany
	"sukces-tlo": true, "sukces-obrys": true, "sukces-tekst": true,
	"ostrzezenie-tlo": true, "ostrzezenie-obrys": true, "ostrzezenie-tekst": true,
	"blad-tlo": true, "blad-obrys": true, "blad-tekst": true,
	"informacja-tlo": true, "informacja-obrys": true, "informacja-tekst": true,
	// Rama kokpitu i atrament
	"rama": true, "rama-tekst": true, "rama-tekst-2": true, "rama-hover": true,
	"rama-obrys": true, "atrament": true, "atrament-hover": true, "atrament-tekst": true,
	// Kroje pisma
	"ff-naglowek": true, "ff-bazowa": true, "ff-mono": true,
	// Stopnie pisma
	"fs-2xs": true, "fs-xs": true, "fs-sm": true, "fs-base": true, "fs-md": true,
	"fs-lg": true, "fs-xl": true, "fs-2xl": true, "fs-3xl": true, "fs-display": true,
	// Grubości i interlinia
	"fw-normalna": true, "fw-srednia": true, "fw-polgruba": true, "fw-gruba": true,
	"lh-ciasny": true, "lh-bazowy": true, "lh-luzny": true,
	// Odstępy liter
	"ls-naglowek": true, "ls-wersaliki": true, "ls-mono-wersaliki": true,
	// Przestrzeń
	"od-0": true, "od-1": true, "od-2": true, "od-3": true, "od-4": true, "od-5": true,
	"od-6": true, "od-8": true, "od-10": true, "od-12": true, "od-16": true,
	// Promienie
	"r-xs": true, "r-sm": true, "r-md": true, "r-lg": true, "r-xl": true, "r-pill": true,
	// Wymiary gęstości zwartej
	"wym-kontrolka": true, "wym-ikonowy": true, "wym-pasek": true, "wym-pas-kart": true,
	"wym-wiersz": true, "wym-boczna": true, "wym-pas-komunikacji": true, "wym-modal": true,
	"wym-check": true, "wym-kropka": true,
	// Ruch
	"czas-1": true, "czas-2": true, "czas-3": true, "czas-tetno": true,
	// Cienie
	"cien-1": true, "cien-2": true, "cien-3": true, "cien-lg": true, "cien-sygnal": true,
}

// ZapiszZestawZetonow utrwala zestaw żetonów — obsługuje
// `design.tokenset.save`, sprawdzając rodzaj każdego żetonu wobec kontraktu.
func (a *adapterDesignu) ZapiszZestawZetonow(ctx context.Context,
	z shared.DesignTokensetSaveRequest) (shared.DesignTokensetSaveResponse, error) {

	if z.WindowId == "" {
		return shared.DesignTokensetSaveResponse{}, bladWskazaniaDesignu(
			"komenda design.tokenset.save bez wskazania okna")
	}
	if strings.TrimSpace(z.Name) == "" {
		return shared.DesignTokensetSaveResponse{}, bladWskazaniaDesignu(
			"komenda design.tokenset.save bez nazwy zestawu")
	}
	if len(z.Tokens) == 0 {
		return shared.DesignTokensetSaveResponse{}, bladWskazaniaDesignu(
			"komenda design.tokenset.save bez ani jednego żetonu — zestaw pusty nie jest systemem " +
				"projektowym, a zapis nadpisałby zestaw zastany pustką")
	}
	for _, zeton := range z.Tokens {
		if err := sprawdzZetonDesignu("design.tokenset.save", zeton); err != nil {
			return shared.DesignTokensetSaveResponse{}, err
		}
	}

	kod := nowyIdentyfikator(przedrostekZestawuZetonowDesign)
	if z.TokenSetId != nil && strings.TrimSpace(*z.TokenSetId) != "" {
		kod = strings.TrimSpace(*z.TokenSetId)
		zastany, err := a.repozytorium.ZestawZetonowDesignuPoKodzie(ctx, kod)
		if err != nil {
			return shared.DesignTokensetSaveResponse{}, bladNieznanegoZestawuZetonowDesignu(kod, err)
		}
		if zastany.Okno != z.WindowId {
			return shared.DesignTokensetSaveResponse{}, bladNieznanegoBytuDesignu(
				"zestawu żetonów " + kod + " nie ma w oknie " + z.WindowId)
		}
	}

	zapisany, err := a.zapiszZestawZetonowDesignu(ctx, kod, z.WindowId,
		strings.TrimSpace(z.Name), z.Theme, z.Tokens)
	if err != nil {
		return shared.DesignTokensetSaveResponse{}, err
	}
	return shared.DesignTokensetSaveResponse{TokenSet: zapisany}, nil
}

// ZestawyZetonow zwraca zestawy żetonów okna wraz z ich żetonami — obsługuje
// `design.tokenset.list`. Zestawy wracają wraz z zawartością, bez osobnego
// wołania po żetony każdego z nich.
func (a *adapterDesignu) ZestawyZetonow(ctx context.Context,
	z shared.DesignTokensetListRequest) (shared.DesignTokensetListResponse, error) {

	if z.WindowId == "" {
		return shared.DesignTokensetListResponse{}, bladWskazaniaDesignu(
			"komenda design.tokenset.list bez wskazania okna")
	}
	wiersze, err := a.repozytorium.ZestawyZetonowDesignu(ctx, z.WindowId, z.TokenSetId)
	if err != nil {
		return shared.DesignTokensetListResponse{}, bladDesignu(err)
	}
	zestawy := make([]shared.DesignTokenSet, 0, len(wiersze))
	for _, wiersz := range wiersze {
		zetony, err := a.repozytorium.ZetonyZestawuDesignu(ctx, wiersz.ID)
		if err != nil {
			return shared.DesignTokensetListResponse{}, bladDesignu(err)
		}
		zestawy = append(zestawy, zestawZetonowKontraktuDesignu(wiersz, zetony))
	}
	return shared.DesignTokensetListResponse{TokenSets: zestawy, Total: len(zestawy)}, nil
}

// WczytajZestawZetonow zakłada zestaw z zapisu zewnętrznego — obsługuje
// `design.tokenset.import`, rozpoznając postać zapisu po treści, nie po
// nazwie pliku.
func (a *adapterDesignu) WczytajZestawZetonow(ctx context.Context,
	z shared.DesignTokensetImportRequest) (shared.DesignTokensetImportResponse, error) {

	if z.WindowId == "" {
		return shared.DesignTokensetImportResponse{}, bladWskazaniaDesignu(
			"komenda design.tokenset.import bez wskazania okna")
	}
	if strings.TrimSpace(z.Name) == "" {
		return shared.DesignTokensetImportResponse{}, bladWskazaniaDesignu(
			"komenda design.tokenset.import bez nazwy zakładanego zestawu")
	}
	if strings.TrimSpace(z.ContentBase64) == "" {
		return shared.DesignTokensetImportResponse{}, bladWskazaniaDesignu(
			"komenda design.tokenset.import bez treści zapisu")
	}
	bajty, err := base64.StdEncoding.DecodeString(z.ContentBase64)
	if err != nil {
		return shared.DesignTokensetImportResponse{}, bladWskazaniaDesignu(
			"treść zapisu żetonów nie jest poprawnym base64: " + err.Error())
	}
	if len(bajty) == 0 {
		return shared.DesignTokensetImportResponse{}, bladWskazaniaDesignu(
			"treść zapisu żetonów jest pusta")
	}

	zetony, err := rozlozZapisZetonowDesignu(string(bajty), z.Format)
	if err != nil {
		return shared.DesignTokensetImportResponse{}, bladWskazaniaDesignu(err.Error())
	}
	if len(zetony) == 0 {
		return shared.DesignTokensetImportResponse{}, bladWskazaniaDesignu(
			"w zapisie nie ma ani jednej roli, którą serwer umiałby odczytać — zestaw pusty " +
				"nie powstaje, bo wyglądałby na wczytany")
	}

	zapisany, err := a.zapiszZestawZetonowDesignu(ctx,
		nowyIdentyfikator(przedrostekZestawuZetonowDesign), z.WindowId,
		strings.TrimSpace(z.Name), nil, zetony)
	if err != nil {
		return shared.DesignTokensetImportResponse{}, err
	}

	odpowiedz := shared.DesignTokensetImportResponse{TokenSet: zapisany}
	nieznane := roleNieznaneSystemowiDesignu(zetony)
	if len(nieznane) > 0 {
		odpowiedz.UnknownNames = nieznane
	}
	return odpowiedz, nil
}

// zapiszZestawZetonowDesignu jest wspólną drogą zapisu dla `save` i `import` —
// obie utrwalają dokładnie to samo, więc druga ścieżka byłaby drugą prawdą
// o tym, czym jest zestaw.
func (a *adapterDesignu) zapiszZestawZetonowDesignu(ctx context.Context,
	kod, okno, nazwa string, motyw *string, zetony []shared.DesignToken) (shared.DesignTokenSet, error) {

	wiersze := make([]dane.ZetonDesignu, 0, len(zetony))
	for numer, zeton := range zetony {
		wiersze = append(wiersze, dane.ZetonDesignu{
			Nazwa:      strings.TrimSpace(zeton.Name),
			Rodzaj:     string(zeton.Kind),
			Wartosc:    zeton.Value,
			Opis:       zeton.Description,
			OdsylaczDo: zeton.AliasOf,
			Kolejnosc:  numer + 1,
		})
	}
	zapisany, err := a.repozytorium.ZapiszZestawZetonowDesignu(ctx, dane.ZestawZetonowDesignu{
		Kod: kod, Okno: okno, Nazwa: nazwa, Motyw: motyw,
	}, wiersze)
	if err != nil {
		return shared.DesignTokenSet{}, bladDesignu(err)
	}
	odczytane, err := a.repozytorium.ZetonyZestawuDesignu(ctx, zapisany.ID)
	if err != nil {
		return shared.DesignTokenSet{}, bladDesignu(err)
	}
	return zestawZetonowKontraktuDesignu(zapisany, odczytane), nil
}

// sprawdzZetonDesignu odrzuca żeton bez nazwy roli i żeton o rodzaju spoza
// kontraktu. Wykaz rodzajów pochodzi z kontraktu i nie jest tu przepisywany —
// rodzaj dołożony do kontraktu wchodzi do sprawdzenia sam.
func sprawdzZetonDesignu(komenda string, zeton shared.DesignToken) error {
	if strings.TrimSpace(zeton.Name) == "" {
		return bladWskazaniaDesignu(fmt.Sprintf(
			"komenda %s z żetonem bez nazwy roli — rola bez nazwy nie ma jak trafić do kodu", komenda))
	}
	dopuszczalne := shared.WartosciDesignTokenKind()
	for _, znany := range dopuszczalne {
		if zeton.Kind == znany {
			return nil
		}
	}
	nazwy := make([]string, 0, len(dopuszczalne))
	for _, znany := range dopuszczalne {
		nazwy = append(nazwy, string(znany))
	}
	return bladWskazaniaDesignu(fmt.Sprintf(
		"komenda %s z żetonem %q o rodzaju %q, którego kontrakt nie zna; rodzaje dopuszczalne: %s",
		komenda, zeton.Name, string(zeton.Kind), strings.Join(nazwy, ", ")))
}

// roleNieznaneSystemowiDesignu wylicza role zapisu, których system wizualny
// produktu nie nazywa. Kolejność jest kolejnością zapisu, bez powtórzeń.
func roleNieznaneSystemowiDesignu(zetony []shared.DesignToken) []string {
	widziane := map[string]bool{}
	nieznane := []string{}
	for _, zeton := range zetony {
		nazwa := strings.TrimPrefix(strings.TrimSpace(zeton.Name), rolePrzedrostekDesignu)
		if roleSystemuWizualnegoDesignu[nazwa] || widziane[nazwa] {
			continue
		}
		widziane[nazwa] = true
		nieznane = append(nieznane, zeton.Name)
	}
	return nieznane
}

// rozlozZapisZetonowDesignu rozkłada zapis zewnętrzny na żetony. Rdzeń czyta
// dwie postacie: zapis JSON (płaski albo zagnieżdżony) oraz zmienne CSS.
func rozlozZapisZetonowDesignu(tresc string, postac *string) ([]shared.DesignToken, error) {
	nazwaPostaci := ""
	if postac != nil {
		nazwaPostaci = strings.ToLower(strings.TrimSpace(*postac))
	}
	if nazwaPostaci == "" {
		if strings.HasPrefix(strings.TrimSpace(tresc), "{") {
			nazwaPostaci = "json"
		} else if strings.Contains(tresc, rolePrzedrostekDesignu) || strings.Contains(tresc, "--") {
			nazwaPostaci = "css"
		}
	}
	switch nazwaPostaci {
	case "json":
		return rozlozZetonyZJsonDesignu(tresc)
	case "css", "cssvariables":
		return rozlozZetonyZCssDesignu(tresc), nil
	}
	return nil, fmt.Errorf("postaci zapisu żetonów serwer nie rozpoznał; serwer czyta zapis json " +
		"(płaski albo zagnieżdżony) oraz zmienne css — wskaż postać polem format")
}

// rozlozZetonyZJsonDesignu czyta zapis JSON. Struktura zagnieżdżona składa
// nazwę roli z kolejnych poziomów rozdzielonych myślnikiem — tak, jak nazywa je
// system wizualny produktu, a nie kropką, której CSS w nazwie własności nie
// przyjmie.
func rozlozZetonyZJsonDesignu(tresc string) ([]shared.DesignToken, error) {
	var drzewo map[string]any
	if err := json.Unmarshal([]byte(tresc), &drzewo); err != nil {
		return nil, fmt.Errorf("zapis żetonów nie jest poprawnym JSON-em: %v", err)
	}
	zetony := []shared.DesignToken{}
	zejdzPoDrzewieZetonowDesignu("", drzewo, &zetony)
	return zetony, nil
}

// zejdzPoDrzewieZetonowDesignu schodzi po zapisie zagnieżdżonym i dokłada
// żetony. Węzeł z polem `value` (albo `$value`, jak stanowi W3C) jest żetonem;
// pozostałe są grupami.
func zejdzPoDrzewieZetonowDesignu(sciezka string, wezel map[string]any, zetony *[]shared.DesignToken) {
	if wartosc, rodzaj, jest := wartoscWezlaZetonuDesignu(wezel); jest && sciezka != "" {
		zeton := shared.DesignToken{
			Name:  sciezka,
			Kind:  rodzaj,
			Value: wartosc,
		}
		if opis, znany := wezel["description"].(string); znany {
			zeton.Description = &opis
		}
		*zetony = append(*zetony, zeton)
		return
	}
	for klucz, poddrzewo := range wezel {
		if strings.HasPrefix(klucz, "$") {
			continue
		}
		nazwa := klucz
		if sciezka != "" {
			nazwa = sciezka + "-" + klucz
		}
		switch wartosc := poddrzewo.(type) {
		case map[string]any:
			zejdzPoDrzewieZetonowDesignu(nazwa, wartosc, zetony)
		case string:
			// Zapis płaski: nazwa roli wprost na wartość.
			*zetony = append(*zetony, shared.DesignToken{
				Name:  nazwa,
				Kind:  rodzajZetonuZWartosciDesignu(wartosc),
				Value: wartosc,
			})
		}
	}
}

// wartoscWezlaZetonuDesignu rozstrzyga, czy węzeł jest żetonem, i oddaje jego
// wartość wraz z rodzajem. Rodzaj zapisany wprost (`type` albo `$type`) bije
// rozpoznanie z wartości — zapis mówi o sobie prawdę lepiej niż zgadywacz.
func wartoscWezlaZetonuDesignu(wezel map[string]any) (string, shared.DesignTokenKind, bool) {
	surowa, jest := wezel["value"]
	if !jest {
		surowa, jest = wezel["$value"]
	}
	if !jest {
		return "", "", false
	}
	wartosc := fmt.Sprintf("%v", surowa)
	rodzaj := rodzajZetonuZWartosciDesignu(wartosc)
	nazwaRodzaju, znany := wezel["type"].(string)
	if !znany {
		nazwaRodzaju, znany = wezel["$type"].(string)
	}
	if znany {
		if przelozony, rozpoznany := rodzajZetonuZNazwyDesignu(nazwaRodzaju); rozpoznany {
			rodzaj = przelozony
		}
	}
	return wartosc, rodzaj, true
}

// rozlozZetonyZCssDesignu czyta zmienne CSS. Rdzeń bierze wyłącznie własności
// niestandardowe (`--nazwa: wartość`) — reszta arkusza opisuje wygląd
// konkretnych elementów, a nie role systemu.
func rozlozZetonyZCssDesignu(tresc string) []shared.DesignToken {
	zetony := []shared.DesignToken{}
	for _, wiersz := range strings.Split(tresc, "\n") {
		wiersz = strings.TrimSpace(wiersz)
		if !strings.HasPrefix(wiersz, "--") {
			continue
		}
		numer := strings.Index(wiersz, ":")
		if numer < 0 {
			continue
		}
		nazwa := strings.TrimSpace(wiersz[:numer])
		wartosc := strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(wiersz[numer+1:]), ";"))
		if wartosc == "" {
			continue
		}
		zetony = append(zetony, shared.DesignToken{
			Name:  strings.TrimPrefix(nazwa, rolePrzedrostekDesignu),
			Kind:  rodzajZetonuZWartosciDesignu(wartosc),
			Value: wartosc,
		})
	}
	return zetony
}

// rodzajZetonuZNazwyDesignu przekłada nazwę rodzaju z zapisu zewnętrznego na
// rodzaj kontraktu. Nazwy spoza wykazu nie są tłumaczone na siłę — wołający
// dostaje wtedy rodzaj rozpoznany z wartości.
func rodzajZetonuZNazwyDesignu(nazwa string) (shared.DesignTokenKind, bool) {
	switch strings.ToLower(strings.TrimSpace(nazwa)) {
	case "color":
		return shared.DesignTokenKindColor, true
	case "dimension", "spacing", "size", "borderradius":
		return shared.DesignTokenKindDimension, true
	case "fontfamily", "font-family":
		return shared.DesignTokenKindFontFamily, true
	case "fontweight", "font-weight":
		return shared.DesignTokenKindFontWeight, true
	case "duration":
		return shared.DesignTokenKindDuration, true
	case "shadow", "boxshadow":
		return shared.DesignTokenKindShadow, true
	}
	return "", false
}

// rodzajZetonuZWartosciDesignu rozpoznaje rodzaj z samej wartości. Rozpoznanie
// jest zachowawcze i przy niepewności oddaje rodzaj miary.
func rodzajZetonuZWartosciDesignu(wartosc string) shared.DesignTokenKind {
	oczyszczona := strings.ToLower(strings.TrimSpace(wartosc))
	switch {
	case strings.HasPrefix(oczyszczona, "#"),
		strings.HasPrefix(oczyszczona, "rgb"),
		strings.HasPrefix(oczyszczona, "hsl"),
		strings.HasPrefix(oczyszczona, "oklch"),
		strings.HasPrefix(oczyszczona, "color-mix("):
		return shared.DesignTokenKindColor
	case strings.HasSuffix(oczyszczona, "ms"), strings.HasSuffix(oczyszczona, "s") &&
		!strings.HasSuffix(oczyszczona, "px") && !strings.HasSuffix(oczyszczona, "rems"):
		return shared.DesignTokenKindDuration
	case strings.Contains(oczyszczona, "inset") || strings.Count(oczyszczona, "px") > 2:
		return shared.DesignTokenKindShadow
	case strings.Contains(oczyszczona, "serif"), strings.Contains(oczyszczona, "monospace"),
		strings.Contains(oczyszczona, "system-ui"):
		return shared.DesignTokenKindFontFamily
	}
	return shared.DesignTokenKindDimension
}

// zestawZetonowKontraktuDesignu składa `DesignTokenSet` kontraktu z wiersza
// zestawu i jego żetonów. `TokenCount` bierze się z licznika policzonego przez
// bazę, nie z długości wykazu — prawdą o zestawie jest to, co w nim leży.
func zestawZetonowKontraktuDesignu(z dane.ZestawZetonowDesignu,
	zetony []dane.ZetonDesignu) shared.DesignTokenSet {

	lista := make([]shared.DesignToken, 0, len(zetony))
	for _, zeton := range zetony {
		lista = append(lista, shared.DesignToken{
			Name:        zeton.Nazwa,
			Kind:        shared.DesignTokenKind(zeton.Rodzaj),
			Value:       zeton.Wartosc,
			Description: zeton.Opis,
			AliasOf:     zeton.OdsylaczDo,
		})
	}
	return shared.DesignTokenSet{
		Id:         z.Kod,
		WindowId:   z.Okno,
		Name:       z.Nazwa,
		Theme:      z.Motyw,
		Tokens:     lista,
		TokenCount: z.Liczba,
		UpdatedAt:  chwilaBazy(z.Zaktualizowano),
	}
}

// bladNieznanegoZestawuZetonowDesignu nazywa zestaw, którego rdzeń nie zna —
// zapis z kodem wskazującym zestaw nieistniejący jest odmową, nie założeniem
// nowego zestawu pod cudzym kodem.
func bladNieznanegoZestawuZetonowDesignu(kod string, err error) error {
	if czyBrakZasobuDesignu(err) {
		return bladNieznanegoBytuDesignu("zestawu żetonów " + kod + " nie ma w module Design")
	}
	return bladDesignu(err)
}
