// Skutek pól dokumentu: czy pole wchodzi policzone, czy pole obliczane liczy
// naprawdę i czy pole bez czego policzyć mówi to wprost.
package core

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"danacoconsole/shared"
)

// poleUprzazSprawdzianu składa adapter modułu nad bazą sprawdzianu, gotowy do
// wywołania komend odczytu i wstawiania pól dokumentu.
func poleUprzazSprawdzianu(t *testing.T) (*adapterStudia, context.Context, string) {
	t.Helper()
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)
	dokument := dokumentZTrescia(t, zmontowany, zycie, "okno-pola",
		"Pismo w sprawie.\nTreść pisma.")
	return nowyAdapterStudia(zmontowany.dane.Studio), zycie, dokument.Id
}

// poleZWykazu odczytuje pole osobnym wywołaniem wykazu, niezależnym od
// odpowiedzi wstawienia, żeby sprawdzian mierzył stan zapisany.
func poleZWykazu(t *testing.T, adapter *adapterStudia, zycie context.Context,
	dokument, pole string) shared.StudioDocumentField {
	t.Helper()

	wykaz, err := adapter.WykazPol(zycie, shared.StudioFieldListRequest{DocumentId: dokument})
	if err != nil {
		t.Fatalf("wykaz pól odmówił: %v", err)
	}
	for _, pozycja := range wykaz.Fields {
		if pozycja.Id == pole {
			return pozycja
		}
	}
	t.Fatalf("wykaz pól nie niesie pola %s", pole)
	return shared.StudioDocumentField{}
}

// TestPoleLiczbyStronWchodziPoliczone mierzy, że pole nie wchodzi puste,
// tylko z wartością rzeczywiście policzoną przez rdzeń.
func TestPoleLiczbyStronWchodziPoliczone(t *testing.T) {
	adapter, zycie, dokument := poleUprzazSprawdzianu(t)

	wstawione, err := adapter.WstawPole(zycie, shared.StudioFieldInsertRequest{
		DocumentId: dokument, Offset: 0, Kind: shared.StudioFieldKindPageCount,
	})
	if err != nil {
		t.Fatalf("wstawienie pola liczby stron odmówiło: %v", err)
	}
	pole := poleZWykazu(t, adapter, zycie, dokument, wstawione.Field.Id)
	if pole.Value == nil || *pole.Value != "1" {
		t.Errorf("pole liczby stron niesie wartość %v, oczekiwano 1", pole.Value)
	}
	if pole.Stale == nil || *pole.Stale {
		t.Error("pole policzone przy wstawieniu stoi jako wymagające odświeżenia")
	}
}

// TestPoleNumeruStronyIWlasciwosciLiczaSie mierzy numer strony i właściwość
// dokumentu — obie wartości mają być policzone przez rdzeń, nie zgadnięte.
func TestPoleNumeruStronyIWlasciwosciLiczaSie(t *testing.T) {
	adapter, zycie, dokument := poleUprzazSprawdzianu(t)

	numer, err := adapter.WstawPole(zycie, shared.StudioFieldInsertRequest{
		DocumentId: dokument, Offset: 3, Kind: shared.StudioFieldKindPageNumber,
	})
	if err != nil {
		t.Fatalf("wstawienie pola numeru strony odmówiło: %v", err)
	}
	if numer.Field.Value == nil || *numer.Field.Value != "1" {
		t.Errorf("pole numeru strony niesie %v, oczekiwano 1", numer.Field.Value)
	}

	znaki, err := adapter.WstawPole(zycie, shared.StudioFieldInsertRequest{
		DocumentId: dokument, Offset: 0, Kind: shared.StudioFieldKindDocumentProperty,
		PropertyName: wskaznik("liczba znaków"),
	})
	if err != nil {
		t.Fatalf("wstawienie pola właściwości odmówiło: %v", err)
	}
	pole := poleZWykazu(t, adapter, zycie, dokument, znaki.Field.Id)
	if pole.Value == nil || *pole.Value == "" || *pole.Value == "0" {
		t.Errorf("pole liczby znaków niesie %v — wartość ma być policzona z treści",
			pole.Value)
	}
}

