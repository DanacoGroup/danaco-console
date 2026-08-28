package konfiguracja

import "path/filepath"

const (
	// katalogKlientaZrodlo to katalog projektu interfejsu w drzewie budowy,
	// wskazywany względem katalogu roboczego procesu.
	katalogKlientaZrodlo = "klient"
	// katalogKlientaWydanie to podkatalog, do którego program Vite składa
	// gotowy pakiet interfejsu przy budowie.
	katalogKlientaWydanie = "dist"
)

// KatalogKlientaDomyslny wskazuje pakiet interfejsu obok miejsca
// uruchomienia rdzenia, w podkatalogu klient/dist. Ścieżka jest względna
// i podąża za katalogiem roboczym procesu. Katalog nieistniejący nie
// wstrzymuje startu.
func KatalogKlientaDomyslny() string {
	return filepath.Join(katalogKlientaZrodlo, katalogKlientaWydanie)
}
