// Rozstrzyga dostęp okna do plików i do sieci: własny katalog okna jest
// dozwolony zawsze, poza nim wymaga aktywnego nadania w trybie zgodnym
// z wykonywaną operacją, a wyłączony zakres przepuszcza wszystko.
package core

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// nadanieSprawdzianu składa parę punkt + nadanie dla katalogu lokalnego,
// gotową do użycia w przypadkach sprawdzających granicę dostępu plikowego.
func nadanieSprawdzianu(korzenie []string, tryb shared.AccessMode,
	punktCzynny, nadanieCzynne bool) NadanieMostu {

	return NadanieMostu{
		Punkt: shared.AccessPoint{
			Id:      "punkt-katalog",
			Name:    "Katalog projektu",
			Kind:    shared.AccessPointKindLocalDirectory,
			Roots:   korzenie,
			Enabled: punktCzynny,
		},
		Nadanie: shared.AccessGrant{
			Id:            "nadanie-katalog",
			WindowId:      "okno-pierwsze",
			AccessPointId: "punkt-katalog",
			Mode:          tryb,
			Enabled:       nadanieCzynne,
		},
	}
}

// TestStrazPlikowPrzepuszczaWlasnyKatalogOkna sprawdza punkt wyjścia: własny
// katalog roboczy okna jest dozwolony bez żadnego nadania.
func TestStrazPlikowPrzepuszczaWlasnyKatalogOkna(t *testing.T) {
	wlasny := t.TempDir()
	straz := nowaStrazPlikow(session.Zasady{Pliki: true}, wlasny, nil)

	for _, sciezka := range []string{
		wlasny,
		filepath.Join(wlasny, "plik.txt"),
		filepath.Join(wlasny, "podkatalog", "glebiej", "plik.txt"),
	} {
		if err := straz.sprawdzOdczyt(sciezka); err != nil {
			t.Errorf("odczyt %q odrzucony: %v", sciezka, err)
		}
		if err := straz.sprawdzZapis(sciezka); err != nil {
			t.Errorf("zapis %q odrzucony: %v", sciezka, err)
		}
	}
}

// TestStrazPlikowOdrzucaSciezkeSpozaObszaru mierzy odmowę zwykłą oraz odmowę
// przy wyjściu w górę drzewa. Wyjście w górę jest przypadkiem właściwym: ścieżka
// złożona z członów modelu dochodzi do straży dokładnie w tej postaci.
func TestStrazPlikowOdrzucaSciezkeSpozaObszaru(t *testing.T) {
	wlasny := t.TempDir()
	straz := nowaStrazPlikow(session.Zasady{Pliki: true}, wlasny, nil)

	przypadki := map[string]string{
		"katalog obcy":            filepath.Join(t.TempDir(), "cudze.txt"),
		"wyjście w górę drzewa":   filepath.Join(wlasny, "..", "cudze.txt"),
		"wyjście dwa razy w górę": filepath.Join(wlasny, "..", "..", "etc", "passwd"),
		"ścieżka pusta":           "",
		"same odstępy":            "   ",
	}
	for nazwa, sciezka := range przypadki {
		t.Run(nazwa, func(t *testing.T) {
			if err := straz.sprawdzOdczyt(sciezka); err == nil {
				t.Errorf("odczyt %q przeszedł", sciezka)
			}
			if err := straz.sprawdzZapis(sciezka); err == nil {
				t.Errorf("zapis %q przeszedł", sciezka)
			}
		})
	}
}

// TestNazwaPodobnaDoKorzeniaNieWchodzi pilnuje porównania po członach ścieżki,
// nie po przedrostku napisu. Katalog `dane-cudze` obok `dane` nie leży w `dane`,
// choć zaczyna się tak samo.
func TestNazwaPodobnaDoKorzeniaNieWchodzi(t *testing.T) {
	nadrzedny := t.TempDir()
	wlasny := filepath.Join(nadrzedny, "dane")
	podobny := filepath.Join(nadrzedny, "dane-cudze", "plik.txt")

	straz := nowaStrazPlikow(session.Zasady{Pliki: true}, wlasny, nil)
	if err := straz.sprawdzOdczyt(podobny); err == nil {
		t.Errorf("ścieżka %q przeszła przez porównanie po przedrostku napisu", podobny)
	}
}

