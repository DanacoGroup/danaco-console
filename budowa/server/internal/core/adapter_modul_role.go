// Odpowiedzialność pliku: rodzina `role.*` — dwie komendy nadające i zmieniające
// rolę okna w pętli koordynator–wykonawca.
//
// Rodzina jest fasadą, nie drugą prawdą o roli. Rola okna ma w rdzeniu jednego
// właściciela: pakiet `session` (rejestr okien i `rola_okna.go`, który zna
// słownik ról i normalizuje więź), a jej ślad trwały — kolumny
// `okno_komunikacji.rola_okna` i `okno_komunikacji.okno_koordynatora_id`
// z `migracja_002_okna.sql`. Rodzina `role.*` jedzie dokładnie na nich, tak samo
// jak `window.update`. Gdyby założyła własny zapis roli, produkt miałby dwie
// odpowiedzi na pytanie „jaką rolę ma to okno” — jedną z `window.list`, drugą
// z `role.*`.
//
// Wobec `window.update` rodzina dokłada trzy rzeczy, i tylko trzy:
//
//  1. Wcielenie (`persona`) — `window.update` nie zna tego pola wcale.
//  2. Odmowę zamiast cichego pominięcia. `window.update` ze wskazaniem
//     koordynatora dla okna, które wykonawcą nie jest, cicho zdejmuje wskazanie
//     (normalizujRole w `session/rola_okna.go`). Rodzina `role.*` odmawia
//     z powodem, bo jej jedynym tematem jest właśnie rola i jej więź.
//  3. Ślad trwały. `window.update` zmienia wyłącznie rejestr pamięciowy;
//     `role.*` zapisuje rolę także do wiersza okna, gdy wiersz istnieje — bez
//     tego `window.state.get` po restarcie rdzenia oddawałby rolę sprzed
//     nadania (adapter_stan_okna.go czyta wiersz, gdy rejestr okna nie zna).
//
// Wcielenie mieszka tam, gdzie już mieszka. Klient utrwala wcielenie okna
// komendą `config.set` na poziomie zasięgu `window` pod kluczem
// `multitasking.wcielenie` (`client/src/moduly/multitasking/wcielenia-analizy.ts`,
// `okno-results-analyzer.ts`). Rdzeń pisze i czyta ten sam adres, więc wcielenie
// nadane komendą `role.update` widzi selektor analityka i odwrotnie. Własny
// klucz albo własna tabela byłyby drugim wcieleniem tego samego okna.
package core

