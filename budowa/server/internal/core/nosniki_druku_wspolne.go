// Plik zawiera jedyny w rdzeniu wykaz nośników druku, ze wszystkimi wymiarami w milimetrach zapisanymi w układzie pionowym, współdzielony przez moduł Design i moduł Studio.
package core

import (
	"strings"

	"danacoconsole/shared"
)

// Rodziny nośników. Napis, nie wyliczenie kontraktu, bo rodzina jest
// zawężeniem wykazu dla Operatora, a nie bytem, o który pyta komenda:
// `design.print.paper.list` przyjmuje ją polem `family` jako napis wolny.
const (
	rodzinaNosnikaISOA     = "iso-a"
	rodzinaNosnikaISOB     = "iso-b"
	rodzinaNosnikaKoperta  = "koperta"
	rodzinaNosnikaUS       = "us"
	rodzinaNosnikaHandlowy = "handlowy"
	rodzinaNosnikaWielki   = "wielkiformat"
)

// nosnikDruku to jedna pozycja wykazu wspólnego: nazwa, wymiary w milimetrach
// w układzie pionowym i rodzina, którą wykaz daje się zawężać.
type nosnikDruku struct {
	Nazwa       string
	SzerokoscMm float64
	WysokoscMm  float64
	Rodzina     string
}

// wykazNosnikowDruku jest JEDYNYM wykazem nośników w rdzeniu.
//
// Poprawka wymiaru wchodzi tutaj i obowiązuje wszystkich wołających naraz.
// Pozycja dopisana tu pojawia się i w wykazie Designu, i w nastawach strony
// Studia — bez pamiętania o drugim miejscu.
var wykazNosnikowDruku = []nosnikDruku{
	{Nazwa: "A0", SzerokoscMm: 841, WysokoscMm: 1189, Rodzina: rodzinaNosnikaISOA},
	{Nazwa: "A1", SzerokoscMm: 594, WysokoscMm: 841, Rodzina: rodzinaNosnikaISOA},
	{Nazwa: "A2", SzerokoscMm: 420, WysokoscMm: 594, Rodzina: rodzinaNosnikaISOA},
	{Nazwa: "A3", SzerokoscMm: 297, WysokoscMm: 420, Rodzina: rodzinaNosnikaISOA},
	{Nazwa: "A4", SzerokoscMm: 210, WysokoscMm: 297, Rodzina: rodzinaNosnikaISOA},
	{Nazwa: "A5", SzerokoscMm: 148, WysokoscMm: 210, Rodzina: rodzinaNosnikaISOA},
	{Nazwa: "A6", SzerokoscMm: 105, WysokoscMm: 148, Rodzina: rodzinaNosnikaISOA},
	{Nazwa: "B1", SzerokoscMm: 707, WysokoscMm: 1000, Rodzina: rodzinaNosnikaISOB},
	{Nazwa: "B2", SzerokoscMm: 500, WysokoscMm: 707, Rodzina: rodzinaNosnikaISOB},
	{Nazwa: "B3", SzerokoscMm: 353, WysokoscMm: 500, Rodzina: rodzinaNosnikaISOB},
	// B4 i B5 należą do szeregu ISO B, tak samo jak pozycje B1–B3 powyżej.
	{Nazwa: "B4", SzerokoscMm: 250, WysokoscMm: 353, Rodzina: rodzinaNosnikaISOB},
	{Nazwa: "B5", SzerokoscMm: 176, WysokoscMm: 250, Rodzina: rodzinaNosnikaISOB},
	// Koperta DL ma wymiar normy ISO 269, zapisany w układzie pionowym jako 110 na 220 milimetrów.
	{Nazwa: "DL", SzerokoscMm: 110, WysokoscMm: 220, Rodzina: rodzinaNosnikaKoperta},
	{Nazwa: "C4", SzerokoscMm: 229, WysokoscMm: 324, Rodzina: rodzinaNosnikaKoperta},
	{Nazwa: "C5", SzerokoscMm: 162, WysokoscMm: 229, Rodzina: rodzinaNosnikaKoperta},
	{Nazwa: "C6", SzerokoscMm: 114, WysokoscMm: 162, Rodzina: rodzinaNosnikaKoperta},
	{Nazwa: "Letter", SzerokoscMm: 215.9, WysokoscMm: 279.4, Rodzina: rodzinaNosnikaUS},
	{Nazwa: "Legal", SzerokoscMm: 215.9, WysokoscMm: 355.6, Rodzina: rodzinaNosnikaUS},
	{Nazwa: "Tabloid", SzerokoscMm: 279.4, WysokoscMm: 431.8, Rodzina: rodzinaNosnikaUS},
	// Wizytówka ma wymiar własny pionowy; orientację pozioma nakłada materiał.
	{Nazwa: "wizytowka", SzerokoscMm: 50, WysokoscMm: 90, Rodzina: rodzinaNosnikaHandlowy},
	{Nazwa: "rolka-1000", SzerokoscMm: 1000, WysokoscMm: 3000, Rodzina: rodzinaNosnikaWielki},
	{Nazwa: "rolka-1370", SzerokoscMm: 1370, WysokoscMm: 5000, Rodzina: rodzinaNosnikaWielki},
}

// nosnikDrukuONazwie znajduje nośnik po nazwie, nie bacząc na wielkość liter —
// Operator wpisuje „a4" tak samo często jak „A4".
func nosnikDrukuONazwie(nazwa string) (nosnikDruku, bool) {
	szukana := strings.ToUpper(strings.TrimSpace(nazwa))
	if szukana == "" {
		return nosnikDruku{}, false
	}
	for _, nosnik := range wykazNosnikowDruku {
		if strings.ToUpper(nosnik.Nazwa) == szukana {
			return nosnik, true
		}
	}
	return nosnikDruku{}, false
}

// nosnikiDrukuJakoFormatyStudia przekłada wykaz wspólny na formaty nośnika kontraktu Studia, wyznaczając kopertę dla rodziny koperta i arkusz dla pozostałych rodzin.
func nosnikiDrukuJakoFormatyStudia() []shared.StudioPaperFormat {
	wykaz := make([]shared.StudioPaperFormat, 0, len(wykazNosnikowDruku))
	for _, nosnik := range wykazNosnikowDruku {
		rodzaj := shared.StudioPaperKind(shared.StudioPaperKindSheet)
		if nosnik.Rodzina == rodzinaNosnikaKoperta {
			rodzaj = shared.StudioPaperKindEnvelope
		}
		wykaz = append(wykaz, shared.StudioPaperFormat{
			Name:     nosnik.Nazwa,
			Kind:     rodzaj,
			WidthMm:  nosnik.SzerokoscMm,
			HeightMm: nosnik.WysokoscMm,
			Source:   postacWskaznikTekstu(zrodloWykazuNosnikowRdzenia),
		})
	}
	return wykaz
}

// zrodloWykazuNosnikowRdzenia mówi Operatorowi i prowadzącemu, skąd pozycja
// wykazu pochodzi. Jedno źródło, bo jeden wykaz — dawne rozróżnienie „wykaz
// rdzenia" i „uzupełnienie Studia" opisywało dwa miejsca, których już nie ma.
const zrodloWykazuNosnikowRdzenia = "wspólny wykaz nośników druku rdzenia"
