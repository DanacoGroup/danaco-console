// Odpowiedzialność pliku: port i uchwyty doraźnego dołożenia narzędzia na czas
// sesji — `session.tool.attach`, `session.tool.detach`, `session.tool.list`
// i `tools.catalog.list` — wraz z rozgłoszeniem `session.tool.attached`
// (migracja 122).
//
// Ekspert dostaje dobrany podzbiór narzędzi; komenda po ukośniku wstrzykuje jedno
// narzędzie na żądanie, nie ruszając definicji eksperta ani biegu rozmowy. Jest to
// jedyne miejsce, w którym zestaw narzędzi rośnie w trakcie pracy, więc
// `session.tool.attach` i `session.tool.detach` nie są narzędziami modelu:
// narzędzie, którym model dokłada sobie narzędzia, byłoby drugą prawdą o tym, czym
// model dysponuje. Odczyt model ma — `session.tool.list` i `tools.catalog.list`.
package core

import (
	"context"
	"sort"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// zrodloPlatformy jest przedrostkiem źródła pozycji pochodzących z kontraktu.
// Operator widzi po nim, że pozycja jest czynnością samej platformy, a nie
// rzeczą doniesioną — dokładnie tak, jak `anthropic-skills:` mówi o skillu.
const zrodloPlatformy = "danaco"

// NarzedziaSesji jest portem doraźnych dołożeń i wykazu po ukośniku.
type NarzedziaSesji interface {
	Doloz(ctx context.Context, z shared.SessionToolAttachRequest) (shared.SessionToolAttachResponse, error)
	// Zdejmij oddaje ODPOWIEDŹ oraz pozycje faktycznie zdjęte. Druga wartość
	// istnieje po to, by uchwyt miał czym rozgłosić `session.tool.detached`:
	// odpowiedź kontraktu niesie zestaw PO czynności, a z niego nie widać, co
	// zniknęło ani pod jaką nazwą pełną.
	Zdejmij(ctx context.Context, z shared.SessionToolDetachRequest) (shared.SessionToolDetachResponse, []shared.SessionTool, error)
	Wykaz(ctx context.Context, z shared.SessionToolListRequest) (shared.SessionToolListResponse, error)
	Katalog(ctx context.Context, z shared.ToolsCatalogListRequest) (shared.ToolsCatalogListResponse, error)
}

// Rozjazd portu z adapterem zatrzymuje kompilację tutaj, nie na martwej komendzie.
var _ NarzedziaSesji = (*adapterNarzedziSesji)(nil)

// zarejestrujNarzedziaSesji wpina cztery uchwyty rodziny.
func zarejestrujNarzedziaSesji(r *Rejestr, ns NarzedziaSesji, e *emiter) {
	if r == nil || ns == nil {
		return
	}

	r.Zarejestruj(shared.CommandSessionToolAttach,
		obsluz(func(ctx context.Context, z shared.SessionToolAttachRequest) (shared.SessionToolAttachResponse, error) {
			odpowiedz, err := ns.Doloz(ctx, z)
			if err != nil {
				return shared.SessionToolAttachResponse{}, err
			}
			// Rozgłoszenie idzie wyłącznie przy dołożeniu nowym. Powtórzona
			// komenda po ukośniku niczego modelowi nie dodała, więc zdarzenie
			// „model dostał narzędzie, którego nie miał" byłoby wtedy nieprawdą.
			if !odpowiedz.AlreadyAttached {
				rozglosDolozenieNarzedzia(ctx, e, z.SessionId, odpowiedz.Tool)
			}
			return odpowiedz, nil
		}))

	r.Zarejestruj(shared.CommandSessionToolDetach,
		obsluz(func(ctx context.Context, z shared.SessionToolDetachRequest) (shared.SessionToolDetachResponse, error) {
			odpowiedz, zdjete, err := ns.Zdejmij(ctx, z)
			if err != nil {
				return shared.SessionToolDetachResponse{}, err
			}
			// Jedno zdarzenie na jedną pozycję, tak samo jak przy dołożeniu.
			// Nazwa skrócona może trafić w kilka dołożeń z różnych źródeł, a
			// sąsiednie okno usuwa wiersze po nazwie pełnej — zbiorcze zdarzenie
			// kazałoby mu zgadywać, które z nich zeszły.
			for _, narzedzie := range zdjete {
				rozglosZdjecieNarzedzia(ctx, e, z.SessionId, narzedzie)
			}
			return odpowiedz, nil
		}))

	r.Zarejestruj(shared.CommandSessionToolList,
		obsluz(func(ctx context.Context, z shared.SessionToolListRequest) (shared.SessionToolListResponse, error) {
			return ns.Wykaz(ctx, z)
		}))

	r.Zarejestruj(shared.CommandToolsCatalogList,
		obsluz(func(ctx context.Context, z shared.ToolsCatalogListRequest) (shared.ToolsCatalogListResponse, error) {
			return ns.Katalog(ctx, z)
		}))
}

// rozglosDolozenieNarzedzia rozgłasza `session.tool.attached`.
//
// Skoro model dostaje narzędzie, którego nie miał, Operator ma to zobaczyć na
// ekranie. Zdarzenie niesie całą pozycję — nazwę pełną ze źródłem, skróconą
// i opis — bo sama nazwa `:design-audit` nie mówi Operatorowi nic. Sprawcę
// bierzemy z kontekstu (jak zdarzenia z `zdarzenia.go`): bez niego dołożenie
// zrobione ręką asystenta wyglądałoby jak własne.
func rozglosDolozenieNarzedzia(ctx context.Context, e *emiter, idSesji string, narzedzie shared.SessionTool) {
	if e == nil {
		return
	}
	zdarzenie := shared.SessionToolAttachedEvent{SessionId: idSesji, Tool: narzedzie}
	zdarzenie.Actor, zdarzenie.ActorClientId = sprawca(ctx)
	e.wyslij(shared.EventSessionToolAttached, idSesji, zdarzenie)
}

// rozglosZdjecieNarzedzia rozgłasza `session.tool.detached`.
//
// Rozgłaszanie samego dołożenia zostawiłoby sąsiednie okno z wierszem narzędzia,
// którego model już nie ma — zestaw pokazany szerszym, niż jest, kłamie tak samo
// jak poszerzony po cichu. Zdarzenie niesie całą pozycję, nie samą nazwę: okno
// usuwa wiersz po
// nazwie pełnej, a gdy dołożenia jeszcze nie widziało, ma z czego złożyć wpis
// dziennika mówiący Operatorowi, co dokładnie zeszło.
func rozglosZdjecieNarzedzia(ctx context.Context, e *emiter, idSesji string, narzedzie shared.SessionTool) {
	if e == nil {
		return
	}
	zdarzenie := shared.SessionToolDetachedEvent{SessionId: idSesji, Tool: narzedzie}
	zdarzenie.Actor, zdarzenie.ActorClientId = sprawca(ctx)
	e.wyslij(shared.EventSessionToolDetached, idSesji, zdarzenie)
}

// narzedzieKontraktu przekłada wiersz dołożenia na kształt kontraktu.
func narzedzieKontraktu(n dane.NarzedzieSesji) shared.SessionTool {
	return shared.SessionTool{
		Name:        n.NazwaPelna,
		ShortName:   n.NazwaSkrocona,
		Description: n.Opis,
		Kind:        shared.SlashEntryKind(n.Rodzaj),
		Group:       n.Grupa,
		Source:      shared.SessionToolSource(n.Zrodlo),
		AttachedAt:  n.Dolozono,
	}
}

// nazwyTrafionych wylicza nazwy pełne pozycji, w które trafiła nazwa skrócona.
func nazwyTrafionych(trafione []*shared.ToolCatalogEntry) string {
	nazwy := make([]string, 0, len(trafione))
	for _, pozycja := range trafione {
		nazwy = append(nazwy, pozycja.Name)
	}
	sort.Strings(nazwy)
	return strings.Join(nazwy, ", ")
}

// ── odmowy ───────────────────────────────────────────────────────────────────
// Każda odmowa mówi trzy rzeczy: co odmówiło (przedrostek), dlaczego (powód
// nazwany) i czym Operator to zmieni (czynność po średniku).

// bladNarzedziaSesji składa odmowę żądania niezgodnego z rodziną.
func bladNarzedziaSesji(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"dołożenia sesji: "+powod))
}

