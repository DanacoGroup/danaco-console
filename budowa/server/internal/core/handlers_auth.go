// Plik wpina sześć komend rodziny auth.* — bramki Operatora; cała bramka stoi na jednym porcie, bo sześć komend dotyka dwóch tabel i jednego sejfu.
package core

import (
	"context"

	"danacoconsole/shared"
)

// Uwierzytelnianie jest portem rodziny `auth.*`. Mówi wyłącznie
// typami kontraktu; wiersze tabel `metoda_uwierzytelnienia` i `sesja_bramki`
// oraz wpisy sejfu poświadczeń leżą po drugiej stronie adaptera.
type Uwierzytelnianie interface {
	// ZalozBramke obsługuje `auth.register`.
	ZalozBramke(ctx context.Context, z shared.AuthRegisterRequest) (shared.AuthRegisterResponse, error)
	// PotwierdzAdres obsługuje `auth.verify` — zamyka rejestrację i wydaje
	// urządzeniu token dostępu.
	PotwierdzAdres(ctx context.Context, z shared.AuthVerifyRequest) (shared.AuthVerifyResponse, error)
	// RozpocznijOdzyskanie obsługuje `auth.recover`.
	RozpocznijOdzyskanie(ctx context.Context, z shared.AuthRecoverRequest) (shared.AuthRecoverResponse, error)
	// UstawNoweHaslo obsługuje `auth.reset`.
	UstawNoweHaslo(ctx context.Context, z shared.AuthResetRequest) (shared.AuthResetResponse, error)
	// WejdzPrzezBramke obsługuje `auth.login`.
	WejdzPrzezBramke(ctx context.Context, z shared.AuthLoginRequest) (shared.AuthLoginResponse, error)
	// ZalozMetodeWejscia obsługuje `auth.method.add`.
	ZalozMetodeWejscia(ctx context.Context, z shared.AuthMethodAddRequest) (shared.AuthMethodAddResponse, error)
	// ZdejmijMetodeWejscia obsługuje `auth.method.remove`.
	ZdejmijMetodeWejscia(ctx context.Context, z shared.AuthMethodRemoveRequest) (shared.AuthMethodRemoveResponse, error)
	// ZmienHasloBramki obsługuje `auth.password.reset`.
	ZmienHasloBramki(ctx context.Context, z shared.AuthPasswordResetRequest) (shared.AuthPasswordResetResponse, error)
	// PrzedluzSesjeBramki obsługuje `auth.token.refresh`.
	PrzedluzSesjeBramki(ctx context.Context, z shared.AuthTokenRefreshRequest) (shared.AuthTokenRefreshResponse, error)

	// MetodyWejscia oddaje wykaz metod po zmianie, potrzebny ładunkowi zdarzenia auth.changed.
	MetodyWejscia(ctx context.Context) ([]shared.AuthMethod, error)

	// BramkaZalozona mówi, czy sekret bramki istnieje, by klient odróżnił zaloguj się od ustaw hasło.
	BramkaZalozona(ctx context.Context) (bool, error)

	// RozpoznajSesjeBramki sprawdza token z powitania i oddaje skrót sesji; fałsz znaczy token nieznany.
	RozpoznajSesjeBramki(ctx context.Context, token string) (string, bool, error)
}

// Adapter wypełnia port w całości. Gdyby port i adapter się rozjechały,
// kompilacja stanie tutaj, a nie dopiero na martwej komendzie u Operatora.
var _ Uwierzytelnianie = (*adapterUwierzytelnienia)(nil)