// TestTrybNadaniaRozdzielaOdczytOdZapisu jest sprawdzianem samego trybu: korzeń
// nadany do odczytu nie staje się korzeniem do zapisu. Odmowa ma przy tym
// nazywać powód wprost — Operator ma wiedzieć, że brakuje trybu, a nie nadania.
func TestTrybNadaniaRozdzielaOdczytOdZapisu(t *testing.T) {
	wlasny := t.TempDir()
	doOdczytu := t.TempDir()
	doZapisu := t.TempDir()

	straz := nowaStrazPlikow(session.Zasady{Pliki: true}, wlasny, []NadanieMostu{
		nadanieSprawdzianu([]string{doOdczytu}, shared.AccessModeRead, true, true),
		nadanieSprawdzianu([]string{doZapisu}, shared.AccessModeWrite, true, true),
	})

	plikDoOdczytu := filepath.Join(doOdczytu, "ustalenie.md")
	if err := straz.sprawdzOdczyt(plikDoOdczytu); err != nil {
		t.Errorf("odczyt korzenia nadanego do odczytu odrzucony: %v", err)
	}
	err := straz.sprawdzZapis(plikDoOdczytu)
	if err == nil {
		t.Fatal("zapis w korzeniu nadanym wyłącznie do odczytu przeszedł")
	}
	if !strings.Contains(err.Error(), "wyłącznie do odczytu") {
		t.Errorf("odmowa nie nazywa braku trybu: %v", err)
	}

	plikDoZapisu := filepath.Join(doZapisu, "wynik.md")
	if err := straz.sprawdzOdczyt(plikDoZapisu); err != nil {
		t.Errorf("odczyt korzenia nadanego do zapisu odrzucony: %v", err)
	}
	if err := straz.sprawdzZapis(plikDoZapisu); err != nil {
		t.Errorf("zapis w korzeniu nadanym do zapisu odrzucony: %v", err)
	}
}

// TestNadanieNieczynneNiczegoNieOtwiera sprawdza trzy drogi wygaszenia nadania.
// Brak wiersza znaczy brak uprawnienia — każda z tych dróg ma zamykać tak samo
// szczelnie jak brak nadania w ogóle.
func TestNadanieNieczynneNiczegoNieOtwiera(t *testing.T) {
	wlasny := t.TempDir()
	korzen := t.TempDir()
	plik := filepath.Join(korzen, "ustalenie.md")

	przypadki := map[string]NadanieMostu{
		"nadanie wygaszone": nadanieSprawdzianu([]string{korzen}, shared.AccessModeWrite, true, false),
		"punkt wygaszony":   nadanieSprawdzianu([]string{korzen}, shared.AccessModeWrite, false, true),
	}
	for nazwa, nadanie := range przypadki {
		t.Run(nazwa, func(t *testing.T) {
			straz := nowaStrazPlikow(session.Zasady{Pliki: true}, wlasny, []NadanieMostu{nadanie})
			if err := straz.sprawdzOdczyt(plik); err == nil {
				t.Error("odczyt przeszedł mimo wygaszenia")
			}
		})
	}

	// Punkt innego rodzaju nie wnosi korzeni plikowych, choćby je miał.
	mostowe := nadanieSprawdzianu([]string{korzen}, shared.AccessModeWrite, true, true)
	mostowe.Punkt.Kind = shared.AccessPointKindMcpBridge
	straz := nowaStrazPlikow(session.Zasady{Pliki: true}, wlasny, []NadanieMostu{mostowe})
	if err := straz.sprawdzOdczyt(plik); err == nil {
		t.Error("korzeń punktu mostowego otworzył dostęp plikowy")
	}
}

