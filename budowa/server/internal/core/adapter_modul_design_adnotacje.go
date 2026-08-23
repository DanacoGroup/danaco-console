// Odpowiedzialność pliku: adnotacje kompozycji Design Board wraz z wątkami —
// `design.annotation.set` i `design.annotation.list`. Metody stoją na
// `*adapterDesignu` (`adapter_modul_design.go`).
//
// Pole `note` warstwy niesie JEDNO zdanie bez autora i bez wątku, a przy każdym
// `design.board.update` jedzie z całym układem i wraca przepisane od nowa —
// uwaga jednej osoby znikała więc przy pierwszym przesunięciu warstwy przez
// drugą. Adnotacja ma własny wiersz, własny czas i własnego autora.
//
// Autora bierzemy z kontekstu wywołania (`sprawca`), nie z żądania. Autor
// podany przez wołającego byłby polem, w które da się wpisać cudze nazwisko,
// a oznaczenia osób w wątku mają znaczyć to, co znaczą.
package core

import (
	"context"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// przedrostekAdnotacjiDesign znakuje identyfikatory zewnętrzne adnotacji.
const przedrostekAdnotacjiDesign = "adnotacja-"

// UstawAdnotacje zakłada albo zmienia adnotację — obsługuje
// `design.annotation.set`.
//
// Adnotacja nadrzędna musi należeć do TEJ SAMEJ kompozycji. Wątek rozpięty
// między dwiema tablicami nie jest wątkiem: druga tablica pokazywałaby
// odpowiedź na uwagę, której u siebie nie ma.
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
		// Autor zostaje ten, kto adnotację założył — zmiana treści przez drugą
		// osobę nie czyni jej autorką cudzej uwagi.
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
// `design.annotation.list`.
//
// Zawężenie do wątków niezamkniętych działa po adnotacji, nie po wątku: wątek
// zamyka się adnotacja po adnotacji, a odpowiedź otwarta pod zamkniętą uwagą
// jest sprawą wciąż otwartą i ma zostać widoczna.
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
// odczytanym. Wykaz kompozycji czytamy raz i sprawdzamy nim oba wskazania
// (zmienianą i nadrzędną) — dwa odczyty po kodzie robiłyby tę samą pracę drugi
// raz i mogły trafić na stan zmieniony w międzyczasie.
func adnotacjaZWykazuDesignu(wykaz []dane.AdnotacjaDesignu, kod string) (dane.AdnotacjaDesignu, bool) {
	for _, adnotacja := range wykaz {
		if adnotacja.Kod == kod {
			return adnotacja, true
		}
	}
	return dane.AdnotacjaDesignu{}, false
}

// autorAdnotacjiDesignu bierze autora z kontekstu wywołania. Brak wskazania
// klienta zostawia pole puste — „nieznany" wpisany w kolumnę wyglądałby przy
// uwadze jak czyjś podpis.
func autorAdnotacjiDesignu(ctx context.Context) *string {
	_, klient := sprawca(ctx)
	if klient == nil || strings.TrimSpace(*klient) == "" {
		return nil
	}
	return klient
}

// adnotacjaKontraktuDesignu składa `DesignAnnotation` kontraktu z wiersza
// repozytorium.
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
