package core

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"danacoconsole/shared"
)

// Skutek czynności na tabelach: czy po odpowiedzi udanej w dokumencie NAPRAWDĘ
// stoi tabela, jaką Operator zamówił.
//
// Szkody, które ten plik ma wykluczyć:
//  1. szerokości kolumn zerowe po scaleniu — tabela wygląda na złożoną,
//     a w wydaniu ma kolumny niewidzialne (wymaganie zlecenia wprost);
//  2. sortowanie, które rozrywa wiersze albo przestawia wiersz nagłówkowy;
//  3. zamiana tekstu na tabelę, która zostawia tekst w treści i daje dokument
//     niosący to samo dwa razy;
//  4. usunięcie ostatniego wiersza albo kolumny, po którym zostaje tabela
//     o zerowej siatce.
//
// Miara jest brana OSOBNYM wywołaniem wykazu tabel, a nie z odpowiedzi
// czynności: Operator otworzy dokument ponownie, a nie przeczyta odpowiedź.

// tabelaUprzazSprawdzianu składa adapter modułu nad bazą sprawdzianu.
//
// Adapter woła się wprost, bo wpięcie tych komend do rejestru robi prowadzący
// odcinek — sprawdzian nie ma prawa czekać na cudzą robotę, żeby zmierzyć swoją.
func tabelaUprzazSprawdzianu(t *testing.T) (*adapterStudia, context.Context, string) {
	t.Helper()
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)
	dokument := dokumentZTrescia(t, zmontowany, zycie, "okno-tabele",
		strings.Join([]string{"Akapit pierwszy.", "Akapit drugi."}, "\n"))
	return nowyAdapterStudia(zmontowany.dane.Studio), zycie, dokument.Id
}

// tabelaWykazSprawdzianu odczytuje tabelę osobnym wywołaniem.
func tabelaWykazSprawdzianu(t *testing.T, adapter *adapterStudia, zycie context.Context,
	dokument, tabela string) shared.StudioDocumentTable {
	t.Helper()

	wykaz, err := adapter.WykazTabel(zycie, shared.StudioTableListRequest{
		DocumentId: dokument, TableId: &tabela,
	})
	if err != nil {
		t.Fatalf("wykaz tabel odmówił: %v", err)
	}
	if len(wykaz.Tables) != 1 {
		t.Fatalf("wykaz niesie %d tabel, oczekiwano jednej", len(wykaz.Tables))
	}
	return wykaz.Tables[0]
}

// tabelaKomorkaSprawdzianu odczytuje treść komórki wykazu.
func tabelaKomorkaSprawdzianu(t *testing.T, tabela shared.StudioDocumentTable,
	wiersz, kolumna int) string {
	t.Helper()

	for _, komorka := range tabela.Cells {
		if komorka.Row == wiersz && komorka.Column == kolumna {
			if komorka.Text == nil {
				return ""
			}
			return *komorka.Text
		}
	}
	t.Fatalf("tabela nie ma komórki (%d, %d)", wiersz, kolumna)
	return ""
}

// tabelaSprawdzSzerokosci pilnuje wymagania zlecenia: szerokości policzone, nie
// zerowe, i zgodne z liczbą kolumn.
func tabelaSprawdzSzerokosci(t *testing.T, tabela shared.StudioDocumentTable, gdzie string) {
	t.Helper()

	if len(tabela.ColumnWidthsMm) != tabela.Columns {
		t.Fatalf("%s: tabela ma %d kolumn, a szerokości %d", gdzie, tabela.Columns,
			len(tabela.ColumnWidthsMm))
	}
	suma := 0.0
	for numer, szerokosc := range tabela.ColumnWidthsMm {
		if szerokosc <= 0 {
			t.Errorf("%s: kolumna %d ma szerokość %v — szerokość zerowa znaczy kolumnę "+
				"niewidzialną", gdzie, numer, szerokosc)
		}
		suma += szerokosc
	}
	if tabela.WidthMm == nil {
		t.Fatalf("%s: tabela bez szerokości całości", gdzie)
	}
	if tabelaRoznicaMiar(suma, *tabela.WidthMm) > 0.5 {
		t.Errorf("%s: suma szerokości kolumn %v nie zgadza się ze szerokością tabeli %v",
			gdzie, suma, *tabela.WidthMm)
	}
}

