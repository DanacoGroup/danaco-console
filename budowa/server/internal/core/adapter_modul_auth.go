// Adapter bramki dostępu Operatora: zakłada jedyne konto komendą auth.register
// i otwiera je komendą auth.login; pozostałe czynności rodziny i postać
// sekretu leżą w sąsiednich plikach pakietu.
package core

import (
	"strings"
	"context"
	"errors"
	"sync"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/nadajnik"
	"danacoconsole/shared"
)

// Czas życia sesji bramki bez odnowienia; kontrakt wymaga wygasania
// przesuwnego, lecz nie ustala długości, więc wartość przyjmuje rdzeń.
const trwanieSesjiBramki = 12 * time.Hour

// Czas życia sesji bramki po zaznaczeniu opcji trwałego logowania; wygasanie
// nadal obowiązuje, tylko z dłuższym oknem odnawianym każdym wejściem.
const trwanieSesjiBramkiDlugie = 365 * 24 * time.Hour

// przedrostekBytuSejfu znakuje wpisy bramki w sejfie poświadczeń, żeby nie
// mieszały się z poświadczeniami kont i punktów dostępu w tym samym pliku.
const przedrostekBytuSejfu = "auth:"

// adapterUwierzytelnienia wypełnia port Uwierzytelnianie.
//
// `sejf` bywa zerowy, gdy montaż go nie wpiął. Nie jest to cichy brak skutku:
// każda czynność dotykająca sekretu odmawia wtedy wprost, bo hasła nie ma gdzie
// odłożyć ani z czym porównać.
type adapterUwierzytelnienia struct {
	repozytorium dane.RepozytoriumUwierzytelnienia
	sejf         SejfPoswiadczen

	// konto trzyma tożsamość właściciela i drogi potwierdzenia; zerowe konto
	// odmawia rejestracji wprost.
	konto dane.RepozytoriumKontaWlasciciela

	// nadajnik to konto nadawcze platformy podane przy starcie; nastawy
	// nakładają na nie kontoNadawcze.
	nadajnik nadajnik.Nastawy

	// nastawy daje odczyt konta nadawczego z okna Konfiguracji; zerowe —
	// obowiązuje samo konto startowe.
	nastawy NastawyPlatformy

	// zamekZmiany szereguje sprawdzenie i zapis bramki, żeby równoległe zmiany
	// hasła się nie nadpisały.
	zamekZmiany sync.Mutex

	// dlawik spowalnia zgadywanie sekretu zwłoką rosnącą po niepowodzeniach
	// i zerowaną udanym wejściem.
	dlawik *dlawikWejscia
}

// nowyAdapterUwierzytelnienia składa adapter bramki nad repozytorium
// uwierzytelnienia i zakłada dławik ograniczający tempo prób wejścia.
func nowyAdapterUwierzytelnienia(repozytorium dane.RepozytoriumUwierzytelnienia) *adapterUwierzytelnienia {
	return &adapterUwierzytelnienia{repozytorium: repozytorium, dlawik: nowyDlawikWejscia()}
}

// ZSejfem wpina magazyn sekretów do adaptera; bez wpiętego sejfu żadna
// czynność dotykająca hasła bramki się nie wykona.
func (a *adapterUwierzytelnienia) ZSejfem(sejf SejfPoswiadczen) *adapterUwierzytelnienia {
	a.sejf = sejf
	return a
}

// ZKontem wpina trwałość tożsamości właściciela. Bez niej rejestracja odmawia,
// bo konta nie ma gdzie zapisać.
func (a *adapterUwierzytelnienia) ZKontem(
	konto dane.RepozytoriumKontaWlasciciela) *adapterUwierzytelnienia {

	a.konto = konto
	return a
}

// ZNadajnikiem wpina konto nadawcze platformy — to, którym aplikacja pisze
// do Operatora przy rejestracji i przy odzyskiwaniu konta.
func (a *adapterUwierzytelnienia) ZNadajnikiem(n nadajnik.Nastawy) *adapterUwierzytelnienia {
	a.nadajnik = n
	return a
}

