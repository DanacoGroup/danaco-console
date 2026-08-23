package injection

import (
	"fmt"
	"os"
	"strings"
)

// Materializacja domyka lukę między treścią ustawień a argumentami procesu.
// Przełączniki --settings i --mcp-config programu `claude` wskazują plik na
// dysku. Warstwa sesji składa jednak wartości tych pól z obszarów konfiguracji
// (tools, permissions, mcp) i przekazuje je jako napis JSON — treść, nie
// ścieżkę. Gdyby taki napis trafił wprost do argv, program szukałby pliku o tej
// nazwie, nie znalazłby go i całe przekierowanie byłoby bezczynne.
//
// Dlatego przed uruchomieniem procesu zapisujemy treść do pliku tymczasowego
// i podmieniamy wartość na jego ścieżkę. Plik żyje tylko przez jeden przebieg
// tury — sprząta go zwrócony domknięciem porządek (defer w wykonajPrzebieg).
//
// Rozpoznanie treści od ścieżki jest jednoznaczne: konfiguracja JSON zaczyna
// się od `{` albo `[` (po odcięciu białych znaków), a ścieżka pliku nigdy tak
// nie zaczyna. Wartość rozpoznaną jako ścieżkę zostawiamy nietkniętą — plik
// zapisała już inna warstwa, a my tylko wskazujemy go procesowi.

// zmaterializujUstawienia zapisuje na dysk te pola ustawień, które niosą treść
// JSON (PlikUstawien, wpisy KonfiguracjaMCP), i podmienia je na ścieżki
// utworzonych plików tymczasowych. Zwraca kopię ustawień gotową do złożenia
// argumentów oraz porządek sprzątający pliki po turze. Kopia nie narusza
// oryginału — prowenancja nadal pokazuje pierwotną treść, a argv realną ścieżkę.
func zmaterializujUstawienia(u Ustawienia) (Ustawienia, func(), error) {
	var pliki []string
	sprzataj := func() {
		for _, sciezka := range pliki {
			_ = os.Remove(sciezka)
		}
	}

	if jestTrescJSON(u.PlikUstawien) {
		sciezka, err := zapiszTymczasowy("danaco-ustawienia-*.json", u.PlikUstawien)
		if err != nil {
			sprzataj()
			return u, func() {}, err
		}
		pliki = append(pliki, sciezka)
		u.PlikUstawien = sciezka
	}

	if len(u.KonfiguracjaMCP) > 0 {
		// Nowa tablica, bo kopia ustawień dzieli tablicę bazową z oryginałem —
		// podmiana w miejscu przeciekłaby do prowenancji i do wywołującego.
		zmienione := make([]string, len(u.KonfiguracjaMCP))
		for i, wpis := range u.KonfiguracjaMCP {
			if !jestTrescJSON(wpis) {
				zmienione[i] = wpis
				continue
			}
			sciezka, err := zapiszTymczasowy("danaco-mcp-*.json", wpis)
			if err != nil {
				sprzataj()
				return u, func() {}, err
			}
			pliki = append(pliki, sciezka)
			zmienione[i] = sciezka
		}
		u.KonfiguracjaMCP = zmienione
	}

	return u, sprzataj, nil
}

// jestTrescJSON rozstrzyga, czy wartość jest treścią konfiguracji (napisem
// JSON), czy ścieżką do pliku. Obiekt lub tablica JSON zaczyna się od `{`
// albo `[`; ścieżka pliku — nigdy. Wartość pusta nie jest treścią (przełącznika
// i tak nie będzie).
func jestTrescJSON(wartosc string) bool {
	przyciete := strings.TrimSpace(wartosc)
	if przyciete == "" {
		return false
	}
	return przyciete[0] == '{' || przyciete[0] == '['
}

// zapiszTymczasowy zapisuje treść do świeżego pliku tymczasowego i zwraca jego
// bezwzględną ścieżkę. Katalog tymczasowy systemu daje ścieżkę widoczną dla
// procesu niezależnie od jego katalogu roboczego. Przy każdym potknięciu
// sprzątamy po sobie, żeby nie zostawić pliku bez właściciela.
func zapiszTymczasowy(wzorzec, tresc string) (string, error) {
	plik, err := os.CreateTemp("", wzorzec)
	if err != nil {
		return "", fmt.Errorf("injection: utworzenie pliku tymczasowego %q: %w", wzorzec, err)
	}
	if _, err := plik.WriteString(tresc); err != nil {
		_ = plik.Close()
		_ = os.Remove(plik.Name())
		return "", fmt.Errorf("injection: zapis pliku tymczasowego %s: %w", plik.Name(), err)
	}
	if err := plik.Close(); err != nil {
		_ = os.Remove(plik.Name())
		return "", fmt.Errorf("injection: domknięcie pliku tymczasowego %s: %w", plik.Name(), err)
	}
	return plik.Name(), nil
}
