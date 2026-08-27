package core

import (
	"context"
	"sort"

	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// rejestrObecnosci składa żywy stan sesji trwających na rdzeniu dla kontrolki powrotu na stronie głównej. Rejestr niczego nie posiada — składa odpis kontraktu z czterech źródeł: nadzorcy, pętli, warstwy danych i telemetrii.
type rejestrObecnosci struct {
	kontekst  context.Context
	nadzorca  *session.Nadzorca
	biegi     *rejestrBiegow
	czynnosc  *pamiecCzynnosci
	osadzenie zrodloOsadzenia
	nadania   zrodloNadan
	strumien  zrodloStrumienia
	nadawca   *emiter
}

// nowyRejestrObecnosci wiąże składacz z nadzorcą i rejestrem biegów. Kontekst
// służy odczytom wykonywanym poza obsługą komendy — zdarzenie nie ma własnego
// żądania, więc nie ma też jego kontekstu.
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

// ZeZrodlami dokłada osadzenie sesji w środowisku i nadania dostępu okien, potrzebne do złożenia pełnego odpisu.
func (r *rejestrObecnosci) ZeZrodlami(osadzenie zrodloOsadzenia, nadania zrodloNadan) *rejestrObecnosci {
	r.osadzenie, r.nadania = osadzenie, nadania
	return r
}

// ZeStrumieniem dokłada wiedzę o turze trwającej w oknie. Bez niej odpis mówi
// o zerze okien strumieniujących, a nie o braku sesji.
func (r *rejestrObecnosci) ZeStrumieniem(strumien zrodloStrumienia) *rejestrObecnosci {
	r.strumien = strumien
	return r
}

// Odpisy zwracają żywy stan wszystkich sesji czynnych, uporządkowany po identyfikatorze. Przy okazji jednego pełnego przejścia rejestr wykreśla ślady okien już zamkniętych, których nie ma w rejestrze nadzorcy.
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

// Odpis zwraca żywy stan jednej sesji. Sesja nieznana rejestrowi nadzorcy daje
// fałsz: wiersz w bazie bez bytu żywego znaczy sesję do odtworzenia, a nie
// sesję trwającą w tle.
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

// zloz składa odpis jednej sesji ze wszystkich czterech źródeł: nadzorcy, pętli, warstwy danych i telemetrii.
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

// naniesOgnisko wskazuje okno, do którego prowadzi powrót, i moduł, w którym to
// okno pracuje. Sesja bez otwartych okien zostaje bez wskazania — powrót
// prowadzi wtedy do samej karty.
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

// oknaOtwarte zwraca okna sesji przyjmujące pracę. Okno zamknięte nie należy do
// obrazu sesji trwającej w tle.
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

// nadaniaOkna odczytuje nadania dostępu okna przez port. Brak podłączonego portu daje wykaz pusty, nie odmowę.
func (r *rejestrObecnosci) nadaniaOkna(ctx context.Context, idOkna string) []string {
	if r.nadania == nil {
		return nil
	}
	return r.nadania.NadaniaOkna(ctx, idOkna)
}

// kodSrodowiska odczytuje osadzenie sesji przez port, zwracając pusty napis, gdy port nie jest podłączony.
func (r *rejestrObecnosci) kodSrodowiska(ctx context.Context, idSesji string) string {
	if r.osadzenie == nil {
		return ""
	}
	return r.osadzenie.KodSrodowiska(ctx, idSesji)
}