// ZNastawamiPlatformy wpina drogę odczytu konta nadawczego zapisanego przez
// Operatora w oknie Konfiguracji. Bez niej obowiązuje samo konto ze startu.
func (a *adapterUwierzytelnienia) ZNastawamiPlatformy(n NastawyPlatformy) *adapterUwierzytelnienia {
	a.nastawy = n
	return a
}

// kontoNadawcze oddaje konto nadawcze obowiązujące w tej chwili: podane przy
// starcie i przesłonięte nastawami Operatora.
func (a *adapterUwierzytelnienia) kontoNadawcze(ctx context.Context) nadajnik.Nastawy {
	return kontoNadawcze(ctx, a.nadajnik, a.nastawy)
}

// ── auth.register ────────────────────────────────────────────────────────────

// ZalozBramke zakłada jedyne konto właściciela: login, adres e-mail i hasło.
// Bez skonfigurowanej poczty pomija wysyłkę listu potwierdzającego i zostawia
// adres niepotwierdzony; z pocztą wysyła list i czeka na auth.verify.
func (a *adapterUwierzytelnienia) ZalozBramke(ctx context.Context,
	z shared.AuthRegisterRequest) (shared.AuthRegisterResponse, error) {

	if err := a.gotowa(); err != nil {
		return shared.AuthRegisterResponse{}, err
	}
	if err := a.kontoGotowe(); err != nil {
		return shared.AuthRegisterResponse{}, err
	}
	login, email, err := daneRejestracji(z)
	if err != nil {
		return shared.AuthRegisterResponse{}, err
	}
	if z.Password == "" {
		return shared.AuthRegisterResponse{}, bladBramki(shared.ErrorCodeValidationFailed,
			"Podaj hasło.")
	}
	// Konto nadawcze nie jest warunkiem założenia bramki, tylko warunkiem
	// wysyłki listu.
	pocztaJest := a.kontoNadawcze(ctx).Brak() == nil

	/* Zamek obejmuje zapis konta i jego hasła: dwa równoległe żądania o tym
	   samym loginie mają dać jedno konto i jedną odmowę, a nie dwa wiersze.
	   Sprawdzenia „czy platforma ma już konto" tu nie ma — platforma prowadzi
	   dowolną liczbę kont, a jednoznaczności pilnują wskaźniki na loginie
	   i adresie (migracja 406, `decyzje.md` poz. 22). */
	a.zamekZmiany.Lock()
	defer a.zamekZmiany.Unlock()

	teraz := time.Now().UnixMilli()
	kontoId, err := a.konto.ZalozKonto(ctx, dane.KontoWlasciciela{
		Login: login, Email: email, Potwierdzone: false, Utworzono: teraz,
	})
	if err != nil {
		if errors.Is(err, dane.ErrKolizjaWiersza) {
			return shared.AuthRegisterResponse{}, bladBramkiZPowodem(shared.ErrorCodeConflict,
				PowodKolizjaDanych,
				"Podany login albo adres e-mail należy już do istniejącego konta.")
		}
		return shared.AuthRegisterResponse{}, err
	}
	kotwica, err := a.zalozMetode(ctx, dane.MetodaUwierzytelnienia{
		Rodzaj:          shared.AuthMethodKindPassword,
		Etykieta:        wskaznikTekstu("Hasło konta"),
		NazwaUrzadzenia: niepustyTekst(z.DeviceName),
		Kotwica:         true,
		KontoId:         kontoId,
		Utworzono:       teraz,
	}, z.Password)
	if err != nil {
		a.cofnijRejestracje(ctx, "", kontoId)
		return shared.AuthRegisterResponse{}, err
	}
	// Bez poczty rejestracja kończy się tutaj: konto jest, hasło otwiera
	// bramkę, adres niepotwierdzony.
	if !pocztaJest {
		if err := a.zapiszZnacznikBezPoczty(ctx, email); err != nil {
			a.cofnijRejestracje(ctx, kotwica.Kod, kontoId)
			return shared.AuthRegisterResponse{}, err
		}
		return shared.AuthRegisterResponse{Registered: true, PendingVerification: false}, nil
	}
	// Niepowodzenie wysyłki nie cofa rejestracji — czynność schodzi na drogę
	// bez poczty.
	if err := a.wyslijDrogePotwierdzenia(ctx, dane.CelWeryfikacja, email, login, kontoId); err != nil {
		if err := a.zapiszZnacznikBezPoczty(ctx, email); err != nil {
			a.cofnijRejestracje(ctx, kotwica.Kod, kontoId)
			return shared.AuthRegisterResponse{}, err
		}
		return shared.AuthRegisterResponse{Registered: true, PendingVerification: false}, nil
	}
	return shared.AuthRegisterResponse{Registered: true, PendingVerification: true}, nil
}

