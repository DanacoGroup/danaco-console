// Adapter obsługuje `document.convert`: zamienia dokument między formatami i odkłada wynik
// w magazynie rdzenia, dobierając LibreOffice albo skład Pandoc+typst zależnie od formatu
// źródłowego i docelowego.
package core

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"

	"danacoconsole/server/internal/zewnetrzne"
	"danacoconsole/shared"
)

// narzedzieTypst opisuje silnik składu typst: jeden plik wykonywalny bez własnego drzewa
// zasobów, mieszczący się w płaskim układzie pomocniczych programów pakowania produktu.
var narzedzieTypst = zewnetrzne.Narzedzie{
	Nazwa: "typst", Program: "typst", Pakiet: "typst (jeden plik wykonywalny z wydania projektu)",
}

// czcionkaSkladuTypst jest rodziną czcionek narzucaną składowi typst przez zmienną szablonu
// `mainfont`; bez niej wykaz czcionek jest pusty i typst odmawia złożenia dokumentu.
const czcionkaSkladuTypst = "DejaVu Serif"

// Przeksztalc obsługuje komendę `document.convert`: rozpoznaje format źródłowy i docelowy,
// przeprowadza konwersję i odkłada gotowy dokument w magazynie rdzenia.
func (a *adapterNarzedziDokumentu) Przeksztalc(ctx context.Context,
	z shared.DocumentConvertRequest) (shared.DocumentConvertResponse, error) {

	docelowy := normalizujFormatDokumentu(z.ToFormat)
	if docelowy == "" {
		return shared.DocumentConvertResponse{}, bladZadaniaDokumentu(
			"komenda document.convert z formatem docelowym " +
				opisWskazaniaDokumentu(z.ToFormat) + ", którego rdzeń nie zna; " +
				"formaty znane: " + wykazFormatowDokumentu())
	}

	katalogPracy, posprzataj, err := katalogPracyDokumentu()
	if err != nil {
		return shared.DocumentConvertResponse{}, err
	}
	defer posprzataj()

	zrodlo, err := a.ustalZrodlo(ctx, katalogPracy, z.AssetId, z.SourcePath, z.Content,
		z.FromFormat, "document.convert")
	if err != nil {
		return shared.DocumentConvertResponse{}, err
	}
	if zrodlo.format == "" {
		return shared.DocumentConvertResponse{}, bladZadaniaDokumentu(
			"formatu źródła nie da się rozpoznać ani z pliku, ani z żądania — " +
				"naprawa: podać pole fromFormat; formaty znane: " + wykazFormatowDokumentu())
	}
	if obrazyDokumentu[zrodlo.format] {
		return shared.DocumentConvertResponse{}, bladZadaniaDokumentu(
			"źródło jest obrazem (" + zrodlo.format + "), a nie dokumentem — " +
				"zamiana formatów obrazu nie należy do tej komendy; " +
				"treść widoczną na obrazie odczytuje document.text.extract")
	}
	if zrodlo.format == "pdf" {
		return shared.DocumentConvertResponse{}, bladZadaniaDokumentu(
			"rdzeń nie czyta PDF-u jako dokumentu źródłowego: PDF niesie instrukcje " +
				"rysowania, a nie strukturę, z której da się złożyć inny format; " +
				"naprawa: wyciągnąć treść komendą document.text.extract, a jej wynik " +
				"przekonwertować dalej")
	}

	wyniki, err := a.przeprowadzKonwersje(ctx, katalogPracy, zrodlo, docelowy)
	if err != nil {
		return shared.DocumentConvertResponse{}, err
	}

	asset, rozmiar, err := a.odlozWynik(ctx, wyniki, oknoWynikuArsenalu(z.WindowId, zrodlo.okno), docelowy)
	if err != nil {
		return shared.DocumentConvertResponse{}, err
	}
	return shared.DocumentConvertResponse{Asset: asset, SizeBytes: rozmiar}, nil
}

