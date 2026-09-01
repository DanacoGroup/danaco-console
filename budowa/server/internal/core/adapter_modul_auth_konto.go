// Odpowiedzialność pliku: tożsamość właściciela i dwie drogi, które prowadzą
// przez jego skrzynkę — potwierdzenie adresu po rejestracji (auth.verify)
// oraz odzyskanie konta (auth.recover, auth.reset). Droga zawsze idzie
// listem na adres.
package core

import (
	"context"
	"errors"
	"fmt"
	netmail "net/mail"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/konfiguracja"
	"danacoconsole/server/internal/mail"
	"danacoconsole/server/internal/nadajnik"
	"danacoconsole/shared"

	"danacoconsole/server/internal/core/listy"
)

// trwanieDrogiPotwierdzenia — jak długo ważny jest materiał wysłany listem.
// Godzina jest kompromisem między skrzynką sprawdzaną rzadko a materiałem
// ważnym bez końca; droga wygasła nie zamyka niczego trwale, Operator prosi
// o nową.
const trwanieDrogiPotwierdzenia = time.Hour

// pulapProbDrogi — po tylu pomyłkach droga zostaje zamknięta. Sześciocyfrowy kod
// ma milion wartości, więc dziesięć prób na jedną wydaną drogę zostawia
// zgadującemu szansę jedną na sto tysięcy, a Operatorowi zapas na pomyłkę przy
// przepisywaniu kodu z ekranu.
const pulapProbDrogi = 10

// odstepMiedzyListami — najkrótszy odstęp między dwoma listami z drogą do tego
// samego konta i celu. Kod wysłany chwilę temu jest jeszcze ważny; odstęp
// krótszy mnożyłby drogi czynne i dawałby pytającemu tyle listów na cudzy adres,
// ile zdąży poprosić.
const odstepMiedzyListami = time.Minute

// kontoGotowe odmawia, gdy montaż nie wpiął trwałości tożsamości.
//
// Odmowa jest wprost, nie cicha: rejestracja bez zapisu konta wyglądałaby jak
// udana, a Operator zostałby przed platformą, która o nim nie wie.
func (a *adapterUwierzytelnienia) kontoGotowe() error {
	if a.konto == nil {
		return bladBramki(shared.ErrorCodeInternalError,
			"Nie można zapisać konta — magazyn danych jest niedostępny.")
	}
	return nil
}

// kontoPotwierdzone zamyka bramkę przed kontem, którego adresu nikt nie
// potwierdził. Brak trwałości konta i brak wiersza konta nie zamykają bramki,
// bo o potwierdzeniu nie ma wtedy co rozstrzygać.
func (a *adapterUwierzytelnienia) kontoPotwierdzone(ctx context.Context,
	konto dane.KontoWlasciciela) error {

	// Konto puste znaczy platformę bez trwałości konta albo przed rejestracją —
	// o potwierdzeniu nie ma wtedy co rozstrzygać.
	if a.konto == nil || konto.Id == 0 {
		return nil
	}
	if konto.Potwierdzone {
		return nil
	}
	// Bramki nie zamyka potwierdzenie, gdy konto założono bez poczty — wejście idzie wtedy hasłem.
	if _, bezPoczty := a.znacznikBezPoczty(ctx); bezPoczty {
		return nil
	}
	return bladBramkiZPowodem(shared.ErrorCodeNotAuthenticated, PowodAdresNiepotwierdzony,
		"Konto oczekuje na potwierdzenie adresu "+konto.Email+
			". Wprowadź kod potwierdzający z wiadomości — do tego czasu wejście "+
			"jest zamknięte, bo adresem odzyskuje się konto po utracie hasła.")
}

