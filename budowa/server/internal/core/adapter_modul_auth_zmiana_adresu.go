// Odpowiedzialność pliku: zmiana adresu e-mail uwierzytelniającego — para listów,
// potwierdzenie kodem z nowego adresu i wycofanie drogą z adresu dotychczasowego.
package core

import (
	"context"
	"errors"
	"fmt"
	netmail "net/mail"
	"strconv"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfiguracja"
	"danacoconsole/server/internal/mail"
	"danacoconsole/server/internal/nadajnik"
	"danacoconsole/shared"
)

/*
trwanieDrogiWycofania to termin drogi z listu ostrzegawczego.

Trzy doby, nie godzina jak kod. List na adres dotychczasowy jest jedyną obroną
Operatora, któremu przejęto konto, a ten czyta pocztę wtedy, kiedy czyta —
termin krótszy niż weekend zostawiałby go bez obrony. Żadne źródło dostawy
wartości nie podaje; ta stoi tutaj i nazywa powód.
*/
const trwanieDrogiWycofania = 72 * time.Hour

// probyZmianyAdresu to sufit pomyłek kodu; po nim zamówienie schodzi.
const probyZmianyAdresu = 5

// RozpocznijZmianeAdresu obsługuje `auth.email.change.start`: zapisuje zamówienie
// i wysyła parę listów — kod na adres nowy, ostrzeżenie na dotychczasowy.
func (a *adapterUwierzytelnienia) RozpocznijZmianeAdresu(ctx context.Context,
	z shared.AuthEmailChangeStartRequest) (shared.AuthEmailChangeStartResponse, error) {

	if err := a.gotowa(); err != nil {
		return shared.AuthEmailChangeStartResponse{}, err
	}
	if err := a.kontoGotowe(); err != nil {
		return shared.AuthEmailChangeStartResponse{}, err
	}
	adresNowy := strings.TrimSpace(z.NewEmail)
	if adresNowy == "" || z.CurrentPassword == "" {
		return shared.AuthEmailChangeStartResponse{}, bladBramki(shared.ErrorCodeValidationFailed,
			"Podaj nowy adres oraz hasło bieżące.")
	}
	if _, err := netmail.ParseAddress(adresNowy); err != nil {
		return shared.AuthEmailChangeStartResponse{}, bladBramki(shared.ErrorCodeValidationFailed,
			"Podany adres nie jest poprawnym adresem e-mail.")
	}

	a.zamekZmiany.Lock()
	defer a.zamekZmiany.Unlock()

	kontoId, err := kontoWolajacego(ctx)
	if err != nil {
		return shared.AuthEmailChangeStartResponse{}, err
	}
	konto, err := a.konto.KontoPoId(ctx, kontoId)
	if err != nil {
		return shared.AuthEmailChangeStartResponse{}, err
	}
	// Hasło bieżące jest wymogiem czynności, bo adres to jedyna droga odzyskania
	// konta: kto zmieni adres, ten przejmuje konto na stałe.
	kotwica, err := a.kotwicaKonta(ctx, kontoId)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return shared.AuthEmailChangeStartResponse{}, bladBramki(shared.ErrorCodeNotFound,
			"Hasło dostępu nie zostało jeszcze ustawione.")
	}
	if err != nil {
		return shared.AuthEmailChangeStartResponse{}, err
	}
	zgadza, err := a.sekretZgadzaSieZWpisem(ctx, kotwica, z.CurrentPassword)
	if err != nil {
		return shared.AuthEmailChangeStartResponse{}, err
	}
	if !zgadza {
		return shared.AuthEmailChangeStartResponse{}, bladBramki(shared.ErrorCodeValidationFailed,
			"Podane hasło bieżące jest nieprawidłowe.")
	}
	if strings.EqualFold(adresNowy, konto.Email) {
		return shared.AuthEmailChangeStartResponse{}, bladBramki(shared.ErrorCodeValidationFailed,
			"Konto stoi już przy tym adresie.")
	}
	/* Adres zajęty odpowiada tak samo jak wolny. Różnica mówiłaby pytającemu,
	   które adresy mają konto na tej instalacji — tą samą drogą, którą zamyka
	   odpowiedź auth.recover. Listy w tym wypadku nie wychodzą. */
	if _, err := a.konto.KontoPoTozsamosci(ctx, adresNowy); err == nil {
		return shared.AuthEmailChangeStartResponse{Sent: true}, nil
	} else if !errors.Is(err, dane.ErrBrakWiersza) {
		return shared.AuthEmailChangeStartResponse{}, err
	}

	kod, err := nowyKodPotwierdzenia()
	if err != nil {
		return shared.AuthEmailChangeStartResponse{}, err
	}
	droga, err := nowyTokenBramki()
	if err != nil {
		return shared.AuthEmailChangeStartResponse{}, err
	}
	teraz := time.Now()
	zamowienie := dane.ZmianaAdresuKonta{
		KontoId:         kontoId,
		AdresNowy:       adresNowy,
		AdresPoprzedni:  konto.Email,
		SkrotKodu:       skrotTokenu(bezOdstepow(kod)),
		SkrotWycofania:  skrotTokenu(droga),
		WygasaKod:       teraz.Add(trwanieDrogiPotwierdzenia).UnixMilli(),
		WygasaWycofanie: teraz.Add(trwanieDrogiWycofania).UnixMilli(),
		Utworzono:       teraz.UnixMilli(),
		// Adresu źródłowego rdzeń nie zna: gniazdo WebSocket kończy się na
		// warstwie transportu, a kontekst żądania go nie niesie. List pokazuje
		// wtedy myślnik zamiast wiersza wymyślonego.
		Urzadzenie: a.urzadzenieSesji(ctx),
	}
	if err := a.zmianaAdresu.ZamowZmianeAdresu(ctx, zamowienie); err != nil {
		return shared.AuthEmailChangeStartResponse{}, err
	}
	if err := a.wyslijListyZmianyAdresu(ctx, zamowienie, kod, droga, teraz); err != nil {
		return shared.AuthEmailChangeStartResponse{}, err
	}
	return shared.AuthEmailChangeStartResponse{Sent: true}, nil
}

