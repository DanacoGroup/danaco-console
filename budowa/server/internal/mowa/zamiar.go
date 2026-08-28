// Odpowiedzialność pliku: rozróżnienie, czym jest wypowiedź Operatora —
// poleceniem dla modelu czy akcją platformy; plik rozstrzyga i nazywa
// podstawę rozstrzygnięcia, nie wykonuje ani jednego, ani drugiego.
package mowa

import (
	"sort"
	"strconv"
	"strings"
	"unicode"
)

// Zamiar mówi, czym wypowiedź jest dla rdzenia. Wartość zerowa jest
// nierozstrzygnięciem: struktura powstała przez pomyłkę ma znaczyć
// „nikt niczego nie zlecił", a nie „jedź do modelu".
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
// kontraktu w polu `intent`. Napisy są te same co wartości `Zamiar`, bo pole
// kontraktu ma nieść rozstrzygnięcie Operatora w tym samym słowniku, co rdzeń.
const (
	// WskazanieBrak — żądanie pola "intent" nie niesie wartości, bo pole nie
	// istnieje w kontrakcie albo wołający go nie wypełnił; rozstrzygnięcie
	// traktuje to tak samo jak wskazanie akcji platformy.
	WskazanieBrak = ""
	// WskazanieModelu — Operator zadeklarował wprost w polu żądania, że
	// wypowiedź jest poleceniem dla modelu; rozstrzygnięcie bierze tę deklarację
	// jako rozstrzygającą, choćby treść wypowiedzi brzmiała jak fraza akcji.
	WskazanieModelu = "model"
	// WskazaniePlatformy — Operator zadeklarował wprost w polu żądania, że
	// wypowiedź jest akcją platformy; rozstrzygnięcie szuka wtedy dopasowania
	// frazy w leksykonie akcji i przy jego braku zwraca wynik nierozstrzygnięty.
	WskazaniePlatformy = "platform"
)

// Podstawa nazywa, czym rozstrzygnięto — bo „platforma" bez odpowiedzi na
// pytanie „skąd wiesz" jest nie do zdiagnozowania, gdy Operator zgłosi, że
// rdzeń zrobił nie to, co powiedział.
type Podstawa string

const (
	// PodstawaBrak — nie rozstrzygnięto; wartość zerowa pola Podstawa towarzyszy
	// zamiarowi nierozstrzygniętemu i nie ma innego znaczenia w wyniku funkcji
	// Rozstrzygnij.
	PodstawaBrak Podstawa = ""
	// PodstawaWskazania — rozstrzygnęło pole żądania "intent"; Operator
	// zadeklarował wprost, czy wypowiedź jest poleceniem dla modelu, i ta
	// deklaracja przesądziła wynik.
	PodstawaWskazania Podstawa = "wskazanie"
	// PodstawaFrazy — rozstrzygnęła cała wypowiedź równa frazie leksykonu akcji
	// platformy; dopasowanie poszło po całym tekście, znak w znak po
	// normalizacji, a nie po fragmencie.
	PodstawaFrazy Podstawa = "fraza"
	// PodstawaDomyslna — nic nie wskazało akcji platformy, więc jedzie model; to
	// droga dotychczasowa, brana wtedy, gdy brak deklaracji Operatora i brak
	// dopasowanej frazy.
	PodstawaDomyslna Podstawa = "domyslnie"
)

// AkcjaPlatformy wiąże frazy Operatora z jedną komendą kontraktu. Nastawy są
// dosłowne, nie wyliczane: fraza niesie wartość dlatego, że tak zapisano
// w leksykonie, a nie dlatego, że coś przeczytało słowo z wypowiedzi.
type AkcjaPlatformy struct {
	// Frazy — całe wypowiedzi (nie fragmenty) uruchamiające tę akcję.
	Frazy []string
	// Komenda — nazwa komendy kontraktu, na przykład `assistant.action.status`.
	Komenda string
	// Nastawy — pola żądania o wartościach dosłownych, wpisanych w leksykonie.
	Nastawy map[string]string
	// Wymaga — pola żądania spoza leksykonu; wołający dokłada je z kontekstu.
	Wymaga []string
	// Opis — zdanie dla Operatora, czym ta akcja jest.
	Opis string
}

// Wypowiedz jest wejściem rozstrzygnięcia: to, co przyszło komendą głosową —
// sama treść wypowiedzi razem z deklaracją Operatora z pola żądania, jeśli
// taka deklaracja została przysłana.
type Wypowiedz struct {
	// Tekst — transkrypcja (własna albo poprawiona przez Operatora).
	Tekst string
	// Wskazanie — deklaracja Operatora z pola żądania; puste znaczy brak pola.
	Wskazanie string
}

