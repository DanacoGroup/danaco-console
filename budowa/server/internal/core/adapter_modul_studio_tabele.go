// Odpowiedzialność pliku: sześć czynności na tabelach dokumentu — założenie
// tabeli o wskazanym rozmiarze, wykaz tabel, zmiana budowy (wiersze, kolumny,
// scalenie, podział), postać tabeli i komórek wraz z powtarzaniem wiersza
// nagłówkowego, sortowanie zawartości oraz zamiana tekstu na tabelę i tabeli
// na tekst.
//
// Rachunek siatki stoi w `_tabele_pomocniki.go`, a droga zapisu, zmiany
// śledzonej i dziennika — w `_postac.go` (`postacWczytaj`, `postacZakoncz`).
// Ten plik nie liczy ani jednej z tych rzeczy drugi raz.
//
// ── Dlaczego tabela nie zajmuje ani jednego znaku treści ────────────────────
// Tabela wchodzi do drzewa blokiem nietekstowym, więc jej wstawienie nie
// przesuwa ani jednego zakresu zaznaczenia, przypisu ani blokady. Gdyby tabela
// zajmowała znaki, każde jej wstawienie rozjeżdżałoby wszystko, co wisi na
// miejscu w treści — a to jest właśnie ta cicha szkoda, której zlecenie zakazuje.
package core

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"danacoconsole/shared"
)

// tabelaRozdzielnikDomyslny — brak wskazania znaczy tabulator, tak jak
// w pakiecie biurowym.
const tabelaRozdzielnikDomyslny = "\t"

