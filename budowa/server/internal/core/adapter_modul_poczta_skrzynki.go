// Odpowiedzialność pliku: trzy czynności nad katalogiem skrzynek — podpięcie,
// odpięcie i rozpoznanie nastaw z urządzenia. Żadna z nich nie jest narzędziem
// modelu: kontrakt trzyma je poza wykazem narzędzi, poza dostępem modelu.
package core

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/poczta"
	"danacoconsole/shared"
)

// Podepnij zapisuje skrzynkę i odkłada jej sekret w sejfie — obsługuje
// `mail.account.add`. Pierwsza podpięta skrzynka jest domyślna, bez osobnego
// pytania.
func (a *adapterPoczty) Podepnij(ctx context.Context,
	z shared.MailAccountAddRequest) (shared.MailAccountAddResponse, error) {

	if a.skrzynki == nil {
		return shared.MailAccountAddResponse{}, bladPoczty(
			errors.New("repozytorium skrzynek nie jest wpięte"))
	}
	adres := strings.TrimSpace(z.Address)
	if adres == "" {
		return shared.MailAccountAddResponse{}, bladWskazaniaPoczty(
			"komenda mail.account.add bez adresu skrzynki — serwer nie zgaduje, którą pocztę podpiąć")
	}
	host := strings.TrimSpace(wartoscLubPustka(z.IncomingHost))
	if host == "" {
		// Zgadywanie hosta po domenie adresu byłoby zmyśleniem nastawy, trafnym
		// tylko u części dostawców.
		return shared.MailAccountAddResponse{}, bladWskazaniaPoczty(
			"komenda mail.account.add bez serwera poczty przychodzącej dla " + adres +
				" — serwer nie zgaduje hosta dostawcy; podpowiedzi z urządzenia oddaje mail.account.discover")
	}

	istniejace, err := a.skrzynki.Skrzynki(ctx)
	if err != nil {
		return shared.MailAccountAddResponse{}, bladPoczty(err)
	}
	kod := kodSkrzynki(istniejace, adres)

	// Sekret idzie do sejfu przed zapisem wiersza, przed skrzynką bez poświadczenia.
	odwolanie, err := odwolaniePoswiadczenia(ctx, a.sejf, przedrostekBytuSejfuPoczty+kod, z.Secret)
	if err != nil {
		return shared.MailAccountAddResponse{}, bladPoczty(err)
	}

	wiersz := dane.SkrzynkaOperatora{
		Kod:              kod,
		Adres:            adres,
		NazwaWyswietlana: z.DisplayName,
		Protokol:         protokolSkrzynki(z.Protocol),
		Zrodlo:           string(shared.MailAccountSourceOperator),
		HostOdbioru:      host,
		PortOdbioru:      wartoscLubZero(z.IncomingPort),
		HostWysylki:      strings.TrimSpace(wartoscLubPustka(z.OutgoingHost)),
		PortWysylki:      wartoscLubZero(z.OutgoingPort),
		Uzytkownik:       pierwszyTekstNiepusty(wartoscLubPustka(z.Username), adres),
		HasloOdwolanie:   odwolanie,
		// Tryb szyfrowania czyta się z numeru portu wg IANA, port pominięty
		// daje szyfrowanie.
		SzyfrujOdbior:  czySzyfrowanyPort(wartoscLubZero(z.IncomingPort), 993, 995),
		SzyfrujWysylke: czySzyfrowanyPort(wartoscLubZero(z.OutgoingPort), 465),
		WeryfikujTLS:   true,
		Domyslna:       len(istniejace) == 0,
	}
	zapisana, err := a.skrzynki.Zapisz(ctx, wiersz)
	if err != nil {
		return shared.MailAccountAddResponse{}, bladPoczty(err)
	}

	lacznosc := a.czyLacznosc(ctx, zapisana)
	return shared.MailAccountAddResponse{
		Account:   skrzynkaKontraktu(zapisana, lacznosc),
		Connected: lacznosc,
	}, nil
}

// Odepnij usuwa skrzynkę z platformy i kasuje jej sekret — obsługuje
// `mail.account.remove`. Skrzynki u dostawcy nie tyka: rdzeń nie zna komendy,
// którą kasuje się cudze konto pocztowe.
func (a *adapterPoczty) Odepnij(ctx context.Context,
	z shared.MailAccountRemoveRequest) (shared.MailAccountRemoveResponse, error) {

	if a.skrzynki == nil {
		return shared.MailAccountRemoveResponse{}, bladPoczty(
			errors.New("repozytorium skrzynek nie jest wpięte"))
	}
	kod := strings.TrimSpace(z.AccountId)
	if kod == "" {
		return shared.MailAccountRemoveResponse{}, bladWskazaniaPoczty(
			"komenda mail.account.remove bez wskazania skrzynki")
	}
	odpiete, err := a.skrzynki.Usun(ctx, kod)
	if err != nil {
		return shared.MailAccountRemoveResponse{}, bladPoczty(err)
	}
	if odpiete {
		// Sekret ginie razem z wierszem, nie zostaje hasłem do zapomnianej skrzynki.
		usunPoswiadczenie(ctx, a.sejf, przedrostekBytuSejfuPoczty+kod)
	}
	return shared.MailAccountRemoveResponse{Removed: odpiete}, nil
}

// Rozpoznaj czyta nastawy klientów poczty zainstalowanych na urządzeniu —
// obsługuje `mail.account.discover`. Niczego nie podpina i po żadne hasło
// nie sięga.
func (a *adapterPoczty) Rozpoznaj(_ context.Context,
	_ shared.MailAccountDiscoverRequest) (shared.MailAccountDiscoverResponse, error) {

	katalog, err := os.UserHomeDir()
	if err != nil {
		return shared.MailAccountDiscoverResponse{}, bladPoczty(
			fmt.Errorf("nie da się ustalić katalogu domowego Operatora: %w", err))
	}
	rozpoznane := poczta.Rozpoznaj(katalog)
	skrzynki := make([]shared.MailAccount, 0, len(rozpoznane))
	for _, r := range rozpoznane {
		skrzynki = append(skrzynki, rozpoznanaKontraktu(r))
	}
	return shared.MailAccountDiscoverResponse{Accounts: skrzynki, Total: len(skrzynki)}, nil
}

// czySzyfrowanyPort rozstrzyga tryb gniazda po numerze portu, tym samym
// przydziałem IANA, którym kieruje się `Podepnij`. Port zerowy znaczy szyfrowanie.
func czySzyfrowanyPort(port int, szyfrowane ...int) bool {
	if port <= 0 {
		return true
	}
	for _, s := range szyfrowane {
		if port == s {
			return true
		}
	}
	return false
}

// kodSkrzynki oddaje kod skrzynki o tym adresie, jeśli już jest podpięta,
// żeby ponowne podpięcie poprawiło nastawy, a nie założyło drugiego wiersza.
func kodSkrzynki(istniejace []dane.SkrzynkaOperatora, adres string) string {
	for _, w := range istniejace {
		if strings.EqualFold(w.Adres, adres) {
			return w.Kod
		}
	}
	return nowyIdentyfikator(przedrostekSkrzynki)
}

// protokolSkrzynki bierze protokół żądania, a przy jego braku IMAP, tak jak
// stanowi kontrakt: brak pola bierze wartość imap.
func protokolSkrzynki(wskazany *shared.MailProtocol) string {
	if wskazany != nil && strings.TrimSpace(string(*wskazany)) != "" {
		return string(*wskazany)
	}
	return poczta.ProtokolImap
}
