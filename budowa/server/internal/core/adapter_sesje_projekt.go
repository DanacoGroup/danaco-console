package core

import (
	"context"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// przedrostekProjektu znakuje kod projektu założonego przy przenoszeniu sesji.
// Projekt zakładany z historii nie ma jeszcze wiersza w module Workspace, więc
// kod nadaje rdzeń — nazwa pozostaje tym, co Operator wpisał.
const przedrostekProjektu = "prj-"

// ZProjektami wpina repozytorium projektów. Bez niego przenoszenie sesji do
// projektu odmawia wprost, zamiast zapisywać wskazanie, którego nikt nie zna.
func (a *adapterSesji) ZProjektami(p dane.RepozytoriumPrzestrzeniRoboczej) *adapterSesji {
	a.projekty = p
	return a
}

// PrzypiszProjekt przenosi sesję do wskazanego projektu, a przy podaniu samej nazwy zamiast
// istniejącego kodu zakłada nowy projekt i przenosi sesję do niego.
func (a *adapterSesji) PrzypiszProjekt(ctx context.Context,
	z shared.SessionProjectSetRequest) (shared.SessionProjectSetResponse, error) {

	kod, err := a.ustalProjekt(ctx, z)
	if err != nil {
		return shared.SessionProjectSetResponse{}, err
	}
	return shared.SessionProjectSetResponse{
		ProjectId: kod,
		MovedIds:  a.przestawProjekt(ctx, z.SessionIds, kod),
	}, nil
}

// OdepnijProjekt wyjmuje sesje z projektu.
//
// Sesja zostaje w historii — wyjęcie z projektu NIE JEST usunięciem. To jest
// osobna czynność Operatora („usuń z projektu" wobec „usuń"), więc i osobna
// komenda; pomylenie ich kosztowałoby zapis.
func (a *adapterSesji) OdepnijProjekt(ctx context.Context,
	z shared.SessionProjectClearRequest) (shared.SessionProjectClearResponse, error) {

	return shared.SessionProjectClearResponse{
		ClearedIds: a.przestawProjekt(ctx, z.SessionIds, ""),
	}, nil
}

// ustalProjekt rozstrzyga, do którego projektu trafiają sesje.
//
// Wskazanie istniejącego projektu jest sprawdzane, a nie przyjmowane na słowo:
// zapisanie sesji do projektu, którego nie ma, dałoby wykaz wskazujący w pustkę.
func (a *adapterSesji) ustalProjekt(ctx context.Context,
	z shared.SessionProjectSetRequest) (string, error) {

	if a.projekty == nil {
		return "", protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"serwer nie ma repozytorium projektów — przeniesienie sesji nie ma dokąd trafić"))
	}
	if z.ProjectId != nil && *z.ProjectId != "" {
		if _, err := a.projekty.Projekt(ctx, *z.ProjectId); err != nil {
			return "", protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
				"projekt "+*z.ProjectId+" nie istnieje"))
		}
		return *z.ProjectId, nil
	}
	if z.ProjectName == nil || *z.ProjectName == "" {
		return "", protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
			"przeniesienie wymaga wskazania projektu albo nazwy nowego"))
	}
	projekt, _, err := a.projekty.ZapewnijProjekt(ctx, nowyIdentyfikator(przedrostekProjektu), *z.ProjectName)
	if err != nil {
		return "", err
	}
	return projekt.Kod, nil
}

// OdepnijProjektZywych zdejmuje kod projektu z sesji żywych rejestru nadzorcy
// i zwraca ich identyfikatory. Idzie po `project.delete`: baza odpięła sesje
// w transakcji kasowania, a rejestr w pamięci sam tego nie widzi i niósłby kod
// projektu, którego nie ma.
func (a *adapterSesji) OdepnijProjektZywych(kod string) []string {
	odpiete := []string{}
	if a == nil || a.nadzorca == nil || kod == "" {
		return odpiete
	}
	rejestr := a.nadzorca.Rejestr()
	for _, sesja := range rejestr.Sesje() {
		if sesja.IdProjektu != kod {
			continue
		}
		if _, err := rejestr.ZmienProjektSesji(sesja.Id, ""); err != nil {
			continue
		}
		odpiete = append(odpiete, sesja.Id)
	}
	return odpiete
}

// przestawProjekt zapisuje przynależność sesji i zwraca te, które faktycznie
// przestawiono. Wskazanie bez odpowiednika jest pomijane — czynność zbiorcza nie
// może paść przez jedną pozycję.
func (a *adapterSesji) przestawProjekt(ctx context.Context, wskazania []string, kod string) []string {
	przestawione := make([]string, 0, len(wskazania))
	for _, identyfikator := range wskazania {
		if a.trwalosc == nil || a.trwalosc.sesje == nil {
			continue
		}
		wiersz, err := a.trwalosc.sesje.PoIdentyfikatorze(ctx, identyfikator)
		if err != nil {
			continue
		}
		if err := a.trwalosc.sesje.ZmienProjekt(ctx, wiersz.ID, kod); err != nil {
			continue
		}
		_, _ = a.nadzorca.Rejestr().ZmienProjektSesji(identyfikator, kod)
		przestawione = append(przestawione, identyfikator)
	}
	return przestawione
}
