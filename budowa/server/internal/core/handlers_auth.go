// Plik wpina sześć komend rodziny auth.* — bramki Operatora; cała bramka stoi na jednym porcie, bo sześć komend dotyka dwóch tabel i jednego sejfu.
package core

import (
	"context"

	"danacoconsole/server/internal/dane"
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
	// RozpocznijZmianeAdresu obsługuje `auth.email.change.start`.
	RozpocznijZmianeAdresu(ctx context.Context, z shared.AuthEmailChangeStartRequest) (shared.AuthEmailChangeStartResponse, error)
	// PotwierdzZmianeAdresu obsługuje `auth.email.change.confirm`.
	PotwierdzZmianeAdresu(ctx context.Context, z shared.AuthEmailChangeConfirmRequest) (shared.AuthEmailChangeConfirmResponse, error)
	// WycofajZmianeAdresu obsługuje `auth.email.change.revoke`.
	WycofajZmianeAdresu(ctx context.Context, z shared.AuthEmailChangeRevokeRequest) (shared.AuthEmailChangeRevokeResponse, error)
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

	// Odzyskanie unieważnia tokeny wydane wcześniej: rozgłasza zmianę urządzeniom konta i zrywa ich gniazda.
	r.Zarejestruj(shared.CommandAuthReset,
		obsluz(func(ctx context.Context, z shared.AuthResetRequest) (shared.AuthResetResponse, error) {
			odpowiedz, kontoId, err := ustawNoweHaslo(ctx, u, z)
			if err == nil {
				// Żądanie idzie z urządzenia bez sesji, więc konto adresata wchodzi z drogi odzyskania, nie z gniazda.
				ctx = zKontemOdzyskania(ctx, kontoId)
				e.wyslijDoKonta(ctx, shared.EventAuthChanged, "", shared.AuthChangedEvent{
					Reason: shared.AuthChangeReasonPasswordReset,
				})
				rozlaczPoUniewaznieniu(ctx, e, "sesja bramki unieważniona po ustawieniu nowego hasła")
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandAuthEmailChangeStart,
		obsluz(func(ctx context.Context, z shared.AuthEmailChangeStartRequest) (shared.AuthEmailChangeStartResponse, error) {
			return u.RozpocznijZmianeAdresu(ctx, z)
		}))

	// Zmiana adresu przenosi tożsamość konta, więc urządzenia dowiadują się o niej
	// tak samo jak o zmianie hasła.
	r.Zarejestruj(shared.CommandAuthEmailChangeConfirm,
		obsluz(func(ctx context.Context, z shared.AuthEmailChangeConfirmRequest) (shared.AuthEmailChangeConfirmResponse, error) {
			odpowiedz, err := u.PotwierdzZmianeAdresu(ctx, z)
			if err == nil {
				e.wyslijDoKonta(ctx, shared.EventAuthChanged, "", shared.AuthChangedEvent{
					Reason: shared.AuthChangeReasonEmailChanged,
				})
			}
			return odpowiedz, err
		}))

	// Wycofanie idzie drogą z listu, bez sesji: rozgłoszenia nie ma czym adresować.
	r.Zarejestruj(shared.CommandAuthEmailChangeRevoke,
		obsluz(func(ctx context.Context, z shared.AuthEmailChangeRevokeRequest) (shared.AuthEmailChangeRevokeResponse, error) {
			return u.WycofajZmianeAdresu(ctx, z)
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
				rozglosZmianeBramki(ctx, e, shared.AuthChangeReasonMethodAdded,
					odpowiedz.Methods, &z.DeviceId)
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandAuthMethodRemove,
		obsluz(func(ctx context.Context, z shared.AuthMethodRemoveRequest) (shared.AuthMethodRemoveResponse, error) {
			odpowiedz, err := u.ZdejmijMetodeWejscia(ctx, z)
			// Bez zdjęcia nie ma zmiany, więc nie ma czego rozgłaszać.
			if err == nil && odpowiedz.Removed {
				rozglosZmianeBramki(ctx, e, shared.AuthChangeReasonMethodRemoved,
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
				rozglosZmianeBramki(ctx, e, shared.AuthChangeReasonPasswordChanged, metody, nil)
				rozlaczPoUniewaznieniu(ctx, e, "sesja bramki unieważniona po zmianie hasła")
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandAuthTokenRefresh,
		obsluz(func(ctx context.Context, z shared.AuthTokenRefreshRequest) (shared.AuthTokenRefreshResponse, error) {
			// Więź połączenia wchodzi kontekstem tak samo jak przy zmianie hasła: puste Token znaczy bieżącą.
			return u.PrzedluzSesjeBramki(zSesjaBiezaca(ctx, wiez.SkrotKontekstu(ctx)), z)
		}))
}

// rozglosZmianeBramki wysyła `auth.changed` do urządzeń konta wołającego.
// Zdarzenie idzie bez sesji komunikatu: bramka nie należy do żadnej karty
// sesji — to stan platformy, a nie stan pracy. Nadajnik niepodłączony nie
// jest błędem.
func rozglosZmianeBramki(ctx context.Context, e *emiter, powod shared.AuthChangeReason,
	metody []shared.AuthMethod, urzadzenie *string) {

	if e == nil {
		return
	}
	e.wyslijDoKonta(ctx, shared.EventAuthChanged, "", shared.AuthChangedEvent{
		Reason:   powod,
		Methods:  metody,
		DeviceId: niepustyTekst(urzadzenie),
	})
}

// odzyskanieZKontem jest rozszerzeniem nieobowiązkowym portu: oddaje konto,
// którego drogą odzyskania ustawiono hasło. Odpowiedź kontraktu konta nie
// niesie, a rozgłoszenie i zerwanie gniazd potrzebują adresata.
type odzyskanieZKontem interface {
	ustawNoweHasloKonta(ctx context.Context, z shared.AuthResetRequest) (shared.AuthResetResponse, int64, error)
}

// ustawNoweHaslo wykonuje `auth.reset` i oddaje konto drogi, gdy adapter je zna; zero znaczy konto nieznane.
func ustawNoweHaslo(ctx context.Context, u Uwierzytelnianie,
	z shared.AuthResetRequest) (shared.AuthResetResponse, int64, error) {

	if adapter, umie := u.(odzyskanieZKontem); umie {
		return adapter.ustawNoweHasloKonta(ctx, z)
	}
	odpowiedz, err := u.UstawNoweHaslo(ctx, z)
	return odpowiedz, 0, err
}

// zKontemOdzyskania wpisuje do kontekstu konto z drogi odzyskania; konto nieznane zostawia kontekst bez zmiany.
func zKontemOdzyskania(ctx context.Context, kontoId int64) context.Context {
	if kontoId == 0 {
		return ctx
	}
	return dane.ZKontemOperatora(ctx, kontoId)
}

// rozlaczPoUniewaznieniu zrywa gniazda konta wołającego, których sesja bramki
// przestała nadawać. Idzie po rozgłoszeniu, żeby zdarzenie doszło przed
// zerwaniem. Nadajnik bez zrywania — sprawdzian — nie jest błędem.
func rozlaczPoUniewaznieniu(ctx context.Context, e *emiter, powod string) {
	rozlaczanie := e.rozlaczanie()
	if rozlaczanie == nil {
		return
	}
	rozlaczanie.RozlaczPoUniewaznieniu(ctx, kontoAdresata(ctx), powod)
}
