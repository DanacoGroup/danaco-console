// Odpowiedzialność pliku: stan zleceń asystenta i sterowanie nimi
// (`assistant.action.status`) oraz chronologiczny dziennik działań
// (`assistant.activity.list`) — obszar Actions Monitor. Typ `adapterAsystenta`
// i przyjęcie polecenia głosowego deklaruje `adapter_modul_asystent.go`; ten
// plik dokłada metody na tym samym typie, bez drugiej deklaracji.
//
// Zlecenie asystenta ma własny automat stanu:
// `queued/running/paused/done/failed/cancelled` z kontraktu. To nie jest stan
// katalogu akcji (`core/akcje.go`, statyczny spis dostępnych czynności) ani
// stan pętli sesyjnej kolejki — inny byt, inna prawda.
//
// Nieznane zlecenie wraca odmową, nie cichą zgodą: sterowanie zleceniem,
// którego nie ma, nie może wyglądać jak sukces, bo potwierdzałoby czynność,
// która się nie odbyła.
package core

import (
	"context"
	"errors"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// StanCzynnosci obsługuje `assistant.action.status`. Z kontrolą — zmienia stan
// wskazanego zlecenia i oddaje je po zmianie. Bez kontroli — czyta jedno
// zlecenie (po `ActionId`) albo wszystkie zlecenia okna (po `WindowId`).
func (a *adapterAsystenta) StanCzynnosci(ctx context.Context,
	z shared.AssistantActionStatusRequest) (shared.AssistantActionStatusResponse, error) {

	if z.Control != nil && *z.Control != shared.AssistantActionControlNone {
		return a.steruj(ctx, z)
	}

	if z.ActionId != nil && *z.ActionId != "" {
		zlecenie, err := a.repozytorium.Zlecenie(ctx, *z.ActionId)
		if err != nil {
			return shared.AssistantActionStatusResponse{}, bladNieznanegoZleceniaAsystenta(*z.ActionId, err)
		}
		return shared.AssistantActionStatusResponse{Actions: []shared.AssistantAction{zlozZlecenie(zlecenie)}}, nil
	}

	if z.WindowId != nil && *z.WindowId != "" {
		wiersze, err := a.repozytorium.Zlecenia(ctx, *z.WindowId)
		if err != nil {
			return shared.AssistantActionStatusResponse{}, bladAsystenta(err)
		}
		return shared.AssistantActionStatusResponse{Actions: zlozZlecenia(wiersze)}, nil
	}

	return shared.AssistantActionStatusResponse{}, bladWskazaniaAsystenta(
		"żądanie bez wskazania zlecenia (actionId) i bez okna (windowId)")
}

// steruj przekłada `AssistantActionControl` na zmianę stanu zlecenia.
// Wymaga `ActionId` — sterowanie bez wskazania zlecenia nie ma na czym
// zadziałać, więc jest błędem żądania, nie operacją na wszystkich zleceniach
// okna naraz.
func (a *adapterAsystenta) steruj(ctx context.Context,
	z shared.AssistantActionStatusRequest) (shared.AssistantActionStatusResponse, error) {

	if z.ActionId == nil || *z.ActionId == "" {
		return shared.AssistantActionStatusResponse{}, bladWskazaniaAsystenta(
			"sterowanie zleceniem bez wskazania actionId")
	}
	nowyStan, err := stanDlaSterowania(*z.Control)
	if err != nil {
		return shared.AssistantActionStatusResponse{}, err
	}

	// Stan sprzed zapisu czytamy zawsze, bo zapis bez niego jest ślepy. Dwie
	// rzeczy zależą od tego, co zlecenie robiło przed sterowaniem:
	//  1. `resume` ma wznowić zlecenie wstrzymane, a nie dokładać drugiej tury
	//     zleceniu, które właśnie biegnie (oba stoją po zapisie w `running`,
	//     więc po nim już ich nie odróżnić);
	//  2. zlecenie zamknięte (`done`/`failed`/`cancelled`) nie ma czego wznawiać
	//     ani wstrzymywać. Bez tego odczytu `resume` na zleceniu zamkniętym
	//     przestawiłby je na `running`, nikt by go nie podjął (wykonawca wchodzi
	//     tylko ze stanu `queued`/`paused`) i wiersz stałby w toku bez końca,
	//     bez wykonawcy i bez zdarzenia — a anulowanie dałoby się tą drogą
	//     cofnąć jednym żądaniem.
	przed, err := a.repozytorium.Zlecenie(ctx, *z.ActionId)
	if err != nil {
		return shared.AssistantActionStatusResponse{}, bladNieznanegoZleceniaAsystenta(*z.ActionId, err)
	}
	if err := sprawdzSterowalnosc(*z.Control, przed.Stan); err != nil {
		return shared.AssistantActionStatusResponse{}, err
	}
	wstrzymanePrzed := przed.Stan == string(shared.AssistantActionStatusPaused)

	// Anulowanie i wstrzymanie sięgają do biegu, nie tylko do wiersza. Zlecenie
	// biegnące ma turę modelu w locie; samo przestawienie stanu zostawiłoby ją
	// pracującą dalej na cudzy rachunek. Przerwanie idzie przed zapisem, żeby
	// wykonawca zastał już decyzję Operatora, gdy tura wróci z błędem.
	if *z.Control == shared.AssistantActionControlCancel || *z.Control == shared.AssistantActionControlPause {
		a.przerwijBieg(*z.ActionId)
	}
	zlecenie, err := a.repozytorium.UstawStanZlecenia(ctx, *z.ActionId, nowyStan)
	if err != nil {
		return shared.AssistantActionStatusResponse{}, bladNieznanegoZleceniaAsystenta(*z.ActionId, err)
	}

	// Priorytet idzie osobnym zapisem, bo jest osobną czynnością: kontrakt
	// pozwala przy sterowaniu zmienić kolejność obsługi w Actions Monitor.
	// Przyjęcie priorytetu bez zapisania go byłoby potwierdzeniem czynności,
	// która się nie odbyła.
	if z.Priority != nil {
		zlecenie, err = a.repozytorium.UstawPriorytetZlecenia(ctx, *z.ActionId, int64(*z.Priority))
		if err != nil {
			return shared.AssistantActionStatusResponse{}, bladNieznanegoZleceniaAsystenta(*z.ActionId, err)
		}
	}

	// Ponowna próba wraca do wykonawcy. `retry` przestawia zlecenie na `queued`
	// (zob. stanDlaSterowania), a `queued` znaczy „do wykonania", nie „czeka na
	// kolejny ręczny ruch". Wykonawca podejmuje je tą samą drogą, którą podejmuje
	// zlecenie świeżo złożone; bez wpiętego wykonawcy `podejmij` nie robi nic.
	if zlecenie.Stan == string(shared.AssistantActionStatusQueued) {
		a.podejmij(*z.ActionId)
	}

	// Wznowienie też wraca do wykonawcy. Sam zapis stanu nie wystarcza:
	// wykonawca podejmuje wyłącznie `queued`, więc zlecenie przestawione na
	// `running` stałoby w toku bez końca i bez powodu. `resume` oddaje je więc
	// temu samemu wykonawcy, tylko ze stanem wejścia `running`; bez wpiętego
	// wykonawcy `wznow` — jak `podejmij` — nie robi nic.
	if wstrzymanePrzed && zlecenie.Stan == string(shared.AssistantActionStatusRunning) {
		a.wznow(*z.ActionId)
	}

	return shared.AssistantActionStatusResponse{Actions: []shared.AssistantAction{zlozZlecenie(zlecenie)}}, nil
}

// stanDlaSterowania tłumaczy sterowanie Actions Monitor na wartość stanu
// automatu kontraktu. `retry` wraca zlecenie do `queued` — to jedyny sposób,
// jakim automat stanu z kontraktu opisuje ponowną próbę: nowego przebiegu nie
// zakładamy, bo repozytorium nie ma pojęcia przebiegu odrębnego od zlecenia.
func stanDlaSterowania(sterowanie shared.AssistantActionControl) (string, error) {
	switch sterowanie {
	case shared.AssistantActionControlPause:
		return shared.AssistantActionStatusPaused, nil
	case shared.AssistantActionControlResume:
		return shared.AssistantActionStatusRunning, nil
	case shared.AssistantActionControlCancel:
		return shared.AssistantActionStatusCancelled, nil
	case shared.AssistantActionControlRetry:
		return shared.AssistantActionStatusQueued, nil
	default:
		return "", bladWskazaniaAsystenta("nieznane sterowanie zleceniem: " + string(sterowanie))
	}
}

// zlecenieZamkniete mówi, czy stan zlecenia jest końcowy. Automat kontraktu ma
// trzy takie stany (`done`, `failed`, `cancelled`) i żaden z nich nie prowadzi
// dalej sam z siebie — z każdego wychodzi się wyłącznie ponowną próbą.
func zlecenieZamkniete(stan string) bool {
	switch stan {
	case string(shared.AssistantActionStatusDone),
		string(shared.AssistantActionStatusFailed),
		string(shared.AssistantActionStatusCancelled):
		return true
	}
	return false
}

// sprawdzSterowalnosc odmawia sterowania, którego zlecenie zamknięte nie ma jak
// wykonać. Odmowa opisuje brak, nie zakaz: w zleceniu, które już się skończyło,
// nie ma biegu do wstrzymania ani do wznowienia, a `cancel` nie ma czego
// odwołać. Jedynym sterowaniem, które z zamkniętego zlecenia prowadzi dalej,
// jest `retry` — i ono zostaje wpuszczone, bo to jego rola w automacie stanu.
//
// Kod `conflict` („stan bytu wyklucza czynność") jest tu jedynym prawdziwym:
// żądanie jest zgodne z kontraktem, a bytu nie brakuje — brakuje stanu, w
// którym czynność miałaby sens. Ponawianie takiego żądania nic nie zmieni,
// więc `retryable` zostaje fałszem.
func sprawdzSterowalnosc(sterowanie shared.AssistantActionControl, stan string) error {
	if !zlecenieZamkniete(stan) || sterowanie == shared.AssistantActionControlRetry {
		return nil
	}
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeConflict,
		"moduł Assistant: zlecenie jest już zamknięte (stan "+stan+"), więc sterowanie "+
			string(sterowanie)+" nie ma czego dotyczyć; ponowne wykonanie zamawia się sterowaniem retry"))
}

