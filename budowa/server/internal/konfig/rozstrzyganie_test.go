package konfig

import (
	"errors"
	"testing"

	"danacoconsole/shared"
)

func kontekstPelny() Kontekst {
	return Kontekst{
		Srodowisko:  "srodowisko-talkin",
		Modul:       "modul-assistant",
		ParaModulow: "para-assistant-terminal",
		Projekt:     "projekt-oferta",
		KartaSesji:  "sesja-pierwsza",
		Rola:        "rola-wykonawca",
		Okno:        "okno-pierwsze",
		Model:       "model-glowny",
		Konto:       "konto-pierwsze",
	}
}

// TestPoziomyIdaOdNajwezszegoDoNajszerszego utrwala kolejność, na której stoi
// cała reguła. Wykaz jest wiążący: przestawienie dwóch pozycji zmienia zachowanie
// każdego ustawienia, a kompilacja tego nie zauważy.
func TestPoziomyIdaOdNajwezszegoDoNajszerszego(t *testing.T) {
	oczekiwane := []Poziom{
		shared.ConfigScopeWindow,
		shared.ConfigScopeRole,
		shared.ConfigScopeSession,
		shared.ConfigScopeProject,
		shared.ConfigScopeModulePair,
		shared.ConfigScopeModule,
		shared.ConfigScopeEnvironment,
		shared.ConfigScopeGlobal,
		shared.ConfigScopeApplication,
	}
	if len(poziomyOdNajwezszego) != len(oczekiwane) {
		t.Fatalf("poziomów jest %d, oczekiwane %d", len(poziomyOdNajwezszego), len(oczekiwane))
	}
	for indeks, poziom := range oczekiwane {
		if poziomyOdNajwezszego[indeks] != poziom {
			t.Errorf("na miejscu %d stoi %q, oczekiwane %q", indeks, poziomyOdNajwezszego[indeks], poziom)
		}
	}

	// Aplikacja przegrywa ze wszystkim: opisuje sam program, nie treść w nim
	// prowadzoną.
	if poziomyOdNajwezszego[len(poziomyOdNajwezszego)-1] != shared.ConfigScopeApplication {
		t.Error("poziom aplikacji nie stoi na końcu — zapis o programie bije zapis o treści")
	}
	if poziomyOdNajwezszego[0] != shared.ConfigScopeWindow {
		t.Error("okno komunikacji nie jest poziomem najwęższym")
	}
}

func TestOsieIdaOdKontaDoPlatformy(t *testing.T) {
	oczekiwane := []Os{OsKonta, OsModelu, OsPlatformy}
	if len(osieOdNajwezszej) != len(oczekiwane) {
		t.Fatalf("osi jest %d, oczekiwane %d", len(osieOdNajwezszej), len(oczekiwane))
	}
	for indeks, os := range oczekiwane {
		if osieOdNajwezszej[indeks] != os {
			t.Errorf("na miejscu %d stoi %q, oczekiwane %q", indeks, osieOdNajwezszej[indeks], os)
		}
	}
}

func TestOsPustaZnaczyPlatforme(t *testing.T) {
	if OsLubPlatforma("") != OsPlatformy {
		t.Error("oś pusta nie została odczytana jako platforma")
	}
	if !ZnanaOs("") {
		t.Error("oś pusta uznana za nieznaną")
	}
	for _, os := range osieOdNajwezszej {
		if !ZnanaOs(os) {
			t.Errorf("oś %q kontraktu uznana za nieznaną", os)
		}
	}
	if ZnanaOs("os-wymyslona") {
		t.Error("oś spoza kontraktu uznana za znaną")
	}
}

func TestZnanyPoziomObejmujeKomplet(t *testing.T) {
	for _, poziom := range poziomyOdNajwezszego {
		if !Znany(poziom) {
			t.Errorf("poziom %q kontraktu uznany za nieznany", poziom)
		}
	}
	if Znany(PoziomBrak) {
		t.Error("brak poziomu uznany za poziom znany")
	}
	if Znany("poziom-wymyslony") {
		t.Error("poziom spoza kontraktu uznany za znany")
	}
}

