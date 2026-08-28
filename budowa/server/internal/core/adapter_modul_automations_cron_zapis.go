// Odpowiedzialność pliku: rozkład zapisu cron na zbiory dopuszczalnych wartości
// pięciu pól. Rachunek terminu stoi w ..._cron.go — tu leży wyłącznie odczyt
// zapisu.
package core

import (
	"strconv"
	"strings"
)

// polaCronZTekstu rozkłada zapis pięciopolowy cron na zbiory dopuszczalnych wartości wszystkich pięciu pól.
func polaCronZTekstu(zapis string) (polaCron, bool) {
	czesci := strings.Fields(strings.TrimSpace(zapis))
	if len(czesci) != 5 {
		return polaCron{}, false
	}
	minuty, ok1 := zbiorPolaCron(czesci[0], 0, 59)
	godziny, ok2 := zbiorPolaCron(czesci[1], 0, 23)
	dni, ok3 := zbiorPolaCron(czesci[2], 1, 31)
	miesiace, ok4 := zbiorPolaCron(czesci[3], 1, 12)
	tygodnie, ok5 := zbiorPolaCron(czesci[4], 0, 6)
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 {
		return polaCron{}, false
	}
	return polaCron{
		minuty: minuty, godziny: godziny, dniMiesiaca: dni, miesiace: miesiace,
		dniTygodnia:    tygodnie,
		dzienDowolny:   strings.TrimSpace(czesci[2]) == "*",
		tydzienDowolny: strings.TrimSpace(czesci[4]) == "*",
	}, true
}

// zbiorPolaCron rozkłada jedno pole zapisu cron na zbiór dopuszczalnych wartości liczbowych tego pola.
func zbiorPolaCron(pole string, dolna, gorna int) (map[int]bool, bool) {
	zbior := map[int]bool{}
	for _, czlon := range strings.Split(strings.TrimSpace(pole), ",") {
		if !dopiszCzlonCron(zbior, czlon, dolna, gorna) {
			return nil, false
		}
	}
	if len(zbior) == 0 {
		return nil, false
	}
	return zbior, true
}

// dopiszCzlonCron dokłada do zbioru wartości jednego członu pola zapisu cron tego harmonogramu automatyki.
func dopiszCzlonCron(zbior map[int]bool, czlon string, dolna, gorna int) bool {
	zakres, krokTekst, maKrok := strings.Cut(strings.TrimSpace(czlon), "/")
	krok := 1
	if maKrok {
		wartosc, err := strconv.Atoi(krokTekst)
		if err != nil || wartosc <= 0 {
			return false
		}
		krok = wartosc
	}
	od, do := dolna, gorna
	if zakres != "*" {
		poczatek, koniec, maZakres := strings.Cut(zakres, "-")
		wartoscOd, err := strconv.Atoi(strings.TrimSpace(poczatek))
		if err != nil {
			return false
		}
		od = wartoscOd
		do = wartoscOd
		if maZakres {
			wartoscDo, err := strconv.Atoi(strings.TrimSpace(koniec))
			if err != nil {
				return false
			}
			do = wartoscDo
		}
	}
	if od < dolna || do > gorna || od > do {
		return false
	}
	for wartosc := od; wartosc <= do; wartosc += krok {
		zbior[wartosc] = true
	}
	return true
}
