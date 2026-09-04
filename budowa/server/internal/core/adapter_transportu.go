// Plik jest stykiem rdzenia z warstwą transportu: wnosi do kontekstu żądania tożsamość i ujście gniazda, przypisuje gniazdu konto rozpoznane przez bramkę, kieruje fragmenty strumienia do gniazda zamawiającego turę, zrywa gniazda po unieważnieniu sesji i oddaje serwer transportu jako nasłuch oraz nadajnik zdarzeń.
package core

import (
	"context"
	"strconv"
	"sync"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/transport"
	"danacoconsole/shared"
)

// wejscieTransportu podaje rdzeń warstwie transportu w kształcie, którego transport oczekuje, a rdzeniowi podaje transport jako nadajnik i nasłuch. Zależność idzie w jedną stronę: transport nie zna rdzenia.
type wejscieTransportu struct {
	rdzen *Rdzen
	tor   *torStrumieni
}

// Obsluz wypełnia interfejs transport.Rdzen. Odpowiedź wraca wynikiem, zdarzenia idą rozgłoszeniem do urządzeń konta, a strumień tury do gniazda, które ją zamówiło.
func (w wejscieTransportu) Obsluz(kontekst context.Context, zadanie protocol.Request, ujscie transport.Ujscie) protocol.Koperta {
	if ujscie != nil {
		kontekst = zUjsciem(zPolaczeniem(kontekst, ujscie.Tozsamosc()), ujscie)
	}
	return w.rdzen.WykonajZadanie(kontekst, zadanie)
}

// Przylaczono wypełnia nieobowiązkowe rozszerzenie transport.ObserwatorPolaczen. Konto przedstawione przy nawiązaniu jest słowem klienta, a nie ustaleniem rdzenia, więc świeże gniazdo wraca na konto domyślne: przynależność do konta nadaje wyłącznie bramka — `connection.hello` z ważnym tokenem albo udane wejście.
func (w wejscieTransportu) Przylaczono(ujscie transport.Ujscie) {
	if ujscie == nil {
		return
	}
	ujscie.PrzypiszKonto(transport.KontoDomyslne)
}

// Odlaczono zdejmuje więź rozłączonego urządzenia i jego wpisy w torze strumieni. Bez tego obie mapy rosłyby o wpis na każde nawiązanie przez całe życie procesu.
func (w wejscieTransportu) Odlaczono(ujscie transport.Ujscie) {
	if ujscie == nil {
		return
	}
	w.tor.odlacz(ujscie.Id())
	if w.rdzen != nil {
		w.rdzen.wiez.Rozwiaz(ujscie.Id())
	}
}

// kluczUjscia niesie ujście gniazda, z którego przyszło żądanie. Wpis stoi obok kluczPolaczenia, a nie zamiast niego: tożsamość jest kopią faktów o wołającym, ujście jest żywym kanałem, któremu rdzeń przypisuje konto.
const kluczUjscia kluczKontekstu = "danaco:ujscie"

// zUjsciem dokłada do kontekstu ujście gniazda. Wpina to `wejscieTransportu.Obsluz` — jedyne miejsce, przez które przechodzi każde żądanie z gniazda.
func zUjsciem(ctx context.Context, ujscie transport.Ujscie) context.Context {
	if ctx == nil || ujscie == nil {
		return ctx
	}
	return context.WithValue(ctx, kluczUjscia, ujscie)
}

// ujscieZKontekstu oddaje ujście wołającego. Zero znaczy żądanie spoza gniazda — tak wygląda bieg wewnętrzny rdzenia i wywołanie ze sprawdzianu; to nie jest błąd.
func ujscieZKontekstu(ctx context.Context) transport.Ujscie {
	if ctx == nil {
		return nil
	}
	ujscie, jest := ctx.Value(kluczUjscia).(transport.Ujscie)
	if !jest {
		return nil
	}
	return ujscie
}

// przypiszKontoGniazda wiąże gniazdo z kontem, któremu wydano sesję bramki. Transport adresuje rozgłoszenia kontem połączenia, a wie o nim tyle, ile powie rdzeń — bez tego wywołania konto gniazda zostaje takie, jakim przedstawił je klient. Skrót pusty odsyła gniazdo na konto domyślne: połączenie bez sesji nie należy do żadnego konta.
// żadnego konta.
func przypiszKontoGniazda(ctx context.Context, konta RozpoznanieKontaSesji, skrot string) {
	ujscie := ujscieZKontekstu(ctx)
	if ujscie == nil {
		return
	}
	ujscie.PrzypiszKonto(nazwaKontaGniazda(ctx, konta, skrot))
}

