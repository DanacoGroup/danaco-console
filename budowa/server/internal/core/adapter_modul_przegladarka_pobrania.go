// Odpowiedzialność pliku: menedżer pobrań i nagrywarka makr oraz granice
// działania Wykonawcy — `browser.download.list`, `.control`,
// `browser.macro.record`, `browser.executor.limits.set`, `.get`.
//
// Pobranie naprawdę ściąga plik. Ponowienie (`retry`) idzie po treść spod adresu
// pobrania i odkłada ją w magazynie modułu, a postęp w wierszu jest liczbą
// bajtów, które na dysku leżą — nie deklaracją. Wstrzymanie i wznowienie
// zmieniają stan kolejki, przerwanie ją kończy, zdjęcie usuwa wpis.
//
// Makro zapisuje kroki w kształcie kontraktu (`AutomationStep`), tym samym,
// którym jedzie moduł Automations. Dzięki temu przekazanie scenariusza do
// Automations jest przełożeniem wiersza, a nie tłumaczeniem jednego kształtu na
// drugi.
//
// Granice Wykonawcy mają wartość domyślną w kodzie, nie w schemacie. Brak
// wiersza znaczy „granice domyślne rdzenia" i tak też odpowiada odczyt —
// zamiast odmawiać, że nikt jeszcze niczego nie ustawił.
package core

