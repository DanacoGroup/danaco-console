// Odpowiedzialność pliku: wypełnienie portu Wiedza dwiema komendami tekstowymi
// rodziny `knowledge.*` — przełożenie żądań kontraktu na zlecenia pakietu
// `wiedza` i, co ważniejsze, przełożenie jego typowanych odmów na kody
// kontraktu. Oś obrazu stoi obok, w `adapter_modul_wiedza_obraz.go`.
//
// Wyszukiwanie po słowach działa w module Library (indeks FTS5,
// `migracja_111_indeks_tresci_biblioteki.sql`). Ta rodzina wnosi wyszukiwanie
// po znaczeniu, a różnica jest cała w tym, czego tamto nie umie: pytanie „co
// robić, gdy padła maszyna" nie ma z dokumentem o awarii węzła ani jednego
// wspólnego słowa, więc indeks liter go nie znajdzie.
//
// Pomocnik osadzeń startuje tym samym uruchamiaczem i przez tę samą bramę
// izolacji, co każdy inny proces drzewa, więc potrzebuje trójki
// `session.Okno` + `session.Zasady` + `session.Obszar`. Żądania `knowledge.*`
// niosą co najwyżej `windowId`, i to po to, żeby wskazać przestrzeń do
// przeszukania — nie po to, żeby w tym oknie liczyć. Wskaźnik jest jeden na
// maszynę i wspólny dla wszystkich okien
// (`migracja_115_wskaznik_znaczenia.sql`), więc trójka składa się w zasięgu
// platformy, tą samą drogą co w silniku mowy:
//
//	zasady := ZasadyIzolacji(rozstrzygacz, konfig.Kontekst{})
//	obszar := ObszarOkna(katalog.Ustal(konfig.Kontekst{}, ""), "")
//
// Pusty kontekst zasięgu jest poprawnym adresem najszerszego z poziomów, a nie
// podstawieniem pustych struktur po cichu.
//
// Silnik powstaje na każde wywołanie, tak samo jak w module mowy:
// `ZUstawieniami` mutuje byt, więc jedna instancja współdzielona przez
// równoległe żądania oznaczałaby wyścig o nastawy. Nastawy są przy tym świeże —
// zmiana `wiedza_model` komendą `config.set` wchodzi w następnym przebiegu, bez
// restartu rdzenia.
package core

import (
	"context"
	"errors"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/session"
	"danacoconsole/server/internal/wiedza"
	"danacoconsole/shared"
)

// domyslnaLiczbaTrafien wchodzi, gdy żądanie nie poda `limit`. Dziesięć
// fragmentów to tyle, ile model zmieści w kontekście bez wypierania rozmowy —
// wyszukiwanie ma dać modelowi podstawę do odpowiedzi, a nie zająć całą jego uwagę.
const domyslnaLiczbaTrafien = 10

// granicaTrafien chroni przed żądaniem, które chce cały wskaźnik. Setka
// fragmentów po 700 znaków to już 70 tysięcy znaków odpowiedzi.
const granicaTrafien = 100

// adapterWiedzy wypełnia port Wiedza.
type adapterWiedzy struct {
	// uruchamiacz jest portem warstwy kanału — jedyną drogą startu procesu
	// w drzewie. Bez niego moduł nie ruszy pomocnika osadzeń.
	uruchamiacz session.Uruchamiacz
	// katalogDanych — ten sam katalog, w którym leżą baza, sejf poświadczeń
	// i magazyn treści biblioteki. Tam wykłada się pomocnik i tam lądują wagi.
	katalogDanych string
	// rozstrzygacz i katalog składają nastawy oraz zasady izolacji zasięgu
	// platformy — te same dwa źródła, którymi jadą Terminal, Developer i mowa.
	rozstrzygacz *konfig.Rozstrzygacz
	katalog      *KatalogRoboczy
	// skladnica trzyma wektory przy bazie rdzenia. Zależność obowiązkowa:
	// wskaźnik bez miejsca zapisu nie jest wskaźnikiem, a odmowa nazywa to
	// wprost zamiast oddawać pustkę udającą wynik.
	skladnica *wiedza.Skladnica
	// biblioteka i historia to źródła treści czytane przez repozytoria, a nie
	// własnym SQL-em: druga droga do tych samych wierszy byłaby drugą prawdą
	// o tym, co Operator ma w bibliotece.
	biblioteka dane.RepozytoriumBiblioteki
	historia   dane.RepozytoriumHistorii
	// teraz oddaje czas w milisekundach epoki — jeden zegar na byt.
	teraz func() int64
}