// cofnijRejestracje zdejmuje to, co rejestracja zdążyła zapisać: metodę,
// poświadczenie w sejfie, konto i znacznik bramki bez poczty; niepowodzenie
// cofnięcia trafia do dziennika, nie do odpowiedzi.
func (a *adapterUwierzytelnienia) cofnijRejestracje(ctx context.Context, kodKotwicy string, kontoId int64) {
	if kodKotwicy != "" {
		_, _ = a.repozytorium.UsunMetode(ctx, kodKotwicy)
		usunPoswiadczenie(ctx, a.sejf, przedrostekBytuSejfu+kodKotwicy)
	}
	if a.konto != nil && kontoId != 0 {
		_ = a.konto.UsunKonto(ctx, kontoId)
	}
	// Znacznik bramki bez poczty odchodzi z kontem, żeby nie przeżył w sejfie
	// usuniętej rejestracji.
	a.zdejmijZnacznikBezPoczty(ctx)
}

// ── auth.login ───────────────────────────────────────────────────────────────

// WejdzPrzezBramke otwiera bramkę hasłem albo PIN-em i zakłada sesję; konto
// niepotwierdzone bramki nie otwiera, bo adres jest jedyną drogą odzyskania
// dostępu.
func (a *adapterUwierzytelnienia) WejdzPrzezBramke(ctx context.Context,
	z shared.AuthLoginRequest) (shared.AuthLoginResponse, error) {

	if err := a.gotowa(); err != nil {
		return shared.AuthLoginResponse{}, err
	}
	/* Konto wskazuje login albo adres z żądania. Przy wielu kontach hasło samo
	   nie mówi, czyje jest — bez tego wskazania wejście otwierałoby konto
	   najstarsze niezależnie od tego, kto się loguje. */
	konto, err := a.kontoWejscia(ctx, z)
	if err != nil {
		return shared.AuthLoginResponse{}, err
	}
	if err := a.kontoPotwierdzone(ctx, konto); err != nil {
		return shared.AuthLoginResponse{}, err
	}
	// Zwłoka nakłada się przed sprawdzeniem sekretu, żeby czas odpowiedzi nie
	// zdradzał wyniku.
	droga := drogaWejscia(string(z.Method), wartoscTekstu(z.DeviceId))
	a.dlawik.Zaczekaj(ctx, droga)

	metoda, err := a.metodaWejscia(ctx, z, konto)
	if err != nil {
		return shared.AuthLoginResponse{}, err
	}
	if z.Secret == nil || *z.Secret == "" {
		return shared.AuthLoginResponse{}, bladBramki(shared.ErrorCodeValidationFailed,
			"wejście metodą "+string(z.Method)+" bez sekretu")
	}
	zgadza, err := a.sekretZgadzaSieZWpisem(ctx, metoda, *z.Secret)
	if err != nil {
		return shared.AuthLoginResponse{}, err
	}
	if !zgadza {
		// Nieudana próba podnosi zwłokę następnej; ta odpowiada od razu, bo
		// wynik już jest znany.
		a.dlawik.Niepowodzenie(droga)
		// Sekret niezgodny to nieudane wejście, nie wadliwe żądanie, więc kod
		// jest not_authenticated.
		return shared.AuthLoginResponse{}, bladBramki(shared.ErrorCodeNotAuthenticated,
			"Nie rozpoznano danych logowania.")
	}
	// Wejście udane zeruje licznik zwłoki, więc kolejne pomyłki liczą się
	// od nowa.
	a.dlawik.Wyzeruj(droga)
	if err := a.repozytorium.OdnotujUzycie(ctx, metoda.Kod, time.Now().UnixMilli()); err != nil {
		return shared.AuthLoginResponse{}, err
	}
	// Urządzenie sesji bierze się z żądania; wiersz metody uzupełnia je tylko,
	// gdy żądanie go nie poda.
	urzadzenie := niepustyTekst(z.DeviceId)
	if urzadzenie == nil {
		urzadzenie = metoda.UrzadzenieKod
	}
	sesja, err := a.zalozSesje(ctx, metoda.Rodzaj, urzadzenie, wartoscPrawdy(z.KeepSignedIn), konto.Id)
	if err != nil {
		return shared.AuthLoginResponse{}, err
	}
	metody, err := a.MetodyWejscia(ctx)
	if err != nil {
		return shared.AuthLoginResponse{}, err
	}
	return shared.AuthLoginResponse{Session: sesja, Methods: metody}, nil
}

