// Odpowiedzialność pliku: dwie komendy rodziny `config.*`, które nie dotykają
// rejestru ustawień — `config.explain.get` (prowenancja wywołania modelu)
// i `config.window.open` (zakres, na którym otwiera się okno konfiguracji).
//
// Stoją osobno od `handlers_config.go` i `handlers_sesja_konfiguracja.go`, bo
// tamte pliki wpinają komendy czytające i piszące rejestr ustawień. Te dwie nie
// zapisują niczego: pierwsza składa wywołanie modelu tak, jak zrobiłaby to tura,
// i pokazuje je Operatorowi; druga rozstrzyga zakres wejściowy okna. Wspólny
// mają wyłącznie przedrostek nazwy.
//
// Prowenancja idzie tą samą drogą co tura, nie drugą: `config.explain.get` nie
// ma własnego składacza wywołania, tylko bierze dokładnie te funkcje, którymi
// jedzie `message.send`:
//
//	zapytanieKanalu → nakladkaOkna → uzupelnijSrodowisko → uzupelnijKonfiguracje
//	→ ustawieniaKanaluGlownego → nakladkaKanaluGlownego → injection.Argumenty
//
// Gdyby prowenancja składała wywołanie po swojemu, pokazywałaby wiersz, którego
// tura nigdy nie wykona. Ten plik niczego nie dopuszcza i niczego nie wstrzymuje.
//
// Od wiersza procesu prowenancja różni się świadomie w dwóch miejscach:
//
//  1. Treść `--settings` i `--mcp-config` jedzie tu jako treść, a nie jako
//     ścieżka pliku tymczasowego. Tura materializuje te napisy na dysku tuż
//     przed uruchomieniem procesu (`injection/materializacja.go`) i sprząta je
//     po sobie, więc pliku o tej ścieżce jeszcze nie ma. Tak samo rozstrzyga sam
//     pakiet injection: jego prowenancja niesie treść pierwotną, a argv realną
//     ścieżkę (`injection/przebieg.go`).
//  2. `argv[0]` jest programem. Kontrakt tej komendy nie ma osobnego pola
//     `program`, więc bez tego nie byłoby widać, jaki plik wykonywalny rusza —
//     a `argv` ma być wierszem wywołania, nie samą listą przełączników.
package core

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"danacoconsole/server/internal/injection"
	"danacoconsole/server/internal/konfig"
	"danacoconsole/server/internal/models"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/server/internal/session"
	"danacoconsole/shared"
)

// adapterProwenancjiKonfiguracji wypełnia port ProwenancjaKonfiguracji: adapter
// ustawień świadomy osi rozszerzony o dwie komendy, które potrzebują drogi tury,
// a nie rejestru ustawień.
//
// To rozszerzenie portu, a nie drugi port: typ osadza adapter ustawień, więc
// niesie komplet jego metod i jest tym samym bytem, którym jadą pozostałe
// komendy rodziny — konfiguracja ma w rdzeniu jedną bramę.
type adapterProwenancjiKonfiguracji struct {
	*adapterUstawienOsi
	// rozmowa jest jedynym składaczem wywołania modelu w rdzeniu. Zerowy znaczy
	// odmowę z kodem `internal_error`, nie ciszę — bez niego nie ma czego
	// wyjaśniać, a odpowiedź z pustym argv wyglądałaby na wywołanie bez programu.
	rozmowa *adapterRozmowy
}

// ZProwenancja rozszerza adapter ustawień o `config.explain.get`
// i `config.window.open`. Bierze adapter rozmowy — ten sam egzemplarz, którym
// prowadzone są tury, więc prowenancja pokazuje wywołanie tego okna, a nie
// wywołanie złożone z wartości domyślnych.
func (a *adapterUstawienOsi) ZProwenancja(rozmowa *adapterRozmowy) *adapterProwenancjiKonfiguracji {
	return &adapterProwenancjiKonfiguracji{adapterUstawienOsi: a, rozmowa: rozmowa}
}

// ── config.explain.get ───────────────────────────────────────────────────────