/*
wyslijListyZmianyAdresu składa i nadaje obie wiadomości.

Para jest jedną czynnością: list na adres dotychczasowy jest wymogiem
bezpieczeństwa, nie uprzejmością, więc niepowodzenie któregokolwiek z nadań
przerywa całość i zamyka zamówienie. Zamówienie bez ostrzeżenia zmieniałoby
adres tak samo cicho, jak gdyby listu nie było w dostawie wcale.
*/
func (a *adapterUwierzytelnienia) wyslijListyZmianyAdresu(ctx context.Context,
	z dane.ZmianaAdresuKonta, kod, droga string, teraz time.Time) error {

	nadawca := a.kontoNadawcze(ctx)
	if err := nadawca.Brak(); err != nil {
		return bladBramki(shared.ErrorCodeInternalError, err.Error()+"; "+dwieDrogiKontaNadawczego)
	}
	komplet, err := kompletListow()
	if err != nil {
		return err
	}
	wspolne := map[string]string{
		"support_address":   adresWsparcia,
		"year":              strconv.Itoa(teraz.Year()),
		"new_address":       z.AdresNowy,
		"previous_address":  z.AdresPoprzedni,
		"requested_at":      teraz.Format(postacCzasuListu),
		"recipient_address": z.AdresNowy,
	}
	naNowy := kopiaWartosci(wspolne)
	naNowy["code"] = rozdzielony(kod)
	naNowy["expiry_minutes"] = strconv.Itoa(int(trwanieDrogiPotwierdzenia.Minutes()))

	naPoprzedni := kopiaWartosci(wspolne)
	naPoprzedni["recipient_address"] = z.AdresPoprzedni
	naPoprzedni["revoke_url"] = konfiguracja.AdresWycofaniaZmiany(a.adresKonsoli, droga)
	naPoprzedni["revoke_expiry_hours"] = strconv.Itoa(int(trwanieDrogiWycofania.Hours()))
	naPoprzedni["ip_address"] = wartoscAlboMyslnik(z.ZrodloIP)
	naPoprzedni["device"] = wartoscAlboMyslnik(z.Urzadzenie)

	listy := []struct {
		rodzaj   mail.Kind
		odbiorca string
		wartosci map[string]string
	}{
		{mail.KindAddressChangeCode, z.AdresNowy, naNowy},
		{mail.KindAddressChangeAlert, z.AdresPoprzedni, naPoprzedni},
	}
	for _, l := range listy {
		wiadomosc, err := komplet.Build(l.rodzaj, mail.Envelope{
			From: nadawca.Nadawca(),
			To:   netmail.Address{Address: l.odbiorca},
			Date: teraz,
		}, mail.Content{Values: l.wartosci}, znakiMarki())
		if err != nil {
			a.zamknijPoNiepowodzeniuListu(ctx, z)
			return bladBramki(shared.ErrorCodeInternalError,
				fmt.Sprintf("nie udało się złożyć listu na %s: %v; %s", l.odbiorca, err, dwieDrogiKontaNadawczego))
		}
		if _, err := nadajnik.Wyslij(nadawca, wiadomosc); err != nil {
			a.zamknijPoNiepowodzeniuListu(ctx, z)
			return bladBramki(shared.ErrorCodeInternalError,
				fmt.Sprintf("nie udało się wysłać listu na %s: %v; %s", l.odbiorca, err, dwieDrogiKontaNadawczego))
		}
	}
	return nil
}

