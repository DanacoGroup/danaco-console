// Odpowiedzialność pliku: rozróżnienie, czym jest wypowiedź Operatora —
// poleceniem dla modelu czy akcją platformy (komendą kontraktu wykonywaną przez
// rdzeń). Plik rozstrzyga i nazywa podstawę rozstrzygnięcia; nie wykonuje ani
// jednego, ani drugiego.
//
// Dlaczego to nie jest zgadywanie po treści.
// Rdzeń nie ma prawa domyślać się z wolnego tekstu, że Operator chciał wykonać
// komendę platformy. Ten plik tego nie łamie, bo rozstrzyga wyłącznie na
// deklaracjach, nigdy na podobieństwie:
//
//  1. Wskazanie Operatora — pole żądania kontraktu (`intent`). Deklaracja
//     wprost; nie ma czego zgadywać.
//  2. Leksykon akcji — zamknięty wykaz fraz, z których każda jest przypisana
//     jednej komendzie kontraktu. Dopasowanie idzie po całej wypowiedzi, znak
//     w znak po normalizacji, a nie po fragmencie (patrz `dopasujAkcje`).
//  3. Domyślnie — model. Wolny tekst bez deklaracji jedzie do modelu, więc to
//     rozstrzygnięcie niczego nie zmienia za plecami Operatora.
//
// Gdzie deklaracji brak, a domyślnej drogi wziąć nie wolno (Operator zadeklarował
// akcję platformy, lecz nie nazwał żadnej), wynikiem jest zamiar nierozstrzygnięty
// z nazwanym powodem. Nierozstrzygnięcie jest tu wynikiem pełnoprawnym: lepsze
// od wybrania komendy za Operatora.
//
// Czego ten plik nie robi.
// Nie woła komend, nie zna rdzenia i — tak jak reszta pakietu — nie zna
// kontraktu: mówi własnymi typami, żeby dało się go wpiąć niezależnie od tego,
// kiedy pole `intent` wejdzie do `shared/contract.json`.
package mowa

import (
	"sort"
	"strconv"
	"strings"
	"unicode"
)

// Zamiar mówi, czym wypowiedź jest dla rdzenia.
//
// Wartość zerowa jest nierozstrzygnięciem i to jest wybór rozmyślny: struktura
// wynikowa, która powstała przez pomyłkę zamiast przez rozstrzygnięcie, ma
// znaczyć „nikt niczego nie zlecił", a nie „jedź do modelu" ani tym bardziej
// „wykonaj komendę". Najbezpieczniejszy wynik ma być najtańszy do przypadkowego
// wytworzenia.
type Zamiar string

const (
	// ZamiarNierozstrzygniety — nie wiadomo, czego Operator chciał, i nie wolno
	// tego dobrać za niego. Wynik niesie wtedy niepusty `Powod`.
	ZamiarNierozstrzygniety Zamiar = ""
	// ZamiarModelu — wypowiedź jest poleceniem dla modelu; jedzie turą kanału,
	// dokładnie tak jak dziś jedzie każde zlecenie asystenta.
	ZamiarModelu Zamiar = "model"
	// ZamiarPlatformy — wypowiedź jest akcją platformy: nazwaną komendą
	// kontraktu, którą ma wykonać rdzeń, a nie model.
	ZamiarPlatformy Zamiar = "platform"
)

// Trzy wartości pola wskazania — dokładnie te, o które prosi zgłoszenie
// kontraktu (`intent` w `assistant.voice.command`).
//
// Napisy są te same co wartości `Zamiar` i to nie jest przypadek: pole kontraktu
// ma nieść rozstrzygnięcie Operatora w tym samym słowniku, w którym rdzeń o nim
// mówi. Drugi słownik na to samo byłby drugą prawdą.
const (
	// WskazanieBrak — żądanie pola nie niesie (albo kontrakt jeszcze go nie ma).
	WskazanieBrak = ""
	// WskazanieModelu — Operator zadeklarował polecenie dla modelu.
	WskazanieModelu = "model"
	// WskazaniePlatformy — Operator zadeklarował akcję platformy.
	WskazaniePlatformy = "platform"
)