// ProwenancjaWywolania oddaje warunki, w jakich pójdzie najbliższa tura okna:
// wiersz wywołania, prompt systemowy złożony z warstw, ustawienia przekazane
// procesowi i sumę kontrolną konstytucji.
//
// Rzecz liczona jest na świeżo, a nie odczytywana z zapisu poprzedniej tury:
// zapis mówiłby, co poszło ostatnio, a Operator pyta okno konfiguracji o to, co
// pójdzie po jego zmianach. Wywołanie niczego nie utrwala i niczego nie
// uruchamia.
func (a *adapterProwenancjiKonfiguracji) ProwenancjaWywolania(ctx context.Context,
	z shared.ConfigExplainGetRequest) (shared.ConfigExplainGetResponse, error) {

	if a == nil || a.rozmowa == nil || a.rozmowa.nadzorca == nil || a.rozmowa.kanaly == nil {
		return shared.ConfigExplainGetResponse{}, bladProwenancji(shared.ErrorCodeInternalError,
			"droga tury niewpięta; prowenancji wywołania nie ma z czego złożyć")
	}
	okno, err := a.oknoProwenancji(z.WindowId, z.SessionId)
	if err != nil {
		return shared.ConfigExplainGetResponse{}, err
	}
	kanal := strings.TrimSpace(okno.KanalModelu)
	if kanal == "" {
		return shared.ConfigExplainGetResponse{}, bladProwenancji(shared.ErrorCodeConflict,
			"okno "+okno.Id+" nie wskazuje kanału modelu; tura nie ma dokąd pójść")
	}
	// Rejestr oddaje wyłącznie kanały czynne. Kanał wskazany oknem, a nieobecny
	// w rejestrze, wywróciłby także turę (models/wysylka.go) — odmawiamy tym
	// samym rozpoznaniem, zamiast pokazywać wiersz wywołania, którego nikt nie
	// wykona.
	definicja, jest := a.rozmowa.kanaly.Definicja(kanal)
	if !jest {
		return shared.ConfigExplainGetResponse{}, bladProwenancji(shared.ErrorCodeNotFound,
			"kanał modelu "+kanal+" nie istnieje w rejestrze albo jest nieczynny")
	}

	zapytanie := a.zapytanieProwenancji(ctx, okno)
	nakladka := nakladkaKanaluGlownego(zapytanie)
	return shared.ConfigExplainGetResponse{
		Argv:         argvProwenancji(definicja, ustawieniaKanaluGlownego(definicja, zapytanie), nakladka),
		SystemPrompt: nakladka.Prompt(),
		Settings:     ustawieniaProwenancji(zapytanie.PlikUstawien),
		// Suma kontrolna samej konstytucji, nie całej nakładki: profil roli
		// i ekspertyza zadaniowa zmieniają się z oknem, a konstytucja jest
		// warstwą, której niezmienność się sprawdza. Pusta konstytucja daje
		// pusty skrót, bo nie ma czego sumować.
		ConstitutionHash: models.SkrotTresci(zapytanie.Nakladka.Konstytucja),
	}, nil
}

// zapytanieProwenancji składa zapytanie kanału dokładnie tak, jak robi to
// `prowadzTure` — w tej samej kolejności i tymi samymi funkcjami. Nie ma tu
// treści wypowiedzi ani historii: żadna z nich nie wchodzi do wiersza wywołania
// ani do promptu systemowego, a zmyślona treść pytania byłaby atrapą.
func (a *adapterProwenancjiKonfiguracji) zapytanieProwenancji(ctx context.Context,
	okno session.Okno) models.Zapytanie {

	// Załączników tu nie ma i nic ich nie udaje: prowenancja pokazuje wiersz
	// wywołania i prompt systemowy, a nie treść pytania; wymyślony załącznik
	// byłby ścieżką, której Operator nie dołączył.
	zapytanie := zapytanieKanalu(okno, shared.Message{}, shared.Message{}, nil)
	// Wznowienie rozmowy wchodzi do argv przełącznikiem --resume, więc bez niego
	// wiersz byłby wierszem innej tury niż ta, która pójdzie.
	if a.rozmowa.ciaglosc != nil {
		zapytanie.Wznowienie = a.rozmowa.ciaglosc.Przypomnij(ctx, okno.Id)
	}
	zapytanie.Nakladka = a.rozmowa.nakladkaOkna(ctx, okno)
	a.rozmowa.uzupelnijSrodowisko(ctx, okno, &zapytanie)
	a.rozmowa.uzupelnijKonfiguracje(ctx, okno, &zapytanie)
	// Ekspert okna nakłada model, ustawienia i mosty MCP — w tej samej kolejności
	// co w turze. Pominięcie go tutaj pokazywałoby Operatorowi wiersz sprzed
	// wyboru eksperta i nazywałoby go prowenancją bieżącej tury.
	a.rozmowa.uzupelnijAgenta(ctx, okno, &zapytanie)
	return zapytanie
}

