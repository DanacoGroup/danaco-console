// Odpowiedzialność pliku: centrum powiadomień — jedyny mechanizm powiadamiania
// platformy (karta komponentu, rozdz. 11.6; encja, rozdz. 18.4).
//
// ── dwa wejścia, jedno wyjście ─────────────────────────────────────────────
// Zdarzenie wchodzi do rejestru WYŁĄCZNIE od strony rdzenia (`Zglos`), bo tylko
// rdzeń wie, że coś zaszło. Kontrakt niesie sam odczyt i zmianę stanu — komendy
// zgłaszającej nie ma i mieć nie powinien, inaczej klient wpisywałby do rejestru
// zdarzenia, które nigdy nie zaszły.
//
// ── nastawy rozstrzygają, nie kod ──────────────────────────────────────────
// O tym, czy klasa zdarzenia w ogóle wchodzi do rejestru i czy idzie dalej na
// telefon albo listem, rozstrzygają nastawy sekcji „Powiadomienia" okna Ustawień
// (migracja 377). Adapter ich nie powtarza i nie zna ani jednej wartości
// domyślnej — czyta je tym samym rozstrzygaczem, którym idzie każde inne
// ustawienie platformy.
//
// ── kolejka doręczeń dostaje wołacza ───────────────────────────────────────
// Silnik `zdalne.Zglos` stał zbudowany i nieużywany. Tu jest jego jedyne
// wywołanie: zdarzenie klasy dopuszczonej do kanału `mobile` wchodzi do rejestru
// centrum i zaraz potem do kolejki doręczeń. Niepowodzenie kolejki NIE cofa
// zapisu w rejestrze — zdarzenie zaszło niezależnie od tego, czy telefon je
// odebrał, a rejestr centrum jest kanałem podstawowym.
package core

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/zdalne"
	"danacoconsole/shared"
)

// Kształt centrum ma jedno źródło prawdy — `shared/contract.json`.
type (
	// ZadanieRejestruPowiadomien jest żądaniem odczytu rejestru.
	ZadanieRejestruPowiadomien = shared.NotificationListRequest
	// WynikRejestruPowiadomien jest odpowiedzią odczytu rejestru.
	WynikRejestruPowiadomien = shared.NotificationListResponse
)

// CentrumPowiadomien obsługuje rodzinę `notification.*`.
type CentrumPowiadomien interface {
	Wykaz(ctx context.Context, z shared.NotificationListRequest) (shared.NotificationListResponse, error)
	Odczytaj(ctx context.Context, z shared.NotificationAcknowledgeRequest) (shared.NotificationAcknowledgeResponse, error)
	Zamknij(ctx context.Context, z shared.NotificationResolveRequest) (shared.NotificationResolveResponse, error)
	Odloz(ctx context.Context, z shared.NotificationSnoozeRequest) (shared.NotificationSnoozeResponse, error)
}

// ZgloszenieCentrum to zdarzenie wnoszone do rejestru przez rdzeń.
//
// Waga nie jest polem dowolnym: taksonomia rozdz. 11.6 wiąże ją z klasą i to
// wiązanie stoi w `wagaKlasy`. Pole zostaje w kształcie zgłoszenia wyłącznie po
// to, żeby wołający mógł podnieść wagę zdarzenia szczególnego, nigdy po to, by
// ją obniżyć poniżej klasy.
type ZgloszenieCentrum struct {
	Klasa         shared.NotificationClass
	Waga          shared.NotificationWeight
	Tresc         string
	ZrodloTyp     shared.NotificationSourceKind
	ZrodloID      string
	SrodowiskoKod string
	SesjaID       string
	Akcje         []shared.NotificationAction
}

type adapterCentrumPowiadomien struct {
	rejestr dane.RepozytoriumCentrumPowiadomien
	nadawca *emiter
	nastawy NastawyPlatformy
}

var _ CentrumPowiadomien = (*adapterCentrumPowiadomien)(nil)