import (
	"context"
	"errors"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// kluczWcieleniaOkna to adres wcielenia roli w tabeli `ustawienie`, poziom
// zasięgu `window`, byt zasięgu — identyfikator okna. Wartość jest napisem, bo
// napisem jest w kontrakcie: katalogu wcieleń nikt nie rozstrzygnął, a rdzeń
// katalogów sobie nie wymyśla.
const kluczWcieleniaOkna = "multitasking.wcielenie"

// adapterRolOkien wypełnia port RoleOkien: adapter okien komunikacji
// rozszerzony o dwie komendy roli i o jej ślad trwały.
type adapterRolOkien struct {
	*adapterOkien
	// zestaw daje trzy drogi, których adapter okien nie miał: zapis roli do
	// wiersza (RoleOkien), zapis więzi koordynatora (Przekazania) i adres
	// wcielenia (Konfiguracja). Zerowy znaczy pracę na samej pamięci rejestru.
	zestaw *dane.Zestaw
	role   dane.RepozytoriumRolOkien
}

// ZRolami rozszerza adapter okien o rodzinę `role.*`.
func (a *adapterOkien) ZRolami(zestaw *dane.Zestaw) *adapterRolOkien {
	rozszerzony := &adapterRolOkien{adapterOkien: a, zestaw: zestaw}
	if zestaw != nil {
		rozszerzony.role = zestaw.RoleOkien()
	}
	return rozszerzony
}

// ── role.assign ──────────────────────────────────────────────────────────────

// NadajRole nadaje oknu rolę w pętli koordynator–wykonawca.
//
// Rola jest w tym żądaniu wymagana, więc żądanie bez niej albo z wartością
// spoza słownika kontraktu jest odmawiane. Sprawdza to `session.RolaZnana` —
// właściciel słownika ról; drugiego wykazu ról ten plik nie trzyma.
func (a *adapterRolOkien) NadajRole(ctx context.Context,
	z shared.RoleAssignRequest) (shared.RoleAssignResponse, error) {

	idOkna := strings.TrimSpace(z.WindowId)
	if idOkna == "" {
		return shared.RoleAssignResponse{}, bladRoli("żądanie nadania roli bez wskazania okna")
	}
	if !session.RolaZnana(z.Role) {
		return shared.RoleAssignResponse{}, bladRoli("rola " + nazwaRoli(z.Role) +
			" nie należy do słownika WindowRole")
	}
	stan, err := a.nadajRole(ctx, idOkna, &z.Role, z.CoordinatorWindowId)
	if err != nil {
		return shared.RoleAssignResponse{}, err
	}
	return shared.RoleAssignResponse{
		WindowId:            idOkna,
		Role:                stan.Rola,
		CoordinatorWindowId: wskaznikPolaRoli(stan.Koordynator),
	}, nil
}

// ── role.update ──────────────────────────────────────────────────────────────

// ZmienRole zmienia rolę okna albo jej wcielenie. Pole niewskazane zostaje bez
// zmiany — tak mówi kontrakt tej komendy o roli i tak samo traktujemy wcielenie.
//
// Żądanie bez ani jednego pola zmiany nie jest błędem: oddaje rolę i wcielenie
// obowiązujące. Odpowiedź niesie wtedy stan odczytany, a nie ciszę.
func (a *adapterRolOkien) ZmienRole(ctx context.Context,
	z shared.RoleUpdateRequest) (shared.RoleUpdateResponse, error) {

	idOkna := strings.TrimSpace(z.WindowId)
	if idOkna == "" {
		return shared.RoleUpdateResponse{}, bladRoli("żądanie zmiany roli bez wskazania okna")
	}
	if z.Role != nil && !session.RolaZnana(*z.Role) {
		return shared.RoleUpdateResponse{}, bladRoli("rola " + nazwaRoli(*z.Role) +
			" nie należy do słownika WindowRole")
	}
	stan, err := a.nadajRole(ctx, idOkna, z.Role, z.CoordinatorWindowId)
	if err != nil {
		return shared.RoleUpdateResponse{}, err
	}
	wcielenie, err := a.zapiszWcielenie(ctx, idOkna, z.Persona)
	if err != nil {
		return shared.RoleUpdateResponse{}, err
	}
	return shared.RoleUpdateResponse{
		WindowId: idOkna,
		Role:     stan.Rola,
		Persona:  wskaznikPolaRoli(wcielenie),
	}, nil
}

// ── rola i więź ──────────────────────────────────────────────────────────────

// nadajRole nanosi rolę i więź koordynatora na okno, gdziekolwiek okno żyje.
//
// Drogi są dwie, bo okno ma dwa życia. Okno otwarte w tym uruchomieniu rdzenia
// stoi w rejestrze pamięciowym i to on jest jego prawdą bieżącą; okno sprzed
// restartu zostało wyłącznie wierszem (adapter_stan_okna.go opisuje ten sam
// podział przy `window.state.get`). Rodzina `role.*` obsługuje oba, zamiast
// odmawiać oknu, o którym Operator wie z `window.state.get`, że istnieje.
//
// Wskazania niewypełnione (rola nil, koordynator nil) niczego nie zmieniają —
// wtedy obie drogi sprowadzają się do odczytu stanu obowiązującego.
func (a *adapterRolOkien) nadajRole(ctx context.Context, idOkna string,
	rola *shared.WindowRole, koordynator *string) (dane.RolaOkna, error) {

	if a.nadzorca != nil {
		if okno, err := a.nadzorca.Rejestr().Okno(idOkna); err == nil {
			return a.nadajRoleWRejestrze(ctx, okno, rola, koordynator)
		}
	}
	return a.nadajRoleWWierszu(ctx, idOkna, rola, koordynator)
}

// nadajRoleWRejestrze zmienia okno żywe. Zapis do wiersza idzie po zmianie
// w rejestrze i nie jest jej warunkiem — z jednym wyjątkiem opisanym
// przy `utrwalRole`.
func (a *adapterRolOkien) nadajRoleWRejestrze(ctx context.Context, okno session.Okno,
	rola *shared.WindowRole, koordynator *string) (dane.RolaOkna, error) {

	if err := sprawdzWiezRoli(rolaObowiazujaca(okno.RolaOkna, rola), koordynator); err != nil {
		return dane.RolaOkna{}, err
	}
	if rola == nil && koordynator == nil {
		return dane.RolaOkna{Rola: okno.RolaOkna, Koordynator: okno.OknoKoordynatora}, nil
	}
	zmienione, err := a.nadzorca.Rejestr().ZmienOkno(okno.Id,
		session.Zmiana{RolaOkna: rola, OknoKoordynatora: koordynator})
	if err != nil {
		return dane.RolaOkna{}, bladSesji(err)
	}
	stan := dane.RolaOkna{Rola: zmienione.RolaOkna, Koordynator: zmienione.OknoKoordynatora}
	if err := a.utrwalRole(ctx, okno.Id, stan); err != nil {
		return dane.RolaOkna{}, err
	}
	return stan, nil
}

// nadajRoleWWierszu zmienia okno, którego rejestr pamięciowy nie zna. Tu wiersz
// jest jedyną prawdą, więc nieudany zapis jest odmową komendy, a nie dodatkiem:
// bez wiersza nie zostałoby po nadaniu roli nic.
func (a *adapterRolOkien) nadajRoleWWierszu(ctx context.Context, idOkna string,
	rola *shared.WindowRole, koordynator *string) (dane.RolaOkna, error) {

	if a.role == nil {
		return dane.RolaOkna{}, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"rola okna: zapis roli niewpięty, okna "+idOkna+" nie ma gdzie zmienić"))
	}
	stan, err := a.role.RolaOkna(ctx, idOkna)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return dane.RolaOkna{}, bladBrakuOknaRoli(idOkna)
	}
	if err != nil {
		return dane.RolaOkna{}, err
	}
	docelowa := rolaObowiazujaca(stan.Rola, rola)
	if err := sprawdzWiezRoli(docelowa, koordynator); err != nil {
		return dane.RolaOkna{}, err
	}
	if rola == nil && koordynator == nil {
		return stan, nil
	}
	if err := a.role.ZapiszRoleOkna(ctx, idOkna, docelowa); err != nil {
		return dane.RolaOkna{}, bladZapisuRoli(idOkna, err)
	}
	stan.Rola = docelowa
	// Zapis roli spoza pętli zdjął więź w tej samej instrukcji UPDATE, więc
	// odpowiedź musi to powiedzieć, a nie oddawać koordynatora sprzed zmiany.
	if docelowa != shared.WindowRoleExecutor {
		stan.Koordynator = ""
		return stan, nil
	}
	wskazany := strings.TrimSpace(wartoscTekstu(koordynator))
	if wskazany == "" {
		return stan, nil
	}
	// Więź z oknem, którego nie ma, byłaby potwierdzeniem relacji, która nie
	// powstała. Rejestr pamięciowy sprawdza to sam (ErrBrakKoordynatora); na
	// drodze wiersza sprawdzamy to tutaj, bo tam nikt inny tego nie zrobi.
	if _, err := a.role.RolaOkna(ctx, wskazany); err != nil {
		if errors.Is(err, dane.ErrBrakWiersza) {
			return dane.RolaOkna{}, protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
				"rola okna: okno koordynatora "+wskazany+" nie istnieje"))
		}
		return dane.RolaOkna{}, err
	}
	if err := a.polaczZKoordynatorem(ctx, idOkna, wskazany); err != nil {
		return dane.RolaOkna{}, err
	}
	stan.Koordynator = wskazany
	return stan, nil
}

