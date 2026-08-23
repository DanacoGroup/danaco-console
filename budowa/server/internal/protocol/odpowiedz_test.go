package protocol

import (
	"errors"
	"testing"

	"danacoconsole/shared"
)

// TestSukcesNiesieStanKontraktu sprawdza, że odpowiedź udana ma stan wprost
// z kontraktu, a potwierdzenie bez treści zostaje bez ładunku — zamiast
// niepotrzebnego `null`.
func TestSukcesNiesieStanKontraktu(t *testing.T) {
	bezTresci, err := Sukces(nil)
	if err != nil {
		t.Fatalf("odpowiedź bez treści zgłosiła błąd: %v", err)
	}
	if bezTresci.Status != shared.EnvelopeStatusOk {
		t.Errorf("odpowiedź udana ma stan %q", bezTresci.Status)
	}
	if bezTresci.Wynik != nil {
		t.Errorf("potwierdzenie bez treści niesie ładunek: %s", bezTresci.Wynik)
	}
	if bezTresci.Blad != nil {
		t.Error("odpowiedź udana niesie błąd")
	}

	zTrescia, err := Sukces(map[string]int{"count": 2})
	if err != nil {
		t.Fatalf("odpowiedź z treścią zgłosiła błąd: %v", err)
	}
	if string(zTrescia.Wynik) != `{"count":2}` {
		t.Errorf("wynik zapisany jako %s", zTrescia.Wynik)
	}
}

// TestSukcesOdmawiaWynikuNiedajacegoSieZakodowac pilnuje, by niepowodzenie
// zapisu wyszło przy budowie odpowiedzi, a nie przy pakowaniu w kopertę —
// KopertaOdpowiedzi nie ma jak zawieść i nie zwraca błędu.
func TestSukcesOdmawiaWynikuNiedajacegoSieZakodowac(t *testing.T) {
	if _, err := Sukces(make(chan int)); err == nil {
		t.Error("wynik niedający się zakodować został przyjęty")
	}
}

// TestPorazkaKodemBierzePonawialnoscZKontraktu sprawdza, że o ponawialności
// rozstrzyga katalog kontraktu, a nie wywołujący. Kod przejściowy oznaczony
// jako trwały kazałby klientowi poddać się przy usterce, która minie sama.
func TestPorazkaKodemBierzePonawialnoscZKontraktu(t *testing.T) {
	for kod, ponawialny := range shared.KodyPonawialne {
		odpowiedz := PorazkaKodem(kod, "treść dla Operatora")
		if odpowiedz.Status != shared.EnvelopeStatusError {
			t.Errorf("odmowa kodem %q ma stan %q", kod, odpowiedz.Status)
		}
		if odpowiedz.Blad == nil {
			t.Fatalf("odmowa kodem %q nie niesie błędu", kod)
		}
		if odpowiedz.Blad.Retryable != ponawialny {
			t.Errorf("kod %q oznaczony jako ponawialny=%t, kontrakt mówi %t",
				kod, odpowiedz.Blad.Retryable, ponawialny)
		}
	}
}

// TestKopertaOdpowiedziWiazeOdpowiedzZWywolaniem sprawdza to, po czym klient
// koreluje odpowiedź: ten sam typ, ten sam identyfikator, ta sama sesja.
func TestKopertaOdpowiedziWiazeOdpowiedzZWywolaniem(t *testing.T) {
	sesja := "sesja-pierwsza"
	zadanie := Koperta{
		Type:      shared.CommandSessionCreate,
		Id:        "jeden",
		SessionId: &sesja,
	}
	wynik, err := Sukces(map[string]string{"sessionId": "sesja-nowa"})
	if err != nil {
		t.Fatalf("nie można złożyć odpowiedzi: %v", err)
	}

	odpowiedz := KopertaOdpowiedzi(zadanie, wynik)

	if odpowiedz.Type != zadanie.Type {
		t.Errorf("odpowiedź wraca pod typem %q zamiast %q", odpowiedz.Type, zadanie.Type)
	}
	if odpowiedz.Id != zadanie.Id {
		t.Errorf("odpowiedź wraca z identyfikatorem %q zamiast %q", odpowiedz.Id, zadanie.Id)
	}
	if IdSesji(odpowiedz) != sesja {
		t.Errorf("odpowiedź wraca z sesją %q zamiast %q", IdSesji(odpowiedz), sesja)
	}
	if odpowiedz.Status == nil || *odpowiedz.Status != shared.EnvelopeStatusOk {
		t.Error("odpowiedź udana bez stanu")
	}
	if odpowiedz.Timestamp == 0 {
		t.Error("odpowiedź bez znacznika czasu")
	}
}

