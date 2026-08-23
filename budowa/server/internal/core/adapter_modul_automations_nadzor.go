// Odpowiedzialność pliku: pięć czynności rodziny `schedule.*` należących do
// okna Scheduler — okna wykonania, uruchomienie wsteczne, historia wyzwoleń,
// nadzór obecności uruchomień i adres wejściowy wyzwalacza webhook.
//
// Wszystkie pięć zapisują stan DO BAZY, bo wszystkie pięć muszą obowiązywać po
// ponownym złożeniu rdzenia. Okno wykonania trzymane w pamięci procesu
// przestałoby ograniczać budzik po pierwszym restarcie, a Scheduler pokazywałby
// ograniczenie, które już niczego nie ogranicza. Klucz podpisu webhooka
// w pamięci unieważniałby każde wywołanie przychodzące po restarcie.
//
// Klucz podpisu leży w sejfie, nie w bazie i nie w odpowiedzi. Kontrakt oddaje
// `signatureSecretRef` — „referencja klucza podpisu w skarbcu; nigdy sama
// wartość”. Wymiana klucza (`rotateSecret`) zapisuje w sejfie wartość nową
// i oddaje referencję nową; starej nie da się odczytać ani przez tę komendę,
// ani przez żadną inną.
package core

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// oknoDeduplikacjiDomyslne odpowiada wartości domyślnej kolumny schematu —
// pięć minut, w których to samo wywołanie przychodzące nie ruszy automatyki
// dwa razy.
const oknoDeduplikacjiDomyslne = 300

// granicaPrzebiegowWstecznych chroni przed uruchomieniem wstecznym, które
// zakolejkowałoby tysiące przebiegów z jednego nieuważnego zakresu dat.
// Kontrakt zna pole `maxRuns`; ta granica obowiązuje, gdy Operator go nie podał.
const granicaPrzebiegowWstecznych = 100

// UstawOknaWykonania zapisuje przedziały czasu, w których uruchomienie następuje.
func (a *adapterAutomatyk) UstawOknaWykonania(ctx context.Context,
	z shared.ScheduleWindowSetRequest) (shared.ScheduleWindowSetResponse, error) {

	harmonogram, err := a.harmonogramZadania(ctx, z.ScheduleId)
	if err != nil {
		return shared.ScheduleWindowSetResponse{}, err
	}
	okna := make([]dane.OknoWykonania, 0, len(z.Windows))
	for _, okno := range z.Windows {
		dni := zapisStrukturalny(okno.DaysOfWeek)
		zapisDni := "[]"
		if dni != nil {
			zapisDni = *dni
		}
		okna = append(okna, dane.OknoWykonania{
			DniTygodnia: zapisDni, MinutaOd: okno.FromMinuteOfDay,
			MinutaDo: okno.ToMinuteOfDay, StrefaCzasowa: okno.TimeZone,
		})
	}
	if err := a.repozytorium.ZapiszOknaWykonania(ctx, harmonogram.ID, okna); err != nil {
		return shared.ScheduleWindowSetResponse{}, bladAutomatyki(err)
	}
	a.zapisAudytu(ctx, &harmonogram.AutomatykaID, "zapis okien wykonania harmonogramu",
		map[string]any{"harmonogram": harmonogram.Kod, "okien": len(okna)})
	zapisany, err := a.harmonogramPoZmianie(ctx, harmonogram.Kod)
	if err != nil {
		return shared.ScheduleWindowSetResponse{}, err
	}
	return shared.ScheduleWindowSetResponse{Schedule: zapisany}, nil
}

