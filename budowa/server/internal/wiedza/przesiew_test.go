// Sprawdziany drugiego przebiegu i osi obrazu — tego, co da się zmierzyć bez
// wag na dysku: układania kolejności, granic żądania i treści odmów.
package wiedza

import "testing"

// TestPrzesiewUkladaKolejnoscOdNowa pilnuje rzeczy, dla której przesiew istnieje:
// ocena kodera ma zastąpić kosinus, a nie dołożyć się do niego. Kandydat
// ostatni po wektorach ma prawo wyjść na czoło.
func TestPrzesiewUkladaKolejnoscOdNowa(t *testing.T) {
	kandydaci := []Trafienie{
		{Pozycja: Pozycja{Tresc: "pierwszy po wektorach"}, Podobienst: 0.9},
		{Pozycja: Pozycja{Tresc: "drugi po wektorach"}, Podobienst: 0.8},
		{Pozycja: Pozycja{Tresc: "trzeci po wektorach"}, Podobienst: 0.7},
	}

	przesiane := PoPrzesiewie(kandydaci, []float32{0.01, 0.02, 0.99}, 0)

	if len(przesiane) != 3 {
		t.Fatalf("przesiew oddał %d kandydatów zamiast trzech", len(przesiane))
	}
	if przesiane[0].Pozycja.Tresc != "trzeci po wektorach" {
		t.Fatalf("na czele stoi %q, choć koder ocenił najwyżej trzeciego — "+
			"przesiew nie ułożył kolejności od nowa", przesiane[0].Pozycja.Tresc)
	}
	if przesiane[0].Podobienst != 0.99 {
		t.Fatalf("trafność na czele wynosi %v zamiast oceny kodera 0.99 — "+
			"do odpowiedzi poszedł kosinus pierwszego przebiegu", przesiane[0].Podobienst)
	}
}

// TestPrzesiewPrzycinaDoZadanejLiczby pilnuje, że kandydatów wolno mieć więcej
// niż fragmentów w odpowiedzi — w tym cała rzecz: przesiew ma z czego wybierać.
func TestPrzesiewPrzycinaDoZadanejLiczby(t *testing.T) {
	kandydaci := []Trafienie{
		{Pozycja: Pozycja{Tresc: "a"}, Podobienst: 0.9},
		{Pozycja: Pozycja{Tresc: "b"}, Podobienst: 0.8},
		{Pozycja: Pozycja{Tresc: "c"}, Podobienst: 0.7},
	}

	przesiane := PoPrzesiewie(kandydaci, []float32{0.1, 0.9, 0.5}, 2)

	if len(przesiane) != 2 {
		t.Fatalf("odpowiedź niesie %d fragmentów zamiast dwóch żądanych", len(przesiane))
	}
	if przesiane[0].Pozycja.Tresc != "b" || przesiane[1].Pozycja.Tresc != "c" {
		t.Fatalf("po przycięciu zostały %q i %q zamiast dwóch najwyżej ocenionych",
			przesiane[0].Pozycja.Tresc, przesiane[1].Pozycja.Tresc)
	}
}

// TestOcenyNiezgodneZKandydatamiNiePrzestawiajaNiczego pilnuje warunku wiązania
// po pozycji: gdyby ocen było mniej niż kandydatów, przypisanie ich po kolei
// przesunęłoby oceny o jedną pozycję i dałoby ranking, który wygląda jak wynik.
func TestOcenyNiezgodneZKandydatamiNiePrzestawiajaNiczego(t *testing.T) {
	kandydaci := []Trafienie{
		{Pozycja: Pozycja{Tresc: "a"}, Podobienst: 0.9},
		{Pozycja: Pozycja{Tresc: "b"}, Podobienst: 0.8},
	}

	wynik := PoPrzesiewie(kandydaci, []float32{0.1}, 0)

	if len(wynik) != 2 || wynik[0].Pozycja.Tresc != "a" {
		t.Fatalf("wykaz ocen krótszy od wykazu kandydatów przestawił kolejność: %+v", wynik)
	}
}

// TestGranicaKandydatowTrzymaSieWidelek pilnuje, że żądanie nie wciągnie
// całego wskaźnika do kodera i że brak wskazania ma wartość własną.
func TestGranicaKandydatowTrzymaSieWidelek(t *testing.T) {
	if ile := GranicaKandydatowZadania(nil); ile != KandydaciDomyslni {
		t.Fatalf("brak wskazania dał %d kandydatów zamiast %d", ile, KandydaciDomyslni)
	}
	zero := 0
	if ile := GranicaKandydatowZadania(&zero); ile != KandydaciDomyslni {
		t.Fatalf("zero kandydatów dało %d zamiast wartości domyślnej", ile)
	}
	duzo := GranicaKandydatow + 1000
	if ile := GranicaKandydatowZadania(&duzo); ile != GranicaKandydatow {
		t.Fatalf("żądanie %d kandydatów przeszło jako %d — sufit nie działa", duzo, ile)
	}
	piec := 5
	if ile := GranicaKandydatowZadania(&piec); ile != 5 {
		t.Fatalf("żądanie pięciu kandydatów dało %d", ile)
	}
}

