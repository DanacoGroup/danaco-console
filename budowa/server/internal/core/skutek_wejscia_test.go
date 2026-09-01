// Sprawdziany drogi wejścia do aplikacji mierzą skutek trwały w bazie i w
// skrzynce, a nie samą odpowiedź udaną komendy.
package core

import (
	"fmt"
	"strings"
	"testing"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/transport"
	"danacoconsole/shared"
)

// ── rejestracja ──────────────────────────────────────────────────────────────

// ── droga bez poczty ─────────────────────────────────────────────────────────

// Cztery sprawdziany poniżej trzymają drogę pierwszego uruchomienia na maszynie
// bez konta nadawczego.

// TestRejestracjaBezPocztyZakladaKontoIStawiaZnacznik mierzy skutek trwały
// rejestracji bez poczty: wiersz konta, kotwicę hasła i wpis w sejfie, nie samo
// pole `registered`.
func TestRejestracjaBezPocztyZakladaKontoIStawiaZnacznik(t *testing.T) {
	u := zmontujDrogeWejscia(t, pocztaBrak)

	var tresc shared.AuthRegisterResponse
	wykonajUdana(t, u.rdzen, u.zycie, shared.CommandAuthRegister,
		shared.AuthRegisterRequest{
			Login:    loginSprawdzianu,
			Email:    adresSprawdzianu,
			Password: hasloPierwsze,
		}, &tresc)

	if !tresc.Registered {
		t.Error("rejestracja bez poczty oddała stan udany i registered=false naraz")
	}
	if tresc.PendingVerification {
		t.Error("rejestracja bez poczty oznaczyła konto jako czekające na potwierdzenie —" +
			" klient pokaże okno wpisania drogi z listu, którego nikt nie wysłał")
	}

	if ile := liczbaWierszy(t, u, `SELECT COUNT(*) FROM konto_wlasciciela`); ile != 1 {
		t.Errorf("po rejestracji bez poczty w bazie leży %d kont właściciela, oczekiwane 1", ile)
	}
	if ile := liczbaWierszy(t, u,
		`SELECT COUNT(*) FROM konto_wlasciciela WHERE potwierdzone = 0`); ile != 1 {
		t.Errorf("kont niepotwierdzonych w bazie: %d, oczekiwane 1 — adresu nikt nie"+
			" sprawdził i nikt tego stanu nie ma prawa udawać", ile)
	}
	if ile := liczbaWierszy(t, u,
		`SELECT COUNT(*) FROM metoda_uwierzytelnienia WHERE kotwica = 1`); ile != 1 {
		t.Errorf("kotwic hasła w bazie: %d, oczekiwana 1 — bez niej bramki nie otworzy nic", ile)
	}
	// Droga potwierdzenia nie powstaje bez czynności nadania listu, która się
	// tu nie wykonuje.
	if ile := liczbaWierszy(t, u, `SELECT COUNT(*) FROM potwierdzenie_tozsamosci`); ile != 0 {
		t.Errorf("po rejestracji bez poczty leży %d dróg potwierdzenia, oczekiwane 0", ile)
	}
	// Rejestracja nie wpuszcza — ani z pocztą, ani bez niej.
	if ile := liczbaWierszy(t, u, `SELECT COUNT(*) FROM sesja_bramki`); ile != 0 {
		t.Errorf("rejestracja bez poczty założyła %d sesji bramki", ile)
	}

	adres, jest := znacznikWSejfie(t, u)
	if !jest {
		t.Fatal("bramka nie zapamiętała, że powstała bez poczty — przy następnym wejściu" +
			" zamknie się przed potwierdzeniem, którego platforma nie miała czym wysłać")
	}
	if adres != adresSprawdzianu {
		t.Errorf("znacznik niesie adres %q, rejestracja podała %q — bez właściwego adresu"+
			" odmowy nie powiedzą, o który adres idzie", adres, adresSprawdzianu)
	}
}

// TestBezPocztyBramkeOtwieraSamoHaslo pilnuje, że konto założone bez listu
// naprawdę wpuszcza: dowodem jest token sesji i wiersz w `sesja_bramki`, nie
// samo pole `ok`.
func TestBezPocztyBramkeOtwieraSamoHaslo(t *testing.T) {
	u := zmontujDrogeWejscia(t, pocztaBrak)

	var rejestracja shared.AuthRegisterResponse
	wykonajUdana(t, u.rdzen, u.zycie, shared.CommandAuthRegister,
		shared.AuthRegisterRequest{
			Login:    loginSprawdzianu,
			Email:    adresSprawdzianu,
			Password: hasloPierwsze,
		}, &rejestracja)

	var wejscie shared.AuthLoginResponse
	wykonajUdana(t, u.rdzen, u.zycie, shared.CommandAuthLogin, shared.AuthLoginRequest{
		Method: shared.AuthMethodKindPassword,
		Login:  wskaznik(loginSprawdzianu),
		Secret: wskaznik(hasloPierwsze),
	}, &wejscie)

	if strings.TrimSpace(wejscie.Session.Token) == "" {
		t.Error("wejście hasłem oddało sesję bez tokenu — bramka wpuściła i nie dała klucza")
	}
	if ile := liczbaWierszy(t, u,
		`SELECT COUNT(*) FROM sesja_bramki WHERE uniewazniono IS NULL`); ile != 1 {
		t.Errorf("po wejściu hasłem czynnych sesji w bazie: %d, oczekiwana 1", ile)
	}
	// Wejście niczego nie udaje: adres pozostaje niepotwierdzony, a znacznik
	// bramki nadal stoi.
	if ile := liczbaWierszy(t, u,
		`SELECT COUNT(*) FROM konto_wlasciciela WHERE potwierdzone = 0`); ile != 1 {
		t.Errorf("po wejściu hasłem kont niepotwierdzonych: %d, oczekiwane 1", ile)
	}
	if _, jest := znacznikWSejfie(t, u); !jest {
		t.Error("wejście hasłem zdjęło znacznik bramki bez poczty — następne wejście" +
			" odbije się o brak potwierdzenia")
	}
	// Złe hasło ma dalej odmawiać: wyjątek zdejmuje wyłącznie warunek
	// potwierdzenia adresu.
	blad := wykonajOdmowna(t, u.rdzen, u.zycie, shared.CommandAuthLogin,
		shared.AuthLoginRequest{
			Method: shared.AuthMethodKindPassword,
			Login:  wskaznik(loginSprawdzianu),
			Secret: wskaznik(hasloDrugie),
		})
	if blad.Code != shared.ErrorCodeNotAuthenticated {
		t.Errorf("wejście złym hasłem niesie kod %q, oczekiwany %q",
			blad.Code, shared.ErrorCodeNotAuthenticated)
	}
}

