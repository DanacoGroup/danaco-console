// Odpowiedzialność pliku: dwie komendy rodziny `queue.*` — `queue.list` (wykaz
// kolejek jednego silnika pętli sesyjnej i MultitaskingAI) oraz `queue.link`
// (wiązanie kolejki z oknami, ekspertem, projektem albo automatyką).
//
// To rozszerzenie adaptera, nie drugi adapter. Metody wiszą na `adapterKolejek`
// z `adapter_kolejki.go`, więc jadą tym samym repozytorium, tym samym silnikiem
// wykonania i tym samym przekładem kolejki na kontrakt, co `queue.create`
// i `queue.action`. Drugiego silnika kolejek nie ma nigdzie.
//
// Okna kolejki mają dwa źródła: okna podane przy `queue.create` leżą w pamięci
// powiązań rdzenia (`pamiec_sesji_kolejek.go`), a okna podane przy `queue.link`
// w tabeli `powiazanie_kolejki`. Kolejka kontraktu oddaje sumę obu — inaczej po
// ponownym uruchomieniu rdzenia okna z `queue.create` znikałyby bez śladu,
// a okna z `queue.link` nie byłyby widoczne w ogóle.
package core

import (
	"context"
	"strconv"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// repozytoriumWiazanKolejek to rozszerzenie repozytorium kolejek o czynności,
// których `queue.create` i `queue.action` nie potrzebowały: wykaz kolejek,
// sprawdzenie istnienia i powiązania kolejki.
//
// Interfejs stoi po stronie czytelnika. Deklaracja mieszka tutaj, a nie
// w `dane.RepozytoriumKolejek`, bo wymagają jej wyłącznie dwie komendy tej
// rodziny; port kolejek tych czynności nie zna i nie musi.
type repozytoriumWiazanKolejek interface {
	ListaKolejek(ctx context.Context, stan *shared.QueueStatus) ([]dane.Kolejka, error)
	CzyKolejkaIstnieje(ctx context.Context, id int64) (bool, error)
	PowiazaniaKolejki(ctx context.Context, kolejkaID int64) ([]dane.PowiazanieKolejki, error)
	ZwiazKolejke(ctx context.Context, kolejkaID int64, powiazania []dane.PowiazanieKolejki) error
}

// wiazania oddaje repozytorium rozszerzone albo odmowę. Repozytorium, które
// nie niesie tych czynności, jest usterką montażu, a nie powodem do ciszy:
// pusty wykaz kolejek wyglądałby na „nie ma żadnej kolejki", a to nieprawda.
func (a *adapterKolejek) wiazania() (repozytoriumWiazanKolejek, error) {
	rozszerzone, ok := a.repozytorium.(repozytoriumWiazanKolejek)
	if !ok {
		return nil, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"kolejki: repozytorium nie niesie wykazu ani powiązań kolejki"))
	}
	return rozszerzone, nil
}

// ── queue.list ───────────────────────────────────────────────────────────────

// Wykaz zwraca kolejki spełniające warunki żądania.
//
// Stan zawęża w bazie, sesja i okno — w rdzeniu. Stan jest kolumną, więc sito
// idzie do SQL. Sesja kolejki bywa znana wyłącznie z pamięci powiązań (kolejka
// założona przed pierwszą utrwaloną wiadomością sesji nie ma czym wypełnić
// `kolejka.sesja_id`), a okna leżą w dwóch źródłach opisanych w nagłówku pliku
// — sito po nich składa się więc na kolejce kontraktu, po przekładzie. Sito
// w SQL milczałoby o kolejkach, które warunek spełniają.
//
// Żądanie bez ani jednego warunku jest zgodne z kontraktem: wszystkie trzy pola
// są opcjonalne, więc oddajemy wszystkie kolejki, a nie odmowę.
//
// Pusty wykaz nie jest odmową. Sesja bez kolejek to prawdziwa odpowiedź „zero
// kolejek", nie „nie ma takiej sesji" — sesja żyje w pamięci nadzorcy i bywa
// bez wiersza w bazie, więc rdzeń nie ma jak odróżnić sesji nieistniejącej od
// sesji bez kolejek i nie udaje, że ma.
func (a *adapterKolejek) Wykaz(ctx context.Context,
	z shared.QueueListRequest) (shared.QueueListResponse, error) {

	rozszerzone, err := a.wiazania()
	if err != nil {
		return shared.QueueListResponse{}, err
	}
	wiersze, err := rozszerzone.ListaKolejek(ctx, z.Status)
	if err != nil {
		return shared.QueueListResponse{}, err
	}
	zadanaSesja := wartoscTekstu(z.SessionId)
	zadaneOkno := wartoscTekstu(z.WindowId)
	kolejki := make([]shared.Queue, 0, len(wiersze))
	for _, wiersz := range wiersze {
		kolejka, err := a.kolejkaZPowiazaniami(ctx, rozszerzone, wiersz)
		if err != nil {
			return shared.QueueListResponse{}, err
		}
		if zadanaSesja != "" && kolejka.SessionId != zadanaSesja {
			continue
		}
		if zadaneOkno != "" && !niesieOkno(kolejka.WindowIds, zadaneOkno) {
			continue
		}
		kolejki = append(kolejki, kolejka)
	}
	return shared.QueueListResponse{Queues: kolejki}, nil
}

// niesieOkno mówi, czy kolejka obsługuje wskazane okno.
func niesieOkno(okna []string, szukane string) bool {
	for _, okno := range okna {
		if okno == szukane {
			return true
		}
	}
	return false
}

