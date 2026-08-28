// Odpowiedzialność pliku: klient skrzynki poczty Operatora — typ adaptera,
// jego montaż, wybór skrzynki, otwieranie połączenia i wykaz skrzynek
// (`mail.account.list`). Podpięcie i odczyt wiadomości leżą w plikach
// sąsiednich wedle odpowiedzialności.
package core

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfiguracja"
	"danacoconsole/server/internal/poczta"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

const (
	// przedrostekSkrzynki znakuje identyfikatory zewnętrzne skrzynek — wzorem
	// `zasob-` w Designie. Ten sam tekst jest kodem wiersza i częścią bytu
	// sejfu, więc odwołanie do sekretu da się odtworzyć z samej skrzynki.
	przedrostekSkrzynki = "skrzynka-"

	// przedrostekBytuSejfuPoczty znakuje wpis sejfu należący do poczty — tak samo,
	// jak robią to konta (`konto:<kod>`). Dzięki temu jeden plik sejfu unosi
	// sekrety wielu rodzajów bytów bez zderzenia kluczy.
	przedrostekBytuSejfuPoczty = "poczta:"

	// Magazyn załączników leży w katalogu danych rdzenia, obok magazynu
	// biblioteki i zasobów Designu. Podkatalog własny, bo moduły nie dzielą
	// stanu i odpięcie skrzynki nie ma prawa ruszyć cudzych bajtów.
	podkatalogPoczty      = "poczta"
	podkatalogZalacznikow = "zalaczniki"
)

// adapterPoczty wypełnia port Poczta. Cztery zależności niesie łańcuch
// przeczytania listu, analizy załącznika i odpisu: wiersz skrzynki, sekret z
// sejfu, miejsce na bajty załączników i wiersz zasobu dla arsenału obrazu i
// dokumentów.
type adapterPoczty struct {
	skrzynki dane.RepozytoriumSkrzynek
	sejf     SejfPoswiadczen
	// zasoby jest repozytorium Designu, tym samym do którego pisze
	// `design.asset.upload`.
	zasoby dane.RepozytoriumDesignu
	// magazyn jest miejscem na bajty załączników, tym samym typem blobu co
	// w Designie i Bibliotece.
	magazyn *magazynTresciBiblioteki
}

// nowyAdapterPoczty wiąże port z repozytorium skrzynek i wpina magazyn
// załączników oparty o katalog danych rdzenia. Katalog obowiązujący wchodzi
// montażem (`ZKatalogiemDanych`); domyślny zostaje dla wywołania bez montażu.
func nowyAdapterPoczty(skrzynki dane.RepozytoriumSkrzynek) *adapterPoczty {
	return &adapterPoczty{
		skrzynki: skrzynki,
		magazyn:  magazynZalacznikowPoczty(konfiguracja.KatalogDanychDomyslny()),
	}
}

// ZSejfem wpina sejf poświadczeń — TEN SAM, którym jadą konta, punkty dostępu
// i bramka. Bez niego `mail.account.add` nie ma gdzie odłożyć sekretu, a każde
// późniejsze połączenie odmawia, nazywając brak poświadczenia.
func (a *adapterPoczty) ZSejfem(sejf SejfPoswiadczen) *adapterPoczty {
	a.sejf = sejf
	return a
}

// ZKatalogiemDanych przestawia magazyn załączników na katalog wskazany
// konfiguracją procesu — tą samą zmienną, którą jadą sejf, magazyn biblioteki
// i zasoby Designu.
func (a *adapterPoczty) ZKatalogiemDanych(katalog string) *adapterPoczty {
	if katalog != "" {
		a.magazyn = magazynZalacznikowPoczty(katalog)
	}
	return a
}

// ZZasobami wpina repozytorium Designu jako magazyn WIERSZY zasobów. Bez niego
// `mail.message.get` odda treść listu, ale załączników nie wciągnie — i powie
// to wprost, zamiast milczeć.
func (a *adapterPoczty) ZZasobami(zasoby dane.RepozytoriumDesignu) *adapterPoczty {
	a.zasoby = zasoby
	return a
}

// magazynZalacznikowPoczty składa magazyn bajtów załączników nad katalogiem
// danych, tym samym mechanizmem blobów pod sumą sha256 co magazyn biblioteki.
func magazynZalacznikowPoczty(katalogDanych string) *magazynTresciBiblioteki {
	if strings.TrimSpace(katalogDanych) == "" {
		return nil
	}
	return &magazynTresciBiblioteki{
		katalog: katalogDanych + string(os.PathSeparator) + podkatalogPoczty +
			string(os.PathSeparator) + podkatalogZalacznikow,
	}
}

// Skrzynki oddaje podpięte skrzynki wraz ze stanem łączności — obsługuje
// `mail.account.list`. Łączność jest mierzona, a nie deklarowana: połączenie
// do każdej skrzynki zostaje otwarte i natychmiast zamknięte, zamiast
// wpisywać stałą wartość.
func (a *adapterPoczty) Skrzynki(ctx context.Context,
	_ shared.MailAccountListRequest) (shared.MailAccountListResponse, error) {

	if a.skrzynki == nil {
		return shared.MailAccountListResponse{}, bladPoczty(
			errors.New("repozytorium skrzynek nie jest wpięte"))
	}
	wiersze, err := a.skrzynki.Skrzynki(ctx)
	if err != nil {
		return shared.MailAccountListResponse{}, bladPoczty(err)
	}
	skrzynki := make([]shared.MailAccount, 0, len(wiersze))
	for _, wiersz := range wiersze {
		skrzynki = append(skrzynki, skrzynkaKontraktu(wiersz, a.czyLacznosc(ctx, wiersz)))
	}
	return shared.MailAccountListResponse{Accounts: skrzynki, Total: len(skrzynki)}, nil
}