func TestPoziomBezBytuJestPomijany(t *testing.T) {
	pusty := Kontekst{}

	for _, poziom := range []Poziom{shared.ConfigScopeGlobal, shared.ConfigScopeApplication} {
		klucz, obowiazuje := pusty.Adres(poziom)
		if !obowiazuje {
			t.Errorf("poziom %q nie obowiązuje w pustym kontekście", poziom)
		}
		if klucz != "" {
			t.Errorf("poziom %q ma byt %q — platformy ani programu nie ma czym zawęzić", poziom, klucz)
		}
	}

	for _, poziom := range []Poziom{
		shared.ConfigScopeWindow, shared.ConfigScopeRole, shared.ConfigScopeSession,
		shared.ConfigScopeProject, shared.ConfigScopeModulePair, shared.ConfigScopeModule,
		shared.ConfigScopeEnvironment,
	} {
		if _, obowiazuje := pusty.Adres(poziom); obowiazuje {
			t.Errorf("poziom %q obowiązuje mimo pustego bytu", poziom)
		}
	}

	if _, obowiazuje := pusty.Adres("poziom-wymyslony"); obowiazuje {
		t.Error("poziom spoza kontraktu obowiązuje")
	}
}

func TestAdresyPustegoKontekstuToDwaPoziomyPlatformy(t *testing.T) {
	adresy := Kontekst{}.Adresy()
	oczekiwane := []Adres{
		{Poziom: shared.ConfigScopeGlobal, Os: OsPlatformy},
		{Poziom: shared.ConfigScopeApplication, Os: OsPlatformy},
	}
	if len(adresy) != len(oczekiwane) {
		t.Fatalf("pusty kontekst dał %d adresów, oczekiwane %d: %+v", len(adresy), len(oczekiwane), adresy)
	}
	for indeks, adres := range oczekiwane {
		if adresy[indeks] != adres {
			t.Errorf("na miejscu %d stoi %+v, oczekiwane %+v", indeks, adresy[indeks], adres)
		}
	}
}

func TestAdresyIdaPoziomamiAWRamachPoziomuOsiami(t *testing.T) {
	adresy := kontekstPelny().Adresy()

	if oczekiwane := len(poziomyOdNajwezszego) * len(osieOdNajwezszej); len(adresy) != oczekiwane {
		t.Fatalf("kontekst pełny dał %d adresów, oczekiwane %d", len(adresy), oczekiwane)
	}

	// Poziom zmienia się co trzy pozycje, oś co jedną: sprawdzian idzie po obu wykazach naraz.
	pozycja := 0
	for _, poziom := range poziomyOdNajwezszego {
		for _, os := range osieOdNajwezszej {
			adres := adresy[pozycja]
			if adres.Poziom != poziom {
				t.Errorf("na miejscu %d stoi poziom %q, oczekiwany %q", pozycja, adres.Poziom, poziom)
			}
			if adres.Os != os {
				t.Errorf("na miejscu %d stoi oś %q, oczekiwana %q", pozycja, adres.Os, os)
			}
			pozycja++
		}
	}
}

func TestOsBezBytuJestPomijana(t *testing.T) {
	kontekst := kontekstPelny()
	kontekst.Model = ""
	kontekst.Konto = ""

	adresy := kontekst.Adresy()
	if len(adresy) != len(poziomyOdNajwezszego) {
		t.Fatalf("kontekst bez osi dał %d adresów, oczekiwane %d", len(adresy), len(poziomyOdNajwezszego))
	}
	for _, adres := range adresy {
		if adres.Os != OsPlatformy {
			t.Errorf("adres %+v niesie oś inną niż platforma", adres)
		}
		if adres.KluczOsi != "" {
			t.Errorf("adres %+v niesie byt osi mimo braku wskazania", adres)
		}
	}
}

func TestNajwezszyZapisWygrywa(t *testing.T) {
	const klucz = kluczTrybUprawnien

	zrodlo := NoweZrodloPamieciowe()
	kontekst := kontekstPelny()
	kontekst.Model = ""
	kontekst.Konto = ""

	for _, poziom := range poziomyOdNajwezszego {
		kluczZasiegu, _ := kontekst.Adres(poziom)
		zrodlo.Ustaw(poziom, kluczZasiegu, klucz, string(poziom), RodzajTekst)
	}

	rozstrzygacz := Nowy(zrodlo, RejestrWbudowany())

	for _, poziom := range poziomyOdNajwezszego {
		wynik := rozstrzygacz.Rozstrzygnij(kontekst, klucz)
		if wynik.Poziom != poziom {
			t.Fatalf("wygrał poziom %q, oczekiwany %q", wynik.Poziom, poziom)
		}
		if wynik.Wartosc != string(poziom) {
			t.Errorf("poziom %q oddał wartość %q", poziom, wynik.Wartosc)
		}
		if wynik.Pochodzenie != PochodzenieZapis {
			t.Errorf("poziom %q oddał pochodzenie %q", poziom, wynik.Pochodzenie)
		}

		kluczZasiegu, _ := kontekst.Adres(poziom)
		zrodlo.Usun(poziom, kluczZasiegu, klucz)
	}

	wynik := rozstrzygacz.Rozstrzygnij(kontekst, klucz)
	if wynik.Poziom != PoziomBrak {
		t.Errorf("po zdjęciu wszystkich zapisów wygrał poziom %q", wynik.Poziom)
	}
	if wynik.Pochodzenie != pochodzenieDomyslna {
		t.Errorf("po zdjęciu wszystkich zapisów pochodzenie to %q", wynik.Pochodzenie)
	}
}

