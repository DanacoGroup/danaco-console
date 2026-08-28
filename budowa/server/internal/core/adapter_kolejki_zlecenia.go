// Ten plik obsługuje dwanaście czynności Queue Managera dla zleceń kolejki,
// polityki, zadań martwych i głębokości w czasie: `queue.item.*`,
// `queue.policy.set`, `queue.dead.list`, `queue.depth.get`. Cykl życia
// kolejki prowadzi `kolejka_silnik.go`.
package core

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// odcinekGlebokosciDomyslny odpowiada godzinie — rozdzielczości, w której
// wykres głębokości Queue Managera pokazuje tydzień obciążenia.
const odcinekGlebokosciDomyslny = 3600

// DodajZlecenie dokłada zlecenie do kolejki poza harmonogramem, a klucz
// idempotencji już użyty oddaje zlecenie zastane jako duplikat.
func (a *adapterKolejek) DodajZlecenie(ctx context.Context,
	z shared.QueueItemEnqueueRequest) (shared.QueueItemEnqueueResponse, error) {

	kolejka, err := a.wierszKolejkiZlecen(ctx, z.QueueId)
	if err != nil {
		return shared.QueueItemEnqueueResponse{}, err
	}
	if klucz := wartoscTekstu(z.IdempotencyKey); klucz != "" {
		zastane, err := a.repozytorium.ZlecenieKluczem(ctx, kolejka.ID, klucz)
		if err == nil {
			return shared.QueueItemEnqueueResponse{
				Item: zlecenieKontraktu(zastane), Duplicate: true,
			}, nil
		}
	}
	zlecenie := dane.Zlecenie{
		Kod: nowyIdentyfikator(przedrostekZleceniaKolejki), KolejkaID: kolejka.ID,
		Stan: stanZleceniaBazy(shared.QueueItemStatusPending), Ladunek: zapisLadunku(z.Payload),
		Priorytet: wartoscLiczby(z.Priority), KluczIdempotencji: z.IdempotencyKey,
	}
	if z.ScheduledAt != nil && *z.ScheduledAt > 0 {
		termin := znacznikChwiliAutomatyzacji(*z.ScheduledAt)
		zlecenie.Termin = &termin
		zlecenie.Stan = stanZleceniaBazy(shared.QueueItemStatusDelayed)
	}
	zapisane, err := a.repozytorium.DodajZlecenie(ctx, zlecenie)
	if err != nil {
		return shared.QueueItemEnqueueResponse{}, bladZlecenia(err)
	}
	return shared.QueueItemEnqueueResponse{Item: zlecenieKontraktu(zapisane), Duplicate: false}, nil
}

// ZdejmijZlecenie zdejmuje zlecenie z kolejki przed jego wykonaniem. Zlecenie
// już przetworzone nie daje się zdjąć — zdjęcie po fakcie byłoby kłamstwem
// o tym, że praca się nie odbyła.
func (a *adapterKolejek) ZdejmijZlecenie(ctx context.Context,
	z shared.QueueItemDequeueRequest) (shared.QueueItemDequeueResponse, error) {

	_, zlecenie, err := a.zlecenieZadania(ctx, z.QueueId, z.ItemId)
	if err != nil {
		return shared.QueueItemDequeueResponse{}, err
	}
	zdjete := false
	if !czyStanKoncowyZlecenia(zlecenie.Stan) {
		if err := a.repozytorium.ZmienStanZlecenia(ctx, zlecenie.ID,
			stanZleceniaBazy(shared.QueueItemStatusRemoved)); err != nil {
			return shared.QueueItemDequeueResponse{}, bladZlecenia(err)
		}
		zdjete = true
	}
	kolejka, err := a.Kolejka(ctx, z.QueueId)
	if err != nil {
		return shared.QueueItemDequeueResponse{}, err
	}
	return shared.QueueItemDequeueResponse{Removed: zdjete, Queue: kolejka}, nil
}

