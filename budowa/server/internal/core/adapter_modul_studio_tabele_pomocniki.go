// Odpowiedzialność pliku: rachunek tabeli dokumentu — siatka komórek, wstawienie
// i usunięcie wiersza oraz kolumny, scalenie i podział komórek, przeliczenie
// szerokości kolumn, zamiana tabeli na tekst i tekstu na tabelę.
//
// ── Gdzie stoi prawda o tabeli ──────────────────────────────────────────────
// Tabela siedzi w drzewie postaci (`stanPostaci.forma.Tables`), a jej miejsce
// w kolejności czytania niesie blok rodzaju `tabela` ze wskazaniem tabeli.
// Osobnego wiersza w bazie tabela nie ma i mieć nie powinna: nie pyta się „które
// tabele są nieświeże", tak jak pyta się o przypisy — tabelę czyta się razem
// z dokumentem, więc jedzie jego drzewem.
//
// ── Dlaczego siatka jest pełna, a scalenie znaczy się flagą ─────────────────
// Komórki trzymane rzadko (tylko te wypełnione) wymagałyby przy każdym odczycie
// zgadywania, czy komórki nie ma, bo jest pusta, czy bo została wchłonięta
// scaleniem. Dlatego siatka jest PEŁNA: wiersze razy kolumny, a komórka
// wchłonięta scaleniem niesie `merged` i zostaje na swoim miejscu. Wstawienie
// wiersza w środek scalenia ma wtedy co przesunąć, a nie zgaduje.
//
// ── Szerokości nigdy zerowe ─────────────────────────────────────────────────
// Wymaganie zlecenia mówi wprost: tabela po scaleniu komórek ma mieć POLICZONE
// szerokości, nie zerowe. Dlatego każda czynność zmieniająca budowę woła
// `tabelaPrzeliczSzerokosci`, a nie zostawia tablicy krótszej niż liczba kolumn.
package core

import (
	"sort"
	"strconv"
	"strings"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// tabelaSzerokoscNosnikaDomyslna to szerokość arkusza A4 w milimetrach — miara
// dokumentu, którego nastawy strony nie nazywają nośnika ani wymiaru.
//
// Nośnik NAZWANY bierze wymiary ze wspólnego wykazu rdzenia
// (`nosniki_druku_wspolne.go`), a nie z drugiej tablicy w tym pliku: dwa wykazy
// rozmiarów arkusza dałyby tabelę mierzoną szerokością innego nośnika, niż
// Operator ustawił.
const tabelaSzerokoscNosnikaDomyslna = 210.0

// tabelaMarginesDomyslny to margines pisma urzędowego w milimetrach.
const tabelaMarginesDomyslny = 25.0

// tabelaSzerokoscKolumnyNajmniejsza chroni kolumnę przed szerokością, przy
// której nie zmieści się ani jedna litera.
const tabelaSzerokoscKolumnyNajmniejsza = 5.0

// tabelaGranicaWymiaru chroni przed tabelą, której nikt nie zamawiał: siatka
// stu tysięcy komórek zajęłaby pamięć rdzenia i nie weszłaby na żadną stronę.
const (
	tabelaGranicaWierszy = 2000
	tabelaGranicaKolumn  = 128
)

// tabelaBladWskazania nazywa brak po stronie żądania.
func tabelaBladWskazania(powod string) error {
	return bladWskazaniaStudio(powod)
}

// tabelaBladNieznanej mówi wprost, której tabeli w dokumencie nie ma.
func tabelaBladNieznanej(kod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
		"moduł Studio, tabele: dokument nie ma tabeli o wskazaniu "+kod))
}

// tabelaSzerokoscTekstu liczy szerokość kolumny tekstu dokumentu — tyle miejsca
// ma tabela, która nie podaje własnej szerokości.
func tabelaSzerokoscTekstu(forma *shared.StudioDocumentForm) float64 {
	szerokosc := tabelaSzerokoscNosnikaDomyslna
	lewy, prawy := tabelaMarginesDomyslny, tabelaMarginesDomyslny
	if forma != nil && forma.PageSetup != nil {
		nastawy := forma.PageSetup
		poziomo := nastawy.Orientation != nil &&
			*nastawy.Orientation == shared.StudioPageOrientationPozioma
		switch {
		case nastawy.WidthMm != nil && *nastawy.WidthMm > 0:
			szerokosc = *nastawy.WidthMm
		case nastawy.PageSize != nil:
			// Nośnik nazwany: wymiary z wykazu wspólnego. Orientacja pozioma
			// zamienia wymiary arkusza — bez tego tabela w załączniku poziomym
			// mierzyłaby się szerokością strony pionowej.
			if nosnik, jest := nosnikDrukuONazwie(*nastawy.PageSize); jest {
				szerokosc = nosnik.SzerokoscMm
				if poziomo {
					szerokosc = nosnik.WysokoscMm
				}
			}
		case poziomo:
			// Orientacja pozioma bez nazwy nośnika: arkusz domyślny położony.
			if a4, jest := nosnikDrukuONazwie("A4"); jest {
				szerokosc = a4.WysokoscMm
			}
		}
		if nastawy.MarginLeft != nil {
			lewy = float64(*nastawy.MarginLeft)
		}
		if nastawy.MarginRight != nil {
			prawy = float64(*nastawy.MarginRight)
		}
	}
	uzyteczna := szerokosc - lewy - prawy
	if uzyteczna < 10 {
		uzyteczna = 10
	}
	return uzyteczna
}