// nowyAdapterWiedzy wiąże port z uruchamiaczem procesów i katalogiem danych.
func nowyAdapterWiedzy(uruchamiacz session.Uruchamiacz, katalogDanych string) *adapterWiedzy {
	return &adapterWiedzy{
		uruchamiacz:   uruchamiacz,
		katalogDanych: katalogDanych,
		teraz:         func() int64 { return time.Now().UTC().UnixMilli() },
	}
}

// skladnicaWiedzy składa trwałość wskaźnika nad bazą montażu.
//
// Osobna funkcja, a nie wyrażenie w miejscu wpięcia: montaż bez bazy jest
// stanem, który zdarza się przy rdzeniu składanym do sprawdzenia transportu,
// i `nil` przechodzi tędy bez warunku po stronie wołającego. Składnica pusta
// nie udaje wtedy, że działa — odmawia zdaniem nazywającym brak.
func skladnicaWiedzy(m Montaz) *wiedza.Skladnica {
	if m.Baza == nil {
		return nil
	}
	return wiedza.NowaSkladnica(m.Baza.DB)
}

// ZIzolacja podpina rozstrzygacz zasięgu i ustalacz katalogu roboczego.
func (a *adapterWiedzy) ZIzolacja(rozstrzygacz *konfig.Rozstrzygacz, katalog *KatalogRoboczy) *adapterWiedzy {
	a.rozstrzygacz, a.katalog = rozstrzygacz, katalog
	return a
}

// ZeSkladnica podpina trwałość wskaźnika nad bazą rdzenia.
func (a *adapterWiedzy) ZeSkladnica(s *wiedza.Skladnica) *adapterWiedzy {
	a.skladnica = s
	return a
}

// ZeZrodlami podpina repozytoria, z których czytana jest treść Operatora.
func (a *adapterWiedzy) ZeZrodlami(biblioteka dane.RepozytoriumBiblioteki,
	historia dane.RepozytoriumHistorii) *adapterWiedzy {

	a.biblioteka, a.historia = biblioteka, historia
	return a
}

// Wskaznik obsługuje `knowledge.index`.
//
// Kolejność kroków jest rozstrzygnięciem: najpierw pytanie o gotowość silnika,
// dopiero potem odczyt treści. Odwrotnie, przy brakującym silniku, rdzeń
// przeczytałby całą bibliotekę z dysku, podzielił ją na fragmenty i dopiero
// wtedy powiedział „nie ma czym liczyć". Sprawdzenie gotowości kosztuje jedno
// uruchomienie pomocnika.
func (a *adapterWiedzy) Wskaznik(ctx context.Context,
	z shared.KnowledgeIndexRequest) (shared.KnowledgeIndexResponse, error) {

	zakresy, err := zakresyZadania(z.Scope)
	if err != nil {
		return shared.KnowledgeIndexResponse{}, err
	}
	ustawienia := a.ustawienia()
	silnik := a.silnik(ustawienia)
	okno, zasady, obszar := a.zasiegPlatformy()

	if err := silnik.Gotowy(ctx, okno, zasady, obszar, wiedza.LimitBudowania); err != nil {
		return shared.KnowledgeIndexResponse{}, bladWiedzy(err)
	}
	if z.Rebuild != nil && *z.Rebuild {
		if err := a.skladnica.UsunZakres(ctx, zakresy); err != nil {
			return shared.KnowledgeIndexResponse{}, bladWiedzy(err)
		}
	}

	wniesione := 0
	for _, zakres := range zakresy {
		dokumenty, err := a.dokumenty(ctx, zakres, z.WindowId)
		if err != nil {
			return shared.KnowledgeIndexResponse{}, err
		}
		ile, err := a.wniesDokumenty(ctx, silnik, okno, zasady, obszar, ustawienia, dokumenty)
		if err != nil {
			return shared.KnowledgeIndexResponse{}, bladWiedzy(err)
		}
		wniesione += ile
	}

	wszystkie, err := a.skladnica.Policz(ctx, ustawienia.Model)
	if err != nil {
		return shared.KnowledgeIndexResponse{}, bladWiedzy(err)
	}
	model := ustawienia.Model
	return shared.KnowledgeIndexResponse{Indexed: wniesione, Total: wszystkie, Model: &model}, nil
}