// OdlozZlecenie odkłada wykonanie zlecenia o wskazany czas, przyjmując
// wyłącznie dodatnią liczbę sekund odłożenia.
func (a *adapterKolejek) OdlozZlecenie(ctx context.Context,
	z shared.QueueItemDelayRequest) (shared.QueueItemDelayResponse, error) {

	_, zlecenie, err := a.zlecenieZadania(ctx, z.QueueId, z.ItemId)
	if err != nil {
		return shared.QueueItemDelayResponse{}, err
	}
	if z.DelaySeconds <= 0 {
		return shared.QueueItemDelayResponse{},
			bladWskazaniaZlecenia("odłożenie wymaga dodatniej liczby sekund")
	}
	termin := time.Now().UTC().Add(time.Duration(z.DelaySeconds) * time.Second).
		Format(formatZnacznikaBazy)
	if err := a.repozytorium.OdlozZlecenie(ctx, zlecenie.ID, termin); err != nil {
		return shared.QueueItemDelayResponse{}, bladZlecenia(err)
	}
	odlozone, err := a.repozytorium.Zlecenie(ctx, zlecenie.Kod)
	if err != nil {
		return shared.QueueItemDelayResponse{}, bladZlecenia(err)
	}
	return shared.QueueItemDelayResponse{Item: zlecenieKontraktu(odlozone)}, nil
}

// PodzielZlecenie dzieli zlecenie na podzadania. Zlecenie źródłowe zostaje
// zamknięte — tak mówi kontrakt.
func (a *adapterKolejek) PodzielZlecenie(ctx context.Context,
	z shared.QueueItemSplitRequest) (shared.QueueItemSplitResponse, error) {

	kolejka, zlecenie, err := a.zlecenieZadania(ctx, z.QueueId, z.ItemId)
	if err != nil {
		return shared.QueueItemSplitResponse{}, err
	}
	ladunki, err := ladunkiPodzadan(z.Payloads)
	if err != nil {
		return shared.QueueItemSplitResponse{}, err
	}
	powstale := make([]shared.QueueItem, 0, len(ladunki))
	for _, ladunek := range ladunki {
		zapis := string(ladunek)
		podzadanie, err := a.repozytorium.DodajZlecenie(ctx, dane.Zlecenie{
			Kod: nowyIdentyfikator(przedrostekZleceniaKolejki), KolejkaID: kolejka.ID,
			Stan: stanZleceniaBazy(shared.QueueItemStatusPending), Ladunek: &zapis,
			Priorytet: zlecenie.Priorytet, ZlecenieZrodloweID: &zlecenie.ID,
			PrzebiegID: zlecenie.PrzebiegID,
		})
		if err != nil {
			return shared.QueueItemSplitResponse{}, bladZlecenia(err)
		}
		powstale = append(powstale, zlecenieKontraktu(podzadanie))
	}
	if err := a.repozytorium.ZmienStanZlecenia(ctx, zlecenie.ID,
		stanZleceniaBazy(shared.QueueItemStatusSucceeded)); err != nil {
		return shared.QueueItemSplitResponse{}, bladZlecenia(err)
	}
	return shared.QueueItemSplitResponse{Items: powstale}, nil
}

// ScalZlecenia scala kilka zleceń w jedno. Brak ładunku składa ładunki źródłowe
// w wykaz — tak mówi kontrakt.
func (a *adapterKolejek) ScalZlecenia(ctx context.Context,
	z shared.QueueItemMergeRequest) (shared.QueueItemMergeResponse, error) {

	kolejka, err := a.wierszKolejkiZlecen(ctx, z.QueueId)
	if err != nil {
		return shared.QueueItemMergeResponse{}, err
	}
	if len(z.ItemIds) == 0 {
		return shared.QueueItemMergeResponse{},
			bladWskazaniaZlecenia("scalenie bez wskazania zleceń źródłowych")
	}
	zrodlowe := make([]dane.Zlecenie, 0, len(z.ItemIds))
	skladowe := make([]json.RawMessage, 0, len(z.ItemIds))
	for _, kod := range z.ItemIds {
		zlecenie, err := a.zlecenieKolejki(ctx, kolejka.ID, kod)
		if err != nil {
			return shared.QueueItemMergeResponse{}, err
		}
		zrodlowe = append(zrodlowe, zlecenie)
		if zlecenie.Ladunek != nil {
			skladowe = append(skladowe, json.RawMessage(*zlecenie.Ladunek))
		}
	}
	ladunek := zapisLadunku(z.Payload)
	if ladunek == nil {
		ladunek = zapisStrukturalny(skladowe)
	}
	scalone, err := a.repozytorium.DodajZlecenie(ctx, dane.Zlecenie{
		Kod: nowyIdentyfikator(przedrostekZleceniaKolejki), KolejkaID: kolejka.ID,
		Stan: stanZleceniaBazy(shared.QueueItemStatusPending), Ladunek: ladunek,
		Priorytet: zrodlowe[0].Priorytet, PrzebiegID: zrodlowe[0].PrzebiegID,
	})
	if err != nil {
		return shared.QueueItemMergeResponse{}, bladZlecenia(err)
	}
	for _, zlecenie := range zrodlowe {
		if err := a.repozytorium.ZmienStanZlecenia(ctx, zlecenie.ID,
			stanZleceniaBazy(shared.QueueItemStatusSucceeded)); err != nil {
			return shared.QueueItemMergeResponse{}, bladZlecenia(err)
		}
	}
	return shared.QueueItemMergeResponse{Item: zlecenieKontraktu(scalone)}, nil
}