// tabelaZnajdz odnajduje tabelę dokumentu po wskazaniu.
func tabelaZnajdz(forma *shared.StudioDocumentForm, kod string) (*shared.StudioDocumentTable, error) {
	szukany := strings.TrimSpace(kod)
	if szukany == "" {
		return nil, tabelaBladWskazania("czynność na tabeli bez wskazania tabeli")
	}
	for i := range forma.Tables {
		if forma.Tables[i].Id == szukany {
			return &forma.Tables[i], nil
		}
	}
	return nil, tabelaBladNieznanej(szukany)
}

// tabelaSiatkaPelna uzupełnia siatkę komórek do pełnego prostokąta i porządkuje
// ją po wierszach i kolumnach.
func tabelaSiatkaPelna(tabela *shared.StudioDocumentTable) {
	if tabela.Rows < 0 {
		tabela.Rows = 0
	}
	if tabela.Columns < 0 {
		tabela.Columns = 0
	}
	zastane := map[[2]int]shared.StudioTableCell{}
	for _, komorka := range tabela.Cells {
		if komorka.Row < 0 || komorka.Column < 0 {
			continue
		}
		if komorka.Row >= tabela.Rows || komorka.Column >= tabela.Columns {
			continue
		}
		zastane[[2]int{komorka.Row, komorka.Column}] = komorka
	}
	pelna := make([]shared.StudioTableCell, 0, tabela.Rows*tabela.Columns)
	for wiersz := 0; wiersz < tabela.Rows; wiersz++ {
		for kolumna := 0; kolumna < tabela.Columns; kolumna++ {
			if komorka, jest := zastane[[2]int{wiersz, kolumna}]; jest {
				komorka.Row, komorka.Column = wiersz, kolumna
				pelna = append(pelna, komorka)
				continue
			}
			pelna = append(pelna, shared.StudioTableCell{Row: wiersz, Column: kolumna})
		}
	}
	tabela.Cells = pelna
}

// tabelaKomorka oddaje wskaźnik na komórkę siatki. Siatka jest pełna, więc
// wskaźnik jest zawsze, o ile wskazanie leży w tabeli.
func tabelaKomorka(tabela *shared.StudioDocumentTable, wiersz, kolumna int) *shared.StudioTableCell {
	if wiersz < 0 || kolumna < 0 || wiersz >= tabela.Rows || kolumna >= tabela.Columns {
		return nil
	}
	miejsce := wiersz*tabela.Columns + kolumna
	if miejsce >= len(tabela.Cells) {
		return nil
	}
	return &tabela.Cells[miejsce]
}

// tabelaPrzeliczSzerokosci doprowadza tablicę szerokości do liczby kolumn.
//
// Kolumna bez podanej szerokości bierze resztę miejsca rozdzieloną równo, a
// suma nigdy nie schodzi do zera — zerowa szerokość znaczyłaby kolumnę
// niewidzialną, a to jest ta sama cicha szkoda, przed którą stoi wymaganie
// „szerokości policzone, nie zerowe".
func tabelaPrzeliczSzerokosci(tabela *shared.StudioDocumentTable, szerokoscTekstu float64) {
	if tabela.Columns <= 0 {
		tabela.ColumnWidthsMm = nil
		return
	}
	calosc := szerokoscTekstu
	if tabela.WidthMm != nil && *tabela.WidthMm > 0 {
		calosc = *tabela.WidthMm
	}
	najmniejsza := float64(tabela.Columns) * tabelaSzerokoscKolumnyNajmniejsza
	if calosc < najmniejsza {
		calosc = najmniejsza
	}

	stare := tabela.ColumnWidthsMm
	nowe := make([]float64, tabela.Columns)
	suma := 0.0
	znane := 0
	for i := 0; i < tabela.Columns; i++ {
		if i < len(stare) && stare[i] > 0 {
			nowe[i] = stare[i]
			suma += stare[i]
			znane++
		}
	}
	if znane == 0 {
		rowna := calosc / float64(tabela.Columns)
		for i := range nowe {
			nowe[i] = tabelaZaokraglenieMiary(rowna)
		}
		tabela.ColumnWidthsMm = nowe
		tabela.WidthMm = postacWskaznikMiary(tabelaSumaMiar(nowe))
		return
	}
	if znane < tabela.Columns {
		reszta := calosc - suma
		naKolumne := reszta / float64(tabela.Columns-znane)
		if naKolumne < tabelaSzerokoscKolumnyNajmniejsza {
			naKolumne = tabelaSzerokoscKolumnyNajmniejsza
		}
		for i := range nowe {
			if nowe[i] <= 0 {
				nowe[i] = tabelaZaokraglenieMiary(naKolumne)
			}
		}
	}
	tabela.ColumnWidthsMm = nowe
	tabela.WidthMm = postacWskaznikMiary(tabelaSumaMiar(nowe))
}

