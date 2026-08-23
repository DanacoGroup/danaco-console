// Odpowiedzialność pliku: bramka rozruchu — kolejkowanie założenia pracy
// podagentów, żeby wszyscy powołani naprawdę zaczęli pracować.
//
// Powołanie kilkunastu podagentów jednym `subagent.spawn` puszcza tyleż
// goroutine naraz; każda zakłada kolejkę, dokłada pozycję, wiąże ją z wierszem
// i przestawia stan. Baza stoi w WAL z `busy_timeout`; zwykły zapis równoległy
// nie zawodzi, ale transakcja, która najpierw czyta, a potem pisze, zawodzi:
// podniesienie blokady odczytu do zapisu nie jest objęte `busy_timeout` i wraca
// natychmiast jako SQLITE_BUSY (5) albo BUSY_SNAPSHOT (517). Bramka szereguje
// te transakcje, więc nie rywalizują o blokadę i nie zawodzą.
//
// Bramka obejmuje wyłącznie założenie pracy — kilka zapisów trwających
// milisekundy. Samej pracy podagenta (tura modelu, sekundy albo minuty) nie
// obejmuje i obejmować nie może: podagenci mają pracować równolegle, a bramka
// rozciągnięta na turę zamieniłaby sieć kilkunastu w gęsiego idącą jedynkę.
// Wołający wchodzi w bramkę przed pierwszym zapisem i wychodzi z niej przed
// wywołaniem silnika.
package podagenci

import "context"

// MiejscRozruchu to liczba założeń pracy, które mogą iść jednocześnie.
//
// Jedno, nie więcej. Pisarz w SQLite jest jeden — drugie miejsce nie dokłada
// przepustowości, dokłada rywalizację, czyli dokładnie to, co ta bramka usuwa.
// Wartość jest stałą, a nie nastawą: nastawa bez pytania, które by ją
// rozstrzygało, byłaby pokrętłem bez skali.
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

// NowaBramkaRozruchu zakłada bramkę o liczbie miejsc rozruchu (MiejscRozruchu).
func NowaBramkaRozruchu() *Bramka { return NowaBramka(MiejscRozruchu) }

// Wpusc czeka na wolne miejsce i oddaje funkcję zwalniającą je z powrotem.
//
// Zwolnienie oddaje się zawsze — także po błędzie założenia pracy; wołający
// stawia `defer zwolnij()` zaraz po wejściu. Miejsce niezwrócone zabrałoby
// sieci przepustowość na stałe.
//
// Kontekst zerwany przerywa czekanie. Podagent odwołany w kolejce do bramki
// nie ma po co dostać miejsca — `wpuszczony` jest wtedy fałszem, a zwolnienie
// mimo to wolno wywołać (nic nie robi), żeby wołający nie musiał rozgałęziać
// `defer`.
//
// Bramka pusta wpuszcza natychmiast.
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