// wniesDokumenty dzieli treść na fragmenty, osadza je i zapisuje.
//
// Dokument po dokumencie, a nie wszystko naraz w jednej transakcji: biblioteka
// Operatora bywa gigabajtem tekstu, a jedna transakcja na całość znaczyłaby
// komplet wektorów w pamięci rdzenia. Przerwanie w połowie przebiegu zostawia
// wskaźnik niepełny, ale spójny — każdy dokument, który wszedł, wszedł w całości
// (patrz `wiedza.Skladnica.Zapisz`), a powtórzony przebieg dokończy resztę.
func (a *adapterWiedzy) wniesDokumenty(ctx context.Context, silnik *wiedza.Silnik,
	okno session.Okno, zasady session.Zasady, obszar session.Obszar,
	ustawienia wiedza.Ustawienia, dokumenty []dokumentWiedzy) (int, error) {

	wniesione := 0
	for _, dokument := range dokumenty {
		fragmenty := wiedza.Podziel(dokument.Tresc, ustawienia.DlugoscFragmentu)
		if len(fragmenty) == 0 {
			continue
		}
		teksty := make([]string, len(fragmenty))
		for i, fragment := range fragmenty {
			teksty[i] = fragment.Tresc
		}
		wektory, err := silnik.Osadz(ctx, okno, zasady, obszar, teksty, wiedza.LimitBudowania)
		if err != nil {
			return wniesione, err
		}

		// Kasowanie przed zapisem: dokument skrócony od poprzedniego przebiegu
		// zostawiłby inaczej fragmenty treści, której już nie ma — a wracałyby
		// jako cytat z dokumentu, w którym ich nie ma.
		if err := a.skladnica.UsunZrodlo(ctx, dokument.Zakres, dokument.ZrodloKod); err != nil {
			return wniesione, err
		}
		pozycje := make([]wiedza.Pozycja, len(fragmenty))
		for i, fragment := range fragmenty {
			pozycje[i] = wiedza.Pozycja{
				Zakres:    dokument.Zakres,
				Zrodlo:    dokument.Zrodlo,
				ZrodloKod: dokument.ZrodloKod,
				Kolejnosc: fragment.Kolejnosc,
				Tresc:     fragment.Tresc,
				Model:     ustawienia.Model,
				Wektor:    wiedza.Znormalizuj(wektory[i]),
			}
		}
		ile, err := a.skladnica.Zapisz(ctx, pozycje, a.teraz())
		if err != nil {
			return wniesione, err
		}
		wniesione += ile
	}
	return wniesione, nil
}

