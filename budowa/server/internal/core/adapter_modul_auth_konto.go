// Odpowiedzialność pliku: tożsamość właściciela i dwie drogi, które prowadzą
// przez jego skrzynkę — potwierdzenie adresu po rejestracji (auth.verify)
// oraz odzyskanie konta (auth.recover, auth.reset). Droga zawsze idzie
// listem na adres.
package core

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/nadajnik"
	"danacoconsole/shared"

	"danacoconsole/server/internal/core/listy"
)

// trwanieDrogiPotwierdzenia — jak długo ważny jest materiał wysłany listem.
// Godzina jest kompromisem między skrzynką sprawdzaną rzadko a materiałem
// ważnym bez końca; droga wygasła nie zamyka niczego trwale, Operator prosi
// o nową.
const trwanieDrogiPotwierdzenia = time.Hour

// kontoGotowe odmawia, gdy montaż nie wpiął trwałości tożsamości.
//
// Odmowa jest wprost, nie cicha: rejestracja bez zapisu konta wyglądałaby jak
// udana, a Operator zostałby przed platformą, która o nim nie wie.
func (a *adapterUwierzytelnienia) kontoGotowe() error {
	if a.konto == nil {
		return bladBramki(shared.ErrorCodeInternalError,
			"Nie można zapisać konta — magazyn danych jest niedostępny.")
	}
	return nil
}

// kontoPotwierdzone zamyka bramkę przed kontem, którego adresu nikt nie
// potwierdził. Brak trwałości konta i brak wiersza konta nie zamykają bramki,
// bo o potwierdzeniu nie ma wtedy co rozstrzygać.
func (a *adapterUwierzytelnienia) kontoPotwierdzone(ctx context.Context,
	konto dane.KontoWlasciciela) error {

	// Konto puste znaczy platformę bez trwałości konta albo przed rejestracją —
	// o potwierdzeniu nie ma wtedy co rozstrzygać.
	if a.konto == nil || konto.Id == 0 {
		return nil
	}
	if konto.Potwierdzone {
		return nil
	}
	// Bramki nie zamyka potwierdzenie, gdy konto założono bez poczty — wejście idzie wtedy hasłem.
	if _, bezPoczty := a.znacznikBezPoczty(ctx); bezPoczty {
		return nil
	}
	return bladBramki(shared.ErrorCodeNotAuthenticated,
		"Konto oczekuje na potwierdzenie adresu "+konto.Email+
			". Wprowadź kod potwierdzający z wiadomości — do tego czasu wejście "+
			"jest zamknięte, bo adresem odzyskuje się konto po utracie hasła.")
}

// daneRejestracji sprawdza login i adres podane przy rejestracji. Sprawdzenie
// adresu jest celowo płytkie, bo głębsza kontrola postaci adresu odrzuca
// adresy poprawne i tak nie rozstrzyga, czy skrzynka istnieje.
func daneRejestracji(z shared.AuthRegisterRequest) (string, string, error) {
	login := strings.TrimSpace(z.Login)
	if login == "" {
		return "", "", bladBramki(shared.ErrorCodeValidationFailed,
			"Podaj login — to nazwa, pod którą będziesz się logować.")
	}
	email := strings.TrimSpace(z.Email)
	malpa := strings.LastIndex(email, "@")
	if malpa <= 0 || !strings.Contains(email[malpa+1:], ".") {
		return "", "", bladBramki(shared.ErrorCodeValidationFailed,
			"adres e-mail uwierzytelniający nie ma postaci adresu — "+
				"jest jedyną drogą odzyskania konta, więc musi być adresem działającym")
	}
	return login, email, nil
}

