// Odpowiedzialność pliku: czynność `document.convert` — zamiana dokumentu
// między formatami i odłożenie wyniku w magazynie rdzenia.
//
// Drogi są dwie, bo żadna pojedyncza nie wystarcza. Pandoc zamienia struktury
// tekstowe (markdown, html, docx, odt, rtf, epub, csv, tekst czysty), ale PDF-u
// nie zapisze bez silnika składu — a instalka rdzenia żadnego nie niesie, więc
// `pandoc -t pdf` skończyłoby się odmową w połowie czynności.
//
// Dlatego PDF powstaje LibreOffice'em w trybie bez okna, a materiał podaje mu
// Pandoc: markdown → html (Pandoc) → pdf (LibreOffice). Gdy źródło samo jest
// strawne dla LibreOffice'a (docx, odt, rtf, html, csv, txt), kroku Pandoca nie
// ma — jeden przebieg zamiast dwóch.
//
// Komenda nie czyta PDF-u: jego treść jest ciągiem instrukcji rysowania, z
// którego Pandoc nie złoży struktury dokumentu. Żądanie `pdf → cokolwiek`
// kończy się odmową wskazującą drogę — treść PDF-u wyciąga
// `document.text.extract`, a jej wynik da się przekonwertować dalej.
//
// Wyniku zastępczego komenda nie podstawia: każda droga bez bajtów kończy się
// błędem, a pusty plik na wyjściu binarium też jest odmową, bo dokument
// o zerowej długości wygląda w panelu jak dokument.
package core

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"

	"danacoconsole/shared"
)

// Przeksztalc obsługuje `document.convert`.
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

// przeprowadzKonwersje dobiera drogę i oddaje ścieżkę gotowego pliku wyniku.
//
// Format docelowy równy źródłowemu nie jest ani błędem, ani pracą: wraca plik
// źródłowy bez zmiany. Przepuszczenie go przez Pandoca przepisałoby dokument
// i po cichu zgubiło to, czego Pandoc nie odwzorowuje.
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

// doPdf prowadzi jedyną drogę do PDF-u, jaką rdzeń ma: LibreOffice bez okna.
// Materiał niestrawny dla niego wprost idzie najpierw przez Pandoca do
// samodzielnego HTML-a — samodzielnego (`--standalone`), bo dopiero wtedy plik
// niesie nagłówek kodowania i polskie znaki nie rozsypują się w składzie.
func (a *adapterNarzedziDokumentu) doPdf(ctx context.Context, katalogPracy string,
	zrodlo zrodloDokumentu) (string, error) {

	material := zrodlo.sciezka
	if !formatyDokumentu[zrodlo.format].strawnyDlaLibre {
		if !formatyDokumentu[zrodlo.format].czytaPandoc {
			return "", bladZadaniaDokumentu("rdzeń nie umie zamienić formatu " +
				zrodlo.format + " na PDF: ani LibreOffice nie otworzy go wprost, " +
				"ani Pandoc nie złoży z niego materiału pośredniego")
		}
		posredni, err := a.pandokiem(ctx, katalogPracy, zrodlo.sciezka,
			formatyDokumentu[zrodlo.format].pandoc, "html", "posredni.html")
		if err != nil {
			return "", err
		}
		material = posredni
	}

	// Własny profil użytkownika na jedno wywołanie: LibreOffice trzyma stan
	// w profilu, a drugie równoległe uruchomienie na wspólnym profilu kończy się
	// cichym zwarciem — jeden przebieg oddaje plik, drugi nic. Profil
	// w katalogu roboczym czynności znika razem z nim.
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

	// LibreOffice nie pyta o nazwę wyniku, tylko składa ją z nazwy materiału,
	// więc szukamy jej po fakcie. Nazwa zgadnięta z góry byłaby pudłem przy
	// pliku z kropkami, a odmowa mówiłaby o braku wyniku, którego nie ma tylko
	// dlatego, że patrzymy w złe miejsce.
	nazwa := strings.TrimSuffix(filepath.Base(material), filepath.Ext(material)) + ".pdf"
	plik := filepath.Join(wyjscie, nazwa)
	if _, err := os.Stat(plik); err != nil {
		return "", odmowaDokumentu(shared.ErrorCodeInternalError,
			"LibreOffice zakończył pracę, ale pliku PDF nie ma pod "+plik+
				" — dokument nie powstał: "+err.Error())
	}
	return plik, nil
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

// odlozWynik utrwala bajty dokumentu w magazynie rdzenia i zakłada wiersz
// zasobu — tą samą drogą i tym samym magazynem, co `design.asset.upload`.
//
// Kolejność jest zamierzona: najpierw bajty, potem wiersz. Wiersz wskazujący
// odwołanie, za którym nic nie leży, byłby dokumentem nie do otwarcia, a model
// zacytowałby go jako gotowy.
//
// Plik pusty jest odmową: binarium potrafi skończyć się kodem zero i zostawić
// zero bajtów (uszkodzone wejście, wyczerpany nośnik), a zero bajtów nie jest
// dokumentem.
//
// Okno puste znaczy „bez wiersza", a nie „błąd" — wspólnie dla wszystkich
// rodzin arsenału, jak opisano przy `oknoWynikuArsenalu`.
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
