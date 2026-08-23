// Plik wpina sześć komend rodziny `auth.*` — bramki Operatora. Tyle właśnie
// komend niesie kontrakt w tej rodzinie i każda ma tu swój uchwyt.
//
// Cała bramka stoi na jednym porcie. Sześć komend dotyka dwóch tabel i jednego
// sejfu; osobne porty wejścia, metod i sesji byłyby trzema prawdami o jednej
// bramce.
//
// Zdarzenie `auth.changed` rozgłaszają trzy komendy, każda swoim powodem:
// założenie metody, jej zdjęcie i zmiana hasła. `auth.login` nie rozgłasza —
// wejście nie zmienia ani składu metod, ani hasła, a sekcja Uwierzytelnianie
// dostaje wykaz metod wprost w odpowiedzi. `auth.register` też nie: przy
// pierwszym uruchomieniu nie ma komu rozgłosić zmiany. `auth.token.refresh` nie
// zmienia stanu uwierzytelnienia w ogóle.
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

	// MetodyWejscia oddaje wykaz metod po zmianie. Potrzebuje go ładunek
	// zdarzenia `auth.changed` po zmianie hasła — odpowiedź tej komendy wykazu
	// nie niesie, a sekcja Uwierzytelnianie ma się odświeżyć bez odpytywania.
	MetodyWejscia(ctx context.Context) ([]shared.AuthMethod, error)

	// BramkaZalozona mówi, czy sekret bramki w ogóle istnieje. Potrzebuje tego
	// powitanie, żeby klient odróżnił „zaloguj się" od „ustaw hasło po raz
	// pierwszy" bez wyprowadzania stanu z odmowy.
	BramkaZalozona(ctx context.Context) (bool, error)

	// RozpoznajSesjeBramki sprawdza token przedstawiony w powitaniu i oddaje
	// skrót rozpoznanej sesji — skrót, nie token, bo więź trzyma w pamięci
	// wyłącznie to, co i tak leży w bazie. Fałsz bez błędu znaczy „token
	// nieznany, unieważniony albo wygasły" — to odpowiedź, nie awaria.
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

	// Rejestracja połączenia NIE wiąże i sesji nie zakłada: konto powstaje
	// niepotwierdzone, a token dostępu wydaje dopiero potwierdzenie adresu.
	r.Zarejestruj(shared.CommandAuthRegister,
		obsluz(func(ctx context.Context, z shared.AuthRegisterRequest) (shared.AuthRegisterResponse, error) {
			return u.ZalozBramke(ctx, z)
		}))

	// Potwierdzenie adresu wiąże połączenie od razu — to jest pierwsze wejście
	// Operatora do platformy. Bez związania byłby dla rdzenia nierozpoznany aż
	// do następnego powitania, a powitanie idzie raz, przy nawiązaniu gniazda.
	r.Zarejestruj(shared.CommandAuthVerify,
		obsluz(func(ctx context.Context, z shared.AuthVerifyRequest) (shared.AuthVerifyResponse, error) {
			odpowiedz, err := u.PotwierdzAdres(ctx, z)
			if err == nil {
				wiez.Zwiaz(polaczenieZKontekstu(ctx), skrotTokenu(odpowiedz.Session.Token))
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandAuthRecover,
		obsluz(func(ctx context.Context, z shared.AuthRecoverRequest) (shared.AuthRecoverResponse, error) {
			return u.RozpocznijOdzyskanie(ctx, z)
		}))

	// Odzyskanie unieważnia tokeny wydane wcześniej, więc rozgłasza zmianę —
	// pozostałe urządzenia mają się dowiedzieć, że wylatują, a nie odkryć tego
	// przy następnej komendzie.
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
				wiez.Zwiaz(polaczenieZKontekstu(ctx), skrotTokenu(odpowiedz.Session.Token))
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
			// Sesja wołającego wchodzi kontekstem, żeby zmiana hasła nie wyrzuciła
			// za drzwi tego, kto ją właśnie wykonał. Połączenie niezwiązane daje
			// skrót pusty, a wtedy unieważniane są wszystkie sesje.
			odpowiedz, err := u.ZmienHasloBramki(zSesjaBiezaca(ctx, wiez.SkrotKontekstu(ctx)), z)
			if err == nil && odpowiedz.Changed {
				// Wykaz metod bierzemy osobno — odpowiedź tej komendy go nie
				// niesie. Nieudany odczyt kończy wyłącznie ładunek zdarzenia,
				// nie komendę: hasło jest już zmienione.
				metody, _ := u.MetodyWejscia(ctx)
				rozglosZmianeBramki(e, shared.AuthChangeReasonPasswordChanged, metody, nil)
			}
			return odpowiedz, err
		}))

	r.Zarejestruj(shared.CommandAuthTokenRefresh,
		obsluz(func(ctx context.Context, z shared.AuthTokenRefreshRequest) (shared.AuthTokenRefreshResponse, error) {
			// Więź połączenia wchodzi kontekstem tak samo jak przy zmianie
			// hasła: kontrakt mówi przy `AuthTokenRefreshRequest.Token`, że
			// puste znaczy sesję bieżącego połączenia, a wskazać ją da się
			// wyłącznie stąd — żądanie tokenu wtedy nie niesie.
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