// wyslijDrogePotwierdzenia losuje materiał, zapisuje jego skrót i nadaje list.
// Kolejność jest wiążąca: najpierw zapis skrótu, potem nadanie, inaczej list
// mógłby dojść, zanim droga stałaby się ważna.
func (a *adapterUwierzytelnienia) wyslijDrogePotwierdzenia(ctx context.Context,
	cel, email, login string, kontoId int64) error {

	// Brak konta nadawczego nazywa się przed zapisaniem drogi, żeby baza nie trzymała drogi bez listu.
	if err := a.kontoNadawcze(ctx).Brak(); err != nil {
		return bladBramki(shared.ErrorCodeInternalError,
			err.Error()+"; "+dwieDrogiKontaNadawczego)
	}
	/* Kod, nie token sesji: okno przyjmuje go w sześciu polach, a Operator
	   przepisuje go z wiadomości ręcznie. */
	droga, err := nowyKodPotwierdzenia()
	if err != nil {
		return err
	}
	teraz := time.Now()
	if err := a.konto.ZalozPotwierdzenie(ctx, dane.PotwierdzenieTozsamosci{
		Skrot:     skrotTokenu(droga),
		Cel:       cel,
		Wygasa:    teraz.Add(trwanieDrogiPotwierdzenia).UnixMilli(),
		Utworzono: teraz.UnixMilli(),
		// Konto, do którego droga prowadzi; przy dwóch kontach naraz to jedyne,
		// co rozstrzyga, które z nich potwierdza przepisany kod.
		KontoId: kontoId,
	}); err != nil {
		return err
	}
	if _, err := nadajnik.Wyslij(a.kontoNadawcze(ctx), listPotwierdzenia(cel, login, droga, email)); err != nil {
		return bladBramki(shared.ErrorCodeInternalError,
			fmt.Sprintf("nie udało się wysłać listu na %s: %v; %s", email, err, dwieDrogiKontaNadawczego))
	}
	return nil
}

// listPotwierdzenia składa treść jednego z dwóch listów systemowych. Treść
// jest zwięzła i mówi wprost, co się stało i co zrobić, bo rozwlekły list
// systemowy nakłania do zignorowania go.
func listPotwierdzenia(cel, login, droga, email string) nadajnik.List {
	/* Nagłówki wymagane przez opracowanie poczty transakcyjnej (rozdz. 5.3):
	   list z kodem jest wytworem programu, nie rozmową — bez tych nagłówków
	   autorespondery po drugiej stronie odpisują na niego w kółko. */
	naglowkiTransakcyjne := []string{
		"Auto-Submitted: auto-generated",
		"X-Auto-Response-Suppress: All",
	}
	teraz := time.Now()
	minutyWaznosci := int(trwanieDrogiPotwierdzenia.Minutes())

	if cel == dane.CelOdzyskanie {
		return nadajnik.List{
			Do:       email,
			Temat:    "Danaco Console — kod odzyskania konta",
			Naglowki: naglowkiTransakcyjne,
			Znak:     listy.ZnakMarki,
			IdZnaku:  listy.IdZnaku,
			TrescHtml: listy.ZlozLogowanie(listy.Logowanie{
				Odbiorca:      email,
				Kod:           droga,
				WaznoscMinuty: minutyWaznosci,
				Zadano:        teraz,
			}),
			Tresc: "Otrzymaliśmy prośbę o ustawienie nowego hasła do konta " + login + ".\n\n" +
				"Kod potwierdzający:\n\n    " + droga + "\n\n" +
				"Wprowadź go w oknie odzyskiwania dostępu, aby ustawić nowe hasło. " +
				"Kod jest jednorazowy i zachowuje ważność przez godzinę.\n\n" +
				"Po ustawieniu nowego hasła wszystkie urządzenia będą wymagały ponownego zalogowania.\n\n" +
				"Jeżeli prośba nie pochodzi od Ciebie, nie podejmuj żadnych czynności. " +
				"Bez tego kodu hasło pozostaje bez zmian.\n",
		}
	}
	return nadajnik.List{
		Do:       email,
		Temat:    "Danaco Console — kod potwierdzający adres",
		Naglowki: naglowkiTransakcyjne,
		Znak:     listy.ZnakMarki,
		IdZnaku:  listy.IdZnaku,
		TrescHtml: listy.ZlozAktywacje(listy.Aktywacja{
			Odbiorca:      email,
			Kod:           droga,
			WaznoscMinuty: minutyWaznosci,
			Zadano:        teraz,
		}),
		Tresc: "Konto " + login + " zostało założone i oczekuje na potwierdzenie tego adresu.\n\n" +
			"Kod potwierdzający:\n\n    " + droga + "\n\n" +
			"Wprowadź go w oknie rejestracji, aby zakończyć zakładanie konta i wejść do platformy. " +
			"Kod jest jednorazowy i zachowuje ważność przez godzinę.\n\n" +
			"Tym adresem odzyskasz konto, jeżeli zapomnisz hasła.\n",
	}
}

