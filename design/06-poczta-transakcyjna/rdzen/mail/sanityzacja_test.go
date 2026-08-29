package mail

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestCleanUntrusted(t *testing.T) {
	przypadki := []struct {
		nazwa, wejscie, oczekiwane string
	}{
		{"złamania wiersza na spacje", "Temat\r\nBcc: kto@przyklad.pl", "Temat Bcc: kto@przyklad.pl"},
		{"znaki sterujące usunięte", "Temat\x00\x07 sprawy", "Temat sprawy"},
		{"tabulator na spację", "Temat\tsprawy", "Temat sprawy"},
		{"białe znaki obcięte", "   Temat   ", "Temat"},
		{"diakrytyki zachowane", "Prośba o zwiększenie limitu", "Prośba o zwiększenie limitu"},
	}
	for _, p := range przypadki {
		t.Run(p.nazwa, func(t *testing.T) {
			if got := cleanUntrusted(p.wejscie); got != p.oczekiwane {
				t.Errorf("cleanUntrusted(%q) = %q, oczekiwano %q", p.wejscie, got, p.oczekiwane)
			}
		})
	}
}

// TestCleanUntrustedSkracaPoZnakach pilnuje, żeby cięcie nie rozbiło
// polskiego znaku diakrytycznego na końcu wartości.
func TestCleanUntrustedSkracaPoZnakach(t *testing.T) {
	wejscie := strings.Repeat("ą", 200)
	got := cleanUntrusted(wejscie)
	if n := utf8.RuneCountInString(got); n != maxUntrustedLen {
		t.Errorf("długość = %d znaków, oczekiwano %d", n, maxUntrustedLen)
	}
	if !utf8.ValidString(got) {
		t.Error("wynik nie jest poprawnym ciągiem UTF-8")
	}
}

func TestEscapeForHTML(t *testing.T) {
	got := escapeForHTML(`<img src=x onerror="alert(1)">`)
	if strings.Contains(got, "<img") || strings.Contains(got, `"`) {
		t.Errorf("znaczniki bez ucieczki: %q", got)
	}
}

func TestSubstituteBrakWartosci(t *testing.T) {
	_, err := substitute("kod {{code}} i {{expiry_minutes}}", map[string]string{"code": "418402"})
	if err == nil {
		t.Fatal("oczekiwano błędu")
	}
	if !strings.Contains(err.Error(), "expiry_minutes") {
		t.Errorf("komunikat nie wskazuje brakującej zmiennej: %v", err)
	}
}

func TestSubstitutePodstawia(t *testing.T) {
	got, err := substitute("kod {{code}}", map[string]string{"code": "418 402"})
	if err != nil {
		t.Fatalf("substitute: %v", err)
	}
	if got != "kod 418 402" {
		t.Errorf("wynik = %q", got)
	}
}

// TestSubstituteNiePodstawiaRekurencyjnie sprawdza, że wartość zawierająca
// ciąg wyglądający jak zmienna nie jest podstawiana ponownie.
func TestSubstituteNiePodstawiaRekurencyjnie(t *testing.T) {
	got, err := substitute("{{a}}", map[string]string{"a": "{{b}}", "b": "x"})
	if err != nil {
		t.Fatalf("substitute: %v", err)
	}
	if got != "{{b}}" {
		t.Errorf("wynik = %q, oczekiwano {{b}}", got)
	}
}

func TestPlaceholders(t *testing.T) {
	got := placeholders("{{b}} {{a}} {{a}}")
	if len(got) != 2 || got[0] != "a" || got[1] != "b" {
		t.Errorf("placeholders = %v", got)
	}
}
