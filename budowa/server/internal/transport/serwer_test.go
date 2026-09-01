package transport

import (
	"context"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/coder/websocket"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// Sprawdziany warstwy nasłuchu idą przez prawdziwe gniazdo, nie przez atrapę biblioteki.

// rdzenAtrapa jest realizacją interfejsu Rdzen po stronie sprawdzianu. Transport
// nie zna rdzenia — zna wyłącznie ten interfejs — więc atrapa jest tu bytem
// pełnoprawnym, a nie namiastką.
type rdzenAtrapa struct {
	// obsluga rozstrzyga odpowiedź. Zerowa oddaje potwierdzenie bez treści.
	obsluga func(ctx context.Context, z protocol.Request, u Ujscie) protocol.Koperta

	zamek       sync.Mutex
	przylaczone []string
	odlaczone   []string
}

func (r *rdzenAtrapa) Obsluz(ctx context.Context, z protocol.Request, u Ujscie) protocol.Koperta {
	if r.obsluga != nil {
		return r.obsluga(ctx, z, u)
	}
	return protocol.KopertaOdpowiedzi(z.Koperta(), protocol.Odpowiedz{Status: shared.EnvelopeStatusOk})
}

// kontaDoPrzypisania niesie konto, które polacz zamówił dla najbliższego gniazda;
// kontaPrzypisane potwierdza przypisanie. Konto nadaje rdzeń po bramce, nie
// parametr zapytania, więc sprawdzian nadaje je tą samą drogą — z Przylaczono.
var (
	kontaDoPrzypisania = make(chan string, 8)
	kontaPrzypisane    = make(chan struct{}, 8)
)

func (r *rdzenAtrapa) Przylaczono(u Ujscie) {
	r.zamek.Lock()
	r.przylaczone = append(r.przylaczone, u.Id())
	r.zamek.Unlock()
	select {
	case konto := <-kontaDoPrzypisania:
		u.PrzypiszKonto(konto)
		kontaPrzypisane <- struct{}{}
	default:
	}
}

func (r *rdzenAtrapa) Odlaczono(u Ujscie) {
	r.zamek.Lock()
	r.odlaczone = append(r.odlaczone, u.Id())
	r.zamek.Unlock()
}

func (r *rdzenAtrapa) zdarzeniaPolaczen() (int, int) {
	r.zamek.Lock()
	defer r.zamek.Unlock()
	return len(r.przylaczone), len(r.odlaczone)
}

// podnies uruchamia serwer sprawdzianu na porcie wskazanym przez system i zwraca go wraz z adresem gniazda nasłuchu.
func podnies(t *testing.T, rdzen Rdzen) (*Serwer, string) {
	t.Helper()

	ustawienia := Domyslne()
	ustawienia.Adres = "127.0.0.1"
	ustawienia.PortDowolny = true
	ustawienia.Dziennik = log.New(io.Discard, "", 0)

	serwer := Nowy(ustawienia)
	if rdzen != nil {
		serwer.PodlaczRdzen(rdzen)
	}

	zycie, zakoncz := context.WithCancel(context.Background())
	t.Cleanup(zakoncz)
	if err := serwer.Uruchom(zycie); err != nil {
		t.Fatalf("nie można uruchomić nasłuchu: %v", err)
	}
	t.Cleanup(func() { _ = serwer.Zamknij() })

	return serwer, "ws://" + serwer.Adres() + SciezkaGniazdaDomyslna
}

// polacz nawiązuje gniazdo klienta pod wskazanym adresem. Konto niepuste
// przypisuje atrapa rdzenia w Przylaczono; puste zostawia konto domyślne serwera.
// Origin jest pochodzeniem powłoki bez portu, bo wzorce własne nie mają portu dowolnego.
func polacz(t *testing.T, adres, konto string) *websocket.Conn {
	t.Helper()

	if konto != "" {
		kontaDoPrzypisania <- konto
	}
	ctx, przerwij := context.WithTimeout(context.Background(), 5*time.Second)
	defer przerwij()

	gniazdo, _, err := websocket.Dial(ctx, adres, &websocket.DialOptions{
		HTTPHeader: http.Header{"Origin": []string{"http://localhost"}},
	})
	if err != nil {
		t.Fatalf("nie można nawiązać gniazda: %v", err)
	}
	t.Cleanup(func() { _ = gniazdo.CloseNow() })
	if konto != "" {
		select {
		case <-kontaPrzypisane:
		case <-ctx.Done():
			t.Fatalf("rdzeń nie przypisał konta %q gniazdu", konto)
		}
	}
	return gniazdo
}

// wyslij koduje kopertę kontraktu i podaje ją do gniazda klienta w ramach jednego wywołania sprawdzianu.
func wyslij(t *testing.T, gniazdo *websocket.Conn, k protocol.Koperta) {
	t.Helper()
	dane, err := protocol.Zakoduj(k)
	if err != nil {
		t.Fatalf("nie można zakodować koperty: %v", err)
	}
	ctx, przerwij := context.WithTimeout(context.Background(), 5*time.Second)
	defer przerwij()
	if err := gniazdo.Write(ctx, websocket.MessageText, dane); err != nil {
		t.Fatalf("nie można wysłać koperty: %v", err)
	}
}

// odbierz czyta jedną kopertę z gniazda klienta i dekoduje ją do postaci kontraktu przed jej zwrotem sprawdzianowi.
func odbierz(t *testing.T, gniazdo *websocket.Conn) protocol.Koperta {
	t.Helper()
	ctx, przerwij := context.WithTimeout(context.Background(), 5*time.Second)
	defer przerwij()
	_, dane, err := gniazdo.Read(ctx)
	if err != nil {
		t.Fatalf("nie można odczytać koperty: %v", err)
	}
	koperta, err := protocol.Odkoduj(dane)
	if err != nil {
		t.Fatalf("odpowiedź nie jest kopertą kontraktu: %v", err)
	}
	return koperta
}

// TestOdpowiedzWracaZTymSamymIdentyfikatorem sprawdza obieg zamknięty przez
// prawdziwe gniazdo. Klient koreluje odpowiedź po identyfikatorze, więc rozjazd
// tego pola zostawia okno w wiecznym ładowaniu, choć odpowiedź przyszła.
func TestOdpowiedzWracaZTymSamymIdentyfikatorem(t *testing.T) {
	_, adres := podnies(t, &rdzenAtrapa{})
	gniazdo := polacz(t, adres, "")

	wyslij(t, gniazdo, protocol.Koperta{Type: shared.CommandSessionList, Id: "zadanie-pierwsze"})
	odpowiedz := odbierz(t, gniazdo)

	if odpowiedz.Id != "zadanie-pierwsze" {
		t.Errorf("odpowiedź wraca z identyfikatorem %q", odpowiedz.Id)
	}
	if odpowiedz.Type != shared.CommandSessionList {
		t.Errorf("odpowiedź wraca pod typem %q", odpowiedz.Type)
	}
	if odpowiedz.Status == nil || *odpowiedz.Status != shared.EnvelopeStatusOk {
		t.Errorf("odpowiedź niesie stan %+v", odpowiedz.Status)
	}
}

// TestStrumienIdzieJednymIdentyfikatoremIWKolejnosci przechodzi turę strumieniową:
// rdzeń odsyła fragmenty przez ujście, a odpowiedź na komendę nie idzie wcale.
// Sprawdza jeden identyfikator, rosnące numery i jedno domknięcie.
func TestStrumienIdzieJednymIdentyfikatoremIWKolejnosci(t *testing.T) {
	const fragmentow = 5

	rdzen := &rdzenAtrapa{
		obsluga: func(_ context.Context, z protocol.Request, u Ujscie) protocol.Koperta {
			for i := 1; i <= fragmentow; i++ {
				koperta, err := protocol.KopertaFragmentu(z.Id, z.Zasieg.Sesja, i, i == fragmentow,
					protocol.ChunkTekstu("okno-pierwsze", "wiadomosc-pierwsza", "porcja"))
				if err != nil {
					return protocol.KopertaBledu(z.Koperta(),
						protocol.NowyBlad(shared.ErrorCodeInternalError, err.Error()))
				}
				if err := u.Wyslij(koperta); err != nil {
					return protocol.KopertaBledu(z.Koperta(),
						protocol.NowyBlad(shared.ErrorCodeInternalError, err.Error()))
				}
			}
			// Pusty typ znaczy „odpowiedź poszła osobno" — transport nic wtedy
			// nie odsyła.
			return protocol.Koperta{}
		},
	}

	_, adres := podnies(t, rdzen)
	gniazdo := polacz(t, adres, "")

	wyslij(t, gniazdo, protocol.Koperta{Type: shared.CommandMessageSend, Id: "tura-pierwsza"})

	poprzedni := 0
	domkniec := 0
	for i := 0; i < fragmentow; i++ {
		koperta := odbierz(t, gniazdo)
		if koperta.Type != shared.EventStreamChunk {
			t.Fatalf("fragment %d przyszedł pod typem %q", i+1, koperta.Type)
		}
		if koperta.Id != "tura-pierwsza" {
			t.Errorf("fragment %d niesie identyfikator %q zamiast identyfikatora żądania",
				i+1, koperta.Id)
		}
		if numer := protocol.Numer(koperta); numer != poprzedni+1 {
			t.Errorf("fragment przyszedł z numerem %d, oczekiwany %d", numer, poprzedni+1)
		} else {
			poprzedni = numer
		}
		if protocol.Ostatni(koperta) {
			domkniec++
		}
	}
	if domkniec != 1 {
		t.Errorf("strumień domknięty %d razy, oczekiwane raz", domkniec)
	}
}

// TestRozgloszenieNieWychodziPozaKonto sprawdza wyciek między kontami. Rozgłoszenie
// jest jedyną drogą synchronizacji wielourządzeniowej, więc pomyłka w doborze
// odbiorców oddaje treść jednego konta urządzeniom drugiego.
func TestRozgloszenieNieWychodziPozaKonto(t *testing.T) {
	serwer, adres := podnies(t, &rdzenAtrapa{})

	pierwszeKonta := polacz(t, adres, "konto-pierwsze")
	drugieKonta := polacz(t, adres, "konto-pierwsze")
	obce := polacz(t, adres, "konto-drugie")

	poczekajNaPolaczenia(t, serwer, 3)

	zdarzenie, err := protocol.NowaKoperta(shared.EventSessionChanged, "rozgloszenie", "sesja-pierwsza", nil)
	if err != nil {
		t.Fatalf("nie można złożyć zdarzenia: %v", err)
	}
	if odbiorcow := serwer.Rozglos("konto-pierwsze", zdarzenie); odbiorcow != 2 {
		t.Errorf("rozgłoszenie przyjęło %d urządzeń, oczekiwane 2", odbiorcow)
	}

	for nazwa, gniazdo := range map[string]*websocket.Conn{
		"pierwsze urządzenie konta": pierwszeKonta,
		"drugie urządzenie konta":   drugieKonta,
	} {
		koperta := odbierz(t, gniazdo)
		if koperta.Type != shared.EventSessionChanged {
			t.Errorf("%s dostało %q zamiast zdarzenia", nazwa, koperta.Type)
		}
	}

	// Urządzenie obcego konta nie ma dostać niczego; krótki termin odczytu jest tu miarą ciszy.
	ctx, przerwij := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer przerwij()
	if _, dane, err := obce.Read(ctx); err == nil {
		t.Errorf("zdarzenie konta pierwszego wyciekło do konta drugiego: %s", dane)
	}
}

// TestRozgloszenieBezKontaIdzieDoWszystkich sprawdza drugą stronę tej samej
// reguły: puste konto znaczy wszystkie połączenia rdzenia.
func TestRozgloszenieBezKontaIdzieDoWszystkich(t *testing.T) {
	serwer, adres := podnies(t, &rdzenAtrapa{})
	polacz(t, adres, "konto-pierwsze")
	polacz(t, adres, "konto-drugie")
	poczekajNaPolaczenia(t, serwer, 2)

	zdarzenie, err := protocol.NowaKoperta(shared.EventConfigChanged, "rozgloszenie", "", nil)
	if err != nil {
		t.Fatalf("nie można złożyć zdarzenia: %v", err)
	}
	if odbiorcow := serwer.Rozglos("", zdarzenie); odbiorcow != 2 {
		t.Errorf("rozgłoszenie bez konta przyjęło %d urządzeń, oczekiwane 2", odbiorcow)
	}
}

// TestRozlaczenieWykreslaPolaczenieZRejestru pilnuje sieroty w rejestrze
// połączeń. Sierota nie wywraca niczego od razu — rośnie po cichu i zabiera
// rozgłoszenia adresowane do konta, którego już nikt nie słucha.
func TestRozlaczenieWykreslaPolaczenieZRejestru(t *testing.T) {
	rdzen := &rdzenAtrapa{}
	serwer, adres := podnies(t, rdzen)

	gniazdo := polacz(t, adres, "konto-pierwsze")
	poczekajNaPolaczenia(t, serwer, 1)

	if err := gniazdo.Close(websocket.StatusNormalClosure, "koniec sprawdzianu"); err != nil {
		t.Fatalf("nie można zamknąć gniazda: %v", err)
	}
	poczekajNaPolaczenia(t, serwer, 0)

	przylaczonych, odlaczonych := rdzen.zdarzeniaPolaczen()
	if przylaczonych != 1 || odlaczonych != 1 {
		t.Errorf("rdzeń zawiadomiony o %d przyłączeniach i %d odłączeniach, oczekiwane po jednym",
			przylaczonych, odlaczonych)
	}
}

// TestRozlaczenieWTrakcieStrumieniaNieZostawiaSieroty jest tym samym
// sprawdzianem w warunkach, w których wyciek jest najbardziej prawdopodobny:
// urządzenie odchodzi, gdy rdzeń jeszcze pisze.
func TestRozlaczenieWTrakcieStrumieniaNieZostawiaSieroty(t *testing.T) {
	ruszyl := make(chan struct{})
	skonczyl := make(chan struct{})

	rdzen := &rdzenAtrapa{
		obsluga: func(_ context.Context, z protocol.Request, u Ujscie) protocol.Koperta {
			close(ruszyl)
			defer close(skonczyl)
			for i := 1; i <= 200; i++ {
				koperta, err := protocol.KopertaFragmentu(z.Id, "", i, false,
					protocol.ChunkTekstu("okno-pierwsze", "wiadomosc-pierwsza", "porcja"))
				if err != nil {
					break
				}
				// Błąd wysyłki do urządzenia, które odeszło, jest błędem jednego wywołania, a nie awarią rdzenia.
				_ = u.Wyslij(koperta)
			}
			return protocol.Koperta{}
		},
	}

	serwer, adres := podnies(t, rdzen)
	gniazdo := polacz(t, adres, "konto-pierwsze")
	poczekajNaPolaczenia(t, serwer, 1)

	wyslij(t, gniazdo, protocol.Koperta{Type: shared.CommandMessageSend, Id: "tura-przerwana"})
	<-ruszyl
	if err := gniazdo.CloseNow(); err != nil {
		t.Fatalf("nie można zerwać gniazda: %v", err)
	}

	<-skonczyl
	poczekajNaPolaczenia(t, serwer, 0)
}

// TestRamkaNieczytelnaNieZrywaPolaczenia sprawdza obietnicę z opisu: komunikat
// niepoprawny strukturalnie wraca odmową z kodem kontraktu, a gniazdo zostaje
// otwarte. Sprawdzian dowodzi tego kolejną komendą po odmowie.
func TestRamkaNieczytelnaNieZrywaPolaczenia(t *testing.T) {
	_, adres := podnies(t, &rdzenAtrapa{})
	gniazdo := polacz(t, adres, "")

	ctx, przerwij := context.WithTimeout(context.Background(), 5*time.Second)
	defer przerwij()
	if err := gniazdo.Write(ctx, websocket.MessageText, []byte("to nie jest koperta")); err != nil {
		t.Fatalf("nie można wysłać śmieci: %v", err)
	}

	odmowa := odbierz(t, gniazdo)
	if odmowa.Status == nil || *odmowa.Status != shared.EnvelopeStatusError {
		t.Error("komunikat niepoprawny nie dostał stanu błędu")
	}
	if odmowa.Error == nil || odmowa.Error.Code != shared.ErrorCodeValidationFailed {
		t.Errorf("komunikat niepoprawny dostał %+v, oczekiwany kod %q",
			odmowa.Error, shared.ErrorCodeValidationFailed)
	}

	wyslij(t, gniazdo, protocol.Koperta{Type: shared.CommandSessionList, Id: "po-smieciach"})
	if dalsza := odbierz(t, gniazdo); dalsza.Id != "po-smieciach" {
		t.Errorf("po odmowie połączenie nie obsłużyło kolejnej komendy: %+v", dalsza)
	}
}

// TestZalamanieObslugiWracaKodemZamiastZabicProcesu sprawdza zachowanie fail-open:
// usterka obsługiwacza kończy jedno wywołanie, a nie proces rdzenia.
func TestZalamanieObslugiWracaKodemZamiastZabicProcesu(t *testing.T) {
	rdzen := &rdzenAtrapa{
		obsluga: func(_ context.Context, z protocol.Request, _ Ujscie) protocol.Koperta {
			if z.Id == "zadanie-wywracajace" {
				panic("obsługiwacz sprawdzianu")
			}
			return protocol.KopertaOdpowiedzi(z.Koperta(), protocol.Odpowiedz{Status: shared.EnvelopeStatusOk})
		},
	}
	_, adres := podnies(t, rdzen)
	gniazdo := polacz(t, adres, "")

	wyslij(t, gniazdo, protocol.Koperta{Type: shared.CommandSessionList, Id: "zadanie-wywracajace"})
	odpowiedz := odbierz(t, gniazdo)

	if odpowiedz.Error == nil || odpowiedz.Error.Code != shared.ErrorCodeInternalError {
		t.Errorf("załamanie obsługi dało %+v, oczekiwany kod %q",
			odpowiedz.Error, shared.ErrorCodeInternalError)
	}

	wyslij(t, gniazdo, protocol.Koperta{Type: shared.CommandSessionList, Id: "po-zalamaniu"})
	if dalsza := odbierz(t, gniazdo); dalsza.Id != "po-zalamaniu" {
		t.Errorf("po załamaniu obsługi połączenie nie obsłużyło kolejnej komendy: %+v", dalsza)
	}
}

// TestRdzenNiepodlaczonyOddajeOdmoweZamiastCiszy sprawdza drugie zachowanie fail-open:
// połączenie nawiązane przed podłączeniem rdzenia żyje, a jego komendy dostają
// odpowiedź zamiast ciszy.
func TestRdzenNiepodlaczonyOddajeOdmoweZamiastCiszy(t *testing.T) {
	serwer, adres := podnies(t, nil)
	gniazdo := polacz(t, adres, "")

	wyslij(t, gniazdo, protocol.Koperta{Type: shared.CommandSessionList, Id: "przed-rdzeniem"})
	odpowiedz := odbierz(t, gniazdo)
	if odpowiedz.Id != "przed-rdzeniem" {
		t.Errorf("odmowa zgubiła identyfikator żądania: %q", odpowiedz.Id)
	}
	if odpowiedz.Status == nil || *odpowiedz.Status != shared.EnvelopeStatusError {
		t.Error("komenda bez rdzenia nie dostała stanu błędu")
	}
	if odpowiedz.Error == nil || odpowiedz.Error.Code != shared.ErrorCodeNotFound {
		t.Errorf("komenda bez rdzenia dostała %+v, oczekiwany kod %q",
			odpowiedz.Error, shared.ErrorCodeNotFound)
	}

	// Typ spoza kontraktu idzie drogą zdarzenia obszaru.
	wyslij(t, gniazdo, protocol.Koperta{Type: "session.wymyslona", Id: "spoza-kontraktu"})
	spoza := odbierz(t, gniazdo)
	if spoza.Type != shared.ZdarzenieNieznanej("session.wymyslona") {
		t.Errorf("typ spoza kontraktu dostał %q", spoza.Type)
	}

	// Po podłączeniu rdzenia połączenie zaczyna być obsługiwane: rdzeń jest pobierany przy komunikacie.
	serwer.PodlaczRdzen(&rdzenAtrapa{})
	wyslij(t, gniazdo, protocol.Koperta{Type: shared.CommandSessionList, Id: "po-rdzeniu"})
	dalsza := odbierz(t, gniazdo)
	if dalsza.Status == nil || *dalsza.Status != shared.EnvelopeStatusOk {
		t.Errorf("po podłączeniu rdzenia komenda nie została obsłużona: %+v", dalsza)
	}
}

// TestPochodzenieObceJestOdrzucane sprawdza jedyną zaporę nawiązania. Nie jest
// to kontrola dostępu — odcina stronę trzecią, która namówiła przeglądarkę
// Operatora, żeby otworzyła gniazdo do jego rdzenia.
func TestPochodzenieObceJestOdrzucane(t *testing.T) {
	_, adres := podnies(t, &rdzenAtrapa{})

	ctx, przerwij := context.WithTimeout(context.Background(), 5*time.Second)
	defer przerwij()
	gniazdo, _, err := websocket.Dial(ctx, adres, &websocket.DialOptions{
		HTTPHeader: http.Header{"Origin": []string{"https://strona-obca.example"}},
	})
	if err == nil {
		_ = gniazdo.CloseNow()
		t.Fatal("gniazdo z obcym pochodzeniem zostało nawiązane")
	}
}

// TestPochodzeniePetliZwrotnejSzostejWersjiPrzechodzi pilnuje wzorca `[::1]` bez
// portu, którego brak odcinałby powłokę podającą stronę z adresu szóstej wersji.
func TestPochodzeniePetliZwrotnejSzostejWersjiPrzechodzi(t *testing.T) {
	_, adres := podnies(t, &rdzenAtrapa{})

	ctx, przerwij := context.WithTimeout(context.Background(), 5*time.Second)
	defer przerwij()
	gniazdo, _, err := websocket.Dial(ctx, adres, &websocket.DialOptions{
		HTTPHeader: http.Header{"Origin": []string{"http://[::1]"}},
	})
	if err != nil {
		t.Fatalf("pochodzenie pętli zwrotnej szóstej wersji zostało odrzucone: %v", err)
	}
	_ = gniazdo.CloseNow()
}

// TestZajetyPortZatrzymujeStartOdmowa sprawdza jedyny błąd, który Uruchom
// zwraca — i to, że zwraca go zamiast wstać po cichu na innym porcie.
func TestZajetyPortZatrzymujeStartOdmowa(t *testing.T) {
	pierwszy, _ := podnies(t, &rdzenAtrapa{})

	_, zapisPortu, err := net.SplitHostPort(pierwszy.Adres())
	if err != nil {
		t.Fatalf("adres nasłuchu bez portu: %v", err)
	}
	port, err := strconv.Atoi(zapisPortu)
	if err != nil {
		t.Fatalf("nieczytelny numer portu: %v", err)
	}

	ustawienia := Domyslne()
	ustawienia.Adres = "127.0.0.1"
	ustawienia.Port = port
	ustawienia.Dziennik = log.New(io.Discard, "", 0)

	drugi := Nowy(ustawienia)
	zycie, zakoncz := context.WithCancel(context.Background())
	defer zakoncz()
	if err := drugi.Uruchom(zycie); err == nil {
		_ = drugi.Zamknij()
		t.Error("drugi nasłuch wstał na zajętym porcie")
	}
}

// TestNiekompletnaParaTlsNieWstaje sprawdza, że sprawdzenie pary idzie przed
// otwarciem nasłuchu — czyli że rdzeń ze wskazaniem połowicznym nie zdąży
// przyjąć ani jednego połączenia otwartym tekstem.
func TestNiekompletnaParaTlsNieWstaje(t *testing.T) {
	ustawienia := Domyslne()
	ustawienia.Adres = "127.0.0.1"
	ustawienia.PortDowolny = true
	ustawienia.CertyfikatTLS = "cert.pem"
	ustawienia.Dziennik = log.New(io.Discard, "", 0)

	serwer := Nowy(ustawienia)
	zycie, zakoncz := context.WithCancel(context.Background())
	defer zakoncz()
	if err := serwer.Uruchom(zycie); err == nil {
		_ = serwer.Zamknij()
		t.Error("nasłuch wstał ze wskazaniem połowicznej pary TLS")
	}
}

// poczekajNaPolaczenia czeka, aż rejestr osiągnie oczekiwany rozmiar. Rejestracja
// i wykreślenie idą w biegu obsługi gniazda, więc czekanie na stan zamiast stałej
// przerwy znosi zależność od chwili odczytu.
func poczekajNaPolaczenia(t *testing.T, serwer *Serwer, oczekiwane int) {
	t.Helper()

	termin := time.Now().Add(5 * time.Second)
	for time.Now().Before(termin) {
		if serwer.polaczenia.liczba() == oczekiwane {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatalf("rejestr połączeń stanął na %d, oczekiwane %d",
		serwer.polaczenia.liczba(), oczekiwane)
}
