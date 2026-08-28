// Plik liczy pola dokumentu: numer strony, liczbę stron, datę, godzinę, tytuł i właściwości dokumentu, pole obliczane oraz pole szablonu, przez czynności wstawienia, wykazu i odświeżenia.
package core

import (
	"context"
	"strconv"
	"strings"
	"time"
	"unicode"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// WstawPole wstawia pole dokumentu i od razu liczy jego wartość, żeby odpowiedź nie zostawiała pustego miejsca w treści (`studio.field.insert`).
func (a *adapterStudia) WstawPole(ctx context.Context,
	z shared.StudioFieldInsertRequest) (shared.StudioFieldInsertResponse, error) {

	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioFieldInsertResponse{}, err
	}
	autor := postacAutor(z.Author)
	if err := poleSprawdzRodzaj(z.Kind); err != nil {
		return shared.StudioFieldInsertResponse{}, err
	}
	if err := poleSprawdzWymagania(z); err != nil {
		return shared.StudioFieldInsertResponse{}, err
	}

	dlugosc := postacDlugosc(&stan.forma)
	miejsce := z.Offset
	if miejsce < 0 {
		miejsce = 0
	}
	if miejsce > dlugosc {
		miejsce = dlugosc
	}
	if odcinki, pominiete := postacOdcinkiDozwolone(&stan.forma, miejsce, miejsce, autor); len(odcinki) == 0 {
		return shared.StudioFieldInsertResponse{}, bladWskazaniaStudio(
			"wstawienie pola " + postacZapisZakresu(miejsce, miejsce) +
				" zatrzymane przez blokadę fragmentu: " + postacNazwaBlokad(pominiete))
	}

	pole := shared.StudioDocumentField{
		Id: nowyIdentyfikator(przedrostekPolaPostaci), Kind: z.Kind,
		AnchorOffset: postacWskaznikLiczby(miejsce), Format: z.Format,
		Expression: z.Expression, PropertyName: z.PropertyName,
	}
	// Pole wchodzi policzone od razu, inaczej dokument pokazałby puste miejsce mimo odpowiedzi „wstawiono”
	strony, stron, err := aparatStronyAkapitow(postacTekstFormy(&stan.forma), &stan.forma)
	if err != nil {
		return shared.StudioFieldInsertResponse{}, err
	}
	bilans := shared.StudioActionBalance{Skipped: []shared.StudioSkippedItem{}}
	wartosc, pominiecie := a.polePolicz(ctx, stan, pole, strony, stron)
	if pominiecie != nil {
		bilans.Skipped = append(bilans.Skipped, *pominiecie)
		pole.Stale = postacWskaznikPrawdy(true)
	} else {
		pole.Value = postacWskaznikTekstu(wartosc)
		pole.Stale = postacWskaznikPrawdy(false)
	}

	skladnica, err := a.postacSkladnica()
	if err != nil {
		return shared.StudioFieldInsertResponse{}, err
	}
	zapisane, err := skladnica.ZapiszPoleDokumentu(ctx, poleDoWiersza(stan.dokument.ID, pole))
	if err != nil {
		return shared.StudioFieldInsertResponse{}, bladStudio(err)
	}
	if err := a.postacWczytajWiersze(ctx, stan); err != nil {
		return shared.StudioFieldInsertResponse{}, err
	}

	bilans.Applied = 1
	bilans.Note = postacWskaznikTekstu(poleNazwaRodzaju(z.Kind) + " wstawione " +
		postacZapisZakresu(miejsce, miejsce) + "; wartość policzona: " +
		poleZapisWartosci(zapisane.Wartosc))
	stan.opisCzynnosci = "wstawienie pola: " + poleNazwaRodzaju(z.Kind)

	forma, bilansGotowy, zmiana, err := a.postacZakoncz(ctx, stan, autor,
		shared.StudioChangeKindWstawienie, shared.StudioActionKindObjectChange,
		miejsce, miejsce, bilans)
	if err != nil {
		return shared.StudioFieldInsertResponse{}, err
	}
	zlozone := postacZlozPola([]dane.PoleDokumentuStudia{zapisane})
	if len(zlozone) == 0 {
		return shared.StudioFieldInsertResponse{}, postacBladZaplecza(
			"pole zapisane, ale nie da się go złożyć do odpowiedzi")
	}
	return shared.StudioFieldInsertResponse{
		Form: forma, Balance: bilansGotowy, Change: zmiana,
		ActionId: stan.czynnosc, Field: zlozone[0],
	}, nil
}