// SkierujZlecenie kieruje zlecenie do innej kolejki albo do innego wykonawcy,
// odmawiając, gdy skierowanie niczego nie zmienia.
func (a *adapterKolejek) SkierujZlecenie(ctx context.Context,
	z shared.QueueItemRouteRequest) (shared.QueueItemRouteResponse, error) {

	kolejka, zlecenie, err := a.zlecenieZadania(ctx, z.QueueId, z.ItemId)
	if err != nil {
		return shared.QueueItemRouteResponse{}, err
	}
	docelowa := kolejka.ID
	if kod := wartoscTekstu(z.TargetQueueId); kod != "" {
		wiersz, err := a.wierszKolejkiZlecen(ctx, kod)
		if err != nil {
			return shared.QueueItemRouteResponse{}, err
		}
		docelowa = wiersz.ID
	}
	if docelowa == kolejka.ID && wartoscTekstu(z.TargetAgentId) == "" {
		return shared.QueueItemRouteResponse{}, bladWskazaniaZlecenia(
			"skierowanie bez kolejki docelowej i bez eksperta docelowego niczego nie zmienia")
	}
	if err := a.repozytorium.SkierujZlecenie(ctx, zlecenie.ID, docelowa, z.TargetAgentId); err != nil {
		return shared.QueueItemRouteResponse{}, bladZlecenia(err)
	}
	skierowane, err := a.repozytorium.Zlecenie(ctx, zlecenie.Kod)
	if err != nil {
		return shared.QueueItemRouteResponse{}, bladZlecenia(err)
	}
	return shared.QueueItemRouteResponse{Item: zlecenieKontraktu(skierowane)}, nil
}

// RozgalezZlecenie rozgałęzia przetwarzanie zlecenia na tory równoległe,
// z których każdy zakłada osobne zlecenie potomne.
func (a *adapterKolejek) RozgalezZlecenie(ctx context.Context,
	z shared.QueueItemBranchRequest) (shared.QueueItemBranchResponse, error) {

	kolejka, zlecenie, err := a.zlecenieZadania(ctx, z.QueueId, z.ItemId)
	if err != nil {
		return shared.QueueItemBranchResponse{}, err
	}
	if len(z.Branches) == 0 {
		return shared.QueueItemBranchResponse{},
			bladWskazaniaZlecenia("rozgałęzienie bez ani jednego toru")
	}
	tory := make([]shared.QueueItem, 0, len(z.Branches))
	for _, tor := range z.Branches {
		docelowa := kolejka.ID
		if kod := wartoscTekstu(tor.TargetQueueId); kod != "" {
			wiersz, err := a.wierszKolejkiZlecen(ctx, kod)
			if err != nil {
				return shared.QueueItemBranchResponse{}, err
			}
			docelowa = wiersz.ID
		}
		ladunek := zapisLadunku(tor.Payload)
		if ladunek == nil {
			ladunek = zlecenie.Ladunek
		}
		powstale, err := a.repozytorium.DodajZlecenie(ctx, dane.Zlecenie{
			Kod: nowyIdentyfikator(przedrostekZleceniaKolejki), KolejkaID: docelowa,
			Stan: stanZleceniaBazy(shared.QueueItemStatusPending), Ladunek: ladunek,
			Warunek: tor.Condition, Priorytet: zlecenie.Priorytet,
			ZlecenieZrodloweID: &zlecenie.ID, PrzebiegID: zlecenie.PrzebiegID,
		})
		if err != nil {
			return shared.QueueItemBranchResponse{}, bladZlecenia(err)
		}
		tory = append(tory, zlecenieKontraktu(powstale))
	}
	return shared.QueueItemBranchResponse{Items: tory}, nil
}

