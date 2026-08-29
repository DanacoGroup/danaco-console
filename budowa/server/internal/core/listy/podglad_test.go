package listy

import (
	"encoding/base64"
	"os"
	"testing"
)

// TestPodgladListu odkłada postać graficzną listu aktywacji do pliku, żeby dało
// się ją obejrzeć przeglądarką. Nie sprawdza niczego — jest narzędziem
// oglądania, nie sprawdzianem: wysyłka dołącza znak częścią listu, a przeglądarka
// odwołania `cid:` nie zna, więc podgląd wpisuje obraz w treść.
func TestPodgladListu(t *testing.T) {
	html := PodgladDanymiPrzykladowymi(
		"data:image/png;base64," + base64.StdEncoding.EncodeToString(ZnakMarki))
	if err := os.WriteFile("/srv/podglad/wykaz-list.html", []byte(html), 0o644); err != nil {
		t.Log("nie odłożono podglądu:", err)
	}
}
