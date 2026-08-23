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

// Sprawdziany drogi wejścia do aplikacji.
//
// Mierzony jest SKUTEK, nie koperta. Odpowiedź udana nie jest dowodem niczego —
// w tym produkcie zdarzało się `status: ok` przy pustym wyniku. Dowodem jest
// wiersz w bazie, list w skrzynce i to, że wywołanie następne zachowuje się
// inaczej niż pierwsze.

// ── rejestracja ──────────────────────────────────────────────────────────────

// TestRejestracjaBezKontaNadawczegoOdmawiaINieZakladaKonta pilnuje kolejności,
// od której zależy, czy Operator w ogóle wejdzie: nadajnik sprawdza się PRZED
// zapisaniem czegokolwiek.
//
// Konto założone bez wysłanego listu byłoby kontem, do którego nikt nie ma jak
// wejść: droga potwierdzenia idzie wyłącznie listem, a rejestracja wykonuje się
// raz i drugi raz odmawia konfliktem. Operator zostałby więc z platformą, która
// o nim wie, i bez jednej ścieżki do środka.
//
// Sprawdzian mierzy trzy tabele, nie samą odmowę: odmowa oddana po zapisaniu
// konta wygląda w kopercie identycznie jak odmowa oddana przed nim.
func TestRejestracjaBezKontaNadawczegoOdmawiaINieZakladaKonta(t *testing.T) {
	u := zmontujDrogeWejscia(t, pocztaBrak)

	blad := wykonajOdmowna(t, u.rdzen, u.zycie, shared.CommandAuthRegister,
		shared.AuthRegisterRequest{
			Login:    loginSprawdzianu,
			Email:    adresSprawdzianu,
			Password: hasloPierwsze,
		})
	if blad.Code != shared.ErrorCodeInternalError {
		t.Errorf("odmowa niesie kod %q, oczekiwany %q — brak konta nadawczego jest brakiem"+
			" po stronie platformy, nie pomyłką Operatora", blad.Code, shared.ErrorCodeInternalError)
	}
	if strings.TrimSpace(blad.Message) == "" {
		t.Error("odmowa bez treści — Operator nie dowie się, czego brakuje")
	}

	if ile := liczbaWierszy(t, u, `SELECT COUNT(*) FROM konto_wlasciciela`); ile != 0 {
		t.Errorf("po odmowie w bazie leży %d kont właściciela — konto założone bez listu"+
			" jest kontem bez drogi wejścia, a rejestracji nie da się powtórzyć", ile)
	}
	if ile := liczbaWierszy(t, u, `SELECT COUNT(*) FROM metoda_uwierzytelnienia`); ile != 0 {
		t.Errorf("po odmowie w bazie leży %d metod wejścia", ile)
	}
	if ile := liczbaWierszy(t, u, `SELECT COUNT(*) FROM potwierdzenie_tozsamosci`); ile != 0 {
		t.Errorf("po odmowie w bazie leży %d dróg potwierdzenia — droga zapisana bez"+
			" wysłanego listu nie jest znana nikomu", ile)
	}
}

// TestRejestracjaNieZakladaSesji pilnuje, że rejestracja nie wpuszcza.
//
// Gdyby wpuszczała, potwierdzenie adresu byłoby ozdobą: konto działałoby bez
// niego, a adres — jedyna droga odzyskania dostępu — zostawałby niesprawdzony aż
// do dnia, w którym trzeba nim odzyskać konto.
//
// Dowodem jest brak wiersza w `sesja_bramki`, a nie brak pola w odpowiedzi:
// sesja założona i przemilczana w kopercie byłaby sesją tak samo ważną.
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

// TestDrugaRejestracjaOdmawiaKodemConflict pilnuje jednorazowości rejestracji.
//
// Powtórzone żądanie jest próbą podmiany hasła bez znajomości starego — od tego
// jest odzyskanie konta. Odmowa musi być nazwana kodem `conflict`, bo tylko
// wtedy klient odróżni „konto już jest" od awarii wartej ponowienia.
//
// Sprawdzian mierzy też, czego druga rejestracja NIE zrobiła: nie podmieniła
// konta, nie założyła drugiej kotwicy i nie wysłała drugiego listu.
func TestDrugaRejestracjaOdmawiaKodemConflict(t *testing.T) {
	u := zmontujDrogeWejscia(t, pocztaDziala)
	zarejestrujWlasciciela(t, u)

	blad := wykonajOdmowna(t, u.rdzen, u.zycie, shared.CommandAuthRegister,
		shared.AuthRegisterRequest{
			Login:    "podszywacz",
			Email:    "podszywacz@danaco.sprawdzian",
			Password: hasloDrugie,
		})
	if blad.Code != shared.ErrorCodeConflict {
		t.Errorf("druga rejestracja niesie kod %q, oczekiwany %q", blad.Code, shared.ErrorCodeConflict)
	}

	if ile := liczbaWierszy(t, u,
		`SELECT COUNT(*) FROM konto_wlasciciela WHERE login = ? AND email = ?`,
		loginSprawdzianu, adresSprawdzianu); ile != 1 {
		t.Errorf("konto właściciela po drugiej rejestracji nie jest tym z pierwszej (pasujących wierszy: %d)", ile)
	}
	if ile := liczbaWierszy(t, u, `SELECT COUNT(*) FROM konto_wlasciciela`); ile != 1 {
		t.Errorf("kont właściciela w bazie: %d, oczekiwane 1", ile)
	}
	if ile := liczbaWierszy(t, u,
		`SELECT COUNT(*) FROM metoda_uwierzytelnienia WHERE kotwica = 1`); ile != 1 {
		t.Errorf("kotwic bramki w bazie: %d, oczekiwana 1", ile)
	}
	if listy := u.poczta.Listy(); len(listy) != 1 {
		t.Errorf("do skrzynki przyszło %d listów, oczekiwany 1 — odmówiona rejestracja"+
			" nie ma prawa wysyłać drugiej drogi potwierdzenia", len(listy))
	}
}

