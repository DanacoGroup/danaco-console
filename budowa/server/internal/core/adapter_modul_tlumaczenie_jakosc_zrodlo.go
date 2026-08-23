// Odpowiedzialność pliku: moduł Translate — kontrola jakości porównująca panel
// ze źródłem. Metody i funkcje na `*adapterTlumaczenia` (typ deklaruje
// `adapter_modul_tlumaczenie.go`); kontrole liczone z samej treści panelu —
// znaczniki niedomknięte i pusty panel „gotowy" — zostają
// w `adapter_modul_tlumaczenie_jakosc.go`.
//
// Porównanie ze źródłem jest wykonalne, bo panel widzi swoje okno:
// `quality.check` niesie sam `PanelId`, a `PanelTlumaczenia.OknoKod`
// (doczytywany złączeniem) daje kod zewnętrzny okna, którego wymaga
// `RepozytoriumTlumaczen.Okno`.
//
// Czego nie sprawdzamy i dlaczego nie udajemy, że sprawdzamy:
//
//   - `date` — data poprawnie przełożona zmienia zapis („3 maja 2024" ↔
//     „May 3, 2024" ↔ „2024-05-03"). Mechaniczne porównanie napisów uznałoby
//     poprawny przekład za wadę, a fałszywy alarm w kontroli jakości jest
//     gorszy od braku kontroli: uczy Operatora ignorować wynik. Uczciwe
//     sprawdzenie wymaga rozpoznania daty w obu językach, czego rdzeń nie ma.
//     Zasada `date` idzie za to do polecenia dla modelu (`zasadyJakosci`) —
//     żądać można więcej, niż da się zweryfikować.
//   - `omission` — pominięcie zdania to rzecz znaczeniowa, nie napisowa:
//     stwierdzenie, czy segment źródłowy ma swój odpowiednik w przekładzie,
//     wymaga rozumienia treści. Zgrubny zastępnik (liczba zdań w źródle vs
//     w panelu) mylnie oskarżałby każdy przekład, który łączy albo dzieli
//     zdania — a to jest normalna, dobra robota tłumacza. Sprawdzenie
//     przybliżone długością zostaje pod rodzajem `length`, gdzie jest uczciwe.
package core

