// Wypełnienie komendy knowledge.image.search: obrazy bierze się z biblioteki po rodzaju
// treści zapisanym przy wgraniu, nie po rozszerzeniu nazwy pliku.
package core

import (
	"context"
	"errors"
	"log"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/wiedza"
	"danacoconsole/shared"
)

const (
	// domyslnaLiczbaObrazow jest niższa niż dla fragmentów tekstu: obraz wskazuje plik, nie niesie odpowiedzi.
	domyslnaLiczbaObrazow    = 5
	granicaLiczbyObrazow     = 50
	przedrostekRodzajuObrazu = "image/"
)

// Kolejność jak przy budowie wskaźnika: gotowość modelu przed odczytem biblioteki.
func (a *adapterWiedzy) SzukajObrazu(ctx context.Context,
	z shared.KnowledgeImageSearchRequest) (shared.KnowledgeImageSearchResponse, error) {

	pytanie := strings.TrimSpace(z.Query)
	if pytanie == "" {
		return shared.KnowledgeImageSearchResponse{}, bladZadaniaWiedzy(
			"komenda bez zdania opisującego obraz — oś obrazu nie ma czego porównać " +
				"z obrazami Operatora; naprawa: podać opis w polu `query`")
	}

	ustawienia := a.ustawienia()
	silnik := wiedza.NowySilnikObrazu(a.uruchamiacz, a.katalogDanych).
		ZUstawieniami(ustawienia).ZeSkladnica(a.skladnica)
	okno, zasady, obszar := a.zasiegPlatformy()
	if err := silnik.Gotowy(ctx, okno, zasady, obszar, wiedza.LimitOsiObrazu); err != nil {
		return shared.KnowledgeImageSearchResponse{}, bladWiedzy(err)
	}

	obrazy, obciete, err := a.obrazyBiblioteki(ctx, z.ProjectId)
	if err != nil {
		return shared.KnowledgeImageSearchResponse{}, err
	}
	if obciete {
		log.Printf("moduł Wiedza: oś obrazu porównała sufit %d obrazów biblioteki "+
			"(GranicaObrazow) — biblioteka niesie więcej obrazów niż to weszło do porównania",
			wiedza.GranicaObrazow)
	}
	model := silnik.Model()
	if len(obrazy) == 0 {
		zero := 0
		return shared.KnowledgeImageSearchResponse{
			Results: []shared.KnowledgeImageHit{}, Total: 0,
			Examined: &zero, Model: &model,
		}, nil
	}

	sciezki := make([]string, len(obrazy))
	for i, obraz := range obrazy {
		sciezki[i] = obraz.Sciezka
	}
	oceny, _, err := silnik.Dopasuj(ctx, okno, zasady, obszar, pytanie, sciezki,
		wiedza.LimitOsiObrazu)
	if errors.Is(err, dane.ErrKolizjaWiersza) {
		return shared.KnowledgeImageSearchResponse{}, odmowaWiedzy(shared.ErrorCodeConflict, err.Error())
	}
	if err != nil {
		return shared.KnowledgeImageSearchResponse{}, bladWiedzy(err)
	}

	najblizsze := wiedza.NajblizszeObrazy(obrazy, oceny, granicaObrazowZadania(z.Limit))
	wyniki := make([]shared.KnowledgeImageHit, 0, len(najblizsze))
	for _, obraz := range najblizsze {
		wyniki = append(wyniki, przelozObraz(obraz))
	}
	przejrzane := len(obrazy)
	return shared.KnowledgeImageSearchResponse{
		Results: wyniki, Total: len(wyniki), Examined: &przejrzane, Model: &model,
	}, nil
}

func (a *adapterWiedzy) obrazyBiblioteki(ctx context.Context,
	projekt *string) ([]wiedza.Obraz, bool, error) {

	if a.biblioteka == nil {
		return nil, false, bladWiedzyBezZrodla("biblioteka")
	}
	filtr := dane.FiltrPlikow{Limit: granicaDokumentowBiblioteki}
	if projekt != nil && strings.TrimSpace(*projekt) != "" {
		filtr.ProjektID = projekt
	}
	pliki, _, err := a.biblioteka.Pliki(ctx, filtr)
	if err != nil {
		return nil, false, bladWiedzy(err)
	}

	obrazy := make([]wiedza.Obraz, 0, len(pliki))
	obciete := false
	for _, plik := range pliki {
		if plik.MimeType == nil || !strings.HasPrefix(*plik.MimeType, przedrostekRodzajuObrazu) {
			continue
		}
		if plik.TrescOdwolanie == nil || strings.TrimSpace(*plik.TrescOdwolanie) == "" {
			continue
		}
		if len(obrazy) >= wiedza.GranicaObrazow {
			obciete = true
			continue
		}
		obrazy = append(obrazy, wiedza.Obraz{
			Zrodlo:    plik.Nazwa,
			ZrodloKod: plik.Kod,
			Rodzaj:    *plik.MimeType,
			Sciezka:   *plik.TrescOdwolanie,
		})
	}
	return obrazy, obciete, nil
}

// Ścieżka na dysku nie wchodzi do odpowiedzi: wynosiłaby układ dysku Operatora.
func przelozObraz(obraz wiedza.Obraz) shared.KnowledgeImageHit {
	trafnosc := wiedza.WSetnych(obraz.Podobienst)
	pozycja := shared.KnowledgeImageHit{Source: obraz.Zrodlo, Score: &trafnosc}
	if kod := obraz.ZrodloKod; kod != "" {
		pozycja.SourceId = &kod
	}
	if rodzaj := obraz.Rodzaj; rodzaj != "" {
		pozycja.MimeType = &rodzaj
	}
	return pozycja
}

func granicaObrazowZadania(limit *int) int {
	if limit == nil || *limit <= 0 {
		return domyslnaLiczbaObrazow
	}
	if *limit > granicaLiczbyObrazow {
		return granicaLiczbyObrazow
	}
	return *limit
}
