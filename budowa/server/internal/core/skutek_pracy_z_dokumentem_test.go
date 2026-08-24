package core

import (
	"encoding/json"
	"strings"
	"testing"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// Skutek okna pracy z dokumentem: czy dwie zmiany rdzenia NAPRAWDĘ coś robią.
//
// Szkody, które ten plik ma wykluczyć:
//   1. przyjęcie wskazanych fragmentów propozycji, które podmienia treść CAŁĄ —
//      Operator wybiera fragment drugi, a dostaje wszystko;
//   2. zakres zmiany śledzonej liczony w bajtach — dokument polski ma litery
//      dwubajtowe i decyzja rozcinałaby je w środku;
//   3. nastawy suwaków i polecenie Operatora, które nie dojeżdżają do modelu —
//      suwak przestawiałby wtedy pole bez skutku.
//
// Operacja kontekstowa jest tu mierzona OD KOŃCA DO KOŃCA, jednym przebiegiem
// (`TestOperacjaKontekstowaOdKoncaDoKonca`). Do 17.08.2026 stało w tym miejscu
// zdanie, że zmierzyć jej nie sposób, bo wymaga kanału modelu, którego uprząż
// nie stawia — i było nieprawdą: uprząż niesie kanał `echo`, adapter bez sieci
// wkompilowany w rdzeń, tym samym wpięciem, którym mierzy się prowenancja
// i zajętość okna kontekstu. Droga „operacja → zmiana śledzona w dokumencie" nie
// jest więc mierzona po częściach; części zostają obok jako sprawdziany rachunku.

// TestPrzyjecieWskazanychFragmentowZmieniaTylkoJe jest sednem pola hunkIndexes:
// fragment wskazany bierze się ze strony propozycji, a fragment pominięty
// zostaje taki, jaki stoi w dokumencie.
func TestPrzyjecieWskazanychFragmentowZmieniaTylkoJe(t *testing.T) {
	bazowa := strings.Join([]string{
		"Wstęp bez zmian.",
		"Zdanie pierwsze do poprawy.",
		"Zakończenie bez zmian.",
	}, "\n")
	docelowa := strings.Join([]string{
		"Wstęp bez zmian.",
		"Zdanie pierwsze poprawione przez model.",
		"Zakończenie bez zmian.",
	}, "\n")

	// Rachunek fragmentów jest ten sam, który widzi Operator w oknie — z niego
	// bierze się numer, który potem wskazuje.
	fragmenty := policzFragmentyRoznicy(bazowa, docelowa)
	var numerZmiany int
	for _, fragment := range fragmenty {
		if fragment.Kind != shared.DiffHunkKindContext {
			numerZmiany = fragment.Index
		}
	}
	if numerZmiany == 0 {
		t.Fatalf("różnica dwóch treści nie ma ani jednego fragmentu zmiany: %+v", fragmenty)
	}

	zlozona, err := zlozTrescZFragmentow(bazowa, docelowa, []int{numerZmiany})
	if err != nil {
		t.Fatalf("złożenie treści z fragmentu %d odmówiło: %v", numerZmiany, err)
	}
	if zlozona != docelowa {
		t.Errorf("przyjęcie jedynego fragmentu zmiany nie dało treści docelowej:\n%q", zlozona)
	}

	// Wykaz bez tego fragmentu ma zostawić treść bazową w całości — to jest
	// właściwa miara „zmieniło się tylko wskazane".
	bezZmiany, err := zlozTrescZFragmentow(bazowa, docelowa, []int{fragmenty[0].Index})
	if err != nil {
		t.Fatalf("złożenie treści z fragmentu kontekstowego odmówiło: %v", err)
	}
	if bezZmiany != bazowa {
		t.Errorf("wskazanie fragmentu kontekstowego zmieniło treść:\n%q", bezZmiany)
	}
}

// TestFragmentSpozaRachunkuOdmawia pilnuje granicy: przyjęcie fragmentu, którego
// w różnicy nie ma, nie może wrócić jako wykonane.
func TestFragmentSpozaRachunkuOdmawia(t *testing.T) {
	if _, err := zlozTrescZFragmentow("jedno zdanie", "inne zdanie", []int{7}); err == nil {
		t.Error("wskazanie fragmentu spoza różnicy wróciło bez odmowy")
	}
}

// TestZakresZmianySledzonejLiczySieWZnakach mierzy szkodę po polskich literach:
// zakres liczony w bajtach wstawiłby treść w środek znaku dwubajtowego.
func TestZakresZmianySledzonejLiczySieWZnakach(t *testing.T) {
	tresc := "Zażółć gęślą jaźń."
	przed := "gęślą"
	po := "gęślą jaźń"
	// „gęślą" zaczyna się na siódmym ZNAKU treści; w bajtach byłby to znak
	// dziewiąty, bo „ż" i „ó" są dwubajtowe.
	od := int64(strings.Count(string([]rune(tresc)[:7]), "") - 1)

	zmiana := dane.ZmianaSledzona{
		Kod: "studio-zm-1", Rodzaj: string(shared.StudioChangeKindWstawienie),
		Autor: string(shared.StudioAuthorModel), ZakresOd: od, ZakresDo: od + int64(len([]rune(przed))),
		TrescPrzed: &przed, TrescPo: &po,
	}

	nowa, zmieniono := zastosujZmianeSledzona(tresc, zmiana, true)
	if !zmieniono {
		t.Fatal("przyjęcie zmiany o zakresie w granicach treści nie zmieniło treści")
	}
	if nowa != "Zażółć gęślą jaźń jaźń." {
		t.Errorf("przyjęcie wstawiło treść w niewłaściwe miejsce: %q", nowa)
	}

	wrocona, _ := zastosujZmianeSledzona(nowa, dane.ZmianaSledzona{
		Kod: "studio-zm-1", Rodzaj: string(shared.StudioChangeKindWstawienie),
		Autor: string(shared.StudioAuthorModel), ZakresOd: od, ZakresDo: od + int64(len([]rune(po))),
		TrescPrzed: &przed, TrescPo: &po,
	}, false)
	if wrocona != tresc {
		t.Errorf("odrzucenie nie przywróciło treści sprzed zmiany: %q", wrocona)
	}
}

// TestZakresOperacjiSchodziNaCalyDokument sprawdza, że zakres bez zaznaczenia
// i zakres spoza treści dają cały dokument, a nie odmowę po wykonanej pracy.
func TestZakresOperacjiSchodziNaCalyDokument(t *testing.T) {
	dlugosc := 40

	od, do_ := zakresOperacjiStudia(shared.StudioContextualOpRequest{
		Scope: shared.StudioOperationScopeDocument,
	}, dlugosc)
	if od != 0 || do_ != dlugosc {
		t.Errorf("zakres dokumentu wyszedł %d–%d, oczekiwano 0–%d", od, do_, dlugosc)
	}

	od, do_ = zakresOperacjiStudia(shared.StudioContextualOpRequest{
		Scope: shared.StudioOperationScopeSelection, SelectionStart: wskaznikZakresuPracy(5),
		SelectionEnd: wskaznikZakresuPracy(12),
	}, dlugosc)
	if od != 5 || do_ != 12 {
		t.Errorf("zakres zaznaczenia wyszedł %d–%d, oczekiwano 5–12", od, do_)
	}

	od, do_ = zakresOperacjiStudia(shared.StudioContextualOpRequest{
		Scope: shared.StudioOperationScopeSelection, SelectionStart: wskaznikZakresuPracy(5),
		SelectionEnd: wskaznikZakresuPracy(500),
	}, dlugosc)
	if od != 0 || do_ != dlugosc {
		t.Errorf("zakres spoza treści wyszedł %d–%d, oczekiwano zejścia na całość", od, do_)
	}
}

// TestPolecenieOperatoraDojezdzaDoModelu pilnuje wiersza polecenia i suwaków:
// nastawy jadą polem `params`, więc treść polecenia dla modelu musi je nieść.
func TestPolecenieOperatoraDojezdzaDoModelu(t *testing.T) {
	tresc := "Strony ustalają zakres współpracy."
	ladunek, err := json.Marshal(map[string]any{
		"polecenie": "skróć to zdanie o połowę",
		"objetosc":  -2,
	})
	if err != nil {
		t.Fatalf("nie można złożyć ładunku sprawdzianu: %v", err)
	}

	polecenie := trescOperacjiStudia(shared.StudioContextualOpRequest{
		ActionId: "studio.styl.skrocenie", Scope: shared.StudioOperationScopeDocument,
		Params: ladunek,
	}, dane.DokumentStudia{Tresc: &tresc})

	if !strings.Contains(polecenie, "skróć to zdanie o połowę") {
		t.Error("polecenie Operatora nie weszło do treści wysyłanej modelowi")
	}
	if !strings.Contains(polecenie, "objetosc") {
		t.Error("nastawa suwaka nie weszła do treści wysyłanej modelowi")
	}
	if !strings.Contains(polecenie, tresc) {
		t.Error("treść dokumentu wypadła z polecenia po dołożeniu nastaw")
	}
}

// wskaznikZakresuPracy oddaje wskaźnik na liczbę — pola zakresu kontraktu są
// nieobowiązkowe, a sprawdzian potrzebuje ich wypełnionych. Przedrostek obszaru
// jest wymagany: przestrzeń nazw pakietu `core` jest dzielona.
func wskaznikZakresuPracy(wartosc int) *int {
	return &wartosc
}

// TestOperacjaKontekstowaOdKoncaDoKonca mierzy CAŁĄ drogę operacji: żądanie
// okna, wywołanie kanału modelu, wpisanie wyniku w treść dokumentu, zmianę
// śledzoną i wersję.
//
// ── Dlaczego jednym przebiegiem, a nie po częściach ─────────────────────────
// Części tej drogi były mierzone osobno: złożenie polecenia dla modelu i rachunek
// zakresu. Obie mogą być poprawne, a droga nadal zerwana — wynik modelu może nie
// dojść do treści, zmiana śledzona może się nie odłożyć, wersja może nie powstać.
// Kontrakt obiecuje, że po operacji kontekstowej w dokumencie STOI zmiana
// oznaczona autorstwem modelu, i to jest twierdzenie o skutku, nie o rachunku.
//
// ── Czym jest tu kanał próbny ───────────────────────────────────────────────
// Kanał `echo` odsyła treść zapytania i nie sięga do sieci ani do żadnego
// programu — jest wkompilowany w rdzeń. Wynik operacji jest więc znany z góry:
// jest nim polecenie złożone przez `trescOperacjiStudia`, a w nim wiersz
// „Czynnosc: <pozycja rejestru>". Sprawdzian szuka właśnie jego, bo to jedyny
// znak, który mógł przyjść WYŁĄCZNIE od modelu — treści dokumentu nie było
// w nim ani jednego takiego wiersza.
func TestOperacjaKontekstowaOdKoncaDoKonca(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)
	kanal := kanalEchoSprawdzianu(t, zmontowany, zycie, 8192)
	sesja := zalozSesjeSprawdzianu(t, zmontowany, zycie)
	okno := zalozOknoSprawdzianu(t, zmontowany, zycie, sesja, kanal)

	// Dokument zakłada się w TYM oknie, bo operacja kontekstowa czyta kanał
	// modelu z okna, w którego imieniu przyszło żądanie.
	var otwarcie shared.StudioDocumentOpenResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioDocumentOpen,
		shared.StudioDocumentOpenRequest{WindowId: okno}, &otwarcie)

	zastana := "Strony ustalają zakres współpracy.\nTermin: koniec miesiąca."
	var zapis shared.StudioDocumentSaveResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioDocumentSave,
		shared.StudioDocumentSaveRequest{
			DocumentId: otwarcie.Document.Id,
			Content:    zastana,
			Title:      wskaznik("umowa sprawdzianu"),
		}, &zapis)
	kodDokumentu := zapis.Document.Id

	const pozycjaRejestru = "studio.styl.uproszczenie"
	var operacja shared.StudioContextualOpResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioContextualOp,
		shared.StudioContextualOpRequest{
			WindowId:       okno,
			DocumentId:     kodDokumentu,
			ActionId:       pozycjaRejestru,
			Scope:          shared.StudioOperationScopeSelection,
			SelectionStart: wskaznikZakresuPracy(0),
			SelectionEnd:   wskaznikZakresuPracy(len([]rune("Strony ustalają zakres współpracy."))),
		}, &operacja)

	if operacja.ProposalId == nil || *operacja.ProposalId == "" {
		t.Fatal("operacja kontekstowa nie oddała propozycji — a wynik modelu przyszedł")
	}
	if operacja.ResultText == nil || !strings.Contains(*operacja.ResultText, pozycjaRejestru) {
		t.Fatalf("wynik operacji nie niesie znaku wywołania modelu: %+v", operacja.ResultText)
	}

	// ── Skutek pierwszy: wynik STOI w treści dokumentu ──────────────────────
	var poOperacji shared.StudioDocumentOpenResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioDocumentOpen,
		shared.StudioDocumentOpenRequest{
			WindowId: okno, DocumentId: wskaznik(kodDokumentu),
		}, &poOperacji)

	tresc := ""
	if poOperacji.Document.Content != nil {
		tresc = *poOperacji.Document.Content
	}
	if !strings.Contains(tresc, "Czynnosc: "+pozycjaRejestru) {
		t.Fatalf("treść dokumentu nie niesie wyniku operacji — droga „operacja → "+
			"dokument"+"\" jest zerwana. Treść: %q", tresc)
	}
	if !strings.Contains(tresc, "Termin: koniec miesiąca.") {
		t.Error("operacja na ZAZNACZENIU nadpisała treść poza zaznaczeniem")
	}

	// ── Skutek drugi: zmiana śledzona autorstwa modelu ──────────────────────
	var sledzone shared.StudioTrackingListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioTrackingList,
		shared.StudioTrackingListRequest{DocumentId: kodDokumentu}, &sledzone)

	if len(sledzone.Changes) == 0 {
		t.Fatal("po operacji kontekstowej dokument nie ma ani jednej zmiany śledzonej — " +
			"a kontrakt obiecuje zmianę oznaczoną autorstwem modelu")
	}
	odModelu := 0
	for _, zmiana := range sledzone.Changes {
		if zmiana.Author != shared.StudioAuthorModel {
			continue
		}
		odModelu++
		// Zakres liczy się w ZNAKACH i musi wskazywać miejsce wyniku w treści
		// NOWEJ — inaczej decyzja o zmianie odtworzyłaby nie ten fragment.
		if zmiana.RangeEnd <= zmiana.RangeStart {
			t.Errorf("zmiana %s ma zakres pusty albo odwrócony: %d–%d",
				zmiana.Id, zmiana.RangeStart, zmiana.RangeEnd)
		}
		if zmiana.RangeEnd > len([]rune(tresc)) {
			t.Errorf("zmiana %s wskazuje znak %d, a treść ma znaków %d",
				zmiana.Id, zmiana.RangeEnd, len([]rune(tresc)))
		}
	}
	if odModelu == 0 {
		t.Error("żadna zmiana śledzona nie jest oznaczona autorstwem modelu")
	}

	// ── Skutek trzeci: wersja z odwołaniem do propozycji ────────────────────
	var wersje shared.StudioRepositoryListResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioRepositoryList,
		shared.StudioRepositoryListRequest{DocumentId: kodDokumentu}, &wersje)

	wersjaModelu := false
	for _, wersja := range wersje.Versions {
		if wersja.Author != nil && *wersja.Author == shared.StudioAuthorModel {
			wersjaModelu = true
		}
	}
	if !wersjaModelu {
		t.Error("operacja kontekstowa nie założyła wersji autorstwa modelu — bez niej " +
			"nie ma do czego wrócić po odrzuceniu wyniku")
	}
}

