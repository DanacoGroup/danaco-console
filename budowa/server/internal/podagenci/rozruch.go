// Bramka rozruchu kolejkuje założenie pracy podagentów, żeby transakcje
// SQLite szeregowane w kolejce nie rywalizowały o blokadę i nie zawodziły.
package podagenci

import "context"

// MiejscRozruchu to liczba założeń pracy, które mogą iść jednocześnie: jedno,
// bo pisarz w SQLite jest jeden.
const MiejscRozruchu = 1

// Bramka wpuszcza do założenia pracy najwyżej MiejscRozruchu wołających naraz.
// Bramka pusta (wskaźnik nil) przepuszcza wszystkich — brak bramki nie
// zatrzymuje powołania.
type Bramka struct {
	miejsca chan struct{}
}

// NowaBramka zakłada bramkę o zadanej liczbie miejsc. Liczba poniżej jedynki
// znaczy jedno miejsce: bramka o zerze miejsc nie wpuściłaby nikogo nigdy,
// czyli byłaby zatrzymaniem platformy pod nazwą przepustowości.
func NowaBramka(miejsc int) *Bramka {
	if miejsc < 1 {
		miejsc = 1
	}
	return &Bramka{miejsca: make(chan struct{}, miejsc)}
}

// NowaBramkaRozruchu zakłada bramkę o liczbie miejsc rozruchu wskazanej stałą
// MiejscRozruchu tego pakietu.
func NowaBramkaRozruchu() *Bramka { return NowaBramka(MiejscRozruchu) }

// Wpusc czeka na wolne miejsce i oddaje funkcję zwalniającą je z powrotem;
// kontekst zerwany przerywa czekanie.
func (b *Bramka) Wpusc(ctx context.Context) (zwolnij func(), wpuszczony bool) {
	if b == nil || b.miejsca == nil {
		return func() {}, true
	}
	select {
	case b.miejsca <- struct{}{}:
		zwolniono := false
		return func() {
			if zwolniono {
				return
			}
			zwolniono = true
			<-b.miejsca
		}, true
	case <-ctx.Done():
		return func() {}, false
	}
}

// Zajete mówi, ile miejsc bramki jest w tej chwili zajętych — do meldunku
// i do testu, nie do rozstrzygania. Bramka pusta ma zajętych zero.
func (b *Bramka) Zajete() int {
	if b == nil || b.miejsca == nil {
		return 0
	}
	return len(b.miejsca)
}
