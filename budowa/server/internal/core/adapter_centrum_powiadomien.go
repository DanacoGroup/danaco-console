// Odpowiedzialność pliku: centrum powiadomień jest jedynym mechanizmem
// powiadamiania platformy — zdarzenie wchodzi do rejestru wyłącznie od strony
// rdzenia, a doręczenie na telefon idzie kolejką doręczeń dopiero po zapisie
// w rejestrze.
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

// Kształt centrum powiadomień ma jedno źródło prawdy — plik
// `shared/contract.json`, z którego pochodzą typy żądania i odpowiedzi
// rejestru poniżej.
type (
	// ZadanieRejestruPowiadomien jest żądaniem odczytu rejestru centrum
	// powiadomień, niosącym filtry klasy, wagi, stanu, środowiska i sesji.
	ZadanieRejestruPowiadomien = shared.NotificationListRequest
	// WynikRejestruPowiadomien jest odpowiedzią odczytu rejestru centrum
	// powiadomień, niosącą wykaz pozycji oraz liczbę pozycji nieprzeczytanych.
	WynikRejestruPowiadomien = shared.NotificationListResponse
)

// CentrumPowiadomien obsługuje rodzinę komend `notification.*`: odczyt
// wykazu, potwierdzenie odczytania, zamknięcie oraz odłożenie zdarzenia
// w czasie.
type CentrumPowiadomien interface {
	Wykaz(ctx context.Context, z shared.NotificationListRequest) (shared.NotificationListResponse, error)
	Odczytaj(ctx context.Context, z shared.NotificationAcknowledgeRequest) (shared.NotificationAcknowledgeResponse, error)
	Zamknij(ctx context.Context, z shared.NotificationResolveRequest) (shared.NotificationResolveResponse, error)
	Odloz(ctx context.Context, z shared.NotificationSnoozeRequest) (shared.NotificationSnoozeResponse, error)
}

// ZgloszenieCentrum to zdarzenie wnoszone do rejestru przez rdzeń. Pole Waga
// nie jest dowolne: wiąże je z klasą funkcja `wagaKlasy`, a nadanie własnej
// wagi pozwala jedynie podnieść wagę zdarzenia szczególnego, nigdy obniżyć
// jej poniżej wagi klasy.
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
	// Zdarzenia odłożone wracają przy odczycie, a nie osobnym zegarem.
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
	a.rozglosZmiane(ctx, z.Ids, shared.NotificationStateOdczytane, nowe)
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
		a.rozglosZmiane(ctx, []string{z.Id}, shared.NotificationStateObsluzone, nowe)
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
	// Odłożenie w przeszłość wróciłoby przy najbliższym odczycie, czyli natychmiast.
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
		a.rozglosZmiane(ctx, []string{z.Id}, shared.NotificationStateOdlozone, nowe)
	}
	return shared.NotificationSnoozeResponse{Snoozed: zmienione, Unread: nowe}, nil
}

// ── wejście od strony rdzenia ────────────────────────────────────────────────