// TestPorownanieBezStronWracaOdmowa pilnuje granicy, na której `studio.diff.compare`
// meldował powodzenie kopertą pustą.
//
// Koperta pusta ze stanem `ok` mówi oknu „porównałem i nie ma czego pokazać",
// a rdzeń nie porównał niczego: fragmenty różnicy potrzebują dwóch stron,
// a wzorzec potrzebuje strony, po której ma szukać. Odmowa nazywająca brakujące
// pole jest tu jedyną odpowiedzią prawdziwą — po pustej kopercie okno nie ma jak
// odróżnić „wersje są zgodne" od „nie podałeś, co z czym porównać".
func TestPorownanieBezStronWracaOdmowa(t *testing.T) {
	zmontowany, zycie, _ := zmontujDoPomiaruSkutku(t)
	sesja := zalozSesjeSprawdzianu(t, zmontowany, zycie)
	okno := zalozOknoSprawdzianu(t, zmontowany, zycie, sesja, "kanal-sprawdzianu")

	var otwarcie shared.StudioDocumentOpenResponse
	wykonajUdana(t, zmontowany, zycie, shared.CommandStudioDocumentOpen,
		shared.StudioDocumentOpenRequest{WindowId: okno}, &otwarcie)
	dokument := otwarcie.Document.Id

	zalozWersje := func(tresc string) string {
		t.Helper()
		var zapis shared.StudioDocumentSaveResponse
		wykonajUdana(t, zmontowany, zycie, shared.CommandStudioDocumentSave,
			shared.StudioDocumentSaveRequest{
				DocumentId:    dokument,
				Content:       tresc,
				CreateVersion: wskaznik(true),
			}, &zapis)
		if zapis.Version == nil {
			t.Fatalf("zapis z createVersion nie oddał wersji: %+v", zapis)
		}
		return zapis.Version.Id
	}
	pierwsza := zalozWersje("ala ma psa i kota")
	druga := zalozWersje("ala ma psa oraz kota")

	t.Run("bez wskazania stron", func(t *testing.T) {
		blad := wykonajOdmowna(t, zmontowany, zycie, shared.CommandStudioDiffCompare,
			shared.StudioDiffCompareRequest{DocumentId: dokument})
		odmowaNazywa(t, blad, shared.ErrorCodeValidationFailed,
			"baseVersionId", "targetVersionId")
	})

	t.Run("sam wzorzec bez strony przeszukiwanej", func(t *testing.T) {
		blad := wykonajOdmowna(t, zmontowany, zycie, shared.CommandStudioDiffCompare,
			shared.StudioDiffCompareRequest{DocumentId: dokument, Pattern: wskaznik("ma")})
		odmowaNazywa(t, blad, shared.ErrorCodeValidationFailed, "wzorzec")
	})

	t.Run("jedna strona bez drugiej", func(t *testing.T) {
		blad := wykonajOdmowna(t, zmontowany, zycie, shared.CommandStudioDiffCompare,
			shared.StudioDiffCompareRequest{DocumentId: dokument, BaseVersionId: &pierwsza})
		odmowaNazywa(t, blad, shared.ErrorCodeValidationFailed, "targetVersionId")
	})

	// Odmowa nie ma prawa objąć żądania, z którego da się coś policzyć —
	// dwie wersje mają dać różnicę, a wzorzec przy nich trafienia.
	t.Run("dwie strony dają różnicę", func(t *testing.T) {
		var porownanie shared.StudioDiffCompareResponse
		wykonajUdana(t, zmontowany, zycie, shared.CommandStudioDiffCompare,
			shared.StudioDiffCompareRequest{
				DocumentId:      dokument,
				BaseVersionId:   &pierwsza,
				TargetVersionId: &druga,
				Pattern:         wskaznik("ma"),
			}, &porownanie)
		if len(porownanie.Hunks) == 0 {
			t.Error("porównanie dwóch różnych wersji nie oddało ani jednego fragmentu")
		}
		if len(porownanie.Matches) == 0 {
			t.Error("wzorzec obecny w treści strony nie dał ani jednego trafienia")
		}
	})
}
