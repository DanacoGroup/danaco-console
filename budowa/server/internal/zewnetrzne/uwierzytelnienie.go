// Jedna prawda o tym, co Operator czyta, gdy zewnętrzny dostawca odmówi
// kanałowi modelu — i o tym, skąd kanał bierze poświadczenie.
package zewnetrzne

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// PrzedrostekSejfu znakuje odwołanie wskazujące wpis sejfu poświadczeń,
// a nie zmienną środowiskową otoczenia procesu.
const PrzedrostekSejfu = "sejf:"

// Rodzaje źródła poświadczenia. Wartości są częścią zdania dla Operatora, więc
// brzmią po polsku i w dopełniaczu — wchodzą w „poświadczeniem z ...".
const (
	// ZrodloSejf — odwołanie wskazuje wpis sejfu poświadczeń, a nie nazwę
	// zmiennej środowiskowej otoczenia procesu.
	ZrodloSejf = "sejf"
	// ZrodloSrodowisko — odwołanie jest nazwą zmiennej środowiskowej,
	// odczytywanej z otoczenia procesu w chwili wywołania.
	ZrodloSrodowisko = "środowisko"
	// ZrodloBrak — kanał nie ma żadnego odwołania; punkt końcowy działa bez
	// uwierzytelnienia po stronie dostawcy.
	ZrodloBrak = "brak"
)

// Poswiadczenie opisuje, skąd kanał bierze sekret, i samego sekretu nie niesie.
// Tyle wystarcza, żeby powiedzieć Operatorowi, który klucz zawiódł, i żeby nie
// wnieść wartości klucza do komunikatu ani do dziennika.
type Poswiadczenie struct {
	// Odwolanie jest zapisem z wiersza rejestru, w całości („sejf:konto-firmowe"
	// albo „OPENAI_API_KEY").
	Odwolanie string
	// Zrodlo mówi, jaką drogą odwołanie się rozwiązuje: ZrodloSejf,
	// ZrodloSrodowisko albo ZrodloBrak.
	Zrodlo string
	// Byt jest nazwą samego wpisu, bez przedrostka sejfu.
	Byt string
}

// RozpoznajPoswiadczenie rozbiera odwołanie wiersza rejestru na źródło i byt.
// Odwołanie puste znaczy punkt końcowy bez uwierzytelnienia — to nie jest błąd,
// tylko trzecia możliwość, którą trzeba umieć nazwać.
func RozpoznajPoswiadczenie(odwolanie string) Poswiadczenie {
	oczyszczone := strings.TrimSpace(odwolanie)
	if oczyszczone == "" {
		return Poswiadczenie{Zrodlo: ZrodloBrak}
	}
	if byt, jest := strings.CutPrefix(oczyszczone, PrzedrostekSejfu); jest {
		return Poswiadczenie{
			Odwolanie: oczyszczone,
			Zrodlo:    ZrodloSejf,
			Byt:       strings.TrimSpace(byt),
		}
	}
	return Poswiadczenie{Odwolanie: oczyszczone, Zrodlo: ZrodloSrodowisko, Byt: oczyszczone}
}

// ZSejfu odpowiada, czy poświadczenie idzie z sejfu. Osobna metoda, bo pytanie
// pada w kilku miejscach, a porównanie z napisem rozjeżdża się przy literówce.
func (p Poswiadczenie) ZSejfu() bool { return p.Zrodlo == ZrodloSejf }

// Opis nazywa źródło poświadczenia tak, jak ma je przeczytać człowiek.
// Wchodzi w środek zdania, więc zaczyna się małą literą i nie ma kropki.
func (p Poswiadczenie) Opis() string {
	switch p.Zrodlo {
	case ZrodloSejf:
		return "poświadczeniem z sejfu, wpis „" + p.Byt + "”"
	case ZrodloSrodowisko:
		return "poświadczeniem ze zmiennej środowiskowej " + p.Byt
	default:
		return "bez poświadczenia (wiersz kanału nie ma odwołania)"
	}
}

// WskazanieOdwolania nazywa miejsce, w którym Operator przestawia odwołanie
// kanału — parametr wewnątrz pola konfiguracji, nie osobne pole żądania.
const WskazanieOdwolania = "channel.update (pole config, klucz „credentialRef”)"

