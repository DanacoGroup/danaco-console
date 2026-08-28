package konfiguracja

import (
	"fmt"
	"strings"
)

// Rola wyznacza, którą część rdzenia uruchamia proces.
// Ta sama binarka pracuje lokalnie w trakcie budowy i na serwerze po przeniesieniu.
type Rola string

const (
	// RolaHub oznacza rdzeń serwerowy prowadzący kontrakt, magazyn danych
	// oraz kolejki dla podłączonych agentów lokalnych.
	RolaHub Rola = "hub"
	// RolaAgent oznacza agenta lokalnego urządzenia wykonującego pracę na
	// plikach w imieniu podłączonego huba serwerowego.
	RolaAgent Rola = "agent"
	// RolaWszystko oznacza hub i agenta uruchomione razem w jednym procesie,
	// bez podziału na dwie osobne role.
	RolaWszystko Rola = "all"
)

// Role zwraca dopuszczalne wartości przełącznika roli w kolejności, w jakiej
// są prezentowane operatorowi.
func Role() []Rola {
	return []Rola{RolaHub, RolaAgent, RolaWszystko}
}

// NazwyRol zwraca dopuszczalne wartości roli jako teksty, w tej samej
// kolejności, w jakiej zwraca je funkcja Role.
func NazwyRol() []string {
	nazwy := make([]string, 0, len(Role()))
	for _, rola := range Role() {
		nazwy = append(nazwy, string(rola))
	}
	return nazwy
}

// RolaZTekstu zamienia tekst wskazany przełącznikiem na rolę i odrzuca
// wartość spoza katalogu ról zwróconego przez Role.
func RolaZTekstu(tekst string) (Rola, error) {
	for _, rola := range Role() {
		if string(rola) == tekst {
			return rola, nil
		}
	}
	return "", fmt.Errorf("nieznana rola %q; dopuszczalne: %s", tekst, strings.Join(NazwyRol(), "|"))
}