func TestPoziomBijeOsNiezaleznieOdSzerokosci(t *testing.T) {
	const klucz = kluczTrybUprawnien

	zrodlo := NoweZrodloPamieciowe()
	kontekst := kontekstPelny()

	zrodlo.ustawWOsi(Adres{
		Poziom: shared.ConfigScopeGlobal,
		Os:     OsKonta, KluczOsi: kontekst.Konto,
	}, klucz, "zapis konta na poziomie globalnym", RodzajTekst)

	zrodlo.ustawWOsi(Adres{
		Poziom: shared.ConfigScopeWindow, KluczZasiegu: kontekst.Okno,
		Os: OsPlatformy,
	}, klucz, "zapis platformy w oknie", RodzajTekst)

	wynik := Nowy(zrodlo, RejestrWbudowany()).Rozstrzygnij(kontekst, klucz)

	if wynik.Poziom != shared.ConfigScopeWindow {
		t.Fatalf("wygrał poziom %q — oś przebiła poziom", wynik.Poziom)
	}
	if wynik.Os != OsPlatformy {
		t.Errorf("wygrała oś %q", wynik.Os)
	}
	if wynik.Wartosc != "zapis platformy w oknie" {
		t.Errorf("wygrała wartość %q", wynik.Wartosc)
	}
}

func TestOsRozstrzygaWRamachJednegoPoziomu(t *testing.T) {
	const klucz = kluczTrybUprawnien
	kontekst := kontekstPelny()

	przypadki := []struct {
		nazwa    string
		zapisane []Os
		wygrywa  Os
	}{
		{"komplet osi", []Os{OsPlatformy, OsModelu, OsKonta}, OsKonta},
		{"model i platforma", []Os{OsPlatformy, OsModelu}, OsModelu},
		{"sama platforma", []Os{OsPlatformy}, OsPlatformy},
		{"konto i platforma", []Os{OsPlatformy, OsKonta}, OsKonta},
	}

	for _, przypadek := range przypadki {
		t.Run(przypadek.nazwa, func(t *testing.T) {
			zrodlo := NoweZrodloPamieciowe()
			for _, os := range przypadek.zapisane {
				kluczOsi, _ := kontekst.adresOsi(os)
				zrodlo.ustawWOsi(Adres{
					Poziom: shared.ConfigScopeSession, KluczZasiegu: kontekst.KartaSesji,
					Os: os, KluczOsi: kluczOsi,
				}, klucz, string(os), RodzajTekst)
			}

			wynik := Nowy(zrodlo, RejestrWbudowany()).Rozstrzygnij(kontekst, klucz)
			if wynik.Os != przypadek.wygrywa {
				t.Errorf("wygrała oś %q, oczekiwana %q", wynik.Os, przypadek.wygrywa)
			}
			if wynik.Poziom != shared.ConfigScopeSession {
				t.Errorf("wygrał poziom %q", wynik.Poziom)
			}
		})
	}
}

// TestZapisPozaKontekstemNieWchodzi pilnuje szczelności w drugą stronę:
// ustawienie zapisane dla innej sesji nie może wejść do rozstrzygnięcia tej.
// To jest dokładnie ten przeciek, który widać dopiero po zachowaniu modelu.
func TestZapisPozaKontekstemNieWchodzi(t *testing.T) {
	const klucz = kluczTrybUprawnien

	zrodlo := NoweZrodloPamieciowe()
	zrodlo.Ustaw(shared.ConfigScopeSession, "sesja-cudza", klucz, "wartość cudzej sesji", RodzajTekst)
	zrodlo.Ustaw(shared.ConfigScopeWindow, "okno-cudze", klucz, "wartość cudzego okna", RodzajTekst)

	kontekst := Kontekst{KartaSesji: "sesja-pierwsza", Okno: "okno-pierwsze"}
	wynik := Nowy(zrodlo, RejestrWbudowany()).Rozstrzygnij(kontekst, klucz)

	if wynik.Pochodzenie == PochodzenieZapis {
		t.Errorf("do rozstrzygnięcia weszła wartość spoza kontekstu: %+v", wynik)
	}
	if wynik.Pochodzenie != pochodzenieDomyslna {
		t.Errorf("pochodzenie to %q, oczekiwana wartość domyślna", wynik.Pochodzenie)
	}
}

