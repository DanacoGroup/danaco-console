// Odpowiedzialność pliku: przegląd startowy zleceń asystenta — jedna czynność,
// wykonywana raz przy montażu rdzenia: domknięcie zleceń zastanych w stanie
// `running`. Wykonawca zleceń (`adapter_modul_asystent_wykonawca.go`) prowadzi
// zlecenia biegnące teraz; ten plik zajmuje się wyłącznie spadkiem po procesie,
// który już nie żyje. Osobna odpowiedzialność, osobny plik, ten sam
// typ `adapterAsystenta`.
package core

import (
	"context"

	"danacoconsole/shared"
)

// ZDomknieciemZastanych przegląda zlecenia zastane w stanie `running` i domyka
// je powodem nazwanym.
//
// Rdzeń zatrzymany w połowie tury zostawia zlecenie w `running`, a wykonawca
// wchodzi wyłącznie drogą świeżego polecenia, `retry` albo `resume` — bez tego
// przeglądu Actions Monitor pokazywałby wiersz „w toku" bez wyniku i bez powodu
// na zawsze. Tura tamtego zlecenia nie istnieje, bo umarła z procesem, więc
// uczciwym stanem jest `failed` z powodem, a nie dalsze udawanie pracy.
//
// Przegląd idzie gorutyną, bo montaż nie ma czekać na obejście wszystkich okien;
// bez nadzorcy sesji albo kontekstu życia nie ma czego przeglądać.
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
					"rdzeń został zatrzymany w trakcie tury tego zlecenia, a tura nie przeżywa"+
						" zatrzymania procesu; naprawa: ponowić zlecenie (control: retry)")
			}
		}
	}
}
