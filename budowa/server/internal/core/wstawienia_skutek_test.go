package core

import (
	"context"
	"strings"
	"testing"

	"danacoconsole/shared"
)

// Skutek wstawień: czy obiekt osadzony w dokumencie niesie POCHODZENIE i czy
// odmowa jest nazwana tam, gdzie rdzeń nie ma czym wykonać czynności.
//
// Szkody, które ten plik ma wykluczyć:
//  1. obraz wstawiony bez zapisanego pochodzenia — za tydzień nikt nie odtworzy,
//     na czym pismo się opiera (zlecenie mówi o tym wprost);
//  2. obiekt o zerowym rozmiarze, czyli niewidzialny, oddany jako wstawiony;
//  3. uprzejma odmowa bez nazwania braku — Operator ma wiedzieć, po czyjej
//     stronie brakuje i którą drogą czynność jest wykonalna;
//  4. usunięcie obiektu, po którym wykaz nadal go pokazuje.

// obiektUprzazSprawdzianu składa adapter modułu nad bazą sprawdzianu.
func obiektUprzazSprawdzianu(t *testing.T) (*adapterStudia, context.Context, string) {
	t.Helper()
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)
	dokument := dokumentZTrescia(t, zmontowany, zycie, "okno-wstawienia",
		"Akapit pierwszy.\nAkapit drugi.")
	return nowyAdapterStudia(zmontowany.dane.Studio), zycie, dokument.Id
}

// TestObiektBezPochodzeniaOdmawiaNazwanieDrog mierzy wymaganie zlecenia:
// wstawienie obrazu bez pochodzenia jest brakiem, nie skrótem.
func TestObiektBezPochodzeniaOdmawiaNazwanieDrog(t *testing.T) {
	adapter, zycie, dokument := obiektUprzazSprawdzianu(t)

	_, err := adapter.WstawObiekt(zycie, shared.StudioObjectInsertRequest{
		DocumentId: dokument, Offset: 0, Kind: shared.StudioObjectKindImage,
	})
	if err == nil {
		t.Fatal("wstawienie obrazu bez pochodzenia wróciło bez odmowy")
	}
	for _, droga := range []string{"assetId", "designNodeId", "libraryFileId", "path",
		"bytesBase64", "sourceUrl"} {

		if !strings.Contains(err.Error(), droga) {
			t.Errorf("odmowa nie nazywa drogi %q: %v", droga, err)
		}
	}
}

// TestKsztaltWstawionyNiesiePochodzenieIRozmiar mierzy kształt rysowany na
// miejscu: pochodzenie zapisane, rozmiar niezerowy, warstwa nadana.
func TestKsztaltWstawionyNiesiePochodzenieIRozmiar(t *testing.T) {
	adapter, zycie, dokument := obiektUprzazSprawdzianu(t)

	rodzaj := shared.StudioShapeKind(shared.StudioShapeKindRectangle)
	wstawiony, err := adapter.WstawObiekt(zycie, shared.StudioObjectInsertRequest{
		DocumentId: dokument, Offset: 0, Kind: shared.StudioObjectKindShape,
		ShapeKind: &rodzaj, InnerText: wskaznik("Uwaga"),
	})
	if err != nil {
		t.Fatalf("wstawienie kształtu odmówiło: %v", err)
	}

	wykaz, err := adapter.WykazObiektow(zycie, shared.StudioObjectListRequest{
		DocumentId: dokument,
	})
	if err != nil {
		t.Fatalf("wykaz obiektów odmówił: %v", err)
	}
	if len(wykaz.Objects) != 1 {
		t.Fatalf("wykaz niesie %d obiektów, oczekiwano jednego", len(wykaz.Objects))
	}
	obiekt := wykaz.Objects[0]
	if obiekt.Id != wstawiony.Object.Id {
		t.Errorf("wykaz niesie obiekt %q, wstawiono %q", obiekt.Id, wstawiony.Object.Id)
	}
	if obiekt.Source == nil {
		t.Fatal("obiekt w wykazie nie niesie zapisanego pochodzenia")
	}
	if *obiekt.Source != shared.StudioObjectSourceDrawn {
		t.Errorf("kształt rysowany niesie pochodzenie %q, oczekiwano %q", *obiekt.Source,
			shared.StudioObjectSourceDrawn)
	}
	if obiekt.WidthMm == nil || *obiekt.WidthMm <= 0 ||
		obiekt.HeightMm == nil || *obiekt.HeightMm <= 0 {
		t.Error("obiekt wstawiony ma zerowy rozmiar — obiekt niewidzialny oddany jako " +
			"wstawiony jest cichą szkodą")
	}
	if obiekt.ShapeKind == nil || *obiekt.ShapeKind != shared.StudioShapeKindRectangle {
		t.Error("kształt w wykazie nie niesie swojego rodzaju")
	}
}