func tabelaRoznicaMiar(pierwsza, druga float64) float64 {
	if pierwsza > druga {
		return pierwsza - druga
	}
	return druga - pierwsza
}

// TestTabelaWstawionaMaSiatkeISzerokosci mierzy założenie tabeli: pełna siatka,
// policzone szerokości i powtarzanie wiersza nagłówkowego.
func TestTabelaWstawionaMaSiatkeISzerokosci(t *testing.T) {
	adapter, zycie, dokument := tabelaUprzazSprawdzianu(t)

	wstawiona, err := adapter.WstawTabele(zycie, shared.StudioTableInsertRequest{
		DocumentId: dokument, Offset: 0, Rows: 3, Columns: 4,
	})
	if err != nil {
		t.Fatalf("wstawienie tabeli odmówiło: %v", err)
	}
	tabela := tabelaWykazSprawdzianu(t, adapter, zycie, dokument, wstawiona.Table.Id)

	if len(tabela.Cells) != 12 {
		t.Errorf("tabela 3×4 ma %d komórek, oczekiwano 12 — siatka ma być pełna",
			len(tabela.Cells))
	}
	tabelaSprawdzSzerokosci(t, tabela, "tabela świeżo założona")
	if tabela.RepeatHeader == nil || !*tabela.RepeatHeader {
		t.Error("tabela założona nie powtarza wiersza nagłówkowego na kolejnych stronach")
	}
	if tabela.HeaderRows == nil || *tabela.HeaderRows != 1 {
		t.Error("tabela założona nie ma wskazanego wiersza nagłówkowego")
	}
}

// TestTabelaPoScaleniuMaPoliczoneSzerokosci jest sednem wymagania zlecenia:
// scalenie komórek nie ma prawa zostawić szerokości zerowych, a treść komórek
// wchłoniętych nie ma prawa przepaść.
func TestTabelaPoScaleniuMaPoliczoneSzerokosci(t *testing.T) {
	adapter, zycie, dokument := tabelaUprzazSprawdzianu(t)

	komorki, err := json.Marshal([]shared.StudioTableCell{
		{Row: 1, Column: 0, Text: wskaznik("lewa")},
		{Row: 1, Column: 1, Text: wskaznik("prawa")},
	})
	if err != nil {
		t.Fatalf("nie można złożyć komórek żądania: %v", err)
	}
	wstawiona, err := adapter.WstawTabele(zycie, shared.StudioTableInsertRequest{
		DocumentId: dokument, Offset: 0, Rows: 3, Columns: 3, Cells: komorki,
	})
	if err != nil {
		t.Fatalf("wstawienie tabeli odmówiło: %v", err)
	}

	if _, err := adapter.ZmienBudoweTabeli(zycie, shared.StudioTableStructureEditRequest{
		DocumentId: dokument, TableId: wstawiona.Table.Id,
		Operation: shared.StudioTableStructureOpMergeCells,
		Row:       wskaznik(1), Column: wskaznik(0),
		RowSpan: wskaznik(1), ColumnSpan: wskaznik(2),
	}); err != nil {
		t.Fatalf("scalenie komórek odmówiło: %v", err)
	}

	tabela := tabelaWykazSprawdzianu(t, adapter, zycie, dokument, wstawiona.Table.Id)
	tabelaSprawdzSzerokosci(t, tabela, "tabela po scaleniu")

	wiodaca := tabelaKomorkaSprawdzianu(t, tabela, 1, 0)
	if !strings.Contains(wiodaca, "lewa") || !strings.Contains(wiodaca, "prawa") {
		t.Errorf("scalenie zgubiło treść komórki wchłoniętej: komórka wiodąca niesie %q",
			wiodaca)
	}
	for _, komorka := range tabela.Cells {
		if komorka.Row != 1 || komorka.Column != 1 {
			continue
		}
		if komorka.Merged == nil || !*komorka.Merged {
			t.Error("komórka wchłonięta scaleniem nie jest oznaczona jako wchłonięta")
		}
	}
}