func TestKluczSpozaRejestruNieJestBledem(t *testing.T) {
	wynik := Nowy(nil, nil).Rozstrzygnij(Kontekst{}, "klucz.spoza.rejestru")

	if wynik.Pochodzenie != PochodzenieNieznane {
		t.Errorf("klucz spoza rejestru dał pochodzenie %q", wynik.Pochodzenie)
	}
	if wynik.Wartosc != "" {
		t.Errorf("klucz spoza rejestru dał wartość %q", wynik.Wartosc)
	}
	if wynik.Rodzaj != RodzajTekst {
		t.Errorf("klucz spoza rejestru dał rodzaj %q", wynik.Rodzaj)
	}
}

func TestKluczSpozaRejestruZZapisemWraca(t *testing.T) {
	const klucz = "klucz.spoza.rejestru"

	zrodlo := NoweZrodloPamieciowe()
	zrodlo.Ustaw(shared.ConfigScopeGlobal, "", klucz, "wartość zapisana", RodzajTekst)

	wynik := Nowy(zrodlo, RejestrWbudowany()).Rozstrzygnij(Kontekst{}, klucz)
	if wynik.Pochodzenie != PochodzenieZapis {
		t.Errorf("klucz spoza rejestru z zapisem dał pochodzenie %q", wynik.Pochodzenie)
	}
	if wynik.Wartosc != "wartość zapisana" {
		t.Errorf("wartość odczytana jako %q", wynik.Wartosc)
	}
}

type zrodloWadliwe struct {
	wpisy []Wpis
	blad  error
}

func (z zrodloWadliwe) Wpisy(int64, []Adres) ([]Wpis, error) { return z.wpisy, z.blad }

func TestBladZrodlaNieZatrzymujeRozstrzygania(t *testing.T) {
	const klucz = kluczTrybUprawnien

	zrodlo := zrodloWadliwe{
		wpisy: []Wpis{{
			Poziom: shared.ConfigScopeGlobal, Os: OsPlatformy,
			Klucz: klucz, Wartosc: "wartość mimo błędu", Rodzaj: RodzajTekst,
		}},
		blad: errors.New("warstwa trwałości nie odpowiada"),
	}

	wynik := Nowy(zrodlo, RejestrWbudowany()).Rozstrzygnij(Kontekst{}, klucz)
	if wynik.Wartosc != "wartość mimo błędu" {
		t.Errorf("wpis odczytany mimo błędu nie wszedł do rozstrzygnięcia: %+v", wynik)
	}

	polityka := Nowy(zrodlo, RejestrWbudowany()).PolitykaEfektywna(Kontekst{})
	if polityka.BladZrodla == nil {
		t.Error("polityka nie niesie błędu źródła — diagnostyka go nie zobaczy")
	}
	if len(polityka.Pozycje) == 0 {
		t.Error("polityka jest pusta mimo dostarczonych wpisów")
	}
}

func TestWpisBezKluczaJestPomijany(t *testing.T) {
	zrodlo := zrodloWadliwe{wpisy: []Wpis{
		{Poziom: shared.ConfigScopeGlobal, Os: OsPlatformy, Klucz: "", Wartosc: "bez klucza"},
	}}
	rozstrzygacz := Nowy(zrodlo, RejestrWbudowany())

	if wynik := rozstrzygacz.Rozstrzygnij(Kontekst{}, ""); wynik.Pochodzenie == PochodzenieZapis {
		t.Errorf("wpis bez klucza wszedł do rozstrzygnięcia: %+v", wynik)
	}
}

func TestOdbiornikZerowyRozstrzygaDomyslnie(t *testing.T) {
	var zerowy *Rozstrzygacz
	wynik := zerowy.Rozstrzygnij(Kontekst{}, kluczTrybUprawnien)
	if wynik.Pochodzenie != PochodzenieNieznane {
		t.Errorf("rozstrzygacz zerowy dał pochodzenie %q", wynik.Pochodzenie)
	}

	bezZaleznosci := Nowy(nil, nil)
	if bezZaleznosci.Rejestr() == nil {
		t.Error("rozstrzygacz bez rejestru nie sięgnął po rejestr wbudowany")
	}
	zDomyslna := bezZaleznosci.Rozstrzygnij(Kontekst{}, kluczTrybUprawnien)
	if zDomyslna.Pochodzenie != pochodzenieDomyslna {
		t.Errorf("rozstrzygacz bez źródła dał pochodzenie %q zamiast wartości domyślnej",
			zDomyslna.Pochodzenie)
	}
}