// ── queue.link ───────────────────────────────────────────────────────────────

// Zwiaz wiąże kolejkę z oknami, ekspertem, projektem albo automatyką.
//
// Żądanie bez ani jednego bytu jest odmawiane. Cztery pola wiązania są
// opcjonalne z osobna, ale wszystkie naraz puste znaczą żądanie bez treści:
// odpowiedź „kolejka po powiązaniu" byłaby wtedy kolejką po niczym, czyli ciszą
// udającą skutek. Kod odmowy to `validation_failed`.
//
// Powiązania się dokładają. Kontrakt zna wyłącznie wiązanie — komendy
// rozwiązującej nie ma, więc `queue.link` niczego nie zdejmuje. Powiązanie
// powtórzone nie jest drugim faktem i nie jest błędem; zapis je pomija.
//
// Istnienia bytu wiązanego rdzeń tu nie sprawdza. Adapter kolejek nie ma
// dostępu do repozytoriów ekspertów, projektów ani automatyk, a kolejka nie
// może zależeć od tego, czy inny moduł zdążył się utrwalić. Zapisany zostaje
// identyfikator w kształcie, w jakim przyszedł.
func (a *adapterKolejek) Zwiaz(ctx context.Context,
	z shared.QueueLinkRequest) (shared.QueueLinkResponse, error) {

	rozszerzone, err := a.wiazania()
	if err != nil {
		return shared.QueueLinkResponse{}, err
	}
	id, err := strconv.ParseInt(z.QueueId, 10, 64)
	if err != nil {
		return shared.QueueLinkResponse{}, bladBrakuKolejki(z.QueueId)
	}
	powiazania := powiazaniaZadania(z)
	if len(powiazania) == 0 {
		return shared.QueueLinkResponse{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeValidationFailed,
			"kolejki: żądanie queue.link nie wskazuje ani okna, ani eksperta, "+
				"ani projektu, ani automatyki"))
	}
	jest, err := rozszerzone.CzyKolejkaIstnieje(ctx, id)
	if err != nil {
		return shared.QueueLinkResponse{}, err
	}
	if !jest {
		return shared.QueueLinkResponse{}, bladBrakuKolejki(z.QueueId)
	}
	if err := rozszerzone.ZwiazKolejke(ctx, id, powiazania); err != nil {
		return shared.QueueLinkResponse{}, err
	}
	zapisana, err := a.repozytorium.PobierzKolejke(ctx, id)
	if err != nil {
		return shared.QueueLinkResponse{}, err
	}
	kolejka, err := a.kolejkaZPowiazaniami(ctx, rozszerzone, zapisana)
	if err != nil {
		return shared.QueueLinkResponse{}, err
	}
	return shared.QueueLinkResponse{Queue: kolejka}, nil
}

// powiazaniaZadania składa wykaz bytów wskazanych żądaniem. Puste wskazanie
// jest pomijane, nie zapisywane jako powiązanie z bytem o pustej nazwie.
func powiazaniaZadania(z shared.QueueLinkRequest) []dane.PowiazanieKolejki {
	powiazania := []dane.PowiazanieKolejki{}
	for _, okno := range z.WindowIds {
		if okno != "" {
			powiazania = append(powiazania,
				dane.PowiazanieKolejki{Rodzaj: dane.RodzajPowiazaniaOkno, Byt: okno})
		}
	}
	pojedyncze := []struct {
		rodzaj string
		byt    *string
	}{
		{dane.RodzajPowiazaniaEkspert, z.AgentId},
		{dane.RodzajPowiazaniaProjekt, z.ProjectId},
		{dane.RodzajPowiazaniaAutomatyka, z.WorkflowId},
	}
	for _, wskazanie := range pojedyncze {
		if byt := wartoscTekstu(wskazanie.byt); byt != "" {
			powiazania = append(powiazania,
				dane.PowiazanieKolejki{Rodzaj: wskazanie.rodzaj, Byt: byt})
		}
	}
	return powiazania
}

// ── wspólne ──────────────────────────────────────────────────────────────────

// kolejkaZPowiazaniami składa kolejkę kontraktu i dokłada do jej wykazu okien
// te, które zapisała komenda `queue.link`. Suma dwóch źródeł opisana jest
// w nagłówku pliku; powtórzeń w wykazie nie ma.
func (a *adapterKolejek) kolejkaZPowiazaniami(ctx context.Context,
	rozszerzone repozytoriumWiazanKolejek, wiersz dane.Kolejka) (shared.Queue, error) {

	kolejka := a.kolejkaKontraktu(ctx, wiersz)
	powiazania, err := rozszerzone.PowiazaniaKolejki(ctx, wiersz.ID)
	if err != nil {
		return shared.Queue{}, err
	}
	for _, powiazanie := range powiazania {
		if powiazanie.Rodzaj != dane.RodzajPowiazaniaOkno {
			continue
		}
		if !niesieOkno(kolejka.WindowIds, powiazanie.Byt) {
			kolejka.WindowIds = append(kolejka.WindowIds, powiazanie.Byt)
		}
	}
	return kolejka, nil
}

// bladBrakuKolejki składa odmowę wskazującą kolejkę, której nie ma. Ten sam
// kształt, którym odmawia `queue.action` przy identyfikatorze nie do odczytania.
func bladBrakuKolejki(identyfikator string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
		"kolejki: kolejka "+identyfikator+" nie istnieje"))
}

// bladMontazuKolejek składa odmowę usterki montażu obszaru kolejek.
func bladMontazuKolejek(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError, powod))
}