// WstawTabele zakłada tabelę o wskazanym rozmiarze (`studio.table.insert`).
func (a *adapterStudia) WstawTabele(ctx context.Context,
	z shared.StudioTableInsertRequest) (shared.StudioTableInsertResponse, error) {

	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioTableInsertResponse{}, err
	}
	autor := postacAutor(z.Author)

	if z.Rows <= 0 || z.Columns <= 0 {
		return shared.StudioTableInsertResponse{}, tabelaBladWskazania(
			"założenie tabeli o " + strconv.Itoa(z.Rows) + " wierszach i " +
				strconv.Itoa(z.Columns) + " kolumnach — tabela wymaga co najmniej " +
				"jednego wiersza i jednej kolumny")
	}
	if z.Rows > tabelaGranicaWierszy || z.Columns > tabelaGranicaKolumn {
		return shared.StudioTableInsertResponse{}, tabelaBladWskazania(
			"założenie tabeli o " + strconv.Itoa(z.Rows) + " wierszach i " +
				strconv.Itoa(z.Columns) + " kolumnach przekracza granicę rdzenia: " +
				strconv.Itoa(tabelaGranicaWierszy) + " wierszy, " +
				strconv.Itoa(tabelaGranicaKolumn) + " kolumn")
	}

	dlugosc := postacDlugosc(&stan.forma)
	miejsce := z.Offset
	if miejsce < 0 {
		miejsce = 0
	}
	if miejsce > dlugosc {
		miejsce = dlugosc
	}
	// Blokada obowiązuje w rdzeniu, przed dotknięciem dokumentu: tabela
	// wstawiona w zablokowany fragment omijałaby blokadę bez wysiłku.
	if odcinki, pominiete := postacOdcinkiDozwolone(&stan.forma, miejsce, miejsce, autor); len(odcinki) == 0 {
		return shared.StudioTableInsertResponse{}, tabelaBladWskazania(
			"wstawienie tabeli " + postacZapisZakresu(miejsce, miejsce) +
				" zatrzymane przez blokadę fragmentu: " + postacNazwaBlokad(pominiete))
	}

	tabela := shared.StudioDocumentTable{
		Id: nowyIdentyfikator(przedrostekTabeliPostaci), Rows: z.Rows, Columns: z.Columns,
		StyleName: z.StyleName, HeaderRows: z.HeaderRows, RepeatHeader: z.RepeatHeader,
		WidthMm: z.WidthMm,
	}
	if tabela.HeaderRows == nil {
		tabela.HeaderRows = postacWskaznikLiczby(1)
	}
	if *tabela.HeaderRows > tabela.Rows {
		tabela.HeaderRows = postacWskaznikLiczby(tabela.Rows)
	}
	if tabela.RepeatHeader == nil {
		// Powtarzanie wiersza nagłówkowego na kolejnych stronach jest wymienione
		// w zleceniu wprost i jest nastawą, której Operator oczekuje domyślnie —
		// tabela wielostronicowa bez nagłówka na drugiej stronie jest nieczytelna.
		tabela.RepeatHeader = postacWskaznikPrawdy(*tabela.HeaderRows > 0)
	}
	tabelaSiatkaPelna(&tabela)

	bilans := shared.StudioActionBalance{Skipped: []shared.StudioSkippedItem{}}
	if len(z.Cells) > 0 {
		var przyslane []shared.StudioTableCell
		if err := json.Unmarshal(z.Cells, &przyslane); err != nil {
			return shared.StudioTableInsertResponse{}, tabelaBladWskazania(
				"treść komórek w żądaniu jest nieczytelna: " + err.Error())
		}
		for _, przyslana := range przyslane {
			komorka := tabelaKomorka(&tabela, przyslana.Row, przyslana.Column)
			if komorka == nil {
				bilans.Skipped = append(bilans.Skipped, shared.StudioSkippedItem{
					Reason: "komórka poza tabelą",
					Detail: postacWskaznikTekstu("komórka (" + strconv.Itoa(przyslana.Row) +
						", " + strconv.Itoa(przyslana.Column) + ") leży poza tabelą o " +
						strconv.Itoa(tabela.Rows) + " wierszach i " +
						strconv.Itoa(tabela.Columns) + " kolumnach"),
				})
				continue
			}
			przyslana.Row, przyslana.Column = komorka.Row, komorka.Column
			*komorka = przyslana
		}
	}
	tabelaPrzeliczSzerokosci(&tabela, tabelaSzerokoscTekstu(&stan.forma))

	stan.forma.Tables = append(stan.forma.Tables, tabela)
	tabelaWstawBlokWMiejscu(&stan.forma, shared.StudioDocumentBlock{
		Id: nowyIdentyfikator(przedrostekBlokuPostaci), Kind: blokPostaciTabela,
		TableId: postacWskaznikTekstu(tabela.Id),
	}, miejsce)

	bilans.Applied = 1
	bilans.Note = postacWskaznikTekstu("tabela o " + strconv.Itoa(tabela.Rows) +
		" wierszach i " + strconv.Itoa(tabela.Columns) + " kolumnach założona " +
		postacZapisZakresu(miejsce, miejsce))
	stan.opisCzynnosci = "wstawienie tabeli " + postacZapisZakresu(miejsce, miejsce)

	forma, bilansGotowy, zmiana, err := a.postacZakoncz(ctx, stan, autor,
		shared.StudioChangeKindWstawienie, shared.StudioActionKindTableChange,
		miejsce, miejsce, bilans)
	if err != nil {
		return shared.StudioTableInsertResponse{}, err
	}
	zalozona, err := tabelaZnajdz(&forma, tabela.Id)
	if err != nil {
		return shared.StudioTableInsertResponse{}, err
	}
	return shared.StudioTableInsertResponse{
		Form: forma, Balance: bilansGotowy, Change: zmiana,
		ActionId: stan.czynnosc, Table: *zalozona,
	}, nil
}

// WykazTabel oddaje tabele dokumentu (`studio.table.list`).
func (a *adapterStudia) WykazTabel(ctx context.Context,
	z shared.StudioTableListRequest) (shared.StudioTableListResponse, error) {

	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioTableListResponse{}, err
	}
	if z.TableId != nil && strings.TrimSpace(*z.TableId) != "" {
		tabela, err := tabelaZnajdz(&stan.forma, *z.TableId)
		if err != nil {
			return shared.StudioTableListResponse{}, err
		}
		kopia := *tabela
		tabelaSiatkaPelna(&kopia)
		tabelaPrzeliczSzerokosci(&kopia, tabelaSzerokoscTekstu(&stan.forma))
		return shared.StudioTableListResponse{Tables: []shared.StudioDocumentTable{kopia}}, nil
	}
	tabele := make([]shared.StudioDocumentTable, 0, len(stan.forma.Tables))
	for _, tabela := range stan.forma.Tables {
		kopia := tabela
		tabelaSiatkaPelna(&kopia)
		tabelaPrzeliczSzerokosci(&kopia, tabelaSzerokoscTekstu(&stan.forma))
		tabele = append(tabele, kopia)
	}
	return shared.StudioTableListResponse{Tables: tabele}, nil
}

