package konfiguracja

import "path/filepath"

const (
	// katalogKlientaZrodlo to katalog projektu interfejsu w drzewie budowy.
	katalogKlientaZrodlo = "klient"
	// katalogKlientaWydanie to podkatalog, do którego Vite składa pakiet.
	katalogKlientaWydanie = "dist"
)

// KatalogKlientaDomyslny wskazuje pakiet interfejsu obok miejsca uruchomienia
// rdzenia (klient/dist). Ścieżka jest względna świadomie: rdzeń zmienia
// umiejscowienie w czasie, więc katalog pakietu podąża za katalogiem
// roboczym procesu, dopóki Operator nie wskaże innego.
//
// Nazwa jest nazwą katalogu projektu interfejsu w drzewie budowy
// (`budowa/klient`), więc rdzeń uruchomiony z korzenia tego drzewa znajduje
// pakiet bez przełącznika. Ta sama nazwa wiąże wdrożenie: pakowanie wydania
// kładzie pakiet obok binarium pod nazwą stąd, a nie pod własną — dwie nazwy
// jednego katalogu wracają odmową pakietu przy pierwszym uruchomieniu.
//
// Katalog nieistniejący nie wstrzymuje startu — gniazdo pracuje bez plików
// statycznych.
func KatalogKlientaDomyslny() string {
	return filepath.Join(katalogKlientaZrodlo, katalogKlientaWydanie)
}
