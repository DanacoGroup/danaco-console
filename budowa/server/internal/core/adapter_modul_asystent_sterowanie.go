// Odpowiedzialność pliku: doprowadzenie do tury zlecenia asystenta tego, bez
// czego asystent nie jest klawiaturą Operatora, tylko drugim rozmówcą —
// konfiguracji MCP okna wraz z wpisem serwera narzędzi modelu.
//
// `zapytanieZlecenia` (adapter_modul_asystent_wykonawca.go) składa wywołanie
// kanału z pól okna — kanał, katalogi, środowisko, tryb, rola — i nie dokłada
// `KonfiguracjaMCP`. Rozmowa dokłada ją zawsze (`adapter_rozmowa_srodowisko.go`
// → `mosty.tekstZNarzedziami`), bo dopiero wpis `danaco` w pliku `--mcp-config`
// mówi procesowi modelu, że narzędzia sterowania platformą istnieją. Tura
// zlecenia bez tego wpisu nie miałaby ani jednego narzędzia — czym otworzyć
// okna ani czym wpisać do niego promptu — i zostawałaby przy samej odpowiedzi
// tekstem.
//
// Stąd ta sama tura, ten sam rejestr kanałów i ten sam składacz mostów co
// w rozmowie; drugiej drogi do narzędzi nie ma. Model zlecenia dostaje wykaz
// narzędzi kontraktu, więc ciąg `window.create` → `message.send` w oknie
// docelowym wykonuje się przez rdzeń, komendami kontraktu.
//
// Bez wpiętego składacza mostów tura idzie samym tekstem. Bez binarium serwera
// narzędzi odmawia `most_narzedzi.go`, meldując powód do dziennika rdzenia; ten
// plik tego nie powtarza ani nie obchodzi.
//
// Ślad jest skutkiem ubocznym drogi, nie dopiskiem: każde posunięcie asystenta
// idzie komendą kontraktu przez rdzeń, a rdzeń rozgłasza je tak samo jak
// posunięcie Operatora — `window.create` → `window.changed`, `message.send` →
// `message.changed` + `stream.chunk` + `progress.changed`, `session.focus` →
// `session.focus.changed`. Ekran Operatora dostaje komplet zdarzeń, choć akcje
// szły innym połączeniem.
//
// Ślad nie niesie sprawcy: ani `Window`, ani `Message`, ani `Session`, ani
// `AssistantAction` nie mają takiego pola, a żadne zdarzenie nie niesie
// identyfikatora połączenia, które czynność wywołało. Ekran Operatora widzi
// więc, że okno się otworzyło i że prompt poszedł, ale nie odróżni ruchu
// asystenta od własnego. Dopisanie pola byłoby zmianą zamrożonego kontraktu.
package core

import (
	"context"
	"encoding/json"

	"danacoconsole/server/internal/models"
	"danacoconsole/server/internal/narzedzia"
	"danacoconsole/server/internal/session"
)

// ZMostami wpina składacz konfiguracji mostów MCP okna — ten sam, którym jedzie
// rozmowa. Wpięcie osobnego składacza dałoby asystentowi drugi wykaz narzędzi
// i drugą prawdę o tym, czym model steruje platformą.
//
// Strażnik na nil jest tu z tego samego powodu co przy `ZMowa`: montaż podaje
// wskaźnik, a wskaźnik nil wpięty bez sprawdzenia zamieniłby brak składacza
// w panikę przy pierwszym zleceniu zamiast w pracę w zakresie niepełnym.
func (a *adapterAsystenta) ZMostami(m *mostyOkna) *adapterAsystenta {
	if m == nil {
		return a
	}
	a.mosty = m
	return a
}

