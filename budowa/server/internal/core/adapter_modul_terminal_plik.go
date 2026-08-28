// Komenda `terminal.file.read` — odczyt treści pliku z katalogu roboczego
// karty. Odczyt pliku nie jest czynnością powłoki, tylko czynnością rdzenia,
// niezależną od programów dostępnych na maszynie.
package core

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/session"
	"danacoconsole/server/internal/zewnetrzne"
	"danacoconsole/shared"
)

// granicaOdczytuPliku jest górną granicą jednego odczytu. Ta sama wartość co
// granica wyjścia procesu (`terminal.output.read`): plik czyta ten sam model
// i mieści go w tym samym oknie kontekstu.
const granicaOdczytuPliku = 64 * 1024

// granicaCzasuOdczytuZdalnego ogranicza odczyt pliku z maszyny zdalnej. Odczyt
// idzie łączem, więc granica jest hojniejsza niż lokalna, ale nie nieskończona:
// łącze, które przestało odpowiadać, nie ma prawa trzymać żądania bez końca.
const granicaCzasuOdczytuZdalnego = 30 * time.Second

// narzedzieSSH jest programem powłoki zdalnej. Deklaracja stoi tutaj, bo tu
// leży pierwsza czynność, która woła `ssh` poza kartą; sonda startowa bierze ją
// z tego miejsca (`zaleznosci_zewnetrzne.go`), zamiast wypisywać nazwę drugi raz.
var narzedzieSSH = zewnetrzne.Narzedzie{
	Nazwa:   "OpenSSH",
	Program: "ssh",
	Pakiet:  "openssh-client",
}

// OdczytajPlik obsługuje `terminal.file.read`, wybierając odczyt lokalny
// albo zdalny wedle powłoki karty.
func (a *adapterTerminala) OdczytajPlik(ctx context.Context,
	z shared.TerminalFileReadRequest) (shared.TerminalFileReadResponse, error) {

	karta, err := a.kartaZadania(z.SessionId)
	if err != nil {
		return shared.TerminalFileReadResponse{}, err
	}
	wskazanie := strings.TrimSpace(z.Path)
	if wskazanie == "" {
		return shared.TerminalFileReadResponse{}, bladZadaniaTerminala(
			"odczyt pliku wymaga wskazania ścieżki")
	}
	okno, err := a.oknoWykonania(karta.oknoKod)
	if err != nil {
		return shared.TerminalFileReadResponse{}, err
	}
	granica := granicaOdczytu(z.MaxBytes)

	if karta.powloka == shared.TerminalShellSsh {
		return a.odczytZdalny(ctx, karta, okno, wskazanie, granica, z.Tail)
	}
	return a.odczytLokalny(karta, okno, wskazanie, granica, z.Tail)
}

// odczytLokalny czyta plik z dysku maszyny rdzenia przez pakiet `os`, bez
// uruchamiania programu powłoki.
func (a *adapterTerminala) odczytLokalny(karta *kartaTerminala, okno session.Okno,
	wskazanie string, granica int, ogon *int) (shared.TerminalFileReadResponse, error) {

	pelna, err := a.sciezkaWObszarze(karta, okno, wskazanie)
	if err != nil {
		return shared.TerminalFileReadResponse{}, err
	}
	opis, err := os.Stat(pelna)
	if err != nil {
		return shared.TerminalFileReadResponse{}, bladBrakuZasobuTerminala("plik " + wskazanie)
	}
	if opis.IsDir() {
		return shared.TerminalFileReadResponse{}, bladZadaniaTerminala(
			"ścieżka " + wskazanie + " wskazuje katalog, a nie plik")
	}
	tresc, err := os.ReadFile(pelna)
	if err != nil {
		return shared.TerminalFileReadResponse{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeInternalError, "moduł Terminal: nie można odczytać pliku "+
				wskazanie+": "+err.Error()))
	}
	rozmiar := opis.Size()
	odpowiedz := zlozOdczytPliku(string(tresc), pelna, granica, ogon)
	odpowiedz.SizeBytes = &rozmiar
	return odpowiedz, nil
}

