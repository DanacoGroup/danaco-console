// Adapter bramki dostępu Operatora: zakłada jedyne konto komendą auth.register
// i otwiera je komendą auth.login; pozostałe czynności rodziny i postać
// sekretu leżą w sąsiednich plikach pakietu.
package core

import (
	"context"
	"errors"
	"strconv"
	"strings"
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

	// zmianaAdresu prowadzi zamówione zmiany adresu uwierzytelniającego.
	zmianaAdresu dane.RepozytoriumZmianyAdresu

	// nadajnik to konto nadawcze platformy podane przy starcie; nastawy
	// nakładają na nie kontoNadawcze.
	nadajnik nadajnik.Nastawy

	// adresKonsoli to publiczny adres, pod który kieruje odsyłacz z listu
	// aktywacji. Pusty znaczy adres wbudowany w pakiet.
	adresKonsoli string

	// nastawy daje odczyt konta nadawczego z okna Konfiguracji; zerowe —
	// obowiązuje samo konto startowe.
	nastawy NastawyPlatformy

	// zamekZmiany szereguje sprawdzenie i zapis bramki, żeby równoległe zmiany
	// hasła się nie nadpisały.
	zamekZmiany sync.Mutex

	// dlawik spowalnia zgadywanie sekretu zwłoką rosnącą po niepowodzeniach
	// i zerowaną udanym wejściem.
	dlawik *dlawikDrog
}