// tabelaSumaMiar sumuje szerokości kolumn.
func tabelaSumaMiar(miary []float64) float64 {
	suma := 0.0
	for _, miara := range miary {
		suma += miara
	}
	return tabelaZaokraglenieMiary(suma)
}

// tabelaZaokraglenieMiary ucina miarę do dziesiątych części milimetra — okno
// pokazuje szerokość kolumny liczbą, a nie ułamkiem o piętnastu cyfrach.
func tabelaZaokraglenieMiary(miara float64) float64 {
	zaokraglona := float64(int64(miara*10+0.5)) / 10
	if zaokraglona < tabelaSzerokoscKolumnyNajmniejsza {
		return tabelaSzerokoscKolumnyNajmniejsza
	}
	return zaokraglona
}

// ── Budowa tabeli ───────────────────────────────────────────────────────────

// tabelaWstawWiersze wstawia wiersze przed albo za wskazanym.
func tabelaWstawWiersze(tabela *shared.StudioDocumentTable, wiersz int, przed bool, ile int) error {
	if ile <= 0 {
		ile = 1
	}
	if tabela.Rows+ile > tabelaGranicaWierszy {
		return tabelaBladWskazania("tabela o " + strconv.Itoa(tabela.Rows+ile) +
			" wierszach przekracza granicę " + strconv.Itoa(tabelaGranicaWierszy) +
			" wierszy przyjętą w rdzeniu")
	}
	miejsce := wiersz
	if !przed {
		miejsce = wiersz + 1
	}
	if miejsce < 0 {
		miejsce = 0
	}
	if miejsce > tabela.Rows {
		miejsce = tabela.Rows
	}
	// Wiersz wstawiony w środek scalenia pionowego rozcina to scalenie: dwa
	// wiersze pod jedną komórką, z których jeden nie należy do niczego, byłyby
	// tabelą niespójną.
	tabelaRozetnijScaleniaPionowe(tabela, miejsce)

	nowe := make([]shared.StudioTableCell, 0, len(tabela.Cells)+ile*tabela.Columns)
	for _, komorka := range tabela.Cells {
		if komorka.Row >= miejsce {
			komorka.Row += ile
		}
		nowe = append(nowe, komorka)
	}
	for i := 0; i < ile; i++ {
		for kolumna := 0; kolumna < tabela.Columns; kolumna++ {
			nowe = append(nowe, shared.StudioTableCell{Row: miejsce + i, Column: kolumna})
		}
	}
	tabela.Rows += ile
	tabela.Cells = nowe
	tabelaSiatkaPelna(tabela)
	return nil
}

// tabelaUsunWiersze usuwa wskazane wiersze.
func tabelaUsunWiersze(tabela *shared.StudioDocumentTable, wiersz, ile int) error {
	if ile <= 0 {
		ile = 1
	}
	if wiersz < 0 || wiersz >= tabela.Rows {
		return tabelaBladWskazania("usunięcie wiersza " + strconv.Itoa(wiersz) +
			", a tabela ma wierszy " + strconv.Itoa(tabela.Rows))
	}
	if wiersz+ile > tabela.Rows {
		ile = tabela.Rows - wiersz
	}
	if ile >= tabela.Rows {
		return tabelaBladWskazania("usunięcie wszystkich wierszy tabeli — tabelę bez ani " +
			"jednego wiersza usuwa się w całości zamianą na tekst (studio.table.convert)")
	}
	tabelaRozetnijScaleniaPionowe(tabela, wiersz)
	tabelaRozetnijScaleniaPionowe(tabela, wiersz+ile)

	nowe := make([]shared.StudioTableCell, 0, len(tabela.Cells))
	for _, komorka := range tabela.Cells {
		switch {
		case komorka.Row >= wiersz && komorka.Row < wiersz+ile:
			continue
		case komorka.Row >= wiersz+ile:
			komorka.Row -= ile
		}
		nowe = append(nowe, komorka)
	}
	tabela.Rows -= ile
	tabela.Cells = nowe
	if tabela.HeaderRows != nil && *tabela.HeaderRows > tabela.Rows {
		tabela.HeaderRows = postacWskaznikLiczby(tabela.Rows)
	}
	tabelaSiatkaPelna(tabela)
	return nil
}

