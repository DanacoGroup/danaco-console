// Odpowiedzialność pliku: wypełnienie portu Agenci biblioteką ekspertów
// z bazy — założenie, wykaz, zmiana tożsamości i usunięcie. Ekspert nie ma
// własnego rejestru modeli.
package core

import (
	"context"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/models"
	"danacoconsole/shared"
)

// Zgodność struktury adapterAgentow z portem Agenci jest sprawdzana już przy
// samej kompilacji pakietu core.
var _ Agenci = (*adapterAgentow)(nil)

// adapterAgentow wypełnia port Agenci tabelami rodziny `agent*`, obsługując
// bibliotekę ekspertów rdzenia.
type adapterAgentow struct {
	repozytorium dane.RepozytoriumAgentow
	// kanaly jest rejestrem kanałów modelu rdzenia, do sprawdzenia wskazania
	// kanału.
	kanaly *models.Rejestr
	// punkty są katalogiem punktów dostępu; konektor `mcp` wskazuje most
	// z tego katalogu.
	punkty dane.RepozytoriumPunktowDostepu
	// warstwy jest repozytorium tożsamości własnej eksperta, wpinanym metodą
	// `ZWarstwami`.
	warstwy dane.RepozytoriumWarstwAgenta
}

// nowyAdapterAgentow wiąże port Agenci z repozytorium biblioteki ekspertów,
// bez wpięcia kanałów, punktów dostępu ani warstw tożsamości.
func nowyAdapterAgentow(repozytorium dane.RepozytoriumAgentow) *adapterAgentow {
	return &adapterAgentow{repozytorium: repozytorium}
}

// ZKanalami wpina rejestr kanałów modelu. Bez niego wskazanie kanału przechodzi
// bez sprawdzenia — biblioteka ekspertów pracuje dalej.
func (a *adapterAgentow) ZKanalami(kanaly *models.Rejestr) *adapterAgentow {
	a.kanaly = kanaly
	return a
}

// ZPunktamiDostepu wpina katalog mostów obsługujący konektory rodzaju
// `mcp` w bibliotece ekspertów rdzenia.
func (a *adapterAgentow) ZPunktamiDostepu(punkty dane.RepozytoriumPunktowDostepu) *adapterAgentow {
	a.punkty = punkty
	return a
}

// Utworz zakłada eksperta wraz z wyjściowym kompletem uprawnień i
// tożsamością zapisaną osobną czynnością.
func (a *adapterAgentow) Utworz(ctx context.Context,
	z shared.AgentCreateRequest) (shared.AgentCreateResponse, error) {

	if a == nil || a.repozytorium == nil {
		return shared.AgentCreateResponse{}, bladBrakuKatalogu("ekspertów")
	}
	nazwa := strings.TrimSpace(z.Name)
	if nazwa == "" {
		return shared.AgentCreateResponse{}, bladZadaniaEksperta("ekspert wymaga nazwy")
	}
	kanal, err := a.kanalWskazany(z.ChannelId)
	if err != nil {
		return shared.AgentCreateResponse{}, err
	}
	if err := sprawdzWidocznosc(z.Visibility); err != nil {
		return shared.AgentCreateResponse{}, err
	}
	// Pominięte pole pamięci znaczy komplet: nil to brak deklaracji, lista
	// pusta wyłącza pamięć.
	poziomy := poziomyPamieciWyjsciowe
	if z.MemoryLevels != nil {
		poziomy, err = sprawdzPoziomyPamieci(z.MemoryLevels)
		if err != nil {
			return shared.AgentCreateResponse{}, err
		}
	}
	zapisany, err := a.repozytorium.Dodaj(ctx, dane.Agent{
		Kod:                 nowyIdentyfikator(przedrostekEksperta),
		Nazwa:               nazwa,
		Opis:                wartoscTekstu(z.Description),
		InstrukcjeSystemowe: wartoscTekstu(z.SystemPrompt),
		KanalKod:            kanal,
		Model:               a.modelWskazany(z.Model, kanal),
		Aktywny:             true,
		Widocznosc:          widocznoscZadania(z.Visibility),
		PoziomyPamieci:      poziomy,
	})
	if err != nil {
		return shared.AgentCreateResponse{}, err
	}
	// Imię własne i favikon idą tą samą drogą co przy zmianie tożsamości
	// eksperta.
	if err := a.zapiszTozsamosc(ctx, zapisany.Kod, z.DisplayName, z.Favicon); err != nil {
		return shared.AgentCreateResponse{}, err
	}
	if err := a.zapiszTrybNakladki(ctx, zapisany.Kod, z.Mode); err != nil {
		return shared.AgentCreateResponse{}, err
	}
	pelny, err := a.EkspertPelny(ctx, zapisany.Kod)
	if err != nil {
		return shared.AgentCreateResponse{}, err
	}
	return shared.AgentCreateResponse{Agent: pelny}, nil
}

