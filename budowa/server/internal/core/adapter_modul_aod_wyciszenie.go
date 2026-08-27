// Plik obsługuje cztery komendy wyciszenia nakładki Always On Display:
// aod.mute.get, aod.mute.set, aod.signal.report oraz aod.signal.list,
// prowadzone w rdzeniu niezależnie od okna nakładki.
package core

import (
	"context"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

const (
	// przedrostekWyciszeniaNakladki znakuje identyfikator wyciszenia nadany przez rdzeń,
	// odróżniając go od identyfikatora sygnału przy odczycie dziennika zdarzeń.
	przedrostekWyciszeniaNakladki = "wyc-"
	// przedrostekSygnaluNakladki znakuje identyfikator sygnału nadany przez rdzeń,
	// odróżniając go od identyfikatora wyciszenia przy odczycie dziennika zdarzeń.
	przedrostekSygnaluNakladki = "syg-"
)

// ZWyciszeniami wpina magazyn wyciszeń i sygnałów nakładki. Bez niego cztery
// komendy wyciszenia odmawiają z nazwą niewpiętego składnika — pusty wykaz
// wyciszeń byłby ciszą mówiącą Operatorowi, że nic nie jest wyciszone.
func (a *adapterNakladkiAod) ZWyciszeniami(
	w dane.RepozytoriumWyciszenNakladki) *adapterNakladkiAod {

	a.wyciszenia = w
	return a
}

// ── aod.mute.get ────────────────────────────────────────────────────────────

// WyciszeniaNakladki obsługuje `aod.mute.get`: oddaje wyciszenia czynne w tej
// chwili. Wyciszenie czasowe, którego czas upłynął, do wykazu nie wchodzi
// i przestaje istnieć — Operator nie ma odklikiwać ciszy, która się skończyła.
func (a *adapterNakladkiAod) WyciszeniaNakladki(ctx context.Context,
	z shared.AodMuteGetRequest) (shared.AodMuteGetResponse, error) {

	if _, err := a.urzadzenieNakladki(ctx, z.DeviceId); err != nil {
		return shared.AodMuteGetResponse{}, err
	}
	wykaz, err := a.wyciszeniaCzynne(ctx)
	if err != nil {
		return shared.AodMuteGetResponse{}, err
	}
	return shared.AodMuteGetResponse{Mutes: wyciszeniaKontraktu(wykaz)}, nil
}

// ── aod.mute.set ────────────────────────────────────────────────────────────

// PrzestawWyciszenieNakladki obsługuje `aod.mute.set`: zakłada wyciszenie
// nakładki albo je znosi, wskazane identyfikatorem albo rodzajem i zakresem,
// i oddaje wykaz wyciszeń po zmianie razem z polem `changed`.
func (a *adapterNakladkiAod) PrzestawWyciszenieNakladki(ctx context.Context,
	z shared.AodMuteSetRequest) (shared.AodMuteSetResponse, shared.AodMute, error) {

	urzadzenie, err := a.urzadzenieNakladki(ctx, z.DeviceId)
	if err != nil {
		return shared.AodMuteSetResponse{}, shared.AodMute{}, err
	}
	if err := a.sprawdzMagazynWyciszen(); err != nil {
		return shared.AodMuteSetResponse{}, shared.AodMute{}, err
	}
	if z.Muted {
		return a.zalozWyciszenie(ctx, z, urzadzenie)
	}
	return a.zniesWyciszenie(ctx, z)
}

// zalozWyciszenie nadaje wyciszeniu identyfikator, wiąże je z urządzeniem
// żądania, zapisuje je w magazynie i składa odpowiedź komendy z wykazem
// wyciszeń po zmianie.
func (a *adapterNakladkiAod) zalozWyciszenie(ctx context.Context,
	z shared.AodMuteSetRequest, urzadzenie string) (shared.AodMuteSetResponse, shared.AodMute, error) {

	wzor, err := wzorWyciszenia(z)
	if err != nil {
		return shared.AodMuteSetResponse{}, shared.AodMute{}, err
	}
	wzor.Identyfikator = nowyIdentyfikator(przedrostekWyciszeniaNakladki)
	wzor.Urzadzenie = urzadzenie
	zapisane, zmiana, err := a.wyciszenia.ZapiszWyciszenieNakladki(ctx, wzor)
	if err != nil {
		return shared.AodMuteSetResponse{}, shared.AodMute{}, protocol.JakoError(
			protocol.BladZeZrodla(shared.ErrorCodeValidationFailed, err))
	}
	wykaz, err := a.wyciszeniaCzynne(ctx)
	if err != nil {
		return shared.AodMuteSetResponse{}, shared.AodMute{}, err
	}
	return shared.AodMuteSetResponse{
		Mutes:   wyciszeniaKontraktu(wykaz),
		Changed: zmiana,
	}, wyciszenieKontraktu(zapisane), nil
}

// zniesWyciszenie znosi wyciszenie wskazane identyfikatorem albo rodzajem
// i zakresem, po czym składa odpowiedź komendy z wykazem wyciszeń czynnych
// po zmianie.
func (a *adapterNakladkiAod) zniesWyciszenie(ctx context.Context,
	z shared.AodMuteSetRequest) (shared.AodMuteSetResponse, shared.AodMute, error) {

	znoszone, err := a.wyciszenieDoZniesienia(ctx, z)
	if err != nil {
		return shared.AodMuteSetResponse{}, shared.AodMute{}, err
	}
	zniesione, err := a.wyciszenia.ZniesWyciszenieNakladki(ctx, znoszone.Identyfikator)
	if err != nil {
		return shared.AodMuteSetResponse{}, shared.AodMute{}, err
	}
	wykaz, err := a.wyciszeniaCzynne(ctx)
	if err != nil {
		return shared.AodMuteSetResponse{}, shared.AodMute{}, err
	}
	return shared.AodMuteSetResponse{
		Mutes:   wyciszeniaKontraktu(wykaz),
		Changed: zniesione,
	}, wyciszenieKontraktu(znoszone), nil
}

// wyciszenieDoZniesienia odnajduje wyciszenie znoszone żądaniem. Wyciszenie,
// którego nie ma, kończy się odmową `not_found` — wykaz oddany bez zmiany
// wyglądałby na wykonaną czynność.
func (a *adapterNakladkiAod) wyciszenieDoZniesienia(ctx context.Context,
	z shared.AodMuteSetRequest) (dane.WyciszenieNakladki, error) {

	if identyfikator := wartoscTekstu(z.MuteId); identyfikator != "" {
		wykaz, err := a.wyciszeniaCzynne(ctx)
		if err != nil {
			return dane.WyciszenieNakladki{}, err
		}
		for _, wyciszenie := range wykaz {
			if wyciszenie.Identyfikator == identyfikator {
				return wyciszenie, nil
			}
		}
		return dane.WyciszenieNakladki{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeNotFound, "nakładka AOD: wyciszenia "+identyfikator+
				" nie ma; nie ma czego znieść"))
	}
	wzor, err := wzorWyciszenia(z)
	if err != nil {
		return dane.WyciszenieNakladki{}, err
	}
	zastane, jest, err := a.wyciszenia.WyciszenieNakladkiPoBycie(ctx, wzor)
	if err != nil {
		return dane.WyciszenieNakladki{}, err
	}
	if !jest {
		return dane.WyciszenieNakladki{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeNotFound, "nakładka AOD: "+opisBytuWyciszenia(wzor)+
				" nie jest wyciszone; nie ma czego znieść"))
	}
	return zastane, nil
}

