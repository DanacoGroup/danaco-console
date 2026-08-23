// Odpowiedzialność pliku: składanie tożsamości modelu z danych katalogu.
//
// Tożsamość modelu jest zamieniana, nie dołączana do ustawień fabrycznych.
// Silnik nakładki już to umie —
// `internal/injection` niesie trzy warstwy wg krytyczności oraz dwa
// tryby podania: TrybZastap (`--system-prompt`) i TrybDopisz
// (`--append-system-prompt`). Ten moduł nie powtarza silnika — daje mu
// sterowanie z danych:
// bierze kategorie z katalogu, treść z osi platformy, modelu i konta, układa
// warstwy i oddaje prompt wraz z trybem.
//
// Rdzeń nie zna ani jednego zdania promptu. Zna wyłącznie porządek
// składania i regułę wyboru osi.
package core

import (
	"context"
	"strings"

	"danacoconsole/shared"
)

// KatalogTozsamosci jest portem danych tożsamości modelu. Mówi
// wyłącznie typami kontraktu — warstwa trwałości leży po drugiej stronie
// adaptera (`tozsamosc_zrodlo.go`).
type KatalogTozsamosci interface {
	// Kategorie zwraca katalog kategorii zasad; kolejność wewnątrz warstwy
	// wnosi katalog, kolejność warstw nakłada składacz.
	Kategorie(ctx context.Context, tylkoAktywne bool) ([]shared.IdentityCategory, error)
	// Dokumenty zwraca treści zapisane dla jednej osi i jednego jej bytu.
	Dokumenty(ctx context.Context, os shared.ConfigAxis, bytOsi string) ([]shared.IdentityDocument, error)
}

// ZapytanieTozsamosci opisuje, DLA CZEGO liczona jest nakładka: dla jakiego
// modelu i jakiego konta. Puste wskazanie znaczy „ta oś nie obowiązuje”, nie
// „brak wyniku” — zostaje wtedy sama oś platformy.
type ZapytanieTozsamosci struct {
	// Model jest bytem osi `model` — identyfikatorem modelu okna.
	Model string
	// Konto jest bytem osi `account` — identyfikatorem konta z rejestru kont.
	Konto string
	// TrybDomyslny pochodzi z klucza `tozsamosc.tryb_domyslny`. Pusty znaczy
	// ZASTAP.
	TrybDomyslny shared.IdentityMode
}

// SkladaczTozsamosci buduje nakładkę obowiązującą z wierszy katalogu i treści.
// Ta sama konfiguracja daje bajtowo ten sam prompt — porządek jest w całości
// wyznaczony danymi, nigdy kolejnością odczytu z mapy. Bez tej własności
// pamięć podręczna promptu po stronie kanału byłaby bezużyteczna.
type SkladaczTozsamosci struct {
	katalog KatalogTozsamosci
}

// NowySkladaczTozsamosci wiąże składacz z katalogiem.
func NowySkladaczTozsamosci(katalog KatalogTozsamosci) *SkladaczTozsamosci {
	return &SkladaczTozsamosci{katalog: katalog}
}

// Zloz zwraca nakładkę obowiązującą: tryb, warstwy w kolejności krytyczności,
// złożony prompt, jego odcisk oraz wykaz kategorii obowiązkowych bez treści.
//
// Brak katalogu, brak treści i brak kategorii obowiązkowej nie wstrzymują
// niczego: wynikiem jest pusty prompt, a kanał rusza wtedy z samą powłoką.
// Wykaz braków jest informacją dla Operatora, nie bramą.
func (s *SkladaczTozsamosci) Zloz(ctx context.Context, z ZapytanieTozsamosci) (shared.IdentityEffectiveGetResponse, error) {
	if s == nil || s.katalog == nil {
		return pustaTozsamosc(z.TrybDomyslny), nil
	}
	kategorie, err := s.katalog.Kategorie(ctx, true)
	if err != nil {
		return shared.IdentityEffectiveGetResponse{}, err
	}
	tresci, err := s.tresciOsi(ctx, z)
	if err != nil {
		return shared.IdentityEffectiveGetResponse{}, err
	}
	return zlozNakladke(kategorie, tresci, z.TrybDomyslny), nil
}

// tresciOsi zbiera treści trzech osi w kolejności od najszerszej. Oś bez
// wskazanego bytu jest pomijana — nie ma czego rozstrzygać.
func (s *SkladaczTozsamosci) tresciOsi(ctx context.Context, z ZapytanieTozsamosci) ([]shared.IdentityDocument, error) {
	osie := []struct {
		os  shared.ConfigAxis
		byt string
	}{
		{shared.ConfigAxisPlatform, ""},
		{shared.ConfigAxisModel, strings.TrimSpace(z.Model)},
		{shared.ConfigAxisAccount, strings.TrimSpace(z.Konto)},
	}
	zebrane := []shared.IdentityDocument{}
	for _, os := range osie {
		if os.os != shared.ConfigAxisPlatform && os.byt == "" {
			continue
		}
		dokumenty, err := s.katalog.Dokumenty(ctx, os.os, os.byt)
		if err != nil {
			return nil, err
		}
		zebrane = append(zebrane, dokumenty...)
	}
	return zebrane, nil
}

// pustaTozsamosc jest wynikiem dla konfiguracji, w której nie ma czego składać.
// Tryb wraca mimo pustego promptu, bo Operator ma widzieć, co obowiązywałoby po
// wpisaniu treści.
func pustaTozsamosc(trybDomyslny shared.IdentityMode) shared.IdentityEffectiveGetResponse {
	return shared.IdentityEffectiveGetResponse{
		Mode:   trybObowiazujacy(nil, trybDomyslny),
		Layers: []shared.IdentityLayerContent{},
	}
}