// TestTabelaSortowanieOmijaNaglowek mierzy sortowanie: wiersz nagłówkowy stoi na
// miejscu, a wiersze zawartości idą w porządku liczbowym CAŁE.
func TestTabelaSortowanieOmijaNaglowek(t *testing.T) {
	adapter, zycie, dokument := tabelaUprzazSprawdzianu(t)

	komorki, err := json.Marshal([]shared.StudioTableCell{
		{Row: 0, Column: 0, Text: wskaznik("Pozycja")},
		{Row: 0, Column: 1, Text: wskaznik("Kwota")},
		{Row: 1, Column: 0, Text: wskaznik("druga")},
		{Row: 1, Column: 1, Text: wskaznik("120,50")},
		{Row: 2, Column: 0, Text: wskaznik("pierwsza")},
		{Row: 2, Column: 1, Text: wskaznik("9,90")},
	})
	if err != nil {
		t.Fatalf("nie można złożyć komórek żądania: %v", err)
	}
	wstawiona, err := adapter.WstawTabele(zycie, shared.StudioTableInsertRequest{
		DocumentId: dokument, Offset: 0, Rows: 3, Columns: 2, Cells: komorki,
		HeaderRows: wskaznik(1),
	})
	if err != nil {
		t.Fatalf("wstawienie tabeli odmówiło: %v", err)
	}

	if _, err := adapter.SortujTabele(zycie, shared.StudioTableSortRequest{
		DocumentId: dokument, TableId: wstawiona.Table.Id, Column: 1,
		Numeric: wskaznik(true),
	}); err != nil {
		t.Fatalf("sortowanie tabeli odmówiło: %v", err)
	}

	tabela := tabelaWykazSprawdzianu(t, adapter, zycie, dokument, wstawiona.Table.Id)
	if naglowek := tabelaKomorkaSprawdzianu(t, tabela, 0, 0); naglowek != "Pozycja" {
		t.Errorf("sortowanie przestawiło wiersz nagłówkowy: w nagłówku stoi %q", naglowek)
	}
	if pierwsza := tabelaKomorkaSprawdzianu(t, tabela, 1, 1); pierwsza != "9,90" {
		t.Errorf("po sortowaniu rosnącym w pierwszym wierszu zawartości stoi %q, "+
			"oczekiwano 9,90", pierwsza)
	}
	// Wiersz ma iść CAŁY: kwota bez swojej pozycji znaczyłaby tabelę rozerwaną.
	if opis := tabelaKomorkaSprawdzianu(t, tabela, 1, 0); opis != "pierwsza" {
		t.Errorf("sortowanie rozerwało wiersz: przy kwocie 9,90 stoi pozycja %q", opis)
	}
}

