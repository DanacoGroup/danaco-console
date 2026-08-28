// Odpowiedzialność pliku: moduł Browser — typ adaptera, konstruktor i dwie
// komendy migawki strony (`browser.navigate`, `browser.snapshot.get`). Źródła
// i notatki zebrane w toku przeglądania mają własny plik adaptera i własny
// port.
package core

import (
	"context"
	"errors"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// Przedrostki identyfikatorów bytów modułu. Migawka wychodzi kontraktem pod
// własnym kodem, nie pod kluczem wiersza (wzór: automations).
const (
	przedrostekMigawki = "migawka-"
	przedrostekZrodla  = "zrodloprz-"
	przedrostekNotatki = "notatkaprz-"
)

// adapterPrzegladarki wypełnia port Przegladarka. Jedna zależność:
// repozytorium modułu, wspólne dla migawek, źródeł i notatek — jedno
// repozytorium, trzy pliki.
type adapterPrzegladarki struct {
	repozytorium dane.RepozytoriumPrzegladania
	// magazyn trzyma bajty materiału sesji: zrzutów, archiwów i odniesień
	// monitorów, nie ich kopię.
	magazyn *magazynTresciBiblioteki
	// katalogDanych jest korzeniem, względem którego liczone są odwołania.
	katalogDanych string
	// silnik uruchamia stronę, gdy czynność wymaga strony wykonanej, a nie jej
	// źródła.
	silnik *silnikPrzegladarki
}

// nowyAdapterPrzegladarki wiąże port z repozytorium modułu przeglądania, bez
// którego adapter nie ma skąd czytać ani gdzie zapisać migawki.
func nowyAdapterPrzegladarki(repozytorium dane.RepozytoriumPrzegladania) *adapterPrzegladarki {
	return &adapterPrzegladarki{repozytorium: repozytorium}
}

// Nawiguj pobiera stronę spod adresu wskazanego przez Operatora i zapisuje
// migawkę jej treści realnym HTTP GET-em; nieudane pobranie zwraca błąd bez
// migawki.
func (a *adapterPrzegladarki) Nawiguj(ctx context.Context,
	z shared.BrowserNavigateRequest) (shared.BrowserNavigateResponse, error) {

	okno := wartoscTekstu(&z.WindowId)
	if okno == "" || z.Url == "" {
		return shared.BrowserNavigateResponse{}, bladWskazaniaPrzegladarki(
			"komenda navigate bez okna lub adresu")
	}
	// Adres jest obcinany przed pobraniem, żeby w migawce wylądował ten sam
	// łańcuch pobrania.
	adres := strings.TrimSpace(z.Url)
	tresc, err := pobierzStrone(ctx, adres)
	if err != nil {
		// Zasób, którego nie da się pokazać jako strony, nie kończy drogi odmową:
		// przeglądarka pobiera plik.
		if errors.Is(err, errZasobNieJestStrona) {
			return shared.BrowserNavigateResponse{}, a.odlozPobranie(ctx, okno, adres, err)
		}
		return shared.BrowserNavigateResponse{}, bladPobraniaStrony(adres, err)
	}
	zapisana, err := a.repozytorium.ZapiszMigawke(ctx, dane.MigawkaStrony{
		Kod: nowyIdentyfikator(przedrostekMigawki), Okno: okno, Url: adres,
		Tytul:           wskaznikTekstu(tresc.Tytul),
		TekstOdwolanie:  wskaznikTekstu(tresc.Tekst),
		ZrodloOdwolanie: wskaznikTekstu(tresc.Html),
	})
	if err != nil {
		return shared.BrowserNavigateResponse{}, bladPrzegladarki(err)
	}
	return shared.BrowserNavigateResponse{Snapshot: migawkaKontraktu(zapisana, true, true)}, nil
}

// Migawka oddaje bieżącą treść strony widoczną Operatorowi i modelowi,
// sięgając po najświeższy wiersz historii nawigacji okna wskazanego żądaniem.
func (a *adapterPrzegladarki) Migawka(ctx context.Context,
	z shared.BrowserSnapshotGetRequest) (shared.BrowserSnapshotGetResponse, error) {

	okno := wartoscTekstu(&z.WindowId)
	if okno == "" {
		return shared.BrowserSnapshotGetResponse{}, bladWskazaniaPrzegladarki(
			"komenda snapshot.get bez okna")
	}
	wiersz, err := a.repozytorium.OstatniaMigawka(ctx, okno)
	if err != nil {
		return shared.BrowserSnapshotGetResponse{}, bladBrakuMigawki(okno, err)
	}
	dolaczHtml := z.IncludeHtml != nil && *z.IncludeHtml
	dolaczZrzut := z.IncludeScreenshot != nil && *z.IncludeScreenshot
	return shared.BrowserSnapshotGetResponse{Snapshot: migawkaKontraktu(wiersz, dolaczHtml, dolaczZrzut)}, nil
}

// migawkaKontraktu składa `BrowserSnapshot` z wiersza repozytorium. Tekst
// renderowany wychodzi zawsze; HTML i zrzut ekranu wychodzą tylko na żądanie
// flag.
func migawkaKontraktu(wiersz dane.MigawkaStrony, dolaczHtml, dolaczZrzut bool) shared.BrowserSnapshot {
	migawka := shared.BrowserSnapshot{
		Id: wiersz.Kod, WindowId: wiersz.Okno, Url: wiersz.Url,
		Title: wiersz.Tytul, Text: wiersz.TekstOdwolanie,
		CapturedAt: chwilaBazy(wiersz.Utworzono),
	}
	if dolaczHtml {
		migawka.Html = wiersz.ZrodloOdwolanie
	}
	if dolaczZrzut {
		migawka.ScreenshotRef = wiersz.ZrzutOdwolanie
	}
	return migawka
}

// bladPrzegladarki znakuje usterkę wewnętrzną modułu kodem kontraktu
// `internal_error`, odróżnioną od winy żądania Operatora.
func bladPrzegladarki(err error) error {
	if err == nil {
		return nil
	}
	return protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeInternalError, err))
}

