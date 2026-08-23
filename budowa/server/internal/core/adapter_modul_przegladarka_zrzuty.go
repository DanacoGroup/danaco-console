// Odpowiedzialność pliku: materiał wizualny i wytwory sesji przeglądania —
// `browser.screenshot.capture`, `browser.snapshot.screenshot.get`,
// `browser.artifact.add`.
//
// Za każdym zrzutem leżą bajty. Komenda wykonuje zrzut silnikiem przeglądarki,
// odkłada obraz w magazynie modułu pod sumą sha256 i dopiero wtedy zapisuje
// wiersz z odwołaniem. Odpowiedź udana z pustym `ref` — albo z odwołaniem
// wskazującym nic — jest dokładnie tą szkodą, którą ten produkt już raz
// popełnił w module Design, i której pilnują sprawdziany skutku.
//
// Migawka dostaje odsyłacz do zrzutu. Po `browser.screenshot.capture` wiersz
// migawki okna niesie `zrzut_odwolanie`, więc `browser.snapshot.get`
// z `includeScreenshot` przestaje oddawać pustkę: to jest ta sama treść widziana
// dwiema drogami, a nie dwa niezależne byty.
package core

import (
	"context"
	"encoding/base64"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// WykonajZrzut obsługuje `browser.screenshot.capture`.
func (a *adapterPrzegladarki) WykonajZrzut(ctx context.Context,
	z shared.BrowserScreenshotCaptureRequest) (shared.BrowserScreenshotCaptureResponse, error) {

	migawka, err := a.adresOstatniejStrony(ctx, z.WindowId, "screenshot.capture")
	if err != nil {
		return shared.BrowserScreenshotCaptureResponse{}, err
	}
	format := shared.BrowserImageFormat(shared.BrowserImageFormatPng)
	if z.Format != nil {
		format = *z.Format
	}

	sesja, wynik, err := a.otworzStrone(ctx, "browser.screenshot.capture", nastawyStrony{Url: migawka.Url})
	if err != nil {
		return shared.BrowserScreenshotCaptureResponse{}, err
	}
	defer sesja.Zamknij()

	zapis, szerokosc, wysokosc, err := sesja.zrzut(ctx, string(z.Mode), string(format),
		wartoscTekstuLubPusta(z.Selector), wartoscLiczby(z.Quality),
		wartoscLiczby(z.X), wartoscLiczby(z.Y), wartoscLiczby(z.Width), wartoscLiczby(z.Height),
		[2]int{wynik.Szerokosc, wynik.Wysokosc})
	if err != nil {
		return shared.BrowserScreenshotCaptureResponse{}, bladSilnikaPrzegladarki("browser.screenshot.capture", err)
	}
	bajty, err := base64.StdEncoding.DecodeString(zapis)
	if err != nil {
		return shared.BrowserScreenshotCaptureResponse{}, bladPrzegladarki(err)
	}
	odwolanie, err := a.zapiszTresc(bajty)
	if err != nil {
		return shared.BrowserScreenshotCaptureResponse{}, err
	}

	rozmiar := int64(len(bajty))
	zrzut, err := a.repozytorium.ZapiszZrzut(ctx, dane.ZrzutPrzegladania{
		Kod:                 nowyIdentyfikator(przedrostekZrzutu),
		Okno:                z.WindowId,
		MigawkaZewnetrznaID: &migawka.Kod,
		TrescOdwolanie:      odwolanie,
		Tryb:                string(z.Mode),
		Format:              string(format),
		Szerokosc:           int64(szerokosc),
		Wysokosc:            int64(wysokosc),
		RozmiarBajtow:       &rozmiar,
	})
	if err != nil {
		return shared.BrowserScreenshotCaptureResponse{}, bladPrzegladarki(err)
	}

	// Migawka okna dostaje odsyłacz do zrzutu: bez tego `browser.snapshot.get`
	// z `includeScreenshot` oddawałby pustkę tuż po wykonaniu zrzutu tej samej
	// strony — i Operator miałby dwa sprzeczne zdania o jednym stanie okna.
	migawka.ZrzutOdwolanie = &odwolanie
	if _, err := a.repozytorium.ZapiszMigawke(ctx, dane.MigawkaStrony{
		Kod: nowyIdentyfikator(przedrostekMigawki), Okno: migawka.Okno, Url: migawka.Url,
		Tytul: migawka.Tytul, TekstOdwolanie: migawka.TekstOdwolanie,
		ZrodloOdwolanie: migawka.ZrodloOdwolanie, ZrzutOdwolanie: &odwolanie,
	}); err != nil {
		return shared.BrowserScreenshotCaptureResponse{}, bladPrzegladarki(err)
	}

	// Zrzut jest też wytworem sesji — Capture & Monitor Panel pokazuje jedną
	// listę materiału, a nie osobną listę na każdy rodzaj.
	tytul := "Zrzut — " + migawka.Url
	mime := "image/" + string(format)
	if _, err := a.repozytorium.ZapiszWytwor(ctx, dane.WytworPrzegladania{
		Kod: nowyIdentyfikator(przedrostekWytworu), Okno: z.WindowId,
		Rodzaj: string(shared.BrowserArtifactKindScreenshot), Tytul: &tytul,
		TrescOdwolanie: odwolanie, TypMime: &mime, RozmiarBajtow: &rozmiar,
		UrlZrodla: &migawka.Url, MigawkaZewnetrznaID: &migawka.Kod,
	}); err != nil {
		return shared.BrowserScreenshotCaptureResponse{}, bladPrzegladarki(err)
	}

	return shared.BrowserScreenshotCaptureResponse{Screenshot: zrzutKontraktu(zrzut, nil)}, nil
}

// OdczytajZrzut obsługuje `browser.snapshot.screenshot.get` — odczyt zrzutu
// zapisanego przy migawce, wskazanego migawką albo odwołaniem.
func (a *adapterPrzegladarki) OdczytajZrzut(ctx context.Context,
	z shared.BrowserSnapshotScreenshotGetRequest) (shared.BrowserSnapshotScreenshotGetResponse, error) {

	migawkaKod := wartoscTekstuLubPusta(z.SnapshotId)
	odwolanie := wartoscTekstuLubPusta(z.ScreenshotRef)
	if migawkaKod == "" && odwolanie == "" {
		return shared.BrowserSnapshotScreenshotGetResponse{}, bladWskazaniaPrzegladarki(
			"komenda snapshot.screenshot.get bez migawki i bez odwołania — nie ma czego szukać")
	}

	var wiersz dane.ZrzutPrzegladania
	var err error
	if odwolanie != "" {
		wiersz, err = a.repozytorium.ZrzutPoOdwolaniu(ctx, odwolanie)
		if err != nil {
			return shared.BrowserSnapshotScreenshotGetResponse{}, bladWierszaPrzegladania("zrzutu", odwolanie, err)
		}
	} else {
		wiersz, err = a.repozytorium.ZrzutMigawki(ctx, migawkaKod)
		if err != nil {
			return shared.BrowserSnapshotScreenshotGetResponse{}, bladWierszaPrzegladania("zrzutu migawki", migawkaKod, err)
		}
	}

	var tresc *string
	if z.IncludeContent != nil && *z.IncludeContent {
		bajty, err := a.odczytajTresc(wiersz.TrescOdwolanie)
		if err != nil {
			return shared.BrowserSnapshotScreenshotGetResponse{}, err
		}
		zapis := base64.StdEncoding.EncodeToString(bajty)
		tresc = &zapis
	}
	return shared.BrowserSnapshotScreenshotGetResponse{Screenshot: zrzutKontraktu(wiersz, tresc)}, nil
}

// DodajWytwor obsługuje `browser.artifact.add` — zapisanie wytworu sesji:
// zrzutu, archiwum, wyodrębnionych danych albo adnotacji.
//
// Treść przychodzi wprost (`contentBase64`) albo odwołaniem do materiału już
// leżącego w magazynie (`contentRef`). Żądanie bez jednego i drugiego jest
// odmową: wytwór bez bajtów byłby wpisem wskazującym nic.
func (a *adapterPrzegladarki) DodajWytwor(ctx context.Context,
	z shared.BrowserArtifactAddRequest) (shared.BrowserArtifactAddResponse, error) {

	if strings.TrimSpace(z.WindowId) == "" {
		return shared.BrowserArtifactAddResponse{}, bladWskazaniaPrzegladarki("komenda artifact.add bez okna")
	}
	odwolanie := wartoscTekstuLubPusta(z.ContentRef)
	zapis := wartoscTekstuLubPusta(z.ContentBase64)
	var rozmiar int64

	switch {
	case zapis != "":
		bajty, err := base64.StdEncoding.DecodeString(zapis)
		if err != nil {
			return shared.BrowserArtifactAddResponse{}, bladWskazaniaPrzegladarki(
				"treść wytworu nie jest zapisem base64: " + err.Error())
		}
		odwolanie, err = a.zapiszTresc(bajty)
		if err != nil {
			return shared.BrowserArtifactAddResponse{}, err
		}
		rozmiar = int64(len(bajty))
	case odwolanie != "":
		// Odwołanie wskazane w żądaniu jest sprawdzane, a nie przyjmowane na
		// słowo: wiersz wskazujący plik, którego nie ma, przeszedłby każdy
		// sprawdzian istnienia wiersza i nie miałby czego pokazać.
		bajty, err := a.odczytajTresc(odwolanie)
		if err != nil {
			return shared.BrowserArtifactAddResponse{}, err
		}
		rozmiar = int64(len(bajty))
	default:
		return shared.BrowserArtifactAddResponse{}, bladWskazaniaPrzegladarki(
			"komenda artifact.add bez treści i bez odwołania do niej")
	}

	wytwor, err := a.repozytorium.ZapiszWytwor(ctx, dane.WytworPrzegladania{
		Kod:                 nowyIdentyfikator(przedrostekWytworu),
		Okno:                z.WindowId,
		Rodzaj:              string(z.Kind),
		Tytul:               z.Title,
		TrescOdwolanie:      odwolanie,
		TypMime:             z.MimeType,
		RozmiarBajtow:       &rozmiar,
		UrlZrodla:           z.SourceUrl,
		MigawkaZewnetrznaID: z.SnapshotId,
	})
	if err != nil {
		return shared.BrowserArtifactAddResponse{}, bladPrzegladarki(err)
	}
	return shared.BrowserArtifactAddResponse{Artifact: wytworKontraktu(wytwor)}, nil
}

// zrzutKontraktu przekłada wiersz zrzutu na byt kontraktu.
func zrzutKontraktu(w dane.ZrzutPrzegladania, tresc *string) shared.BrowserScreenshot {
	zrzut := shared.BrowserScreenshot{
		Ref:           w.TrescOdwolanie,
		SnapshotId:    w.MigawkaZewnetrznaID,
		Mode:          shared.BrowserScreenshotMode(w.Tryb),
		Format:        shared.BrowserImageFormat(w.Format),
		Width:         int(w.Szerokosc),
		Height:        int(w.Wysokosc),
		SizeBytes:     w.RozmiarBajtow,
		ContentBase64: tresc,
		CapturedAt:    chwilaBazy(w.Utworzono),
	}
	return zrzut
}

// wytworKontraktu przekłada wiersz wytworu na byt kontraktu.
func wytworKontraktu(w dane.WytworPrzegladania) shared.BrowserArtifact {
	return shared.BrowserArtifact{
		Id:         w.Kod,
		WindowId:   w.Okno,
		Kind:       shared.BrowserArtifactKind(w.Rodzaj),
		Title:      w.Tytul,
		ContentRef: w.TrescOdwolanie,
		MimeType:   w.TypMime,
		SizeBytes:  w.RozmiarBajtow,
		SourceUrl:  w.UrlZrodla,
		SnapshotId: w.MigawkaZewnetrznaID,
		CreatedAt:  chwilaBazy(w.Utworzono),
	}
}