// TestKorzenieNadaniaZawezajaSieDoKorzeniPunktu sprawdza regułę zawężenia:
// nadanie nie sięga poza korzenie punktu. Bez tego wskazanie korzenia w nadaniu
// byłoby drogą obejścia punktu — Operator ustawia zakres w jednym miejscu,
// a rozszerza go w drugim.
func TestKorzenieNadaniaZawezajaSieDoKorzeniPunktu(t *testing.T) {
	wlasny := t.TempDir()
	korzenPunktu := t.TempDir()
	korzenObcy := t.TempDir()

	nadanie := nadanieSprawdzianu([]string{korzenPunktu}, shared.AccessModeWrite, true, true)
	nadanie.Nadanie.Roots = []string{korzenPunktu, korzenObcy}

	straz := nowaStrazPlikow(session.Zasady{Pliki: true}, wlasny, []NadanieMostu{nadanie})

	if err := straz.sprawdzOdczyt(filepath.Join(korzenPunktu, "plik.md")); err != nil {
		t.Errorf("korzeń wskazany w punkcie i w nadaniu odrzucony: %v", err)
	}
	if err := straz.sprawdzOdczyt(filepath.Join(korzenObcy, "plik.md")); err == nil {
		t.Error("korzeń wskazany wyłącznie w nadaniu rozszerzył zakres punktu")
	}
}

// TestPunktBezKorzeniOddajeZawezenieNadaniu utrwala drugą stronę tej samej
// reguły: punkt bez korzeni obejmuje cały system plików maszyny, więc korzenie
// nadania są jedynym zawężeniem i wchodzą w całości.
func TestPunktBezKorzeniOddajeZawezenieNadaniu(t *testing.T) {
	wlasny := t.TempDir()
	korzenNadania := t.TempDir()

	nadanie := nadanieSprawdzianu(nil, shared.AccessModeRead, true, true)
	nadanie.Nadanie.Roots = []string{korzenNadania}

	straz := nowaStrazPlikow(session.Zasady{Pliki: true}, wlasny, []NadanieMostu{nadanie})
	if err := straz.sprawdzOdczyt(filepath.Join(korzenNadania, "plik.md")); err != nil {
		t.Errorf("korzeń nadania przy punkcie bez korzeni odrzucony: %v", err)
	}
	if err := straz.sprawdzOdczyt(filepath.Join(t.TempDir(), "plik.md")); err == nil {
		t.Error("punkt bez korzeni otworzył cały system plików mimo zawężenia w nadaniu")
	}
}

// TestZakresWylaczonyPrzepuszczaWszystko utrwala stan wyjściowy platformy. Jest
// to zachowanie zamierzone i sprawdzian ma je trzymać — nie po to, by je
// chwalić, lecz po to, by jego zmiana była zmianą świadomą.
func TestZakresWylaczonyPrzepuszczaWszystko(t *testing.T) {
	straz := nowaStrazPlikow(session.Zasady{Pliki: false}, t.TempDir(), nil)
	obca := filepath.Join(t.TempDir(), "cudze.txt")

	if err := straz.sprawdzOdczyt(obca); err != nil {
		t.Errorf("zakres wyłączony odrzucił odczyt: %v", err)
	}
	if err := straz.sprawdzZapis(obca); err != nil {
		t.Errorf("zakres wyłączony odrzucił zapis: %v", err)
	}

	var zerowa *strazPlikow
	if err := zerowa.sprawdzOdczyt(obca); err != nil {
		t.Errorf("straż niewpięta odrzuciła odczyt: %v", err)
	}
}

