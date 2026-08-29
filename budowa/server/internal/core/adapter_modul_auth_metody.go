// Cztery czynności bramki wykonywane po wejściu: założenie i zdjęcie metody
// szybkiego wejścia, zmiana hasła i przedłużenie sesji. Założenie bramki
// i wejście leżą w `adapter_modul_auth.go`. Hasło jest kotwicą bramki
// w całym pliku.
package core

import (
	"context"
	"errors"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// ── auth.method.add ──────────────────────────────────────────────────────────

// ZalozMetodeWejscia zakłada PIN na wskazanym urządzeniu. Trzy przypadki
// kończą się odmową, każdy z własnym powodem: rodzaj `password`, rodzaj
// `hello` i PIN już założony na urządzeniu.
func (a *adapterUwierzytelnienia) ZalozMetodeWejscia(ctx context.Context,
	z shared.AuthMethodAddRequest) (shared.AuthMethodAddResponse, error) {

	if err := a.gotowa(); err != nil {
		return shared.AuthMethodAddResponse{}, err
	}
	// Sprawdzenie „wolno” i założenie idą pod jednym zamkiem, by dwa żądania
	// nie przeszły oba sprawdzenia.
	a.zamekZmiany.Lock()
	defer a.zamekZmiany.Unlock()
	if err := a.wolnoZalozyc(ctx, z); err != nil {
		return shared.AuthMethodAddResponse{}, err
	}
	_, err := a.zalozMetode(ctx, dane.MetodaUwierzytelnienia{
		Rodzaj:          shared.AuthMethodKindPin,
		Etykieta:        niepustyTekst(z.Label),
		UrzadzenieKod:   &z.DeviceId,
		NazwaUrzadzenia: niepustyTekst(z.DeviceName),
		Utworzono:       time.Now().UnixMilli(),
	}, *z.Secret)
	if err != nil {
		return shared.AuthMethodAddResponse{}, err
	}
	metody, err := a.MetodyWejscia(ctx)
	if err != nil {
		return shared.AuthMethodAddResponse{}, err
	}
	return shared.AuthMethodAddResponse{Methods: metody}, nil
}

// wolnoZalozyc rozstrzyga, czy żądanie założenia metody w ogóle ma prawo dojść
// do skutku. Wydzielone, żeby sama czynność została czynnością, a nie ciągiem
// warunków.
func (a *adapterUwierzytelnienia) wolnoZalozyc(ctx context.Context,
	z shared.AuthMethodAddRequest) error {

	switch z.Kind {
	case shared.AuthMethodKindPin:
	case shared.AuthMethodKindHello:
		return bladBramki(shared.ErrorCodeValidationFailed, odmowaHello)
	case shared.AuthMethodKindPassword:
		return bladBramki(shared.ErrorCodeValidationFailed,
			"Pierwsze hasło ustawia się podczas rejestracji.")
	default:
		return bladBramki(shared.ErrorCodeValidationFailed,
			"rodzaj metody "+string(z.Kind)+" nie należy do kontraktu")
	}
	if z.DeviceId == "" {
		return bladBramki(shared.ErrorCodeValidationFailed,
			"Wskaż urządzenie — kod PIN obowiązuje na jednym urządzeniu.")
	}
	if z.Secret == nil || *z.Secret == "" {
		return bladBramki(shared.ErrorCodeValidationFailed, "Podaj kod PIN.")
	}
	// Metoda szybkiego wejścia bez kotwicy zdjęta z urządzeniem zostawiłaby
	// bramkę bez hasła.
	if _, err := a.kotwica(ctx); errors.Is(err, dane.ErrBrakWiersza) {
		return bladBramki(shared.ErrorCodeConflict,
			"Kod PIN zakłada się dopiero po ustawieniu hasła.")
	} else if err != nil {
		return err
	}
	_, err := a.repozytorium.MetodaUrzadzenia(ctx, shared.AuthMethodKindPin, z.DeviceId)
	if err == nil {
		return bladBramki(shared.ErrorCodeConflict,
			"urządzenie "+z.DeviceId+" ma już PIN; zdejmij go komendą auth.method.remove, "+
				"zanim założysz nowy")
	}
	if !errors.Is(err, dane.ErrBrakWiersza) {
		return err
	}
	return nil
}

// ── auth.method.remove ───────────────────────────────────────────────────────

// ZdejmijMetodeWejscia zdejmuje metodę szybkiego wejścia z urządzenia.
// Odpowiedź niesie `removed`, ale nigdy w postaci „false, bo nie było czego
// zdjąć" — brak metody jest odmową `not_found` z jej nazwą.
func (a *adapterUwierzytelnienia) ZdejmijMetodeWejscia(ctx context.Context,
	z shared.AuthMethodRemoveRequest) (shared.AuthMethodRemoveResponse, error) {

	if err := a.gotowa(); err != nil {
		return shared.AuthMethodRemoveResponse{}, err
	}
	if z.MethodId == "" {
		return shared.AuthMethodRemoveResponse{}, bladBramki(shared.ErrorCodeValidationFailed,
			"Wskaż metodę logowania, którą chcesz usunąć.")
	}
	metoda, err := a.repozytorium.MetodaPoKodzie(ctx, z.MethodId)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return shared.AuthMethodRemoveResponse{}, bladBramki(shared.ErrorCodeNotFound,
			"metoda wejścia "+z.MethodId+" nie istnieje")
	}
	if err != nil {
		return shared.AuthMethodRemoveResponse{}, err
	}
	if err := a.wolnoZdjac(ctx, metoda, z); err != nil {
		return shared.AuthMethodRemoveResponse{}, err
	}
	zdjeta, err := a.repozytorium.UsunMetode(ctx, metoda.Kod)
	if err != nil {
		return shared.AuthMethodRemoveResponse{}, err
	}
	if !zdjeta {
		return shared.AuthMethodRemoveResponse{}, bladBramki(shared.ErrorCodeNotFound,
			"metoda wejścia "+z.MethodId+" nie istnieje")
	}
	// Sekret bez właściciela byłby śmieciem w sejfie; sprzątanie nie cofa
	// zdjęcia metody.
	usunPoswiadczenie(ctx, a.sejf, przedrostekBytuSejfu+metoda.Kod)
	metody, err := a.MetodyWejscia(ctx)
	if err != nil {
		return shared.AuthMethodRemoveResponse{}, err
	}
	return shared.AuthMethodRemoveResponse{Methods: metody, Removed: true}, nil
}

