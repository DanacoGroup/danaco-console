// Odpowiedzialność pliku: zasoby eksperta — umiejętności (Skills Manager),
// konektory (Connectors Manager) i uprawnienia (Permissions Center).
package core

import (
	"context"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// DodajUmiejetnosc przypisuje ekspertowi umiejętność, obsługując komendę agent.skill.add kontraktu Agents.
func (a *adapterAgentow) DodajUmiejetnosc(ctx context.Context,
	z shared.AgentSkillAddRequest) (shared.AgentSkillAddResponse, error) {

	if a == nil || a.repozytorium == nil {
		return shared.AgentSkillAddResponse{}, bladBrakuKatalogu("ekspertów")
	}
	umiejetnosc := strings.TrimSpace(z.SkillId)
	if umiejetnosc == "" {
		return shared.AgentSkillAddResponse{}, bladZadaniaEksperta("przypisanie wymaga wskazania umiejętności")
	}
	zapisany, err := a.repozytorium.DodajUmiejetnosc(ctx, z.AgentId, umiejetnosc)
	if err != nil {
		return shared.AgentSkillAddResponse{}, bladWskazania(err, "ekspert", z.AgentId)
	}
	return shared.AgentSkillAddResponse{Agent: ekspertKontraktu(zapisany)}, nil
}

// DodajKonektor podłącza ekspertowi konektor, wtyczkę albo serwer MCP wskazany żądaniem komendy kontraktu.
func (a *adapterAgentow) DodajKonektor(ctx context.Context,
	z shared.AgentConnectorAddRequest) (shared.AgentConnectorAddResponse, error) {

	if a == nil || a.repozytorium == nil {
		return shared.AgentConnectorAddResponse{}, bladBrakuKatalogu("ekspertów")
	}
	nazwa := strings.TrimSpace(z.Name)
	if nazwa == "" {
		return shared.AgentConnectorAddResponse{}, bladZadaniaEksperta("konektor wymaga nazwy")
	}
	if _, jest := rodzajeKonektora[z.Kind]; !jest {
		return shared.AgentConnectorAddResponse{},
			bladZadaniaEksperta("rodzaj konektora " + string(z.Kind) + " nie należy do kontraktu")
	}
	konfiguracja, err := sprawdzParametry(z.Config, "konfiguracja integracji")
	if err != nil {
		return shared.AgentConnectorAddResponse{}, err
	}
	kodPunktu := strings.TrimSpace(wartoscTekstu(z.AccessPointId))
	numerPunktu, err := a.punktMostu(ctx, z.Kind, kodPunktu)
	if err != nil {
		return shared.AgentConnectorAddResponse{}, err
	}
	zapisany, err := a.repozytorium.DodajKonektor(ctx, dane.KonektorAgenta{
		Kod: nowyIdentyfikator(przedrostekKonektora), AgentKod: z.AgentId, Nazwa: nazwa,
		Rodzaj: string(z.Kind), PunktDostepuID: numerPunktu, Konfiguracja: konfiguracja,
		Aktywny: true,
	})
	if err != nil {
		return shared.AgentConnectorAddResponse{}, bladWskazania(err, "ekspert", z.AgentId)
	}
	return shared.AgentConnectorAddResponse{Connector: konektorKontraktu(zapisany, kodPunktu)}, nil
}

// UstawUprawnienie przyznaje albo odbiera uprawnienie eksperta w jednej z czterech grup zakresu kontraktu.
func (a *adapterAgentow) UstawUprawnienie(ctx context.Context,
	z shared.AgentPermissionSetRequest) (shared.AgentPermissionSetResponse, error) {

	if a == nil || a.repozytorium == nil {
		return shared.AgentPermissionSetResponse{}, bladBrakuKatalogu("ekspertów")
	}
	if _, jest := grupyUprawnien[z.Group]; !jest {
		return shared.AgentPermissionSetResponse{},
			bladZadaniaEksperta("grupa zakresu " + string(z.Group) + " nie należy do kontraktu")
	}
	zapisany, err := a.repozytorium.UstawUprawnienie(ctx, z.AgentId, dane.UprawnienieAgenta{
		Grupa: string(z.Group), Zakres: wartoscTekstu(z.Scope), Przyznane: z.Granted,
	})
	if err != nil {
		return shared.AgentPermissionSetResponse{}, bladWskazania(err, "ekspert", z.AgentId)
	}
	return shared.AgentPermissionSetResponse{Permissions: uprawnieniaKontraktu(zapisany.Uprawnienia)}, nil
}

// punktMostu przekłada kod punktu dostępu na numer wiersza. Konektor rodzaju
// mcp bez wskazania mostu jest odrzucany, bo adres serwera nie ma gdzie
// zamieszkać poza katalogiem punktów dostępu.
func (a *adapterAgentow) punktMostu(ctx context.Context,
	rodzaj shared.AgentConnectorKind, kod string) (*int64, error) {

	if kod == "" {
		if rodzaj == shared.AgentConnectorKindMcp {
			return nil, bladZadaniaEksperta(
				"konektor MCP wymaga wskazania mostu z katalogu punktów dostępu")
		}
		return nil, nil
	}
	if a.punkty == nil {
		return nil, bladBrakuKatalogu("punktów dostępu")
	}
	punkt, err := a.punkty.PoKodzie(ctx, kod)
	if err != nil {
		return nil, bladWskazania(err, "punkt dostępu", kod)
	}
	if rodzaj == shared.AgentConnectorKindMcp && punkt.Rodzaj != shared.AccessPointKindMcpBridge {
		return nil, bladZadaniaEksperta("punkt dostępu " + kod + " nie jest mostem MCP")
	}
	numer := punkt.ID
	return &numer, nil
}
