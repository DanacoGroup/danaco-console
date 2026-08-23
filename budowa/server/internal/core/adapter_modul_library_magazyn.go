// Moduł Library — magazyn treści: miejsce na dysku, w którym lądują bajty
// każdego pliku wgranego do repozytorium wiedzy, przysłane w żądaniu
// (`contentBase64`) albo wciągnięte spod wskazanej ścieżki (`sourcePath`,
// `ZapiszZePliku`). Obie drogi kończą się blobem w magazynie, bo treść, którą
// ktoś z zewnątrz może nadpisać, nie jest treścią wersji. Dzięki temu czytelnik
// podglądu (`adapter_modul_library_podglad.go`) nie musi wiedzieć, którą drogą
// plik przyszedł: zawsze czyta ścieżkę z dysku.
//
// Magazyn leży w katalogu danych rdzenia — tym samym, w którym leży baza
// (`konfiguracja.Konfiguracja.SciezkaBazy`) i sejf poświadczeń
// (`dane/sejf_poswiadczen.go`) — w podkatalogu `biblioteka/tresc`. Nie leży
// w katalogu roboczym sesji, bo przeżywa restart rdzenia tak samo jak wiersz
// w bazie, który go wskazuje.
//
// Nazwą pliku jest suma kontrolna. Adapter i tak liczy sha256 każdej przysłanej
// treści (`trescZadania`), a nazwa z tej sumy niesie trzy własności naraz: dwa
// wgrania tej samej treści dzielą jeden blob zamiast dwóch kopii; ponowny zapis
// tej samej treści jest bezczynnością, a nie nadpisaniem; nazwa nie zależy od
// nazwy pliku z żądania, więc nie da się nią wyjść z katalogu magazynu. Pierwsze
// dwa znaki sumy tworzą podkatalog, żeby jeden katalog nie urósł do dziesiątek
// tysięcy wpisów.
//
// Zapis jest niepodzielny. Plik tymczasowy powstaje w katalogu docelowym (ten
// sam nośnik, więc `os.Rename` jest przemianowaniem, nie kopiowaniem między
// dyskami) i dopiero przemianowanie czyni treść widoczną pod odwołaniem. Awaria
// w połowie zostawia plik tymczasowy, nigdy pliku obciętego, który podgląd
// pokazałby jako pełną treść. Niepowodzenie zapisu jest odmową komendy
// u wołającego (`trescWgrania`, `adapter_modul_library_tresc.go`), nie pustym
// odwołaniem podanym jako powodzenie.
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
	// podkatalogBiblioteki oddziela magazyn modułu od reszty katalogu danych.
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

// ErrZrodloTresciNieczytelne oddziela pomyłkę wołającego od awarii magazynu.
// Ścieżka źródłowa przychodzi z żądania (`sourcePath`), więc literówka, plik
// usunięty, brak praw i wskazanie katalogu są brakami po stronie wołającego —
// żądanie nie uda się przy żadnym ponowieniu. Bez tego rozróżnienia wszystko
// wraca jako `internal_error` z `retryable:true`, a klient z pętlą ponowień
// powtarza żądanie bez końca.
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

// Zapisz utrwala bajty treści i zwraca ścieżkę, która staje się odwołaniem do
// treści w bazie. Treść już leżąca pod tą sumą nie jest zapisywana po raz drugi:
// nazwa jest sumą jej zawartości, więc plik o tej nazwie ma tę zawartość.
//
// Zapis bez sumy kontrolnej jest odmówiony, a nie zgadywany — bez nazwy nie ma
// gdzie odłożyć bajtów, a nazwa wymyślona (losowa) rozjechałaby się z sumą, którą
// ten sam wgrany plik niesie do bazy.
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

	// Plik tymczasowy leży w katalogu docelowym — ten sam nośnik, więc
	// przemianowanie jest niepodzielne. Plik tymczasowy w katalogu systemowym
	// dałby przemianowanie między wolumenami, czyli kopiowanie, czyli okno,
	// w którym pod odwołaniem leży treść obcięta.
	tymczasowy, err := os.CreateTemp(katalogBloku, "tresc-*.czesciowa")
	if err != nil {
		return "", fmt.Errorf("magazyn treści biblioteki: plik tymczasowy w %s: %w", katalogBloku, err)
	}
	nazwaTymczasowa := tymczasowy.Name()
	if _, err := tymczasowy.Write(bajty); err != nil {
		zamknijIUsun(tymczasowy, nazwaTymczasowa)
		return "", fmt.Errorf("magazyn treści biblioteki: zapis %s: %w", nazwaTymczasowa, err)
	}
	// Zrzut na nośnik przed przemianowaniem (`domknijTymczasowy`): inaczej po
	// utracie zasilania odwołanie w bazie wskazywałoby plik pusty, a nie żaden.
	if err := domknijTymczasowy(tymczasowy, nazwaTymczasowa); err != nil {
		return "", err
	}
	if err := os.Rename(nazwaTymczasowa, docelowa); err != nil {
		_ = os.Remove(nazwaTymczasowa)
		return "", fmt.Errorf("magazyn treści biblioteki: przemianowanie na %s: %w", docelowa, err)
	}
	return docelowa, nil
}

