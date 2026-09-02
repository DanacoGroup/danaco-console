// Skład debaty i jego odczyt: roundtable.debate.get, roundtable.model.*, roundtable.team.*
// i roundtable.role.list; okno otwarte w trakcie debaty poznaje stan sprzed otwarcia.
package core

import (
	"context"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

const (
	przedrostekZespoluDebaty = "zespol-"
)

func (a *adapterDebaty) StanDebaty(ctx context.Context,
	z shared.RoundtableDebateGetRequest) (shared.RoundtableDebateGetResponse, error) {

	okno := strings.TrimSpace(z.WindowId)
	if okno == "" {
		return shared.RoundtableDebateGetResponse{},
			bladWskazaniaDebaty("komenda bez wskazania okna debaty")
	}
	uczestnicy, err := a.repozytorium.Uczestnicy(ctx, okno)
	if err != nil {
		return shared.RoundtableDebateGetResponse{}, bladDebaty(err)
	}
	tury, err := a.turyOknaOdPierwszej(ctx, okno)
	if err != nil {
		return shared.RoundtableDebateGetResponse{}, err
	}

	turaKod := strings.TrimSpace(wartoscTekstu(z.TurnId))
	wypowiedzi, err := a.wypowiedziZakresu(ctx, okno, turaKod)
	if err != nil {
		return shared.RoundtableDebateGetResponse{}, err
	}
	lacznie := len(wypowiedzi)
	wycinek := wytnijWypowiedziDebaty(wypowiedzi, wartoscLiczby(z.Offset), wartoscLiczby(z.Limit))

	migawka := shared.RoundtableDebateSnapshot{
		WindowId:     okno,
		Participants: uczestnicyKontraktu(uczestnicy),
		Turns:        turyKontraktu(tury),
		Statements:   wypowiedziKontraktu(wycinek),
	}
	// Stanowisko wchodzi do migawki tylko wtedy, gdy naprawdę leży w bazie.
	if stanowisko, err := a.repozytorium.Stanowisko(ctx, okno, ""); err == nil {
		wynik := stanowiskoKontraktu(stanowisko, nil)
		migawka.Consensus = &wynik
	}
	return shared.RoundtableDebateGetResponse{Snapshot: migawka, Total: &lacznie}, nil
}

func (a *adapterDebaty) wypowiedziZakresu(ctx context.Context,
	okno, turaKod string) ([]dane.WypowiedzDebaty, error) {

	if turaKod == "" {
		wypowiedzi, err := a.repozytorium.WypowiedziOkna(ctx, okno)
		if err != nil {
			return nil, bladDebaty(err)
		}
		return wypowiedzi, nil
	}
	tura, err := a.repozytorium.Tura(ctx, turaKod)
	if err != nil {
		return nil, bladNieznanejTury(turaKod, err)
	}
	if tura.Okno != okno {
		return nil, bladWskazaniaDebaty("tura " + turaKod + " nie należy do okna " + okno)
	}
	wypowiedzi, err := a.repozytorium.Wypowiedzi(ctx, turaKod)
	if err != nil {
		return nil, bladDebaty(err)
	}
	return wypowiedzi, nil
}

// Repozytorium oddaje tury od najnowszej (tak czyta Debate Panel); migawka czyta się od pierwszej.
func (a *adapterDebaty) turyOknaOdPierwszej(ctx context.Context,
	okno string) ([]dane.TuraDebaty, error) {

	tury, err := a.repozytorium.Tury(ctx, okno, 0)
	if err != nil {
		return nil, bladDebaty(err)
	}
	odwrocone := make([]dane.TuraDebaty, 0, len(tury))
	for i := len(tury) - 1; i >= 0; i-- {
		odwrocone = append(odwrocone, tury[i])
	}
	return odwrocone, nil
}

// Granica niedodatnia znaczy „bez granicy”, przesunięcie poza wykaz oddaje wykaz pusty.
func wytnijWypowiedziDebaty(wypowiedzi []dane.WypowiedzDebaty,
	przesuniecie, granica int) []dane.WypowiedzDebaty {

	if przesuniecie < 0 {
		przesuniecie = 0
	}
	if przesuniecie >= len(wypowiedzi) {
		return nil
	}
	wycinek := wypowiedzi[przesuniecie:]
	if granica > 0 && granica < len(wycinek) {
		wycinek = wycinek[:granica]
	}
	return wycinek
}

func (a *adapterDebaty) Uczestnicy(ctx context.Context,
	z shared.RoundtableModelListRequest) (shared.RoundtableModelListResponse, error) {

	okno := strings.TrimSpace(z.WindowId)
	if okno == "" {
		return shared.RoundtableModelListResponse{},
			bladWskazaniaDebaty("komenda bez wskazania okna debaty")
	}
	uczestnicy, err := a.repozytorium.Uczestnicy(ctx, okno)
	if err != nil {
		return shared.RoundtableModelListResponse{}, bladDebaty(err)
	}
	return shared.RoundtableModelListResponse{Participants: uczestnicyKontraktu(uczestnicy)}, nil
}

// Wypowiedzi usuniętego zostają w zapisie tury: transkrypt jest zapisem tego, co padło.
func (a *adapterDebaty) UsunModel(ctx context.Context,
	z shared.RoundtableModelRemoveRequest) (shared.RoundtableModelRemoveResponse, error) {

	okno := strings.TrimSpace(z.WindowId)
	kod := strings.TrimSpace(z.ParticipantId)
	if okno == "" {
		return shared.RoundtableModelRemoveResponse{},
			bladWskazaniaDebaty("komenda bez wskazania okna debaty")
	}
	if kod == "" {
		return shared.RoundtableModelRemoveResponse{},
			bladWskazaniaDebaty("usunięcie bez wskazania uczestnika")
	}
	uczestnik, err := a.repozytorium.Uczestnik(ctx, kod)
	if err != nil {
		return shared.RoundtableModelRemoveResponse{}, bladNieznanegoUczestnika(kod, err)
	}
	if uczestnik.Okno != okno {
		return shared.RoundtableModelRemoveResponse{},
			bladWskazaniaDebaty("uczestnik " + kod + " nie należy do okna " + okno)
	}
	if err := a.repozytorium.UsunUczestnika(ctx, kod); err != nil {
		return shared.RoundtableModelRemoveResponse{}, bladNieznanegoUczestnika(kod, err)
	}
	uczestnicy, err := a.repozytorium.Uczestnicy(ctx, okno)
	if err != nil {
		return shared.RoundtableModelRemoveResponse{}, bladDebaty(err)
	}
	return shared.RoundtableModelRemoveResponse{Participants: uczestnicyKontraktu(uczestnicy)}, nil
}

// Pole niepodane zostaje takie, jakie było; zmiana wagi nie czyści promptu systemowego.
func (a *adapterDebaty) ZmienModel(ctx context.Context,
	z shared.RoundtableModelUpdateRequest) (shared.RoundtableModelUpdateResponse, error) {

	okno := strings.TrimSpace(z.WindowId)
	kod := strings.TrimSpace(z.ParticipantId)
	if okno == "" {
		return shared.RoundtableModelUpdateResponse{},
			bladWskazaniaDebaty("komenda bez wskazania okna debaty")
	}
	if kod == "" {
		return shared.RoundtableModelUpdateResponse{},
			bladWskazaniaDebaty("zmiana bez wskazania uczestnika")
	}
	uczestnik, err := a.repozytorium.Uczestnik(ctx, kod)
	if err != nil {
		return shared.RoundtableModelUpdateResponse{}, bladNieznanegoUczestnika(kod, err)
	}
	if uczestnik.Okno != okno {
		return shared.RoundtableModelUpdateResponse{},
			bladWskazaniaDebaty("uczestnik " + kod + " nie należy do okna " + okno)
	}

	if z.PersonaName != nil {
		uczestnik.NazwaTozsamosci = wskaznikTekstu(strings.TrimSpace(*z.PersonaName))
	}
	if z.SystemPrompt != nil {
		uczestnik.PromptSystemowy = wskaznikTekstu(strings.TrimSpace(*z.SystemPrompt))
	}
	if z.Avatar != nil {
		uczestnik.Awatar = wskaznikTekstu(strings.TrimSpace(*z.Avatar))
	}
	if z.RoleDescription != nil {
		uczestnik.OpisRoli = wskaznikTekstu(strings.TrimSpace(*z.RoleDescription))
	}
	if z.Role != nil {
		uczestnik.Rola = strings.TrimSpace(string(*z.Role))
	}
	if z.Weight != nil {
		if *z.Weight < 0 {
			return shared.RoundtableModelUpdateResponse{},
				bladWskazaniaDebaty("waga kompetencji nie może być ujemna")
		}
		uczestnik.Waga = *z.Weight
	}
	if z.Key != nil {
		uczestnik.Kluczowy = *z.Key
	}
	if z.AgentId != nil {
		uczestnik.Agent = wskaznikTekstu(strings.TrimSpace(*z.AgentId))
	}

	if err := a.repozytorium.ZmienUczestnika(ctx, uczestnik); err != nil {
		return shared.RoundtableModelUpdateResponse{}, bladNieznanegoUczestnika(kod, err)
	}
	po, err := a.repozytorium.Uczestnik(ctx, kod)
	if err != nil {
		return shared.RoundtableModelUpdateResponse{}, bladNieznanegoUczestnika(kod, err)
	}
	return shared.RoundtableModelUpdateResponse{Participant: uczestnikKontraktu(po)}, nil
}

func (a *adapterDebaty) ZapiszZespol(ctx context.Context,
	z shared.RoundtableTeamSaveRequest) (shared.RoundtableTeamSaveResponse, error) {

	okno := strings.TrimSpace(z.WindowId)
	nazwa := strings.TrimSpace(z.Name)
	if okno == "" {
		return shared.RoundtableTeamSaveResponse{},
			bladWskazaniaDebaty("komenda bez wskazania okna debaty")
	}
	if nazwa == "" {
		return shared.RoundtableTeamSaveResponse{}, bladWskazaniaDebaty("zespół bez nazwy")
	}
	uczestnicy, err := a.sklad(ctx, okno)
	if err != nil {
		return shared.RoundtableTeamSaveResponse{}, err
	}

	format := ""
	if z.IncludeFormat != nil && *z.IncludeFormat {
		// Format bierze się z tury ostatniej — w niej ten skład debatował.
		if tury, err := a.repozytorium.Tury(ctx, okno, 1); err == nil && len(tury) > 0 {
			format = tury[0].Format
		}
	}

	skladZespolu := make([]dane.UczestnikZespoluDebaty, 0, len(uczestnicy))
	for _, uczestnik := range uczestnicy {
		skladZespolu = append(skladZespolu, dane.UczestnikZespoluDebaty{
			KanalModelu:     uczestnik.KanalModelu,
			NazwaTozsamosci: uczestnik.NazwaTozsamosci,
			PromptSystemowy: uczestnik.PromptSystemowy,
			Rola:            uczestnik.Rola,
			Waga:            uczestnik.Waga,
			Awatar:          uczestnik.Awatar,
			OpisRoli:        uczestnik.OpisRoli,
		})
	}
	zespol, err := a.repozytorium.ZapiszZespol(ctx, dane.ZespolDebaty{
		Kod: nowyIdentyfikator(przedrostekZespoluDebaty), Nazwa: nazwa, Format: format,
		Uczestnicy: skladZespolu,
	})
	if err != nil {
		return shared.RoundtableTeamSaveResponse{}, bladDebaty(err)
	}
	return shared.RoundtableTeamSaveResponse{Team: zespolKontraktu(zespol)}, nil
}

func (a *adapterDebaty) Zespoly(ctx context.Context,
	z shared.RoundtableTeamListRequest) (shared.RoundtableTeamListResponse, error) {

	zespoly, err := a.repozytorium.Zespoly(ctx, wartoscTekstu(z.Query), wartoscLiczby(z.Limit))
	if err != nil {
		return shared.RoundtableTeamListResponse{}, bladDebaty(err)
	}
	wykaz := make([]shared.RoundtableTeam, 0, len(zespoly))
	for _, zespol := range zespoly {
		wykaz = append(wykaz, zespolKontraktu(zespol))
	}
	return shared.RoundtableTeamListResponse{Teams: wykaz}, nil
}

// Kanał nieobecny w rejestrze odmawia całości: zespół wniesiony po połowie zostawiłby skład, którego nikt nie wybrał.
func (a *adapterDebaty) WniesZespol(ctx context.Context,
	z shared.RoundtableTeamApplyRequest) (shared.RoundtableTeamApplyResponse, error) {

	okno := strings.TrimSpace(z.WindowId)
	kod := strings.TrimSpace(z.TeamId)
	if okno == "" {
		return shared.RoundtableTeamApplyResponse{},
			bladWskazaniaDebaty("komenda bez wskazania okna debaty")
	}
	if kod == "" {
		return shared.RoundtableTeamApplyResponse{},
			bladWskazaniaDebaty("wniesienie bez wskazania zespołu")
	}
	if a.kanaly == nil {
		return shared.RoundtableTeamApplyResponse{}, bladBrakuKanalow()
	}
	zespol, err := a.repozytorium.Zespol(ctx, kod)
	if err != nil {
		return shared.RoundtableTeamApplyResponse{}, bladNieznanegoZespolu(kod, err)
	}
	for _, uczestnik := range zespol.Uczestnicy {
		if _, jest := kanalKonta(ctx, a.repozytoriumKanalow, a.kanaly, uczestnik.KanalModelu); !jest {
			return shared.RoundtableTeamApplyResponse{}, bladNieznanegoKanalu(uczestnik.KanalModelu)
		}
	}

	biezacy, err := a.repozytorium.Uczestnicy(ctx, okno)
	if err != nil {
		return shared.RoundtableTeamApplyResponse{}, bladDebaty(err)
	}
	if z.Replace != nil && *z.Replace {
		for _, uczestnik := range biezacy {
			if err := a.repozytorium.UsunUczestnika(ctx, uczestnik.Kod); err != nil {
				return shared.RoundtableTeamApplyResponse{}, bladDebaty(err)
			}
		}
		biezacy = nil
	}

	for pozycja, wpis := range zespol.Uczestnicy {
		zapisany, err := a.repozytorium.ZapiszUczestnika(ctx, dane.UczestnikDebaty{
			Kod:             nowyIdentyfikator(przedrostekUczestnika),
			Okno:            okno,
			KanalModelu:     wpis.KanalModelu,
			NazwaTozsamosci: wpis.NazwaTozsamosci,
			PromptSystemowy: wpis.PromptSystemowy,
			Kolejnosc:       len(biezacy) + pozycja + 1,
		})
		if err != nil {
			return shared.RoundtableTeamApplyResponse{}, bladDebaty(err)
		}
		// Rola, waga i opis idą osobnym zapisem po dopisaniu uczestnika.
		zapisany.Rola, zapisany.Waga = wpis.Rola, wpis.Waga
		zapisany.Awatar, zapisany.OpisRoli = wpis.Awatar, wpis.OpisRoli
		if err := a.repozytorium.ZmienUczestnika(ctx, zapisany); err != nil {
			return shared.RoundtableTeamApplyResponse{}, bladDebaty(err)
		}
	}

	uczestnicy, err := a.repozytorium.Uczestnicy(ctx, okno)
	if err != nil {
		return shared.RoundtableTeamApplyResponse{}, bladDebaty(err)
	}
	return shared.RoundtableTeamApplyResponse{Participants: uczestnicyKontraktu(uczestnicy)}, nil
}

func (a *adapterDebaty) Role(ctx context.Context,
	z shared.RoundtableRoleListRequest) (shared.RoundtableRoleListResponse, error) {

	role, err := a.repozytorium.Role(ctx, wartoscTekstu(z.Query))
	if err != nil {
		return shared.RoundtableRoleListResponse{}, bladDebaty(err)
	}
	wykaz := make([]shared.RoundtableRole, 0, len(role))
	for _, rola := range role {
		fabryczna := rola.Fabryczna
		wykaz = append(wykaz, shared.RoundtableRole{
			Id: rola.Kod, Name: rola.Nazwa, SystemPrompt: rola.PromptSystemowy,
			Description: rola.Opis, BuiltIn: &fabryczna,
		})
	}
	return shared.RoundtableRoleListResponse{Roles: wykaz}, nil
}

// Uczestnicy zespołu nie mają okna ani identyfikatora uczestnika: są kopią składu, nie składem w oknie.
func zespolKontraktu(z dane.ZespolDebaty) shared.RoundtableTeam {
	sklad := make([]shared.RoundtableParticipant, 0, len(z.Uczestnicy))
	for _, wpis := range z.Uczestnicy {
		waga := wpis.Waga
		uczestnik := shared.RoundtableParticipant{
			ChannelId:       wpis.KanalModelu,
			PersonaName:     wpis.NazwaTozsamosci,
			SystemPrompt:    wpis.PromptSystemowy,
			Weight:          &waga,
			Avatar:          wpis.Awatar,
			RoleDescription: wpis.OpisRoli,
		}
		if rola := strings.TrimSpace(wpis.Rola); rola != "" {
			wartosc := shared.RoundtableParticipantRole(rola)
			uczestnik.Role = &wartosc
		}
		sklad = append(sklad, uczestnik)
	}
	zespol := shared.RoundtableTeam{
		Id: z.Kod, Name: z.Nazwa, Participants: sklad, CreatedAt: chwilaBazy(z.Utworzono),
	}
	if format := strings.TrimSpace(z.Format); format != "" {
		wartosc := shared.RoundtableFormat(format)
		zespol.Format = &wartosc
	}
	return zespol
}

func turyKontraktu(tury []dane.TuraDebaty) []shared.RoundtableTurn {
	wykaz := make([]shared.RoundtableTurn, 0, len(tury))
	for _, tura := range tury {
		wykaz = append(wykaz, turaKontraktu(tura))
	}
	return wykaz
}

func wypowiedziKontraktu(wypowiedzi []dane.WypowiedzDebaty) []shared.RoundtableStatement {
	wykaz := make([]shared.RoundtableStatement, 0, len(wypowiedzi))
	for _, wypowiedz := range wypowiedzi {
		wykaz = append(wykaz, wypowiedzKontraktu(wypowiedz))
	}
	return wykaz
}
