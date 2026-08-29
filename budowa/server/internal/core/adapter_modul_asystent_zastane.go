// Plik obsługuje przegląd startowy zleceń asystenta: jedną czynność wykonywaną raz przy montażu rdzenia, domykającą zlecenia zastane w stanie `running` po procesie, który już nie żyje.
package core

import (
	"context"

	"danacoconsole/shared"
)

// ZDomknieciemZastanych przegląda zlecenia zastane w stanie `running` i domyka je powodem nazwanym. Rdzeń zatrzymany w połowie tury zostawia zlecenie w `running`; uczciwym stanem jest wtedy `failed` z powodem, nie dalsze udawanie pracy.
func (a *adapterAsystenta) ZDomknieciemZastanych() *adapterAsystenta {
	if a.nadzorca == nil || a.zycie == nil {
		return a
	}
	go a.domknijZastaneZlecenia(a.zycie)
	return a
}

// domknijZastaneZlecenia obchodzi okna odtworzone z bazy i zamyka każde
// zlecenie stojące w `running`. Droga domknięcia jest ta sama co przy zerwaniu
// w trakcie pracy (`zerwijZlecenie`) — jedno domknięcie, nie dwa.
func (a *adapterAsystenta) domknijZastaneZlecenia(ctx context.Context) {
	rejestr := a.nadzorca.Rejestr()
	for _, sesja := range rejestr.Sesje() {
		okna, err := rejestr.OknaSesji(sesja.Id)
		if err != nil {
			continue
		}
		for _, okno := range okna {
			zlecenia, err := a.repozytorium.Zlecenia(ctx, okno.Id)
			if err != nil {
				continue
			}
			for _, zlecenie := range zlecenia {
				if zlecenie.Stan != string(shared.AssistantActionStatusRunning) {
					continue
				}
				a.zerwijZlecenie(ctx, zlecenie, zlecenie.Kod, okno.IdSesji,
					"serwer został zatrzymany w trakcie tury tego zlecenia, a tura nie przeżywa"+
						" zatrzymania procesu; naprawa: ponowić zlecenie (control: retry)")
			}
		}
	}
}
