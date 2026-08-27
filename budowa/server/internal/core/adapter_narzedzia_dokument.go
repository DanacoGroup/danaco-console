// Plik niesie podstawę rodziny narzędzi dokumentowych modelu: typ adaptera,
// jego montaż, zasięg uruchamiania binariów, magazyn bajtów wyniku,
// rozpoznanie formatu źródła i wspólne odmowy.
package core

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/konfiguracja"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/session"
	"danacoconsole/server/internal/zewnetrzne"
	"danacoconsole/shared"
)

const (
	granicaKonwersjiDokumentu   = 3 * time.Minute
	granicaOdczytuDokumentu     = 90 * time.Second
	granicaRozpoznaniaDokumentu = 2 * time.Minute
)

// Pięć binariów arsenału, którymi pracuje ta rodzina. Nazwy są dla CZŁOWIEKA
// i wchodzą wprost do treści odmowy „nie ma czym", razem z pakietem do
// dociągnięcia — Operator ma przeczytać, co zainstalować, a nie szukać sam.
var (
	narzedziePandoc = zewnetrzne.Narzedzie{
		Nazwa: "Pandoc", Program: "pandoc", Pakiet: "pandoc",
	}
	narzedziePdfDoTekstu = zewnetrzne.Narzedzie{
		Nazwa: "poppler (pdftotext)", Program: "pdftotext", Pakiet: "poppler-utils",
	}
	narzedziePdfDoObrazu = zewnetrzne.Narzedzie{
		Nazwa: "poppler (pdftoppm)", Program: "pdftoppm", Pakiet: "poppler-utils",
	}
	narzedzieTesseract = zewnetrzne.Narzedzie{
		Nazwa: "Tesseract OCR", Program: "tesseract", Pakiet: "tesseract-ocr tesseract-ocr-pol",
	}
	narzedzieLibreOffice = zewnetrzne.Narzedzie{
		Nazwa: "LibreOffice", Program: "libreoffice", Pakiet: "libreoffice",
	}
)

// adapterNarzedziDokumentu wypełnia port Dokumenty i składa zależności, przez
// które przechodzą obie komendy rodziny narzędzi dokumentowych.
type adapterNarzedziDokumentu struct {
	// uruchamiacz jest jedyną portową drogą startu procesu w drzewie.
	uruchamiacz session.Uruchamiacz
	// rozstrzygacz i katalog składają zasady izolacji i obszar zasięgu.
	rozstrzygacz *konfig.Rozstrzygacz
	katalog      *KatalogRoboczy
	// magazyn trzyma bajty wyników pod sumą sha256, wspólnie z zasobami Designu.
	magazyn *magazynTresciBiblioteki
	// zasobyDesignu czyta zasób wskazany `assetId` i zapisuje wiersz wyniku.
	zasobyDesignu dane.RepozytoriumDesignu
	// biblioteka jest drugim magazynem, w którym `assetId` może wskazywać
	// materiał.
	biblioteka dane.RepozytoriumBiblioteki
}

// nowyAdapterNarzedziDokumentu wiąże port z uruchamiaczem procesów i wpina
// magazyn wyników nad domyślnym katalogiem danych, zastępowanym montażem
// `ZKatalogiemDanych`.
func nowyAdapterNarzedziDokumentu(uruchamiacz session.Uruchamiacz) *adapterNarzedziDokumentu {
	return &adapterNarzedziDokumentu{
		uruchamiacz: uruchamiacz,
		magazyn:     magazynZasobowDesignu(konfiguracja.KatalogDanychDomyslny()),
	}
}

// ZKatalogiemDanych przestawia magazyn wyników na katalog wskazany
// konfiguracją procesu (przełącznik `-dane`, zmienna `DANACO_KATALOG_DANYCH`).
func (a *adapterNarzedziDokumentu) ZKatalogiemDanych(katalog string) *adapterNarzedziDokumentu {
	if strings.TrimSpace(katalog) != "" {
		a.magazyn = magazynZasobowDesignu(katalog)
	}
	return a
}

