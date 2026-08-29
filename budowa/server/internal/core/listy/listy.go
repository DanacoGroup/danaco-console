// Odpowiedzialność pliku: obie postacie listów transakcyjnych — graficzna
// i tekstowa — składane z szablonów marki wraz ze znakami obu wariantów.
package listy

import (
	_ "embed"
	"html"
	"strconv"
	"strings"
	"time"
)

/*
Szablony pochodzą z `design/06-poczta-transakcyjna/` i są kopiowane do tego
katalogu bez zmian — plik źródłowy jest opracowaniem Właściciela, a nie miejscem
na poprawki wykonawcy. Zmienne opisuje `opracowanie-techniczne.md`.

Każdy list ma dwie części: `text/html` i `text/plain`. Obie stoją w szablonach,
żadnej nie składa kod — inaczej rozeszłyby się przy pierwszej poprawce jednej
z nich, a Operator z klientem bez grafiki czytałby co innego niż reszta.

Szablony i znaki idą w binarce, nie z dysku: rdzeń wysyła listy z maszyny
wdrożenia, na której katalogu `design/` nie ma.
*/

//go:embed list-01-aktywacja-konta.html
var szablonAktywacjiHtml string

//go:embed list-01-aktywacja-konta.txt
var szablonAktywacjiTekst string

//go:embed list-02-uwierzytelnienie-logowania.html
var szablonLogowaniaHtml string

//go:embed list-02-uwierzytelnienie-logowania.txt
var szablonLogowaniaTekst string

//go:embed list-03-autoresponder-brak-skrzynki.html
var szablonAutoresponderaHtml string

//go:embed list-03-autoresponder-brak-skrzynki.txt
var szablonAutoresponderaTekst string

/*
Znaki obu wariantów. Nagłówek listu niesie dwa znaczniki obrazu: jasny widoczny
domyślnie, ciemny odsłaniany zapytaniem medialnym (opracowanie, rozdz. 4.3).
Lockup jasny ma atrament `#181818` i na tle motywu ciemnego znika, więc jeden
znak nie wystarczy.
*/

//go:embed znak-marki.png
var ZnakJasny []byte

//go:embed znak-marki-ciemny.png
var ZnakCiemny []byte

// Odwołania, którymi szablon wskazuje znaki dołączone do listu. Opracowanie
// stanowi, że w wysyłce produkcyjnej znak idzie przez `cid:`, a `data:`
// występuje wyłącznie w pliku podglądu.
const (
	IdZnakuJasnego  = "danaco-lockup"
	IdZnakuCiemnego = "danaco-lockup-dark"
)

// AdresWsparcia to skrzynka obsługiwana, do której listy odsyłają Operatora.
// Skrzynka nadawcza odpowiedzi nie przyjmuje.
const AdresWsparcia = "support@danaco-group.pl"

/** Ile znaków kodu rozdziela spacja — opracowanie: „rozdzielany spacją co trzy znaki”. */
const grupaKodu = 3

// Postacie niesie obie części jednej wiadomości.
type Postacie struct {
	Html  string
	Tekst string
}

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

// ZlozAktywacje składa obie postacie listu aktywacji konta.
func ZlozAktywacje(a Aktywacja) Postacie {
	return zloz(szablonAktywacjiHtml, szablonAktywacjiTekst, a.Odbiorca, a.Zadano, []string{
		"{{code}}", rozdzielony(a.Kod),
		"{{expiry_minutes}}", strconv.Itoa(a.WaznoscMinuty),
		"{{requested_at}}", czas(a.Zadano),
		"{{activation_url}}", a.AdresAktywacji,
	})
}

// ZlozLogowanie składa obie postacie listu z kodem logowania.
func ZlozLogowanie(l Logowanie) Postacie {
	return zloz(szablonLogowaniaHtml, szablonLogowaniaTekst, l.Odbiorca, l.Zadano, []string{
		"{{code}}", rozdzielony(l.Kod),
		"{{expiry_minutes}}", strconv.Itoa(l.WaznoscMinuty),
		"{{requested_at}}", czas(l.Zadano),
		"{{ip_address}}", l.AdresZrodlowy,
		"{{device}}", l.Urzadzenie,
	})
}

// ZlozAutoresponder składa obie postacie odpowiedzi na wiadomość nadesłaną
// na skrzynkę bez odbioru.
func ZlozAutoresponder(a Autoresponder) Postacie {
	return zloz(szablonAutoresponderaHtml, szablonAutoresponderaTekst, a.Odbiorca, a.Odebrano, []string{
		"{{original_subject}}", oczyszczonyTemat(a.TematNadeslany),
		"{{received_at}}", czas(a.Odebrano),
	})
}

/*
zloz podstawia wartości w obie postacie listu.

Ucieczka znaczników obejmuje wyłącznie postać graficzną: w części tekstowej
`&amp;` czytałoby się dosłownie, a nie jako znak.
*/
func zloz(szablonHtml, szablonTekst, odbiorca string, chwila time.Time, wlasne []string) Postacie {
	wspolneHtml := []string{
		"{{logo_src}}", "cid:" + IdZnakuJasnego,
		"{{logo_src_dark}}", "cid:" + IdZnakuCiemnego,
		"{{recipient_address}}", html.EscapeString(odbiorca),
		"{{year}}", strconv.Itoa(chwila.Year()),
		"{{support_address}}", AdresWsparcia,
	}
	wspolneTekst := []string{
		"{{recipient_address}}", odbiorca,
		"{{year}}", strconv.Itoa(chwila.Year()),
		"{{support_address}}", AdresWsparcia,
	}
	return Postacie{
		Html:  strings.NewReplacer(append(zUcieczka(wlasne), wspolneHtml...)...).Replace(szablonHtml),
		Tekst: strings.NewReplacer(append(append([]string{}, wlasne...), wspolneTekst...)...).Replace(szablonTekst),
	}
}

// zUcieczka przepuszcza wartości przez ucieczkę znaczników na potrzeby postaci
// graficznej; nazwy zmiennych zostają bez zmian.
func zUcieczka(pary []string) []string {
	wynik := make([]string, len(pary))
	for i := 0; i+1 < len(pary); i += 2 {
		wynik[i] = pary[i]
		wynik[i+1] = html.EscapeString(pary[i+1])
	}
	return wynik
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
zmienną spoza rdzenia. Reguła jest nienegocjowalna: obcięcie do 120 znaków
i usunięcie znaków sterujących; ucieczkę znaczników dokłada złożenie postaci
graficznej.

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
	return bezSterujacych
}

// PodgladDanymiPrzykladowymi składa list aktywacji ze znakami wpisanymi
// w treść, żeby dało się go obejrzeć przeglądarką bez części powiązanych listu.
func PodgladDanymiPrzykladowymi(jasny, ciemny string) string {
	list := ZlozAktywacje(Aktywacja{
		Odbiorca:       "a.kowalska@przyklad.pl",
		Kod:            "418402",
		WaznoscMinuty:  60,
		Zadano:         time.Now(),
		AdresAktywacji: "https://console.danaco-group.pl/aktywacja",
	}).Html
	return strings.NewReplacer(
		"cid:"+IdZnakuJasnego, jasny,
		"cid:"+IdZnakuCiemnego, ciemny,
	).Replace(list)
}