// TestBezPocztyPotwierdzenieAdresuCzekaNaDrogeZListu pilnuje, że ustawienie
// nadajnika w oknie Konfiguracji nie potwierdza adresu samo z siebie.
func TestBezPocztyPotwierdzenieAdresuCzekaNaDrogeZListu(t *testing.T) {
	u := zmontujDrogeWejscia(t, pocztaBrak)

	var rejestracja shared.AuthRegisterResponse
	wykonajUdana(t, u.rdzen, u.zycie, shared.CommandAuthRegister,
		shared.AuthRegisterRequest{
			Login:    loginSprawdzianu,
			Email:    adresSprawdzianu,
			Password: hasloPierwsze,
		}, &rejestracja)

	// Nadajnik pojawia się po rejestracji tą samą drogą, którą ustawia go
	// Operator: komendą konfiguracji.
	odbiornik := podnieOdbiornikSMTP(t)
	ustawNadajnik(t, u, odbiornik)

	// Drugiej rejestracji nie ma, więc drugiego listu z drogą weryfikacji też nie.
	powtorka := wykonajOdmowna(t, u.rdzen, u.zycie, shared.CommandAuthRegister,
		shared.AuthRegisterRequest{
			Login:    loginSprawdzianu,
			Email:    adresSprawdzianu,
			Password: hasloPierwsze,
		})
	if powtorka.Code != shared.ErrorCodeConflict {
		t.Errorf("druga rejestracja niesie kod %q, oczekiwany %q",
			powtorka.Code, shared.ErrorCodeConflict)
	}

	if ile := liczbaWierszy(t, u,
		`SELECT COUNT(*) FROM konto_wlasciciela WHERE potwierdzone = 0`); ile != 1 {
		t.Errorf("po ustawieniu nadajnika kont niepotwierdzonych: %d, oczekiwane 1", ile)
	}
	if _, jest := znacznikWSejfie(t, u); !jest {
		t.Error("znacznik zniknął po ustawieniu nadajnika, a adres pozostał niepotwierdzony —" +
			" bramka zamknie się przed potwierdzeniem, którego nie da się zdobyć")
	}
	// Bramka ma być otwarta dalej: nadajnik ustawiony nie może zamknąć wejścia,
	// które przed nim działało.
	var wejscie shared.AuthLoginResponse
	wykonajUdana(t, u.rdzen, u.zycie, shared.CommandAuthLogin, shared.AuthLoginRequest{
		Method: shared.AuthMethodKindPassword,
		Login:  wskaznik(loginSprawdzianu),
		Secret: wskaznik(hasloPierwsze),
	}, &wejscie)
	if strings.TrimSpace(wejscie.Session.Token) == "" {
		t.Error("po ustawieniu nadajnika wejście hasłem oddało sesję bez tokenu")
	}

	// Konto niepotwierdzone dostaje z `auth.recover` drogę weryfikacji, bo
	// jedyną jego przeszkodą jest brak aktywacji. Sam list adresu nie potwierdza —
	// potwierdza dopiero droga z niego podana do `auth.verify`.
	var odzyskanie shared.AuthRecoverResponse
	wykonajUdana(t, u.rdzen, u.zycie, shared.CommandAuthRecover,
		shared.AuthRecoverRequest{Email: adresSprawdzianu}, &odzyskanie)
	droga := drogaZListu(t, odbiornik.Ostatni(t))
	if ile := liczbaWierszy(t, u,
		`SELECT COUNT(*) FROM konto_wlasciciela WHERE potwierdzone = 0`); ile != 1 {
		t.Errorf("po wysłaniu listu kont niepotwierdzonych: %d, oczekiwane 1 —"+
			" list wysłany nie jest listem przeczytanym", ile)
	}
	if ile := liczbaWierszy(t, u,
		`SELECT COUNT(*) FROM potwierdzenie_tozsamosci WHERE cel = ?`, dane.CelWeryfikacja); ile != 1 {
		t.Errorf("dróg weryfikacji w bazie: %d, oczekiwana 1 — konto niepotwierdzone"+
			" ma dostać drogę aktywacji, nie odzyskania", ile)
	}

	var potwierdzenie shared.AuthVerifyResponse
	wykonajUdana(t, u.rdzen, u.zycie, shared.CommandAuthVerify,
		shared.AuthVerifyRequest{Token: droga}, &potwierdzenie)
	if !potwierdzenie.Verified {
		t.Fatal("potwierdzenie drogą z listu oddało stan udany i verified=false naraz")
	}
	if ile := liczbaWierszy(t, u,
		`SELECT COUNT(*) FROM konto_wlasciciela WHERE potwierdzone = 1`); ile != 1 {
		t.Errorf("po drodze z listu kont potwierdzonych: %d, oczekiwane 1", ile)
	}
	if _, jest := znacznikWSejfie(t, u); jest {
		t.Error("znacznik bramki bez poczty stoi po potwierdzeniu adresu —" +
			" bramkę trzyma odtąd wiersz konta, nie znacznik")
	}
}