// UruchomWstecznie wykonuje przebiegi dla przeszłych, pominiętych terminów.
// Termin, dla którego przebieg już był, zostaje pominięty i policzony — nie
// odtwarza się pracy, która się odbyła.
func (a *adapterAutomatyk) UruchomWstecznie(ctx context.Context,
	z shared.ScheduleBackfillRunRequest) (shared.ScheduleBackfillRunResponse, error) {

	harmonogram, err := a.harmonogramZadania(ctx, z.ScheduleId)
	if err != nil {
		return shared.ScheduleBackfillRunResponse{}, err
	}
	automatyka, err := a.repozytorium.AutomatykaPoID(ctx, harmonogram.AutomatykaID)
	if err != nil {
		return shared.ScheduleBackfillRunResponse{}, bladAutomatyki(err)
	}
	wyzwalacze, err := a.repozytorium.Wyzwalacze(ctx, harmonogram.ID)
	if err != nil {
		wyzwalacze = nil
	}
	granica := wartoscLiczbyLub(z.MaxRuns, granicaPrzebiegowWstecznych)
	terminy := terminyWsteczne(harmonogram, wyzwalacze, z.FromAt, z.ToAt, granica)

	zaszle, err := a.chwileWyzwolen(ctx, harmonogram)
	if err != nil {
		return shared.ScheduleBackfillRunResponse{}, err
	}
	zakolejkowane := []int64{}
	pominiete := 0
	for _, termin := range terminy {
		if zaszle[termin.Format(formatZnacznikaBazy)] {
			pominiete++
			continue
		}
		if _, err := a.UruchomAutomatyke(ctx, automatyka); err != nil {
			return shared.ScheduleBackfillRunResponse{}, err
		}
		a.odnotujWyzwolenie(ctx, harmonogram, shared.AutomationTriggerCauseBackfill)
		zakolejkowane = append(zakolejkowane, termin.UnixMilli())
	}
	a.zapisAudytu(ctx, &automatyka.ID, "uruchomienie wsteczne harmonogramu", map[string]any{
		"harmonogram": harmonogram.Kod, "zakolejkowano": len(zakolejkowane), "pominieto": pominiete,
	})
	return shared.ScheduleBackfillRunResponse{QueuedAt: zakolejkowane, Skipped: pominiete}, nil
}

// terminyWsteczne wylicza terminy cykliczności mieszczące się w zakresie.
// Rachunek idzie tą samą maszynerią, którą liczy się najbliższe uruchomienie —
// drugiego czytnika notacji cron w rdzeniu nie ma.
func terminyWsteczne(harmonogram dane.Harmonogram, wyzwalacze []dane.WyzwalaczAutomatyki,
	odMilisekund, doMilisekund int64, granica int) []time.Time {

	terminy := []time.Time{}
	chwila := time.UnixMilli(odMilisekund).UTC()
	koniec := time.UnixMilli(doMilisekund).UTC()
	for len(terminy) < granica {
		// Harmonogram wyłączony też ma terminy przeszłe: uruchomienie wsteczne
		// jest jawnym poleceniem Operatora, a nie pracą budzika, więc stan
		// przełącznika cykliczności go nie wstrzymuje.
		znacznik := chwilaNastepnegoUruchomienia(harmonogram.Cron, wyzwalacze, true, chwila)
		if znacznik == nil {
			break
		}
		nastepny, err := time.Parse(formatZnacznikaBazy, *znacznik)
		if err != nil || nastepny.After(koniec) {
			break
		}
		terminy = append(terminy, nastepny)
		chwila = nastepny
	}
	return terminy
}

// chwileWyzwolen zbiera znaczniki wyzwoleń, które już zaszły — po nich poznaje
// się termin, dla którego przebieg już był.
func (a *adapterAutomatyk) chwileWyzwolen(ctx context.Context,
	harmonogram dane.Harmonogram) (map[string]bool, error) {

	wiersze, err := a.repozytorium.Wyzwolenia(ctx, 0, harmonogram.ID, 0)
	if err != nil {
		return nil, bladAutomatyki(err)
	}
	zaszle := make(map[string]bool, len(wiersze))
	for _, wiersz := range wiersze {
		zaszle[wiersz.Chwila] = true
	}
	return zaszle, nil
}