// Szukaj obsługuje `knowledge.search`.
//
// Pytanie osadza się tym samym modelem, co dokumenty: wektor pytania z modelu
// innego niż wektory wskaźnika daje iloczyn skalarny, który jest liczbą i nie
// znaczy nic. Zawężenie odczytu po nazwie modelu (`wiedza.Skladnica.Pozycje`)
// jest jedyną obroną przed tym po zmianie ustawienia.
//
// Wynik pusty jest odpowiedzią, nie odmową: wskaźnik pusty albo wiedza bez
// związku z pytaniem znaczą „nie mam na to nic" i model ma to usłyszeć wprost.
// Inaczej brak silnika, który jest odmową — wtedy rdzeń nie wie, czy ma coś,
// czy nie ma.
//
// Pole `rerank` dokłada drugi przebieg (`zPrzesiewem` niżej). Pierwszy zostaje
// niezmieniony i wykonuje się zawsze: przesiew nie ZASTĘPUJE kosinusa, tylko
// układa na nowo tych kandydatów, których kosinus wybrał — krzyżowym koderem nie
// da się przejrzeć całego wskaźnika (uzasadnienie liczbami w `wiedza/przesiew.go`).
func (a *adapterWiedzy) Szukaj(ctx context.Context,
	z shared.KnowledgeSearchRequest) (shared.KnowledgeSearchResponse, error) {

	pytanie := strings.TrimSpace(z.Query)
	if pytanie == "" {
		return shared.KnowledgeSearchResponse{}, bladZadaniaWiedzy(
			"komenda bez pytania — wyszukiwanie po znaczeniu nie ma czego porównać " +
				"z wiedzą Operatora; naprawa: podać treść pytania w polu `query`")
	}
	zakresy, err := zakresyZadania(z.Scope)
	if err != nil {
		return shared.KnowledgeSearchResponse{}, err
	}

	ustawienia := a.ustawienia()
	silnik := a.silnik(ustawienia)
	okno, zasady, obszar := a.zasiegPlatformy()

	wektory, err := silnik.Osadz(ctx, okno, zasady, obszar,
		[]string{pytanie}, wiedza.LimitZapytania)
	if err != nil {
		return shared.KnowledgeSearchResponse{}, bladWiedzy(err)
	}
	if len(wektory) != 1 {
		return shared.KnowledgeSearchResponse{}, bladWiedzy(errors.New(
			"pomocnik osadzeń nie oddał wektora pytania; naprawa: zgłosić usterkę pomocnika"))
	}

	pozycje, err := a.skladnica.Pozycje(ctx, zakresy, ustawienia.Model)
	if err != nil {
		return shared.KnowledgeSearchResponse{}, bladWiedzy(err)
	}
	granica := granicaZadania(z.Limit)
	wektorPytania := wiedza.Znormalizuj(wektory[0])
	if przesiewZadany(z.Rerank) {
		return a.zPrzesiewem(ctx, okno, zasady, obszar, ustawienia, pytanie,
			wektorPytania, pozycje, z, granica)
	}
	trafienia := wiedza.Najblizsze(pozycje, wektorPytania, granica)
	return odpowiedzSzukania(trafienia, false), nil
}

// zPrzesiewem przeprowadza drugi przebieg: bierze kandydatów pierwszego
// przebiegu i oddaje ich w kolejności ułożonej przez krzyżowy koder.
//
// Kandydatów jest więcej niż oddawanych fragmentów i to jest sens rzeczy —
// przesiew może wynieść na czoło fragment, który po samych wektorach był
// dwudziesty. Gdyby kandydatami było dokładnie tyle, ile fragmentów wraca,
// przesiew przestawiałby wyłącznie kolejność wewnątrz zbioru już wybranego.
//
// Odmowa przesiewu jest odmową całego żądania, a nie zejściem na wynik
// pierwszego przebiegu. Wołający prosił o kolejność ułożoną na nowo; oddanie mu
// po cichu tej samej kolejności, którą miał bez pytania, byłoby odpowiedzią
// nierozpoznawalnie gorszą — pole `reranked` istnieje właśnie po to, żeby
// odróżnienie było możliwe, ale milcząca podmiana czyniłaby je kłamstwem.
func (a *adapterWiedzy) zPrzesiewem(ctx context.Context, okno session.Okno,
	zasady session.Zasady, obszar session.Obszar, ustawienia wiedza.Ustawienia,
	pytanie string, wektorPytania []float32, pozycje []wiedza.Pozycja,
	z shared.KnowledgeSearchRequest, granica int) (shared.KnowledgeSearchResponse, error) {

	kandydaci := wiedza.Najblizsze(pozycje, wektorPytania,
		wiedza.GranicaKandydatowZadania(z.RerankCandidates))
	if len(kandydaci) == 0 {
		return odpowiedzSzukania(nil, true), nil
	}

	teksty := make([]string, len(kandydaci))
	for i, kandydat := range kandydaci {
		teksty[i] = kandydat.Pozycja.Tresc
	}
	oceny, err := wiedza.NowySilnikPrzesiewu(a.uruchamiacz, a.katalogDanych).
		ZUstawieniami(ustawienia).
		Przesiej(ctx, okno, zasady, obszar, pytanie, teksty, wiedza.LimitPrzesiewu)
	if err != nil {
		return shared.KnowledgeSearchResponse{}, bladWiedzy(err)
	}
	return odpowiedzSzukania(wiedza.PoPrzesiewie(kandydaci, oceny, granica), true), nil
}

