// Egzekucja trzech punktów izolacji, rozstrzyganych przydziałem zasobu
// wykonawczego: konta i tokenu, instancji procesu modelu oraz serwera
// wykonania. Zasób, którego właścicielem nie jest to okno, zostaje
// odrzucony.
package session

import (
	"strings"

	"danacoconsole/server/internal/konfig"
)

// Zasob to jeden zasób wykonawczy, po który sięga okno, wraz z informacją
// o tym, do jakiego okna należy.
type Zasob struct {
	// Identyfikator zasobu — kod konta, identyfikator procesu, nazwa serwera.
	Identyfikator string
	// Wlasciciel wskazuje okno, dla którego zasób został wydzielony; puste
	// znaczy zasób wspólny platformy.
	Wlasciciel string
}

// Przydzial opisuje zasoby jednego uruchomienia okna komunikacji: konto,
// proces oraz serwer wykonania modelu.
type Przydzial struct {
	// IdOkna — okno, dla którego przydział powstał.
	IdOkna string
	// Konto niesie profil uwierzytelnienia kanału — odwołanie, nigdy sekret.
	Konto Zasob
	// Proces niesie instancję procesu wykonawczego modelu.
	Proces Zasob
	// Serwer niesie serwer, na którym proces ma pracować.
	Serwer Zasob
}

// SprawdzPrzydzial egzekwuje wyłączność zasobów wymaganą przez zasady izolacji.
// Zakres wyłączony nie ogranicza niczego — zasób wspólny jest wtedy stanem
// prawidłowym.
func SprawdzPrzydzial(zasady Zasady, przydzial Przydzial) error {
	sprawdzenia := []struct {
		wymagana bool
		klucz    string
		nazwa    string
		zasob    Zasob
	}{
		{zasady.KontoIToken, konfig.KluczIzolacjaKontoIToken, "konto kanału modelu", przydzial.Konto},
		{zasady.ModelProcesu, konfig.KluczIzolacjaModelProcesu, "instancja procesu modelu", przydzial.Proces},
		{zasady.SerwerWykonania, konfig.KluczIzolacjaSerwerWykonania, "serwer wykonania", przydzial.Serwer},
	}
	for _, sprawdzenie := range sprawdzenia {
		if !sprawdzenie.wymagana {
			continue
		}
		if err := sprawdzWylacznosc(sprawdzenie.klucz, sprawdzenie.nazwa,
			przydzial.IdOkna, sprawdzenie.zasob); err != nil {
			return err
		}
	}
	return nil
}

// sprawdzWylacznosc odrzuca zasób wspólny platformy oraz zasób wydzielony
// innemu oknu niż to sprawdzane.
func sprawdzWylacznosc(klucz, nazwa, idOkna string, zasob Zasob) error {
	wlasciciel := strings.TrimSpace(zasob.Wlasciciel)
	if wlasciciel == "" {
		return NoweNaruszenie(klucz, nazwa+" "+opisZasobu(zasob)+
			" jest zasobem wspólnym platformy, a okno "+idOkna+" ma mieć własny")
	}
	if wlasciciel != strings.TrimSpace(idOkna) {
		return NoweNaruszenie(klucz, nazwa+" "+opisZasobu(zasob)+
			" należy do okna "+wlasciciel+", a sięga po niego okno "+idOkna)
	}
	return nil
}

// opisZasobu zwraca identyfikator zasobu w postaci czytelnej, nadającej się
// do komunikatu naruszenia izolacji.
func opisZasobu(zasob Zasob) string {
	if identyfikator := strings.TrimSpace(zasob.Identyfikator); identyfikator != "" {
		return identyfikator
	}
	return "bez identyfikatora"
}
