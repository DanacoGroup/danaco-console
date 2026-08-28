// Odpowiedzialność pliku: wydanie zestawu żetonów w postaci przyjmowanej przez
// kod (`design.tokenset.export`) i wydanie przewodnika systemu projektowego do
// modułu docelowego (`design.styleguide.publish`).
package core

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// WydajZestawZetonow wydaje zestaw żetonów w postaci przyjmowanej przez kod —
// obsługuje `design.tokenset.export`.
func (a *adapterDesignu) WydajZestawZetonow(ctx context.Context,
	z shared.DesignTokensetExportRequest) (shared.DesignTokensetExportResponse, error) {

	if strings.TrimSpace(z.TokenSetId) == "" {
		return shared.DesignTokensetExportResponse{}, bladWskazaniaDesignu(
			"komenda design.tokenset.export bez wskazania zestawu")
	}
	if err := sprawdzPostacWydaniaZetonowDesignu(z.Target); err != nil {
		return shared.DesignTokensetExportResponse{}, err
	}

	zestaw, zetony, err := a.zestawZZetonamiDesignu(ctx, strings.TrimSpace(z.TokenSetId))
	if err != nil {
		return shared.DesignTokensetExportResponse{}, err
	}

	tresc := zlozWydanieZetonowDesignu(zestaw, zetony, z.Target)
	nazwa := oczyscNazwePlikuDesignu(zestaw.Nazwa) + "." + rozszerzenieWydaniaZetonowDesignu(z.Target)

	odpowiedz := shared.DesignTokensetExportResponse{Content: tresc, FileName: nazwa}
	if z.TargetModuleId != nil && strings.TrimSpace(*z.TargetModuleId) != "" {
		if _, err := a.odlozWydanieDesignu(ctx, zestaw.Okno, nazwa, []byte(tresc),
			shared.DesignAssetKindDocument, rozszerzenieWydaniaZetonowDesignu(z.Target)); err != nil {
			return shared.DesignTokensetExportResponse{}, err
		}
		dostarczone := true
		odpowiedz.Delivered = &dostarczone
	}
	return odpowiedz, nil
}

// WydajPrzewodnikStylu wydaje przewodnik systemu projektowego do modułu
// docelowego — obsługuje `design.styleguide.publish`.
func (a *adapterDesignu) WydajPrzewodnikStylu(ctx context.Context,
	z shared.DesignStyleguidePublishRequest) (shared.DesignStyleguidePublishResponse, error) {

	if strings.TrimSpace(z.TokenSetId) == "" {
		return shared.DesignStyleguidePublishResponse{}, bladWskazaniaDesignu(
			"komenda design.styleguide.publish bez wskazania zestawu żetonów")
	}
	if strings.TrimSpace(z.TargetModuleId) == "" {
		return shared.DesignStyleguidePublishResponse{}, bladWskazaniaDesignu(
			"komenda design.styleguide.publish bez modułu docelowego — przewodnik wydany donikąd " +
				"jest plikiem, o którym nikt się nie dowie")
	}

	zestaw, zetony, err := a.zestawZZetonamiDesignu(ctx, strings.TrimSpace(z.TokenSetId))
	if err != nil {
		return shared.DesignStyleguidePublishResponse{}, err
	}

	var kolekcja *dane.KolekcjaDesignu
	if z.CollectionId != nil && strings.TrimSpace(*z.CollectionId) != "" {
		odczytana, err := a.repozytorium.KolekcjaDesignuPoKodzie(ctx, strings.TrimSpace(*z.CollectionId))
		if err != nil {
			return shared.DesignStyleguidePublishResponse{}, bladNieznanejKolekcjiDesignu(
				*z.CollectionId, err)
		}
		kolekcja = &odczytana
	}

	dokument := zlozPrzewodnikStyluDesignu(zestaw, zetony, strings.TrimSpace(z.TargetModuleId))
	nazwa := oczyscNazwePlikuDesignu("przewodnik-"+zestaw.Nazwa) + ".html"
	zasob, err := a.odlozWydanieDesignu(ctx, zestaw.Okno, nazwa, []byte(dokument),
		shared.DesignAssetKindDocument, "html")
	if err != nil {
		return shared.DesignStyleguidePublishResponse{}, err
	}

	if kolekcja != nil {
		if _, err := a.repozytorium.ZmienPrzypisaniaKolekcjiDesignu(ctx, kolekcja.ID,
			[]string{zasob.Kod}, false); err != nil {
			return shared.DesignStyleguidePublishResponse{}, bladDesignu(err)
		}
	}
	return shared.DesignStyleguidePublishResponse{AssetId: zasob.Kod, Published: true}, nil
}

