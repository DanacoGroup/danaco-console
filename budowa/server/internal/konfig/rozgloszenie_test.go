package konfig

import (
	"encoding/json"
	"sync"
	"testing"

	"danacoconsole/shared"
)

// Droga na żywo od zapisu nastawy do tego, kto z niej korzysta, bez zbędnych doręczeń.

// odbiorSprawdzianu zbiera doręczenia. Rozgłośnia doręcza w wątku ogłaszającego,
// ale nasłuch może być zapisany z innego, więc zbieranie idzie pod zamkiem.
type odbiorSprawdzianu struct {
	zamek      sync.Mutex
	doreczenia []Zmiana
}

func (o *odbiorSprawdzianu) przyjmij(z Zmiana) {
	o.zamek.Lock()
	o.doreczenia = append(o.doreczenia, z)
	o.zamek.Unlock()
}

func (o *odbiorSprawdzianu) wykaz() []Zmiana {
	o.zamek.Lock()
	defer o.zamek.Unlock()
	return append([]Zmiana(nil), o.doreczenia...)
}

func (o *odbiorSprawdzianu) liczba() int {
	o.zamek.Lock()
	defer o.zamek.Unlock()
	return len(o.doreczenia)
}

// TestPierwszeDoreczenieIdzieOdRazu sprawdza część umowy, nie uprzejmość:
// nasłuchujący dostaje punkt wyjścia przy zapisaniu się, więc nie musi po niego
// sięgać drugą drogą — a druga droga znaczyłaby wyścig z zapisem, który zdążył
// się wcisnąć pomiędzy.
func TestPierwszeDoreczenieIdzieOdRazu(t *testing.T) {
	rozstrzygacz := Nowy(NoweZrodloPamieciowe(), RejestrWbudowany())
	odbior := &odbiorSprawdzianu{}

	odwolaj := rozstrzygacz.Sledz(Kontekst{}, []string{kluczTrybUprawnien}, odbior.przyjmij)
	defer odwolaj()

	doreczenia := odbior.wykaz()
	if len(doreczenia) != 1 {
		t.Fatalf("doręczeń po zapisaniu się: %d, oczekiwane jedno", len(doreczenia))
	}
	if doreczenia[0].Klucz != kluczTrybUprawnien {
		t.Errorf("doręczono klucz %q", doreczenia[0].Klucz)
	}
	if doreczenia[0].Wynik.Pochodzenie != pochodzenieDomyslna {
		t.Errorf("pierwsze doręczenie niesie pochodzenie %q", doreczenia[0].Wynik.Pochodzenie)
	}
}

// TestOgloszenieDorczaZmianeRzeczywista sprawdza drogę zwykłą: zapis, ogłoszenie,
// doręczenie nowej wartości.
func TestOgloszenieDorczaZmianeRzeczywista(t *testing.T) {
	zrodlo := NoweZrodloPamieciowe()
	rozstrzygacz := Nowy(zrodlo, RejestrWbudowany())
	odbior := &odbiorSprawdzianu{}

	odwolaj := rozstrzygacz.Sledz(Kontekst{KartaSesji: "sesja-pierwsza"},
		[]string{kluczTrybUprawnien}, odbior.przyjmij)
	defer odwolaj()

	zrodlo.Ustaw(shared.ConfigScopeSession, "sesja-pierwsza", kluczTrybUprawnien,
		shared.PermissionModePlan, RodzajTekst)
	rozstrzygacz.Oglos(kluczTrybUprawnien)

	doreczenia := odbior.wykaz()
	if len(doreczenia) != 2 {
		t.Fatalf("doręczeń: %d, oczekiwane dwa (stan wyjściowy i zmiana)", len(doreczenia))
	}
	ostatnie := doreczenia[len(doreczenia)-1]
	if ostatnie.Wynik.Wartosc != shared.PermissionModePlan {
		t.Errorf("doręczono wartość %q", ostatnie.Wynik.Wartosc)
	}
	if ostatnie.Wynik.Poziom != shared.ConfigScopeSession {
		t.Errorf("doręczono poziom %q", ostatnie.Wynik.Poziom)
	}
}

