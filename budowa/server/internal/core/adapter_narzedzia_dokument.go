// Odpowiedzialność pliku: podstawa rodziny narzędzi dokumentowych modelu —
// typ adaptera, jego montaż, zasięg uruchamiania binariów, magazyn bajtów
// wyniku, rozpoznanie formatu źródła i wspólne odmowy. Same czynności leżą
// osobno wedle odpowiedzialności: zamiana formatu w
// `adapter_narzedzia_dokument_konwersja.go`, odczyt treści w
// `adapter_narzedzia_dokument_tekst.go`, wpięcie do rejestru w
// `handlers_narzedzia_dokument.go`.
//
// Obie komendy rodziny są narzędziami modelu, nie panelem Operatora:
//
//   - `document.convert` — model oddaje Operatorowi dokument, a nie tekst
//     w oknie rozmowy;
//   - `document.text.extract` — model czyta plik, który dostał od Operatora:
//     PDF, skan, zdjęcie kartki. Bez tej komendy model nie ma żadnej drogi
//     do treści pliku leżącego na dysku Operatora.
//
// Jedna droga do binarium: wszystkie uruchomienia idą przez `zewnetrzne.Wolaj`
// — port `session.Uruchamiacz`, brama izolacji okna, objęcie drzewa potomstwa,
// obowiązkowa granica czasu. Własnego `exec.Command` ten moduł nie ma;
// w całym drzewie stoi dokładnie jedno, w `injection/rozruch.go`.
//
// Katalog uruchomienia zostaje pusty, a ścieżki są bezwzględne. Gdyby moduł
// podał bramie własny katalog tymczasowy, punkt izolacji katalogu roboczego
// odrzuciłby uruchomienie jako wyjście poza katalog okna
// (`session/izolacja_polecenie.go`). Pusty katalog pozwala bramie wstawić
// katalog własny zasięgu, a bezwzględne ścieżki argumentów sprawiają, że
// wybór katalogu nie zmienia wyniku pracy.
//
// Zasięg jest zasięgiem platformy. Żądania obu komend nie niosą okna — tak samo
// jak w rodzinie `speech.*` (`adapter_modul_mowa.go`) — bo narzędzie
// dokumentowe jest zdolnością platformy. Zasady i obszar składają się więc dla
// pustego `konfig.Kontekst{}`, czyli adresu najszerszego poziomu. Podstawienie
// `session.Zasady{}` z ręki znaczyłoby „izolacja wyłączona" niezależnie od
// tego, co Operator ustawił.
//
// Wynik jest widoczny dla Operatora, gdy żądanie poda `windowId`: powstaje
// wtedy wiersz zasobu w wykazie okna. Brak `windowId` nie wstrzymuje czynności
// — bajty i tak idą do magazynu pod sumą kontrolną, tylko zasób nie pojawia się
// w wykazie okna (patrz `oknoWynikuArsenalu`). Rodzaj zasobu rozstrzyga format
// wyniku, nie nazwa komendy: PDF jest dokumentem, a konwersja do PNG obrazem
// (`rodzajZasobuArsenalu`).
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

// adapterNarzedziDokumentu wypełnia port Dokumenty.
type adapterNarzedziDokumentu struct {
	// uruchamiacz jest portem warstwy kanału — jedyną drogą startu procesu
	// w drzewie. Bez niego rodzina odmawia, zamiast startować
	// binaria własnym `exec.Command`.
	uruchamiacz session.Uruchamiacz
	// rozstrzygacz i katalog składają zasady izolacji i obszar zasięgu
	// platformy — te same dwa źródła, którymi jadą Terminal, Developer i mowa.
	rozstrzygacz *konfig.Rozstrzygacz
	katalog      *KatalogRoboczy
	// magazyn trzyma BAJTY wyników konwersji pod sumą sha256 — i jest TYM SAMYM
	// magazynem zasobów Designu, którym jedzie `design.asset.upload`. Wynik jest
	// typu `DesignAsset`, więc Assets Panel znajduje go tam, gdzie szuka; skład
	// zasobu poza bazą jest jeden.
	magazyn *magazynTresciBiblioteki
	// zasobyDesignu daje DWIE rzeczy: odczyt zasobu wskazanego `assetId`
	// (materiał źródłowy) i zapis wiersza zasobu wynikowego, gdy żądanie poda
	// `windowId`.
	zasobyDesignu dane.RepozytoriumDesignu
	// biblioteka jest DRUGIM magazynem, w którym `assetId` może wskazywać
	// materiał. Kontrakt obiecuje przy tym polu „zasób z magazynu rdzenia
	// (Design albo Library)", a plik biblioteki trzyma bajty tak samo jak zasób
	// designu — bezwzględną ścieżką do bloku pod sumą sha256, różnica jest
	// wyłącznie w korzeniu katalogu. Bez tej zależności identyfikator pliku
	// biblioteki odmawiał, choć opis kontraktu go dopuszczał.
	biblioteka dane.RepozytoriumBiblioteki
}