// ── aod.signal.report ───────────────────────────────────────────────────────

// ZglosSygnalNakladki obsługuje `aod.signal.report`: odkłada w rdzeniu sygnał
// klasy zdarzeń wyzwalających i mówi, czy sygnał wpadł w wyciszenie czynne
// w tej chwili.
func (a *adapterNakladkiAod) ZglosSygnalNakladki(ctx context.Context,
	z shared.AodSignalReportRequest) (shared.AodSignalReportResponse, error) {

	if err := a.sprawdzMagazynWyciszen(); err != nil {
		return shared.AodSignalReportResponse{}, err
	}
	if z.EventClass == "" {
		return shared.AodSignalReportResponse{}, bladZadaniaNakladki(
			"sygnał bez wskazania klasy zdarzeń — klasa mówi, co się stało i co da się wyciszyć")
	}
	if z.Text == "" {
		return shared.AodSignalReportResponse{}, bladZadaniaNakladki("sygnał bez treści")
	}
	wiersz := dane.SygnalNakladki{
		Identyfikator:   nowyIdentyfikator(przedrostekSygnaluNakladki),
		KlasaZdarzen:    z.EventClass,
		Tresc:           z.Text,
		ModulKod:        wartoscTekstu(z.ModuleId),
		SesjaKod:        wartoscTekstu(z.SessionId),
		LiczbaWystapien: wartoscLiczby(z.OccurrenceCount),
	}
	zapisany, err := a.wyciszenia.ZapiszSygnalNakladki(ctx, wiersz)
	if err != nil {
		return shared.AodSignalReportResponse{}, protocol.JakoError(
			protocol.BladZeZrodla(shared.ErrorCodeValidationFailed, err))
	}
	czynne, err := a.wyciszeniaCzynne(ctx)
	if err != nil {
		return shared.AodSignalReportResponse{}, err
	}
	odpowiedz := shared.AodSignalReportResponse{Signal: sygnalKontraktu(zapisany)}
	if wyciszenie, wstrzymany := wyciszenieObejmujaceSygnal(czynne, zapisany); wstrzymany {
		obejmujace := wyciszenieKontraktu(wyciszenie)
		odpowiedz.Suppressed = true
		odpowiedz.Mute = &obejmujace
	}
	return odpowiedz, nil
}