import (
	"context"
	"encoding/json"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

const (
	// domyslneKrokiWykonawcy i domyslnyCzasWykonawcy są granicami obowiązującymi
	// wtedy, gdy Operator nie ustawił własnych. Wartości skończone, bo pętla
	// wykonawcza bez granicy chodziłaby po stronie bez końca.
	domyslneKrokiWykonawcy = 40
	domyslnyCzasWykonawcy  = 600
)

// WykazPobran obsługuje `browser.download.list`.
func (a *adapterPrzegladarki) WykazPobran(ctx context.Context,
	z shared.BrowserDownloadListRequest) (shared.BrowserDownloadListResponse, error) {

	stan := ""
	if z.Status != nil {
		stan = string(*z.Status)
	}
	wiersze, err := a.repozytorium.Pobrania(ctx, wartoscTekstuLubPusta(z.WindowId), stan, wartoscLiczby(z.Limit))
	if err != nil {
		return shared.BrowserDownloadListResponse{}, bladPrzegladarki(err)
	}
	pobrania := make([]shared.BrowserDownload, 0, len(wiersze))
	for _, wiersz := range wiersze {
		pobrania = append(pobrania, pobranieKontraktu(wiersz))
	}
	return shared.BrowserDownloadListResponse{Downloads: pobrania}, nil
}

// SterujPobraniem obsługuje `browser.download.control`.
func (a *adapterPrzegladarki) SterujPobraniem(ctx context.Context,
	z shared.BrowserDownloadControlRequest) (shared.BrowserDownloadControlResponse, error) {

	if strings.TrimSpace(z.DownloadId) == "" {
		return shared.BrowserDownloadControlResponse{}, bladWskazaniaPrzegladarki(
			"komenda download.control bez pobrania")
	}
	pobranie, err := a.repozytorium.Pobranie(ctx, z.DownloadId)
	if err != nil {
		return shared.BrowserDownloadControlResponse{}, bladWierszaPrzegladania("pobrania", z.DownloadId, err)
	}

	switch z.Action {
	case shared.BrowserDownloadActionPause:
		pobranie.Stan = string(shared.BrowserDownloadStatusPaused)
	case shared.BrowserDownloadActionResume:
		pobranie.Stan = string(shared.BrowserDownloadStatusRunning)
	case shared.BrowserDownloadActionCancel:
		pobranie.Stan = string(shared.BrowserDownloadStatusCancelled)
		pobranie.Zakonczono = terazWBazie()
	case shared.BrowserDownloadActionRemove:
		usuniete, err := a.repozytorium.UsunPobranie(ctx, z.DownloadId)
		if err != nil {
			return shared.BrowserDownloadControlResponse{}, bladPrzegladarki(err)
		}
		if !usuniete {
			return shared.BrowserDownloadControlResponse{}, bladNieznanegoBytu("pobrania", z.DownloadId)
		}
		// Wiersza już nie ma, ale odpowiedź niesie to, co zdjęto: klient ma
		// z czego zdjąć pozycję z wykazu, zamiast zgadywać, która zniknęła.
		return shared.BrowserDownloadControlResponse{Download: pobranieKontraktu(pobranie)}, nil
	case shared.BrowserDownloadActionRetry:
		pobrany, err := a.sciagnijPobranie(ctx, pobranie)
		if err != nil {
			return shared.BrowserDownloadControlResponse{}, err
		}
		pobranie = pobrany
	default:
		return shared.BrowserDownloadControlResponse{}, bladWskazaniaPrzegladarki(
			"nieznane działanie menedżera pobrań: " + string(z.Action))
	}

	zapisane, err := a.repozytorium.ZapiszPobranie(ctx, pobranie)
	if err != nil {
		return shared.BrowserDownloadControlResponse{}, bladPrzegladarki(err)
	}
	return shared.BrowserDownloadControlResponse{Download: pobranieKontraktu(zapisane)}, nil
}

// sciagnijPobranie ściąga treść spod adresu pobrania i odkłada ją w magazynie
// modułu. Postęp i rozmiar biorą się z bajtów, które naprawdę przyszły.
func (a *adapterPrzegladarki) sciagnijPobranie(ctx context.Context,
	pobranie dane.PobraniePrzegladania) (dane.PobraniePrzegladania, error) {

	tresc, err := pobierzStrone(ctx, pobranie.Url)
	if err != nil {
		komunikat := err.Error()
		pobranie.Stan = string(shared.BrowserDownloadStatusFailed)
		pobranie.KomunikatBledu = &komunikat
		pobranie.Zakonczono = terazWBazie()
		zapisane, zapis := a.repozytorium.ZapiszPobranie(ctx, pobranie)
		if zapis != nil {
			return dane.PobraniePrzegladania{}, bladPrzegladarki(zapis)
		}
		return zapisane, nil
	}
	bajty := []byte(tresc.Html)
	odwolanie, err := a.zapiszTresc(bajty)
	if err != nil {
		return dane.PobraniePrzegladania{}, err
	}
	rozmiar := int64(len(bajty))
	pobranie.SciezkaDocelowa = &odwolanie
	pobranie.OdebranoBajtow = &rozmiar
	pobranie.RazemBajtow = &rozmiar
	pobranie.KomunikatBledu = nil
	pobranie.Stan = string(shared.BrowserDownloadStatusCompleted)
	pobranie.Zakonczono = terazWBazie()
	return pobranie, nil
}

// NagrywajMakro obsługuje `browser.macro.record` — start zapisu, jego koniec
// i doklejenie pojedynczego kroku.
func (a *adapterPrzegladarki) NagrywajMakro(ctx context.Context,
	z shared.BrowserMacroRecordRequest) (shared.BrowserMacroRecordResponse, error) {

	if strings.TrimSpace(z.WindowId) == "" {
		return shared.BrowserMacroRecordResponse{}, bladWskazaniaPrzegladarki("komenda macro.record bez okna")
	}

	if z.Action == shared.BrowserMacroActionStart {
		nazwa := wartoscTekstuLubPusta(z.Name)
		if nazwa == "" {
			nazwa = "Makro przeglądania"
		}
		kod := wartoscTekstuLubPusta(z.MacroId)
		if kod == "" {
			kod = nowyIdentyfikator(przedrostekMakra)
		}
		zapisane, err := a.repozytorium.ZapiszMakro(ctx, dane.MakroPrzegladania{
			Kod: kod, Okno: z.WindowId, Nazwa: nazwa, Nagrywanie: true,
		})
		if err != nil {
			return shared.BrowserMacroRecordResponse{}, bladPrzegladarki(err)
		}
		return shared.BrowserMacroRecordResponse{Macro: makroKontraktu(zapisane), Recording: true}, nil
	}

	kod := wartoscTekstuLubPusta(z.MacroId)
	if kod == "" {
		return shared.BrowserMacroRecordResponse{}, bladWskazaniaPrzegladarki(
			"czynność " + string(z.Action) + " bez wskazania makra — nie wiadomo, do czego dopisać krok")
	}
	makro, err := a.repozytorium.Makro(ctx, kod)
	if err != nil {
		return shared.BrowserMacroRecordResponse{}, bladWierszaPrzegladania("makra przeglądania", kod, err)
	}

	switch z.Action {
	case shared.BrowserMacroActionStep:
		if z.Step == nil {
			return shared.BrowserMacroRecordResponse{}, bladWskazaniaPrzegladarki(
				"doklejenie kroku bez samego kroku")
		}
		if !makro.Nagrywanie {
			return shared.BrowserMacroRecordResponse{}, bladWskazaniaPrzegladarki(
				"makro " + kod + " nie jest w trakcie nagrywania — krok nie ma się gdzie dopisać")
		}
		kroki := krokiMakra(makro.KrokiJson)
		kroki = append(kroki, *z.Step)
		zapis, err := json.Marshal(kroki)
		if err != nil {
			return shared.BrowserMacroRecordResponse{}, bladPrzegladarki(err)
		}
		tekst := string(zapis)
		makro.KrokiJson = &tekst
	case shared.BrowserMacroActionStop:
		makro.Nagrywanie = false
	default:
		return shared.BrowserMacroRecordResponse{}, bladWskazaniaPrzegladarki(
			"nieznana czynność nagrywarki makra: " + string(z.Action))
	}
	if nazwa := wartoscTekstuLubPusta(z.Name); nazwa != "" {
		makro.Nazwa = nazwa
	}

	zapisane, err := a.repozytorium.ZapiszMakro(ctx, makro)
	if err != nil {
		return shared.BrowserMacroRecordResponse{}, bladPrzegladarki(err)
	}
	return shared.BrowserMacroRecordResponse{
		Macro: makroKontraktu(zapisane), Recording: zapisane.Nagrywanie,
	}, nil
}

// UstawGraniceWykonawcy obsługuje `browser.executor.limits.set`.
func (a *adapterPrzegladarki) UstawGraniceWykonawcy(ctx context.Context,
	z shared.BrowserExecutorLimitsSetRequest) (shared.BrowserExecutorLimitsSetResponse, error) {

	zasieg, wskazanie := zasiegGranic(z.WindowId, z.SessionId)
	granica := dane.GranicaWykonawcy{
		Zasieg:                string(zasieg),
		ZasiegID:              wskazanie,
		MaxKrokow:             int64(domyslneKrokiWykonawcy),
		MaxCzasSekund:         int64(domyslnyCzasWykonawcy),
		DomenyDozwoloneJson:   wykazJson(z.AllowedDomains),
		DomenyZablokowaneJson: wykazJson(z.BlockedDomains),
		PotwierdzajWyslanie:   z.ConfirmBeforeSubmit == nil || *z.ConfirmBeforeSubmit,
	}
	// Zastane granice są punktem wyjścia: żądanie zmieniające sam limit kroków
	// nie ma kasować wykazu domen ustawionego wcześniej.
	if zastane, err := a.repozytorium.Granice(ctx, granica.Zasieg, wskazanie); err == nil {
		granica.MaxKrokow, granica.MaxCzasSekund = zastane.MaxKrokow, zastane.MaxCzasSekund
		if z.AllowedDomains == nil {
			granica.DomenyDozwoloneJson = zastane.DomenyDozwoloneJson
		}
		if z.BlockedDomains == nil {
			granica.DomenyZablokowaneJson = zastane.DomenyZablokowaneJson
		}
		if z.ConfirmBeforeSubmit == nil {
			granica.PotwierdzajWyslanie = zastane.PotwierdzajWyslanie
		}
	} else if !isBrakWiersza(err) {
		return shared.BrowserExecutorLimitsSetResponse{}, bladPrzegladarki(err)
	}
	if z.MaxSteps != nil && *z.MaxSteps > 0 {
		granica.MaxKrokow = int64(*z.MaxSteps)
	}
	if z.MaxDurationSeconds != nil && *z.MaxDurationSeconds > 0 {
		granica.MaxCzasSekund = int64(*z.MaxDurationSeconds)
	}

	zapisana, err := a.repozytorium.ZapiszGranice(ctx, granica)
	if err != nil {
		return shared.BrowserExecutorLimitsSetResponse{}, bladPrzegladarki(err)
	}
	return shared.BrowserExecutorLimitsSetResponse{Limits: graniceKontraktu(zapisana)}, nil
}

// OdczytajGraniceWykonawcy obsługuje `browser.executor.limits.get`.
func (a *adapterPrzegladarki) OdczytajGraniceWykonawcy(ctx context.Context,
	z shared.BrowserExecutorLimitsGetRequest) (shared.BrowserExecutorLimitsGetResponse, error) {

	zasieg, wskazanie := zasiegGranic(z.WindowId, z.SessionId)
	zapisana, err := a.repozytorium.Granice(ctx, string(zasieg), wskazanie)
	if err != nil {
		if !isBrakWiersza(err) {
			return shared.BrowserExecutorLimitsGetResponse{}, bladPrzegladarki(err)
		}
		// Brak wiersza nie jest brakiem odpowiedzi: obowiązują wtedy granice
		// domyślne rdzenia i to je odczyt oddaje, wraz ze wskazaniem zasięgu.
		return shared.BrowserExecutorLimitsGetResponse{Limits: shared.BrowserExecutorLimits{
			Scope: zasieg, ScopeId: wskaznikNiepustyPrzegladania(wskazanie),
			MaxSteps: domyslneKrokiWykonawcy, MaxDurationSeconds: domyslnyCzasWykonawcy,
			ConfirmBeforeSubmit: true,
		}}, nil
	}
	return shared.BrowserExecutorLimitsGetResponse{Limits: graniceKontraktu(zapisana)}, nil
}

// zasiegGranic rozstrzyga, czego granice dotyczą: okna przeglądarki czy karty
// sesji. Okno ma pierwszeństwo, bo jest węższe — granice okna są tym, co
// Operator widzi przed sobą.
func zasiegGranic(okno, sesja *string) (shared.ConfigScope, string) {
	if wskazanie := wartoscTekstuLubPusta(okno); wskazanie != "" {
		return shared.ConfigScopeWindow, wskazanie
	}
	if wskazanie := wartoscTekstuLubPusta(sesja); wskazanie != "" {
		return shared.ConfigScopeSession, wskazanie
	}
	// Żądanie bez wskazania dotyczy całej aplikacji — najszerszego z dziewięciu
	// poziomów zasięgu, tego, który ustępuje każdemu węższemu.
	return shared.ConfigScopeApplication, ""
}

// wskaznikNiepustyPrzegladania oddaje wskaźnik na tekst albo brak dla tekstu pustego.
func wskaznikNiepustyPrzegladania(tekst string) *string {
	if tekst == "" {
		return nil
	}
	kopia := tekst
	return &kopia
}

// krokiMakra odczytuje kroki zapisane kolumną JSON.
func krokiMakra(zapis *string) []shared.AutomationStep {
	if zapis == nil || *zapis == "" {
		return nil
	}
	var kroki []shared.AutomationStep
	if err := json.Unmarshal([]byte(*zapis), &kroki); err != nil {
		return nil
	}
	return kroki
}

// pobranieKontraktu przekłada wiersz pobrania na byt kontraktu.
func pobranieKontraktu(w dane.PobraniePrzegladania) shared.BrowserDownload {
	return shared.BrowserDownload{
		Id:            w.Kod,
		WindowId:      w.Okno,
		Url:           w.Url,
		FileName:      w.NazwaPliku,
		TargetPath:    w.SciezkaDocelowa,
		MimeType:      w.TypMime,
		Status:        shared.BrowserDownloadStatus(w.Stan),
		ReceivedBytes: w.OdebranoBajtow,
		TotalBytes:    w.RazemBajtow,
		ErrorMessage:  w.KomunikatBledu,
		StartedAt:     chwilaBazy(w.Rozpoczeto),
		FinishedAt:    chwilaZeZnacznika(w.Zakonczono),
	}
}

// makroKontraktu przekłada wiersz makra na byt kontraktu.
func makroKontraktu(w dane.MakroPrzegladania) shared.BrowserMacro {
	makro := shared.BrowserMacro{
		Id:         w.Kod,
		WindowId:   w.Okno,
		Name:       w.Nazwa,
		Recording:  w.Nagrywanie,
		WorkflowId: w.AutomatykaKod,
		CreatedAt:  chwilaBazy(w.Utworzono),
		UpdatedAt:  chwilaBazy(w.Zaktualizowano),
	}
	if kroki := krokiMakra(w.KrokiJson); len(kroki) > 0 {
		makro.Steps = kroki
	}
	return makro
}

// graniceKontraktu przekładają wiersz granic na byt kontraktu.
func graniceKontraktu(w dane.GranicaWykonawcy) shared.BrowserExecutorLimits {
	granice := shared.BrowserExecutorLimits{
		Scope:               shared.ConfigScope(w.Zasieg),
		ScopeId:             wskaznikNiepustyPrzegladania(w.ZasiegID),
		MaxSteps:            int(w.MaxKrokow),
		MaxDurationSeconds:  int(w.MaxCzasSekund),
		ConfirmBeforeSubmit: w.PotwierdzajWyslanie,
		UpdatedAt:           chwilaBazy(w.Zaktualizowano),
	}
	if dozwolone := wykazZJson(w.DomenyDozwoloneJson); len(dozwolone) > 0 {
		granice.AllowedDomains = dozwolone
	}
	if zablokowane := wykazZJson(w.DomenyZablokowaneJson); len(zablokowane) > 0 {
		granice.BlockedDomains = zablokowane
	}
	return granice
}

// odlozPobranie ściąga zasób, którego nie da się pokazać jako strony, i zakłada
// dla niego wiersz w menedżerze pobrań. Oddaje odmowę komendy `browser.navigate`
// — bo migawki strony z tego nie ma — ale odmowa nazywa skutek, który naprawdę
// zaszedł: pobranie o podanym identyfikatorze, z bajtami leżącymi w magazynie.
//
// To jest jedyna droga, którą pobrania powstają, i jest to droga naturalna:
// w przeglądarce plik pobiera się przez wejście pod jego adres, a nie osobnym
// poleceniem „dodaj pobranie" — takiego kontrakt zresztą nie niesie.
func (a *adapterPrzegladarki) odlozPobranie(ctx context.Context, okno, adres string, powod error) error {
	plik, err := pobierzPlik(ctx, adres)
	if err != nil {
		return bladPobraniaStrony(adres, err)
	}
	odwolanie, err := a.zapiszTresc(plik.Bajty)
	if err != nil {
		return err
	}
	rozmiar := int64(len(plik.Bajty))
	pobranie := dane.PobraniePrzegladania{
		Kod:             nowyIdentyfikator(przedrostekPobrania),
		Okno:            okno,
		Url:             adres,
		NazwaPliku:      &plik.Nazwa,
		SciezkaDocelowa: &odwolanie,
		Stan:            string(shared.BrowserDownloadStatusCompleted),
		OdebranoBajtow:  &rozmiar,
		RazemBajtow:     &rozmiar,
		Zakonczono:      terazWBazie(),
	}
	if plik.Typ != "" {
		typ := plik.Typ
		pobranie.TypMime = &typ
	}
	zapisane, err := a.repozytorium.ZapiszPobranie(ctx, pobranie)
	if err != nil {
		return bladPrzegladarki(err)
	}
	return bladWskazaniaPrzegladarki(powod.Error() + " — rdzeń odłożył go jako pobranie " +
		zapisane.Kod + " (" + jakoLiczba(int(rozmiar)) + " B); wykaz: browser.download.list")
}