// UwarunkujZlecenie ustala warunek przetworzenia zlecenia; warunek pusty
// zdejmuje warunek już zapisany.
func (a *adapterKolejek) UwarunkujZlecenie(ctx context.Context,
	z shared.QueueItemConditionRequest) (shared.QueueItemConditionResponse, error) {

	_, zlecenie, err := a.zlecenieZadania(ctx, z.QueueId, z.ItemId)
	if err != nil {
		return shared.QueueItemConditionResponse{}, err
	}
	var warunek *string
	if z.Condition != "" {
		tresc := z.Condition
		warunek = &tresc
	}
	if err := a.repozytorium.UstawWarunekZlecenia(ctx, zlecenie.ID, warunek); err != nil {
		return shared.QueueItemConditionResponse{}, bladZlecenia(err)
	}
	zapisane, err := a.repozytorium.Zlecenie(ctx, zlecenie.Kod)
	if err != nil {
		return shared.QueueItemConditionResponse{}, bladZlecenia(err)
	}
	return shared.QueueItemConditionResponse{Item: zlecenieKontraktu(zapisane)}, nil
}

// WykazZlecen oddaje zlecenia kolejki wraz z liczbą wszystkich, z opcjonalnym
// zawężeniem po stanie i limitem wierszy.
func (a *adapterKolejek) WykazZlecen(ctx context.Context,
	z shared.QueueItemListRequest) (shared.QueueItemListResponse, error) {

	kolejka, err := a.wierszKolejkiZlecen(ctx, z.QueueId)
	if err != nil {
		return shared.QueueItemListResponse{}, err
	}
	stan := ""
	if z.Status != nil {
		stan = stanZleceniaBazy(*z.Status)
	}
	wiersze, wszystkich, err := a.repozytorium.ZleceniaKolejki(ctx, kolejka.ID, stan,
		wartoscLiczby(z.Limit))
	if err != nil {
		return shared.QueueItemListResponse{}, bladZlecenia(err)
	}
	zlecenia := make([]shared.QueueItem, 0, len(wiersze))
	for _, wiersz := range wiersze {
		zlecenia = append(zlecenia, zlecenieKontraktu(wiersz))
	}
	return shared.QueueItemListResponse{Items: zlecenia, Total: wszystkich}, nil
}

// UstawPolitykeKolejki zapisuje zasięg, współbieżność, przepustowość
// i politykę ponawiania, nanosząc tylko pola obecne w żądaniu.
func (a *adapterKolejek) UstawPolitykeKolejki(ctx context.Context,
	z shared.QueuePolicySetRequest) (shared.QueuePolicySetResponse, error) {

	kolejka, err := a.wierszKolejkiZlecen(ctx, z.QueueId)
	if err != nil {
		return shared.QueuePolicySetResponse{}, err
	}
	zastana, err := a.repozytorium.PolitykaKolejki(ctx, kolejka.ID)
	if err != nil {
		return shared.QueuePolicySetResponse{}, bladZlecenia(err)
	}
	if err := a.repozytorium.ZapiszPolitykeKolejki(ctx, kolejka.ID,
		naniesPolityke(zastana, z.Policy)); err != nil {
		return shared.QueuePolicySetResponse{}, bladZlecenia(err)
	}
	zapisana, err := a.Kolejka(ctx, z.QueueId)
	if err != nil {
		return shared.QueuePolicySetResponse{}, err
	}
	return shared.QueuePolicySetResponse{Queue: zapisana}, nil
}

