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

type adapterSesji struct {
	nadzorca     *session.Nadzorca
	trwalosc     *utrwalaczStanow
	obecnosc     *rejestrObecnosci
	zapewnienie  ZapewnienieSesji
	projekty     dane.RepozytoriumPrzestrzeniRoboczej
	zestaw       *dane.Zestaw
	przerwijTure PrzerwanieTury
	rozmowa      zrodloOstatniejWypowiedzi
}

type zrodloOstatniejWypowiedzi interface {
	Wykaz(idOkna string, przed *string, ograniczenie *int) ([]shared.Message, bool)
}

func (a *adapterSesji) ZRozmowa(d zrodloOstatniejWypowiedzi) *adapterSesji {
	a.rozmowa = d
	return a
}

// Sesja czeka na Operatora, gdy ostatnia wypowiedź któregoś okna należy do modelu i jest domknięta.
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

func nowyAdapterSesji(nadzorca *session.Nadzorca) *adapterSesji {
	return &adapterSesji{nadzorca: nadzorca}
}

func (a *adapterSesji) ZTrwaloscia(u *utrwalaczStanow) *adapterSesji {
	a.trwalosc = u
	return a
}

func (a *adapterSesji) ZObecnoscia(o *rejestrObecnosci) *adapterSesji {
	a.obecnosc = o
	return a
}

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
	sesja := a.nadzorca.ZalozSesjeSrodowiska(tytul, projekt, srodowisko, kontoRejestru(ctx, a.zestaw))
	// Sesja idzie do bazy od razu, nie dopiero z pierwszą wiadomością, bo zapis ma trwać.
	a.utrwalZalozona(ctx, sesja)
	return shared.SessionCreateResponse{Session: sesjaKontraktu(sesja)}, nil
}

// Menu okna roboczego niesie tę czynność pod nazwą „Eksportuj sesję”.
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

func (a *adapterSesji) Wykaz(ctx context.Context, z shared.SessionListRequest) (shared.SessionListResponse, error) {
	wszystkie := a.nadzorca.Rejestr().SesjeKonta(kontoRejestru(ctx, a.zestaw))
	wybrane := make([]shared.Session, 0, len(wszystkie))
	for _, sesja := range wszystkie {
		if z.Status != nil && sesja.Stan != *z.Status {
			continue
		}
		// Wykaz bez wskazania stanu jest historią bieżącą, a sesja archiwalna
		// z niej znika; jej wgląd daje wskazanie stanu i session.archive.list.
		if z.Status == nil && sesja.Stan == shared.SessionStatusArchived {
			continue
		}
		opis := sesjaKontraktu(sesja)
		czeka := a.czekaNaReakcje(sesja)
		opis.AwaitingReaction = &czeka
		wybrane = append(wybrane, opis)
	}
	wynik := shared.SessionListResponse{Sessions: wycinek(wybrane, z.Offset, z.Limit), Total: len(wybrane)}
	if z.IncludePresence != nil && *z.IncludePresence {
		wynik.Presence = a.obecnosc.Odpisy(ctx, kontoRejestru(ctx, a.zestaw))
	}
	return wynik, nil
}

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

func (a *adapterSesji) Zamknij(ctx context.Context, z shared.SessionCloseRequest) (shared.SessionCloseResponse, error) {
	zamkniete, err := a.nadzorca.ZamknijSesje(z.SessionId)
	if err != nil {
		return shared.SessionCloseResponse{}, bladSesji(err)
	}
	sesja, err := a.nadzorca.Rejestr().Sesja(z.SessionId)
	if err != nil {
		return shared.SessionCloseResponse{}, bladSesji(err)
	}
	a.trwalosc.StanOkien(ctx, identyfikatoryOkienSesji(zamkniete), shared.WindowStatusClosed)
	a.trwalosc.StanSesji(ctx, sesja.Id, sesja.Stan)
	return shared.SessionCloseResponse{Session: sesjaKontraktu(sesja)}, nil
}

// Usuwanie sesji mieszka w osobnym pliku jako operacja nieodwracalna i zbiorcza.

func identyfikatoryOkienSesji(okna []session.Okno) []string {
	identyfikatory := make([]string, 0, len(okna))
	for _, okno := range okna {
		identyfikatory = append(identyfikatory, okno.Id)
	}
	return identyfikatory
}

// Wartości spoza zakresu dają wykaz pusty, nigdy błąd.
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

/*
Rejestr sesji trzyma konto liczbą, a `WarunekKonta` w bazie tłumaczy zero

	i NULL na konto najstarsze. Żeby wykaz z rejestru wypadał tak samo jak
	zapytanie, konto operatora rozstrzyga się tu raz i tą samą regułą.
*/
func kontoRejestru(ctx context.Context, zestaw *dane.Zestaw) int64 {
	if zestaw == nil {
		return kontoZRepozytorium(ctx, nil)
	}
	return kontoZRepozytorium(ctx, zestaw.KontoWlasciciela)
}

// kontoZRepozytorium rozstrzyga to samo tam, gdzie adapter nie ma całego zestawu.
func kontoZRepozytorium(ctx context.Context, konta dane.RepozytoriumKontaWlasciciela) int64 {
	if id := dane.KontoOperatora(ctx); id != 0 {
		return id
	}
	if konta == nil {
		return 0
	}
	konto, err := konta.Konto(ctx)
	if err != nil {
		return 0
	}
	return konto.Id
}