// TestKopertaBleduNiesieStanIError sprawdza drugą drogę pakowania.
func TestKopertaBleduNiesieStanIError(t *testing.T) {
	zadanie := Koperta{Type: shared.CommandSessionCreate, Id: "jeden"}
	odpowiedz := KopertaBledu(zadanie, NowyBlad(shared.ErrorCodeValidationFailed, "brak tytułu"))

	if odpowiedz.Status == nil || *odpowiedz.Status != shared.EnvelopeStatusError {
		t.Error("koperta błędu bez stanu błędu")
	}
	if odpowiedz.Error == nil || odpowiedz.Error.Code != shared.ErrorCodeValidationFailed {
		t.Errorf("koperta błędu niesie %+v", odpowiedz.Error)
	}
	if odpowiedz.Payload != nil {
		t.Errorf("koperta błędu niesie ładunek: %s", odpowiedz.Payload)
	}
}

// TestBladPrzechodziPrzezWarstwyGoBezTlumaczenia sprawdza obieg zamknięty
// kodu: nadany raz, przeniesiony jako zwykły błąd Go i odczytany z powrotem
// bez podmiany na kod zastępczy.
func TestBladPrzechodziPrzezWarstwyGoBezTlumaczenia(t *testing.T) {
	zrodlowy := NowyBlad(shared.ErrorCodeConflict, "sesja jest już otwarta")
	przeniesiony := JakoError(zrodlowy)

	odczytany := BladZeZrodla(shared.ErrorCodeInternalError, przeniesiony)
	if odczytany.Code != shared.ErrorCodeConflict {
		t.Errorf("kod nadany w warstwie niższej zmieniono na %q", odczytany.Code)
	}
	if odczytany.Message != zrodlowy.Message {
		t.Errorf("treść błędu zmieniła się: %q", odczytany.Message)
	}
}

// TestBladOpakowanyNadalNiesieSwojKod sprawdza, że przeniesienie przez
// fmt.Errorf z %w nie gubi kodu kontraktu — tak wygląda droga błędu przez
// warstwy repozytoriów.
func TestBladOpakowanyNadalNiesieSwojKod(t *testing.T) {
	zrodlowy := NowyBlad(shared.ErrorCodePermissionDenied, "brak nadania")
	opakowany := errors.Join(errors.New("warstwa wyższa"), JakoError(zrodlowy))

	odczytany := BladZeZrodla(shared.ErrorCodeInternalError, opakowany)
	if odczytany.Code != shared.ErrorCodePermissionDenied {
		t.Errorf("kod zgubiony przy opakowaniu: %q", odczytany.Code)
	}
}

// TestBladObcyDostajeKodWskazany sprawdza drogę zwykłego błędu Go — takiego,
// któremu kodu nikt nie nadał.
func TestBladObcyDostajeKodWskazany(t *testing.T) {
	odczytany := BladZeZrodla(shared.ErrorCodeInternalError, errors.New("plik nie istnieje"))
	if odczytany.Code != shared.ErrorCodeInternalError {
		t.Errorf("błąd obcy dostał kod %q", odczytany.Code)
	}
	if odczytany.Message != "plik nie istnieje" {
		t.Errorf("treść błędu obcego zmieniona na %q", odczytany.Message)
	}
}

// TestBrakBleduDajeBladZerowy pilnuje, by nil nie zamienił się w odmowę.
func TestBrakBleduDajeBladZerowy(t *testing.T) {
	odczytany := BladZeZrodla(shared.ErrorCodeInternalError, nil)
	if odczytany.Code != "" || odczytany.Message != "" || odczytany.Retryable {
		t.Errorf("brak błędu zamieniony na %+v", odczytany)
	}
}

// TestOpisSkladaKodZTrescia sprawdza tekst, który idzie do dziennika.
func TestOpisSkladaKodZTrescia(t *testing.T) {
	opis := Opis(NowyBlad(shared.ErrorCodeNotFound, "sesji nie ma"))
	if opis != "not_found: sesji nie ma" {
		t.Errorf("opis błędu brzmi %q", opis)
	}
}
