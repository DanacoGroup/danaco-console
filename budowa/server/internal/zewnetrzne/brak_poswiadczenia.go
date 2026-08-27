// Odmowa, która pada zanim żądanie wyjdzie w sieć, gdy kanał nie ma czym się
// uwierzytelnić; nazywa brak i drogę jego uzupełnienia.
package zewnetrzne

import "strings"

// Powody braku poświadczenia. Rozróżnienie jest po to, żeby każdy z nich miał
// własną drogę naprawy — kilka różnych usterek pod jednym zdaniem zostawia
// Operatora ze zgadywaniem.
const (
	// PowodSejfNiewpiety — wiersz wskazuje sejf poświadczeń, a rdzeń złożono bez
	// wpiętego sejfu do warstwy modeli.
	PowodSejfNiewpiety = "sejf-niewpiety"
	// PowodBrakWpisu — sejf poświadczeń stoi, ale wpisu o tej nazwie w nim brak
	// wśród zapisanych sekretów.
	PowodBrakWpisu = "brak-wpisu"
	// PowodPustyWpis — wpis w sejfie poświadczeń jest, lecz pusty; sekret
	// zgubiono gdzieś po drodze zapisu.
	PowodPustyWpis = "pusty-wpis"
	// PowodPustaZmienna — zmienna środowiskowa wskazana odwołaniem jest
	// nieustawiona albo pusta w otoczeniu.
	PowodPustaZmienna = "pusta-zmienna"
)

// BrakPoswiadczenia mówi, że kanał nie miał czym się uwierzytelnić i dlaczego;
// osobny typ odróżnia ten brak od odmowy dostawcy, która bywa chwilowa.
type BrakPoswiadczenia struct {
	// Kanal jest kodem kanału z wiersza rejestru.
	Kanal string
	// Powod jest rozpoznaniem braku (PowodSejfNiewpiety i dalsze).
	Powod string
	// Poswiadczenie opisuje odwołanie, które nie dało sekretu. Bez wartości.
	Poswiadczenie Poswiadczenie
}

// NowyBrakPoswiadczenia składa odmowę z kodu kanału, odwołania i rozpoznanego
// powodu; odwołanie przechodzi przez RozpoznajPoswiadczenie, które ustala
// jego rodzaj.
func NowyBrakPoswiadczenia(kanal, odwolanie, powod string) *BrakPoswiadczenia {
	return &BrakPoswiadczenia{
		Kanal:         strings.TrimSpace(kanal),
		Powod:         powod,
		Poswiadczenie: RozpoznajPoswiadczenie(odwolanie),
	}
}

// DotyczyPoswiadczenia jest zawsze prawdą — ten typ z definicji mówi o kluczu.
// Metoda istnieje po to, by wołacz mógł zadać jedno pytanie obu odmowom,
// zamiast rozgałęziać się po typie w każdym miejscu.
func (b *BrakPoswiadczenia) DotyczyPoswiadczenia() bool { return true }

// Error składa zdanie dla Operatora w tej samej kolejności, co odmowa dostawcy:
// co się stało, czego dotyczy, co z tym zrobić.
func (b *BrakPoswiadczenia) Error() string {
	return "kanał " + b.Kanal + ": " + b.opisBraku() + "; " + b.naprawaBraku()
}

// opisBraku nazywa sam brak, dobierając zdanie do rozpoznanego Powodu; każda
// gałąź switch opisuje inną usterkę źródła sekretu.
func (b *BrakPoswiadczenia) opisBraku() string {
	switch b.Powod {
	case PowodSejfNiewpiety:
		return "wiersz kanału wskazuje sejf poświadczeń (odwołanie „" + b.Poswiadczenie.Odwolanie +
			"”), a rdzeń pracuje BEZ WPIĘTEGO SEJFU — sekretu nie ma skąd wziąć " +
			"i żądanie nie zostało wysłane"
	case PowodBrakWpisu:
		return "w sejfie NIE MA WPISU „" + b.Poswiadczenie.Byt +
			"”, na który powołuje się wiersz kanału — żądanie nie zostało wysłane"
	case PowodPustyWpis:
		return "wpis sejfu „" + b.Poswiadczenie.Byt +
			"” jest PUSTY — żądanie nie zostało wysłane, bo poszłoby bez klucza"
	case PowodPustaZmienna:
		return "zmienna środowiskowa " + b.Poswiadczenie.Byt +
			" jest pusta albo nieustawiona w otoczeniu procesu rdzenia — " +
			"żądanie nie zostało wysłane"
	default:
		return "nie ma czym się uwierzytelnić pod odwołaniem „" +
			b.Poswiadczenie.Odwolanie + "” — żądanie nie zostało wysłane"
	}
}

// naprawaBraku dobiera drogę wyjścia. Przy niewpiętym sejfie droga nie należy
// do Operatora i zdanie mówi to wprost: gdyby zapisał sekret przez
// `account.update`, wpis powstałby poprawnie, a kanał i tak by go nie zobaczył.
func (b *BrakPoswiadczenia) naprawaBraku() string {
	if b.Powod == PowodSejfNiewpiety {
		return "to usterka SKŁADANIA RDZENIA, nie nastawa Operatora: montaż ma wpiąć " +
			"sejf do warstwy modeli (core/montaz_porty.go, models.UstawSejfPoswiadczen). " +
			"Obejście na teraz: wskazać kanałowi odwołanie do zmiennej środowiskowej " +
			"komendą " + WskazanieOdwolania
	}
	return b.Poswiadczenie.Naprawa()
}