// odczytZdalny czyta plik na maszynie karty zdalnej.
//
// Ścieżka idzie POJEDYNCZYM argumentem programu `ssh`, a nie sklejona w wiersz
// powłoki, i poprzedza ją `--`: nazwa pliku zaczynająca się od myślnika ma być
// nazwą pliku, a nie przełącznikiem.
func (a *adapterTerminala) odczytZdalny(ctx context.Context, karta *kartaTerminala,
	okno session.Okno, wskazanie string, granica int,
	ogon *int) (shared.TerminalFileReadResponse, error) {

	if a.uruchamiacz == nil {
		return shared.TerminalFileReadResponse{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeInternalError, "moduł Terminal: rdzeń nie ma uruchamiacza procesów"))
	}
	zdalne, err := argumentyPowlokiZdalnej(karta)
	if err != nil {
		return shared.TerminalFileReadResponse{}, bladZadaniaTerminala(err.Error())
	}
	sciezka := wskazanie
	if !strings.HasPrefix(sciezka, "/") && strings.TrimSpace(karta.katalog) != "" {
		sciezka = strings.TrimRight(karta.katalog, "/") + "/" + sciezka
	}
	argumenty := append([]string{"-T"}, zdalne...)
	argumenty = append(argumenty, "cat", "--", sciezka)

	wynik, err := zewnetrzne.Wolaj(ctx, a.uruchamiacz, okno, a.zasadyOkna(okno),
		a.obszarOkna(okno), narzedzieSSH, argumenty, karta.katalog, granicaCzasuOdczytuZdalnego)
	if err != nil {
		var brak *zewnetrzne.BrakNarzedzia
		if errors.As(err, &brak) {
			return shared.TerminalFileReadResponse{}, protocol.JakoError(protocol.NowyBlad(
				shared.ErrorCodeChannelUnavailable, "moduł Terminal: "+brak.Error()))
		}
		if errors.Is(err, session.ErrIzolacja) {
			return shared.TerminalFileReadResponse{}, protocol.JakoError(protocol.NowyBlad(
				shared.ErrorCodePermissionDenied, "moduł Terminal: "+err.Error()))
		}
		return shared.TerminalFileReadResponse{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeNotFound, "moduł Terminal: nie można odczytać pliku "+wskazanie+
				" na maszynie karty: "+err.Error()))
	}
	return zlozOdczytPliku(string(wynik.Wyjscie), sciezka, granica, ogon), nil
}

// sciezkaWObszarze sprowadza wskazanie do ścieżki bezwzględnej i egzekwuje na
// niej punkt izolacji „pliki” okna.
func (a *adapterTerminala) sciezkaWObszarze(karta *kartaTerminala, okno session.Okno,
	wskazanie string) (string, error) {

	pelna := filepath.FromSlash(wskazanie)
	if !filepath.IsAbs(pelna) {
		korzen := strings.TrimSpace(karta.katalog)
		if korzen == "" {
			korzen = a.obszarOkna(okno).KatalogRoboczy
		}
		if korzen == "" {
			return "", bladZadaniaTerminala("karta " + karta.kod +
				" nie ma katalogu roboczego, więc ścieżka względna nie ma od czego się liczyć")
		}
		pelna = filepath.Join(korzen, pelna)
	}
	pelna = filepath.Clean(pelna)

	zasady := a.zasadyOkna(okno)
	if !zasady.Pliki {
		return pelna, nil
	}
	obszar := a.obszarOkna(okno).KatalogRoboczy
	if obszar == "" || !session.SciezkaWewnatrz(obszar, pelna) {
		return "", protocol.JakoError(protocol.NowyBlad(shared.ErrorCodePermissionDenied,
			"moduł Terminal: okno ma włączony punkt izolacji „pliki”, a ścieżka "+wskazanie+
				" leży poza obszarem okna"))
	}
	return pelna, nil
}

// zasadyOkna składa zasady izolacji obowiązujące w oknie. Rdzeń bez
// rozstrzygacza pracuje na zasadach pustych — stan wyjściowy platformy.
func (a *adapterTerminala) zasadyOkna(okno session.Okno) session.Zasady {
	if a.rozstrzygacz == nil {
		return session.Zasady{}
	}
	return ZasadyIzolacji(a.rozstrzygacz, konfig.Kontekst{Okno: okno.Id})
}

// zlozOdczytPliku przycina treść ogonem i granicą rozmiaru, w tej
// kolejności: ogon idzie pierwszy, bo jest wskazaniem Operatora, a granica
// rozmiaru zabezpieczeniem rdzenia.
func zlozOdczytPliku(tresc, sciezka string, granica int, ogon *int) shared.TerminalFileReadResponse {
	odciete := 0
	if ogon != nil && *ogon >= 0 {
		wiersze := strings.Split(tresc, "\n")
		if len(wiersze) > *ogon {
			pominiete := strings.Join(wiersze[:len(wiersze)-*ogon], "\n")
			odciete += len(pominiete) + 1
			wiersze = wiersze[len(wiersze)-*ogon:]
		}
		tresc = strings.Join(wiersze, "\n")
	}
	if granica > 0 && len(tresc) > granica {
		odciete += len(tresc) - granica
		tresc = tresc[len(tresc)-granica:]
	}
	odpowiedz := shared.TerminalFileReadResponse{
		Content:   tresc,
		Path:      filepath.ToSlash(sciezka),
		Truncated: odciete > 0,
	}
	if odciete > 0 {
		ile := int64(odciete)
		odpowiedz.TruncatedBytes = &ile
	}
	return odpowiedz
}

// granicaOdczytu czyta granicę rozmiaru z żądania. Brak wskazania i wartość
// większa od granicy rdzenia schodzą na granicę rdzenia.
func granicaOdczytu(bajty *int) int {
	if bajty == nil || *bajty <= 0 || *bajty > granicaOdczytuPliku {
		return granicaOdczytuPliku
	}
	return *bajty
}
