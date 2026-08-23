// Odpowiedzialność pliku: wpięcie komendy `advisor.consult` — jedynego odbiorcy
// konsultacji u doradcy — oraz rozgłoszenie zdarzenia `advisor.consulted`.
// Adapter wraz z rozstrzygnięciami leży w `adapter_doradcy.go`, pojęcie doradcy
// i sufit siły w `podagenci/doradca.go` oraz `podagenci/doradca_wybor.go`.
//
// Odbiorcą rady jest model w trakcie tury, bo to on staje przed
// rozstrzygnięciem, w którym rada się przydaje.
//
// Rodzina ma jedną komendę: kontrakt nie zna ani `advisor.list`, ani
// `advisor.history`. Wykaz kandydatów oddaje `Doradcy()` adaptera, a historia
// konsultacji leży w dzienniku `konsultacja_doradcy` i czeka na własną komendę.
// Dopisanie nazw, których klient nie zna, byłoby rozrostem kontraktu bez
// odbiorcy.
//
// Zdarzenie jest częścią czynności: rada ma być jawna w strumieniu
// i w prowenancji, inaczej agent działa na przesłance, której w aktach nie ma.
// Dlatego port oddaje obok wyniku gotowy ładunek zdarzenia — skrót rady liczy
// `podagenci.ZlozRade`, więc obsługiwacz go nie przelicza.
package core

import (
	"context"

	"danacoconsole/shared"
)

// JawnoscKonsultacji niesie to, czego wynik komendy nie unosi, a jawność
// wymaga: kartę sesji do koperty zdarzenia i ładunek `advisor.consulted`.
//
// Sesja nie stoi w ładunku zdarzenia, bo kontrakt jej tam nie ma — zdarzenie
// wskazuje okno, a koperta wskazuje kartę sesji, w której to okno pracuje.
// Adapter zna oba fakty z danych okna, obsługiwacz żadnego z nich nie zna.
type JawnoscKonsultacji struct {
	IdSesji   string
	Zdarzenie shared.AdvisorConsultedEvent
}

// Doradcy jest portem rodziny `advisor.*`. Mówi wyłącznie typami
// kontraktu; kanał doradcy, dobór pod sufitem siły i wpis dziennika leżą po
// drugiej stronie adaptera.
type Doradcy interface {
	// Konsultacja obsługuje `advisor.consult`. Drugi wynik niesie jawność:
	// odbytą konsultację się rozgłasza, a nie tylko oddaje pytającemu.
	Konsultacja(ctx context.Context, z shared.AdvisorConsultRequest) (
		shared.AdvisorConsultResponse, JawnoscKonsultacji, error)
}

// zarejestrujDoradcow wpina jedyną komendę rodziny `advisor.*`.
//
// Zdarzenie idzie wyłącznie po konsultacji udanej. Odmowa doboru — a taką jest
// każda próba sięgnięcia przez model po model silniejszy bez wyraźnego
// wskazania — wraca pytającemu błędem i ląduje w dzienniku konsultacji jako
// wpis `odmowa`. Rozgłoszenie „skonsultowano" po konsultacji, która się nie
// odbyła, byłoby fałszywym powiadomieniem.
func zarejestrujDoradcow(r *Rejestr, d Doradcy, e *emiter) {
	if r == nil || d == nil {
		return
	}

	r.Zarejestruj(shared.CommandAdvisorConsult,
		obsluz(func(ctx context.Context, z shared.AdvisorConsultRequest) (shared.AdvisorConsultResponse, error) {
			odpowiedz, jawnosc, err := d.Konsultacja(ctx, z)
			if err == nil {
				e.konsultacja(jawnosc)
			}
			return odpowiedz, err
		}))
}

// konsultacja rozgłasza odbytą konsultację u doradcy. Zdarzenie jedzie z kartą
// sesji okna pytającego — tak samo jak zmiana okna czy zlecenie asystenta —
// żeby powierzchnia obserwująca tę sesję dostała radę tam, gdzie pracuje agent.
func (e *emiter) konsultacja(j JawnoscKonsultacji) {
	e.wyslij(shared.EventAdvisorConsulted, j.IdSesji, j.Zdarzenie)
}
