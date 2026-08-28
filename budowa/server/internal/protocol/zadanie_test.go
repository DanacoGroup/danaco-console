package protocol

import (
	"encoding/json"
	"testing"

	"danacoconsole/shared"
)

// rejestrSprawdzianu buduje rozpoznanie z pełnego kontraktu — tak samo, jak
// robi to punkt wejścia rdzenia.
func rejestrSprawdzianu() *RejestrKomend {
	return NowyRejestrKomend(shared.WszystkieKomendy()...)
}

// TestZadanieRozpoznajeKomendeKontraktu sprawdza drogę zwykłą: nazwa
// z kontraktu przechodzi nietknięta i jest oznaczona jako znana.
func TestZadanieRozpoznajeKomendeKontraktu(t *testing.T) {
	zadanie := ZbudujZadanie(Koperta{
		Type: shared.CommandSessionCreate,
		Id:   "jeden",
	}, rejestrSprawdzianu())

	if !zadanie.Znana {
		t.Error("komenda kontraktu została uznana za nieznaną")
	}
	if zadanie.Komenda != shared.CommandSessionCreate {
		t.Errorf("komenda zmieniła się w rozpoznaniu: %q", zadanie.Komenda)
	}
	if zadanie.TypZadany != shared.CommandSessionCreate {
		t.Errorf("typ zadany nie został zachowany: %q", zadanie.TypZadany)
	}
}

// TestZadanieSpozaKontraktuDostajeZdarzenieObszaru pilnuje zachowania
// fail-open: nazwa nieznana nie przerywa przetwarzania, tylko zamienia się
// w zdarzenie `*.unknown` swojego obszaru.
func TestZadanieSpozaKontraktuDostajeZdarzenieObszaru(t *testing.T) {
	zadanie := ZbudujZadanie(Koperta{
		Type: "session.wymyslona",
		Id:   "jeden",
	}, rejestrSprawdzianu())

	if zadanie.Znana {
		t.Error("nazwa spoza kontraktu została uznana za znaną")
	}
	if oczekiwane := shared.ZdarzenieNieznanej("session.wymyslona"); zadanie.Komenda != oczekiwane {
		t.Errorf("nazwa nieznana skierowana na %q, kontrakt wskazuje %q", zadanie.Komenda, oczekiwane)
	}
	if zadanie.TypZadany != "session.wymyslona" {
		t.Errorf("typ zadany przez klienta został zgubiony: %q", zadanie.TypZadany)
	}
}

// TestObszarSpozaKontraktuTrafiaDoZdarzeniaPolaczenia sprawdza wyjście
// ostateczne: nazwa bez rozpoznawalnego obszaru idzie na zdarzenie połączenia.
func TestObszarSpozaKontraktuTrafiaDoZdarzeniaPolaczenia(t *testing.T) {
	zadanie := ZbudujZadanie(Koperta{Type: "obszarnieistniejacy.cokolwiek"}, rejestrSprawdzianu())
	if zadanie.Komenda != shared.EventConnectionUnknown {
		t.Errorf("nazwa z obszaru spoza kontraktu skierowana na %q, oczekiwane %q",
			zadanie.Komenda, shared.EventConnectionUnknown)
	}
}

// TestZasiegWychodziZLadunku sprawdza odczyt zasięgu na czterech poziomach.
// Okno jest poziomem najwęższym i wygrywa z pozostałymi w rozstrzyganiu
// konfiguracji, więc jego odczyt musi być pewny.
func TestZasiegWychodziZLadunku(t *testing.T) {
	ladunek := json.RawMessage(`{
		"environmentId":"talkin",
		"projectId":"oferta",
		"sessionId":"sesja-z-ladunku",
		"windowId":"okno-pierwsze"
	}`)
	zadanie := ZbudujZadanie(Koperta{
		Type:    shared.CommandConfigGet,
		Payload: ladunek,
	}, rejestrSprawdzianu())

	if zadanie.Zasieg.Srodowisko != "talkin" {
		t.Errorf("środowisko odczytane jako %q", zadanie.Zasieg.Srodowisko)
	}
	if zadanie.Zasieg.Projekt != "oferta" {
		t.Errorf("projekt odczytany jako %q", zadanie.Zasieg.Projekt)
	}
	if zadanie.Zasieg.Sesja != "sesja-z-ladunku" {
		t.Errorf("sesja odczytana jako %q", zadanie.Zasieg.Sesja)
	}
	if zadanie.Zasieg.Okno != "okno-pierwsze" {
		t.Errorf("okno odczytane jako %q", zadanie.Zasieg.Okno)
	}
}