// ZmienBudoweTabeli wstawia i usuwa wiersze oraz kolumny, scala i dzieli
// komórki (`studio.table.structure.edit`).
func (a *adapterStudia) ZmienBudoweTabeli(ctx context.Context,
	z shared.StudioTableStructureEditRequest) (shared.StudioTableStructureEditResponse, error) {

	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioTableStructureEditResponse{}, err
	}
	autor := postacAutor(z.Author)
	tabela, err := tabelaZnajdz(&stan.forma, z.TableId)
	if err != nil {
		return shared.StudioTableStructureEditResponse{}, err
	}
	tabelaSiatkaPelna(tabela)

	miejsce := tabelaMiejsceTabeli(&stan.forma, tabela.Id)
	if odcinki, pominiete := postacOdcinkiDozwolone(&stan.forma, miejsce, miejsce, autor); len(odcinki) == 0 {
		return shared.StudioTableStructureEditResponse{}, tabelaBladWskazania(
			"zmiana budowy tabeli zatrzymana przez blokadę fragmentu: " +
				postacNazwaBlokad(pominiete))
	}

	wiersz, kolumna := 0, 0
	if z.Row != nil {
		wiersz = *z.Row
	}
	if z.Column != nil {
		kolumna = *z.Column
	}
	ile := 1
	if z.Count != nil && *z.Count > 0 {
		ile = *z.Count
	}
	przed := true
	if z.Before != nil {
		przed = *z.Before
	}
	objeteWiersze, objeteKolumny := 1, 1
	if z.RowSpan != nil {
		objeteWiersze = *z.RowSpan
	}
	if z.ColumnSpan != nil {
		objeteKolumny = *z.ColumnSpan
	}

	var opis string
	switch z.Operation {
	case shared.StudioTableStructureOpInsertRow:
		if err := tabelaWstawWiersze(tabela, wiersz, przed, ile); err != nil {
			return shared.StudioTableStructureEditResponse{}, err
		}
		opis = "wstawienie " + strconv.Itoa(ile) + " wierszy tabeli"
	case shared.StudioTableStructureOpDeleteRow:
		if err := tabelaUsunWiersze(tabela, wiersz, ile); err != nil {
			return shared.StudioTableStructureEditResponse{}, err
		}
		opis = "usunięcie " + strconv.Itoa(ile) + " wierszy tabeli"
	case shared.StudioTableStructureOpInsertColumn:
		if err := tabelaWstawKolumny(tabela, kolumna, przed, ile); err != nil {
			return shared.StudioTableStructureEditResponse{}, err
		}
		opis = "wstawienie " + strconv.Itoa(ile) + " kolumn tabeli"
	case shared.StudioTableStructureOpDeleteColumn:
		if err := tabelaUsunKolumny(tabela, kolumna, ile); err != nil {
			return shared.StudioTableStructureEditResponse{}, err
		}
		opis = "usunięcie " + strconv.Itoa(ile) + " kolumn tabeli"
	case shared.StudioTableStructureOpMergeCells:
		if err := tabelaScalKomorki(tabela, wiersz, kolumna, objeteWiersze, objeteKolumny); err != nil {
			return shared.StudioTableStructureEditResponse{}, err
		}
		opis = "scalenie komórek tabeli"
	case shared.StudioTableStructureOpSplitCell:
		if err := tabelaPodzielKomorke(tabela, wiersz, kolumna, objeteWiersze, objeteKolumny); err != nil {
			return shared.StudioTableStructureEditResponse{}, err
		}
		opis = "podział komórki tabeli"
	default:
		return shared.StudioTableStructureEditResponse{}, tabelaBladWskazania(
			"czynność na budowie tabeli o nazwie „" + string(z.Operation) +
				"”, której rdzeń nie zna; wykaz: insertRow, deleteRow, insertColumn, " +
				"deleteColumn, mergeCells, splitCell")
	}
	// Szerokości liczą się po KAŻDEJ zmianie budowy — to jest wymaganie
	// zlecenia: tabela po scaleniu komórek ma szerokości policzone, nie zerowe.
	tabelaPrzeliczSzerokosci(tabela, tabelaSzerokoscTekstu(&stan.forma))

	bilans := shared.StudioActionBalance{
		Applied: 1, Skipped: []shared.StudioSkippedItem{},
		Note: postacWskaznikTekstu(opis + "; tabela ma teraz " + strconv.Itoa(tabela.Rows) +
			" wierszy i " + strconv.Itoa(tabela.Columns) + " kolumn"),
	}
	stan.opisCzynnosci = opis

	forma, bilansGotowy, zmiana, err := a.postacZakoncz(ctx, stan, autor,
		shared.StudioChangeKindFormatowanie, shared.StudioActionKindTableChange,
		miejsce, miejsce, bilans)
	if err != nil {
		return shared.StudioTableStructureEditResponse{}, err
	}
	zmieniona, err := tabelaZnajdz(&forma, z.TableId)
	if err != nil {
		return shared.StudioTableStructureEditResponse{}, err
	}
	return shared.StudioTableStructureEditResponse{
		Form: forma, Balance: bilansGotowy, Change: zmiana,
		ActionId: stan.czynnosc, Table: *zmieniona,
	}, nil
}

