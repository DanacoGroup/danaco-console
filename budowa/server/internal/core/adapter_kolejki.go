package core

import (
	"context"
	"strconv"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/models"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// adapterKolejek wypełnia port Kolejki repozytorium kolejek jednym silnikiem pętli sesyjnej i MultitaskingAI. Kolejka bez wiersza sesji zapisuje się bez dowiązania — nie jest to błąd.
type adapterKolejek struct {
	repozytorium dane.RepozytoriumKolejek
	silnik       silnikKolejki
	sesje        dane.RepozytoriumSesji
	telemetria   *telemetriaPostepu
	sesjeKolejek *pamiecSesjiKolejek
	// kanaly jest rejestrem kanałów modelu, którym wykonawca kroku prowadzi realną turę pozycji.
	kanaly *models.Rejestr
}

// nowyAdapterKolejek wiąże port z repozytorium kolejek i z silnikiem wykonania
// pozycji (kolejka_silnik.go). Adapter przekłada kontrakt; cykl życia zlecenia
// prowadzi silnik i nikt poza nim.
func nowyAdapterKolejek(repozytorium dane.RepozytoriumKolejek) *adapterKolejek {
	return &adapterKolejek{
		repozytorium: repozytorium,
		silnik:       silnikKolejki{repozytorium: repozytorium},
		sesjeKolejek: nowaPamiecSesjiKolejek(),
	}
}

// ZSesjami dokłada repozytorium sesji, którym adapter dowiązuje kolejkę do
// wiersza sesji. Bez niego adapter pracuje jak dotąd — dowiązanie jest
// dodatkiem, nie warunkiem.
func (a *adapterKolejek) ZSesjami(repozytorium dane.RepozytoriumSesji) *adapterKolejek {
	a.sesje = repozytorium
	return a
}

// ZWykonawcaModelu wpina wykonawcę kroku, który pozycję wchodzącą w stan wykonywana uruchamia turą kanału modelu; bez tej metody silnik zostaje czystą maszyną stanów.
func (a *adapterKolejek) ZWykonawcaModelu(kanaly *models.Rejestr, nadajnik Nadajnik) *adapterKolejek {
	a.kanaly = kanaly
	a.silnik = a.silnik.ZWykonawca(nowyWykonawcaModelu(kanaly, nadajnik, a.rozwiazKanalPozycji))
	return a
}

// UstawUjscieWyniku wpina odbiorcę zebranej treści tury pozycji. Silnik jest jeden, więc ujście dostaje każdą pozycję z treścią; pozycja nienależąca do podagenta kończy się zapisem donikąd.
func (a *adapterKolejek) UstawUjscieWyniku(ujscie func(ctx context.Context, pozycjaID int64, tresc string)) {
	a.silnik = a.silnik.ZUjsciemWyniku(ujscie)
}

// rozwiazKanalPozycji ustala kanał modelu i zasięgi wykonania pozycji kolejki; pozycja nie niesie kanału, więc krok jedzie domyślnym czynnym kanałem rejestru — pierwszym wierszem gotowym do pracy.
func (a *adapterKolejek) rozwiazKanalPozycji(_ context.Context, _ dane.Pozycja) (string, models.Zasiegi, bool) {
	if a.kanaly == nil {
		return "", models.Zasiegi{}, false
	}
	czynne := a.kanaly.Kontrakt(true)
	if len(czynne) == 0 {
		return "", models.Zasiegi{}, false
	}
	return czynne[0].Id, models.Zasiegi{}, true
}

// ZTelemetria dokłada producenta telemetrii postępu. Kolejka jest procesem: ma
// etapy równe pozycjom, stan i bieg naprawczy bez limitu obiegów, więc zasila
// Process Monitor tym samym zdarzeniem, co tura okna.
func (a *adapterKolejek) ZTelemetria(telemetria *telemetriaPostepu) *adapterKolejek {
	a.telemetria = telemetria
	return a
}

// Utworz zakłada kolejkę sesji wraz ze zleceniami początkowymi z ładunku.
// Wykaz zleceń czytany jest przed założeniem kolejki: ładunek uszkodzony ma
// zakończyć się odmową, a nie kolejką założoną w połowie.
func (a *adapterKolejek) Utworz(ctx context.Context, z shared.QueueCreateRequest) (shared.QueueCreateResponse, error) {
	pozycje, err := odczytajZlecenia(z.Items)
	if err != nil {
		return shared.QueueCreateResponse{}, err
	}
	kolejka := dane.Kolejka{
		Nazwa: wartoscTekstu(z.Name), Stan: shared.QueueStatusIdle,
		SesjaID: a.wierszSesji(ctx, z.SessionId),
	}
	if kolejka.Nazwa == "" {
		kolejka.Nazwa = z.SessionId
	}
	id, err := a.repozytorium.UtworzKolejke(ctx, kolejka)
	if err != nil {
		return shared.QueueCreateResponse{}, err
	}
	if err := a.silnik.Zasil(ctx, id, pozycje); err != nil {
		return shared.QueueCreateResponse{}, err
	}
	a.sesjeKolejek.Zapamietaj(id, z.SessionId, z.WindowIds)
	kolejkaKontraktu, err := a.odpowiedz(ctx, id, etapZalozenieKolejki)
	if err != nil {
		return shared.QueueCreateResponse{}, err
	}
	return shared.QueueCreateResponse{Queue: kolejkaKontraktu}, nil
}

// odpowiedz składa kolejkę kontraktu z zapisanego wiersza i odnotowuje postęp
// procesu kolejki. Jedno miejsce dla obu komend — drugiego składania nie ma.
func (a *adapterKolejek) odpowiedz(ctx context.Context, id int64, etap string) (shared.Queue, error) {
	zapisana, err := a.repozytorium.PobierzKolejke(ctx, id)
	if err != nil {
		return shared.Queue{}, err
	}
	kolejkaKontraktu := a.kolejkaKontraktu(ctx, zapisana)
	a.odnotujPostep(ctx, id, kolejkaKontraktu, etap)
	return kolejkaKontraktu, nil
}

// wierszSesji odnajduje wiersz sesji po identyfikatorze rdzenia. Brak
// repozytorium, brak wiersza i błąd odczytu dają kolejkę bez dowiązania —
// założenie kolejki nie może zależeć od tego, czy sesja zdążyła się utrwalić.
func (a *adapterKolejek) wierszSesji(ctx context.Context, idSesji string) *int64 {
	if a.sesje == nil || idSesji == "" {
		return nil
	}
	sesja, err := a.sesje.PoIdentyfikatorze(ctx, idSesji)
	if err != nil {
		return nil
	}
	id := sesja.ID
	return &id
}

// sesjaZWiersza odtwarza identyfikator sesji rdzenia z dowiązanego wiersza.
// Tą drogą kolejka zapisana przed ponownym uruchomieniem rdzenia wraca z sesją,
// choć pamięć powiązań jest już pusta.
func (a *adapterKolejek) sesjaZWiersza(ctx context.Context, wiersz *int64) string {
	if a.sesje == nil || wiersz == nil {
		return ""
	}
	sesja, err := a.sesje.Pobierz(ctx, *wiersz)
	if err != nil || sesja.IdentyfikatorZewnetrzny == nil {
		return ""
	}
	return *sesja.IdentyfikatorZewnetrzny
}

// Wykonaj prowadzi działanie kontraktu przez silnik wykonania: najpierw pozycje
// przechodzą stany, potem kolejka dostaje stan z nich wyprowadzony. Powtórzenie
// kroku nie ma limitu obiegów — adapter niczego nie zlicza i niczego
// nie odmawia.
func (a *adapterKolejek) Wykonaj(ctx context.Context, z shared.QueueActionRequest) (shared.QueueActionResponse, error) {
	id, err := strconv.ParseInt(z.QueueId, 10, 64)
	if err != nil {
		return shared.QueueActionResponse{}, protocol.JakoError(
			protocol.NowyBlad(shared.ErrorCodeNotFound, "kolejka nie istnieje: "+z.QueueId))
	}
	stan, err := a.silnik.Wykonaj(ctx, id, z.Action, z.ItemId)
	if err != nil {
		return shared.QueueActionResponse{}, err
	}
	if err := a.repozytorium.ZmienStanKolejki(ctx, id, stan, z.Action); err != nil {
		return shared.QueueActionResponse{}, err
	}
	kolejkaKontraktu, err := a.odpowiedz(ctx, id, etapDzialania(z.Action))
	if err != nil {
		return shared.QueueActionResponse{}, err
	}
	return shared.QueueActionResponse{Queue: kolejkaKontraktu}, nil
}

// etapDzialania nazywa etap telemetrii odpowiadający działaniu na kolejce.
// Powtórzenie kroku jest biegiem naprawczym i tak też się nazywa.
func etapDzialania(dzialanie shared.QueueAction) string {
	if dzialanie == shared.QueueActionRetry {
		return etapBiegNaprawczy
	}
	return etapDzialanieNaKolejce
}

// kolejkaKontraktu przekłada wiersz kolejki na kolejkę kontraktu; identyfikator sesji bierze z pamięci powiązań, a licznik obiegów pochodzi z pozycji, na której kolejka stoi.
func (a *adapterKolejek) kolejkaKontraktu(ctx context.Context, k dane.Kolejka) shared.Queue {
	idSesji, okna := a.sesjeKolejek.Odczytaj(k.ID)
	if idSesji == "" {
		idSesji = a.sesjaZWiersza(ctx, k.SesjaID)
	}
	kolejka := shared.Queue{
		Id: strconv.FormatInt(k.ID, 10), SessionId: idSesji, Status: k.Stan,
		WindowIds: okna, Cycle: a.silnik.Cykl(a.silnik.Pozycje(ctx, k.ID)),
		CreatedAt: chwilaBazy(k.Utworzono), UpdatedAt: chwilaBazy(k.Zaktualizowano),
	}
	if k.Nazwa != "" {
		nazwa := k.Nazwa
		kolejka.Name = &nazwa
	}
	// Polityka i liczba zleceń idą razem z kolejką, bo Queue Manager pokazuje oba pola w nagłówku okna.
	if polityka, err := a.repozytorium.PolitykaKolejki(ctx, k.ID); err == nil {
		zapis := politykaKontraktu(polityka)
		kolejka.Policy = &zapis
	}
	// Licznik bierze się z długości wykazu zawężonego stanem, nie z licznika wszystkich zleceń kolejki.
	if oczekujace, _, err := a.repozytorium.ZleceniaKolejki(ctx, k.ID,
		stanZleceniaBazy(shared.QueueItemStatusPending), 0); err == nil {
		oczekujacych := len(oczekujace)
		kolejka.PendingCount = &oczekujacych
	}
	return kolejka
}

// Kolejka oddaje kolejkę kontraktu po jej identyfikatorze. Służy czynnościom
// rodziny `queue.item.*`, których wynikiem jest kolejka po zmianie, a nie samo
// zlecenie — drugiego składania kolejki w rdzeniu nie ma.
func (a *adapterKolejek) Kolejka(ctx context.Context, id string) (shared.Queue, error) {
	wiersz, err := a.wierszKolejkiZlecen(ctx, id)
	if err != nil {
		return shared.Queue{}, err
	}
	return a.kolejkaKontraktu(ctx, wiersz), nil
}