// Zglos wnosi zdarzenie do rejestru centrum i, gdy nastawy na to pozwalają,
// do kolejki doręczeń funkcji Mobile. Zwraca `false`, gdy klasa jest wygaszona
// nastawą — to wola Operatora, nie usterka do zgłoszenia.
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
		a.nadawca.wyslijDoKonta(ctx, shared.EventNotificationRaised, "", shared.NotificationRaisedEvent{
			Notification: pozycja,
			Unread:       nowe,
		})
	}

	if naTelefon {
		// Jedyne wywołanie silnika kolejki doręczeń; niepowodzenie nie cofa zapisu.
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

// rozglosZmiane niesie zmianę stanu zdarzenia do wszystkich połączeń
// Operatora, wraz z nowym stanem i liczbą pozycji nieprzeczytanych w rejestrze.
func (a *adapterCentrumPowiadomien) rozglosZmiane(ctx context.Context, identyfikatory []string,
	stan shared.NotificationState, nowe int) {

	if a.nadawca == nil {
		return
	}
	a.nadawca.wyslijDoKonta(ctx, shared.EventNotificationChanged, "", shared.NotificationChangedEvent{
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

// kanalDopuszczony mówi, czy zdarzenie danej klasy ma iść, obok rejestru
// centrum, także kanałem dodatkowym wskazanym w nastawach Operatora.
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

// wygaszone rozpoznaje wyłącznie jawne „nie" w zapisie nastawy; zapis pusty
// albo nierozpoznany zostawia w mocy wartość domyślną katalogu.
func wygaszone(zapis string) bool {
	wartosc, wskazana := wartoscWymoguLogowania(zapis)
	return wskazana && !wartosc
}

// bladCentrum składa odmowę rodziny komend `notification.*`, wiążąc
// przekazany kod błędu z powodem opisującym centrum powiadomień.
func bladCentrum(kod protocol.KodBledu, powod string) error {
	return protocol.JakoError(protocol.NowyBlad(kod, "centrum powiadomień: "+powod))
}

// ── przekłady ────────────────────────────────────────────────────────────────

// powiadomienieKontraktu przekłada wiersz rejestru centrum na kształt
// `shared.Notification` zgodny z kontraktem, z polami opcjonalnymi
// ustawionymi warunkowo.
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

// wagaZdarzenia bierze wagę wprost ze zgłoszenia, gdy ją niesie, a przy jej
// braku wyprowadza wagę z taksonomii przypisanej klasie zdarzenia.
func wagaZdarzenia(z ZgloszenieCentrum) shared.NotificationWeight {
	if z.Waga != "" {
		return z.Waga
	}
	return wagaKlasy(z.Klasa)
}

// wagaKlasy wiąże klasę zdarzenia z wagą domyślną: decyzja i błąd niosą wagę
// wymagającą decyzji, zakończenie i wzmianka wagę normalną, reszta informacyjną.
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

// tytulZgloszenia składa tytuł listu kolejki doręczeń z pierwszego zdania
// treści zgłoszenia, a gdy zdanie jest zbyt długie, ucina je do
// osiemdziesięciu znaków.
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

// priorytetZgloszenia przekłada wagę centrum powiadomień na priorytet kolejki
// doręczeń: waga wymagająca decyzji niesie priorytet pilny, reszta zwykły.
func priorytetZgloszenia(waga shared.NotificationWeight) string {
	if waga == shared.NotificationWeightWymagajacaDecyzji {
		return "pilny"
	}
	return "zwykly"
}

// akcjeAlbatPuste oddaje wykaz akcji pusty zamiast wartości `null`, ponieważ
// kolumna rejestru trzyma zapis JSON i nie przyjmuje wartości pustej.
func akcjeAlbatPuste(akcje []shared.NotificationAction) []shared.NotificationAction {
	if akcje == nil {
		return []shared.NotificationAction{}
	}
	return akcje
}

// napisyWyliczenia przenosi wykaz wartości typu wyliczeniowego na wykaz
// napisów, pomijając po drodze wartości puste, do filtrowania zapytania rejestru.
func napisyWyliczenia[T ~string](wartosci []T) []string {
	napisy := make([]string, 0, len(wartosci))
	for _, wartosc := range wartosci {
		if wartosc != "" {
			napisy = append(napisy, string(wartosc))
		}
	}
	return napisy
}

// numeryZdarzen przekłada wykaz identyfikatorów kontraktu na wykaz numerów
// wierszy rejestru, odmawiając przy pierwszym identyfikatorze nierozpoznanym.
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

// numerZdarzenia rozpoznaje identyfikator zdarzenia zapisany jako liczba
// dodatnia; identyfikator pusty, ujemny albo nieliczbowy kończy się odmową.
func numerZdarzenia(identyfikator string) (int64, error) {
	numer, err := strconv.ParseInt(strings.TrimSpace(identyfikator), 10, 64)
	if err != nil || numer <= 0 {
		return 0, bladCentrum(shared.ErrorCodeValidationFailed,
			"identyfikator zdarzenia „"+identyfikator+"” nie jest identyfikatorem rejestru centrum")
	}
	return numer, nil
}
