package core

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"danacoconsole/shared"
)

// Sprawdziany centrum powiadomień mierzą SKUTEK: co zostaje w rejestrze po
// zgłoszeniu i co widzi Operator, a nie to, że wywołanie wróciło bez błędu.
//
// Uprząż jest ta sama, co dla pozostałych sprawdzianów skutku
// (`zmontujDoPomiaruSkutku`): świeża baza, pełny montaż, komendy przez rejestr.

// zglosDoCentrum wnosi zdarzenie drogą wewnętrzną — tą, którą idzie rdzeń.
// Komendy zgłaszającej kontrakt nie ma z zamysłu, więc sprawdzian sięga po
// adapter tak samo jak reszta rdzenia.
func zglosDoCentrum(t *testing.T, zmontowany *Zmontowany, zycie context.Context,
	zgloszenie ZgloszenieCentrum) bool {

	t.Helper()
	// Adapter składa się nad tym samym repozytorium i tym samym adapterem
	// ustawień, co w montażu — droga wewnętrzna nie ma innego wejścia.
	centrum := nowyAdapterCentrumPowiadomien(zmontowany.dane.CentrumPowiadomien).
		ZNastawami(zmontowany.ustawienia)
	wniesione, err := centrum.Zglos(zycie, zgloszenie)
	if err != nil {
		t.Fatalf("zgłoszenie do centrum nie powiodło się: %v", err)
	}
	return wniesione
}

func TestZgloszenieWchodziDoRejestruIPodnosiLicznik(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	zglosDoCentrum(t, zmontowany, zycie, ZgloszenieCentrum{
		Klasa: shared.NotificationClassZakonczenie,
		Tresc: "Przebieg pętli wykonawczej zakończony",
	})

	var wykaz shared.NotificationListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandNotificationList,
		shared.NotificationListRequest{}, &wykaz)

	if len(wykaz.Notifications) != 1 {
		t.Fatalf("rejestr niesie %d zdarzeń, oczekiwane 1", len(wykaz.Notifications))
	}
	if wykaz.Unread != 1 {
		t.Errorf("licznik zdarzeń nowych: %d, oczekiwany 1", wykaz.Unread)
	}
	zdarzenie := wykaz.Notifications[0]
	if zdarzenie.State != shared.NotificationStateNowe {
		t.Errorf("stan zdarzenia: %q, oczekiwany nowe", zdarzenie.State)
	}
	// Waga bierze się z taksonomii klasy, a nie ze zgłoszenia — zgłoszenie jej
	// nie podało.
	if zdarzenie.Weight != shared.NotificationWeightNormalna {
		t.Errorf("waga zdarzenia klasy zakonczenie: %q, oczekiwana normalna", zdarzenie.Weight)
	}
	// Klasa „zakończenie" ma w katalogu (migracja 377, makieta 5 rozdz. 7.4)
	// kanał Mobile włączony domyślnie, więc zdarzenie idzie obiema drogami.
	// Nastawy rozstrzygają o kanale, nie kod adaptera.
	if zdarzenie.Delivery != shared.NotificationDeliveryCentrumIPush {
		t.Errorf("kanał dostarczenia: %q, oczekiwany centrum_i_push", zdarzenie.Delivery)
	}
}

func TestKlasaWygaszonaNastawaNieWchodziDoRejestru(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	// Operator wyłącza klasę „termin" w sekcji Powiadomień — tą samą drogą, co
	// każdą inną nastawę platformy.
	wykonajUdana(t, zmontowany, zycie, shared.CommandConfigSet, shared.ConfigSetRequest{
		Key:   "powiadomienia.klasa.termin",
		Value: json.RawMessage("false"),
		Scope: shared.ConfigScopeGlobal,
	}, nil)

	wniesione := zglosDoCentrum(t, zmontowany, zycie, ZgloszenieCentrum{
		Klasa: shared.NotificationClassTermin,
		Tresc: "Termin zadania mija jutro",
	})

	if wniesione {
		t.Error("zdarzenie klasy wygaszonej weszło do rejestru")
	}

	var wykaz shared.NotificationListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandNotificationList,
		shared.NotificationListRequest{}, &wykaz)
	if len(wykaz.Notifications) != 0 {
		t.Errorf("rejestr niesie %d zdarzeń, oczekiwane 0", len(wykaz.Notifications))
	}
}

