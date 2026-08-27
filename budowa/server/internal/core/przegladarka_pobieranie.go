// Plik niesie klienta pobierania stron modułu Browser: `browser.navigate`
// sięga po stronę HTTP GET-em i wydobywa z HTML-a tytuł oraz tekst
// renderowany, którymi wypełnia migawkę.
package core

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode"
)

const (
	// czasPobrania ogranicza całe pobranie — połączenie, nagłówki i czytanie
	// treści łącznie. Operator nie ma czekać w nieskończoność na stronę, która
	// milczy.
	czasPobrania = 15 * time.Second

	// limitTresci tnie odpowiedź na tej granicy bajtów. Rdzeń trzyma migawkę
	// strony, nie kopię całego internetu — strona większa zostaje ucięta.
	limitTresci = 2 << 20 // 2 MiB

	// limitPrzekierowan tnie łańcuch przekierowań. Bez tej granicy dwie strony
	// odsyłające do siebie nawzajem trzymałyby rdzeń aż do limitu czasu; pięć
	// skoków wystarcza każdej uczciwej witrynie (http→https, adres kanoniczny).
	limitPrzekierowan = 5

	// znacznikRdzenia przedstawia rdzeń serwerom witryn, z których część
	// odrzuca żądanie bez nagłówka User-Agent.
	znacznikRdzenia = "DanacoConsole/1.0 (+rdzen przegladarki)"
)

// errZaDuzaStrona nazywa przekroczenie granicy rozmiaru osobnym błędem, bo to
// jedyny przypadek, w którym rdzeń ma już bajty strony i mimo to odmawia —
// odmowa musi powiedzieć „strona jest za duża", a nie udać ucięty sukces.
var errZaDuzaStrona = errors.New("strona przekracza granicę rozmiaru")

// errStronaNieodpowiedziala nazywa niepowodzenie transportu: gospodarz nie
// przyjął połączenia, nazwa się nie rozwiązała, albo czas pobrania minął.
var errStronaNieodpowiedziala = errors.New("strona nie odpowiedziała")

// bladStanuStrony niesie stan HTTP odpowiedzi, której nie da się przerobić
// na migawkę, wraz ze zdaniem opisującym odmowę po polsku.
type bladStanuStrony struct {
	Kod    int
	Zdanie string
}

func (b bladStanuStrony) Error() string { return b.Zdanie }

// bladAdresuStrony nazywa odmowę wynikającą z samego adresu: zły protokół,
// brak gospodarza albo łańcuch przekierowań dłuższy niż granica.
type bladAdresuStrony struct {
	Zdanie string
}

func (b bladAdresuStrony) Error() string { return b.Zdanie }

// errZasobNieJestStrona nazywa zasób, którego nie da się pokazać jako
// stronę: obraz, dokument PDF, archiwum, i który moduł zamiast tego pobiera.
var errZasobNieJestStrona = errors.New("zasób nie jest stroną do odczytu")

// plikZeStrony niesie wynik pobrania zasobu, którego nie da się pokazać jako
// strony: same bajty, deklarowany typ i nazwę wziętą z adresu.
type plikZeStrony struct {
	Bajty []byte
	Typ   string
	Nazwa string
}

