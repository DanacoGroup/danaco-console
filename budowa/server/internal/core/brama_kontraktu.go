// Odpowiedzialność pliku: brama kontraktu — sprawdzenie treści żądania wobec
// kontraktu, zanim wejdzie ona do czynności domeny.
//
// ── Dlaczego brama stoi tu, a nie w każdym obsługiwaczu z osobna ────────────
// Kontrakt niesie dla każdej komendy wykaz pól wraz z oznaczeniem `wymagane`,
// a dla pól o typie wyliczenia — komplet dopuszczalnych wartości. Rdzeń bez
// bramy przyjmował żądanie niepełne i wartość spoza wyliczenia, po czym
// uzupełniał brak wartością domyślną i meldował powodzenie. Odpowiedź
// wyglądająca dobrze jest w tym miejscu groźniejsza od odmowy, bo nie wzywa
// wołającego do sprawdzenia: okno dostawało punkt dostępu rodzaju, o który nie
// prosiło, i tryb uprawnień, którego nie ustawiło.
//
// Brama nie ma własnego wykazu pól ani własnego wykazu wartości — oba czyta
// z artefaktu kontraktu (`shared/contract.go`), więc pole dołożone do kontraktu
// jest pilnowane od razu, bez zmiany w tym pliku.
//
// ── Jedyna komenda spod bramy wyjęta ────────────────────────────────────────
// Powitanie kanału (`connection.hello`) bramie nie podlega i odpowiada zawsze,
// także na żądanie niepełne — rejestr decyzji, pozycja 10. Powitanie jest
// jedynym miejscem, w którym klient odczytuje `protocolVersion` rdzenia, czyli
// jedynym, w którym rozpoznaje, że jest starszy. Brama sprawdzająca je wobec
// kontraktu zakłada, że obie strony znają już ten sam kontrakt — zakłada więc
// to, co powitanie ma dopiero ustalić, i klientowi sprzed wprowadzenia pola
// oddaje odmowę zamiast wersji, po której ten rozpoznałby rozjazd.
//
// Braki pól powitania idą do dziennika rdzenia, nie do treści odpowiedzi:
// `ConnectionHelloResponse` nie ma pola, w którym mogłyby wrócić wołającemu,
// a dołożenie takiego pola jest zmianą kontraktu. Celowi wyjątku to wystarcza —
// klient starszy ma odczytać wersję protokołu i sam rozpoznać rozjazd, a do
// tego potrzebuje odpowiedzi, nie wykazu swoich braków.
//
// Wyjątek jest jeden i pozostaje jeden. Wynika z roli powitania w uzgodnieniu,
// nie z wygody, więc każda inna komenda przechodzi bramę bez ustępstw.
//
// ── Czego brama NIE obejmuje ────────────────────────────────────────────────
// Obecność pól wymaganych sprawdzana jest dla każdej komendy: oznaczenie
// `wymagane` niesie znacznik struktury żądania, a struktury ma każda komenda.
// Wartości wyliczeń sprawdzane są węziej — komplet dopuszczalnych wartości
// stoi w artefakcie Go wyłącznie przy polach komend wystawionych jako narzędzia
// modelu (`shared.NarzedziaModelu`). Wyliczenie samo w sobie ma w artefakcie
// funkcję `Wartosci<Nazwa>`, ale nie ma odwzorowania „typ pola → ta funkcja",
// więc pole wyliczeniowe komendy spoza tego zbioru przechodzi bez sprawdzenia
// wartości. Pełne pokrycie wymaga tabeli wyprowadzonej z kontraktu przy
// generowaniu artefaktu, a nie przepisanej tutaj — drugi wykaz wartości
// rozjechałby się z kontraktem przy pierwszej dołożonej wartości.
package core

