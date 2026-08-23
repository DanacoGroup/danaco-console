// Odpowiedzialność pliku: bramka Operatora — założenie bramki (`auth.register`)
// i wejście przez nią (`auth.login`). Cztery pozostałe czynności rodziny leżą
// w `adapter_modul_auth_metody.go`, postać sekretu — w `adapter_modul_auth_sekret.go`.
//
// Bramka jest jedna, a Operator bezimienny: żądania rodziny nie niosą ani
// nazwy, ani adresu e-mail, a `auth.register` nie zakłada konta — ustawia sekret
// bramki przy pierwszym uruchomieniu i od razu wpuszcza. Konto
// z `migracja_014_katalog_kont.sql` to poświadczenie do kanału modelu i z bramką
// nie ma nic wspólnego.
//
// Metoda `hello` (Windows Hello przez WebAuthn) nie jest zbudowana, a miejsce po
// niej nie jest ciszą: `auth.login` i `auth.method.add` odmawiają jej wprost
// i z powodem. WebAuthn wywodzi `rp_id` z pochodzenia dokumentu, a powłoka
// podaje interfejs z `http://127.0.0.1`.
//
// Sekret przechodzi wyłącznie przez sejf. Bez wpiętego sejfu bramka nie
// powstaje — hasła nie ma gdzie odłożyć, więc odmowa niesie `internal_error`.
package core

import (
	"context"
	"errors"
	"sync"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/nadajnik"
	"danacoconsole/shared"
)

// trwanieSesjiBramki — jak długo żyje wejście bez odnowienia.
//
// Wartość jest własna rdzenia: kontrakt mówi „wygasanie przesuwne", ale długości
// nie podaje, a katalog ustawień (`migracja_012_katalog_ustawien.sql`) nie ma
// pozycji na czas życia sesji bramki. Dwanaście godzin to jedna doba robocza —
// Operator, który rano wszedł, nie loguje się w połowie dnia, a maszyna
// zostawiona na noc bramkę zamyka.
const trwanieSesjiBramki = 12 * time.Hour

// trwanieSesjiBramkiDlugie obowiązuje po zaznaczeniu „nie wyloguj mnie".
//
// Nie znosi to wygasania: sesja bez końca byłaby wpisem, którego nic nigdy nie
// sprząta, a Operator nie miałby jak zobaczyć, że coś jeszcze trwa.
// Rok to długość, po której zapomniane urządzenie przestaje wchodzić samo,
// a Operator pracujący codziennie nie zobaczy logowania ani razu — bo wygasanie
// jest przesuwne i każde wejście je odnawia.
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

	// konto trzyma tożsamość właściciela — login, adres e-mail i stan
	// potwierdzenia — oraz jednorazowe drogi potwierdzenia tożsamości. Bywa
	// zerowe, gdy montaż go nie wpiął; rejestracja odmawia wtedy wprost, bo
	// konta nie ma gdzie zapisać.
	konto dane.RepozytoriumKontaWlasciciela

	// nadajnik jest kontem nadawczym PLATFORMY, nie skrzynką Operatora. Pisze
	// dwa listy: potwierdzenie adresu przy rejestracji i drogę odzyskania konta.
	// To jest konto podane przy starcie; nastawy Operatora nakłada na nie
	// `kontoNadawcze` (nastawy_nadajnika.go).
	nadajnik nadajnik.Nastawy

	// nastawy daje drogę odczytu konta nadawczego zapisanego w oknie
	// Konfiguracji. Bywa zerowe — wtedy obowiązuje samo konto ze startu.
	nastawy NastawyPlatformy

	// zamekZmiany szereguje czynności, które sprawdzają stan bramki i zaraz
	// potem go zmieniają: założenie bramki, założenie metody wejścia i zmianę
	// hasła. Bez niego sprawdzenie i zapis są dwiema czynnościami, a między nie
	// wchodzi drugie żądanie z tego samego gniazda — dwie równoległe zmiany hasła
	// odpowiadają wtedy obie `changed: true`, a bramkę otwiera tylko jedno
	// z dwóch nowych haseł, bo drugi zapis do sejfu nadpisuje pierwszy. Sejf jest
	// plikiem z wpisami i transakcji nie zna, więc niepodzielność musi stanąć
	// tutaj. Zamek trzyma się całej czynności, nie samego zapisu.
	zamekZmiany sync.Mutex

	// dlawik spowalnia zgadywanie sekretu bramki (adapter_modul_auth_dlawik.go).
	// Nie odmawia ani jednej próby — nakłada na nią zwłokę rosnącą po
	// niepowodzeniach i zerowaną pierwszym wejściem udanym.
	dlawik *dlawikWejscia
}