// Naprawa mówi, co zrobić, żeby poświadczenie znów było ważne — komendami,
// które kontrakt rzeczywiście ma.
func (p Poswiadczenie) Naprawa() string {
	switch p.Zrodlo {
	case ZrodloSejf:
		return "naprawa: zapisać ważny klucz w sejfie pod wpisem „" + p.Byt +
			"” komendą account.update (pole credential; przy koncie nowym account.add), " +
			"albo wskazać kanałowi inne odwołanie komendą " + WskazanieOdwolania
	case ZrodloSrodowisko:
		return "naprawa: ustawić ważny klucz w zmiennej środowiskowej " + p.Byt +
			" i uruchomić serwer ponownie — zmienna czytana jest przy wysyłce, " +
			"z otoczenia procesu serwera; trwalej: przenieść klucz do sejfu " +
			"(account.update, pole credential) i wskazać kanałowi odwołanie „" +
			PrzedrostekSejfu + "<nazwa konta>” komendą " + WskazanieOdwolania
	default:
		return "naprawa: wskazać kanałowi odwołanie do poświadczenia komendą " +
			WskazanieOdwolania + " — punkt końcowy żąda uwierzytelnienia, " +
			"a wiersz kanału nie wskazuje żadnego klucza"
	}
}

// Rodzaje odmowy dostawcy. Rozróżnienie istnieje po to, żeby Operator nie
// szukał złego klucza, gdy skończyły się środki, i nie doładowywał konta,
// gdy klucz jest odwołany.
const (
	// OdmowaUwierzytelnienia — kod stanu 401: dostawca nie uznał podanego
	// poświadczenia za ważne w tym żądaniu.
	OdmowaUwierzytelnienia = "uwierzytelnienie"
	// OdmowaUprawnienia — kod stanu 403: poświadczenie zostało uznane, ale
	// sama czynność jest niedozwolona.
	OdmowaUprawnienia = "uprawnienie"
	// OdmowaPlatnosci — kod stanu 402: konto dostawcy jest bez środków albo
	// bez czynnego, opłaconego abonamentu.
	OdmowaPlatnosci = "płatność"
	// OdmowaNatezenia — kod stanu 429: przekroczone dopuszczalne natężenie
	// żądań albo przyznany przydział.
	OdmowaNatezenia = "natężenie"
	// OdmowaPoDostawcy — kody stanu 5xx: usterka leży po stronie zewnętrznego
	// dostawcy, nie samego kanału.
	OdmowaPoDostawcy = "dostawca"
	// OdmowaZadania — pozostałe kody: żądanie zostało odrzucone co do swojej
	// treści, nie z powodu poświadczenia.
	OdmowaZadania = "żądanie"
)

// RodzajOdmowy przekłada kod stanu odpowiedzi dostawcy na rodzaj odmowy,
// znany dalszym warstwom rdzenia.
func RodzajOdmowy(status int) string {
	switch {
	case status == http.StatusUnauthorized:
		return OdmowaUwierzytelnienia
	case status == http.StatusForbidden:
		return OdmowaUprawnienia
	case status == http.StatusPaymentRequired:
		return OdmowaPlatnosci
	case status == http.StatusTooManyRequests:
		return OdmowaNatezenia
	case status >= 500:
		return OdmowaPoDostawcy
	default:
		return OdmowaZadania
	}
}

// dotyczyPoswiadczenia mówi, czy przy tym rodzaju odmowy warto Operatorowi
// wskazywać klucz i drogę jego wymiany. Przy 5xx i przy odrzuconym żądaniu
// wskazywanie klucza wysłałoby go w złą stronę.
func dotyczyPoswiadczenia(rodzaj string) bool {
	switch rodzaj {
	case OdmowaUwierzytelnienia, OdmowaUprawnienia, OdmowaPlatnosci:
		return true
	default:
		return false
	}
}

// naglowekOdmowy jest pierwszym zdaniem komunikatu dla Operatora — mówi, co
// się stało, zanim padnie którykolwiek szczegół sprawy.
func naglowekOdmowy(rodzaj string) string {
	switch rodzaj {
	case OdmowaUwierzytelnienia:
		return "dostawca NIE UZNAŁ poświadczenia"
	case OdmowaUprawnienia:
		return "dostawca uznał poświadczenie, ale ODMÓWIŁ TEJ CZYNNOŚCI"
	case OdmowaPlatnosci:
		return "dostawca odmówił z powodu ROZLICZENIA konta"
	case OdmowaNatezenia:
		return "dostawca odrzucił żądanie z powodu NATĘŻENIA (za dużo wywołań albo wyczerpany przydział)"
	case OdmowaPoDostawcy:
		return "usterka PO STRONIE DOSTAWCY — poświadczenie nie jest tu winne"
	default:
		return "dostawca odrzucił żądanie"
	}
}