// WykazPol oddaje pola dokumentu, na żądanie zawężone do wskazanego rodzaju pola (`studio.field.list`).
func (a *adapterStudia) WykazPol(ctx context.Context,
	z shared.StudioFieldListRequest) (shared.StudioFieldListResponse, error) {

	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioFieldListResponse{}, err
	}
	if z.Kind != nil {
		if err := poleSprawdzRodzaj(*z.Kind); err != nil {
			return shared.StudioFieldListResponse{}, err
		}
	}
	pola := make([]shared.StudioDocumentField, 0, len(stan.forma.Fields))
	for _, pole := range stan.forma.Fields {
		if z.Kind != nil && pole.Kind != *z.Kind {
			continue
		}
		pola = append(pola, pole)
	}
	return shared.StudioFieldListResponse{Fields: pola}, nil
}

// OdswiezPola liczy wartości pól od nowa, całego dokumentu albo pola wskazanego kodem (`studio.field.refresh`).
func (a *adapterStudia) OdswiezPola(ctx context.Context,
	z shared.StudioFieldRefreshRequest) (shared.StudioFieldRefreshResponse, error) {

	stan, err := a.postacWczytaj(ctx, z.DocumentId)
	if err != nil {
		return shared.StudioFieldRefreshResponse{}, err
	}
	autor := postacAutor(z.Author)
	if len(stan.forma.Fields) == 0 {
		return shared.StudioFieldRefreshResponse{}, bladWskazaniaStudio(
			"odświeżenie pól dokumentu, który nie ma ani jednego pola — pole trzeba " +
				"najpierw wstawić (studio.field.insert)")
	}
	wskazane := strings.TrimSpace(wartoscTekstu(z.FieldId))
	if wskazane != "" {
		if _, jest := poleZnajdz(&stan.forma, wskazane); !jest {
			return shared.StudioFieldRefreshResponse{}, bladWskazaniaStudio(
				"odświeżenie pola „" + wskazane + "”, którego dokument nie ma")
		}
	}

	strony, stron, err := aparatStronyAkapitow(postacTekstFormy(&stan.forma), &stan.forma)
	if err != nil {
		return shared.StudioFieldRefreshResponse{}, err
	}
	skladnica, err := a.postacSkladnica()
	if err != nil {
		return shared.StudioFieldRefreshResponse{}, err
	}

	bilans := shared.StudioActionBalance{Skipped: []shared.StudioSkippedItem{}}
	odswiezone := make([]shared.StudioDocumentField, 0, len(stan.forma.Fields))
	for _, pole := range stan.forma.Fields {
		if wskazane != "" && pole.Id != wskazane {
			continue
		}
		wartosc, pominiecie := a.polePolicz(ctx, stan, pole, strony, stron)
		if pominiecie != nil {
			bilans.Skipped = append(bilans.Skipped, *pominiecie)
			pole.Stale = postacWskaznikPrawdy(true)
		} else {
			pole.Value = postacWskaznikTekstu(wartosc)
			pole.Stale = postacWskaznikPrawdy(false)
			bilans.Applied++
		}
		zapisane, err := skladnica.ZapiszPoleDokumentu(ctx, poleDoWiersza(stan.dokument.ID, pole))
		if err != nil {
			return shared.StudioFieldRefreshResponse{}, bladStudio(err)
		}
		if zlozone := postacZlozPola([]dane.PoleDokumentuStudia{zapisane}); len(zlozone) > 0 {
			odswiezone = append(odswiezone, zlozone[0])
		}
	}
	if len(odswiezone) == 0 {
		return shared.StudioFieldRefreshResponse{}, bladWskazaniaStudio(
			"odświeżenie pól nie objęło ani jednego pola — zawężenie nie pasuje do " +
				"niczego w dokumencie")
	}
	if err := a.postacWczytajWiersze(ctx, stan); err != nil {
		return shared.StudioFieldRefreshResponse{}, err
	}
	bilans.Note = postacWskaznikTekstu("pola odświeżone; policzonych " +
		strconv.Itoa(bilans.Applied) + " z " + strconv.Itoa(len(odswiezone)) +
		", stron dokumentu " + strconv.Itoa(stron))
	stan.opisCzynnosci = "odświeżenie pól dokumentu"

	forma, bilansGotowy, zmiana, err := a.postacZakoncz(ctx, stan, autor,
		shared.StudioChangeKindFormatowanie, shared.StudioActionKindObjectChange,
		0, postacDlugosc(&stan.forma), bilans)
	if err != nil {
		return shared.StudioFieldRefreshResponse{}, err
	}
	return shared.StudioFieldRefreshResponse{
		Form: forma, Balance: bilansGotowy, Change: zmiana,
		ActionId: stan.czynnosc, Fields: odswiezone,
	}, nil
}

