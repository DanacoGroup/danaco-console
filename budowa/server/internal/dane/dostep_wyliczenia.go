// Odpowiedzialność pliku: wartości wyliczeniowe obszaru dostępów — zbiory sprawdzające przynależność kolumn
// punktu dostępu i nadania do wartości kontraktu.
package dane

import (
	"fmt"

	"danacoconsole/shared"
)

var (
	// trybyDostepu zawiera komplet wartości wyliczenia AccessMode, dopuszczalnych w kolumnach trybu dostępu i nadania.
	trybyDostepu = []shared.AccessMode{shared.AccessModeRead, shared.AccessModeWrite}

	// rodzajePunktuDostepu zawiera komplet wartości wyliczenia AccessPointKind, dopuszczalnych w kolumnie rodzaju punktu.
	rodzajePunktuDostepu = []shared.AccessPointKind{
		shared.AccessPointKindMcpBridge, shared.AccessPointKindLocalDirectory,
	}

	// stanyPunktuDostepu zawiera komplet wartości wyliczenia AccessPointStatus, dopuszczalnych w kolumnie stanu punktu.
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