// bladNosnikaNarzedziSesji nazywa brak trwałości dołożeń albo usterkę zapisu.
//
// Dołożenie, które nie doszło do bazy, znika przy pierwszym rozłączeniu klienta,
// a ma przetrwać. Zestaw pusty oddany zamiast odmowy wyglądałby jak
// „nic nie dołożono", choć komenda zameldowała powodzenie.
func bladNosnikaNarzedziSesji(przyczyna error) error {
	if przyczyna == nil {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"dołożenia sesji: trwałość dołożeń niewpięta, migracja 122 nie ma nośnika; "+
				"Operator uruchomi rdzeń z bazą danych — dołożenie ma przeżyć rozłączenie "+
				"klienta, a w pamięci gniazda by nie przeżyło"))
	}
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"dołożenia sesji: "+przyczyna.Error()+"; Operator powtórzy czynność, a przy "+
			"nawrocie sprawdzi dziennik rdzenia — żądanie było poprawne, zawiódł nośnik"))
}

// niedostepnaPozycja nazywa powód, dla którego pozycji nie da się dołożyć.
func niedostepnaPozycja(p shared.ToolCatalogEntry) string {
	if p.Kind == shared.SlashEntryKindAction {
		return "pozycja " + p.Name + " jest KOMENDĄ AKCJI, nie narzędziem: wykonuje " +
			"czynność aplikacji i zestawu narzędzi nie zmienia (rozdz. 6); Operator " +
			"wywoła ją wprost komendą " + p.ShortName + ", a dołożyć może pozycje " +
			"rodzaju tool — te z katalogu rozszerzeń"
	}
	return "pozycja " + p.Name + " nie jest zainstalowana albo jest wyłączona, więc " +
		"nie ma czego podać modelowi; Operator zainstaluje ją i włączy komendami " +
		"extension.install oraz extension.toggle — dołożenie pozycji nieobecnej " +
		"meldowałoby poszerzenie zestawu, którego nie ma"
}