// naniesPolityke nakłada pola żądania na politykę zastaną. Pole nieobecne
// zostawia wartość zastaną: kontrakt ma same pola opcjonalne, więc żądanie
// zmieniające jedną liczbę nie może zerować dziewięciu pozostałych.
func naniesPolityke(zastana dane.PolitykaKolejki, zadana shared.QueuePolicy) dane.PolitykaKolejki {
	wynik := zastana
	if zadana.Scope != nil {
		wynik.Zasieg = zasiegBazy(*zadana.Scope)
	}
	if zadana.ScopeId != nil {
		wynik.ZasiegID = zadana.ScopeId
	}
	if zadana.MaxConcurrent != nil {
		wynik.LimitRownoleglych = *zadana.MaxConcurrent
	}
	if zadana.RatePerMinute != nil {
		wynik.TempoNaMinute = *zadana.RatePerMinute
	}
	if zadana.MaxAttempts != nil {
		wynik.LimitProb = *zadana.MaxAttempts
	}
	if zadana.Backoff != nil {
		wynik.Wycofanie = string(*zadana.Backoff)
	}
	if zadana.BackoffSeconds != nil {
		wynik.WycofanieSekundy = *zadana.BackoffSeconds
	}
	if zadana.Jitter != nil {
		wynik.Rozproszenie = *zadana.Jitter
	}
	if zadana.DeadLetterEnabled != nil {
		wynik.ZadaniaMartwe = *zadana.DeadLetterEnabled
	}
	if zadana.IdempotencyTtlSeconds != nil {
		wynik.IdempotencjaZycie = *zadana.IdempotencyTtlSeconds
	}
	return wynik
}

// WykazZadanMartwych oddaje zlecenia trwale nieudane — jednej kolejki, gdy
// wskazana, albo wszystkich kolejek rdzenia.
func (a *adapterKolejek) WykazZadanMartwych(ctx context.Context,
	z shared.QueueDeadListRequest) (shared.QueueDeadListResponse, error) {

	var kolejkaID int64
	if kod := wartoscTekstu(z.QueueId); kod != "" {
		wiersz, err := a.wierszKolejkiZlecen(ctx, kod)
		if err != nil {
			return shared.QueueDeadListResponse{}, err
		}
		kolejkaID = wiersz.ID
	}
	wiersze, err := a.repozytorium.ZleceniaMartwe(ctx, kolejkaID, wartoscLiczby(z.Limit))
	if err != nil {
		return shared.QueueDeadListResponse{}, bladZlecenia(err)
	}
	zlecenia := make([]shared.QueueItem, 0, len(wiersze))
	for _, wiersz := range wiersze {
		zlecenia = append(zlecenia, zlecenieKontraktu(wiersz))
	}
	return shared.QueueDeadListResponse{Items: zlecenia}, nil
}

// GlebokoscKolejki oddaje liczbę zleceń oczekujących w kolejnych odcinkach
// czasu, domyślnie liczonych godziną.
func (a *adapterKolejek) GlebokoscKolejki(ctx context.Context,
	z shared.QueueDepthGetRequest) (shared.QueueDepthGetResponse, error) {

	var kolejkaID int64
	if kod := wartoscTekstu(z.QueueId); kod != "" {
		wiersz, err := a.wierszKolejkiZlecen(ctx, kod)
		if err != nil {
			return shared.QueueDepthGetResponse{}, err
		}
		kolejkaID = wiersz.ID
	}
	odcinek := z.BucketSeconds
	if odcinek <= 0 {
		odcinek = odcinekGlebokosciDomyslny
	}
	wiersze, err := a.repozytorium.GlebokoscKolejki(ctx, kolejkaID, odcinek,
		znacznikChwiliAutomatyzacji(z.FromAt))
	if err != nil {
		return shared.QueueDepthGetResponse{}, bladZlecenia(err)
	}
	punkty := make([]shared.QueueDepthPoint, 0, len(wiersze))
	for _, wiersz := range wiersze {
		sekundy, err := strconv.ParseInt(wiersz.Chwila, 10, 64)
		if err != nil {
			continue
		}
		punkty = append(punkty, shared.QueueDepthPoint{
			At: sekundy * 1000, Pending: wiersz.Oczekujacych, Running: wiersz.Pracujacych,
		})
	}
	return shared.QueueDepthGetResponse{Points: punkty}, nil
}

// wierszKolejkiZlecen odnajduje kolejkę po jej identyfikatorze kontraktu,
// odmawiając, gdy identyfikator nie jest liczbą.
func (a *adapterKolejek) wierszKolejkiZlecen(ctx context.Context, id string) (dane.Kolejka, error) {
	numer, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return dane.Kolejka{}, bladWskazaniaZlecenia("kolejka o nieznanym identyfikatorze: " + id)
	}
	kolejka, err := a.repozytorium.PobierzKolejke(ctx, numer)
	if err != nil {
		return dane.Kolejka{}, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"kolejki: kolejka nie istnieje: "+id))
	}
	return kolejka, nil
}