// metodaWejscia odnajduje metodę wejścia wskazaną przez żądanie — hasło, PIN
// urządzenia albo odrzuconą metodę hello — i zwraca błąd kontraktu, gdy
// metoda nie istnieje.
/* Konto wskazane loginem albo adresem z żądania. Metody właściwe urządzeniu
   (PIN, klucz) loginu nie niosą — dla nich konto wskazuje sama metoda, więc
   tutaj zostaje konto najstarsze, tak jak przed dołożeniem wielu kont. */
func (a *adapterUwierzytelnienia) kontoWejscia(ctx context.Context,
	z shared.AuthLoginRequest) (dane.KontoWlasciciela, error) {

	if a.konto == nil {
		return dane.KontoWlasciciela{}, nil
	}
	wskazanie := strings.TrimSpace(wartoscTekstu(z.Login))
	if wskazanie == "" {
		konto, err := a.konto.Konto(ctx)
		if errors.Is(err, dane.ErrBrakWiersza) {
			return dane.KontoWlasciciela{}, nil
		}
		return konto, err
	}
	konto, err := a.konto.KontoPoTozsamosci(ctx, wskazanie)
	if errors.Is(err, dane.ErrBrakWiersza) {
		/* Tożsamość nieznana odpowiada tak samo jak hasło niezgodne — inaczej
		   pytanie o kolejne loginy wskazywałoby, które z nich istnieją. */
		return dane.KontoWlasciciela{}, bladBramki(shared.ErrorCodeNotAuthenticated,
			"Nie rozpoznano danych logowania.")
	}
	return konto, err
}

func (a *adapterUwierzytelnienia) metodaWejscia(ctx context.Context,
	z shared.AuthLoginRequest, konto dane.KontoWlasciciela) (dane.MetodaUwierzytelnienia, error) {

	switch z.Method {
	case shared.AuthMethodKindPassword:
		metoda, err := a.repozytorium.KotwicaKonta(ctx, konto.Id)
		if errors.Is(err, dane.ErrBrakWiersza) {
			/* To samo zdanie co przy sekrecie niezgodnym: konto bez hasła i hasło
			   niezgodne muszą wyglądać dla wołającego identycznie, inaczej
			   odpowiedź zdradza, które loginy istnieją. */
			return metoda, bladBramki(shared.ErrorCodeNotAuthenticated,
				"Nie rozpoznano danych logowania.")
		}
		return metoda, err
	case shared.AuthMethodKindPin:
		urzadzenie := wartoscTekstu(z.DeviceId)
		if urzadzenie == "" {
			return dane.MetodaUwierzytelnienia{}, bladBramki(shared.ErrorCodeValidationFailed,
				"Wskaż urządzenie — kod PIN obowiązuje na jednym urządzeniu.")
		}
		metoda, err := a.repozytorium.MetodaUrzadzenia(ctx, shared.AuthMethodKindPin, urzadzenie)
		if errors.Is(err, dane.ErrBrakWiersza) {
			return metoda, bladBramki(shared.ErrorCodeNotFound,
				"PIN urządzenia "+urzadzenie+" nie istnieje")
		}
		return metoda, err
	case shared.AuthMethodKindHello:
		return dane.MetodaUwierzytelnienia{}, bladBramki(shared.ErrorCodeValidationFailed,
			odmowaHello)
	default:
		return dane.MetodaUwierzytelnienia{}, bladBramki(shared.ErrorCodeValidationFailed,
			"metoda wejścia "+string(z.Method)+" nie należy do kontraktu")
	}
}

