package core

import (
	"strings"
	"testing"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// Skutek zdjęcia znacznika bramki bez poczty.
//
// Znacznik `auth:bramka-bez-poczty` zdejmuje jeden warunek wejścia i tylko
// jeden: bramki nie zamyka potwierdzenie, którego platforma nie miała czym
// wysłać (rejestr decyzji, pozycja 11). Wyjątek ma więc trwać dokładnie tak
// długo, jak trwa jego powód — a powód kończy się w chwili, w której adres
// zostaje potwierdzony.
//
// Znacznik, który przeżyje potwierdzenie, przestaje być wyjątkiem pierwszego
// uruchomienia i staje się trwałym obejściem bramki: od tej chwili konto
// cofnięte do stanu niepotwierdzonego wchodziłoby hasłem mimo działającej
// poczty. Dlatego mierzona jest tu strona, której nie mierzy granica istnienia
// znacznika w `skutek_wejscia_test.go`: tamta prowadzi drogę Z POCZTĄ, gdzie
// znacznik nie powstaje w ogóle, więc przechodzi także wtedy, gdy rdzeń
// znacznika nie zdejmuje.

// TestPotwierdzenieAdresuZdejmujeZnacznikPostawionyBezPoczty prowadzi całą drogę
// bez poczty do końca: znacznik powstaje przy rejestracji, a potwierdzenie
// adresu go kasuje.
//
// Dowodem jest odczyt sejfu po potwierdzeniu, nie `verified: true` — bramka
// wpuszczająca hasłem i bramka trzymana wierszem konta oddają w kopercie to samo.
func TestPotwierdzenieAdresuZdejmujeZnacznikPostawionyBezPoczty(t *testing.T) {
	u := zmontujDrogeWejscia(t, pocztaBrak)

	var rejestracja shared.AuthRegisterResponse
	wykonajUdana(t, u.rdzen, u.zycie, shared.CommandAuthRegister,
		shared.AuthRegisterRequest{
			Login:    loginSprawdzianu,
			Email:    adresSprawdzianu,
			Password: hasloPierwsze,
		}, &rejestracja)

	if _, jest := znacznikWSejfie(t, u); !jest {
		t.Fatal("rejestracja bez poczty nie postawiła znacznika — dalszy pomiar mierzyłby" +
			" zdejmowanie czegoś, czego nie ma, i meldował to jako wynik")
	}

	droga := wydajDrogeWeryfikacji(t, u)

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
		t.Errorf("po potwierdzeniu adresu w sejfie dalej leży znacznik bramki bez poczty"+
			" (adres %q) — wyjątek pierwszego uruchomienia przeżył swój powód i zdejmuje"+
			" warunek potwierdzenia na instalce, która ten adres właśnie potwierdziła", adres)
	}

	// Bramka zostaje otwarta: potwierdzenie adresu przestawia ją z wyjątku na
	// wiersz konta, a nie zamyka wejścia, które przed nim działało.
	var wejscie shared.AuthLoginResponse
	wykonajUdana(t, u.rdzen, u.zycie, shared.CommandAuthLogin, shared.AuthLoginRequest{
		Method: shared.AuthMethodKindPassword,
		Login:  wskaznik(loginSprawdzianu),
		Secret: wskaznik(hasloPierwsze),
	}, &wejscie)
	if strings.TrimSpace(wejscie.Session.Token) == "" {
		t.Error("po potwierdzeniu adresu wejście hasłem oddało sesję bez tokenu")
	}
}

// wydajDrogeWeryfikacji zakłada drogę potwierdzenia adresu i oddaje materiał,
// który Operator przepisałby z listu.
//
// Idzie warstwą danych, nie komendą, bo na instalce bez poczty żadna komenda tej
// drogi nie wydaje: rejestracja wykonuje się raz i list wyszedłby tylko z niej,
// a `auth.recover` wydaje drogę do innego celu, której `auth.verify` nie
// przyjmuje (wykazuje to `TestBezPocztyPotwierdzenieAdresuCzekaNaDrogeZListu`).
// Brak tej drogi jest zgłoszeniem rejestru terenów, nie przedmiotem tego pomiaru
// — mierzone jest to, co rdzeń robi, GDY adres zostaje potwierdzony.
//
// Zapis idzie dokładnie tak, jak robi to `wyslijDrogePotwierdzenia`: skrót
// materiału, cel weryfikacji, godzina ważności. Sprawdzian trzyma sam materiał,
// bo z bazy odczytać się go nie da — leży tam wyłącznie jego skrót.
func wydajDrogeWeryfikacji(t *testing.T, u uprzazWejscia) string {
	t.Helper()

	droga, err := nowyTokenBramki()
	if err != nil {
		t.Fatalf("nie można wylosować drogi potwierdzenia: %v", err)
	}
	teraz := time.Now()
	_, err = u.baza.DB.ExecContext(u.zycie,
		`INSERT INTO potwierdzenie_tozsamosci (skrot, cel, wygasa, utworzono) VALUES (?, ?, ?, ?)`,
		skrotTokenu(droga), dane.CelWeryfikacja,
		teraz.Add(trwanieDrogiPotwierdzenia).UnixMilli(), teraz.UnixMilli())
	if err != nil {
		t.Fatalf("nie można założyć drogi potwierdzenia: %v", err)
	}
	return droga
}
