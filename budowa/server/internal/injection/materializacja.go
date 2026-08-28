// Materializacja domyka lukę między treścią ustawień a argumentami procesu
// programu `claude`, zapisując treść JSON do pliku tymczasowego przed
// uruchomieniem.
package injection

import (
	"fmt"
	"os"
	"strings"
)

// zmaterializujUstawienia zapisuje na dysk te pola ustawień, które niosą treść
// JSON, i podmienia je na ścieżki utworzonych plików tymczasowych.
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
		// Nowa tablica, bo kopia ustawień dzieli tablicę bazową z oryginałem.
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
// bezwzględną ścieżkę widoczną dla procesu.
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