// zestawZZetonamiDesignu czyta zestaw wraz z jego żetonami i odmawia zestawu
// pustego, bo taki wygląda jak wydanie udane bez ani jednej roli.
func (a *adapterDesignu) zestawZZetonamiDesignu(ctx context.Context,
	kod string) (dane.ZestawZetonowDesignu, []dane.ZetonDesignu, error) {

	zestaw, err := a.repozytorium.ZestawZetonowDesignuPoKodzie(ctx, kod)
	if err != nil {
		return dane.ZestawZetonowDesignu{}, nil, bladNieznanegoZestawuZetonowDesignu(kod, err)
	}
	zetony, err := a.repozytorium.ZetonyZestawuDesignu(ctx, zestaw.ID)
	if err != nil {
		return dane.ZestawZetonowDesignu{}, nil, bladDesignu(err)
	}
	if len(zetony) == 0 {
		return dane.ZestawZetonowDesignu{}, nil, bladWskazaniaDesignu(
			"zestaw żetonów " + kod + " nie ma ani jednego żetonu — nie ma czego wydać")
	}
	return zestaw, zetony, nil
}

// odlozWydanieDesignu utrwala bajty wydania w magazynie rdzenia i zakłada
// wiersz zasobu, przez który moduł docelowy po nie sięgnie.
func (a *adapterDesignu) odlozWydanieDesignu(ctx context.Context, okno, nazwa string,
	bajty []byte, rodzaj shared.DesignAssetKind, format string) (dane.ZasobDesignu, error) {

	if a.magazyn == nil {
		return dane.ZasobDesignu{}, bladZapisuZasobuDesignu(
			"magazyn treści nie jest wpięty — nie ma gdzie odłożyć bajtów wydania")
	}
	suma := sha256.Sum256(bajty)
	odwolanie, err := a.magazyn.Zapisz(bajty, hex.EncodeToString(suma[:]))
	if err != nil {
		return dane.ZasobDesignu{}, bladZapisuZasobuDesignu(
			"nie można utrwalić bajtów wydania: " + err.Error())
	}
	nazwaZasobu := nazwa
	postac := format
	zapisany, err := a.repozytorium.ZapiszZasob(ctx, dane.ZasobDesignu{
		Kod:    nowyIdentyfikator(przedrostekZasobuDesign),
		Okno:   okno,
		Nazwa:  &nazwaZasobu,
		Rodzaj: string(rodzaj),
		Format: &postac,
		URI:    &odwolanie,
	})
	if err != nil {
		return dane.ZasobDesignu{}, bladDesignu(err)
	}
	return zapisany, nil
}

// sprawdzPostacWydaniaZetonowDesignu odrzuca postać spoza kontraktu. Wykaz
// pochodzi z kontraktu i nie jest tu przepisywany.
func sprawdzPostacWydaniaZetonowDesignu(postac shared.DesignTokenTarget) error {
	dopuszczalne := shared.WartosciDesignTokenTarget()
	for _, znana := range dopuszczalne {
		if postac == znana {
			return nil
		}
	}
	nazwy := make([]string, 0, len(dopuszczalne))
	for _, znana := range dopuszczalne {
		nazwy = append(nazwy, string(znana))
	}
	return bladWskazaniaDesignu(fmt.Sprintf(
		"komenda design.tokenset.export z postacią wydania %q, której kontrakt nie zna; "+
			"postacie dopuszczalne: %s", string(postac), strings.Join(nazwy, ", ")))
}

// rozszerzenieWydaniaZetonowDesignu oddaje rozszerzenie pliku dla postaci
// wydania, zgodne z formatem, jaki narzędzie odbierające go rozpozna.
func rozszerzenieWydaniaZetonowDesignu(postac shared.DesignTokenTarget) string {
	switch postac {
	case shared.DesignTokenTargetCssVariables:
		return "css"
	case shared.DesignTokenTargetScss:
		return "scss"
	case shared.DesignTokenTargetTailwind, shared.DesignTokenTargetJavascript:
		return "js"
	case shared.DesignTokenTargetIos:
		return "swift"
	case shared.DesignTokenTargetAndroid:
		return "kt"
	}
	return "txt"
}