// ── Rachunek wartości ───────────────────────────────────────────────────────

// polePolicz liczy wartość pola. Pominięcie zamiast wartości znaczy, że rdzeń
// nie ma z czego policzyć — i mówi to wprost, zamiast wstawiać puste miejsce.
func (a *adapterStudia) polePolicz(ctx context.Context, stan *stanPostaci,
	pole shared.StudioDocumentField, strony []int,
	stron int) (string, *shared.StudioSkippedItem) {

	teraz := time.Now()
	switch pole.Kind {
	case shared.StudioFieldKindPageNumber:
		miejsce := aparatWartoscLiczby(pole.AnchorOffset)
		return strconv.Itoa(aparatStronaMiejsca(&stan.forma, strony, miejsce)), nil

	case shared.StudioFieldKindPageCount:
		return strconv.Itoa(stron), nil

	case shared.StudioFieldKindDate:
		return poleZapisChwili(teraz, wartoscTekstu(pole.Format), "dd.MM.yyyy"), nil

	case shared.StudioFieldKindTime:
		return poleZapisChwili(teraz, wartoscTekstu(pole.Format), "HH:mm"), nil

	case shared.StudioFieldKindDocumentTitle:
		tytul := strings.TrimSpace(wartoscTekstu(stan.dokument.Tytul))
		if tytul == "" {
			return "", poleBrak(pole, "dokument bez tytułu",
				"pole tytułu nie ma czego pokazać, dopóki dokument nie dostanie tytułu "+
					"(studio.document.save, pole title)")
		}
		return tytul, nil

	case shared.StudioFieldKindDocumentAuthor:
		return a.poleAutorDokumentu(ctx, stan, pole)

	case shared.StudioFieldKindDocumentProperty:
		return a.poleWlasciwoscDokumentu(stan, pole, stron)

	case shared.StudioFieldKindCalculated:
		wyrazenie := strings.TrimSpace(wartoscTekstu(pole.Expression))
		if wyrazenie == "" {
			return "", poleBrak(pole, "pole obliczane bez wyrażenia",
				"pole obliczane wymaga wyrażenia (pole expression)")
		}
		wynik, err := poleObliczWyrazenie(wyrazenie, &stan.forma)
		if err != nil {
			return "", poleBrak(pole, "wyrażenia nie da się policzyć", err.Error())
		}
		return poleZapisLiczby(wynik), nil

	case shared.StudioFieldKindTemplateField:
		// Pole szablonu wypełnia się przy zakładaniu z szablonu; odświeżenie nie ma skąd wziąć jego wartości
		if pole.Value != nil && strings.TrimSpace(*pole.Value) != "" {
			return *pole.Value, nil
		}
		return "", poleBrak(pole, "pole szablonu bez podstawionej wartości",
			"wartość pola szablonu „"+wartoscTekstu(pole.PropertyName)+
				"” podstawia się przy zakładaniu dokumentu z szablonu "+
				"(studio.template.apply); odświeżenie pól nie ma jej skąd wziąć")
	}
	return "", poleBrak(pole, "rodzaj pola bez rachunku",
		"rdzeń nie ma czym policzyć pola rodzaju „"+string(pole.Kind)+"”")
}

