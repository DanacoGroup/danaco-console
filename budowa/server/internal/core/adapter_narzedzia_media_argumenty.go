// Odpowiedzialność pliku: układanie wiersza wywołania `ffmpeg` dla pięciu
// czynności wyliczenia `MediaOperationKind` oraz dobór kontenera. Każda
// odmowa wychodzi stąd przed uruchomieniem programu, nazywając brak.
package core

import (
	"strconv"
	"strings"

	"danacoconsole/shared"
)

// argumentyPrzetworzeniaMediow układa wiersz wywołania `ffmpeg` dla wskazanej
// czynności, nazywając brak parametru przed uruchomieniem programu. `-y`
// stoi przy każdej czynności, bo wynik leży w świeżym katalogu tymczasowym.
func argumentyPrzetworzeniaMediow(z shared.MediaTranscodeRequest, zrodlo,
	wynik string, strumienie []strumienMediow) ([]string, error) {

	switch z.Operation {
	case shared.MediaOperationKindConvert:
		return []string{"-y", "-i", zrodlo, wynik}, nil

	case shared.MediaOperationKindTrim:
		return argumentyWycieciaMediow(z, zrodlo, wynik)

	case shared.MediaOperationKindExtractAudio:
		return argumentyDzwiekuMediow(z, zrodlo, wynik, strumienie), nil

	case shared.MediaOperationKindResize:
		return argumentyRozmiaruMediow(z, zrodlo, wynik)

	case shared.MediaOperationKindFrame:
		return argumentyKlatkiMediow(z, zrodlo, wynik), nil
	}

	return nil, odmowaNieznanejCzynnosciMediow(z.Operation)
}

// argumentyWycieciaMediow składa wycięcie fragmentu. Jedna granica wystarczy:
// sam `startMs` to „od tego miejsca do końca", sam `endMs` odwrotnie. `-ss`
// i `-to` stoją po wejściu, dla granicy dokładnej co do klatki.
func argumentyWycieciaMediow(z shared.MediaTranscodeRequest, zrodlo,
	wynik string) ([]string, error) {

	if z.StartMs == nil && z.EndMs == nil {
		return nil, odmowaNarzedziMediow(shared.ErrorCodeValidationFailed,
			"czynność trim bez granic fragmentu: brakuje pól startMs i endMs — "+
				"wycięcie nie ma czego wyciąć, a wycięciem całości nie jest")
	}
	if z.StartMs != nil && *z.StartMs < 0 {
		return nil, odmowaNarzedziMediow(shared.ErrorCodeValidationFailed,
			"pole startMs jest ujemne — materiał nie zaczyna się przed swoim początkiem")
	}
	if z.EndMs != nil && *z.EndMs <= 0 {
		return nil, odmowaNarzedziMediow(shared.ErrorCodeValidationFailed,
			"pole endMs nie jest dodatnie — fragment kończyłby się przed początkiem materiału")
	}
	if z.StartMs != nil && z.EndMs != nil && *z.EndMs <= *z.StartMs {
		return nil, odmowaNarzedziMediow(shared.ErrorCodeValidationFailed,
			"koniec fragmentu ("+strconv.Itoa(*z.EndMs)+" ms) nie leży po jego "+
				"początku ("+strconv.Itoa(*z.StartMs)+" ms)")
	}

	argumenty := []string{"-y", "-i", zrodlo}
	if z.StartMs != nil {
		argumenty = append(argumenty, "-ss", sekundyMediow(*z.StartMs))
	}
	if z.EndMs != nil {
		argumenty = append(argumenty, "-to", sekundyMediow(*z.EndMs))
	}
	return append(argumenty, "-c", "copy", wynik), nil
}

// argumentyDzwiekuMediow składa wyodrębnienie ścieżki dźwiękowej: `-vn`
// odrzuca obraz, a strumień idzie kopią, gdy kontener wynika z kodeka.
func argumentyDzwiekuMediow(z shared.MediaTranscodeRequest, zrodlo, wynik string,
	strumienie []strumienMediow) []string {

	argumenty := []string{"-y", "-i", zrodlo, "-vn"}
	if bezWartosci(z.Format) && kodekDzwiekuMediow(strumienie) != "" {
		argumenty = append(argumenty, "-c:a", "copy")
	}
	return append(argumenty, wynik)
}