// HistoriaWyzwolen oddaje rzeczywiste momenty wyzwolenia wraz z przyczyną.
func (a *adapterAutomatyk) HistoriaWyzwolen(ctx context.Context,
	z shared.ScheduleTriggerHistoryRequest) (shared.ScheduleTriggerHistoryResponse, error) {

	var automatykaID, harmonogramID int64
	if kod := wartoscTekstu(z.WorkflowId); kod != "" {
		wiersz, err := a.wiersz(ctx, kod)
		if err != nil {
			return shared.ScheduleTriggerHistoryResponse{}, err
		}
		automatykaID = wiersz.ID
	}
	if kod := wartoscTekstu(z.ScheduleId); kod != "" {
		harmonogram, err := a.repozytorium.HarmonogramPoKodzie(ctx, kod)
		if err != nil {
			return shared.ScheduleTriggerHistoryResponse{},
				bladNieznanegoBytuAutomatyki(err, "harmonogram nie istnieje: ", kod)
		}
		harmonogramID = harmonogram.ID
	}
	wiersze, err := a.repozytorium.Wyzwolenia(ctx, automatykaID, harmonogramID, wartoscLiczby(z.Limit))
	if err != nil {
		return shared.ScheduleTriggerHistoryResponse{}, bladAutomatyki(err)
	}
	wpisy := make([]shared.AutomationTriggerHistoryEntry, 0, len(wiersze))
	for _, wiersz := range wiersze {
		wpisy = append(wpisy, shared.AutomationTriggerHistoryEntry{
			WorkflowId: wiersz.AutomatykaKod, ScheduleId: wiersz.HarmonogramKod,
			Cause:       shared.AutomationTriggerCause(wiersz.Przyczyna),
			ExecutionId: wiersz.PrzebiegKod, At: chwilaBazy(wiersz.Chwila),
		})
	}
	return shared.ScheduleTriggerHistoryResponse{Entries: wpisy}, nil
}

// UstawNadzorUruchomien zapisuje okno tolerancji nadzoru obecności uruchomień.
// Zero wyłącza nadzór — tak mówi kontrakt.
func (a *adapterAutomatyk) UstawNadzorUruchomien(ctx context.Context,
	z shared.ScheduleHeartbeatSetRequest) (shared.ScheduleHeartbeatSetResponse, error) {

	harmonogram, err := a.harmonogramZadania(ctx, z.ScheduleId)
	if err != nil {
		return shared.ScheduleHeartbeatSetResponse{}, err
	}
	err = a.repozytorium.UstawNadzorHarmonogramu(ctx, harmonogram.ID,
		z.ToleranceSeconds, z.AlertRuleId)
	if err != nil {
		return shared.ScheduleHeartbeatSetResponse{}, bladAutomatyki(err)
	}
	a.zapisAudytu(ctx, &harmonogram.AutomatykaID, "zapis nadzoru obecności uruchomień",
		map[string]any{"harmonogram": harmonogram.Kod, "tolerancja": z.ToleranceSeconds})
	zapisany, err := a.harmonogramPoZmianie(ctx, harmonogram.Kod)
	if err != nil {
		return shared.ScheduleHeartbeatSetResponse{}, err
	}
	return shared.ScheduleHeartbeatSetResponse{Schedule: zapisany}, nil
}

// AdresWebhooka oddaje unikatowy adres wejściowy automatyki wraz z referencją
// klucza podpisu. Wartość klucza nie wraca ani tu, ani nigdzie indziej.
func (a *adapterAutomatyk) AdresWebhooka(ctx context.Context,
	z shared.ScheduleWebhookEndpointGetRequest) (shared.ScheduleWebhookEndpointGetResponse, error) {

	wiersz, err := a.wiersz(ctx, z.WorkflowId)
	if err != nil {
		return shared.ScheduleWebhookEndpointGetResponse{}, err
	}
	harmonogram, err := a.repozytorium.Harmonogram(ctx, wiersz.ID)
	if err != nil {
		return shared.ScheduleWebhookEndpointGetResponse{}, bladNieznanegoBytuAutomatyki(err,
			"automatyka nie ma harmonogramu, więc nie ma wyzwalacza webhook: ", wiersz.Kod)
	}
	wymiana := z.RotateSecret != nil && *z.RotateSecret
	odwolanie, err := a.kluczPodpisu(ctx, harmonogram, wymiana)
	if err != nil {
		return shared.ScheduleWebhookEndpointGetResponse{}, err
	}
	okno := harmonogram.OknoDeduplikacjiSek
	if okno <= 0 {
		okno = oknoDeduplikacjiDomyslne
	}
	if wymiana {
		a.zapisAudytu(ctx, &wiersz.ID, "wymiana klucza podpisu webhooka",
			map[string]any{"harmonogram": harmonogram.Kod})
	}
	return shared.ScheduleWebhookEndpointGetResponse{
		EndpointUrl:                adresWejsciowyAutomatyki(wiersz.Kod),
		SignatureSecretRef:         odwolanie,
		DeduplicationWindowSeconds: okno,
	}, nil
}

