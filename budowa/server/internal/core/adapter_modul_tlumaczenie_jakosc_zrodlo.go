// Modul Translate — kontrola jakosci porownujaca panel ze zrodlem. Metody
// i funkcje na adapterTlumaczenia; kontrole liczone z samej tresci panelu
// zostaja w osobnym pliku modulu.
package core

import (
	"context"
	"regexp"
	"strconv"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// tekstZrodlowyPanelu doczytuje tekst zrodlowy okna, do ktorego nalezy panel.
// Pusty wynik oznacza brak zapisanego zrodla — kontrola porownawcza wtedy sie
// nie odbywa; brak samego okna jest usterka i konczy sie odmowa.
func (a *adapterTlumaczenia) tekstZrodlowyPanelu(ctx context.Context,
	panel dane.PanelTlumaczenia) (string, error) {

	if panel.OknoKod == "" {
		return "", nil
	}
	okno, err := a.repozytorium.Okno(ctx, panel.OknoKod)
	if err != nil {
		return "", bladNieznanegoOkna(panel.OknoKod, err)
	}
	if okno.TekstZrodlowy == nil {
		return "", nil
	}
	return *okno.TekstZrodlowy, nil
}

// walutyKontrolowane to znaki i kody walut, ktorych obecnosc w zrodle sprawdzamy
// w przekladzie; wykaz jest zamkniety i krotki, zeby nie brac za kod waluty
// kazdego skrotowca z trzech wielkich liter.
var walutyKontrolowane = []string{
	"zł", "€", "$", "£", "¥", "₴", "₽", "₺",
	"PLN", "EUR", "USD", "GBP", "CHF", "JPY", "CZK", "SEK", "NOK", "DKK", "HUF", "UAH", "RON", "BGN",
}

// zapisDaty rozpoznaje date zapisana cyframi z mysnikiem, kropka albo ukosnikiem.
// Sluzy wylacznie do wyciecia daty przed liczeniem liczb — kontrola dat
// pozostaje niesprawdzana.
var zapisDaty = regexp.MustCompile(`\d{1,4}[-./]\d{1,2}[-./]\d{1,4}`)

// proporcjeDlugosciPrzekladu wyznacza pas, w ktorym dlugosc przekladu jest
// uznawana za normalna. Pas jest szeroki, bo jezyki roznia sie rozwleklosica;
// lapie tylko przypadki razace.
const (
	najkrotszaProporcjaPrzekladu = 0.4
	najdluzszaProporcjaPrzekladu = 2.5
)

// zbadajPanel skada komplet niezgodnosci panelu — jedno miejsce, z ktorego
// korzysta komenda quality.check oraz migawka zdejmowana zaraz po przekladzie
// modelu. Pusty tekst zrodlowy wylacza kontrole porownawcze.
func zbadajPanel(panel dane.PanelTlumaczenia, tekstZrodlowy string) []dane.NiezgodnoscTlumaczenia {
	var niezgodnosci []dane.NiezgodnoscTlumaczenia
	niezgodnosci = append(niezgodnosci, sprawdzZnaczniki(panel)...)
	niezgodnosci = append(niezgodnosci, sprawdzDlugosc(panel)...)

	tresc := ""
	if panel.Tresc != nil {
		tresc = *panel.Tresc
	}
	zrodlo := strings.TrimSpace(tekstZrodlowy)
	if zrodlo == "" || strings.TrimSpace(tresc) == "" {
		return niezgodnosci
	}
	niezgodnosci = append(niezgodnosci, sprawdzLiczbyWzgledemZrodla(zrodlo, tresc)...)
	niezgodnosci = append(niezgodnosci, sprawdzWalutyWzgledemZrodla(zrodlo, tresc)...)
	niezgodnosci = append(niezgodnosci, sprawdzZnacznikiWzgledemZrodla(zrodlo, tresc)...)
	niezgodnosci = append(niezgodnosci, sprawdzProporcjeDlugosci(zrodlo, tresc)...)
	return niezgodnosci
}

// sprawdzLiczbyWzgledemZrodla zglasza rodzaj number dla kazdej liczby obecnej
// w zrodle, a nieobecnej w przekladzie; porownanie idzie po wartosci, nie po
// napisie. Liczba dolozona w przekladzie nie jest zglaszana.
func sprawdzLiczbyWzgledemZrodla(zrodlo, tresc string) []dane.NiezgodnoscTlumaczenia {
	wPrzekladzie := wartosciLiczbowe(tresc)
	var wynik []dane.NiezgodnoscTlumaczenia
	zgloszone := map[string]bool{}
	for wartosc := range wartosciLiczbowe(zrodlo) {
		if wPrzekladzie[wartosc] || zgloszone[wartosc] {
			continue
		}
		zgloszone[wartosc] = true
		wynik = append(wynik, niezgodnoscRodzaju(shared.TranslationIssueKindNumber,
			"liczba "+wartosc+" jest w tekście źródłowym, ale nie ma jej w przekładzie"))
	}
	return wynik
}

// wartosciLiczbowe wylawia z tekstu liczby i sprowadza je do postaci
// porownywalnej: bez separatorow tysiecy, z kropka dziesietna, bez zer
// nieznaczacych. Wynikiem jest zbior wartosci bez powtorzen.
func wartosciLiczbowe(tekst string) map[string]bool {
	wynik := map[string]bool{}
	// Daty wycinamy przed liczeniem liczb, zeby zapis daty nie rozpadal sie na fałszywe braki liczb.
	tekst = zapisDaty.ReplaceAllString(tekst, " ")
	znaki := []rune(tekst)
	for i := 0; i < len(znaki); {
		if !cyfra(znaki[i]) {
			i++
			continue
		}
		poczatek := i
		for i < len(znaki) && (cyfra(znaki[i]) || separatorLiczby(znaki, i)) {
			i++
		}
		surowa := string(znaki[poczatek:i])
		if znormalizowana := znormalizujLiczbe(surowa); znormalizowana != "" {
			wynik[znormalizowana] = true
		}
	}
	return wynik
}

func cyfra(r rune) bool { return r >= '0' && r <= '9' }

// separatorLiczby mówi, czy znak na pozycji `i` jest separatorem wewnątrz
// liczby — kropką albo przecinkiem mającym cyfrę po obu stronach. Kropka
// kończąca zdanie („kosztuje 20. Następnie…") nie ma cyfry z prawej, więc nie
// zostanie wciągnięta do liczby.
func separatorLiczby(znaki []rune, i int) bool {
	if znaki[i] != '.' && znaki[i] != ',' && znaki[i] != ' ' {
		return false
	}
	return i > 0 && cyfra(znaki[i-1]) && i+1 < len(znaki) && cyfra(znaki[i+1])
}

// znormalizujLiczbe sprowadza zapis liczby do jednej postaci wedle reguly
// opisanej przy wartosciLiczbowe. Napis nieczytelny jako liczba oddaje pusty
// wynik i nie wchodzi do porownania.
func znormalizujLiczbe(surowa string) string {
	surowa = strings.ReplaceAll(surowa, " ", "")
	ostatni := strings.LastIndexAny(surowa, ".,")
	calosc := surowa
	ulamek := ""
	if ostatni >= 0 {
		ogon := surowa[ostatni+1:]
		if len(ogon) == 3 {
			// Trzy cyfry po ostatnim separatorze — separator tysięcy.
			calosc = strings.NewReplacer(".", "", ",", "").Replace(surowa)
		} else {
			calosc = strings.NewReplacer(".", "", ",", "").Replace(surowa[:ostatni])
			ulamek = ogon
		}
	}
	liczba, err := strconv.ParseFloat(calosc+"."+ulamek+"0", 64)
	if err != nil {
		return ""
	}
	return strconv.FormatFloat(liczba, 'f', -1, 64)
}

// sprawdzWalutyWzgledemZrodla zgłasza rodzaj `currency` dla waluty obecnej
// w źródle, a nieobecnej w przekładzie — podmieniona albo zgubiona waluta
// zmienia znaczenie kwoty, czego przekład robić nie może.
func sprawdzWalutyWzgledemZrodla(zrodlo, tresc string) []dane.NiezgodnoscTlumaczenia {
	var wynik []dane.NiezgodnoscTlumaczenia
	for _, waluta := range walutyKontrolowane {
		if !strings.Contains(zrodlo, waluta) || strings.Contains(tresc, waluta) {
			continue
		}
		wynik = append(wynik, niezgodnoscRodzaju(shared.TranslationIssueKindCurrency,
			"waluta „"+waluta+"” jest w tekście źródłowym, ale nie ma jej w przekładzie"))
	}
	return wynik
}

// sprawdzZnacznikiWzgledemZrodla zglasza rodzaj placeholder dla znacznika
// podstawienia obecnego w zrodle, a nieobecnego w przekladzie: model tlumaczy
// nazwe w srodku znacznika albo gubi znacznik zupelnie.
func sprawdzZnacznikiWzgledemZrodla(zrodlo, tresc string) []dane.NiezgodnoscTlumaczenia {
	var wynik []dane.NiezgodnoscTlumaczenia
	for _, znacznik := range znacznikiPodstawienia(zrodlo) {
		if strings.Contains(tresc, znacznik) {
			continue
		}
		wynik = append(wynik, niezgodnoscRodzaju(shared.TranslationIssueKindPlaceholder,
			"znacznik podstawienia "+znacznik+" jest w tekście źródłowym, ale nie ma go w przekładzie"))
	}
	return wynik
}

// znacznikiPodstawienia wyławia z tekstu domknięte znaczniki `{…}` wraz
// z klamrami. Znaczniki zagnieżdżone i klamry niedomknięte pomija — od nich
// jest `sprawdzZnaczniki`, kontrola składniowa z drugiego pliku.
func znacznikiPodstawienia(tekst string) []string {
	var wynik []string
	widziane := map[string]bool{}
	for {
		od := strings.Index(tekst, "{")
		if od < 0 {
			return wynik
		}
		do_ := strings.Index(tekst[od:], "}")
		if do_ < 0 {
			return wynik
		}
		znacznik := tekst[od : od+do_+1]
		if !strings.Contains(znacznik[1:], "{") && !widziane[znacznik] {
			widziane[znacznik] = true
			wynik = append(wynik, znacznik)
		}
		tekst = tekst[od+do_+1:]
	}
}

// sprawdzProporcjeDlugosci zglasza rodzaj length, gdy przeklad jest razaco
// krotszy albo dluzszy od zrodla. Liczymy w runach, nie w bajtach, zeby
// diakrytyka i alfabet nielacinski nie liczyly sie jako nadmiar.
func sprawdzProporcjeDlugosci(zrodlo, tresc string) []dane.NiezgodnoscTlumaczenia {
	dlugoscZrodla := len([]rune(strings.TrimSpace(zrodlo)))
	dlugoscPrzekladu := len([]rune(strings.TrimSpace(tresc)))
	if dlugoscZrodla == 0 {
		return nil
	}
	proporcja := float64(dlugoscPrzekladu) / float64(dlugoscZrodla)
	if proporcja >= najkrotszaProporcjaPrzekladu && proporcja <= najdluzszaProporcjaPrzekladu {
		return nil
	}
	opis := "przekład jest rażąco krótszy od źródła"
	if proporcja > najdluzszaProporcjaPrzekladu {
		opis = "przekład jest rażąco dłuższy od źródła"
	}
	return []dane.NiezgodnoscTlumaczenia{niezgodnoscRodzaju(shared.TranslationIssueKindLength,
		opis+" ("+strconv.Itoa(dlugoscZrodla)+" znaków źródła wobec "+
			strconv.Itoa(dlugoscPrzekladu)+" znaków przekładu)")}
}

// niezgodnoscRodzaju skada jeden wiersz niezgodnosci dowolnego rodzaju. Pole
// Segment zostaje puste: kontrole tego pliku porownuja cale teksty, nie
// segment po segmencie.
func niezgodnoscRodzaju(rodzaj shared.TranslationIssueKind, szczegol string) dane.NiezgodnoscTlumaczenia {
	s := szczegol
	return dane.NiezgodnoscTlumaczenia{Rodzaj: string(rodzaj), Szczegol: &s}
}
