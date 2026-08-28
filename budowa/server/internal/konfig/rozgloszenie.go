// Plik prowadzi drogę na żywo od zapisu nastawy do jej odbiorców. Rozgłośnia
// nie przechowuje wartości ustawień; przy ogłoszeniu pyta rozstrzygacz
// o wynik dla kontekstu nasłuchu. Ogłoszenie następuje po udanym zapisie
// w tabeli ustawienie.
package konfig

import "sync"

// Zmiana to jedno doręczenie: klucz oraz jego rozstrzygnięcie w kontekście
// nasłuchującego w chwili doręczenia. Zapis na poziomie szerszym może nie
// zmienić wartości nasłuchującego, jeżeli ma on wartość z poziomu węższego,
// i wtedy doręczenia nie będzie.
type Zmiana struct {
	Klucz string
	Wynik Wynik
}

// Odbiorca przyjmuje doręczenie. Wołany jest w wątku ogłaszającego (drogi
// zapisu), synchronicznie, więc ma wyłącznie wziąć wartość — nie wykonywać
// w nim pracy długiej i nie wołać zwrotnie rozgłośni tym samym nasłuchem.
type Odbiorca func(Zmiana)

// nasluch to jedna zapisana chęć bycia informowanym: czyj kontekst, które
// klucze, komu doręczać i co już powiedziano.
type nasluch struct {
	zamek    sync.Mutex
	kontekst Kontekst
	klucze   map[string]struct{}
	odbiorca Odbiorca
	ostatnie map[string]Wynik
}

// dotyczy odpowiada, czy nasłuch zarejestrował zainteresowanie danym
// kluczem wśród kluczy, na które czeka.
func (n *nasluch) dotyczy(klucz string) bool {
	_, jest := n.klucze[klucz]
	return jest
}

// odswiez rozstrzyga wskazane klucze w kontekście nasłuchu i doręcza te, które
// naprawdę się zmieniły wobec ostatniego doręczenia.
func (n *nasluch) odswiez(r *Rozstrzygacz, klucze []string) {
	n.zamek.Lock()
	defer n.zamek.Unlock()
	for _, klucz := range klucze {
		if !n.dotyczy(klucz) {
			continue
		}
		wynik := r.Rozstrzygnij(n.kontekst, klucz)
		if poprzedni, znany := n.ostatnie[klucz]; znany && poprzedni == wynik {
			continue
		}
		n.ostatnie[klucz] = wynik
		n.odbiorca(Zmiana{Klucz: klucz, Wynik: wynik})
	}
}

// rozglosnia trzyma wykaz nasłuchów. Wartość zerowa jest zdatna do pracy —
// wykaz powstaje przy pierwszym zapisaniu się nasłuchującego.
type rozglosnia struct {
	zamek    sync.Mutex
	kolejny  uint64
	nasluchy map[uint64]*nasluch
}

// dodaj zapisuje nasłuch na liście rozgłośni i zwraca jego numer, używany
// później do jego wykreślenia.
func (g *rozglosnia) dodaj(n *nasluch) uint64 {
	g.zamek.Lock()
	defer g.zamek.Unlock()
	if g.nasluchy == nil {
		g.nasluchy = make(map[uint64]*nasluch)
	}
	g.kolejny++
	g.nasluchy[g.kolejny] = n
	return g.kolejny
}

// usun wykreśla nasłuch o podanym numerze z listy rozgłośni; wykreślenie
// nasłuchu nieistniejącego nie jest błędem.
func (g *rozglosnia) usun(numer uint64) {
	g.zamek.Lock()
	defer g.zamek.Unlock()
	delete(g.nasluchy, numer)
}

// wykaz zwraca migawkę nasłuchów. Doręczanie dzieje się poza zamkiem wykazu,
// bo odbiorca może zapisać albo wykreślić nasłuch, a wtedy zamek trzymany przez
// doręczyciela zamknąłby rdzeń na sobie samym.
func (g *rozglosnia) wykaz() []*nasluch {
	g.zamek.Lock()
	defer g.zamek.Unlock()
	migawka := make([]*nasluch, 0, len(g.nasluchy))
	for _, n := range g.nasluchy {
		migawka = append(migawka, n)
	}
	return migawka
}