// daneRejestracji sprawdza login i adres podane przy rejestracji. Sprawdzenie
// adresu jest celowo płytkie, bo głębsza kontrola postaci adresu odrzuca
// adresy poprawne i tak nie rozstrzyga, czy skrzynka istnieje.
func daneRejestracji(z shared.AuthRegisterRequest) (string, string, error) {
	login := strings.TrimSpace(z.Login)
	if login == "" {
		return "", "", bladBramki(shared.ErrorCodeValidationFailed,
			"Podaj login — to nazwa, pod którą będziesz się logować.")
	}
	email := strings.TrimSpace(z.Email)
	malpa := strings.LastIndex(email, "@")
	if malpa <= 0 || !strings.Contains(email[malpa+1:], ".") {
		return "", "", bladBramki(shared.ErrorCodeValidationFailed,
			"adres e-mail uwierzytelniający nie ma postaci adresu — "+
				"jest jedyną drogą odzyskania konta, więc musi być adresem działającym")
	}
	return login, email, nil
}

// wyslijDrogePotwierdzenia losuje materiał, zapisuje jego skrót i nadaje list.
// Kolejność jest wiążąca: najpierw zapis skrótu, potem nadanie, inaczej list
// mógłby dojść, zanim droga stałaby się ważna.
func (a *adapterUwierzytelnienia) wyslijDrogePotwierdzenia(ctx context.Context,
	cel, email string, kontoId int64) error {

	// Konto nadawcze czytane raz: to samo, którym list wyjdzie, nazywa nadawcę w nagłówku wiadomości.
	nadawca := a.kontoNadawcze(ctx)
	// Brak konta nadawczego nazywa się przed zapisaniem drogi, żeby baza nie trzymała drogi bez listu.
	if err := nadawca.Brak(); err != nil {
		return bladBramki(shared.ErrorCodeInternalError,
			err.Error()+"; "+dwieDrogiKontaNadawczego)
	}
	/* Kod, nie token sesji: okno przyjmuje go w sześciu polach, a Operator
	   przepisuje go z wiadomości ręcznie. */
	droga, err := nowyKodPotwierdzenia()
	if err != nil {
		return err
	}
	teraz := time.Now()
	/* Droga nowa zamyka poprzednie drogi tego konta i celu. Bez tego kody się
	   kumulują: każdy list zostawia w obiegu drogę ważną godzinę, więc zgadujący
	   ma tyle celów naraz, ile razy Operator poprosił o list. */
	if _, err := a.repozytorium.ZamknijDrogiKonta(ctx, kontoId, cel); err != nil {
		return err
	}
	if err := a.konto.ZalozPotwierdzenie(ctx, dane.PotwierdzenieTozsamosci{
		Skrot:     skrotTokenu(bezOdstepow(droga)),
		Cel:       cel,
		Wygasa:    teraz.Add(trwanieDrogiPotwierdzenia).UnixMilli(),
		Utworzono: teraz.UnixMilli(),
		// Konto, do którego droga prowadzi; przy dwóch kontach naraz to jedyne,
		// co rozstrzyga, które z nich potwierdza przepisany kod.
		KontoId: kontoId,
	}); err != nil {
		return err
	}
	wiadomosc, err := listPotwierdzenia(nadawca, a.adresKonsoli, cel, droga, email)
	if err != nil {
		return bladBramki(shared.ErrorCodeInternalError,
			fmt.Sprintf("nie udało się złożyć listu na %s: %v; %s", email, err, dwieDrogiKontaNadawczego))
	}
	if _, err := nadajnik.Wyslij(nadawca, wiadomosc); err != nil {
		return bladBramki(shared.ErrorCodeInternalError,
			fmt.Sprintf("nie udało się wysłać listu na %s: %v; %s", email, err, dwieDrogiKontaNadawczego))
	}
	return nil
}

// adresWsparcia to skrzynka obsługiwana, do której listy odsyłają Operatora.
// Skrzynka nadawcza odpowiedzi nie przyjmuje.
const adresWsparcia = "support@danaco-group.pl"

