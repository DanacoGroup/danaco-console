package core

import (
	"strings"
	"testing"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// Skutek zdjęcia znacznika bramki bez poczty: warunek wejścia znika, gdy adres zostaje potwierdzony.

// TestPotwierdzenieAdresuZdejmujeZnacznikPostawionyBezPoczty prowadzi całą drogę bez poczty do końca: znacznik powstaje przy rejestracji, a potwierdzenie adresu go kasuje. Dowodem jest odczyt sejfu po potwierdzeniu, nie pole verified.
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

	// Bramka zostaje otwarta: potwierdzenie adresu przestawia ją z wyjątku na wiersz konta.
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

// wydajDrogeWeryfikacji zakłada drogę potwierdzenia adresu i oddaje materiał, który trzeba przepisać z listu, zapisując go tak samo jak funkcja wyslijDrogePotwierdzenia: skrót materiału, cel weryfikacji, godzinę ważności.
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