// UstawPostacTabeli ustawia postać tabeli i komórek (`studio.table.format.set`).
//
// Brak wskazania wiersza i kolumny znaczy CAŁĄ tabelę, wskazanie samego wiersza
// — cały wiersz, samej kolumny — całą kolumnę. Tak samo jak w pakiecie biurowym.
func (a *adapterStudia) UstawPostacTabeli(ctx context.Context,
	z shared.StudioTableFormatSetRequest) (shared.StudioTableFormatSetResponse, error) {

	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioTableFormatSetResponse{}, err
	}
	autor := postacAutor(z.Author)
	tabela, err := tabelaZnajdz(&stan.forma, z.TableId)
	if err != nil {
		return shared.StudioTableFormatSetResponse{}, err
	}
	tabelaSiatkaPelna(tabela)

	var obramowanie *shared.StudioBorder
	if len(z.Border) > 0 {
		var wczytane shared.StudioBorder
		if err := json.Unmarshal(z.Border, &wczytane); err != nil {
			return shared.StudioTableFormatSetResponse{}, tabelaBladWskazania(
				"obramowanie w żądaniu jest nieczytelne: " + err.Error())
		}
		obramowanie = &wczytane
	}
	cechaTabeli := len(z.ColumnWidthsMm) > 0 || z.WidthMm != nil || z.StyleName != nil ||
		z.HeaderRows != nil || z.RepeatHeader != nil || z.Caption != nil
	cechaKomorki := obramowanie != nil || z.ShadingColor != nil || z.VerticalAlign != nil ||
		z.Align != nil
	if !cechaTabeli && !cechaKomorki {
		return shared.StudioTableFormatSetResponse{}, tabelaBladWskazania(
			"ustawienie postaci tabeli bez ani jednej cechy do ustawienia")
	}

	miejsce := tabelaMiejsceTabeli(&stan.forma, tabela.Id)
	if odcinki, pominiete := postacOdcinkiDozwolone(&stan.forma, miejsce, miejsce, autor); len(odcinki) == 0 {
		return shared.StudioTableFormatSetResponse{}, tabelaBladWskazania(
			"zmiana postaci tabeli zatrzymana przez blokadę fragmentu: " +
				postacNazwaBlokad(pominiete))
	}

	bilans := shared.StudioActionBalance{Skipped: []shared.StudioSkippedItem{}}

	if len(z.ColumnWidthsMm) > 0 {
		if len(z.ColumnWidthsMm) != tabela.Columns {
			bilans.Skipped = append(bilans.Skipped, shared.StudioSkippedItem{
				Reason: "szerokości kolumn niezgodne z tabelą",
				Detail: postacWskaznikTekstu("podano " + strconv.Itoa(len(z.ColumnWidthsMm)) +
					" szerokości, a tabela ma kolumn " + strconv.Itoa(tabela.Columns) +
					"; szerokości podane zostały użyte od pierwszej kolumny, reszta " +
					"policzona przez rdzeń"),
			})
		}
		szerokosci := make([]float64, tabela.Columns)
		for i := range szerokosci {
			if i < len(z.ColumnWidthsMm) && z.ColumnWidthsMm[i] > 0 {
				szerokosci[i] = z.ColumnWidthsMm[i]
			}
		}
		tabela.ColumnWidthsMm = szerokosci
		tabela.WidthMm = nil
		bilans.Applied++
	}
	if z.WidthMm != nil {
		tabela.WidthMm = z.WidthMm
		if len(z.ColumnWidthsMm) == 0 {
			// Szerokość całej tabeli rozdziela się na kolumny od nowa: kolumny
			// zostawione bez zmiany nie sumowałyby się do nowej szerokości.
			tabela.ColumnWidthsMm = nil
		}
		bilans.Applied++
	}
	if z.StyleName != nil {
		tabela.StyleName = z.StyleName
		bilans.Applied++
	}
	if z.HeaderRows != nil {
		ile := *z.HeaderRows
		if ile < 0 {
			ile = 0
		}
		if ile > tabela.Rows {
			bilans.Skipped = append(bilans.Skipped, shared.StudioSkippedItem{
				Reason: "wiersz nagłówkowy poza tabelą",
				Detail: postacWskaznikTekstu("wskazano " + strconv.Itoa(ile) +
					" wierszy nagłówkowych, a tabela ma wierszy " +
					strconv.Itoa(tabela.Rows)),
			})
			ile = tabela.Rows
		}
		tabela.HeaderRows = postacWskaznikLiczby(ile)
		bilans.Applied++
	}
	if z.RepeatHeader != nil {
		if *z.RepeatHeader && (tabela.HeaderRows == nil || *tabela.HeaderRows == 0) {
			// Powtarzanie nagłówka bez wskazania, który wiersz jest nagłówkiem,
			// nie miałoby czego powtarzać — rdzeń bierze wiersz pierwszy i mówi
			// o tym w bilansie, zamiast milczeć.
			tabela.HeaderRows = postacWskaznikLiczby(1)
			bilans.Note = postacWskaznikTekstu("powtarzanie wiersza nagłówkowego włączone; " +
				"za nagłówek przyjęty wiersz pierwszy, bo tabela nie miała go wskazanego")
		}
		tabela.RepeatHeader = z.RepeatHeader
		bilans.Applied++
	}
	if z.Caption != nil {
		tabela.Caption = z.Caption
		// Podpis tabeli zmieniony znaczy spis tabel nieświeży — inaczej spis
		// mówiłby co innego niż tabela pod nim.
		if err := a.aparatZnaczNieswiezoscRodzaju(ctx, stan,
			shared.StudioApparatusKindTableIndex); err != nil {

			return shared.StudioTableFormatSetResponse{}, err
		}
		bilans.Applied++
	}

	if cechaKomorki {
		zmienione := 0
		for wiersz := 0; wiersz < tabela.Rows; wiersz++ {
			if z.Row != nil && *z.Row != wiersz {
				continue
			}
			for kolumna := 0; kolumna < tabela.Columns; kolumna++ {
				if z.Column != nil && *z.Column != kolumna {
					continue
				}
				komorka := tabelaKomorka(tabela, wiersz, kolumna)
				if komorka == nil {
					continue
				}
				if obramowanie != nil {
					kopia := *obramowanie
					komorka.Border = &kopia
				}
				if z.ShadingColor != nil {
					if *z.ShadingColor == "" {
						komorka.ShadingColor = nil
					} else {
						komorka.ShadingColor = z.ShadingColor
					}
				}
				if z.VerticalAlign != nil {
					komorka.VerticalAlign = z.VerticalAlign
				}
				if z.Align != nil {
					// Wyrównanie w komórce jest cechą jej akapitu — tak samo jak
					// w treści dokumentu, więc idzie tą samą strukturą.
					komorka.Paragraph = postacScalAkapit(komorka.Paragraph,
						shared.StudioParagraphFormat{Align: z.Align})
				}
				zmienione++
			}
		}
		if zmienione == 0 {
			return shared.StudioTableFormatSetResponse{}, tabelaBladWskazania(
				"postać komórek nie miała na czym stanąć: wskazanie wiersza i kolumny " +
					"nie obejmuje ani jednej komórki tabeli o " + strconv.Itoa(tabela.Rows) +
					" wierszach i " + strconv.Itoa(tabela.Columns) + " kolumnach")
		}
		if z.Row == nil && z.Column == nil && obramowanie != nil {
			tabela.Border = obramowanie
		}
		if z.Row == nil && z.Column == nil && z.Align != nil {
			tabela.Align = z.Align
		}
		bilans.Applied += zmienione
	}
	tabelaPrzeliczSzerokosci(tabela, tabelaSzerokoscTekstu(&stan.forma))

	if bilans.Note == nil {
		bilans.Note = postacWskaznikTekstu("postać tabeli ustawiona; miejsc zmienionych " +
			strconv.Itoa(bilans.Applied))
	}
	stan.opisCzynnosci = "zmiana postaci tabeli"

	forma, bilansGotowy, zmiana, err := a.postacZakoncz(ctx, stan, autor,
		shared.StudioChangeKindFormatowanie, shared.StudioActionKindTableChange,
		miejsce, miejsce, bilans)
	if err != nil {
		return shared.StudioTableFormatSetResponse{}, err
	}
	zmieniona, err := tabelaZnajdz(&forma, z.TableId)
	if err != nil {
		return shared.StudioTableFormatSetResponse{}, err
	}
	return shared.StudioTableFormatSetResponse{
		Form: forma, Balance: bilansGotowy, Change: zmiana,
		ActionId: stan.czynnosc, Table: *zmieniona,
	}, nil
}