// nowyAdapterUwierzytelnienia składa adapter bramki nad repozytorium.
func nowyAdapterUwierzytelnienia(repozytorium dane.RepozytoriumUwierzytelnienia) *adapterUwierzytelnienia {
	return &adapterUwierzytelnienia{repozytorium: repozytorium, dlawik: nowyDlawikWejscia()}
}

// ZSejfem wpina magazyn sekretów. Bez niego bramki założyć się nie da.
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

// ZalozBramke zakłada jedyne konto właściciela: login, adres e-mail
// uwierzytelniający i hasło.
//
// Poczty NIE wymaga. Na świeżej instalacji konta nadawczego platformy nie ma
// jeszcze czym wskazać, a ustawia się je w oknie Konfiguracji — czyli za bramką.
// Rejestracja idzie więc dwiema drogami: z pocztą nadaje list z drogą
// potwierdzenia i zostawia konto niepotwierdzone, bez poczty zakłada konto
// i zapamiętuje w sejfie, że adresu nikt nie potwierdził
// (`adapter_modul_auth_pierwsze_uruchomienie.go`). W obu razach hasło jest
// jedynym, co otwiera bramkę.
//
// Sesji NIE zakłada. Konto powstaje w stanie niepotwierdzonym i pozostaje w nim
// do chwili potwierdzenia adresu komendą `auth.verify` — dopiero potwierdzenie
// wydaje urządzeniu token dostępu. Rejestracja, która wpuszczałaby od razu,
// czyniłaby weryfikację adresu ozdobą: konto działałoby bez niej, a adres
// zostawał niesprawdzony aż do dnia, w którym trzeba nim odzyskać dostęp.
//
// Wykonalna tylko raz: istniejąca kotwica daje `conflict`, nie ciche
// `registered: false`. Powtórzone żądanie jest próbą podmiany hasła bez
// znajomości starego, a od tego jest odzyskanie konta.
//
// List wysyła się PRZED oddaniem odpowiedzi i jego niepowodzenie jest odmową
// całej czynności. Konto założone bez wysłanego listu byłoby kontem, do którego
// nikt nie ma jak wejść — a Operator zobaczyłby „zarejestrowano" i czekał na
// wiadomość, która nigdy nie przyszła.
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
			"rejestracja bez hasła")
	}
	// Konto nadawcze NIE jest warunkiem postawienia bramki — jest warunkiem
	// pisania listów. Świeża instalka poczty nie ma (`mailer.host`
	// i `mailer.address` bez wartości domyślnej, powłoka nie podaje żadnego
	// `DANACO_NADAWCA_*`), a jedyna droga do okna Konfiguracji prowadzi przez
	// bramkę — więc wymóg poczty tutaj zamykał produkt na pierwszym ekranie:
	// „Załóż konto" oddawało `internal_error`, a odmowa radziła naprawę
	// nieosiągalną. Poczta zostaje wymogiem tam, gdzie list jest treścią
	// czynności: przy potwierdzaniu adresu i przy odzyskiwaniu hasła.
	pocztaJest := a.kontoNadawcze(ctx).Brak() == nil

	// Sprawdzenie kotwicy i jej założenie są tu jedną czynnością — inaczej dwa
	// równoległe `auth.register` obchodzą odmowę „konto już jest" i drugie
	// rozbija się dopiero o warunek bazy, oddając surowy błąd SQLite.
	a.zamekZmiany.Lock()
	defer a.zamekZmiany.Unlock()
	if _, err := a.kotwica(ctx); err == nil {
		return shared.AuthRegisterResponse{}, bladBramki(shared.ErrorCodeConflict,
			"konto właściciela jest już założone; rejestracja wykonuje się raz — "+
				"utracone hasło odzyskuje się komendą auth.recover")
	} else if !errors.Is(err, dane.ErrBrakWiersza) {
		return shared.AuthRegisterResponse{}, err
	}

	teraz := time.Now().UnixMilli()
	if err := a.konto.ZalozKonto(ctx, dane.KontoWlasciciela{
		Login: login, Email: email, Potwierdzone: false, Utworzono: teraz,
	}); err != nil {
		if errors.Is(err, dane.ErrKolizjaWiersza) {
			return shared.AuthRegisterResponse{}, bladBramki(shared.ErrorCodeConflict,
				"konto właściciela jest już założone; rejestracja wykonuje się raz")
		}
		return shared.AuthRegisterResponse{}, err
	}
	kotwica, err := a.zalozMetode(ctx, dane.MetodaUwierzytelnienia{
		Rodzaj:          shared.AuthMethodKindPassword,
		Etykieta:        wskaznikTekstu("Hasło konta"),
		NazwaUrzadzenia: niepustyTekst(z.DeviceName),
		Kotwica:         true,
		Utworzono:       teraz,
	}, z.Password)
	if err != nil {
		a.cofnijRejestracje(ctx, "")
		return shared.AuthRegisterResponse{}, err
	}
	// Bez poczty rejestracja KOŃCZY SIĘ TUTAJ i kończy się prawdą: konto jest,
	// hasło otwiera bramkę, a adres pozostaje niepotwierdzony, bo listu nie było
	// czym nadać. Znacznik zapamiętuje to na trwałe — inaczej ustawienie poczty
	// w oknie Konfiguracji zamykałoby bramkę przy następnym wejściu, przed
	// listem, którego żadna komenda nie wysyła drugi raz.
	//
	// `pendingVerification` idzie fałszem, bo kontrakt wiąże prawdę z faktem
	// nadania listu („Prawda oznacza, ze list z droga potwierdzenia zostal
	// wyslany"), a tu żaden list nie wyszedł. Pola na trzeci stan — „konto jest,
	// adres niepotwierdzony, listu nie wysłano" — kontrakt nie ma.
	if !pocztaJest {
		if err := a.zapiszZnacznikBezPoczty(ctx, email); err != nil {
			a.cofnijRejestracje(ctx, kotwica.Kod)
			return shared.AuthRegisterResponse{}, err
		}
		return shared.AuthRegisterResponse{Registered: true, PendingVerification: false}, nil
	}
	// Nadanie idzie ostatnie, a jego niepowodzenie NIE COFA rejestracji — schodzi
	// na drogę bez poczty.
	//
	// Sprawdzenie nastaw wyżej mówi tylko tyle, że konto nadawcze jest wskazane —
	// nie, że serwer odpowiada. Cofnięcie rejestracji było tu wcześniej ratunkiem
	// przed platformą NIE DO OTWARCIA (rejestracja wykonuje się raz, a wejść nie
	// było czym, bo bramkę zamykał brak potwierdzenia). Odkąd konto bez
	// potwierdzonego adresu wchodzi hasłem, ratunek jest zbędny, a sam był
	// pułapką: literówka w `DANACO_NADAWCA_HOST` po stronie powłoki zamykała
	// pierwsze uruchomienie równie szczelnie jak brak poczty w ogóle. Zostaje więc
	// konto, hasło otwiera bramkę, a adres czeka na potwierdzenie — dokładnie ten
	// sam stan, co przy instalce bez poczty.
	if err := a.wyslijDrogePotwierdzenia(ctx, dane.CelWeryfikacja, email, login); err != nil {
		if err := a.zapiszZnacznikBezPoczty(ctx, email); err != nil {
			a.cofnijRejestracje(ctx, kotwica.Kod)
			return shared.AuthRegisterResponse{}, err
		}
		return shared.AuthRegisterResponse{Registered: true, PendingVerification: false}, nil
	}
	return shared.AuthRegisterResponse{Registered: true, PendingVerification: true}, nil
}

