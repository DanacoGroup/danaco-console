// Odpowiedzialność pliku: moduł Browser — typ adaptera, konstruktor i dwie
// komendy migawki strony (`browser.navigate`, `browser.snapshot.get`). Źródła
// i notatki zebrane w toku przeglądania mają własny plik adaptera i własny port.
//
// Przeglądarka jest częścią każdego środowiska. `Nawiguj` sięga po stronę
// realnym HTTP GET-em (`przegladarka_pobieranie.go`, biblioteka standardowa)
// i wypełnia migawkę tytułem, tekstem renderowanym i HTML-em pobranej strony.
// Kontrakt `BrowserSnapshot` niesie te pola (`title`, `text`, `html`), więc
// adapter je wypełnia; puste zostaje tylko `screenshotRef`, bo zrzut ekranu
// wymaga silnika przeglądarki spoza `net/http`. Flagi
// `IncludeHtml`/`IncludeScreenshot` w `browser.snapshot.get` sterują tym, co
// migawka oddaje z tego, co ma — HTML bywa ciężki, więc wychodzi na żądanie.
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

// adapterPrzegladarki wypełnia port Przegladarka. Jedna zależność: repozytorium
// modułu, wspólne dla migawek, źródeł i notatek (jedno repozytorium, trzy pliki
// — patrz `dane/przegladarka.go`).
type adapterPrzegladarki struct {
	repozytorium dane.RepozytoriumPrzegladania
	// magazyn trzyma bajty materiału sesji: zrzutów, archiwów, odniesień
	// monitorów i rejestrów sieciowych. Wiersz w bazie jest wskazaniem na nie,
	// nie ich kopią (`adapter_modul_przegladarka_zaplecze.go`).
	magazyn *magazynTresciBiblioteki
	// katalogDanych jest korzeniem, względem którego liczone są odwołania
	// magazynu wychodzące kontraktem.
	katalogDanych string
	// silnik uruchamia stronę, gdy czynność wymaga strony wykonanej, a nie
	// samego jej źródła: zrzutu, drzewa DOM, konsoli, rejestru sieciowego,
	// emulacji urządzenia i przewinięcia.
	silnik *silnikPrzegladarki
}

// nowyAdapterPrzegladarki wiąże port z repozytorium modułu.
func nowyAdapterPrzegladarki(repozytorium dane.RepozytoriumPrzegladania) *adapterPrzegladarki {
	return &adapterPrzegladarki{repozytorium: repozytorium}
}

// Nawiguj pobiera stronę spod adresu wskazanego przez Operatora i zapisuje
// migawkę jej treści. Pobranie idzie realnym HTTP GET-em (`pobierzStrone`);
// gdy strona milczy, odpowiada błędem albo nie jest do odczytu, Nawiguj wraca
// uczciwym błędem, a migawki nie tworzy — historia nawigacji nie zapisuje
// przejść, które się nie odbyły.
func (a *adapterPrzegladarki) Nawiguj(ctx context.Context,
	z shared.BrowserNavigateRequest) (shared.BrowserNavigateResponse, error) {

	okno := wartoscTekstu(&z.WindowId)
	if okno == "" || z.Url == "" {
		return shared.BrowserNavigateResponse{}, bladWskazaniaPrzegladarki(
			"komenda navigate bez okna lub adresu")
	}
	// Adres obcinamy tutaj, przed pobraniem, żeby w migawce wylądował dokładnie
	// ten łańcuch, którym rdzeń pobierał stronę. Obcięcie zostawione samemu
	// sprawdzianowi protokołu dałoby dwie prawdy o jednym adresie.
	adres := strings.TrimSpace(z.Url)
	tresc, err := pobierzStrone(ctx, adres)
	if err != nil {
		// Zasób, którego nie da się pokazać jako strony (dokument, obraz,
		// archiwum), nie kończy drogi odmową „to nie strona": przeglądarka
		// w takiej sytuacji POBIERA plik i tak samo robi moduł. Odmowa zostaje —
		// migawki z tego nie ma — ale niesie identyfikator pobrania, które
		// naprawdę powstało i którego bajty leżą w magazynie.
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

// Migawka oddaje bieżącą treść strony widoczną Operatorowi i modelowi.
// `BrowserSnapshotGetRequest` nie wskazuje kodu migawki — pyta o okno, więc
// adapter sięga po jej najświeższy wiersz (`OstatniaMigawka`), zgodnie z tym,
// że historia nawigacji to kolejne wiersze `migawka_strony`.
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
// renderowany wychodzi zawsze — to podstawowa treść widoczna modelowi; HTML i
// zrzut ekranu wychodzą tylko na żądanie flag żądania, bo bywają ciężkie, a
// domyślnie nikt o nie nie prosił.
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

// bladPrzegladarki znakuje usterkę wewnętrzną kodem kontraktu (wzór: automations).
func bladPrzegladarki(err error) error {
	if err == nil {
		return nil
	}
	return protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeInternalError, err))
}

// bladPobraniaStrony nazywa nieudane pobranie strony — Operator dostaje wprost
// powód (strona milczy, odpowiedziała błędem, nie jest do odczytu), a nie pustą
// migawkę podszytą pod sukces. Kod odmowy dobiera `kodOdmowyPobrania`
// — bo powód rozstrzyga, czy ponowienie ma sens.
//
// Jeden kod na wszystkie nieszczęścia nie wystarcza: `validation_failed`
// z `retryable:false` mówiłby o niezgodności żądania także wtedy, gdy żądanie
// było zgodne z kontraktem, a gospodarz milczał albo witryna oddała 503.
// Klient z pętlą ponowień dostawałby wtedy „nie ponawiaj" przy usterce z natury
// przemijającej.
func bladPobraniaStrony(url string, err error) error {
	kod, zdanie := kodOdmowyPobrania(err)
	return protocol.JakoError(protocol.NowyBlad(kod,
		"moduł Browser: nie udało się pobrać strony "+url+": "+zdanie))
}

// kodOdmowyPobrania przekłada powód nieudanego pobrania na kod kontraktu.
// Rozstrzyga, czyj to brak: wołającego (zły adres, żądanie odrzucone przez
// witrynę), strony (nie ma jej pod tym adresem) czy drogi (transport zerwany,
// witryna chwilowo padła, tempo ograniczone).
//
// Kontrakt nie ma kodu „zasób zewnętrzny chwilowo niedostępny"; jedynym kodem
// ponawialnym, który nie kłamie o usterce rdzenia (`internal_error`), jest
// `channel_unavailable` — droga na zewnątrz jest niedostępna dla tego jednego
// wywołania.
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
	// Reszta — treść nie do odczytu, strona ponad granicą rozmiaru, czytanie
	// przerwane — jest własnością wskazanego zasobu: ponowienie nic nie zmieni.
	return shared.ErrorCodeValidationFailed, err.Error()
}

// bladWskazaniaPrzegladarki nazywa brak danych w żądaniu — błąd Operatora, nie rdzenia.
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