func nowyAdapterCentrumPowiadomien(
	rejestr dane.RepozytoriumCentrumPowiadomien) *adapterCentrumPowiadomien {

	return &adapterCentrumPowiadomien{rejestr: rejestr}
}

// ZNadawca wpina drogę rozgłoszenia. Bez niej rejestr działa, ale plakietka
// i kolumna centrum dowiadują się o zmianie dopiero przy następnym odczycie.
func (a *adapterCentrumPowiadomien) ZNadawca(n *emiter) *adapterCentrumPowiadomien {
	a.nadawca = n
	return a
}

// ZNastawami wpina drogę odczytu nastaw sekcji „Powiadomienia". Bez niej
// obowiązują wartości domyślne katalogu, bo brak drogi odczytu nie jest
// wskazaniem Operatora.
func (a *adapterCentrumPowiadomien) ZNastawami(n NastawyPlatformy) *adapterCentrumPowiadomien {
	a.nastawy = n
	return a
}

// ── odczyt i zmiana stanu ────────────────────────────────────────────────────

func (a *adapterCentrumPowiadomien) Wykaz(ctx context.Context,
	z shared.NotificationListRequest) (shared.NotificationListResponse, error) {

	if a.rejestr == nil {
		return shared.NotificationListResponse{Notifications: []shared.Notification{}}, nil
	}
	// Zdarzenia odłożone wracają przy odczycie, a nie budzikiem. Centrum czyta
	// się wtedy, gdy Operator na nie patrzy; osobny takt dla rejestru, który i tak
	// pyta się przy każdym otwarciu kolumny, byłby drugim zegarem bez odbiorcy.
	if _, err := a.rejestr.Przywroc(ctx, time.Now()); err != nil {
		return shared.NotificationListResponse{}, err
	}

	wiersze, err := a.rejestr.Wykaz(ctx, dane.FiltrCentrum{
		Klasy:         napisyWyliczenia(z.Classes),
		Wagi:          napisyWyliczenia(z.Weights),
		Stany:         napisyWyliczenia(z.States),
		SrodowiskoKod: wartoscTekstu(z.EnvironmentId),
		SesjaID:       wartoscTekstu(z.SessionId),
		Limit:         wartoscLiczby(z.Limit),
	})
	if err != nil {
		return shared.NotificationListResponse{}, err
	}
	nowe, err := a.rejestr.Nowe(ctx)
	if err != nil {
		return shared.NotificationListResponse{}, err
	}

	pozycje := make([]shared.Notification, 0, len(wiersze))
	for _, wiersz := range wiersze {
		pozycje = append(pozycje, powiadomienieKontraktu(wiersz))
	}
	return shared.NotificationListResponse{Notifications: pozycje, Unread: nowe}, nil
}

func (a *adapterCentrumPowiadomien) Odczytaj(ctx context.Context,
	z shared.NotificationAcknowledgeRequest) (shared.NotificationAcknowledgeResponse, error) {

	if a.rejestr == nil {
		return shared.NotificationAcknowledgeResponse{}, nil
	}
	numery, err := numeryZdarzen(z.Ids)
	if err != nil {
		return shared.NotificationAcknowledgeResponse{}, err
	}
	ile, err := a.rejestr.Odczytaj(ctx, numery)
	if err != nil {
		return shared.NotificationAcknowledgeResponse{}, err
	}
	nowe, err := a.rejestr.Nowe(ctx)
	if err != nil {
		return shared.NotificationAcknowledgeResponse{}, err
	}
	a.rozglosZmiane(z.Ids, shared.NotificationStateOdczytane, nowe)
	return shared.NotificationAcknowledgeResponse{Acknowledged: ile, Unread: nowe}, nil
}