// tabelaWstawKolumny wstawia kolumny przed albo za wskazaną.
func tabelaWstawKolumny(tabela *shared.StudioDocumentTable, kolumna int, przed bool, ile int) error {
	if ile <= 0 {
		ile = 1
	}
	if tabela.Columns+ile > tabelaGranicaKolumn {
		return tabelaBladWskazania("tabela o " + strconv.Itoa(tabela.Columns+ile) +
			" kolumnach przekracza granicę " + strconv.Itoa(tabelaGranicaKolumn) +
			" kolumn przyjętą w rdzeniu")
	}
	miejsce := kolumna
	if !przed {
		miejsce = kolumna + 1
	}
	if miejsce < 0 {
		miejsce = 0
	}
	if miejsce > tabela.Columns {
		miejsce = tabela.Columns
	}
	tabelaRozetnijScaleniaPoziome(tabela, miejsce)

	nowe := make([]shared.StudioTableCell, 0, len(tabela.Cells)+ile*tabela.Rows)
	for _, komorka := range tabela.Cells {
		if komorka.Column >= miejsce {
			komorka.Column += ile
		}
		nowe = append(nowe, komorka)
	}
	for wiersz := 0; wiersz < tabela.Rows; wiersz++ {
		for i := 0; i < ile; i++ {
			nowe = append(nowe, shared.StudioTableCell{Row: wiersz, Column: miejsce + i})
		}
	}
	tabela.Columns += ile

	// Szerokości idą za kolumnami: kolumna wstawiona bez szerokości byłaby
	// kolumną zerowej szerokości, czyli niewidzialną.
	szerokosci := make([]float64, 0, tabela.Columns)
	szerokosci = append(szerokosci, tabelaWytnijMiary(tabela.ColumnWidthsMm, 0, miejsce)...)
	for i := 0; i < ile; i++ {
		szerokosci = append(szerokosci, 0)
	}
	szerokosci = append(szerokosci, tabelaWytnijMiary(tabela.ColumnWidthsMm, miejsce,
		len(tabela.ColumnWidthsMm))...)
	tabela.ColumnWidthsMm = szerokosci
	// Szerokości kolumn zastanych zostają, a kolumna nowa dostaje resztę
	// miejsca przy przeliczeniu. Szerokość CAŁEJ tabeli idzie do przeliczenia od
	// nowa, bo tabela rozszerzona poza kolumnę tekstu nie weszłaby na stronę.
	tabela.WidthMm = nil
	tabelaSiatkaPelna(tabela)
	return nil
}

// tabelaUsunKolumny usuwa wskazane kolumny.
func tabelaUsunKolumny(tabela *shared.StudioDocumentTable, kolumna, ile int) error {
	if ile <= 0 {
		ile = 1
	}
	if kolumna < 0 || kolumna >= tabela.Columns {
		return tabelaBladWskazania("usunięcie kolumny " + strconv.Itoa(kolumna) +
			", a tabela ma kolumn " + strconv.Itoa(tabela.Columns))
	}
	if kolumna+ile > tabela.Columns {
		ile = tabela.Columns - kolumna
	}
	if ile >= tabela.Columns {
		return tabelaBladWskazania("usunięcie wszystkich kolumn tabeli — tabelę bez ani " +
			"jednej kolumny usuwa się w całości zamianą na tekst (studio.table.convert)")
	}
	tabelaRozetnijScaleniaPoziome(tabela, kolumna)
	tabelaRozetnijScaleniaPoziome(tabela, kolumna+ile)

	nowe := make([]shared.StudioTableCell, 0, len(tabela.Cells))
	for _, komorka := range tabela.Cells {
		switch {
		case komorka.Column >= kolumna && komorka.Column < kolumna+ile:
			continue
		case komorka.Column >= kolumna+ile:
			komorka.Column -= ile
		}
		nowe = append(nowe, komorka)
	}
	szerokosci := make([]float64, 0, tabela.Columns)
	szerokosci = append(szerokosci, tabelaWytnijMiary(tabela.ColumnWidthsMm, 0, kolumna)...)
	szerokosci = append(szerokosci, tabelaWytnijMiary(tabela.ColumnWidthsMm, kolumna+ile,
		len(tabela.ColumnWidthsMm))...)
	tabela.Columns -= ile
	tabela.Cells = nowe
	tabela.ColumnWidthsMm = szerokosci
	// Szerokość całej tabeli maleje wraz z kolumnami — tabela po usunięciu
	// kolumny nie ma prawa rosnąć w pozostałych.
	tabela.WidthMm = nil
	tabelaSiatkaPelna(tabela)
	return nil
}

// tabelaWytnijMiary oddaje wycinek tablicy szerokości, znosząc wskazania poza
// jej długością.
func tabelaWytnijMiary(miary []float64, od, do int) []float64 {
	if od < 0 {
		od = 0
	}
	if do > len(miary) {
		do = len(miary)
	}
	if od >= do {
		return nil
	}
	return append([]float64(nil), miary[od:do]...)
}