// Zmien zapisuje zmienione pola tożsamości eksperta, poziomy pamięci i tryb
// nakładki, podnosząc licznik wersji.
func (a *adapterAgentow) Zmien(ctx context.Context,
	z shared.AgentUpdateRequest) (shared.AgentUpdateResponse, error) {

	if a == nil || a.repozytorium == nil {
		return shared.AgentUpdateResponse{}, bladBrakuKatalogu("ekspertów")
	}
	biezacy, err := a.repozytorium.PoKodzie(ctx, z.AgentId)
	if err != nil {
		return shared.AgentUpdateResponse{}, bladWskazania(err, "ekspert", z.AgentId)
	}
	if z.Name != nil {
		nazwa := strings.TrimSpace(*z.Name)
		if nazwa == "" {
			return shared.AgentUpdateResponse{}, bladZadaniaEksperta("nazwa eksperta nie może być pusta")
		}
		biezacy.Nazwa = nazwa
	}
	if z.Description != nil {
		biezacy.Opis = *z.Description
	}
	if z.SystemPrompt != nil {
		biezacy.InstrukcjeSystemowe = *z.SystemPrompt
	}
	if z.Enabled != nil {
		biezacy.Aktywny = *z.Enabled
	}
	if err := sprawdzWidocznosc(z.Visibility); err != nil {
		return shared.AgentUpdateResponse{}, err
	}
	if z.Visibility != nil {
		biezacy.Widocznosc = string(*z.Visibility)
	}
	zapisany, err := a.repozytorium.Aktualizuj(ctx, biezacy)
	if err != nil {
		return shared.AgentUpdateResponse{}, err
	}
	// Imię własne i favikon idą osobną drogą, bo zapytanie aktualizujące
	// eksperta nie zna tych kolumn.
	if err := a.zapiszTozsamosc(ctx, zapisany.Kod, z.DisplayName, z.Favicon); err != nil {
		return shared.AgentUpdateResponse{}, err
	}
	if err := a.zapiszTrybNakladki(ctx, zapisany.Kod, z.Mode); err != nil {
		return shared.AgentUpdateResponse{}, err
	}
	// Poziomy pamięci idą osobną drogą, bo siedzą w tabeli podrzędnej; nil
	// zostawia zastane poziomy.
	if z.MemoryLevels != nil {
		poziomy, err := sprawdzPoziomyPamieci(z.MemoryLevels)
		if err != nil {
			return shared.AgentUpdateResponse{}, err
		}
		if _, err := a.repozytorium.UstawPoziomyPamieci(ctx, zapisany.Kod, poziomy); err != nil {
			return shared.AgentUpdateResponse{}, bladWskazania(err, "ekspert", zapisany.Kod)
		}
	}
	pelny, err := a.EkspertPelny(ctx, zapisany.Kod)
	if err != nil {
		return shared.AgentUpdateResponse{}, err
	}
	return shared.AgentUpdateResponse{Agent: pelny}, nil
}

// Wykaz zwraca bibliotekę ekspertów. Biblioteka pusta nie jest odmową — okno
// pokazuje wtedy stan pusty i zachętę do założenia pierwszego eksperta.
func (a *adapterAgentow) Wykaz(ctx context.Context,
	z shared.AgentListRequest) (shared.AgentListResponse, error) {

	if a == nil || a.repozytorium == nil {
		return shared.AgentListResponse{Agents: []shared.Agent{}}, nil
	}
	// Wskazany projekt zawęża wykaz do ekspertów w nim widocznych.
	filtr := dane.FiltrAgentow{
		Fraza:       wartoscTekstu(z.Query),
		KodProjektu: strings.TrimSpace(wartoscTekstu(z.ProjectId)),
	}
	if z.EnabledOnly != nil {
		filtr.TylkoAktywne = *z.EnabledOnly
	}
	if z.Limit != nil {
		filtr.Granica = *z.Limit
	}
	wiersze, razem, err := a.repozytorium.Lista(ctx, filtr)
	if err != nil {
		return shared.AgentListResponse{}, err
	}
	// Warstwy i wtyczki dochodzą dwoma zapytaniami na cały wykaz, nie dwoma
	// na eksperta.
	warstwy, wtyczki := a.tozsamosciWykazu(ctx)
	eksperci := make([]shared.Agent, 0, len(wiersze))
	for _, wiersz := range wiersze {
		ekspert := ekspertKontraktu(wiersz)
		ekspert.Layers = warstwyKontraktu(warstwy[wiersz.Kod])
		ekspert.PluginIds = kodyWtyczek(wtyczki[wiersz.Kod])
		eksperci = append(eksperci, ekspert)
	}
	liczba := razem
	return shared.AgentListResponse{Agents: eksperci, Total: &liczba}, nil
}

// widocznoscZadania czyta zasięg widoczności z żądania. Pominięty znaczy
// `global` — stan wyjściowy platformy jest szeroki, a zawężenie do projektu
// wymaga wskazania.
func widocznoscZadania(widocznosc *shared.AgentVisibility) string {
	if widocznosc == nil {
		return shared.AgentVisibilityGlobal
	}
	return string(*widocznosc)
}

// Usun kasuje eksperta wraz z powiązaniami. Brak wiersza nie jest odmową:
// wynik `deleted=false` mówi, że nie było czego usuwać.
func (a *adapterAgentow) Usun(ctx context.Context,
	z shared.AgentDeleteRequest) (shared.AgentDeleteResponse, error) {

	if a == nil || a.repozytorium == nil {
		return shared.AgentDeleteResponse{}, bladBrakuKatalogu("ekspertów")
	}
	usuniety, err := a.repozytorium.Usun(ctx, z.AgentId)
	if err != nil {
		return shared.AgentDeleteResponse{}, err
	}
	return shared.AgentDeleteResponse{AgentId: z.AgentId, Deleted: usuniety}, nil
}

// Pobierz oddaje jednego eksperta rozgłoszeniu zmiany biblioteki, w postaci
// pełnego bytu kontraktu do wysłania.
func (a *adapterAgentow) Pobierz(ctx context.Context, idEksperta string) (shared.Agent, error) {
	if a == nil || a.repozytorium == nil {
		return shared.Agent{}, bladBrakuKatalogu("ekspertów")
	}
	wiersz, err := a.repozytorium.PoKodzie(ctx, idEksperta)
	if err != nil {
		return shared.Agent{}, bladWskazania(err, "ekspert", idEksperta)
	}
	return ekspertKontraktu(wiersz), nil
}