// TestSesjaZKopertyUzupelniaBrakWLadunku pilnuje pierwszeństwa: ładunek
// rozstrzyga, a koperta uzupełnia. Odwrócenie tej kolejności kazałoby komendzie
// pracować na innej sesji, niż wskazuje jej własna treść.
func TestSesjaZKopertyUzupelniaBrakWLadunku(t *testing.T) {
	zKoperty := "sesja-z-koperty"

	bezSesjiWLadunku := ZbudujZadanie(Koperta{
		Type:      shared.CommandSessionFocus,
		SessionId: &zKoperty,
		Payload:   json.RawMessage(`{"windowId":"okno-pierwsze"}`),
	}, rejestrSprawdzianu())
	if bezSesjiWLadunku.Zasieg.Sesja != zKoperty {
		t.Errorf("brak sesji w ładunku nie został uzupełniony z koperty: %q",
			bezSesjiWLadunku.Zasieg.Sesja)
	}

	zSesjaWLadunku := ZbudujZadanie(Koperta{
		Type:      shared.CommandSessionFocus,
		SessionId: &zKoperty,
		Payload:   json.RawMessage(`{"sessionId":"sesja-z-ladunku"}`),
	}, rejestrSprawdzianu())
	if zSesjaWLadunku.Zasieg.Sesja != "sesja-z-ladunku" {
		t.Errorf("sesja z koperty przykryła sesję z ładunku: %q", zSesjaWLadunku.Zasieg.Sesja)
	}
}

// TestLadunekInnegoKsztaltuNieWywracaZasiegu sprawdza, że ładunek, który nie
// jest obiektem, daje zasięg pusty zamiast przerwać rozpoznanie.
func TestLadunekInnegoKsztaltuNieWywracaZasiegu(t *testing.T) {
	zadanie := ZbudujZadanie(Koperta{
		Type:    shared.CommandSessionList,
		Payload: json.RawMessage(`[1,2,3]`),
	}, rejestrSprawdzianu())

	if zadanie.Zasieg != (Zasieg{}) {
		t.Errorf("ładunek innego kształtu dał zasięg %+v, oczekiwany pusty", zadanie.Zasieg)
	}
}

// TestOdpowiedzNieznanejNiesieStanIKod sprawdza kształt odmowy, na którym
// klient opiera rozstrzygnięcie obietnicy wywołania. Koperta bez pola `status`
// zostawiłaby okno w wiecznym ładowaniu.
func TestOdpowiedzNieznanejNiesieStanIKod(t *testing.T) {
	zadanie := ZbudujZadanie(Koperta{
		Type: "session.wymyslona",
		Id:   "jeden",
	}, rejestrSprawdzianu())

	odpowiedz := OdpowiedzNieznanej(zadanie)

	if odpowiedz.Id != "jeden" {
		t.Errorf("odmowa zgubiła identyfikator żądania: %q", odpowiedz.Id)
	}
	if odpowiedz.Status == nil || *odpowiedz.Status != shared.EnvelopeStatusError {
		t.Error("odmowa bez stanu błędu")
	}
	if odpowiedz.Error == nil {
		t.Fatal("odmowa bez pola error")
	}
	if odpowiedz.Error.Code != shared.ErrorCodeNotFound {
		t.Errorf("odmowa niesie kod %q, oczekiwany %q", odpowiedz.Error.Code, shared.ErrorCodeNotFound)
	}

	var tresc shared.UnknownCommandPayload
	if err := LadunekDo(odpowiedz, &tresc); err != nil {
		t.Fatalf("nieczytelny ładunek odmowy: %v", err)
	}
	if tresc.RequestedType != "session.wymyslona" {
		t.Errorf("odmowa nie mówi, czego klient zażądał: %q", tresc.RequestedType)
	}
	if tresc.RequestId == nil || *tresc.RequestId != "jeden" {
		t.Error("odmowa nie wiąże się z identyfikatorem żądania")
	}
}