// utrwalRole odkłada rolę okna żywego do wiersza.
//
// Brak wiersza nie jest awarią. Okno komunikacji dostaje wiersz leniwie, przy
// pierwszej wiadomości, więc okno świeżo otwarte nie ma go jeszcze wcale —
// rola żyje wtedy w rejestrze i wykaz okien czyta ją stamtąd (dolozWiezi
// w `adapter_okna.go` opisuje ten sam podział dla więzi). Każda inna usterka
// zapisu kończy komendę: zapis, który się nie udał, a został przemilczany,
// oddałby Operatorowi rolę, której po restarcie nie będzie.
func (a *adapterRolOkien) utrwalRole(ctx context.Context, idOkna string, stan dane.RolaOkna) error {
	if a.role == nil {
		return nil
	}
	err := a.role.ZapiszRoleOkna(ctx, idOkna, stan.Rola)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return nil
	}
	if err != nil {
		return bladZapisuRoli(idOkna, err)
	}
	if stan.Koordynator == "" {
		return nil
	}
	return a.polaczZKoordynatorem(ctx, idOkna, stan.Koordynator)
}

// polaczZKoordynatorem zapisuje więź koordynator–wykonawca w kolumnie, którą
// wypełnia także `window.handoff` — drugiej drogi do tej więzi nie ma.
//
// Okna bez wiersza pomija z tego samego powodu, co `utrwalRole`: więź żyje wtedy
// w rejestrze pamięciowym, a wykaz okien zostawia wartość z pamięci dokładnie
// wtedy, gdy wiersza nie ma.
func (a *adapterRolOkien) polaczZKoordynatorem(ctx context.Context, idOkna, koordynator string) error {
	if a.zestaw == nil || a.zestaw.Przekazania == nil {
		return nil
	}
	err := a.zestaw.Przekazania.UstawKoordynatora(ctx, idOkna, koordynator)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return nil
	}
	if err != nil {
		return bladZapisuRoli(idOkna, err)
	}
	return nil
}

