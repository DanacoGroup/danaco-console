// Plik obsługuje provenance.call.replay: powtórzenie wywołania modelu i zestawienie odpowiedzi z pierwowzorem; port bierze rejestr kanałów, bo powtórzenie jest nowym wywołaniem, nie odczytem.
package core

import (
	"context"
	"errors"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/models"
	"danacoconsole/shared"
)

// przedrostekPowtorzenia znakuje wiersz śladu wywołania powtórzonego w tym dzienniku prowenancji rdzenia.
const przedrostekPowtorzenia = "powtorzenie-"

// ZKanalami wpina rejestr kanałów — jedyną drogę, którą rdzeń wykonuje
// wywołanie modelu. Zależność opcjonalna: bez niej cztery komendy odczytu
// pracują bez zmian, a powtórzenie odmawia, nazywając brak.
func (a *adapterProwenancji) ZKanalami(kanaly *models.Rejestr) *adapterProwenancji {
	a.kanaly = kanaly
	return a
}

// PowtorzWywolanie obsługuje provenance.call.replay, powtarzając wywołanie modelu i zestawiając wynik.
func (a *adapterProwenancji) PowtorzWywolanie(ctx context.Context,
	z shared.ProvenanceCallReplayRequest) (shared.ProvenanceCallReplayResponse, error) {

	kod := strings.TrimSpace(z.CallId)
	if kod == "" {
		return shared.ProvenanceCallReplayResponse{},
			odmowaSladu(shared.ErrorCodeValidationFailed, "powtórzenie bez wskazania wywołania")
	}
	if a.repozytorium == nil {
		return shared.ProvenanceCallReplayResponse{},
			odmowaSladu(shared.ErrorCodeInternalError, "rdzeń nie ma wpiętego magazynu śladu")
	}
	pierwowzor, err := a.repozytorium.Wywolanie(ctx, kod)
	if err != nil {
		if errors.Is(err, dane.ErrBrakWiersza) {
			return shared.ProvenanceCallReplayResponse{},
				odmowaSladu(shared.ErrorCodeNotFound, "nie ma wywołania o identyfikatorze "+kod)
		}
		return shared.ProvenanceCallReplayResponse{},
			odmowaSladu(shared.ErrorCodeInternalError, err.Error())
	}

	prompt := strings.TrimSpace(wartoscTekstu(pierwowzor.Prompt))
	if !pierwowzor.TrescZapisana || prompt == "" {
		// To nie jest awaria: ślad bez treści jest śladem poprawnym, zgodnym z nastawą prywatności platformy.
		return shared.ProvenanceCallReplayResponse{
			Replay: shared.ModelCallReplay{
				OriginalCallId: kod, ContentAvailable: false,
			},
		}, nil
	}
	if a.kanaly == nil {
		return shared.ProvenanceCallReplayResponse{}, odmowaSladu(shared.ErrorCodeChannelUnavailable,
			"rdzeń nie ma wpiętego rejestru kanałów — powtórzenie jest nowym wywołaniem "+
				"kanału modelu, a nie odczytem śladu, więc bez rejestru nie ma czym go wykonać")
	}

	kanal := strings.TrimSpace(wartoscTekstu(z.ChannelId))
	if kanal == "" {
		kanal = strings.TrimSpace(wartoscTekstu(pierwowzor.KanalKod))
	}
	if kanal == "" {
		return shared.ProvenanceCallReplayResponse{}, odmowaSladu(shared.ErrorCodeValidationFailed,
			"pierwowzór nie ma zapisanego kanału, a żądanie go nie wskazuje — "+
				"rdzeń nie dobiera kanału za Operatora, bo powtórzenie na innym kanale "+
				"jest innym doświadczeniem niż to, które miało zostać powtórzone")
	}
	if _, jest := a.kanaly.Kanal(kanal); !jest {
		return shared.ProvenanceCallReplayResponse{}, odmowaSladu(shared.ErrorCodeChannelUnavailable,
			"kanału „"+kanal+"” nie ma w rejestrze kanałów rdzenia albo jest wyłączony")
	}

	kodPowtorzenia := nowyIdentyfikator(przedrostekPowtorzenia)
	var odpowiedz strings.Builder
	ujscie := models.UjscieFunkcji(func(_ context.Context, f models.Fragment) error {
		if f.Kind == shared.ChunkKindText {
			odpowiedz.WriteString(models.TrescFragmentu(f))
		}
		return nil
	})
	zapytanie := models.Zapytanie{
		Zasiegi:   zasiegiPowtorzenia(pierwowzor),
		Wiadomosc: kodPowtorzenia,
		Tresc:     prompt,
		Kanal:     kanal,
	}
	if model := strings.TrimSpace(wartoscTekstu(z.Model)); model != "" {
		zapytanie.Model = model
	}

	poczatek := time.Now()
	blad := a.kanaly.Wyslij(ctx, zapytanie, ujscie)
	opoznienie := int(time.Since(poczatek).Milliseconds())

	tresc := odpowiedz.String()
	// Stan zapisuje się jako baza, przez odwzorowanie kontraktu, jak czyta stanKontraktuWywolania.
	stan := shared.WartosciBazyModelCallStatus[shared.ModelCallStatusOk]
	if blad != nil {
		stan = shared.WartosciBazyModelCallStatus[shared.ModelCallStatusFailed]
	}

	// Ślad powtórzenia zapisuje się zawsze, nawet nieudany: poszło do kanału bez względu na wynik.
	powtorzenie := dane.WywolanieModelu{
		Kod:           kodPowtorzenia,
		RodzicKod:     &pierwowzor.Kod,
		SladKod:       pierwowzor.SladKod,
		SesjaKod:      pierwowzor.SesjaKod,
		OknoKod:       pierwowzor.OknoKod,
		KanalKod:      &kanal,
		Model:         modelPowtorzenia(z.Model, pierwowzor),
		Stan:          stan,
		Poczatek:      poczatek.UnixMilli(),
		OpoznienieMs:  &opoznienie,
		TrescZapisana: true,
		Prompt:        &prompt,
		Odpowiedz:     &tresc,
	}
	koniec := time.Now().UnixMilli()
	powtorzenie.Koniec = &koniec
	if blad != nil {
		komunikat := blad.Error()
		powtorzenie.KodBledu = &komunikat
	}
	if _, err := a.repozytorium.ZapiszWywolanie(ctx, powtorzenie); err != nil {
		return shared.ProvenanceCallReplayResponse{},
			odmowaSladu(shared.ErrorCodeInternalError,
				"powtórzenie się wykonało, ale jego śladu nie dało się zapisać: "+err.Error())
	}
	if blad != nil {
		return shared.ProvenanceCallReplayResponse{}, odmowaSladu(shared.ErrorCodeChannelUnavailable,
			"kanał „"+kanal+"” nie wykonał powtórzenia: "+blad.Error())
	}

	pierwotna := strings.TrimSpace(wartoscTekstu(pierwowzor.Odpowiedz))
	wynik := shared.ModelCallReplay{
		OriginalCallId:   kod,
		ReplayCallId:     kodPowtorzenia,
		ResponseChanged:  strings.TrimSpace(tresc) != pierwotna,
		ContentAvailable: true,
	}
	if roznica := roznicaOdpowiedzi(pierwotna, strings.TrimSpace(tresc)); roznica != "" {
		wynik.Diff = &roznica
	}
	if pierwowzor.OpoznienieMs != nil {
		delta := opoznienie - *pierwowzor.OpoznienieMs
		wynik.LatencyDeltaMs = &delta
	}
	return shared.ProvenanceCallReplayResponse{Replay: wynik}, nil
}