// pobierzPlik ściąga zasób spod adresu w całości i oddaje jego bajty. Granice są
// te same, co przy stronie: protokół, czas, rozmiar i łańcuch przekierowań —
// menedżer pobrań nie jest drogą obejścia zabezpieczeń pobierania stron.
func pobierzPlik(ctx context.Context, adres string) (plikZeStrony, error) {
	adres = strings.TrimSpace(adres)
	if err := sprawdzProtokol(adres); err != nil {
		return plikZeStrony{}, err
	}
	kontekst, anuluj := context.WithTimeout(ctx, czasPobrania)
	defer anuluj()

	zadanie, err := http.NewRequestWithContext(kontekst, http.MethodGet, adres, nil)
	if err != nil {
		return plikZeStrony{}, fmt.Errorf("adres nie da się złożyć w żądanie: %w", err)
	}
	zadanie.Header.Set("User-Agent", znacznikRdzenia)

	klient := &http.Client{Timeout: czasPobrania, CheckRedirect: pilnujPrzekierowan}
	odpowiedz, err := klient.Do(zadanie)
	if err != nil {
		var bladOtoczki *url.Error
		if errors.As(err, &bladOtoczki) && bladOtoczki.Err != nil {
			err = bladOtoczki.Err
		}
		return plikZeStrony{}, fmt.Errorf("%w: %v", errStronaNieodpowiedziala, err)
	}
	defer odpowiedz.Body.Close()

	if err := sprawdzStanOdpowiedzi(odpowiedz.StatusCode); err != nil {
		return plikZeStrony{}, err
	}
	if odpowiedz.ContentLength > limitTresci {
		return plikZeStrony{}, fmt.Errorf("%w (%d B przy granicy %d B)",
			errZaDuzaStrona, odpowiedz.ContentLength, int64(limitTresci))
	}
	bajty, err := io.ReadAll(io.LimitReader(odpowiedz.Body, limitTresci+1))
	if err != nil {
		return plikZeStrony{}, fmt.Errorf("pobieranie przerwane: %w", err)
	}
	if len(bajty) > limitTresci {
		return plikZeStrony{}, fmt.Errorf("%w (granica %d B)", errZaDuzaStrona, int64(limitTresci))
	}
	if len(bajty) == 0 {
		return plikZeStrony{}, fmt.Errorf("zasób nie oddał ani jednego bajtu")
	}
	return plikZeStrony{
		Bajty: bajty,
		Typ:   odpowiedz.Header.Get("Content-Type"),
		Nazwa: nazwaZAdresu(adres),
	}, nil
}

// nazwaZAdresu wyjmuje nazwę pliku z adresu. Adres bez ostatniego segmentu daje
// nazwę zastępczą — pobranie ma nazwę zawsze, bo bez niej wykaz pokazywałby
// pustą pozycję.
func nazwaZAdresu(adres string) string {
	rozebrany, err := url.Parse(adres)
	if err != nil {
		return "pobranie"
	}
	segmenty := strings.Split(strings.Trim(rozebrany.Path, "/"), "/")
	nazwa := segmenty[len(segmenty)-1]
	if strings.TrimSpace(nazwa) == "" {
		return "pobranie"
	}
	return nazwa
}

// trescStrony niesie wynik pobrania: tytuł, tekst renderowany i surowe HTML.
// Puste pole znaczy „strona tego nie miała", nie „nie pobrano" — brak pobrania
// wraca błędem, nie pustą treścią.
type trescStrony struct {
	Tytul string
	Tekst string
	Html  string
}

// pobierzStrone pobiera stronę spod adresu przez HTTP GET i wydobywa z niej
// tytuł oraz tekst renderowany, w granicach czasu i rozmiaru pobrania.
func pobierzStrone(ctx context.Context, adres string) (trescStrony, error) {
	// Adres obcinamy raz, na wejściu, i to jest cała prawda o nim dalej.
	adres = strings.TrimSpace(adres)
	if err := sprawdzProtokol(adres); err != nil {
		return trescStrony{}, err
	}

	kontekst, anuluj := context.WithTimeout(ctx, czasPobrania)
	defer anuluj()

	zadanie, err := http.NewRequestWithContext(kontekst, http.MethodGet, adres, nil)
	if err != nil {
		return trescStrony{}, fmt.Errorf("adres nie da się złożyć w żądanie: %w", err)
	}
	zadanie.Header.Set("User-Agent", znacznikRdzenia)
	zadanie.Header.Set("Accept",
		"text/html,application/xhtml+xml,text/plain;q=0.9,*/*;q=0.8")

	klient := &http.Client{Timeout: czasPobrania, CheckRedirect: pilnujPrzekierowan}
	odpowiedz, err := klient.Do(zadanie)
	if err != nil {
		// Rozwijamy `*url.Error`, żeby zdanie odmowy niosło samą przyczynę.
		var bladOtoczki *url.Error
		if errors.As(err, &bladOtoczki) && bladOtoczki.Err != nil {
			err = bladOtoczki.Err
		}
		// Odmowa CheckRedirect jest wadą adresu, nie milczeniem gospodarza.
		var bladAdresu bladAdresuStrony
		if errors.As(err, &bladAdresu) {
			return trescStrony{}, bladAdresu
		}
		return trescStrony{}, fmt.Errorf("%w: %w", errStronaNieodpowiedziala, err)
	}
	defer odpowiedz.Body.Close()

	if err := sprawdzStanOdpowiedzi(odpowiedz.StatusCode); err != nil {
		return trescStrony{}, err
	}
	if typ := odpowiedz.Header.Get("Content-Type"); !typTekstowy(typ) {
		return trescStrony{}, fmt.Errorf("%w, tylko treścią typu %q", errZasobNieJestStrona, typ)
	}
	// Deklarowany rozmiar sprawdzamy przed czytaniem odpowiedzi.
	if odpowiedz.ContentLength > limitTresci {
		return trescStrony{}, fmt.Errorf("%w (%d B przy granicy %d B)",
			errZaDuzaStrona, odpowiedz.ContentLength, int64(limitTresci))
	}

	// Czytamy o bajt więcej niż granica: nadmiarowy bajt jest dowodem, że
	// strona się nie zmieściła.
	surowe, err := io.ReadAll(io.LimitReader(odpowiedz.Body, limitTresci+1))
	if err != nil {
		return trescStrony{}, fmt.Errorf("czytanie strony przerwane: %w", err)
	}
	if len(surowe) > limitTresci {
		return trescStrony{}, fmt.Errorf("%w (granica %d B)", errZaDuzaStrona, int64(limitTresci))
	}
	html := string(surowe)
	return trescStrony{
		Tytul: wydobadzTytul(html),
		Tekst: wydobadzTekst(html),
		Html:  html,
	}, nil
}

