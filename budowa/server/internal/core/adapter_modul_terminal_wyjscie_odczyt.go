// Komenda `terminal.output.read` — jednorazowy odczyt wyjścia jednego procesu
// terminala wraz z jego stanem i kodem wyjścia.
//
// Czynność jest osobną komendą, a nie polami w `terminal.command.exec`, bo
// tamta kończy się w chwili, gdy proces ruszy, a nie gdy się skończy
// (`adapter_modul_terminal_bieg.go`). Dołożenie do jej wyniku pól
// `stdout`/`stderr`/`exitCode` znaczyłoby jedno z dwojga:
//
//   - odpowiedź czeka na koniec procesu — wtedy `go build ./...` trzyma żądanie
//     gniazda kilka minut, a zerwane gniazdo zabiera wynik wykonanej pracy;
//   - odpowiedź wraca od razu z polami pustymi — czyli rdzeń zgłasza brak
//     wyjścia dla polecenia, które dopiero zaczęło pisać.
//
// Rozdzielenie czynności rozdziela też ich czasy: uruchomienie jest
// natychmiastowe, odczyt następuje wtedy, kiedy jest co czytać.
//
// Źródłem jest ten sam dziennik zbiorczego wyjścia, na którym stoi
// `terminal.output.stream` (`adapter_modul_terminal_wyjscie_dziennik.go`); ta
// komenda zawęża go do jednego procesu i składa wiersze z powrotem w tekst.
// Dziennik jest pierścieniem w pamięci, bo schemat bazy nie ma tabeli wyjścia —
// po ponownym uruchomieniu rdzenia proces nie występuje już w rejestrze
// i odczyt kończy się odmową `not_found`.
package core

import (
	"context"
	"strings"
	"time"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// granicaWyjsciaKomendy jest górną granicą jednego strumienia (osobno stdout,
// osobno stderr) oddawanego w odpowiedzi na `terminal.output.read`.
//
// 65 536 bajtów to około 16 000 znaczników — ułamek okna kontekstu modelu,
// a zarazem wielokrotność wyników, dla których ta komenda powstała: pełny
// przebieg `go build ./...` z błędami mieści się w kilku kilobajtach, `go test
// ./...` całego drzewa w kilkudziesięciu.
//
// Przycięcie nie jest ciche: wychodzi w polach `truncated` i `truncatedBytes`.
const granicaWyjsciaKomendy = 64 * 1024

// granicaCzekaniaOdczytu jest górną granicą pola `waitMs`. Ogranicza żądanie,
// nie proces: po upływie czekania proces biegnie dalej, a odczyt oddaje wyjście
// dotychczasowe wraz ze stanem `running`. Wskazanie większe schodzi do tej
// wartości zamiast kończyć się odmową.
const granicaCzekaniaOdczytu = 60 * time.Second

// OdczytajWyjscie obsługuje `terminal.output.read`.
func (a *adapterWyjsciaTerminala) OdczytajWyjscie(ctx context.Context,
	z shared.TerminalOutputReadRequest) (shared.TerminalOutputReadResponse, error) {

	if a == nil || a.dziennik == nil {
		return shared.TerminalOutputReadResponse{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeInternalError,
			"moduł Terminal: rdzeń nie ma nadajnika wyjścia, więc nie prowadzi dziennika wyjścia "+
				"procesów — wyjścia polecenia nie ma skąd wziąć"))
	}
	kod := strings.TrimSpace(z.ProcessId)
	if kod == "" {
		return shared.TerminalOutputReadResponse{}, bladZadaniaTerminala(
			"odczyt wyjścia wymaga wskazania procesu; identyfikator wraca z terminal.command.exec")
	}
	proces, jest := a.rejestr.Proces(kod)
	if !jest {
		return shared.TerminalOutputReadResponse{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeNotFound,
			"moduł Terminal: proces "+kod+" nie występuje w rejestrze rdzenia — "+
				"albo nigdy nie ruszył, albo wypadł z historii (pojemność 500 przebiegów), "+
				"albo rdzeń był od tego czasu uruchomiony ponownie"))
	}

	// Czekanie idzie przed odczytem, żeby wiersze dopisane w jego trakcie
	// weszły do odpowiedzi.
	if czekanie := czekanieOdczytu(z.WaitMs); czekanie > 0 {
		proces.CzekajNaKoniec(czekanie)
	}

	stan, kodWyjscia, _ := proces.Migawka()
	zwykle, diagnostyka := a.dziennik.WyjscieProcesu(kod, ogonOdczytu(z.Tail))
	odpowiedz := shared.TerminalOutputReadResponse{Status: stan}
	odpowiedz.Stdout, odpowiedz.Stderr, odpowiedz.Truncated, odpowiedz.TruncatedBytes =
		zlozStrumienie(zwykle, diagnostyka)
	// Kod wyjścia wchodzi do odpowiedzi wyłącznie wtedy, gdy proces go ma. Zero
	// wpisane przy procesie biegnącym albo ubitym sygnałem oznaczałoby udane
	// zakończenie polecenia, które się jeszcze nie skończyło.
	if kodWyjscia != nil {
		odpowiedz.ExitCode = kodWyjscia
	}
	return odpowiedz, nil
}