// grupaKodu — ile znaków kodu rozdziela spacja; opracowanie poczty:
// „rozdzielany spacją co trzy znaki”.
const grupaKodu = 3

// postacCzasuListu zapisuje chwilę tak, jak żąda opracowanie (rozdz. 6.2):
// `29.08.2026, 17:22 CEST`.
const postacCzasuListu = "02.01.2006, 15:04 MST"

// kompletListow wczytuje szablony raz na proces. Wczytanie przechodzi siedem
// par plików i zdejmuje z nich komentarze wewnętrzne — powtarzanie tego przy
// każdym liście byłoby pracą bez skutku.
var kompletListow = sync.OnceValues(mail.WbudowanyKomplet)

// listPotwierdzenia składa jeden z dwóch listów systemowych w gotowe bajty
// wiadomości. Nagłówki, kodowanie tematu i części oraz próg przycięcia
// rozstrzyga pakiet mail; rdzeń podaje wyłącznie wartości zmiennych.
func listPotwierdzenia(nadawca nadajnik.Nastawy, adresKonsoli, cel, droga, email string) (*mail.Message, error) {
	komplet, err := kompletListow()
	if err != nil {
		return nil, err
	}
	teraz := time.Now()
	wartosci := map[string]string{
		"recipient_address": email,
		"support_address":   adresWsparcia,
		"year":              strconv.Itoa(teraz.Year()),
		"code":              rozdzielony(droga),
		"expiry_minutes":    strconv.Itoa(int(trwanieDrogiPotwierdzenia.Minutes())),
		"requested_at":      teraz.Format(postacCzasuListu),
	}
	rodzaj := mail.KindAccountActivation
	if cel == dane.CelOdzyskanie {
		/* Odzyskanie dostaje list resetu hasła, nie kodu logowania: to on nazywa
		   czynność, którą Operator właśnie prowadzi, i tylko on mówi, że hasło
		   dotychczasowe zostaje ważne do ustawienia nowego. */
		rodzaj = mail.KindPasswordReset
	} else {
		// Odsyłacz do okna, w którym Operator wprowadza kod; zna go tylko ta gałąź.
		wartosci["activation_url"] = konfiguracja.AdresAktywacji(adresKonsoli)
	}

	return komplet.Build(rodzaj, mail.Envelope{
		From: nadawca.Nadawca(),
		To:   netmail.Address{Address: email},
		Date: teraz,
	}, mail.Content{Values: wartosci}, znakiMarki())
}

// znakiMarki podaje oba warianty znaku dołączane częścią listu. Nazwy plików
// widzi Operator wtedy, gdy jego klient pocztowy potraktuje znak jak załącznik.
func znakiMarki() mail.Logos {
	return mail.Logos{
		Light: listy.ZnakJasny, LightName: "danaco.png",
		Dark: listy.ZnakCiemny, DarkName: "danaco-ciemny.png",
	}
}

// rozdzielony wstawia spację co trzy znaki kodu — czyta się go wtedy z ekranu
// bez gubienia miejsca, a przepisuje bez pomyłki.
func rozdzielony(kod string) string {
	var wynik strings.Builder
	for i, znak := range kod {
		if i > 0 && i%grupaKodu == 0 {
			wynik.WriteByte(' ')
		}
		wynik.WriteRune(znak)
	}
	return wynik.String()
}

/*
bezOdstepow zdejmuje z kodu wszystkie odstępy, nie tylko brzegowe.

List pokazuje kod rozdzielony spacją co trzy znaki — `418 402` — żeby dało się
go przeczytać z ekranu bez gubienia miejsca. Operator, który go stamtąd skopiuje
i wklei, poda właśnie taką postać. Odstęp jest sposobem zapisu, nie częścią kodu.

Sito dotyczy wyłącznie kodu potwierdzenia. Token sesji przechodzi przez ten sam
`skrotTokenu`, ale jest nieprzezroczysty i odstępu usuwać w nim nie wolno.
*/
func bezOdstepow(kod string) string {
	return strings.Map(func(znak rune) rune {
		if unicode.IsSpace(znak) {
			return -1
		}
		return znak
	}, kod)
}