func (a *adapterCentrumPowiadomien) Zamknij(ctx context.Context,
	z shared.NotificationResolveRequest) (shared.NotificationResolveResponse, error) {

	if a.rejestr == nil {
		return shared.NotificationResolveResponse{}, nil
	}
	numer, err := numerZdarzenia(z.Id)
	if err != nil {
		return shared.NotificationResolveResponse{}, err
	}
	zmienione, err := a.rejestr.Zamknij(ctx, numer)
	if err != nil {
		return shared.NotificationResolveResponse{}, err
	}
	nowe, err := a.rejestr.Nowe(ctx)
	if err != nil {
		return shared.NotificationResolveResponse{}, err
	}
	if zmienione {
		a.rozglosZmiane([]string{z.Id}, shared.NotificationStateObsluzone, nowe)
	}
	return shared.NotificationResolveResponse{Resolved: zmienione, Unread: nowe}, nil
}

func (a *adapterCentrumPowiadomien) Odloz(ctx context.Context,
	z shared.NotificationSnoozeRequest) (shared.NotificationSnoozeResponse, error) {

	if a.rejestr == nil {
		return shared.NotificationSnoozeResponse{}, nil
	}
	numer, err := numerZdarzenia(z.Id)
	if err != nil {
		return shared.NotificationSnoozeResponse{}, err
	}
	// Odłożenie w przeszłość jest odczytaniem przebranym za odłożenie: zdarzenie
	// wróciłoby przy najbliższym odczycie, czyli natychmiast.
	if int64(z.Until) <= time.Now().UnixMilli() {
		return shared.NotificationSnoozeResponse{}, bladCentrum(shared.ErrorCodeValidationFailed,
			"chwila powrotu leży w przeszłości — zdarzenie wróciłoby natychmiast; "+
				"od zamknięcia zdarzenia jest notification.resolve")
	}
	zmienione, err := a.rejestr.Odloz(ctx, numer, time.UnixMilli(int64(z.Until)))
	if err != nil {
		return shared.NotificationSnoozeResponse{}, err
	}
	nowe, err := a.rejestr.Nowe(ctx)
	if err != nil {
		return shared.NotificationSnoozeResponse{}, err
	}
	if zmienione {
		a.rozglosZmiane([]string{z.Id}, shared.NotificationStateOdlozone, nowe)
	}
	return shared.NotificationSnoozeResponse{Snoozed: zmienione, Unread: nowe}, nil
}

// ── wejście od strony rdzenia ────────────────────────────────────────────────

// Zglos wnosi zdarzenie do rejestru centrum i — gdy nastawy na to pozwalają —
// do kolejki doręczeń funkcji Mobile.
//
// Zwraca `false`, gdy klasa jest wygaszona nastawą: to nie jest usterka, tylko
// wola Operatora, więc wołający nie ma czego zgłaszać jako błąd.
func (a *adapterCentrumPowiadomien) Zglos(ctx context.Context, z ZgloszenieCentrum) (bool, error) {
	if a.rejestr == nil || strings.TrimSpace(z.Tresc) == "" {
		return false, nil
	}
	if !a.klasaCzynna(ctx, z.Klasa) {
		return false, nil
	}

	naTelefon := a.kanalDopuszczony(ctx, z.Klasa, kanalMobile)
	kanal := shared.NotificationDeliveryCentrum
	if naTelefon {
		kanal = shared.NotificationDeliveryCentrumIPush
	}

	akcje, err := json.Marshal(akcjeAlbatPuste(z.Akcje))
	if err != nil {
		return false, err
	}

	wiersz, err := a.rejestr.Zapisz(ctx, dane.ZdarzenieCentrum{
		Klasa:         string(z.Klasa),
		Waga:          string(wagaZdarzenia(z)),
		Tresc:         z.Tresc,
		ZrodloTyp:     string(z.ZrodloTyp),
		ZrodloID:      z.ZrodloID,
		SrodowiskoKod: z.SrodowiskoKod,
		SesjaID:       z.SesjaID,
		Akcje:         string(akcje),
		Kanal:         string(kanal),
	})
	if err != nil {
		return false, err
	}

	nowe, err := a.rejestr.Nowe(ctx)
	if err != nil {
		return false, err
	}
	pozycja := powiadomienieKontraktu(wiersz)
	if a.nadawca != nil {
		a.nadawca.wyslij(shared.EventNotificationRaised, "", shared.NotificationRaisedEvent{
			Notification: pozycja,
			Unread:       nowe,
		})
	}

	if naTelefon {
		// Jedyne wywołanie silnika kolejki doręczeń w całym rdzeniu. Niepowodzenie
		// nie cofa zapisu w rejestrze: zdarzenie zaszło niezależnie od tego, czy
		// telefon je odebrał.
		_, _ = zdalne.Zglos(ctx, zdalne.Zgloszenie{
			Tytul:     tytulZgloszenia(z),
			Tresc:     z.Tresc,
			Priorytet: priorytetZgloszenia(pozycja.Weight),
			Powod:     "centrum powiadomień: " + string(z.Klasa),
			BytRodzaj: string(z.ZrodloTyp),
			BytID:     z.ZrodloID,
		})
	}
	return true, nil
}

