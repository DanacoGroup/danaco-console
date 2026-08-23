package core

import (
	"context"

	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// parametryOknaModulu składa ustawienia okna zakładanego przez `workspace.enter`.
//
// Powód istnienia tej funkcji: wejście do modułu niesie w kontrakcie wyłącznie
// sesję, moduł i ewentualne okno do ponownego użycia. Gdyby okno powstawało
// z samego modułu, rodziłoby się bez kanału modelu — a więc jako okno, które
// się rysuje i nie umie rozmawiać.
//
// Kanał bierzemy stąd, skąd Operator by go oczekiwał — z tego, czym ta sesja
// już rozmawia. Dopiero gdy sesja nie rozmawiała jeszcze niczym, sięgamy po
// pierwszy czynny kanał rejestru. Brak obu daje okno bez kanału: rdzeń nie
// odmawia wejścia do modułu z powodu nieskonfigurowanego kanału,
// a okno zgłosi brak dopiero przy próbie tury.
func (a *adapterNawigacji) parametryOknaModulu(ctx context.Context, idSesji, kodModulu string) session.Ustawienia {
	u := session.Ustawienia{Modul: kodModulu, RolaOkna: shared.WindowRoleStandalone}
	a.dziedziczPoOknachSesji(idSesji, &u)
	if u.KanalModelu == "" {
		u.KanalModelu = a.pierwszyCzynnyKanal(ctx)
	}
	return u
}

// dziedziczPoOknachSesji przepisuje parametry wykonania z okna, które w tej
// sesji już rozmawia. Bierzemy okno ostatnie z niepustym kanałem — jest nim
// okno najświeższe, a więc to, którego ustawienia Operator widział ostatnio.
//
// Dziedziczy się komplet parametrów wykonania, nie sam kanał: nowe okno modułu
// ma pracować w tych samych katalogach, w tym samym środowisku i na tym samym
// trybie uprawnień co reszta sesji. Rola zostaje samodzielna — okno modułu nie
// wchodzi do cudzej pętli koordynator–wykonawca.
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