// ZIzolacja podpina rozstrzygacz zasięgu i ustalacz katalogu roboczego — dwa
// źródła, z których powstają zasady i obszar egzekwowane przy uruchomieniu.
func (a *adapterNarzedziDokumentu) ZIzolacja(rozstrzygacz *konfig.Rozstrzygacz,
	katalog *KatalogRoboczy) *adapterNarzedziDokumentu {

	a.rozstrzygacz, a.katalog = rozstrzygacz, katalog
	return a
}

// ZZasobamiDesignu podpina repozytorium zasobów, przez które przechodzi
// odczyt `assetId` oraz zapis wiersza wyniku; zależność jest opcjonalna.
func (a *adapterNarzedziDokumentu) ZZasobamiDesignu(r dane.RepozytoriumDesignu) *adapterNarzedziDokumentu {
	a.zasobyDesignu = r
	return a
}

// ZBiblioteka podpina repozytorium plików biblioteki, drugi magazyn, w którym
// `assetId` może wskazywać materiał; zależność jest opcjonalna.
func (a *adapterNarzedziDokumentu) ZBiblioteka(r dane.RepozytoriumBiblioteki) *adapterNarzedziDokumentu {
	a.biblioteka = r
	return a
}

// zasiegDokumentu składa trójkę okno-zasady-obszar dla zasięgu platformy,
// wpisując środowisko wykonania jawnie zamiast pozostawiać je puste.
func (a *adapterNarzedziDokumentu) zasiegDokumentu() (session.Okno, session.Zasady, session.Obszar) {
	okno := session.Okno{Ustawienia: session.Ustawienia{
		SrodowiskoWykonania: shared.ExecutionEnvCore,
	}}
	zasady := session.Zasady{}
	if a.rozstrzygacz != nil {
		zasady = ZasadyIzolacji(a.rozstrzygacz, konfig.Kontekst{})
	}
	obszar := session.Obszar{}
	if a.katalog != nil {
		obszar = ObszarOkna(a.katalog.Ustal(konfig.Kontekst{}, ""), "")
	}
	return okno, zasady, obszar
}

// wolaj uruchamia jedno binarium arsenału i oddaje jego wyjście; brak
// uruchamiacza kończy się nazwaną odmową zamiast paniką procesu.
func (a *adapterNarzedziDokumentu) wolaj(ctx context.Context, n zewnetrzne.Narzedzie,
	argumenty []string, granica time.Duration) ([]byte, error) {

	if a.uruchamiacz == nil {
		return nil, odmowaDokumentu(shared.ErrorCodeInternalError,
			"rdzeń nie ma uruchamiacza procesów, więc narzędzie "+n.Nazwa+
				" nie ma czym wystartować; naprawa: podpiąć warstwę kanału przy składaniu rdzenia")
	}
	okno, zasady, obszar := a.zasiegDokumentu()
	wynik, err := zewnetrzne.Wolaj(ctx, a.uruchamiacz, okno, zasady, obszar, n, argumenty, "", granica)
	if err != nil {
		return wynik.Wyjscie, bladNarzedziaDokumentu(err)
	}
	return wynik.Wyjscie, nil
}

// bladNarzedziaDokumentu przekłada odmowy pakietu `zewnetrzne` na kody
// kontraktu: brak narzędzia na `channel_unavailable`, naruszenie izolacji na
// `permission_denied`, a resztę na `internal_error`.
func bladNarzedziaDokumentu(err error) error {
	var brak *zewnetrzne.BrakNarzedzia
	if errors.As(err, &brak) {
		return odmowaDokumentu(shared.ErrorCodeChannelUnavailable, brak.Error())
	}
	if errors.Is(err, session.ErrIzolacja) {
		return odmowaDokumentu(shared.ErrorCodePermissionDenied, err.Error())
	}
	return odmowaDokumentu(shared.ErrorCodeInternalError, err.Error())
}