// zlozWydanieZetonowDesignu składa treść wydania w żądanej postaci, rozdzielając
// pracę do funkcji właściwej dla każdej z sześciu postaci kontraktu.
func zlozWydanieZetonowDesignu(zestaw dane.ZestawZetonowDesignu,
	zetony []dane.ZetonDesignu, postac shared.DesignTokenTarget) string {

	switch postac {
	case shared.DesignTokenTargetCssVariables:
		return zlozZetonyCssDesignu(zestaw, zetony)
	case shared.DesignTokenTargetScss:
		return zlozZetonyScssDesignu(zestaw, zetony)
	case shared.DesignTokenTargetTailwind:
		return zlozZetonyTailwindDesignu(zestaw, zetony)
	case shared.DesignTokenTargetJavascript:
		return zlozZetonyJavascriptDesignu(zestaw, zetony)
	case shared.DesignTokenTargetIos:
		return zlozZetonySwiftDesignu(zestaw, zetony)
	case shared.DesignTokenTargetAndroid:
		return zlozZetonyKotlinDesignu(zestaw, zetony)
	}
	return ""
}

// naglowekWydaniaZetonowDesignu składa komentarz otwierający każde wydanie.
// Wydanie ma powiedzieć o sobie, z czego powstało — plik żetonów bez tego
// jest w cudzym repozytorium plikiem nieznanego pochodzenia.
func naglowekWydaniaZetonowDesignu(zestaw dane.ZestawZetonowDesignu, znacznik string) string {
	motyw := ""
	if zestaw.Motyw != nil && strings.TrimSpace(*zestaw.Motyw) != "" {
		motyw = ", motyw " + *zestaw.Motyw
	}
	return fmt.Sprintf("%s Zestaw żetonów %q modułu Design Danaco Console (%s%s).\n"+
		"%s Wydanie z %s. Plik powstaje z zestawu — poprawki nanoś w module, nie tutaj.\n\n",
		znacznik, zestaw.Nazwa, zestaw.Kod, motyw, znacznik, zestaw.Zaktualizowano)
}

// zlozZetonyCssDesignu składa zmienne CSS pod korzeniem dokumentu — tam, gdzie
// kładzie je motyw produktu.
func zlozZetonyCssDesignu(zestaw dane.ZestawZetonowDesignu, zetony []dane.ZetonDesignu) string {
	var b strings.Builder
	b.WriteString("/*\n")
	b.WriteString(naglowekWydaniaZetonowDesignu(zestaw, " *"))
	b.WriteString(" */\n:root {\n")
	for _, zeton := range zetony {
		if zeton.Opis != nil && strings.TrimSpace(*zeton.Opis) != "" {
			fmt.Fprintf(&b, "  /* %s */\n", strings.TrimSpace(*zeton.Opis))
		}
		fmt.Fprintf(&b, "  %s%s: %s;\n", rolePrzedrostekDesignu, zeton.Nazwa,
			wartoscZetonuWydaniaDesignu(zeton))
	}
	b.WriteString("}\n")
	return b.String()
}

// zlozZetonyScssDesignu składa zmienne SCSS, jedną na żeton, poprzedzone
// nagłówkiem wydania w postaci komentarza SCSS.
func zlozZetonyScssDesignu(zestaw dane.ZestawZetonowDesignu, zetony []dane.ZetonDesignu) string {
	var b strings.Builder
	b.WriteString(naglowekWydaniaZetonowDesignu(zestaw, "//"))
	for _, zeton := range zetony {
		fmt.Fprintf(&b, "$%s: %s;\n", zeton.Nazwa, wartoscZetonuWydaniaDesignu(zeton))
	}
	return b.String()
}

// zlozZetonyTailwindDesignu składa konfigurację Tailwind. Żetony wchodzą
// rozdzielone po rodzaju, bo Tailwind ma osobne gałęzie dla barw, odstępów,
// krojów i cieni.
func zlozZetonyTailwindDesignu(zestaw dane.ZestawZetonowDesignu, zetony []dane.ZetonDesignu) string {
	galezie := map[string][]dane.ZetonDesignu{}
	for _, zeton := range zetony {
		galezie[galazTailwindDesignu(zeton.Rodzaj)] = append(galezie[galazTailwindDesignu(zeton.Rodzaj)], zeton)
	}
	nazwy := make([]string, 0, len(galezie))
	for nazwa := range galezie {
		nazwy = append(nazwy, nazwa)
	}
	sort.Strings(nazwy)

	var b strings.Builder
	b.WriteString(naglowekWydaniaZetonowDesignu(zestaw, "//"))
	b.WriteString("module.exports = {\n  theme: {\n    extend: {\n")
	for _, nazwa := range nazwy {
		fmt.Fprintf(&b, "      %s: {\n", nazwa)
		for _, zeton := range galezie[nazwa] {
			fmt.Fprintf(&b, "        %q: %q,\n", zeton.Nazwa, wartoscZetonuWydaniaDesignu(zeton))
		}
		b.WriteString("      },\n")
	}
	b.WriteString("    },\n  },\n};\n")
	return b.String()
}