// TestOgloszenieBezZmianyNieBudziNasluchu jest sprawdzianem drugiej połowy
// umowy. Rozgłośnia liczy skutek osobno dla kontekstu każdego nasłuchu, więc
// zapis, który dla tego nasłuchu niczego nie zmienia, nie może go obudzić.
func TestOgloszenieBezZmianyNieBudziNasluchu(t *testing.T) {
	zrodlo := NoweZrodloPamieciowe()
	rozstrzygacz := Nowy(zrodlo, RejestrWbudowany())
	odbior := &odbiorSprawdzianu{}

	odwolaj := rozstrzygacz.Sledz(Kontekst{}, []string{kluczTrybUprawnien}, odbior.przyjmij)
	defer odwolaj()

	poPierwszym := odbior.liczba()

	// Ogłoszenie bez żadnego zapisu.
	rozstrzygacz.Oglos(kluczTrybUprawnien)
	if odbior.liczba() != poPierwszym {
		t.Error("ogłoszenie bez zapisu obudziło nasłuch")
	}

	// Zapis tej samej wartości, którą nasłuch już ma.
	zrodlo.Ustaw(shared.ConfigScopeGlobal, "", kluczTrybUprawnien,
		shared.PermissionModeManual, RodzajTekst)
	rozstrzygacz.Oglos(kluczTrybUprawnien)
	if odbior.liczba() == poPierwszym {
		t.Skip("wartość domyślna i zapisana różnią się poziomem, więc doręczenie jest zasadne")
	}
}

// TestZapisSzerszyNieBudziNasluchuZWezszym sprawdza regułę wprost z opisu:
// zapis na poziomie aplikacji nie budzi nikogo, kto ma wartość z poziomu
// węższego. Doręczenie w tej sytuacji podstawiłoby nasłuchującemu wartość, która
// go nie obowiązuje.
func TestZapisSzerszyNieBudziNasluchuZWezszym(t *testing.T) {
	zrodlo := NoweZrodloPamieciowe()
	rozstrzygacz := Nowy(zrodlo, RejestrWbudowany())

	// Nasłuch ma wartość z poziomu okna.
	zrodlo.Ustaw(shared.ConfigScopeWindow, "okno-pierwsze", kluczTrybUprawnien,
		shared.PermissionModePlan, RodzajTekst)

	odbior := &odbiorSprawdzianu{}
	odwolaj := rozstrzygacz.Sledz(Kontekst{Okno: "okno-pierwsze"},
		[]string{kluczTrybUprawnien}, odbior.przyjmij)
	defer odwolaj()

	poPierwszym := odbior.liczba()

	// Zapis na poziomie globalnym — szerszym, więc przegrywa z oknem.
	zrodlo.Ustaw(shared.ConfigScopeGlobal, "", kluczTrybUprawnien,
		shared.PermissionModeAuto, RodzajTekst)
	rozstrzygacz.Oglos(kluczTrybUprawnien)

	if odbior.liczba() != poPierwszym {
		t.Errorf("zapis na poziomie szerszym obudził nasłuch mający wartość z okna: %+v",
			odbior.wykaz())
	}
}

// TestOgloszenieDotyczyWylacznieKluczyNasluchu pilnuje zawężenia: nasłuch pytał
// o jeden klucz i nie ma dostawać doręczeń o innych.
func TestOgloszenieDotyczyWylacznieKluczyNasluchu(t *testing.T) {
	zrodlo := NoweZrodloPamieciowe()
	rozstrzygacz := Nowy(zrodlo, RejestrWbudowany())
	odbior := &odbiorSprawdzianu{}

	odwolaj := rozstrzygacz.Sledz(Kontekst{}, []string{kluczTrybUprawnien}, odbior.przyjmij)
	defer odwolaj()

	poPierwszym := odbior.liczba()

	zrodlo.Ustaw(shared.ConfigScopeGlobal, "", KluczPulapKosztu, "12", RodzajLiczba)
	rozstrzygacz.Oglos(KluczPulapKosztu)

	if odbior.liczba() != poPierwszym {
		t.Errorf("doręczono zmianę klucza spoza nasłuchu: %+v", odbior.wykaz())
	}
}