// zuzyjDroge sprawdza drogę i zamyka ją w jednej czynności. Sprawdzenie
// i zamknięcie są niepodzielne warunkiem w bazie, więc dwa żądania z tym
// samym materiałem nie zastają obie drogi ważnej.
func (a *adapterUwierzytelnienia) zuzyjDroge(ctx context.Context, cel, droga string) error {
	if strings.TrimSpace(droga) == "" {
		return bladBramki(shared.ErrorCodeValidationFailed, "Podaj kod potwierdzający.")
	}
	skrot := skrotTokenu(droga)
	zapis, err := a.konto.PotwierdzeniePoSkrocie(ctx, skrot)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return bladBramki(shared.ErrorCodeNotAuthenticated,
			"Ten kod potwierdzający nie jest znany.")
	}
	if err != nil {
		return err
	}
	if zapis.Cel != cel {
		return bladBramki(shared.ErrorCodeNotAuthenticated,
			"Ten kod potwierdzający dotyczy innej czynności.")
	}
	zamknieta, err := a.konto.ZuzyjPotwierdzenie(ctx, skrot, time.Now().UnixMilli())
	if err != nil {
		return err
	}
	if !zamknieta {
		return bladBramki(shared.ErrorCodeNotAuthenticated,
			"Ten kod potwierdzający został już użyty lub wygasł. Poproś o nowy.")
	}
	return nil
}

// ── auth.verify ──────────────────────────────────────────────────────────────

// PotwierdzAdres zamyka rejestrację: przenosi konto do stanu potwierdzonego
// i wydaje urządzeniu token dostępu.
//
// To jest moment pierwszego wejścia Operatora do platformy — dlatego sesja
// powstaje tutaj, a nie przy `auth.register`.
func (a *adapterUwierzytelnienia) PotwierdzAdres(ctx context.Context,
	z shared.AuthVerifyRequest) (shared.AuthVerifyResponse, error) {

	if err := a.kontoGotowe(); err != nil {
		return shared.AuthVerifyResponse{}, err
	}
	a.zamekZmiany.Lock()
	defer a.zamekZmiany.Unlock()

	konto, err := a.konto.Konto(ctx)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return shared.AuthVerifyResponse{}, bladBramki(shared.ErrorCodeConflict,
			"Konto Operatora nie zostało jeszcze założone. Zarejestruj się.")
	}
	if err != nil {
		return shared.AuthVerifyResponse{}, err
	}
	if konto.Potwierdzone {
		return shared.AuthVerifyResponse{}, bladBramki(shared.ErrorCodeConflict,
			"Adres jest już potwierdzony. Zaloguj się.")
	}
	if err := a.zuzyjDroge(ctx, dane.CelWeryfikacja, z.Token); err != nil {
		return shared.AuthVerifyResponse{}, err
	}
	if err := a.konto.PotwierdzKonto(ctx, konto.Id); err != nil {
		return shared.AuthVerifyResponse{}, err
	}
	// Znacznik bramki bez poczty przestał być prawdą — bramkę trzyma odtąd sam wiersz konta.
	a.zdejmijZnacznikBezPoczty(ctx)
	sesja, err := a.zalozSesje(ctx, shared.AuthMethodKindPassword,
		niepustyTekst(z.DeviceId), wartoscPrawdy(z.KeepSignedIn), konto.Id)
	if err != nil {
		return shared.AuthVerifyResponse{}, err
	}
	return shared.AuthVerifyResponse{Verified: true, Session: sesja}, nil
}