// galazTailwindDesignu przekłada rodzaj żetonu na gałąź konfiguracji Tailwind,
// zgodną z podziałem, jakiego oczekuje narzędzie po drugiej stronie.
func galazTailwindDesignu(rodzaj string) string {
	switch shared.DesignTokenKind(rodzaj) {
	case shared.DesignTokenKindColor:
		return "colors"
	case shared.DesignTokenKindFontFamily:
		return "fontFamily"
	case shared.DesignTokenKindFontWeight:
		return "fontWeight"
	case shared.DesignTokenKindDuration:
		return "transitionDuration"
	case shared.DesignTokenKindShadow:
		return "boxShadow"
	}
	return "spacing"
}

// zlozZetonyJavascriptDesignu składa moduł JavaScript. Nazwy ról zostają
// dosłowne, w kluczach napisowych.
func zlozZetonyJavascriptDesignu(zestaw dane.ZestawZetonowDesignu, zetony []dane.ZetonDesignu) string {
	var b strings.Builder
	b.WriteString(naglowekWydaniaZetonowDesignu(zestaw, "//"))
	b.WriteString("export const zetony = {\n")
	for _, zeton := range zetony {
		fmt.Fprintf(&b, "  %q: %q,\n", zeton.Nazwa, wartoscZetonuWydaniaDesignu(zeton))
	}
	b.WriteString("};\n\nexport default zetony;\n")
	return b.String()
}

// zlozZetonySwiftDesignu składa zasoby dla systemu iOS. Wartości zostają
// napisami, także dla barw zapisanych jako `#rrggbb`.
func zlozZetonySwiftDesignu(zestaw dane.ZestawZetonowDesignu, zetony []dane.ZetonDesignu) string {
	var b strings.Builder
	b.WriteString(naglowekWydaniaZetonowDesignu(zestaw, "//"))
	b.WriteString("import Foundation\n\npublic enum Zetony {\n")
	for _, zeton := range zetony {
		fmt.Fprintf(&b, "    public static let %s = %q\n",
			identyfikatorKoduDesignu(zeton.Nazwa), wartoscZetonuWydaniaDesignu(zeton))
	}
	b.WriteString("}\n")
	return b.String()
}

// zlozZetonyKotlinDesignu składa zasoby dla systemu Android. Powód napisów jest
// ten sam, co przy wydaniu Swift.
func zlozZetonyKotlinDesignu(zestaw dane.ZestawZetonowDesignu, zetony []dane.ZetonDesignu) string {
	var b strings.Builder
	b.WriteString(naglowekWydaniaZetonowDesignu(zestaw, "//"))
	b.WriteString("object Zetony {\n")
	for _, zeton := range zetony {
		fmt.Fprintf(&b, "    const val %s: String = %q\n",
			identyfikatorKoduDesignu(zeton.Nazwa), wartoscZetonuWydaniaDesignu(zeton))
	}
	b.WriteString("}\n")
	return b.String()
}

// identyfikatorKoduDesignu przekłada nazwę roli na identyfikator przyjmowany
// przez języki programowania.
func identyfikatorKoduDesignu(nazwa string) string {
	czesci := strings.FieldsFunc(nazwa, func(znak rune) bool {
		return znak == '-' || znak == '.' || znak == '_' || znak == ' '
	})
	var b strings.Builder
	for numer, czesc := range czesci {
		if czesc == "" {
			continue
		}
		if numer == 0 {
			b.WriteString(strings.ToLower(czesc))
			continue
		}
		b.WriteString(strings.ToUpper(czesc[:1]))
		b.WriteString(strings.ToLower(czesc[1:]))
	}
	wynik := b.String()
	if wynik == "" {
		return "zeton"
	}
	if wynik[0] >= '0' && wynik[0] <= '9' {
		return "zeton" + strings.ToUpper(wynik[:1]) + wynik[1:]
	}
	return wynik
}

// wartoscZetonuWydaniaDesignu oddaje wartość żetonu do wydania. Żeton będący
// odsyłaczem wydaje się jako odwołanie do roli, na którą wskazuje.
func wartoscZetonuWydaniaDesignu(zeton dane.ZetonDesignu) string {
	if zeton.OdsylaczDo != nil && strings.TrimSpace(*zeton.OdsylaczDo) != "" {
		return fmt.Sprintf("var(%s%s)", rolePrzedrostekDesignu, strings.TrimSpace(*zeton.OdsylaczDo))
	}
	return zeton.Wartosc
}

