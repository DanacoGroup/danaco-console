// Odpowiedzialność pliku: JEDEN wykaz nośników druku dla całego rdzenia —
// arkusze, koperty i nośniki wielkoformatowe wraz z wymiarami w milimetrach.
//
// ── Dlaczego wykaz jest jeden i stoi poza modułem ───────────────────────────
// Nośnik czytają dziś dwa moduły: Design wykazem `design.print.paper.list`
// i wydaniem do druku, Studio nastawami strony (`studio.page.paper.list`,
// `studio.page.setup.set`) oraz podglądem wydania. Dopóki każdy z nich miał
// wykaz własny, wykazy się rozjeżdżały — i rozjechały się naprawdę: koperta DL
// miała w rdzeniu 99 na 210 mm, a w oknie Studia 220 na 110 mm, szereg B
// urywał się w rdzeniu na B3, a w oknie sięgał B5. Rozstrzygnięcie Właściciela
// z 17.08.2026 mówi wprost: wykaz stoi w miejscu wspólnym, bo dwa wykazy przy
// pierwszej poprawce dadzą Operatorowi A3 o wymiarach A4.
//
// Plik nie należy więc do żadnego modułu i nie woła ani jednej ich funkcji:
// jest wiedzą rdzenia o materiale, tak samo jak przelicznik milimetrów na cale.
//
// ── Dlaczego wymiary są w układzie pionowym ─────────────────────────────────
// Nośnik ma jeden wymiar własny; obrót jest rozstrzygnięciem MATERIAŁU, nie
// nośnika. Koperta zapisana szerokością większą od wysokości niosłaby obrót
// w samej definicji i nikt nie umiałby powiedzieć, czy C4 stoi, czy leży.
// Wszystkie pozycje są więc pionowe (szerokość nie większa od wysokości),
// a orientację nakłada nastawa strony.
//
// ── Dlaczego koperty C mają wymiary normy, a nie zaokrąglone ────────────────
// Norma ISO 269 wiąże koperty C z szeregiem A: C4 bierze arkusz A4 bez
// zginania, C5 arkusz zgięty raz, C6 zgięty dwa razy. Właściciel wymienił C4,
// C5 i C6 wprost wraz z wymiarami i te wymiary tu stoją — 229 na 324, 162 na
// 229 oraz 114 na 162 mm.
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
	// B4 i B5 stały dotąd wyłącznie w wykazie okna Studia. Wykaz rdzenia ich nie
	// niósł, więc Operator, który wybrał B5 w oknie, dostawał od rdzenia odmowę
	// „nośnika nie znam" — a wybrał go z wykazu, który okno mu pokazało.
	{Nazwa: "B4", SzerokoscMm: 250, WysokoscMm: 353, Rodzina: rodzinaNosnikaISOB},
	{Nazwa: "B5", SzerokoscMm: 176, WysokoscMm: 250, Rodzina: rodzinaNosnikaISOB},
	// Koperty. DL jest wymiarem normy ISO 269 (110 na 220 mm) zapisanym pionowo;
	// wykaz Designu niósł 99 na 210, a wykaz okna Studia 220 na 110 — obie
	// liczby były inne i obie były wykazem tego samego nośnika.
	{Nazwa: "DL", SzerokoscMm: 110, WysokoscMm: 220, Rodzina: rodzinaNosnikaKoperta},
	{Nazwa: "C4", SzerokoscMm: 229, WysokoscMm: 324, Rodzina: rodzinaNosnikaKoperta},
	{Nazwa: "C5", SzerokoscMm: 162, WysokoscMm: 229, Rodzina: rodzinaNosnikaKoperta},
	{Nazwa: "C6", SzerokoscMm: 114, WysokoscMm: 162, Rodzina: rodzinaNosnikaKoperta},
	{Nazwa: "Letter", SzerokoscMm: 215.9, WysokoscMm: 279.4, Rodzina: rodzinaNosnikaUS},
	{Nazwa: "Legal", SzerokoscMm: 215.9, WysokoscMm: 355.6, Rodzina: rodzinaNosnikaUS},
	{Nazwa: "Tabloid", SzerokoscMm: 279.4, WysokoscMm: 431.8, Rodzina: rodzinaNosnikaUS},
	// Wizytówka jest pozioma w użyciu, ale wymiar własny ma pionowy — obrót
	// nakłada materiał, tak samo jak przy kopercie.
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

// nosnikiDrukuJakoFormatyStudia przekłada wykaz wspólny na formaty nośnika
// kontraktu Studia.
//
// Rodzaj bierze się z rodziny: rodzina `koperta` daje kopertę, wszystko inne
// arkusz. Rodzaj `custom` w wykazie nie stoi z zamysłu — wymiar własny podaje
// Operator polami `widthMm` i `heightMm`, więc pozycja katalogowa o takim
// rodzaju byłaby pozycją, której nie da się wybrać.
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
