package konfiguracja

import (
	"fmt"
	"strings"
)

// Rola wyznacza, którą część rdzenia uruchamia proces.
// Ta sama binarka pracuje lokalnie w trakcie budowy i na serwerze po przeniesieniu.
type Rola string

const (
	// RolaHub oznacza rdzeń serwerowy — kontrakt, magazyn, kolejki.
	RolaHub Rola = "hub"
	// RolaAgent oznacza agenta lokalnego urządzenia wykonującego pracę na plikach.
	RolaAgent Rola = "agent"
	// RolaWszystko oznacza hub i agenta w jednym procesie.
	RolaWszystko Rola = "all"
)

// Role zwraca dopuszczalne wartości przełącznika roli w kolejności prezentacji.
func Role() []Rola {
	return []Rola{RolaHub, RolaAgent, RolaWszystko}
}

// NazwyRol zwraca dopuszczalne wartości roli jako teksty.
func NazwyRol() []string {
	nazwy := make([]string, 0, len(Role()))
	for _, rola := range Role() {
		nazwy = append(nazwy, string(rola))
	}
	return nazwy
}

// RolaZTekstu zamienia tekst na rolę i odrzuca wartość spoza katalogu ról.
func RolaZTekstu(tekst string) (Rola, error) {
	for _, rola := range Role() {
		if string(rola) == tekst {
			return rola, nil
		}
	}
	return "", fmt.Errorf("nieznana rola %q; dopuszczalne: %s", tekst, strings.Join(NazwyRol(), "|"))
}