// TestWykazKatalogowNieJestObcinanyPoCichu sprawdza obietnicę z opisu: katalog
// spoza obszaru okna zatrzymuje uruchomienie, zamiast zniknąć z wykazu. Wykaz
// obcięty po cichu dawałby Operatorowi obraz niezgodny z tym, co dostał proces.
func TestWykazKatalogowNieJestObcinanyPoCichu(t *testing.T) {
	wlasny := t.TempDir()
	obcy := t.TempDir()
	straz := nowaStrazPlikow(session.Zasady{Pliki: true}, wlasny, nil)

	if wykaz, err := straz.Katalogi([]string{wlasny}); err != nil {
		t.Errorf("wykaz z samym katalogiem własnym odrzucony: %v", err)
	} else if len(wykaz) != 1 {
		t.Errorf("wykaz oddał %d pozycji, oczekiwana jedna", len(wykaz))
	}

	wykaz, err := straz.Katalogi([]string{wlasny, obcy})
	if err == nil {
		t.Fatalf("wykaz z katalogiem obcym przeszedł: %v", wykaz)
	}
	if wykaz != nil {
		t.Errorf("odmowa oddała wykaz obcięty: %v", wykaz)
	}
}

// mostSprawdzianu składa parę punkt + nadanie dla mostu MCP, gotową do
// użycia w przypadkach sprawdzających granicę dostępu sieciowego.
func mostSprawdzianu(host string, czynne bool) NadanieMostu {
	return NadanieMostu{
		Punkt: shared.AccessPoint{
			Id:      "punkt-most-" + host,
			Name:    "Most " + host,
			Kind:    shared.AccessPointKindMcpBridge,
			Host:    &host,
			Enabled: czynne,
		},
		Nadanie: shared.AccessGrant{
			Id:            "nadanie-most-" + host,
			WindowId:      "okno-pierwsze",
			AccessPointId: "punkt-most-" + host,
			Mode:          shared.AccessModeRead,
			Enabled:       czynne,
		},
	}
}

// TestStrazSieciPrzepuszczaWylacznieMaszynyNadane jest sprawdzianem granicy
// wyjścia rdzenia na świat. Wykaz maszyn bierze się z nadań okna, więc maszyna
// spoza nadań ma odpaść niezależnie od tego, w jakim zapisie przyszła.
func TestStrazSieciPrzepuszczaWylacznieMaszynyNadane(t *testing.T) {
	straz := nowaStrazSieci(session.Zasady{DostepSieciowy: true},
		[]NadanieMostu{mostSprawdzianu("maszyna-nadana", true)})

	przechodzace := []string{
		"maszyna-nadana",
		"maszyna-nadana:8080",
		"https://maszyna-nadana/sciezka",
		"https://maszyna-nadana:443/sciezka?pytanie=1",
		"ssh://operator@maszyna-nadana:22",
		"MASZYNA-NADANA",
	}
	for _, adres := range przechodzace {
		t.Run("przechodzi "+adres, func(t *testing.T) {
			if err := straz.sprawdzAdres(adres); err != nil {
				t.Errorf("adres maszyny nadanej odrzucony: %v", err)
			}
		})
	}

	odpadajace := []string{
		"maszyna-obca",
		"https://maszyna-obca/sciezka",
		"maszyna-nadana.example",
		"operator@maszyna-obca:22",
		"",
		"   ",
	}
	for _, adres := range odpadajace {
		t.Run("odpada "+adres, func(t *testing.T) {
			if err := straz.sprawdzAdres(adres); err == nil {
				t.Errorf("adres %q spoza nadań przeszedł", adres)
			}
		})
	}
}

// TestMaszynaZWygaszonegoNadaniaOdpada pilnuje tej samej reguły co przy plikach:
// wygaszenie zamyka tak samo szczelnie jak brak nadania.
func TestMaszynaZWygaszonegoNadaniaOdpada(t *testing.T) {
	straz := nowaStrazSieci(session.Zasady{DostepSieciowy: true},
		[]NadanieMostu{mostSprawdzianu("maszyna-wygaszona", false)})

	if err := straz.sprawdzAdres("maszyna-wygaszona"); err == nil {
		t.Error("maszyna z wygaszonego nadania przeszła")
	}
}

// TestZakresSieciowyWylaczonyPrzepuszczaWszystko utrwala stan wyjściowy —
// wspólny dostęp sieciowy serwera.
func TestZakresSieciowyWylaczonyPrzepuszczaWszystko(t *testing.T) {
	straz := nowaStrazSieci(session.Zasady{DostepSieciowy: false}, nil)
	if err := straz.sprawdzAdres("maszyna-obca"); err != nil {
		t.Errorf("zakres wyłączony odrzucił adres: %v", err)
	}

	var zerowa *strazSieci
	if err := zerowa.sprawdzAdres("maszyna-obca"); err != nil {
		t.Errorf("straż niewpięta odrzuciła adres: %v", err)
	}
}