// ── potwierdzenie adresu ─────────────────────────────────────────────────────

// TestPotwierdzenieDrogaZListuWydajeSesjeAPowtorzenieOdmawia sprawdza obie
// połowy jednorazowości: droga wpuszcza raz i tylko raz.
//
// Droga wpuszczająca dwa razy jest drogą wpuszczającą każdego, kto zajrzy do
// skrzynki później — list zostaje w niej na zawsze.
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

	// Druga połowa jednorazowości, mierzona na drodze odzyskania.
	//
	// Powtórzone `auth.verify` odbija się o stan konta („adres jest już
	// potwierdzony") ZANIM dojdzie do drogi, więc samo w sobie nie dowodzi, że
	// zamknięcie drogi działa. Odzyskanie konta takiej zapory przed sobą nie ma:
	// jedynym, co zatrzymuje drugie użycie, jest zamknięcie wiersza drogi.
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

// TestDrogaWydanaDoOdzyskaniaNieDzialaJakoDrogaWeryfikacji pilnuje pola `cel`.
//
// Obie drogi wyglądają tak samo i leżą w jednej tabeli. Bez rozdziału po celu
// droga wysłana na prośbę o nowe hasło potwierdzałaby adres, a droga wysłana
// przy rejestracji ustawiałaby hasło — czyli list o jednej treści wykonywałby
// czynność, o którą nikt nie prosił.
//
// Sprawdzian idzie w obie strony i kończy dowodem, że próba użycia drogi nie
// w swojej czynności jej NIE zużyła: droga odrzucona i zarazem spalona byłaby
// gorsza od samej odmowy.
func TestDrogaWydanaDoOdzyskaniaNieDzialaJakoDrogaWeryfikacji(t *testing.T) {
	u := zmontujDrogeWejscia(t, pocztaDziala)
	drogaWeryfikacji := zarejestrujWlasciciela(t, u)

	var odzyskanie shared.AuthRecoverResponse
	wykonajUdana(t, u.rdzen, u.zycie, shared.CommandAuthRecover,
		shared.AuthRecoverRequest{Email: adresSprawdzianu}, &odzyskanie)
	drogaOdzyskania := drogaZListu(t, u.poczta.Ostatni(t))

	if drogaOdzyskania == drogaWeryfikacji {
		t.Fatal("obie drogi są tym samym materiałem — jedna droga na dwie czynności" +
			" znosi rozdział celów niezależnie od pola w bazie")
	}

	// Droga weryfikacji podana do ustawienia hasła.
	blad := wykonajOdmowna(t, u.rdzen, u.zycie, shared.CommandAuthReset,
		shared.AuthResetRequest{Token: drogaWeryfikacji, NewPassword: hasloDrugie})
	if blad.Code != shared.ErrorCodeNotAuthenticated {
		t.Errorf("odmowa niesie kod %q, oczekiwany %q", blad.Code, shared.ErrorCodeNotAuthenticated)
	}

	// Droga odzyskania podana do potwierdzenia adresu.
	blad = wykonajOdmowna(t, u.rdzen, u.zycie, shared.CommandAuthVerify,
		shared.AuthVerifyRequest{Token: drogaOdzyskania})
	if blad.Code != shared.ErrorCodeNotAuthenticated {
		t.Errorf("odmowa niesie kod %q, oczekiwany %q", blad.Code, shared.ErrorCodeNotAuthenticated)
	}

	if ile := liczbaWierszy(t, u,
		`SELECT COUNT(*) FROM potwierdzenie_tozsamosci WHERE uzyte = 1`); ile != 0 {
		t.Errorf("odrzucone użycie zamknęło %d dróg — droga odmówiona i zarazem spalona"+
			" zabiera Operatorowi jedyny materiał, jaki dostał listem", ile)
	}

	// Obie drogi mają dalej działać w swoich czynnościach.
	var potwierdzenie shared.AuthVerifyResponse
	wykonajUdana(t, u.rdzen, u.zycie, shared.CommandAuthVerify,
		shared.AuthVerifyRequest{Token: drogaWeryfikacji}, &potwierdzenie)
	if !potwierdzenie.Verified {
		t.Error("droga weryfikacji przestała działać po odrzuconym użyciu w innej czynności")
	}

	var zmiana shared.AuthResetResponse
	wykonajUdana(t, u.rdzen, u.zycie, shared.CommandAuthReset,
		shared.AuthResetRequest{Token: drogaOdzyskania, NewPassword: hasloDrugie}, &zmiana)
	if !zmiana.Changed {
		t.Error("droga odzyskania przestała działać po odrzuconym użyciu w innej czynności")
	}
}