// TestObrazBezOcenyDodatniejNieWracaJakoTrafienie pilnuje, że plik, którego
// pomocnik nie otworzył, nie wraca jako odpowiedź. Ocena zerowa trzyma
// pozycję w wykazie, a nie stawia w wyniku.
func TestObrazBezOcenyDodatniejNieWracaJakoTrafienie(t *testing.T) {
	obrazy := []Obraz{
		{Zrodlo: "uszkodzony.png"},
		{Zrodlo: "plot.jpg"},
		{Zrodlo: "kot.jpg"},
	}

	najblizsze := NajblizszeObrazy(obrazy, []float32{0, 0.11, 0.31}, 0)

	if len(najblizsze) != 2 {
		t.Fatalf("wynik niesie %d obrazów zamiast dwóch policzonych", len(najblizsze))
	}
	if najblizsze[0].Zrodlo != "kot.jpg" {
		t.Fatalf("na czele stoi %q zamiast obrazu o najwyższej ocenie", najblizsze[0].Zrodlo)
	}
}

// TestOdmowaPrzesiewuNazywaBrakINaprawe pilnuje, że odmowa dołożonego silnika
// nie odsyła Operatora po bibliotekę osadzeń — brakuje czego innego i naprawia
// się to czym innym.
func TestOdmowaPrzesiewuNazywaBrakINaprawe(t *testing.T) {
	zdanie := (&BrakSilnika{
		Silnik: SilnikDlaPrzesiewu, Rodzaj: brakBiblioteki,
		Model: ModelPrzesiewuDomyslny, WagaMb: WagaPrzesiewuMb,
	}).Error()

	for _, czlon := range []string{"przesiew wyników", "torch", "transformers",
		ModelPrzesiewuDomyslny, "2,2 GB", "naprawa:"} {

		if !zawiera(zdanie, czlon) {
			t.Fatalf("odmowa przesiewu nie niesie członu %q: %s", czlon, zdanie)
		}
	}
	if zawiera(zdanie, "fastembed") {
		t.Fatalf("odmowa przesiewu odsyła po bibliotekę osadzeń, której brak nie dotyczy: %s",
			zdanie)
	}
}

// TestOdmowaOsiObrazuNazywaSwojeUstawienie pilnuje, że odmowa wskazuje katalog
// wag TEJ osi — wskazanie katalogu osadzarki kierowałoby Operatora do miejsca,
// w którym leży inny model.
func TestOdmowaOsiObrazuNazywaSwojeUstawienie(t *testing.T) {
	zdanie := (&BrakSilnika{
		Silnik: SilnikDlaObrazu, Rodzaj: brakWagStojacych,
		Model: ModelObrazuDomyslny, WagaMb: WagaObrazuMb,
	}).Error()

	if !zawiera(zdanie, KluczKatalogObrazu) {
		t.Fatalf("odmowa osi obrazu nie nazywa ustawienia %s: %s", KluczKatalogObrazu, zdanie)
	}
	if zawiera(zdanie, KluczKatalogModeli) {
		t.Fatalf("odmowa osi obrazu wskazuje katalog osadzarki: %s", zdanie)
	}
	if !zawiera(zdanie, "NAZW plików") {
		t.Fatalf("odmowa osi obrazu nie mówi, że rdzeń nie zejdzie na nazwy plików: %s", zdanie)
	}
}

// TestOdmowaOsadzarkiMowiToSamoCoMowila pilnuje, że dołożenie pola nie zmieniło
// odmowy silnika, który stał tu pierwszy.
func TestOdmowaOsadzarkiMowiToSamoCoMowila(t *testing.T) {
	zdanie := (&BrakSilnika{Rodzaj: brakBiblioteki, Model: ModelDomyslny,
		WagaMb: WagaModeluMb}).Error()

	if !zawiera(zdanie, "fastembed") {
		t.Fatalf("odmowa osadzarki przestała nazywać bibliotekę osadzeń: %s", zdanie)
	}
	if !zawiera(zdanie, "po SŁOWACH") {
		t.Fatalf("odmowa osadzarki przestała mówić, na co rdzeń nie zejdzie: %s", zdanie)
	}
}

// zawiera odpowiada na pytanie, czy zdanie niesie człon — bez sięgania po
// strings w każdym sprawdzianie z osobna.
func zawiera(zdanie, czlon string) bool {
	return len(czlon) == 0 || indeks(zdanie, czlon) >= 0
}

func indeks(zdanie, czlon string) int {
	for i := 0; i+len(czlon) <= len(zdanie); i++ {
		if zdanie[i:i+len(czlon)] == czlon {
			return i
		}
	}
	return -1
}