// ── aod.signal.list ────────────────────────────────────────────────────────

// SygnalyNakladki obsługuje `aod.signal.list`: oddaje sygnały ujawnialne,
// a wstrzymane nazywa liczbą i wyciszeniem, które je wstrzymało.
func (a *adapterNakladkiAod) SygnalyNakladki(ctx context.Context,
	z shared.AodSignalListRequest) (shared.AodSignalListResponse, error) {

	if _, err := a.urzadzenieNakladki(ctx, z.DeviceId); err != nil {
		return shared.AodSignalListResponse{}, err
	}
	if err := a.sprawdzMagazynWyciszen(); err != nil {
		return shared.AodSignalListResponse{}, err
	}
	czynne, err := a.wyciszeniaCzynne(ctx)
	if err != nil {
		return shared.AodSignalListResponse{}, err
	}
	// Zawężenie do klas idzie po odczycie, nie zapytaniem — zapytanie na
	// jedną klasę zwęziłoby wykaz.
	wiersze, err := a.wyciszenia.SygnalyNakladki(ctx, dane.SygnalNakladki{
		ModulKod: wartoscTekstu(z.ModuleId),
		SesjaKod: wartoscTekstu(z.SessionId),
	}, 0)
	if err != nil {
		return shared.AodSignalListResponse{}, err
	}
	zadaneKlasy := map[shared.AodEventClass]struct{}{}
	for _, klasa := range z.Classes {
		zadaneKlasy[klasa] = struct{}{}
	}
	granica := wartoscLiczby(z.Limit)
	sygnaly := make([]shared.AodSignal, 0, len(wiersze))
	wstrzymujace := map[string]dane.WyciszenieNakladki{}
	wstrzymanych := 0
	for _, wiersz := range wiersze {
		if len(zadaneKlasy) > 0 {
			if _, zadana := zadaneKlasy[wiersz.KlasaZdarzen]; !zadana {
				continue
			}
		}
		if wyciszenie, wstrzymany := wyciszenieObejmujaceSygnal(czynne, wiersz); wstrzymany {
			wstrzymanych++
			wstrzymujace[wyciszenie.Identyfikator] = wyciszenie
			continue
		}
		if granica > 0 && len(sygnaly) == granica {
			continue
		}
		sygnaly = append(sygnaly, sygnalKontraktu(wiersz))
	}
	// Wyciszenia idą w kolejności wykazu czynnych, nie napotkania, żeby
	// odpowiedź była powtarzalna.
	wyciszenia := make([]shared.AodMute, 0, len(wstrzymujace))
	for _, wyciszenie := range czynne {
		if _, wstrzymalo := wstrzymujace[wyciszenie.Identyfikator]; wstrzymalo {
			wyciszenia = append(wyciszenia, wyciszenieKontraktu(wyciszenie))
		}
	}
	return shared.AodSignalListResponse{
		Signals:         sygnaly,
		SuppressedCount: wstrzymanych,
		Mutes:           wyciszenia,
	}, nil
}