// poleAutorDokumentu oddaje autora dokumentu jako rodzaj autora ostatniej wersji, Operatora albo model, bo rdzeń nie trzyma imienia i nazwiska autora, tylko rodzaj autora wersji.
func (a *adapterStudia) poleAutorDokumentu(ctx context.Context, stan *stanPostaci,
	pole shared.StudioDocumentField) (string, *shared.StudioSkippedItem) {

	kod := stan.dokument.WersjaBiezacaKod
	if kod == nil || strings.TrimSpace(*kod) == "" {
		return "", poleBrak(pole, "dokument bez wersji",
			"rdzeń wie o autorze dokumentu tyle, ile niesie autor jego wersji; dokument "+
				"bez ani jednej wersji nie ma czego pokazać")
	}
	wersja, err := a.repozytorium.Wersja(ctx, strings.TrimSpace(*kod))
	if err != nil {
		return "", poleBrak(pole, "wersji dokumentu nie da się odczytać", err.Error())
	}
	if wersja.Autor == nil || strings.TrimSpace(*wersja.Autor) == "" {
		return "", poleBrak(pole, "wersja bez zapisanego autora",
			"wersja bieżąca dokumentu nie niesie autora")
	}
	if *wersja.Autor == string(shared.StudioAuthorModel) {
		return "model", nil
	}
	return "Operator", nil
}

// poleWlasciwoscDokumentu oddaje wartość wskazanej właściwości dokumentu, na przykład tytuł, liczbę stron albo liczbę słów.
func (a *adapterStudia) poleWlasciwoscDokumentu(stan *stanPostaci,
	pole shared.StudioDocumentField, stron int) (string, *shared.StudioSkippedItem) {

	nazwa := strings.ToLower(strings.TrimSpace(wartoscTekstu(pole.PropertyName)))
	if nazwa == "" {
		return "", poleBrak(pole, "właściwość dokumentu bez nazwy",
			"pole właściwości wymaga jej nazwy (pole propertyName); wykaz: "+
				poleWykazWlasciwosci())
	}
	tresc := postacTekstFormy(&stan.forma)
	switch nazwa {
	case "tytuł", "tytul":
		tytul := strings.TrimSpace(wartoscTekstu(stan.dokument.Tytul))
		if tytul == "" {
			return "", poleBrak(pole, "dokument bez tytułu",
				"właściwość „tytuł” nie ma czego pokazać, dopóki dokument nie ma tytułu")
		}
		return tytul, nil
	case "wskazanie", "kod":
		return stan.dokument.Kod, nil
	case "okno":
		return stan.dokument.Okno, nil
	case "format":
		return string(stan.dokument.Format), nil
	case "utworzono":
		return stan.dokument.Utworzono, nil
	case "zaktualizowano":
		return stan.dokument.Zaktualizowano, nil
	case "wersja":
		if stan.dokument.WersjaBiezacaKod == nil {
			return "", poleBrak(pole, "dokument bez wersji",
				"właściwość „wersja” nie ma czego pokazać: dokument nie ma ani jednej wersji")
		}
		return *stan.dokument.WersjaBiezacaKod, nil
	case "liczba znaków", "liczbaznakow":
		return strconv.Itoa(len([]rune(tresc))), nil
	case "liczba słów", "liczbaslow":
		return strconv.Itoa(len(strings.Fields(tresc))), nil
	case "liczba akapitów", "liczbaakapitow":
		return strconv.Itoa(len(strings.Split(tresc, "\n"))), nil
	case "liczba stron", "liczbastron":
		return strconv.Itoa(stron), nil
	case "liczba tabel", "liczbatabel":
		return strconv.Itoa(len(stan.forma.Tables)), nil
	case "liczba obiektów", "liczbaobiektow":
		return strconv.Itoa(len(stan.forma.Objects)), nil
	}
	return "", poleBrak(pole, "właściwość dokumentu nieznana rdzeniowi",
		"właściwości „"+wartoscTekstu(pole.PropertyName)+"” rdzeń nie zna; wykaz: "+
			poleWykazWlasciwosci())
}

// poleWykazWlasciwosci wymienia właściwości dokumentu pełnymi nazwami —
// odmowa ma powiedzieć Operatorowi, o co wolno pytać.
func poleWykazWlasciwosci() string {
	return "tytuł, wskazanie, okno, format, utworzono, zaktualizowano, wersja, " +
		"liczba znaków, liczba słów, liczba akapitów, liczba stron, liczba tabel, " +
		"liczba obiektów"
}