// bladZadaniaDokumentu znakuje wadę ŻĄDANIA kodem kontraktu. Kod nieponawialny
// i słusznie: to samo żądanie powtórzone da to samo.
func bladZadaniaDokumentu(powod string) error {
	return odmowaDokumentu(shared.ErrorCodeValidationFailed, powod)
}

// odmowaDokumentu składa odmowę rodziny z kodem kontraktu. Przedrostek nazywa
// AUTORA odmowy, żeby Operator czytający Errors Panel wiedział, kto odmówił,
// zanim przeczyta dlaczego.
func odmowaDokumentu(kod shared.ErrorCode, powod string) error {
	return protocol.JakoError(protocol.NowyBlad(kod, "narzędzia dokumentowe: "+powod))
}

// zrodloDokumentu jest rozstrzygniętym wejściem czynności: ścieżka
// BEZWZGLĘDNA do bajtów oraz format rozpoznany albo wskazany.
type zrodloDokumentu struct {
	sciezka string
	format  string
	// okno niesie identyfikator okna materiału źródłowego, gdy jest nim zasób.
	okno string
}

// ustalZrodlo rozstrzyga trzy drogi wejścia kontraktu w jedną ścieżkę na
// dysku, w pierwszeństwie `assetId`, potem `sourcePath`, na końcu `content`,
// i odmawia, gdy żądanie nie poda żadnego wskazania.
func (a *adapterNarzedziDokumentu) ustalZrodlo(ctx context.Context, katalogPracy string,
	zasob, sciezka, tresc, format *string, komenda string) (zrodloDokumentu, error) {

	if !bezWartosci(zasob) {
		return a.zrodloZZasobu(ctx, *zasob, format)
	}

	if !bezWartosci(sciezka) {
		pelna, err := filepath.Abs(strings.TrimSpace(*sciezka))
		if err != nil {
			return zrodloDokumentu{}, bladZadaniaDokumentu(
				"ścieżki źródłowej nie da się rozwinąć do bezwzględnej: " + err.Error())
		}
		if _, err := os.Stat(pelna); err != nil {
			return zrodloDokumentu{}, odmowaDokumentu(shared.ErrorCodeNotFound,
				"pod ścieżką "+pelna+" nie ma pliku do przetworzenia: "+err.Error())
		}
		return zrodloDokumentu{
			sciezka: pelna,
			// Kolejność rozpoznania: wskazanie wołającego, rozszerzenie pliku,
			// na końcu pomiar bajtów.
			format: pierwszyFormatDokumentu(wartoscTekstu(format), formatZeSciezki(pelna),
				formatZNaglowka(pelna)),
		}, nil
	}

	if !bezWartosci(tresc) {
		// Treść wprost trafia na dysk, bo binaria arsenału czytają pliki, nie
		// pamięć rdzenia.
		wskazany := normalizujFormatDokumentu(wartoscTekstu(format))
		if wskazany == "" {
			// Domyślny markdown jest wskazaniem kontraktu, nie zgadywaniem
			// formatu.
			wskazany = "markdown"
		}
		plik := filepath.Join(katalogPracy, "zrodlo."+rozszerzenieFormatuDokumentu(wskazany))
		if err := os.WriteFile(plik, []byte(*tresc), 0o600); err != nil {
			return zrodloDokumentu{}, odmowaDokumentu(shared.ErrorCodeInternalError,
				"nie można odłożyć treści źródłowej na dysk: "+err.Error())
		}
		return zrodloDokumentu{sciezka: plik, format: wskazany}, nil
	}

	return zrodloDokumentu{}, bladZadaniaDokumentu("komenda " + komenda +
		" bez materiału: brakuje pola assetId, sourcePath i content — " +
		"czynność nie ma czego przetworzyć")
}

