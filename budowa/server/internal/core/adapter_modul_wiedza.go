// Plik wypełnia port Wiedza dwiema komendami rodziny `knowledge.*`: przekłada
// żądania kontraktu na zlecenia pakietu `wiedza` i jego typowane odmowy na
// kody kontraktu kanału; oś obrazu stoi w innym pliku tej samej warstwy.
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

type adapterWiedzy struct {
	uruchamiacz   session.Uruchamiacz
	katalogDanych string
	rozstrzygacz  *konfig.Rozstrzygacz
	katalog       *KatalogRoboczy
	skladnica     *wiedza.Skladnica
	biblioteka    dane.RepozytoriumBiblioteki
	historia      dane.RepozytoriumHistorii
	teraz         func() int64
}

func nowyAdapterWiedzy(uruchamiacz session.Uruchamiacz, katalogDanych string) *adapterWiedzy {
	return &adapterWiedzy{
		uruchamiacz:   uruchamiacz,
		katalogDanych: katalogDanych,
		teraz:         func() int64 { return time.Now().UTC().UnixMilli() },
	}
}

func skladnicaWiedzy(m Montaz) *wiedza.Skladnica {
	if m.Baza == nil {
		return nil
	}
	return wiedza.NowaSkladnica(m.Baza.DB)
}

func (a *adapterWiedzy) ZIzolacja(rozstrzygacz *konfig.Rozstrzygacz, katalog *KatalogRoboczy) *adapterWiedzy {
	a.rozstrzygacz, a.katalog = rozstrzygacz, katalog
	return a
}

func (a *adapterWiedzy) ZeSkladnica(s *wiedza.Skladnica) *adapterWiedzy {
	a.skladnica = s
	return a
}

func (a *adapterWiedzy) ZeZrodlami(biblioteka dane.RepozytoriumBiblioteki,
	historia dane.RepozytoriumHistorii) *adapterWiedzy {

	a.biblioteka, a.historia = biblioteka, historia
	return a
}

func (a *adapterWiedzy) Wskaznik(ctx context.Context,
	z shared.KnowledgeIndexRequest) (shared.KnowledgeIndexResponse, error) {

	zakresy, err := zakresyZadania(z.Scope)
	if err != nil {
		return shared.KnowledgeIndexResponse{}, err
	}
	ustawienia := a.ustawienia(ctx)
	silnik := a.silnik(ustawienia)
	okno, zasady, obszar := a.zasiegPlatformy(ctx)

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

		// Kasowanie przed zapisem usuwa fragmenty dokumentu skróconego od
		// poprzedniego przebiegu.
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

	ustawienia := a.ustawienia(ctx)
	silnik := a.silnik(ustawienia)
	okno, zasady, obszar := a.zasiegPlatformy(ctx)

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

func przesiewZadany(rerank *bool) bool {
	return rerank != nil && *rerank
}

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

func granicaZadania(limit *int) int {
	if limit == nil || *limit <= 0 {
		return domyslnaLiczbaTrafien
	}
	if *limit > granicaTrafien {
		return granicaTrafien
	}
	return *limit
}

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

func (a *adapterWiedzy) ustawienia(ctx context.Context) wiedza.Ustawienia {
	komplet := wiedza.UstawieniaDomyslne()
	if a.rozstrzygacz == nil {
		return komplet
	}
	kontekst := konfig.Kontekst{KontoOperatora: dane.KontoOperatora(ctx)}
	for _, klucz := range []string{wiedza.KluczProgram, wiedza.KluczModel,
		wiedza.KluczKatalogModeli, wiedza.KluczDlugoscFragmentu,
		wiedza.KluczModelPrzesiewu, wiedza.KluczKatalogPrzesiewu,
		wiedza.KluczModelObrazu, wiedza.KluczKatalogObrazu} {

		wynik := a.rozstrzygacz.Rozstrzygnij(kontekst, klucz)
		if wynik.Pochodzenie == konfig.PochodzenieNieznane {
			continue
		}
		komplet.Nanies(klucz, wynik.Wartosc)
	}
	return komplet
}

// silnik składa silnik osadzeń na nastawach danego wywołania, tworzony na
// nowo za każdym razem, bo `ZUstawieniami` mutuje jego stan wewnętrzny.
func (a *adapterWiedzy) silnik(ustawienia wiedza.Ustawienia) *wiedza.Silnik {
	return wiedza.NowySilnik(a.uruchamiacz, a.katalogDanych).ZUstawieniami(ustawienia)
}

func (a *adapterWiedzy) zasiegPlatformy(ctx context.Context) (session.Okno, session.Zasady, session.Obszar) {
	okno := session.Okno{Ustawienia: session.Ustawienia{
		SrodowiskoWykonania: shared.ExecutionEnvCore,
	}}
	kontekst := konfig.Kontekst{KontoOperatora: dane.KontoOperatora(ctx)}
	zasady := session.Zasady{}
	if a.rozstrzygacz != nil {
		zasady = ZasadyIzolacji(a.rozstrzygacz, kontekst)
	}
	obszar := session.Obszar{}
	if a.katalog != nil {
		obszar = ObszarOkna(a.katalog.Ustal(kontekst, ""), "")
	}
	return okno, zasady, obszar
}

// bladWiedzy znakuje odmowę pakietu kodem kontraktu: brak silnika daje
// `channel_unavailable`, naruszenie izolacji `permission_denied`, a każda
// inna odmowa najostrzejszy kod `internal_error`.
func bladWiedzy(err error) error {
	var brakSilnika *wiedza.BrakSilnika
	if errors.As(err, &brakSilnika) {
		return odmowaWiedzy(shared.ErrorCodeChannelUnavailable, err.Error())
	}
	if errors.Is(err, session.ErrIzolacja) {
		return odmowaWiedzy(shared.ErrorCodePermissionDenied, err.Error())
	}
	if errors.Is(err, dane.ErrKolizjaWiersza) {
		return odmowaWiedzy(shared.ErrorCodeConflict, err.Error())
	}
	return odmowaWiedzy(shared.ErrorCodeInternalError, err.Error())
}

func bladZadaniaWiedzy(powod string) error {
	return odmowaWiedzy(shared.ErrorCodeValidationFailed, powod)
}

func odmowaWiedzy(kod shared.ErrorCode, powod string) error {
	return protocol.JakoError(protocol.NowyBlad(kod, "wyszukiwanie po znaczeniu: "+powod))
}
