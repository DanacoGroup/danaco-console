// Odpowiedzialność pliku: egzekucja trzech wymiarów izolacji kontekstu —
// historii wymiany wiadomości, pamięci długoterminowej i bieżącego stanu
// roboczego.
//
// Wymiar odrębny znaczy, że treść należąca do innego okna nie wchodzi do tury
// tego okna: nowa karta zaczyna z pustą historią, a pamięć jednego zasięgu
// pozostaje niewidoczna w innym. Egzekucja polega więc na odrzuceniu źródła
// cudzego, nie na cichym pominięciu go — Operator ma wiedzieć, że tura miała
// zaciągnąć treść spoza okna.
//
// Wymiar współdzielony nie ogranicza niczego; jest jawną decyzją Operatora
// o tym, że ta sama treść zasila kilka zasięgów.
package core

import (
	"strings"

	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/session"
)

// zasiegKontekstu wskazuje odbiorcę treści: okno komunikacji i jego sesję.
type zasiegKontekstu struct {
	// IdOkna — okno, dla którego składana jest tura.
	IdOkna string
	// IdSesji — sesja okna.
	IdSesji string
}

// zrodloKontekstu to jedna pozycja dokładana do tury wraz ze wskazaniem, czyja
// jest. Wymiar bierze się ze stałych pakietu konfig.
type zrodloKontekstu struct {
	// Wymiar izolacji: historia, pamięć albo kontekst.
	Wymiar string
	// IdOkna — okno, do którego treść należy. Puste znaczy treść spoza okien:
	// zasób zasięgu szerszego, wspólny dla kilku okien.
	IdOkna string
	// IdSesji — sesja, do której treść należy.
	IdSesji string
	// Odwolanie nazywa treść w komunikacie o naruszeniu; nie niesie jej samej.
	Odwolanie string
}

// sprawdzZrodlaKontekstu odrzuca pierwsze źródło naruszające odrębność wymiaru.
// Źródło o wymiarze spoza trzech punktów izolacji przechodzi bez zmian: nie ma
// ustawienia, które by je opisywało, a odczyt nie jest bramą.
func sprawdzZrodlaKontekstu(zasady session.Zasady, cel zasiegKontekstu,
	zrodla []zrodloKontekstu) error {

	for _, zrodlo := range zrodla {
		klucz, odrebny := wymiarOdrebny(zasady, zrodlo.Wymiar)
		if klucz == "" || !odrebny {
			continue
		}
		if err := sprawdzPrzynaleznosc(klucz, cel, zrodlo); err != nil {
			return err
		}
	}
	return nil
}

// sprawdzPrzynaleznosc odrzuca źródło należące do innego okna oraz źródło
// wspólne, którego przy wymiarze odrębnym w turze być nie może.
func sprawdzPrzynaleznosc(klucz string, cel zasiegKontekstu, zrodlo zrodloKontekstu) error {
	wlasciciel := strings.TrimSpace(zrodlo.IdOkna)
	if wlasciciel == strings.TrimSpace(cel.IdOkna) && wlasciciel != "" {
		return nil
	}
	if wlasciciel == "" {
		return session.NoweNaruszenie(klucz, opisZrodla(zrodlo)+
			" jest zasobem wspólnym zasięgu, a okno "+cel.IdOkna+" ma ten wymiar odrębny")
	}
	return session.NoweNaruszenie(klucz, opisZrodla(zrodlo)+" należy do okna "+
		wlasciciel+", a wchodzi do tury okna "+cel.IdOkna)
}

// wymiarOdrebny zwraca klucz punktu izolacji odpowiadający wymiarowi oraz
// odpowiedź, czy wymiar pozostaje odrębny.
func wymiarOdrebny(zasady session.Zasady, wymiar string) (string, bool) {
	switch strings.TrimSpace(wymiar) {
	case konfig.KluczIzolacjaHistoria:
		return konfig.KluczIzolacjaHistoria, zasady.HistoriaOdrebna
	case konfig.KluczIzolacjaPamiec:
		return konfig.KluczIzolacjaPamiec, zasady.PamiecOdrebna
	case konfig.KluczIzolacjaKontekst:
		return konfig.KluczIzolacjaKontekst, zasady.KontekstOdrebny
	}
	return "", false
}

// opisZrodla nazywa źródło w komunikacie o naruszeniu.
func opisZrodla(zrodlo zrodloKontekstu) string {
	if odwolanie := strings.TrimSpace(zrodlo.Odwolanie); odwolanie != "" {
		return "źródło " + odwolanie
	}
	return "źródło bez odwołania"
}
