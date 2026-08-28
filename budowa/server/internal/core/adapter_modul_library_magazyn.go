// Modul Library — magazyn tresci: katalog na dysku, pod ktorym rdzen zapisuje
// bajty kazdego pliku wgranego do repozytorium wiedzy, niezaleznie od tego, czy
// przyszly w zadaniu, czy zostaly wciagniete ze wskazanej sciezki zrodlowej.
package core

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	// podkatalogBiblioteki oddziela magazyn modulu od reszty katalogu danych
	// rdzenia, w ktorym leza takze baza i sejf poswiadczen.
	podkatalogBiblioteki = "biblioteka"
	// podkatalogTresciBiblioteki mieści same bajty treści. Osobny poziom, bo
	// moduł może kiedyś potrzebować w swoim katalogu czegoś jeszcze (miniatury,
	// wyciągi tekstu) i nie ma mieszać tego z blobami.
	podkatalogTresciBiblioteki = "tresc"
	// prawaKataloguTresci i prawaPlikuTresci: treść biblioteki należy do
	// użytkownika, który uruchomił rdzeń — tak samo jak sejf poświadczeń.
	prawaKataloguTresci = 0o700
	prawaPlikuTresci    = 0o600
)

// ErrZrodloTresciNieczytelne oddziela pomylke wolajacego od awarii magazynu:
// literowka, plik usuniety, brak praw i wskazanie katalogu sa brakami po stronie
// wolajacego, wiec zadanie nie uda sie przy zadnym ponowieniu.
var ErrZrodloTresciNieczytelne = errors.New("magazyn treści biblioteki")

// magazynTresciBiblioteki jest magazynem bajtów treści pod katalogiem danych.
// Trzyma wyłącznie korzeń — reszta ścieżki wynika z sumy kontrolnej treści.
type magazynTresciBiblioteki struct {
	katalog string
}

// nowyMagazynTresciBiblioteki składa magazyn nad katalogiem danych rdzenia.
// Katalog nie musi jeszcze istnieć — powstanie przy pierwszym zapisie, tak jak
// przy sejfie poświadczeń.
func nowyMagazynTresciBiblioteki(katalogDanych string) *magazynTresciBiblioteki {
	return &magazynTresciBiblioteki{
		katalog: filepath.Join(katalogDanych, podkatalogBiblioteki, podkatalogTresciBiblioteki),
	}
}

// Zapisz utrwala bajty tresci i zwraca sciezke, ktora staje sie odwolaniem do
// tresci w bazie; tresc juz lezaca pod dana suma nie jest zapisywana po raz drugi.
func (m *magazynTresciBiblioteki) Zapisz(bajty []byte, sumaKontrolna string) (string, error) {
	if m == nil || m.katalog == "" {
		return "", fmt.Errorf("magazyn treści biblioteki nie ma wskazanego katalogu")
	}
	if sumaKontrolna == "" {
		return "", fmt.Errorf("magazyn treści biblioteki: zapis bez sumy kontrolnej")
	}

	katalogBloku := filepath.Join(m.katalog, sumaKontrolna[:2])
	docelowa := filepath.Join(katalogBloku, sumaKontrolna)
	if stan, err := os.Stat(docelowa); err == nil && !stan.IsDir() {
		return docelowa, nil
	}
	if err := os.MkdirAll(katalogBloku, prawaKataloguTresci); err != nil {
		return "", fmt.Errorf("magazyn treści biblioteki: katalog %s: %w", katalogBloku, err)
	}

	// Plik tymczasowy powstaje w katalogu docelowym, zeby przemianowanie bylo niepodzielne.
	tymczasowy, err := os.CreateTemp(katalogBloku, "tresc-*.czesciowa")
	if err != nil {
		return "", fmt.Errorf("magazyn treści biblioteki: plik tymczasowy w %s: %w", katalogBloku, err)
	}
	nazwaTymczasowa := tymczasowy.Name()
	if _, err := tymczasowy.Write(bajty); err != nil {
		zamknijIUsun(tymczasowy, nazwaTymczasowa)
		return "", fmt.Errorf("magazyn treści biblioteki: zapis %s: %w", nazwaTymczasowa, err)
	}
	// Zrzut na nosnik przed przemianowaniem chroni przed odwolaniem do pliku pustego po utracie zasilania.
	if err := domknijTymczasowy(tymczasowy, nazwaTymczasowa); err != nil {
		return "", err
	}
	if err := os.Rename(nazwaTymczasowa, docelowa); err != nil {
		_ = os.Remove(nazwaTymczasowa)
		return "", fmt.Errorf("magazyn treści biblioteki: przemianowanie na %s: %w", docelowa, err)
	}
	return docelowa, nil
}