// TestOdkodujZadanieOddzielaDwaNiepowodzenia sprawdza wejście warstwy
// transportu: błąd oznacza wyłącznie komunikat niepoprawny strukturalnie,
// a nieznana komenda błędem nie jest.
func TestOdkodujZadanieOddzielaDwaNiepowodzenia(t *testing.T) {
	if _, err := OdkodujZadanie([]byte("to nie jest koperta"), rejestrSprawdzianu()); err == nil {
		t.Error("bajty niebędące kopertą zostały przyjęte")
	}

	zadanie, err := OdkodujZadanie(
		[]byte(`{"type":"session.wymyslona","id":"jeden"}`), rejestrSprawdzianu())
	if err != nil {
		t.Fatalf("nieznana komenda zgłoszona jako błąd dekodowania: %v", err)
	}
	if zadanie.Znana {
		t.Error("nieznana komenda została uznana za znaną")
	}
}

// TestKopertaZadaniaOdtwarzaWywolanie sprawdza drogę powrotną: żądanie
// rozłożone na części składa się z powrotem w kopertę, z której klient
// rozpozna swoje wywołanie.
func TestKopertaZadaniaOdtwarzaWywolanie(t *testing.T) {
	zrodlowa := Koperta{
		Type:      shared.CommandSessionCreate,
		Id:        "jeden",
		Payload:   json.RawMessage(`{"sessionId":"sesja-pierwsza"}`),
		Timestamp: 1234,
	}
	odtworzona := ZbudujZadanie(zrodlowa, rejestrSprawdzianu()).Koperta()

	if odtworzona.Type != zrodlowa.Type || odtworzona.Id != zrodlowa.Id {
		t.Errorf("odtworzona koperta rozeszła się ze źródłem: %+v", odtworzona)
	}
	if odtworzona.Timestamp != zrodlowa.Timestamp {
		t.Errorf("znacznik czasu zmienił się przy odtworzeniu: %d", odtworzona.Timestamp)
	}
	if IdSesji(odtworzona) != "sesja-pierwsza" {
		t.Errorf("sesja zgubiona przy odtworzeniu: %q", IdSesji(odtworzona))
	}
}

// TestRejestrPustyNieZnaNiczego sprawdza zachowanie odbiornika zerowego
// i rejestru pustego — obie drogi mają kierować na `*.unknown`, a nie wywracać
// rozpoznania.
func TestRejestrPustyNieZnaNiczego(t *testing.T) {
	var zerowy *RejestrKomend
	if _, znana := zerowy.Rozpoznaj(shared.CommandConnectionHello); znana {
		t.Error("rejestr zerowy uznał komendę za znaną")
	}
	if zerowy.Liczba() != 0 {
		t.Errorf("rejestr zerowy podaje rozmiar %d", zerowy.Liczba())
	}

	pusty := NowyRejestrKomend()
	if _, znana := pusty.Rozpoznaj(shared.CommandConnectionHello); znana {
		t.Error("rejestr pusty uznał komendę za znaną")
	}
}

// TestRejestrPomijaNazwyPuste pilnuje, by pusty napis nie wszedł do zbioru
// nazw znanych — koperta bez typu trafiłaby wtedy na ścieżkę komendy znanej.
func TestRejestrPomijaNazwyPuste(t *testing.T) {
	rejestr := NowyRejestrKomend("", shared.CommandConnectionHello, "")
	if rejestr.Liczba() != 1 {
		t.Errorf("rejestr przyjął nazwy puste: rozmiar %d", rejestr.Liczba())
	}
	if _, znana := rejestr.Rozpoznaj(""); znana {
		t.Error("pusty typ został uznany za komendę znaną")
	}
}