// sprawdzStanOdpowiedzi rozstrzyga, czy z odpowiedzi o danym stanie może
// powstać migawka strony: tylko odpowiedź 2xx niosąca treść, reszta odmawia.
func sprawdzStanOdpowiedzi(kod int) error {
	switch {
	case kod == http.StatusNoContent || kod == http.StatusResetContent:
		return bladStanuStrony{Kod: kod, Zdanie: fmt.Sprintf(
			"strona odpowiedziała stanem %d (%s), czyli bez żadnej treści do pokazania",
			kod, strings.ToLower(http.StatusText(kod)))}
	case kod >= 300 && kod < 400:
		return bladStanuStrony{Kod: kod, Zdanie: fmt.Sprintf(
			"strona odpowiedziała przekierowaniem %d, ale nie powiedziała dokąd "+
				"(brak nagłówka Location), więc nie ma czego pobrać", kod)}
	case kod < 200 || kod >= 400:
		return bladStanuStrony{Kod: kod, Zdanie: fmt.Sprintf(
			"strona odpowiedziała stanem %d (%s)", kod,
			strings.ToLower(http.StatusText(kod)))}
	}
	return nil
}

// sprawdzProtokol przepuszcza wyłącznie `http` i `https` ze wskazanym
// gospodarzem, odmawiając zdaniem nazywającym brak drogi do strony.
func sprawdzProtokol(adres string) error {
	// Odstępy zdejmuje `pobierzStrone` na wejściu, więc adres tu już taki
	// pojedzie do żądania.
	rozlozony, err := url.Parse(adres)
	if err != nil {
		return bladAdresuStrony{Zdanie: "adres nie jest poprawnym adresem strony: " + err.Error()}
	}
	schemat := strings.ToLower(rozlozony.Scheme)
	if schemat != "http" && schemat != "https" {
		if schemat == "" {
			return bladAdresuStrony{Zdanie: "adres nie mówi, jakim protokołem pobrać stronę; " +
				"rdzeń pobiera tylko przez http i https"}
		}
		return bladAdresuStrony{Zdanie: fmt.Sprintf(
			"rdzeń pobiera strony tylko przez http i https, a ten adres wskazuje protokół %q", schemat)}
	}
	if rozlozony.Host == "" {
		return bladAdresuStrony{Zdanie: "adres nie wskazuje gospodarza, spod którego pobrać stronę"}
	}
	return nil
}

// pilnujPrzekierowan tnie łańcuch przekierowań na `limitPrzekierowan` skokach
// i sprawdza protokół każdego kolejnego adresu w łańcuchu.
func pilnujPrzekierowan(zadanie *http.Request, przebyte []*http.Request) error {
	if len(przebyte) > limitPrzekierowan {
		return bladAdresuStrony{Zdanie: fmt.Sprintf(
			"strona przekierowuje dalej niż %d razy; "+
				"rdzeń przerywa łańcuch, zamiast krążyć w nim bez końca", limitPrzekierowan)}
	}
	return sprawdzProtokol(zadanie.URL.String())
}