// TestZnacznikBezPocztyStoiTylkoTamGdzieListuNieBylo pilnuje granicy istnienia
// znacznika bramki: na drodze z pocztą znacznik nie powstaje ani przed
// potwierdzeniem adresu, ani po nim.
func TestZnacznikBezPocztyStoiTylkoTamGdzieListuNieBylo(t *testing.T) {
	u := zmontujDrogeWejscia(t, pocztaDziala)

	droga := zarejestrujWlasciciela(t, u)
	if adres, jest := znacznikWSejfie(t, u); jest {
		t.Fatalf("rejestracja z wysłanym listem postawiła znacznik bramki bez poczty (adres %q) —"+
			" bramka przestałaby czekać na potwierdzenie, które właśnie poszło listem", adres)
	}

	var potwierdzenie shared.AuthVerifyResponse
	wykonajUdana(t, u.rdzen, u.zycie, shared.CommandAuthVerify,
		shared.AuthVerifyRequest{Token: droga}, &potwierdzenie)
	if !potwierdzenie.Verified {
		t.Fatal("potwierdzenie adresu oddało stan udany i verified=false naraz")
	}

	if ile := liczbaWierszy(t, u,
		`SELECT COUNT(*) FROM konto_wlasciciela WHERE potwierdzone = 1`); ile != 1 {
		t.Errorf("po potwierdzeniu kont potwierdzonych: %d, oczekiwane 1", ile)
	}
	if adres, jest := znacznikWSejfie(t, u); jest {
		t.Errorf("po potwierdzeniu adresu w sejfie leży znacznik bramki bez poczty (adres %q) —"+
			" bramkę trzyma odtąd sam wiersz konta i drugiej pamięci o niej nie ma", adres)
	}
}

// TestRejestracjaNieZakladaSesjiIOddajePendingVerification pilnuje, że
// rejestracja nie wpuszcza: dowodem jest brak wiersza w `sesja_bramki`, a nie
// brak pola w odpowiedzi.
func TestRejestracjaNieZakladaSesjiIOddajePendingVerification(t *testing.T) {
	u := zmontujDrogeWejscia(t, pocztaDziala)

	odpowiedz := wykonajKomende(t, u.rdzen, u.zycie, shared.CommandAuthRegister,
		shared.AuthRegisterRequest{
			Login:    loginSprawdzianu,
			Email:    adresSprawdzianu,
			Password: hasloPierwsze,
		})
	if odpowiedz.Error != nil {
		t.Fatalf("rejestracja odmówiła: kod=%s treść=%s", odpowiedz.Error.Code, odpowiedz.Error.Message)
	}

	var tresc shared.AuthRegisterResponse
	if err := protocol.LadunekDo(odpowiedz, &tresc); err != nil {
		t.Fatalf("nieczytelny ładunek rejestracji: %v", err)
	}
	if !tresc.Registered {
		t.Error("rejestracja oddała stan udany i registered=false naraz")
	}
	if !tresc.PendingVerification {
		t.Error("rejestracja nie oznaczyła konta jako czekającego na potwierdzenie adresu —" +
			" klient nie pokaże okna wpisania drogi z listu")
	}
	if strings.Contains(string(odpowiedz.Payload), "token") {
		t.Errorf("odpowiedź rejestracji niesie token: %s", string(odpowiedz.Payload))
	}

	if ile := liczbaWierszy(t, u, `SELECT COUNT(*) FROM sesja_bramki`); ile != 0 {
		t.Errorf("rejestracja założyła %d sesji bramki — wejście ma wydawać dopiero"+
			" potwierdzenie adresu", ile)
	}
	if ile := liczbaWierszy(t, u,
		`SELECT COUNT(*) FROM konto_wlasciciela WHERE potwierdzone = 0`); ile != 1 {
		t.Errorf("kont niepotwierdzonych w bazie: %d, oczekiwane 1", ile)
	}
}

