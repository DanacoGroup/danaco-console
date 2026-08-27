// Odpowiedzialność pliku: adnotacje kompozycji Design Board wraz z wątkami —
// design.annotation.set i design.annotation.list; adnotacja ma własny
// wiersz, czas i autora, branego z kontekstu wywołania.
package core

import (
	"context"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// przedrostekAdnotacjiDesign znakuje identyfikatory zewnętrzne adnotacji,
// żeby dało się je odróżnić od pozostałych bytów kompozycji.
const przedrostekAdnotacjiDesign = "adnotacja-"

// UstawAdnotacje zakłada albo zmienia adnotację — obsługuje
// design.annotation.set; adnotacja nadrzędna musi należeć do TEJ SAMEJ
// kompozycji.
func (a *adapterDesignu) UstawAdnotacje(ctx context.Context,
	z shared.DesignAnnotationSetRequest) (shared.DesignAnnotationSetResponse, error) {

	if strings.TrimSpace(z.BoardId) == "" {
		return shared.DesignAnnotationSetResponse{}, bladWskazaniaDesignu(
			"komenda design.annotation.set bez wskazania kompozycji")
	}
	if strings.TrimSpace(z.Text) == "" {
		return shared.DesignAnnotationSetResponse{}, bladWskazaniaDesignu(
			"komenda design.annotation.set bez treści adnotacji — pusta uwaga nie mówi niczego, " +
				"a zajmuje miejsce w wątku")
	}
	kompozycja, err := a.repozytorium.Kompozycja(ctx, strings.TrimSpace(z.BoardId))
	if err != nil {
		return shared.DesignAnnotationSetResponse{}, bladNieznanejKompozycjiDesignu(z.BoardId, err)
	}

	zastane, err := a.repozytorium.AdnotacjeKompozycjiDesignu(ctx, kompozycja.ID, false)
	if err != nil {
		return shared.DesignAnnotationSetResponse{}, bladDesignu(err)
	}

	kod := nowyIdentyfikator(przedrostekAdnotacjiDesign)
	var autor *string
	if z.AnnotationId != nil && strings.TrimSpace(*z.AnnotationId) != "" {
		kod = strings.TrimSpace(*z.AnnotationId)
		poprzednia, znana := adnotacjaZWykazuDesignu(zastane, kod)
		if !znana {
			return shared.DesignAnnotationSetResponse{}, bladNieznanegoBytuDesignu(
				"adnotacji " + kod + " nie ma w kompozycji " + kompozycja.Kod)
		}
		// Autor zostaje ten, kto adnotację założył — zmiana treści nie czyni
		// jej autorką cudzej uwagi.
		autor = poprzednia.Autor
	} else {
		autor = autorAdnotacjiDesignu(ctx)
	}

	if z.ParentId != nil && strings.TrimSpace(*z.ParentId) != "" {
		nadrzedna := strings.TrimSpace(*z.ParentId)
		if nadrzedna == kod {
			return shared.DesignAnnotationSetResponse{}, bladWskazaniaDesignu(
				"komenda design.annotation.set czyni adnotację " + kod + " odpowiedzią na samą siebie")
		}
		if _, znana := adnotacjaZWykazuDesignu(zastane, nadrzedna); !znana {
			return shared.DesignAnnotationSetResponse{}, bladNieznanegoBytuDesignu(
				"adnotacji nadrzędnej " + nadrzedna + " nie ma w kompozycji " + kompozycja.Kod)
		}
	}

	zapisana, err := a.repozytorium.ZapiszAdnotacjeDesignu(ctx, dane.AdnotacjaDesignu{
		Kod:          kod,
		KompozycjaID: kompozycja.ID,
		WarstwaID:    z.LayerId,
		NadrzednaID:  z.ParentId,
		Autor:        autor,
		Tresc:        z.Text,
		Zamknieta:    z.Resolved != nil && *z.Resolved,
	})
	if err != nil {
		return shared.DesignAnnotationSetResponse{}, bladDesignu(err)
	}
	return shared.DesignAnnotationSetResponse{
		Annotation: adnotacjaKontraktuDesignu(zapisana, kompozycja.Kod),
	}, nil
}

// Adnotacje zwraca adnotacje kompozycji wraz z wątkami — obsługuje
// design.annotation.list; zawężenie do wątków niezamkniętych działa po
// adnotacji, nie po wątku.
func (a *adapterDesignu) Adnotacje(ctx context.Context,
	z shared.DesignAnnotationListRequest) (shared.DesignAnnotationListResponse, error) {

	if strings.TrimSpace(z.BoardId) == "" {
		return shared.DesignAnnotationListResponse{}, bladWskazaniaDesignu(
			"komenda design.annotation.list bez wskazania kompozycji")
	}
	kompozycja, err := a.repozytorium.Kompozycja(ctx, strings.TrimSpace(z.BoardId))
	if err != nil {
		return shared.DesignAnnotationListResponse{}, bladNieznanejKompozycjiDesignu(z.BoardId, err)
	}
	wiersze, err := a.repozytorium.AdnotacjeKompozycjiDesignu(ctx, kompozycja.ID,
		z.OpenOnly != nil && *z.OpenOnly)
	if err != nil {
		return shared.DesignAnnotationListResponse{}, bladDesignu(err)
	}
	adnotacje := make([]shared.DesignAnnotation, 0, len(wiersze))
	for _, wiersz := range wiersze {
		adnotacje = append(adnotacje, adnotacjaKontraktuDesignu(wiersz, kompozycja.Kod))
	}
	return shared.DesignAnnotationListResponse{Annotations: adnotacje, Total: len(adnotacje)}, nil
}

// adnotacjaZWykazuDesignu odszukuje adnotację po kodzie w wykazie już
// odczytanym, sprawdzanym oboma wskazaniami: zmienianym i nadrzędnym.
func adnotacjaZWykazuDesignu(wykaz []dane.AdnotacjaDesignu, kod string) (dane.AdnotacjaDesignu, bool) {
	for _, adnotacja := range wykaz {
		if adnotacja.Kod == kod {
			return adnotacja, true
		}
	}
	return dane.AdnotacjaDesignu{}, false
}

// autorAdnotacjiDesignu bierze autora z kontekstu wywołania; brak wskazania
// klienta zostawia pole puste, żeby nie wpisywać podpisu, którego nie było.
func autorAdnotacjiDesignu(ctx context.Context) *string {
	_, klient := sprawca(ctx)
	if klient == nil || strings.TrimSpace(*klient) == "" {
		return nil
	}
	return klient
}

// adnotacjaKontraktuDesignu składa DesignAnnotation kontraktu z wiersza
// repozytorium, w kształcie oczekiwanym przez odpowiedź komendy.
func adnotacjaKontraktuDesignu(a dane.AdnotacjaDesignu, kodKompozycji string) shared.DesignAnnotation {
	zamknieta := a.Zamknieta
	return shared.DesignAnnotation{
		Id:        a.Kod,
		BoardId:   kodKompozycji,
		LayerId:   a.WarstwaID,
		ParentId:  a.NadrzednaID,
		Author:    a.Autor,
		Text:      a.Tresc,
		Resolved:  &zamknieta,
		CreatedAt: chwilaBazy(a.Utworzono),
	}
}