// SortujTabele sortuje zawartość tabeli po wskazanej kolumnie
// (`studio.table.sort`).
func (a *adapterStudia) SortujTabele(ctx context.Context,
	z shared.StudioTableSortRequest) (shared.StudioTableSortResponse, error) {

	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioTableSortResponse{}, err
	}
	autor := postacAutor(z.Author)
	tabela, err := tabelaZnajdz(&stan.forma, z.TableId)
	if err != nil {
		return shared.StudioTableSortResponse{}, err
	}
	tabelaSiatkaPelna(tabela)

	miejsce := tabelaMiejsceTabeli(&stan.forma, tabela.Id)
	if odcinki, pominiete := postacOdcinkiDozwolone(&stan.forma, miejsce, miejsce, autor); len(odcinki) == 0 {
		return shared.StudioTableSortResponse{}, tabelaBladWskazania(
			"sortowanie tabeli zatrzymane przez blokadę fragmentu: " +
				postacNazwaBlokad(pominiete))
	}

	malejaco := z.Descending != nil && *z.Descending
	liczbowo := z.Numeric != nil && *z.Numeric
	pomijaj := z.SkipHeader == nil || *z.SkipHeader

	ile, err := tabelaSortujWiersze(tabela, z.Column, malejaco, liczbowo, pomijaj)
	if err != nil {
		return shared.StudioTableSortResponse{}, err
	}
	tabelaPrzeliczSzerokosci(tabela, tabelaSzerokoscTekstu(&stan.forma))

	porzadek := "rosnąco"
	if malejaco {
		porzadek = "malejąco"
	}
	sposob := "jak tekst"
	if liczbowo {
		sposob = "jak liczby"
	}
	bilans := shared.StudioActionBalance{
		Applied: ile, Skipped: []shared.StudioSkippedItem{},
		Note: postacWskaznikTekstu("tabela posortowana po kolumnie " +
			strconv.Itoa(z.Column) + " " + porzadek + ", " + sposob +
			"; wierszy przestawionych " + strconv.Itoa(ile)),
	}
	stan.opisCzynnosci = "sortowanie tabeli po kolumnie " + strconv.Itoa(z.Column)

	forma, bilansGotowy, zmiana, err := a.postacZakoncz(ctx, stan, autor,
		shared.StudioChangeKindFormatowanie, shared.StudioActionKindTableChange,
		miejsce, miejsce, bilans)
	if err != nil {
		return shared.StudioTableSortResponse{}, err
	}
	posortowana, err := tabelaZnajdz(&forma, z.TableId)
	if err != nil {
		return shared.StudioTableSortResponse{}, err
	}
	return shared.StudioTableSortResponse{
		Form: forma, Balance: bilansGotowy, Change: zmiana,
		ActionId: stan.czynnosc, Table: *posortowana,
	}, nil
}