import (
	"context"
	"encoding/json"
	"reflect"
	"sort"
	"strings"
	"sync"

	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// znacznikPolaKontraktu nazywa znacznik struktury niosący nazwę pola kontraktu.
const znacznikPolaKontraktu = "json"

// wskazaniePolaOpcjonalnego to człon znacznika, którym generator kontraktu
// oznacza pole opcjonalne. Pole wymagane znacznika bez tego członu — to jest
// jedyne rozróżnienie „wymagane / opcjonalne" po stronie Go i pochodzi wprost
// z oznaczenia `wymagane` w `contract.json`.
const wskazaniePolaOpcjonalnego = "omitempty"

// komendaPozaBrama nazywa jedyną komendę, której brama nie sprawdza. Stała
// zamiast warunku wpisanego w gałąź, bo wyjątek ma być odczytywalny w jednym
// miejscu — wyjątek rozsypany po warunkach przestaje być jedynym.
const komendaPozaBrama = shared.CommandConnectionHello

// sprawdzZadanieWobecKontraktu odmawia żądaniu niezgodnemu z kontraktem.
//
// Sprawdzane są dwie rzeczy, obie wyczytane z kontraktu:
//
//  1. Obecność pól wymaganych. Sprawdzana jest OBECNOŚĆ klucza w treści, nie
//     jego zawartość: kontrakt mówi „pole ma być", nie „pole ma być niepuste".
//     Treść pusta bywa treścią prawdziwą — `studio.document.save` z pustym
//     `content` zapisuje dokument opróżniony i jest żądaniem poprawnym. Wartość
//     pustą, tam gdzie dziedzina jej nie zniesie, odrzuca obsługiwacz komendy,
//     bo tylko on wie, czy pustka coś znaczy.
//
//  2. Przynależność wartości do wyliczenia kontraktu — dla tych pól, dla których
//     kontrakt niesie komplet wartości w deklaracji narzędzia modelu.
//
// Żądanie zgodne przechodzi bez śladu; niezgodne wraca kodem `validation_failed`
// z treścią nazywającą brak albo wartość spoza zakresu.
//
// Powitanie kanału przechodzi zawsze — jest jedyną komendą spod bramy wyjętą,
// a jego braki idą do dziennika rdzenia (nagłówek pliku).
func sprawdzZadanieWobecKontraktu(ctx context.Context, komenda shared.MessageType,
	ladunek json.RawMessage, wzor any) error {

	pola, sa := polaTresci(ladunek)
	if !sa {
		return nil
	}
	brakujace := brakujacePolaWymagane(wzor, pola)
	if komenda == komendaPozaBrama {
		odnotujBrakiPowitania(ctx, komenda, brakujace)
		return nil
	}
	if len(brakujace) > 0 {
		return bladZgodnosciZKontraktem(komenda,
			"żądanie bez pól wymaganych kontraktem: "+strings.Join(brakujace, ", "))
	}
	if powod := wartosciSpozaWyliczen(komenda, pola); powod != "" {
		return bladZgodnosciZKontraktem(komenda, powod)
	}
	return nil
}

// odnotujBrakiPowitania kładzie braki powitania w dzienniku rdzenia — jedynym
// miejscu, do którego mają dojść.
//
// Zapis idzie wyłącznie przy brakach: powitanie pełne jest przypadkiem zwykłym
// i wpis o nim zasypywałby dziennik przy każdym nawiązaniu połączenia. Brak
// dziennika nie zmienia zachowania rdzenia — znika sam zapis, a powitanie
// odpowiada tak samo.
func odnotujBrakiPowitania(ctx context.Context, komenda shared.MessageType, brakujace []string) {
	if len(brakujace) == 0 {
		return
	}
	dziennik := dziennikZKontekstu(ctx)
	if dziennik == nil {
		return
	}
	dziennik.Printf("brama kontraktu: %s bez pól wymaganych kontraktem: %s; "+
		"powitanie odpowiada mimo to, bo klient odczytuje z niego wersję protokołu",
		komenda, strings.Join(brakujace, ", "))
}

// polaTresci rozkłada treść żądania na pola wierzchnie. Drugi wynik odróżnia
// „treści nie da się rozłożyć na pola" od „treść nie ma ani jednego pola":
// pierwsze nie jest przedmiotem tej bramy (odczyt ładunku w kształcie kontraktu
// odmówił już wcześniej albo komenda ładunku nie ma), drugie jest żądaniem
// pustym i podlega sprawdzeniu jak każde inne.
func polaTresci(ladunek json.RawMessage) (map[string]json.RawMessage, bool) {
	if len(ladunek) == 0 {
		return map[string]json.RawMessage{}, true
	}
	var pola map[string]json.RawMessage
	if err := json.Unmarshal(ladunek, &pola); err != nil {
		return nil, false
	}
	if pola == nil {
		return map[string]json.RawMessage{}, true
	}
	return pola, true
}

// brakujacePolaWymagane zwraca nazwy pól, które kontrakt oznacza jako wymagane,
// a których treść żądania nie niesie. Nazwy idą w kolejności kontraktu, bo
// w tej kolejności czyta je człowiek w opisie komendy.
func brakujacePolaWymagane(wzor any, pola map[string]json.RawMessage) []string {
	ksztalt := reflect.TypeOf(wzor)
	for ksztalt != nil && ksztalt.Kind() == reflect.Pointer {
		ksztalt = ksztalt.Elem()
	}
	if ksztalt == nil || ksztalt.Kind() != reflect.Struct {
		return nil
	}
	var brakujace []string
	for i := range ksztalt.NumField() {
		nazwa, wymagane := polePola(ksztalt.Field(i))
		if !wymagane {
			continue
		}
		// Sprawdzana jest OBECNOŚĆ klucza. `null` obecnością jest i przechodzi:
		// pole wymagane o typie tablicy wychodzi z niepustego wykazu pustego
		// jako `null` (tak koduje pusty wycinek biblioteka standardowa Go), więc
		// odmowa w tym miejscu odrzucałaby żądania składane przez sam rdzeń.
		// O tym, czy `null` niesie treść, rozstrzyga obsługiwacz komendy.
		if _, jest := pola[nazwa]; !jest {
			brakujace = append(brakujace, nazwa)
		}
	}
	return brakujace
}

// polePola odczytuje ze znacznika struktury nazwę pola kontraktu oraz to, czy
// pole jest wymagane. Pole nieeksportowane i pole wyłączone znacznikiem „-"
// nie należą do treści kontraktu.
func polePola(pole reflect.StructField) (string, bool) {
	if !pole.IsExported() {
		return "", false
	}
	znacznik, jest := pole.Tag.Lookup(znacznikPolaKontraktu)
	if !jest {
		return "", false
	}
	nazwa, reszta, _ := strings.Cut(znacznik, ",")
	if nazwa == "" || nazwa == "-" {
		return "", false
	}
	for _, czlon := range strings.Split(reszta, ",") {
		if czlon == wskazaniePolaOpcjonalnego {
			return nazwa, false
		}
	}
	return nazwa, true
}

// wartosciSpozaWyliczen zwraca opis pierwszej wartości, która nie należy do
// wyliczenia kontraktu, albo pusty napis, gdy wszystkie należą.
//
// Zakres sprawdzenia wyznacza to, co kontrakt udostępnia rdzeniowi w czasie
// pracy: komplet wartości wyliczenia stoi przy polach komend wystawionych jako
// narzędzia modelu (`shared.NarzedziaModelu`). Pola pozostałych komend brama
// przepuszcza bez sprawdzenia wartości — nie dlatego, że wolno im mieć wartość
// dowolną, lecz dlatego, że artefakt kontraktu nie niesie dla nich wykazu
// wartości w postaci czytelnej w czasie pracy rdzenia.
func wartosciSpozaWyliczen(komenda shared.MessageType, pola map[string]json.RawMessage) string {
	dopuszczalne := wyliczeniaKomendy(komenda)
	if len(dopuszczalne) == 0 {
		return ""
	}
	nazwy := make([]string, 0, len(pola))
	for nazwa := range pola {
		nazwy = append(nazwy, nazwa)
	}
	sort.Strings(nazwy)
	for _, nazwa := range nazwy {
		zakres, jest := dopuszczalne[nazwa]
		if !jest {
			continue
		}
		if powod := wartoscPozaZakresem(nazwa, pola[nazwa], zakres); powod != "" {
			return powod
		}
	}
	return ""
}

// wartoscPozaZakresem sprawdza jedną wartość — napis albo tablicę napisów —
// wobec zakresu wyliczenia. Wartość `null` i wartość innego kształtu nie są
// przedmiotem tego sprawdzenia: pierwszą rozstrzyga sprawdzenie pól wymaganych,
// drugą odczyt ładunku w kształcie kontraktu.
//
// Napis pusty przechodzi, tak samo jak przy sprawdzeniu pól wymaganych: pusto
// znaczy „wartości nie podano", a brama orzeka o obecności pola, nie o jego
// zawartości. Rozstrzygnięcie, czy brak wyboru jest w danej komendzie dopuszczalny,
// należy do jej obsługiwacza — to on wie, czy rodzaj nienazwany coś znaczy.
func wartoscPozaZakresem(nazwa string, wartosc json.RawMessage, zakres []string) string {
	var napis string
	if err := json.Unmarshal(wartosc, &napis); err == nil {
		if napis == "" || naliscie(napis, zakres) {
			return ""
		}
		return "pole " + nazwa + ": wartość " + napis + " nie należy do wyliczenia kontraktu (" +
			strings.Join(zakres, ", ") + ")"
	}
	var lista []string
	if err := json.Unmarshal(wartosc, &lista); err != nil {
		return ""
	}
	for _, element := range lista {
		if element == "" || naliscie(element, zakres) {
			continue
		}
		return "pole " + nazwa + ": wartość " + element + " nie należy do wyliczenia kontraktu (" +
			strings.Join(zakres, ", ") + ")"
	}
	return ""
}

// naliscie odpowiada, czy wartość stoi w zakresie wyliczenia.
func naliscie(wartosc string, zakres []string) bool {
	for _, dopuszczalna := range zakres {
		if wartosc == dopuszczalna {
			return true
		}
	}
	return false
}

// wykazWyliczenKomend składa raz odwzorowanie „komenda → pole → dopuszczalne
// wartości" z deklaracji narzędzi modelu. Składanie przy każdym żądaniu
// przeglądałoby trzysta trzydzieści dziewięć deklaracji dla jednej komendy.
var wykazWyliczenKomend = sync.OnceValue(func() map[shared.MessageType]map[string][]string {
	wykaz := make(map[shared.MessageType]map[string][]string)
	for _, narzedzie := range shared.NarzedziaModelu() {
		for _, parametr := range narzedzie.Parameters {
			if len(parametr.Enum) == 0 {
				continue
			}
			pola, jest := wykaz[narzedzie.Command]
			if !jest {
				pola = make(map[string][]string)
				wykaz[narzedzie.Command] = pola
			}
			pola[parametr.Name] = parametr.Enum
		}
	}
	return wykaz
})

// wyliczeniaKomendy oddaje dopuszczalne wartości pól jednej komendy.
func wyliczeniaKomendy(komenda shared.MessageType) map[string][]string {
	return wykazWyliczenKomend()[komenda]
}

// bladZgodnosciZKontraktem nazywa niezgodność żądania z kontraktem. Komenda
// stoi w treści odmowy, bo odmowa idzie do Operatora oderwana od żądania.
func bladZgodnosciZKontraktem(komenda shared.MessageType, powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		string(komenda)+": "+powod))
}
