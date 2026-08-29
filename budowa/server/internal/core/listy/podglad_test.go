package listy

import (
	"os"
	"testing"
)

// TestPodgladListu odkłada postać graficzną listu do pliku, żeby dało się ją
// obejrzeć przeglądarką. Nie sprawdza niczego — jest narzędziem oglądania,
// nie sprawdzianem.
func TestPodgladListu(t *testing.T) {
	html := Zloz(TrescListu{
		Naglowek:     "Potwierdzenie adresu e-mail",
		Wstep:        "Konto operator zostało założone i oczekuje na potwierdzenie tego adresu.",
		EtykietaKodu: "Kod potwierdzający",
		Kod:          "482913",
		Polecenie:    "Wprowadź go w oknie rejestracji, aby zakończyć zakładanie konta i wejść do platformy. Kod jest jednorazowy i zachowuje ważność przez godzinę.",
		Nota:         "Ten adres będzie później jedyną drogą odzyskania konta.",
	})
	if err := os.WriteFile("/srv/podglad/wykaz-list.html", []byte(html), 0o644); err != nil {
		t.Log("nie odłożono podglądu:", err)
	}
}
