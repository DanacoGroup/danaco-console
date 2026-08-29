// Wypełnia `speech.audio.upload` i `speech.audio.fetch`: przyjmuje bajty
// nagrania mowy, utrwala je w magazynie rdzenia pod katalogiem danych
// i oddaje z powrotem wyłącznie nagrania, które rdzeń sam wystawił.
package core

import (
	"context"
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/mowa"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

const (
	// katalogNagranMowy to podkatalog katalogu danych rdzenia, w którym lądują
	// bajty przyjęte od okna. Obok niego leży `synteza-mowy` — oba są
	// nagraniami i oba odsłuchuje ta sama komenda.
	katalogNagranMowy = "nagrania-mowy"

	przedrostekNagraniaMowy = "nagranie-"

	// granicaNagraniaBajtow domyka jedno przyjęcie. Dwadzieścia pięć megabajtów
	// to około pół godziny mowy w Opusie — więcej niż każde polecenie głosowe
	// i każde dyktowanie, a mniej niż wielkość, przy której zapis w pamięci
	// wywraca proces.
	granicaNagraniaBajtow = 25 * 1024 * 1024

	// zycieNagraniaTymczasowego mówi, jak długo żyje nagranie, które nie ma
	// przeżyć transkrypcji. Godzina, a nie minuta: transkrypcja bywa
	// w kolejce, a nagranie skasowane przed nią zamieniłoby zlecenie
	// w odmowę „pliku nie ma".
	zycieNagraniaTymczasowego = time.Hour

	// kluczZapisuNagran to nastawa mówiąca, czy nagrania poleceń mają być
	// zapisywane trwale (Activity Feed). Kontrakt `speech.audio.upload` odsyła
	// do niej wprost przy pustym polu `retain`.
	kluczZapisuNagran = "mowa_zapis_nagran"
)

// PrzyjmijNagranie obsługuje `speech.audio.upload`: przyjmuje bajty nagrania
// zakodowane base64, utrwala je w magazynie i zakłada wiersz rejestru
// nagrań mowy.
func (a *adapterMowy) PrzyjmijNagranie(ctx context.Context,
	z shared.SpeechAudioUploadRequest) (shared.SpeechAudioUploadResponse, error) {

	if a.nagrania == nil {
		return shared.SpeechAudioUploadResponse{}, bladZapleczaNagran(
			"serwer nie ma wpiętego rejestru nagrań")
	}
	if strings.TrimSpace(z.Audio) == "" {
		return shared.SpeechAudioUploadResponse{}, bladWskazaniaNagrania(
			"żądanie bez bajtów nagrania — nie ma czego przyjąć")
	}
	bajty, err := base64.StdEncoding.DecodeString(strings.TrimSpace(z.Audio))
	if err != nil {
		return shared.SpeechAudioUploadResponse{}, bladWskazaniaNagrania(
			"bajtów nagrania nie da się odczytać jako base64: " + err.Error() +
				"; pole audio ma nieść sam zapis, bez przedrostka schematu danych")
	}
	if len(bajty) == 0 {
		return shared.SpeechAudioUploadResponse{}, bladWskazaniaNagrania(
			"nagranie ma zerowy rozmiar, więc nie niesie dźwięku do przepisania")
	}
	if len(bajty) > granicaNagraniaBajtow {
		return shared.SpeechAudioUploadResponse{}, bladWskazaniaNagrania(
			"nagranie przekracza granicę jednego przyjęcia (25 MB); " +
				"naprawa: podzielić wypowiedź na krótsze albo nagrywać w Opusie")
	}

	rozszerzenie, err := rozszerzenieNagrania(z.ContentType)
	if err != nil {
		return shared.SpeechAudioUploadResponse{}, err
	}

	katalog, err := a.katalogNagran()
	if err != nil {
		return shared.SpeechAudioUploadResponse{}, err
	}
	kod := nowyIdentyfikator(przedrostekNagraniaMowy)
	sciezka := filepath.Join(katalog, kod+rozszerzenie)
	if err := os.WriteFile(sciezka, bajty, 0o600); err != nil {
		return shared.SpeechAudioUploadResponse{}, bladZapleczaNagran(
			"nie można utrwalić bajtów nagrania: " + err.Error())
	}

	teraz := time.Now()
	trwale := a.czyZapisywacNagrania(z.Retain)
	wiersz := dane.NagranieMowy{
		Kod:           kod,
		Sciezka:       sciezka,
		TypTresci:     strings.TrimSpace(z.ContentType),
		RozmiarBajtow: int64(len(bajty)),
		SesjaKod:      z.SessionId,
		OknoKod:       z.WindowId,
		Trwale:        trwale,
		Utworzono:     teraz.UnixMilli(),
	}
	if !trwale {
		wygasa := teraz.Add(zycieNagraniaTymczasowego).UnixMilli()
		wiersz.Wygasa = &wygasa
	}
	zapisane, err := a.nagrania.ZapiszNagranieMowy(ctx, wiersz)
	if err != nil {
		_ = os.Remove(sciezka)
		return shared.SpeechAudioUploadResponse{}, bladZapleczaNagran(err.Error())
	}

	// Sprzątanie idzie przy okazji przyjęcia — jedyna pewna chwila na to.
	a.posprzatajNagrania(ctx, teraz.UnixMilli())

	// Odcinek nasłuchu idzie od razu do rozpoznania tą samą komendą.
	if z.WindowId != nil {
		a.ogloszOdcinekNasluchu(ctx, *z.WindowId, zapisane.Sciezka)
	}

	odpowiedz := shared.SpeechAudioUploadResponse{
		AudioRef:  zapisane.Sciezka,
		SizeBytes: int(zapisane.RozmiarBajtow),
	}
	odpowiedz.ExpiresAt = zapisane.Wygasa
	return odpowiedz, nil
}

// OddajNagranie obsługuje `speech.audio.fetch`. Wycinka czasowego rdzeń nie
// wykonuje: żądanie z polami `rangeStartMs` i `rangeEndMs` kończy się odmową
// nazywającą powód, nie całym nagraniem podanym jako wycinek.
func (a *adapterMowy) OddajNagranie(ctx context.Context,
	z shared.SpeechAudioFetchRequest) (shared.SpeechAudioFetchResponse, error) {

	odnosnik := strings.TrimSpace(z.AudioRef)
	if odnosnik == "" {
		return shared.SpeechAudioFetchResponse{}, bladWskazaniaNagrania(
			"żądanie bez odnośnika — nie wiadomo, czego odsłuchać")
	}
	if (z.RangeStartMs != nil && *z.RangeStartMs > 0) || (z.RangeEndMs != nil && *z.RangeEndMs > 0) {
		return shared.SpeechAudioFetchResponse{}, bladWycinkaNagrania()
	}

	sciezka, typTresci, err := a.rozstrzygnijOdnosnikNagrania(ctx, odnosnik)
	if err != nil {
		return shared.SpeechAudioFetchResponse{}, err
	}
	bajty, err := os.ReadFile(sciezka)
	if err != nil {
		return shared.SpeechAudioFetchResponse{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeNotFound, "nagrania mowy: bajtów nagrania nie ma pod odnośnikiem "+
				odnosnik+": "+err.Error()))
	}
	return shared.SpeechAudioFetchResponse{
		Audio:       base64.StdEncoding.EncodeToString(bajty),
		ContentType: typTresci,
		SizeBytes:   len(bajty),
	}, nil
}