// zlozPrzewodnikStyluDesignu składa przewodnik systemu projektowego jako
// dokument samowystarczalny: barwa pokazana próbką, miara — paskiem o tej
// szerokości, krój — zdaniem w tym kroju.
func zlozPrzewodnikStyluDesignu(zestaw dane.ZestawZetonowDesignu,
	zetony []dane.ZetonDesignu, modul string) string {

	var b strings.Builder
	b.WriteString("<!doctype html>\n<html lang=\"pl\">\n<head>\n<meta charset=\"utf-8\">\n")
	fmt.Fprintf(&b, "<title>Przewodnik systemu projektowego — %s</title>\n",
		zabezpieczHtmlDesignu(zestaw.Nazwa))
	b.WriteString("<style>\n")
	b.WriteString("body{font-family:system-ui,sans-serif;margin:2rem;background:#fff;color:#111}\n")
	b.WriteString("table{border-collapse:collapse;width:100%}\n")
	b.WriteString("th,td{border:1px solid #ddd;padding:.5rem .75rem;text-align:left;vertical-align:middle}\n")
	b.WriteString("th{background:#f4f4f4}\n")
	b.WriteString(".probka{display:inline-block;width:2.5rem;height:1.5rem;border:1px solid #999}\n")
	b.WriteString(".pasek{display:inline-block;height:.75rem;background:#333}\n")
	b.WriteString("</style>\n</head>\n<body>\n")
	fmt.Fprintf(&b, "<h1>%s</h1>\n", zabezpieczHtmlDesignu(zestaw.Nazwa))
	fmt.Fprintf(&b, "<p>Zestaw %s modułu Design, wydany do modułu %s. Żetonów: %d. "+
		"Stan z %s.</p>\n",
		zabezpieczHtmlDesignu(zestaw.Kod), zabezpieczHtmlDesignu(modul), len(zetony),
		zabezpieczHtmlDesignu(zestaw.Zaktualizowano))
	b.WriteString("<table>\n<thead><tr><th>Rola</th><th>Rodzaj</th><th>Wartość</th>" +
		"<th>Próbka</th><th>Do czego służy</th></tr></thead>\n<tbody>\n")
	for _, zeton := range zetony {
		opis := ""
		if zeton.Opis != nil {
			opis = *zeton.Opis
		}
		fmt.Fprintf(&b, "<tr><td><code>%s</code></td><td>%s</td><td><code>%s</code></td>"+
			"<td>%s</td><td>%s</td></tr>\n",
			zabezpieczHtmlDesignu(zeton.Nazwa), zabezpieczHtmlDesignu(zeton.Rodzaj),
			zabezpieczHtmlDesignu(zeton.Wartosc), probkaZetonuDesignu(zeton),
			zabezpieczHtmlDesignu(opis))
	}
	b.WriteString("</tbody>\n</table>\n</body>\n</html>\n")
	return b.String()
}

// probkaZetonuDesignu składa komórkę pokazującą wartość żetonu naocznie.
// Wartość wchodzi tu w atrybut stylu, więc przechodzi przez to samo
// zabezpieczenie, co treść — inaczej wartość z apostrofem rozerwałaby atrybut.
func probkaZetonuDesignu(zeton dane.ZetonDesignu) string {
	wartosc := zabezpieczHtmlDesignu(zeton.Wartosc)
	switch shared.DesignTokenKind(zeton.Rodzaj) {
	case shared.DesignTokenKindColor:
		return fmt.Sprintf(`<span class="probka" style="background:%s"></span>`, wartosc)
	case shared.DesignTokenKindDimension:
		return fmt.Sprintf(`<span class="pasek" style="width:%s"></span>`, wartosc)
	case shared.DesignTokenKindFontFamily:
		return fmt.Sprintf(`<span style="font-family:%s">Danaco Console</span>`, wartosc)
	case shared.DesignTokenKindFontWeight:
		return fmt.Sprintf(`<span style="font-weight:%s">Danaco Console</span>`, wartosc)
	case shared.DesignTokenKindShadow:
		return fmt.Sprintf(`<span class="probka" style="box-shadow:%s"></span>`, wartosc)
	}
	return "&nbsp;"
}

// zabezpieczHtmlDesignu zdejmuje z tekstu znaczenie składniowe HTML. Wartości
// żetonów pochodzą z zapisu wczytanego z zewnątrz, więc przewodnik składany
// bez tego byłby stroną, w którą da się wstrzyknąć cudzą treść.
func zabezpieczHtmlDesignu(tekst string) string {
	zamiennik := strings.NewReplacer(
		"&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&quot;", "'", "&#39;")
	return zamiennik.Replace(tekst)
}