// sprawdzWiezRoli odmawia więzi roli, która więzi nie niesie.
//
// Koordynatora ma wyłącznie okno wykonawcy (`session/rola_okna.go`).
// Pakiet sesji takie wskazanie po cichu zdejmuje, bo jego zadaniem jest otworzyć
// okno mimo wszystko; komenda, której jedynym tematem jest rola, musi powiedzieć
// wprost, że o wskazany skutek nie prosiła sama siebie.
func sprawdzWiezRoli(rola shared.WindowRole, koordynator *string) error {
	if strings.TrimSpace(wartoscTekstu(koordynator)) == "" {
		return nil
	}
	if rola == shared.WindowRoleExecutor {
		return nil
	}
	return bladRoli("koordynatora niesie wyłącznie okno wykonawcy; rola " +
		nazwaRoli(rola) + " go nie ma")
}

// rolaObowiazujaca zwraca rolę po zmianie: wskazaną, a bez wskazania — bieżącą.
func rolaObowiazujaca(biezaca shared.WindowRole, wskazana *shared.WindowRole) shared.WindowRole {
	if wskazana != nil {
		return *wskazana
	}
	return biezaca
}

// wskaznikPolaRoli oddaje pole opcjonalne odpowiedzi: napis pusty wychodzi jako
// brak pola, nie jako pole o pustej treści.
//
// Pomocnik jest własny, a nie wspólny, bo pakiet niesie dwa o nazwie
// `wskaznikTekstu` i o różnym znaczeniu pustego napisu (`sesja_konfiguracja.go`
// oddaje brak, `adapter_modul_auth.go` — wskaźnik na pustkę). Rodzina `role.*`
// mówi wprost, której zasady trzyma się jej odpowiedź: koordynator pusty
// i wcielenie puste są brakiem, a klient odróżnia brak od wartości.
func wskaznikPolaRoli(wartosc string) *string {
	if wartosc == "" {
		return nil
	}
	return &wartosc
}

// nazwaRoli oddaje rolę w postaci czytelnej w komunikacie odmowy. Rola pusta ma
// nazwę własną — „(pusta)” mówi więcej niż pustka w środku zdania.
func nazwaRoli(rola shared.WindowRole) string {
	if strings.TrimSpace(string(rola)) == "" {
		return "(pusta)"
	}
	return string(rola)
}

// ── wcielenie roli ───────────────────────────────────────────────────────────