// zasiegiPowtorzenia odtwarza kontekst pierwowzoru, żeby powtórzenie poszło
// tam, gdzie poszło oryginalne wywołanie. Zasięg wzięty z powietrza dałby
// wywołanie z innymi parametrami wykonania, czyli nieporównywalne.
func zasiegiPowtorzenia(pierwowzor dane.WywolanieModelu) models.Zasiegi {
	return models.Zasiegi{
		Srodowisko: wartoscTekstu(pierwowzor.Srodowisko),
		Projekt:    wartoscTekstu(pierwowzor.ProjektKod),
		Sesja:      wartoscTekstu(pierwowzor.SesjaKod),
		Okno:       wartoscTekstu(pierwowzor.OknoKod),
	}
}

// modelPowtorzenia rozstrzyga model zapisywany w śladzie powtórzenia: wskazany
// w żądaniu, a w jego braku model pierwowzoru.
func modelPowtorzenia(zadany *string, pierwowzor dane.WywolanieModelu) *string {
	if model := strings.TrimSpace(wartoscTekstu(zadany)); model != "" {
		kopia := model
		return &kopia
	}
	return pierwowzor.Model
}

// roznicaOdpowiedzi opisuje, czym odpowiedź powtórzenia różni się od pierwowzoru: zestawieniem długości i wspólnego przedrostka, nie różnicą wierszową, bo odpowiedzi bywają jednym akapitem.
func roznicaOdpowiedzi(pierwotna, powtorzona string) string {
	if pierwotna == powtorzona {
		return ""
	}
	if pierwotna == "" {
		return "pierwowzór nie ma zapisanej odpowiedzi; powtórzenie oddało " +
			zapisLiczbyMiary(float64(len([]rune(powtorzona)))) + " znaków"
	}
	wspolny := 0
	pierwsze := []rune(pierwotna)
	drugie := []rune(powtorzona)
	for wspolny < len(pierwsze) && wspolny < len(drugie) && pierwsze[wspolny] == drugie[wspolny] {
		wspolny++
	}
	return "odpowiedzi rozchodzą się od znaku " + zapisLiczbyMiary(float64(wspolny)) +
		"; pierwowzór ma " + zapisLiczbyMiary(float64(len(pierwsze))) +
		" znaków, powtórzenie " + zapisLiczbyMiary(float64(len(drugie)))
}