// TestTabelaZamianaTekstuNaTabeleZdejmujeTekst mierzy, że zamiana nie zostawia
// dokumentu niosącego to samo dwa razy.
func TestTabelaZamianaTekstuNaTabeleZdejmujeTekst(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)
	dokument := dokumentZTrescia(t, zmontowany, zycie, "okno-tabele",
		"Wstęp.\nJan\t42\nEwa\t7\nZakończenie.")
	adapter := nowyAdapterStudia(zmontowany.dane.Studio)

	tresc := "Wstęp.\nJan\t42\nEwa\t7\nZakończenie."
	poczatek := strings.Index(tresc, "Jan")
	koniec := strings.Index(tresc, "\nZakończenie.")

	zamieniona, err := adapter.ZamienTabeleITekst(zycie, shared.StudioTableConvertRequest{
		DocumentId: dokument.Id, Direction: shared.StudioTableConvertTextToTable,
		RangeStart: wskaznik(len([]rune(tresc[:poczatek]))),
		RangeEnd:   wskaznik(len([]rune(tresc[:koniec]))),
	})
	if err != nil {
		t.Fatalf("zamiana tekstu na tabelę odmówiła: %v", err)
	}
	if zamieniona.Table == nil {
		t.Fatal("zamiana tekstu na tabelę nie oddała tabeli")
	}
	if zamieniona.Table.Rows != 2 || zamieniona.Table.Columns != 2 {
		t.Errorf("tabela z tekstu ma %d×%d, oczekiwano 2×2",
			zamieniona.Table.Rows, zamieniona.Table.Columns)
	}
	tabela := tabelaWykazSprawdzianu(t, adapter, zycie, dokument.Id, zamieniona.Table.Id)
	if wartosc := tabelaKomorkaSprawdzianu(t, tabela, 0, 1); wartosc != "42" {
		t.Errorf("komórka tabeli niesie %q, oczekiwano 42", wartosc)
	}
	tabelaSprawdzSzerokosci(t, tabela, "tabela z tekstu")

	// Treść dokumentu nie ma już niosić zamienionych wierszy.
	pozostala := trescDokumentu(t, zmontowany, zycie, "okno-tabele", dokument.Id)
	if strings.Contains(pozostala, "Jan\t42") {
		t.Errorf("po zamianie tekst został w treści: %q", pozostala)
	}
	if !strings.Contains(pozostala, "Zakończenie.") {
		t.Errorf("zamiana zjadła treść poza zakresem: %q", pozostala)
	}
}

// TestTabelaZamianaNaTekstOddajeBilansPominiec pilnuje uczciwości: scalenie
// i postać komórek nie przechodzą do tekstu i nie wolno tego przemilczeć.
func TestTabelaZamianaNaTekstOddajeBilansPominiec(t *testing.T) {
	adapter, zycie, dokument := tabelaUprzazSprawdzianu(t)

	wstawiona, err := adapter.WstawTabele(zycie, shared.StudioTableInsertRequest{
		DocumentId: dokument, Offset: 0, Rows: 2, Columns: 2,
	})
	if err != nil {
		t.Fatalf("wstawienie tabeli odmówiło: %v", err)
	}
	if _, err := adapter.ZmienBudoweTabeli(zycie, shared.StudioTableStructureEditRequest{
		DocumentId: dokument, TableId: wstawiona.Table.Id,
		Operation: shared.StudioTableStructureOpMergeCells,
		Row:       wskaznik(1), Column: wskaznik(0),
		RowSpan: wskaznik(1), ColumnSpan: wskaznik(2),
	}); err != nil {
		t.Fatalf("scalenie komórek odmówiło: %v", err)
	}

	zamieniona, err := adapter.ZamienTabeleITekst(zycie, shared.StudioTableConvertRequest{
		DocumentId: dokument, Direction: shared.StudioTableConvertTableToText,
		TableId: &wstawiona.Table.Id,
	})
	if err != nil {
		t.Fatalf("zamiana tabeli na tekst odmówiła: %v", err)
	}
	if zamieniona.Balance.SkippedCount == 0 {
		t.Error("zamiana tabeli ze scaleniem na tekst nie oddała ani jednego pominięcia — " +
			"przemilczenie straty jest zakazane")
	}
	wykaz, err := adapter.WykazTabel(zycie, shared.StudioTableListRequest{DocumentId: dokument})
	if err != nil {
		t.Fatalf("wykaz tabel odmówił: %v", err)
	}
	if len(wykaz.Tables) != 0 {
		t.Errorf("po zamianie na tekst dokument niesie %d tabel, oczekiwano zera",
			len(wykaz.Tables))
	}
}