// TestPoleObliczaneSumujeKolumneTabeli mierzy sedno pola obliczanego: rachunek
// stoi w rdzeniu i liczy na zawartości tabeli dokumentu.
func TestPoleObliczaneSumujeKolumneTabeli(t *testing.T) {
	adapter, zycie, dokument := poleUprzazSprawdzianu(t)

	komorki, err := json.Marshal([]shared.StudioTableCell{
		{Row: 0, Column: 0, Text: wskaznik("Pozycja")},
		{Row: 0, Column: 1, Text: wskaznik("Kwota")},
		{Row: 1, Column: 0, Text: wskaznik("pierwsza")},
		{Row: 1, Column: 1, Text: wskaznik("120,50")},
		{Row: 2, Column: 0, Text: wskaznik("druga")},
		{Row: 2, Column: 1, Text: wskaznik("9,50")},
	})
	if err != nil {
		t.Fatalf("nie można złożyć komórek żądania: %v", err)
	}
	tabela, err := adapter.WstawTabele(zycie, shared.StudioTableInsertRequest{
		DocumentId: dokument, Offset: 0, Rows: 3, Columns: 2, Cells: komorki,
		HeaderRows: wskaznik(1),
	})
	if err != nil {
		t.Fatalf("wstawienie tabeli odmówiło: %v", err)
	}

	wstawione, err := adapter.WstawPole(zycie, shared.StudioFieldInsertRequest{
		DocumentId: dokument, Offset: 0, Kind: shared.StudioFieldKindCalculated,
		Expression: wskaznik("SUMA(" + tabela.Table.Id + "; 1) * 2"),
	})
	if err != nil {
		t.Fatalf("wstawienie pola obliczanego odmówiło: %v", err)
	}
	pole := poleZWykazu(t, adapter, zycie, dokument, wstawione.Field.Id)
	if pole.Value == nil || *pole.Value != "260" {
		t.Errorf("pole obliczane niesie %v, oczekiwano 260 (suma 130 razy dwa)", pole.Value)
	}

	// Wiersz nagłówkowy nie wchodzi do rachunku, a odświeżenie liczy od nowa.
	if _, err := adapter.OdswiezPola(zycie, shared.StudioFieldRefreshRequest{
		DocumentId: dokument, FieldId: &wstawione.Field.Id,
	}); err != nil {
		t.Fatalf("odświeżenie pola odmówiło: %v", err)
	}
	poOdswiezeniu := poleZWykazu(t, adapter, zycie, dokument, wstawione.Field.Id)
	if poOdswiezeniu.Value == nil || *poOdswiezeniu.Value != "260" {
		t.Errorf("pole obliczane po odświeżeniu niesie %v, oczekiwano 260",
			poOdswiezeniu.Value)
	}
}

// TestPoleObliczaneNazwaneOdmowyPrzyBledzie pilnoje, że rachunek nie udaje
// wyniku: wyrażenie, którego nie zna, wraca odmową nazywającą miejsce.
func TestPoleObliczaneNazwaneOdmowyPrzyBledzie(t *testing.T) {
	adapter, zycie, dokument := poleUprzazSprawdzianu(t)

	if _, err := adapter.WstawPole(zycie, shared.StudioFieldInsertRequest{
		DocumentId: dokument, Offset: 0, Kind: shared.StudioFieldKindCalculated,
	}); err == nil {
		t.Error("pole obliczane bez wyrażenia wróciło bez odmowy")
	}

	wstawione, err := adapter.WstawPole(zycie, shared.StudioFieldInsertRequest{
		DocumentId: dokument, Offset: 0, Kind: shared.StudioFieldKindCalculated,
		Expression: wskaznik("2 + nieznane_dzialanie"),
	})
	if err != nil {
		t.Fatalf("wstawienie pola z wyrażeniem błędnym odmówiło całą czynnością: %v", err)
	}
	if wstawione.Balance.SkippedCount == 0 {
		t.Error("pole z wyrażeniem, którego rachunek nie zna, weszło bez ani jednego " +
			"pominięcia w bilansie")
	}
	pole := poleZWykazu(t, adapter, zycie, dokument, wstawione.Field.Id)
	if pole.Stale == nil || !*pole.Stale {
		t.Error("pole bez policzonej wartości nie stoi jako wymagające odświeżenia")
	}
	if pole.Value != nil && *pole.Value != "" {
		t.Errorf("pole bez policzonej wartości niesie %q — udawana wartość jest gorsza "+
			"niż brak", *pole.Value)
	}
}