// rozstrzygnijOdnosnikNagrania sprawdza, czy odnośnik wolno odsłuchać, i oddaje
// ścieżkę wraz z typem treści. Dwie drogi są zamknięte: wiersz rejestru nagrań
// albo plik pod katalogiem danych rdzenia; ścieżka spoza obu jest odmową
// uprawnienia.
func (a *adapterMowy) rozstrzygnijOdnosnikNagrania(ctx context.Context,
	odnosnik string) (string, string, error) {

	if a.nagrania != nil {
		wiersz, err := a.nagrania.NagranieMowyPoSciezce(ctx, odnosnik)
		if err == nil {
			return wiersz.Sciezka, wiersz.TypTresci, nil
		}
	}

	sciezka, err := filepath.Abs(odnosnik)
	if err != nil {
		return "", "", bladWskazaniaNagrania("odnośnika nie da się rozwinąć do ścieżki: " + err.Error())
	}
	podstawa := strings.TrimSpace(a.katalogDanych)
	if podstawa == "" {
		return "", "", bladZapleczaNagran("serwer nie zna katalogu danych — nie ma jak " +
			"rozstrzygnąć, czy ten odnośnik wolno odsłuchać")
	}
	korzen, err := filepath.Abs(podstawa)
	if err != nil {
		return "", "", bladZapleczaNagran("katalogu danych nie da się rozwinąć: " + err.Error())
	}
	wzgledna, err := filepath.Rel(korzen, sciezka)
	if err != nil || wzgledna == ".." || strings.HasPrefix(wzgledna, ".."+string(filepath.Separator)) {
		return "", "", protocol.JakoError(protocol.NowyBlad(shared.ErrorCodePermissionDenied,
			"nagrania mowy: odnośnik "+odnosnik+" nie jest nagraniem serwera — "+
				"odsłuch oddaje wyłącznie nagrania przyjęte przez speech.audio.upload "+
				"oraz wytworzone przez translate.speech.synthesize, a nie dowolny plik dysku"))
	}
	if !mowa.FormatPrzyjmowany(sciezka) {
		return "", "", bladWskazaniaNagrania("plik " + odnosnik +
			" nie ma rozszerzenia nagrania przyjmowanego przez serwer")
	}
	return sciezka, typTresciZeSciezki(sciezka), nil
}