// zrodloZZasobu rozwiązuje `assetId` na odwołanie do bajtów w magazynie,
// szukając najpierw wśród zasobów designu, a dopiero potem w bibliotece.
func (a *adapterNarzedziDokumentu) zrodloZZasobu(ctx context.Context, kod string,
	format *string) (zrodloDokumentu, error) {

	if a.zasobyDesignu == nil && a.biblioteka == nil {
		return zrodloDokumentu{}, odmowaDokumentu(shared.ErrorCodeInternalError,
			"wskazano assetId, a rdzeń nie ma wpiętego żadnego magazynu zasobów — "+
				"naprawa: podpiąć repozytorium zasobów designu albo biblioteki przy "+
				"składaniu rdzenia; do czasu naprawy materiał podaje się polem "+
				"sourcePath albo content")
	}

	kod = strings.TrimSpace(kod)
	if a.zasobyDesignu != nil {
		if wiersz, err := a.zasobyDesignu.Zasob(ctx, kod); err == nil {
			return a.zrodloZOdwolania(kod, wiersz.URI, wartoscTekstu(wiersz.Format),
				wiersz.Okno, format)
		}
	}
	if a.biblioteka != nil {
		if plik, err := a.biblioteka.Plik(ctx, kod); err == nil {
			// Plik biblioteki nie należy do okna, więc wynik zostaje zasobem
			// wolnym.
			return a.zrodloZOdwolania(kod, plik.TrescOdwolanie, wartoscTekstu(plik.MimeType),
				"", format)
		}
	}
	return zrodloDokumentu{}, odmowaDokumentu(shared.ErrorCodeNotFound,
		"zasobu "+kod+" nie ma ani wśród zasobów designu, ani wśród plików biblioteki")
}

// zrodloZOdwolania składa materiał z wiersza magazynu: rozwija odwołanie do
// ścieżki bezwzględnej i rozpoznaje format. Wspólne dla obu magazynów, bo
// różnią się tylko nazwami kolumn, którymi wołający już się posłużył.
func (a *adapterNarzedziDokumentu) zrodloZOdwolania(kod string, odwolanie *string,
	formatWiersza string, okno string, format *string) (zrodloDokumentu, error) {

	if odwolanie == nil || strings.TrimSpace(*odwolanie) == "" {
		return zrodloDokumentu{}, odmowaDokumentu(shared.ErrorCodeNotFound,
			"zasób "+kod+" nie ma odwołania do bajtów — nie ma czego przetworzyć")
	}
	pelna, err := filepath.Abs(*odwolanie)
	if err != nil {
		return zrodloDokumentu{}, odmowaDokumentu(shared.ErrorCodeInternalError,
			"odwołania zasobu "+kod+" nie da się rozwinąć do ścieżki bezwzględnej: "+err.Error())
	}
	rozpoznany := pierwszyFormatDokumentu(wartoscTekstu(format), formatWiersza,
		formatZNaglowka(pelna))
	if rozpoznany == "" {
		return zrodloDokumentu{}, bladZadaniaDokumentu(
			"zasób " + kod + " nie niesie formatu, a żądanie go nie wskazuje — " +
				"naprawa: podać pole fromFormat")
	}
	return zrodloDokumentu{sciezka: pelna, format: rozpoznany, okno: okno}, nil
}

// katalogPracyDokumentu zakłada katalog tymczasowy na materiał pośredni
// jednego wywołania i oddaje funkcję jego usunięcia niezależnie od powodzenia.
func katalogPracyDokumentu() (string, func(), error) {
	katalog, err := os.MkdirTemp("", "danaco-dokument-")
	if err != nil {
		return "", func() {}, odmowaDokumentu(shared.ErrorCodeInternalError,
			"nie można założyć katalogu roboczego czynności: "+err.Error())
	}
	return katalog, func() { _ = os.RemoveAll(katalog) }, nil
}