// nowyAdapterUwierzytelnienia składa adapter bramki nad repozytorium
// uwierzytelnienia i zakłada dławik ograniczający tempo prób wejścia.
func nowyAdapterUwierzytelnienia(repozytorium dane.RepozytoriumUwierzytelnienia) *adapterUwierzytelnienia {
	return &adapterUwierzytelnienia{repozytorium: repozytorium, dlawik: nowyDlawikDrog()}
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

// ZZmianaAdresu wpina trwałość zamówionych zmian adresu. Bez niej zmiana adresu
// odmawia, bo zamówienia nie ma gdzie zapisać.
func (a *adapterUwierzytelnienia) ZZmianaAdresu(
	zmiana dane.RepozytoriumZmianyAdresu) *adapterUwierzytelnienia {

	a.zmianaAdresu = zmiana
	return a
}

// ZNadajnikiem wpina konto nadawcze platformy — to, którym aplikacja pisze
// do Operatora przy rejestracji i przy odzyskiwaniu konta.
func (a *adapterUwierzytelnienia) ZNadajnikiem(n nadajnik.Nastawy) *adapterUwierzytelnienia {
	a.nadajnik = n
	return a
}

// ZAdresemKonsoli wpina publiczny adres Konsoli. Odsyłacz z listu ma otworzyć
// okno na maszynie Operatora, więc nie może być adresem nasłuchu rdzenia.
func (a *adapterUwierzytelnienia) ZAdresemKonsoli(adres string) *adapterUwierzytelnienia {
	a.adresKonsoli = adres
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
		if err := a.zapiszZnacznikBezPoczty(ctx, kontoId, email); err != nil {
			a.cofnijRejestracje(ctx, kotwica.Kod, kontoId)
			return shared.AuthRegisterResponse{}, err
		}
		return shared.AuthRegisterResponse{Registered: true, PendingVerification: false}, nil
	}
	/* Niepowodzenie wysyłki nie cofa rejestracji — czynność schodzi na drogę
	   bez poczty. Powód idzie do dziennika: odpowiedź kontraktu niesie samo
	   `pendingVerification: false`, więc bez tego zapisu przyczyna — brak konta
	   nadawczego czy odmowa serwera pocztowego — przepadałaby bez śladu. */
	if err := a.wyslijDrogePotwierdzenia(ctx, dane.CelWeryfikacja, email, kontoId); err != nil {
		if dziennik := dziennikZKontekstu(ctx); dziennik != nil {
			dziennik.Printf("rejestracja: list z kodem nie wyszedł na %s: %v", email, err)
		}
		if err := a.zapiszZnacznikBezPoczty(ctx, kontoId, email); err != nil {
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
	// Znacznik bramki bez poczty odchodzi z kontem, żeby nie przeżył w sejfie
	// usuniętej rejestracji.
	a.zdejmijZnacznikBezPoczty(ctx, kontoId)
	if a.konto != nil && kontoId != 0 {
		_ = a.konto.UsunKonto(ctx, kontoId)
	}
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
		a.wyrownajNieznanyLogin(ctx, z)
		return shared.AuthLoginResponse{}, err
	}
	if err := a.kontoPotwierdzone(ctx, konto); err != nil {
		return shared.AuthLoginResponse{}, err
	}
	// Zwłoka nakłada się przed sprawdzeniem sekretu, żeby czas odpowiedzi nie
	// zdradzał wyniku, i trzyma drogę zajętą do rozstrzygnięcia próby.
	proba := a.dlawik.Podejdz(ctx, kluczDlawika(ctx, czynnoscWejscia, konto.Id))
	defer proba.Zwolnij()

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
		proba.Niepowodzenie()
		// Sekret niezgodny to nieudane wejście, nie wadliwe żądanie, więc kod
		// jest not_authenticated.
		return shared.AuthLoginResponse{}, bladBramki(shared.ErrorCodeNotAuthenticated,
			"Nie rozpoznano danych logowania.")
	}
	// Wejście udane zeruje licznik zwłoki, więc kolejne pomyłki liczą się
	// od nowa.
	proba.Wyzeruj()
	if err := a.repozytorium.OdnotujUzycie(ctx, metoda.Kod, time.Now().UnixMilli()); err != nil {
		return shared.AuthLoginResponse{}, err
	}
	// Urządzenie sesji bierze się z żądania; wiersz metody uzupełnia je tylko,
	// gdy żądanie go nie poda.
	urzadzenie := niepustyTekst(z.DeviceId)
	if urzadzenie == nil {
		urzadzenie = metoda.UrzadzenieKod
	}
	/* Konto sesji wskazuje metoda, o ile je niesie: PIN należy do urządzenia
	   i loginu nie niesie, więc bez tego wskazania wejście PIN-em zakładałoby
	   sesję konta najstarszego, czyli cudzą. */
	kontoSesji := konto.Id
	if metoda.KontoId != 0 {
		kontoSesji = metoda.KontoId
	}
	sesja, err := a.zalozSesje(ctx, metoda.Rodzaj, urzadzenie, wartoscPrawdy(z.KeepSignedIn), kontoSesji)
	if err != nil {
		return shared.AuthLoginResponse{}, err
	}
	/* Wykaz metod idzie kontem, które właśnie weszło: rozpoznanie połączenia
	   stanęło przed tym logowaniem, więc kontekst żądania niesie jeszcze stan
	   sprzed wejścia. */
	metody, err := a.MetodyWejscia(dane.ZKontemOperatora(ctx, kontoSesji))
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

// MetodyWejscia oddaje komplet metod konta wołającego w kształcie kontraktu.
// Zasila odpowiedzi trzech komend i ładunek zdarzenia `auth.changed`.
//
// Bez rozpoznanego konta wykaz jest pusty, nie pełny: niesie etykiety, kody
// metod i nazwy urządzeń, a po samym kodzie metodę da się zdjąć.
func (a *adapterUwierzytelnienia) MetodyWejscia(ctx context.Context) ([]shared.AuthMethod, error) {
	if dane.KontoOperatora(ctx) == 0 {
		return []shared.AuthMethod{}, nil
	}
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

// kotwicaKonta zwraca hasło wskazanego konta. Wynik `dane.ErrBrakWiersza`
// przechodzi bez przekładu — wołający rozstrzyga, czy brak kotwicy jest tu
// odmową, czy zgodą.
//
// Konto podaje się zawsze, także zerem: kotwicy bez wskazania konta nie ma,
// bo hasło należy do konta, a nie do platformy.
func (a *adapterUwierzytelnienia) kotwicaKonta(ctx context.Context,
	kontoId int64) (dane.MetodaUwierzytelnienia, error) {

	return a.repozytorium.KotwicaKonta(ctx, kontoId)
}

/*
kontoWolajacego oddaje konto, do którego należy żądanie. Rozpoznanie stawia
rdzeń z sesji bramki związanej z gniazdem — tutaj zostaje sam odczyt.

Zero jest odmową nazwaną, nie wartością zastępczą: czynność, która zmienia
metodę wejścia albo hasło, musi wiedzieć, czyje one są, a zgadnięcie „konto
najstarsze" oddawało bramkę Właściciela pierwszemu, kto o nią poprosił.
*/
func kontoWolajacego(ctx context.Context) (int64, error) {
	kontoId := dane.KontoOperatora(ctx)
	if kontoId == 0 {
		return 0, bladBramki(shared.ErrorCodeNotAuthenticated,
			"Nie wiadomo, czyje jest to połączenie — zaloguj się ponownie.")
	}
	return kontoId, nil
}

// czynnoscWejscia i czynnoscDrogi nazywają dwie drogi dławika: wejście sekretem
// i wpisanie kodu z listu. Osobne, bo zgadywanie jednego nie ma opóźniać
// drugiego — Operator, któremu ktoś obcy dławi logowanie, ma dalej móc
// potwierdzić adres.
const (
	czynnoscWejscia = "wejscie"
	czynnoscDrogi   = "droga"
)

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

// ── dławik prób ──────────────────────────────────────────────────────────────

// wygasanieDrogiDlawika — po tym czasie bez próby wpis drogi znika wraz z jej
// licznikiem. Mapa bez wygaszania rośnie o wpis na każde połączenie i każde
// konto przez całe życie procesu.
const wygasanieDrogiDlawika = 15 * time.Minute

/*
dlawikDrog nakłada zwłokę na próby zgadywania sekretu i kodu z listu.

Klucz drogi nie pochodzi z żądania: każda wartość podana przez proszącego —
`deviceId` na czele — zakłada przy podmianie drogę nową, z licznikiem od zera,
więc zwłoka nie obowiązuje nikogo. Kluczem jest wskazane konto albo połączenie,
czyli to, czego proszący sam nie podmienia.

Zwłoka nakłada się pod zamkiem drogi. Próby puszczone naraz czekają inaczej
wspólnie tyle, ile jedna, i cała zwłoka sprowadza się do jednego odstępu.
*/
type dlawikDrog struct {
	licznik *dlawikWejscia

	zamek sync.Mutex
	drogi map[string]*drogaDlawiona
}

// drogaDlawiona to jedna droga prób: zamek szeregujący próby i chwila ostatniej
// z nich, po której wpis wygasa.
type drogaDlawiona struct {
	zamek    sync.Mutex
	wZajeciu int
	ostatnia time.Time
}

// zajecieDrogi jest uchwytem próby trwającej. Wołający oddaje go zawsze —
// zwolnienie zdejmuje zamek drogi, a rozstrzygnięcie próby idzie osobno.
type zajecieDrogi struct {
	dlawik *dlawikDrog
	klucz  string
	droga  *drogaDlawiona
}

// nowyDlawikDrog zakłada dławik z pustą mapą dróg.
func nowyDlawikDrog() *dlawikDrog {
	return &dlawikDrog{licznik: nowyDlawikWejscia(), drogi: map[string]*drogaDlawiona{}}
}

// Podejdz zajmuje drogę i nakłada należną jej zwłokę. Wraca dopiero wtedy, gdy
// próba wolno się rozstrzygnąć.
func (d *dlawikDrog) Podejdz(ctx context.Context, klucz string) zajecieDrogi {
	if d == nil {
		return zajecieDrogi{}
	}
	teraz := time.Now()
	d.zamek.Lock()
	d.wygas(teraz)
	droga, jest := d.drogi[klucz]
	if !jest {
		droga = &drogaDlawiona{}
		d.drogi[klucz] = droga
	}
	droga.wZajeciu++
	droga.ostatnia = teraz
	d.zamek.Unlock()

	droga.zamek.Lock()
	d.licznik.Zaczekaj(ctx, klucz)
	return zajecieDrogi{dlawik: d, klucz: klucz, droga: droga}
}

// Zwolnij oddaje drogę następnej próbie.
func (z zajecieDrogi) Zwolnij() {
	if z.dlawik == nil {
		return
	}
	z.droga.zamek.Unlock()
	z.dlawik.zamek.Lock()
	defer z.dlawik.zamek.Unlock()
	z.droga.wZajeciu--
	z.droga.ostatnia = time.Now()
}

// Niepowodzenie podnosi zwłokę należną próbie następnej na tej samej drodze.
func (z zajecieDrogi) Niepowodzenie() {
	if z.dlawik == nil {
		return
	}
	z.dlawik.licznik.Niepowodzenie(z.klucz)
}

// Wyzeruj kasuje licznik drogi; woła się po próbie udanej, nie po każdej.
func (z zajecieDrogi) Wyzeruj() {
	if z.dlawik == nil {
		return
	}
	z.dlawik.licznik.Wyzeruj(z.klucz)
}

// wygas zdejmuje drogi bez próby dłużej niż wygasanieDrogiDlawika. Droga zajęta
// zostaje: jej próba właśnie trwa. Wołane pod zamkiem dławika.
func (d *dlawikDrog) wygas(teraz time.Time) {
	for klucz, droga := range d.drogi {
		if droga.wZajeciu > 0 || teraz.Sub(droga.ostatnia) < wygasanieDrogiDlawika {
			continue
		}
		delete(d.drogi, klucz)
		d.licznik.Wyzeruj(klucz)
	}
}

/*
kluczDlawika składa drogę licznika z czynności i z tego, kogo próba dotyczy.

Konto wskazane żądaniem jest tu wartością pewną — rdzeń odczytał je z bazy,
zanim doszło do sprawdzenia sekretu. Gdy czynność konta jeszcze nie zna, zostaje
tożsamość połączenia nadana przez transport. Żądanie spoza gniazda ma jedną
drogę wspólną: praca wewnętrzna rdzenia nie zgaduje sekretów, a rozdzielenie jej
na drogi wymagałoby wartości, której nie ma.
*/
// wyrownajNieznanyLogin nakłada na login nieznany tę samą zwłokę i ten sam koszt
// sprawdzenia sekretu, co na login znany — inaczej czas odpowiedzi i brak
// dławika zdradzałyby, które loginy mają konto.
func (a *adapterUwierzytelnienia) wyrownajNieznanyLogin(ctx context.Context, z shared.AuthLoginRequest) {
	proba := a.dlawik.Podejdz(ctx, kluczDlawika(ctx, czynnoscWejscia, 0))
	defer proba.Zwolnij()
	sekret := ""
	if z.Secret != nil {
		sekret = *z.Secret
	}
	_, _ = sekretZgadzaSie(zapisSekretuWzorcowego(), sekret)
	proba.Niepowodzenie()
}

var (
	zapisWzorcowyRaz   sync.Once
	zapisWzorcowyTekst string
)

// zapisSekretuWzorcowego oddaje zapis PBKDF2 sekretu pustego, liczony raz, do
// porównań o tym samym koszcie przy loginie bez konta.
func zapisSekretuWzorcowego() string {
	zapisWzorcowyRaz.Do(func() {
		zapisWzorcowyTekst, _ = zapisSekretu("")
	})
	return zapisWzorcowyTekst
}

func kluczDlawika(ctx context.Context, czynnosc string, kontoId int64) string {
	if kontoId != 0 {
		return czynnosc + "\x00konto:" + strconv.FormatInt(kontoId, 10)
	}
	if polaczenie := polaczenieZKontekstu(ctx); polaczenie != "" {
		return czynnosc + "\x00polaczenie:" + polaczenie
	}
	return czynnosc + "\x00bez zrodla"
}