/*
TestRejestracjaOdmawiaTylkoPrzyKolizji pilnuje reguły z rejestru decyzji, poz. 22:
platforma przyjmuje dowolną liczbę kont, a odmawia wyłącznie wtedy, gdy zajęty
jest login albo adres. Tożsamość jest własnością wiersza konta, nie instalacji.

Sprawdzian idzie trzema drogami, bo odmowa i przyjęcie różnią się tu jedną
wartością: konto obce przechodzi, kolizja adresu i kolizja loginu odmawiają.
*/
func TestRejestracjaOdmawiaTylkoPrzyKolizji(t *testing.T) {
	u := zmontujDrogeWejscia(t, pocztaDziala)
	zarejestrujWlasciciela(t, u)

	// Konto o wolnym loginie i wolnym adresie wchodzi — to jest sedno poz. 22.
	var drugie shared.AuthRegisterResponse
	wykonajUdana(t, u.rdzen, u.zycie, shared.CommandAuthRegister,
		shared.AuthRegisterRequest{
			Login:    "drugi",
			Email:    "drugi@danaco.sprawdzian",
			Password: hasloDrugie,
		}, &drugie)
	if !drugie.Registered {
		t.Error("konto o wolnym loginie i adresie nie zostało założone")
	}

	// Adres zajęty odmawia: jedno konto na jeden adres.
	blad := wykonajOdmowna(t, u.rdzen, u.zycie, shared.CommandAuthRegister,
		shared.AuthRegisterRequest{
			Login:    "trzeci",
			Email:    adresSprawdzianu,
			Password: hasloDrugie,
		})
	if blad.Code != shared.ErrorCodeConflict {
		t.Errorf("rejestracja na zajęty adres niesie kod %q, oczekiwany %q",
			blad.Code, shared.ErrorCodeConflict)
	}

	// Login zajęty odmawia tak samo — obie wartości są jednoznaczne.
	blad = wykonajOdmowna(t, u.rdzen, u.zycie, shared.CommandAuthRegister,
		shared.AuthRegisterRequest{
			Login:    loginSprawdzianu,
			Email:    "czwarty@danaco.sprawdzian",
			Password: hasloDrugie,
		})
	if blad.Code != shared.ErrorCodeConflict {
		t.Errorf("rejestracja na zajęty login niesie kod %q, oczekiwany %q",
			blad.Code, shared.ErrorCodeConflict)
	}

	if ile := liczbaWierszy(t, u, `SELECT COUNT(*) FROM konto_wlasciciela`); ile != 2 {
		t.Errorf("kont w bazie: %d, oczekiwane 2 — weszły oba wolne, odpadły obie kolizje", ile)
	}
	// Kotwica jest jedna NA KONTO, nie jedna w tabeli: tak stanowi migracja 406.
	if ile := liczbaWierszy(t, u,
		`SELECT COUNT(*) FROM metoda_uwierzytelnienia WHERE kotwica = 1`); ile != 2 {
		t.Errorf("kotwic bramki w bazie: %d, oczekiwane 2 — po jednej na konto", ile)
	}
	if listy := u.poczta.Listy(); len(listy) != 2 {
		t.Errorf("do skrzynki przyszło %d listów, oczekiwane 2 — odmówiona rejestracja"+
			" nie ma prawa wysyłać drogi potwierdzenia", len(listy))
	}
}

// ── potwierdzenie adresu ─────────────────────────────────────────────────────

// TestPotwierdzenieDrogaZListuWydajeSesjeAPowtorzenieOdmawia sprawdza
// jednorazowość drogi z listu: wpuszcza raz, a powtórzone użycie odmawia.
func TestPotwierdzenieDrogaZListuWydajeSesjeAPowtorzenieOdmawia(t *testing.T) {
	u := zmontujDrogeWejscia(t, pocztaDziala)
	droga := zarejestrujWlasciciela(t, u)

	var potwierdzenie shared.AuthVerifyResponse
	wykonajUdana(t, u.rdzen, u.zycie, shared.CommandAuthVerify, shared.AuthVerifyRequest{
		Token:    droga,
		DeviceId: wskaznik("maszyna-pierwsza"),
	}, &potwierdzenie)

	if !potwierdzenie.Verified {
		t.Error("potwierdzenie oddało stan udany i verified=false naraz")
	}
	if potwierdzenie.Session.Token == "" {
		t.Fatal("potwierdzenie nie wydało tokenu — Operator zostaje przed platformą")
	}
	if ile := liczbaWierszy(t, u, `SELECT COUNT(*) FROM sesja_bramki WHERE token_skrot = ?`,
		skrotTokenu(potwierdzenie.Session.Token)); ile != 1 {
		t.Errorf("wydany token nie ma wiersza sesji w bazie (pasujących: %d) —"+
			" token bez wiersza otworzy platformę dokładnie zero razy", ile)
	}
	if ile := liczbaWierszy(t, u,
		`SELECT COUNT(*) FROM konto_wlasciciela WHERE potwierdzone = 1`); ile != 1 {
		t.Errorf("konto po potwierdzeniu nie jest potwierdzone w bazie (pasujących: %d)", ile)
	}
	if ile := liczbaWierszy(t, u,
		`SELECT COUNT(*) FROM potwierdzenie_tozsamosci WHERE skrot = ? AND uzyte = 1`,
		skrotTokenu(droga)); ile != 1 {
		t.Errorf("droga po użyciu nie została zamknięta w bazie (pasujących: %d)", ile)
	}

	// Ta sama droga drugi raz.
	blad := wykonajOdmowna(t, u.rdzen, u.zycie, shared.CommandAuthVerify,
		shared.AuthVerifyRequest{Token: droga, DeviceId: wskaznik("maszyna-druga")})
	if strings.TrimSpace(blad.Message) == "" {
		t.Error("druga próba odmówiła bez treści")
	}
	if ile := liczbaWierszy(t, u, `SELECT COUNT(*) FROM sesja_bramki`); ile != 1 {
		t.Errorf("po drugiej próbie w bazie leży %d sesji bramki, oczekiwana 1 —"+
			" droga zużyta wydała drugi token", ile)
	}

	// Druga połowa jednorazowości mierzona na drodze odzyskania, bez zapory
	// stanu konta przed drogą.
	var odzyskanie shared.AuthRecoverResponse
	wykonajUdana(t, u.rdzen, u.zycie, shared.CommandAuthRecover,
		shared.AuthRecoverRequest{Email: adresSprawdzianu}, &odzyskanie)
	drogaOdzyskania := drogaZListu(t, u.poczta.Ostatni(t))

	var zmiana shared.AuthResetResponse
	wykonajUdana(t, u.rdzen, u.zycie, shared.CommandAuthReset,
		shared.AuthResetRequest{Token: drogaOdzyskania, NewPassword: hasloDrugie}, &zmiana)

	powtorka := wykonajOdmowna(t, u.rdzen, u.zycie, shared.CommandAuthReset,
		shared.AuthResetRequest{Token: drogaOdzyskania, NewPassword: "haslo-trzecie-3310"})
	if powtorka.Code != shared.ErrorCodeNotAuthenticated {
		t.Errorf("druga zmiana tą samą drogą odmówiła kodem %q, oczekiwany %q",
			powtorka.Code, shared.ErrorCodeNotAuthenticated)
	}

	// Skutek odmowy, nie sama odmowa: hasło ma zostać tym z pierwszej zmiany.
	var wejscie shared.AuthLoginResponse
	wykonajUdana(t, u.rdzen, u.zycie, shared.CommandAuthLogin, shared.AuthLoginRequest{
		Method: shared.AuthMethodKindPassword,
		Login:  wskaznik(loginSprawdzianu),
		Secret: wskaznik(hasloDrugie),
	}, &wejscie)
	if wejscie.Session.Token == "" {
		t.Error("po odmówionej drugiej zmianie hasło z pierwszej zmiany przestało otwierać bramkę")
	}
}