// TestOdwolanieNasluchuJestSkuteczneIPowtarzalne sprawdza wykreślenie i to, że
// wołanie go wielokrotnie jest bezpieczne.
func TestOdwolanieNasluchuJestSkuteczneIPowtarzalne(t *testing.T) {
	zrodlo := NoweZrodloPamieciowe()
	rozstrzygacz := Nowy(zrodlo, RejestrWbudowany())
	odbior := &odbiorSprawdzianu{}

	odwolaj := rozstrzygacz.Sledz(Kontekst{}, []string{kluczTrybUprawnien}, odbior.przyjmij)
	poPierwszym := odbior.liczba()

	odwolaj()
	odwolaj()

	zrodlo.Ustaw(shared.ConfigScopeGlobal, "", kluczTrybUprawnien,
		shared.PermissionModeAuto, RodzajTekst)
	rozstrzygacz.Oglos(kluczTrybUprawnien)

	if odbior.liczba() != poPierwszym {
		t.Error("wykreślony nasłuch nadal dostaje doręczenia")
	}
}

// TestSledzenieBezOdbiorcyAlboBezKluczyDajeFunkcjePusta pilnuje reguły
// fail-open: brak wskazania nie jest odmową.
func TestSledzenieBezOdbiorcyAlboBezKluczyDajeFunkcjePusta(t *testing.T) {
	rozstrzygacz := Nowy(NoweZrodloPamieciowe(), RejestrWbudowany())

	if odwolaj := rozstrzygacz.Sledz(Kontekst{}, []string{kluczTrybUprawnien}, nil); odwolaj == nil {
		t.Fatal("śledzenie bez odbiorcy nie oddało funkcji wykreślającej")
	} else {
		odwolaj()
	}

	odbior := &odbiorSprawdzianu{}
	for _, klucze := range [][]string{nil, {}, {""}} {
		odwolaj := rozstrzygacz.Sledz(Kontekst{}, klucze, odbior.przyjmij)
		odwolaj()
	}
	if odbior.liczba() != 0 {
		t.Errorf("śledzenie bez kluczy doręczyło %d zmian", odbior.liczba())
	}

	var zerowy *Rozstrzygacz
	zerowy.Sledz(Kontekst{}, []string{kluczTrybUprawnien}, odbior.przyjmij)()
	zerowy.Oglos(kluczTrybUprawnien)
}

// TestNastawaNaZywoNieStarzejeSie sprawdza uchwyt z drogi gorącej. Straż bramki
// rozstrzyga przy każdym pakiecie z gniazda, więc uchwyt zdejmuje koszt zejścia
// do trwałości — i nie wolno mu przy tym zdjąć prawdy.
func TestNastawaNaZywoNieStarzejeSie(t *testing.T) {
	zrodlo := NoweZrodloPamieciowe()
	rozstrzygacz := Nowy(zrodlo, RejestrWbudowany())

	nastawa := rozstrzygacz.NastawaNaZywo(Kontekst{}, kluczTrybUprawnien)
	defer nastawa.Zamknij()

	if nastawa.Wartosc() != shared.PermissionModeManual {
		t.Errorf("uchwyt wystartował z wartością %q", nastawa.Wartosc())
	}

	zrodlo.Ustaw(shared.ConfigScopeGlobal, "", kluczTrybUprawnien,
		shared.PermissionModeAuto, RodzajTekst)
	rozstrzygacz.Oglos(kluczTrybUprawnien)

	if nastawa.Wartosc() != shared.PermissionModeAuto {
		t.Errorf("uchwyt nie przyjął zmiany: %q", nastawa.Wartosc())
	}
	if nastawa.Wynik().Poziom != shared.ConfigScopeGlobal {
		t.Errorf("uchwyt niesie poziom %q", nastawa.Wynik().Poziom)
	}

	nastawa.Zamknij()
	zrodlo.Ustaw(shared.ConfigScopeGlobal, "", kluczTrybUprawnien,
		shared.PermissionModePlan, RodzajTekst)
	rozstrzygacz.Oglos(kluczTrybUprawnien)
	if nastawa.Wartosc() != shared.PermissionModeAuto {
		t.Errorf("uchwyt zamknięty nadal przyjmuje zmiany: %q", nastawa.Wartosc())
	}
}

// TestUchwytZerowyNieWywracaWywolania pilnuje odbiornika zerowego — uchwyt bywa
// polem struktury złożonej przed nasłuchem.
func TestUchwytZerowyNieWywracaWywolania(t *testing.T) {
	var zerowa *Nastawa
	if zerowa.Wartosc() != "" {
		t.Error("uchwyt zerowy oddał wartość")
	}
	zerowa.Zamknij()
}