func TestPrzelacznikGlownyWygaszaWszystkieKlasy(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	wykonajUdana(t, zmontowany, zycie, shared.CommandConfigSet, shared.ConfigSetRequest{
		Key:   "powiadomienia.wlaczone",
		Value: json.RawMessage("false"),
		Scope: shared.ConfigScopeGlobal,
	}, nil)

	for _, klasa := range []shared.NotificationClass{
		shared.NotificationClassBlad,
		shared.NotificationClassDecyzja,
		shared.NotificationClassSystem,
	} {
		if zglosDoCentrum(t, zmontowany, zycie, ZgloszenieCentrum{Klasa: klasa, Tresc: "cokolwiek"}) {
			t.Errorf("klasa %q weszła do rejestru mimo wyłączonego przełącznika głównego", klasa)
		}
	}
}

func TestOdczytanieZbiorczeZerujeLicznik(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	for _, tresc := range []string{"pierwsze", "drugie", "trzecie"} {
		zglosDoCentrum(t, zmontowany, zycie, ZgloszenieCentrum{
			Klasa: shared.NotificationClassSystem, Tresc: tresc,
		})
	}

	var potwierdzenie shared.NotificationAcknowledgeResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandNotificationAcknowledge,
		shared.NotificationAcknowledgeRequest{}, &potwierdzenie)

	if potwierdzenie.Acknowledged != 3 {
		t.Errorf("oznaczono %d zdarzeń, oczekiwane 3", potwierdzenie.Acknowledged)
	}
	if potwierdzenie.Unread != 0 {
		t.Errorf("licznik po oznaczeniu: %d, oczekiwany 0", potwierdzenie.Unread)
	}

	// Zdarzenia zostają w rejestrze — odczytanie nie jest skasowaniem.
	var wykaz shared.NotificationListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandNotificationList,
		shared.NotificationListRequest{}, &wykaz)
	if len(wykaz.Notifications) != 3 {
		t.Errorf("rejestr niesie %d zdarzeń, oczekiwane 3", len(wykaz.Notifications))
	}
}

func TestZamknieteZdarzenieZnikaZWidokuDomyslnego(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	zglosDoCentrum(t, zmontowany, zycie, ZgloszenieCentrum{
		Klasa: shared.NotificationClassDecyzja, Tresc: "Krok czeka na zatwierdzenie",
	})

	var wykaz shared.NotificationListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandNotificationList,
		shared.NotificationListRequest{}, &wykaz)

	var zamkniecie shared.NotificationResolveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandNotificationResolve,
		shared.NotificationResolveRequest{Id: wykaz.Notifications[0].Id}, &zamkniecie)

	if !zamkniecie.Resolved {
		t.Error("zamknięcie zdarzenia nie zmieniło stanu")
	}

	var poZamknieciu shared.NotificationListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandNotificationList,
		shared.NotificationListRequest{}, &poZamknieciu)
	if len(poZamknieciu.Notifications) != 0 {
		t.Errorf("widok domyślny niesie %d zdarzeń zamkniętych, oczekiwane 0",
			len(poZamknieciu.Notifications))
	}

	// Zapytane wprost — zdarzenie nadal jest. Zamknięcie nie kasuje wiersza.
	var wprost shared.NotificationListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandNotificationList,
		shared.NotificationListRequest{States: []shared.NotificationState{
			shared.NotificationStateObsluzone,
		}}, &wprost)
	if len(wprost.Notifications) != 1 {
		t.Errorf("rejestr zgubił zdarzenie zamknięte: %d", len(wprost.Notifications))
	}
}