// tabelaScalKomorki scala prostokąt komórek w jedną.
//
// Treść komórek wchłoniętych nie przepada: wchodzi do komórki wiodącej
// oddzielona odstępem. Cicha utrata treści przy scaleniu byłaby dokładnie tą
// szkodą, po której Operator nie wie, kiedy dokument się rozjechał.
func tabelaScalKomorki(tabela *shared.StudioDocumentTable, wiersz, kolumna,
	wierszy, kolumn int) error {

	if wierszy <= 0 {
		wierszy = 1
	}
	if kolumn <= 0 {
		kolumn = 1
	}
	if wierszy == 1 && kolumn == 1 {
		return tabelaBladWskazania("scalenie jednej komórki z niczym — scalenie wymaga " +
			"objęcia więcej niż jednego wiersza albo więcej niż jednej kolumny")
	}
	wiodaca := tabelaKomorka(tabela, wiersz, kolumna)
	if wiodaca == nil {
		return tabelaBladWskazania("scalenie od komórki (" + strconv.Itoa(wiersz) + ", " +
			strconv.Itoa(kolumna) + "), której tabela nie ma")
	}
	if wiersz+wierszy > tabela.Rows || kolumna+kolumn > tabela.Columns {
		return tabelaBladWskazania("scalenie wychodzi poza tabelę: tabela ma " +
			strconv.Itoa(tabela.Rows) + " wierszy i " + strconv.Itoa(tabela.Columns) +
			" kolumn")
	}

	czesci := make([]string, 0, wierszy*kolumn)
	if wiodaca.Text != nil && strings.TrimSpace(*wiodaca.Text) != "" {
		czesci = append(czesci, *wiodaca.Text)
	}
	for w := wiersz; w < wiersz+wierszy; w++ {
		for k := kolumna; k < kolumna+kolumn; k++ {
			if w == wiersz && k == kolumna {
				continue
			}
			komorka := tabelaKomorka(tabela, w, k)
			if komorka == nil {
				continue
			}
			if komorka.Text != nil && strings.TrimSpace(*komorka.Text) != "" {
				czesci = append(czesci, *komorka.Text)
			}
			pusta := shared.StudioTableCell{Row: w, Column: k, Merged: postacWskaznikPrawdy(true)}
			*komorka = pusta
		}
	}
	wiodaca.RowSpan = postacWskaznikLiczby(wierszy)
	wiodaca.ColumnSpan = postacWskaznikLiczby(kolumn)
	wiodaca.Merged = nil
	if len(czesci) > 0 {
		wiodaca.Text = postacWskaznikTekstu(strings.Join(czesci, " "))
	}
	return nil
}

// tabelaPodzielKomorke znosi scalenie komórki.
func tabelaPodzielKomorke(tabela *shared.StudioDocumentTable, wiersz, kolumna,
	wierszy, kolumn int) error {

	wiodaca := tabelaKomorka(tabela, wiersz, kolumna)
	if wiodaca == nil {
		return tabelaBladWskazania("podział komórki (" + strconv.Itoa(wiersz) + ", " +
			strconv.Itoa(kolumna) + "), której tabela nie ma")
	}
	objeteWiersze, objeteKolumny := 1, 1
	if wiodaca.RowSpan != nil && *wiodaca.RowSpan > 1 {
		objeteWiersze = *wiodaca.RowSpan
	}
	if wiodaca.ColumnSpan != nil && *wiodaca.ColumnSpan > 1 {
		objeteKolumny = *wiodaca.ColumnSpan
	}
	if objeteWiersze == 1 && objeteKolumny == 1 {
		// Podział komórki niescalonej znaczy w pakiecie biurowym wstawienie
		// kolumn albo wierszy wewnątrz niej; rdzeń tego nie udaje i mówi wprost,
		// którą czynnością się to robi.
		return tabelaBladWskazania("komórka (" + strconv.Itoa(wiersz) + ", " +
			strconv.Itoa(kolumna) + ") nie jest scalona, więc nie ma czego podzielić; " +
			"tabelę rozdziela się wstawieniem wiersza albo kolumny " +
			"(operation insertRow, insertColumn)")
	}
	if wierszy > objeteWiersze || kolumn > objeteKolumny {
		return tabelaBladWskazania("podział na " + strconv.Itoa(wierszy) + " wierszy i " +
			strconv.Itoa(kolumn) + " kolumn, a scalenie obejmuje " +
			strconv.Itoa(objeteWiersze) + " wierszy i " + strconv.Itoa(objeteKolumny) +
			" kolumn")
	}
	for w := wiersz; w < wiersz+objeteWiersze; w++ {
		for k := kolumna; k < kolumna+objeteKolumny; k++ {
			komorka := tabelaKomorka(tabela, w, k)
			if komorka == nil {
				continue
			}
			komorka.Merged = nil
			if w == wiersz && k == kolumna {
				komorka.RowSpan, komorka.ColumnSpan = nil, nil
			}
		}
	}
	return nil
}