// zlozStrumienie skleja wiersze obu strumieni w tekst i przycina każdy z nich
// do granicy rozmiaru.
//
// Przycięcie idzie od początku tekstu, bo ostatnie bajty niosą podsumowanie
// budowania, ostatni błąd i wiersz zamykający. Tak samo działa sam dziennik,
// z którego wiersze pochodzą — pierścień wypycha najstarsze.
func zlozStrumienie(zwykle, diagnostyka []string) (string, string, bool, *int) {
	stdout, uciete1 := przytnijDoGranicy(strings.Join(zwykle, "\n"))
	stderr, uciete2 := przytnijDoGranicy(strings.Join(diagnostyka, "\n"))
	razem := uciete1 + uciete2
	if razem == 0 {
		return stdout, stderr, false, nil
	}
	odciete := razem
	return stdout, stderr, true, &odciete
}

// przytnijDoGranicy zostawia ostatnie `granicaWyjsciaKomendy` bajtów tekstu
// i mówi, ile bajtów odcięto z początku.
func przytnijDoGranicy(tekst string) (string, int) {
	if len(tekst) <= granicaWyjsciaKomendy {
		return tekst, 0
	}
	odciete := len(tekst) - granicaWyjsciaKomendy
	return tekst[odciete:], odciete
}

// ogonOdczytu czyta liczbę wierszy z żądania. Brak wskazania znaczy komplet
// zapamiętanych wierszy — inaczej niż w `terminal.output.stream`, bo odczyt
// dotyczy wyniku polecenia, a nie podglądu na żywo. Zero daje sam stan procesu
// bez treści, czyli tanie zapytanie o zakończenie.
func ogonOdczytu(ile *int) int {
	if ile == nil {
		return pojemnoscDziennikaWyjscia
	}
	if *ile < 0 {
		return 0
	}
	if *ile > pojemnoscDziennikaWyjscia {
		return pojemnoscDziennikaWyjscia
	}
	return *ile
}

// czekanieOdczytu czyta granicę czekania z żądania. Wartość niedodatnia znaczy
// odczyt natychmiastowy; wartość większa od granicy schodzi do granicy.
func czekanieOdczytu(milisekundy *int) time.Duration {
	if milisekundy == nil || *milisekundy <= 0 {
		return 0
	}
	czekanie := time.Duration(*milisekundy) * time.Millisecond
	if czekanie > granicaCzekaniaOdczytu {
		return granicaCzekaniaOdczytu
	}
	return czekanie
}

// WyjscieProcesu zwraca ostatnie `ile` wierszy jednego procesu, rozdzielone na
// wyjście zwykłe i diagnostyczne, w kolejności wypisania. Na tym stoi komenda
// `terminal.output.read`.
//
// Zawężenie idzie po procesie, nie po karcie: karta prowadzi wiele poleceń po
// kolei, a pytanie dotyczy wyniku jednego z nich. Pole `processId` wiersza jest
// w kontrakcie opcjonalne, więc wiersz bez procesu nie pasuje do żadnego
// zawężenia i nie wchodzi do wyniku.
//
// Liczba `ile` dotyczy każdego strumienia z osobna. Wspólny licznik pozwoliłby
// wyjściu zwykłemu wypchnąć z wyniku komunikaty błędów.
func (d *dziennikWyjscia) WyjscieProcesu(kodProcesu string, ile int) ([]string, []string) {
	zwykle := make([]string, 0, 16)
	diagnostyka := make([]string, 0, 16)
	if d == nil || ile <= 0 || kodProcesu == "" {
		return zwykle, diagnostyka
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	// Przebieg od końca bierze ostatnie wiersze; odwrócenie wyniku przywraca
	// kolejność wypisania.
	for i := len(d.wpisy) - 1; i >= 0 && (len(zwykle) < ile || len(diagnostyka) < ile); i-- {
		wiersz := d.wpisy[i].wiersz
		if wiersz.ProcessId == nil || *wiersz.ProcessId != kodProcesu {
			continue
		}
		if wiersz.Channel == shared.TerminalOutputChannelStderr {
			if len(diagnostyka) < ile {
				diagnostyka = append(diagnostyka, wiersz.Text)
			}
			continue
		}
		if len(zwykle) < ile {
			zwykle = append(zwykle, wiersz.Text)
		}
	}
	return odwroc(zwykle), odwroc(diagnostyka)
}

// odwroc odwraca wykaz wierszy w miejscu i oddaje go z powrotem.
func odwroc(wiersze []string) []string {
	for lewy, prawy := 0, len(wiersze)-1; lewy < prawy; lewy, prawy = lewy+1, prawy-1 {
		wiersze[lewy], wiersze[prawy] = wiersze[prawy], wiersze[lewy]
	}
	return wiersze
}