// drogaKonta rozpoznaje drogę kodem i sprawdza jej cel, nie zamykając jej.
// Rozpoznanie stoi osobno od zamknięcia, bo wywołujący musi znać konto, zanim
// zdecyduje, czy drogę zużyć — kod odrzuconego żądania ma zostać ważny.
func (a *adapterUwierzytelnienia) drogaKonta(ctx context.Context,
	cel, droga string) (dane.PotwierdzenieTozsamosci, error) {

	droga = bezOdstepow(droga)
	if droga == "" {
		return dane.PotwierdzenieTozsamosci{},
			bladBramki(shared.ErrorCodeValidationFailed, "Podaj kod potwierdzający.")
	}
	zapis, err := a.konto.PotwierdzeniePoSkrocie(ctx, skrotTokenu(droga))
	if errors.Is(err, dane.ErrBrakWiersza) {
		return dane.PotwierdzenieTozsamosci{}, bladBramki(shared.ErrorCodeNotAuthenticated,
			"Ten kod potwierdzający nie jest znany.")
	}
	if err != nil {
		return dane.PotwierdzenieTozsamosci{}, err
	}
	if zapis.Cel != cel {
		return dane.PotwierdzenieTozsamosci{}, bladBramki(shared.ErrorCodeNotAuthenticated,
			"Ten kod potwierdzający dotyczy innej czynności.")
	}
	return zapis, nil
}

// zamknijDroge zużywa drogę rozpoznaną wcześniej przez `drogaKonta`.
func (a *adapterUwierzytelnienia) zamknijDroge(ctx context.Context,
	zapis dane.PotwierdzenieTozsamosci) error {

	zamknieta, err := a.konto.ZuzyjPotwierdzenie(ctx, zapis.Skrot, time.Now().UnixMilli())
	if err != nil {
		return err
	}
	if !zamknieta {
		return bladBramki(shared.ErrorCodeNotAuthenticated,
			"Ten kod potwierdzający został już użyty lub wygasł. Poproś o nowy.")
	}
	return nil
}

// kontoDrogi oddaje konto, do którego droga prowadzi. Droga zastana, bez
// wskazania konta, pochodzi sprzed migracji 406 i prowadzi do konta
// najstarszego — instalacja z jednym kontem rozstrzyga jednoznacznie.
func (a *adapterUwierzytelnienia) kontoDrogi(ctx context.Context,
	zapis dane.PotwierdzenieTozsamosci) (int64, error) {

	if zapis.KontoId != 0 {
		return zapis.KontoId, nil
	}
	konto, err := a.konto.Konto(ctx)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return 0, bladBramki(shared.ErrorCodeConflict,
			"Konto Operatora nie zostało jeszcze założone. Zarejestruj się.")
	}
	if err != nil {
		return 0, err
	}
	return konto.Id, nil
}

/*
odnotujPomylkeDrogi zapisuje kod, który nie trafił w żadną drogę: podnosi zwłokę
należną próbie następnej i dolicza pomyłkę drogom czynnym tego celu, zamykając
te, które doszły do pułapu.

Pomyłka liczy się drogom, a nie jednej z nich, bo kod nietrafiony żadnej nie
nazywa. Niepowodzenie zapisu idzie do dziennika, nie do odpowiedzi — czynność
i tak kończy się odmową, a wołający ma dostać jedno zdanie o kodzie, nie o bazie.
*/
func (a *adapterUwierzytelnienia) odnotujPomylkeDrogi(ctx context.Context,
	cel string, proba zajecieDrogi) {

	proba.Niepowodzenie()
	zamkniete, err := a.repozytorium.OdnotujPomylkeDrogi(ctx, cel, pulapProbDrogi,
		time.Now().UnixMilli())
	dziennik := dziennikZKontekstu(ctx)
	if dziennik == nil {
		return
	}
	if err != nil {
		dziennik.Printf("bramka: nie można doliczyć pomyłki drogom celu %s: %v", cel, err)
		return
	}
	dziennik.Printf("bramka: kod celu %s nie trafił w żadną drogę; dróg zamkniętych po %d próbach: %d",
		cel, pulapProbDrogi, zamkniete)
}

