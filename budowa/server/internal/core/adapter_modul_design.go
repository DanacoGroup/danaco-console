// Moduł Design: typ adaptera, konstruktor, prompt strukturalny i wykaz zasobów;
// wytworzenie, kanał obrazowy, wniesienie, kompozycje i etykiety leżą w osobnych plikach.
package core

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/konfiguracja"
	"danacoconsole/server/internal/models"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

type adapterDesignu struct {
	repozytorium dane.RepozytoriumDesignu
	kanaly       *models.Rejestr
	// repozytoriumKanalow rozstrzyga własność kanału z żądania (decyzja 34).
	repozytoriumKanalow dane.RepozytoriumKanalow
	magazyn             *magazynTresciBiblioteki
	obecnosc            *rejestrObecnosciDesignu
	sejf                SejfPoswiadczen
	uruchamiacz         session.Uruchamiacz
	rozstrzygacz        *konfig.Rozstrzygacz
	katalogRoboczy      *KatalogRoboczy
}

func (a *adapterDesignu) ZOdczytemPisma(uruchamiacz session.Uruchamiacz,
	rozstrzygacz *konfig.Rozstrzygacz, katalog *KatalogRoboczy) *adapterDesignu {

	a.uruchamiacz, a.rozstrzygacz, a.katalogRoboczy = uruchamiacz, rozstrzygacz, katalog
	return a
}

func (a *adapterDesignu) ZSejfem(sejf SejfPoswiadczen) *adapterDesignu {
	a.sejf = sejf
	return a
}

func nowyAdapterDesignu(repozytorium dane.RepozytoriumDesignu) *adapterDesignu {
	return &adapterDesignu{
		repozytorium: repozytorium,
		magazyn:      magazynZasobowDesignu(konfiguracja.KatalogDanychDomyslny()),
		obecnosc:     nowyRejestrObecnosciDesignu(),
	}
}

// Bez tego zasoby leżałyby w katalogu domyślnym, a baza tam, gdzie wskazuje konfiguracja.
func (a *adapterDesignu) ZKatalogiemDanych(katalog string) *adapterDesignu {
	if katalog != "" {
		a.magazyn = magazynZasobowDesignu(katalog)
	}
	return a
}

func (a *adapterDesignu) ZKanalami(kanaly *models.Rejestr,
	repozytorium dane.RepozytoriumKanalow) *adapterDesignu {

	a.kanaly, a.repozytoriumKanalow = kanaly, repozytorium
	return a
}

func zlozPolecenieObrazu(p shared.DesignPrompt) string {
	var b strings.Builder
	b.WriteString(p.Subject)
	dopiszCeche(&b, "styl", p.Style)
	dopiszCeche(&b, "kompozycja", p.Composition)
	dopiszCeche(&b, "oświetlenie", p.Lighting)
	dopiszCeche(&b, "paleta", p.Palette)
	dopiszCeche(&b, "proporcje kadru", p.AspectRatio)
	dopiszCeche(&b, "bez", p.Exclusions)
	return b.String()
}

func dopiszCeche(b *strings.Builder, nazwa string, wartosc *string) {
	if wartosc == nil || strings.TrimSpace(*wartosc) == "" {
		return
	}
	b.WriteString("; ")
	b.WriteString(nazwa)
	b.WriteString(": ")
	b.WriteString(strings.TrimSpace(*wartosc))
}

// Total niesie liczbę spełniających warunki, nie długość strony; PromptId zostaje puste.
func (a *adapterDesignu) Zasoby(ctx context.Context,
	z shared.DesignAssetListRequest) (shared.DesignAssetListResponse, error) {

	filtr := dane.FiltrZasobow{
		Okno:          z.WindowId,
		Etykiety:      z.Tags,
		TylkoUlubione: z.FavoriteOnly != nil && *z.FavoriteOnly,
		Limit:         wartoscLiczby(z.Limit),
	}
	if z.Kind != nil {
		rodzaj := string(*z.Kind)
		filtr.Rodzaj = &rodzaj
	}

	wiersze, razem, err := a.repozytorium.Zasoby(ctx, filtr)
	if err != nil {
		return shared.DesignAssetListResponse{}, bladDesignu(err)
	}
	zasoby := make([]shared.DesignAsset, 0, len(wiersze))
	for _, wiersz := range wiersze {
		etykiety, err := a.repozytorium.EtykietyZasobu(ctx, wiersz.ID)
		if err != nil {
			return shared.DesignAssetListResponse{}, bladDesignu(err)
		}
		zasoby = append(zasoby, zasobKontraktu(wiersz, etykiety))
	}
	return shared.DesignAssetListResponse{Assets: zasoby, Total: &razem}, nil
}

// Uri wychodzi jako ścieżka względna magazynu, nie ścieżka na dysku.
func zasobKontraktu(z dane.ZasobDesignu, etykiety []string) shared.DesignAsset {
	var odwolanie *string
	if z.URI != nil {
		if wzgledne := odwolanieZasobuDesignu(*z.URI); wzgledne != "" {
			odwolanie = &wzgledne
		}
	}
	return shared.DesignAsset{
		Id: z.Kod, WindowId: z.Okno, Name: z.Nazwa, Kind: shared.DesignAssetKind(z.Rodzaj),
		Format: z.Format, Uri: odwolanie, PromptId: z.PromptKod,
		VariantOfAssetId: z.WariantZasobuID,
		Tags:             etykiety, Favorite: &z.Ulubiony, Width: z.Szerokosc, Height: z.Wysokosc,
		CreatedAt: chwilaBazy(z.Utworzono),
	}
}

func sprawdzWyliczenieDesignu[T ~string](komenda, pole string, wartosc T, dopuszczalne []T) error {
	for _, znana := range dopuszczalne {
		if wartosc == znana {
			return nil
		}
	}
	nazwy := make([]string, 0, len(dopuszczalne))
	for _, znana := range dopuszczalne {
		nazwy = append(nazwy, string(znana))
	}
	return bladWskazaniaDesignu(fmt.Sprintf(
		"komenda %s z wartością %q pola %s, której kontrakt nie zna; wartości dopuszczalne: %s",
		komenda, string(wartosc), pole, strings.Join(nazwy, ", ")))
}

func bladDesignu(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, dane.ErrKolizjaWiersza) {
		return protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeConflict, err))
	}
	return protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeInternalError, err))
}

func bladWskazaniaDesignu(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		fmt.Sprintf("moduł Design: %s", powod)))
}

func bladNieznanegoBytuDesignu(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound, "moduł Design: "+powod))
}

func sprawdzRodzajZasobu(komenda string, rodzaj shared.DesignAssetKind) error {
	dopuszczalne := shared.WartosciDesignAssetKind()
	for _, znany := range dopuszczalne {
		if rodzaj == znany {
			return nil
		}
	}
	nazwy := make([]string, 0, len(dopuszczalne))
	for _, znany := range dopuszczalne {
		nazwy = append(nazwy, string(znany))
	}
	return bladWskazaniaDesignu(fmt.Sprintf(
		"komenda %s z rodzajem zasobu %q, którego kontrakt nie zna; rodzaje dopuszczalne: %s",
		komenda, string(rodzaj), strings.Join(nazwy, ", ")))
}