// TestDwaNasluchyDostajaWlasneRozstrzygniecia sprawdza, że rozgłośnia liczy
// skutek osobno dla każdego kontekstu, a nie raz dla wszystkich.
func TestDwaNasluchyDostajaWlasneRozstrzygniecia(t *testing.T) {
	zrodlo := NoweZrodloPamieciowe()
	rozstrzygacz := Nowy(zrodlo, RejestrWbudowany())

	zPierwszej := &odbiorSprawdzianu{}
	zDrugiej := &odbiorSprawdzianu{}

	odwolajPierwszy := rozstrzygacz.Sledz(Kontekst{KartaSesji: "sesja-pierwsza"},
		[]string{kluczTrybUprawnien}, zPierwszej.przyjmij)
	defer odwolajPierwszy()
	odwolajDrugi := rozstrzygacz.Sledz(Kontekst{KartaSesji: "sesja-druga"},
		[]string{kluczTrybUprawnien}, zDrugiej.przyjmij)
	defer odwolajDrugi()

	drugiPrzed := zDrugiej.liczba()

	zrodlo.Ustaw(shared.ConfigScopeSession, "sesja-pierwsza", kluczTrybUprawnien,
		shared.PermissionModePlan, RodzajTekst)
	rozstrzygacz.Oglos(kluczTrybUprawnien)

	ostatnie := zPierwszej.wykaz()
	if ostatnie[len(ostatnie)-1].Wynik.Wartosc != shared.PermissionModePlan {
		t.Errorf("nasłuch własnej sesji nie dostał zmiany: %+v", ostatnie)
	}
	if zDrugiej.liczba() != drugiPrzed {
		t.Errorf("zapis sesji pierwszej obudził nasłuch sesji drugiej: %+v", zDrugiej.wykaz())
	}
}

// TestWartoscWchodziIWychodziZKopertyBezZmiany sprawdza obieg zamknięty
// kodowania wartości. Rodzaj rozstrzyga o postaci w kopercie kontraktu, więc
// pomyłka tutaj wysyła klientowi liczbę jako napis albo odwrotnie.
func TestWartoscWchodziIWychodziZKopertyBezZmiany(t *testing.T) {
	przypadki := []struct {
		nazwa     string
		wartosc   string
		rodzaj    Rodzaj
		wKopercie string
	}{
		{"tekst", "tryb ręczny", RodzajTekst, `"tryb ręczny"`},
		{"liczba", "12", RodzajLiczba, `12`},
		{"liczba ułamkowa", "0.5", RodzajLiczba, `0.5`},
		{"logiczna prawda", "true", RodzajLogiczna, `true`},
		{"logiczna fałsz", "false", RodzajLogiczna, `false`},
		{"json tablicowy", `["a","b"]`, RodzajJSON, `["a","b"]`},
	}

	for _, przypadek := range przypadki {
		t.Run(przypadek.nazwa, func(t *testing.T) {
			zakodowana := KodujJSON(przypadek.wartosc, przypadek.rodzaj)
			if string(zakodowana) != przypadek.wKopercie {
				t.Errorf("w kopercie stoi %s, oczekiwane %s", zakodowana, przypadek.wKopercie)
			}
			wartosc, rodzaj := DekodujJSON(zakodowana)
			if wartosc != przypadek.wartosc {
				t.Errorf("z koperty wróciła wartość %q", wartosc)
			}
			if rodzaj != przypadek.rodzaj {
				t.Errorf("z koperty wrócił rodzaj %q", rodzaj)
			}
		})
	}
}

// TestWartoscUszkodzonaIdzieJakoNapis pilnuje obietnicy z opisu: kodowanie nigdy
// nie zawodzi. Wartość, która miała być liczbą, a liczbą nie jest, jedzie jako
// napis — nie jako błąd i nie jako niepoprawny JSON w kopercie.
func TestWartoscUszkodzonaIdzieJakoNapis(t *testing.T) {
	zakodowana := KodujJSON("nie liczba", RodzajLiczba)
	if !json.Valid(zakodowana) {
		t.Fatalf("koperta niesie treść spoza JSON: %s", zakodowana)
	}
	if string(zakodowana) != `"nie liczba"` {
		t.Errorf("wartość uszkodzona zakodowana jako %s", zakodowana)
	}
}

