// Odpowiedzialność pliku: monitory zmian strony — `browser.monitor.add`,
// `.list`, `.check`, `.remove`.
//
// Monitor mierzy, a nie melduje. Założenie monitora POBIERA stronę i odkłada jej
// treść jako odniesienie; sprawdzenie pobiera ją ponownie i zestawia obie treści
// wiersz po wierszu. Monitor bez odniesienia oddawałby zawsze „bez zmian" albo
// zawsze „zmiana" — jedno i drugie jest meldunkiem bez pomiaru, a to wzorzec
// szkody, którego ten produkt już raz doświadczył.
//
// Różnica jest liczbą, nie wrażeniem: `BrowserContentDiff` niesie liczbę wierszy
// dodanych, usuniętych, numer pierwszego wiersza różnicy i przyrost znaków.
// Wszystkie cztery liczone są z dwóch treści, nie z niczego.
//
// Próg zmiany odsiewa drgania. Strona z zegarem albo licznikiem odwiedzin różni
// się przy każdym pobraniu; próg podany w znakach mówi, od jakiej różnicy zmiana
// jest zmianą. Bez progu każdy taki monitor alarmowałby co godzinę.
package core

import (
	"context"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// domyslnyInterwalMonitora jest cyklem sprawdzania przyjmowanym, gdy żądanie go
// nie podaje: godzina. Wartość domyślna należy do kodu, nie do schematu.
const domyslnyInterwalMonitora = 3600

// ZalozMonitor obsługuje `browser.monitor.add`. Odniesienie powstaje od razu:
// monitor bez treści początkowej nie miałby czego porównywać przy pierwszym
// sprawdzeniu i pierwsza jego odpowiedź byłaby zgadywaniem.
func (a *adapterPrzegladarki) ZalozMonitor(ctx context.Context,
	z shared.BrowserMonitorAddRequest) (shared.BrowserMonitorAddResponse, error) {

	if strings.TrimSpace(z.WindowId) == "" || strings.TrimSpace(z.Url) == "" {
		return shared.BrowserMonitorAddResponse{}, bladWskazaniaPrzegladarki(
			"komenda monitor.add bez okna albo adresu")
	}
	adres := strings.TrimSpace(z.Url)
	tresc, err := pobierzStrone(ctx, adres)
	if err != nil {
		return shared.BrowserMonitorAddResponse{}, bladPobraniaStrony(adres, err)
	}
	pilnowana := trescPilnowana(tresc.Tekst, z.Selector)
	odwolanie, err := a.zapiszTresc([]byte(pilnowana))
	if err != nil {
		return shared.BrowserMonitorAddResponse{}, err
	}
	dlugosc := int64(len(pilnowana))

	interwal := int64(wartoscLiczby(z.IntervalSeconds))
	if interwal <= 0 {
		interwal = domyslnyInterwalMonitora
	}
	wlaczony := z.Enabled == nil || *z.Enabled
	monitor := dane.MonitorPrzegladania{
		Kod:                  nowyIdentyfikator(przedrostekMonitora),
		Okno:                 z.WindowId,
		Url:                  adres,
		Selektor:             z.Selector,
		InterwalSekund:       interwal,
		Wlaczony:             wlaczony,
		Stan:                 string(shared.BrowserMonitorStatusPending),
		OdniesienieOdwolanie: &odwolanie,
		OdniesienieDlugosc:   &dlugosc,
	}
	if z.ChangeThreshold != nil {
		prog := int64(*z.ChangeThreshold)
		monitor.ProgZmiany = &prog
	}
	if z.NotifyChannel != nil {
		kanal := string(*z.NotifyChannel)
		monitor.KanalPowiadomienia = &kanal
	}
	zapisany, err := a.repozytorium.ZapiszMonitor(ctx, monitor)
	if err != nil {
		return shared.BrowserMonitorAddResponse{}, bladPrzegladarki(err)
	}
	return shared.BrowserMonitorAddResponse{Monitor: monitorKontraktu(zapisany)}, nil
}

// WykazMonitorow obsługuje `browser.monitor.list`.
func (a *adapterPrzegladarki) WykazMonitorow(ctx context.Context,
	z shared.BrowserMonitorListRequest) (shared.BrowserMonitorListResponse, error) {

	// Karta sesji zawęża tak samo jak okno: monitory modułu Browser stoją przy
	// oknie operacyjnym, a `sessionId` żądania jest drugą drogą wskazania tego
	// samego wykazu, nie drugim wykazem.
	okno := wartoscTekstuLubPusta(z.WindowId)
	if okno == "" {
		okno = wartoscTekstuLubPusta(z.SessionId)
	}
	tylkoWlaczone := z.EnabledOnly != nil && *z.EnabledOnly
	wiersze, err := a.repozytorium.Monitory(ctx, okno, tylkoWlaczone, wartoscLiczby(z.Limit))
	if err != nil {
		return shared.BrowserMonitorListResponse{}, bladPrzegladarki(err)
	}
	monitory := make([]shared.BrowserMonitor, 0, len(wiersze))
	for _, wiersz := range wiersze {
		monitory = append(monitory, monitorKontraktu(wiersz))
	}
	return shared.BrowserMonitorListResponse{Monitors: monitory}, nil
}

// SprawdzMonitor obsługuje `browser.monitor.check` — sprawdzenie natychmiastowe,
// poza cyklem, wraz z różnicą wobec odniesienia.
func (a *adapterPrzegladarki) SprawdzMonitor(ctx context.Context,
	z shared.BrowserMonitorCheckRequest) (shared.BrowserMonitorCheckResponse, error) {

	if strings.TrimSpace(z.MonitorId) == "" {
		return shared.BrowserMonitorCheckResponse{}, bladWskazaniaPrzegladarki(
			"komenda monitor.check bez monitora")
	}
	monitor, err := a.repozytorium.Monitor(ctx, z.MonitorId)
	if err != nil {
		return shared.BrowserMonitorCheckResponse{}, bladWierszaPrzegladania("monitora zmian", z.MonitorId, err)
	}

	tresc, err := pobierzStrone(ctx, monitor.Url)
	if err != nil {
		// Nieudane sprawdzenie jest sprawdzeniem: monitor odnotowuje stan
		// `failed` i czas próby, zamiast udawać, że niczego nie było.
		monitor.Stan = string(shared.BrowserMonitorStatusFailed)
		monitor.Sprawdzono = terazWBazie()
		if _, zapis := a.repozytorium.ZapiszMonitor(ctx, monitor); zapis != nil {
			return shared.BrowserMonitorCheckResponse{}, bladPrzegladarki(zapis)
		}
		return shared.BrowserMonitorCheckResponse{}, bladPobraniaStrony(monitor.Url, err)
	}

	biezaca := trescPilnowana(tresc.Tekst, monitor.Selektor)
	var odniesienie string
	if monitor.OdniesienieOdwolanie != nil {
		bajty, err := a.odczytajTresc(*monitor.OdniesienieOdwolanie)
		if err != nil {
			return shared.BrowserMonitorCheckResponse{}, err
		}
		odniesienie = string(bajty)
	}

	roznica := roznicaTresci(odniesienie, biezaca)
	prog := int64(0)
	if monitor.ProgZmiany != nil {
		prog = *monitor.ProgZmiany
	}
	zmiana := roznica.AddedLines+roznica.RemovedLines > 0 && bezwzglednaLiczba(int64(roznica.CharDelta)) >= prog

	monitor.Sprawdzono = terazWBazie()
	if zmiana {
		monitor.Stan = string(shared.BrowserMonitorStatusChanged)
		monitor.Zmieniono = terazWBazie()
	} else {
		monitor.Stan = string(shared.BrowserMonitorStatusUnchanged)
	}
	// Odniesienie przesuwa się tylko na żądanie. Bez tego pierwsze sprawdzenie
	// po zmianie zjadałoby ją: kolejne sprawdzenie porównywałoby nową treść
	// z nową treścią i meldowało spokój na stronie, która właśnie się zmieniła.
	if z.UpdateBaseline != nil && *z.UpdateBaseline {
		odwolanie, err := a.zapiszTresc([]byte(biezaca))
		if err != nil {
			return shared.BrowserMonitorCheckResponse{}, err
		}
		dlugosc := int64(len(biezaca))
		monitor.OdniesienieOdwolanie = &odwolanie
		monitor.OdniesienieDlugosc = &dlugosc
	}

	zapisany, err := a.repozytorium.ZapiszMonitor(ctx, monitor)
	if err != nil {
		return shared.BrowserMonitorCheckResponse{}, bladPrzegladarki(err)
	}
	return shared.BrowserMonitorCheckResponse{
		Monitor: monitorKontraktu(zapisany), Changed: zmiana, Diff: &roznica,
	}, nil
}

// ZdejmijMonitor obsługuje `browser.monitor.remove`.
func (a *adapterPrzegladarki) ZdejmijMonitor(ctx context.Context,
	z shared.BrowserMonitorRemoveRequest) (shared.BrowserMonitorRemoveResponse, error) {

	if strings.TrimSpace(z.MonitorId) == "" {
		return shared.BrowserMonitorRemoveResponse{}, bladWskazaniaPrzegladarki(
			"komenda monitor.remove bez monitora")
	}
	usuniety, err := a.repozytorium.UsunMonitor(ctx, z.MonitorId)
	if err != nil {
		return shared.BrowserMonitorRemoveResponse{}, bladPrzegladarki(err)
	}
	if !usuniety {
		return shared.BrowserMonitorRemoveResponse{}, bladNieznanegoBytu("monitora zmian", z.MonitorId)
	}
	return shared.BrowserMonitorRemoveResponse{Removed: true}, nil
}

// trescPilnowana zawęża treść strony do fragmentu wskazanego selektorem.
//
// Selektor jest tu wskazaniem TEKSTOWYM, nie selektorem CSS wykonywanym na
// drzewie: monitor pobiera stronę HTTP-em, bez uruchamiania jej, więc zawęża
// treść do wierszy niosących wskazany napis. To jest granica nazwana wprost,
// a nie udawanie zawężenia po drzewie DOM: wskazanie, którego na stronie nie
// ma, daje treść pustą i monitor pilnuje wtedy pustki — co widać w wykazie po
// zerowej długości odniesienia.
func trescPilnowana(tekst string, selektor *string) string {
	wskazanie := wartoscTekstuLubPusta(selektor)
	if wskazanie == "" {
		return tekst
	}
	var wybrane []string
	for _, wiersz := range strings.Split(tekst, "\n") {
		if strings.Contains(wiersz, wskazanie) {
			wybrane = append(wybrane, wiersz)
		}
	}
	return strings.Join(wybrane, "\n")
}

// roznicaTresci mierzy różnicę dwóch treści wiersz po wierszu.
//
// Miara jest prosta i uczciwa: wiersze porównywane po kolei, a nadmiar po
// jednej ze stron liczy się jako dopisanie albo usunięcie. Nie jest to
// najkrótsza ścieżka edycji — i nie ma być, bo monitor odpowiada na pytanie
// „czy i jak bardzo się zmieniło", a nie „jaka jest najoszczędniejsza łata".
func roznicaTresci(odniesienie, biezaca string) shared.BrowserContentDiff {
	stare := strings.Split(odniesienie, "\n")
	nowe := strings.Split(biezaca, "\n")
	if odniesienie == "" {
		stare = nil
	}
	if biezaca == "" {
		nowe = nil
	}

	roznica := shared.BrowserContentDiff{
		CharDelta: len(biezaca) - len(odniesienie),
	}
	wspolne := len(stare)
	if len(nowe) < wspolne {
		wspolne = len(nowe)
	}
	for i := 0; i < wspolne; i++ {
		if stare[i] == nowe[i] {
			continue
		}
		roznica.AddedLines++
		roznica.RemovedLines++
		if roznica.FirstChangedLine == nil {
			numer := i + 1
			roznica.FirstChangedLine = &numer
		}
		if roznica.Preview == nil {
			podglad := nowe[i]
			if len(podglad) > 200 {
				podglad = podglad[:200]
			}
			roznica.Preview = &podglad
		}
	}
	if len(nowe) > wspolne {
		roznica.AddedLines += len(nowe) - wspolne
		if roznica.FirstChangedLine == nil {
			numer := wspolne + 1
			roznica.FirstChangedLine = &numer
		}
		if roznica.Preview == nil {
			podglad := nowe[wspolne]
			if len(podglad) > 200 {
				podglad = podglad[:200]
			}
			roznica.Preview = &podglad
		}
	}
	if len(stare) > wspolne {
		roznica.RemovedLines += len(stare) - wspolne
		if roznica.FirstChangedLine == nil {
			numer := wspolne + 1
			roznica.FirstChangedLine = &numer
		}
	}
	return roznica
}

// bezwzglednaLiczba oddaje wartość bezwzględną — próg zmiany dotyczy przyrostu
// treści tak samo jak jej ubytku.
func bezwzglednaLiczba(wartosc int64) int64 {
	if wartosc < 0 {
		return -wartosc
	}
	return wartosc
}

// monitorKontraktu przekłada wiersz monitora na byt kontraktu.
func monitorKontraktu(w dane.MonitorPrzegladania) shared.BrowserMonitor {
	monitor := shared.BrowserMonitor{
		Id:              w.Kod,
		WindowId:        w.Okno,
		Url:             w.Url,
		Selector:        w.Selektor,
		IntervalSeconds: int(w.InterwalSekund),
		Enabled:         w.Wlaczony,
		Status:          shared.BrowserMonitorStatus(w.Stan),
		BaselineRef:     w.OdniesienieOdwolanie,
		LastCheckedAt:   chwilaZeZnacznika(w.Sprawdzono),
		LastChangedAt:   chwilaZeZnacznika(w.Zmieniono),
		CreatedAt:       chwilaBazy(w.Utworzono),
	}
	if w.ProgZmiany != nil {
		prog := int(*w.ProgZmiany)
		monitor.ChangeThreshold = &prog
	}
	if w.OdniesienieDlugosc != nil {
		dlugosc := int(*w.OdniesienieDlugosc)
		monitor.BaselineLength = &dlugosc
	}
	if w.KanalPowiadomienia != nil {
		kanal := shared.BrowserNotifyChannel(*w.KanalPowiadomienia)
		monitor.NotifyChannel = &kanal
	}
	return monitor
}