// zlecenieZadania dobiera kolejkę i jej zlecenie wskazane żądaniem,
// sprawdzając, że zlecenie stoi w tej kolejce.
func (a *adapterKolejek) zlecenieZadania(ctx context.Context, idKolejki,
	idZlecenia string) (dane.Kolejka, dane.Zlecenie, error) {

	kolejka, err := a.wierszKolejkiZlecen(ctx, idKolejki)
	if err != nil {
		return dane.Kolejka{}, dane.Zlecenie{}, err
	}
	zlecenie, err := a.zlecenieKolejki(ctx, kolejka.ID, idZlecenia)
	if err != nil {
		return dane.Kolejka{}, dane.Zlecenie{}, err
	}
	return kolejka, zlecenie, nil
}

// zlecenieKolejki odnajduje zlecenie i sprawdza, czy stoi we wskazanej kolejce.
// Zlecenie z innej kolejki jest błędem żądania, nie cichym powodzeniem: żądanie
// mówi o parze (kolejka, zlecenie) i para ma się zgadzać.
func (a *adapterKolejek) zlecenieKolejki(ctx context.Context, kolejkaID int64,
	kod string) (dane.Zlecenie, error) {

	if kod == "" {
		return dane.Zlecenie{}, bladWskazaniaZlecenia("czynność wymaga wskazania zlecenia (itemId)")
	}
	zlecenie, err := a.repozytorium.Zlecenie(ctx, kod)
	if err != nil {
		return dane.Zlecenie{}, bladNieznanegoBytuAutomatyki(err, "zlecenie nie istnieje: ", kod)
	}
	if zlecenie.KolejkaID != kolejkaID {
		return dane.Zlecenie{}, bladWskazaniaZlecenia("zlecenie " + kod + " stoi w innej kolejce")
	}
	return zlecenie, nil
}

// ladunkiPodzadan rozbiera ładunki podziału. Kontrakt niesie je jednym zapisem
// strukturalnym typu `json`, a podział zakłada po jednym zleceniu na ładunek,
// więc zapis musi być wykazem.
func ladunkiPodzadan(zapis json.RawMessage) ([]json.RawMessage, error) {
	if len(zapis) == 0 {
		return nil, bladWskazaniaZlecenia("podział bez ładunków podzadań")
	}
	ladunki := []json.RawMessage{}
	if err := json.Unmarshal(zapis, &ladunki); err != nil {
		return nil, bladWskazaniaZlecenia(
			"ładunki podzadań mają być wykazem — każda pozycja zakłada jedno zlecenie")
	}
	if len(ladunki) == 0 {
		return nil, bladWskazaniaZlecenia("podział bez ładunków podzadań")
	}
	return ladunki, nil
}

// zapisLadunku przekłada ładunek kontraktu na kolumnę bazy; ładunek pusty
// w żądaniu daje w kolumnie brak, nie pusty ciąg.
func zapisLadunku(ladunek json.RawMessage) *string {
	if len(ladunek) == 0 {
		return nil
	}
	zapis := string(ladunek)
	return &zapis
}

// zlecenieKontraktu przekłada wiersz zlecenia zapisany w bazie na byt
// kontraktu zwracany w odpowiedzi komendy.
func zlecenieKontraktu(wiersz dane.Zlecenie) shared.QueueItem {
	priorytet, proby := wiersz.Priorytet, wiersz.Proby
	zlecenie := shared.QueueItem{
		Id: wiersz.Kod, QueueId: strconv.FormatInt(wiersz.KolejkaID, 10),
		Status: stanZleceniaKontraktu(wiersz.Stan), Priority: &priorytet, Attempts: &proby,
		Condition: wiersz.Warunek, IdempotencyKey: wiersz.KluczIdempotencji,
		ExecutionId: wiersz.PrzebiegKod,
		CreatedAt:   chwilaBazy(wiersz.Utworzono), UpdatedAt: chwilaBazy(wiersz.Zaktualizowano),
	}
	if wiersz.Ladunek != nil {
		zlecenie.Payload = json.RawMessage(*wiersz.Ladunek)
	}
	if wiersz.Termin != nil {
		termin := chwilaBazy(*wiersz.Termin)
		zlecenie.ScheduledAt = &termin
	}
	return zlecenie
}

