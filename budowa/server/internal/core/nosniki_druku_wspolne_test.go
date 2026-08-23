package core

import (
	"strings"
	"testing"

	"danacoconsole/shared"
)

// Sprawdziany wspólnego wykazu nośników druku.
//
// Szkoda, którą ten plik ma wykluczyć, jest jedna i nazwał ją Właściciel wprost:
// dwa wykazy nośników rozjadą się przy pierwszej poprawce i Operator dostanie
// A3 o wymiarach A4. Rozjechały się już raz — koperta DL miała w wykazie
// Designu 99 na 210 mm, a w wykazie okna Studia 220 na 110 mm — więc pilnowanie
// nie jest tu przewidywaniem, a zapisem tego, co się stało.

// TestWykazNosnikowJestJedenIPionowy mierzy dwie rzeczy naraz: że nazwa nie
// powtarza się w wykazie i że każdy wymiar jest zapisany pionowo. Nośnik
// zapisany poziomo niósłby obrót w definicji, a obrót jest rozstrzygnięciem
// materiału.
func TestWykazNosnikowJestJedenIPionowy(t *testing.T) {
	widziane := map[string]int{}
	for _, nosnik := range wykazNosnikowDruku {
		klucz := strings.ToUpper(strings.TrimSpace(nosnik.Nazwa))
		if klucz == "" {
			t.Fatal("wykaz nośników niesie pozycję bez nazwy")
		}
		widziane[klucz]++
		if nosnik.SzerokoscMm <= 0 || nosnik.WysokoscMm <= 0 {
			t.Errorf("nośnik %s ma wymiar niedodatni: %v na %v",
				nosnik.Nazwa, nosnik.SzerokoscMm, nosnik.WysokoscMm)
		}
		if nosnik.SzerokoscMm > nosnik.WysokoscMm {
			t.Errorf("nośnik %s jest zapisany poziomo (%v na %v) — wykaz trzyma "+
				"wymiary pionowe, a obrót nakłada nastawa strony",
				nosnik.Nazwa, nosnik.SzerokoscMm, nosnik.WysokoscMm)
		}
	}
	for nazwa, ile := range widziane {
		if ile > 1 {
			t.Errorf("nośnik %s stoi w wykazie wspólnym %d razy", nazwa, ile)
		}
	}
}

// TestWykazNosnikowNiesieKopertyWlasciciela pilnuje trzech kopert wymienionych
// przez Właściciela wraz z wymiarami normy ISO 269.
func TestWykazNosnikowNiesieKopertyWlasciciela(t *testing.T) {
	wymagane := map[string][2]float64{
		"C4": {229, 324},
		"C5": {162, 229},
		"C6": {114, 162},
	}
	for nazwa, wymiar := range wymagane {
		nosnik, jest := nosnikDrukuONazwie(nazwa)
		if !jest {
			t.Fatalf("wykaz wspólny nie niesie koperty %s", nazwa)
		}
		if nosnik.SzerokoscMm != wymiar[0] || nosnik.WysokoscMm != wymiar[1] {
			t.Errorf("koperta %s ma wymiary %v na %v mm, a Właściciel wskazał %v na %v",
				nazwa, nosnik.SzerokoscMm, nosnik.WysokoscMm, wymiar[0], wymiar[1])
		}
		if nosnik.Rodzina != rodzinaNosnikaKoperta {
			t.Errorf("koperta %s stoi w rodzinie %q, a nie w rodzinie kopert",
				nazwa, nosnik.Rodzina)
		}
	}
}