// oknoProwenancji ustala okno, dla którego liczona jest prowenancja. Wskazanie
// wprost wygrywa; karta sesji zastępuje je tylko wtedy, gdy rozstrzyga
// jednoznacznie.
//
// Karta z wieloma oknami jest odmową, nie zgadywaniem. Sesja niesie okna
// w relacji 1:N, a każde okno ma własny kanał, własny katalog i własną
// konfigurację — wybranie „pierwszego z brzegu" pokazałoby Operatorowi
// prowenancję cudzego okna i wyglądałoby na odpowiedź. Rdzeń nie zna tu ogniska
// klienta, więc mówi wprost, czego mu brakuje.
func (a *adapterProwenancjiKonfiguracji) oknoProwenancji(idOkna, idSesji *string) (session.Okno, error) {
	rejestr := a.rozmowa.nadzorca.Rejestr()
	if wskazane := strings.TrimSpace(wartoscTekstu(idOkna)); wskazane != "" {
		okno, err := rejestr.Okno(wskazane)
		if err != nil {
			return session.Okno{}, bladSesji(err)
		}
		return okno, nil
	}
	karta := strings.TrimSpace(wartoscTekstu(idSesji))
	if karta == "" {
		return session.Okno{}, bladProwenancji(shared.ErrorCodeValidationFailed,
			"żądanie bez wskazania okna i bez karty sesji")
	}
	okna, err := rejestr.OknaSesji(karta)
	if err != nil {
		return session.Okno{}, bladSesji(err)
	}
	if len(okna) == 0 {
		return session.Okno{}, bladProwenancji(shared.ErrorCodeNotFound,
			"karta sesji "+karta+" nie ma otwartego okna komunikacji")
	}
	if len(okna) > 1 {
		return session.Okno{}, bladProwenancji(shared.ErrorCodeValidationFailed,
			"karta sesji "+karta+" niesie okien: "+strconv.Itoa(len(okna))+
				"; prowenancja jest własnością jednego okna — wskaż windowId")
	}
	return okna[0], nil
}

// argvProwenancji składa wiersz wywołania procesu modelu.
//
// Kanał bez procesu nie ma wiersza i nic tego nie udaje. Adaptery `echo` i `api`
// nie uruchamiają programu — dla nich `argv` jest pustą tablicą (nie brakiem
// pola i nie zmyślonym wierszem `claude …`), dokładnie tak, jak rozstrzyga
// prowenancja tych kanałów w `models/prowenancja.go`. Pozostałe pola
// odpowiedzi — prompt systemowy, ustawienia, suma konstytucji — są dla nich
// równie prawdziwe jak dla kanału głównego.
func argvProwenancji(d models.Definicja, u injection.Ustawienia, n injection.Nakladka) []string {
	if !kanalZProcesem(d) {
		return []string{}
	}
	argumenty := injection.Argumenty(u, n)
	argv := make([]string, 0, len(argumenty)+1)
	argv = append(argv, u.Program)
	return append(argv, argumenty...)
}

// kanalZProcesem mówi, czy wiersz rejestru jedzie kanałem uruchamiającym proces.
// Oba klucze wskazują tę samą fabrykę co montaż (`montaz_zrodla.go`), więc
// rozpoznanie nie rozjedzie się z tym, co naprawdę startuje proces.
func kanalZProcesem(d models.Definicja) bool {
	switch d.KluczAdaptera() {
	case models.AdapterCLI, rodzajKanaluLokalny:
		return true
	default:
		return false
	}
}

// ustawieniaProwenancji podaje treść pliku ustawień jako wartość JSON.
//
// Trzy stany, każdy nazwany wprost: treść JSON jedzie sobą; wartość, która jest
// ścieżką pliku zapisanego wcześniej przez inną warstwę, jedzie jako napis
// (surowa wywróciłaby kodowanie odpowiedzi); brak wartości jedzie jako `null`,
// co znaczy „proces nie dostanie przełącznika --settings" — a nie „przekazano
// pusty zestaw reguł".
func ustawieniaProwenancji(tresc string) json.RawMessage {
	przyciete := strings.TrimSpace(tresc)
	if przyciete == "" {
		return json.RawMessage("null")
	}
	if json.Valid([]byte(przyciete)) {
		return json.RawMessage(przyciete)
	}
	zakodowana, err := json.Marshal(przyciete)
	if err != nil {
		return json.RawMessage("null")
	}
	return zakodowana
}

// ── config.window.open ───────────────────────────────────────────────────────

