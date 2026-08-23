package core

import (
	"strings"
	"testing"

	"danacoconsole/shared"
)

// STRAŻE DROGI WEJŚCIA.
//
// Trzy sprawdziany poniżej powstały jako ZAPORY opisujące usterki: wypadały
// niepomyślnie i wtedy, gdy usterka się pogłębi, i wtedy, gdy zostanie
// naprawiona. Wszystkie trzy usterki są dziś naprawione, więc — zgodnie
// z własnym poleceniem tamtych zapór — zostały ODWRÓCONE w straże pilnujące
// stanu naprawionego.
//
// Zapisu nie usunięto, bo trzy rzeczy warto trzymać razem: co produkt obiecuje,
// czym się to złamało i czym jest trzymane teraz. Straż bez tej pamięci wygląda
// na przesadę i pierwszy, kto ją zobaczy, uzna ją za zbędną.

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

// TestNieudaneNadanieListuCofaRejestracje pilnuje, że zatrzymany przekaźnik
// poczty nie zamienia się w trwałą utratę platformy.
//
// CZYM SIĘ TO ZŁAMAŁO. Sprawdzana była wyłącznie OBECNOŚĆ nastaw konta
// nadawczego — i owszem, przed zapisem. Samo nadanie szło ostatnie, już po
// zapisaniu konta, kotwicy i drogi, a cofnięcia nie było.
//
// DLACZEGO TO WAŻY. Rejestracja wykonuje się raz. Konto zostawione po nieudanym
// nadaniu było platformą nie do otwarcia: drugiej rejestracji nie ma, wejść nie
// ma czym, bo adresu nikt nie potwierdził, a nowej drogi weryfikacji nie wysyła
// żadna komenda. Naprawa samej straży pierwszej BEZ tej zamieniłaby zatrzymany
// serwer poczty w trwałą utratę produktu przy pierwszym uruchomieniu — te dwie
// trzymają się nawzajem i tak trzeba je czytać.
func TestNieudaneNadanieListuCofaRejestracje(t *testing.T) {
	u := zmontujDrogeWejscia(t, pocztaNieosiagalna)

	blad := wykonajOdmowna(t, u.rdzen, u.zycie, shared.CommandAuthRegister,
		shared.AuthRegisterRequest{
			Login:    loginSprawdzianu,
			Email:    adresSprawdzianu,
			Password: hasloPierwsze,
		})
	if !strings.Contains(blad.Message, "nie udało się wysłać listu") {
		t.Fatalf("odmowa nie mówi o nieudanym nadaniu: %q", blad.Message)
	}

	konta := liczbaWierszy(t, u, `SELECT COUNT(*) FROM konto_wlasciciela`)
	kotwice := liczbaWierszy(t, u, `SELECT COUNT(*) FROM metoda_uwierzytelnienia WHERE kotwica = 1`)

	// Te dwa wiersze rozstrzygają o tym, czy Operator może spróbować ponownie:
	// konto blokuje drugą rejestrację warunkiem schematu, kotwica — odmową
	// „konto już założone". Oba muszą zniknąć.
	if konta != 0 {
		t.Errorf("po nieudanym nadaniu w bazie leży %d kont — druga rejestracja odbije się o nie", konta)
	}
	if kotwice != 0 {
		t.Errorf("po nieudanym nadaniu w bazie leży %d kotwic — druga rejestracja odmówi konfliktem", kotwice)
	}

	// Droga potwierdzenia MOŻE zostać i to nie jest usterka: leży jako sam skrót
	// materiału, który do nikogo nie dojechał, wygasa po godzinie, a bez konta
	// nie ma czego otworzyć. Kasowanie jej wymagałoby czwartej czynności
	// repozytorium dla stanu, który sam się kończy.

	// Dowód właściwy: droga powrotu MA BYĆ OTWARTA. Bez tego straż pilnowałaby
	// pustych tabel, a nie tego, po co je opróżniono.
	//
	// Poczta w tej uprzęży pozostaje nieosiągalna, więc druga rejestracja też
	// odmówi — ale MUSI odmówić z powodu nadania, nie konfliktem. Konflikt
	// znaczyłby, że konto wciąż stoi i Operator nie ma jak spróbować ponownie po
	// naprawieniu serwera poczty.
	powtorka := wykonajOdmowna(t, u.rdzen, u.zycie, shared.CommandAuthRegister,
		shared.AuthRegisterRequest{
			Login:    loginSprawdzianu,
			Email:    adresSprawdzianu,
			Password: hasloPierwsze,
		})
	if powtorka.Code == shared.ErrorCodeConflict {
		t.Error("druga rejestracja odmawia konfliktem — cofnięcie nie otworzyło drogi powrotu")
	}
	if !strings.Contains(powtorka.Message, "nie udało się wysłać listu") {
		t.Errorf("druga rejestracja odmawia z innego powodu niż nadanie: %q", powtorka.Message)
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