// poleBrak składa pominięcie bilansu dla pola, którego nie da się policzyć, wraz z powodem i wskazaniem miejsca pola w dokumencie.
func poleBrak(pole shared.StudioDocumentField, powod, szczegol string) *shared.StudioSkippedItem {
	pominiecie := shared.StudioSkippedItem{
		Reason: powod,
		Detail: postacWskaznikTekstu(poleNazwaRodzaju(pole.Kind) + " „" + pole.Id + "”: " +
			szczegol),
	}
	if pole.AnchorOffset != nil {
		pominiecie.RangeStart = pole.AnchorOffset
		pominiecie.RangeEnd = pole.AnchorOffset
	}
	return &pominiecie
}

// ── Zapis daty i godziny ────────────────────────────────────────────────────

// poleZapisChwili przekłada wzór Operatora na zapis chwili.
//
// Wzór jest pisany znakami, które Operator zna z pakietu biurowego (dd, MM,
// yyyy, HH, mm, ss), a nie układem odniesienia biblioteki Go: nikt nie wpisze
// „2006-01-02" jako wzoru daty.
func poleZapisChwili(chwila time.Time, wzor, domyslny string) string {
	zapis := strings.TrimSpace(wzor)
	if zapis == "" {
		zapis = domyslny
	}
	zamiany := []struct{ z, na string }{
		{"yyyy", "2006"}, {"yy", "06"},
		{"MMMM", "January"}, {"MMM", "Jan"}, {"MM", "01"},
		{"dddd", "Monday"}, {"ddd", "Mon"}, {"dd", "02"},
		{"HH", "15"}, {"mm", "04"}, {"ss", "05"},
	}
	uklad := zapis
	for _, zamiana := range zamiany {
		uklad = strings.ReplaceAll(uklad, zamiana.z, zamiana.na)
	}
	return chwila.Format(uklad)
}

// ── Pole obliczane ──────────────────────────────────────────────────────────

// poleObliczWyrazenie liczy wyrażenie pola obliczanego i zwraca liczbę albo błąd wskazujący miejsce, w którym rachunek się zatrzymał.
func poleObliczWyrazenie(wyrazenie string, forma *shared.StudioDocumentForm) (float64, error) {
	rachunek := &poleRachunek{znaki: []rune(wyrazenie), forma: forma}
	wynik, err := rachunek.suma()
	if err != nil {
		return 0, err
	}
	rachunek.omijOdstepy()
	if rachunek.miejsce < len(rachunek.znaki) {
		return 0, bladWskazaniaStudio("wyrażenie pola obliczanego: rachunek stanął na " +
			"znaku " + strconv.Itoa(rachunek.miejsce+1) + " („" +
			string(rachunek.znaki[rachunek.miejsce]) + "”); rachunek zna liczby, działania " +
			"+ - * /, nawiasy oraz działania na kolumnie tabeli: SUMA, ŚREDNIA, MIN, MAKS, " +
			"LICZBA")
	}
	return wynik, nil
}

// poleRachunek to rachunek wyrażenia pola obliczanego, czytany znak po znaku od lewej strony wyrażenia.
type poleRachunek struct {
	znaki   []rune
	miejsce int
	forma   *shared.StudioDocumentForm
}

func (r *poleRachunek) omijOdstepy() {
	for r.miejsce < len(r.znaki) && unicode.IsSpace(r.znaki[r.miejsce]) {
		r.miejsce++
	}
}

// suma liczy dodawanie i odejmowanie w wyrażeniu pola obliczanego, wywołując iloczyn dla składników silniej wiążących.
func (r *poleRachunek) suma() (float64, error) {
	wynik, err := r.iloczyn()
	if err != nil {
		return 0, err
	}
	for {
		r.omijOdstepy()
		if r.miejsce >= len(r.znaki) {
			return wynik, nil
		}
		znak := r.znaki[r.miejsce]
		if znak != '+' && znak != '-' {
			return wynik, nil
		}
		r.miejsce++
		prawy, err := r.iloczyn()
		if err != nil {
			return 0, err
		}
		if znak == '+' {
			wynik += prawy
		} else {
			wynik -= prawy
		}
	}
}