import (
	"context"
	"regexp"
	"strconv"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// tekstZrodlowyPanelu doczytuje tekst źródłowy okna, do którego należy panel —
// drogą przez kod zewnętrzny okna, jedyny, jaki przyjmuje `Okno`.
//
// Treść żyjąca w pliku (`TekstZrodlowyOdwolanie`) nie jest czytana:
// ta warstwa nie sięga po pliki spoza repozytorium. Wtedy wynik jest pusty,
// a kontrola porównawcza po prostu się nie odbywa — pusty tekst źródłowy nie
// udaje źródła, przez co żadna liczba nie zostanie fałszywie zgłoszona jako
// zgubiona. Brak okna jest natomiast usterką, nie normalnym stanem:
// panel bez okna nie miałby prawa istnieć, więc idzie odmową.
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

// walutyKontrolowane to znaki i kody walut, których obecność w źródle
// sprawdzamy w przekładzie. Wykaz jest zamknięty i krótki: „każde trzy wielkie
// litery" brałoby za kod waluty każdy skrótowiec (`API`, `VAT`, `PDF`)
// i zasypywało Operatora fałszywymi zastrzeżeniami.
var walutyKontrolowane = []string{
	"zł", "€", "$", "£", "¥", "₴", "₽", "₺",
	"PLN", "EUR", "USD", "GBP", "CHF", "JPY", "CZK", "SEK", "NOK", "DKK", "HUF", "UAH", "RON", "BGN",
}

// zapisDaty rozpoznaje datę zapisaną cyframi z myślnikiem, kropką albo ukośnikiem
// („2024-05-03", „3.05.2024", „05/03/2024"). Służy wyłącznie do wycięcia daty
// przed liczeniem liczb (`wartosciLiczbowe`) — nie jest kontrolą dat i nie
// udaje jej; rodzaj `date` pozostaje niesprawdzany (nagłówek pliku).
var zapisDaty = regexp.MustCompile(`\d{1,4}[-./]\d{1,2}[-./]\d{1,4}`)

// proporcjeDlugosciPrzekladu wyznacza pas, w którym długość przekładu jest
// uznawana za normalną. Języki różnią się rozwlekłością — przekład z angielskiego
// na polski bywa o połowę dłuższy, na węgierski krótszy — więc pas jest
// szeroki z rozmysłem. Nie mierzy stylu; łapie tylko przypadki rażące:
// przekład ucięty w połowie albo model, który zamiast tłumaczyć napisał
// wypracowanie.
const (
	najkrotszaProporcjaPrzekladu = 0.4
	najdluzszaProporcjaPrzekladu = 2.5
)

// zbadajPanel składa komplet niezgodności panelu — jedno miejsce, z którego
// korzystają obie drogi kontroli: komenda `quality.check` wołana przez
// Operatora i migawka zdejmowana zaraz po przekładzie modelu (`target.add`).
// Gdyby każda liczyła po swojemu, Operator dostawałby dwa różne wykazy wad
// tego samego panelu.
//
// Puste `tekstZrodlowy` (okno bez zapisanego źródła) wyłącza kontrole
// porównawcze — nie ma z czym porównywać. Kontrole z samej treści panelu
// działają zawsze.
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

// sprawdzLiczbyWzgledemZrodla zgłasza rodzaj `number` dla każdej liczby
// obecnej w źródle, a nieobecnej w przekładzie. Porównanie idzie po wartości,
// nie po napisie: „1 234,50", „1.234,50" i „1234.5" to ta sama liczba zapisana
// zwyczajem trzech różnych języków (`wartosciLiczbowe`), a przekład ma prawo
// zmienić zapis — nie ma prawa zmienić wartości.
//
// Kierunek sprawdzenia jest jednostronny: liczba dołożona w przekładzie nie
// jest zgłaszana. Bywa dorobiona uczciwie (rozwinięcie „dwa" na „2")
// i zgłaszanie jej dawałoby fałszywe alarmy; liczba zgubiona jest natomiast
// zawsze wadą.
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

// wartosciLiczbowe wyławia z tekstu liczby i sprowadza je do postaci
// porównywalnej: bez separatorów tysięcy, z kropką dziesiętną, bez zer
// nieznaczących. Wynikiem jest zbiór wartości (mapa na `bool`, bo Go nie ma
// typu zbioru) — powtórzenie tej samej liczby w tekście nie jest osobnym
// faktem do sprawdzenia.
//
// Ograniczenie: rozstrzygnięcie, czy przecinek w „1,5" oddziela
// część dziesiętną (zwyczaj polski), czy tysiące (zwyczaj angielski), jest
// niejednoznaczne z samego napisu. Przyjmujemy regułę: separator, po którym
// idą dokładnie trzy cyfry i nic więcej z rzędu, to separator tysięcy; każdy
// inny — dziesiętny. Reguła myli się na „1,500" oznaczającym półtora tysiąca
// z przecinkiem dziesiętnym, ale trafia w zdecydowanej większości zapisów.
func wartosciLiczbowe(tekst string) map[string]bool {
	wynik := map[string]bool{}
	// Daty wycinamy przed liczeniem liczb. Bez tego „2024-05-03" rozpada się na
	// trzy liczby (2024, 5, 3), a poprawny przekład „w maju 2024" gubi dwie
	// z nich — kontrola zgłaszałaby dwa fałszywe braki liczb za jedną poprawnie
	// przełożoną datę. Skoro rodzaj `date` świadomie nie jest sprawdzany
	// (nagłówek pliku), to jego składniki nie wchodzą bocznymi drzwiami jako
	// `number`.
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

// znormalizujLiczbe sprowadza zapis liczby do jednej postaci wedle reguły
// opisanej przy `wartosciLiczbowe`. Napis, którego nie da się odczytać jako
// liczby, oddaje pusty wynik — wtedy nie wchodzi do porównania wcale, zamiast
// wchodzić jako wartość zmyślona.
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

// sprawdzZnacznikiWzgledemZrodla zgłasza rodzaj `placeholder` dla znacznika
// podstawienia obecnego w źródle, a nieobecnego w przekładzie. Kontrola z samej
// treści panelu widzi wyłącznie klamrę niedomkniętą, a najczęstsza wada
// znacznika jest inna — model tłumaczy nazwę w środku („{userName}" →
// „{nazwaUżytkownika}") albo gubi znacznik zupełnie, i program podstawia wtedy
// w pustkę.
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

// sprawdzProporcjeDlugosci zgłasza rodzaj `length`, gdy przekład jest rażąco
// krótszy albo dłuższy od źródła — pas dopuszczalny i jego uzasadnienie przy
// stałych `najkrotszaProporcjaPrzekladu`/`najdluzszaProporcjaPrzekladu`.
// Liczymy w runach, nie w bajtach: w bajtach każdy przekład na język
// z diakrytyką albo alfabetem niełacińskim wychodziłby „za długi".
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

// niezgodnoscRodzaju składa jeden wiersz niezgodności dowolnego rodzaju.
// `Segment` zostaje pusty: kontrole tego pliku porównują całe teksty, nie
// segment po segmencie (przypisanie wady do segmentu wymagałoby dopasowania
// segmentów źródła do segmentów przekładu, czyli tego samego rozumienia
// treści, którego brak wyklucza sprawdzanie `omission` — nagłówek pliku).
func niezgodnoscRodzaju(rodzaj shared.TranslationIssueKind, szczegol string) dane.NiezgodnoscTlumaczenia {
	s := szczegol
	return dane.NiezgodnoscTlumaczenia{Rodzaj: string(rodzaj), Szczegol: &s}
}