// TestRodzajNieznanyStajeSieTekstem sprawdza drogę wartości nierozpoznanego rodzaju przez kodowanie kontraktu.
func TestRodzajNieznanyStajeSieTekstem(t *testing.T) {
	if RodzajLubTekst("rodzaj-wymyslony") != RodzajTekst {
		t.Error("rodzaj spoza schematu nie zszedł na tekst")
	}
	for _, rodzaj := range []Rodzaj{RodzajTekst, RodzajLiczba, RodzajLogiczna, RodzajJSON} {
		if !rodzaj.Znany() {
			t.Errorf("rodzaj %q ze schematu uznany za nieznany", rodzaj)
		}
	}
	if string(KodujJSON("cokolwiek", "rodzaj-wymyslony")) != `"cokolwiek"` {
		t.Error("rodzaj spoza schematu nie został potraktowany jak tekst")
	}
}

// TestPustaTrescKopertyDajeTekstPusty sprawdza dolny kraniec drogi: pustą treść koperty bez awarii dekodowania.
func TestPustaTrescKopertyDajeTekstPusty(t *testing.T) {
	wartosc, rodzaj := DekodujJSON(nil)
	if wartosc != "" || rodzaj != RodzajTekst {
		t.Errorf("pusta treść dała %q rodzaju %q", wartosc, rodzaj)
	}
}

// TestWpisKontraktuNieNiesieOsiPlatformy pilnuje reguły koperty: brak osi znaczy
// platformę, więc oś platformy nie ma prawa jechać wprost. Wysłanie jej
// kazałoby klientowi odróżniać dwa zapisy tego samego.
func TestWpisKontraktuNieNiesieOsiPlatformy(t *testing.T) {
	wpis := Wynik{
		Klucz: kluczTrybUprawnien, Wartosc: shared.PermissionModePlan, Rodzaj: RodzajTekst,
		Poziom: shared.ConfigScopeSession, KluczZasiegu: "sesja-pierwsza",
		Os: OsPlatformy, Pochodzenie: PochodzenieZapis,
	}.WpisKontraktu()

	if wpis.Axis != nil {
		t.Errorf("wpis niesie oś platformy wprost: %v", *wpis.Axis)
	}
	if wpis.ScopeId == nil || *wpis.ScopeId != "sesja-pierwsza" {
		t.Errorf("wpis zgubił byt poziomu: %v", wpis.ScopeId)
	}
	if wpis.Scope != shared.ConfigScopeSession {
		t.Errorf("wpis niesie poziom %q", wpis.Scope)
	}
}

// TestWpisKontraktuNiesieOsWezsza sprawdza drugą stronę: oś konta i oś modelu
// jadą w kopercie razem ze swoim bytem.
func TestWpisKontraktuNiesieOsWezsza(t *testing.T) {
	wpis := Wynik{
		Klucz: kluczTrybUprawnien, Rodzaj: RodzajTekst,
		Poziom: shared.ConfigScopeGlobal,
		Os:     OsKonta, KluczOsi: "konto-pierwsze",
	}.WpisKontraktu()

	if wpis.Axis == nil || *wpis.Axis != OsKonta {
		t.Fatalf("wpis nie niesie osi konta: %v", wpis.Axis)
	}
	if wpis.AxisId == nil || *wpis.AxisId != "konto-pierwsze" {
		t.Errorf("wpis zgubił byt osi: %v", wpis.AxisId)
	}
}

// TestWartoscDomyslnaJedzieBezPoziomu sprawdza, że pusty poziom w kopercie jest
// informacją, nie brakiem: obowiązuje warstwa definicji.
func TestWartoscDomyslnaJedzieBezPoziomu(t *testing.T) {
	wynik := Nowy(nil, RejestrWbudowany()).Rozstrzygnij(Kontekst{}, kluczTrybUprawnien)
	wpis := wynik.WpisKontraktu()

	if wpis.Scope != PoziomBrak {
		t.Errorf("wartość domyślna niesie poziom %q", wpis.Scope)
	}
	if wpis.ScopeId != nil {
		t.Errorf("wartość domyślna niesie byt poziomu: %v", *wpis.ScopeId)
	}
}

