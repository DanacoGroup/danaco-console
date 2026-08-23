// Plik wpina pięć komend historii i archiwum eksperta — `agent.version.list`,
// `agent.version.restore`, `agent.archive`, `agent.restore`,
// `agent.archive.list` — wraz z portem WersjeEksperta i jego wypełnieniem.
//
// Historia i archiwum stoją w osobnym porcie, a nie w porcie Agenci: siedzą
// w innych tabelach, mają własne repozytoria i wchodzą do rdzenia jednym
// wywołaniem `zarejestrujWersjeEksperta`. Wtopienie ich w interfejs Agenci
// rozdęłoby port biblioteki o pięć czynności z biblioteką niezwiązanych.
//
// Nazwy komend i kształty pól pochodzą wyłącznie z pakietu `shared`; powielenie
// literału nazwy po którejkolwiek stronie jest błędem.
//
// Kształt `AgentVersion` jest ten sam co w `StudioVersion` i `LibraryVersion`:
// `id · agentId · label · suma kontrolna · createdAt`. Migawka tożsamości nie
// jedzie w wierszu wykazu — wykaz trzyma sumę kontrolną, a treść pobiera się
// osobno przy przywróceniu. `mode` jest polem niewymaganym i niesie trzeci stan
// promptu: brak wartości znaczy prompt globalny, `DOLACZ` prompt dopisywany,
// `ZASTAP` odstępstwo jawne.
//
// Kontrakt daje modułowi Agents wyłącznie zdarzenie `agent.changed`, więc
// przywrócenie wersji, archiwizacja i powrót z archiwum rozgłaszają się
// rodzajem `updated` wraz z ekspertem po zmianie.
package core

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// WersjeEksperta jest portem historii i archiwum eksperta.
type WersjeEksperta interface {
	WykazWersji(ctx context.Context, z shared.AgentVersionListRequest) (shared.AgentVersionListResponse, error)
	PrzywrocWersje(ctx context.Context, z shared.AgentVersionRestoreRequest) (shared.AgentVersionRestoreResponse, error)
	Zarchiwizuj(ctx context.Context, z shared.AgentArchiveRequest) (shared.AgentArchiveResponse, error)
	PrzywrocZArchiwum(ctx context.Context, z shared.AgentRestoreRequest) (shared.AgentRestoreResponse, error)
	WykazArchiwum(ctx context.Context, z shared.AgentArchiveListRequest) (shared.AgentArchiveListResponse, error)
}

// zarejestrujWersjeEksperta wpina pięć komend historii i archiwum eksperta.
// Woła ją składanie modułu Agents, obok `zarejestrujWarstwyAgenta`.
func zarejestrujWersjeEksperta(r *Rejestr, wersje WersjeEksperta, e *emiter) {
	if r == nil || wersje == nil {
		return
	}

	// Dwa odczyty bez zdarzenia — wykaz historii i wykaz archiwum niczego nie
	// zmieniają, więc nie ma czego rozgłaszać.
	r.Zarejestruj(shared.CommandAgentVersionList, obsluz(wersje.WykazWersji))
	r.Zarejestruj(shared.CommandAgentArchiveList, obsluz(wersje.WykazArchiwum))

	r.Zarejestruj(shared.CommandAgentVersionRestore,
		obsluz(func(ctx context.Context, z shared.AgentVersionRestoreRequest) (shared.AgentVersionRestoreResponse, error) {
			w, err := wersje.PrzywrocWersje(ctx, z)
			if err == nil {
				e.agent(shared.ChangeKindUpdated, w.Agent)
			}
			return w, err
		}))

	r.Zarejestruj(shared.CommandAgentArchive,
		obsluz(func(ctx context.Context, z shared.AgentArchiveRequest) (shared.AgentArchiveResponse, error) {
			w, err := wersje.Zarchiwizuj(ctx, z)
			// Zdarzenie idzie rodzajem `updated`, nie `deleted`: ekspert nadal
			// istnieje wraz z definicją i historią, tylko zszedł z wykazu
			// czynnych. `deleted` kazałoby oknom zapomnieć byt, który wróci.
			if err == nil && w.Archived && w.Agent != nil {
				e.agent(shared.ChangeKindUpdated, *w.Agent)
			}
			return w, err
		}))

	r.Zarejestruj(shared.CommandAgentRestore,
		obsluz(func(ctx context.Context, z shared.AgentRestoreRequest) (shared.AgentRestoreResponse, error) {
			w, err := wersje.PrzywrocZArchiwum(ctx, z)
			if err == nil && w.Restored && w.Agent != nil {
				e.agent(shared.ChangeKindUpdated, *w.Agent)
			}
			return w, err
		}))
}

