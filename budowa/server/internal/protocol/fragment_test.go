package protocol

import (
	"context"
	"testing"

	"danacoconsole/shared"
)

// Strumień jest jedyną drogą, którą odpowiedź modelu dociera do okna, i jedyną,
// w której klient składa treść z kawałków. Trzy rzeczy muszą się w nim zgadzać
// zawsze: identyfikator równy identyfikatorowi żądania, numer rosnący od
// jedynki oraz jedno domknięcie na końcu. Rozjazd któregokolwiek zostawia okno
// w ładowaniu albo składa treść z dwóch tur naraz.

// TestFragmentWracaZKopertyBezZmiany sprawdza obieg zamknięty fragmentu.
func TestFragmentWracaZKopertyBezZmiany(t *testing.T) {
	zrodlowy := ChunkTekstu("okno-pierwsze", "wiadomosc-pierwsza", "Zaczynam pracę")

	koperta, err := KopertaFragmentu("jeden", "sesja-pierwsza", 1, false, zrodlowy)
	if err != nil {
		t.Fatalf("nie można spakować fragmentu: %v", err)
	}
	odczytany, err := FragmentZKoperty(koperta)
	if err != nil {
		t.Fatalf("nie można odczytać fragmentu: %v", err)
	}

	if odczytany.WindowId != zrodlowy.WindowId || odczytany.MessageId != zrodlowy.MessageId {
		t.Errorf("fragment zgubił przypisanie do okna albo wiadomości: %+v", odczytany)
	}
	if odczytany.Kind != shared.ChunkKindText {
		t.Errorf("rodzaj fragmentu zmienił się na %q", odczytany.Kind)
	}
	if Tresc(odczytany) != "Zaczynam pracę" {
		t.Errorf("treść fragmentu zmieniła się na %q", Tresc(odczytany))
	}
}

// TestKopertaFragmentuNiesieTypZdarzeniaKontraktu pilnuje, by strumień szedł
// pod nazwą, której klient nasłuchuje — a nie pod nazwą komendy, która turę
// otworzyła.
func TestKopertaFragmentuNiesieTypZdarzeniaKontraktu(t *testing.T) {
	koperta, err := KopertaFragmentu("jeden", "sesja-pierwsza", 1, false,
		ChunkTekstu("okno-pierwsze", "wiadomosc-pierwsza", "treść"))
	if err != nil {
		t.Fatalf("nie można spakować fragmentu: %v", err)
	}
	if koperta.Type != shared.EventStreamChunk {
		t.Errorf("fragment idzie pod typem %q, kontrakt wskazuje %q",
			koperta.Type, shared.EventStreamChunk)
	}
}

// TestNumerIDomkniecieStojaWKopercie sprawdza rozkład pól wskazany przez
// kontrakt: numer i znacznik końca należą do koperty, nie do ładunku. Fragment
// niedomykający nie ma pola `done` wcale — klient odróżnia w ten sposób
// „jeszcze nie koniec" od „koniec równy fałsz".
func TestNumerIDomkniecieStojaWKopercie(t *testing.T) {
	fragment := ChunkTekstu("okno-pierwsze", "wiadomosc-pierwsza", "treść")

	posrodku, err := KopertaFragmentu("jeden", "sesja-pierwsza", 2, false, fragment)
	if err != nil {
		t.Fatalf("nie można spakować fragmentu: %v", err)
	}
	if Numer(posrodku) != 2 {
		t.Errorf("numer fragmentu odczytany jako %d", Numer(posrodku))
	}
	if posrodku.Done != nil {
		t.Error("fragment niedomykający niesie pole done")
	}

	domykajacy, err := KopertaFragmentu("jeden", "sesja-pierwsza", 3, true, fragment)
	if err != nil {
		t.Fatalf("nie można spakować fragmentu domykającego: %v", err)
	}
	if !Ostatni(domykajacy) {
		t.Error("fragment domykający nie domyka strumienia")
	}
}

// TestStrumienNiesieJedenIdentyfikatorPrzezCalyBieg sprawdza obietnicę
// z kontraktu wprost: odpowiedź i wszystkie fragmenty powtarzają identyfikator
// żądania. Sprawdzian przechodzi całą turę — żądanie, trzy fragmenty,
// domknięcie — i porównuje identyfikatory oraz kolejność numerów.
func TestStrumienNiesieJedenIdentyfikatorPrzezCalyBieg(t *testing.T) {
	const idZadania = "jeden"
	zadanie := Koperta{Type: shared.CommandMessageSend, Id: idZadania}

	odpowiedz := KopertaOdpowiedzi(zadanie, Odpowiedz{Status: shared.EnvelopeStatusOk})
	if odpowiedz.Id != idZadania {
		t.Errorf("odpowiedź na komendę niesie identyfikator %q", odpowiedz.Id)
	}

	tresci := []string{"Za", "czy", "nam"}
	poprzedniNumer := 0
	domkniec := 0
	for i, tekst := range tresci {
		ostatni := i == len(tresci)-1
		koperta, err := KopertaFragmentu(idZadania, "sesja-pierwsza", i+1, ostatni,
			ChunkTekstu("okno-pierwsze", "wiadomosc-pierwsza", tekst))
		if err != nil {
			t.Fatalf("nie można spakować fragmentu %d: %v", i+1, err)
		}
		if koperta.Id != idZadania {
			t.Errorf("fragment %d niesie identyfikator %q zamiast %q", i+1, koperta.Id, idZadania)
		}
		if Numer(koperta) != poprzedniNumer+1 {
			t.Errorf("fragment %d ma numer %d, oczekiwany %d", i+1, Numer(koperta), poprzedniNumer+1)
		}
		poprzedniNumer = Numer(koperta)
		if Ostatni(koperta) {
			domkniec++
		}
	}
	if domkniec != 1 {
		t.Errorf("strumień domknięty %d razy, oczekiwane raz", domkniec)
	}
}