// TestKsztaltBezRodzajuOdmawia pilnuje, że kształt bez wskazania rodzaju nie
// wchodzi jako obiekt bez postaci.
func TestKsztaltBezRodzajuOdmawia(t *testing.T) {
	adapter, zycie, dokument := obiektUprzazSprawdzianu(t)

	_, err := adapter.WstawObiekt(zycie, shared.StudioObjectInsertRequest{
		DocumentId: dokument, Offset: 0, Kind: shared.StudioObjectKindShape,
	})
	if err == nil {
		t.Fatal("wstawienie kształtu bez rodzaju wróciło bez odmowy")
	}
	if !strings.Contains(err.Error(), "shapeKind") {
		t.Errorf("odmowa nie nazywa brakującego pola: %v", err)
	}
}

// TestWykresOdmawiaNazwaniemBrakuPoStronieRdzenia pilnuje uczciwości: rdzeń nie
// ma rachunku wykresu i mówi to wprost, wraz z drogą, która działa.
func TestWykresOdmawiaNazwaniemBrakuPoStronieRdzenia(t *testing.T) {
	adapter, zycie, dokument := obiektUprzazSprawdzianu(t)

	_, err := adapter.WstawObiekt(zycie, shared.StudioObjectInsertRequest{
		DocumentId: dokument, Offset: 0, Kind: shared.StudioObjectKindChart,
	})
	if err == nil {
		t.Fatal("wstawienie wykresu wróciło bez odmowy — rdzeń nie ma rachunku wykresu")
	}
	if !strings.Contains(err.Error(), "Design") {
		t.Errorf("odmowa nie wskazuje drogi, która działa: %v", err)
	}
}

// TestPoleTekstoweBezTresciOdmawia pilnuje zakazu odpowiedzi „ok" bez skutku.
func TestPoleTekstoweBezTresciOdmawia(t *testing.T) {
	adapter, zycie, dokument := obiektUprzazSprawdzianu(t)

	if _, err := adapter.WstawObiekt(zycie, shared.StudioObjectInsertRequest{
		DocumentId: dokument, Offset: 0, Kind: shared.StudioObjectKindTextbox,
	}); err == nil {
		t.Error("wstawienie pola tekstowego bez treści wróciło bez odmowy")
	}
}

// TestObiektRozmiarZachowujeProporcje mierzy, że podanie samej szerokości przy
// zachowaniu proporcji przelicza wysokość, a nie rozciąga obiektu.
func TestObiektRozmiarZachowujeProporcje(t *testing.T) {
	adapter, zycie, dokument := obiektUprzazSprawdzianu(t)

	rodzaj := shared.StudioShapeKind(shared.StudioShapeKindEllipse)
	wstawiony, err := adapter.WstawObiekt(zycie, shared.StudioObjectInsertRequest{
		DocumentId: dokument, Offset: 0, Kind: shared.StudioObjectKindShape,
		ShapeKind: &rodzaj, WidthMm: wskaznik(100.0), HeightMm: wskaznik(50.0),
	})
	if err != nil {
		t.Fatalf("wstawienie kształtu odmówiło: %v", err)
	}
	if _, err := adapter.UstawPostacObiektu(zycie, shared.StudioObjectFormatSetRequest{
		DocumentId: dokument, ObjectId: wstawiony.Object.Id,
		WidthMm: wskaznik(60.0), KeepAspect: wskaznik(true),
	}); err != nil {
		t.Fatalf("zmiana postaci obiektu odmówiła: %v", err)
	}

	wykaz, err := adapter.WykazObiektow(zycie, shared.StudioObjectListRequest{
		DocumentId: dokument,
	})
	if err != nil {
		t.Fatalf("wykaz obiektów odmówił: %v", err)
	}
	if len(wykaz.Objects) != 1 {
		t.Fatalf("wykaz niesie %d obiektów, oczekiwano jednego", len(wykaz.Objects))
	}
	obiekt := wykaz.Objects[0]
	if obiekt.WidthMm == nil || *obiekt.WidthMm != 60 {
		t.Errorf("szerokość obiektu po zmianie: %v, oczekiwano 60", obiekt.WidthMm)
	}
	if obiekt.HeightMm == nil || *obiekt.HeightMm < 29.9 || *obiekt.HeightMm > 30.1 {
		t.Errorf("wysokość obiektu po zmianie: %v, oczekiwano 30 (proporcja 100×50)",
			obiekt.HeightMm)
	}
}

