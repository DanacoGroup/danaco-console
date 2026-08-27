// Odpowiedzialność pliku: wypełnienie komendy `knowledge.image.search` — osi
// obrazu rodziny `knowledge.*`. Obrazy bierze się z biblioteki, nie
// z przestrzeni roboczej, po rodzaju treści zapisanym przy wgraniu, nie po
// rozszerzeniu nazwy pliku.
package core

import (
	"context"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/wiedza"
	"danacoconsole/shared"
)

const (
	// domyslnaLiczbaObrazow wchodzi, gdy żądanie nie poda `limit`, mniej niż
	// fragmentów tekstu, bo obraz jest wskazaniem pliku, nie podstawą odpowiedzi.
	domyslnaLiczbaObrazow = 5
	// granicaLiczbyObrazow chroni przed żądaniem, które chce całą bibliotekę
	// obrazów naraz, jednym przebiegiem wyszukiwania.
	granicaLiczbyObrazow = 50
	// przedrostekRodzajuObrazu rozpoznaje rodzaj treści będący obrazem, po
// prefiksie MIME zapisanym przy wgraniu pliku do biblioteki.
	przedrostekRodzajuObrazu = "image/"
)

// SzukajObrazu obsługuje `knowledge.image.search`. Kolejność kroków jest ta
// sama co przy budowaniu wskaźnika: najpierw pytanie o gotowość modelu,
// dopiero potem odczyt biblioteki. Wynik pusty jest odpowiedzią, nie odmową.
func (a *adapterWiedzy) SzukajObrazu(ctx context.Context,
	z shared.KnowledgeImageSearchRequest) (shared.KnowledgeImageSearchResponse, error) {

	pytanie := strings.TrimSpace(z.Query)
	if pytanie == "" {
		return shared.KnowledgeImageSearchResponse{}, bladZadaniaWiedzy(
			"komenda bez zdania opisującego obraz — oś obrazu nie ma czego porównać " +
				"z obrazami Operatora; naprawa: podać opis w polu `query`")
	}

	ustawienia := a.ustawienia()
	silnik := wiedza.NowySilnikObrazu(a.uruchamiacz, a.katalogDanych).ZUstawieniami(ustawienia)
	okno, zasady, obszar := a.zasiegPlatformy()
	if err := silnik.Gotowy(ctx, okno, zasady, obszar, wiedza.LimitOsiObrazu); err != nil {
		return shared.KnowledgeImageSearchResponse{}, bladWiedzy(err)
	}

	obrazy, err := a.obrazyBiblioteki(ctx, z.ProjectId)
	if err != nil {
		return shared.KnowledgeImageSearchResponse{}, err
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

// obrazyBiblioteki wybiera z repozytorium wiedzy pliki będące obrazami. Wykaz
// idzie przez to samo repozytorium, którym czyta bibliotekę budowanie
// wskaźnika. Plik bez odwołania do bajtów jest pomijany.
func (a *adapterWiedzy) obrazyBiblioteki(ctx context.Context,
	projekt *string) ([]wiedza.Obraz, error) {

	if a.biblioteka == nil {
		return nil, bladWiedzyBezZrodla("biblioteka")
	}
	filtr := dane.FiltrPlikow{Limit: granicaDokumentowBiblioteki}
	if projekt != nil && strings.TrimSpace(*projekt) != "" {
		filtr.ProjektID = projekt
	}
	pliki, _, err := a.biblioteka.Pliki(ctx, filtr)
	if err != nil {
		return nil, bladWiedzy(err)
	}

	obrazy := make([]wiedza.Obraz, 0, len(pliki))
	for _, plik := range pliki {
		if len(obrazy) >= wiedza.GranicaObrazow {
			break
		}
		if plik.MimeType == nil || !strings.HasPrefix(*plik.MimeType, przedrostekRodzajuObrazu) {
			continue
		}
		if plik.TrescOdwolanie == nil || strings.TrimSpace(*plik.TrescOdwolanie) == "" {
			continue
		}
		obrazy = append(obrazy, wiedza.Obraz{
			Zrodlo:    plik.Nazwa,
			ZrodloKod: plik.Kod,
			Rodzaj:    *plik.MimeType,
			Sciezka:   *plik.TrescOdwolanie,
		})
	}
	return obrazy, nil
}

// przelozObraz składa pozycję kontraktu ze wskazaniem źródła. Ścieżka na dysku
// do odpowiedzi nie wchodzi, wynosiłaby wołającemu układ dysku Operatora.
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

// granicaObrazowZadania rozstrzyga liczbę oddawanych obrazów, między
// domyślną a górnym sufitem chroniącym przed pełnym przebiegiem biblioteki.
func granicaObrazowZadania(limit *int) int {
	if limit == nil || *limit <= 0 {
		return domyslnaLiczbaObrazow
	}
	if *limit > granicaLiczbyObrazow {
		return granicaLiczbyObrazow
	}
	return *limit
}
