// Odpowiedzialność pliku: zaplecze wspólne wszystkim rodzinom modułu Browser
// dołożonym ponad migawkę i notatkę — magazyn bajtów sesji przeglądania, silnik
// przeglądarki, przekłady czasu oraz odmowy nazywające brak.
//
// Magazyn jest ten sam co w Library i Design co do mechaniki (blob pod sumą
// sha256, zapis niepodzielny), a inny co do miejsca: bajty modułu Browser leżą
// w `<dane>/przegladarka/tresc`. Wspólny katalog z biblioteką mieszałby materiał
// trwały (dokument wniesiony do repozytorium wiedzy) z materiałem sesji (zrzut
// strony, archiwum, rejestr sieciowy) — a te dwa mają różny cykl życia.
//
// Silnik przeglądarki wchodzi tu jednym polem, nie jednym na rodzinę: Chromium
// startuje ten sam dla zrzutu, dla drzewa DOM i dla konsoli, więc drugie pole
// byłoby drugą prawdą o tym, czym rdzeń renderuje stronę.
package core

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/session"
	"danacoconsole/server/internal/zewnetrzne"
	"danacoconsole/shared"
)

const (
	// podkatalogPrzegladarki oddziela materiał sesji przeglądania od reszty
	// katalogu danych rdzenia.
	podkatalogPrzegladarki = "przegladarka"
	// podkatalogTresciPrzegladarki mieści same bajty — zrzuty, archiwa,
	// odniesienia monitorów i rejestry sieciowe.
	podkatalogTresciPrzegladarki = "tresc"
	// korzenTresciPrzegladarki jest tą samą parą podkatalogów liczoną od katalogu
	// danych; postać odwołania i miejsce zapisu nie mają jak się rozjechać.
	korzenTresciPrzegladarki = podkatalogPrzegladarki + "/" + podkatalogTresciPrzegladarki
)

// Przedrostki identyfikatorów bytów dołożonych do modułu. Rdzeń nadaje je sam,
// bo kontrakt nie niesie ich w żądaniach (wzór: źródło i notatka).
const (
	przedrostekKartyPrzegladania  = "tab-"
	przedrostekGrupyKart          = "tabgrp-"
	przedrostekPrzestrzeni        = "wksp-"
	przedrostekMonitora           = "mon-"
	przedrostekKanaluPrzegladania = "feed-"
	przedrostekWpisuKanalu        = "feeditem-"
	przedrostekPozycjiCzytania    = "read-"
	przedrostekZakladki           = "bmk-"
	przedrostekWytworu            = "bart-"
	przedrostekZrzutu             = "shot-"
	przedrostekPobrania           = "dl-"
	przedrostekMakra              = "macro-"
	przedrostekZestawuZrodel      = "srcgrp-"
	przedrostekWatkuNotatek       = "thread-"
)

// magazynTresciPrzegladarki składa magazyn bajtów modułu nad katalogiem danych.
func magazynTresciPrzegladarki(katalogDanych string) *magazynTresciBiblioteki {
	if strings.TrimSpace(katalogDanych) == "" {
		return nil
	}
	return &magazynTresciBiblioteki{
		katalog: filepath.Join(katalogDanych, podkatalogPrzegladarki, podkatalogTresciPrzegladarki),
	}
}

// ZMagazynem przestawia adapter na wskazany katalog danych i wpina magazyn
// bajtów sesji przeglądania. Wołane przy montażu — bez niego rodziny wytwarzające
// materiał (zrzut, archiwum, rejestr) nie mają gdzie odłożyć bajtów i mówią to
// wprost, zamiast meldować powodzenie bez treści.
func (a *adapterPrzegladarki) ZMagazynem(katalogDanych string) *adapterPrzegladarki {
	a.katalogDanych = katalogDanych
	a.magazyn = magazynTresciPrzegladarki(katalogDanych)
	return a
}

// ZSilnikiem wpina silnik przeglądarki wraz z jego zasięgiem izolacji. Bez niego
// rodziny wymagające uruchomienia strony (zrzut, DOM, konsola, sieć, emulacja,
// przewinięcie) odmawiają zdaniem nazywającym brak.
func (a *adapterPrzegladarki) ZSilnikiem(uruchamiacz session.Uruchamiacz,
	rozstrzygacz *konfig.Rozstrzygacz, katalog *KatalogRoboczy) *adapterPrzegladarki {

	a.silnik = &silnikPrzegladarki{uruchamiacz: uruchamiacz, rozstrzygacz: rozstrzygacz, katalog: katalog}
	return a
}