// ── wspólne ──────────────────────────────────────────────────────────────────

// odmowaHello jest jednym powodem odmowy metody `hello` dla obu komend, które
// ją przyjmują. Jedno zdanie w jednym miejscu, żeby dwie odmowy nie zaczęły
// mówić dwóch rzeczy o tej samej niedostępności.
const odmowaHello = "metoda hello (Windows Hello przez WebAuthn) nie jest zbudowana: " +
	"rozstrzygnięcie Właściciela z 13.08.2026 odkłada ją do chwili wystawienia " +
	"platformy na serwerze, a WebAuthn i tak wywodzi rp_id z pochodzenia dokumentu — " +
	"pod http://127.0.0.1 poprawnego rp_id nie ma. Dziś działają hasło i PIN"

// MetodyWejscia oddaje komplet metod w kształcie kontraktu. Zasila odpowiedzi
// trzech komend i ładunek zdarzenia `auth.changed`.
func (a *adapterUwierzytelnienia) MetodyWejscia(ctx context.Context) ([]shared.AuthMethod, error) {
	wiersze, err := a.repozytorium.Metody(ctx)
	if err != nil {
		return nil, err
	}
	metody := make([]shared.AuthMethod, 0, len(wiersze))
	for _, wiersz := range wiersze {
		metody = append(metody, metodaBramkiKontraktu(wiersz))
	}
	return metody, nil
}

// gotowa odmawia pracy bramce bez podłoża. Repozytorium zerowe znaczy montaż
// bez bazy, sejf zerowy — montaż bez magazynu sekretów; w obu razach hasła nie
// ma gdzie odłożyć, więc czynność kończy się odmową, a nie pozorem.
func (a *adapterUwierzytelnienia) gotowa() error {
	if a == nil || a.repozytorium == nil {
		return bladBramki(shared.ErrorCodeInternalError,
			"Magazyn kont jest niedostępny.")
	}
	if a.sejf == nil {
		return bladBramki(shared.ErrorCodeInternalError,
			"Magazyn haseł jest niedostępny.")
	}
	return nil
}

// kotwica zwraca hasło bramki. Wynik `dane.ErrBrakWiersza` przechodzi bez
// przekładu — wołający rozstrzyga, czy brak kotwicy jest tu odmową, czy zgodą.
func (a *adapterUwierzytelnienia) kotwica(ctx context.Context) (dane.MetodaUwierzytelnienia, error) {
	return a.repozytorium.Kotwica(ctx)
}

// zalozMetode kładzie sekret w sejfie i dopiero potem wstawia wiersz metody,
// żeby wiersz nie wskazywał na nieistniejący sekret; nieudane wstawienie
// usuwa zapisany sekret.
func (a *adapterUwierzytelnienia) zalozMetode(ctx context.Context,
	metoda dane.MetodaUwierzytelnienia, sekret string) (dane.MetodaUwierzytelnienia, error) {

	metoda.Kod = nowyIdentyfikator("am-")
	zapis, err := zapisSekretu(sekret)
	if err != nil {
		return dane.MetodaUwierzytelnienia{}, bladBramki(shared.ErrorCodeInternalError, err.Error())
	}
	byt := przedrostekBytuSejfu + metoda.Kod
	odwolanie, err := a.sejf.Zapisz(ctx, byt, zapis)
	if err != nil {
		return dane.MetodaUwierzytelnienia{}, err
	}
	metoda.OdwolanieSekretu = odwolanie
	zapisana, err := a.repozytorium.ZalozMetode(ctx, metoda)
	if err != nil {
		usunPoswiadczenie(ctx, a.sejf, byt)
		return dane.MetodaUwierzytelnienia{}, kolizjaMetody(metoda, err)
	}
	return zapisana, nil
}