// ── auth.verify ──────────────────────────────────────────────────────────────

// PotwierdzAdres zamyka rejestrację: przenosi konto do stanu potwierdzonego
// i wydaje urządzeniu token dostępu.
//
// To jest moment pierwszego wejścia Operatora do platformy — dlatego sesja
// powstaje tutaj, a nie przy `auth.register`.
func (a *adapterUwierzytelnienia) PotwierdzAdres(ctx context.Context,
	z shared.AuthVerifyRequest) (shared.AuthVerifyResponse, error) {

	if err := a.kontoGotowe(); err != nil {
		return shared.AuthVerifyResponse{}, err
	}
	a.zamekZmiany.Lock()
	defer a.zamekZmiany.Unlock()

	// Zwłoka nakłada się przed sprawdzeniem kodu, żeby czas odpowiedzi nie
	// zdradzał wyniku, i trzyma drogę zajętą do rozstrzygnięcia próby.
	proba := a.dlawik.Podejdz(ctx, kluczDlawika(ctx, czynnoscDrogi, 0))
	defer proba.Zwolnij()

	/* Konto rozstrzyga droga, nie kolejność założenia: przy dwóch rejestracjach
	   naraz kod z listu prowadzi do konta, na które ten list poszedł. */
	zapis, err := a.drogaKonta(ctx, dane.CelWeryfikacja, z.Token)
	if err != nil {
		a.odnotujPomylkeDrogi(ctx, dane.CelWeryfikacja, proba)
		return shared.AuthVerifyResponse{}, err
	}
	proba.Wyzeruj()

	kontoId, err := a.kontoDrogi(ctx, zapis)
	if err != nil {
		return shared.AuthVerifyResponse{}, err
	}
	konto, err := a.konto.KontoPoId(ctx, kontoId)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return shared.AuthVerifyResponse{}, bladBramki(shared.ErrorCodeConflict,
			"Konto Operatora nie zostało jeszcze założone. Zarejestruj się.")
	}
	if err != nil {
		return shared.AuthVerifyResponse{}, err
	}
	if konto.Potwierdzone {
		return shared.AuthVerifyResponse{}, bladBramki(shared.ErrorCodeConflict,
			"Adres jest już potwierdzony. Zaloguj się.")
	}
	if err := a.zamknijDroge(ctx, zapis); err != nil {
		return shared.AuthVerifyResponse{}, err
	}
	if err := a.konto.PotwierdzKonto(ctx, konto.Id); err != nil {
		return shared.AuthVerifyResponse{}, err
	}
	// Znacznik bramki bez poczty przestał być prawdą — bramkę trzyma odtąd sam wiersz konta.
	a.zdejmijZnacznikBezPoczty(ctx)
	sesja, err := a.zalozSesje(ctx, shared.AuthMethodKindPassword,
		niepustyTekst(z.DeviceId), wartoscPrawdy(z.KeepSignedIn), konto.Id)
	if err != nil {
		return shared.AuthVerifyResponse{}, err
	}
	return shared.AuthVerifyResponse{Verified: true, Session: sesja}, nil
}

// ── auth.recover ─────────────────────────────────────────────────────────────

