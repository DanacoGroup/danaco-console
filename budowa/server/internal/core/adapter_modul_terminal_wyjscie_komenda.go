// Komenda `terminal.output.stream` — zapis okna na zbiorcze wyjście wszystkich
// otwartych kart terminala wraz z ogonem historii.
//
// Typ poniżej rozszerza port zamiast zakładać drugi adapter: osadza adapter
// modułu Terminal, więc niesie komplet jego komend i jest tym samym bytem,
// którym idą karty i procesy. Dziennik zbiorczy podpina się przy montażu,
// owijając nadajnik istniejącej pompy (`ZDziennikiemWyjscia`), a nie zakładając
// drugiej pompy.
//
// Odpowiedź niesie dwa niezależne pola: `subscribed` mówi o zapisie okna na
// strumień, `lines` o ogonie historii. Żądanie bez `windowId` jest w kontrakcie
// dopuszczone i znaczy sam odczyt ogona — odpowiedź niesie wtedy
// `subscribed: false` wraz z wierszami.
package core

import (
	"context"
	"errors"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// Zgodność rozszerzonego adaptera z rozszerzonym portem sprawdza kompilator.
var _ WyjscieTerminala = (*adapterWyjsciaTerminala)(nil)

// adapterWyjsciaTerminala wypełnia port WyjscieTerminala: adapter modułu
// Terminal rozszerzony o zbiorcze wyjście.
type adapterWyjsciaTerminala struct {
	*adapterTerminala
	// dziennik jest niezerowy, gdy pompa wyjścia ma nadajnik, który dało się
	// owinąć. Zerowy powoduje odmowę komendy z kodem `internal_error`.
	dziennik *dziennikWyjscia
}

// ZDziennikiemWyjscia rozszerza adapter modułu Terminal o rodzinę zbiorczego
// wyjścia. Owija nadajnik pompy, więc dotychczasowa droga fragmentów zostaje
// nietknięta, a dziennik dostaje ich kopię.
//
// Wywołanie jest bezpieczne przy powtórzeniu: nadajnik owinięty raz nie owija
// się drugi raz, bo każde owinięcie dokładałoby kolejną kopię tych samych
// wierszy do tego samego dziennika.
func (a *adapterTerminala) ZDziennikiemWyjscia() *adapterWyjsciaTerminala {
	if a == nil {
		return nil
	}
	if a.wyjscie == nil || a.wyjscie.nadajnik == nil {
		// Moduł bez nadajnika pracuje dalej — karty i procesy działają — lecz
		// zbiorczego wyjścia nie ma z czego złożyć. Komenda powie to wprost.
		return &adapterWyjsciaTerminala{adapterTerminala: a}
	}
	if opakowany, juz := a.wyjscie.nadajnik.(*nadajnikZDziennikiem); juz {
		return &adapterWyjsciaTerminala{adapterTerminala: a, dziennik: opakowany.dziennik}
	}
	dziennik := nowyDziennikWyjscia(a.rejestr)
	dziennik.dalej = a.wyjscie.nadajnik
	a.wyjscie.nadajnik = &nadajnikZDziennikiem{dalej: dziennik.dalej, dziennik: dziennik}
	return &adapterWyjsciaTerminala{adapterTerminala: a, dziennik: dziennik}
}

// StrumienWyjscia obsługuje `terminal.output.stream`.
func (a *adapterWyjsciaTerminala) StrumienWyjscia(ctx context.Context,
	z shared.TerminalOutputStreamRequest) (shared.TerminalOutputStreamResponse, error) {

	if a == nil || a.dziennik == nil {
		return shared.TerminalOutputStreamResponse{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeInternalError,
			"moduł Terminal: rdzeń nie ma nadajnika wyjścia, więc nie prowadzi zbiorczego strumienia kart"))
	}

	filtr, err := a.filtrZadania(ctx, z)
	if err != nil {
		return shared.TerminalOutputStreamResponse{}, err
	}
	ile, err := ogonZadania(z.Tail)
	if err != nil {
		return shared.TerminalOutputStreamResponse{}, err
	}

	zapisane := false
	if oknoKod := strings.TrimSpace(wartoscTekstu(z.WindowId)); oknoKod != "" {
		// Okno zamknięte nie odbierze niczego, więc zapis na strumień kończy
		// się odmową zamiast milczącego przyjęcia.
		okno, err := a.oknoWykonania(oknoKod)
		if err != nil {
			return shared.TerminalOutputStreamResponse{}, err
		}
		a.dziennik.Obserwuj(oknoKod, okno.IdSesji, filtr)
		zapisane = true
	}
	a.odsiejZamknieteOkna()

	return shared.TerminalOutputStreamResponse{
		Lines:      a.dziennik.Ogon(filtr, ile),
		Subscribed: zapisane,
	}, nil
}

