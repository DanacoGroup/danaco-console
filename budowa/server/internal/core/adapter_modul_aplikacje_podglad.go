// Moduł Apps — podgląd na żywo warstwy produktu: `apps.preview.start`
// i `apps.preview.stop`.
//
// PODGLĄD JEST SERWEREM, KTÓRY NAPRAWDĘ STOI. `AppsPreviewStartResponse` niesie
// `previewUrl` — adres, pod który Operator ma wejść. Adres wymyślony byłby
// dokładnie tym wzorcem szkody, którego pilnują sprawdziany skutku: odpowiedzią
// udaną, za którą nie ma niczego. Dlatego `apps.preview.start` podnosi nasłuch
// `net/http` na pętli zwrotnej, oddaje pod nim treść plików warsztatu wskazanej
// warstwy i zwraca adres wydany przez system operacyjny.
//
// PORT WYDAJE SYSTEM, NIE KONWENCJA. Nasłuch idzie na `127.0.0.1:0`, więc dwa
// okna podglądane naraz nie walczą o ten sam numer, a rdzeń nie musi zgadywać,
// co na tej maszynie jest wolne. Adres jest zawsze pętlą zwrotną: podgląd służy
// Operatorowi tej maszyny, a wystawienie warsztatu na świat byłoby udostępnieniem
// kodu produktu bez czyjejkolwiek zgody.
//
// SERWER NIE PRZEŻYWA RESTARTU RDZENIA i wiersz w bazie tego nie ukrywa: rejestr
// nasłuchów żyje w pamięci, a `apps.preview.start` po restarcie podnosi nowy
// serwer pod nowym adresem. Odtwarzanie nasłuchów przy starcie stawiałoby
// podgląd, którego nikt w tej sesji nie zamówił.
//
// ZATRZYMANIE ZAMYKA NASŁUCH. `apps.preview.stop` woła `Shutdown` i dopiero po
// jego powrocie melduje `stopped` — inaczej pole `stopped: true` znaczyłoby
// „poprosiliśmy", a port zostawałby zajęty.
package core