// TestDrogaWydanaDoOdzyskaniaNieDzialaJakoDrogaWeryfikacji pilnuje pola `cel`:
// droga wydana do odzyskania nie potwierdza adresu, a droga weryfikacji nie
// ustawia hasła. Drogę odzyskania wydaje dopiero konto potwierdzone — konto
// niepotwierdzone dostaje z `auth.recover` drogę weryfikacji.
func TestDrogaWydanaDoOdzyskaniaNieDzialaJakoDrogaWeryfikacji(t *testing.T) {
	u := zmontujDrogeWejscia(t, pocztaDziala)
	drogaWeryfikacji := zarejestrujWlasciciela(t, u)

	// Droga weryfikacji podana do ustawienia hasła.
	blad := wykonajOdmowna(t, u.rdzen, u.zycie, shared.CommandAuthReset,
		shared.AuthResetRequest{Token: drogaWeryfikacji, NewPassword: hasloDrugie})
	if blad.Code != shared.ErrorCodeNotAuthenticated {
		t.Errorf("odmowa niesie kod %q, oczekiwany %q", blad.Code, shared.ErrorCodeNotAuthenticated)
	}
	if ile := liczbaWierszy(t, u,
		`SELECT COUNT(*) FROM potwierdzenie_tozsamosci WHERE uzyte = 1`); ile != 0 {
		t.Errorf("odrzucone użycie zamknęło %d dróg — droga odmówiona i zarazem spalona"+
			" zabiera Operatorowi jedyny materiał, jaki dostał listem", ile)
	}

	// Droga weryfikacji ma dalej działać w swojej czynności.
	var potwierdzenie shared.AuthVerifyResponse
	wykonajUdana(t, u.rdzen, u.zycie, shared.CommandAuthVerify,
		shared.AuthVerifyRequest{Token: drogaWeryfikacji}, &potwierdzenie)
	if !potwierdzenie.Verified {
		t.Fatal("droga weryfikacji przestała działać po odrzuconym użyciu w innej czynności")
	}

	// Konto potwierdzone: odzyskanie wydaje drogę odzyskania. Cel odzyskania
	// nie miał wcześniejszej drogi, więc odstęp między listami nie blokuje.
	u.poczta.Wyczysc()
	var odzyskanie shared.AuthRecoverResponse
	wykonajUdana(t, u.rdzen, u.zycie, shared.CommandAuthRecover,
		shared.AuthRecoverRequest{Email: adresSprawdzianu}, &odzyskanie)
	drogaOdzyskania := drogaZListu(t, u.poczta.Ostatni(t))

	if drogaOdzyskania == drogaWeryfikacji {
		t.Fatal("obie drogi są tym samym materiałem — jedna droga na dwie czynności" +
			" znosi rozdział celów niezależnie od pola w bazie")
	}
	if ile := liczbaWierszy(t, u,
		`SELECT COUNT(*) FROM potwierdzenie_tozsamosci WHERE cel = ?`, dane.CelOdzyskanie); ile != 1 {
		t.Errorf("dróg odzyskania w bazie: %d, oczekiwana 1", ile)
	}

	// Droga odzyskania podana do potwierdzenia adresu.
	blad = wykonajOdmowna(t, u.rdzen, u.zycie, shared.CommandAuthVerify,
		shared.AuthVerifyRequest{Token: drogaOdzyskania})
	if blad.Code != shared.ErrorCodeNotAuthenticated {
		t.Errorf("odmowa niesie kod %q, oczekiwany %q", blad.Code, shared.ErrorCodeNotAuthenticated)
	}
	if !strings.Contains(blad.Message, "innej czynności") {
		t.Errorf("odmowa nie mówi, że droga wydana jest do innej czynności: %q", blad.Message)
	}
	if ile := liczbaWierszy(t, u,
		`SELECT COUNT(*) FROM potwierdzenie_tozsamosci WHERE uzyte = 1 AND cel = ?`,
		dane.CelOdzyskanie); ile != 0 {
		t.Errorf("odrzucone użycie zamknęło %d dróg odzyskania", ile)
	}

	// Droga odzyskania ma dalej działać w swojej czynności.
	var zmiana shared.AuthResetResponse
	wykonajUdana(t, u.rdzen, u.zycie, shared.CommandAuthReset,
		shared.AuthResetRequest{Token: drogaOdzyskania, NewPassword: hasloDrugie}, &zmiana)
	if !zmiana.Changed {
		t.Error("droga odzyskania przestała działać po odrzuconym użyciu w innej czynności")
	}
}