// ── zawężenie i porządek wykazu ──────────────────────────────────────────────
// Cztery zachowania: dopasowanie w środku nazwy, szukanie obejmujące obie nazwy,
// porządek według trafności przed alfabetem, filtrowanie od pierwszego znaku.
// Wytłuszczenie trafionego fragmentu należy do widoku — rdzeń oddaje pozycje
// i ich porządek.

// tekstZawezenia normalizuje tekst szukany do porównania bez względu na
// wielkość liter; puste znaczy „bez zawężenia tekstem".
func tekstZawezenia(zapytanie *string) string {
	if zapytanie == nil {
		return ""
	}
	return strings.ToLower(strings.TrimSpace(*zapytanie))
}

// zawezPozycje zostawia pozycje spełniające zawężenie żądania.
func zawezPozycje(pozycje []*shared.ToolCatalogEntry,
	z shared.ToolsCatalogListRequest) []*shared.ToolCatalogEntry {

	szukane := tekstZawezenia(z.Query)
	grupa := ""
	if z.Group != nil {
		grupa = strings.TrimSpace(*z.Group)
	}
	wybrane := make([]*shared.ToolCatalogEntry, 0, len(pozycje))
	for _, pozycja := range pozycje {
		if z.Kind != nil && *z.Kind != "" && pozycja.Kind != *z.Kind {
			continue
		}
		if grupa != "" && !strings.EqualFold(pozycja.Group, grupa) {
			continue
		}
		if szukane != "" && trafnoscPozycji(pozycja, szukane) == 0 {
			continue
		}
		wybrane = append(wybrane, pozycja)
	}
	return wybrane
}

// trafnoscPozycji ocenia dopasowanie; zero znaczy „pozycja nie trafia".
//
// Dopasowanie idzie środkiem, nie od początku: przy przedrostkach źródła każda
// pozycja zaczyna się tak samo, więc szukanie od początku nazwy pełnej byłoby
// bezużyteczne. Stopnie są cztery, bo `ski` ma trafić w `skill-creator` wyżej
// niż w `template-skill`, a w oba wyżej niż w pozycję, która ma `ski` wyłącznie
// w opisie.
func trafnoscPozycji(pozycja *shared.ToolCatalogEntry, szukane string) int {
	skrocona := strings.ToLower(pozycja.ShortName)
	switch {
	case strings.HasPrefix(skrocona, szukane):
		return 4
	case strings.Contains(skrocona, szukane):
		return 3
	case strings.Contains(strings.ToLower(pozycja.Name), szukane):
		return 2
	case strings.Contains(strings.ToLower(pozycja.Description), szukane):
		return 1
	}
	return 0
}

// uporzadkujPozycje układa wykaz: trafność przed alfabetem, a przy braku
// zawężenia tekstem — grupa, potem nazwa pełna. Sortowanie jest stabilne, więc
// remis trafności zachowuje porządek źródła zamiast rozstrzygać się losowo.
func uporzadkujPozycje(pozycje []*shared.ToolCatalogEntry, szukane string) {
	sort.SliceStable(pozycje, func(i, j int) bool {
		if szukane != "" {
			pierwsza, druga := trafnoscPozycji(pozycje[i], szukane), trafnoscPozycji(pozycje[j], szukane)
			if pierwsza != druga {
				return pierwsza > druga
			}
		}
		if pozycje[i].Group != pozycje[j].Group {
			return pozycje[i].Group < pozycje[j].Group
		}
		return pozycje[i].Name < pozycje[j].Name
	})
}

// przytnijPozycje wycina okno wykazu i oddaje pozycje po wartości — wskaźniki
// służyły wyłącznie do oznaczenia dołożeń i dalej nie wychodzą.
func przytnijPozycje(pozycje []*shared.ToolCatalogEntry, limit, przesuniecie *int) []shared.ToolCatalogEntry {
	od := 0
	if przesuniecie != nil && *przesuniecie > 0 {
		od = min(*przesuniecie, len(pozycje))
	}
	doWyl := len(pozycje)
	if limit != nil && *limit > 0 {
		doWyl = min(od+*limit, doWyl)
	}
	okno := make([]shared.ToolCatalogEntry, 0, doWyl-od)
	for _, pozycja := range pozycje[od:doWyl] {
		okno = append(okno, *pozycja)
	}
	return okno
}

// grupyPozycji wylicza grupy obecne w wykazie PO zawężeniu — nośnik
// pogrupowania po przeznaczeniu, które obowiązuje tak samo jak w menu wyboru.
func grupyPozycji(pozycje []*shared.ToolCatalogEntry) []string {
	widziane := make(map[string]bool, 16)
	grupy := make([]string, 0, 16)
	for _, pozycja := range pozycje {
		if pozycja.Group == "" || widziane[pozycja.Group] {
			continue
		}
		widziane[pozycja.Group] = true
		grupy = append(grupy, pozycja.Group)
	}
	sort.Strings(grupy)
	return grupy
}
