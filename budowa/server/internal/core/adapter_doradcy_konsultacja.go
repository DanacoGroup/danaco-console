// Odpowiedzialność pliku: komenda `advisor.consult` — jedyna droga, którą model
// prosi o radę, i miejsce, w którym rada staje się widoczna w oknie.
//
// Adapter doradcy (`adapter_doradcy.go`) wykonuje konsultację tą samą drogą, co
// każdy inny wołacz kanału. Tu stoi to, co należy wyłącznie do komendy: odczyt
// okna, złożenie pytania z danych okna, strumień jawności i przełożenie wyniku
// na kontrakt.
//
// Jawność nie jest buforem: ujście konsultacji jest nadawcą strumienia, więc
// fragmenty jadą kopertami `stream.chunk` do okna pytającego — tą samą drogą,
// którą płynie odpowiedź agenta, debata Roundtable i wyjście Terminala. Bez
// tego Operator widziałby wyłącznie skrót rady w zdarzeniu `advisor.consulted`,
// a pytanie, doradca i rada nie stanęłyby w jednym oknie. Drugiej drogi do okna
// rdzeń nie ma i tutaj też jej nie zakładamy.
package core

import (
	"context"
	"strings"

	"danacoconsole/server/internal/models"
	"danacoconsole/server/internal/podagenci"
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// Konsultacja obsługuje `advisor.consult`: przekłada żądanie modelu na pytanie
// do doradcy, wykonuje konsultację tą samą drogą co każdy inny wołacz, pokazuje
// ją w oknie i oddaje radę wraz z gotowym ładunkiem zdarzenia jawności.
//
// Kim jest pytający, rozstrzyga okno, a nie żądanie. Model podaje wyłącznie
// sprawę, kontekst i prośbę o doradcę; gdyby `Pytajacy` dało się podać
// żądaniem, model podniósłby sobie sufit jednym polem.
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
	// Okno zamknięte nie konsultuje. Rada ma być widoczna tam, gdzie pracuje
	// agent — w oknie zamkniętym nie pracuje nikt, więc konsultacja szłaby
	// w miejsce, którego Operator już nie ogląda, a zdarzenie `advisor.consulted`
	// rozgłaszałoby pracę okna, które stanęło. Sprawdzenie jest tym samym, które
	// robią pozostałe moduły oknowe (adapter_modul_model.go, adapter_modul_aod.go).
	// Kod `conflict` mówi prawdę o powodzie: stan zasobu, nie brak kanału — i nie
	// jest ponawialny, bo okno samo się nie otworzy.
	if !okno.CzyOtwarte() {
		return shared.AdvisorConsultResponse{}, JawnoscKonsultacji{}, bladOknaDoradcy(
			shared.ErrorCodeConflict, "okno "+okno.Id+" jest zamknięte — konsultacja nie ma się gdzie odbyć ani komu pokazać")
	}

	// Strumień jawności zakładamy przed konsultacją i domykamy zawsze — także
	// po odmowie doboru. Identyfikator strumienia jest tożsamością tej jednej
	// konsultacji: wchodzi do koperty jako identyfikator żądania i do fragmentów
	// jako `messageId`, więc klient wie, że wszystkie te fragmenty należą do
	// jednego wywołania (wzorem wypowiedzi Roundtable).
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

// Podstawa doboru jest jedna i dlatego stoi wyżej stałą, a nie funkcją.
// Kontrakt zna dwie wartości `selection`, ale `operatorIndication` nie ma
// w produkcie danych pod sobą: Operator nie ma czym wskazać doradcy, a pole
// `modelChannelId` okna mówi, którym modelem pracuje okno, nie kogo Operator
// wyznaczył na doradcę. Podstawienie jednego pod drugie ogłaszałoby „wskazanie
// Operatora" przy każdej konsultacji bez prośby. Wartość kontraktu zostaje na
// czas, gdy wskazanie Operatora będzie miało skąd pochodzić (nagłówek
// `podagenci/doradca_wybor.go`).

// pytanieZOkna składa pytanie do doradcy: sprawa i kontekst od modelu, kanał
// pytającego z okna, tożsamość strumienia od komendy.
//
// Pole `Wiadomosc` niesie identyfikator strumienia, bo to ono zasila `messageId`
// wszystkich fragmentów kanału (`models/fragment.go`) — bez niego fragmenty rady
// jechałyby do okna bez wskazania, do czego należą.
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
// tekst, którym strumień zostanie domknięty.
//
// Zbiera cały tekst, łącznie z blokiem jawności doklejanym przez `Skonsultuj` —
// fragment domykający ma być tym, co Operator widzi w oknie po konsultacji,
// czyli radą wraz z podpisem, kto jej udzielił.
type ujscieKonsultacji struct {
	strumien *nadawcaStrumienia
	zebrane  strings.Builder
}

// Fragment przyjmuje jeden fragment strumienia doradcy.
func (u *ujscieKonsultacji) Fragment(ctx context.Context, f models.Fragment) error {
	if f.Kind == shared.ChunkKindText {
		u.zebrane.WriteString(models.TrescFragmentu(f))
	}
	if u.strumien == nil {
		return nil
	}
	return u.strumien.Fragment(ctx, f)
}