// bladPobraniaStrony nazywa nieudane pobranie strony — Operator dostaje wprost
// powód, a nie pustą migawkę podszytą pod sukces. Kod odmowy dobiera
// `kodOdmowyPobrania`, bo powód rozstrzyga, czy ponowienie ma sens.
func bladPobraniaStrony(url string, err error) error {
	kod, zdanie := kodOdmowyPobrania(err)
	return protocol.JakoError(protocol.NowyBlad(kod,
		"moduł Browser: nie udało się pobrać strony "+url+": "+zdanie))
}

// kodOdmowyPobrania przekłada powód nieudanego pobrania na kod kontraktu,
// rozstrzygając, czyj to brak: wołającego, strony, czy drogi między nimi.
func kodOdmowyPobrania(err error) (protocol.KodBledu, string) {
	var bladAdresu bladAdresuStrony
	if errors.As(err, &bladAdresu) {
		return shared.ErrorCodeValidationFailed, err.Error()
	}
	if errors.Is(err, errStronaNieodpowiedziala) {
		return shared.ErrorCodeChannelUnavailable, err.Error()
	}
	var bladStanu bladStanuStrony
	if errors.As(err, &bladStanu) {
		switch {
		case bladStanu.Kod == 404 || bladStanu.Kod == 410:
			return shared.ErrorCodeNotFound, err.Error()
		case bladStanu.Kod == 401 || bladStanu.Kod == 403:
			return shared.ErrorCodePermissionDenied, err.Error()
		case bladStanu.Kod == 429:
			return shared.ErrorCodeRateLimited, err.Error()
		case bladStanu.Kod >= 500:
			return shared.ErrorCodeChannelUnavailable, err.Error()
		}
	}
	// Reszta jest własnością wskazanego zasobu: ponowienie nic nie zmieni.
	return shared.ErrorCodeValidationFailed, err.Error()
}

// bladWskazaniaPrzegladarki nazywa brak danych wymaganych w żądaniu — to błąd
// wołającego, nie usterka rdzenia.
func bladWskazaniaPrzegladarki(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"moduł Browser: "+powod))
}

// bladBrakuMigawki odróżnia „okno nie miało jeszcze migawki” od usterki odczytu —
// Operator dostaje odmowę zrozumiałą wprost, nie pustą migawkę udającą sukces.
func bladBrakuMigawki(okno string, err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Browser: okno nie ma jeszcze migawki: "+okno))
	}
	return bladPrzegladarki(err)
}
