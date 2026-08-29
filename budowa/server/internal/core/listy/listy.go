// Odpowiedzialność pliku: postać graficzna listów transakcyjnych — szablony
// marki wraz ze znakiem, w które rdzeń wpisuje kod i okoliczności żądania.
package listy

import (
	_ "embed"
	"html"
	"strconv"
	"strings"
	"time"
)

/*
Szablony pochodzą z `design/06-poczta-transakcyjna/html/` i są kopiowane do tego
katalogu bez zmian — plik źródłowy jest opracowaniem Właściciela, a nie miejscem
na poprawki wykonawcy. Zmienne opisuje `opracowanie-techniczne.md`, rozdz. 4.

Szablony i znak idą w binarce, nie z dysku: rdzeń wysyła listy z maszyny
wdrożenia, na której katalogu `design/` nie ma. Plik odczytywany w chwili
wysyłki byłby drogą, na której brak jednego pliku zamienia list w pustą kartkę.
*/

//go:embed list-01-aktywacja-konta.html
var szablonAktywacji string

//go:embed list-02-uwierzytelnienie-logowania.html
var szablonLogowania string

//go:embed list-03-autoresponder-brak-skrzynki.html
var szablonAutorespondera string

// ZnakMarki niesie obraz znaku dołączany do listu jako część powiązana.
//
//go:embed znak-marki.png
var ZnakMarki []byte

/*
IdZnaku jest odwołaniem, którym szablon wskazuje znak dołączony do listu.
Opracowanie stanowi (rozdz. 4.1), że w wysyłce produkcyjnej znak idzie przez
`cid:`, a `data:` występuje wyłącznie w pliku podglądu — klienty pocztowe
odrzucają albo blokują obrazy wpisane w treść, a znak dołączony częścią listu
pokazują bez pytania.
*/
const IdZnaku = "danaco-lockup"

// AdresWsparcia to skrzynka obsługiwana, do której listy odsyłają Operatora.
// Skrzynka nadawcza odpowiedzi nie przyjmuje (opracowanie, rozdz. 5.2).
const AdresWsparcia = "support@danaco-group.pl"

/** Ile znaków kodu rozdziela spacja — opracowanie, rozdz. 4.2: „rozdzielany spacją co trzy znaki”. */
const grupaKodu = 3

// Aktywacja niesie dane listu zakładającego konto.
type Aktywacja struct {
	Odbiorca       string
	Kod            string
	WaznoscMinuty  int
	Zadano         time.Time
	AdresAktywacji string
}

// Logowanie niesie dane listu z kodem drugiego składnika wejścia.
type Logowanie struct {
	Odbiorca      string
	Kod           string
	WaznoscMinuty int
	Zadano        time.Time
	AdresZrodlowy string
	Urzadzenie    string
}

// Autoresponder niesie dane odpowiedzi na wiadomość przysłaną na skrzynkę,
// która odbioru nie prowadzi.
type Autoresponder struct {
	Odbiorca       string
	TematNadeslany string
	Odebrano       time.Time
}

// ZlozAktywacje składa list aktywacji konta.
func ZlozAktywacje(a Aktywacja) string {
	return strings.NewReplacer(
		wspolne(a.Odbiorca, a.Zadano)...,
	).Replace(strings.NewReplacer(
		"{{code}}", rozdzielony(a.Kod),
		"{{expiry_minutes}}", strconv.Itoa(a.WaznoscMinuty),
		"{{requested_at}}", czas(a.Zadano),
		"{{activation_url}}", html.EscapeString(a.AdresAktywacji),
	).Replace(szablonAktywacji))
}

// ZlozLogowanie składa list z kodem logowania.
func ZlozLogowanie(l Logowanie) string {
	return strings.NewReplacer(
		wspolne(l.Odbiorca, l.Zadano)...,
	).Replace(strings.NewReplacer(
		"{{code}}", rozdzielony(l.Kod),
		"{{expiry_minutes}}", strconv.Itoa(l.WaznoscMinuty),
		"{{requested_at}}", czas(l.Zadano),
		"{{ip_address}}", html.EscapeString(l.AdresZrodlowy),
		"{{device}}", html.EscapeString(l.Urzadzenie),
	).Replace(szablonLogowania))
}

// ZlozAutoresponder składa odpowiedź na wiadomość nadesłaną na skrzynkę bez odbioru.
func ZlozAutoresponder(a Autoresponder) string {
	return strings.NewReplacer(
		wspolne(a.Odbiorca, a.Odebrano)...,
	).Replace(strings.NewReplacer(
		"{{original_subject}}", oczyszczonyTemat(a.TematNadeslany),
		"{{received_at}}", czas(a.Odebrano),
	).Replace(szablonAutorespondera))
}

// wspolne oddaje pary podstawień wspólne wszystkim trzem listom (rozdz. 4.1).
func wspolne(odbiorca string, chwila time.Time) []string {
	return []string{
		"{{logo_src}}", "cid:" + IdZnaku,
		"{{recipient_address}}", html.EscapeString(odbiorca),
		"{{year}}", strconv.Itoa(chwila.Year()),
		"{{support_address}}", AdresWsparcia,
	}
}

// czas zapisuje chwilę w postaci z opracowania: `29.08.2026, 17:22 CEST`.
func czas(chwila time.Time) string {
	return chwila.Format("02.01.2006, 15:04 MST")
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
oczyszczonyTemat przygotowuje temat wiadomości nadesłanej z zewnątrz — jedyną
zmienną spoza rdzenia. Reguła jest nienegocjowalna (opracowanie, rozdz. 4.5):
ucieczka znaczników, obcięcie do 120 znaków i usunięcie znaków sterujących.

Bez tego temat nadesłany przez obcego staje się drogą wstrzyknięcia znaczników
do listu wychodzącego, a złamanie wiersza — drogą wstrzyknięcia nagłówka.
*/
func oczyszczonyTemat(temat string) string {
	bezSterujacych := strings.Map(func(znak rune) rune {
		if znak < 0x20 || znak == 0x7F {
			return -1
		}
		return znak
	}, temat)
	if len(bezSterujacych) > 120 {
		bezSterujacych = bezSterujacych[:120]
	}
	return html.EscapeString(bezSterujacych)
}

// PodgladDanymiPrzykladowymi składa list aktywacji ze znakiem wpisanym w treść,
// żeby dało się go obejrzeć przeglądarką bez części powiązanej listu.
func PodgladDanymiPrzykladowymi(znakDataURI string) string {
	list := ZlozAktywacje(Aktywacja{
		Odbiorca:       "a.kowalska@przyklad.pl",
		Kod:            "418402",
		WaznoscMinuty:  60,
		Zadano:         time.Now(),
		AdresAktywacji: "https://console.danaco-group.pl/aktywacja",
	})
	/* Wszystkie wystąpienia, nie pierwsze: szablon wskazuje znak w kilku
	   miejscach — w nagłówku listu i w ukrytym wierszu podglądu. */
	return strings.ReplaceAll(list, "cid:"+IdZnaku, znakDataURI)
}
