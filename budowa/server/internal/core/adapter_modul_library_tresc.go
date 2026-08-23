// Odpowiedzialność pliku: moduł Library — rozstrzygnięcie treści komendy, czyli
// jedno miejsce odpowiadające na pytanie „co dokładnie zostaje utrwalone i czym
// to zmierzyć". Obsługuje `library.file.upload` i `library.version.add` razem,
// bo obie komendy przyjmują tę samą parę pól kontraktu (`contentBase64`,
// `sourcePath`) i muszą ją rozstrzygać tak samo — dwa rachunki nad jedną treścią
// byłyby dwiema prawdami o jej integralności. Dostęp do nośnika i miara treści
// to osobna odpowiedzialność od składania odpowiedzi kontraktu, stąd podział
// wobec `adapter_modul_library.go`.
//
// Treść zawsze ląduje w magazynie rdzenia. Gdyby `sourcePath` stawała się
// odwołaniem do treści, plik wgrany ścieżką i potem nadpisany na dysku
// zmieniałby treść swojej rzekomo utrwalonej wersji, a `library.version.restore`
// przywracałby wskaźnik do treści, której już nie ma. Wersja, po której nie da
// się odtworzyć zawartości, nie jest wersją — więc ścieżka jest źródłem bajtów,
// a miejscem ich składowania jest magazyn treści rdzenia
// (`adapter_modul_library_magazyn.go`), tak samo jak dla treści przysłanej
// base64. Ścieżka źródłowa nie znika bez śladu: zostaje przy pliku w kolumnie
// `sciezka` (`migracja_045_biblioteka.sql`) jako informacja, skąd plik przyszedł.
//
// Granicy rozmiaru nie ma. Obie drogi idą strumieniem albo jednym buforem już
// przysłanym przez klienta, więc duży plik nie ma jak przewrócić rdzenia
// pamięcią; odmowa opisuje brak (nie da się otworzyć, nie da się zapisać),
// nigdy zakaz („plik za duży").
package core

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// trescWgrania utrwala treść żądania i oddaje trójkę, którą wiersz pliku oraz
// wiersz wersji opisują tę samą treść: odwołanie do bajtów, ich rozmiar i sumę
// kontrolną.
//
// Ścieżka ma pierwszeństwo przed treścią: gdy żądanie niesie oba pola, bajty
// spod ścieżki są tym, co Operator wskazał, a base64 bywa wtedy skrótem
// podglądu przysłanym przez okno.
//
// Żądanie bez jednego i drugiego nie jest błędem — wraca sama pustka. Taka
// komenda oznacza wersję-znacznik (kamień milowy na treści bieżącej,
// `trescNowejWersji`) albo plik zakładany bez zawartości; odmowa należy tu
// wyłącznie do zapisu, który się nie powiódł.
func (a *adapterBiblioteki) trescWgrania(trescBase64, sciezkaZrodlowa *string) (*string, *int64, *string, error) {
	if !bezWartosci(sciezkaZrodlowa) {
		if a.magazyn == nil {
			return nil, nil, nil, bladZapisuTresciBiblioteki(
				"magazyn treści nie jest wpięty — nie ma gdzie wciągnąć pliku ze ścieżki")
		}
		odwolanie, rozmiar, suma, err := a.magazyn.ZapiszZePliku(*sciezkaZrodlowa)
		if err != nil {
			// Pomyłka wskazującego nie jest awarią rdzenia. Ścieżka przychodzi
			// z żądania: literówka, plik usunięty, brak praw albo wskazany
			// katalog to brak po stronie wołającego — żądanie nie ma prawa się
			// udać przy żadnym ponowieniu. Kod `internal_error` z
			// `retryable:true` zapętliłby klienta ponawiającego żądanie.
			if errors.Is(err, ErrZrodloTresciNieczytelne) {
				return nil, nil, nil, bladWskazaniaBiblioteki(
					"nie da się wciągnąć treści spod wskazanej ścieżki: " + err.Error())
			}
			return nil, nil, nil, bladZapisuTresciBiblioteki(
				"nie można wciągnąć treści spod ścieżki: " + err.Error())
		}
		return &odwolanie, &rozmiar, &suma, nil
	}

	bajty, rozmiar, suma, err := trescZadania(trescBase64)
	if err != nil || len(bajty) == 0 {
		return nil, rozmiar, suma, err
	}
	if a.magazyn == nil {
		return nil, nil, nil, bladZapisuTresciBiblioteki(
			"magazyn treści nie jest wpięty — nie ma gdzie zapisać wgranych bajtów")
	}
	odwolanie, err := a.magazyn.Zapisz(bajty, *suma)
	if err != nil {
		return nil, nil, nil, bladZapisuTresciBiblioteki("nie można utrwalić treści pliku: " + err.Error())
	}
	return &odwolanie, rozmiar, suma, nil
}

// trescZadania rozbiera treść przysłaną base64 na bajty, rozmiar i sumę
// kontrolną. Suma jest zarazem nazwą bloba w magazynie, więc jedno wgranie
// liczy sha256 dokładnie raz.
//
// Żądanie bez treści nie jest błędem — wraca sama pustka, bo treść mogła przyjść
// ścieżką (`sourcePath`).
func trescZadania(trescBase64 *string) ([]byte, *int64, *string, error) {
	if bezWartosci(trescBase64) {
		return nil, nil, nil, nil
	}
	bajty, err := base64.StdEncoding.DecodeString(*trescBase64)
	if err != nil {
		return nil, nil, nil, bladWskazaniaBiblioteki("treść pliku nie jest poprawnym base64: " + err.Error())
	}
	dlugosc := int64(len(bajty))
	suma := sha256.Sum256(bajty)
	tekst := hex.EncodeToString(suma[:])
	return bajty, &dlugosc, &tekst, nil
}

// bladZapisuTresciBiblioteki nazywa niepowodzenie utrwalenia treści w magazynie
// rdzenia — nośnik pełny, brak praw do katalogu danych, magazyn niewpięty.
// Ścieżki źródłowej wskazanej przez wołającego to nie dotyczy: jej brak jest
// pomyłką żądania i wraca `bladWskazaniaBiblioteki` (wyżej). Kod jest
// wewnętrzny, bo to awaria zapisu po stronie rdzenia. Komenda odmawia
// w całości: plik, którego bajtów nie ma nigdzie, nie może trafić do wykazu
// jako wgrany.
func bladZapisuTresciBiblioteki(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError, "moduł Library: "+powod))
}

// bladBrakuTresciBiblioteki nazywa sytuację, w której podgląd nie ma z czego
// się zbudować — plik bez zapisanego odwołania do treści. Z wpiętym magazynem
// wgranie zawsze zostawia odwołanie, więc tu trafiają wiersze starsze oraz
// wersje-znaczniki założone na pliku bez treści. Podgląd zgłasza to wprost
// zamiast zmyślać zawartość.
func bladBrakuTresciBiblioteki(kod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"moduł Library: plik "+kod+" nie ma odwołania do treści — podgląd niemożliwy"))
}