// ── reguła wyciszenia ──────────────────────────────────────────────────────

// wyciszenieObejmujaceSygnal wskazuje, które z wyciszeń czynnych wstrzymuje
// dany sygnał, uwzględniając rodzaj wyciszenia, zakres oraz klasę zdarzeń
// sygnału. Regułę powtarza client/src/aod/wyciszenie-aod.ts po stronie okna.
func wyciszenieObejmujaceSygnal(czynne []dane.WyciszenieNakladki,
	sygnal dane.SygnalNakladki) (dane.WyciszenieNakladki, bool) {

	for _, wyciszenie := range czynne {
		switch wyciszenie.Rodzaj {
		case shared.AodMuteKindTimed:
			// Wyciszenie czasowe obejmuje każdy sygnał, niezależnie od modułu
			// i klasy zdarzeń.
			return wyciszenie, true
		case shared.AodMuteKindContextual:
			byt := sygnal.ModulKod
			if wyciszenie.Zakres == shared.AodMuteScopeSession {
				byt = sygnal.SesjaKod
			}
			if byt != "" && byt == wyciszenie.KluczZakresu {
				return wyciszenie, true
			}
		case shared.AodMuteKindEventClass:
			if sygnal.KlasaZdarzen == wyciszenie.KlasaZdarzen {
				return wyciszenie, true
			}
		}
	}
	return dane.WyciszenieNakladki{}, false
}

// ── ustalenia wspólne ──────────────────────────────────────────────────────

// wyciszeniaCzynne odczytuje wyciszenia czynne o chwili bieżącej, pomijając
// wyciszenia czasowe, których czas trwania już upłynął.
func (a *adapterNakladkiAod) wyciszeniaCzynne(
	ctx context.Context) ([]dane.WyciszenieNakladki, error) {

	if err := a.sprawdzMagazynWyciszen(); err != nil {
		return nil, err
	}
	return a.wyciszenia.WyciszeniaNakladki(ctx, time.Now().UTC().Format(formatZnacznikaBazy))
}

// wzorWyciszenia składa wzór wyciszenia z żądania i odmawia, gdy żądanie nie
// wskazuje kompletu potrzebnego rodzajowi. Odmowa nazywa brak, bo wyciszenie
// niepełne wyciszyłoby coś innego, niż Operator prosił.
func wzorWyciszenia(z shared.AodMuteSetRequest) (dane.WyciszenieNakladki, error) {
	rodzaj := shared.AodMuteKind("")
	if z.Kind != nil {
		rodzaj = *z.Kind
	}
	wzor := dane.WyciszenieNakladki{Rodzaj: rodzaj, NazwaZakresu: wartoscTekstu(z.ScopeName)}
	if z.Scope != nil {
		wzor.Zakres = *z.Scope
	}
	if z.EventClass != nil {
		wzor.KlasaZdarzen = *z.EventClass
	}
	wzor.KluczZakresu = wartoscTekstu(z.ScopeId)

	switch rodzaj {
	case shared.AodMuteKindTimed:
		if z.EndsAt == nil || *z.EndsAt <= 0 {
			return dane.WyciszenieNakladki{}, bladZadaniaNakladki(
				"wyciszenie czasowe bez chwili końca — cisza bez terminu nie mówi Operatorowi, " +
					"do kiedy trwa")
		}
		wzor.Zakres = ""
		wzor.KlasaZdarzen = ""
		wzor.KluczZakresu = ""
		wzor.KonczySie = time.UnixMilli(*z.EndsAt).UTC().Format(formatZnacznikaBazy)
	case shared.AodMuteKindContextual:
		if wzor.Zakres != shared.AodMuteScopeModule && wzor.Zakres != shared.AodMuteScopeSession {
			return dane.WyciszenieNakladki{}, bladZadaniaNakladki(
				"wyciszenie kontekstowe bez zakresu — rozdz. 3.5 zna bieżący moduł albo bieżącą " +
					"kartę sesji")
		}
		if wzor.KluczZakresu == "" {
			return dane.WyciszenieNakladki{}, bladZadaniaNakladki(
				"wyciszenie kontekstowe bez wskazania bytu zakresu")
		}
		wzor.KlasaZdarzen = ""
	case shared.AodMuteKindEventClass:
		if wzor.KlasaZdarzen == "" {
			return dane.WyciszenieNakladki{}, bladZadaniaNakladki(
				"wyciszenie klasy zdarzeń bez wskazania klasy")
		}
		wzor.Zakres = shared.AodMuteScopeEventClass
		wzor.KluczZakresu = ""
	default:
		return dane.WyciszenieNakladki{}, bladZadaniaNakladki(
			"wyciszenie bez rodzaju; rodzajami są wyciszenie czasowe, kontekstowe i klasy zdarzeń")
	}
	return wzor, nil
}