// iloczyn liczy mnożenie i dzielenie w wyrażeniu pola obliczanego, wywołując składnik dla pojedynczych wartości.
func (r *poleRachunek) iloczyn() (float64, error) {
	wynik, err := r.skladnik()
	if err != nil {
		return 0, err
	}
	for {
		r.omijOdstepy()
		if r.miejsce >= len(r.znaki) {
			return wynik, nil
		}
		znak := r.znaki[r.miejsce]
		if znak != '*' && znak != '/' {
			return wynik, nil
		}
		r.miejsce++
		prawy, err := r.skladnik()
		if err != nil {
			return 0, err
		}
		if znak == '/' {
			if prawy == 0 {
				return 0, bladWskazaniaStudio("wyrażenie pola obliczanego dzieli przez zero")
			}
			wynik /= prawy
			continue
		}
		wynik *= prawy
	}
}

// skladnik liczy pojedynczy składnik wyrażenia: liczbę, nawias, znak jednoargumentowy albo działanie na kolumnie tabeli.
func (r *poleRachunek) skladnik() (float64, error) {
	r.omijOdstepy()
	if r.miejsce >= len(r.znaki) {
		return 0, bladWskazaniaStudio("wyrażenie pola obliczanego urywa się przed końcem")
	}
	switch znak := r.znaki[r.miejsce]; {
	case znak == '-':
		r.miejsce++
		wartosc, err := r.skladnik()
		return -wartosc, err
	case znak == '+':
		r.miejsce++
		return r.skladnik()
	case znak == '(':
		r.miejsce++
		wartosc, err := r.suma()
		if err != nil {
			return 0, err
		}
		r.omijOdstepy()
		if r.miejsce >= len(r.znaki) || r.znaki[r.miejsce] != ')' {
			return 0, bladWskazaniaStudio("wyrażenie pola obliczanego ma nawias otwarty " +
				"i niezamknięty")
		}
		r.miejsce++
		return wartosc, nil
	case unicode.IsDigit(znak) || znak == '.' || znak == ',':
		return r.liczba()
	case unicode.IsLetter(znak):
		return r.dzialanieNaTabeli()
	}
	return 0, bladWskazaniaStudio("wyrażenie pola obliczanego: znaku „" +
		string(r.znaki[r.miejsce]) + "” rachunek nie zna")
}

// liczba czyta zapis liczby z wyrażenia, przyjmując przecinek dziesiętny właściwy pismu polskiemu obok kropki.
func (r *poleRachunek) liczba() (float64, error) {
	poczatek := r.miejsce
	for r.miejsce < len(r.znaki) {
		znak := r.znaki[r.miejsce]
		if unicode.IsDigit(znak) || znak == '.' || znak == ',' {
			r.miejsce++
			continue
		}
		break
	}
	zapis := strings.ReplaceAll(string(r.znaki[poczatek:r.miejsce]), ",", ".")
	wartosc, err := strconv.ParseFloat(zapis, 64)
	if err != nil {
		return 0, bladWskazaniaStudio("wyrażenie pola obliczanego: „" + zapis +
			"” nie jest liczbą")
	}
	return wartosc, nil
}