// TestFragmentBleduNiesiePrzyczyneDwiemaDrogami sprawdza kształt fragmentu
// błędu: przyczyna jedzie i jako treść nietekstowa dla klienta, i jako tekst
// dla Operatora. Brak jednej z dróg zostawiłby okno z pustym dymkiem.
func TestFragmentBleduNiesiePrzyczyneDwiemaDrogami(t *testing.T) {
	blad := NowyBlad(shared.ErrorCodeChannelUnavailable, "kanał modelu nie odpowiada")
	fragment := ChunkBledu("okno-pierwsze", "wiadomosc-pierwsza", blad)

	if fragment.Kind != shared.ChunkKindError {
		t.Errorf("fragment błędu ma rodzaj %q", fragment.Kind)
	}
	if Tresc(fragment) != blad.Message {
		t.Errorf("fragment błędu nie niesie treści dla Operatora: %q", Tresc(fragment))
	}
	if len(fragment.Data) == 0 {
		t.Fatal("fragment błędu nie niesie przyczyny maszynowej")
	}

	koperta, err := KopertaFragmentu("jeden", "sesja-pierwsza", 1, true, fragment)
	if err != nil {
		t.Fatalf("nie można spakować fragmentu błędu: %v", err)
	}
	odczytany, err := FragmentZKoperty(koperta)
	if err != nil {
		t.Fatalf("nie można odczytać fragmentu błędu: %v", err)
	}
	if odczytany.Kind != shared.ChunkKindError || Tresc(odczytany) != blad.Message {
		t.Errorf("fragment błędu zmienił się w obiegu: %+v", odczytany)
	}
}

// TestWersjaOstatecznaJestOsobnymRodzajem pilnuje rozróżnienia, na którym stoi
// składanie treści po stronie klienta: fragment ostateczny zastępuje treść
// złożoną z fragmentów tekstowych, więc nie może przyjść pod rodzajem `text`.
func TestWersjaOstatecznaJestOsobnymRodzajem(t *testing.T) {
	ostateczna := ChunkWersjiOstatecznej("okno-pierwsze", "wiadomosc-pierwsza", "Cała odpowiedź")
	if ostateczna.Kind != shared.ChunkKindFinal {
		t.Errorf("wersja ostateczna ma rodzaj %q, oczekiwany %q",
			ostateczna.Kind, shared.ChunkKindFinal)
	}
	if ostateczna.Kind == shared.ChunkKindText {
		t.Error("wersja ostateczna nie odróżnia się od porcji tekstu")
	}
	if Tresc(ostateczna) != "Cała odpowiedź" {
		t.Errorf("wersja ostateczna niesie %q", Tresc(ostateczna))
	}
}

// TestFragmentTrescNietekstowaNieDajeSieZakodowac sprawdza drogę odmowy przy
// budowie fragmentu dowolnego rodzaju.
func TestFragmentTrescNietekstowaNieDajeSieZakodowac(t *testing.T) {
	if _, err := NowyChunk(shared.ChunkKindImage, "okno-pierwsze", "wiadomosc-pierwsza",
		make(chan int)); err == nil {
		t.Error("treść niedająca się zakodować została przyjęta")
	}
}

// TestTozsamoscStrumieniaBierzeIdentyfikatorZadania sprawdza rozstrzygnięcie,
// które trzyma jeden identyfikator przez całą turę.
func TestTozsamoscStrumieniaBierzeIdentyfikatorZadania(t *testing.T) {
	zZadaniem := ZIdZadania(context.Background(), "jeden")
	if id := TozsamoscStrumienia(zZadaniem, "zastępcza"); id != "jeden" {
		t.Errorf("strumień tury otwartej żądaniem dostał tożsamość %q", id)
	}
}

// TestTozsamoscStrumieniaBezZadaniaBierzeZastepcza sprawdza turę powołaną poza
// drogą komendy — koordynator pętli, bieg naprawczy, podagent. Brak wpisu nie
// jest błędem i ma dawać tożsamość podaną przez wywołującego.
func TestTozsamoscStrumieniaBezZadaniaBierzeZastepcza(t *testing.T) {
	if id := TozsamoscStrumienia(context.Background(), "zastępcza"); id != "zastępcza" {
		t.Errorf("strumień bez żądania dostał tożsamość %q", id)
	}
}

// TestPustyIdentyfikatorNieZakladaWpisu pilnuje rozróżnienia „żądania nie było"
// od „żądanie miało puste pole". Wpis pusty kazałby warstwie wyżej uznać, że
// tożsamość jest, i wysłać strumień z pustym identyfikatorem.
func TestPustyIdentyfikatorNieZakladaWpisu(t *testing.T) {
	ctx := ZIdZadania(context.Background(), "")
	if id := TozsamoscStrumienia(ctx, "zastępcza"); id != "zastępcza" {
		t.Errorf("pusty identyfikator założył wpis: tożsamość %q", id)
	}
}
