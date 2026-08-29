package listy

import (
	"encoding/base64"
	"os"
	"testing"
)

// TestPodgladListu odkłada postać graficzną listu aktywacji do pliku, żeby dało
// się ją obejrzeć przeglądarką. Nie sprawdza niczego — jest narzędziem
// oglądania: wysyłka dołącza znaki częściami listu, a przeglądarka odwołania
// `cid:` nie zna, więc podgląd wpisuje oba obrazy w treść.
func TestPodgladListu(t *testing.T) {
	wTresc := func(dane []byte) string {
		return "data:image/png;base64," + base64.StdEncoding.EncodeToString(dane)
	}
	html := PodgladDanymiPrzykladowymi(wTresc(ZnakJasny), wTresc(ZnakCiemny))
	if err := os.WriteFile("/srv/podglad/wykaz-list.html", []byte(html), 0o644); err != nil {
		t.Log("nie odłożono podglądu:", err)
	}
}