// ZamienTabeleITekst zamienia tekst na tabelę albo tabelę na tekst
// (`studio.table.convert`).
func (a *adapterStudia) ZamienTabeleITekst(ctx context.Context,
	z shared.StudioTableConvertRequest) (shared.StudioTableConvertResponse, error) {

	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioTableConvertResponse{}, err
	}
	autor := postacAutor(z.Author)
	rozdzielnik := tabelaRozdzielnikDomyslny
	if z.Separator != nil && *z.Separator != "" {
		rozdzielnik = *z.Separator
	}

	switch z.Direction {
	case shared.StudioTableConvertTextToTable:
		return a.tabelaZamienTekstNaTabele(ctx, stan, z, autor, rozdzielnik)
	case shared.StudioTableConvertTableToText:
		return a.tabelaZamienTabeleNaTekst(ctx, stan, z, autor, rozdzielnik)
	default:
		return shared.StudioTableConvertResponse{}, tabelaBladWskazania(
			"kierunek zamiany „" + string(z.Direction) + "”, którego rdzeń nie zna; " +
				"wykaz: textToTable, tableToText")
	}
}

// tabelaZamienTekstNaTabele zakłada tabelę z zaznaczonego tekstu i zdejmuje ten
// tekst z treści — inaczej dokument niósłby to samo dwa razy.
func (a *adapterStudia) tabelaZamienTekstNaTabele(ctx context.Context, stan *stanPostaci,
	z shared.StudioTableConvertRequest, autor shared.StudioAuthor,
	rozdzielnik string) (shared.StudioTableConvertResponse, error) {

	dlugosc := postacDlugosc(&stan.forma)
	od, do := postacZakres(z.RangeStart, z.RangeEnd, dlugosc)
	if od == do {
		return shared.StudioTableConvertResponse{}, tabelaBladWskazania(
			"zamiana tekstu na tabelę bez wskazania zakresu tekstu (rangeStart, rangeEnd)")
	}
	odcinki, pominiete := postacOdcinkiDozwolone(&stan.forma, od, do, autor)
	if len(odcinki) == 0 {
		return shared.StudioTableConvertResponse{}, tabelaBladWskazania(
			"zamiana tekstu na tabelę zatrzymana przez blokadę fragmentu: " +
				postacNazwaBlokad(pominiete))
	}
	znaki := []rune(postacTekstFormy(&stan.forma))
	roboczyOd, roboczyDo := odcinki[0][0], odcinki[0][1]
	wybrany := string(znaki[roboczyOd:roboczyDo])

	tabela, err := tabelaZTekstu(strings.Split(wybrany, "\n"), rozdzielnik)
	if err != nil {
		return shared.StudioTableConvertResponse{}, err
	}
	tabelaPrzeliczSzerokosci(&tabela, tabelaSzerokoscTekstu(&stan.forma))
	if tabela.Rows > 1 {
		tabela.HeaderRows = postacWskaznikLiczby(1)
		tabela.RepeatHeader = postacWskaznikPrawdy(true)
	}

	// Tekst zamieniony schodzi z treści, a tabela wchodzi w jego miejsce.
	postacZamienTresc(&stan.forma, roboczyOd, roboczyDo, "", nil, &autor)
	stan.forma.Tables = append(stan.forma.Tables, tabela)
	tabelaWstawBlokWMiejscu(&stan.forma, shared.StudioDocumentBlock{
		Id: nowyIdentyfikator(przedrostekBlokuPostaci), Kind: blokPostaciTabela,
		TableId: postacWskaznikTekstu(tabela.Id),
	}, roboczyOd)

	bilans := shared.StudioActionBalance{Applied: 1, Skipped: pominiete}
	bilans.Note = postacWskaznikTekstu("tekst " + postacZapisZakresu(roboczyOd, roboczyDo) +
		" zamieniony na tabelę o " + strconv.Itoa(tabela.Rows) + " wierszach i " +
		strconv.Itoa(tabela.Columns) + " kolumnach")
	stan.opisCzynnosci = "zamiana tekstu na tabelę"

	forma, bilansGotowy, zmiana, err := a.postacZakoncz(ctx, stan, autor,
		shared.StudioChangeKindWstawienie, shared.StudioActionKindTableChange,
		roboczyOd, roboczyDo, bilans)
	if err != nil {
		return shared.StudioTableConvertResponse{}, err
	}
	zalozona, err := tabelaZnajdz(&forma, tabela.Id)
	if err != nil {
		return shared.StudioTableConvertResponse{}, err
	}
	return shared.StudioTableConvertResponse{
		Form: forma, Balance: bilansGotowy, Change: zmiana,
		ActionId: stan.czynnosc, Table: zalozona,
	}, nil
}

