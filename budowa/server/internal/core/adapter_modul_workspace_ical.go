// Odpowiedzialność pliku: zapis iCalendar (RFC 5545) w zakresie, którego
// używa kalendarz projektu — czytanie i składanie wydarzeń `VEVENT`.
//
// ── Dlaczego własny czytnik, a nie biblioteka z sieci ──────────────────────
// Instalka produktu niesie jedno binarium i nie wolno jej rozszerzać
// o zależność, której nie ma na maszynie budującej. Zakres potrzebny
// kalendarzowi projektu to sześć pól jednego składnika (`UID`, `SUMMARY`,
// `DTSTART`, `DTEND`, `RRULE`, `DTSTAMP`) wraz z rozwijaniem złamanych wierszy.
// To jest czytnik na dwieście wierszy, a nie warstwa kalendarza — i jest
// wkompilowany w rdzeń, więc wciągnięcie pliku `.ics` działa na maszynie
// Operatora tak samo jak na serwerze.
//
// ── Pominięcie mówi, dlaczego ──────────────────────────────────────────────
// Wydarzenie bez tytułu albo bez czytelnego początku nie wchodzi do projektu,
// ale wraca powodem w polu `skippedReasons`. Milczące pomijanie zostawiłoby
// Operatora z kalendarzem niepełnym i bez śladu, czego w nim brakuje.
package core

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// wydarzenieIcalWorkspace to jedno wydarzenie odczytane z pliku `.ics`.
type wydarzenieIcalWorkspace struct {
	Uid        string
	Tytul      string
	PoczatekMs int64
	KoniecMs   int64
	CalyDzien  bool
	Regula     string
}

// rozbierzIcalWorkspace czyta wydarzenia z treści pliku iCal. Drugi wynik
// niesie powody pominięcia wydarzeń, których nie dało się przyjąć.
func rozbierzIcalWorkspace(tresc string) ([]wydarzenieIcalWorkspace, []string) {
	wiersze := rozwinWierszeIcalWorkspace(tresc)
	wydarzenia := []wydarzenieIcalWorkspace{}
	powody := []string{}

	wSkladniku := false
	biezace := wydarzenieIcalWorkspace{}
	for _, wiersz := range wiersze {
		gorny := strings.ToUpper(strings.TrimSpace(wiersz))
		switch {
		case gorny == "BEGIN:VEVENT":
			wSkladniku, biezace = true, wydarzenieIcalWorkspace{}
			continue
		case gorny == "END:VEVENT":
			if !wSkladniku {
				continue
			}
			wSkladniku = false
			if strings.TrimSpace(biezace.Tytul) == "" {
				powody = append(powody, "wydarzenie bez pola SUMMARY — pozycja bez tytułu")
				continue
			}
			if biezace.PoczatekMs == 0 {
				powody = append(powody, "wydarzenie „"+biezace.Tytul+
					"” bez czytelnego pola DTSTART")
				continue
			}
			wydarzenia = append(wydarzenia, biezace)
			continue
		}
		if !wSkladniku {
			continue
		}
		nazwa, parametry, wartosc := rozbierzWierszIcalWorkspace(wiersz)
		switch nazwa {
		case "UID":
			biezace.Uid = wartosc
		case "SUMMARY":
			biezace.Tytul = odkodujTekstIcalWorkspace(wartosc)
		case "DTSTART":
			chwila, calyDzien, err := chwilaIcalWorkspace(parametry, wartosc)
			if err == nil {
				biezace.PoczatekMs, biezace.CalyDzien = chwila, calyDzien
			}
		case "DTEND":
			if chwila, _, err := chwilaIcalWorkspace(parametry, wartosc); err == nil {
				biezace.KoniecMs = chwila
			}
		case "RRULE":
			biezace.Regula = wartosc
		}
	}
	return wydarzenia, powody
}

// zlozIcalWorkspace składa plik iCal z wydarzeń projektu.
func zlozIcalWorkspace(nazwaKalendarza string, wydarzenia []wydarzenieIcalWorkspace) string {
	var zapis strings.Builder
	zapis.WriteString("BEGIN:VCALENDAR\r\nVERSION:2.0\r\n")
	zapis.WriteString("PRODID:-//Danaco Holding Group//Danaco Console//PL\r\n")
	zapis.WriteString("CALSCALE:GREGORIAN\r\n")
	zapis.WriteString("X-WR-CALNAME:" + zakodujTekstIcalWorkspace(nazwaKalendarza) + "\r\n")
	znacznik := time.Now().UTC().Format("20060102T150405Z")
	for _, wydarzenie := range wydarzenia {
		zapis.WriteString("BEGIN:VEVENT\r\n")
		zapis.WriteString("UID:" + wydarzenie.Uid + "\r\n")
		zapis.WriteString("DTSTAMP:" + znacznik + "\r\n")
		zapis.WriteString("SUMMARY:" + zakodujTekstIcalWorkspace(wydarzenie.Tytul) + "\r\n")
		if wydarzenie.CalyDzien {
			zapis.WriteString("DTSTART;VALUE=DATE:" + dataIcalWorkspace(wydarzenie.PoczatekMs) + "\r\n")
			if wydarzenie.KoniecMs != 0 {
				zapis.WriteString("DTEND;VALUE=DATE:" + dataIcalWorkspace(wydarzenie.KoniecMs) + "\r\n")
			}
		} else {
			zapis.WriteString("DTSTART:" + znacznikIcalWorkspace(wydarzenie.PoczatekMs) + "\r\n")
			if wydarzenie.KoniecMs != 0 {
				zapis.WriteString("DTEND:" + znacznikIcalWorkspace(wydarzenie.KoniecMs) + "\r\n")
			}
		}
		if wydarzenie.Regula != "" {
			zapis.WriteString("RRULE:" + wydarzenie.Regula + "\r\n")
		}
		zapis.WriteString("END:VEVENT\r\n")
	}
	zapis.WriteString("END:VCALENDAR\r\n")
	return zapis.String()
}