// odpowiedzSzukania składa odpowiedź kontraktu z wykazu trafień.
//
// Pole `reranked` wchodzi zawsze, gdy przesiew się odbył, także przy wykazie
// pustym: „nie mam na to nic" po przesiewie i „nie mam na to nic" bez niego są
// dwiema różnymi odpowiedziami i wołający ma prawo je rozróżnić.
func odpowiedzSzukania(trafienia []wiedza.Trafienie, przesiane bool) shared.KnowledgeSearchResponse {
	wyniki := make([]shared.KnowledgeHit, 0, len(trafienia))
	for _, trafienie := range trafienia {
		wyniki = append(wyniki, przelozTrafienie(trafienie))
	}
	odpowiedz := shared.KnowledgeSearchResponse{Results: wyniki, Total: len(wyniki)}
	if przesiane {
		tak := true
		odpowiedz.Reranked = &tak
	}
	return odpowiedz
}

// przesiewZadany czyta wskazanie żądania. Brak pola znaczy pierwszy przebieg
// sam — przesiew kosztuje wczytanie drugiego modelu, więc wchodzi wyłącznie na
// wyraźne żądanie.
func przesiewZadany(rerank *bool) bool {
	return rerank != nil && *rerank
}

// przelozTrafienie składa pozycję kontraktu zawsze ze wskazaniem źródła.
// Kod źródła pusty oddaje `sourceId` niewypełnione — pole jest opcjonalne
// i pusty napis udawałby identyfikator, którego nie ma.
func przelozTrafienie(trafienie wiedza.Trafienie) shared.KnowledgeHit {
	trafnosc := wiedza.WSetnych(trafienie.Podobienst)
	pozycja := shared.KnowledgeHit{
		Text:   trafienie.Pozycja.Tresc,
		Source: trafienie.Pozycja.Zrodlo,
		Score:  &trafnosc,
		Scope:  shared.KnowledgeScope(trafienie.Pozycja.Zakres),
	}
	if kod := trafienie.Pozycja.ZrodloKod; kod != "" {
		pozycja.SourceId = &kod
	}
	return pozycja
}

// granicaZadania rozstrzyga liczbę oddawanych fragmentów.
func granicaZadania(limit *int) int {
	if limit == nil || *limit <= 0 {
		return domyslnaLiczbaTrafien
	}
	if *limit > granicaTrafien {
		return granicaTrafien
	}
	return *limit
}

// zakresyZadania rozwija wskazanie kontraktu na wykaz zakresów wskaźnika.
//
// Brak wskazania bierze bibliotekę — tak stanowi kontrakt wprost („Zakres
// wskaznika; brak bierze biblioteke"). `all` rozwija się na trzy zakresy, bo
// „wszystko" nie jest czwartym miejscem, z którego coś pochodzi.
func zakresyZadania(zakres *shared.KnowledgeScope) ([]string, error) {
	if zakres == nil {
		return []string{shared.KnowledgeScopeLibrary}, nil
	}
	switch *zakres {
	case shared.KnowledgeScopeLibrary, shared.KnowledgeScopeHistory, shared.KnowledgeScopeWorkspace:
		return []string{string(*zakres)}, nil
	case shared.KnowledgeScopeAll:
		return []string{shared.KnowledgeScopeLibrary, shared.KnowledgeScopeHistory,
			shared.KnowledgeScopeWorkspace}, nil
	default:
		return nil, bladZadaniaWiedzy("zakres `" + string(*zakres) +
			"` nie jest znany — wykaz kontraktu to `library`, `history`, `workspace` i `all`")
	}
}