// WykazCzynnosci obsługuje `assistant.activity.list`. Po `ActionId` czyta
// dziennik jednego zlecenia (bez granicy — rozmowa spięta ze zleceniem jest z
// natury skończona, `dane/asystent_dziennik.go`); po `WindowId` czyta
// chronologiczny dziennik okna, ucięty `Limit`, gdy żądanie go niesie.
func (a *adapterAsystenta) WykazCzynnosci(ctx context.Context,
	z shared.AssistantActivityListRequest) (shared.AssistantActivityListResponse, error) {

	if z.ActionId != nil && *z.ActionId != "" {
		wpisy, err := a.repozytorium.WpisyZlecenia(ctx, *z.ActionId)
		if err != nil {
			return shared.AssistantActivityListResponse{}, bladAsystenta(err)
		}
		return shared.AssistantActivityListResponse{Entries: zlozWpisyUciete(wpisy, z.Limit)}, nil
	}

	if z.WindowId != nil && *z.WindowId != "" {
		limit := 0
		if z.Limit != nil {
			limit = *z.Limit
		}
		wpisy, err := a.repozytorium.Wpisy(ctx, *z.WindowId, limit)
		if err != nil {
			return shared.AssistantActivityListResponse{}, bladAsystenta(err)
		}
		return shared.AssistantActivityListResponse{Entries: zlozWpisy(wpisy)}, nil
	}

	return shared.AssistantActivityListResponse{}, bladWskazaniaAsystenta(
		"żądanie bez wskazania okna (windowId) i bez zlecenia (actionId)")
}