// polacz otwiera połączenie do wskazanej skrzynki. Wołający ZAWSZE domyka je
// `Zamknij` — połączenie do cudzej skrzynki żyje tyle, co jedna komenda
// (uzasadnienie przy typie `poczta.Klient`).
func (a *adapterPoczty) polacz(ctx context.Context, wskazana *string) (*poczta.Klient, dane.SkrzynkaOperatora, error) {
	wiersz, err := a.wybierzSkrzynke(ctx, wskazana)
	if err != nil {
		return nil, dane.SkrzynkaOperatora{}, err
	}
	klient, err := poczta.Polacz(a.nastawy(ctx, wiersz))
	if err != nil {
		// Kod `channel_unavailable`: niedostępność skrzynki bywa chwilowa,
		// odmowa jest ponawialna.
		return nil, wiersz, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeChannelUnavailable, "moduł poczty: "+err.Error()))
	}
	return klient, wiersz, nil
}

// wybierzSkrzynke odnajduje skrzynkę wskazaną żądaniem albo, gdy żądanie jej
// nie wskazuje, skrzynkę domyślną Operatora.
func (a *adapterPoczty) wybierzSkrzynke(ctx context.Context, wskazana *string) (dane.SkrzynkaOperatora, error) {
	if a.skrzynki == nil {
		return dane.SkrzynkaOperatora{}, bladPoczty(errors.New("repozytorium skrzynek nie jest wpięte"))
	}
	kod := strings.TrimSpace(wartoscLubPustka(wskazana))
	if kod != "" {
		wiersz, err := a.skrzynki.Skrzynka(ctx, kod)
		if errors.Is(err, dane.ErrBrakWiersza) {
			return dane.SkrzynkaOperatora{}, protocol.JakoError(protocol.NowyBlad(
				shared.ErrorCodeNotFound, "moduł poczty: platforma nie zna skrzynki "+kod+
					" — wykaz podpiętych oddaje mail.account.list"))
		}
		if err != nil {
			return dane.SkrzynkaOperatora{}, bladPoczty(err)
		}
		return wiersz, nil
	}
	wiersz, err := a.skrzynki.SkrzynkaDomyslna(ctx)
	if errors.Is(err, dane.ErrBrakWiersza) {
		// Pierwsza z trzech odmów nazywających brak: nie ma skrzynki.
		return dane.SkrzynkaOperatora{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeNotFound, "moduł poczty: Operator nie ma podpiętej ani jednej skrzynki — "+
				"rdzeń nie stawia serwera poczty i nie zakłada kont, więc nie ma czego czytać; "+
				"podepnij skrzynkę komendą mail.account.add albo poproś o podpowiedzi z urządzenia "+
				"komendą mail.account.discover"))
	}
	if err != nil {
		return dane.SkrzynkaOperatora{}, bladPoczty(err)
	}
	return wiersz, nil
}

// nastawy składa nastawy połączenia, dobierając sekret z sejfu. Sekret żyje
// tylko w tej strukturze i tylko do końca komendy; nie wraca do bazy, do
// odpowiedzi ani do dziennika.
func (a *adapterPoczty) nastawy(ctx context.Context, w dane.SkrzynkaOperatora) poczta.Nastawy {
	sekret := ""
	if a.sejf != nil && w.HasloOdwolanie != nil {
		sekret, _ = a.sejf.Odczytaj(ctx, przedrostekBytuSejfuPoczty+w.Kod)
	}
	return poczta.Nastawy{
		Adres:            w.Adres,
		NazwaWyswietlana: wartoscLubPustka(w.NazwaWyswietlana),
		Uzytkownik:       w.Uzytkownik,
		Sekret:           sekret,
		Protokol:         w.Protokol,
		HostOdbioru:      w.HostOdbioru,
		PortOdbioru:      w.PortOdbioru,
		HostWysylki:      w.HostWysylki,
		PortWysylki:      w.PortWysylki,
		SzyfrujOdbior:    w.SzyfrujOdbior,
		SzyfrujWysylke:   w.SzyfrujWysylke,
		WeryfikujTLS:     w.WeryfikujTLS,
	}
}

// czyLacznosc sprawdza, czy rdzeń rzeczywiście dosięga skrzynki. Nieudane
// połączenie nie jest tu błędem komendy: brak łączności jest faktem o
// skrzynce, który wykaz ma pokazać.
func (a *adapterPoczty) czyLacznosc(ctx context.Context, w dane.SkrzynkaOperatora) bool {
	klient, err := poczta.Polacz(a.nastawy(ctx, w))
	if err != nil {
		return false
	}
	klient.Zamknij()
	return true
}

// bladPoczty nazywa awarię po stronie rdzenia — bazy, sejfu albo magazynu
// załączników, odróżnioną od winy żądania Operatora.
func bladPoczty(przyczyna error) error {
	return protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeInternalError,
		fmt.Errorf("moduł poczty: %w", przyczyna)))
}

// bladWskazaniaPoczty nazywa żądanie niekompletne — wina jest po stronie
// wołającego, więc kod jest walidacyjny i odmowa NIE jest ponawialna.
func bladWskazaniaPoczty(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed, "moduł poczty: "+powod))
}