// przeprowadzKonwersje dobiera drogę zamiany formatu i oddaje ścieżkę gotowego pliku wyniku;
// format docelowy równy źródłowemu zwraca plik źródłowy bez zmiany, z pominięciem Pandoca.
func (a *adapterNarzedziDokumentu) przeprowadzKonwersje(ctx context.Context, katalogPracy string,
	zrodlo zrodloDokumentu, docelowy string) (string, error) {

	if zrodlo.format == docelowy {
		return zrodlo.sciezka, nil
	}
	if docelowy == "pdf" {
		return a.doPdf(ctx, katalogPracy, zrodlo)
	}

	opisZrodla, opisCelu := formatyDokumentu[zrodlo.format], formatyDokumentu[docelowy]
	if !opisZrodla.czytaPandoc {
		return "", bladZadaniaDokumentu("rdzeń nie umie odczytać formatu " +
			zrodlo.format + " jako dokumentu źródłowego")
	}
	if !opisCelu.piszePandoc {
		return "", bladZadaniaDokumentu("rdzeń nie umie zapisać dokumentu w formacie " +
			docelowy + " na tej maszynie")
	}
	return a.pandokiem(ctx, katalogPracy, zrodlo.sciezka, opisZrodla.pandoc, opisCelu.pandoc,
		"wynik."+opisCelu.rozszerzenie)
}

// doPdf rozstrzyga drogę powstania PDF-u: materiał, który LibreOffice otwiera wprost, idzie
// LibreOffice'em; materiał niestrawny dla niego idzie przez Pandoca do składu typst albo do
// pośredniego HTML-a.
func (a *adapterNarzedziDokumentu) doPdf(ctx context.Context, katalogPracy string,
	zrodlo zrodloDokumentu) (string, error) {

	material := zrodlo.sciezka
	if !formatyDokumentu[zrodlo.format].strawnyDlaLibre {
		if !formatyDokumentu[zrodlo.format].czytaPandoc {
			return "", bladZadaniaDokumentu("rdzeń nie umie zamienić formatu " +
				zrodlo.format + " na PDF: ani LibreOffice nie otworzy go wprost, " +
				"ani Pandoc nie złoży z niego materiału pośredniego")
		}
		if zewnetrzne.Stoi(narzedzieTypst) {
			return a.skladTypstem(ctx, katalogPracy, zrodlo)
		}
		posredni, err := a.pandokiem(ctx, katalogPracy, zrodlo.sciezka,
			formatyDokumentu[zrodlo.format].pandoc, "html", "posredni.html")
		if err != nil {
			return "", err
		}
		material = posredni
	}

	// Profil własny na wywołanie unika zwarcia dwóch równoległych uruchomień na wspólnym profilu.
	profil := "-env:UserInstallation=file://" + filepath.Join(katalogPracy, "profil-libre")
	wyjscie := filepath.Join(katalogPracy, "pdf")
	if err := os.MkdirAll(wyjscie, 0o700); err != nil {
		return "", odmowaDokumentu(shared.ErrorCodeInternalError,
			"nie można założyć katalogu wyniku: "+err.Error())
	}
	if _, err := a.wolaj(ctx, narzedzieLibreOffice, []string{
		profil, "--headless", "--norestore", "--convert-to", "pdf",
		"--outdir", wyjscie, material,
	}, granicaKonwersjiDokumentu); err != nil {
		return "", err
	}

	// LibreOffice nazywa wynik po materiale, więc szuka się pliku po fakcie zamiast zgadywać nazwę.
	nazwa := strings.TrimSuffix(filepath.Base(material), filepath.Ext(material)) + ".pdf"
	plik := filepath.Join(wyjscie, nazwa)
	if _, err := os.Stat(plik); err != nil {
		return "", odmowaDokumentu(shared.ErrorCodeInternalError,
			"LibreOffice zakończył pracę, ale pliku PDF nie ma pod "+plik+
				" — dokument nie powstał: "+err.Error())
	}
	return plik, nil
}

// skladTypstem prowadzi drogę składu: Pandoc zamienia dokument na źródło typsta, a typst
// składa z niego PDF osobnym wywołaniem, aby odmowa silnika składu docierała bez zawinięcia
// w komunikat Pandoca.
func (a *adapterNarzedziDokumentu) skladTypstem(ctx context.Context, katalogPracy string,
	zrodlo zrodloDokumentu) (string, error) {

	zrodloSkladu := filepath.Join(katalogPracy, "sklad.typ")
	if _, err := a.wolaj(ctx, narzedziePandoc, []string{
		"--standalone", "--variable", "mainfont=" + czcionkaSkladuTypst,
		"--from", formatyDokumentu[zrodlo.format].pandoc, "--to", "typst",
		"--output", zrodloSkladu, zrodlo.sciezka,
	}, granicaKonwersjiDokumentu); err != nil {
		return "", err
	}
	if _, err := os.Stat(zrodloSkladu); err != nil {
		return "", odmowaDokumentu(shared.ErrorCodeInternalError,
			"Pandoc zakończył pracę, ale źródła składu nie ma pod "+zrodloSkladu+
				" — nie ma czego złożyć: "+err.Error())
	}

	wynik := filepath.Join(katalogPracy, "sklad.pdf")
	if _, err := a.wolaj(ctx, narzedzieTypst, []string{
		"compile", zrodloSkladu, wynik,
	}, granicaKonwersjiDokumentu); err != nil {
		return "", err
	}
	if _, err := os.Stat(wynik); err != nil {
		return "", odmowaDokumentu(shared.ErrorCodeInternalError,
			"typst zakończył pracę, ale pliku PDF nie ma pod "+wynik+
				" — dokument nie powstał: "+err.Error())
	}
	return wynik, nil
}

