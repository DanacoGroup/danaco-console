// Odpowiedzialność pliku: zapis iCalendar (RFC 5545) w zakresie, którego używa kalendarz
// projektu — własny czytnik i składacz wydarzeń `VEVENT`, wkompilowany w rdzeń.
package core

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// wydarzenieIcalWorkspace to jedno wydarzenie odczytane z pliku `.ics` przed przełożeniem na kalendarz.
type wydarzenieIcalWorkspace struct {
	Uid        string
	Tytul      string
	PoczatekMs int64
	KoniecMs   int64
	CalyDzien  bool
	Regula     string
}

// rozbierzIcalWorkspace czyta wydarzenia z treści pliku iCal, a drugi wynik niesie powody pominięcia wydarzeń.
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

// zlozIcalWorkspace składa plik iCal z wydarzeń projektu, gotowy do wydania Operatorowi na eksport z rdzenia.
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

// rozwinWierszeIcalWorkspace skleja wiersze złamane zapisem RFC 5545, gdzie ciąg dalszy zaczyna spacja.
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

// rozbierzWierszIcalWorkspace rozdziela wiersz na nazwę pola, jego parametry i wartość zapisu pliku iCal.
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

// chwilaIcalWorkspace czyta wartość czasu wraz z rozpoznaniem pozycji całodniowej, w strefie czasu UTC.
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

// odkodujTekstIcalWorkspace zdejmuje znaki chronione zapisu iCal, przywracając tekst pierwotny treści wydarzenia.
func odkodujTekstIcalWorkspace(wartosc string) string {
	zamiennik := strings.NewReplacer(`\n`, "\n", `\N`, "\n", `\,`, ",", `\;`, ";", `\\`, `\`)
	return zamiennik.Replace(wartosc)
}

// zakodujTekstIcalWorkspace chroni znaki, które w zapisie iCal mają znaczenie składniowe całego pliku.
func zakodujTekstIcalWorkspace(wartosc string) string {
	zamiennik := strings.NewReplacer(`\`, `\\`, ";", `\;`, ",", `\,`, "\n", `\n`)
	return zamiennik.Replace(wartosc)
}

// znacznikIcalWorkspace zapisuje chwilę w mierze kontraktu jako znacznik czasu UTC formatu zapisu iCal.
func znacznikIcalWorkspace(milisekundy int64) string {
	return time.UnixMilli(milisekundy).UTC().Format("20060102T150405Z")
}

// dataIcalWorkspace zapisuje chwilę jako samą datę, właściwą dla pozycji całodniowej zapisu pliku iCal.
func dataIcalWorkspace(milisekundy int64) string {
	return time.UnixMilli(milisekundy).UTC().Format("20060102")
}

// uidIcalWorkspace nadaje identyfikator wydarzeniu bez UID, żeby powtórne wciągnięcie nie podwoiło kalendarza.
func uidIcalWorkspace(tytul string, poczatek int64) string {
	odcisk := odciskTresci(pierwszaNiepustaWorkspace(tytul, "bez tytulu"))
	return "danaco-" + strconv.FormatInt(poczatek, 36) + "-" + odcisk[:12]
}
