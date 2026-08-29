// Straże drogi wejścia: trzy sprawdziany trzymają zachowania, których
// zerwanie zamyka Operatorowi drogę do platformy albo otwiera ją komuś, kto
// nie powinien wejść.
package core

import (
	"strings"
	"testing"

	"danacoconsole/shared"
)

// TestBramkaZamknietaDoPotwierdzeniaAdresu pilnuje, że weryfikacja adresu nie
// jest ozdobą: adres jest jedyną drogą odzyskania konta.
func TestBramkaZamknietaDoPotwierdzeniaAdresu(t *testing.T) {
	u := zmontujDrogeWejscia(t, pocztaDziala)
	zarejestrujWlasciciela(t, u)

	if ile := liczbaWierszy(t, u,
		`SELECT COUNT(*) FROM konto_wlasciciela WHERE potwierdzone = 0`); ile != 1 {
		t.Fatalf("konto po samej rejestracji nie jest niepotwierdzone (pasujących: %d)", ile)
	}

	blad := wykonajOdmowna(t, u.rdzen, u.zycie, shared.CommandAuthLogin,
		shared.AuthLoginRequest{
			Method: shared.AuthMethodKindPassword,
			Login:  wskaznik(loginSprawdzianu),
			Secret: wskaznik(hasloPierwsze),
		})

	if blad.Code != shared.ErrorCodeNotAuthenticated {
		t.Errorf("odmowa wejścia przy koncie niepotwierdzonym niesie kod %q, oczekiwany %q",
			blad.Code, shared.ErrorCodeNotAuthenticated)
	}
	// Odmowa ma prowadzić do naprawy: mówić, czego brakuje i czym to zrobić.

	/* „Nie wolno" bez drogi dalszej zostawia Operatora przed zamkniętym wejściem.
	   Drogę nazywa się czynnością, którą ma wykonać, nie nazwą komendy: nazwa
	   komendy jest pojęciem budowy i w zdaniu dla Operatora nie stoi. */
	for _, slowo := range []string{"kod", "wiadomoś"} {
		if !strings.Contains(strings.ToLower(blad.Message), slowo) {
			t.Errorf("odmowa nie wskazuje drogi potwierdzenia (brak %q): %q", slowo, blad.Message)
		}
	}

	// Odmowa nie ma prawa niczego zakładać po drodze.
	if ile := liczbaWierszy(t, u,
		`SELECT COUNT(*) FROM sesja_bramki WHERE uniewazniono IS NULL`); ile != 0 {
		t.Errorf("po odmówionym wejściu leży %d czynnych sesji", ile)
	}
}

// TestNieudaneNadanieListuSchodziNaDrogeBezPoczty pilnuje, że zatrzymany
// przekaźnik poczty nie zamienia się w trwałą utratę platformy: rejestracja
// wykonuje się raz i nie ma prawa cofać się bez otwartej drogi powrotu.
func TestNieudaneNadanieListuSchodziNaDrogeBezPoczty(t *testing.T) {
	u := zmontujDrogeWejscia(t, pocztaNieosiagalna)

	var tresc shared.AuthRegisterResponse
	wykonajUdana(t, u.rdzen, u.zycie, shared.CommandAuthRegister,
		shared.AuthRegisterRequest{
			Login:    loginSprawdzianu,
			Email:    adresSprawdzianu,
			Password: hasloPierwsze,
		}, &tresc)

	if !tresc.Registered {
		t.Error("rejestracja oddała stan udany i registered=false naraz")
	}
	if tresc.PendingVerification {
		t.Error("rejestracja po nieudanym nadaniu oznaczyła konto jako czekające na" +
			" potwierdzenie — klient każe przepisać drogę z listu, który nie wyszedł")
	}

	konta := liczbaWierszy(t, u, `SELECT COUNT(*) FROM konto_wlasciciela`)
	kotwice := liczbaWierszy(t, u, `SELECT COUNT(*) FROM metoda_uwierzytelnienia WHERE kotwica = 1`)
	if konta != 1 {
		t.Errorf("po nieudanym nadaniu kont właściciela w bazie: %d, oczekiwane 1", konta)
	}
	if kotwice != 1 {
		t.Errorf("po nieudanym nadaniu kotwic hasła w bazie: %d, oczekiwana 1 —"+
			" bez kotwicy konto zostaje bez klucza, a drugiej rejestracji nie ma", kotwice)
	}

	// Znacznik jest tu jedyną rzeczą, która trzyma bramkę otwartą.

	// Adres został niepotwierdzony, bo list nie doszedł.
	adres, jest := znacznikWSejfie(t, u)
	if !jest {
		t.Fatal("bramka nie zapamiętała nieudanego nadania — zamknie się przed" +
			" potwierdzeniem, którego nie miała jak wysłać")
	}
	if adres != adresSprawdzianu {
		t.Errorf("znacznik niesie adres %q, rejestracja podała %q", adres, adresSprawdzianu)
	}
	if ile := liczbaWierszy(t, u,
		`SELECT COUNT(*) FROM konto_wlasciciela WHERE potwierdzone = 0`); ile != 1 {
		t.Errorf("kont niepotwierdzonych po nieudanym nadaniu: %d, oczekiwane 1 —"+
			" adresu nikt nie sprawdził i nikt tego stanu nie ma prawa udawać", ile)
	}

	// Droga potwierdzenia zostaje i to nie jest usterka, sama wygasa po godzinie.

	// Dowód właściwy: platforma jest do otwarcia, nie tylko wiersze w bazie.
	var wejscie shared.AuthLoginResponse
	wykonajUdana(t, u.rdzen, u.zycie, shared.CommandAuthLogin, shared.AuthLoginRequest{
		Method: shared.AuthMethodKindPassword,
		Login:  wskaznik(loginSprawdzianu),
		Secret: wskaznik(hasloPierwsze),
	}, &wejscie)
	if strings.TrimSpace(wejscie.Session.Token) == "" {
		t.Error("wejście hasłem po nieudanym nadaniu oddało sesję bez tokenu")
	}

	// Druga rejestracja odmawia konfliktem i tak ma być: konto stoi, hasło je otwiera.

	// Powtórzone żądanie jest próbą podmiany hasła bez znajomości starego.
	powtorka := wykonajOdmowna(t, u.rdzen, u.zycie, shared.CommandAuthRegister,
		shared.AuthRegisterRequest{
			Login:    loginSprawdzianu,
			Email:    adresSprawdzianu,
			Password: hasloDrugie,
		})
	if powtorka.Code != shared.ErrorCodeConflict {
		t.Errorf("druga rejestracja niesie kod %q, oczekiwany %q",
			powtorka.Code, shared.ErrorCodeConflict)
	}
}