// argumentyRozmiaruMediow składa zmianę rozdzielczości. Wymiar niepodany idzie
// jako `-2`, co dobiera proporcję z zachowaniem parzystości wymaganej przez
// kodeki obrazu. Dźwięk idzie bez kodowania (`-c:a copy`).
func argumentyRozmiaruMediow(z shared.MediaTranscodeRequest, zrodlo,
	wynik string) ([]string, error) {

	if z.Width == nil && z.Height == nil {
		return nil, odmowaNarzedziMediow(shared.ErrorCodeValidationFailed,
			"czynność resize bez wymiarów: brakuje pól width i height — "+
				"nie wiadomo, do jakiego rozmiaru sprowadzić obraz")
	}
	if (z.Width != nil && *z.Width <= 0) || (z.Height != nil && *z.Height <= 0) {
		return nil, odmowaNarzedziMediow(shared.ErrorCodeValidationFailed,
			"wymiar niedodatni — obraz o szerokości albo wysokości zero nie istnieje")
	}

	return []string{"-y", "-i", zrodlo,
		"-vf", "scale=" + wymiarSkaliMediow(z.Width) + ":" + wymiarSkaliMediow(z.Height),
		"-c:a", "copy", wynik}, nil
}

// argumentyKlatkiMediow składa zrzut pojedynczej klatki: `-ss` stoi przed
// wejściem dla szybkiego przeskoku, brak `startMs` znaczy pierwszą klatkę,
// `-frames:v 1` ogranicza wynik, `-an` odrzuca dźwięk.
func argumentyKlatkiMediow(z shared.MediaTranscodeRequest, zrodlo,
	wynik string) []string {

	poczatek := 0
	if z.StartMs != nil && *z.StartMs > 0 {
		poczatek = *z.StartMs
	}
	argumenty := []string{"-y", "-ss", sekundyMediow(poczatek), "-i", zrodlo,
		"-frames:v", "1", "-an"}
	if z.Width != nil || z.Height != nil {
		argumenty = append(argumenty, "-vf",
			"scale="+wymiarSkaliMediow(z.Width)+":"+wymiarSkaliMediow(z.Height))
	}
	return append(argumenty, wynik)
}

// odmowaNieznanejCzynnosciMediow nazywa rodzaj przetworzenia spoza wyliczenia.
// Wykaz w treści jest celowy: wołający ma przeczytać, czego mógł chcieć,
// zamiast szukać wyliczenia w kontrakcie.
func odmowaNieznanejCzynnosciMediow(operacja shared.MediaOperationKind) error {
	return odmowaNarzedziMediow(shared.ErrorCodeValidationFailed,
		"nieznany rodzaj przetworzenia \""+string(operacja)+"\" — "+
			"serwer zna convert, trim, extractAudio, resize i frame")
}

// sekundyMediow przekłada milisekundy kontraktu na zapis czasu, który rozumie
// `ffmpeg`. Trzy miejsca po przecinku, bo milisekunda jest jednostką kontraktu
// i zaokrąglenie do sekundy przesuwałoby granice fragmentu.
func sekundyMediow(milisekundy int) string {
	return strconv.FormatFloat(float64(milisekundy)/1000, 'f', 3, 64)
}

// wymiarSkaliMediow oddaje wymiar dla filtra skalującego. Uzasadnienie „-2"
// stoi przy `argumentyRozmiaruMediow`.
func wymiarSkaliMediow(wymiar *int) string {
	if wymiar == nil || *wymiar <= 0 {
		return "-2"
	}
	return strconv.Itoa(*wymiar)
}

// kontenerDzwiekuMediow dobiera nośnik do zmierzonego kodeka dźwięku, żeby
// ścieżkę dało się przepisać bez ponownego kodowania. Kodek spoza wykazu
// oddaje pustkę zamiast nośnika dobranego na chybił trafił.
func kontenerDzwiekuMediow(kodek string) string {
	if strings.HasPrefix(kodek, "pcm_") {
		return "wav"
	}
	switch kodek {
	case "aac", "alac":
		return "m4a"
	case "mp3":
		return "mp3"
	case "flac":
		return "flac"
	case "opus":
		return "opus"
	case "vorbis":
		return "ogg"
	case "ac3":
		return "ac3"
	case "eac3":
		return "eac3"
	}
	return ""
}

// kontenerZrodlaMediow wybiera rozszerzenie z nazwy formatu oddanej przez
// `ffprobe`, po własnej kolejności pierwszeństwa, nie po kolejności programu.
func kontenerZrodlaMediow(nazwa string) string {
	czlony := strings.Split(strings.ToLower(strings.TrimSpace(nazwa)), ",")
	nalezy := func(szukany string) bool {
		for _, czlon := range czlony {
			if strings.TrimSpace(czlon) == szukany {
				return true
			}
		}
		return false
	}
	// Kolejność jest treścią: kontenery powszechne przed niszowymi.

	// Wynik ma się otworzyć u Operatora bez dokładania odtwarzacza.
	for _, kandydat := range []string{"mp4", "webm", "mp3", "wav", "flac", "ogg",
		"avi", "mov", "mpegts"} {

		if nalezy(kandydat) {
			return kandydat
		}
	}
	if nalezy("matroska") {
		return "mkv"
	}
	if len(czlony) > 0 && strings.TrimSpace(czlony[0]) != "" {
		return strings.TrimSpace(czlony[0])
	}
	return ""
}
