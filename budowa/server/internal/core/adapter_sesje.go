package core

import (
	"context"
	"encoding/json"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"

	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// adapterSesji wypełnia port Sesje nadzorcą pakietu sesji. Pakiet session jest
// właścicielem pojęcia sesji i okna, więc rdzeń niczego tu
// nie rozstrzyga — przekłada żądanie kontraktu na wywołanie nadzorcy i wynik
// z powrotem na kontrakt.
type adapterSesji struct {
	nadzorca    *session.Nadzorca
	trwalosc    *utrwalaczStanow
	obecnosc    *rejestrObecnosci
	zapewnienie ZapewnienieSesji
	// projekty daje przynależność sesji do projektu (czynności historii).
	projekty dane.RepozytoriumPrzestrzeniRoboczej
	// zestaw daje komplet repozytoriow — potrzebny wylacznie kopiowaniu sesji.
	zestaw *dane.Zestaw
	// przerwijTure przerywa ture okna — potrzebne zatrzymaniu calej sesji.
	przerwijTure PrzerwanieTury
	// rozmowa daje ostatnią wypowiedź okna; po niej wykaz sesji rozstrzyga,
	// która sesja czeka na reakcję Operatora.
	rozmowa zrodloOstatniejWypowiedzi
}

// zrodloOstatniejWypowiedzi oddaje ostatnią wypowiedź okna komunikacji.
type zrodloOstatniejWypowiedzi interface {
	Wykaz(idOkna string, przed *string, ograniczenie *int) ([]shared.Message, bool)
}

// ZRozmowa wpina dziennik rozmowy; bez niego wykaz sesji nie rozstrzyga
// oczekiwania na reakcję i pole zostaje puste.
func (a *adapterSesji) ZRozmowa(d zrodloOstatniejWypowiedzi) *adapterSesji {
	a.rozmowa = d
	return a
}

// czekaNaReakcje mówi, czy sesja czeka na Operatora: ostatnia wypowiedź
// któregokolwiek z jej okien należy do modelu i jest domknięta.
func (a *adapterSesji) czekaNaReakcje(sesja session.Sesja) bool {
	if a.rozmowa == nil {
		return false
	}
	for _, idOkna := range sesja.IdOkien {
		wykaz, _ := a.rozmowa.Wykaz(idOkna, nil, nil)
		if len(wykaz) == 0 {
			continue
		}
		ostatnia := wykaz[len(wykaz)-1]
		if ostatnia.Role == shared.MessageRoleAssistant && ostatnia.Status == shared.MessageStatusComplete {
			return true
		}
	}
	return false
}

// nowyAdapterSesji wiąże port z nadzorcą pakietu sesji, jedynym źródłem
// prawdy o żywym stanie sesji rdzenia.
func nowyAdapterSesji(nadzorca *session.Nadzorca) *adapterSesji {
	return &adapterSesji{nadzorca: nadzorca}
}

// ZTrwaloscia dokłada zapis stanu sesji i jej okien do bazy. Bez niego sesja
// zamknięta wraca po restarcie rdzenia jako czynna.
func (a *adapterSesji) ZTrwaloscia(u *utrwalaczStanow) *adapterSesji {
	a.trwalosc = u
	return a
}

// ZObecnoscia dokłada żywy stan sesji trwających na rdzeniu. Bez niego wykaz
// sesji odpowiada bez obecności, a kontrolka powrotu bierze ją ze zdarzenia
// session.changed.
func (a *adapterSesji) ZObecnoscia(o *rejestrObecnosci) *adapterSesji {
	a.obecnosc = o
	return a
}

// Utworz zakłada sesję wspólną dla plików, pamięci, projektu i agentów
// w ramach jednego środowiska platformy.
func (a *adapterSesji) Utworz(ctx context.Context, z shared.SessionCreateRequest) (shared.SessionCreateResponse, error) {
	tytul, projekt := "", ""
	if z.Title != nil {
		tytul = *z.Title
	}
	if z.ProjectId != nil {
		projekt = *z.ProjectId
	}
	srodowisko := ""
	if z.EnvironmentCode != nil {
		srodowisko = *z.EnvironmentCode
	}
	sesja := a.nadzorca.ZalozSesjeSrodowiska(tytul, projekt, srodowisko)
	// Sesja idzie do bazy od razu, nie dopiero z pierwszą wiadomością, bo zapis ma trwać.
	a.utrwalZalozona(ctx, sesja)
	return shared.SessionCreateResponse{Session: sesjaKontraktu(sesja)}, nil
}

// Wydaj oddaje zapis sesji wraz z jej oknami. Menu okna roboczego niesie tę
// czynność pod nazwą „Eksportuj sesję"; zapis maszynowy niesie sesję i okna,
// zapis do czytania — te same rzeczy zdaniami.
func (a *adapterSesji) Wydaj(_ context.Context, z shared.SessionExportRequest) (shared.SessionExportResponse, error) {
	sesja, err := a.nadzorca.Rejestr().Sesja(z.SessionId)
	if err != nil {
		return shared.SessionExportResponse{}, bladSesji(err)
	}
	okna, err := a.nadzorca.Rejestr().OknaSesji(z.SessionId)
	if err != nil {
		return shared.SessionExportResponse{}, bladSesji(err)
	}
	postac := shared.SessionExportFormat(shared.SessionExportFormatJson)
	if z.Format != nil && *z.Format != "" {
		postac = *z.Format
	}
	opis := sesjaKontraktu(sesja)
	if postac == shared.SessionExportFormatMarkdown {
		return shared.SessionExportResponse{
			Content:  zapisSesjiDoCzytania(opis, oknaKontraktu(okna)),
			Format:   postac,
			FileName: z.SessionId + ".md",
		}, nil
	}
	tresc, err := json.MarshalIndent(struct {
		Session shared.Session  `json:"session"`
		Windows []shared.Window `json:"windows"`
	}{Session: opis, Windows: oknaKontraktu(okna)}, "", "  ")
	if err != nil {
		return shared.SessionExportResponse{}, protocol.JakoError(protocol.NowyBlad(
			shared.ErrorCodeInternalError, "sesje: zapisu sesji nie da się złożyć: "+err.Error()))
	}
	return shared.SessionExportResponse{
		Content:  string(tresc),
		Format:   postac,
		FileName: z.SessionId + ".json",
	}, nil
}

// zapisSesjiDoCzytania składa zapis sesji zdaniami: nazwa, stan i wykaz okien.
func zapisSesjiDoCzytania(sesja shared.Session, okna []shared.Window) string {
	nazwa := sesja.Id
	if sesja.Title != nil && *sesja.Title != "" {
		nazwa = *sesja.Title
	}
	wiersze := []string{"# " + nazwa, "", "Stan: " + string(sesja.Status), ""}
	if sesja.EnvironmentCode != nil && *sesja.EnvironmentCode != "" {
		wiersze = append(wiersze, "Środowisko: "+*sesja.EnvironmentCode, "")
	}
	wiersze = append(wiersze, "## Okna")
	for _, okno := range okna {
		wiersze = append(wiersze, "- "+okno.Id+" — moduł "+okno.ModuleId+", stan "+string(okno.Status))
	}
	return strings.Join(wiersze, "\n") + "\n"
}

// Wykaz zwraca sesje konta w kolejności ich świeżości, opcjonalnie zawężone
// stanem i wycinkiem wyniku.
func (a *adapterSesji) Wykaz(ctx context.Context, z shared.SessionListRequest) (shared.SessionListResponse, error) {
	wszystkie := a.nadzorca.Rejestr().Sesje()
	wybrane := make([]shared.Session, 0, len(wszystkie))
	for _, sesja := range wszystkie {
		if z.Status != nil && sesja.Stan != *z.Status {
			continue
		}
		opis := sesjaKontraktu(sesja)
		czeka := a.czekaNaReakcje(sesja)
		opis.AwaitingReaction = &czeka
		wybrane = append(wybrane, opis)
	}
	wynik := shared.SessionListResponse{Sessions: wycinek(wybrane, z.Offset, z.Limit), Total: len(wybrane)}
	if z.IncludePresence != nil && *z.IncludePresence {
		wynik.Presence = a.obecnosc.Odpisy(ctx)
	}
	return wynik, nil
}

// Otworz zwraca sesję wraz z jej oknami komunikacji, w relacji jednej sesji
// do wielu okien tej rozmowy.
func (a *adapterSesji) Otworz(_ context.Context, z shared.SessionOpenRequest) (shared.SessionOpenResponse, error) {
	sesja, err := a.nadzorca.Rejestr().Sesja(z.SessionId)
	if err != nil {
		return shared.SessionOpenResponse{}, bladSesji(err)
	}
	okna, err := a.nadzorca.Rejestr().OknaSesji(z.SessionId)
	if err != nil {
		return shared.SessionOpenResponse{}, bladSesji(err)
	}
	return shared.SessionOpenResponse{Session: sesjaKontraktu(sesja), Windows: oknaKontraktu(okna)}, nil
}

// Zamknij zamyka sesję wraz z jej oknami i ubija wszystkie ich procesy, nie
// zostawiając żadnego osieroconego.
func (a *adapterSesji) Zamknij(_ context.Context, z shared.SessionCloseRequest) (shared.SessionCloseResponse, error) {
	zamkniete, err := a.nadzorca.ZamknijSesje(z.SessionId)
	if err != nil {
		return shared.SessionCloseResponse{}, bladSesji(err)
	}
	sesja, err := a.nadzorca.Rejestr().Sesja(z.SessionId)
	if err != nil {
		return shared.SessionCloseResponse{}, bladSesji(err)
	}
	a.trwalosc.StanOkien(identyfikatoryOkienSesji(zamkniete), shared.WindowStatusClosed)
	a.trwalosc.StanSesji(sesja.Id, sesja.Stan)
	return shared.SessionCloseResponse{Session: sesjaKontraktu(sesja)}, nil
}

// Usuwanie sesji mieszka w osobnym pliku jako operacja nieodwracalna i zbiorcza.

// identyfikatoryOkienSesji wylicza identyfikatory okien zamkniętych wraz
// z sesją, potrzebne do ubicia ich procesów.
func identyfikatoryOkienSesji(okna []session.Okno) []string {
	identyfikatory := make([]string, 0, len(okna))
	for _, okno := range okna {
		identyfikatory = append(identyfikatory, okno.Id)
	}
	return identyfikatory
}

// wycinek stosuje przesunięcie i ograniczenie wykazu. Wartości spoza zakresu
// dają wykaz pusty, nigdy błąd.
func wycinek(wykaz []shared.Session, przesuniecie, ograniczenie *int) []shared.Session {
	poczatek := 0
	if przesuniecie != nil && *przesuniecie > 0 {
		poczatek = *przesuniecie
	}
	if poczatek >= len(wykaz) {
		return []shared.Session{}
	}
	koniec := len(wykaz)
	if ograniczenie != nil && *ograniczenie > 0 && poczatek+*ograniczenie < koniec {
		koniec = poczatek + *ograniczenie
	}
	return wykaz[poczatek:koniec]
}