// wolnoZdjac broni dwóch zdań kontraktu: kotwicy zdjąć się nie da i ostatniej
// metody zdjąć się nie da.
func (a *adapterUwierzytelnienia) wolnoZdjac(ctx context.Context,
	metoda dane.MetodaUwierzytelnienia, z shared.AuthMethodRemoveRequest) error {

	if metoda.Kotwica {
		return bladBramki(shared.ErrorCodeConflict,
			"metoda wejścia "+metoda.Kod+" jest kotwicą bramki; hasła zdjąć się nie da")
	}
	wskazane := wartoscTekstu(z.DeviceId)
	if wskazane != "" && (metoda.UrzadzenieKod == nil || *metoda.UrzadzenieKod != wskazane) {
		return bladBramki(shared.ErrorCodeConflict,
			"metoda wejścia "+metoda.Kod+" nie należy do urządzenia "+wskazane)
	}
	metody, err := a.repozytorium.Metody(ctx)
	if err != nil {
		return err
	}
	if len(metody) <= 1 {
		return bladBramki(shared.ErrorCodeConflict,
			"metoda wejścia "+metoda.Kod+" jest ostatnią metodą bramki; ostatniej zdjąć się nie da")
	}
	return nil
}

// ── auth.password.reset ──────────────────────────────────────────────────────

// ZmienHasloBramki podmienia hasło ze znajomością hasła bieżącego. Zmiana
// hasła unieważnia sesje bramki poza bieżącą; `revokedSessions` liczy to,
// co naprawdę unieważniono.
func (a *adapterUwierzytelnienia) ZmienHasloBramki(ctx context.Context,
	z shared.AuthPasswordResetRequest) (shared.AuthPasswordResetResponse, error) {

	if err := a.gotowa(); err != nil {
		return shared.AuthPasswordResetResponse{}, err
	}
	// Sprawdzenie hasła bieżącego i podmiana sekretu idą pod jednym zamkiem jako
	// jedna czynność.
	a.zamekZmiany.Lock()
	defer a.zamekZmiany.Unlock()
	if z.CurrentPassword == "" || z.NewPassword == "" {
		return shared.AuthPasswordResetResponse{}, bladBramki(shared.ErrorCodeValidationFailed,
			"Podaj hasło bieżące oraz nowe.")
	}
	kotwica, err := a.kotwica(ctx)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return shared.AuthPasswordResetResponse{}, bladBramki(shared.ErrorCodeNotFound,
			"Hasło dostępu nie zostało jeszcze ustawione.")
	}
	if err != nil {
		return shared.AuthPasswordResetResponse{}, err
	}
	zgadza, err := a.sekretZgadzaSieZWpisem(ctx, kotwica, z.CurrentPassword)
	if err != nil {
		return shared.AuthPasswordResetResponse{}, err
	}
	if !zgadza {
		return shared.AuthPasswordResetResponse{}, bladBramki(shared.ErrorCodeValidationFailed,
			"Podane hasło bieżące jest nieprawidłowe.")
	}
	if err := a.podmienSekret(ctx, kotwica, z.NewPassword); err != nil {
		return shared.AuthPasswordResetResponse{}, err
	}
	uniewaznione, err := a.repozytorium.UniewaznijSesjeBramkiPoza(ctx,
		sesjaBiezacaZKontekstu(ctx), time.Now().UnixMilli())
	if err != nil {
		return shared.AuthPasswordResetResponse{}, err
	}
	liczba := uniewaznione
	return shared.AuthPasswordResetResponse{Changed: true, RevokedSessions: &liczba}, nil
}

