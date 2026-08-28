// Plik wpina komendę `advisor.consult`, jedynego odbiorcę konsultacji
// u doradcy, i rozgłasza zdarzenie `advisor.consulted` po udanej konsultacji.
package core

import (
	"context"

	"danacoconsole/shared"
)

// JawnoscKonsultacji niesie kartę sesji do koperty zdarzenia i ładunek
// `advisor.consulted`, czego wynik komendy nie unosi.
type JawnoscKonsultacji struct {
	IdSesji   string
	Zdarzenie shared.AdvisorConsultedEvent
}

// Doradcy jest portem rodziny `advisor.*`. Mówi wyłącznie typami
// kontraktu; kanał doradcy, dobór pod sufitem siły i wpis dziennika leżą po
// drugiej stronie adaptera.
type Doradcy interface {
	// Konsultacja obsługuje `advisor.consult` i oddaje jawność: odbytą
	// konsultację się rozgłasza.
	Konsultacja(ctx context.Context, z shared.AdvisorConsultRequest) (
		shared.AdvisorConsultResponse, JawnoscKonsultacji, error)
}

// zarejestrujDoradcow wpina jedyną komendę rodziny `advisor.*` i rozgłasza
// konsultację wyłącznie po jej udanym przebiegu.
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