// ── auth.recover ─────────────────────────────────────────────────────────────

// RozpocznijOdzyskanie wysyła drogę potwierdzenia na adres uwierzytelniający.
// Odpowiedź jest zawsze taka sama, niezależnie od tego, czy adres pasuje do
// konta, dla właściciela czynność jest wykonana i list wychodzi.
func (a *adapterUwierzytelnienia) RozpocznijOdzyskanie(ctx context.Context,
	z shared.AuthRecoverRequest) (shared.AuthRecoverResponse, error) {

	if err := a.kontoGotowe(); err != nil {
		return shared.AuthRecoverResponse{}, err
	}
	email := strings.TrimSpace(z.Email)
	if email == "" {
		return shared.AuthRecoverResponse{}, bladBramki(shared.ErrorCodeValidationFailed,
			"Podaj adres e-mail konta.")
	}
	konto, err := a.konto.Konto(ctx)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return shared.AuthRecoverResponse{Sent: true}, nil
	}
	if err != nil {
		return shared.AuthRecoverResponse{}, err
	}
	if !strings.EqualFold(konto.Email, email) {
		return shared.AuthRecoverResponse{Sent: true}, nil
	}
	if err := a.wyslijDrogePotwierdzenia(ctx, dane.CelOdzyskanie, konto.Email, konto.Login, konto.Id); err != nil {
		return shared.AuthRecoverResponse{}, err
	}
	return shared.AuthRecoverResponse{Sent: true}, nil
}

// ── auth.reset ───────────────────────────────────────────────────────────────

// UstawNoweHaslo zamyka odzyskanie konta: podmienia hasło i unieważnia tokeny
// wydane przed zmianą. Unieważnienie jest częścią czynności, bo dostęp mógł
// mieć ktoś jeszcze.
func (a *adapterUwierzytelnienia) UstawNoweHaslo(ctx context.Context,
	z shared.AuthResetRequest) (shared.AuthResetResponse, error) {

	if err := a.gotowa(); err != nil {
		return shared.AuthResetResponse{}, err
	}
	if err := a.kontoGotowe(); err != nil {
		return shared.AuthResetResponse{}, err
	}
	if strings.TrimSpace(z.NewPassword) == "" {
		return shared.AuthResetResponse{}, bladBramki(shared.ErrorCodeValidationFailed,
			"Podaj nowe hasło.")
	}
	a.zamekZmiany.Lock()
	defer a.zamekZmiany.Unlock()

	kotwica, err := a.kotwica(ctx)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return shared.AuthResetResponse{}, bladBramki(shared.ErrorCodeConflict,
			"Konto Operatora nie zostało jeszcze założone. Zarejestruj się.")
	}
	if err != nil {
		return shared.AuthResetResponse{}, err
	}
	if err := a.zuzyjDroge(ctx, dane.CelOdzyskanie, z.Token); err != nil {
		return shared.AuthResetResponse{}, err
	}

	zapis, err := zapisSekretu(z.NewPassword)
	if err != nil {
		return shared.AuthResetResponse{}, err
	}
	odwolanie, err := a.sejf.Zapisz(ctx, przedrostekBytuSejfu+kotwica.Kod, zapis)
	if err != nil {
		return shared.AuthResetResponse{}, err
	}
	if err := a.repozytorium.ZapiszOdwolanieSekretu(ctx, kotwica.Kod, odwolanie); err != nil {
		return shared.AuthResetResponse{}, err
	}
	// Pusty skrót zachowany znaczy unieważnij wszystkie — odzyskanie idzie z urządzenia bez sesji.
	uniewaznione, err := a.repozytorium.UniewaznijSesjeBramkiPoza(ctx, "", time.Now().UnixMilli())
	if err != nil {
		return shared.AuthResetResponse{}, err
	}
	return shared.AuthResetResponse{Changed: true, RevokedDevices: uniewaznione}, nil
}