// zapiszWcielenie zapisuje wcielenie roli i oddaje wcielenie obowiązujące.
//
// Żądanie bez pola `persona` niczego nie zapisuje — oddaje wcielenie zastane.
// Pole wypełnione napisem pustym zdejmuje wcielenie: to jedyny sposób, jaki
// kontrakt daje na powrót do roli bez wcielenia, a zapis pustego napisu
// zostawiałby wiersz ustawienia udający wcielenie o pustej nazwie.
func (a *adapterRolOkien) zapiszWcielenie(ctx context.Context, idOkna string,
	wcielenie *string) (string, error) {

	if a.zestaw == nil || a.zestaw.Konfiguracja == nil {
		if wcielenie == nil {
			return "", nil
		}
		return "", protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
			"rola okna: zapis konfiguracji niewpięty, wcielenia okna "+idOkna+" nie ma gdzie zapisać"))
	}
	if wcielenie == nil {
		return a.wcielenieObowiazujace(ctx, idOkna)
	}
	wartosc := strings.TrimSpace(*wcielenie)
	if wartosc == "" {
		if err := a.zestaw.Konfiguracja.Usun(ctx, shared.ConfigScopeWindow, idOkna, kluczWcieleniaOkna); err != nil {
			return "", err
		}
		return "", nil
	}
	ustawienie := dane.Ustawienie{
		Poziom:       shared.ConfigScopeWindow,
		KluczZasiegu: idOkna,
		Klucz:        kluczWcieleniaOkna,
		Wartosc:      &wartosc,
	}
	if err := a.zestaw.Konfiguracja.Ustaw(ctx, ustawienie); err != nil {
		return "", err
	}
	return wartosc, nil
}

// wcielenieObowiazujace czyta wcielenie zapisane przy oknie. Brak wpisu znaczy
// rolę bez wcielenia i jest stanem poprawnym.
func (a *adapterRolOkien) wcielenieObowiazujace(ctx context.Context, idOkna string) (string, error) {
	wpis, jest, err := a.zestaw.Konfiguracja.Odczytaj(ctx, shared.ConfigScopeWindow, idOkna, kluczWcieleniaOkna)
	if err != nil {
		return "", err
	}
	if !jest {
		return "", nil
	}
	return tekstLubPusty(wpis.Wartosc), nil
}

// ── okno po zmianie ──────────────────────────────────────────────────────────

// OknoRoli oddaje okno po zmianie roli — treść zdarzenia `window.changed`.
// Drugi wynik mówi, czy okno w ogóle udało się złożyć; nieudane złożenie kończy
// wyłącznie rozgłoszenie, nigdy komendę, która już się powiodła.
func (a *adapterRolOkien) OknoRoli(ctx context.Context, idOkna string) (shared.Window, bool) {
	if a.nadzorca != nil {
		if okno, err := a.nadzorca.Rejestr().Okno(idOkna); err == nil {
			return oknoKontraktu(okno), true
		}
	}
	if a.zestaw == nil {
		return shared.Window{}, false
	}
	okno, jest, err := oknoUtrwalone(ctx, a.zestaw, idOkna)
	if err != nil {
		return shared.Window{}, false
	}
	return okno, jest
}

// ── odmowy ───────────────────────────────────────────────────────────────────

// bladRoli składa odmowę żądania niezgodnego z kontraktem rodziny.
func bladRoli(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		"rola okna: "+powod))
}

// bladBrakuOknaRoli nazywa okno, którego nie ma — ani w rejestrze, ani w bazie.
func bladBrakuOknaRoli(idOkna string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
		"rola okna: okno komunikacji "+idOkna+" nie istnieje"))
}

// bladZapisuRoli niesie usterkę zapisu wraz z oknem, którego dotyczy. Odmowa
// wskazująca okno pozwala Operatorowi powtórzyć czynność; „nie udało się”
// bez nazwy bytu nie pozwala nic.
func bladZapisuRoli(idOkna string, err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return bladBrakuOknaRoli(idOkna)
	}
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		"rola okna: nie można zapisać roli okna "+idOkna+": "+err.Error()))
}