// ── odzyskanie konta ─────────────────────────────────────────────────────────

// TestOdzyskanieDlaAdresuObcegoOdpowiadaTakSamoINieWysylaListu pilnuje, że
// `auth.recover` odpowiada identycznie dla adresu własnego i obcego oraz nie
// wysyła listu do adresu obcego. Konto jest potwierdzone drogą z listu
// rejestracyjnego, więc odzyskanie wydaje drogę odzyskania, a odstęp między
// listami nie blokuje — cel odzyskania nie miał wcześniejszej drogi.
func TestOdzyskanieDlaAdresuObcegoOdpowiadaTakSamoINieWysylaListu(t *testing.T) {
	u := zmontujDrogeWejscia(t, pocztaDziala)
	drogaWeryfikacji := zarejestrujWlasciciela(t, u)
	var potwierdzenie shared.AuthVerifyResponse
	wykonajUdana(t, u.rdzen, u.zycie, shared.CommandAuthVerify,
		shared.AuthVerifyRequest{Token: drogaWeryfikacji}, &potwierdzenie)
	if !potwierdzenie.Verified {
		t.Fatal("potwierdzenie adresu oddało stan udany i verified=false naraz")
	}
	u.poczta.Wyczysc()

	odpowiedzWlasna := wykonajKomende(t, u.rdzen, u.zycie, shared.CommandAuthRecover,
		shared.AuthRecoverRequest{Email: adresSprawdzianu})
	if odpowiedzWlasna.Error != nil {
		t.Fatalf("odzyskanie dla adresu konta odmówiło: %+v", *odpowiedzWlasna.Error)
	}
	listyWlasne := u.poczta.Listy()
	if len(listyWlasne) != 1 {
		t.Fatalf("dla adresu konta wyszło %d listów, oczekiwany 1 — bez listu"+
			" Operator nie ma czym odzyskać konta", len(listyWlasne))
	}
	if listyWlasne[0].Odbiorca != adresSprawdzianu {
		t.Errorf("list poszedł na %q zamiast na %q", listyWlasne[0].Odbiorca, adresSprawdzianu)
	}

	u.poczta.Wyczysc()
	odpowiedzObca := wykonajKomende(t, u.rdzen, u.zycie, shared.CommandAuthRecover,
		shared.AuthRecoverRequest{Email: "ktos.obcy@gdzie.indziej"})
	if odpowiedzObca.Error != nil {
		t.Fatalf("odzyskanie dla adresu obcego odmówiło: %+v — odmowa sama w sobie"+
			" jest odpowiedzią „to nie ten adres"+`"`, *odpowiedzObca.Error)
	}

	if string(odpowiedzWlasna.Payload) != string(odpowiedzObca.Payload) {
		t.Errorf("odpowiedzi różnią się: dla adresu konta %s, dla obcego %s —"+
			" różnica mówi pytającemu, jaki adres ma Operator",
			string(odpowiedzWlasna.Payload), string(odpowiedzObca.Payload))
	}
	if (odpowiedzWlasna.Status == nil) != (odpowiedzObca.Status == nil) ||
		(odpowiedzWlasna.Status != nil && *odpowiedzWlasna.Status != *odpowiedzObca.Status) {
		t.Errorf("odpowiedzi różnią się stanem: %v wobec %v",
			odpowiedzWlasna.Status, odpowiedzObca.Status)
	}

	if listy := u.poczta.Listy(); len(listy) != 0 {
		t.Errorf("dla adresu obcego wyszło %d listów — platforma pisze do osoby,"+
			" która o nic nie prosiła", len(listy))
	}
	if ile := liczbaWierszy(t, u,
		`SELECT COUNT(*) FROM potwierdzenie_tozsamosci WHERE cel = ?`, dane.CelOdzyskanie); ile != 1 {
		t.Errorf("dróg odzyskania w bazie: %d, oczekiwana 1 — żądanie z adresem obcym"+
			" założyło drogę, której nikt nie dostał", ile)
	}
}