// naprawaOdmowy dobiera drogę naprawy do rodzaju odmowy. Dla odmów, które
// z poświadczeniem nie mają nic wspólnego, mówi wprost, czego nie robić — bez
// tego Operator zaczyna od wymiany działającego klucza.
func naprawaOdmowy(rodzaj string, p Poswiadczenie) string {
	switch rodzaj {
	case OdmowaUwierzytelnienia:
		return p.Naprawa()
	case OdmowaUprawnienia:
		return "naprawa: sprawdzić u dostawcy uprawnienia i zasięg tego klucza " +
			"(model, punkt końcowy, organizacja); " + p.Naprawa()
	case OdmowaPlatnosci:
		return "naprawa: uzupełnić środki albo abonament u dostawcy — wymiana klucza tu nie pomoże"
	case OdmowaNatezenia:
		return "naprawa: ponowić za chwilę albo zmniejszyć natężenie; " +
			"gdy wraca stale — podnieść przydział u dostawcy"
	case OdmowaPoDostawcy:
		return "naprawa: ponowić później; gdy wraca — sprawdzić stan usługi dostawcy"
	default:
		return "naprawa: sprawdzić parametry wiersza kanału (model, base_url, cialo_dodatkowe) " +
			"komendą channel.list i poprawić komendą channel.update (pole config)"
	}
}

// OdmowaKanaluZewnetrznego jest odmową zewnętrznego dostawcy w postaci, którą
// Operator umie przeczytać, a rdzeń rozpoznać po typie, nie po treści napisu.
type OdmowaKanaluZewnetrznego struct {
	// Kanal jest kodem kanału z wiersza rejestru.
	Kanal string
	// Status jest kodem stanu odpowiedzi HTTP.
	Status int
	// Rodzaj jest rozpoznaniem odmowy (OdmowaUwierzytelnienia i dalsze).
	Rodzaj string
	// Poswiadczenie opisuje źródło klucza. Bez wartości klucza.
	Poswiadczenie Poswiadczenie
	// KomunikatDostawcy jest tym, co dostawca powiedział o sobie sam —
	// oczyszczony i pozbawiony sekretu.
	KomunikatDostawcy string
}

// NowaOdmowaKanaluZewnetrznego składa odmowę z tego, co kanał ma pod ręką
// w chwili niepowodzenia: kodu kanału, odwołania z wiersza, kodu stanu i ciała
// odpowiedzi.
func NowaOdmowaKanaluZewnetrznego(kanal, odwolanie string, status int,
	cialoOdpowiedzi []byte, sekret string) *OdmowaKanaluZewnetrznego {

	rodzaj := RodzajOdmowy(status)
	return &OdmowaKanaluZewnetrznego{
		Kanal:             strings.TrimSpace(kanal),
		Status:            status,
		Rodzaj:            rodzaj,
		Poswiadczenie:     RozpoznajPoswiadczenie(odwolanie),
		KomunikatDostawcy: BezSekretu(KomunikatZCiala(cialoOdpowiedzi), sekret),
	}
}

// ZKomunikatem podmienia komunikat dostawcy na wyjęty ścieżką z wiersza
// rejestru. Pusty napis nie podmienia niczego.
func (o *OdmowaKanaluZewnetrznego) ZKomunikatem(komunikat, sekret string) *OdmowaKanaluZewnetrznego {
	if oczyszczony := BezSekretu(strings.TrimSpace(komunikat), sekret); oczyszczony != "" {
		o.KomunikatDostawcy = oczyszczony
	}
	return o
}

// DotyczyPoswiadczenia mówi wołającemu, czy tę odmowę naprawia się kluczem,
// czy przyczyna leży całkiem gdzie indziej.
func (o *OdmowaKanaluZewnetrznego) DotyczyPoswiadczenia() bool {
	return dotyczyPoswiadczenia(o.Rodzaj)
}

// Error składa zdanie dla Operatora w stałej kolejności: co się stało, czym
// kanał się uwierzytelniał, co powiedział dostawca, co z tym zrobić. Czytający
// pierwszy człon ma już wiedzieć, czy zawiódł jego klucz, czy cudza usługa.
func (o *OdmowaKanaluZewnetrznego) Error() string {
	czlony := []string{"kanał " + o.Kanal + ": " + naglowekOdmowy(o.Rodzaj) +
		" (odpowiedź " + strconv.Itoa(o.Status) + ")"}
	if dotyczyPoswiadczenia(o.Rodzaj) {
		czlony = append(czlony, "kanał uwierzytelnia się "+o.Poswiadczenie.Opis())
	}
	if o.KomunikatDostawcy != "" {
		czlony = append(czlony, "dostawca powiedział: "+o.KomunikatDostawcy)
	}
	czlony = append(czlony, naprawaOdmowy(o.Rodzaj, o.Poswiadczenie))
	return strings.Join(czlony, "; ")
}

// granicaKomunikatu ucina komunikat dostawcy. Ciało odpowiedzi bywa stroną
// błędu w HTML albo śladem stosu; powód stoi w pierwszym zdaniu, a reszta
// zalewa Operatora i dziennik.
const granicaKomunikatu = 400

// kluczeKomunikatu są miejscami, w których dostawcy trzymają zdanie o błędzie,
// a kolejność jest kolejnością szukania: od najbardziej szczegółowego.
var kluczeKomunikatu = []string{"message", "error_description", "detail", "error", "msg"}