// katalogNagran zakłada, gdy trzeba, katalog na przyjęte nagrania. Pusty
// katalog danych jest odmową, a nie powodem do wybrania czegoś z własnej
// głowy, tak samo jak przy syntezie mowy.
func (a *adapterMowy) katalogNagran() (string, error) {
	podstawa := strings.TrimSpace(a.katalogDanych)
	if podstawa == "" {
		return "", bladZapleczaNagran("serwer nie zna katalogu danych — nie ma gdzie zapisać " +
			"nagrania; naprawa: wskazać katalog danych przełącznikiem -dane albo zmienną " +
			"DANACO_KATALOG_DANYCH")
	}
	katalog := filepath.Join(podstawa, katalogNagranMowy)
	if err := os.MkdirAll(katalog, 0o700); err != nil {
		return "", bladZapleczaNagran("nie można założyć katalogu nagrań: " + err.Error())
	}
	return katalog, nil
}

// czyZapisywacNagrania rozstrzyga trwałość nagrania: pole żądania, a w jego
// braku nastawa `mowa_zapis_nagran`. Nastawa nierozstrzygnięta znaczy „nie
// zapisuj trwale" — zapis nagrania głosu jest decyzją, którą Operator podejmuje
// jawnie, a nie domyślnie.
func (a *adapterMowy) czyZapisywacNagrania(zZadania *bool) bool {
	if zZadania != nil {
		return *zZadania
	}
	if a.rozstrzygacz == nil {
		return false
	}
	wynik := a.rozstrzygacz.Rozstrzygnij(kontekstZasieguMowy(), kluczZapisuNagran)
	switch strings.ToLower(strings.TrimSpace(wynik.Wartosc)) {
	case "1", "true", "tak":
		return true
	default:
		return false
	}
}