// uzupelnijNarzedzia dokłada do zapytania zlecenia konfigurację MCP okna wraz
// z wpisem serwera narzędzi modelu.
//
// Osobno od `zapytanieZlecenia`, bo tamta funkcja jest czysta — bierze okno
// i polecenie, nic nie czyta i niczego nie woła. Wyliczenie konfiguracji sięga
// do repozytoriów nadań i punktów dostępu, więc potrzebuje kontekstu; ta sama
// granica przebiega w rozmowie (`zapytanieKanalu` obok `uzupelnijSrodowisko`).
//
// Wartość pusta nie nadpisuje niczego: „okno bez nadań i bez narzędzi" znaczy
// tu turę bez przełącznika `--mcp-config`, a nie wyczyszczenie czegoś, co ktoś
// wcześniej ustalił.
func (a *adapterAsystenta) uzupelnijNarzedzia(ctx context.Context, okno session.Okno,
	zapytanie *models.Zapytanie) {

	if a == nil || a.mosty == nil || zapytanie == nil || okno.Id == "" {
		return
	}
	tekst := a.mosty.tekstZNarzedziami(ctx, okno.Id)
	if tekst == "" {
		return
	}
	zapytanie.KonfiguracjaMCP = zZasiegiemKlawiatury(tekst, okno)
}

// modulAsystenta jest identyfikatorem modułu, którego okno jest klawiaturą
// Operatora. Wartość jest wartością katalogu modułów, tą samą, którą oddaje
// `window.list` w polu `moduleId`.
const modulAsystenta = "assistant"

// zZasiegiemKlawiatury dopisuje do wpisu `danaco` rolę okna asystenta.
//
// Asystent ustawia konfigurację zlecenia i wybiera model. Obie te czynności to
// komendy, które wykaz narzędzi kontraktu pomija (`config.session.set`,
// `model.channel.set`), bo model nie przestawia sobie własnego wyposażenia. Dla
// okna roboczego ten zakaz zostaje nietknięty. Okno asystenta jest innym bytem:
// nie pracuje nad zadaniem, tylko nastawia okno docelowe, w którym pracować
// będzie model docelowy — tak jak robiłby to Operator ręką na klawiaturze.
//
// Rozszerzenie wchodzi wyłącznie wtedy, gdy okno tury jest oknem modułu
// Assistant wedle wiersza okna w rdzeniu, i wchodzi do wpisu MCP, czyli zanim
// proces modelu wystartuje. Nie ma komendy, którą model poprosiłby o szerszy
// zasięg, ani pola żądania, które by go niosło. Okno cudze — także okno
// docelowe, do którego asystent zaraz napisze — dostaje zasięg zwykły.
//
// Konfigurację składa `mostyOkna.tekstZNarzedziami` — ta sama i jedyna droga,
// którą jedzie rozmowa. Ten kod jej nie powtarza: bierze jej wynik
// i dokłada do gotowego wpisu argument roli. Wynik nieczytelny albo wpis
// `danaco` nieobecny (brak binarium serwera narzędzi — `most_narzedzi.go`
// melduje powód) zostawia tekst nietknięty: tura idzie w zasięgu, jaki jest,
// zamiast paść na składaniu konfiguracji.
func zZasiegiemKlawiatury(tekst string, okno session.Okno) string {
	if okno.Modul != modulAsystenta {
		return tekst
	}
	argumenty := narzedzia.ArgumentyZasiegu(narzedzia.ZasiegKlawiatury)
	if len(argumenty) == 0 {
		return tekst
	}
	var konfiguracja KonfiguracjaMostu
	if err := json.Unmarshal([]byte(tekst), &konfiguracja); err != nil {
		return tekst
	}
	wpis, jest := konfiguracja.McpServers[narzedzia.KluczWpisu]
	if !jest {
		return tekst
	}
	wpis.Args = append(append([]string{}, wpis.Args...), argumenty...)
	konfiguracja.McpServers[narzedzia.KluczWpisu] = wpis
	tresc, err := json.MarshalIndent(konfiguracja, "", "  ")
	if err != nil {
		return tekst
	}
	return string(tresc)
}
