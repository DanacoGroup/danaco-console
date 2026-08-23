// Odpowiedzialność pliku: droga na żywo od zapisu nastawy do tego, kto z niej
// korzysta. Bez tej drogi Operator zmienia nastawę i nic się nie dzieje, dopóki
// czegoś nie przeładuje.
//
// Czym to nie jest.
// To nie jest drugi mechanizm nastaw ani druga tabela. Rozgłośnia nie
// przechowuje ani jednej wartości ustawienia z własnej woli: na każde ogłoszenie
// pyta ten sam rozstrzygacz o rozstrzygnięcie w kontekście nasłuchującego
// i podaje dalej to, co dostała. Jedynym stanem, jaki trzyma, jest zapis tego,
// co już powiedziała — po to, by nie budzić nasłuchującego zmianą, której nie
// było. Ledger „co powiedziano” nie jest źródłem wartości; jest pamięcią rozmowy.
//
// To nie jest własna usługa ani własny wątek.
// Rozgłośnia nie odpala żadnej gorutyny, nie odpytuje niczego w pętli i nie ma
// zegara. Doręczenie dzieje się w wątku tego, kto ogłosił zapis — czyli w torze
// komendy `config.set`. Rdzeń korzysta z zasobów urządzenia i nie stawia
// własnych usług.
//
// Kto ogłasza.
// Droga zapisu: po udanym utrwaleniu wiersza tabeli `ustawienie` woła `Oglos`
// z kluczem, który się zmienił. Ogłoszenie po zapisie, nie przed — nastawa,
// która nie usiadła w bazie, nie jest zmianą nastawy. Punkty wywołania leżą
// poza tym pakietem.
package konfig

import "sync"

// Zmiana to jedno doręczenie: klucz oraz jego rozstrzygnięcie w kontekście
// nasłuchującego w chwili doręczenia. Nasłuchujący nie dostaje surowego wiersza
// zapisu — dostaje wartość, która go obowiązuje po tym zapisie. Różnica jest
// istotna: zapis na poziomie szerszym może nie zmienić nic, jeżeli nasłuchujący
// ma wartość z poziomu węższego, i wtedy doręczenia nie będzie wcale.
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

// dotyczy odpowiada, czy nasłuch pytał o ten klucz.
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

// dodaj zapisuje nasłuch i zwraca jego numer.
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

// usun wykreśla nasłuch. Wykreślenie nieistniejącego nie jest błędem.
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

// Sledz zapisuje nasłuchującego na wskazane klucze i doręcza mu od razu stan
// bieżący, zanim jeszcze cokolwiek się zmieni.
//
// Pierwsze doręczenie jest częścią umowy, nie uprzejmością. Nasłuchujący, który
// nie ma się przeładowywać, musi skądś wziąć punkt wyjścia; gdyby brał go
// osobnym pytaniem, miałby dwie drogi do jednej wartości i wyścig między nimi
// (zapis mieszczący się pomiędzy pytaniem a zapisaniem się zginąłby). Jedna
// droga, jedno źródło.
//
// Zwrócona funkcja wykreśla nasłuch. Wołanie jej wielokrotnie jest bezpieczne.
// Brak odbiorcy albo brak kluczy daje funkcję pustą — nie odmowę.
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

// kluczeWykazem zamienia zbiór kluczy na wykaz do doręczenia.
func kluczeWykazem(zbior map[string]struct{}) []string {
	wykaz := make([]string, 0, len(zbior))
	for klucz := range zbior {
		wykaz = append(wykaz, klucz)
	}
	return wykaz
}

// Oglos zawiadamia nasłuchujących, że wskazane klucze mogły się zmienić.
// Woła się po udanym utrwaleniu zapisu.
//
// „mogły się zmienić", a nie „zmieniły się" — bo ogłaszający zna adres zapisu,
// a nie skutek dla każdego nasłuchującego z osobna. Skutek liczy rozstrzygacz,
// osobno dla kontekstu każdego nasłuchu, i tylko prawdziwa różnica idzie dalej.
// Dzięki temu zapis na poziomie okna nie budzi nasłuchu poziomu aplikacji,
// a zapis na poziomie aplikacji nie budzi nikogo, kto ma wartość z węższego.
func (r *Rozstrzygacz) Oglos(klucze ...string) {
	if r == nil || len(klucze) == 0 {
		return
	}
	for _, n := range r.rozglos.wykaz() {
		n.odswiez(r, klucze)
	}
}

// Nastawa jest uchwytem na żywo do jednego klucza: trzyma rozstrzygnięcie
// obowiązujące teraz i odświeża je samo, ogłoszeniem, bez pytania i bez
// przeładowania.
//
// Po co uchwyt, skoro jest `Rozstrzygnij`. Bo korzystający siedzi na drodze
// gorącej — straż bramki rozstrzyga przy każdym pakiecie z gniazda — a
// `Rozstrzygnij` schodzi po wartość do warstwy trwałości. Uchwyt zdejmuje ten
// koszt, nie zdejmując prawdy: wartość w nim nie starzeje się nigdy, bo każdy
// zapis ją nadpisuje w tej samej chwili, w której siada w bazie.
//
// I dlatego nie jest to druga prawda. Druga prawda to wartość, która
// może się rozjechać ze źródłem. Ta rozjechać się nie może — jedyną drogą jej
// zmiany jest ogłoszenie ze źródła, a własnego zapisu uchwyt nie przyjmuje.
type Nastawa struct {
	zamek   sync.RWMutex
	wynik   Wynik
	odwolaj func()
}

// NastawaNaZywo zapisuje uchwyt na nasłuch klucza i zwraca go już wypełniony
// wartością bieżącą.
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

// Wartosc zwraca samą wartość rozstrzygnięcia.
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
