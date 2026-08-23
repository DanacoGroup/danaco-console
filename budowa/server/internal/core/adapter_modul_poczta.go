// Odpowiedzialność pliku: RDZEŃ JAKO KLIENT SKRZYNKI OPERATORA — typ adaptera,
// jego montaż, wybór skrzynki, otwieranie połączenia i wykaz skrzynek
// (`mail.account.list`). Podpięcie, odpięcie i rozpoznanie z urządzenia leżą
// w `adapter_modul_poczta_skrzynki.go`, odczyt wiadomości w `_wiadomosci.go`,
// szkic i wysyłka w `_wysylka.go`, przekład na kontrakt w `_przeklad.go` —
// plik wedle odpowiedzialności.
//
// ── CZYM TEN MODUŁ JEST, A CZYM NIE JEST ────────────────────────────────────
// Aplikacja nie ma własnego serwera poczty — używa skrzynki, którą Operator już
// ma skonfigurowaną na urządzeniu albo w chmurze. Rdzeń jest więc klientem
// cudzej skrzynki i nikim
// więcej: nie stawia serwera, nie zakłada kont pocztowych, nie pośredniczy
// przez żadną infrastrukturę Danaco i nie trzyma cudzej poczty u siebie.
// Podpina skrzynkę, którą Operator już ma, i rozmawia z nią jej protokołem.
//
// ── POŚWIADCZENIE IDZIE WYŁĄCZNIE DO SEJFU ────────────────────────
// Hasło albo token przychodzi polem `secret` żądania `mail.account.add`
// i natychmiast ląduje w sejfie poświadczeń — tym samym, którym jadą konta
// i punkty dostępu. Baza dostaje odwołanie („sejf:poczta:<kod>"), nigdy sekret;
// kolumny na sekret nie ma w schemacie w ogóle. Do dziennika sekret nie trafia,
// bo ten moduł nie loguje niczego. Do odpowiedzi komendy nie trafia, bo
// `MailAccount` w kontrakcie nie ma pola, w które dałoby się go włożyć.
//
// ── ODMOWA NAZYWA BRAK, KAŻDY INNYM ZDANIEM ──────────────
// Brak skrzynki, brak poświadczenia i brak łączności to trzy różne rzeczy,
// z których każdą Operator naprawia inaczej: pierwszą podpięciem skrzynki,
// drugą podaniem hasła, trzecią zajrzeniem do sieci albo do dostawcy. Jedno
// wspólne „poczta niedostępna" kazałoby mu zgadywać, którą.
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

// adapterPoczty wypełnia port Poczta. Cztery zależności, bo pełny łańcuch
// „przeczytaj list, przeanalizuj załącznik, odpisz" potrzebuje wszystkich:
// wiersza skrzynki, sekretu z sejfu, miejsca na bajty załączników i wiersza
// zasobu, po którym model sięgnie po nie arsenałem obrazu i dokumentów.
type adapterPoczty struct {
	skrzynki dane.RepozytoriumSkrzynek
	sejf     SejfPoswiadczen
	// zasoby jest repozytorium Designu — TYM SAMYM, do którego pisze
	// `design.asset.upload`. Drugiego magazynu zasobów rdzeń nie ma i mieć nie
	// będzie: załącznik listu wciągnięty osobną drogą byłby zasobem,
	// którego narzędzia obrazu i dokumentów nie widzą.
	zasoby dane.RepozytoriumDesignu
	// magazyn jest miejscem na BAJTY załączników. Ten sam typ, co w Designie
	// i Bibliotece — blob pod sumą sha256, zapis atomowy.
	magazyn *magazynTresciBiblioteki
}

// nowyAdapterPoczty wiąże port z repozytorium skrzynek i wpina magazyn
// załączników oparty o katalog danych rdzenia — wzorem `nowyAdapterDesignu`.
// Katalog obowiązujący wchodzi montażem (`ZKatalogiemDanych`); domyślny zostaje
// dla wywołania bez montażu, żeby konstruktor nigdy nie oddał adaptera bez
// magazynu.
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

// magazynZalacznikowPoczty składa magazyn bajtów nad katalogiem danych.
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
// `mail.account.list`.
//
// Łączność mierzymy, a nie deklarujemy. Pole `connected` znaczy „rdzeń ma z nią
// łączność", więc wypełnienie go na stałe prawdą byłoby powodzeniem czynności,
// której nikt nie wykonał. Otwieramy więc połączenie do każdej
// skrzynki i natychmiast je zamykamy. Cena jest widoczna: wykaz trwa tyle, co
// suma uścisków dłoni. Płacimy ją, bo pytanie „czy moja poczta działa" bez
// sprawdzenia nie ma sensu.
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
		// Kod `channel_unavailable`, nie `internal_error`: skrzynka jest po
		// drugiej stronie sieci i jej niedostępność bywa chwilowa, więc odmowa
		// jest PONAWIALNA (`shared.KodyPonawialne`). Rdzeń nie zawinił i nic tu
		// nie naprawi — a klient, który ponowi za minutę, ma szansę trafić.
		return nil, wiersz, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeChannelUnavailable, "moduł poczty: "+err.Error()))
	}
	return klient, wiersz, nil
}

// wybierzSkrzynke odnajduje skrzynkę wskazaną albo domyślną.
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

// nastawy składa nastawy połączenia, dobierając sekret z sejfu.
//
// Sekret żyje tylko w tej strukturze i tylko do końca komendy. Nie wraca do
// bazy (kolumny nie ma), nie wraca do odpowiedzi (pola w kontrakcie nie ma),
// nie idzie do dziennika (ten moduł nie loguje). Brak sekretu nie jest tu
// błędem: `poczta.Polacz` nazwie go drugą z trzech odmów — brakiem
// poświadczenia, odróżnionym od braku skrzynki i braku łączności.
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

// czyLacznosc sprawdza, czy rdzeń NAPRAWDĘ dosięga skrzynki — patrz komentarz
// przy `Skrzynki`. Nieudane połączenie nie jest tu błędem komendy: brak
// łączności jest FAKTEM o skrzynce, który wykaz ma pokazać, a nie powodem, dla
// którego wykaz miałby nie powstać.
func (a *adapterPoczty) czyLacznosc(ctx context.Context, w dane.SkrzynkaOperatora) bool {
	klient, err := poczta.Polacz(a.nastawy(ctx, w))
	if err != nil {
		return false
	}
	klient.Zamknij()
	return true
}

// bladPoczty nazywa awarię po stronie rdzenia — bazy, sejfu, magazynu.
func bladPoczty(przyczyna error) error {
	return protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeInternalError,
		fmt.Errorf("moduł poczty: %w", przyczyna)))
}

// bladWskazaniaPoczty nazywa żądanie niekompletne — wina jest po stronie
// wołającego, więc kod jest walidacyjny i odmowa NIE jest ponawialna.
func bladWskazaniaPoczty(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed, "moduł poczty: "+powod))
}