// tabelaRozetnijScaleniaPionowe znosi scalenia przekraczające granicę wiersza.
func tabelaRozetnijScaleniaPionowe(tabela *shared.StudioDocumentTable, granica int) {
	if granica <= 0 || granica >= tabela.Rows {
		return
	}
	for wiersz := 0; wiersz < granica; wiersz++ {
		for kolumna := 0; kolumna < tabela.Columns; kolumna++ {
			komorka := tabelaKomorka(tabela, wiersz, kolumna)
			if komorka == nil || komorka.RowSpan == nil || *komorka.RowSpan <= 1 {
				continue
			}
			if wiersz+*komorka.RowSpan > granica {
				kolumn := 1
				if komorka.ColumnSpan != nil && *komorka.ColumnSpan > 1 {
					kolumn = *komorka.ColumnSpan
				}
				_ = tabelaPodzielKomorke(tabela, wiersz, kolumna, *komorka.RowSpan, kolumn)
			}
		}
	}
}

// tabelaRozetnijScaleniaPoziome znosi scalenia przekraczające granicę kolumny.
func tabelaRozetnijScaleniaPoziome(tabela *shared.StudioDocumentTable, granica int) {
	if granica <= 0 || granica >= tabela.Columns {
		return
	}
	for wiersz := 0; wiersz < tabela.Rows; wiersz++ {
		for kolumna := 0; kolumna < granica; kolumna++ {
			komorka := tabelaKomorka(tabela, wiersz, kolumna)
			if komorka == nil || komorka.ColumnSpan == nil || *komorka.ColumnSpan <= 1 {
				continue
			}
			if kolumna+*komorka.ColumnSpan > granica {
				wierszy := 1
				if komorka.RowSpan != nil && *komorka.RowSpan > 1 {
					wierszy = *komorka.RowSpan
				}
				_ = tabelaPodzielKomorke(tabela, wiersz, kolumna, wierszy, *komorka.ColumnSpan)
			}
		}
	}
}

// ── Sortowanie ──────────────────────────────────────────────────────────────

// tabelaSortujWiersze sortuje zawartość tabeli po wskazanej kolumnie.
//
// Sortowanie idzie CAŁYMI wierszami, nie samą kolumną: przestawienie jednej
// kolumny rozerwałoby wiersze i tabela mówiłaby co innego niż przed czynnością.
func tabelaSortujWiersze(tabela *shared.StudioDocumentTable, kolumna int,
	malejaco, liczbowo, pomijajNaglowek bool) (int, error) {

	if kolumna < 0 || kolumna >= tabela.Columns {
		return 0, tabelaBladWskazania("sortowanie po kolumnie " + strconv.Itoa(kolumna) +
			", a tabela ma kolumn " + strconv.Itoa(tabela.Columns))
	}
	pierwszy := 0
	if pomijajNaglowek {
		pierwszy = 1
		if tabela.HeaderRows != nil && *tabela.HeaderRows > 0 {
			pierwszy = *tabela.HeaderRows
		}
	}
	if pierwszy >= tabela.Rows {
		return 0, tabelaBladWskazania("sortowanie tabeli, w której poza wierszem " +
			"nagłówkowym nie ma ani jednego wiersza zawartości")
	}
	// Wiersz objęty scaleniem pionowym nie da się przestawić bez rozerwania
	// scalenia — sortowanie odmawia nazwanie, zamiast wywracać budowę tabeli.
	for wiersz := pierwszy; wiersz < tabela.Rows; wiersz++ {
		for k := 0; k < tabela.Columns; k++ {
			komorka := tabelaKomorka(tabela, wiersz, k)
			if komorka == nil {
				continue
			}
			if komorka.RowSpan != nil && *komorka.RowSpan > 1 {
				return 0, tabelaBladWskazania("sortowanie tabeli ze scaleniem pionowym " +
					"w wierszu " + strconv.Itoa(wiersz) + " — scalenie trzeba najpierw " +
					"znieść (operation splitCell), inaczej wiersze rozerwałyby się")
			}
			if komorka.Merged != nil && *komorka.Merged {
				return 0, tabelaBladWskazania("sortowanie tabeli ze scaleniem pionowym " +
					"w wierszu " + strconv.Itoa(wiersz) + " — scalenie trzeba najpierw " +
					"znieść (operation splitCell)")
			}
		}
	}

	wiersze := make([][]shared.StudioTableCell, 0, tabela.Rows-pierwszy)
	for wiersz := pierwszy; wiersz < tabela.Rows; wiersz++ {
		komorki := make([]shared.StudioTableCell, 0, tabela.Columns)
		for k := 0; k < tabela.Columns; k++ {
			if komorka := tabelaKomorka(tabela, wiersz, k); komorka != nil {
				komorki = append(komorki, *komorka)
			}
		}
		wiersze = append(wiersze, komorki)
	}
	zestawiacz := tabelaZestawiaczPolski()
	sort.SliceStable(wiersze, func(i, j int) bool {
		lewy := tabelaTrescKomorki(wiersze[i], kolumna)
		prawy := tabelaTrescKomorki(wiersze[j], kolumna)
		mniejszy := false
		if liczbowo {
			lewaLiczba, lewaJest := tabelaLiczbaZTekstu(lewy)
			prawaLiczba, prawaJest := tabelaLiczbaZTekstu(prawy)
			switch {
			case lewaJest && prawaJest:
				mniejszy = lewaLiczba < prawaLiczba
			case lewaJest != prawaJest:
				// Wartość nieliczbowa idzie na koniec — inaczej „bez danych"
				// wypadałoby przed najmniejszą liczbą i wyglądało na zero.
				mniejszy = lewaJest
			default:
				mniejszy = tabelaMniejszyTekst(zestawiacz, lewy, prawy)
			}
		} else {
			mniejszy = tabelaMniejszyTekst(zestawiacz, lewy, prawy)
		}
		if malejaco {
			return !mniejszy
		}
		return mniejszy
	})

	for i, komorki := range wiersze {
		docelowy := pierwszy + i
		for k := 0; k < tabela.Columns && k < len(komorki); k++ {
			komorka := tabelaKomorka(tabela, docelowy, k)
			if komorka == nil {
				continue
			}
			przeniesiona := komorki[k]
			przeniesiona.Row, przeniesiona.Column = docelowy, k
			*komorka = przeniesiona
		}
	}
	return len(wiersze), nil
}