// RozpocznijOdzyskanie wysyła drogę potwierdzenia na adres uwierzytelniający.
// Odpowiedź jest zawsze taka sama, niezależnie od tego, czy adres pasuje do
// konta, dla właściciela czynność jest wykonana i list wychodzi.
func (a *adapterUwierzytelnienia) RozpocznijOdzyskanie(ctx context.Context,
	z shared.AuthRecoverRequest) (shared.AuthRecoverResponse, error) {

	if err := a.kontoGotowe(); err != nil {
		return shared.AuthRecoverResponse{}, err
	}
	email := strings.TrimSpace(z.Email)
	if email == "" {
		return shared.AuthRecoverResponse{}, bladBramki(shared.ErrorCodeValidationFailed,
			"Podaj adres e-mail konta.")
	}
	/* Konto wskazuje podany adres, nie kolejność założenia. Odczyt „konta
	   najstarszego” pochodził z czasu, gdy konto było jedno; przy wielu kontach
	   odsyłał każdy adres poza pierwszym z odpowiedzią „wysłano” i nie wysyłał
	   nic — Operator czekał na list, który nigdy nie powstał. */
	konto, err := a.konto.KontoPoTozsamosci(ctx, email)
	if errors.Is(err, dane.ErrBrakWiersza) {
		/* Adres nieznany odpowiada tak samo jak znany: inaczej pytanie o kolejne
		   adresy wskazywałoby, które z nich mają konto. */
		return shared.AuthRecoverResponse{Sent: true}, nil
	}
	if err != nil {
		return shared.AuthRecoverResponse{}, err
	}
	if !strings.EqualFold(konto.Email, email) {
		// Wskazanie trafiło w login, nie w adres — droga odzyskania idzie adresem.
		return shared.AuthRecoverResponse{Sent: true}, nil
	}
	/*
		Konto niepotwierdzone dostaje kod aktywacji, nie kod odzyskania.

		Odzyskać można dostęp do konta, które kiedyś działało; konto bez
		potwierdzonego adresu nigdy nie zostało otwarte, a jego jedyną przeszkodą
		jest brak aktywacji. Kod odzyskania jej nie zdejmuje — droga potwierdzenia
		rozróżnia cele i kodu jednego celu nie przyjmuje w drugim. Bez tego Operator
		stał przed wejściem zamkniętym, a jedyna droga, którą okno mu podawało,
		wydawała kod nieprzydatny do niczego.
	*/
	cel := dane.CelOdzyskanie
	if !konto.Potwierdzone {
		cel = dane.CelWeryfikacja
	}
	/* Odstęp między listami. Droga wydana chwilę temu leży w skrzynce i jest
	   ważna, a każdy list ponad nią to kolejna droga do zgadywania oraz kolejna
	   wiadomość dla właściciela adresu, który o nią nie prosił. Odpowiedź
	   zostaje ta sama, bo różnica mówiłaby pytającemu, że o ten adres ktoś
	   właśnie prosił. */
	ostatnia, err := a.repozytorium.OstatniaDrogaKonta(ctx, konto.Id, cel)
	if err != nil {
		return shared.AuthRecoverResponse{}, err
	}
	if ostatnia != 0 && time.Since(time.UnixMilli(ostatnia)) < odstepMiedzyListami {
		return shared.AuthRecoverResponse{Sent: true}, nil
	}
	if err := a.wyslijDrogePotwierdzenia(ctx, cel, konto.Email, konto.Id); err != nil {
		return shared.AuthRecoverResponse{}, err
	}
	return shared.AuthRecoverResponse{Sent: true}, nil
}

// ── auth.reset ───────────────────────────────────────────────────────────────

// UstawNoweHaslo zamyka odzyskanie konta: podmienia hasło i unieważnia tokeny
// wydane przed zmianą. Unieważnienie jest częścią czynności, bo dostęp mógł
// mieć ktoś jeszcze.
func (a *adapterUwierzytelnienia) UstawNoweHaslo(ctx context.Context,
	z shared.AuthResetRequest) (shared.AuthResetResponse, error) {

	odpowiedz, _, err := a.ustawNoweHasloKonta(ctx, z)
	return odpowiedz, err
}

