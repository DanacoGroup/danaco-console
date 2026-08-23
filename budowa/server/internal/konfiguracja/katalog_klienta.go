package konfiguracja

import "path/filepath"

const (
	// katalogKlientaZrodlo to katalog projektu interfejsu w drzewie budowy.
	katalogKlientaZrodlo = "client"
	// katalogKlientaWydanie to podkatalog, do którego Vite składa pakiet.
	katalogKlientaWydanie = "dist"
)

// KatalogKlientaDomyslny wskazuje pakiet interfejsu obok miejsca uruchomienia
// rdzenia (client/dist). Ścieżka jest względna świadomie: rdzeń zmienia
// umiejscowienie w czasie, więc katalog pakietu podąża za katalogiem
// roboczym procesu, dopóki Operator nie wskaże innego.
//
// Katalog nieistniejący nie wstrzymuje startu — gniazdo pracuje bez plików
// statycznych.
func KatalogKlientaDomyslny() string {
	return filepath.Join(katalogKlientaZrodlo, katalogKlientaWydanie)
}
