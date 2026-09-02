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
	katalogNagranMowy = "nagrania-mowy"

	przedrostekNagraniaMowy = "nagranie-"

	// Dwadzieścia pięć megabajtów to około pół godziny mowy w Opusie.
	granicaNagraniaBajtow = 25 * 1024 * 1024

	zycieNagraniaTymczasowego = time.Hour

	// kluczZapisuNagran: kontrakt `speech.audio.upload` odsyła do tej nastawy przy pustym
	// polu `retain`.
	kluczZapisuNagran = "mowa_zapis_nagran"
)

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
	trwale := a.czyZapisywacNagrania(ctx, z.Retain)
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

func (a *adapterMowy) rozstrzygnijOdnosnikNagrania(ctx context.Context,
	odnosnik string) (string, string, error) {

	if a.nagrania != nil {
		wiersz, err := a.nagrania.NagranieMowyPoSciezce(ctx, odnosnik)
		if err == nil {
			return wiersz.Sciezka, wiersz.TypTresci, nil
		}
	}
	if a.wKataloguNagranPrzyjetych(odnosnik) {
		return "", "", bladObcegoNagrania(odnosnik)
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

// Katalog nagrań przyjętych jest wyłącznością rejestru: bajty spod niego wydaje wyłącznie
// wiersz `nagranie_mowy` zawężony do konta żądania (migracja 484).
func (a *adapterMowy) wKataloguNagranPrzyjetych(odnosnik string) bool {
	podstawa := strings.TrimSpace(a.katalogDanych)
	if podstawa == "" {
		return false
	}
	korzen, err := filepath.Abs(filepath.Join(podstawa, katalogNagranMowy))
	if err != nil {
		return false
	}
	sciezka, err := filepath.Abs(strings.TrimSpace(odnosnik))
	if err != nil {
		return false
	}
	wzgledna, err := filepath.Rel(korzen, sciezka)
	if err != nil {
		return false
	}
	return wzgledna != ".." && !strings.HasPrefix(wzgledna, ".."+string(filepath.Separator))
}

func bladObcegoNagrania(odnosnik string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodePermissionDenied,
		"nagrania mowy: nagranie "+odnosnik+" nie należy do konta tego żądania — "+
			"nagrania przyjęte przez speech.audio.upload wydaje wyłącznie rejestr "+
			"i wyłącznie właścicielowi wiersza"))
}

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

func (a *adapterMowy) czyZapisywacNagrania(ctx context.Context, zZadania *bool) bool {
	if zZadania != nil {
		return *zZadania
	}
	if a.rozstrzygacz == nil {
		return false
	}
	wynik := a.rozstrzygacz.Rozstrzygnij(kontekstZasieguMowy(ctx), kluczZapisuNagran)
	switch strings.ToLower(strings.TrimSpace(wynik.Wartosc)) {
	case "1", "true", "tak":
		return true
	default:
		return false
	}
}

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

// Zbiór zamknięty, pokrywa się z `mowa.FormatyNagran`.
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

// Parametry typu (`;codecs=opus`) odcina się przed dopasowaniem — dokłada je przeglądarka.
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

func bladWskazaniaNagrania(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"nagrania mowy: "+powod))
}

func bladZapleczaNagran(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"nagrania mowy: "+powod))
}

func bladWycinkaNagrania() error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"nagrania mowy: serwer nie wycina fragmentu nagrania — wycięcie z zapisu "+
			"skompresowanego wymaga przekodowania, a przekodowany dźwięk nie jest tym "+
			"samym, który przyszedł; przewijanie odpowiedzi robi okno na pobranym "+
			"nagraniu, u siebie"))
}
