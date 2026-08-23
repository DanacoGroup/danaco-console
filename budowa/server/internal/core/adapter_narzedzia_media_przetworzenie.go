// Przetworzenie materiału: komenda `media.transcode` — jej przebieg i dobór
// kontenera wyniku. Wiersz wywołania `ffmpeg` dla pięciu rodzajów czynności
// (`MediaOperationKind`: convert, trim, extractAudio, resize, frame) układa
// `adapter_narzedzia_media_argumenty.go`; trzon rodziny leży w
// `adapter_narzedzia_media.go`, pomiar w `adapter_narzedzia_media_pomiar.go`.
//
// Każda czynność wytwarza nowy zasób, źródło zostaje nietknięte: `ffmpeg` nigdy
// nie dostaje ścieżki źródła jako celu zapisu. Wejściem jest blob magazynu (plik
// pod sumą kontrolną, którego nadpisanie zerwałoby tożsamość wszystkich zasobów
// o tej treści), wyjściem plik katalogu tymczasowego.
//
// Pomiar idzie przed przetworzeniem, bo daje trzy rozstrzygnięcia, których
// inaczej trzeba by zgadywać: kontener domyślny (gdy wołający nie podał
// formatu), kodek dźwięku (przy wyodrębnianiu ścieżki) oraz to, czy materiał
// w ogóle niesie dźwięk. Zgadnięte, każde z nich kończy się plikiem pustym albo
// odmową samego binarium, czyli odmową bez nazwy braku.
//
// Parametru, którego czynność wymaga, nie domyślamy się: `trim` bez granic daje
// odmowę nazywającą brak, a nie kopię materiału pod nazwą fragmentu.
package core

import (
	"context"
	"os"
	"strings"

	"danacoconsole/shared"
)

// kontenerObrazuMediow jest domyślnym nośnikiem zrzutu klatki. Zrzut ma być
// wierny, a PNG jest bezstratny i czytany przez każdy panel rdzenia. Wołający,
// który chce innego nośnika, podaje `format`.
const kontenerObrazuMediow = "png"

// Przetworz obsługuje `media.transcode`.
//
// Porządek jest ścisły: źródło (bajty pod sumą) → pomiar (rozstrzygnięcia
// zamiast domysłów) → argumenty (odmowa przed uruchomieniem, gdy czegoś brak) →
// binarium → sprawdzenie, czy wynik ma bajty → pomiar wyniku → zasób.
// Sprawdzenie rozmiaru wyprzedza pomiar, bo `ffmpeg` potrafi zakończyć się
// powodzeniem, zostawiając plik pusty, gdy wskazany fragment leży poza
// materiałem, a zasobu bez bajtów magazyn nie przyjmuje.
func (a *adapterNarzedziMediow) Przetworz(ctx context.Context,
	z shared.MediaTranscodeRequest) (shared.MediaTranscodeResponse, error) {

	sciezka, oknoZrodla, err := a.zrodloMediow(ctx, "media.transcode", z.AssetId, z.SourcePath)
	if err != nil {
		return shared.MediaTranscodeResponse{}, err
	}

	opis, surowe, err := a.zmierzMediow(ctx, sciezka)
	if err != nil {
		return shared.MediaTranscodeResponse{}, err
	}
	strumienie := strumienieMediow(surowe)

	rozszerzenie, err := rozszerzenieWynikuMediow(z, opis, strumienie)
	if err != nil {
		return shared.MediaTranscodeResponse{}, err
	}

	katalog, sprzataj, err := katalogPrzejsciowyMediow()
	if err != nil {
		return shared.MediaTranscodeResponse{}, err
	}
	defer sprzataj()

	wynikowy := plikWynikowyMediow(katalog, rozszerzenie)
	argumenty, err := argumentyPrzetworzeniaMediow(z, sciezka, wynikowy, strumienie)
	if err != nil {
		return shared.MediaTranscodeResponse{}, err
	}

	if _, err := a.wolajMediow(ctx, narzedzieFfmpeg, argumenty,
		granicaPrzetworzeniaMediow); err != nil {

		return shared.MediaTranscodeResponse{}, bladNarzedziMediow(err)
	}

	if stan, err := os.Stat(wynikowy); err != nil || stan.Size() == 0 {
		return shared.MediaTranscodeResponse{}, odmowaNarzedziMediow(
			shared.ErrorCodeValidationFailed,
			"ffmpeg zakończył pracę, ale nie zostawił żadnych bajtów wyniku — "+
				"najczęściej znaczy to, że wskazany fragment leży poza materiałem; "+
				"naprawa: sprawdzić granice komendą media.inspect")
	}

	return a.zasobPrzetworzenia(ctx, z, oknoZrodla, wynikowy, rozszerzenie)
}