// TestWykazDesignuJestPrzekladem wykazuje, że obszar druku Designu
// PRZEKŁADA wykaz wspólny, a nie trzyma drugiej kopii liczb.
//
// Poprzednik tego sprawdzianu pilnował zgodności DWÓCH wykazów stojących obok
// siebie i zapalił się przy pierwszym przebiegu, wskazując dwa realne rozjazdy:
// kopertę DL (Design miał 99 na 210 mm — wymiar wkładki, nie koperty normy ISO
// 269, czyli 110 na 220) oraz wizytówkę zapisaną położoną wbrew zdaniu nad
// wykazem. Wykaz własny Designu został potem zdjęty i tę zgodność sprawdza się
// teraz inaczej: przekład ma oddać KAŻDĄ pozycję wykazu wspólnego, z tymi samymi
// wymiarami i z rodziną — a nie ich podzbiór.
//
// Pomiar liczby pozycji jest tu istotny, nie ozdobny: przekład gubiący pozycję
// zabiera ją Operatorowi po cichu, bo wykaz nadal wygląda na pełny.
func TestWykazDesignuJestPrzekladem(t *testing.T) {
	przelozony := nosnikiDruku()
	if len(przelozony) != len(wykazNosnikowDruku) {
		t.Fatalf("przekład Designu ma %d pozycji, a wykaz wspólny %d — przekład gubi "+
			"albo dokłada nośnik", len(przelozony), len(wykazNosnikowDruku))
	}
	for _, nosnik := range przelozony {
		wspolny, jest := nosnikDrukuONazwie(nosnik.Name)
		if !jest {
			t.Errorf("nośnik %s wyszedł z przekładu Designu, a nie stoi w wykazie wspólnym "+
				"— przekład nie ma skąd wziąć pozycji własnej", nosnik.Name)
			continue
		}
		if wspolny.SzerokoscMm != nosnik.WidthMm || wspolny.WysokoscMm != nosnik.HeightMm {
			t.Errorf("nośnik %s ma w przekładzie wymiary %v na %v mm, a w wykazie wspólnym "+
				"%v na %v", nosnik.Name,
				nosnik.WidthMm, nosnik.HeightMm, wspolny.SzerokoscMm, wspolny.WysokoscMm)
		}
		if nosnik.Family == nil || *nosnik.Family != wspolny.Rodzina {
			t.Errorf("nośnik %s stracił w przekładzie rodzinę (%v wobec %q) — bez niej "+
				"zawężenie wykazu rodziną przestaje działać", nosnik.Name,
				nosnik.Family, wspolny.Rodzina)
		}
	}
	// Koperty wskazane przez Właściciela muszą dojść do Designu, bo do tej pory
	// ich tam nie było — to jest cały powód tego przestawienia.
	for _, nazwa := range []string{"C4", "C5", "C6", "B4", "B5"} {
		if _, jest := nosnikPoNazwie(nazwa); !jest {
			t.Errorf("nośnik %s nie doszedł do wykazu Designu po przestawieniu", nazwa)
		}
	}
}

// TestNosnikiStudiaCzytajaWykazWspolny wykazuje, że przestawienie strony Studia
// jest przestawieniem, a nie kopią: format oddany komendą `studio.page.paper.list`
// niesie źródło wykazu wspólnego, a nie dawne „uzupełnienie Studia".
func TestNosnikiStudiaCzytajaWykazWspolny(t *testing.T) {
	wykaz := stronaNosniki()
	if len(wykaz) != len(wykazNosnikowDruku) {
		t.Fatalf("wykaz Studia ma %d pozycji, a wykaz wspólny %d — jedno z dwóch "+
			"dokłada coś od siebie", len(wykaz), len(wykazNosnikowDruku))
	}
	for _, nosnik := range wykaz {
		if nosnik.Source == nil || *nosnik.Source != zrodloWykazuNosnikowRdzenia {
			t.Errorf("nośnik %s nie mówi, że pochodzi z wykazu wspólnego", nosnik.Name)
		}
	}
	c4, jest := stronaNosnik("c4")
	if !jest {
		t.Fatal("nastawy strony Studia nie znajdują koperty C4")
	}
	if c4.Kind != shared.StudioPaperKindEnvelope {
		t.Errorf("koperta C4 ma w Studiu rodzaj %q", c4.Kind)
	}
	if _, jest := stronaNosnik("B5"); !jest {
		t.Error("nastawy strony Studia nie znajdują nośnika B5 — okno pokazywało go " +
			"z wykazu wbudowanego, a rdzeń odmawiał")
	}
}
