// Brama kontraktu sprawdza żądanie wobec kontraktu przed czynnością domeny: pola wymagane i wartości wyliczeń, czytane z artefaktu kontraktu; powitanie kanału jest jedynym wyjątkiem.
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

// znacznikPolaKontraktu nazywa znacznik struktury Go, spod którego brama odczytuje nazwę pola żądania odpowiadającą polu kontraktu.
const znacznikPolaKontraktu = "json"

// wskazaniePolaOpcjonalnego to człon znacznika, którym generator kontraktu oznacza pole opcjonalne; pole wymagane niesie znacznik bez tego członu, wprost z oznaczenia „wymagane” kontraktu.
const wskazaniePolaOpcjonalnego = "omitempty"

// komendaPozaBrama nazywa jedyną komendę, której brama nie sprawdza. Stała
// zamiast warunku wpisanego w gałąź, bo wyjątek ma być odczytywalny w jednym
// miejscu — wyjątek rozsypany po warunkach przestaje być jedynym.
const komendaPozaBrama = shared.CommandConnectionHello

// sprawdzZadanieWobecKontraktu odmawia żądaniu niezgodnemu z kontraktem: sprawdza obecność pól wymaganych oraz przynależność wartości do wyliczenia kontraktu, oba wyczytane z kontraktu.
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

// odnotujBrakiPowitania kładzie braki powitania w dzienniku rdzenia, jedynym miejscu, do którego mają dojść, ponieważ odpowiedź powitania pola na braki nie niesie.
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

// polaTresci rozkłada treść żądania na pola wierzchnie. Drugi wynik odróżnia treść nierozkładalną na pola od treści rozłożonej, lecz pustej.
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
		// Obecność klucza jest sprawdzana, nie jego treść: `null` też jest obecnością i przechodzi.
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

// wartosciSpozaWyliczen zwraca opis pierwszej wartości spoza wyliczenia kontraktu, albo pusty napis, gdy wszystkie wartości do wyliczenia należą.
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

// wartoscPozaZakresem sprawdza jedną wartość, napis albo tablicę napisów, wobec zakresu wyliczenia kontraktu; wartość pusta oraz wartość `null` przechodzą bez sprawdzenia.
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

// naliscie odpowiada, czy podana wartość stoi w zakresie dopuszczalnych wartości wyliczenia kontraktu przekazanym jako parametr.
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

// wyliczeniaKomendy oddaje odwzorowanie pól na dopuszczalne wartości wyliczeń, właściwe jednej wskazanej komendzie kontraktu.
func wyliczeniaKomendy(komenda shared.MessageType) map[string][]string {
	return wykazWyliczenKomend()[komenda]
}

// bladZgodnosciZKontraktem nazywa niezgodność żądania z kontraktem. Komenda
// stoi w treści odmowy, bo odmowa idzie do Operatora oderwana od żądania.
func bladZgodnosciZKontraktem(komenda shared.MessageType, powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		string(komenda)+": "+powod))
}