// TestUstawienieNowegoHaslaUniewazniaTokenyWydaneWczesniej pilnuje, że
// odzyskanie konta unieważnia token wydany wcześniej, zamyka stare hasło i
// otwiera bramkę nowym.
func TestUstawienieNowegoHaslaUniewazniaTokenyWydaneWczesniej(t *testing.T) {
	u := zmontujDrogeWejscia(t, pocztaDziala)
	drogaWeryfikacji := zarejestrujWlasciciela(t, u)

	var potwierdzenie shared.AuthVerifyResponse
	wykonajUdana(t, u.rdzen, u.zycie, shared.CommandAuthVerify, shared.AuthVerifyRequest{
		Token:    drogaWeryfikacji,
		DeviceId: wskaznik("maszyna-pierwsza"),
	}, &potwierdzenie)
	tokenStary := potwierdzenie.Session.Token

	var odzyskanie shared.AuthRecoverResponse
	wykonajUdana(t, u.rdzen, u.zycie, shared.CommandAuthRecover,
		shared.AuthRecoverRequest{Email: adresSprawdzianu}, &odzyskanie)
	drogaOdzyskania := drogaZListu(t, u.poczta.Ostatni(t))

	var zmiana shared.AuthResetResponse
	wykonajUdana(t, u.rdzen, u.zycie, shared.CommandAuthReset, shared.AuthResetRequest{
		Token:       drogaOdzyskania,
		NewPassword: hasloDrugie,
	}, &zmiana)
	if !zmiana.Changed {
		t.Error("ustawienie hasła oddało stan udany i changed=false naraz")
	}
	if zmiana.RevokedDevices < 1 {
		t.Errorf("odzyskanie unieważniło %d tokenów, a wydany był co najmniej jeden",
			zmiana.RevokedDevices)
	}
	if ile := liczbaWierszy(t, u,
		`SELECT COUNT(*) FROM sesja_bramki WHERE token_skrot = ? AND uniewazniono IS NOT NULL`,
		skrotTokenu(tokenStary)); ile != 1 {
		t.Errorf("token wydany przed zmianą hasła nie jest unieważniony w bazie (pasujących: %d)", ile)
	}

	// Stare hasło ma przestać otwierać.
	blad := wykonajOdmowna(t, u.rdzen, u.zycie, shared.CommandAuthLogin, shared.AuthLoginRequest{
		Method: shared.AuthMethodKindPassword,
		Login:  wskaznik(loginSprawdzianu),
		Secret: wskaznik(hasloPierwsze),
	})
	if blad.Code != shared.ErrorCodeNotAuthenticated {
		t.Errorf("wejście starym hasłem odmówiło kodem %q, oczekiwany %q",
			blad.Code, shared.ErrorCodeNotAuthenticated)
	}

	// Nowe hasło ma otwierać.
	var wejscie shared.AuthLoginResponse
	wykonajUdana(t, u.rdzen, u.zycie, shared.CommandAuthLogin, shared.AuthLoginRequest{
		Method: shared.AuthMethodKindPassword,
		Login:  wskaznik(loginSprawdzianu),
		Secret: wskaznik(hasloDrugie),
	}, &wejscie)
	if wejscie.Session.Token == "" {
		t.Fatal("wejście nowym hasłem nie wydało tokenu")
	}
	if wejscie.Session.Token == tokenStary {
		t.Error("wejście po zmianie hasła oddało token sprzed zmiany")
	}
	if ile := liczbaWierszy(t, u,
		`SELECT COUNT(*) FROM sesja_bramki WHERE token_skrot = ? AND uniewazniono IS NULL`,
		skrotTokenu(wejscie.Session.Token)); ile != 1 {
		t.Errorf("token wydany po zmianie hasła nie ma czynnego wiersza sesji (pasujących: %d)", ile)
	}
}

// ── wykaz urządzeń ───────────────────────────────────────────────────────────

// TestWykazUrzadzenOznaczaUrzadzenieBiezace pilnuje pola `current`: wykaz
// oznacza bieżące urządzenie po tożsamości połączenia, a bez niej nie oznacza
// żadnego.
func TestWykazUrzadzenOznaczaUrzadzenieBiezace(t *testing.T) {
	u := zmontujDrogeWejscia(t, pocztaDziala)
	droga := zarejestrujWlasciciela(t, u)

	// Kontekst połączenia — to on wiąże sesję z gniazdem przy potwierdzeniu.
	zGniazda := zPolaczeniem(u.zycie, transport.Tozsamosc{IdPolaczenia: "gniazdo-sprawdzianu"})

	var potwierdzenie shared.AuthVerifyResponse
	wykonajUdana(t, u.rdzen, zGniazda, shared.CommandAuthVerify, shared.AuthVerifyRequest{
		Token:    droga,
		DeviceId: wskaznik("maszyna-tutejsza"),
	}, &potwierdzenie)

	// Druga maszyna: sesja bramki założona wprost, bo wykaz urządzeń bierze się
	// z sesji.
	teraz := potwierdzenie.Session.ExpiresAt
	if _, err := u.rdzen.dane.Uwierzytelnienie.ZalozSesjeBramki(u.zycie, dane.SesjaBramki{
		SkrotTokenu:   skrotTokenu("token-maszyny-obcej"),
		UrzadzenieKod: wskaznik("maszyna-obca"),
		Wygasa:        teraz,
		Utworzono:     teraz,
	}); err != nil {
		t.Fatalf("nie można założyć sesji drugiego urządzenia: %v", err)
	}

	var wykaz shared.DeviceListResponse
	wykonajUdana(t, u.rdzen, zGniazda, shared.CommandDeviceList,
		shared.DeviceListRequest{}, &wykaz)

	if len(wykaz.Devices) != 2 {
		t.Fatalf("wykaz niesie %d urządzeń, oczekiwane 2: %+v", len(wykaz.Devices), wykaz.Devices)
	}
	biezace := make([]string, 0, len(wykaz.Devices))
	for _, urzadzenie := range wykaz.Devices {
		if urzadzenie.Current {
			biezace = append(biezace, urzadzenie.DeviceId)
		}
	}
	if len(biezace) != 1 {
		t.Fatalf("wykaz oznaczył %d urządzeń jako bieżące (%v), oczekiwane 1", len(biezace), biezace)
	}
	if biezace[0] != "maszyna-tutejsza" {
		t.Errorf("jako bieżące oznaczono %q, a połączenie stoi sesją urządzenia %q —"+
			" Operator odebrałby dostęp nie tej maszynie", biezace[0], "maszyna-tutejsza")
	}

	// Żądanie spoza gniazda: rdzeń nie zna bieżącego urządzenia i nie oznacza
	// żadnego.
	var wykazBezGniazda shared.DeviceListResponse
	wykonajUdana(t, u.rdzen, u.zycie, shared.CommandDeviceList,
		shared.DeviceListRequest{}, &wykazBezGniazda)
	for _, urzadzenie := range wykazBezGniazda.Devices {
		if urzadzenie.Current {
			t.Errorf("bez tożsamości połączenia wykaz oznaczył %q jako bieżące", urzadzenie.DeviceId)
		}
	}
}

