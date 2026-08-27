// Plik obsluguje dwie komendy wylaczen pamieci — memory.disable.list i
// memory.disable.set — oraz filtruje wpisy wylaczone przy odczycie pamieci
// w komendzie memory.list.
package core

import (
	"context"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// przedrostekWylaczeniaPamieci znakuje identyfikator wyłączenia nadany przez rdzeń
// przy zapisie nowego wpisu w magazynie.
const przedrostekWylaczeniaPamieci = "wyp-"

// ZWylaczeniami wpina magazyn wyłączeń pamięci. Bez niego komendy wyłączeń odmawiają
// z kodem internal_error, a memory.list nie ma czym odsiać wpisów wstrzymanych.
func (a *adapterPamieciPrzestrzeni) ZWylaczeniami(
	w dane.RepozytoriumWylaczenPamieci) *adapterPamieciPrzestrzeni {

	a.wylaczenia = w
	return a
}

// ── memory.disable.list ──────────────────────────────────────────────────────

// WylaczeniaPamieciZasiegu obsługuje `memory.disable.list`: oddaje wyłączenia
// obowiązujące w zasięgu żądania. Zawężenie składa się z pól niepustych; żądanie
// bez zawężenia pyta o wszystkie wyłączenia i tyle dostaje.
func (a *adapterPamieciPrzestrzeni) WylaczeniaPamieciZasiegu(ctx context.Context,
	z shared.MemoryDisableListRequest) (shared.MemoryDisableListResponse, error) {

	if err := a.sprawdzMagazynWylaczen(); err != nil {
		return shared.MemoryDisableListResponse{}, err
	}
	wzor := dane.WylaczeniePamieci{
		WpisIdentyfikator: wartoscTekstu(z.EntryId),
		KluczZasiegu:      wartoscTekstu(z.ScopeId),
	}
	if z.Level != nil {
		wzor.PoziomPamieci = *z.Level
	}
	if z.Scope != nil {
		wzor.Zasieg = *z.Scope
	}
	wiersze, err := a.wylaczenia.WylaczeniaPamieci(ctx, wzor)
	if err != nil {
		return shared.MemoryDisableListResponse{}, err
	}
	return shared.MemoryDisableListResponse{Disables: wylaczeniaKontraktu(wiersze)}, nil
}

// ── memory.disable.set ───────────────────────────────────────────────────────

// PrzestawWylaczeniePamieci obsługuje memory.disable.set: zakłada wyłączenie
// pamięci albo je znosi tą samą komendą, z polem disabled równym fałszowi
// dla zniesienia.
func (a *adapterPamieciPrzestrzeni) PrzestawWylaczeniePamieci(ctx context.Context,
	z shared.MemoryDisableSetRequest) (shared.MemoryDisableSetResponse, error) {

	if err := a.sprawdzMagazynWylaczen(); err != nil {
		return shared.MemoryDisableSetResponse{}, err
	}
	wzor, err := wzorWylaczenia(z)
	if err != nil {
		return shared.MemoryDisableSetResponse{}, err
	}
	if z.Disabled {
		return a.zalozWylaczenie(ctx, wzor)
	}
	return a.zniesWylaczenie(ctx, wzor, wartoscTekstu(z.DisableId))
}

// zalozWylaczenie zapisuje wyłączenie pamięci w magazynie i składa odpowiedź
// komendy niosącą wykaz wyłączeń po zmianie.
func (a *adapterPamieciPrzestrzeni) zalozWylaczenie(ctx context.Context,
	wzor dane.WylaczeniePamieci) (shared.MemoryDisableSetResponse, error) {

	if wzor.WpisIdentyfikator != "" {
		// Wyłączenie wpisu, którego nie ma, wskazywałoby byt nieistniejący.
		if _, err := a.wpisZadania(ctx, wzor.WpisIdentyfikator); err != nil {
			return shared.MemoryDisableSetResponse{}, err
		}
	}
	wzor.Identyfikator = nowyIdentyfikator(przedrostekWylaczeniaPamieci)
	zapisane, powstalo, err := a.wylaczenia.ZapiszWylaczeniePamieci(ctx, wzor)
	if err != nil {
		return shared.MemoryDisableSetResponse{}, bladZapisuWylaczenia(err)
	}
	wykaz, err := a.wylaczenia.WylaczeniaPamieci(ctx, dane.WylaczeniePamieci{})
	if err != nil {
		return shared.MemoryDisableSetResponse{}, err
	}
	return shared.MemoryDisableSetResponse{
		Disable:  wylaczenieKontraktu(zapisane),
		Disables: wylaczeniaKontraktu(wykaz),
		Changed:  powstalo,
	}, nil
}

// zniesWylaczenie znosi wyłączenie wskazane identyfikatorem albo bytem i zasięgiem,
// gdy identyfikator nie jest znany.
func (a *adapterPamieciPrzestrzeni) zniesWylaczenie(ctx context.Context,
	wzor dane.WylaczeniePamieci, identyfikator string) (shared.MemoryDisableSetResponse, error) {

	znoszone := dane.WylaczeniePamieci{}
	if identyfikator != "" {
		wykaz, err := a.wylaczenia.WylaczeniaPamieci(ctx, dane.WylaczeniePamieci{})
		if err != nil {
			return shared.MemoryDisableSetResponse{}, err
		}
		for _, wylaczenie := range wykaz {
			if wylaczenie.Identyfikator == identyfikator {
				znoszone = wylaczenie
				break
			}
		}
		if znoszone.Identyfikator == "" {
			return shared.MemoryDisableSetResponse{}, protocol.JakoError(protocol.NowyBlad(
				shared.ErrorCodeNotFound, "moduł Workspace: wyłączenia pamięci "+identyfikator+
					" nie ma; nie ma czego znieść"))
		}
	} else {
		zastane, jest, err := a.wylaczenia.WylaczeniePamieciPoBycie(ctx, wzor)
		if err != nil {
			return shared.MemoryDisableSetResponse{}, err
		}
		if !jest {
			return shared.MemoryDisableSetResponse{}, protocol.JakoError(protocol.NowyBlad(
				shared.ErrorCodeNotFound, "moduł Workspace: "+opisBytuWylaczenia(wzor)+
					" nie jest wyłączona w zasięgu "+string(wzor.Zasieg)+
					"; nie ma czego znieść"))
		}
		znoszone = zastane
	}
	zniesione, err := a.wylaczenia.ZniesWylaczeniePamieci(ctx, znoszone.Identyfikator)
	if err != nil {
		return shared.MemoryDisableSetResponse{}, err
	}
	wykaz, err := a.wylaczenia.WylaczeniaPamieci(ctx, dane.WylaczeniePamieci{})
	if err != nil {
		return shared.MemoryDisableSetResponse{}, err
	}
	return shared.MemoryDisableSetResponse{
		Disable:  wylaczenieKontraktu(znoszone),
		Disables: wylaczeniaKontraktu(wykaz),
		Changed:  zniesione,
	}, nil
}

// ── uwzględnienie wyłączeń w odczycie pamięci ───────────────────────────────

// sitoWylaczen składa sito wpisów pamięci z wyłączeń obowiązujących w rdzeniu.
// Zasięg wyłączenia jedzie w odpowiedzi, żeby Operator wiedział, który zasięg
// wpis wyciszył.
func (a *adapterPamieciPrzestrzeni) sitoWylaczen(ctx context.Context) (
	func(dane.WpisPamieciProjektu) (shared.MemoryDisabledEntry, bool), error) {

	if err := a.sprawdzMagazynWylaczen(); err != nil {
		return nil, err
	}
	wykaz, err := a.wylaczenia.WylaczeniaPamieci(ctx, dane.WylaczeniePamieci{})
	if err != nil {
		return nil, err
	}
	if len(wykaz) == 0 {
		return nil, nil
	}
	return func(wpis dane.WpisPamieciProjektu) (shared.MemoryDisabledEntry, bool) {
		for _, wylaczenie := range wykaz {
			if !wylaczenieObejmujeWpis(wylaczenie, wpis) {
				continue
			}
			wstrzymany := shared.MemoryDisabledEntry{
				EntryId:   wpis.Identyfikator,
				DisableId: wylaczenie.Identyfikator,
				Scope:     wylaczenie.Zasieg,
			}
			if wylaczenie.KluczZasiegu != "" {
				byt := wylaczenie.KluczZasiegu
				wstrzymany.ScopeId = &byt
			}
			if wylaczenie.PoziomPamieci != "" {
				poziom := wylaczenie.PoziomPamieci
				wstrzymany.Level = &poziom
			}
			return wstrzymany, true
		}
		return shared.MemoryDisabledEntry{}, false
	}, nil
}

// wylaczenieObejmujeWpis mówi, czy to wyłączenie wstrzymuje ten wpis: wskazaniem
// identyfikatora wpisu albo poziomem pamięci, na którym wpis stoi.
func wylaczenieObejmujeWpis(wylaczenie dane.WylaczeniePamieci,
	wpis dane.WpisPamieciProjektu) bool {

	if wylaczenie.WpisIdentyfikator != "" {
		return wylaczenie.WpisIdentyfikator == wpis.Identyfikator
	}
	poziom, ma := poziomPamieciZasiegu(wpis.Poziom)
	return ma && poziom == wylaczenie.PoziomPamieci
}

// poziomPamieciZasiegu przekłada zasięg wpisu na poziom pamięci. Zasięg spoza
// czterech poziomów pamięci nie jest poziomem pamięci i wyłączenie poziomu go
// nie dotyczy.
func poziomPamieciZasiegu(zasieg shared.ConfigScope) (shared.MemoryLevel, bool) {
	switch zasieg {
	case shared.ConfigScopeGlobal:
		return shared.MemoryLevelGlobal, true
	case shared.ConfigScopeEnvironment:
		return shared.MemoryLevelEnvironment, true
	case shared.ConfigScopeProject:
		return shared.MemoryLevelProject, true
	case shared.ConfigScopeSession:
		return shared.MemoryLevelSession, true
	default:
		return "", false
	}
}

// ── ustalenia wspólne ───────────────────────────────────────────────────────

// wzorWylaczenia składa wzór wyłączenia z żądania i odmawia, gdy żądanie nie
// wskazuje bytu do wyłączenia.
func wzorWylaczenia(z shared.MemoryDisableSetRequest) (dane.WylaczeniePamieci, error) {
	wzor := dane.WylaczeniePamieci{
		WpisIdentyfikator: wartoscTekstu(z.EntryId),
		Zasieg:            z.Scope,
		KluczZasiegu:      wartoscTekstu(z.ScopeId),
	}
	if z.Level != nil {
		wzor.PoziomPamieci = *z.Level
	}
	if wzor.WpisIdentyfikator == "" && wzor.PoziomPamieci == "" {
		return dane.WylaczeniePamieci{}, bladProjektu(
			"wyłączenie pamięci bez wskazania wpisu ani poziomu pamięci")
	}
	if wzor.WpisIdentyfikator != "" && wzor.PoziomPamieci != "" {
		return dane.WylaczeniePamieci{}, bladProjektu(
			"wyłączenie pamięci wskazuje i wpis " + wzor.WpisIdentyfikator +
				", i poziom " + string(wzor.PoziomPamieci) + " — jedno albo drugie")
	}
	if wzor.Zasieg == "" {
		return dane.WylaczeniePamieci{}, bladProjektu(
			"wyłączenie pamięci bez wskazania zasięgu; zasięg mówi, GDZIE pamięć milczy")
	}
	return wzor, nil
}

// opisBytuWylaczenia nazywa byt wyłączenia w treści odmowy: wpis identyfikatorem
// albo poziom pamięci nazwą poziomu.
func opisBytuWylaczenia(wzor dane.WylaczeniePamieci) string {
	if wzor.WpisIdentyfikator != "" {
		return "pamięć wpisu " + wzor.WpisIdentyfikator
	}
	return "pamięć poziomu " + string(wzor.PoziomPamieci)
}

// sprawdzMagazynWylaczen odmawia, gdy magazynu wyłączeń nie wpięto. Cisza byłaby
// tu gorsza od odmowy: Operator dostałby wykaz pamięci, o którym nie wiadomo,
// czy uwzględnia wyłączenia.
func (a *adapterPamieciPrzestrzeni) sprawdzMagazynWylaczen() error {
	if a.wylaczenia != nil {
		return nil
	}
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"moduł Workspace: magazyn wyłączeń pamięci niewpięty — nie ma gdzie zapisać ani czym "+
			"odsiać wyłączeń"))
}