// ustawNoweHasloKonta wykonuje odzyskanie i oddaje obok odpowiedzi konto z drogi
// — adresata rozgłoszenia i zerwania gniazd, którego odpowiedź kontraktu nie
// niesie.
func (a *adapterUwierzytelnienia) ustawNoweHasloKonta(ctx context.Context,
	z shared.AuthResetRequest) (shared.AuthResetResponse, int64, error) {

	if err := a.gotowa(); err != nil {
		return shared.AuthResetResponse{}, 0, err
	}
	if err := a.kontoGotowe(); err != nil {
		return shared.AuthResetResponse{}, 0, err
	}
	if strings.TrimSpace(z.NewPassword) == "" {
		return shared.AuthResetResponse{}, 0, bladBramki(shared.ErrorCodeValidationFailed,
			"Podaj nowe hasło.")
	}
	a.zamekZmiany.Lock()
	defer a.zamekZmiany.Unlock()

	// Zwłoka nakłada się przed sprawdzeniem kodu, żeby czas odpowiedzi nie
	// zdradzał wyniku, i trzyma drogę zajętą do rozstrzygnięcia próby.
	proba := a.dlawik.Podejdz(ctx, kluczDlawika(ctx, czynnoscDrogi, 0))
	defer proba.Zwolnij()

	/* Konto wskazuje droga z listu, nie kolejność założenia. Kotwica czytana bez
	   wskazania konta jest kotwicą konta najstarszego: kod przysłany na własny
	   adres ustawiałby hasło konta założonego przy instalacji, czyli otwierał
	   cudzą bramkę bez znajomości czegokolwiek, co do niej należy. */
	zapis, err := a.drogaKonta(ctx, dane.CelOdzyskanie, z.Token)
	if err != nil {
		a.odnotujPomylkeDrogi(ctx, dane.CelOdzyskanie, proba)
		return shared.AuthResetResponse{}, 0, err
	}
	proba.Wyzeruj()

	kontoId, err := a.kontoDrogi(ctx, zapis)
	if err != nil {
		return shared.AuthResetResponse{}, 0, err
	}
	// Dalsze czynności idą kontem z drogi: żądanie przychodzi z urządzenia bez
	// sesji, więc kontekst połączenia konta nie zna.
	ctx = dane.ZKontemOperatora(ctx, kontoId)

	kotwica, err := a.kotwicaKonta(ctx, kontoId)
	if errors.Is(err, dane.ErrBrakWiersza) {
		return shared.AuthResetResponse{}, 0, bladBramki(shared.ErrorCodeConflict,
			"Konto Operatora nie zostało jeszcze założone. Zarejestruj się.")
	}
	if err != nil {
		return shared.AuthResetResponse{}, 0, err
	}
	if err := a.zamknijDroge(ctx, zapis); err != nil {
		return shared.AuthResetResponse{}, 0, err
	}

	sekret, err := zapisSekretu(z.NewPassword)
	if err != nil {
		return shared.AuthResetResponse{}, 0, err
	}
	odwolanie, err := a.sejf.Zapisz(ctx, przedrostekBytuSejfu+kotwica.Kod, sekret)
	if err != nil {
		return shared.AuthResetResponse{}, 0, err
	}
	if err := a.repozytorium.ZapiszOdwolanieSekretu(ctx, kotwica.Kod, odwolanie); err != nil {
		return shared.AuthResetResponse{}, 0, err
	}
	// Pusty skrót zachowany znaczy unieważnij wszystkie sesje tego konta —
	// odzyskanie idzie z urządzenia bez sesji, a konta obce zmiana nie dotyczy.
	uniewaznione, err := a.repozytorium.UniewaznijSesjeBramkiPoza(ctx, "", time.Now().UnixMilli())
	if err != nil {
		return shared.AuthResetResponse{}, 0, err
	}
	return shared.AuthResetResponse{Changed: true, RevokedDevices: uniewaznione}, kontoId, nil
}
