// Plik doprowadza tożsamość modelu do zapytania kanału, gdzie domyka się
// łańcuch od katalogu przez składacz aż po przełącznik systemowego polecenia
// procesu.
package core

import (
	"context"
	"strings"

	"danacoconsole/server/internal/models"
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// ZTozsamoscia wpina port zasad i tożsamości modelu. Bez niego zapytanie kanału
// idzie z nakładką pustą, a kanał rusza z samą powłoką programu.
func (a *adapterRozmowy) ZTozsamoscia(t Tozsamosc) *adapterRozmowy {
	a.tozsamosc = t
	return a
}

// ZAgentami wpina źródło tożsamości ekspertów. Bez niego okno ze wskazanym
// ekspertem pracuje na samej osi modelu. To jedyne wejście tożsamości eksperta
// do tury.
func (a *adapterRozmowy) ZAgentami(z ZrodloTozsamosciAgenta) *adapterRozmowy {
	a.agenci = z
	return a
}

// nakladkaOkna liczy nakładkę obowiązującą okna rozmowy: oś modelu i oś konta
// bierze wskazanie okna, a błąd odczytu daje nakładkę pustą, nie zerwanie
// tury.
func (a *adapterRozmowy) nakladkaOkna(ctx context.Context, okno session.Okno) models.Nakladka {
	if a == nil || okno.Id == "" {
		return models.Nakladka{}
	}
	nakladka := models.Nakladka{}
	if a.tozsamosc != nil {
		idOkna := okno.Id
		wynik, err := a.tozsamosc.Obowiazujaca(ctx,
			shared.IdentityEffectiveGetRequest{WindowId: &idOkna})
		if err == nil {
			nakladka = nakladkaZapytania(wynik)
		}
	}
	return nakladkaZAgentem(nakladka, a.tozsamoscAgenta(ctx, okno))
}

// tozsamoscAgenta zwraca tożsamość eksperta wskazanego przez okno rozmowy,
// potrzebną do wyliczenia nakładki.
func (a *adapterRozmowy) tozsamoscAgenta(ctx context.Context, okno session.Okno) TozsamoscAgenta {
	if a == nil {
		return TozsamoscAgenta{}
	}
	return tozsamoscAgentaOkna(ctx, a.agenci, okno.Agent)
}

// uzupelnijAgenta nakłada na zapytanie nastawy procesu wnoszone przez eksperta:
// model, ustawienia i mosty MCP, w tej samej funkcji i kolejności co nakładka
// promptu.
func (a *adapterRozmowy) uzupelnijAgenta(ctx context.Context, okno session.Okno,
	zapytanie *models.Zapytanie) {

	nastawyAgenta(zapytanie, a.tozsamoscAgenta(ctx, okno))
}

// nakladkaZapytania przekłada nakładkę obowiązującą na trzy warstwy zapytania
// kanału, tak aby prompt złożony tutaj był tym samym ciągiem bajtów co prompt
// składacza.
func nakladkaZapytania(wynik shared.IdentityEffectiveGetResponse) models.Nakladka {
	warstwy := map[shared.IdentityLayer][]string{}
	for _, warstwa := range wynik.Layers {
		if strings.TrimSpace(warstwa.Content) == "" {
			continue
		}
		warstwy[warstwa.Layer] = append(warstwy[warstwa.Layer], warstwa.Content)
	}
	return models.Nakladka{
		Konstytucja: strings.Join(warstwy[shared.IdentityLayerConstitution], spoinaWarstw),
		ProfilRoli:  strings.Join(warstwy[shared.IdentityLayerProfile], spoinaWarstw),
		Ekspertyza:  strings.Join(warstwy[shared.IdentityLayerExpertise], spoinaWarstw),
		Tryb:        NakladkaSilnika(wynik).Tryb,
	}
}
