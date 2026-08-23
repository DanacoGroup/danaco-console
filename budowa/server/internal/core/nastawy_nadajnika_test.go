package core

import (
	"context"
	"testing"

	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/nadajnik"
)

// nastawyPodstawione oddaje wartości z mapy — bez bazy, bez rozstrzygacza.
type nastawyPodstawione map[string]string

func (n nastawyPodstawione) Nastawa(_ context.Context, klucz string) string { return n[klucz] }

// TestKontoNadawczeBezZrodlaZostajePrzyStarcie pilnuje granicy: brak drogi
// odczytu nie jest wskazaniem Operatora i niczego nie kasuje.
func TestKontoNadawczeBezZrodlaZostajePrzyStarcie(t *testing.T) {
	startu := nadajnik.Nastawy{Host: "smtp.start", Port: 25, Adres: "platforma@start"}

	wynik := kontoNadawcze(context.Background(), startu, nil)

	if wynik != startu {
		t.Errorf("konto zmieniło się bez źródła nastaw: %+v", wynik)
	}
}

// TestNastawaPustaNieKasujeKontaZeStartu jest sednem nakładki: pierwsze wejście
// do okna Konfiguracji zastaje pola puste, a puste pole ma znaczyć „nie
// wskazałem”, nie „skasuj konto podane środowiskiem”.
func TestNastawaPustaNieKasujeKontaZeStartu(t *testing.T) {
	startu := nadajnik.Nastawy{
		Host: "smtp.start", Port: 587, Adres: "platforma@start",
		Uzytkownik: "operator", Sekret: "tajne", SzyfrujStartTLS: true,
	}

	wynik := kontoNadawcze(context.Background(), startu, nastawyPodstawione{})

	if wynik != startu {
		t.Errorf("puste nastawy skasowały konto ze startu: %+v", wynik)
	}
}

// TestNastawaOperatoraPrzeslaniaStart sprawdza drogę, dla której cała ta warstwa
// powstała: pomyłkę w adresie serwera poczty naprawia się bez zatrzymywania
// rdzenia.
func TestNastawaOperatoraPrzeslaniaStart(t *testing.T) {
	startu := nadajnik.Nastawy{Host: "smtp.bledny", Port: 587, Adres: "platforma@start", SzyfrujStartTLS: true}

	wynik := kontoNadawcze(context.Background(), startu, nastawyPodstawione{
		konfig.KluczNadawcaHost:     "smtp.wlasciwy",
		konfig.KluczNadawcaPort:     "2525",
		konfig.KluczNadawcaStartTLS: "false",
	})

	if wynik.Host != "smtp.wlasciwy" {
		t.Errorf("serwer poczty: %q, oczekiwany smtp.wlasciwy", wynik.Host)
	}
	if wynik.Port != 2525 {
		t.Errorf("port: %d, oczekiwany 2525", wynik.Port)
	}
	if wynik.SzyfrujStartTLS {
		t.Error("wskazanie „nie szyfruj” nie zadziałało — nastawa logiczna ma wygrywać w obie strony")
	}
	if wynik.Adres != "platforma@start" {
		t.Errorf("adres nadawcy zgubiony: %q", wynik.Adres)
	}
}

// TestPortNieczytelnyNieRozbrajaKontaNadawczego pilnuje, żeby zapis niebędący
// liczbą zostawił port ze startu zamiast zerować go w locie — port zerowy to
// list, który nigdzie nie idzie.
func TestPortNieczytelnyNieRozbrajaKontaNadawczego(t *testing.T) {
	startu := nadajnik.Nastawy{Host: "smtp.start", Port: 587, Adres: "platforma@start"}

	for _, zapis := range []string{"", "  ", "port", "0", "-1"} {
		wynik := kontoNadawcze(context.Background(), startu, nastawyPodstawione{
			konfig.KluczNadawcaPort: zapis,
		})
		if wynik.Port != 587 {
			t.Errorf("zapis %q dał port %d, oczekiwany 587 ze startu", zapis, wynik.Port)
		}
	}
}