// adapterWersjiEksperta wypełnia port WersjeEksperta repozytoriami danych.
// Stoi w tym samym pliku co port, bo jest jego jedynym wypełnieniem i liczy
// mniej niż sam port — osobny plik adaptera niósłby wyłącznie przekład.
type adapterWersjiEksperta struct {
	historia   dane.RepozytoriumWersjiAgenta
	archiwum   dane.RepozytoriumArchiwumAgentow
	biblioteka dane.RepozytoriumAgentow
}

var _ WersjeEksperta = (*adapterWersjiEksperta)(nil)

// NowyPortWersjiEksperta wiąże port z trzema repozytoriami: historią, archiwum
// i biblioteką ekspertów. Biblioteka służy wyłącznie oddaniu eksperta po
// zmianie — historia własnego widoku eksperta nie składa.
func NowyPortWersjiEksperta(historia dane.RepozytoriumWersjiAgenta,
	archiwum dane.RepozytoriumArchiwumAgentow,
	biblioteka dane.RepozytoriumAgentow) *adapterWersjiEksperta {

	return &adapterWersjiEksperta{historia: historia, archiwum: archiwum, biblioteka: biblioteka}
}

// WykazWersji oddaje historię eksperta od najnowszej wersji.
func (a *adapterWersjiEksperta) WykazWersji(ctx context.Context,
	z shared.AgentVersionListRequest) (shared.AgentVersionListResponse, error) {

	if a == nil || a.historia == nil {
		return shared.AgentVersionListResponse{}, bladBrakuKatalogu("historii ekspertów")
	}
	wiersze, err := a.historia.Wersje(ctx, z.AgentId)
	if err != nil {
		return shared.AgentVersionListResponse{}, bladWskazania(err, "ekspert", z.AgentId)
	}
	wersje := make([]shared.AgentVersion, 0, len(wiersze))
	for _, wiersz := range wiersze {
		wersje = append(wersje, wersjaEkspertaKontraktu(z.AgentId, wiersz))
	}
	return shared.AgentVersionListResponse{AgentId: z.AgentId, Versions: wersje, Total: len(wersje)}, nil
}

// PrzywrocWersje zapisuje wskazaną wersję jako wersję kolejną i oddaje eksperta
// po zmianie wraz z migawką nowo powstałą. Wersję wskazuje `versionId`, nie
// numer — tak wskazuje ją `library.version.restore`. Numer zostaje porządkiem
// historii, tożsamością wersji jest klucz wiersza, bo tylko on jest niezmienny.
func (a *adapterWersjiEksperta) PrzywrocWersje(ctx context.Context,
	z shared.AgentVersionRestoreRequest) (shared.AgentVersionRestoreResponse, error) {

	if a == nil || a.historia == nil {
		return shared.AgentVersionRestoreResponse{}, bladBrakuKatalogu("historii ekspertów")
	}
	numer, err := a.numerWersji(ctx, z.AgentId, z.VersionId)
	if err != nil {
		return shared.AgentVersionRestoreResponse{}, err
	}
	nowa, err := a.historia.Przywroc(ctx, z.AgentId, numer)
	if err != nil {
		return shared.AgentVersionRestoreResponse{}, bladWskazania(err, "wersja eksperta", z.AgentId)
	}
	ekspert, err := a.ekspert(ctx, z.AgentId)
	if err != nil {
		return shared.AgentVersionRestoreResponse{}, err
	}
	return shared.AgentVersionRestoreResponse{Agent: ekspert, Version: wersjaEkspertaKontraktu(z.AgentId, nowa)}, nil
}