// dzialanieNaTabeli liczy działanie na kolumnie tabeli, na przykład SUMA(tabela; kolumna), pomijając wiersz nagłówkowy.
func (r *poleRachunek) dzialanieNaTabeli() (float64, error) {
	poczatek := r.miejsce
	for r.miejsce < len(r.znaki) && (unicode.IsLetter(r.znaki[r.miejsce]) ||
		r.znaki[r.miejsce] == 'Ś' || r.znaki[r.miejsce] == 'ś') {
		r.miejsce++
	}
	nazwa := strings.ToUpper(string(r.znaki[poczatek:r.miejsce]))
	r.omijOdstepy()
	if r.miejsce >= len(r.znaki) || r.znaki[r.miejsce] != '(' {
		return 0, bladWskazaniaStudio("wyrażenie pola obliczanego: po nazwie działania „" +
			nazwa + "” rachunek oczekuje nawiasu ze wskazaniem tabeli i kolumny, " +
			"na przykład SUMA(studio-tab-1; 2)")
	}
	r.miejsce++
	poczatekArgumentow := r.miejsce
	for r.miejsce < len(r.znaki) && r.znaki[r.miejsce] != ')' {
		r.miejsce++
	}
	if r.miejsce >= len(r.znaki) {
		return 0, bladWskazaniaStudio("wyrażenie pola obliczanego: działanie „" + nazwa +
			"” ma nawias otwarty i niezamknięty")
	}
	argumenty := string(r.znaki[poczatekArgumentow:r.miejsce])
	r.miejsce++

	czesci := strings.FieldsFunc(argumenty, func(znak rune) bool {
		return znak == ';' || znak == ','
	})
	if len(czesci) != 2 {
		return 0, bladWskazaniaStudio("wyrażenie pola obliczanego: działanie „" + nazwa +
			"” wymaga dwóch wskazań — tabeli i kolumny liczonej od zera, na przykład " +
			"SUMA(studio-tab-1; 2)")
	}
	kodTabeli := strings.TrimSpace(czesci[0])
	kolumna, err := strconv.Atoi(strings.TrimSpace(czesci[1]))
	if err != nil {
		return 0, bladWskazaniaStudio("wyrażenie pola obliczanego: „" +
			strings.TrimSpace(czesci[1]) + "” nie jest numerem kolumny")
	}
	if r.forma == nil {
		return 0, postacBladZaplecza("rachunek pola obliczanego bez postaci dokumentu")
	}
	tabela, err := tabelaZnajdz(r.forma, kodTabeli)
	if err != nil {
		return 0, err
	}
	kopia := *tabela
	tabelaSiatkaPelna(&kopia)
	if kolumna < 0 || kolumna >= kopia.Columns {
		return 0, bladWskazaniaStudio("wyrażenie pola obliczanego: kolumna " +
			strconv.Itoa(kolumna) + ", a tabela ma kolumn " + strconv.Itoa(kopia.Columns))
	}
	pierwszy := 0
	if kopia.HeaderRows != nil && *kopia.HeaderRows > 0 {
		// Wiersz nagłówkowy nie wchodzi do rachunku: sumowanie nagłówka dałoby wynik trudny do wyjaśnienia.
		pierwszy = *kopia.HeaderRows
	}
	liczby := make([]float64, 0, kopia.Rows)
	for wiersz := pierwszy; wiersz < kopia.Rows; wiersz++ {
		komorka := tabelaKomorka(&kopia, wiersz, kolumna)
		if komorka == nil {
			continue
		}
		if liczba, jest := tabelaLiczbaZTekstu(wartoscTekstu(komorka.Text)); jest {
			liczby = append(liczby, liczba)
		}
	}
	return poleDzialanie(nazwa, liczby)
}

// poleDzialanie liczy wskazane działanie — sumę, średnią, minimum, maksimum albo liczbę wartości — na zebranych liczbach kolumny.
func poleDzialanie(nazwa string, liczby []float64) (float64, error) {
	if nazwa == "LICZBA" {
		return float64(len(liczby)), nil
	}
	if len(liczby) == 0 {
		return 0, bladWskazaniaStudio("wyrażenie pola obliczanego: działanie „" + nazwa +
			"” nie ma na czym stanąć — wskazana kolumna nie niesie ani jednej liczby")
	}
	switch nazwa {
	case "SUMA":
		suma := 0.0
		for _, liczba := range liczby {
			suma += liczba
		}
		return suma, nil
	case "ŚREDNIA", "SREDNIA":
		suma := 0.0
		for _, liczba := range liczby {
			suma += liczba
		}
		return suma / float64(len(liczby)), nil
	case "MIN":
		najmniejsza := liczby[0]
		for _, liczba := range liczby {
			if liczba < najmniejsza {
				najmniejsza = liczba
			}
		}
		return najmniejsza, nil
	case "MAKS", "MAX":
		najwieksza := liczby[0]
		for _, liczba := range liczby {
			if liczba > najwieksza {
				najwieksza = liczba
			}
		}
		return najwieksza, nil
	}
	return 0, bladWskazaniaStudio("wyrażenie pola obliczanego: działania „" + nazwa +
		"” rachunek nie zna; wykaz: SUMA, ŚREDNIA, MIN, MAKS, LICZBA")
}

// poleZapisLiczby zapisuje wynik pola obliczanego tekstem, ucinając zbędne zera po przecinku dziesiętnym.
func poleZapisLiczby(wartosc float64) string {
	zapis := strconv.FormatFloat(wartosc, 'f', -1, 64)
	return zapis
}

// ── Drobne rachunki ─────────────────────────────────────────────────────────