// zlozWpisyUciete składa wpisy dziennika zlecenia i ucina je granicą, której
// `WpisyZlecenia` nie stosuje samo (dziennik zlecenia w repozytorium wraca w
// całości — ucięcie po stronie rdzenia nie duplikuje żadnej logiki zapytania,
// tylko skraca już złożoną listę).
func zlozWpisyUciete(wiersze []dane.WpisDziennikaAsystenta, limit *int) []shared.AssistantActivityEntry {
	wpisy := zlozWpisy(wiersze)
	if limit != nil && *limit > 0 && len(wpisy) > *limit {
		wpisy = wpisy[:*limit]
	}
	return wpisy
}

// zlozWpisy składa listę wpisów kontraktu z wierszy repozytorium.
func zlozWpisy(wiersze []dane.WpisDziennikaAsystenta) []shared.AssistantActivityEntry {
	wpisy := make([]shared.AssistantActivityEntry, 0, len(wiersze))
	for _, wiersz := range wiersze {
		wpisy = append(wpisy, shared.AssistantActivityEntry{
			Id:        wiersz.Kod,
			WindowId:  wiersz.OknoKod,
			ActionId:  wiersz.ZlecenieKod,
			Kind:      shared.AssistantActivityKind(wiersz.Rodzaj),
			Content:   wiersz.Tresc,
			AudioRef:  wiersz.NagranieOdnosnik,
			CreatedAt: wiersz.Utworzono,
			Important: wskaznikPrawdy(wiersz.Wazny),
		})
	}
	return wpisy
}