// tabelaTrescKomorki oddaje treść komórki wiersza.
func tabelaTrescKomorki(wiersz []shared.StudioTableCell, kolumna int) string {
	if kolumna < 0 || kolumna >= len(wiersz) {
		return ""
	}
	if wiersz[kolumna].Text == nil {
		return ""
	}
	return *wiersz[kolumna].Text
}

// tabelaPorzadekPolski jest KOLACJĄ pisma polskiego — tablicą zestawiania,
// którą sortowanie tabeli układa wyrazy.
//
// ── Dlaczego biblioteka, a nie własna tablica znaków ────────────────────────
// Porównanie napisów bajt po bajcie stawia „ł" za „z", bo tak leżą one
// w Unikodzie — wykaz nazwisk wychodził wtedy z Łukasiewiczem na końcu, za
// Zawadzkim. Własna tablica znaków polskich rozwiązałaby to na dziesięć minut
// i rozjechała się na pierwszym wyrazie z „ﬁ", z apostrofem albo z cyfrą.
// `golang.org/x/text/collate` niesie tablice CLDR i stoi w drzewie od dawna —
// funkcja jest więc wzięta z biblioteki, nie napisana od nowa.
//
// Zestawiacz nie jest bezpieczny dla wielu wątków (trzyma bufor rachunku), więc
// nie stoi tu jako zmienna wspólna: każde sortowanie zakłada własny
// (`tabelaZestawiaczPolski`). Koszt złożenia jest jednorazowy na sortowanie, a
// nie na parę wyrazów.
//
// Siła porównania jest DRUGORZĘDNA (`collate.Loose` byłoby za mało, pełna za
// wiele): różnica wielkości liter nie decyduje o kolejności — Operator sortujący
// wykaz nazwisk nie spodziewa się, że „ćma" wyprzedzi „Dom" tylko dlatego, że
// jedno zaczyna się małą literą — a różnica znaku diakrytycznego decyduje, bo
// „laska" i „łaska" są dwoma różnymi wyrazami i mają stanąć osobno.
var tabelaPorzadekPolski = language.Polish

// tabelaZestawiaczPolski składa zestawiacz pisma polskiego do jednego sortowania.
func tabelaZestawiaczPolski() *collate.Collator {
	return collate.New(tabelaPorzadekPolski, collate.IgnoreCase)
}

// tabelaMniejszyTekst porównuje treść komórek porządkiem alfabetycznym pisma
// polskiego, bez względu na wielkość liter.
//
// Zestawiacz podaje się argumentem, bo jeden sort woła tę funkcję kilkadziesiąt
// razy, a zakładanie zestawiacza przy każdym porównaniu byłoby rachunkiem
// wykonanym od nowa dla każdej pary wyrazów.
func tabelaMniejszyTekst(zestawiacz *collate.Collator, lewy, prawy string) bool {
	return zestawiacz.CompareString(strings.TrimSpace(lewy), strings.TrimSpace(prawy)) < 0
}

// tabelaLiczbaZTekstu odczytuje liczbę z treści komórki, znosząc odstępy
// tysięczne i przyjmując przecinek dziesiętny pisma polskiego.
func tabelaLiczbaZTekstu(tekst string) (float64, bool) {
	oczyszczony := strings.TrimSpace(tekst)
	if oczyszczony == "" {
		return 0, false
	}
	oczyszczony = strings.ReplaceAll(oczyszczony, " ", "")
	oczyszczony = strings.ReplaceAll(oczyszczony, " ", "")
	oczyszczony = strings.ReplaceAll(oczyszczony, ",", ".")
	liczba, err := strconv.ParseFloat(oczyszczony, 64)
	if err != nil {
		return 0, false
	}
	return liczba, true
}

// ── Zamiana tabeli i tekstu ─────────────────────────────────────────────────