// tabelaZamienTabeleNaTekst rozkłada tabelę na wiersze tekstu i zdejmuje ją
// z dokumentu.
func (a *adapterStudia) tabelaZamienTabeleNaTekst(ctx context.Context, stan *stanPostaci,
	z shared.StudioTableConvertRequest, autor shared.StudioAuthor,
	rozdzielnik string) (shared.StudioTableConvertResponse, error) {

	if z.TableId == nil || strings.TrimSpace(*z.TableId) == "" {
		return shared.StudioTableConvertResponse{}, tabelaBladWskazania(
			"zamiana tabeli na tekst bez wskazania tabeli (tableId)")
	}
	tabela, err := tabelaZnajdz(&stan.forma, *z.TableId)
	if err != nil {
		return shared.StudioTableConvertResponse{}, err
	}
	tabelaSiatkaPelna(tabela)
	kod := tabela.Id
	miejsce := tabelaMiejsceTabeli(&stan.forma, kod)
	if odcinki, pominiete := postacOdcinkiDozwolone(&stan.forma, miejsce, miejsce, autor); len(odcinki) == 0 {
		return shared.StudioTableConvertResponse{}, tabelaBladWskazania(
			"zamiana tabeli na tekst zatrzymana przez blokadę fragmentu: " +
				postacNazwaBlokad(pominiete))
	}
	wiersze := tabelaTekstZTabeli(tabela, rozdzielnik)

	bilans := shared.StudioActionBalance{Applied: 1, Skipped: []shared.StudioSkippedItem{}}
	if tabelaMaScalenia(tabela) {
		bilans.Skipped = append(bilans.Skipped, shared.StudioSkippedItem{
			Reason: "scalenie komórek nie przechodzi do tekstu",
			Detail: postacWskaznikTekstu("tabela miała komórki scalone; w tekście " +
				"komórki wchłonięte wychodzą jako puste miejsca między znakami " +
				"rozdzielającymi, bo tekst nie niesie pojęcia scalenia"),
		})
	}
	if tabelaMaPostacKomorek(tabela) {
		bilans.Skipped = append(bilans.Skipped, shared.StudioSkippedItem{
			Reason: "postać komórek schodzi",
			Detail: postacWskaznikTekstu("obramowanie, cieniowanie i wyrównanie " +
				"w komórkach nie mają odpowiednika w tekście"),
		})
	}

	// Tabela schodzi z drzewa i z wykazu, a jej zawartość wchodzi w treść
	// w miejscu, w którym tabela stała.
	tabelaUsunBlokTabeli(&stan.forma, kod)
	pozostale := make([]shared.StudioDocumentTable, 0, len(stan.forma.Tables))
	for _, zastana := range stan.forma.Tables {
		if zastana.Id == kod {
			continue
		}
		pozostale = append(pozostale, zastana)
	}
	stan.forma.Tables = pozostale
	postacZamienTresc(&stan.forma, miejsce, miejsce, strings.Join(wiersze, "\n"), nil, &autor)
	if err := a.aparatZnaczNieswiezoscRodzaju(ctx, stan,
		shared.StudioApparatusKindTableIndex); err != nil {

		return shared.StudioTableConvertResponse{}, err
	}

	bilans.Note = postacWskaznikTekstu("tabela zamieniona na " + strconv.Itoa(len(wiersze)) +
		" wierszy tekstu " + postacZapisZakresu(miejsce, miejsce))
	stan.opisCzynnosci = "zamiana tabeli na tekst"

	forma, bilansGotowy, zmiana, err := a.postacZakoncz(ctx, stan, autor,
		shared.StudioChangeKindUsuniecie, shared.StudioActionKindTableChange,
		miejsce, miejsce, bilans)
	if err != nil {
		return shared.StudioTableConvertResponse{}, err
	}
	return shared.StudioTableConvertResponse{
		Form: forma, Balance: bilansGotowy, Change: zmiana, ActionId: stan.czynnosc,
	}, nil
}

// tabelaMaScalenia mówi, czy tabela niesie choć jedno scalenie.
func tabelaMaScalenia(tabela *shared.StudioDocumentTable) bool {
	for _, komorka := range tabela.Cells {
		if komorka.Merged != nil && *komorka.Merged {
			return true
		}
		if komorka.RowSpan != nil && *komorka.RowSpan > 1 {
			return true
		}
		if komorka.ColumnSpan != nil && *komorka.ColumnSpan > 1 {
			return true
		}
	}
	return false
}

// tabelaMaPostacKomorek mówi, czy komórki niosą postać, której tekst nie
// przeniesie.
func tabelaMaPostacKomorek(tabela *shared.StudioDocumentTable) bool {
	for _, komorka := range tabela.Cells {
		if komorka.Border != nil || komorka.ShadingColor != nil ||
			komorka.VerticalAlign != nil || komorka.Character != nil ||
			komorka.Paragraph != nil {
			return true
		}
	}
	return tabela.Border != nil || tabela.StyleName != nil
}