// stanyZlecenia wiąże stan kontraktu ze słownikiem bazy. Mapa jest zupełna po
// wyliczeniu: stan dopisany do kontraktu zatrzyma się tutaj, zamiast po cichu
// zejść na stan domyślny.
var stanyZlecenia = map[shared.QueueItemStatus]string{
	shared.QueueItemStatusPending:   "oczekuje",
	shared.QueueItemStatusDelayed:   "odlozone",
	shared.QueueItemStatusRunning:   "przetwarzane",
	shared.QueueItemStatusSucceeded: "zakonczone",
	shared.QueueItemStatusFailed:    "bledne",
	shared.QueueItemStatusDead:      "martwe",
	shared.QueueItemStatusRemoved:   "zdjete",
}

// stanZleceniaBazy przekłada stan kontraktu na słownik bazy, a stan nieznany
// zapisuje jako oczekujący.
func stanZleceniaBazy(stan shared.QueueItemStatus) string {
	if kolumna, jest := stanyZlecenia[stan]; jest {
		return kolumna
	}
	return "oczekuje"
}

// stanZleceniaKontraktu przekłada słownik bazy na stan kontraktu, a wartość
// nieznaną zwraca jako oczekującą.
func stanZleceniaKontraktu(kolumna string) shared.QueueItemStatus {
	for stan, wartosc := range stanyZlecenia {
		if wartosc == kolumna {
			return stan
		}
	}
	return shared.QueueItemStatusPending
}

// czyStanKoncowyZlecenia mówi, czy ze zleceniem nic już się nie stanie — stan
// końcowy nie przyjmuje kolejnej zmiany.
func czyStanKoncowyZlecenia(kolumna string) bool {
	switch kolumna {
	case "zakonczone", "bledne", "martwe", "zdjete":
		return true
	default:
		return false
	}
}

// zasiegiKolejki wiąże zasięg kontraktu ze słownikiem bazy, jedną mapą
// czytaną w obie strony przekładu.
var zasiegiKolejki = map[shared.QueueScope]string{
	shared.QueueScopeGlobal:  "globalna",
	shared.QueueScopeLocal:   "lokalna",
	shared.QueueScopeModel:   "modelu",
	shared.QueueScopeAgent:   "eksperta",
	shared.QueueScopeProject: "projektu",
}

// zasiegBazy przekłada zasięg kontraktu na słownik bazy, a zasięg nieznany
// zapisuje jako zasięg lokalny.
func zasiegBazy(zasieg shared.QueueScope) string {
	if kolumna, jest := zasiegiKolejki[zasieg]; jest {
		return kolumna
	}
	return "lokalna"
}

// zasiegKontraktu przekłada słownik bazy na zasięg kontraktu, a wartość
// nieznaną zwraca jako zasięg lokalny.
func zasiegKontraktu(kolumna string) shared.QueueScope {
	for zasieg, wartosc := range zasiegiKolejki {
		if wartosc == kolumna {
			return zasieg
		}
	}
	return shared.QueueScopeLocal
}

// politykaKontraktu przekłada politykę kolejki zapisaną w bazie na byt
// kontraktu zwracany w odpowiedzi.
func politykaKontraktu(wiersz dane.PolitykaKolejki) shared.QueuePolicy {
	zasieg := zasiegKontraktu(wiersz.Zasieg)
	wycofanie := shared.QueueBackoffKind(wiersz.Wycofanie)
	limit, tempo := wiersz.LimitRownoleglych, wiersz.TempoNaMinute
	proby, sekundy := wiersz.LimitProb, wiersz.WycofanieSekundy
	rozproszenie, martwe := wiersz.Rozproszenie, wiersz.ZadaniaMartwe
	zycie := wiersz.IdempotencjaZycie
	return shared.QueuePolicy{
		Scope: &zasieg, ScopeId: wiersz.ZasiegID, MaxConcurrent: &limit,
		RatePerMinute: &tempo, MaxAttempts: &proby, Backoff: &wycofanie,
		BackoffSeconds: &sekundy, Jitter: &rozproszenie, DeadLetterEnabled: &martwe,
		IdempotencyTtlSeconds: &zycie,
	}
}

// bladZlecenia znakuje usterkę zapisu zlecenia kodem błędu kontraktu, a błąd
// pusty przepuszcza bez zmian.
func bladZlecenia(err error) error {
	if err == nil {
		return nil
	}
	return protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeInternalError, err))
}

// bladWskazaniaZlecenia nazywa brak danych wymaganych w żądaniu jako błąd
// wołającego, nie odmowę rdzenia.
func bladWskazaniaZlecenia(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed, "kolejki: "+powod))
}