// Podstawa nazywa, czym rozstrzygnięto — bo „platforma" bez odpowiedzi na
// pytanie „skąd wiesz" jest nie do zdiagnozowania, gdy Operator zgłosi, że
// rdzeń zrobił nie to, co powiedział.
type Podstawa string

const (
	// PodstawaBrak — nie rozstrzygnięto.
	PodstawaBrak Podstawa = ""
	// PodstawaWskazania — rozstrzygnęło pole żądania.
	PodstawaWskazania Podstawa = "wskazanie"
	// PodstawaFrazy — rozstrzygnęła cała wypowiedź równa frazie leksykonu.
	PodstawaFrazy Podstawa = "fraza"
	// PodstawaDomyslna — nic nie wskazało akcji platformy, więc jedzie model.
	PodstawaDomyslna Podstawa = "domyslnie"
)

// AkcjaPlatformy wiąże frazy Operatora z jedną komendą kontraktu.
//
// Nastawy są dosłowne, nie wyliczane. Fraza „anuluj zlecenie" niesie
// `control=cancel` dlatego, że tak zapisano w leksykonie — nie dlatego, że coś
// przeczytało słowo „anuluj". Dzięki temu leksykon da się przeczytać jak umowę
// i wskazać palcem, co która wypowiedź robi.
type AkcjaPlatformy struct {
	// Frazy — całe wypowiedzi (nie fragmenty) uruchamiające tę akcję.
	Frazy []string
	// Komenda — nazwa komendy kontraktu, np. `assistant.action.status`.
	Komenda string
	// Nastawy — pola żądania o wartościach dosłownych, wpisanych w leksykonie.
	Nastawy map[string]string
	// Wymaga — pola żądania, których leksykon nie zna i których wołający musi
	// dołożyć z kontekstu wywołania (np. `windowId` z `assistant.voice.command`).
	// Wykaz stoi tu po to, żeby wołający miał czym sprawdzić, czy ma komplet —
	// zamiast wysyłać żądanie niepełne i dowiadywać się tego z odmowy rdzenia.
	Wymaga []string
	// Opis — zdanie dla Operatora, czym ta akcja jest.
	Opis string
}

// Wypowiedz jest wejściem rozstrzygnięcia: to, co przyszło komendą głosową.
type Wypowiedz struct {
	// Tekst — transkrypcja (własna albo poprawiona przez Operatora).
	Tekst string
	// Wskazanie — deklaracja Operatora z pola żądania; puste znaczy brak pola.
	Wskazanie string
}

// Rozstrzygniecie jest wynikiem rozróżnienia.
type Rozstrzygniecie struct {
	// Zamiar — model, platforma albo nierozstrzygnięty (wartość zerowa).
	Zamiar Zamiar
	// Akcja — wypełniona wyłącznie przy `ZamiarPlatformy`; dla modelu pusta,
	// bo model nie dostaje komendy, tylko treść.
	Akcja AkcjaPlatformy
	// Podstawa — czym rozstrzygnięto.
	Podstawa Podstawa
	// Powod — wypełniony wyłącznie przy zamiarze nierozstrzygniętym; trzy
	// części, tak jak każda odmowa tego pakietu: co · dlaczego · czym naprawić.
	Powod string
}

// DoModelu i DoPlatformy odpowiadają na pytanie wołającego wprost, żeby nie
// musiał porównywać napisów u siebie (dwa porównania w dwóch miejscach to dwie
// okazje do rozjazdu).
func (r Rozstrzygniecie) DoModelu() bool    { return r.Zamiar == ZamiarModelu }
func (r Rozstrzygniecie) DoPlatformy() bool { return r.Zamiar == ZamiarPlatformy }

// Rozstrzygniete mówi, czy wynik w ogóle rozstrzyga. Fałsz znaczy: nic nie
// wykonuj, oddaj Operatorowi `Powod`.
func (r Rozstrzygniecie) Rozstrzygniete() bool { return r.Zamiar != ZamiarNierozstrzygniety }