import (
	"context"
	"net"
	"net/http"
	"path"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// czasZamknieciaPodgladuApp jest granicą czekania na zamknięcie nasłuchu.
// Bez granicy `apps.preview.stop` wisiałby tak długo, jak długo trwa najdłuższe
// otwarte pobranie treści.
const czasZamknieciaPodgladuApp = 3 * time.Second

// podgladWarstwyApp to jeden stojący serwer podglądu.
type podgladWarstwyApp struct {
	serwer  *http.Server
	adres   string
	warstwa shared.AppWorkspaceLayer
	start   int64
}

// rejestrPodgladowApp trzyma stojące serwery podglądu — jeden na okno, tak jak
// wiersz `podglad_apps` (czoło migracji 203).
type rejestrPodgladowApp struct {
	zamek    sync.Mutex
	podglady map[string]*podgladWarstwyApp
}

// nowyRejestrPodgladowApp składa pusty rejestr nasłuchów.
func nowyRejestrPodgladowApp() *rejestrPodgladowApp {
	return &rejestrPodgladowApp{podglady: map[string]*podgladWarstwyApp{}}
}

// Zamknij zatrzymuje wszystkie stojące podglądy. Woła to zamknięcie rdzenia:
// gorutyna nasłuchu przeżyłaby zamknięcie procesu testowego i zostawiła port.
func (r *rejestrPodgladowApp) Zamknij() {
	if r == nil {
		return
	}
	r.zamek.Lock()
	stojace := make([]*podgladWarstwyApp, 0, len(r.podglady))
	for okno, podglad := range r.podglady {
		stojace = append(stojace, podglad)
		delete(r.podglady, okno)
	}
	r.zamek.Unlock()

	for _, podglad := range stojace {
		ctx, zakoncz := context.WithTimeout(context.Background(), czasZamknieciaPodgladuApp)
		_ = podglad.serwer.Shutdown(ctx)
		zakoncz()
	}
}

// UruchomPodglad obsługuje `apps.preview.start`. Warstwa pominięta znaczy
// frontend — kontrakt mówi to wprost.
func (a *adapterAplikacji) UruchomPodglad(ctx context.Context,
	z shared.AppsPreviewStartRequest) (shared.AppsPreviewStartResponse, error) {

	okno, err := oknoAplikacji(z.WindowId, "apps.preview.start")
	if err != nil {
		return shared.AppsPreviewStartResponse{}, err
	}
	warstwa := shared.AppWorkspaceLayer(shared.AppWorkspaceLayerFrontend)
	if z.Layer != nil {
		if err := sprawdzWarstweWarsztatu(*z.Layer); err != nil {
			return shared.AppsPreviewStartResponse{}, err
		}
		warstwa = *z.Layer
	}
	if a.podglady == nil {
		a.podglady = nowyRejestrPodgladowApp()
	}

	// Podgląd pustej warstwy jest odmową, nie serwerem oddającym pustkę:
	// Operator ma się dowiedzieć, że nie ma czego pokazać, a nie oglądać białą
	// stronę i zgadywać, czy to wina warsztatu, czy nasłuchu.
	pliki, err := a.plikiWarstwyApp(ctx, okno, warstwa)
	if err != nil {
		return shared.AppsPreviewStartResponse{}, err
	}
	if len(pliki) == 0 {
		a.zapiszStanPodgladuApp(ctx, okno, warstwa, "", shared.AppPreviewStatusFailed, 0)
		return shared.AppsPreviewStartResponse{}, bladWskazaniaAplikacji(
			"warstwa " + string(warstwa) + " okna " + okno +
				" nie ma ani jednego pliku — podgląd nie ma czego pokazać")
	}

	// Nasłuch zastany zatrzymujemy przed podniesieniem nowego: jedno okno ma
	// jeden podgląd (kontrakt `apps.preview.stop` nie ma czym wskazać drugiego),
	// więc drugi nasłuch byłby portem, którego nikt już nie zamknie.
	a.zatrzymajNasluchPodgladuApp(okno)

	nasluch, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		a.zapiszStanPodgladuApp(ctx, okno, warstwa, "", shared.AppPreviewStatusFailed, 0)
		return shared.AppsPreviewStartResponse{}, bladAplikacji(err)
	}
	adres := "http://" + nasluch.Addr().String() + "/"
	start := time.Now().UTC().UnixMilli()

	serwer := &http.Server{
		Handler:           a.obslugaPodgladuApp(okno, warstwa),
		ReadHeaderTimeout: czasZamknieciaPodgladuApp,
	}
	go func() { _ = serwer.Serve(nasluch) }()

	a.podglady.zamek.Lock()
	a.podglady.podglady[okno] = &podgladWarstwyApp{
		serwer: serwer, adres: adres, warstwa: warstwa, start: start,
	}
	a.podglady.zamek.Unlock()

	a.zapiszStanPodgladuApp(ctx, okno, warstwa, adres, shared.AppPreviewStatusRunning, start)
	a.dopiszDziennikApp(ctx, okno, nil, nil,
		"podgląd warstwy "+string(warstwa)+" wstał pod adresem "+adres+
			", plików w warstwie: "+strconv.Itoa(len(pliki)))

	return shared.AppsPreviewStartResponse{
		PreviewUrl: adres,
		Status:     shared.AppPreviewStatusRunning,
		StartedAt:  start,
	}, nil
}

// ZatrzymajPodglad obsługuje `apps.preview.stop`. Zatrzymanie okna, które nic
// nie podgląda, jest odmową `not_found`: pole `stopped: true` bez stojącego
// nasłuchu byłoby meldunkiem o czynności, która się nie odbyła.
func (a *adapterAplikacji) ZatrzymajPodglad(ctx context.Context,
	z shared.AppsPreviewStopRequest) (shared.AppsPreviewStopResponse, error) {

	okno, err := oknoAplikacji(z.WindowId, "apps.preview.stop")
	if err != nil {
		return shared.AppsPreviewStopResponse{}, err
	}
	podglad := a.zatrzymajNasluchPodgladuApp(okno)
	if podglad == nil {
		return shared.AppsPreviewStopResponse{}, bladNieznanegoBytuApp(
			"stojący podgląd okna", okno, dane.ErrBrakWiersza)
	}
	teraz := time.Now().UTC().UnixMilli()
	a.zapiszStanPodgladuApp(ctx, okno, podglad.warstwa, podglad.adres,
		shared.AppPreviewStatusStopped, podglad.start)
	a.dopiszDziennikApp(ctx, okno, nil, nil,
		"podgląd warstwy "+string(podglad.warstwa)+" zatrzymany po "+
			strconv.FormatInt(teraz-podglad.start, 10)+" ms")
	return shared.AppsPreviewStopResponse{Stopped: true}, nil
}