// nazwaKontaGniazda przekłada sesję bramki na nazwę konta, którą posługuje się transport. Sesja nierozpoznana daje konto domyślne, nie napis pusty: pustego konta PrzypiszKonto nie przyjmuje, a gniazdo musi mieć jak zejść z konta przypisanego wcześniej.
func nazwaKontaGniazda(ctx context.Context, konta RozpoznanieKontaSesji, skrot string) string {
	if konta == nil || skrot == "" {
		return transport.KontoDomyslne
	}
	kontoId, err := konta.KontoSesjiBramki(ctx, skrot)
	if err != nil {
		return transport.KontoDomyslne
	}
	return nazwaKontaTransportu(kontoId)
}

// nazwaKontaTransportu przekłada identyfikator konta z bazy na nazwę konta w transporcie — innej nazwy konta transport nie zna. Zero, czyli konto nierozpoznane, daje konto domyślne.
func nazwaKontaTransportu(kontoId int64) string {
	if kontoId == 0 {
		return transport.KontoDomyslne
	}
	return strconv.FormatInt(kontoId, 10)
}

// kontoAdresata nazywa konto, do którego transport adresuje zdarzenie zamówione w tym kontekście: konto rozpoznane z sesji bramki; konto domyślne dla gniazda bez rozpoznanej sesji; konto puste — wszystkie połączenia — dla czynności własnej rdzenia, która nie przyszła z żadnego gniazda.
func kontoAdresata(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if kontoId := dane.KontoOperatora(ctx); kontoId != 0 {
		return nazwaKontaTransportu(kontoId)
	}
	if konto, jest := ctx.Value(kluczKontaZadania).(string); jest && konto != "" {
		return konto
	}
	if polaczenieZKontekstu(ctx) != "" {
		return transport.KontoDomyslne
	}
	return ""
}

// kluczKontaZadania niesie w kontekście pracy w tle nazwę konta adresata zdarzeń. Konto domyślne gniazda bez sesji bramki nie ma identyfikatora w bazie, więc samo wskazanie Operatora by go nie przeniosło.
const kluczKontaZadania kluczKontekstu = "danaco:konto-zadania"

// zKontemZadania przenosi konto wołającego z kontekstu żądania do kontekstu pracy puszczonej w tle. Praca przeżywa żądanie, a jej zdarzenia i zapisy mają należeć do konta, które ją zamówiło.
func zKontemZadania(praca, zadanie context.Context) context.Context {
	if praca == nil || zadanie == nil {
		return praca
	}
	if konto := kontoAdresata(zadanie); konto != "" {
		praca = context.WithValue(praca, kluczKontaZadania, konto)
	}
	kontoId := dane.KontoOperatora(zadanie)
	if kontoId == 0 {
		return praca
	}
	return dane.ZKontemOperatora(praca, kontoId)
}

// torStrumieni pamięta, które gniazdo zamówiło turę: koperta fragmentu powtarza identyfikator żądania, więc identyfikator wystarcza za adres. Wpis żyje od zamówienia do fragmentu ostatniego; gniazdo zerwane w trakcie zostawia wpis z samym kontem, żeby reszta strumienia poszła do urządzeń tego konta.
type torStrumieni struct {
	zamek sync.Mutex
	wpisy map[string]wpisToru
}

// wpisToru niesie ujście zamawiające i konto, na którym stało przy zamówieniu.
type wpisToru struct {
	ujscie transport.Ujscie
	konto  string
}

// nowyTorStrumieni zakłada pusty tor.
func nowyTorStrumieni() *torStrumieni {
	return &torStrumieni{wpisy: make(map[string]wpisToru)}
}

// zwiaz przypisuje strumień żądania gniazdu zamawiającemu. Żądanie spoza gniazda nie zostawia wpisu.
func (t *torStrumieni) zwiaz(idZadania string, ujscie transport.Ujscie) {
	if t == nil || idZadania == "" || ujscie == nil {
		return
	}
	t.zamek.Lock()
	t.wpisy[idZadania] = wpisToru{ujscie: ujscie, konto: ujscie.Konto()}
	t.zamek.Unlock()
}

// odlacz zdejmuje ujście z wpisów rozłączonego gniazda, zostawiając konto.
func (t *torStrumieni) odlacz(idPolaczenia string) {
	if t == nil || idPolaczenia == "" {
		return
	}
	t.zamek.Lock()
	defer t.zamek.Unlock()
	for id, wpis := range t.wpisy {
		if wpis.ujscie != nil && wpis.ujscie.Id() == idPolaczenia {
			t.wpisy[id] = wpisToru{konto: wpis.konto}
		}
	}
}

