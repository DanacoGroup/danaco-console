// Odpowiedzialność pliku: publiczny adres Konsoli, pod którym Operator otwiera
// okna wskazywane odsyłaczem z listu.
package konfiguracja

import (
	"net/url"
	"strings"
)

// AdresKonsoliDomyslny to adres, pod którym stoi Konsola platformy. Wartość jest
// ta sama, którą niosą szablony listów transakcyjnych
// (`design/06-poczta-transakcyjna`) — list i rdzeń muszą wskazywać jedno miejsce.
//
// Adres stoi tu, a nie w pliku środowiska, z tego samego powodu co konto
// nadawcze (`konto_nadawcze.go`): bez niego przycisk w liście aktywacji nie
// prowadzi nigdzie, a świeża instalacja miałaby to naprawiać czynnością
// wdrożeniową. Zmienna `DANACO_ADRES_KONSOLI` tę wartość nadpisuje.
//
// To NIE jest adres nasłuchu rdzenia — `DANACO_ADRES` niesie gniazdo, na którym
// rdzeń słucha, i bywa adresem pętli zwrotnej. Odsyłacz w liście musi działać
// z maszyny Operatora, nie z maszyny wdrożenia.
const AdresKonsoliDomyslny = "https://console.danaco-group.pl"

// SciezkaAktywacji to okno, w którym Operator wprowadza kod z listu aktywacji.
const SciezkaAktywacji = "/aktywacja"

// SciezkaWycofaniaZmiany to trasa, którą otwiera odsyłacz z listu ostrzegającego
// o zamówionej zmianie adresu. Droga wycofania stoi w zapytaniu, bo Operator
// klika ją z poczty, nie przepisuje do okna.
const SciezkaWycofaniaZmiany = "/wycofaj-zmiane-adresu"

/*
AdresAktywacji składa pełny odsyłacz do okna aktywacji.

Ukośnik na styku jest ucinany, żeby adres podany ze znakiem końcowym nie dawał
podwójnego ukośnika — ten w odsyłaczu listu widzi Operator.
*/
func AdresAktywacji(adresKonsoli string) string {
	if adresKonsoli == "" {
		adresKonsoli = AdresKonsoliDomyslny
	}
	return strings.TrimRight(adresKonsoli, "/") + SciezkaAktywacji
}

// AdresWycofaniaZmiany składa odsyłacz wycofania wraz z drogą w zapytaniu.
func AdresWycofaniaZmiany(adresKonsoli, droga string) string {
	if adresKonsoli == "" {
		adresKonsoli = AdresKonsoliDomyslny
	}
	return strings.TrimRight(adresKonsoli, "/") + SciezkaWycofaniaZmiany + "?droga=" + url.QueryEscape(droga)
}

// SciezkaPrzebiegu to okno Execution Monitor otwarte na wskazanym przebiegu.
const SciezkaPrzebiegu = "/automatyki/przebieg"

// AdresPrzebiegu składa odsyłacz do okna przebiegu wskazanego identyfikatorem.
func AdresPrzebiegu(adresKonsoli, kod string) string {
	if adresKonsoli == "" {
		adresKonsoli = AdresKonsoliDomyslny
	}
	return strings.TrimRight(adresKonsoli, "/") + SciezkaPrzebiegu + "/" + url.PathEscape(kod)
}