// ZapiszZePliku wciąga do magazynu treść pliku leżącego już na dysku (żądanie
// z `sourcePath`) i oddaje odwołanie, rozmiar oraz sumę kontrolną tego, co
// zostało wciągnięte.
//
// Kopia powstaje mimo tego, że plik gdzieś już leży: odwołanie do cudzego pliku
// wskazuje treść żywą, a historia wersji wymaga treści zamrożonej. Bez kopii plik
// nadpisany na dysku po wgraniu zmieniałby treść swojej utrwalonej wersji —
// `library.file.preview` oddawałby treść nową, a `library.version.restore` nie
// miałby do czego wrócić. Ścieżka źródłowa jest więc źródłem bajtów, a nie
// miejscem ich składowania.
//
// Suma liczy się w locie, w trakcie przepisywania: bez drugiego przebiegu po
// pliku i bez wciągania całej treści do pamięci. Kopiowanie idzie strumieniem,
// więc plik o dowolnym rozmiarze przechodzi tak samo.
//
// Nazwą bloba jest suma, znana dopiero po przeczytaniu całości, więc plik
// tymczasowy powstaje w korzeniu magazynu i dopiero stamtąd wędruje pod swoją
// sumę (ten sam nośnik, więc przemianowanie zostaje niepodzielne).
func (m *magazynTresciBiblioteki) ZapiszZePliku(sciezka string) (string, int64, string, error) {
	if m == nil || m.katalog == "" {
		return "", 0, "", fmt.Errorf("magazyn treści biblioteki nie ma wskazanego katalogu")
	}
	zrodlo, err := os.Open(sciezka)
	if err != nil {
		return "", 0, "", fmt.Errorf("%w: nie można otworzyć %s: %w", ErrZrodloTresciNieczytelne, sciezka, err)
	}
	defer zrodlo.Close()

	// Wskazanie katalogu rozpoznawane jest przed kopiowaniem. Bez tego
	// sprawdzenia dochodzi ono aż do `io.Copy` i wraca błędem odczytu nie do
	// odróżnienia od usterki nośnika, czyli jako awaria rdzenia, choć jest
	// pomyłką wołającego.
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

// odwolanieMagazynu przekłada ścieżkę bloba na dysku na odwołanie, które wolno
// oddać na zewnątrz rdzenia — ścieżkę względną magazynu, liczoną od jego
// korzenia w katalogu danych, zawsze z ukośnikiem `/`.
//
//	/tmp/dane-operatora/biblioteka/tresc/c4/c414cd0e…  →  biblioteka/tresc/c4/c414cd0e…
//	C:\Users\Jan\dane\design\zasoby\c4\c414cd0e…       →  design/zasoby/c4/c414cd0e…
//
// Ścieżka bezwzględna w polach `preview.imageRef`, `asset.uri` i `file.path`
// wynosiłaby układ katalogów maszyny Operatora poza rdzeń — do klienta i do
// modelu — a klient stojący na innej maszynie i tak pod nią nie sięgnie.
//
// Postacią odwołania jest ścieżka względna magazynu, a nie klucz
// nieprzezroczysty, bo klucz wymaga komendy, którą się go rozwiązuje, a takiej
// kontrakt nie ma ani dla biblioteki, ani dla Designu. Ścieżka względna nie
// niesie ani jednego członu maszyny Operatora, jest ta sama na każdej maszynie
// i po przeniesieniu katalogu danych, a rdzeń rozwiązuje ją z powrotem do
// bajtów jednym `filepath.Join` z katalogiem danych.
//
// Ścieżka spoza magazynu oddaje pustkę: nie da się jej wyrazić względem
// magazynu, a pole kontraktu, będąc niewymaganym, wtedy nie wychodzi. Nazwa
// bloba jest sumą sha256 treści, którą kontrakt i tak niesie jawnie
// w `checksum`, więc odwołanie nie wynosi informacji, której odbiorca by nie
// miał.
func odwolanieMagazynu(sciezka, korzenWzgledny string) string {
	if sciezka == "" || korzenWzgledny == "" {
		return ""
	}
	// Szukamy korzenia magazynu wewnątrz ścieżki zamiast liczyć różnicę względem
	// katalogu magazynu: ta funkcja nie ma dostępu do złożonego magazynu (woła ją
	// `zasobKontraktu`, funkcja bez stanu, dzielona z arsenałem), a odcięcie po
	// znaczniku daje ten sam wynik bez wiązania nowej zależności.
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

// odwolanieTresciBiblioteki oddaje odwołanie do treści pliku biblioteki
// w postaci, którą wolno wypuścić z rdzenia. Pustka znaczy „nie mam odwołania,
// które da się podać" — patrz `odwolanieMagazynu`.
func odwolanieTresciBiblioteki(sciezka string) string {
	return odwolanieMagazynu(sciezka, korzenTresciBiblioteki)
}

// domknijTymczasowy zrzuca plik tymczasowy na nośnik, zamyka go i nadaje mu
// prawa treści. Wspólne dla obu dróg zapisu — zrzut przed przemianowaniem jest
// tym, co odróżnia odwołanie do treści od odwołania do pliku pustego po utracie
// zasilania, więc nie ma prawa istnieć w dwóch wersjach.
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