// dostarcz kieruje fragment do gniazda, które zamówiło turę; fałsz znaczy strumień bez wpisu, który idzie zwykłym rozgłoszeniem. Wysyłka nieudana zdejmuje ujście z wpisu, a fragment idzie do urządzeń konta zamawiającego.
func (t *torStrumieni) dostarcz(serwer *transport.Serwer, k protocol.Koperta) bool {
	if t == nil {
		return false
	}
	t.zamek.Lock()
	wpis, jest := t.wpisy[k.Id]
	if jest && protocol.Ostatni(k) {
		delete(t.wpisy, k.Id)
	}
	t.zamek.Unlock()
	if !jest {
		return false
	}
	if wpis.ujscie != nil {
		if err := wpis.ujscie.Wyslij(k); err == nil {
			return true
		}
		t.odlacz(wpis.ujscie.Id())
	}
	serwer.Rozglos(wpis.konto, k)
	return true
}

// rozlaczanieSesji zrywa gniazda konta, których sesja bramki przestała nadawać. Więź i rozpoznanie ważności wchodzą po złożeniu rdzenia, bo nasłuch powstaje przed nim.
type rozlaczanieSesji struct {
	serwer  *transport.Serwer
	wiez    *wiezBramki
	waznosc RozpoznanieWaznosciSesji
}

// uzupelnij wpina więź gniazd i rozpoznanie ważności sesji złożone razem z rdzeniem.
func (r *rozlaczanieSesji) uzupelnij(wiez *wiezBramki, waznosc RozpoznanieWaznosciSesji) {
	if r == nil {
		return
	}
	r.wiez, r.waznosc = wiez, waznosc
}

// RozlaczPoUniewaznieniu wypełnia port RozlaczanieSesji. O sesji rozstrzyga wiersz w bazie, nie więź: unieważnienie zaszło gdzie indziej i więź o nim nie wie. Ważność nierozstrzygnięta zostawia gniazdo bramce.
func (r *rozlaczanieSesji) RozlaczPoUniewaznieniu(ctx context.Context, konto, powod string) int {
	if r == nil || r.serwer == nil || r.wiez == nil || r.waznosc == nil {
		return 0
	}
	wlasne := polaczenieZKontekstu(ctx)
	return r.serwer.RozlaczKonta(konto, powod, func(id string) bool {
		if id == wlasne {
			return false
		}
		skrot := r.wiez.Skrot(id)
		if skrot == "" {
			return false
		}
		_, wazna, err := r.waznosc.WaznoscSesjiBramki(ctx, skrot)
		if err != nil {
			return false
		}
		return !wazna
	})
}

// zRozlaczaniem przyjmuje zrywanie gniazd dołożone po montażu, bo port składa się piętro niżej niż transport — wzorem zWiezia.
type zRozlaczaniem interface {
	przyjmijRozlaczanie(r RozlaczanieSesji)
}

// nasluchTransportu podaje serwer transportu jako nasłuch rdzenia, nadajnik zdarzeń rozgłaszanych do urządzeń konta i zrywanie gniazd po unieważnieniu sesji. Jest wartością, więc tor i zrywanie stoją za wskaźnikami wspólnymi dla wszystkich kopii.
type nasluchTransportu struct {
	serwer *transport.Serwer
	tor    *torStrumieni
	sesje  *rozlaczanieSesji
}

// nowyNasluchTransportu składa nasłuch nad serwerem transportu z pustym torem strumieni i zrywaniem czekającym na więź rdzenia.
func nowyNasluchTransportu(serwer *transport.Serwer) nasluchTransportu {
	return nasluchTransportu{
		serwer: serwer,
		tor:    nowyTorStrumieni(),
		sesje:  &rozlaczanieSesji{serwer: serwer},
	}
}

// Sluchaj otwiera nasłuch i pracuje do zamknięcia kontekstu, po czym zamyka serwer. Rozłączenie klienta nie kończy pracy rdzenia.
func (n nasluchTransportu) Sluchaj(kontekst context.Context) error {
	if err := n.serwer.Uruchom(kontekst); err != nil {
		return err
	}
	<-kontekst.Done()
	return n.serwer.Zamknij()
}

// Rozglos wypełnia port Nadajnik rdzenia. Fragment strumienia z wpisem w torze idzie do gniazda zamawiającego; pozostałe komunikaty do połączeń konta, a przy koncie pustym do wszystkich połączeń rdzenia.
func (n nasluchTransportu) Rozglos(konto string, k protocol.Koperta) {
	if k.Type == shared.EventStreamChunk && n.tor.dostarcz(n.serwer, k) {
		return
	}
	if k.Type == shared.EventProgressChanged {
		n.serwer.RozglosPoBramce(konto, k)
		return
	}
	n.serwer.Rozglos(konto, k)
}

// RozlaczPoUniewaznieniu wypełnia port RozlaczanieSesji nasłuchem.
func (n nasluchTransportu) RozlaczPoUniewaznieniu(ctx context.Context, konto, powod string) int {
	return n.sesje.RozlaczPoUniewaznieniu(ctx, konto, powod)
}