// pandokiem przeprowadza jedno wywołanie Pandoca i oddaje ścieżkę wyniku.
//
// `--standalone` idzie zawsze, bo formaty niosące własną otoczkę (html, epub,
// rtf) bez niego dostają sam fragment treści — plik otwiera się, ale nie jest
// dokumentem.
func (a *adapterNarzedziDokumentu) pandokiem(ctx context.Context, katalogPracy, material,
	zNazwy, doNazwy, nazwaWyniku string) (string, error) {

	wynik := filepath.Join(katalogPracy, nazwaWyniku)
	if _, err := a.wolaj(ctx, narzedziePandoc, []string{
		"--standalone", "--from", zNazwy, "--to", doNazwy, "--output", wynik, material,
	}, granicaKonwersjiDokumentu); err != nil {
		return "", err
	}
	if _, err := os.Stat(wynik); err != nil {
		return "", odmowaDokumentu(shared.ErrorCodeInternalError,
			"Pandoc zakończył pracę, ale pliku wyniku nie ma pod "+wynik+
				" — dokument nie powstał: "+err.Error())
	}
	return wynik, nil
}

// odlozWynik utrwala bajty dokumentu w magazynie rdzenia i zakłada wiersz zasobu tą samą
// drogą, co przesłanie zasobu; dokument o zerowej długości kończy się odmową zamiast wpisu
// w panelu.
func (a *adapterNarzedziDokumentu) odlozWynik(ctx context.Context,
	sciezka, oknoWyniku, format string) (shared.DesignAsset, int, error) {

	if a.magazyn == nil {
		return shared.DesignAsset{}, 0, odmowaDokumentu(shared.ErrorCodeInternalError,
			"magazyn treści nie jest wpięty — nie ma gdzie odłożyć bajtów dokumentu")
	}
	bajty, err := os.ReadFile(sciezka)
	if err != nil {
		return shared.DesignAsset{}, 0, odmowaDokumentu(shared.ErrorCodeInternalError,
			"nie można odczytać wytworzonego dokumentu: "+err.Error())
	}
	if len(bajty) == 0 {
		return shared.DesignAsset{}, 0, odmowaDokumentu(shared.ErrorCodeInternalError,
			"wytworzony dokument ma zero bajtów — czynność nie dała wyniku, "+
				"mimo że narzędzie zakończyło się bez błędu")
	}
	suma := sha256.Sum256(bajty)
	odwolanie, err := a.magazyn.Zapisz(bajty, hex.EncodeToString(suma[:]))
	if err != nil {
		return shared.DesignAsset{}, 0, odmowaDokumentu(shared.ErrorCodeInternalError,
			"nie można utrwalić bajtów dokumentu: "+err.Error())
	}

	zasob, err := odlozWynikArsenalu(ctx, a.zasobyDesignu, wynikArsenalu{
		odwolanie: odwolanie,
		okno:      oknoWyniku,
		nazwa:     "dokument." + rozszerzenieFormatuDokumentu(format),
		format:    format,
	}, func(powod string) error {
		return odmowaDokumentu(shared.ErrorCodeInternalError, powod)
	})
	if err != nil {
		return shared.DesignAsset{}, 0, err
	}
	return zasob, len(bajty), nil
}

// opisWskazaniaDokumentu ubiera wskazanie wołającego w cudzysłów albo nazywa
// jego brak. Odmowa „format ”" nie mówi nic; „bez formatu docelowego" mówi.
func opisWskazaniaDokumentu(wskazanie string) string {
	if strings.TrimSpace(wskazanie) == "" {
		return "pustym"
	}
	return "„" + wskazanie + "”"
}
