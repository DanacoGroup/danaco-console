package session

import (
	"time"

	"danacoconsole/shared"
)

// Ustawienia niosą komplet parametrów wykonania okna komunikacji.
// Każde okno jednej sesji może mieć inny moduł, inny kanał modelu i inny
// zestaw katalogów roboczych — na tym polega rozdzielenie sesji od wykonania.
type Ustawienia struct {
	// Modul okna; przeniesiony z sesji na okno.
	Modul string
	// KanalModelu wskazuje wiersz rejestru kanałów.
	KanalModelu string
	// KatalogiRobocze są LISTĄ, nie pojedynczą wartością. Pierwszy katalog jest
	// katalogiem uruchomienia procesu okna.
	KatalogiRobocze []string
	// SrodowiskoWykonania mówi, gdzie pracuje model — niezależnie od tego, gdzie
	// stoi rdzeń.
	SrodowiskoWykonania shared.ExecutionEnv
	// TrybUprawnien ze słownika kontraktu.
	TrybUprawnien shared.PermissionMode
	// RolaOkna w pętli koordynator–wykonawca.
	RolaOkna shared.WindowRole
	// OknoKoordynatora wypełnia wyłącznie okno wykonawcy; puste dla pozostałych ról.
	OknoKoordynatora string
	// Tytul okna; pusty jest dopuszczalny.
	Tytul string
	// Agent jest KODEM eksperta nałożonego na kanał modelu tego okna; pusty
	// znaczy model surowy (`store/migracja_084_agent_okna.sql`).
	//
	// Ekspert NIE ZASTĘPUJE kanału — kanał jest drogą do modelu, ekspert
	// tożsamością nałożoną na tę drogę. Dlatego stoi obok `KanalModelu`,
	// a nie zamiast niego: okno bez eksperta rusza dokładnie tak, jak ruszało
	// przedtem, a okno z ekspertem tą samą drogą, tylko z jego warstwami,
	// modelem i nastawami.
	//
	// Kod, nie identyfikator wiersza — ekspert bywa kasowany niezależnie od
	// okien, w których pracował, a kod nierozpoznany jest faktem czytelnym:
	// składacz nakładki nie znajduje eksperta i rusza z samą osią.
	Agent string
}

// Okno komunikacji — byt pośredni między sesją a wiadomością.
type Okno struct {
	// Id okna.
	Id string
	// IdSesji wskazuje sesję nadrzędną.
	IdSesji string
	// Ustawienia wykonania okna.
	Ustawienia
	// Stan okna ze słownika kontraktu.
	Stan shared.WindowStatus
	// Utworzono i Zaktualizowano znakują cykl życia okna.
	Utworzono      time.Time
	Zaktualizowano time.Time
}

// Zmiana opisuje wybiórczą zmianę ustawień okna: pole niewypełnione zostaje
// bez zmiany. Odpowiada komendzie window.update kontraktu.
type Zmiana struct {
	Modul               *string
	KanalModelu         *string
	KatalogiRobocze     []string
	SrodowiskoWykonania *shared.ExecutionEnv
	TrybUprawnien       *shared.PermissionMode
	RolaOkna            *shared.WindowRole
	OknoKoordynatora    *string
	Tytul               *string
	// Agent wskazuje eksperta okna. Wskaźnik pusty zostawia wybór bez zmiany;
	// wskaźnik na pusty napis ZDEJMUJE eksperta i wraca do modelu surowego —
	// tak samo, jak opisuje to kontrakt („puste zdejmuje eksperta").
	Agent *string
}

// noweOkno składa okno otwarte z ustawieniami uzupełnionymi o wartości
// domyślne. Brak ustawienia znaczy wartość domyślną, nigdy odmowę.
func noweOkno(idSesji string, u Ustawienia) *Okno {
	teraz := time.Now().UTC()
	return &Okno{
		Id:             nowyIdentyfikator(przedrostekOkna),
		IdSesji:        idSesji,
		Ustawienia:     uzupelnijUstawienia(u),
		Stan:           shared.WindowStatusOpen,
		Utworzono:      teraz,
		Zaktualizowano: teraz,
	}
}