// ── odzyskanie konta ─────────────────────────────────────────────────────────

// TestOdzyskanieDlaAdresuObcegoOdpowiadaTakSamoINieWysylaListu pilnuje, żeby
// `auth.recover` nie była wyrocznią.
//
// Komenda jest osiągalna przed zalogowaniem, więc pyta ją każdy. Odpowiedź
// różniąca się choć jednym bajtem mówiłaby pytającemu, jaki adres ma Operator —
// a to jest połowa materiału potrzebnego do podszycia się pod niego.
//
// Druga połowa sprawdzianu jest ważniejsza od pierwszej: dla adresu obcego list
// NIE wychodzi. Wysłany szedłby do osoby, która o nic nie prosiła, a rachunek za
// nadania płaciłby Operator.
func TestOdzyskanieDlaAdresuObcegoOdpowiadaTakSamoINieWysylaListu(t *testing.T) {
	u := zmontujDrogeWejscia(t, pocztaDziala)
	zarejestrujWlasciciela(t, u)
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
// odzyskanie konta naprawdę odbiera dostęp.
//
// Odzyskanie zaczyna się od podejrzenia, że dostęp ma ktoś jeszcze. Zostawienie
// mu ważnego tokenu czyniłoby zmianę hasła pozorną — wchodziłby dalej, bez hasła
// i bez śladu.
//
// Mierzone są trzy skutki, bo `changed: true` nie dowodzi żadnego z nich: token
// wydany wcześniej ma być unieważniony w bazie, stare hasło ma przestać otwierać
// bramkę, a nowe ma ją otwierać.
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

// TestWykazUrzadzenOznaczaUrzadzenieBiezace pilnuje pola `current`.
//
// Wiersz własnej maszyny wygląda w wykazie tak samo jak każdy inny. Bez
// oznaczenia Operator odbiera dostęp sobie i traci go w tej samej chwili —
// a rozstrzygnąć to może wyłącznie rdzeń: klient zna identyfikator, który sam
// nadał, ale nie wie, którą sesją stoi jego połączenie.
//
// Sprawdzian wchodzi tą samą drogą, co gniazdo: tożsamością połączenia
// w kontekście. Bez niej pomiar nie dotykałby badanego przypadku.
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

	// Druga maszyna: sesja bramki założona wprost, bo drugiego potwierdzenia
	// adresu już nie będzie, a wykaz urządzeń bierze się z sesji.
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

	// Żądanie spoza gniazda: rdzeń nie wie, które urządzenie jest bieżące,
	// i wtedy nie oznacza ŻADNEGO. Zgadywanie po ostatnim wejściu wskazałoby
	// cudzą maszynę jako własną.
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

// TestDrogaPotwierdzeniaLezyWBazieWylacznieJakoSkrot pilnuje obietnicy, na
// której stoi cała wartość wysyłania drogi listem.
//
// Gdyby materiał z listu leżał w tabeli, kopia bazy pozwalałaby potwierdzić cudzą
// tożsamość i ustawić hasło do konta — czyli byłaby wejściem do platformy, a nie
// zbiorem danych.
//
// Sprawdzian nie ufa nazwom kolumn: czyta CAŁY wiersz i szuka materiału z listu
// w każdej wartości. Kolumna dołożona kiedyś obok `skrot` przechodziłaby pomiar
// pytający wyłącznie o `skrot`.
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

// TestListRejestracyjnyNiesieDrogeINieNiesieHasla czyta to, co naprawdę poszło
// w świat.
//
// List jest jedynym miejscem, w którym materiał potwierdzenia istnieje jawnie,
// i jedynym, którego rdzeń po nadaniu już nie kontroluje. Dwa warunki muszą być
// spełnione naraz: droga MA tam być, bo bez niej list jest bezużyteczny,
// a hasła TAM BYĆ NIE MOŻE — platforma nie zna go jawnie i nie ma prawa odsyłać.
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

	// Droga z listu ma być TĄ drogą, nie dowolnym napisem w treści: dowodem jest
	// wiersz bazy rozpoznający jej skrót.
	if ile := liczbaWierszy(t, u,
		`SELECT COUNT(*) FROM potwierdzenie_tozsamosci WHERE skrot = ? AND cel = ?`,
		skrotTokenu(droga), dane.CelWeryfikacja); ile != 1 {
		t.Errorf("materiał z listu nie ma wiersza w bazie (pasujących: %d) —"+
			" Operator wpisze go i zobaczy odmowę", ile)
	}
}