// zapiszTresc odkłada bajty w magazynie modułu i oddaje odwołanie względne
// katalogu danych — dokładnie w tej postaci, w której wychodzi kontraktem.
// Odwołanie bezwzględne wynosiłoby układ katalogów maszyny do klienta.
func (a *adapterPrzegladarki) zapiszTresc(bajty []byte) (string, error) {
	if a.magazyn == nil {
		return "", bladZapleczaPrzegladarki("rdzeń nie ma magazynu treści przeglądania — " +
			"nie ma gdzie odłożyć bajtów; naprawa: wskazać katalog danych przy składaniu rdzenia")
	}
	if len(bajty) == 0 {
		return "", bladWskazaniaPrzegladarki("zapis treści bez ani jednego bajtu")
	}
	sciezka, err := a.magazyn.Zapisz(bajty, sumaTresciPrzegladania(bajty))
	if err != nil {
		return "", bladPrzegladarki(err)
	}
	odwolanie := odwolanieMagazynu(sciezka, korzenTresciPrzegladarki)
	if odwolanie == "" {
		return "", bladZapleczaPrzegladarki("magazyn oddał ścieżkę spoza katalogu treści przeglądania")
	}
	return odwolanie, nil
}

// odczytajTresc czyta bajty leżące pod odwołaniem magazynu.
func (a *adapterPrzegladarki) odczytajTresc(odwolanie string) ([]byte, error) {
	if a.magazyn == nil {
		return nil, bladZapleczaPrzegladarki("rdzeń nie ma magazynu treści przeglądania")
	}
	sciezka := filepath.Join(a.katalogDanych, filepath.FromSlash(odwolanie))
	bajty, err := odczytajPlikMagazynu(sciezka)
	if err != nil {
		return nil, bladPrzegladarki(err)
	}
	return bajty, nil
}

// upewnijSieOSilniku odmawia wcześnie, gdy nie ma czym uruchomić strony.
func (a *adapterPrzegladarki) upewnijSieOSilniku(komenda string) error {
	if a.silnik == nil {
		return bladZapleczaPrzegladarki("rdzeń nie ma silnika przeglądarki — komenda " + komenda +
			" wymaga uruchomienia strony; naprawa: podpiąć warstwę kanału przy składaniu rdzenia")
	}
	if err := a.silnik.dostepny(); err != nil {
		return bladSilnikaPrzegladarki(komenda, err)
	}
	return nil
}

// otworzStrone otwiera stronę silnikiem i oddaje sesję wraz z jej wynikiem.
// Wołający zamyka sesję — inaczej proces przeglądarki zostaje na maszynie.
func (a *adapterPrzegladarki) otworzStrone(ctx context.Context, komenda string,
	nastawy nastawyStrony) (*sesjaStrony, wynikOtwarcia, error) {

	if err := a.upewnijSieOSilniku(komenda); err != nil {
		return nil, wynikOtwarcia{}, err
	}
	sesja, wynik, err := a.silnik.otworz(ctx, nastawy)
	if err != nil {
		return nil, wynikOtwarcia{}, bladSilnikaPrzegladarki(komenda, err)
	}
	return sesja, wynik, nil
}

// adresOstatniejStrony oddaje adres, pod którym okno stoi. Komendy inspekcyjne
// kontrakt opisuje bez pola adresu — pytają o „bieżącą stronę okna", a bieżącą
// stroną okna jest ostatnia jego migawka. Okno bez migawki dostaje odmowę
// `not_found`, nie pustą inspekcję udającą stronę bez treści.
func (a *adapterPrzegladarki) adresOstatniejStrony(ctx context.Context, okno, komenda string) (dane.MigawkaStrony, error) {
	if strings.TrimSpace(okno) == "" {
		return dane.MigawkaStrony{}, bladWskazaniaPrzegladarki("komenda " + komenda + " bez okna")
	}
	wiersz, err := a.repozytorium.OstatniaMigawka(ctx, okno)
	if err != nil {
		return dane.MigawkaStrony{}, bladBrakuMigawki(okno, err)
	}
	return wiersz, nil
}

// znacznikChwili przekłada milisekundy epoki kontraktu na znacznik czasu bazy.
// Odwrotność `chwilaBazy` — jedno miejsce przekładu w obie strony dla modułu.
func znacznikChwili(chwila int64) *string {
	if chwila <= 0 {
		return nil
	}
	znak := time.UnixMilli(chwila).UTC().Format(formatZnacznikaBazy)
	return &znak
}

// chwilaZeZnacznika przekłada znacznik bazy na wskaźnik chwili kontraktu —
// pola opcjonalne (`lastCheckedAt`, `finishedAt`) niosą brak jako brak.
func chwilaZeZnacznika(znacznik *string) *int64 {
	if znacznik == nil || *znacznik == "" {
		return nil
	}
	chwila := chwilaBazy(*znacznik)
	if chwila == 0 {
		return nil
	}
	return &chwila
}

// terazWBazie oddaje bieżącą chwilę w zapisie kolumn czasu.
func terazWBazie() *string {
	znak := time.Now().UTC().Format(formatZnacznikaBazy)
	return &znak
}