// kolizjaMetody przekłada odmowę bazy na odmowę kontraktu: naruszenie indeksu
// unikalności zwraca conflict z powodem po ludzku zamiast surowego błędu
// SQLite z nazwami tabel i kolumn.
func kolizjaMetody(metoda dane.MetodaUwierzytelnienia, err error) error {
	if !errors.Is(err, dane.ErrKolizjaWiersza) {
		return err
	}
	if metoda.Kotwica {
		return bladBramki(shared.ErrorCodeConflict,
			"Hasło jest już ustawione. Aby je zmienić, użyj odzyskiwania dostępu. "+
				"zmiana hasła idzie komendą auth.password.reset")
	}
	urzadzenie := ""
	if metoda.UrzadzenieKod != nil {
		urzadzenie = *metoda.UrzadzenieKod
	}
	return bladBramki(shared.ErrorCodeConflict,
		"urządzenie "+urzadzenie+" ma już metodę wejścia rodzaju "+metoda.Rodzaj+
			"; zdejmij ją komendą auth.method.remove, zanim założysz nową")
}

// sekretZgadzaSieZWpisem porównuje przedstawiony sekret z zapisem w sejfie.
// Brak wpisu w sejfie nie jest odmową wejścia, lecz awarią rdzenia, bo wiersz
// metody istnieje, więc sekret powinien tam być.
func (a *adapterUwierzytelnienia) sekretZgadzaSieZWpisem(ctx context.Context,
	metoda dane.MetodaUwierzytelnienia, sekret string) (bool, error) {

	zapis, jest := a.sejf.Odczytaj(ctx, przedrostekBytuSejfu+metoda.Kod)
	if !jest || zapis == "" {
		return false, bladBramki(shared.ErrorCodeInternalError,
			"sejf nie niesie sekretu metody wejścia "+metoda.Kod)
	}
	zgadza, err := sekretZgadzaSie(zapis, sekret)
	if err != nil {
		return false, bladBramki(shared.ErrorCodeInternalError, err.Error())
	}
	return zgadza, nil
}

// zalozSesje wytwarza token, zapisuje sesję po skrócie i oddaje ją w kształcie
// kontraktu — razem z tokenem surowym, bo to jedyna chwila, w której rdzeń go
// zna; w bazie zostaje wyłącznie skrót.
func (a *adapterUwierzytelnienia) zalozSesje(ctx context.Context,
	rodzaj string, urzadzenie *string, niewylogowuj bool, kontoId int64) (shared.AuthSession, error) {

	token, err := nowyTokenBramki()
	if err != nil {
		return shared.AuthSession{}, bladBramki(shared.ErrorCodeInternalError, err.Error())
	}
	teraz := time.Now()
	trwanie := trwanieWejscia(niewylogowuj)
	zapisana, err := a.repozytorium.ZalozSesjeBramki(ctx, dane.SesjaBramki{
		SkrotTokenu:   skrotTokenu(token),
		RodzajMetody:  &rodzaj,
		UrzadzenieKod: urzadzenie,
		Wygasa:        teraz.Add(trwanie).UnixMilli(),
		Utworzono:     teraz.UnixMilli(),
		// Trwanie idzie do wiersza, żeby przełącznik „nie wyloguj mnie”
		// obowiązywał też po odnowieniu sesji.
		Trwanie: trwanie.Milliseconds(),
		// Konto, któremu sesja została wydana; bez tego wskazania sesja nie
		// mówi, czyja jest, a przy wielu kontach to jedyne, co je rozróżnia.
		KontoId: kontoId,
	})
	if err != nil {
		return shared.AuthSession{}, err
	}
	return sesjaBramkiKontraktu(token, zapisana), nil
}

// trwanieWejscia rozstrzyga długość sesji z przełącznika „nie wyloguj mnie”,
// który tylko rzadziej pokazuje okno logowania i nie znosi samego wygasania.
func trwanieWejscia(niewylogowuj bool) time.Duration {
	if niewylogowuj {
		return trwanieSesjiBramkiDlugie
	}
	return trwanieSesjiBramki
}