// zasobPrzetworzenia mierzy plik wynikowy i odkłada go jako nowy zasób.
//
// Dwie miary, bo dwa rodzaje wyniku. Zrzut klatki jest obrazem: jego wymiary
// czyta nagłówek pliku (`rozpoznajObrazZasobu` — ta sama droga, którą mierzy je
// wniesienie zasobu), a pole czasu trwania zostaje puste, bo obraz nie trwa.
// Każdy inny wynik jest materiałem czasowym: mierzy go `ffprobe`, a wymiary
// bierzemy ze strumienia obrazu, którego wyodrębniony dźwięk nie ma.
func (a *adapterNarzedziMediow) zasobPrzetworzenia(ctx context.Context,
	z shared.MediaTranscodeRequest, oknoZrodla, wynikowy, rozszerzenie string) (
	shared.MediaTranscodeResponse, error) {

	nazwa := "media-" + string(z.Operation) + "." + rozszerzenie
	// `windowId` żądania bije okno materiału źródłowego, a brak obu znaczy
	// wynik bez wiersza — jedna reguła dla czterech rodzin
	// (`oknoWynikuArsenalu`).
	oknoWyniku := oknoWynikuArsenalu(z.WindowId, oknoZrodla)

	if z.Operation == shared.MediaOperationKindFrame {
		format, szerokosc, wysokosc := rozpoznajObrazZasobu(wynikowy)
		zasob, err := a.odlozWynikMediow(ctx, wynikowy, oknoWyniku, nazwa,
			pierwszyFormatMediow(format, rozszerzenie), szerokosc, wysokosc)
		if err != nil {
			return shared.MediaTranscodeResponse{}, err
		}
		return shared.MediaTranscodeResponse{Asset: zasob}, nil
	}

	opis, surowe, err := a.zmierzMediow(ctx, wynikowy)
	if err != nil {
		return shared.MediaTranscodeResponse{}, err
	}
	szerokosc, wysokosc := wymiaryMediow(strumienieMediow(surowe))

	zasob, err := a.odlozWynikMediow(ctx, wynikowy, oknoWyniku, nazwa, rozszerzenie,
		szerokosc, wysokosc)
	if err != nil {
		return shared.MediaTranscodeResponse{}, err
	}
	trwanie := trwanieMediowMs(opis.Format.Trwanie)
	return shared.MediaTranscodeResponse{Asset: zasob, DurationMs: &trwanie}, nil
}

// rozszerzenieWynikuMediow rozstrzyga kontener wyniku.
//
// Wskazanie wołającego bije wszystko: `format` jest polem kontraktu i tylko temu
// wyborowi służy. Gdy go nie ma, rozstrzyga rodzaj czynności:
//
//   - `convert` bez formatu jest odmową, bo bez wskazania program przepisałby
//     materiał do tego samego kontenera i oddał kopię pod nazwą przekształcenia.
//   - `frame` spada na PNG.
//   - `extractAudio` spada na kontener wyprowadzony z kodeka zmierzonego
//     w materiale, więc ścieżkę da się przepisać bez ponownego kodowania.
//   - `trim` i `resize` zostają w kontenerze źródła, bo żadna z tych czynności
//     nie jest zmianą formatu.
func rozszerzenieWynikuMediow(z shared.MediaTranscodeRequest, opis opisMediow,
	strumienie []strumienMediow) (string, error) {

	if !bezWartosci(z.Format) {
		rozszerzenie := rozszerzenieMediow(*z.Format)
		if rozszerzenie == "" {
			return "", odmowaNarzedziMediow(shared.ErrorCodeValidationFailed,
				"pole format jest puste — format docelowy albo się wskazuje, "+
					"albo pomija, ale nie zostawia pustym")
		}
		return rozszerzenie, nil
	}

	switch z.Operation {
	case shared.MediaOperationKindFrame:
		return kontenerObrazuMediow, nil

	case shared.MediaOperationKindExtractAudio:
		kodek := kodekDzwiekuMediow(strumienie)
		if kodek == "" {
			return "", odmowaNarzedziMediow(shared.ErrorCodeValidationFailed,
				"materiał nie niesie ani jednego strumienia dźwięku — "+
					"nie ma czego wyodrębnić")
		}
		kontener := kontenerDzwiekuMediow(kodek)
		if kontener == "" {
			return "", odmowaNarzedziMediow(shared.ErrorCodeValidationFailed,
				"nie wiadomo, w jakim kontenerze zapisać dźwięk kodeka "+kodek+
					"; naprawa: wskazać format docelowy polem format")
		}
		return kontener, nil

	case shared.MediaOperationKindConvert:
		return "", odmowaNarzedziMediow(shared.ErrorCodeValidationFailed,
			"czynność convert bez pola format nie ma celu — brakuje wskazania "+
				"formatu docelowego")

	case shared.MediaOperationKindTrim, shared.MediaOperationKindResize:
		kontener := kontenerZrodlaMediow(opis.Format.Nazwa)
		if kontener == "" {
			return "", odmowaNarzedziMediow(shared.ErrorCodeValidationFailed,
				"nie da się rozpoznać kontenera materiału źródłowego; "+
					"naprawa: wskazać format docelowy polem format")
		}
		return kontener, nil
	}

	return "", odmowaNieznanejCzynnosciMediow(z.Operation)
}

// pierwszyFormatMediow oddaje format zmierzony z nagłówka pliku, a przy jego
// braku — kontener, o który poproszono. Zmierzone bije zamówione, bo tylko ono
// opisuje bajty leżące w magazynie.
func pierwszyFormatMediow(zmierzony *string, zamowiony string) string {
	if zmierzony != nil && strings.TrimSpace(*zmierzony) != "" {
		return *zmierzony
	}
	return zamowiony
}
