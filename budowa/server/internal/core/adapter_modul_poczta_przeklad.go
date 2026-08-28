// Odpowiedzialność pliku: przekład między kształtem rdzenia a kształtem kontraktu — wiersz skrzynki, nagłówek listu na typy shared, plus pomocniki wartości opcjonalnych.
package core

import (
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/poczta"
	"danacoconsole/shared"
)

// skrzynkaKontraktu przekłada wiersz skrzynki na MailAccount, bez poświadczenia ani odwołania sejfu, sprawy rdzenia, nie klienta.
func skrzynkaKontraktu(w dane.SkrzynkaOperatora, lacznosc bool) shared.MailAccount {
	zrodlo := shared.MailAccountSource(w.Zrodlo)
	skrzynka := shared.MailAccount{
		Id:      w.Kod,
		Address: w.Adres,
		// Server jest polem wymaganym kontraktu, tym samym co IncomingHost; powielenie jest w kontrakcie.
		Server:       w.HostOdbioru,
		Connected:    lacznosc,
		Protocol:     shared.MailProtocol(w.Protokol),
		DisplayName:  w.NazwaWyswietlana,
		IncomingHost: tekstOpcjonalny(w.HostOdbioru),
		IncomingPort: liczbaOpcjonalna(w.PortOdbioru),
		OutgoingHost: tekstOpcjonalny(w.HostWysylki),
		OutgoingPort: liczbaOpcjonalna(w.PortWysylki),
		Username:     tekstOpcjonalny(w.Uzytkownik),
		Source:       &zrodlo,
	}
	return skrzynka
}

// rozpoznanaKontraktu przekłada skrzynkę odczytaną z urządzenia; Id i Connected zostają puste, bo skrzynka nie jest podpięta ani rdzeń się z nią nie łączył.
func rozpoznanaKontraktu(r poczta.Rozpoznana) shared.MailAccount {
	zrodlo := shared.MailAccountSource(shared.MailAccountSourceUrzadzenie)
	return shared.MailAccount{
		Address:      r.Adres,
		Server:       r.HostOdbioru,
		Connected:    false,
		Protocol:     shared.MailProtocol(r.Protokol),
		DisplayName:  tekstOpcjonalny(r.NazwaWyswietlana),
		IncomingHost: tekstOpcjonalny(r.HostOdbioru),
		IncomingPort: liczbaOpcjonalna(r.PortOdbioru),
		OutgoingHost: tekstOpcjonalny(r.HostWysylki),
		OutgoingPort: liczbaOpcjonalna(r.PortWysylki),
		Username:     tekstOpcjonalny(r.Uzytkownik),
		Source:       &zrodlo,
	}
}

// wiadomoscKontraktu przekłada nagłówek listu na `MailMessage`. Treść zostaje
// pusta — wypełnia ją wyłącznie `mail.message.get`, tak jak stanowi kontrakt.
func wiadomoscKontraktu(n poczta.Naglowek) shared.MailMessage {
	nieprzeczytana := n.Nieprzeczytana
	return shared.MailMessage{
		Id:          n.Identyfikator,
		Folder:      n.Folder,
		From:        n.Od,
		To:          n.Do,
		Cc:          n.Kopia,
		Subject:     tekstOpcjonalny(n.Temat),
		Date:        n.Chwila.UnixMilli(),
		Preview:     tekstOpcjonalny(n.Zapowiedz),
		Attachments: n.NazwyZalacznikow,
		Unread:      &nieprzeczytana,
		ThreadId:    tekstOpcjonalny(n.Watek),
	}
}

// Tekst opcjonalny bierze się z tekstOpcjonalny w akcje.go, jedynego pomocnika brakującej wartości.

// liczbaOpcjonalna oddaje wskaźnik na liczbę dodatnią albo nic. Port zerowy
// nie istnieje, więc zero jest tu brakiem wskazania, nie wartością.
func liczbaOpcjonalna(wartosc int) *int {
	if wartosc <= 0 {
		return nil
	}
	return &wartosc
}

// wartoscLubPustka rozpakowuje wskaźnik na tekst — brak znaczy pustkę, wartość niewskazaną w żądaniu klienta.
func wartoscLubPustka(wskaznik *string) string {
	if wskaznik == nil {
		return ""
	}
	return *wskaznik
}

// wartoscLubZero rozpakowuje wskaźnik na liczbę — brak znaczy zero, czyli
// „nie wskazano" (wołający bierze wtedy port domyślny protokołu).
func wartoscLubZero(wskaznik *int) int {
	if wskaznik == nil {
		return 0
	}
	return *wskaznik
}

// pierwszyTekstNiepusty oddaje pierwszą niepustą wartość z podanych, dla ustalenia wartości z priorytetem.
func pierwszyTekstNiepusty(wartosci ...string) string {
	for _, w := range wartosci {
		if strings.TrimSpace(w) != "" {
			return w
		}
	}
	return ""
}