// TestWykazMostowNieJestObcinanyPoCichu sprawdza drugą stronę tej samej
// obietnicy co przy katalogach: most spoza nadań zatrzymuje złożenie
// konfiguracji, zamiast zniknąć z wykazu.
func TestWykazMostowNieJestObcinanyPoCichu(t *testing.T) {
	nadany := mostSprawdzianu("maszyna-nadana", true)
	straz := nowaStrazSieci(session.Zasady{DostepSieciowy: true}, []NadanieMostu{nadany})

	obcyHost := "maszyna-obca"
	obcy := shared.AccessPoint{
		Id: "punkt-most-obcy", Kind: shared.AccessPointKindMcpBridge,
		Host: &obcyHost, Enabled: true,
	}

	if wykaz, err := straz.mosty([]shared.AccessPoint{nadany.Punkt}); err != nil {
		t.Errorf("wykaz z samym mostem nadanym odrzucony: %v", err)
	} else if len(wykaz) != 1 {
		t.Errorf("wykaz oddał %d pozycji, oczekiwana jedna", len(wykaz))
	}

	wykaz, err := straz.mosty([]shared.AccessPoint{nadany.Punkt, obcy})
	if err == nil {
		t.Fatalf("wykaz z mostem spoza nadań przeszedł: %v", wykaz)
	}
	if wykaz != nil {
		t.Errorf("odmowa oddała wykaz obcięty: %v", wykaz)
	}
}

// TestPunktInnegoRodzajuNieWchodziDoWykazuMaszyn sprawdza, że katalog lokalny
// nie wnosi maszyny do wykazu sieciowego — dwa zakresy izolacji nie mieszają się
// przez wspólny wykaz nadań.
func TestPunktInnegoRodzajuNieWchodziDoWykazuMaszyn(t *testing.T) {
	katalogowe := nadanieSprawdzianu([]string{t.TempDir()}, shared.AccessModeRead, true, true)
	host := "maszyna-z-katalogu"
	katalogowe.Punkt.Host = &host

	straz := nowaStrazSieci(session.Zasady{DostepSieciowy: true}, []NadanieMostu{katalogowe})
	if err := straz.sprawdzAdres(host); err == nil {
		t.Error("punkt katalogu lokalnego wniósł maszynę do wykazu sieciowego")
	}
}

// TestOdmowaIzolacjiNazywaZakres pilnuje treści odmowy. Naruszenie bez nazwy
// zakresu nie prowadzi Operatora do nastawy, którą trzeba zmienić — a to jest
// jedyna droga naprawy.
func TestOdmowaIzolacjiNazywaZakres(t *testing.T) {
	strazP := nowaStrazPlikow(session.Zasady{Pliki: true}, t.TempDir(), nil)
	err := strazP.sprawdzOdczyt(filepath.Join(t.TempDir(), "cudze.txt"))
	if err == nil {
		t.Fatal("ścieżka obca przeszła")
	}
	if !strings.Contains(err.Error(), konfig.KluczIzolacjaPliki) {
		t.Errorf("odmowa plikowa nie nazywa zakresu: %v", err)
	}

	strazS := nowaStrazSieci(session.Zasady{DostepSieciowy: true}, nil)
	err = strazS.sprawdzAdres("maszyna-obca")
	if err == nil {
		t.Fatal("adres obcy przeszedł")
	}
	if !strings.Contains(err.Error(), konfig.KluczIzolacjaDostepSieciowy) {
		t.Errorf("odmowa sieciowa nie nazywa zakresu: %v", err)
	}
	if !errors.Is(err, session.ErrIzolacja) {
		t.Error("odmowa izolacji nie wiąże się ze wspólnym korzeniem błędu")
	}
}
