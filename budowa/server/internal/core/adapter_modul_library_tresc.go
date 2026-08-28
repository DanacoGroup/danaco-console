// Odpowiedzialność pliku: moduł Library — rozstrzygnięcie treści komendy:
// co dokładnie zostaje utrwalone i czym to zmierzyć. Obsługuje
// `library.file.upload` i `library.version.add`.
package core

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// trescWgrania utrwala treść żądania i oddaje trójkę: odwołanie do bajtów,
// ich rozmiar i sumę kontrolną. Ścieżka ma pierwszeństwo przed treścią,
// a żądanie bez jednego i drugiego nie jest błędem — wraca sama pustka.
func (a *adapterBiblioteki) trescWgrania(trescBase64, sciezkaZrodlowa *string) (*string, *int64, *string, error) {
	if !bezWartosci(sciezkaZrodlowa) {
		if a.magazyn == nil {
			return nil, nil, nil, bladZapisuTresciBiblioteki(
				"magazyn treści nie jest wpięty — nie ma gdzie wciągnąć pliku ze ścieżki")
		}
		odwolanie, rozmiar, suma, err := a.magazyn.ZapiszZePliku(*sciezkaZrodlowa)
		if err != nil {
			// Pomyłka wskazującego nie jest awarią rdzenia — ponowienie by nie pomogło.
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
// kontrolną. Suma jest zarazem nazwą bloba w magazynie. Żądanie bez treści
// nie jest błędem — wraca sama pustka, bo treść mogła przyjść ścieżką.
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
// Kod jest wewnętrzny, bo to awaria zapisu po stronie rdzenia.
func bladZapisuTresciBiblioteki(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError, "moduł Library: "+powod))
}

// bladBrakuTresciBiblioteki nazywa sytuację, w której podgląd nie ma z czego
// się zbudować — plik bez zapisanego odwołania do treści. Podgląd zgłasza to
// wprost zamiast zmyślać zawartość.
func bladBrakuTresciBiblioteki(kod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"moduł Library: plik "+kod+" nie ma odwołania do treści — podgląd niemożliwy"))
}
