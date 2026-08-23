// Odpowiedzialność pliku: dołożenie wpisu `danaco` — serwera narzędzi modelu —
// do konfiguracji MCP okna rozmowy.
//
// Tędy model dostaje sterowanie platformą. Kontrakt deklaruje narzędzia, pakiet
// `narzedzia` je wystawia, ale bez wpisu w pliku `--mcp-config` proces modelu
// nigdy się o nich nie dowie. Wpis powstaje osobno dla każdego okna i niesie
// jego identyfikator, więc narzędzie zawsze wie, z którego okna przyszło
// wywołanie.
//
// Konfigurację MCP zapytania składa `adapterRozmowy.uzupelnijSrodowisko`
// (`adapter_rozmowa_srodowisko.go`) i woła stamtąd
// `a.mosty.tekstZNarzedziami(ctx, okno.Id)`. Reszta drogi (wykaz, rozdzielnia,
// gniazdo, protokół) stoi w pakiecie `narzedzia`.
//
// Brak binarium serwera narzędzi jest meldowany wprost, raz na powód: bez tego
// model dostałby wpis wskazujący plik, którego nie ma, i sterowanie platformą
// milczałoby bez śladu. Meldunek co turę byłby hałasem, w którym sam by ginął.
package core

import (
	"context"
	"encoding/json"

	"danacoconsole/server/internal/narzedzia"
)

// tekstZNarzedziami zwraca konfigurację MCP okna powiększoną o wpis serwera
// narzędzi.
//
// Każdy krok przepuszcza dalej zamiast przerywać. Okno bez nadań dostaje
// konfigurację z samym wpisem narzędzi — sterowanie platformą nie zależy od tego, czy okno
// ma wgląd w jakąkolwiek maszynę. Konfiguracja nieczytelna albo okno bez
// identyfikatora zostawiają tekst dotychczasowy: rozmowa toczy się dalej, tyle
// że bez narzędzi.
func (m *mostyOkna) tekstZNarzedziami(ctx context.Context, idOkna string) string {
	dotychczasowy := m.Tekst(ctx, idOkna)
	polecenie, argumenty, powod, jest := narzedzia.Wpis(idOkna)
	if !jest {
		m.zglosBrakNarzedzi(idOkna, powod)
		return dotychczasowy
	}
	konfiguracja := KonfiguracjaMostu{McpServers: map[string]WpisMostu{}}
	if dotychczasowy != "" {
		if err := json.Unmarshal([]byte(dotychczasowy), &konfiguracja); err != nil {
			return dotychczasowy
		}
	}
	konfiguracja.McpServers[narzedzia.KluczWpisu] = WpisMostu{
		Type:    RodzajWpisuMostu,
		Command: polecenie,
		Args:    argumenty,
	}
	tresc, err := json.MarshalIndent(konfiguracja, "", "  ")
	if err != nil {
		return dotychczasowy
	}
	return string(tresc)
}

// zglosBrakNarzedzi wpisuje do dziennika rdzenia powód, dla którego okno
// prowadzi turę bez sterowania platformą.
//
// Raz na powód, nie raz na turę: powód nie zmienia się między turami — albo
// binarium w produkcie jest, albo go nie ma — a tur bywa kilkaset dziennie;
// powtarzany meldunek zasłoniłby resztę dziennika i sam przestałby być czytany.
// Okno bez identyfikatora meldunku nie daje: to stan zwykły przy oknie jeszcze
// niezałożonym, a nie usterka wydania.
func (m *mostyOkna) zglosBrakNarzedzi(idOkna, powod string) {
	if m == nil || m.dziennik == nil || powod == "" || idOkna == "" {
		return
	}
	if _, juzBylo := m.zgloszoneBraki.LoadOrStore(powod, true); juzBylo {
		return
	}
	m.dziennik.Printf("serwer narzędzi modelu poza turą: %s — wpisu %q nie ma w konfiguracji MCP,"+
		" rozmowa toczy się bez sterowania platformą", powod, narzedzia.KluczWpisu)
}
