// Plik obsługuje media.transcode: przebieg przetworzenia materiału i dobór kontenera wyniku, gdy wołający formatu nie wskazał.
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

// Przetworz obsługuje media.transcode w ścisłym porządku: źródło, pomiar, argumenty, wywołanie narzędzia, sprawdzenie wyniku i pomiar wyniku.
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

// zasobPrzetworzenia mierzy plik wynikowy dwiema drogami zależnie od rodzaju wyniku i odkłada go jako nowy zasób magazynu.
func (a *adapterNarzedziMediow) zasobPrzetworzenia(ctx context.Context,
	z shared.MediaTranscodeRequest, oknoZrodla, wynikowy, rozszerzenie string) (
	shared.MediaTranscodeResponse, error) {

	nazwa := "media-" + string(z.Operation) + "." + rozszerzenie
	// Pole windowId żądania bije okno materiału źródłowego; brak obu znaczy wynik bez wiersza.
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

// rozszerzenieWynikuMediow rozstrzyga kontener wyniku: wskazanie wołającego bije wszystko, a przy jego braku decyduje rodzaj czynności.
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