// Rozstrzygnij rozróżnia polecenie dla modelu od akcji platformy.
//
// Kolejność reguł nie jest dowolna i wygląda tak:
//
//	tekst pusty                      → nierozstrzygnięty (nie ma czego rozstrzygać)
//	wskazanie nieznane               → nierozstrzygnięty (napis spoza słownika)
//	wskazanie = model                → MODEL, choćby fraza pasowała do akcji
//	wskazanie = platform + fraza     → PLATFORMA (komenda z leksykonu)
//	wskazanie = platform, brak frazy → nierozstrzygnięty (której komendy?)
//	brak wskazania + fraza           → PLATFORMA
//	brak wskazania, brak frazy       → MODEL (droga dotychczasowa)
//
// Deklaracja Operatora bije leksykon. Gdy Operator napisał `intent=model`, a
// wypowiedź brzmi jak fraza akcji, wygrywa `intent`: Operator mógł chcieć, żeby
// model pochylił się nad tym zdaniem (na przykład wyjaśnił je albo przetłumaczył).
// Odwrotne pierwszeństwo znaczyłoby, że pole kontraktu można przegłosować
// leksykonem, czyli że deklaracja nie jest deklaracją.
//
// Leksykon pusty (`akcje == nil`) jest poprawnym wejściem, nie brakiem: rdzeń
// bez wpiętego leksykonu rozstrzyga wtedy wszystko na model — dokładnie tak, jak
// zachowuje się dziś. Podstawienie tu leksykonu domyślnego z własnej głowy
// dałoby zachowanie, którego wołający nie zamówił; kto chce domyślnego, woła
// `AkcjeDomyslne()` u siebie i widzi to w swoim kodzie.
func Rozstrzygnij(w Wypowiedz, akcje []AkcjaPlatformy) Rozstrzygniecie {
	tresc := znormalizuj(w.Tekst)
	if tresc == "" {
		return nierozstrzygniety("wypowiedź jest pusta, więc nie ma czego rozróżnić" +
			"; naprawa: przysłać treść polecenia (pole transcript) albo nagranie do rozpoznania")
	}

	wskazanie := strings.ToLower(strings.TrimSpace(w.Wskazanie))
	switch wskazanie {
	case WskazanieModelu:
		return Rozstrzygniecie{Zamiar: ZamiarModelu, Podstawa: PodstawaWskazania}
	case WskazaniePlatformy, WskazanieBrak:
		// dalej
	default:
		return nierozstrzygniety("wskazanie zamiaru „" + strings.TrimSpace(w.Wskazanie) +
			"” nie jest ani „" + WskazanieModelu + "”, ani „" + WskazaniePlatformy +
			"”, więc rdzeń nie wie, czy uruchomić model, czy komendę platformy" +
			"; naprawa: przysłać jedną z dwóch wartości albo pominąć pole")
	}

	dopasowane := dopasujAkcje(tresc, akcje)
	switch {
	case len(dopasowane) == 1:
		return Rozstrzygniecie{
			Zamiar:   ZamiarPlatformy,
			Akcja:    dopasowane[0],
			Podstawa: PodstawaFrazy,
		}
	case len(dopasowane) > 1:
		// Dwie komendy na jedną frazę to usterka leksykonu, nie wybór. Wzięcie
		// pierwszej z brzegu byłoby zgadywaniem przebranym za rozstrzygnięcie —
		// i to zgadywaniem cichym, bo nikt by się o drugim dopasowaniu nie
		// dowiedział.
		return nierozstrzygniety("wypowiedź „" + tresc + "” pasuje do " +
			strconv.Itoa(len(dopasowane)) + " akcji platformy naraz (" +
			strings.Join(nazwyKomend(dopasowane), ", ") +
			"), więc nie wiadomo, którą Operator miał na myśli" +
			"; naprawa: rozdzielić frazy w leksykonie akcji tak, żeby każda wskazywała jedną komendę")
	}

	if wskazanie == WskazaniePlatformy {
		return nierozstrzygniety("Operator zadeklarował akcję platformy, ale wypowiedź „" +
			tresc + "” nie jest żadną z " + strconv.Itoa(len(akcje)) +
			" znanych fraz akcji, a rdzeń nie dobiera komendy za Operatora" +
			"; naprawa: wydać polecenie frazą z wykazu akcji albo przysłać wskazanie „" +
			WskazanieModelu + "”")
	}

	return Rozstrzygniecie{Zamiar: ZamiarModelu, Podstawa: PodstawaDomyslna}
}

// nierozstrzygniety składa wynik bez zamiaru z nazwanym powodem.
func nierozstrzygniety(powod string) Rozstrzygniecie {
	return Rozstrzygniecie{Zamiar: ZamiarNierozstrzygniety, Podstawa: PodstawaBrak, Powod: powod}
}