// bladZapisuWylaczenia przekłada odmowę schematu na odmowę kontraktu. Naruszenie
// warunku bazy jest tu żądaniem niezgodnym z kontraktem, a nie awarią rdzenia.
func bladZapisuWylaczenia(err error) error {
	return protocol.JakoError(protocol.BladZeZrodla(shared.ErrorCodeValidationFailed, err))
}

// wylaczeniaKontraktu przekłada wiersze na kształt kontraktu. Wykaz jest zawsze
// tablicą: brak wyłączeń to fakt, nie cisza.
func wylaczeniaKontraktu(wiersze []dane.WylaczeniePamieci) []shared.MemoryDisable {
	wykaz := make([]shared.MemoryDisable, 0, len(wiersze))
	for _, wiersz := range wiersze {
		wykaz = append(wykaz, wylaczenieKontraktu(wiersz))
	}
	return wykaz
}

// wylaczenieKontraktu przekłada jeden wiersz wyłączenia pamięci z magazynu na
// kształt odpowiedzi kontraktu.
func wylaczenieKontraktu(w dane.WylaczeniePamieci) shared.MemoryDisable {
	wylaczenie := shared.MemoryDisable{
		Id:        w.Identyfikator,
		Scope:     w.Zasieg,
		CreatedAt: chwilaBazy(w.Utworzono),
		UpdatedAt: chwilaBazy(w.Zaktualizowano),
	}
	if w.WpisIdentyfikator != "" {
		wpis := w.WpisIdentyfikator
		wylaczenie.EntryId = &wpis
	}
	if w.PoziomPamieci != "" {
		poziom := w.PoziomPamieci
		wylaczenie.Level = &poziom
	}
	if w.KluczZasiegu != "" {
		byt := w.KluczZasiegu
		wylaczenie.ScopeId = &byt
	}
	return wylaczenie
}