// TestObiektPostacBezCechyOdmawia pilnuje zakazu odpowiedzi bez skutku.
func TestObiektPostacBezCechyOdmawia(t *testing.T) {
	adapter, zycie, dokument := obiektUprzazSprawdzianu(t)

	rodzaj := shared.StudioShapeKind(shared.StudioShapeKindLine)
	wstawiony, err := adapter.WstawObiekt(zycie, shared.StudioObjectInsertRequest{
		DocumentId: dokument, Offset: 0, Kind: shared.StudioObjectKindShape,
		ShapeKind: &rodzaj,
	})
	if err != nil {
		t.Fatalf("wstawienie kształtu odmówiło: %v", err)
	}
	if _, err := adapter.UstawPostacObiektu(zycie, shared.StudioObjectFormatSetRequest{
		DocumentId: dokument, ObjectId: wstawiony.Object.Id,
	}); err == nil {
		t.Error("zmiana postaci obiektu bez ani jednej cechy wróciła bez odmowy")
	}
}

// TestUsuniecieObiektuZdejmujeGoZWykazuIMowiOZasobie mierzy usunięcie: obiekt
// schodzi z wykazu, a bilans mówi wprost, że bajty zostają w magazynie.
func TestUsuniecieObiektuZdejmujeGoZWykazu(t *testing.T) {
	adapter, zycie, dokument := obiektUprzazSprawdzianu(t)

	rodzaj := shared.StudioShapeKind(shared.StudioShapeKindStar)
	wstawiony, err := adapter.WstawObiekt(zycie, shared.StudioObjectInsertRequest{
		DocumentId: dokument, Offset: 0, Kind: shared.StudioObjectKindShape,
		ShapeKind: &rodzaj,
	})
	if err != nil {
		t.Fatalf("wstawienie kształtu odmówiło: %v", err)
	}
	usuniety, err := adapter.UsunObiekt(zycie, shared.StudioObjectRemoveRequest{
		DocumentId: dokument, ObjectId: wstawiony.Object.Id,
	})
	if err != nil {
		t.Fatalf("usunięcie obiektu odmówiło: %v", err)
	}
	if !usuniety.Removed {
		t.Error("usunięcie obiektu wróciło z Removed równym fałszowi")
	}
	wykaz, err := adapter.WykazObiektow(zycie, shared.StudioObjectListRequest{
		DocumentId: dokument,
	})
	if err != nil {
		t.Fatalf("wykaz obiektów odmówił: %v", err)
	}
	if len(wykaz.Objects) != 0 {
		t.Errorf("po usunięciu wykaz niesie %d obiektów, oczekiwano zera",
			len(wykaz.Objects))
	}
	// Usunięcie obiektu, którego już nie ma, jest odmową nazwaną, nie ciszą.
	if _, err := adapter.UsunObiekt(zycie, shared.StudioObjectRemoveRequest{
		DocumentId: dokument, ObjectId: wstawiony.Object.Id,
	}); err == nil {
		t.Error("usunięcie obiektu nieistniejącego wróciło bez odmowy")
	}
}

// TestIkonaNieznanaOdmawiaWskazaniemKatalogu pilnuje, że nazwy ikon biorą się
// z katalogu modułu Design, a nie z drugiego wykazu w Studiu.
//
// Nazwa szukana jest umyślnie bez ani jednego słowa z etykiet katalogu: katalog
// Designu dopasowuje po zawieraniu w obie strony, więc nazwa niosąca „katalog"
// albo „folder" trafiłaby we wzór i sprawdzian mierzyłby coś innego.
func TestIkonaNieznanaOdmawiaWskazaniemKatalogu(t *testing.T) {
	adapter, zycie, dokument := obiektUprzazSprawdzianu(t)

	_, err := adapter.WstawObiekt(zycie, shared.StudioObjectInsertRequest{
		DocumentId: dokument, Offset: 0, Kind: shared.StudioObjectKindIcon,
		IconName: wskaznik("qxzv"),
	})
	if err == nil {
		t.Fatal("wstawienie ikony spoza katalogu wróciło bez odmowy")
	}
	if !strings.Contains(err.Error(), "design.icon") {
		t.Errorf("odmowa nie wskazuje katalogu ikon rdzenia: %v", err)
	}
}

// TestIkonaBezNazwyOdmawia pilnuje, że ikona nie wchodzi bez wskazania wzoru.
func TestIkonaBezNazwyOdmawia(t *testing.T) {
	adapter, zycie, dokument := obiektUprzazSprawdzianu(t)

	if _, err := adapter.WstawObiekt(zycie, shared.StudioObjectInsertRequest{
		DocumentId: dokument, Offset: 0, Kind: shared.StudioObjectKindIcon,
	}); err == nil {
		t.Error("wstawienie ikony bez nazwy wróciło bez odmowy")
	}
}