// Wywolanie jest gotowym żądaniem akcji platformy: nazwa komendy kontraktu
// i komplet pól, z jakimi ma pojechać.
//
// Typ jest opisem, nie wykonaniem. Pakiet mowy nie ma i nie będzie miał drogi
// do rdzenia — składa wywołanie i oddaje je temu, kto komendy wykonuje. Pola są
// napisami, bo napisami przychodzą i z leksykonu, i z żądania głosowego;
// przełożenie ich na typ kontraktu należy do wołającego, który kontrakt zna.
type Wywolanie struct {
	// Komenda — nazwa komendy kontraktu.
	Komenda string
	// Pola — nastawy leksykonu ORAZ pola dołożone z kontekstu wywołania.
	Pola map[string]string
}

// BrakDanychAkcji jest odmową: akcja jest rozpoznana, lecz nie ma z czym
// pojechać, bo wołający nie dołożył pola, którego leksykon nie zna.
//
// Osobny typ z tego samego powodu, dla którego pakiet ma trzy odmowy w bledy.go:
// wołający musi ODRÓŻNIĆ ten przypadek od nierozstrzygnięcia. Tam Operator ma
// powiedzieć co innego; tutaj Operator powiedział wszystko, a brakuje danych
// po stronie rdzenia. Komunikat trójczęściowy jak każdy w tym pakiecie.
type BrakDanychAkcji struct {
	// Komenda — akcja, której nie da się złożyć.
	Komenda string
	// Pole — nazwa brakującego pola żądania.
	Pole string
}

func (b *BrakDanychAkcji) Error() string {
	return "akcji platformy " + b.Komenda + " nie da się złożyć: wołający nie podał pola " +
		b.Pole + ", którego leksykon akcji nie zna i nie ma prawa zmyślić" +
		"; naprawa: dołożyć " + b.Pole + " z kontekstu wywołania komendy głosowej"
}

// Unwrap nie ma czego oddać — odmowa rozstrzyga się na brakującym kluczu mapy,
// a nie na cudzym błędzie (tak samo jak `BrakSilnika` i `BrakNagrania`).
func (b *BrakDanychAkcji) Unwrap() error { return nil }

// NieAkcjaPlatformy jest odmową: ktoś prosi o złożenie wywołania komendy dla
// wypowiedzi, która akcją platformy nie jest.
//
// Odmowa istnieje po to, żeby pomyłka wołającego kończyła się BŁĘDEM, a nie
// pustym wywołaniem, które dałoby się wysłać. Puste `Wywolanie{}` bez błędu
// byłoby atrapą wyglądającą na wynik.
type NieAkcjaPlatformy struct {
	// Zamiar — co rozstrzygnięcie naprawdę mówi.
	Zamiar Zamiar
	// Powod — powód nierozstrzygnięcia, gdy zamiaru nie ma; pusty przy modelu.
	Powod string
}

func (b *NieAkcjaPlatformy) Error() string {
	if b.Zamiar == ZamiarModelu {
		return "wypowiedź jest poleceniem dla modelu, a nie akcją platformy," +
			" więc nie ma z niej czego złożyć jako komendy" +
			"; naprawa: poprowadzić tę wypowiedź turą kanału modelu"
	}
	powod := strings.TrimSpace(b.Powod)
	if powod == "" {
		powod = "rozstrzygnięcie nie zapadło"
	}
	return "wypowiedzi nie rozstrzygnięto jako akcji platformy: " + powod
}

// Unwrap kończy łańcuch — jak w pozostałych odmowach pakietu.
func (b *NieAkcjaPlatformy) Unwrap() error { return nil }