// zarejestrujUwierzytelnianie wpina sześć komend rodziny `auth.*`. Port
// niewypełniony nie rejestruje niczego: komendy odpowiedzą wtedy
// `auth.unknown`, a pozostałe domeny pracują bez zmian.
func zarejestrujUwierzytelnianie(r *Rejestr, u Uwierzytelnianie, e *emiter, wiez *wiezBramki) {
	if r == nil || u == nil {
		return
	}
	// Konto sesji bramki oddaje ten sam adapter, który sesje wydaje; port
	// Uwierzytelnianie go nie wymienia, bo rozpoznanie konta służy adresowaniu
	// rozgłoszeń, nie żadnej z komend rodziny.
	konta, _ := u.(RozpoznanieKontaSesji)

	// Rejestracja połączenia nie wiąże i sesji nie zakłada: konto powstaje niepotwierdzone.
	r.Zarejestruj(shared.CommandAuthRegister,
		obsluz(func(ctx context.Context, z shared.AuthRegisterRequest) (shared.AuthRegisterResponse, error) {
			return u.ZalozBramke(ctx, z)
		}))

	// Potwierdzenie adresu wiąże połączenie od razu — pierwsze wejście Operatora do platformy.
	r.Zarejestruj(shared.CommandAuthVerify,
		obsluz(func(ctx context.Context, z shared.AuthVerifyRequest) (shared.AuthVerifyResponse, error) {
			odpowiedz, err := u.PotwierdzAdres(ctx, z)
			if err == nil {
				skrot := skrotTokenu(odpowiedz.Session.Token)
				wiez.Zwiaz(polaczenieZKontekstu(ctx), skrot)
				przypiszKontoGniazda(ctx, konta, skrot)
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandAuthRecover,
		obsluz(func(ctx context.Context, z shared.AuthRecoverRequest) (shared.AuthRecoverResponse, error) {
			return u.RozpocznijOdzyskanie(ctx, z)
		}))

	// Odzyskanie unieważnia tokeny wydane wcześniej, więc rozgłasza zmianę pozostałym urządzeniom.
	r.Zarejestruj(shared.CommandAuthReset,
		obsluz(func(ctx context.Context, z shared.AuthResetRequest) (shared.AuthResetResponse, error) {
			odpowiedz, err := u.UstawNoweHaslo(ctx, z)
			if err == nil {
				e.wyslij(shared.EventAuthChanged, "", shared.AuthChangedEvent{
					Reason: shared.AuthChangeReasonPasswordReset,
				})
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandAuthLogin,
		obsluz(func(ctx context.Context, z shared.AuthLoginRequest) (shared.AuthLoginResponse, error) {
			odpowiedz, err := u.WejdzPrzezBramke(ctx, z)
			if err == nil {
				skrot := skrotTokenu(odpowiedz.Session.Token)
				wiez.Zwiaz(polaczenieZKontekstu(ctx), skrot)
				przypiszKontoGniazda(ctx, konta, skrot)
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandAuthMethodAdd,
		obsluz(func(ctx context.Context, z shared.AuthMethodAddRequest) (shared.AuthMethodAddResponse, error) {
			odpowiedz, err := u.ZalozMetodeWejscia(ctx, z)
			if err == nil {
				rozglosZmianeBramki(e, shared.AuthChangeReasonMethodAdded,
					odpowiedz.Methods, &z.DeviceId)
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandAuthMethodRemove,
		obsluz(func(ctx context.Context, z shared.AuthMethodRemoveRequest) (shared.AuthMethodRemoveResponse, error) {
			odpowiedz, err := u.ZdejmijMetodeWejscia(ctx, z)
			// Bez zdjęcia nie ma zmiany, więc nie ma czego rozgłaszać.
			if err == nil && odpowiedz.Removed {
				rozglosZmianeBramki(e, shared.AuthChangeReasonMethodRemoved,
					odpowiedz.Methods, z.DeviceId)
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandAuthPasswordReset,
		obsluz(func(ctx context.Context, z shared.AuthPasswordResetRequest) (shared.AuthPasswordResetResponse, error) {
			// Sesja wołającego wchodzi kontekstem, żeby zmiana hasła nie wyrzuciła tego, kto ją wykonał, za drzwi.
			odpowiedz, err := u.ZmienHasloBramki(zSesjaBiezaca(ctx, wiez.SkrotKontekstu(ctx)), z)
			if err == nil && odpowiedz.Changed {
				// Wykaz metod jest brany osobno — odpowiedź tej komendy go nie niesie, hasło już jest zmienione.
				metody, _ := u.MetodyWejscia(ctx)
				rozglosZmianeBramki(e, shared.AuthChangeReasonPasswordChanged, metody, nil)
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandAuthTokenRefresh,
		obsluz(func(ctx context.Context, z shared.AuthTokenRefreshRequest) (shared.AuthTokenRefreshResponse, error) {
			// Więź połączenia wchodzi kontekstem tak samo jak przy zmianie hasła: puste Token znaczy bieżącą.
			return u.PrzedluzSesjeBramki(zSesjaBiezaca(ctx, wiez.SkrotKontekstu(ctx)), z)
		}))
}

// rozglosZmianeBramki wysyła `auth.changed`. Zdarzenie idzie bez sesji
// komunikatu: bramka nie należy do żadnej karty sesji — to stan platformy,
// a nie stan pracy. Nadajnik niepodłączony nie jest błędem.
func rozglosZmianeBramki(e *emiter, powod shared.AuthChangeReason,
	metody []shared.AuthMethod, urzadzenie *string) {

	if e == nil {
		return
	}
	e.wyslij(shared.EventAuthChanged, "", shared.AuthChangedEvent{
		Reason:   powod,
		Methods:  metody,
		DeviceId: niepustyTekst(urzadzenie),
	})
}