// typTekstowy rozstrzyga, czy nagłówek Content-Type opisuje treść do odczytu.
// Pusty nagłówek dopuszczamy — część serwerów go nie wysyła, a domyślnie
// oczekujemy strony. XHTML i zwykły tekst też przechodzą.
func typTekstowy(typ string) bool {
	t := strings.ToLower(strings.TrimSpace(typ))
	if t == "" {
		return true
	}
	// Kanały RSS, Atom i JSON Feed wchodzą tą samą drogą co strony i deklarują
	// własne typy treści.
	return strings.HasPrefix(t, "text/") ||
		strings.HasPrefix(t, "application/xhtml+xml") ||
		strings.HasPrefix(t, "application/xml") ||
		strings.HasPrefix(t, "application/rss+xml") ||
		strings.HasPrefix(t, "application/atom+xml") ||
		strings.HasPrefix(t, "application/feed+json") ||
		strings.HasPrefix(t, "application/json")
}

// wydobadzTytul zdejmuje zawartość znacznika <title>. Brak znacznika daje
// napis pusty — tytuł jest opcjonalny, jego brak nie jest błędem.
func wydobadzTytul(html string) string {
	dolne := strings.ToLower(html)
	poczatek := strings.Index(dolne, "<title")
	if poczatek < 0 {
		return ""
	}
	przesuniecie := strings.IndexByte(html[poczatek:], '>')
	if przesuniecie < 0 {
		return ""
	}
	od := poczatek + przesuniecie + 1
	koniec := strings.Index(dolne[od:], "</title>")
	if koniec < 0 {
		return ""
	}
	return normalizujOdstepy(odkodujEncje(html[od : od+koniec]))
}

// wydobadzTekst składa tekst renderowany strony: usuwa bloki skryptów
// i stylów, zdejmuje pozostałe znaczniki, odkodowuje encje i zbija odstępy.
func wydobadzTekst(html string) string {
	bez := usunBloki(html, "script")
	bez = usunBloki(bez, "style")

	var b strings.Builder
	wZnaczniku := false
	for i := 0; i < len(bez); i++ {
		switch bez[i] {
		case '<':
			wZnaczniku = true
		case '>':
			wZnaczniku = false
			b.WriteByte(' ') // granica znacznika staje się granicą słowa
		default:
			if !wZnaczniku {
				b.WriteByte(bez[i])
			}
		}
	}
	return normalizujOdstepy(odkodujEncje(b.String()))
}

// usunBloki wycina z HTML-a pary `<znacznik ...>...</znacznik>` wraz z treścią.
// Blok bez zamknięcia znaczy resztę strony jako swoją treść — pomijamy ją,
// zamiast wypuszczać surowy skrypt do tekstu.
func usunBloki(html, znacznik string) string {
	dolne := strings.ToLower(html)
	otw := "<" + znacznik
	zam := "</" + znacznik + ">"

	var b strings.Builder
	i := 0
	for {
		p := strings.Index(dolne[i:], otw)
		if p < 0 {
			b.WriteString(html[i:])
			break
		}
		p += i
		b.WriteString(html[i:p])
		k := strings.Index(dolne[p:], zam)
		if k < 0 {
			break
		}
		i = p + k + len(zam)
	}
	return b.String()
}

// odkodujEncje zamienia najczęstsze encje HTML na znaki. `&amp;` idzie ostatnia,
// by nie rozłożyć encji dopiero co wpisanego znaku `&`.
func odkodujEncje(s string) string {
	zamiany := []struct{ encja, znak string }{
		{"&lt;", "<"}, {"&gt;", ">"}, {"&quot;", "\""},
		{"&#39;", "'"}, {"&apos;", "'"}, {"&nbsp;", " "}, {"&amp;", "&"},
	}
	for _, z := range zamiany {
		s = strings.ReplaceAll(s, z.encja, z.znak)
	}
	return s
}

// normalizujOdstepy zbija każdą sekwencję białych znaków do jednej spacji i
// obcina brzegi. Migawka ma nieść treść, nie drabinę pustych wierszy z układu
// HTML-a.
func normalizujOdstepy(s string) string {
	pola := strings.FieldsFunc(s, func(r rune) bool { return unicode.IsSpace(r) })
	return strings.Join(pola, " ")
}
