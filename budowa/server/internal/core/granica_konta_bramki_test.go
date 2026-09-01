// Plik mierzy granicę konta w bramce: dwa konta na jednej instalacji, z których
// drugie nie sięga hasła, metod ani urządzeń pierwszego — ani drogą listu
// (`auth.reset`), ani z własnej sesji (`auth.password.reset`, `auth.method.remove`,
// `device.list`).
package core

import (
	"context"
	"testing"

	"danacoconsole/server/internal/transport"
	"danacoconsole/shared"
)

// Tożsamość konta drugiego. Konto pierwsze zakłada `zarejestrujWlasciciela`.
const (
	loginDrugiego     = "drugi-operator"
	adresDrugiego     = "drugi@danaco.sprawdzian"
	hasloDrugiego     = "haslo-drugiego-7715"
	hasloDrugiegoNowe = "haslo-drugiego-3182"
	hasloDrugiegoTrz  = "haslo-drugiego-9640"
)

// TestKontoDrugieNieSiegaKontaPierwszego mierzy skutek, nie odpowiedź: po
// wszystkich czynnościach konta drugiego hasło konta pierwszego ma dalej
// otwierać bramkę, a jego sesja ma zostać czynna.
func TestKontoDrugieNieSiegaKontaPierwszego(t *testing.T) {
	u := zmontujDrogeWejscia(t, pocztaDziala)

	gniazdoPierwszego := zPolaczeniem(u.zycie, transport.Tozsamosc{IdPolaczenia: "gniazdo-pierwszego"})
	gniazdoDrugiego := zPolaczeniem(u.zycie, transport.Tozsamosc{IdPolaczenia: "gniazdo-drugiego"})

	u.poczta.Wyczysc()
	drogaPierwszego := zarejestrujWlasciciela(t, u)
	var wejsciePierwszego shared.AuthVerifyResponse
	wykonajUdana(t, u.rdzen, gniazdoPierwszego, shared.CommandAuthVerify, shared.AuthVerifyRequest{
		Token:    drogaPierwszego,
		DeviceId: wskaznik("maszyna-pierwszego"),
	}, &wejsciePierwszego)

	u.poczta.Wyczysc()
	drogaDrugiego := zarejestrujKonto(t, u, loginDrugiego, adresDrugiego, hasloDrugiego)
	var wejscieDrugiego shared.AuthVerifyResponse
	wykonajUdana(t, u.rdzen, gniazdoDrugiego, shared.CommandAuthVerify, shared.AuthVerifyRequest{
		Token:    drogaDrugiego,
		DeviceId: wskaznik("maszyna-drugiego"),
	}, &wejscieDrugiego)

	// ── odzyskanie konta drugiego ────────────────────────────────────────────

	// Kod idzie na adres konta drugiego, więc zmienić ma hasło tego konta.
	u.poczta.Wyczysc()
	var odzyskanie shared.AuthRecoverResponse
	wykonajUdana(t, u.rdzen, u.zycie, shared.CommandAuthRecover,
		shared.AuthRecoverRequest{Email: adresDrugiego}, &odzyskanie)
	list := u.poczta.Ostatni(t)
	if list.Odbiorca != adresDrugiego {
		t.Fatalf("list z drogą poszedł na %q zamiast na %q", list.Odbiorca, adresDrugiego)
	}

	var zmiana shared.AuthResetResponse
	wykonajUdana(t, u.rdzen, u.zycie, shared.CommandAuthReset, shared.AuthResetRequest{
		Token:       drogaZListu(t, list),
		NewPassword: hasloDrugiegoNowe,
	}, &zmiana)
	if !zmiana.Changed {
		t.Fatal("odzyskanie konta drugiego oddało stan udany i changed=false naraz")
	}

	// Hasło konta pierwszego ma dalej otwierać bramkę.
	zalogujHaslem(t, u, u.zycie, loginSprawdzianu, hasloPierwsze)

	// Sesja konta pierwszego ma zostać czynna: unieważnienie obejmuje konto,
	// którego dotyczyła droga.
	if ile := liczbaWierszy(t, u,
		`SELECT COUNT(*) FROM sesja_bramki WHERE token_skrot = ? AND uniewazniono IS NULL`,
		skrotTokenu(wejsciePierwszego.Session.Token)); ile != 1 {
		t.Errorf("po odzyskaniu konta drugiego sesja konta pierwszego jest unieważniona"+
			" (czynnych wierszy: %d) — cudza zmiana hasła wyrzuciła Właściciela z platformy", ile)
	}

	// Konto drugie wchodzi hasłem, które sobie ustawiło, a poprzednie przestaje
	// otwierać — zmiana dosięgła tego konta, którego dotyczyła.
	odmowa := wykonajOdmowna(t, u.rdzen, u.zycie, shared.CommandAuthLogin, shared.AuthLoginRequest{
		Method: shared.AuthMethodKindPassword,
		Login:  wskaznik(loginDrugiego),
		Secret: wskaznik(hasloDrugiego),
	})
	if odmowa.Code != shared.ErrorCodeNotAuthenticated {
		t.Errorf("wejście konta drugiego hasłem sprzed zmiany odmówiło kodem %q, oczekiwany %q",
			odmowa.Code, shared.ErrorCodeNotAuthenticated)
	}
	/* Wejście wiąże gniazdo z nową sesją — dalsze czynności idą już z niej.
	   Gniazdo jest nowe, bo poprzednie stoi sesją unieważnioną przez odzyskanie
	   i każde żądanie z niego kończy się odmową, dopóki transport go nie
	   rozłączy. */
	gniazdoDrugiegoPoZmianie := zPolaczeniem(u.zycie,
		transport.Tozsamosc{IdPolaczenia: "gniazdo-drugiego-po-zmianie"})
	wejscieHaslem := zalogujHaslem(t, u, gniazdoDrugiegoPoZmianie, loginDrugiego, hasloDrugiegoNowe)

	// ── wykaz metod i urządzeń ───────────────────────────────────────────────

	kodMetodyPierwszego := kodKotwicyKonta(t, u, loginSprawdzianu)
	if len(wejscieHaslem.Methods) != 1 {
		t.Fatalf("wejście konta drugiego oddało %d metod, oczekiwana 1: %+v",
			len(wejscieHaslem.Methods), wejscieHaslem.Methods)
	}
	if wejscieHaslem.Methods[0].Id == kodMetodyPierwszego {
		t.Error("wykaz metod konta drugiego niesie kotwicę konta pierwszego —" +
			" po tym kodzie zdejmuje się cudzą metodę wejścia")
	}

	var wykaz shared.DeviceListResponse
	wykonajUdana(t, u.rdzen, gniazdoDrugiegoPoZmianie, shared.CommandDeviceList,
		shared.DeviceListRequest{}, &wykaz)
	for _, urzadzenie := range wykaz.Devices {
		if urzadzenie.DeviceId == "maszyna-pierwszego" {
			t.Errorf("wykaz urządzeń konta drugiego niesie maszynę konta pierwszego —"+
				" jednym `device.revoke` odbiera się cudzy dostęp: %+v", wykaz.Devices)
		}
	}

	// ── zdjęcie cudzej metody i zmiana cudzego hasła ─────────────────────────

	odmowa = wykonajOdmowna(t, u.rdzen, gniazdoDrugiegoPoZmianie, shared.CommandAuthMethodRemove,
		shared.AuthMethodRemoveRequest{MethodId: kodMetodyPierwszego})
	if odmowa.Code != shared.ErrorCodeNotFound {
		t.Errorf("zdjęcie kotwicy konta pierwszego z sesji konta drugiego odmówiło kodem %q,"+
			" oczekiwany %q", odmowa.Code, shared.ErrorCodeNotFound)
	}
	if ile := liczbaWierszy(t, u,
		`SELECT COUNT(*) FROM metoda_uwierzytelnienia WHERE identyfikator_zewnetrzny = ?`,
		kodMetodyPierwszego); ile != 1 {
		t.Fatalf("kotwica konta pierwszego zniknęła z bazy (wierszy: %d)", ile)
	}

	var zmianaHasla shared.AuthPasswordResetResponse
	wykonajUdana(t, u.rdzen, gniazdoDrugiegoPoZmianie, shared.CommandAuthPasswordReset,
		shared.AuthPasswordResetRequest{
			CurrentPassword: hasloDrugiegoNowe,
			NewPassword:     hasloDrugiegoTrz,
		}, &zmianaHasla)
	if !zmianaHasla.Changed {
		t.Fatal("zmiana hasła konta drugiego oddała stan udany i changed=false naraz")
	}

	// Skutek końcowy: każde konto stoi przy swoim haśle.
	zalogujHaslem(t, u, u.zycie, loginSprawdzianu, hasloPierwsze)
	zalogujHaslem(t, u, u.zycie, loginDrugiego, hasloDrugiegoTrz)
}