// urzadzenieSesji oddaje kod urządzenia z sesji bramki żądania; pustka znaczy
// sesję bez zapisanego urządzenia i list pokazuje wtedy myślnik.
func (a *adapterUwierzytelnienia) urzadzenieSesji(ctx context.Context) string {
	skrot := sesjaBiezacaZKontekstu(ctx)
	if skrot == "" {
		return ""
	}
	sesja, err := a.repozytorium.SesjaBramkiPoSkrocie(ctx, skrot)
	if err != nil || sesja.UrzadzenieKod == nil {
		return ""
	}
	return *sesja.UrzadzenieKod
}

// zamknijPoNiepowodzeniuListu zdejmuje zamówienie, którego para listów nie
// wyszła w całości. Błąd zamknięcia nie przesłania błędu nadania.
func (a *adapterUwierzytelnienia) zamknijPoNiepowodzeniuListu(ctx context.Context, z dane.ZmianaAdresuKonta) {
	_ = a.zmianaAdresu.ZamknijZmiane(ctx, z.KontoId, z.Utworzono)
}

// PotwierdzZmianeAdresu obsługuje `auth.email.change.confirm`: przenosi konto na
// adres oczekujący i unieważnia drogę wycofania.
func (a *adapterUwierzytelnienia) PotwierdzZmianeAdresu(ctx context.Context,
	z shared.AuthEmailChangeConfirmRequest) (shared.AuthEmailChangeConfirmResponse, error) {

	if err := a.kontoGotowe(); err != nil {
		return shared.AuthEmailChangeConfirmResponse{}, err
	}
	kod := bezOdstepow(z.Code)
	if kod == "" {
		return shared.AuthEmailChangeConfirmResponse{}, bladBramki(shared.ErrorCodeValidationFailed,
			"Podaj kod z listu wysłanego na nowy adres.")
	}
	a.zamekZmiany.Lock()
	defer a.zamekZmiany.Unlock()

	kontoId, err := kontoWolajacego(ctx)
	if err != nil {
		return shared.AuthEmailChangeConfirmResponse{}, err
	}
	zamowienie, err := a.zmianaAdresu.ZmianaCzynnaKonta(ctx)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return shared.AuthEmailChangeConfirmResponse{}, bladBramki(shared.ErrorCodeNotFound,
			"Nie ma zamówionej zmiany adresu do potwierdzenia.")
	}
	if err != nil {
		return shared.AuthEmailChangeConfirmResponse{}, err
	}
	teraz := time.Now()
	if teraz.UnixMilli() > zamowienie.WygasaKod {
		if err := a.zmianaAdresu.ZamknijZmiane(ctx, zamowienie.KontoId, zamowienie.Utworzono); err != nil {
			return shared.AuthEmailChangeConfirmResponse{}, err
		}
		return shared.AuthEmailChangeConfirmResponse{}, bladBramki(shared.ErrorCodeValidationFailed,
			"Kod stracił ważność. Zamów zmianę adresu jeszcze raz.")
	}
	if skrotTokenu(kod) != zamowienie.SkrotKodu {
		proby, err := a.zmianaAdresu.OdnotujPomylke(ctx, zamowienie.KontoId, zamowienie.Utworzono)
		if err != nil {
			return shared.AuthEmailChangeConfirmResponse{}, err
		}
		// Kod ma sześć cyfr; bez sufitu pomyłek cały ten zakres przechodzi się
		// próbami, a zamówienie stoi otwarte godzinę.
		if proby >= probyZmianyAdresu {
			if err := a.zmianaAdresu.ZamknijZmiane(ctx, zamowienie.KontoId, zamowienie.Utworzono); err != nil {
				return shared.AuthEmailChangeConfirmResponse{}, err
			}
		}
		return shared.AuthEmailChangeConfirmResponse{}, bladBramki(shared.ErrorCodeValidationFailed,
			"Kod jest nieprawidłowy.")
	}
	if err := a.konto.UstawAdresKonta(ctx, kontoId, zamowienie.AdresNowy); err != nil {
		return shared.AuthEmailChangeConfirmResponse{}, err
	}
	// Zamknięcie zdejmuje razem z kodem drogę wycofania: zmiany domkniętej nie
	// wycofuje się listem, tylko kolejną zmianą adresu.
	if err := a.zmianaAdresu.ZamknijZmiane(ctx, zamowienie.KontoId, zamowienie.Utworzono); err != nil {
		return shared.AuthEmailChangeConfirmResponse{}, err
	}
	return shared.AuthEmailChangeConfirmResponse{Changed: true, Email: zamowienie.AdresNowy}, nil
}

