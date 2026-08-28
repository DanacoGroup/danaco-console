// Przekaz zdarzeń zaczepów do modułu Diagnostics jest śladem po zdarzeniu, nie oceną i nie
// blokadą. Każde zdarzenie zostawia wpis dziennika, a niepowodzenie zaczepu dodatkowo wiersz
// błędu.
package core

import (
	"encoding/json"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/injection"
	"danacoconsole/shared"
)

// zrodloZaczepu nazywa źródło wpisu diagnostyki dla zdarzenia zaczepu.
// Kształt „zaczep:<zdarzenie>" odróżnia je od źródeł-komend w tym samym oknie.
func zrodloZaczepu(e injection.ZdarzenieZaczepu) string {
	if e.Zdarzenie == "" {
		return "zaczep"
	}
	return "zaczep:" + e.Zdarzenie
}

// kontekstZaczepu jest treścią pola `DiagnosticError.context` dla niepowodzenia
// zaczepu — wzorem kontekstBledu odmowy komendy (adapter_modul_diagnostics_bledy.go).
type kontekstZaczepu struct {
	HookName  string `json:"hookName"`
	HookEvent string `json:"hookEvent"`
	Outcome   string `json:"outcome,omitempty"`
	ExitCode  *int   `json:"exitCode,omitempty"`
	WindowId  string `json:"windowId,omitempty"`
}

// przekazDiagnostyce oddaje zdarzenie zaczepu diagnostyce. Czynność jest
// bezzwrotna jak ZapiszNiepowodzenie: tura nie czeka na diagnostykę, a jej
// niepowodzenie nie zmienia przebiegu rozmowy.
func (o *zdarzeniaWykonawcze) przekazDiagnostyce(okno string, e injection.ZdarzenieZaczepu) {
	o.mu.RLock()
	d := o.diagnostyka
	o.mu.RUnlock()
	if d == nil {
		return
	}

	chwila := time.Now().UnixMilli()
	zrodlo := zrodloZaczepu(e)
	tresc := trescWpisuZaczepu(e)

	// Wpis dziennika idzie poziomem info dla pracy zwykłej i poziomem warn dla niepowodzenia
	// zaczepu.
	poziom := shared.LogLevel(shared.LogLevelInfo)
	if e.Niepowodzenie() {
		poziom = shared.LogLevelWarn
	}
	d.zakolejkuj(dane.WpisDiagnostyki{
		Kod:     nowyIdentyfikator(przedrostekWpisuDziennika),
		Chwila:  chwila,
		Poziom:  poziom,
		Zrodlo:  wskaznikTekstu(zrodlo),
		Tresc:   tresc,
		OknoKod: wskaznikNiepusty(okno),
		Odcisk:  odciskWpisu(poziom, zrodlo, tresc),
	})

	if !e.Niepowodzenie() {
		return
	}

	// Niepowodzenie zaczepu wchodzi do Errors Panel jako odmowa polityki zaczepu.
	kontekst, err := json.Marshal(kontekstZaczepu{
		HookName: e.Nazwa, HookEvent: e.Zdarzenie, Outcome: e.Wynik,
		ExitCode: e.KodWyjscia, WindowId: okno,
	})
	if err != nil {
		kontekst = nil
	}
	blad := dane.BladDiagnostyczny{
		Kod:       nowyIdentyfikator(przedrostekBleduDiagnozy),
		Odcisk:    odciskBledu(zrodlo, string(shared.ErrorCodePermissionDenied), tresc),
		Tresc:     tresc,
		Zrodlo:    wskaznikTekstu(zrodlo),
		KodBledu:  shared.ErrorCodePermissionDenied,
		Stan:      shared.DiagnosticErrorStatusNew,
		Priorytet: priorytetKoduBledu(shared.ErrorCodePermissionDenied),
		Kontekst:  wskaznikNiepusty(string(kontekst)),
		Pierwsze:  chwila,
		Ostatnie:  chwila,
	}
	if _, err := d.repozytorium.ZapiszBlad(o.zycie, blad); err != nil {
		d.niezapisane.Add(1)
	}
}

// trescWpisuZaczepu opisuje zdarzenie zaczepu jednym zdaniem dziennika, złożonym z rodzaju
// zdarzenia i jego wyniku.
func trescWpisuZaczepu(e injection.ZdarzenieZaczepu) string {
	czesci := make([]string, 0, 4)
	switch {
	case e.Odpowiedz():
		czesci = append(czesci, "zaczep "+e.Nazwa+" odpowiedział")
	default:
		czesci = append(czesci, "zaczep "+e.Nazwa+" uruchomiony")
	}
	if e.Wynik != "" {
		czesci = append(czesci, "wynik: "+e.Wynik)
	}
	if e.KodWyjscia != nil && *e.KodWyjscia != 0 {
		czesci = append(czesci, "kod wyjścia niezerowy")
	}
	if blad := strings.TrimSpace(e.BladWyjscia); blad != "" {
		czesci = append(czesci, "stderr: "+blad)
	}
	return strings.Join(czesci, "; ")
}