// nowyAdapterNarzedziDokumentu wiąże port z uruchamiaczem procesów i wpina
// magazyn wyników nad katalogiem danych DOMYŚLNYM. Katalog obowiązujący wchodzi
// montażem (`ZKatalogiemDanych`) — tak samo jak w module Design; wartość
// domyślna zostaje po to, żeby konstruktor nigdy nie oddał adaptera bez
// magazynu i konwersja nie odmawiała z powodu własnego montażu.
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

// ZZasobamiDesignu podpina repozytorium zasobów: ODCZYT `assetId` i ZAPIS
// wiersza wyniku. Zależność opcjonalna: bez niej obie komendy pracują na
// `sourcePath` i `content`, wskazanie `assetId` odmawia zdaniem nazywającym
// brak, a wynik nie dostaje wiersza — bajty i tak lądują w magazynie. To mniej
// niż pełnia, ale nie jest bramą postawioną przed resztą.
func (a *adapterNarzedziDokumentu) ZZasobamiDesignu(r dane.RepozytoriumDesignu) *adapterNarzedziDokumentu {
	a.zasobyDesignu = r
	return a
}

// ZBiblioteka podpina repozytorium plików biblioteki — drugi magazyn, w którym
// `assetId` może wskazywać materiał. Zależność opcjonalna na tych samych
// zasadach co zasoby designu: bez niej identyfikator pliku biblioteki odmawia
// zdaniem nazywającym brak, a pozostałe drogi materiału pracują dalej.
func (a *adapterNarzedziDokumentu) ZBiblioteka(r dane.RepozytoriumBiblioteki) *adapterNarzedziDokumentu {
	a.biblioteka = r
	return a
}

// zasiegDokumentu składa trójkę okno–zasady–obszar dla zasięgu platformy.
// Uzasadnienie stoi w nagłówku pliku; `SrodowiskoWykonania` wpisane JAWNIE,
// bo plik Operatora leży na hoście rdzenia, więc binarium musi ruszyć tam samo
// — puste pole dałoby ten sam rozruch, ale milczkiem.
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

// wolaj uruchamia jedno binarium arsenału i oddaje jego wyjście.
//
// KATALOG PUSTY — patrz nagłówek pliku. Brak uruchamiacza jest odmową
// nazwaną, a nie panika: rdzeń złożony bez warstwy kanału ma powiedzieć, czego
// mu brakuje, a nie wywrócić się na pierwszym żądaniu modelu.
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
// kontraktu. Rozstrzygnięcie jest to samo, co w silniku mowy:
//
//   - `*zewnetrzne.BrakNarzedzia` → `channel_unavailable`. Nie ma czym wykonać
//     czynności. To nie jest zły argument (żądanie było poprawne) ani usterka
//     rdzenia (produkt działa i mówi wprost, co dociągnąć), tylko ZAPLECZE
//     NIEDOSTĘPNE — kod PONAWIALNY, bo po instalacji to samo żądanie przejdzie.
//   - naruszenie izolacji → `permission_denied`. Punkt izolacji Operatora
//     zatrzymał uruchomienie; tak samo znakuje je Terminal.
//   - reszta → `internal_error`: binarium wystartowało i się wywróciło, a
//     powód niesie własna diagnostyka programu, dołożona przez pakiet.
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
	// okno jest oknem, w którym leży materiał ŹRÓDŁOWY, gdy jest nim zasób.
	// Ścieżka i treść wprost nie należą do żadnego okna i zostawiają tu pustkę.
	// Rozstrzyga to o oknie WYNIKU, gdy żądanie nie poda `windowId`
	// (`oknoWynikuArsenalu`).
	okno string
}

