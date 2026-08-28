// Plik obsługuje komendę advisor.consult: jedyną drogę, którą model prosi o
// radę doradcy, i miejsce, w którym rada staje się widoczna w oknie.
package core

import (
	"context"
	"strings"

	"danacoconsole/server/internal/models"
	"danacoconsole/server/internal/podagenci"
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// Konsultacja obsługuje advisor.consult: przekłada żądanie modelu na pytanie
// do doradcy, wykonuje konsultację tą samą drogą co każdy inny wołacz, pokazuje
// ją w oknie i oddaje radę wraz z gotowym ładunkiem zdarzenia jawności.
func (a *adapterDoradcow) Konsultacja(ctx context.Context, z shared.AdvisorConsultRequest) (
	shared.AdvisorConsultResponse, JawnoscKonsultacji, error) {

	if a.okna == nil {
		return shared.AdvisorConsultResponse{}, JawnoscKonsultacji{}, bladOknaDoradcy(
			shared.ErrorCodeInternalError, "rejestr okien nie jest wpięty — nie ma skąd wziąć kanału pytającego")
	}
	okno, err := a.okna.Okno(strings.TrimSpace(z.WindowId))
	if err != nil {
		return shared.AdvisorConsultResponse{}, JawnoscKonsultacji{}, bladOknaDoradcy(
			shared.ErrorCodeNotFound, "okno "+z.WindowId+" nie istnieje — konsultacja nie ma czyim kanałem pytać")
	}
	// Okno zamknięte nie konsultuje: rada ma być widoczna tam, gdzie pracuje agent.
	if !okno.CzyOtwarte() {
		return shared.AdvisorConsultResponse{}, JawnoscKonsultacji{}, bladOknaDoradcy(
			shared.ErrorCodeConflict, "okno "+okno.Id+" jest zamknięte — konsultacja nie ma się gdzie odbyć ani komu pokazać")
	}

	// Strumień jawności powstaje przed konsultacją i domyka się zawsze, także po
	// odmowie doboru doradcy.
	idStrumienia := nowyIdentyfikator(przedrostekKonsultacji)
	strumien := nowyNadawcaStrumienia(a.nadajnik, idStrumienia, okno.IdSesji)
	ujscie := &ujscieKonsultacji{strumien: strumien}

	rada, err := a.Skonsultuj(ctx, pytanieZOkna(okno, idStrumienia, z), ujscie)
	strumien.Zakoncz(okno.Id, idStrumienia, ujscie.zebrane.String(), err)
	if err != nil {
		return shared.AdvisorConsultResponse{}, JawnoscKonsultacji{}, err
	}

	return shared.AdvisorConsultResponse{
			Advice:         rada.Tresc,
			AdvisorChannel: rada.Doradca.Kanal,
			AdvisorModel:   wskaznikTekstu(rada.Doradca.Model),
			Selection:      shared.AdvisorSelectionStrengthCeiling,
		}, JawnoscKonsultacji{
			IdSesji: okno.IdSesji,
			Zdarzenie: shared.AdvisorConsultedEvent{
				WindowId:       okno.Id,
				AdvisorChannel: rada.Doradca.Kanal,
				AdvisorModel:   wskaznikTekstu(rada.Doradca.Model),
				Selection:      shared.AdvisorSelectionStrengthCeiling,
				AdviceDigest:   rada.Skrot,
			},
		}, nil
}

// Wartość doboru jest stała, ponieważ dane źródłowe wskazania Operatora nie są jeszcze dostępne.

// pytanieZOkna składa pytanie do doradcy: sprawę i kontekst od modelu, kanał
// pytającego z okna oraz tożsamość strumienia od komendy.
func pytanieZOkna(okno session.Okno, idStrumienia string, z shared.AdvisorConsultRequest) podagenci.Pytanie {
	return podagenci.Pytanie{
		Okno:          okno.Id,
		Wiadomosc:     idStrumienia,
		Pytajacy:      strings.TrimSpace(okno.KanalModelu),
		ZadanyDoradca: strings.TrimSpace(wartoscTekstu(z.RequestedAdvisor)),
		Sprawa:        z.Question,
		Kontekst:      wartoscTekstu(z.Context),
	}
}

// ujscieKonsultacji przepuszcza fragmenty konsultacji do okna i zapamiętuje
// tekst, którym strumień zostanie domknięty po zakończeniu.
type ujscieKonsultacji struct {
	strumien *nadawcaStrumienia
	zebrane  strings.Builder
}

// Fragment przyjmuje jeden fragment strumienia doradcy, dopisuje go do zebranego
// tekstu przy rodzaju tekstowym i przekazuje go dalej do nadawcy strumienia okna.
func (u *ujscieKonsultacji) Fragment(ctx context.Context, f models.Fragment) error {
	if f.Kind == shared.ChunkKindText {
		u.zebrane.WriteString(models.TrescFragmentu(f))
	}
	if u.strumien == nil {
		return nil
	}
	return u.strumien.Fragment(ctx, f)
}