// ZlozWywolanie składa żądanie akcji platformy z nastaw leksykonu i pól, które
// wołający bierze z kontekstu komendy głosowej (przede wszystkim `windowId`).
//
// Trzy odmowy, każda inna:
//
//   - rozstrzygnięcie nie jest akcją platformy → odmowa; wołający, który mimo
//     tego złożyłby wywołanie, wykonałby czynność niezleconą;
//   - brak pola z wykazu `Wymaga` → `BrakDanychAkcji` z nazwą pola;
//   - pole puste liczy się jak brak — „windowId: ”" pojechałoby do rdzenia jako
//     żądanie bez okna i wróciło odmową gorzej opisaną niż ta tutaj.
//
// Kontekst nie nadpisuje leksykonu. Wołający dokłada tożsamości (okno, sesja),
// a nie sterowanie: gdyby kontekst mógł podmienić `control`, akcja „wstrzymaj"
// dałaby się w locie zamienić w „anuluj" i leksykon przestałby być umową.
func (r Rozstrzygniecie) ZlozWywolanie(kontekst map[string]string) (Wywolanie, error) {
	if !r.DoPlatformy() {
		return Wywolanie{}, &NieAkcjaPlatformy{Zamiar: r.Zamiar, Powod: r.Powod}
	}

	pola := make(map[string]string, len(r.Akcja.Nastawy)+len(r.Akcja.Wymaga))
	for _, wymagane := range r.Akcja.Wymaga {
		wartosc := strings.TrimSpace(kontekst[wymagane])
		if wartosc == "" {
			return Wywolanie{}, &BrakDanychAkcji{Komenda: r.Akcja.Komenda, Pole: wymagane}
		}
		pola[wymagane] = wartosc
	}
	for pole, wartosc := range r.Akcja.Nastawy {
		pola[pole] = wartosc
	}
	return Wywolanie{Komenda: r.Akcja.Komenda, Pola: pola}, nil
}

// dopasujAkcje szuka akcji, których fraza jest równa całej wypowiedzi.
//
// Równa, a nie zawarta — i to jest sedno tego, czym rozróżnienie różni się od
// zgadywania. Dopasowanie po fragmencie zamieniłoby wzmiankę w wykonanie:
// „przypomnij mi, co robi anuluj zlecenie" anulowałoby zlecenie, choć Operator
// prosił o wyjaśnienie. Wypowiedź równa frazie jest natomiast deklaracją samą
// w sobie — nie da się jej powiedzieć przypadkiem.
//
// Zwracane są wszystkie dopasowania, także sprzeczne: rozstrzyganie sprzeczności
// należy do `Rozstrzygnij`, a nie do wyszukiwania (jedna funkcja, jedna robota).
func dopasujAkcje(tresc string, akcje []AkcjaPlatformy) []AkcjaPlatformy {
	var trafienia []AkcjaPlatformy
	for _, akcja := range akcje {
		for _, fraza := range akcja.Frazy {
			if znormalizuj(fraza) == tresc {
				trafienia = append(trafienia, akcja)
				break
			}
		}
	}
	return trafienia
}

// nazwyKomend wylicza komendy dopasowanych akcji — do komunikatu o sprzeczności.
// Kolejność ustalona sortowaniem, żeby ten sam leksykon dawał ten sam komunikat
// przy każdym uruchomieniu (komunikat zmienny przy stałym wejściu jest mylący).
func nazwyKomend(akcje []AkcjaPlatformy) []string {
	nazwy := make([]string, 0, len(akcje))
	for _, akcja := range akcje {
		nazwy = append(nazwy, akcja.Komenda)
	}
	sort.Strings(nazwy)
	return nazwy
}

// bezOgonkow zdejmuje polskie znaki diakrytyczne.
//
// Potrzebne, bo wejście jest z mowy, a nie z klawiatury. Pomocnik transkrypcji
// oddaje raz „wznów", raz „wznow" — zależnie od modelu i od tego, jak wyraźnie
// Operator mówił. Fraza leksykonu, która rozpoznaje tylko jeden z tych zapisów,
// nie działałaby losowo, a Operator nie miałby jak się dowiedzieć dlaczego.
var bezOgonkow = map[rune]rune{
	'ą': 'a', 'ć': 'c', 'ę': 'e', 'ł': 'l', 'ń': 'n',
	'ó': 'o', 'ś': 's', 'ź': 'z', 'ż': 'z',
}