// poleZnajdz odnajduje pole o wskazanym kodzie w postaci dokumentu i mówi, czy takie pole w niej jest.
func poleZnajdz(forma *shared.StudioDocumentForm, kod string) (shared.StudioDocumentField, bool) {
	for _, pole := range forma.Fields {
		if pole.Id == kod {
			return pole, true
		}
	}
	return shared.StudioDocumentField{}, false
}

// poleSprawdzRodzaj odrzuca rodzaj pola, którego kontrakt nie zna, wskazując w odmowie wykaz rodzajów znanych.
func poleSprawdzRodzaj(rodzaj shared.StudioFieldKind) error {
	for _, znany := range shared.WartosciStudioFieldKind() {
		if rodzaj == znany {
			return nil
		}
	}
	nazwy := make([]string, 0, 9)
	for _, znany := range shared.WartosciStudioFieldKind() {
		nazwy = append(nazwy, string(znany))
	}
	return bladWskazaniaStudio("pole rodzaju „" + string(rodzaj) +
		"”, którego rdzeń nie zna; wykaz: " + strings.Join(nazwy, ", "))
}

// poleSprawdzWymagania pilnuje, żeby pole wchodziło z danymi, bez których rdzeń nie ma z czego policzyć jego wartości.
func poleSprawdzWymagania(z shared.StudioFieldInsertRequest) error {
	switch z.Kind {
	case shared.StudioFieldKindCalculated:
		if z.Expression == nil || strings.TrimSpace(*z.Expression) == "" {
			return bladWskazaniaStudio("pole obliczane bez wyrażenia — pole wymaga " +
				"wyrażenia (pole expression); rachunek zna liczby, działania + - * /, " +
				"nawiasy oraz działania na kolumnie tabeli: SUMA, ŚREDNIA, MIN, MAKS, LICZBA")
		}
	case shared.StudioFieldKindDocumentProperty, shared.StudioFieldKindTemplateField:
		if z.PropertyName == nil || strings.TrimSpace(*z.PropertyName) == "" {
			return bladWskazaniaStudio("pole właściwości bez jej nazwy — pole wymaga " +
				"nazwy (pole propertyName); wykaz właściwości dokumentu: " +
				poleWykazWlasciwosci())
		}
	}
	return nil
}

// poleNazwaRodzaju nazywa rodzaj pola pełnym słowem, zamiast kodu kontraktu, na potrzeby noty bilansu i odmowy.
func poleNazwaRodzaju(rodzaj shared.StudioFieldKind) string {
	switch rodzaj {
	case shared.StudioFieldKindPageNumber:
		return "pole numeru strony"
	case shared.StudioFieldKindPageCount:
		return "pole liczby stron"
	case shared.StudioFieldKindDate:
		return "pole daty"
	case shared.StudioFieldKindTime:
		return "pole godziny"
	case shared.StudioFieldKindDocumentTitle:
		return "pole tytułu dokumentu"
	case shared.StudioFieldKindDocumentAuthor:
		return "pole autora dokumentu"
	case shared.StudioFieldKindDocumentProperty:
		return "pole właściwości dokumentu"
	case shared.StudioFieldKindCalculated:
		return "pole obliczane"
	case shared.StudioFieldKindTemplateField:
		return "pole szablonu"
	default:
		return "pole dokumentu"
	}
}

// poleZapisWartosci opisuje wartość pola w nocie bilansu, nazywając brak wartości wprost, zamiast pomijać go milczeniem.
func poleZapisWartosci(wartosc *string) string {
	if wartosc == nil || *wartosc == "" {
		return "brak — pole wymaga odświeżenia albo uzupełnienia wskazania"
	}
	return "„" + *wartosc + "”"
}

// poleDoWiersza przekłada pole postaci dokumentu na wiersz warstwy danych, gotowy do zapisu w składnicy.
func poleDoWiersza(dokumentID int64, pole shared.StudioDocumentField) dane.PoleDokumentuStudia {
	nieswieze := pole.Stale != nil && *pole.Stale
	return dane.PoleDokumentuStudia{
		Kod: pole.Id, DokumentID: dokumentID, Rodzaj: string(pole.Kind),
		Kotwica: int64(aparatWartoscLiczby(pole.AnchorOffset)),
		Format:  pole.Format, Wyrazenie: pole.Expression,
		NazwaWlasciwosci: pole.PropertyName, Wartosc: pole.Value, Nieswieze: nieswieze,
	}
}