// podmienSekret kładzie nowy zapis pod tym samym bytem sejfu i odświeża
// odwołanie w bazie: sejf nadpisuje wpis bytu, a zapis do bazy jest tu
// asekuracją.
func (a *adapterUwierzytelnienia) podmienSekret(ctx context.Context,
	metoda dane.MetodaUwierzytelnienia, sekret string) error {

	zapis, err := zapisSekretu(sekret)
	if err != nil {
		return bladBramki(shared.ErrorCodeInternalError, err.Error())
	}
	odwolanie, err := a.sejf.Zapisz(ctx, przedrostekBytuSejfu+metoda.Kod, zapis)
	if err != nil {
		return err
	}
	if odwolanie == metoda.OdwolanieSekretu {
		return nil
	}
	return a.repozytorium.ZapiszOdwolanieSekretu(ctx, metoda.Kod, odwolanie)
}

// ── auth.token.refresh ───────────────────────────────────────────────────────

// PrzedluzSesjeBramki przesuwa wygaśnięcie sesji wskazanej tokenem. Pusty
// token znaczy sesję bieżącego połączenia, zgodnie z kontraktem
// `AuthTokenRefreshRequest.Token`.
func (a *adapterUwierzytelnienia) PrzedluzSesjeBramki(ctx context.Context,
	z shared.AuthTokenRefreshRequest) (shared.AuthTokenRefreshResponse, error) {

	if err := a.gotowa(); err != nil {
		return shared.AuthTokenRefreshResponse{}, err
	}
	token := wartoscTekstu(z.Token)
	skrot := skrotTokenu(token)
	if token == "" {
		// Skrót przychodzi kontekstem z więzi gniazda, tak samo jak przy zmianie
		// hasła.
		skrot = sesjaBiezacaZKontekstu(ctx)
		if skrot == "" {
			return shared.AuthTokenRefreshResponse{}, bladBramki(shared.ErrorCodeValidationFailed,
				"Połączenie nie jest zalogowane. "+
					"token wchodzi do serwera powitaniem connection.hello albo wejściem auth.login — "+
					"po jednym z nich pole tokenu wolno zostawić puste")
		}
	}
	sesja, err := a.repozytorium.SesjaBramkiPoSkrocie(ctx, skrot)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return shared.AuthTokenRefreshResponse{}, bladBramki(shared.ErrorCodeNotFound,
			"Ta sesja już nie istnieje — zaloguj się ponownie.")
	}
	if err != nil {
		return shared.AuthTokenRefreshResponse{}, err
	}
	teraz := time.Now()
	if err := sesjaNadaje(sesja, teraz.UnixMilli()); err != nil {
		return shared.AuthTokenRefreshResponse{}, err
	}
	przedluzona, err := a.repozytorium.PrzedluzSesjeBramki(ctx, skrot,
		teraz.Add(trwanieSesji(sesja)).UnixMilli())
	if err != nil {
		return shared.AuthTokenRefreshResponse{}, err
	}
	// Token wraca taki, jaki przyszedł; przy przedłużeniu przez więź wraca
	// pusty.
	return shared.AuthTokenRefreshResponse{Session: sesjaBramkiKontraktu(token, przedluzona)}, nil
}

// trwanieSesji oddaje długość, o którą przesuwa się wygaśnięcie tej sesji.
// Tu działa opcja „nie wyloguj mnie”: przełącznik rozstrzyga przy zakładaniu
// sesji.
func trwanieSesji(sesja dane.SesjaBramki) time.Duration {
	if sesja.Trwanie <= 0 {
		return trwanieSesjiBramki
	}
	return time.Duration(sesja.Trwanie) * time.Millisecond
}

// sesjaNadaje mówi, czy sesję wolno jeszcze przedłużyć. Sesja unieważniona
// i wygasła są dwoma różnymi stanami i dwa różne zdania o nich wracają —
// Operator ma wiedzieć, czy ktoś ją zamknął, czy po prostu minął czas.
func sesjaNadaje(sesja dane.SesjaBramki, teraz int64) error {
	if sesja.Uniewazniono != nil {
		return bladBramki(shared.ErrorCodeConflict,
			"Sesja została zakończona — zaloguj się ponownie.")
	}
	if sesja.Wygasa <= teraz {
		return bladBramki(shared.ErrorCodeConflict,
			"Sesja wygasła — zaloguj się ponownie.")
	}
	return nil
}
