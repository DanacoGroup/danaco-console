package transport

import (
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

// plikWejsciowy jest dokumentem, którym odpowiada się na ścieżki nienależące do
// żadnego pliku. Interfejs jest jedną stroną z własnym trasowaniem, więc adres
// widoku nie ma odpowiednika w pakiecie klienta.
const plikWejsciowy = "index.html"

// Funkcja uchwytStatyki serwuje pakiet interfejsu klienta z katalogu wskazanego jako parametr wywołania.
func uchwytStatyki(katalog string, dziennik *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !katalogIstnieje(katalog) {
			odmowaPakietu(w, katalog)
			return
		}
		zadany := sciezkaWKatalogu(katalog, r.URL.Path)
		if plikIstnieje(zadany) {
			http.ServeFile(w, r, zadany)
			return
		}
		indeks := filepath.Join(katalog, plikWejsciowy)
		if plikIstnieje(indeks) {
			http.ServeFile(w, r, indeks)
			return
		}
		dziennik.Printf("transport: brak pliku klienta %s", r.URL.Path)
		http.NotFound(w, r)
	})
}

// sciezkaWKatalogu zamienia ścieżkę żądania na ścieżkę w katalogu pakietu.
// Czyszczenie względem korzenia odcina wyjście poza katalog przez „..".
func sciezkaWKatalogu(katalog, sciezkaZadania string) string {
	oczyszczona := path.Clean("/" + sciezkaZadania)
	return filepath.Join(katalog, filepath.FromSlash(oczyszczona))
}

// Funkcja katalogIstnieje sprawdza, czy pakiet klienta jest fizycznie na miejscu w systemie plików dysku.
func katalogIstnieje(katalog string) bool {
	if katalog == "" {
		return false
	}
	info, err := os.Stat(katalog)
	return err == nil && info.IsDir()
}

// Funkcja plikIstnieje sprawdza, czy podana ścieżka wskazuje istniejący zwykły plik na tym dysku serwera.
func plikIstnieje(sciezka string) bool {
	info, err := os.Stat(sciezka)
	return err == nil && !info.IsDir()
}

// odmowaPakietu odpowiada, gdy pakietu klienta nie ma. Treść nazywa stan
// katalogu, bo brak plików jest stanem budowy, nie usterką rdzenia.
func odmowaPakietu(w http.ResponseWriter, katalog string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	_, _ = w.Write([]byte("Danaco Console: pakiet klienta niedostępny (" + opisKatalogu(katalog) + ").\n" +
		"Serwer pracuje; kanał WebSocket jest czynny.\n"))
}

// Funkcja opisKatalogu nazywa stan katalogu klienta na potrzeby dziennika oraz treści tej odpowiedzi HTTP.
func opisKatalogu(katalog string) string {
	if katalog == "" {
		return "katalog niewskazany"
	}
	if !katalogIstnieje(katalog) {
		return "brak katalogu " + katalog
	}
	return katalog
}