// numerWersji odnajduje numer migawki po jej identyfikatorze. Odmowa nazywa
// przyczynę: wersja spoza historii tego eksperta nie jest błędem bazy, tylko
// wskazaniem nie do spełnienia.
func (a *adapterWersjiEksperta) numerWersji(ctx context.Context,
	kodAgenta, identyfikator string) (int, error) {

	if identyfikator == "" {
		return 0, bladZadaniaEksperta("wskazanie przywracanej wersji jest puste")
	}
	historia, err := a.historia.Wersje(ctx, kodAgenta)
	if err != nil {
		return 0, bladWskazania(err, "ekspert", kodAgenta)
	}
	for _, wiersz := range historia {
		if strconv.FormatInt(wiersz.Identyfikator, 10) == identyfikator {
			return wiersz.Numer, nil
		}
	}
	return 0, bladZadaniaEksperta("wersja " + identyfikator + " nie należy do historii tego eksperta")
}

// Zarchiwizuj odkłada eksperta poza wykaz czynnych. Ekspert nieodnaleziony albo
// już odłożony wraca z `archived=false`, a nie z odmową.
func (a *adapterWersjiEksperta) Zarchiwizuj(ctx context.Context,
	z shared.AgentArchiveRequest) (shared.AgentArchiveResponse, error) {

	if a == nil || a.archiwum == nil {
		return shared.AgentArchiveResponse{}, bladBrakuKatalogu("archiwum ekspertów")
	}
	odlozony, err := a.archiwum.Zarchiwizuj(ctx, z.AgentId)
	if err != nil {
		return shared.AgentArchiveResponse{}, err
	}
	wynik := shared.AgentArchiveResponse{AgentId: z.AgentId, Archived: odlozony}
	if odlozony {
		ekspert, err := a.ekspert(ctx, z.AgentId)
		if err != nil {
			return shared.AgentArchiveResponse{}, err
		}
		wynik.Agent = &ekspert
	}
	return wynik, nil
}

// PrzywrocZArchiwum oddaje eksperta wykazowi czynnych wraz z zapamiętanym
// stanem czynności.
func (a *adapterWersjiEksperta) PrzywrocZArchiwum(ctx context.Context,
	z shared.AgentRestoreRequest) (shared.AgentRestoreResponse, error) {

	if a == nil || a.archiwum == nil {
		return shared.AgentRestoreResponse{}, bladBrakuKatalogu("archiwum ekspertów")
	}
	wrocil, err := a.archiwum.PrzywrocZArchiwum(ctx, z.AgentId)
	if err != nil {
		return shared.AgentRestoreResponse{}, err
	}
	wynik := shared.AgentRestoreResponse{AgentId: z.AgentId, Restored: wrocil}
	if wrocil {
		ekspert, err := a.ekspert(ctx, z.AgentId)
		if err != nil {
			return shared.AgentRestoreResponse{}, err
		}
		wynik.Agent = &ekspert
	}
	return wynik, nil
}

// WykazArchiwum oddaje ekspertów odłożonych. Archiwum puste nie jest odmową —
// okno pokazuje wtedy stan pusty.
func (a *adapterWersjiEksperta) WykazArchiwum(ctx context.Context,
	_ shared.AgentArchiveListRequest) (shared.AgentArchiveListResponse, error) {

	if a == nil || a.archiwum == nil {
		return shared.AgentArchiveListResponse{Agents: []shared.Agent{}}, nil
	}
	wiersze, err := a.archiwum.Archiwum(ctx)
	if err != nil {
		return shared.AgentArchiveListResponse{}, err
	}
	eksperci := make([]shared.Agent, 0, len(wiersze))
	for _, wiersz := range wiersze {
		eksperci = append(eksperci, ekspertKontraktu(wiersz))
	}
	return shared.AgentArchiveListResponse{Agents: eksperci, Total: len(eksperci)}, nil
}

