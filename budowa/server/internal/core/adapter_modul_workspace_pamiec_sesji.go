// Odpowiedzialność pliku: memory.toggle — poziomy pamięci włączone dla karty
// sesji i zgoda na zapis pamięci, czyli tabela konfiguracja_pamieci_sesji.
package core

import (
	"context"
	"errors"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// poziomyPamieciKontraktu przekłada wyliczenie kontraktu na wartości kolumny
// konfiguracja_pamieci_sesji.poziomy_wlaczone, wiążąc dwa istniejące słowniki
// zamiast zakładać trzeci spis poziomów.
var poziomyPamieciKontraktu = map[shared.ConfigScope]dane.PoziomPamieci{
	shared.ConfigScopeGlobal:      "globalna",
	shared.ConfigScopeEnvironment: "srodowisko",
	shared.ConfigScopeModule:      "modul",
	shared.ConfigScopeProject:     "projekt",
	shared.ConfigScopeSession:     "sesja",
}

// kolejnoscPoziomowPamieci trzyma poziomy od najszerszego do najwęższego.
// Wynik komendy wychodzi w tym porządku niezależnie od kolejności żądania.
var kolejnoscPoziomowPamieci = []shared.ConfigScope{
	shared.ConfigScopeGlobal,
	shared.ConfigScopeEnvironment,
	shared.ConfigScopeModule,
	shared.ConfigScopeProject,
	shared.ConfigScopeSession,
}

// PrzestawPamiecSesji przestawia poziomy pamięci włączone dla karty sesji
// oraz zgodę na zapis pamięci. Pole niewskazane zostawia stan bez zmian.
func (a *adapterPamieciPrzestrzeni) PrzestawPamiecSesji(ctx context.Context,
	z shared.MemoryToggleRequest) (shared.MemoryToggleResponse, error) {

	if z.SessionId == "" {
		return shared.MemoryToggleResponse{}, bladProjektu("przestawienie pamięci bez wskazania karty sesji")
	}
	if a.pamiec == nil || a.sesje == nil {
		return shared.MemoryToggleResponse{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeInternalError,
			"moduł Workspace: pamięć karty sesji niewpięta, przestawienia nie ma gdzie zapisać"))
	}
	sesja, err := a.sesje.PoIdentyfikatorze(ctx, z.SessionId)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return shared.MemoryToggleResponse{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeNotFound, "moduł Workspace: karta sesji "+z.SessionId+" nie istnieje"))
	}
	if err != nil {
		return shared.MemoryToggleResponse{}, err
	}
	konfiguracja, err := a.konfiguracjaPamieciSesji(ctx, sesja.ID)
	if err != nil {
		return shared.MemoryToggleResponse{}, err
	}
	if len(z.Levels) > 0 {
		poziomy, err := poziomyPamieciZadania(z.Levels)
		if err != nil {
			return shared.MemoryToggleResponse{}, err
		}
		konfiguracja.PoziomyWlaczone = poziomy
	}
	if z.WriteEnabled != nil {
		konfiguracja.ZapisWlaczony = *z.WriteEnabled
	}
	if err := a.pamiec.UstawKonfiguracjeSesji(ctx, konfiguracja); err != nil {
		return shared.MemoryToggleResponse{}, err
	}
	poziomy, err := poziomyPamieciWyniku(konfiguracja.PoziomyWlaczone)
	if err != nil {
		return shared.MemoryToggleResponse{}, err
	}
	return shared.MemoryToggleResponse{Levels: poziomy, WriteEnabled: konfiguracja.ZapisWlaczony}, nil
}

// konfiguracjaPamieciSesji zwraca stan zastany albo stan domyślny karty, która
// pamięci jeszcze nie przestawiała.
func (a *adapterPamieciPrzestrzeni) konfiguracjaPamieciSesji(ctx context.Context,
	sesjaID int64) (dane.KonfiguracjaPamieci, error) {

	konfiguracja, jest, err := a.pamiec.KonfiguracjaSesji(ctx, sesjaID)
	if err != nil {
		return dane.KonfiguracjaPamieci{}, err
	}
	if jest {
		return konfiguracja, nil
	}
	// Stan domyślny powtarza wartość domyślną kolumny, wszystkie poziomy i zapis wolny.
	domyslne := make([]dane.PoziomPamieci, 0, len(kolejnoscPoziomowPamieci))
	for _, poziom := range kolejnoscPoziomowPamieci {
		domyslne = append(domyslne, poziomyPamieciKontraktu[poziom])
	}
	return dane.KonfiguracjaPamieci{
		SesjaID:         sesjaID,
		PoziomyWlaczone: domyslne,
		ZapisWlaczony:   true,
	}, nil
}

// poziomyPamieciZadania przekłada poziomy żądania na wartości kolumny. Poziom
// zasięgu, którego obszar pamięci nie zna, kończy komendę odmową.
func poziomyPamieciZadania(zadane []shared.ConfigScope) ([]dane.PoziomPamieci, error) {
	wybrane := map[shared.ConfigScope]bool{}
	for _, poziom := range zadane {
		if _, znany := poziomyPamieciKontraktu[poziom]; !znany {
			return nil, bladProjektu("poziom zasięgu " + string(poziom) + " nie jest poziomem pamięci")
		}
		wybrane[poziom] = true
	}
	poziomy := make([]dane.PoziomPamieci, 0, len(wybrane))
	for _, poziom := range kolejnoscPoziomowPamieci {
		if wybrane[poziom] {
			poziomy = append(poziomy, poziomyPamieciKontraktu[poziom])
		}
	}
	return poziomy, nil
}

// poziomyPamieciWyniku przekłada wartości kolumny z powrotem na wyliczenie
// kontraktu. Wartość, której przekład nie zna, kończy komendę odmową:
// pominięcie jej po cichu oddałoby wykaz mówiący, że poziom jest wyłączony,
// podczas gdy w bazie stoi włączony.
func poziomyPamieciWyniku(zapisane []dane.PoziomPamieci) ([]shared.ConfigScope, error) {
	wykaz := map[dane.PoziomPamieci]bool{}
	for _, poziom := range zapisane {
		wykaz[poziom] = true
	}
	poziomy := make([]shared.ConfigScope, 0, len(zapisane))
	for _, poziom := range kolejnoscPoziomowPamieci {
		kolumna := poziomyPamieciKontraktu[poziom]
		if wykaz[kolumna] {
			poziomy = append(poziomy, poziom)
			delete(wykaz, kolumna)
		}
	}
	for nieznany := range wykaz {
		return nil, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"moduł Workspace: pamięć karty sesji niesie nieznany poziom "+string(nieznany)))
	}
	return poziomy, nil
}