// KomunikatZCiala wyjmuje z ciała odpowiedzi zdanie, które dostawca powiedział
// o sobie sam, gdy wiersz rejestru nie wskazuje własnej ścieżki błędu.
func KomunikatZCiala(cialo []byte) string {
	tresc := strings.TrimSpace(string(cialo))
	if tresc == "" {
		return ""
	}
	if wyjete := zJsonu(cialo); wyjete != "" {
		return skroc(wyjete)
	}
	return skroc(tresc)
}

// zJsonu schodzi po drzewie JSON, szukając pierwszego napisu pod znanym kluczem.
// Zagnieżdżenie („error" → obiekt → „message") jest u dostawców regułą, więc
// schodzenie idzie w głąb, a nie tylko po wierzchu.
func zJsonu(cialo []byte) string {
	var rozebrane any
	if err := json.Unmarshal(cialo, &rozebrane); err != nil {
		return ""
	}
	return szukajKomunikatu(rozebrane, 0)
}

// glebokoscSzukania ogranicza schodzenie w głąb odpowiedzi, żeby zagnieżdżona
// albo złośliwie zbudowana odpowiedź nie kosztowała stosu wywołań.
const glebokoscSzukania = 6

// szukajKomunikatu przegląda rozebrany zapis JSON w poszukiwaniu zdania
// o błędzie, schodząc rekurencyjnie w głąb struktury.
func szukajKomunikatu(wezel any, glebokosc int) string {
	if glebokosc > glebokoscSzukania {
		return ""
	}
	switch typowy := wezel.(type) {
	case map[string]any:
		return zObiektu(typowy, glebokosc)
	case []any:
		for _, element := range typowy {
			if znalezione := szukajKomunikatu(element, glebokosc+1); znalezione != "" {
				return znalezione
			}
		}
	}
	return ""
}

// zObiektu szuka zdania najpierw pod znanymi kluczami tego poziomu, a dopiero
// potem schodzi niżej — komunikat bliższy wierzchu jest zwykle tym właściwym.
func zObiektu(obiekt map[string]any, glebokosc int) string {
	for _, klucz := range kluczeKomunikatu {
		if napis, jest := obiekt[klucz].(string); jest && strings.TrimSpace(napis) != "" {
			return strings.TrimSpace(napis)
		}
	}
	for _, klucz := range kluczeKomunikatu {
		if gniazdo, jest := obiekt[klucz]; jest {
			if znalezione := szukajKomunikatu(gniazdo, glebokosc+1); znalezione != "" {
				return znalezione
			}
		}
	}
	return ""
}

// skroc przycina komunikat do granicy długości, zostawiając na końcu znak,
// że dalsza treść została ucięta.
func skroc(tresc string) string {
	if len(tresc) <= granicaKomunikatu {
		return tresc
	}
	return strings.TrimSpace(tresc[:granicaKomunikatu]) + "…"
}

// najkrotszySekret jest progiem, poniżej którego wycinanie robi więcej szkody
// niż pożytku: krótki napis trafiłby w zwykłe słowo komunikatu i zamazał je.
const najkrotszySekret = 8

// BezSekretu wycina wartość klucza z tekstu idącego do Operatora i do dziennika,
// choć sekretu tu formalnie nie ma.
func BezSekretu(tekst, sekret string) string {
	if tekst == "" {
		return tekst
	}
	for _, kandydat := range kandydaciSekretu(sekret) {
		tekst = strings.ReplaceAll(tekst, kandydat, "‹ukryty klucz›")
	}
	return tekst
}

// kandydaciSekretu rozkłada podaną wartość na napisy warte wycięcia: całość
// oraz człon po ostatniej spacji (klucz bez przedrostka schematu).
func kandydaciSekretu(sekret string) []string {
	oczyszczony := strings.TrimSpace(sekret)
	if len(oczyszczony) < najkrotszySekret {
		return nil
	}
	kandydaci := []string{oczyszczony}
	if miejsce := strings.LastIndex(oczyszczony, " "); miejsce >= 0 {
		if ogon := strings.TrimSpace(oczyszczony[miejsce+1:]); len(ogon) >= najkrotszySekret {
			kandydaci = append(kandydaci, ogon)
		}
	}
	return kandydaci
}

// SekretZOdpowiedzi odzyskuje wartość klucza z żądania, które tę odpowiedź
// wywołało — po to i tylko po to, żeby ją z komunikatu wyciąć.
func SekretZOdpowiedzi(odpowiedz *http.Response, nazwaNaglowka string) string {
	if odpowiedz == nil || odpowiedz.Request == nil {
		return ""
	}
	nazwa := strings.TrimSpace(nazwaNaglowka)
	if nazwa == "" {
		nazwa = "Authorization"
	}
	return odpowiedz.Request.Header.Get(nazwa)
}
