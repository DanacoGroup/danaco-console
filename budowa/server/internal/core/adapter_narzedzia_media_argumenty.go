// Odpowiedzialność pliku: układanie wiersza wywołania `ffmpeg` dla pięciu
// czynności wyliczenia `MediaOperationKind` oraz dobór kontenera. Rozdział
// z `adapter_narzedzia_media_przetworzenie.go` idzie po odpowiedzialności:
// tamten plik prowadzi przebieg komendy (źródło, pomiar, binarium, zasób),
// a ten mówi wyłącznie językiem `ffmpeg`.
//
// Każda odmowa wychodzi stąd przed uruchomieniem programu. Nazwany brak
// parametru („brakuje pól startMs i endMs") jest dla wołającego czymś zupełnie
// innym niż diagnostyka binarium, które dostało wiersz bez sensu i odmówiło
// po swojemu.
package core

import (
	"strconv"
	"strings"

	"danacoconsole/shared"
)

// argumentyPrzetworzeniaMediow układa wiersz wywołania `ffmpeg` dla wskazanej
// czynności. Odmowa wychodzi stąd przed uruchomieniem programu — nazwany brak
// parametru jest dla wołającego czymś zupełnie innym niż diagnostyka binarium.
//
// `-y` stoi przy każdej czynności, bo plik wynikowy leży w świeżym katalogu
// tymczasowym i nadpisać może wyłącznie samego siebie; bez tego przełącznika
// `ffmpeg` czeka na odpowiedź człowieka, którego przy nim nie ma, i kończy się
// dopiero granicą czasu.
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

// argumentyWycieciaMediow składa wycięcie fragmentu.
//
// Brak obu granic jest odmową nazywającą brak — powód w nagłówku pliku. Jedna
// granica wystarczy i znaczy dokładnie tyle, ile mówi: sam `startMs` to
// „od tego miejsca do końca", sam `endMs` to „od początku do tego miejsca".
// To nie jest domyślanie się całości, tylko odczytanie tego, co wskazano.
//
// `-ss` i `-to` stoją po wejściu z zamysłem: przed wejściem `ffmpeg` przeskakuje
// do najbliższej klatki kluczowej i granica przesuwa się o ułamek sekundy, po
// wejściu jest dokładna. `-c copy` przepisuje strumienie bez ponownego
// kodowania — fragment ma być tym samym materiałem, tylko krótszym.
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

// argumentyDzwiekuMediow składa wyodrębnienie ścieżki dźwiękowej.
//
// `-vn` odrzuca obraz — to jest cała treść tej czynności. Wynik nie ma wymiarów
// i mieć ich nie będzie; kontrakt przewiduje to wprost polami opcjonalnymi.
//
// Strumień kopiujemy wtedy, gdy kontener wyprowadziliśmy z kodeka (brak
// `format`): dźwięk trafia do nośnika, który zna ten kodek, więc ponowne
// kodowanie pogorszyłoby materiał bez powodu. Przy formacie wskazanym zostawiamy
// wybór kodeka `ffmpeg`owi — kopia mogłaby do wskazanego kontenera nie pasować,
// a odmowa binarium byłaby wtedy karą za spełnienie prośby wołającego.
func argumentyDzwiekuMediow(z shared.MediaTranscodeRequest, zrodlo, wynik string,
	strumienie []strumienMediow) []string {

	argumenty := []string{"-y", "-i", zrodlo, "-vn"}
	if bezWartosci(z.Format) && kodekDzwiekuMediow(strumienie) != "" {
		argumenty = append(argumenty, "-c:a", "copy")
	}
	return append(argumenty, wynik)
}

// argumentyRozmiaruMediow składa zmianę rozdzielczości.
//
// Brak obu wymiarów jest odmową: „zmień rozmiar" bez podania rozmiaru nie
// niesie żadnego polecenia.
//
// Wymiar niepodany dajemy jako `-2`, a nie jako liczbę wyliczoną samodzielnie.
// `-2` znaczy dla filtra „dobierz z proporcji źródła, zaokrąglając do liczby
// parzystej" — proporcje zostają nietknięte, a parzystość jest wymogiem
// kodeków obrazu, które próbkują chrominancję co dwa piksele i wysokości
// nieparzystej wprost odmawiają.
//
// Dźwięk przepisujemy bez kodowania (`-c:a copy`): zmiana rozdzielczości
// dotyczy obrazu i nie ma prawa dotknąć ścieżki dźwiękowej.
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

// argumentyKlatkiMediow składa zrzut pojedynczej klatki.
//
// `-ss` stoi przed wejściem, odwrotnie niż przy wycięciu, i to jest wybór:
// przeskok do klatki kluczowej jest tu tani i szybki, a różnica ułamka sekundy
// nie ma znaczenia dla zrzutu poglądowego — przy wycięciu miałaby, bo przesuwa
// granice fragmentu.
//
// Brak `startMs` znaczy początek materiału. To nie jest domyślanie się
// parametru, którego zabrakło: zrzut klatki ma sens od pierwszej klatki, a
// „pierwsza" jest wskazaniem tak samo jednoznacznym jak każde inne.
//
// `-frames:v 1` ogranicza wynik do jednej klatki, `-an` odrzuca dźwięk, którego
// obraz nie uniesie.
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
			"rdzeń zna convert, trim, extractAudio, resize i frame")
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

// kontenerDzwiekuMediow dobiera nośnik do zmierzonego kodeka dźwięku, tak żeby
// ścieżkę dało się przepisać bez ponownego kodowania. Kodek spoza wykazu oddaje
// pustkę, a wołający robi z niej odmowę proszącą o wskazanie formatu — bo
// nośnik dobrany na chybił trafił kończy się odmową samego binarium.
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
// `ffprobe`.
//
// Nazwa bywa wykazem, nie pojedynczym słowem: jeden zestaw procedur czyta całą
// rodzinę kontenerów i `ffprobe` oddaje wtedy wszystkie naraz
// („mov,mp4,m4a,3gp,3g2,mj2"). Wybór idzie po naszej kolejności pierwszeństwa,
// a nie po kolejności w wykazie programu: dla materiału MP4 wykaz zaczyna się
// od „mov" i wynik nosiłby rozszerzenie, którego nikt nie zamawiał, choć
// „mp4" stoi w tym samym wykazie o jedną pozycję dalej.
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
	// Kolejność jest treścią: kontenery powszechne przed niszowymi, żeby wynik
	// otwierał się u Operatora bez dokładania odtwarzacza.
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