// TestTabelaUsuniecieCalejSiatkiOdmawia pilnuje granicy: tabela bez ani jednego
// wiersza nie jest tabelą i rdzeń mówi, czym się ją usuwa.
func TestTabelaUsuniecieCalejSiatkiOdmawia(t *testing.T) {
	adapter, zycie, dokument := tabelaUprzazSprawdzianu(t)

	wstawiona, err := adapter.WstawTabele(zycie, shared.StudioTableInsertRequest{
		DocumentId: dokument, Offset: 0, Rows: 2, Columns: 2,
	})
	if err != nil {
		t.Fatalf("wstawienie tabeli odmówiło: %v", err)
	}
	_, err = adapter.ZmienBudoweTabeli(zycie, shared.StudioTableStructureEditRequest{
		DocumentId: dokument, TableId: wstawiona.Table.Id,
		Operation: shared.StudioTableStructureOpDeleteRow,
		Row:       wskaznik(0), Count: wskaznik(2),
	})
	if err == nil {
		t.Fatal("usunięcie wszystkich wierszy tabeli wróciło bez odmowy")
	}
	if !strings.Contains(err.Error(), "convert") {
		t.Errorf("odmowa nie mówi, czym usuwa się tabelę w całości: %v", err)
	}
}

// TestTabelaWstawienieKolumnyZachowujeSzerokosci mierzy, że kolumna wstawiona
// nie jest kolumną zerowej szerokości.
func TestTabelaWstawienieKolumnyZachowujeSzerokosci(t *testing.T) {
	adapter, zycie, dokument := tabelaUprzazSprawdzianu(t)

	wstawiona, err := adapter.WstawTabele(zycie, shared.StudioTableInsertRequest{
		DocumentId: dokument, Offset: 0, Rows: 2, Columns: 2,
	})
	if err != nil {
		t.Fatalf("wstawienie tabeli odmówiło: %v", err)
	}
	if _, err := adapter.ZmienBudoweTabeli(zycie, shared.StudioTableStructureEditRequest{
		DocumentId: dokument, TableId: wstawiona.Table.Id,
		Operation: shared.StudioTableStructureOpInsertColumn,
		Column:    wskaznik(1), Before: wskaznik(false),
	}); err != nil {
		t.Fatalf("wstawienie kolumny odmówiło: %v", err)
	}
	tabela := tabelaWykazSprawdzianu(t, adapter, zycie, dokument, wstawiona.Table.Id)
	if tabela.Columns != 3 {
		t.Fatalf("tabela ma %d kolumn, oczekiwano trzech", tabela.Columns)
	}
	if len(tabela.Cells) != 6 {
		t.Errorf("tabela 2×3 ma %d komórek, oczekiwano sześciu", len(tabela.Cells))
	}
	tabelaSprawdzSzerokosci(t, tabela, "tabela po wstawieniu kolumny")
}

// TestTabelaPostacBezCechyOdmawia pilnuje zakazu odpowiedzi „ok" bez skutku.
func TestTabelaPostacBezCechyOdmawia(t *testing.T) {
	adapter, zycie, dokument := tabelaUprzazSprawdzianu(t)

	wstawiona, err := adapter.WstawTabele(zycie, shared.StudioTableInsertRequest{
		DocumentId: dokument, Offset: 0, Rows: 2, Columns: 2,
	})
	if err != nil {
		t.Fatalf("wstawienie tabeli odmówiło: %v", err)
	}
	if _, err := adapter.UstawPostacTabeli(zycie, shared.StudioTableFormatSetRequest{
		DocumentId: dokument, TableId: wstawiona.Table.Id,
	}); err == nil {
		t.Error("ustawienie postaci tabeli bez ani jednej cechy wróciło bez odmowy")
	}
}