// zatrzymajNasluchPodgladuApp zamyka nasłuch okna i oddaje jego opis; brak
// stojącego podglądu daje nil.
func (a *adapterAplikacji) zatrzymajNasluchPodgladuApp(okno string) *podgladWarstwyApp {
	if a.podglady == nil {
		return nil
	}
	a.podglady.zamek.Lock()
	podglad, stoi := a.podglady.podglady[okno]
	if stoi {
		delete(a.podglady.podglady, okno)
	}
	a.podglady.zamek.Unlock()
	if !stoi {
		return nil
	}
	zamkniecie, zakoncz := context.WithTimeout(context.Background(), czasZamknieciaPodgladuApp)
	defer zakoncz()
	_ = podglad.serwer.Shutdown(zamkniecie)
	return podglad
}

// obslugaPodgladuApp składa obsługiwacza żądań serwera podglądu. Treść bierze
// się z bazy przy każdym żądaniu, nie z migawki z chwili podniesienia: warsztat
// zmienia się w trakcie pracy, a podgląd ma pokazywać stan bieżący („natychmiastowe
// odświeżanie wyniku pracy równolegle z edycją" — opracowanie, Frontend Workspace).
func (a *adapterAplikacji) obslugaPodgladuApp(okno string,
	warstwa shared.AppWorkspaceLayer) http.Handler {

	return http.HandlerFunc(func(odpowiedz http.ResponseWriter, zadanie *http.Request) {
		pliki, err := a.plikiWarstwyApp(zadanie.Context(), okno, warstwa)
		if err != nil {
			http.Error(odpowiedz, "podgląd nie może odczytać warsztatu", http.StatusInternalServerError)
			return
		}
		zadana := strings.TrimPrefix(path.Clean("/"+zadanie.URL.Path), "/")
		if zadana == "" || zadana == "." {
			odpowiedz.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = odpowiedz.Write([]byte(spisPodgladuApp(okno, warstwa, pliki)))
			return
		}
		for _, plik := range pliki {
			if strings.TrimPrefix(plik.Sciezka, "/") != zadana {
				continue
			}
			odpowiedz.Header().Set("Content-Type", typTresciPodgladuApp(plik.Sciezka))
			_, _ = odpowiedz.Write([]byte(plik.Tresc))
			return
		}
		http.NotFound(odpowiedz, zadanie)
	})
}

// spisPodgladuApp składa stronę wejściową podglądu — spis plików warstwy wraz
// z odnośnikami. Gdy warstwa niesie `index.html`, oddajemy jego treść wprost:
// to jest strona produktu, a nie spis warsztatu.
func spisPodgladuApp(okno string, warstwa shared.AppWorkspaceLayer,
	pliki []dane.PlikWarsztatu) string {

	for _, plik := range pliki {
		if strings.TrimPrefix(plik.Sciezka, "/") == "index.html" {
			return plik.Tresc
		}
	}
	var zapis strings.Builder
	zapis.WriteString("<!doctype html><meta charset=\"utf-8\"><title>Podgląd — " +
		tekstSvgApp(okno) + "</title><h1>Podgląd warstwy " + string(warstwa) + "</h1><ul>")
	for _, plik := range pliki {
		sciezka := strings.TrimPrefix(plik.Sciezka, "/")
		zapis.WriteString("<li><a href=\"/" + tekstSvgApp(sciezka) + "\">" +
			tekstSvgApp(sciezka) + "</a></li>")
	}
	zapis.WriteString("</ul>")
	return zapis.String()
}