// rozglosZmiane niesie zmianę stanu do wszystkich połączeń Operatora.
func (a *adapterCentrumPowiadomien) rozglosZmiane(identyfikatory []string,
	stan shared.NotificationState, nowe int) {

	if a.nadawca == nil {
		return
	}
	a.nadawca.wyslij(shared.EventNotificationChanged, "", shared.NotificationChangedEvent{
		Ids:    identyfikatory,
		State:  stan,
		Unread: nowe,
	})
}

// ── nastawy ──────────────────────────────────────────────────────────────────

const kanalMobile = "mobile"

// klasaCzynna czyta przełącznik główny i przełącznik klasy.
//
// Brak drogi odczytu i nastawa niezapisana znaczą to samo: obowiązuje wartość
// domyślna katalogu, czyli „aktywne". Milczenie nie jest wyłączeniem.
func (a *adapterCentrumPowiadomien) klasaCzynna(ctx context.Context,
	klasa shared.NotificationClass) bool {

	if a.nastawy == nil {
		return true
	}
	if wygaszone(a.nastawy.Nastawa(ctx, konfig.KluczPowiadomieniaWlaczone)) {
		return false
	}
	return !wygaszone(a.nastawy.Nastawa(ctx, konfig.KluczKlasyPowiadomien(string(klasa))))
}

// kanalDopuszczony mówi, czy klasa ma iść kanałem dodatkowym.
func (a *adapterCentrumPowiadomien) kanalDopuszczony(ctx context.Context,
	klasa shared.NotificationClass, kanal string) bool {

	if a.nastawy == nil {
		return false
	}
	zapis := a.nastawy.Nastawa(ctx, konfig.KluczKanalowPowiadomien(string(klasa)))
	if strings.TrimSpace(zapis) == "" {
		return false
	}
	var kanaly []string
	if err := json.Unmarshal([]byte(zapis), &kanaly); err != nil {
		return false
	}
	for _, wybrany := range kanaly {
		if wybrany == kanal {
			return true
		}
	}
	return false
}

// wygaszone rozpoznaje wyłącznie jawne „nie"; zapis pusty zostawia domyślną.
func wygaszone(zapis string) bool {
	wartosc, wskazana := wartoscWymoguLogowania(zapis)
	return wskazana && !wartosc
}

// bladCentrum składa odmowę rodziny `notification.*`.
func bladCentrum(kod protocol.KodBledu, powod string) error {
	return protocol.JakoError(protocol.NowyBlad(kod, "centrum powiadomień: "+powod))
}

// ── przekłady ────────────────────────────────────────────────────────────────