func TestOdlozenieWPrzeszloscOdmawia(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	zglosDoCentrum(t, zmontowany, zycie, ZgloszenieCentrum{
		Klasa: shared.NotificationClassBlad, Tresc: "Zadanie zawiodło",
	})
	var wykaz shared.NotificationListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandNotificationList,
		shared.NotificationListRequest{}, &wykaz)

	odpowiedz := wykonajKomende(t, zmontowany, zycie, shared.CommandNotificationSnooze,
		shared.NotificationSnoozeRequest{
			Id:    wykaz.Notifications[0].Id,
			Until: int(time.Now().Add(-time.Hour).UnixMilli()),
		})

	if odpowiedz.Error == nil {
		t.Fatal("odłożenie w przeszłość przeszło — zdarzenie wróciłoby natychmiast")
	}
	if odpowiedz.Error.Code != shared.ErrorCodeValidationFailed {
		t.Errorf("kod odmowy: %q, oczekiwany %q", odpowiedz.Error.Code,
			shared.ErrorCodeValidationFailed)
	}
}

func TestZdarzenieOdlozoneWracaPoTerminie(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	// Klasa „wzmianka" nie ma kanału dodatkowego w katalogu, więc zdarzenie
	// zostaje w samym centrum — sprawdzian odłożenia nie dotyka kolejki doręczeń.
	zglosDoCentrum(t, zmontowany, zycie, ZgloszenieCentrum{
		Klasa: shared.NotificationClassWzmianka, Tresc: "Wzmianka w komentarzu",
	})
	var wykaz shared.NotificationListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandNotificationList,
		shared.NotificationListRequest{}, &wykaz)

	var odlozenie shared.NotificationSnoozeResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandNotificationSnooze,
		shared.NotificationSnoozeRequest{
			Id:    wykaz.Notifications[0].Id,
			Until: int(time.Now().Add(50 * time.Millisecond).UnixMilli()),
		}, &odlozenie)

	if odlozenie.Unread != 0 {
		t.Errorf("zdarzenie odłożone liczy się do plakietki: %d", odlozenie.Unread)
	}

	// Powrót zachodzi przy odczycie rejestru — Operator patrzy, zdarzenie wraca.
	time.Sleep(80 * time.Millisecond)
	var poPowrocie shared.NotificationListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandNotificationList,
		shared.NotificationListRequest{}, &poPowrocie)

	if poPowrocie.Unread != 1 {
		t.Errorf("zdarzenie nie wróciło po terminie: licznik %d, oczekiwany 1", poPowrocie.Unread)
	}
}

func TestFiltrZawezaWykazAleNieLicznik(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)

	zglosDoCentrum(t, zmontowany, zycie, ZgloszenieCentrum{
		Klasa: shared.NotificationClassBlad, Tresc: "Zadanie zawiodło",
	})
	zglosDoCentrum(t, zmontowany, zycie, ZgloszenieCentrum{
		Klasa: shared.NotificationClassSystem, Tresc: "Zmiana urządzeń",
	})

	var wykaz shared.NotificationListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandNotificationList,
		shared.NotificationListRequest{
			Classes: []shared.NotificationClass{shared.NotificationClassBlad},
		}, &wykaz)

	if len(wykaz.Notifications) != 1 {
		t.Errorf("filtr klasy oddał %d zdarzeń, oczekiwane 1", len(wykaz.Notifications))
	}
	// Plakietka liczy CAŁY rejestr, nie widok — inaczej zawężenie filtrem
	// wyglądałoby jak obsłużenie zdarzeń spoza filtru.
	if wykaz.Unread != 2 {
		t.Errorf("licznik przy zawężonym widoku: %d, oczekiwany 2", wykaz.Unread)
	}
}