// znormalizuj sprowadza wypowiedź do postaci porównywalnej: małe litery, bez
// ogonków, bez znaków przestankowych, pojedyncze odstępy, bez brzegów.
//
// Normalizacja nie jest domyślaniem się treści: nie usuwa słów, nie zmienia ich
// kolejności i nie skraca. Zdejmuje wyłącznie to, czego mowa nie niesie —
// wielkość liter i interpunkcję, które dokłada pomocnik transkrypcji.
func znormalizuj(tekst string) string {
	var wynik strings.Builder
	odstep := false
	for _, znak := range strings.ToLower(tekst) {
		if zamiennik, jest := bezOgonkow[znak]; jest {
			znak = zamiennik
		}
		if unicode.IsLetter(znak) || unicode.IsDigit(znak) {
			if odstep && wynik.Len() > 0 {
				wynik.WriteRune(' ')
			}
			odstep = false
			wynik.WriteRune(znak)
			continue
		}
		odstep = true
	}
	return wynik.String()
}

// Nazwy komend leksykonu domyślnego. Napisy, a nie stałe z `shared`, bo pakiet
// mowy kontraktu nie zna (patrz nagłówek pliku).
const komendaStanZlecenia = "assistant.action.status"

// AkcjeDomyslne oddaje leksykon startowy akcji platformy.
//
// Dlaczego taki wąski.
// Wszystkie pięć fraz prowadzi do jednej komendy — `assistant.action.status` —
// i nie jest to niedoróbka, tylko granica postawiona z dwóch własności kontraktu:
//
//  1. Komenda ta ma wszystkie pola żądania nieobowiązkowe poza żadnym: wystarcza
//     `windowId`, a to jedyny identyfikator, który `assistant.voice.command`
//     naprawdę niesie. Akcje w rodzaju `window.create` (wymaga `sessionId`,
//     `moduleId`, `modelChannelId`, `workingDirs`, `executionEnv`,
//     `permissionMode`, `windowRole`) albo `session.stop` (wymaga `sessionId`)
//     musiałyby te pola skądś wziąć — a „skądś" znaczyłoby „z głowy rdzenia".
//  2. Sterowanie własnym zleceniem jest tą czynnością, której droga przez model
//     jest wprost szkodliwa: „anuluj zlecenie" powiedziane do modelu ląduje jako
//     treść tury tego właśnie zlecenia, więc zamiast je przerwać — przedłuża.
//
// Wykaz jest startowy i wołający ma prawo podać własny. To, które frazy platforma
// rozumie, nie należy do tego pliku; wnosi on mechanizm i pięć fraz, przy
// których mechanizm daje się sprawdzić na żywym rdzeniu.
//
// Wynik składany jest przy każdym wywołaniu, bo niesie mapy — wspólna kopia
// pozwoliłaby wołającemu zmienić leksykon wszystkim naraz, nie chcąc tego.
func AkcjeDomyslne() []AkcjaPlatformy {
	return []AkcjaPlatformy{
		{
			Frazy:   []string{"wstrzymaj zlecenie", "pauza zlecenia"},
			Komenda: komendaStanZlecenia,
			Nastawy: map[string]string{"control": "pause"},
			Wymaga:  []string{"windowId"},
			Opis:    "wstrzymuje zlecenie asystenta w tym oknie",
		},
		{
			Frazy:   []string{"wznów zlecenie", "wznów"},
			Komenda: komendaStanZlecenia,
			Nastawy: map[string]string{"control": "resume"},
			Wymaga:  []string{"windowId"},
			Opis:    "wznawia wstrzymane zlecenie asystenta w tym oknie",
		},
		{
			Frazy:   []string{"anuluj zlecenie", "przerwij zlecenie", "stop"},
			Komenda: komendaStanZlecenia,
			Nastawy: map[string]string{"control": "cancel"},
			Wymaga:  []string{"windowId"},
			Opis:    "anuluje zlecenie asystenta w tym oknie",
		},
		{
			Frazy:   []string{"ponów zlecenie", "spróbuj jeszcze raz"},
			Komenda: komendaStanZlecenia,
			Nastawy: map[string]string{"control": "retry"},
			Wymaga:  []string{"windowId"},
			Opis:    "ponawia zlecenie asystenta w tym oknie",
		},
		{
			Frazy:   []string{"stan zlecenia", "co z moim zleceniem"},
			Komenda: komendaStanZlecenia,
			Nastawy: map[string]string{"control": "none"},
			Wymaga:  []string{"windowId"},
			Opis:    "odczytuje stan zleceń asystenta w tym oknie, bez sterowania",
		},
	}
}
