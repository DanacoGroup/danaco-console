package core

import (
	"context"
	"sort"

	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

type rejestrObecnosci struct {
	kontekst  context.Context
	nadzorca  *session.Nadzorca
	biegi     *rejestrBiegow
	czynnosc  *pamiecCzynnosci
	osadzenie zrodloOsadzenia
	nadania   zrodloNadan
	konta     zrodloKontaSesji
	strumien  zrodloStrumienia
	nadawca   *emiter
}

// Kontekst służy odczytom poza obsługą komendy — zdarzenie nie ma własnego żądania.
func nowyRejestrObecnosci(kontekst context.Context, nadzorca *session.Nadzorca,
	biegi *rejestrBiegow) *rejestrObecnosci {

	if kontekst == nil {
		kontekst = context.Background()
	}
	if biegi == nil {
		biegi = nowyRejestrBiegow()
	}
	rejestr := &rejestrObecnosci{
		kontekst: kontekst, nadzorca: nadzorca,
		biegi: biegi, czynnosc: nowaPamiecCzynnosci(),
	}
	biegi.Obserwuj(rejestr.biegZmieniony)
	return rejestr
}

func (r *rejestrObecnosci) ZeZrodlami(osadzenie zrodloOsadzenia, nadania zrodloNadan,
	konta zrodloKontaSesji) *rejestrObecnosci {

	r.osadzenie, r.nadania, r.konta = osadzenie, nadania, konta
	return r
}

// kontekstSesji oddaje kontekst życia rdzenia z kontem sesji dla odczytów bez żądania (decyzja 34).
func (r *rejestrObecnosci) kontekstSesji(idSesji string) context.Context {
	if r.konta == nil {
		return r.kontekst
	}
	return r.konta.KontekstSesji(r.kontekst, idSesji)
}

func (r *rejestrObecnosci) ZeStrumieniem(strumien zrodloStrumienia) *rejestrObecnosci {
	r.strumien = strumien
	return r
}

func (r *rejestrObecnosci) Odpisy(ctx context.Context) []shared.SessionPresence {
	if r == nil || r.nadzorca == nil {
		return nil
	}
	sesje := r.nadzorca.Rejestr().Sesje()
	wykaz := make([]shared.SessionPresence, 0, len(sesje))
	otwarte := map[string]struct{}{}
	for _, sesja := range sesje {
		for _, okno := range r.oknaOtwarte(sesja.Id) {
			otwarte[okno.Id] = struct{}{}
		}
		if !czynna(sesja.Stan) {
			continue
		}
		wykaz = append(wykaz, r.zloz(ctx, sesja))
	}
	sort.Slice(wykaz, func(i, j int) bool { return wykaz[i].SessionId < wykaz[j].SessionId })
	r.biegi.Zachowaj(otwarte)
	r.czynnosc.Zachowaj(otwarte)
	return wykaz
}

// Wiersz w bazie bez bytu żywego znaczy sesję do odtworzenia, nie sesję trwającą w tle.
func (r *rejestrObecnosci) Odpis(ctx context.Context, idSesji string) (shared.SessionPresence, bool) {
	if r == nil || r.nadzorca == nil || idSesji == "" {
		return shared.SessionPresence{}, false
	}
	sesja, err := r.nadzorca.Rejestr().Sesja(idSesji)
	if err != nil {
		return shared.SessionPresence{}, false
	}
	return r.zloz(ctx, sesja), true
}

func (r *rejestrObecnosci) zloz(ctx context.Context, sesja session.Sesja) shared.SessionPresence {
	odpis := shared.SessionPresence{
		SessionId:      sesja.Id,
		Live:           czynna(sesja.Stan),
		LastActivityAt: sesja.Zaktualizowano.UnixMilli(),
	}
	ostatnia, znana := r.czynnosc.Sesja(sesja.Id)
	if znana {
		odpis.LastActivityAt = ostatnia.Chwila.UnixMilli()
	}

	var ognisko session.Okno
	for _, okno := range r.oknaOtwarte(sesja.Id) {
		odpis.OpenWindowCount++
		if r.strumien != nil && r.strumien(okno.Id) {
			odpis.StreamingWindowCount++
		}
		if bieg, ma := r.biegi.Stan(okno.Id); ma && odpis.Loop == nil {
			wskazany := bieg
			odpis.Loop = &wskazany
		}
		odpis.AccessGrantIds = append(odpis.AccessGrantIds, r.nadaniaOkna(ctx, okno.Id)...)
		if ognisko.Id == "" || okno.Id == ostatnia.IdOkna {
			ognisko = okno
		}
	}
	naniesOgnisko(&odpis, ognisko)

	if kod := r.kodSrodowiska(ctx, sesja.Id); kod != "" {
		odpis.EnvironmentCode = &kod
	}
	return odpis
}

// Sesja bez otwartych okien zostaje bez wskazania — powrót prowadzi do samej karty.
func naniesOgnisko(odpis *shared.SessionPresence, ognisko session.Okno) {
	if ognisko.Id == "" {
		return
	}
	idOkna := ognisko.Id
	odpis.FocusedWindowId = &idOkna
	if ognisko.Modul != "" {
		modul := ognisko.Modul
		odpis.ModuleCode = &modul
	}
}

// Okno zamknięte nie należy do obrazu sesji trwającej w tle.
func (r *rejestrObecnosci) oknaOtwarte(idSesji string) []session.Okno {
	if r.nadzorca == nil {
		return nil
	}
	okna, err := r.nadzorca.Rejestr().OknaSesji(idSesji)
	if err != nil {
		return nil
	}
	otwarte := make([]session.Okno, 0, len(okna))
	for _, okno := range okna {
		if okno.CzyOtwarte() {
			otwarte = append(otwarte, okno)
		}
	}
	return otwarte
}

func (r *rejestrObecnosci) nadaniaOkna(ctx context.Context, idOkna string) []string {
	if r.nadania == nil {
		return nil
	}
	return r.nadania.NadaniaOkna(ctx, idOkna)
}

func (r *rejestrObecnosci) kodSrodowiska(ctx context.Context, idSesji string) string {
	if r.osadzenie == nil {
		return ""
	}
	return r.osadzenie.KodSrodowiska(ctx, idSesji)
}