// TestZapisZKontraktuOdrzucaAdresSpozaKontraktu pilnuje jedynej bramy pakietu na
// drodze zapisu: poziom albo oś spoza kontraktu nie mogą wejść do tabeli ustawień,
// a odmowa nazwana jest lepsza od błędu bazy.
func TestZapisZKontraktuOdrzucaAdresSpozaKontraktu(t *testing.T) {
	przypadki := []struct {
		nazwa   string
		klucz   string
		adres   Adres
		wchodzi bool
	}{
		{"poziom i oś kontraktu", kluczTrybUprawnien,
			Adres{Poziom: shared.ConfigScopeSession, KluczZasiegu: "sesja-pierwsza", Os: OsKonta}, true},
		{"oś pusta znaczy platformę", kluczTrybUprawnien,
			Adres{Poziom: shared.ConfigScopeGlobal}, true},
		{"klucz pusty", "",
			Adres{Poziom: shared.ConfigScopeGlobal}, false},
		{"poziom spoza kontraktu", kluczTrybUprawnien,
			Adres{Poziom: "poziom-wymyslony"}, false},
		{"oś spoza kontraktu", kluczTrybUprawnien,
			Adres{Poziom: shared.ConfigScopeGlobal, Os: "os-wymyslona"}, false},
	}

	for _, przypadek := range przypadki {
		t.Run(przypadek.nazwa, func(t *testing.T) {
			wpis, przyjety := ZapisZKontraktuWOsi(przypadek.klucz,
				[]byte(`"tryb ręczny"`), przypadek.adres)
			if przyjety != przypadek.wchodzi {
				t.Fatalf("zapis przyjęty=%t, oczekiwane %t", przyjety, przypadek.wchodzi)
			}
			if !przyjety {
				return
			}
			if wpis.Os != OsLubPlatforma(przypadek.adres.Os) {
				t.Errorf("wpis niesie oś %q", wpis.Os)
			}
			if wpis.Wartosc != "tryb ręczny" || wpis.Rodzaj != RodzajTekst {
				t.Errorf("wpis niesie %q rodzaju %q", wpis.Wartosc, wpis.Rodzaj)
			}
		})
	}
}

// TestPolitykaEfektywnaNiesieKompletIWskazujeZrodlo sprawdza podgląd, na którym
// stoi okno konfiguracji: komplet ustawień rejestru plus klucze zapisane spoza
// rejestru, każde ze wskazaniem, skąd pochodzi.
func TestPolitykaEfektywnaNiesieKompletIWskazujeZrodlo(t *testing.T) {
	const kluczObcy = "klucz.spoza.rejestru"

	zrodlo := NoweZrodloPamieciowe()
	zrodlo.Ustaw(shared.ConfigScopeSession, "sesja-pierwsza", kluczTrybUprawnien,
		shared.PermissionModePlan, RodzajTekst)
	zrodlo.Ustaw(shared.ConfigScopeSession, "sesja-pierwsza", kluczObcy,
		"wartość spoza rejestru", RodzajTekst)

	rozstrzygacz := Nowy(zrodlo, RejestrWbudowany())
	polityka := rozstrzygacz.PolitykaEfektywna(Kontekst{KartaSesji: "sesja-pierwsza"})

	if len(polityka.Pozycje) <= rozstrzygacz.Rejestr().Liczba() {
		t.Errorf("polityka niesie %d pozycji przy %d definicjach — klucz spoza rejestru zniknął",
			len(polityka.Pozycje), rozstrzygacz.Rejestr().Liczba())
	}

	pozycja, jest := polityka.Pozycja(kluczTrybUprawnien)
	if !jest {
		t.Fatal("polityka nie niesie klucza z rejestru")
	}
	if pozycja.Wartosc != shared.PermissionModePlan {
		t.Errorf("polityka niesie wartość %q", pozycja.Wartosc)
	}

	poziom, pochodzenie := polityka.Zrodlo(kluczTrybUprawnien)
	if poziom != shared.ConfigScopeSession || pochodzenie != PochodzenieZapis {
		t.Errorf("polityka wskazuje źródło %q/%q", poziom, pochodzenie)
	}

	if _, jest := polityka.Pozycja(kluczObcy); !jest {
		t.Error("klucz zapisany spoza rejestru nie wszedł do polityki — odczyt stał się bramą")
	}

	poziomBrak, pochodzenieBrak := polityka.Zrodlo("klucz.ktorego.nigdzie.nie.ma")
	if poziomBrak != PoziomBrak || pochodzenieBrak != PochodzenieNieznane {
		t.Errorf("klucz nieobecny dał źródło %q/%q", poziomBrak, pochodzenieBrak)
	}

	if wpisy := polityka.WpisyKontraktu(); len(wpisy) != len(polityka.Pozycje) {
		t.Errorf("koperta niesie %d wpisów przy %d pozycjach", len(wpisy), len(polityka.Pozycje))
	}
}