// typTresciPodgladuApp rozstrzyga nagłówek treści po rozszerzeniu ścieżki.
// Wykaz jest krótki z zamysłu: podgląd oddaje pliki warsztatu, a te są tekstem.
func typTresciPodgladuApp(sciezka string) string {
	switch {
	case strings.HasSuffix(sciezka, ".html"), strings.HasSuffix(sciezka, ".htm"):
		return "text/html; charset=utf-8"
	case strings.HasSuffix(sciezka, ".css"):
		return "text/css; charset=utf-8"
	case strings.HasSuffix(sciezka, ".js"), strings.HasSuffix(sciezka, ".mjs"),
		strings.HasSuffix(sciezka, ".ts"):
		return "text/javascript; charset=utf-8"
	case strings.HasSuffix(sciezka, ".json"):
		return "application/json; charset=utf-8"
	case strings.HasSuffix(sciezka, ".svg"):
		return "image/svg+xml"
	}
	return "text/plain; charset=utf-8"
}

// plikiWarstwyApp zwraca pliki jednej warstwy warsztatu, po ścieżce.
func (a *adapterAplikacji) plikiWarstwyApp(ctx context.Context, okno string,
	warstwa shared.AppWorkspaceLayer) ([]dane.PlikWarsztatu, error) {

	wszystkie, err := a.repozytorium.PlikiWarsztatu(ctx, okno)
	if err != nil {
		return nil, bladAplikacji(err)
	}
	warstwy := make([]dane.PlikWarsztatu, 0, len(wszystkie))
	for _, plik := range wszystkie {
		if plik.Warstwa != warstwa {
			continue
		}
		warstwy = append(warstwy, plik)
	}
	sort.SliceStable(warstwy, func(i, j int) bool { return warstwy[i].Sciezka < warstwy[j].Sciezka })
	return warstwy, nil
}

// zapiszStanPodgladuApp utrwala stan nasłuchu. Nieudany zapis nie przewraca
// komendy: serwer już stoi (albo już nie stoi), a wiersz jest zapisem faktu,
// nie warunkiem jego zajścia.
func (a *adapterAplikacji) zapiszStanPodgladuApp(ctx context.Context, okno string,
	warstwa shared.AppWorkspaceLayer, adres string, stan shared.AppPreviewStatus, start int64) {

	podglad := dane.PodgladApp{
		Okno: okno, Warstwa: string(warstwa), Adres: adres,
		Stan: string(stan), Rozpoczeto: start,
	}
	if stan == shared.AppPreviewStatusStopped {
		teraz := time.Now().UTC().UnixMilli()
		podglad.Zatrzymano = &teraz
	}
	_ = a.repozytorium.ZapiszPodgladApp(ctx, podglad)
}

// dopiszDziennikApp dopisuje wiersz dziennika modułu. Nieudany zapis nie
// przewraca czynności, której dziennik dotyczy — dziennik opisuje pracę, nie
// warunkuje jej.
func (a *adapterAplikacji) dopiszDziennikApp(ctx context.Context, okno string,
	wdrozenie, komponent *string, tresc string) {

	_ = a.repozytorium.DopiszWierszDziennikaApp(ctx, dane.WierszDziennikaApp{
		Okno: okno, WdrozenieKod: wdrozenie, KomponentKod: komponent,
		Chwila: time.Now().UTC().UnixMilli(), Tresc: tresc,
	})
}

// adresPodgladuApp oddaje adres stojącego podglądu okna albo pustkę. Używa go
// `apps.endpoint.probe`, gdy środowisko nie ma nadanej domeny: zapytanie próbne
// idzie wtedy pod ten sam adres, pod którym Operator ogląda produkt.
func (a *adapterAplikacji) adresPodgladuApp(okno string) string {
	if a.podglady == nil {
		return ""
	}
	a.podglady.zamek.Lock()
	defer a.podglady.zamek.Unlock()
	podglad, stoi := a.podglady.podglady[okno]
	if !stoi {
		return ""
	}
	return podglad.adres
}
