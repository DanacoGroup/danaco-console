// Doprowadzenie tożsamości modelu do zapytania kanału.
//
// Tutaj domyka się łańcuch katalog → składacz → zapytanie kanału → silnik
// nakładki → przełącznik `--system-prompt` albo `--append-system-prompt`.
// Rdzeń nie składa promptu samodzielnie — bierze wynik składacza.
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

// nakladkaOkna liczy nakładkę obowiązującą okna rozmowy. Oś modelu i oś konta
// bierze wskazanie okna po stronie adaptera tożsamości — okno zna kanał, kanał
// zna model i konto, więc drugiego wyliczania osi tutaj nie ma.
//
// Po wyliczeniu osi nakładka dostaje warstwy eksperta wskazanego przez okno.
// To jedyne miejsce, w którym tożsamość eksperta wchodzi do promptu — dalej
// jedzie trasą `models.Nakladka` → `nakladkaKanaluGlownego` →
// `injection.Nakladka` → przełącznik CLI.
//
// Okno przychodzi w całości, nie samym identyfikatorem: kod eksperta jest
// nastawą okna, więc pytanie o niego rejestru drugi raz byłoby powtórzeniem
// odczytu, który wywołujący już wykonał.
//
// Błąd odczytu daje nakładkę pustą, nie zerwanie tury: brak zasad jest brakiem
// treści systemowej, a nie przeszkodą w rozmowie.
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

// tozsamoscAgenta zwraca tożsamość eksperta wskazanego przez okno.
func (a *adapterRozmowy) tozsamoscAgenta(ctx context.Context, okno session.Okno) TozsamoscAgenta {
	if a == nil {
		return TozsamoscAgenta{}
	}
	return tozsamoscAgentaOkna(ctx, a.agenci, okno.Agent)
}

// uzupelnijAgenta nakłada na zapytanie nastawy procesu wnoszone przez eksperta:
// model, ustawienia i mosty MCP.
//
// Osobno od nakładki promptu, bo to inne pola zapytania — ale w tej samej
// funkcji składającej turę i w tej samej kolejności co prowenancja, żeby
// podgląd pokazywał wiersz, który tura naprawdę wykona.
func (a *adapterRozmowy) uzupelnijAgenta(ctx context.Context, okno session.Okno,
	zapytanie *models.Zapytanie) {

	nastawyAgenta(zapytanie, a.tozsamoscAgenta(ctx, okno))
}

// nakladkaZapytania przekłada nakładkę obowiązującą na trzy warstwy zapytania
// kanału. Kategorie jednej krytyczności skleja spoiną warstw, żeby prompt
// złożony tutaj był tym samym ciągiem bajtów co prompt składacza.
//
// Tryb pochodzi z jedynego przekładu trybu w rdzeniu — NakladkaSilnika — więc
// rozstrzygnięcie, czy tożsamość zamienia treść systemową, czy się do niej
// dokłada, dojeżdża do procesu modelu tą samą drogą, którą pokazuje je okno
// konfiguracji.
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