// uzupelnijUstawienia wstawia wartości domyślne i normalizuje rolę okna wraz
// z powiązaniem koordynatora.
//
// NORMALIZACJA OBEJMUJE WSZYSTKIE TRZY WYLICZENIA OKNA, NIE TYLKO ROLĘ. Wartość
// spoza słownika kontraktu (np. `permissionMode: "default"`) trafiłaby do
// rejestru nietknięta — rejestr żyje w pamięci i taką wartość zniesie, ale
// wiersz `okno_komunikacji` już nie: przekład na kolumnę (`dane/wyliczenia.go`)
// odmawia słowem „nie należy do słownika kontraktu", zapis okna pada, a dziennik
// rozmowy schodzi CAŁYM OKNEM na bufor pamięci. Wtedy `message.list` pokazuje
// rozmowę (bufor), a `history.load` i `history.delete` widzą pustkę (baza) —
// okno bez wiersza nie ma też historii ani czego przyciąć retencji. Wartość
// nieznana jest więc doprowadzana do domyślnej tak, jak robi to rola okna
// (`rola_okna.go`): bez bramy i bez odmowy, za to z oknem, które da się
// utrwalić. Słowniki pochodzą z kontraktu, nie z literałów tutaj.
func uzupelnijUstawienia(u Ustawienia) Ustawienia {
	if _, znane := shared.WartosciBazyExecutionEnv[u.SrodowiskoWykonania]; !znane {
		u.SrodowiskoWykonania = shared.ExecutionEnvLocal
	}
	if _, znany := shared.WartosciBazyPermissionMode[u.TrybUprawnien]; !znany {
		u.TrybUprawnien = shared.PermissionModeManual
	}
	u.KatalogiRobocze = append([]string(nil), u.KatalogiRobocze...)
	return normalizujRole(u)
}

// zastosuj nanosi wybiórczą zmianę na okno i odświeża znacznik zmiany.
func (o *Okno) zastosuj(z Zmiana) {
	u := o.Ustawienia
	przypiszNapis(&u.Modul, z.Modul)
	przypiszNapis(&u.KanalModelu, z.KanalModelu)
	przypiszNapis(&u.Tytul, z.Tytul)
	przypiszNapis(&u.OknoKoordynatora, z.OknoKoordynatora)
	// Wskaźnik na pusty napis zdejmuje eksperta — `przypiszNapis` przenosi
	// pustkę tak samo jak treść, więc powrót do modelu surowego jest zwykłym
	// zapisem, a nie osobną komendą.
	przypiszNapis(&u.Agent, z.Agent)
	if z.KatalogiRobocze != nil {
		u.KatalogiRobocze = z.KatalogiRobocze
	}
	if z.SrodowiskoWykonania != nil {
		u.SrodowiskoWykonania = *z.SrodowiskoWykonania
	}
	if z.TrybUprawnien != nil {
		u.TrybUprawnien = *z.TrybUprawnien
	}
	if z.RolaOkna != nil {
		u.RolaOkna = *z.RolaOkna
	}
	o.Ustawienia = uzupelnijUstawienia(u)
	o.Zaktualizowano = time.Now().UTC()
}

// przypiszNapis nanosi wartość tekstową, o ile zmiana ją niesie.
func przypiszNapis(cel *string, zrodlo *string) {
	if zrodlo != nil {
		*cel = *zrodlo
	}
}

// Kopia zwraca niezależny odpis okna — łącznie z listą katalogów roboczych.
func (o Okno) Kopia() Okno {
	odpis := o
	if o.KatalogiRobocze != nil {
		odpis.KatalogiRobocze = append([]string(nil), o.KatalogiRobocze...)
	}
	return odpis
}

// KatalogGlowny zwraca katalog uruchomienia procesu okna; pusty napis oznacza
// katalog odziedziczony po rdzeniu.
func (o Okno) KatalogGlowny() string {
	if len(o.KatalogiRobocze) == 0 {
		return ""
	}
	return o.KatalogiRobocze[0]
}

// CzyOtwarte mówi, czy okno przyjmuje pracę.
func (o Okno) CzyOtwarte() bool {
	return o.Stan == shared.WindowStatusOpen
}