// ustalZrodlo rozstrzyga trzy drogi wejścia kontraktu w jedną ścieżkę na dysku.
//
// PIERWSZEŃSTWO: `assetId`, potem `sourcePath`, na końcu `content`. Kolejność
// nie jest dowolna — zasób z magazynu jest treścią ZAMROŻONĄ pod sumą
// kontrolną, ścieżka jest treścią ŻYWĄ (Operator może ją w międzyczasie
// nadpisać), a treść wprost jest tym, co model właśnie napisał. Gdy przyszło
// więcej niż jedno wskazanie, wygrywa to pewniejsze.
//
// ŻADNE WSKAZANIE NIE JEST ODMOWĄ, a nie pustym wynikiem: czynność bez
// materiału nie ma czego przetworzyć, a „udało się, ale nic nie ma" byłoby
// powodzeniem czynności, której nikt nie wykonał.
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
			// Kolejność rozpoznania: wskazanie wołającego, potem rozszerzenie
			// pliku, na końcu POMIAR z pierwszych bajtów. Ostatni krok ratuje
			// materiał bez rozszerzenia — a takim jest każdy blob magazynu
			// (nazwą bloba jest suma sha256).
			format: pierwszyFormatDokumentu(wartoscTekstu(format), formatZeSciezki(pelna),
				formatZNaglowka(pelna)),
		}, nil
	}

	if !bezWartosci(tresc) {
		// Treść wprost trafia na dysk, bo binaria arsenału czytają PLIKI, nie
		// pamięć rdzenia. Rozszerzenie bierze się z formatu, żeby narzędzie
		// miało po czym rozpoznać wejście, gdy nie podamy go jawnie.
		wskazany := normalizujFormatDokumentu(wartoscTekstu(format))
		if wskazany == "" {
			// Domyślny markdown jest tu WSKAZANIEM KONTRAKTU, nie zgadywaniem:
			// pole `fromFormat` opisano jako „brak znaczy rozpoznanie z pliku",
			// a treść wprost pliku nie ma. Markdown to jedyny format, w którym
			// tekst pisany przez model jest sam w sobie poprawnym dokumentem.
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

// zrodloZZasobu rozwiązuje `assetId` na odwołanie do bajtów w magazynie.
//
// MAGAZYNY SĄ DWA i kontrakt mówi to wprost przy tym polu: „zasób z magazynu
// rdzenia (Design albo Library)". Oba trzymają bajty tą samą konwencją —
// bezwzględną ścieżką do bloku pod sumą sha256 — i różnią się wyłącznie
// korzeniem katalogu oraz nazwą kolumny (`URI` zasobu, `TrescOdwolanie` pliku).
// Szukanie idzie najpierw po zasobach designu, bo tam trafiają wyniki własnych
// konwersji, a dopiero potem po bibliotece.
//
// BLOB NIE MA ROZSZERZENIA — jego nazwą jest suma sha256 — więc format bierze
// się z kolumny wiersza, a nie ze ścieżki. Zasób bez formatu i bez wskazania
// w żądaniu jest odmową: zgadnięty format wygląda w wyniku identycznie jak
// rozpoznany i nie da się ich odróżnić.
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
			// Plik biblioteki nie należy do okna — kolumny okna ta tabela nie ma,
			// więc wiersz wyniku nie dostanie przypisania i zostanie zasobem
			// wolnym. Rodzaj treści bywa typem MIME, a `normalizujFormatDokumentu`
			// zna typy MIME obok nazw formatów, więc przechodzi jednym polem.
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

// katalogPracyDokumentu zakłada katalog tymczasowy na materiał pośredni jednego
// wywołania. Sprzątanie zwraca funkcja — wywołanie kończy się jego usunięciem
// NIEZALEŻNIE od powodzenia, bo strony rozłożone na obrazy potrafią zająć setki
// megabajtów, a nikt po rdzeniu tego nie posprząta.
func katalogPracyDokumentu() (string, func(), error) {
	katalog, err := os.MkdirTemp("", "danaco-dokument-")
	if err != nil {
		return "", func() {}, odmowaDokumentu(shared.ErrorCodeInternalError,
			"nie można założyć katalogu roboczego czynności: "+err.Error())
	}
	return katalog, func() { _ = os.RemoveAll(katalog) }, nil
}