// cofnijRejestracje zdejmuje to, co rejestracja zdążyła zapisać.
//
// Niepowodzenie samego cofnięcia nie ma komu wrócić — czynność już odmawia
// z powodu pierwotnego, a druga odmowa przykryłaby ten powód. Idzie więc do
// dziennika, jeżeli jest gdzie, i nie zmienia odpowiedzi.
//
// Drogi potwierdzenia nie kasujemy: leży jako sam skrót, wygasa po godzinie,
// a bez konta nie ma czego otworzyć.
func (a *adapterUwierzytelnienia) cofnijRejestracje(ctx context.Context, kodKotwicy string) {
	if kodKotwicy != "" {
		_, _ = a.repozytorium.UsunMetode(ctx, kodKotwicy)
		usunPoswiadczenie(ctx, a.sejf, przedrostekBytuSejfu+kodKotwicy)
	}
	if a.konto != nil {
		_ = a.konto.UsunKonto(ctx)
	}
	// Znacznik bramki bez poczty odchodzi razem z kontem: został po rejestracji,
	// której już nie ma, a przeżyłby ją w sejfie i mówił o koncie nieistniejącym.
	a.zdejmijZnacznikBezPoczty(ctx)
}

// ── auth.login ───────────────────────────────────────────────────────────────