// kluczPodpisu oddaje referencję klucza HMAC: zastaną albo nową. Klucz powstaje
// z generatora kryptograficznego biblioteki standardowej — programu zewnętrznego
// do tego nie potrzeba i nie wolno go tu wprowadzać.
func (a *adapterAutomatyk) kluczPodpisu(ctx context.Context, harmonogram dane.Harmonogram,
	wymiana bool) (string, error) {

	if !wymiana && harmonogram.OdwolaniePodpisu != nil && *harmonogram.OdwolaniePodpisu != "" {
		return *harmonogram.OdwolaniePodpisu, nil
	}
	if a.sejf == nil {
		return "", errBrakSkarbca
	}
	surowy := make([]byte, 32)
	if _, err := rand.Read(surowy); err != nil {
		return "", bladAutomatyki(err)
	}
	byt := nowyIdentyfikator(przedrostekPoswiadczenia + "podpis-")
	odwolanie, err := a.sejf.Zapisz(ctx, byt, hex.EncodeToString(surowy))
	if err != nil {
		return "", bladAutomatyki(err)
	}
	if err := a.repozytorium.UstawPodpisHarmonogramu(ctx, harmonogram.ID, odwolanie); err != nil {
		_ = a.sejf.Usun(ctx, byt)
		return "", bladAutomatyki(err)
	}
	return odwolanie, nil
}

// adresWejsciowyAutomatyki składa adres wejściowy wyzwalacza webhook. Ścieżka
// jest względna, bo rdzeń nie zna nazwy hosta, pod którą Operator go wystawił —
// adres bezwzględny zmyślony przez rdzeń byłby adresem, pod który nic nie dojdzie.
func adresWejsciowyAutomatyki(kodAutomatyki string) string {
	return "/automations/webhook/" + kodAutomatyki
}

// odnotujWyzwolenie nanosi rzeczywisty moment wyzwolenia. Nieudany zapis nie
// wywraca uruchomienia — przebieg już ruszył.
func (a *adapterAutomatyk) odnotujWyzwolenie(ctx context.Context, harmonogram dane.Harmonogram,
	przyczyna shared.AutomationTriggerCause) {

	id := harmonogram.ID
	_ = a.repozytorium.DopiszWyzwolenie(ctx, harmonogram.AutomatykaID, &id, string(przyczyna), nil)
}

// harmonogramZadania odnajduje harmonogram wskazany żądaniem rodziny
// `schedule.*` i nazywa brak wprost.
func (a *adapterAutomatyk) harmonogramZadania(ctx context.Context,
	kod string) (dane.Harmonogram, error) {

	if kod == "" {
		return dane.Harmonogram{}, bladWskazaniaAutomatyki("komenda bez wskazania harmonogramu")
	}
	harmonogram, err := a.repozytorium.HarmonogramPoKodzie(ctx, kod)
	if err != nil {
		return dane.Harmonogram{}, bladNieznanegoBytuAutomatyki(err, "harmonogram nie istnieje: ", kod)
	}
	return harmonogram, nil
}

// harmonogramPoZmianie oddaje harmonogram odczytany NA NOWO wraz z automatyką,
// której dotyczy. Wiersz zastany sprzed zapisu niósłby stan sprzed niego.
func (a *adapterAutomatyk) harmonogramPoZmianie(ctx context.Context,
	kod string) (shared.AutomationSchedule, error) {

	harmonogram, err := a.repozytorium.HarmonogramPoKodzie(ctx, kod)
	if err != nil {
		return shared.AutomationSchedule{}, bladNieznanegoBytuAutomatyki(err,
			"harmonogram nie istnieje: ", kod)
	}
	automatyka, err := a.repozytorium.AutomatykaPoID(ctx, harmonogram.AutomatykaID)
	if err != nil {
		return shared.AutomationSchedule{}, bladAutomatyki(err)
	}
	zlozony, err := a.harmonogramKontraktu(ctx, automatyka.Kod, harmonogram)
	if err != nil {
		return shared.AutomationSchedule{}, bladAutomatyki(err)
	}
	return zlozony, nil
}