// filtrZadania składa zawężenie strumienia z dwóch wskazań żądania: karty sesji
// i wykazu kart terminala. Oba wskazania działają łącznie: żądanie niosące
// kartę sesji i wykaz kart dostaje wiersze spełniające oba warunki naraz.
func (a *adapterWyjsciaTerminala) filtrZadania(ctx context.Context,
	z shared.TerminalOutputStreamRequest) (filtrWyjscia, error) {

	filtr := filtrWyjscia{}
	if idSesji := strings.TrimSpace(wartoscTekstu(z.SessionId)); idSesji != "" {
		if a.okna == nil {
			return filtrWyjscia{}, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
				"moduł Terminal: rdzeń nie ma rejestru sesji, więc nie sprawdzi karty sesji "+idSesji))
		}
		if _, err := a.okna.Sesja(idSesji); err != nil {
			return filtrWyjscia{}, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
				"moduł Terminal: karta sesji "+idSesji+" nie jest otwarta"))
		}
		filtr.idSesji = idSesji
	}

	if len(z.TerminalSessionIds) == 0 {
		return filtr, nil
	}
	filtr.karty = make(map[string]struct{}, len(z.TerminalSessionIds))
	for _, wskazana := range z.TerminalSessionIds {
		kod := strings.TrimSpace(wskazana)
		if kod == "" {
			return filtrWyjscia{}, bladZadaniaTerminala(
				"wykaz kart terminala niesie pozycję pustą; pusty wykaz znaczy komplet, pusta pozycja nie znaczy nic")
		}
		if err := a.sprawdzKarte(ctx, kod); err != nil {
			return filtrWyjscia{}, err
		}
		filtr.karty[kod] = struct{}{}
	}
	return filtr, nil
}

// sprawdzKarte upewnia się, że wskazana karta terminala w ogóle istnieje.
// Szuka najpierw w rejestrze żywym, potem w tabeli kart terminala: karta
// zamknięta ma historię wyjścia, więc jej wskazanie nie jest błędem.
func (a *adapterWyjsciaTerminala) sprawdzKarte(ctx context.Context, kod string) error {
	if _, jest := a.rejestr.Karta(kod); jest {
		return nil
	}
	if a.repozytorium != nil {
		_, err := a.repozytorium.Karta(ctx, kod)
		if err == nil {
			return nil
		}
		if !errors.Is(err, dane.ErrBrakWiersza) {
			return err
		}
	}
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
		"moduł Terminal: karta "+kod+" nie występuje ani w rejestrze rdzenia, ani w dzienniku kart"))
}

// odsiejZamknieteOkna wykreśla obserwacje okien, których rejestr już nie zna.
// Kontrakt nie ma komendy wypisania się ze strumienia, a okno zamknięte nie
// przestaje być obserwatorem samo z siebie, więc porządkowanie odbywa się przy
// każdym kolejnym żądaniu tej komendy.
func (a *adapterWyjsciaTerminala) odsiejZamknieteOkna() {
	if a.okna == nil {
		return
	}
	a.dziennik.Zapomnij(func(oknoKod string) bool {
		_, err := a.okna.Okno(oknoKod)
		return err == nil
	})
}

// ogonZadania czyta liczbę wierszy ogona. Brak wskazania daje ogon domyślny,
// zero — ogon pusty (żądanie samego zapisu na strumień), a liczba większa od
// pojemności dziennika schodzi do tej pojemności.
func ogonZadania(ile *int) (int, error) {
	if ile == nil {
		return domyslnyOgonWyjscia, nil
	}
	if *ile < 0 {
		return 0, bladZadaniaTerminala("liczba wierszy ogona nie może być ujemna")
	}
	if *ile > pojemnoscDziennikaWyjscia {
		return pojemnoscDziennikaWyjscia, nil
	}
	return *ile, nil
}