// powiadomienieKontraktu przekłada wiersz rejestru na kształt kontraktu.
func powiadomienieKontraktu(w dane.ZdarzenieCentrum) shared.Notification {
	pozycja := shared.Notification{
		Id:        strconv.FormatInt(w.ID, 10),
		Class:     shared.NotificationClass(w.Klasa),
		Weight:    shared.NotificationWeight(w.Waga),
		Text:      w.Tresc,
		State:     shared.NotificationState(w.Stan),
		Delivery:  shared.NotificationDelivery(w.Kanal),
		CreatedAt: int(w.ZnacznikCzasu.UnixMilli()),
	}
	if w.ZrodloTyp != "" {
		rodzaj := shared.NotificationSourceKind(w.ZrodloTyp)
		pozycja.SourceKind = &rodzaj
	}
	pozycja.SourceId = tekstOpcjonalny(w.ZrodloID)
	pozycja.EnvironmentId = tekstOpcjonalny(w.SrodowiskoKod)
	pozycja.SessionId = tekstOpcjonalny(w.SesjaID)
	if !w.OdlozoneDo.IsZero() {
		chwila := int(w.OdlozoneDo.UnixMilli())
		pozycja.SnoozedUntil = &chwila
	}
	var akcje []shared.NotificationAction
	if err := json.Unmarshal([]byte(w.Akcje), &akcje); err == nil && len(akcje) > 0 {
		pozycja.Actions = akcje
	}
	return pozycja
}

// wagaZdarzenia bierze wagę ze zgłoszenia, a przy jej braku z taksonomii klasy.
func wagaZdarzenia(z ZgloszenieCentrum) shared.NotificationWeight {
	if z.Waga != "" {
		return z.Waga
	}
	return wagaKlasy(z.Klasa)
}

// wagaKlasy wiąże klasę z wagą wprost z taksonomii rozdz. 11.6.
func wagaKlasy(klasa shared.NotificationClass) shared.NotificationWeight {
	switch klasa {
	case shared.NotificationClassDecyzja, shared.NotificationClassBlad:
		return shared.NotificationWeightWymagajacaDecyzji
	case shared.NotificationClassZakonczenie, shared.NotificationClassWzmianka:
		return shared.NotificationWeightNormalna
	default:
		return shared.NotificationWeightInformacyjna
	}
}

// tytulZgloszenia składa tytuł listu kolejki doręczeń — pierwsze zdanie treści.
func tytulZgloszenia(z ZgloszenieCentrum) string {
	tresc := strings.TrimSpace(z.Tresc)
	if kropka := strings.IndexAny(tresc, ".!?"); kropka > 0 && kropka < 80 {
		return tresc[:kropka]
	}
	if len(tresc) > 80 {
		return tresc[:80]
	}
	return tresc
}

// priorytetZgloszenia przekłada wagę centrum na priorytet kolejki doręczeń.
func priorytetZgloszenia(waga shared.NotificationWeight) string {
	if waga == shared.NotificationWeightWymagajacaDecyzji {
		return "pilny"
	}
	return "zwykly"
}

// akcjeAlbatPuste oddaje wykaz pusty zamiast `null` — kolumna trzyma JSON.
func akcjeAlbatPuste(akcje []shared.NotificationAction) []shared.NotificationAction {
	if akcje == nil {
		return []shared.NotificationAction{}
	}
	return akcje
}

// napisyWyliczenia przenosi wykaz wartości wyliczenia na wykaz napisów.
func napisyWyliczenia[T ~string](wartosci []T) []string {
	napisy := make([]string, 0, len(wartosci))
	for _, wartosc := range wartosci {
		if wartosc != "" {
			napisy = append(napisy, string(wartosc))
		}
	}
	return napisy
}

// numeryZdarzen przekłada identyfikatory kontraktu na numery wierszy.
func numeryZdarzen(identyfikatory []string) ([]int64, error) {
	numery := make([]int64, 0, len(identyfikatory))
	for _, identyfikator := range identyfikatory {
		numer, err := numerZdarzenia(identyfikator)
		if err != nil {
			return nil, err
		}
		numery = append(numery, numer)
	}
	return numery, nil
}

// numerZdarzenia rozpoznaje identyfikator zdarzenia.
func numerZdarzenia(identyfikator string) (int64, error) {
	numer, err := strconv.ParseInt(strings.TrimSpace(identyfikator), 10, 64)
	if err != nil || numer <= 0 {
		return 0, bladCentrum(shared.ErrorCodeValidationFailed,
			"identyfikator zdarzenia „"+identyfikator+"” nie jest identyfikatorem rejestru centrum")
	}
	return numer, nil
}
