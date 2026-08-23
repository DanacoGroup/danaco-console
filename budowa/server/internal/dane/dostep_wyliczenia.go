// Odpowiedzialność pliku: wartości wyliczeniowe obszaru dostępów.
//
// Kolumny `punkt_dostepu.rodzaj`, `punkt_dostepu.tryb_domyslny`,
// `punkt_dostepu.stan` i `nadanie_dostepu.tryb` niosą wartość kontraktu wprost.
// Kontrakt nie deklaruje przy wyliczeniach AccessPointKind, AccessPointStatus
// ani AccessMode pola `baza` — inaczej niż przy AccountKind, które ma słowniki
// `WartosciBazyAccountKind` / `WartosciKontraktuAccountKind`. Własny przekład
// w warstwie trwałości byłby drugim źródłem przekładu obok
// `shared/contract.json`. Poniższe zbiory są więc sprawdzeniem przynależności,
// nie przekładem: pilnują, żeby do kolumny nie trafiła wartość spoza kontraktu.
package dane

import (
	"fmt"

	"danacoconsole/shared"
)

var (
	// trybyDostepu — komplet wartości wyliczenia AccessMode.
	trybyDostepu = []shared.AccessMode{shared.AccessModeRead, shared.AccessModeWrite}

	// rodzajePunktuDostepu — komplet wartości wyliczenia AccessPointKind.
	rodzajePunktuDostepu = []shared.AccessPointKind{
		shared.AccessPointKindMcpBridge, shared.AccessPointKindLocalDirectory,
	}

	// stanyPunktuDostepu — komplet wartości wyliczenia AccessPointStatus.
	stanyPunktuDostepu = []shared.AccessPointStatus{
		shared.AccessPointStatusUnknown,
		shared.AccessPointStatusReachable,
		shared.AccessPointStatusUnreachable,
	}
)

// wartoscKontraktuNaBaze sprawdza przynależność wartości do wyliczenia kontraktu
// i zwraca ją jako wartość kolumny. Wartość pusta oznacza „nie ustawiono" i daje
// wartość domyślną, nie błąd.
func wartoscKontraktuNaBaze[T ~string](dozwolone []T, wartosc, domyslna T, pole string) (string, error) {
	if wartosc == "" {
		wartosc = domyslna
	}
	for _, znana := range dozwolone {
		if wartosc == znana {
			return string(wartosc), nil
		}
	}
	return "", fmt.Errorf("dane: wartość %q nie należy do wyliczenia kontraktu pola %s",
		string(wartosc), pole)
}

// wartoscKontraktuZBazy sprawdza wartość odczytaną z kolumny. Wartość nieznana
// jest sygnałem rozjazdu schematu z kontraktem, nie stanem normalnym.
func wartoscKontraktuZBazy[T ~string](dozwolone []T, kolumna, pole string) (T, error) {
	for _, znana := range dozwolone {
		if string(znana) == kolumna {
			return znana, nil
		}
	}
	var pusta T
	return pusta, fmt.Errorf("dane: wartość kolumny %q pola %s nie ma odpowiednika w kontrakcie",
		kolumna, pole)
}

func trybDostepuNaBaze(tryb shared.AccessMode, pole string) (string, error) {
	return wartoscKontraktuNaBaze(trybyDostepu, tryb, shared.AccessModeRead, pole)
}

func trybDostepuZBazy(kolumna, pole string) (shared.AccessMode, error) {
	return wartoscKontraktuZBazy(trybyDostepu, kolumna, pole)
}

func rodzajPunktuNaBaze(rodzaj shared.AccessPointKind) (string, error) {
	return wartoscKontraktuNaBaze(rodzajePunktuDostepu, rodzaj,
		shared.AccessPointKindMcpBridge, "punkt_dostepu.rodzaj")
}

func rodzajPunktuZBazy(kolumna string) (shared.AccessPointKind, error) {
	return wartoscKontraktuZBazy(rodzajePunktuDostepu, kolumna, "punkt_dostepu.rodzaj")
}

func stanPunktuNaBaze(stan shared.AccessPointStatus) (string, error) {
	return wartoscKontraktuNaBaze(stanyPunktuDostepu, stan,
		shared.AccessPointStatusUnknown, "punkt_dostepu.stan")
}

func stanPunktuZBazy(kolumna string) (shared.AccessPointStatus, error) {
	return wartoscKontraktuZBazy(stanyPunktuDostepu, kolumna, "punkt_dostepu.stan")
}