// zarejestrujKonto zakłada kolejne konto i oddaje drogę potwierdzenia z listu,
// który po nim przyszedł do skrzynki testowej.
func zarejestrujKonto(t *testing.T, u uprzazWejscia, login, email, haslo string) string {
	t.Helper()

	var odpowiedz shared.AuthRegisterResponse
	wykonajUdana(t, u.rdzen, u.zycie, shared.CommandAuthRegister, shared.AuthRegisterRequest{
		Login:      login,
		Email:      email,
		Password:   haslo,
		DeviceName: wskaznik("maszyna sprawdzianu"),
	}, &odpowiedz)
	if !odpowiedz.Registered {
		t.Fatalf("rejestracja konta %q oddała stan udany i registered=false naraz", login)
	}
	return drogaZListu(t, u.poczta.Ostatni(t))
}

// zalogujHaslem wchodzi przez bramkę hasłem i przerywa sprawdzian, gdy wejście
// nie wydało tokenu — hasło, które przestało otwierać, jest tu niepowodzeniem.
func zalogujHaslem(t *testing.T, u uprzazWejscia, ctx context.Context,
	login, haslo string) shared.AuthLoginResponse {

	t.Helper()

	var wejscie shared.AuthLoginResponse
	wykonajUdana(t, u.rdzen, ctx, shared.CommandAuthLogin, shared.AuthLoginRequest{
		Method: shared.AuthMethodKindPassword,
		Login:  wskaznik(login),
		Secret: wskaznik(haslo),
	}, &wejscie)
	if wejscie.Session.Token == "" {
		t.Fatalf("wejście konta %q oddało sesję bez tokenu", login)
	}
	return wejscie
}

// kodKotwicyKonta czyta z bazy identyfikator hasła wskazanego konta — tę samą
// wartość, którą wykaz metod oddaje oknu.
func kodKotwicyKonta(t *testing.T, u uprzazWejscia, login string) string {
	t.Helper()

	var kod string
	err := u.baza.DB.QueryRowContext(u.zycie,
		`SELECT identyfikator_zewnetrzny FROM metoda_uwierzytelnienia
		 WHERE kotwica = 1
		   AND konto_id = (SELECT id FROM konto_wlasciciela WHERE login = ?)`, login).Scan(&kod)
	if err != nil {
		t.Fatalf("nie można odczytać kotwicy konta %q: %v", login, err)
	}
	return kod
}