// WejdzPrzezBramke otwiera bramkę hasłem albo PIN-em i zakłada sesję.
//
// Konto niepotwierdzone bramki nie otwiera. Bez tego sprawdzenia weryfikacja
// adresu byłaby ozdobą: rejestracja zakłada kotwicę, więc Operator wchodziłby
// hasłem zaraz po niej, nie zaglądając do skrzynki — a literówka w adresie
// wyszłaby na jaw dopiero w dniu, w którym trzeba nim odzyskać konto, czyli za
// późno. Adres jest jedyną drogą odzyskania i musi być sprawdzony ZANIM stanie
// się jedyną drogą.
func (a *adapterUwierzytelnienia) WejdzPrzezBramke(ctx context.Context,
	z shared.AuthLoginRequest) (shared.AuthLoginResponse, error) {

	if err := a.gotowa(); err != nil {
		return shared.AuthLoginResponse{}, err
	}
	if err := a.kontoPotwierdzone(ctx); err != nil {
		return shared.AuthLoginResponse{}, err
	}
	// Zwłoka idzie na wejściu czynności, przed sprawdzeniem sekretu. Nałożona
	// dopiero po rozpoznaniu sekretu jako błędnego, dawałaby wołającemu czas
	// odpowiedzi jako podpowiedź „ten sekret był dobry". Tu czeka każda próba
	// jednakowo i każda przechodzi dalej, bo dławik nie odmawia.
	droga := drogaWejscia(string(z.Method), wartoscTekstu(z.DeviceId))
	a.dlawik.Zaczekaj(ctx, droga)

	metoda, err := a.metodaWejscia(ctx, z)
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
		// Próba nieudana podnosi zwłokę próby następnej; ta odpowiada od razu,
		// bo czekanie ma spowalniać zgadywanie, a nie przetrzymywać wołającego po
		// tym, jak wynik jest już znany.
		a.dlawik.Niepowodzenie(droga)
		// Sekret niezgodny to nieudane wejście, nie wadliwe żądanie, więc kod jest
		// `not_authenticated`. `validation_failed` zostaje przy brakach formularza
		// (wejście bez sekretu, PIN bez urządzenia) — dzięki temu sonda stanu
		// bramki po stronie klienta, która wysyła żądanie bez sekretu, dalej czyta
		// `validation_failed` jako „bramka ustawiona", a Operator dostaje odmowę
		// odróżnialną od literówki w kształcie żądania.
		return shared.AuthLoginResponse{}, bladBramki(shared.ErrorCodeNotAuthenticated,
			"sekret metody "+string(z.Method)+" nie zgadza się z zapisem bramki")
	}
	// Wejście udane zeruje licznik — następne pomyłki liczą się od nowa, a
	// Operator pracujący normalnie nie czeka nigdy.
	a.dlawik.Wyzeruj(droga)
	if err := a.repozytorium.OdnotujUzycie(ctx, metoda.Kod, time.Now().UnixMilli()); err != nil {
		return shared.AuthLoginResponse{}, err
	}
	// Urządzenie sesji bierze się z ŻĄDANIA, a wiersz metody uzupełnia je tylko
	// wtedy, gdy żądanie milczy.
	//
	// Kotwica hasła urządzenia nie ma i mieć nie może — hasło nie jest materiałem
	// jednej maszyny. Branie urządzenia wyłącznie z wiersza metody porzucało więc
	// `deviceId` przy każdym wejściu hasłem, a maszyna nie pojawiała się w wykazie
	// `device.list`. Operator nie widział, co ma dostęp do jego konta, a
	// `device.revoke` z identyfikatorem tej maszyny wracał `revoked: false`
	// i zostawiał token czynny — w oknie, które istnieje po to, żeby dostęp
	// odbierać.
	urzadzenie := niepustyTekst(z.DeviceId)
	if urzadzenie == nil {
		urzadzenie = metoda.UrzadzenieKod
	}
	sesja, err := a.zalozSesje(ctx, metoda.Rodzaj, urzadzenie, wartoscPrawdy(z.KeepSignedIn))
	if err != nil {
		return shared.AuthLoginResponse{}, err
	}
	metody, err := a.MetodyWejscia(ctx)
	if err != nil {
		return shared.AuthLoginResponse{}, err
	}
	return shared.AuthLoginResponse{Session: sesja, Methods: metody}, nil
}