// posprzatajNagrania kasuje nagrania, których czas minął — bajty i wiersz.
//
// Niepowodzenie sprzątania nie przewraca przyjęcia: nagranie właśnie legło na
// dysku i jest Operatorowi potrzebne teraz, a stare pliki poczekają do
// następnego razu.
func (a *adapterMowy) posprzatajNagrania(ctx context.Context, teraz int64) {
	if a.nagrania == nil {
		return
	}
	wygasle, err := a.nagrania.NagraniaMowyWygasle(ctx, teraz)
	if err != nil {
		return
	}
	for _, nagranie := range wygasle {
		_ = os.Remove(nagranie.Sciezka)
		_ = a.nagrania.UsunNagranieMowy(ctx, nagranie.Sciezka)
	}
}

// typyTresciNagran wiąże typ treści przysłany przez okno z rozszerzeniem
// pliku, które przyjmuje silnik mowy. Zbiór jest zamknięty i pokrywa się
// z `mowa.FormatyNagran`.
var typyTresciNagran = map[string]string{
	"audio/wav":              ".wav",
	"audio/x-wav":            ".wav",
	"audio/wave":             ".wav",
	"audio/ogg":              ".ogg",
	"audio/opus":             ".ogg",
	"audio/mp4":              ".m4a",
	"audio/m4a":              ".m4a",
	"audio/x-m4a":            ".m4a",
	"audio/webm":             ".webm",
	"video/webm":             ".webm",
	"audio/webm;codecs=opus": ".webm",
}

// rozszerzenieNagrania przekłada typ treści na rozszerzenie pliku.
//
// Parametry typu (`;codecs=opus`) odcina się przed dopasowaniem: przeglądarka
// dokłada je sama, a rodzaj kontenera rozstrzyga człon przed średnikiem.
func rozszerzenieNagrania(typTresci string) (string, error) {
	klucz := strings.ToLower(strings.TrimSpace(typTresci))
	if klucz == "" {
		return "", bladWskazaniaNagrania("żądanie bez typu treści — serwer nie zgaduje formatu " +
			"nagrania, bo od formatu zależy, czy silnik mowy plik otworzy")
	}
	if rozszerzenie, jest := typyTresciNagran[klucz]; jest {
		return rozszerzenie, nil
	}
	if podzial := strings.IndexByte(klucz, ';'); podzial > 0 {
		if rozszerzenie, jest := typyTresciNagran[strings.TrimSpace(klucz[:podzial])]; jest {
			return rozszerzenie, nil
		}
	}
	return "", bladWskazaniaNagrania("nie znam typu treści „" + typTresci +
		"” — serwer przyjmuje: audio/wav, audio/ogg, audio/mp4, audio/webm")
}

// typTresciZeSciezki odgaduje typ treści z rozszerzenia pliku. Używane wyłącznie
// dla nagrań spoza rejestru (synteza mowy), które własnego typu nie niosą.
func typTresciZeSciezki(sciezka string) string {
	switch strings.ToLower(filepath.Ext(sciezka)) {
	case ".wav":
		return "audio/wav"
	case ".ogg":
		return "audio/ogg"
	case ".m4a":
		return "audio/mp4"
	case ".webm":
		return "audio/webm"
	default:
		return "application/octet-stream"
	}
}

// bladWskazaniaNagrania nazywa niepoprawne żądanie rodziny nagrań: brak
// bajtów, zły format albo odnośnik, którego rdzeń nie rozpoznaje.
func bladWskazaniaNagrania(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"nagrania mowy: "+powod))
}

// bladZapleczaNagran nazywa brak po stronie rdzenia: nie wpięty rejestr
// nagrań, brak katalogu danych albo usterka zapisu na dysku.
func bladZapleczaNagran(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"nagrania mowy: "+powod))
}

// bladWycinkaNagrania nazywa odmowę wycięcia fragmentu — rdzeń oddaje
// wyłącznie nagranie w całości, nigdy zakresu czasowego.
func bladWycinkaNagrania() error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"nagrania mowy: serwer nie wycina fragmentu nagrania — wycięcie z zapisu "+
			"skompresowanego wymaga przekodowania, a przekodowany dźwięk nie jest tym "+
			"samym, który przyszedł; przewijanie odpowiedzi robi okno na pobranym "+
			"nagraniu, u siebie"))
}
