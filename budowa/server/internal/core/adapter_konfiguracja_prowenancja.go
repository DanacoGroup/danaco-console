// Plik obsługuje dwie komendy rodziny `config.*`, które nie dotykają rejestru ustawień: `config.explain.get` (prowenancja wywołania modelu) i `config.window.open` (zakres okna konfiguracji); prowenancja idzie tą samą drogą co tura, nie drugą.
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

// adapterProwenancjiKonfiguracji wypełnia port ProwenancjaKonfiguracji: rozszerza adapter ustawień o dwie komendy potrzebujące drogi tury, a nie rejestru ustawień.
type adapterProwenancjiKonfiguracji struct {
	*adapterUstawienOsi
	// rozmowa jest jedynym składaczem wywołania modelu w rdzeniu; zerowy znaczy odmowę, nie ciszę.
	rozmowa *adapterRozmowy
}

// ZProwenancja rozszerza adapter ustawień o `config.explain.get` i `config.window.open`, biorąc ten sam egzemplarz adaptera rozmowy, którym prowadzone są tury.
func (a *adapterUstawienOsi) ZProwenancja(rozmowa *adapterRozmowy) *adapterProwenancjiKonfiguracji {
	return &adapterProwenancjiKonfiguracji{adapterUstawienOsi: a, rozmowa: rozmowa}
}

// ── config.explain.get ───────────────────────────────────────────────────────

// ProwenancjaWywolania oddaje warunki najbliższej tury okna: wiersz wywołania, prompt systemowy, ustawienia procesu i sumę kontrolną konstytucji, liczone na świeżo, nie z zapisu poprzedniej tury.
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
	// Rejestr oddaje wyłącznie kanały czynne; brak kanału w rejestrze odmawia tym samym rozpoznaniem.
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
		// Suma kontrolna dotyczy samej konstytucji, nie całej nakładki roli i ekspertyzy zadaniowej.
		ConstitutionHash: models.SkrotTresci(zapytanie.Nakladka.Konstytucja),
	}, nil
}

// zapytanieProwenancji składa zapytanie kanału dokładnie tak, jak `prowadzTure` — w tej samej kolejności i tymi samymi funkcjami, bez treści wypowiedzi ani historii.
func (a *adapterProwenancjiKonfiguracji) zapytanieProwenancji(ctx context.Context,
	okno session.Okno) models.Zapytanie {

	// Załączników tu nie ma i nic ich nie udaje: prowenancja pokazuje wiersz wywołania i prompt.
	zapytanie := zapytanieKanalu(okno, shared.Message{}, shared.Message{}, nil)
	// Wznowienie rozmowy wchodzi do argv przełącznikiem --resume.
	if a.rozmowa.ciaglosc != nil {
		zapytanie.Wznowienie = a.rozmowa.ciaglosc.Przypomnij(ctx, okno.Id)
	}
	zapytanie.Nakladka = a.rozmowa.nakladkaOkna(ctx, okno)
	a.rozmowa.uzupelnijSrodowisko(ctx, okno, &zapytanie)
	a.rozmowa.uzupelnijKonfiguracje(ctx, okno, &zapytanie)
	// Ekspert okna nakłada model, ustawienia i mosty MCP w tej samej kolejności co w turze.
	a.rozmowa.uzupelnijAgenta(ctx, okno, &zapytanie)
	return zapytanie
}

// oknoProwenancji ustala okno, dla którego liczona jest prowenancja: wskazanie wprost wygrywa, karta sesji zastępuje je tylko wtedy, gdy rozstrzyga jednoznacznie.
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

// argvProwenancji składa wiersz wywołania procesu modelu; kanał bez procesu nie ma wiersza, a adaptery `echo` i `api` oddają dla niego pustą tablicę.
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

// ustawieniaProwenancji podaje treść pliku ustawień jako wartość JSON, w trzech nazwanych wprost stanach: treść JSON, ścieżka pliku jako napis, albo `null` przy braku wartości.
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

// OtworzOknoKonfiguracji rozstrzyga zakres wejściowy okna konfiguracji: sprawdza istnienie bytu, przynależność obszaru do kontraktu i byt poziomu zasięgu, po czym oddaje rozstrzygnięty zakres.
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

// obszarOtwarcia sprawdza wskazany zakres okna. Brak wskazania nie jest błędem: okno otwiera się na komplecie obszarów, rozpoznawanych wyłącznie po nazwach pól kontraktu.
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

// poziomOtwarcia rozstrzyga poziom zasięgu, na którym okno staje na wejściu: bez wskazania obowiązuje poziom najwęższy, jaki żądanie potrafi wskazać bytem.
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
		// Okno zna swoją kartę, więc jego wskazanie wystarcza także poziomowi karty sesji.
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

// bladProwenancji składa odmowę rodziny `config.*` z zamkniętego katalogu kodów błędów tego kontraktu komendy.
func bladProwenancji(kod shared.ErrorCode, powod string) error {
	return protocol.JakoError(protocol.NowyBlad(kod, "konfiguracja: "+powod))
}