// OtworzOknoKonfiguracji rozstrzyga zakres wejściowy okna konfiguracji.
//
// Kontrakt nie mówi, co rdzeń ma przy niej robić poza podaniem zakresu, więc
// rdzeń robi dokładnie to i nic ponadto: sprawdza, że byt, na którym okno ma
// stanąć, istnieje, że wskazany obszar należy do kontraktu i że wskazany poziom
// zasięgu ma swój byt — a potem oddaje rozstrzygnięty zakres. Okno konfiguracji
// jest bytem klienta; rdzeń nie ma tabeli jego otwarć i nie zakłada jej pod
// jedno pole `opened`.
//
// `opened` nie jest stałą: fałszu ta komenda nie zwraca, zwraca odmowę —
// nieistniejące okno albo karta sesji to `not_found`, obszar spoza kontraktu
// i poziom bez bytu to `validation_failed`. Odpowiedź `opened: false` bez powodu
// byłaby ciszą udającą wynik.
func (a *adapterProwenancjiKonfiguracji) OtworzOknoKonfiguracji(_ context.Context,
	z shared.ConfigWindowOpenRequest) (shared.ConfigWindowOpenResponse, error) {

	if a == nil || a.rozmowa == nil || a.rozmowa.nadzorca == nil {
		return shared.ConfigWindowOpenResponse{}, bladProwenancji(shared.ErrorCodeInternalError,
			"rejestr sesji niewpięty; zakresu okna konfiguracji nie ma na czym sprawdzić")
	}
	obszar, err := obszarOtwarcia(z.Area)
	if err != nil {
		return shared.ConfigWindowOpenResponse{}, err
	}
	poziom, err := a.poziomOtwarcia(z)
	if err != nil {
		return shared.ConfigWindowOpenResponse{}, err
	}
	return shared.ConfigWindowOpenResponse{Opened: true, Area: obszar, Scope: &poziom}, nil
}

// obszarOtwarcia sprawdza wskazany zakres okna. Brak wskazania nie jest błędem
// i nie jest podstawiany żadnym obszarem domyślnym: okno otwiera się wtedy na
// komplecie obszarów, a wynik zostaje bez pola. Rozpoznanie idzie
// jedynym wykazem obszarów w rdzeniu — nazwami pól kontraktu.
func obszarOtwarcia(zadany *shared.SessionConfigArea) (*shared.SessionConfigArea, error) {
	if zadany == nil || *zadany == "" {
		return nil, nil
	}
	if _, znany := znaneObszary[*zadany]; !znany {
		return nil, bladNieznanegoObszaru(*zadany)
	}
	obszar := *zadany
	return &obszar, nil
}

// poziomOtwarcia rozstrzyga poziom zasięgu, na którym okno staje na wejściu.
//
// Bez wskazania obowiązuje poziom najwęższy, jaki żądanie potrafi wskazać bytem:
// okno komunikacji, w jego braku karta sesji, w braku obu poziom globalny.
// Wskazanie wprost jest brane, ale nie na wiarę: poziom okna bez okna i poziom
// karty bez karty ani okna to zakres, którego nie da się pokazać.
func (a *adapterProwenancjiKonfiguracji) poziomOtwarcia(
	z shared.ConfigWindowOpenRequest) (shared.ConfigScope, error) {

	rejestr := a.rozmowa.nadzorca.Rejestr()
	idOkna := strings.TrimSpace(wartoscTekstu(z.WindowId))
	karta := strings.TrimSpace(wartoscTekstu(z.SessionId))
	if idOkna != "" {
		okno, err := rejestr.Okno(idOkna)
		if err != nil {
			return "", bladSesji(err)
		}
		// Okno zna swoją kartę, więc wskazanie okna wystarcza także poziomowi
		// karty sesji — i wystarcza wtedy, gdy klient karty nie podał.
		if karta == "" {
			karta = okno.IdSesji
		}
	}
	if karta != "" {
		if _, err := rejestr.Sesja(karta); err != nil {
			return "", bladSesji(err)
		}
	}
	if z.Scope == nil || *z.Scope == "" {
		switch {
		case idOkna != "":
			return shared.ConfigScopeWindow, nil
		case karta != "":
			return shared.ConfigScopeSession, nil
		default:
			return shared.ConfigScopeGlobal, nil
		}
	}
	poziom := *z.Scope
	if !konfig.Znany(poziom) {
		return "", bladProwenancji(shared.ErrorCodeValidationFailed,
			"poziom zasięgu "+string(poziom)+" nie należy do ośmiu poziomów kontraktu")
	}
	switch poziom {
	case shared.ConfigScopeWindow:
		if idOkna == "" {
			return "", bladProwenancji(shared.ErrorCodeValidationFailed,
				"poziom window bez wskazania okna komunikacji")
		}
	case shared.ConfigScopeSession:
		if karta == "" {
			return "", bladProwenancji(shared.ErrorCodeValidationFailed,
				"poziom session bez wskazania karty sesji ani okna")
		}
	}
	return poziom, nil
}

// bladProwenancji składa odmowę rodziny `config.*` z zamkniętego katalogu kodów.
func bladProwenancji(kod shared.ErrorCode, powod string) error {
	return protocol.JakoError(protocol.NowyBlad(kod, "konfiguracja: "+powod))
}