// Sledz zapisuje nasłuchującego na wskazane klucze i doręcza mu od razu
// stan bieżący, zanim cokolwiek się zmieni. Zwrócona funkcja wykreśla
// nasłuch i jest bezpieczna do wielokrotnego wołania; brak odbiorcy albo
// brak kluczy daje funkcję pustą.
func (r *Rozstrzygacz) Sledz(kontekst Kontekst, klucze []string, odbiorca Odbiorca) func() {
	if r == nil || odbiorca == nil {
		return func() {}
	}
	zbior := make(map[string]struct{}, len(klucze))
	for _, klucz := range klucze {
		if klucz != "" {
			zbior[klucz] = struct{}{}
		}
	}
	if len(zbior) == 0 {
		return func() {}
	}
	n := &nasluch{
		kontekst: kontekst,
		klucze:   zbior,
		odbiorca: odbiorca,
		ostatnie: make(map[string]Wynik, len(zbior)),
	}
	numer := r.rozglos.dodaj(n)
	n.odswiez(r, kluczeWykazem(zbior))

	var raz sync.Once
	return func() { raz.Do(func() { r.rozglos.usun(numer) }) }
}

// kluczeWykazem zamienia zbiór kluczy nasłuchu na wykaz kluczy przekazywany
// dalej przy doręczeniu zmiany nasłuchującemu.
func kluczeWykazem(zbior map[string]struct{}) []string {
	wykaz := make([]string, 0, len(zbior))
	for klucz := range zbior {
		wykaz = append(wykaz, klucz)
	}
	return wykaz
}

// Oglos zawiadamia nasłuchujących, że wskazane klucze mogły się zmienić;
// woła się po udanym utrwaleniu zapisu. Rozstrzygacz liczy skutek osobno dla
// kontekstu każdego nasłuchu, więc dalej idzie tylko prawdziwa różnica.
func (r *Rozstrzygacz) Oglos(klucze ...string) {
	if r == nil || len(klucze) == 0 {
		return
	}
	for _, n := range r.rozglos.wykaz() {
		n.odswiez(r, klucze)
	}
}

// Nastawa jest uchwytem na żywo do jednego klucza: trzyma rozstrzygnięcie
// obowiązujące teraz i odświeża je samo ogłoszeniem, bez pytania i
// przeładowania. Wartość nie starzeje się, bo każdy zapis nadpisuje ją
// w chwili zapisania w bazie.
type Nastawa struct {
	zamek   sync.RWMutex
	wynik   Wynik
	odwolaj func()
}

// NastawaNaZywo zapisuje uchwyt na nasłuch klucza w rozgłośni i zwraca go
// już wypełniony wartością bieżącą rozstrzygnięcia.
func (r *Rozstrzygacz) NastawaNaZywo(kontekst Kontekst, klucz string) *Nastawa {
	n := &Nastawa{}
	n.odwolaj = r.Sledz(kontekst, []string{klucz}, func(z Zmiana) {
		n.zamek.Lock()
		n.wynik = z.Wynik
		n.zamek.Unlock()
	})
	return n
}

// Wynik zwraca rozstrzygnięcie obowiązujące w tej chwili. Nie sięga do źródła —
// wartość jest tu dlatego, że źródło samo ją tu przyniosło.
func (n *Nastawa) Wynik() Wynik {
	if n == nil {
		return Wynik{}
	}
	n.zamek.RLock()
	defer n.zamek.RUnlock()
	return n.wynik
}

// Wartosc zwraca samą wartość rozstrzygnięcia trzymanego przez uchwyt,
// pomijając pozostałe pola struktury Wynik.
func (n *Nastawa) Wartosc() string {
	return n.Wynik().Wartosc
}

// Zamknij wykreśla nasłuch uchwytu. Uchwyt niezamknięty trzyma nasłuch do końca
// biegu procesu — dla nastaw poziomu `aplikacja` jest to stan docelowy, bo tyle
// właśnie żyje warstwa nasłuchu.
func (n *Nastawa) Zamknij() {
	if n == nil || n.odwolaj == nil {
		return
	}
	n.odwolaj()
}