// ZapiszZePliku wciaga do magazynu tresc pliku leżącego juz na dysku i oddaje
// odwolanie, rozmiar oraz sume kontrolna tego, co zostalo wciagniete; kopia
// chroni utrwalona wersje przed zmiana pliku zrodlowego.
func (m *magazynTresciBiblioteki) ZapiszZePliku(sciezka string) (string, int64, string, error) {
	if m == nil || m.katalog == "" {
		return "", 0, "", fmt.Errorf("magazyn treści biblioteki nie ma wskazanego katalogu")
	}
	zrodlo, err := os.Open(sciezka)
	if err != nil {
		return "", 0, "", fmt.Errorf("%w: nie można otworzyć %s: %w", ErrZrodloTresciNieczytelne, sciezka, err)
	}
	defer zrodlo.Close()

	// Wskazanie katalogu jest rozpoznawane przed kopiowaniem, by odroznic pomylke od awarii nosnika.
	if stan, err := zrodlo.Stat(); err == nil && stan.IsDir() {
		return "", 0, "", fmt.Errorf("%w: %s jest katalogiem, a treść pliku bierze się z pliku",
			ErrZrodloTresciNieczytelne, sciezka)
	}

	if err := os.MkdirAll(m.katalog, prawaKataloguTresci); err != nil {
		return "", 0, "", fmt.Errorf("magazyn treści biblioteki: katalog %s: %w", m.katalog, err)
	}
	tymczasowy, err := os.CreateTemp(m.katalog, "wciagana-*.czesciowa")
	if err != nil {
		return "", 0, "", fmt.Errorf("magazyn treści biblioteki: plik tymczasowy w %s: %w", m.katalog, err)
	}
	nazwaTymczasowa := tymczasowy.Name()

	licznik := sha256.New()
	rozmiar, err := io.Copy(io.MultiWriter(tymczasowy, licznik), zrodlo)
	if err != nil {
		zamknijIUsun(tymczasowy, nazwaTymczasowa)
		return "", 0, "", fmt.Errorf("magazyn treści biblioteki: wciąganie %s: %w", sciezka, err)
	}
	if err := domknijTymczasowy(tymczasowy, nazwaTymczasowa); err != nil {
		return "", 0, "", err
	}

	suma := hex.EncodeToString(licznik.Sum(nil))
	docelowa, err := m.podSuma(nazwaTymczasowa, suma)
	if err != nil {
		return "", 0, "", err
	}
	return docelowa, rozmiar, suma, nil
}

// podSuma przenosi gotowy plik tymczasowy pod jego sumę kontrolną. Blob o tej
// nazwie już leżący ma z definicji tę samą zawartość, więc kopia jest zbędna —
// tymczasowy idzie do kosza, a odwołaniem zostaje blob zastany.
func (m *magazynTresciBiblioteki) podSuma(nazwaTymczasowa, suma string) (string, error) {
	katalogBloku := filepath.Join(m.katalog, suma[:2])
	docelowa := filepath.Join(katalogBloku, suma)
	if stan, err := os.Stat(docelowa); err == nil && !stan.IsDir() {
		_ = os.Remove(nazwaTymczasowa)
		return docelowa, nil
	}
	if err := os.MkdirAll(katalogBloku, prawaKataloguTresci); err != nil {
		_ = os.Remove(nazwaTymczasowa)
		return "", fmt.Errorf("magazyn treści biblioteki: katalog %s: %w", katalogBloku, err)
	}
	if err := os.Rename(nazwaTymczasowa, docelowa); err != nil {
		_ = os.Remove(nazwaTymczasowa)
		return "", fmt.Errorf("magazyn treści biblioteki: przemianowanie na %s: %w", docelowa, err)
	}
	return docelowa, nil
}

// odwolanieMagazynu przeklada sciezke bloba na dysku na odwolanie wolne od ukladu
// katalogow maszyny Operatora: sciezke wzgledna magazynu, liczona od jego
// korzenia w katalogu danych, zawsze z ukosnikiem.
func odwolanieMagazynu(sciezka, korzenWzgledny string) string {
	if sciezka == "" || korzenWzgledny == "" {
		return ""
	}
	// Korzen magazynu jest szukany wewnatrz sciezki, bo funkcja nie ma dostepu do zlozonego magazynu.
	zeSlashem := filepath.ToSlash(sciezka)
	znacznik := "/" + korzenWzgledny + "/"
	poczatek := strings.LastIndex(zeSlashem, znacznik)
	if poczatek < 0 {
		return ""
	}
	ogon := zeSlashem[poczatek+len(znacznik):]
	if ogon == "" {
		return ""
	}
	return korzenWzgledny + "/" + ogon
}

// korzenTresciBiblioteki jest korzeniem magazynu treści biblioteki liczonym od
// katalogu danych — ta sama para podkatalogów, którą składa
// `nowyMagazynTresciBiblioteki`, więc postać odwołania i miejsce zapisu nie mają
// jak się rozjechać.
const korzenTresciBiblioteki = podkatalogBiblioteki + "/" + podkatalogTresciBiblioteki

// odwolanieTresciBiblioteki oddaje odwolanie do tresci pliku biblioteki w postaci,
// ktora wolno wypuscic z rdzenia; pustka znaczy brak odwolania mozliwego do podania.
func odwolanieTresciBiblioteki(sciezka string) string {
	return odwolanieMagazynu(sciezka, korzenTresciBiblioteki)
}

// domknijTymczasowy zrzuca plik tymczasowy na nosnik, zamyka go i nadaje mu prawa
// tresci; jest wspolny dla obu drog zapisu, wiec nie ma prawa istniec w dwoch wersjach.
func domknijTymczasowy(plik *os.File, nazwa string) error {
	if err := plik.Sync(); err != nil {
		zamknijIUsun(plik, nazwa)
		return fmt.Errorf("magazyn treści biblioteki: zrzut %s: %w", nazwa, err)
	}
	if err := plik.Close(); err != nil {
		_ = os.Remove(nazwa)
		return fmt.Errorf("magazyn treści biblioteki: zamknięcie %s: %w", nazwa, err)
	}
	if err := os.Chmod(nazwa, prawaPlikuTresci); err != nil {
		_ = os.Remove(nazwa)
		return fmt.Errorf("magazyn treści biblioteki: prawa %s: %w", nazwa, err)
	}
	return nil
}

// zamknijIUsun sprząta plik tymczasowy po nieudanym zapisie. Niepowodzenie
// samego sprzątania nie ma komu nic powiedzieć — komenda i tak już odmawia.
func zamknijIUsun(plik *os.File, nazwa string) {
	_ = plik.Close()
	_ = os.Remove(nazwa)
}
