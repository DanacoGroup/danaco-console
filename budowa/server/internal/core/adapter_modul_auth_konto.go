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
			"trwałość konta właściciela nie jest wpięta — konta nie ma gdzie zapisać")
	}
	return nil
}

// kontoPotwierdzone zamyka bramkę przed kontem, którego adresu nikt nie
// potwierdził. Brak trwałości konta i brak wiersza konta nie zamykają bramki,
// bo o potwierdzeniu nie ma wtedy co rozstrzygać.
func (a *adapterUwierzytelnienia) kontoPotwierdzone(ctx context.Context) error {
	if a.konto == nil {
		return nil
	}
	konto, err := a.konto.Konto(ctx)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return nil
	}
	if err != nil {
		return err
	}
	if konto.Potwierdzone {
		return nil
	}
	// Bramki nie zamyka potwierdzenie, gdy konto założono bez poczty — wejście idzie wtedy hasłem.
	if _, bezPoczty := a.znacznikBezPoczty(ctx); bezPoczty {
		return nil
	}
	return bladBramki(shared.ErrorCodeNotAuthenticated,
		"konto czeka na potwierdzenie adresu "+konto.Email+
			" — przepisz drogę potwierdzenia z listu komendą auth.verify; "+
			"do tego czasu bramka jest zamknięta, bo adres jest jedyną drogą odzyskania konta")
}

// daneRejestracji sprawdza login i adres podane przy rejestracji. Sprawdzenie
// adresu jest celowo płytkie, bo głębsza kontrola postaci adresu odrzuca
// adresy poprawne i tak nie rozstrzyga, czy skrzynka istnieje.
func daneRejestracji(z shared.AuthRegisterRequest) (string, string, error) {
	login := strings.TrimSpace(z.Login)
	if login == "" {
		return "", "", bladBramki(shared.ErrorCodeValidationFailed,
			"rejestracja bez loginu — login jest nazwą, którą Operator się loguje")
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
	cel, email, login string) error {

	// Brak konta nadawczego nazywa się przed zapisaniem drogi, żeby baza nie trzymała drogi bez listu.
	if err := a.kontoNadawcze(ctx).Brak(); err != nil {
		return bladBramki(shared.ErrorCodeInternalError,
			err.Error()+"; "+dwieDrogiKontaNadawczego)
	}
	droga, err := nowyTokenBramki()
	if err != nil {
		return err
	}
	teraz := time.Now()
	if err := a.konto.ZalozPotwierdzenie(ctx, dane.PotwierdzenieTozsamosci{
		Skrot:     skrotTokenu(droga),
		Cel:       cel,
		Wygasa:    teraz.Add(trwanieDrogiPotwierdzenia).UnixMilli(),
		Utworzono: teraz.UnixMilli(),
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
	if cel == dane.CelOdzyskanie {
		return nadajnik.List{
			Do:    email,
			Temat: "Danaco Console — odzyskanie konta",
			Tresc: "Ktoś poprosił o ustawienie nowego hasła do konta " + login + ".\n\n" +
				"Droga potwierdzenia:\n\n    " + droga + "\n\n" +
				"Wpisz ją w oknie odzyskiwania konta, aby ustawić nowe hasło. " +
				"Droga jest jednorazowa i wygasa po godzinie.\n\n" +
				"Po ustawieniu nowego hasła wszystkie urządzenia zalogują się ponownie.\n\n" +
				"Jeżeli to nie Ty prosiłeś o zmianę — nie rób nic. " +
				"Bez tej drogi hasło pozostaje bez zmian.\n",
		}
	}
	return nadajnik.List{
		Do:    email,
		Temat: "Danaco Console — potwierdzenie adresu",
		Tresc: "Konto " + login + " zostało założone i czeka na potwierdzenie tego adresu.\n\n" +
			"Droga potwierdzenia:\n\n    " + droga + "\n\n" +
			"Wpisz ją w oknie rejestracji, aby zakończyć zakładanie konta i wejść do platformy. " +
			"Droga jest jednorazowa i wygasa po godzinie.\n\n" +
			"Ten adres będzie później jedyną drogą odzyskania konta.\n",
	}
}

// zuzyjDroge sprawdza drogę i zamyka ją w jednej czynności. Sprawdzenie
// i zamknięcie są niepodzielne warunkiem w bazie, więc dwa żądania z tym
// samym materiałem nie zastają obie drogi ważnej.
func (a *adapterUwierzytelnienia) zuzyjDroge(ctx context.Context, cel, droga string) error {
	if strings.TrimSpace(droga) == "" {
		return bladBramki(shared.ErrorCodeValidationFailed, "droga potwierdzenia jest pusta")
	}
	skrot := skrotTokenu(droga)
	zapis, err := a.konto.PotwierdzeniePoSkrocie(ctx, skrot)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return bladBramki(shared.ErrorCodeNotAuthenticated,
			"droga potwierdzenia nie jest znana platformie")
	}
	if err != nil {
		return err
	}
	if zapis.Cel != cel {
		return bladBramki(shared.ErrorCodeNotAuthenticated,
			"droga potwierdzenia została wydana do innej czynności")
	}
	zamknieta, err := a.konto.ZuzyjPotwierdzenie(ctx, skrot, time.Now().UnixMilli())
	if err != nil {
		return err
	}
	if !zamknieta {
		return bladBramki(shared.ErrorCodeNotAuthenticated,
			"droga potwierdzenia jest już zużyta albo wygasła — poproś o nową")
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
			"konta właściciela jeszcze nie ma — najpierw rejestracja")
	}
	if err != nil {
		return shared.AuthVerifyResponse{}, err
	}
	if konto.Potwierdzone {
		return shared.AuthVerifyResponse{}, bladBramki(shared.ErrorCodeConflict,
			"adres jest już potwierdzony — wejście idzie komendą auth.login")
	}
	if err := a.zuzyjDroge(ctx, dane.CelWeryfikacja, z.Token); err != nil {
		return shared.AuthVerifyResponse{}, err
	}
	if err := a.konto.PotwierdzKonto(ctx); err != nil {
		return shared.AuthVerifyResponse{}, err
	}
	// Znacznik bramki bez poczty przestał być prawdą — bramkę trzyma odtąd sam wiersz konta.
	a.zdejmijZnacznikBezPoczty(ctx)
	sesja, err := a.zalozSesje(ctx, shared.AuthMethodKindPassword,
		niepustyTekst(z.DeviceId), wartoscPrawdy(z.KeepSignedIn))
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
			"odzyskanie konta bez podanego adresu")
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
	if err := a.wyslijDrogePotwierdzenia(ctx, dane.CelOdzyskanie, konto.Email, konto.Login); err != nil {
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
			"ustawienie nowego hasła bez hasła")
	}
	a.zamekZmiany.Lock()
	defer a.zamekZmiany.Unlock()

	kotwica, err := a.kotwica(ctx)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return shared.AuthResetResponse{}, bladBramki(shared.ErrorCodeConflict,
			"konta właściciela jeszcze nie ma — najpierw rejestracja")
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
