// Plik niesie port i uchwyty doraźnego dołożenia narzędzia na czas sesji: session.tool.attach, session.tool.detach, session.tool.list i tools.catalog.list, wraz z rozgłoszeniem session.tool.attached.
package core

import (
	"context"
	"sort"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// zrodloPlatformy jest przedrostkiem źródła pozycji pochodzących z kontraktu, a nie doniesionych z rozszerzenia.
const zrodloPlatformy = "danaco"

// NarzedziaSesji jest portem doraźnych dołożeń narzędzia do sesji oraz odczytu wykazu po ukośniku i katalogu.
type NarzedziaSesji interface {
	Doloz(ctx context.Context, z shared.SessionToolAttachRequest) (shared.SessionToolAttachResponse, error)
	// Zdejmij oddaje odpowiedź oraz pozycje faktycznie zdjęte, do rozgłoszenia session.tool.detached.
	Zdejmij(ctx context.Context, z shared.SessionToolDetachRequest) (shared.SessionToolDetachResponse, []shared.SessionTool, error)
	Wykaz(ctx context.Context, z shared.SessionToolListRequest) (shared.SessionToolListResponse, error)
	Katalog(ctx context.Context, z shared.ToolsCatalogListRequest) (shared.ToolsCatalogListResponse, error)
}

// Rozjazd portu NarzedziaSesji z adapterem zatrzymuje kompilację w tym miejscu, a nie dopiero na martwej komendzie rejestru.
var _ NarzedziaSesji = (*adapterNarzedziSesji)(nil)

// zarejestrujNarzedziaSesji wpina do rejestru komend cztery uchwyty rodziny doraźnego dołożenia narzędzia.
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
			// Rozgłoszenie idzie wyłącznie przy dołożeniu nowym, nie przy powtórzonej komendzie po ukośniku.
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
			// Jedno zdarzenie na jedną pozycję: nazwa skrócona może trafić w kilka dołożeń z różnych źródeł.
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

// rozglosDolozenieNarzedzia rozgłasza session.tool.attached, gdy model dostaje narzędzie, którego nie miał.
func rozglosDolozenieNarzedzia(ctx context.Context, e *emiter, idSesji string, narzedzie shared.SessionTool) {
	if e == nil {
		return
	}
	zdarzenie := shared.SessionToolAttachedEvent{SessionId: idSesji, Tool: narzedzie}
	zdarzenie.Actor, zdarzenie.ActorClientId = sprawca(ctx)
	e.wyslijDoKonta(ctx, shared.EventSessionToolAttached, idSesji, zdarzenie)
}

// rozglosZdjecieNarzedzia rozgłasza session.tool.detached, niosąc całą zdjętą pozycję, nie samą nazwę.
func rozglosZdjecieNarzedzia(ctx context.Context, e *emiter, idSesji string, narzedzie shared.SessionTool) {
	if e == nil {
		return
	}
	zdarzenie := shared.SessionToolDetachedEvent{SessionId: idSesji, Tool: narzedzie}
	zdarzenie.Actor, zdarzenie.ActorClientId = sprawca(ctx)
	e.wyslijDoKonta(ctx, shared.EventSessionToolDetached, idSesji, zdarzenie)
}

// narzedzieKontraktu przekłada wiersz dołożenia narzędzia z bazy na kształt oczekiwany przez kontrakt.
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

// nazwyTrafionych wylicza posortowane nazwy pełne pozycji katalogu, w które trafiła podana nazwa skrócona.
func nazwyTrafionych(trafione []*shared.ToolCatalogEntry) string {
	nazwy := make([]string, 0, len(trafione))
	for _, pozycja := range trafione {
		nazwy = append(nazwy, pozycja.Name)
	}
	sort.Strings(nazwy)
	return strings.Join(nazwy, ", ")
}

// ── odmowy ──
// Każdą odmowę identyfikuje przedrostek, powód i czynność naprawcza.

// bladNarzedziaSesji składa odmowę żądania rodziny session.tool, gdy treść żądania jest niezgodna z rodziną.
func bladNarzedziaSesji(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"dołożenia sesji: "+powod))
}

// bladNosnikaNarzedziSesji nazywa brak trwałości dołożeń albo usterkę zapisu, gdy dołożenie ma przeżyć rozłączenie klienta.
func bladNosnikaNarzedziSesji(przyczyna error) error {
	if przyczyna == nil {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"dołożenia sesji: trwałość dołożeń niewpięta, migracja 122 nie ma nośnika; "+
				"Operator uruchomi serwer z bazą danych — dołożenie ma przeżyć rozłączenie "+
				"klienta, a w pamięci gniazda by nie przeżyło"))
	}
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"dołożenia sesji: "+przyczyna.Error()+"; Operator powtórzy czynność, a przy "+
			"nawrocie sprawdzi dziennik serwera — żądanie było poprawne, zawiódł nośnik"))
}

// niedostepnaPozycja nazywa powód, dla którego wskazanej pozycji katalogu nie da się dołożyć do sesji.
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

// ── zawężenie i porządek wykazu ──
// Dopasowanie idzie środkiem nazwy, trafność przed alfabetem.

// tekstZawezenia normalizuje tekst szukany do porównania bez względu na
// wielkość liter; puste znaczy „bez zawężenia tekstem".
func tekstZawezenia(zapytanie *string) string {
	if zapytanie == nil {
		return ""
	}
	return strings.ToLower(strings.TrimSpace(*zapytanie))
}

// zawezPozycje zostawia w zwracanym wykazie wyłącznie pozycje spełniające zawężenie z żądania katalogu.
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

// trafnoscPozycji ocenia dopasowanie zapytania do pozycji w czterech stopniach; zero znaczy, że pozycja nie trafia.
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