// TestWykazUrzadzenWidziWejscieHaslem pilnuje, że okno odbierania dostępu ma co
// pokazać: sesja brała urządzenie z wiersza metody, nie z żądania logowania.
func TestWykazUrzadzenWidziWejscieHaslem(t *testing.T) {
	u := zmontujDrogeWejscia(t, pocztaDziala)
	droga := zarejestrujWlasciciela(t, u)

	// Adres potwierdzony, bo bez tego bramka jest zamknięta.

	// Badane jest tu wejście hasłem, nie potwierdzenie.
	var potwierdzenie shared.AuthVerifyResponse
	wykonajUdana(t, u.rdzen, u.zycie, shared.CommandAuthVerify,
		shared.AuthVerifyRequest{Token: droga}, &potwierdzenie)

	const maszyna = "laptop-operatora"
	var wejscie shared.AuthLoginResponse
	wykonajUdana(t, u.rdzen, u.zycie, shared.CommandAuthLogin, shared.AuthLoginRequest{
		Method:   shared.AuthMethodKindPassword,
		Login:    wskaznik(loginSprawdzianu),
		Secret:   wskaznik(hasloPierwsze),
		DeviceId: wskaznik(maszyna),
	}, &wejscie)

	var wykaz shared.DeviceListResponse
	wykonajUdana(t, u.rdzen, u.zycie, shared.CommandDeviceList,
		shared.DeviceListRequest{}, &wykaz)

	znaleziona := false
	for _, urzadzenie := range wykaz.Devices {
		if urzadzenie.DeviceId != maszyna {
			continue
		}
		znaleziona = true
		if !urzadzenie.HasToken {
			t.Error("maszyna po wejściu hasłem jest w wykazie bez ważnego tokenu")
		}
	}
	if !znaleziona {
		t.Fatalf("maszyna %q weszła hasłem, a wykazu urządzeń nie widzi (pozycji: %d)",
			maszyna, len(wykaz.Devices))
	}

	// Skutek, nie koperta: unieważnienie ma naprawdę zamknąć sesję.
	var zdjecie shared.DeviceRevokeResponse
	wykonajUdana(t, u.rdzen, u.zycie, shared.CommandDeviceRevoke,
		shared.DeviceRevokeRequest{DeviceId: maszyna}, &zdjecie)
	if !zdjecie.Revoked {
		t.Error("unieważnienie maszyny z wykazu oddało revoked: false")
	}
	if ile := liczbaWierszy(t, u,
		`SELECT COUNT(*) FROM sesja_bramki WHERE urzadzenie_kod = ? AND uniewazniono IS NULL`,
		maszyna); ile != 0 {
		t.Errorf("po unieważnieniu maszyna ma %d czynnych sesji", ile)
	}
}