// TestPoleObliczaneDzieliPrzezZeroOdmawia pilnuje granicy rachunku: dzielenie
// przez zero ma się skończyć odmową, nie wynikiem zmyślonym.
func TestPoleObliczaneDzieliPrzezZeroOdmawia(t *testing.T) {
	if _, err := poleObliczWyrazenie("4 / 0", nil); err == nil {
		t.Error("dzielenie przez zero w polu obliczanym wróciło bez odmowy")
	}
	wynik, err := poleObliczWyrazenie("(2 + 3,5) * 2", nil)
	if err != nil {
		t.Fatalf("rachunek wyrażenia z nawiasem i przecinkiem dziesiętnym odmówił: %v", err)
	}
	if wynik != 11 {
		t.Errorf("rachunek dał %v, oczekiwano 11", wynik)
	}
}

// TestPoleWlasciwosciNieznanejMowiWykaz pilnuje uczciwości: rdzeń nazywa, czego
// nie zna, i wymienia, o co wolno pytać.
func TestPoleWlasciwosciNieznanejMowiWykaz(t *testing.T) {
	adapter, zycie, dokument := poleUprzazSprawdzianu(t)

	wstawione, err := adapter.WstawPole(zycie, shared.StudioFieldInsertRequest{
		DocumentId: dokument, Offset: 0, Kind: shared.StudioFieldKindDocumentProperty,
		PropertyName: wskaznik("numer sprawy w sądzie"),
	})
	if err != nil {
		t.Fatalf("wstawienie pola właściwości odmówiło całą czynnością: %v", err)
	}
	if wstawione.Balance.SkippedCount == 0 {
		t.Fatal("pole właściwości nieznanej weszło bez pominięcia w bilansie")
	}
	szczegol := ""
	if wstawione.Balance.Skipped[0].Detail != nil {
		szczegol = *wstawione.Balance.Skipped[0].Detail
	}
	if !strings.Contains(szczegol, "liczba stron") {
		t.Errorf("pominięcie nie wymienia właściwości, o które wolno pytać: %q", szczegol)
	}
}

// TestPoleDatyBierzeWzorOperatora mierzy, że wzór daty jest pisany znakami
// z pakietu biurowego, a nie układem odniesienia biblioteki.
func TestPoleDatyBierzeWzorOperatora(t *testing.T) {
	adapter, zycie, dokument := poleUprzazSprawdzianu(t)

	wstawione, err := adapter.WstawPole(zycie, shared.StudioFieldInsertRequest{
		DocumentId: dokument, Offset: 0, Kind: shared.StudioFieldKindDate,
		Format: wskaznik("yyyy-MM-dd"),
	})
	if err != nil {
		t.Fatalf("wstawienie pola daty odmówiło: %v", err)
	}
	pole := poleZWykazu(t, adapter, zycie, dokument, wstawione.Field.Id)
	if pole.Value == nil {
		t.Fatal("pole daty nie niesie wartości")
	}
	zapis := *pole.Value
	if len(zapis) != 10 || zapis[4] != '-' || zapis[7] != '-' {
		t.Errorf("pole daty ze wzorem yyyy-MM-dd niesie %q", zapis)
	}
}

// TestPoleRodzajuNieznanegoOdmawia pilnuje granicy kontraktu: rodzaj pola
// spoza wyliczenia ma dostać odmowę, a nie wejść jako pole zwykłe.
func TestPoleRodzajuNieznanegoOdmawia(t *testing.T) {
	adapter, zycie, dokument := poleUprzazSprawdzianu(t)

	if _, err := adapter.WstawPole(zycie, shared.StudioFieldInsertRequest{
		DocumentId: dokument, Offset: 0, Kind: shared.StudioFieldKind("numerFaktury"),
	}); err == nil {
		t.Error("pole rodzaju, którego kontrakt nie zna, wróciło bez odmowy")
	}
}

// TestOdswiezeniePolBezPolOdmawia pilnuje, że odświeżenie niczego nie wraca
// jako wykonane, gdy dokument nie niesie ani jednego pola obliczanego.
func TestOdswiezeniePolBezPolOdmawia(t *testing.T) {
	adapter, zycie, dokument := poleUprzazSprawdzianu(t)

	if _, err := adapter.OdswiezPola(zycie, shared.StudioFieldRefreshRequest{
		DocumentId: dokument,
	}); err == nil {
		t.Error("odświeżenie pól dokumentu bez ani jednego pola wróciło bez odmowy")
	}
}
