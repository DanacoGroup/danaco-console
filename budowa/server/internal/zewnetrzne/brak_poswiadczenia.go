// Odmowa, która pada zanim żądanie wyjdzie w sieć — gdy kanał nie ma czym się
// uwierzytelnić.
//
// Plik stoi osobno od `uwierzytelnienie.go`, bo to odmowa innej klasy. Tam
// dostawca coś odpowiedział i trzeba tę odpowiedź przetłumaczyć; tutaj nikt nie
// został zapytany, bo pod odwołaniem wskazanym w wierszu kanału nie ma sekretu.
// Zdanie odmowy ma nazwać brak i drogę jego uzupełnienia, bo sam komunikat
// „poświadczenia nie ma" nie mówi Operatorowi, co zrobić.
package zewnetrzne

import "strings"

// Powody braku poświadczenia. Rozróżnienie jest po to, żeby każdy z nich miał
// własną drogę naprawy — kilka różnych usterek pod jednym zdaniem zostawia
// Operatora ze zgadywaniem.
const (
	// PowodSejfNiewpiety — wiersz wskazuje sejf, a rdzeń złożono bez sejfu.
	// To usterka montażu, nie nastawa Operatora.
	PowodSejfNiewpiety = "sejf-niewpiety"
	// PowodBrakWpisu — sejf stoi, ale wpisu o tej nazwie w nim nie ma.
	PowodBrakWpisu = "brak-wpisu"
	// PowodPustyWpis — wpis jest, lecz pusty; sekret zgubiono po drodze.
	PowodPustyWpis = "pusty-wpis"
	// PowodPustaZmienna — zmienna środowiskowa nieustawiona albo pusta.
	PowodPustaZmienna = "pusta-zmienna"
)

// BrakPoswiadczenia mówi, że kanał nie miał czym się uwierzytelnić — i dlaczego.
//
// Osobny typ, bo warstwa wyżej ma odróżnić „nie ma czym wysłać" od „wysłano
// i odbiło się": pierwsze naprawia się w produkcie i ponawianie nic nie da,
// drugie bywa chwilowe. Bez typu obie klasy wyglądają jak ten sam napis.
type BrakPoswiadczenia struct {
	// Kanal jest kodem kanału z wiersza rejestru.
	Kanal string
	// Powod jest rozpoznaniem braku (PowodSejfNiewpiety i dalsze).
	Powod string
	// Poswiadczenie opisuje odwołanie, które nie dało sekretu. Bez wartości.
	Poswiadczenie Poswiadczenie
}

// NowyBrakPoswiadczenia składa odmowę z kodu kanału, odwołania i rozpoznanego
// powodu.
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

// opisBraku nazywa sam brak.
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