// tabelaTekstZTabeli składa tekst z zawartości tabeli: wiersz na wiersz,
// komórki rozdzielone znakiem rozdzielającym.
func tabelaTekstZTabeli(tabela *shared.StudioDocumentTable, rozdzielnik string) []string {
	wiersze := make([]string, 0, tabela.Rows)
	for wiersz := 0; wiersz < tabela.Rows; wiersz++ {
		czesci := make([]string, 0, tabela.Columns)
		for kolumna := 0; kolumna < tabela.Columns; kolumna++ {
			komorka := tabelaKomorka(tabela, wiersz, kolumna)
			if komorka == nil || komorka.Merged != nil && *komorka.Merged {
				czesci = append(czesci, "")
				continue
			}
			czesci = append(czesci, wartoscTekstu(komorka.Text))
		}
		wiersze = append(wiersze, strings.Join(czesci, rozdzielnik))
	}
	return wiersze
}

// tabelaZTekstu składa tabelę z wierszy tekstu.
func tabelaZTekstu(wiersze []string, rozdzielnik string) (shared.StudioDocumentTable, error) {
	if len(wiersze) == 0 {
		return shared.StudioDocumentTable{}, tabelaBladWskazania(
			"zamiana tekstu na tabelę na zakresie, który nie niesie ani jednego wiersza")
	}
	rozbite := make([][]string, 0, len(wiersze))
	kolumn := 0
	for _, wiersz := range wiersze {
		czesci := strings.Split(wiersz, rozdzielnik)
		if len(czesci) > kolumn {
			kolumn = len(czesci)
		}
		rozbite = append(rozbite, czesci)
	}
	if kolumn == 0 {
		kolumn = 1
	}
	if kolumn > tabelaGranicaKolumn {
		return shared.StudioDocumentTable{}, tabelaBladWskazania(
			"zamiana tekstu na tabelę o " + strconv.Itoa(kolumn) + " kolumnach przekracza " +
				"granicę " + strconv.Itoa(tabelaGranicaKolumn) + " kolumn; sprawdź znak " +
				"rozdzielający (pole separator)")
	}
	tabela := shared.StudioDocumentTable{
		Id: nowyIdentyfikator(przedrostekTabeliPostaci), Rows: len(rozbite), Columns: kolumn,
	}
	tabelaSiatkaPelna(&tabela)
	for wiersz, czesci := range rozbite {
		for kolumna, tresc := range czesci {
			komorka := tabelaKomorka(&tabela, wiersz, kolumna)
			if komorka == nil {
				continue
			}
			komorka.Text = postacWskaznikTekstu(tresc)
		}
	}
	return tabela, nil
}

// ── Bloki drzewa ────────────────────────────────────────────────────────────

// tabelaWstawBlokWMiejscu wstawia blok nietekstowy w miejsce wskazane znakiem
// treści, a nie na koniec dokumentu.
//
// Rachunek stoi tutaj i woła go także obszar wstawień: dwa liczenia miejsca
// bloku rozjechałyby tabelę wstawioną w akapicie trzecim z obrazem wstawionym
// w tym samym miejscu.
func tabelaWstawBlokWMiejscu(forma *shared.StudioDocumentForm,
	blok shared.StudioDocumentBlock, miejsce int) {

	postacPrzeliczZakresy(forma)
	wstawiony := false
	nowe := make([]shared.StudioDocumentBlock, 0, len(forma.Blocks)+1)
	for _, biezacy := range forma.Blocks {
		if !wstawiony && postacBlokNiesieTekst(biezacy) && biezacy.RangeEnd != nil &&
			*biezacy.RangeEnd >= miejsce {

			// Blok wchodzi ZA akapitem, w którym stoi miejsce wstawienia:
			// tabela wstawiona w środku akapitu rozcinałaby zdanie, a tego
			// pakiet biurowy też nie robi.
			nowe = append(nowe, biezacy, blok)
			wstawiony = true
			continue
		}
		nowe = append(nowe, biezacy)
	}
	if !wstawiony {
		nowe = append(nowe, blok)
	}
	forma.Blocks = nowe
	postacPrzeliczZakresy(forma)
}

// tabelaUsunBlokTabeli zdejmuje z drzewa blok wskazujący tabelę.
func tabelaUsunBlokTabeli(forma *shared.StudioDocumentForm, kod string) {
	nowe := make([]shared.StudioDocumentBlock, 0, len(forma.Blocks))
	for _, blok := range forma.Blocks {
		if blok.Kind == blokPostaciTabela && blok.TableId != nil && *blok.TableId == kod {
			continue
		}
		nowe = append(nowe, blok)
	}
	forma.Blocks = nowe
}

// tabelaMiejsceTabeli oddaje miejsce tabeli w treści w znakach — spis tabel
// i podpisy potrzebują wiedzieć, gdzie tabela stoi.
func tabelaMiejsceTabeli(forma *shared.StudioDocumentForm, kod string) int {
	for _, blok := range forma.Blocks {
		if blok.Kind == blokPostaciTabela && blok.TableId != nil && *blok.TableId == kod {
			if blok.RangeStart != nil {
				return *blok.RangeStart
			}
			return 0
		}
	}
	return 0
}