// opisBytuWyciszenia nazywa byt wyciszenia w treści odmowy, dobierając opis
// do rodzaju: wyciszenie czasowe, zakres kontekstowy albo klasę zdarzeń.
func opisBytuWyciszenia(wzor dane.WyciszenieNakladki) string {
	switch wzor.Rodzaj {
	case shared.AodMuteKindTimed:
		return "wyciszenie czasowe"
	case shared.AodMuteKindContextual:
		return string(wzor.Zakres) + " " + wzor.KluczZakresu
	default:
		return "klasa zdarzeń " + string(wzor.KlasaZdarzen)
	}
}

// sprawdzMagazynWyciszen odmawia, gdy magazynu wyciszeń nie wpięto do
// adaptera, zamiast dopuścić dalsze odczyty i zapisy do składnika, którego nie ma.
func (a *adapterNakladkiAod) sprawdzMagazynWyciszen() error {
	if a.wyciszenia != nil {
		return nil
	}
	return bladBrakuSkladnikaNakladki("magazyn wyciszeń nakładki")
}

// wyciszeniaKontraktu przekłada wiersze na kształt kontraktu. Wykaz jest zawsze
// tablicą: brak wyciszeń to fakt, nie cisza.
func wyciszeniaKontraktu(wiersze []dane.WyciszenieNakladki) []shared.AodMute {
	wykaz := make([]shared.AodMute, 0, len(wiersze))
	for _, wiersz := range wiersze {
		wykaz = append(wykaz, wyciszenieKontraktu(wiersz))
	}
	return wykaz
}

// wyciszenieKontraktu przekłada jeden wiersz wyciszenia na kształt kontraktu,
// ustawiając pola zakresu, klasy zdarzeń i chwili końca tylko wtedy, gdy
// wiersz je niesie.
func wyciszenieKontraktu(w dane.WyciszenieNakladki) shared.AodMute {
	wyciszenie := shared.AodMute{
		Id:        w.Identyfikator,
		Kind:      w.Rodzaj,
		CreatedAt: chwilaBazy(w.Utworzono),
	}
	if w.Zakres != "" {
		zakres := w.Zakres
		wyciszenie.Scope = &zakres
	}
	wyciszenie.ScopeId = tekstOpcjonalny(w.KluczZakresu)
	wyciszenie.ScopeName = tekstOpcjonalny(w.NazwaZakresu)
	wyciszenie.DeviceId = tekstOpcjonalny(w.Urzadzenie)
	if w.KlasaZdarzen != "" {
		klasa := w.KlasaZdarzen
		wyciszenie.EventClass = &klasa
	}
	if w.KonczySie != "" {
		koniec := chwilaBazy(w.KonczySie)
		wyciszenie.EndsAt = &koniec
	}
	return wyciszenie
}

// sygnalKontraktu przekłada wiersz sygnału na kształt kontraktu, ustawiając
// liczbę wystąpień tylko wtedy, gdy wiersz sygnału ją niesie.
func sygnalKontraktu(s dane.SygnalNakladki) shared.AodSignal {
	sygnal := shared.AodSignal{
		Id:         s.Identyfikator,
		EventClass: s.KlasaZdarzen,
		Text:       s.Tresc,
		OccurredAt: chwilaBazy(s.ZdarzyloSie),
	}
	sygnal.ModuleId = tekstOpcjonalny(s.ModulKod)
	sygnal.SessionId = tekstOpcjonalny(s.SesjaKod)
	if s.LiczbaWystapien > 0 {
		liczba := s.LiczbaWystapien
		sygnal.OccurrenceCount = &liczba
	}
	return sygnal
}