// ekspert dobiera eksperta z biblioteki na potrzeby wyniku i rozgłoszenia.
func (a *adapterWersjiEksperta) ekspert(ctx context.Context, kod string) (shared.Agent, error) {
	if a.biblioteka == nil {
		return shared.Agent{}, bladBrakuKatalogu("ekspertów")
	}
	wiersz, err := a.biblioteka.PoKodzie(ctx, kod)
	if err != nil {
		return shared.Agent{}, bladWskazania(err, "ekspert", kod)
	}
	return ekspertKontraktu(wiersz), nil
}

// wersjaEkspertaKontraktu przekłada migawkę bazy na pozycję historii kontraktu.
// Wiersz wykazu nie niesie migawki — oddaje tożsamość wersji, jej sumę
// kontrolną i czas; treść pobiera się osobno przy przywróceniu.
func wersjaEkspertaKontraktu(kodAgenta string, w dane.WersjaAgenta) shared.AgentVersion {
	wersja := shared.AgentVersion{
		Id:       strconv.FormatInt(w.Identyfikator, 10),
		AgentId:  kodAgenta,
		Label:    wskaznikTekstu("v" + strconv.Itoa(w.Numer)),
		Summary:  wskaznikTekstu(w.Powod),
		Author:   wskaznikTekstu(w.Autor),
		Checksum: wskaznikTekstu(sumaTozsamosci(w)),
		Mode:     trybWarstwyKontraktu(w.TrybNakladki),
	}
	if chwila, ok := chwilaZapisu(w.Zapisano); ok {
		wersja.CreatedAt = chwila
	}
	return wersja
}

// sumaTozsamosci liczy sumę kontrolną migawki. Liczy się ją w rdzeniu, a nie
// trzyma w kolumnie, bo migawki zakłada wyzwalacz bazy — kolumna wymagałaby
// liczenia sumy w SQL, gdzie nie ma po temu narzędzia. Suma obejmuje wyłącznie
// pola tożsamości: dwie wersje o równej sumie są nieodróżnialne w oknie.
func sumaTozsamosci(w dane.WersjaAgenta) string {
	skrot := sha256.Sum256([]byte(strings.Join([]string{
		w.Nazwa, w.Opis, w.InstrukcjeSystemowe,
		wartoscWskaznikaTekstu(w.KanalKod), wartoscWskaznikaTekstu(w.Model),
		wartoscWskaznikaTekstu(w.Transport), w.ParametryJSON,
		w.ImieWlasne, w.Favikon, w.UstawieniaJSON, w.TrybNakladki,
	}, "\x1f")))
	return hex.EncodeToString(skrot[:])
}

// wartoscWskaznikaTekstu znosi wskaźnik pusty do napisu pustego — suma kontrolna
// nie może zależeć od tego, czy kolumna jest NULL, czy pusta.
func wartoscWskaznikaTekstu(tekst *string) string {
	if tekst == nil {
		return ""
	}
	return *tekst
}

// trybWarstwyKontraktu przekłada tryb nakładki promptu na wyliczenie kontraktu.
// Trzeci stan niesie brak wartości, nie trzecia wartość wyliczenia: pusty tryb
// znaczy prompt globalny, więc pole wychodzi puste.
func trybWarstwyKontraktu(tryb string) *shared.IdentityMode {
	switch tryb {
	case string(shared.IdentityModeZASTAP):
		wartosc := shared.IdentityMode(shared.IdentityModeZASTAP)
		return &wartosc
	case string(shared.IdentityModeDOLACZ):
		wartosc := shared.IdentityMode(shared.IdentityModeDOLACZ)
		return &wartosc
	default:
		return nil
	}
}
