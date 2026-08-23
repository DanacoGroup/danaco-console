// Odpowiedzialność pliku: trzy czynności nad katalogiem skrzynek — podpięcie
// (`mail.account.add`), odpięcie (`mail.account.remove`) i rozpoznanie nastaw
// z urządzenia (`mail.account.discover`). Typ adaptera, jego montaż i otwieranie
// połączenia leżą w `adapter_modul_poczta.go`.
//
// Żadna z tych trzech nie jest narzędziem modelu. Kontrakt trzyma je poza
// wykazem narzędzi, bo rozstrzygają, do czego platforma ma dostęp: model, który
// sam podpina albo odpina skrzynkę, decyduje o rzeczy, której mu nie
// powierzono.
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
// `mail.account.add`.
//
// Pierwsza podpięta skrzynka jest domyślna i nie ma o to osobnego pytania: przy
// jednej skrzynce odpowiedź jest jedna, a bez niej wszystkie komendy bez
// `accountId` odmawiałyby mimo podpiętej poczty.
func (a *adapterPoczty) Podepnij(ctx context.Context,
	z shared.MailAccountAddRequest) (shared.MailAccountAddResponse, error) {

	if a.skrzynki == nil {
		return shared.MailAccountAddResponse{}, bladPoczty(
			errors.New("repozytorium skrzynek nie jest wpięte"))
	}
	adres := strings.TrimSpace(z.Address)
	if adres == "" {
		return shared.MailAccountAddResponse{}, bladWskazaniaPoczty(
			"komenda mail.account.add bez adresu skrzynki — rdzeń nie zgaduje, którą pocztę podpiąć")
	}
	host := strings.TrimSpace(wartoscLubPustka(z.IncomingHost))
	if host == "" {
		// Zgadywanie hosta po domenie adresu byłoby zmyśleniem nastawy: „imap."
		// przed domeną trafia u części dostawców i nie trafia u reszty, a skutkiem
		// nietrafienia jest skrzynka podpięta do serwera, którego nie ma.
		return shared.MailAccountAddResponse{}, bladWskazaniaPoczty(
			"komenda mail.account.add bez serwera poczty przychodzącej dla " + adres +
				" — rdzeń nie zgaduje hosta dostawcy; podpowiedzi z urządzenia oddaje mail.account.discover")
	}

	istniejace, err := a.skrzynki.Skrzynki(ctx)
	if err != nil {
		return shared.MailAccountAddResponse{}, bladPoczty(err)
	}
	kod := kodSkrzynki(istniejace, adres)

	// Sekret idzie do sejfu przed zapisem wiersza. Odwrotna kolejność
	// zostawiłaby przy awarii sejfu skrzynkę bez poświadczenia, wyglądającą na
	// gotową, a odmawiającą przy pierwszym użyciu.
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
		// Tryb szyfrowania czytamy z numeru portu. Kontrakt nie ma pola
		// „szyfruj", bo przydział IANA rozstrzyga to jednoznacznie: 993 i 995 są
		// szyfrowane od pierwszego bajtu, 143 i 110 nie są. Port pominięty daje
		// szyfrowanie, bo wartością domyślną ma być bezpieczniejsza z dwóch.
		//
		// Port nieszyfrowany jest wskazaniem, nie furtką: znaczy, że skrzynka
		// stoi na serwerze bez TLS-a. Podniesienie go po cichu zerwałoby uścisk
		// dłoni, a obniżenie szyfrowanego odebrałoby ochronę bez pytania, więc
		// zapisujemy wartość podaną. Wysyłka i tak podnosi się do STARTTLS, gdy
		// serwer go ogłasza (`poczta/smtp.go`).
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
// `mail.account.remove`. Skrzynki u dostawcy nie tyka i nie ma jak tknąć:
// rdzeń nie zna komendy, którą kasuje się cudze konto pocztowe.
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
		// Sekret ginie razem z wierszem. Zostawiony w sejfie byłby hasłem do
		// skrzynki, o której platforma zapomniała, i nie byłoby już czym kazać
		// go skasować.
		usunPoswiadczenie(ctx, a.sejf, przedrostekBytuSejfuPoczty+kod)
	}
	return shared.MailAccountRemoveResponse{Removed: odpiete}, nil
}

// Rozpoznaj czyta nastawy klientów poczty zainstalowanych na urządzeniu —
// obsługuje `mail.account.discover`. Niczego nie podpina i po żadne hasło nie
// sięga; szczegóły i granice tego odczytu opisuje `poczta/rozpoznanie.go`.
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

// czySzyfrowanyPort rozstrzyga tryb gniazda po numerze portu — patrz komentarz
// przy `Podepnij`. Zero (port niewskazany) znaczy szyfrowanie, bo domyślną
// wartością ma być bezpieczniejsza z dwóch.
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

// kodSkrzynki oddaje kod skrzynki o tym adresie, jeśli już jest podpięta —
// żeby ponowne podpięcie tej samej poczty poprawiło nastawy, a nie założyło
// drugiego wiersza o tym samym adresie i innym kodzie.
func kodSkrzynki(istniejace []dane.SkrzynkaOperatora, adres string) string {
	for _, w := range istniejace {
		if strings.EqualFold(w.Adres, adres) {
			return w.Kod
		}
	}
	return nowyIdentyfikator(przedrostekSkrzynki)
}

// protokolSkrzynki bierze protokół żądania, a przy jego braku — IMAP, tak jak
// stanowi kontrakt („brak bierze imap").
func protokolSkrzynki(wskazany *shared.MailProtocol) string {
	if wskazany != nil && strings.TrimSpace(string(*wskazany)) != "" {
		return string(*wskazany)
	}
	return poczta.ProtokolImap
}
