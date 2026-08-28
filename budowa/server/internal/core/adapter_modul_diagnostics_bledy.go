// Odpowiedzialność pliku: błędy modułu Diagnostics — przyjęcie odmowy
// wykonania komendy jako faktu i obsługa komendy diagnostics.error.list
// zasilającej okno Errors Panel.
package core

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strings"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// kontekstBledu jest treścią pola DiagnosticError.context. Kod błędu powtarza
// się tutaj obok pola DiagnosticError.errorCode, bo kontekst opisuje jedno
// wystąpienie odmowy, a pole błędu niesie kod całej grupy spiętej odciskiem.
type kontekstBledu struct {
	ErrorCode string `json:"errorCode"`
	Command   string `json:"command"`
	Retryable bool   `json:"retryable"`
	SessionId string `json:"sessionId,omitempty"`
	WindowId  string `json:"windowId,omitempty"`
}

// ZapiszNiepowodzenie przyjmuje odmowę wykonania komendy. Czynność jest
// bezzwrotna: dyspozytor odpowiada klientowi i nie czeka, aż diagnostyka
// dopisze wiersz.
func (a *adapterDiagnostyki) ZapiszNiepowodzenie(ctx context.Context, n NiepowodzenieKomendy) {
	if a == nil || a.repozytorium == nil || n.Komenda == "" {
		return
	}
	komenda := string(n.Komenda)
	tresc := strings.TrimSpace(n.Blad.Message)
	if tresc == "" {
		tresc = "rdzeń odmówił wykonania bez podania przyczyny"
	}
	chwila := time.Now().UnixMilli()

	// Wpis dziennika powstaje zawsze, nawet gdy zapis błędu się nie powiedzie.
	poziom := poziomKoduBledu(n.Blad.Code)
	a.zakolejkuj(dane.WpisDiagnostyki{
		Kod:      nowyIdentyfikator(przedrostekWpisuDziennika),
		Chwila:   chwila,
		Poziom:   poziom,
		Zrodlo:   wskaznikTekstu(komenda),
		Tresc:    string(n.Blad.Code) + ": " + tresc,
		SesjaKod: wskaznikNiepusty(n.IdSesji),
		OknoKod:  wskaznikNiepusty(n.IdOkna),
		Odcisk:   odciskWpisu(poziom, komenda, string(n.Blad.Code)+": "+tresc),
	})

	kontekst, err := json.Marshal(kontekstBledu{
		ErrorCode: string(n.Blad.Code), Command: komenda, Retryable: n.Blad.Retryable,
		SessionId: n.IdSesji, WindowId: n.IdOkna,
	})
	if err != nil {
		// Kontekst nieserializowalny nie ma prawa zabrać ze sobą błędu.
		kontekst = nil
	}

	blad := dane.BladDiagnostyczny{
		Kod:       nowyIdentyfikator(przedrostekBleduDiagnozy),
		Odcisk:    odciskBledu(komenda, string(n.Blad.Code), tresc),
		Tresc:     tresc,
		Zrodlo:    wskaznikTekstu(komenda),
		KodBledu:  n.Blad.Code,
		Stan:      shared.DiagnosticErrorStatusNew,
		Priorytet: priorytetKoduBledu(n.Blad.Code),
		Kontekst:  wskaznikNiepusty(string(kontekst)),
		Pierwsze:  chwila,
		Ostatnie:  chwila,
	}
	if _, err := a.repozytorium.ZapiszBlad(ctx, blad); err != nil {
		a.niezapisane.Add(1)
	}
}

// WykazBledow obsługuje diagnostics.error.list, zasilając okno Errors Panel wykazem błędów rdzenia platformy.
func (a *adapterDiagnostyki) WykazBledow(ctx context.Context,
	z shared.DiagnosticsErrorListRequest) (shared.DiagnosticsErrorListResponse, error) {

	if a.repozytorium == nil {
		return shared.DiagnosticsErrorListResponse{}, bladBrakuTrwalosci("wykaz błędów")
	}
	filtr := dane.FiltrBledow{
		Stan:      wartoscStanuBledu(z.Status),
		Priorytet: wartoscPriorytetu(z.Priority),
		Od:        wartoscChwili(z.FromTime),
		Do:        wartoscChwili(z.ToTime),
		Granica:   wartoscLiczby(z.Limit),
	}
	wiersze, err := a.repozytorium.Bledy(ctx, filtr)
	if err != nil {
		return shared.DiagnosticsErrorListResponse{}, bladDiagnostyki(err)
	}
	// total liczy się osobno, bez granicy: repozytorium ucina wykaz na 500 wierszach po cichu.
	razem, err := a.repozytorium.LiczbaBledow(ctx, filtr)
	if err != nil {
		return shared.DiagnosticsErrorListResponse{}, bladDiagnostyki(err)
	}
	return shared.DiagnosticsErrorListResponse{
		Errors: bledyKontraktu(wiersze),
		Total:  wskaznikLiczby(razem),
	}, nil
}

// odciskBledu grupuje wystąpienia tej samej odmowy tej samej komendy w jeden byt zbiorczy usterki rdzenia.
func odciskBledu(komenda, kod, tresc string) string {
	suma := sha256.Sum256([]byte(komenda + "\x00" + kod + "\x00" + tresc))
	return hex.EncodeToString(suma[:16])
}

// poziomKoduBledu przekłada kod kontraktu na poziom wpisu dziennika. Odmowa
// z powodu braku uprawnienia albo braku bytu idzie poziomem warn, usterka
// rdzenia i niedostępność kanału idą poziomem error.
func poziomKoduBledu(kod shared.ErrorCode) shared.LogLevel {
	switch kod {
	case shared.ErrorCodeInternalError, shared.ErrorCodeChannelUnavailable:
		return shared.LogLevelError
	default:
		return shared.LogLevelWarn
	}
}

// priorytetKoduBledu nadaje priorytet początkowy wynikający z kodu kontraktu.
// Jest to wyłącznie punkt wyjścia — priorytet ostateczny nadaje Operator.
func priorytetKoduBledu(kod shared.ErrorCode) shared.DiagnosticPriority {
	switch kod {
	case shared.ErrorCodeInternalError:
		return shared.DiagnosticPriorityHigh
	case shared.ErrorCodeChannelUnavailable, shared.ErrorCodeRateLimited:
		return shared.DiagnosticPriorityMedium
	default:
		return shared.DiagnosticPriorityLow
	}
}