// Rozstrzygniecie jest wynikiem rozróżnienia: niesie zamiar, akcję platformy
// przy zamiarze platformowym, podstawę rozstrzygnięcia oraz powód przy
// zamiarze nierozstrzygniętym.
type Rozstrzygniecie struct {
	// Zamiar — model, platforma albo nierozstrzygnięty (wartość zerowa).
	Zamiar Zamiar
	// Akcja — wypełniona wyłącznie przy `ZamiarPlatformy`; dla modelu pusta.
	Akcja AkcjaPlatformy
	// Podstawa — czym rozstrzygnięto.
	Podstawa Podstawa
	// Powod — wypełniony wyłącznie przy zamiarze nierozstrzygniętym.
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

// Rozstrzygnij rozróżnia polecenie dla modelu od akcji platformy. Deklaracja
// Operatora bije leksykon: gdy wypowiedź brzmi jak fraza akcji, a Operator
// napisał `intent=model`, wygrywa `intent`, bo mógł zlecić modelowi wyjaśnienie
// zdania.
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
		// Dwie komendy na jedną frazę to usterka leksykonu, nie wybór między nimi.
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

// nierozstrzygniety składa wynik bez zamiaru z nazwanym powodem; korzystają
// z niej wszystkie gałęzie funkcji Rozstrzygnij, które nie potrafią wskazać
// ani modelu, ani akcji platformy.
func nierozstrzygniety(powod string) Rozstrzygniecie {
	return Rozstrzygniecie{Zamiar: ZamiarNierozstrzygniety, Podstawa: PodstawaBrak, Powod: powod}
}

// Wywolanie jest gotowym żądaniem akcji platformy: nazwa komendy kontraktu
// i komplet pól, z jakimi ma pojechać. Typ jest opisem, nie wykonaniem: pakiet
// mowy nie zna rdzenia i tylko składa wywołanie, które oddaje temu, kto
// komendy wykonuje.
type Wywolanie struct {
	// Komenda — nazwa komendy kontraktu.
	Komenda string
	// Pola — nastawy leksykonu ORAZ pola dołożone z kontekstu wywołania.
	Pola map[string]string
}

// BrakDanychAkcji jest odmową: akcja jest rozpoznana, lecz nie ma z czym
// pojechać, bo wołający nie dołożył pola, którego leksykon nie zna.
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
// wypowiedzi, która akcją platformy nie jest, żeby pomyłka wołającego
// kończyła się błędem, a nie pustym wywołaniem, które dałoby się wysłać.
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

// Unwrap kończy łańcuch — jak w pozostałych odmowach pakietu, NieAkcjaPlatformy
// nie opakowuje cudzego błędu, więc rozwijanie łańcucha błędów zatrzymuje się
// na niej.
func (b *NieAkcjaPlatformy) Unwrap() error { return nil }

// ZlozWywolanie składa żądanie akcji platformy z nastaw leksykonu i pól,
// które wołający bierze z kontekstu komendy głosowej. Kontekst nie
// nadpisuje leksykonu: wołający dokłada tożsamości, a nie sterowanie.
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

// dopasujAkcje szuka akcji, których fraza jest równa całej wypowiedzi, nie
// zawarta w niej, bo dopasowanie po fragmencie zamieniłoby wzmiankę w
// wykonanie. Zwracane są wszystkie dopasowania, także sprzeczne.
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

// bezOgonkow zdejmuje polskie znaki diakrytyczne. Potrzebne, bo wejście jest
// z mowy, a nie z klawiatury: pomocnik transkrypcji oddaje raz „wznów", raz
// „wznow", zależnie od modelu i od tego, jak wyraźnie Operator mówił.
var bezOgonkow = map[rune]rune{
	'ą': 'a', 'ć': 'c', 'ę': 'e', 'ł': 'l', 'ń': 'n',
	'ó': 'o', 'ś': 's', 'ź': 'z', 'ż': 'z',
}

// znormalizuj sprowadza wypowiedź do postaci porównywalnej: małe litery, bez
// ogonków, bez znaków przestankowych, pojedyncze odstępy, bez brzegów.
// Normalizacja nie jest domyślaniem się treści: nie usuwa słów i nie zmienia
// ich kolejności.
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

// Nazwy komend leksykonu domyślnego. Napisy, a nie stałe z pakietu `shared`,
// bo pakiet mowy kontraktu nie zna, tak jak opisuje to nagłówek pliku
// na początku pakietu.
const komendaStanZlecenia = "assistant.action.status"

// AkcjeDomyslne oddaje leksykon startowy akcji platformy: pięć fraz prowadzi
// do jednej komendy assistant.action.status, bo sterowanie zleceniem drogą
// przez model jest wprost szkodliwe. Wynik jest składany przy każdym
// wywołaniu, bo niesie mapy.
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