// metodaWejscia odnajduje metodę, którą żądanie chce otworzyć bramkę.
func (a *adapterUwierzytelnienia) metodaWejscia(ctx context.Context,
	z shared.AuthLoginRequest) (dane.MetodaUwierzytelnienia, error) {

	switch z.Method {
	case shared.AuthMethodKindPassword:
		metoda, err := a.kotwica(ctx)
		if errors.Is(err, dane.ErrBrakWiersza) {
			return metoda, bladBramki(shared.ErrorCodeNotFound,
				"hasło bramki nie istnieje — bramki jeszcze nie ustawiono (auth.register)")
		}
		return metoda, err
	case shared.AuthMethodKindPin:
		urzadzenie := wartoscTekstu(z.DeviceId)
		if urzadzenie == "" {
			return dane.MetodaUwierzytelnienia{}, bladBramki(shared.ErrorCodeValidationFailed,
				"wejście PIN-em bez wskazania urządzenia; PIN jest właściwy urządzeniu")
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
			"repozytorium bramki niewpięte")
	}
	if a.sejf == nil {
		return bladBramki(shared.ErrorCodeInternalError,
			"sejfu poświadczeń nie wpięto; hasła nie ma gdzie odłożyć ani z czym porównać")
	}
	return nil
}

// kotwica zwraca hasło bramki. Wynik `dane.ErrBrakWiersza` przechodzi bez
// przekładu — wołający rozstrzyga, czy brak kotwicy jest tu odmową, czy zgodą.
func (a *adapterUwierzytelnienia) kotwica(ctx context.Context) (dane.MetodaUwierzytelnienia, error) {
	return a.repozytorium.Kotwica(ctx)
}

// zalozMetode kładzie sekret w sejfie i dopiero potem wstawia wiersz metody.
// Kolejność jest z zamysłu: wiersz bez odwołania byłby metodą, którą nie da się
// wejść. Gdy wstawienie wiersza nie wyjdzie, wpis sejfu jest sprzątany —
// inaczej zostałby sekret bez właściciela.
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

// kolizjaMetody przekłada odmowę bazy na odmowę kontraktu.
//
// Sprawdzenia w `wolnoZalozyc` i `ZalozBramke` biegną przed wstawieniem wiersza
// i pod zamkiem, więc do naruszenia indeksu dochodzi już tylko wtedy, gdy stan
// zmienił ktoś spoza tego procesu. Wołający dostaje wtedy to samo zdanie, co
// przy odmowie sprawdzonej wcześniej: `conflict` i powód po ludzku. Przepuszczony
// błąd SQLite mówiłby co innego — `internal_error` z `retryable: true`, czyli
// „spróbuj jeszcze raz" tam, gdzie ponowienie nie wyjdzie, a w komunikacie szłyby
// nazwy tabel i kolumn bazy.
func kolizjaMetody(metoda dane.MetodaUwierzytelnienia, err error) error {
	if !errors.Is(err, dane.ErrKolizjaWiersza) {
		return err
	}
	if metoda.Kotwica {
		return bladBramki(shared.ErrorCodeConflict,
			"sekret bramki jest już ustawiony; hasła nie ustawia się drugi raz — "+
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
// Brak wpisu w sejfie nie jest odmową wejścia, tylko awarią rdzenia: wiersz
// metody istnieje, więc sekret gdzieś być powinien. Odpowiedź „nie zgadza się"
// przemilczałaby zgubiony sejf.
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
	rodzaj string, urzadzenie *string, niewylogowuj bool) (shared.AuthSession, error) {

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
		// Trwanie idzie do wiersza, bo przełącznik „nie wyloguj mnie" ma
		// obowiązywać także po odnowieniu, a sama chwila wygaśnięcia go nie
		// niesie (`migracja_112_trwanie_sesji_bramki.sql`).
		Trwanie: trwanie.Milliseconds(),
	})
	if err != nil {
		return shared.AuthSession{}, err
	}
	return sesjaBramkiKontraktu(token, zapisana), nil
}

// trwanieWejscia rozstrzyga długość sesji z przełącznika „nie wyloguj mnie".
//
// Przełącznik znosi powtarzanie logowania, a nie zakłada bramy: jedynym
// miejscem kontroli jest okno rejestracji i logowania, a ten przełącznik służy
// temu, żeby Operator widywał je jak najrzadziej.
func trwanieWejscia(niewylogowuj bool) time.Duration {
	if niewylogowuj {
		return trwanieSesjiBramkiDlugie
	}
	return trwanieSesjiBramki
}
