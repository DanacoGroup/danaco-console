package core

import (
	"strings"
	"testing"

	"danacoconsole/shared"
)

// STRAŻE DROGI WEJŚCIA.
//
// Trzy sprawdziany poniżej trzymają zachowania, których zerwanie zamyka
// Operatorowi drogę do platformy albo otwiera ją komuś, kto nie powinien wejść.
// Każdy niesie przy sobie zdanie „czym się to łamie" — bo straż bez opisu wagi
// wygląda na przesadę i pierwszy, kto ją zobaczy, uzna ją za zbędną.
//
// Trzy razem, a nie osobno, bo trzymają się nawzajem: bramka zamknięta do
// potwierdzenia adresu bez drogi wyjścia z nieudanego nadania zamieniałaby
// zatrzymany serwer poczty w trwałą utratę produktu, a wykaz urządzeń bez
// wejścia hasłem nie miałby czego pokazać w oknie odbierania dostępu.

// ── STRAŻ PIERWSZA ───────────────────────────────────────────────────────────

// TestBramkaZamknietaDoPotwierdzeniaAdresu pilnuje, że weryfikacja adresu nie
// jest ozdobą.
//
// CZYM SIĘ TO ZŁAMAŁO. Rejestracja nie zakładała sesji — ale zakładała kotwicę,
// a `auth.login` nie pytał o stan potwierdzenia w ogóle. Operator wołał więc
// logowanie zaraz po rejestracji i dostawał pełny token, nie zaglądając do
// skrzynki.
//
// DLACZEGO TO WAŻY. Adres jest JEDYNĄ drogą odzyskania konta. Adres
// niesprawdzony — literówka, cudza skrzynka, domena bez rekordu — wychodziłby na
// jaw dopiero w dniu, w którym trzeba nim odzyskać dostęp, czyli gdy jest już za
// późno: rejestracji nie da się powtórzyć.
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
	// „Nie wolno" bez drogi dalszej zostawia Operatora przed zamkniętą bramką
	// bez klucza.
	if !strings.Contains(blad.Message, "auth.verify") {
		t.Errorf("odmowa nie wskazuje drogi potwierdzenia: %q", blad.Message)
	}

	// Odmowa nie ma prawa niczego zakładać po drodze.
	if ile := liczbaWierszy(t, u,
		`SELECT COUNT(*) FROM sesja_bramki WHERE uniewazniono IS NULL`); ile != 0 {
		t.Errorf("po odmówionym wejściu leży %d czynnych sesji", ile)
	}
}

// ── STRAŻ DRUGA ──────────────────────────────────────────────────────────────

// TestNieudaneNadanieListuSchodziNaDrogeBezPoczty pilnuje, że zatrzymany
// przekaźnik poczty nie zamienia się w trwałą utratę platformy.
//
// CZYM SIĘ TO ŁAMIE. Sprawdzana bywała wyłącznie OBECNOŚĆ nastaw konta
// nadawczego — i owszem, przed zapisem. Samo nadanie idzie ostatnie, już po
// zapisaniu konta, kotwicy i drogi, a nastawa wskazana nie znaczy, że serwer
// odpowiada: przekaźnik bywa zatrzymany, zapora zamknięta, a nazwa hosta
// wpisana z literówką.
//
// CZEGO TU NIE MA I DLACZEGO. Cofnięcia rejestracji. Było ono ratunkiem przed
// platformą NIE DO OTWARCIA — bramkę zamykał wtedy brak potwierdzenia adresu,
// więc konto zostawione po nieudanym nadaniu nie miało czym wejść. Odkąd konto
// bez potwierdzonego adresu wchodzi hasłem (rejestr decyzji, pozycja 11),
// ratunek jest zbędny, a sam był pułapką: literówka w nazwie hosta zamykała
// pierwsze uruchomienie równie szczelnie jak brak poczty w ogóle. Rejestracja
// schodzi więc na drogę bez poczty i kończy się tym samym stanem, co instalka,
// która nadajnika nie ma wcale.
//
// DLACZEGO TO WAŻY. Rejestracja wykonuje się raz. Gdyby nieudane nadanie
// cofało ją bez otwarcia drogi powrotu albo zostawiało konto bez klucza,
// pierwszy Operator tracił platformę na jedną niedostępność serwera poczty.
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

	// Znacznik jest tu jedyną rzeczą, która trzyma bramkę otwartą: adres został
	// niepotwierdzony, bo list nie doszedł.
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

	// Droga potwierdzenia ZOSTAJE i to nie jest usterka: leży jako sam skrót
	// materiału, który do nikogo nie dojechał, i wygasa po godzinie. Kasowanie
	// jej wymagałoby czwartej czynności repozytorium dla stanu, który sam się
	// kończy.

	// Dowód właściwy: platforma jest do otwarcia. Bez tego straż pilnowałaby
	// wierszy w bazie, a nie tego, po co one tam są.
	var wejscie shared.AuthLoginResponse
	wykonajUdana(t, u.rdzen, u.zycie, shared.CommandAuthLogin, shared.AuthLoginRequest{
		Method: shared.AuthMethodKindPassword,
		Login:  wskaznik(loginSprawdzianu),
		Secret: wskaznik(hasloPierwsze),
	}, &wejscie)
	if strings.TrimSpace(wejscie.Session.Token) == "" {
		t.Error("wejście hasłem po nieudanym nadaniu oddało sesję bez tokenu")
	}

	// Druga rejestracja odmawia konfliktem i tak ma być: konto stoi, hasło je
	// otwiera, a powtórzone żądanie jest próbą podmiany hasła bez znajomości
	// starego.
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

// ── STRAŻ TRZECIA ────────────────────────────────────────────────────────────

// TestWykazUrzadzenWidziWejscieHaslem pilnuje, że okno odbierania dostępu ma co
// pokazać.
//
// CZYM SIĘ TO ZŁAMAŁO. Sesja brała urządzenie z wiersza metody, nie z żądania.
// Kotwica hasła urządzenia nie ma i mieć nie może — hasło nie jest materiałem
// jednej maszyny — więc `deviceId` z `auth.login` był porzucany.
//
// DLACZEGO TO WAŻY. Maszyna nie pojawiała się w `device.list`, a `device.revoke`
// z jej identyfikatorem wracał `revoked: false` i zostawiał token czynny —
// w oknie, które istnieje po to, żeby Operator dostęp odbierał.
func TestWykazUrzadzenWidziWejscieHaslem(t *testing.T) {
	u := zmontujDrogeWejscia(t, pocztaDziala)
	droga := zarejestrujWlasciciela(t, u)

	// Adres potwierdzony, bo bez tego bramka jest zamknięta (straż pierwsza),
	// a badane jest wejście hasłem, nie potwierdzenie.
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