// wykazJson zapisuje wykaz tekstów kolumną JSON. Pusty wykaz zapisuje się jako
// brak, nie jako `[]`: kontrakt czyta brak wykazu i wykaz pusty tak samo, a NULL
// w kolumnie mówi wprost „nic tu nie ustawiono".
func wykazJson(wykaz []string) *string {
	if len(wykaz) == 0 {
		return nil
	}
	bajty, err := json.Marshal(wykaz)
	if err != nil {
		return nil
	}
	zapis := string(bajty)
	return &zapis
}

// wykazZJson odczytuje wykaz tekstów z kolumny JSON. Kolumna nieczytelna daje
// wykaz pusty — jedna uszkodzona kolumna nie ma prawa wywrócić całego odczytu.
func wykazZJson(zapis *string) []string {
	if zapis == nil || *zapis == "" {
		return nil
	}
	var wykaz []string
	if err := json.Unmarshal([]byte(*zapis), &wykaz); err != nil {
		return nil
	}
	return wykaz
}

// wartoscTekstuLubPusta oddaje treść wskaźnika albo pusty napis.
func wartoscTekstuLubPusta(wskazanie *string) string {
	if wskazanie == nil {
		return ""
	}
	return strings.TrimSpace(*wskazanie)
}

// bladZapleczaPrzegladarki nazywa brak po stronie montażu rdzenia — nie winę
// Operatora i nie usterkę przemijającą.
func bladZapleczaPrzegladarki(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"moduł Browser: "+powod))
}

// bladNieznanegoBytu nazywa wskazanie, pod którym w module nic nie leży.
func bladNieznanegoBytu(co, wskazanie string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
		"moduł Browser: nie ma "+co+" o wskazaniu: "+wskazanie))
}

// bladSilnikaPrzegladarki przekłada niepowodzenie silnika na odmowę kontraktu.
//
// Brak programu i usterka uruchomienia idą tym samym kodem `channel_unavailable`
// — tak samo jak w pozostałych rodzinach arsenału (obraz, dokument, media).
// Kontrakt nie ma kodu „brakuje programu"; `validation_failed` kłamałby o winie
// żądania, a `internal_error` o usterce rdzenia. Treść odmowy nazywa różnicę
// wprost: przy braku niesie nazwę programu i podpowiedź instalacyjną z samego
// `zewnetrzne.BrakNarzedzia`.
func bladSilnikaPrzegladarki(komenda string, err error) error {
	if err == nil {
		return nil
	}
	var brak *zewnetrzne.BrakNarzedzia
	if errors.As(err, &brak) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
			"moduł Browser: komenda "+komenda+" nie ma czym uruchomić strony: "+err.Error()))
	}
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeChannelUnavailable,
		"moduł Browser: komenda "+komenda+" nie doszła do skutku: "+err.Error()))
}

// sumaTresciPrzegladania liczy nazwę bloba w magazynie — tę samą sumę sha256,
// którą liczy wgranie do biblioteki i wniesienie zasobu Designu.
func sumaTresciPrzegladania(bajty []byte) string {
	suma := sha256.Sum256(bajty)
	return hex.EncodeToString(suma[:])
}

// odczytajPlikMagazynu czyta bajty spod ścieżki magazynu. Plik pusty jest tu
// osobnym niepowodzeniem, nie odmianą braku: wiersz wskazujący plik zerowej
// długości przechodzi każde sprawdzenie istnienia i nie ma czego pokazać.
func odczytajPlikMagazynu(sciezka string) ([]byte, error) {
	bajty, err := os.ReadFile(sciezka)
	if err != nil {
		return nil, fmt.Errorf("treść przeglądania nie leży pod odwołaniem %s: %w", sciezka, err)
	}
	if len(bajty) == 0 {
		return nil, fmt.Errorf("treść przeglądania pod odwołaniem %s ma zerową długość", sciezka)
	}
	return bajty, nil
}

// isBrakWiersza mówi, czy niepowodzenie odczytu jest brakiem wiersza, a nie
// usterką bazy. Jedno miejsce tego rozróżnienia dla całego modułu — rozsiane
// po rodzinach `errors.Is` rozjechałoby się przy pierwszej zmianie warstwy
// danych.
func isBrakWiersza(err error) bool {
	return errors.Is(err, dane.ErrBrakWiersza)
}

// jakoLiteral, jakoLiczba i jakoLogiczna wstawiają wartości rdzenia do wyrażeń
// wykonywanych na stronie.
//
// Wstawienie idzie przez zapis literału (`strconv.Quote`), nie przez sklejenie
// napisów: selektor przychodzi z żądania, a selektor z apostrofem albo
// z domknięciem nawiasu przerwałby wyrażenie i wykonał na stronie coś innego,
// niż rdzeń napisał. To jest ta sama zasada, co parametry zapytania zamiast
// sklejanego SQL-a.
func jakoLiteral(wartosc string) string {
	return strconv.Quote(wartosc)
}

func jakoLiczba(wartosc int) string {
	return strconv.Itoa(wartosc)
}

func jakoLogiczna(wartosc bool) string {
	if wartosc {
		return "true"
	}
	return "false"
}