// ustawienia składa komplet nastaw z konfiguracji zasięgu platformy.
// Klucz bez rozstrzygnięcia zostaje przy wartości domyślnej.
func (a *adapterWiedzy) ustawienia() wiedza.Ustawienia {
	komplet := wiedza.UstawieniaDomyslne()
	if a.rozstrzygacz == nil {
		return komplet
	}
	for _, klucz := range []string{wiedza.KluczProgram, wiedza.KluczModel,
		wiedza.KluczKatalogModeli, wiedza.KluczDlugoscFragmentu} {

		wynik := a.rozstrzygacz.Rozstrzygnij(konfig.Kontekst{}, klucz)
		if wynik.Pochodzenie == konfig.PochodzenieNieznane {
			continue
		}
		komplet.Nanies(klucz, wynik.Wartosc)
	}
	return komplet
}

// silnik składa silnik osadzeń na nastawach tego wywołania.
func (a *adapterWiedzy) silnik(ustawienia wiedza.Ustawienia) *wiedza.Silnik {
	return wiedza.NowySilnik(a.uruchamiacz, a.katalogDanych).ZUstawieniami(ustawienia)
}

// zasiegPlatformy składa trójkę okno–zasady–obszar dla zasięgu platformy.
// Rozstrzygnięcie i jego uzasadnienie stoją w nagłówku pliku.
func (a *adapterWiedzy) zasiegPlatformy() (session.Okno, session.Zasady, session.Obszar) {
	okno := session.Okno{Ustawienia: session.Ustawienia{
		SrodowiskoWykonania: shared.ExecutionEnvCore,
	}}
	zasady := session.Zasady{}
	if a.rozstrzygacz != nil {
		zasady = ZasadyIzolacji(a.rozstrzygacz, konfig.Kontekst{})
	}
	obszar := session.Obszar{}
	if a.katalog != nil {
		obszar = ObszarOkna(a.katalog.Ustal(konfig.Kontekst{}, ""), "")
	}
	return okno, zasady, obszar
}

// bladWiedzy znakuje odmowę pakietu kodem kontraktu.
//
//   - `wiedza.BrakSilnika` → `channel_unavailable`, w praktyce nieponawialny
//     mimo ponawialności kodu: zaplecze liczenia jest niedostępne i to samo
//     żądanie powiedzie się bez zmiany dopiero po naprawie z treści odmowy.
//     Kod odróżniający „nie ma biblioteki" od „nie ma wag" w katalogu kontraktu
//     nie istnieje; rozróżnienie niesie treść, trójczęściowa i różna dla obu.
//   - naruszenie izolacji → `permission_denied`, tak samo jak znakuje je
//     Terminal — dwie reguły dla jednej bramy byłyby rozjazdem.
//   - reszta → `internal_error`. Kod domyślny jest najostrzejszy z zamysłem:
//     nieznana odmowa jest przypadkiem, którego rdzeń nie przewidział.
func bladWiedzy(err error) error {
	var brakSilnika *wiedza.BrakSilnika
	if errors.As(err, &brakSilnika) {
		return odmowaWiedzy(shared.ErrorCodeChannelUnavailable, err.Error())
	}
	if errors.Is(err, session.ErrIzolacja) {
		return odmowaWiedzy(shared.ErrorCodePermissionDenied, err.Error())
	}
	return odmowaWiedzy(shared.ErrorCodeInternalError, err.Error())
}

// bladZadaniaWiedzy znakuje wadę żądania kodem kontraktu.
func bladZadaniaWiedzy(powod string) error {
	return odmowaWiedzy(shared.ErrorCodeValidationFailed, powod)
}

// odmowaWiedzy składa odmowę obszaru. Przedrostek nazywa obszar, żeby czytający
// wiedział, kto odmówił, zanim przeczyta dlaczego.
func odmowaWiedzy(kod shared.ErrorCode, powod string) error {
	return protocol.JakoError(protocol.NowyBlad(kod, "wyszukiwanie po znaczeniu: "+powod))
}