// rozwinWierszeIcalWorkspace skleja wiersze złamane zapisem RFC 5545: wiersz
// zaczynający się od spacji albo tabulatora jest dalszym ciągiem poprzedniego.
func rozwinWierszeIcalWorkspace(tresc string) []string {
	surowe := strings.Split(strings.ReplaceAll(tresc, "\r\n", "\n"), "\n")
	rozwiniete := []string{}
	for _, wiersz := range surowe {
		wiersz = strings.TrimRight(wiersz, "\r")
		if wiersz == "" {
			continue
		}
		if (strings.HasPrefix(wiersz, " ") || strings.HasPrefix(wiersz, "\t")) &&
			len(rozwiniete) > 0 {
			rozwiniete[len(rozwiniete)-1] += wiersz[1:]
			continue
		}
		rozwiniete = append(rozwiniete, wiersz)
	}
	return rozwiniete
}

// rozbierzWierszIcalWorkspace rozdziela wiersz na nazwę pola, jego parametry
// i wartość.
func rozbierzWierszIcalWorkspace(wiersz string) (string, map[string]string, string) {
	dwukropek := strings.Index(wiersz, ":")
	if dwukropek < 0 {
		return "", nil, ""
	}
	naglowek, wartosc := wiersz[:dwukropek], wiersz[dwukropek+1:]
	czesci := strings.Split(naglowek, ";")
	parametry := map[string]string{}
	for _, parametr := range czesci[1:] {
		if znak := strings.Index(parametr, "="); znak > 0 {
			parametry[strings.ToUpper(parametr[:znak])] = strings.ToUpper(parametr[znak+1:])
		}
	}
	return strings.ToUpper(strings.TrimSpace(czesci[0])), parametry, wartosc
}

// chwilaIcalWorkspace czyta wartość czasu wraz z rozpoznaniem pozycji
// całodniowej. Strefa nazwana parametrem TZID czytana jest jako czas UTC —
// rdzeń nie prowadzi bazy stref, a przesunięcie zgadywane byłoby gorsze niż
// przesunięcie nazwane wprost.
func chwilaIcalWorkspace(parametry map[string]string, wartosc string) (int64, bool, error) {
	wartosc = strings.TrimSpace(wartosc)
	if parametry["VALUE"] == "DATE" || len(wartosc) == 8 {
		chwila, err := time.Parse("20060102", wartosc)
		if err != nil {
			return 0, false, err
		}
		return chwila.UTC().UnixMilli(), true, nil
	}
	for _, uklad := range []string{"20060102T150405Z", "20060102T150405"} {
		if chwila, err := time.Parse(uklad, wartosc); err == nil {
			return chwila.UTC().UnixMilli(), false, nil
		}
	}
	return 0, false, fmt.Errorf("nieczytelna chwila iCal %q", wartosc)
}

// odkodujTekstIcalWorkspace zdejmuje znaki chronione zapisu iCal.
func odkodujTekstIcalWorkspace(wartosc string) string {
	zamiennik := strings.NewReplacer(`\n`, "\n", `\N`, "\n", `\,`, ",", `\;`, ";", `\\`, `\`)
	return zamiennik.Replace(wartosc)
}

// zakodujTekstIcalWorkspace chroni znaki, które w zapisie iCal mają znaczenie.
func zakodujTekstIcalWorkspace(wartosc string) string {
	zamiennik := strings.NewReplacer(`\`, `\\`, ";", `\;`, ",", `\,`, "\n", `\n`)
	return zamiennik.Replace(wartosc)
}

// znacznikIcalWorkspace zapisuje chwilę w mierze kontraktu jako znacznik UTC.
func znacznikIcalWorkspace(milisekundy int64) string {
	return time.UnixMilli(milisekundy).UTC().Format("20060102T150405Z")
}

// dataIcalWorkspace zapisuje chwilę jako samą datę — zapis pozycji całodniowej.
func dataIcalWorkspace(milisekundy int64) string {
	return time.UnixMilli(milisekundy).UTC().Format("20060102")
}

// uidIcalWorkspace nadaje identyfikator wydarzeniu, które przyszło bez UID.
// Bez niego powtórne wciągnięcie tego samego pliku podwoiłoby kalendarz.
func uidIcalWorkspace(tytul string, poczatek int64) string {
	odcisk := odciskTresci(pierwszaNiepustaWorkspace(tytul, "bez tytulu"))
	return "danaco-" + strconv.FormatInt(poczatek, 36) + "-" + odcisk[:12]
}