// ── trwałość drogi ───────────────────────────────────────────────────────────

// TestDrogaPotwierdzeniaLezyWBazieWylacznieJakoSkrot pilnuje, że materiał drogi
// z listu leży w bazie wyłącznie jako skrót, czytany z każdej kolumny wiersza.
func TestDrogaPotwierdzeniaLezyWBazieWylacznieJakoSkrot(t *testing.T) {
	u := zmontujDrogeWejscia(t, pocztaDziala)
	droga := zarejestrujWlasciciela(t, u)

	wiersze, err := u.baza.DB.QueryContext(u.zycie, `SELECT * FROM potwierdzenie_tozsamosci`)
	if err != nil {
		t.Fatalf("nie można odczytać tabeli dróg potwierdzenia: %v", err)
	}
	defer wiersze.Close()

	kolumny, err := wiersze.Columns()
	if err != nil {
		t.Fatalf("nie można odczytać nazw kolumn: %v", err)
	}

	ile := 0
	znalezionySkrot := ""
	for wiersze.Next() {
		ile++
		wartosci := make([]any, len(kolumny))
		wskazniki := make([]any, len(kolumny))
		for i := range wartosci {
			wskazniki[i] = &wartosci[i]
		}
		if err := wiersze.Scan(wskazniki...); err != nil {
			t.Fatalf("nieczytelny wiersz drogi potwierdzenia: %v", err)
		}
		for i, kolumna := range kolumny {
			tekst := fmt.Sprintf("%v", wartosci[i])
			if bajty, czy := wartosci[i].([]byte); czy {
				tekst = string(bajty)
			}
			if strings.Contains(tekst, droga) {
				t.Errorf("kolumna %q niesie materiał wysłany listem (%q) —"+
					" kopia bazy wystarczy, żeby potwierdzić cudzą tożsamość", kolumna, tekst)
			}
			if kolumna == "skrot" {
				znalezionySkrot = tekst
			}
		}
	}
	if err := wiersze.Err(); err != nil {
		t.Fatalf("przerwany odczyt tabeli dróg potwierdzenia: %v", err)
	}

	if ile != 1 {
		t.Fatalf("dróg potwierdzenia w bazie: %d, oczekiwana 1", ile)
	}
	if znalezionySkrot != skrotTokenu(droga) {
		t.Errorf("skrót w bazie (%q) nie jest skrótem drogi z listu (%q) —"+
			" wiersz nie rozpozna materiału, który Operator dostał",
			znalezionySkrot, skrotTokenu(droga))
	}
}

// TestListRejestracyjnyNiesieDrogeINieNiesieHasla czyta treść wysłanego listu:
// droga potwierdzenia ma tam być, a hasło Operatora — nie.
func TestListRejestracyjnyNiesieDrogeINieNiesieHasla(t *testing.T) {
	u := zmontujDrogeWejscia(t, pocztaDziala)
	droga := zarejestrujWlasciciela(t, u)
	list := u.poczta.Ostatni(t)

	if list.Odbiorca != adresSprawdzianu {
		t.Errorf("list poszedł na %q zamiast na adres z rejestracji %q",
			list.Odbiorca, adresSprawdzianu)
	}
	if !strings.Contains(list.Dokument, "Subject:") {
		t.Error("list nie ma tematu — w skrzynce wygląda na śmieć")
	}
	if !strings.Contains(list.Dokument, loginSprawdzianu) {
		t.Errorf("list nie mówi, którego konta dotyczy; dokument:\n%s", list.Dokument)
	}
	if skrotTokenu(droga) != skrotTokenu(strings.TrimSpace(droga)) || droga == "" {
		t.Fatal("droga wyjęta z listu jest pusta albo obudowana odstępami")
	}
	if strings.Contains(list.Dokument, hasloPierwsze) {
		t.Errorf("list niesie hasło Operatora — hasła platforma nie zna jawnie"+
			" i nie ma prawa go odsyłać; dokument:\n%s", list.Dokument)
	}

	// Droga z listu ma być tą drogą, nie dowolnym napisem: dowodem jest wiersz
	// bazy jej skrótu.
	if ile := liczbaWierszy(t, u,
		`SELECT COUNT(*) FROM potwierdzenie_tozsamosci WHERE skrot = ? AND cel = ?`,
		skrotTokenu(droga), dane.CelWeryfikacja); ile != 1 {
		t.Errorf("materiał z listu nie ma wiersza w bazie (pasujących: %d) —"+
			" Operator wpisze go i zobaczy odmowę", ile)
	}
}