// WycofajZmianeAdresu obsługuje `auth.email.change.revoke`: zdejmuje zamówienie
// drogą z listu ostrzegawczego, bez sesji.
func (a *adapterUwierzytelnienia) WycofajZmianeAdresu(ctx context.Context,
	z shared.AuthEmailChangeRevokeRequest) (shared.AuthEmailChangeRevokeResponse, error) {

	if err := a.kontoGotowe(); err != nil {
		return shared.AuthEmailChangeRevokeResponse{}, err
	}
	droga := strings.TrimSpace(z.Token)
	if droga == "" {
		return shared.AuthEmailChangeRevokeResponse{}, bladBramki(shared.ErrorCodeValidationFailed,
			"Podaj drogę wycofania z listu.")
	}
	a.zamekZmiany.Lock()
	defer a.zamekZmiany.Unlock()

	/* Droga przychodzi bez sesji — Operator, któremu przejęto konto, sesji mieć
	   nie musi. Wiersz odnajduje sam skrót drogi, a konto bierze się stamtąd,
	   nie z gniazda. */
	zamowienie, err := a.zmianaAdresu.ZmianaPoSkrocieWycofania(ctx, skrotTokenu(droga))
	if errors.Is(err, dane.ErrBrakWiersza) {
		return shared.AuthEmailChangeRevokeResponse{Revoked: false}, nil
	}
	if err != nil {
		return shared.AuthEmailChangeRevokeResponse{}, err
	}
	if zamowienie.Zamkniete || time.Now().UnixMilli() > zamowienie.WygasaWycofanie {
		return shared.AuthEmailChangeRevokeResponse{Revoked: false}, nil
	}
	if err := a.zmianaAdresu.ZamknijZmiane(ctx, zamowienie.KontoId, zamowienie.Utworzono); err != nil {
		return shared.AuthEmailChangeRevokeResponse{}, err
	}
	return shared.AuthEmailChangeRevokeResponse{Revoked: true}, nil
}

// kopiaWartosci oddaje osobną mapę: oba listy pary różnią się kilkoma pozycjami,
// a wspólnych nie wolno im nadpisywać nawzajem.
func kopiaWartosci(z map[string]string) map[string]string {
	kopia := make(map[string]string, len(z)+4)
	for k, v := range z {
		kopia[k] = v
	}
	return kopia
}

// wartoscAlboMyslnik zastępuje pustkę myślnikiem: list pokazuje wiersz zawsze,
// a puste miejsce w tabeli czyta się jak błąd składania.
func wartoscAlboMyslnik(v string) string {
	if strings.TrimSpace(v) == "" {
		return "—"
	}
	return v
}