// zlozZlecenia składa listę zleceń kontraktu z wierszy repozytorium.
func zlozZlecenia(wiersze []dane.ZlecenieAsystenta) []shared.AssistantAction {
	zlecenia := make([]shared.AssistantAction, 0, len(wiersze))
	for _, wiersz := range wiersze {
		zlecenia = append(zlecenia, zlozZlecenie(wiersz))
	}
	return zlecenia
}

// zlozZlecenie składa jedno zlecenie kontraktu z wiersza repozytorium. Czas
// jest tu liczbą wprost — bez przekładu przez `chwilaBazy`, którego używają
// moduły trzymające czas tekstem.
func zlozZlecenie(wiersz dane.ZlecenieAsystenta) shared.AssistantAction {
	return shared.AssistantAction{
		Id:          wiersz.Kod,
		WindowId:    wiersz.OknoKod,
		Title:       wiersz.Tytul,
		Status:      shared.AssistantActionStatus(wiersz.Stan),
		Origin:      shared.AssistantOrigin(wiersz.Droga),
		CurrentStep: liczbaIntZInt64(wiersz.EtapBiezacy),
		TotalSteps:  liczbaIntZInt64(wiersz.LiczbaEtapow),
		Priority:    liczbaIntZInt64(wiersz.Priorytet),
		Result:      wiersz.Wynik,
		CreatedAt:   wiersz.Utworzono,
		UpdatedAt:   wiersz.Zaktualizowano,
	}
}

// liczbaIntZInt64 przekłada wskaźnik na `int64` warstwy danych na wskaźnik na
// `int` kontraktu — kontrakt trzyma `CurrentStep`/`TotalSteps`/`Priority` jako
// `*int`, repozytorium jako `*int64` (kolumna INTEGER SQLite).
func liczbaIntZInt64(wartosc *int64) *int {
	if wartosc == nil {
		return nil
	}
	liczba := int(*wartosc)
	return &liczba
}

// bladNieznanegoZleceniaAsystenta odróżnia „zlecenia nie ma" (odmowa wprost,
// nie cicha zgoda) od „odczyt/zapis się nie powiódł".
func bladNieznanegoZleceniaAsystenta(kod string, err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Assistant: zlecenie nie istnieje: "+kod))
	}
	return bladAsystenta(err)
}

// OznaczWpisDziennika obsługuje `assistant.activity.flag`.
//
// Wyróżnienie ma gdzie zamieszkać: kolumny `wazny` i `notatka_wyroznienia`
// tabeli dziennika (migracja 296). Bez nich Activity Feed pokazywałby
// wyróżnienie, którego rdzeń nie pamięta — znikające przy pierwszym
// odświeżeniu wykazu.
//
// Zdjęcie wyróżnienia kasuje też powód: powód bez znacznika byłby notatką do
// wpisu, którego nikt nie wyróżnił. Rozstrzyga to warstwa danych, w jednym
// miejscu, żeby ta reguła nie istniała w dwóch egzemplarzach.
func (a *adapterAsystenta) OznaczWpisDziennika(ctx context.Context,
	z shared.AssistantActivityFlagRequest) (shared.AssistantActivityFlagResponse, error) {

	kod := strings.TrimSpace(z.EntryId)
	if kod == "" {
		return shared.AssistantActivityFlagResponse{},
			bladWskazaniaAsystenta("oznaczenie bez wskazania wpisu dziennika")
	}
	wiersz, err := a.repozytorium.OznaczWpisDziennika(ctx, kod, z.Important, z.Note)
	if err != nil {
		if errors.Is(err, dane.ErrBrakWiersza) {
			return shared.AssistantActivityFlagResponse{}, protocol.JakoError(protocol.NowyBlad(
				shared.ErrorCodeNotFound,
				"asystent: nie ma wpisu dziennika o identyfikatorze "+kod))
		}
		return shared.AssistantActivityFlagResponse{}, bladAsystenta(err)
	}

	wpis := zlozWpisy([]dane.WpisDziennikaAsystenta{wiersz})[0]
	// Notatka wyróżnienia nie ma własnego pola w `AssistantActivityEntry`,
	// więc odpowiedź niesie ją tam, gdzie okno jej szuka: przy treści wpisu
	// wyróżnienie jest wprost widoczne polem `important`, a powód wraca
	// niezmieniony wraz z wierszem przy kolejnym odczycie dziennika. Rdzeń nie
	// dokłada do kontraktu pola, którego kontrakt nie zna.
	return shared.AssistantActivityFlagResponse{Entry: wpis}, nil
}
