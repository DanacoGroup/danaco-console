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
	// KatalogiRobocze to lista; pierwszy katalog jest katalogiem uruchomienia procesu okna.
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
	// Agent to kod eksperta nałożonego na kanał modelu okna; wartość pusta oznacza model surowy.
	Agent string
}

// Okno komunikacji jest bytem pośrednim między sesją a wiadomością i przechowuje ustawienia
// wykonania, stan bieżący oraz znaczniki czasu utworzenia i ostatniej aktualizacji.
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
	// Agent wskazuje eksperta okna; wskaźnik pusty zostawia wybór, pusty napis zdejmuje eksperta.
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

// Funkcja uzupelnijUstawienia wstawia wartości domyślne i normalizuje rolę okna wraz z powiązaniem koordynatora.
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

// Metoda zastosuj nanosi na okno wybiórczą zmianę przekazaną w strukturze Zmiana i odświeża znacznik czasu ostatniej aktualizacji okna.
func (o *Okno) zastosuj(z Zmiana) {
	u := o.Ustawienia
	przypiszNapis(&u.Modul, z.Modul)
	przypiszNapis(&u.KanalModelu, z.KanalModelu)
	przypiszNapis(&u.Tytul, z.Tytul)
	przypiszNapis(&u.OknoKoordynatora, z.OknoKoordynatora)
	// Wskaźnik na pusty napis zdejmuje eksperta; przypiszNapis przenosi pustkę tak samo jak treść.
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

// Funkcja przypiszNapis nanosi na pole docelowe wartość tekstową ze wskaźnika źródłowego wyłącznie wtedy, gdy zmiana rzeczywiście ją niesie.
func przypiszNapis(cel *string, zrodlo *string) {
	if zrodlo != nil {
		*cel = *zrodlo
	}
}

// Metoda Kopia zwraca niezależny odpis okna, obejmujący również osobną kopię listy katalogów roboczych, aby zmiana odpisu nie naruszała oryginału.
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

// Metoda CzyOtwarte zwraca wartość logiczną informującą, czy okno znajduje się w stanie pozwalającym na przyjmowanie kolejnej pracy.
func (o Okno) CzyOtwarte() bool {
	return o.Stan == shared.WindowStatusOpen
}
