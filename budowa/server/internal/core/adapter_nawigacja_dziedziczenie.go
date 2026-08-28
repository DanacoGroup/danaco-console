package core

import (
	"context"

	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// parametryOknaModulu składa ustawienia okna zakładanego przez workspace.enter,
// dobierając kanał modelu z okna, którym sesja już rozmawia, a bez takiego
// okna sięgając po pierwszy czynny kanał rejestru.
func (a *adapterNawigacji) parametryOknaModulu(ctx context.Context, idSesji, kodModulu string) session.Ustawienia {
	u := session.Ustawienia{Modul: kodModulu, RolaOkna: shared.WindowRoleStandalone}
	a.dziedziczPoOknachSesji(idSesji, &u)
	if u.KanalModelu == "" {
		u.KanalModelu = a.pierwszyCzynnyKanal(ctx)
	}
	return u
}

// dziedziczPoOknachSesji przepisuje komplet parametrów wykonania z najświeższego
// okna sesji z niepustym kanałem, tak aby nowe okno modułu pracowało w tych
// samych katalogach, środowisku i trybie uprawnień co reszta sesji.
func (a *adapterNawigacji) dziedziczPoOknachSesji(idSesji string, u *session.Ustawienia) {
	if a.nadzorca == nil {
		return
	}
	okna, err := a.nadzorca.Rejestr().OknaSesji(idSesji)
	if err != nil {
		return
	}
	for i := len(okna) - 1; i >= 0; i-- {
		if okna[i].KanalModelu == "" {
			continue
		}
		u.KanalModelu = okna[i].KanalModelu
		u.KatalogiRobocze = okna[i].KatalogiRobocze
		u.SrodowiskoWykonania = okna[i].SrodowiskoWykonania
		u.TrybUprawnien = okna[i].TrybUprawnien
		return
	}
}

// pierwszyCzynnyKanal zwraca kod pierwszego czynnego kanału rejestru albo pusty
// napis. Rejestr oddaje kanały w kolejności ustawionej przez Operatora, więc
// „pierwszy" znaczy „ten, który Operator postawił na czele", a nie „dowolny".
func (a *adapterNawigacji) pierwszyCzynnyKanal(ctx context.Context) string {
	if a.zestaw == nil || a.zestaw.Kanaly == nil {
		return ""
	}
	kanaly, err := a.zestaw.Kanaly.Lista(ctx, true)
	if err != nil || len(kanaly) == 0 {
		return ""
	}
	return kanaly[0].Kod
}
