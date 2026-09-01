// Brama kontraktu sprawdza żądanie wobec kontraktu przed czynnością domeny: pola wymagane i wartości wyliczeń, czytane z artefaktu kontraktu; powitanie kanału jest jedynym wyjątkiem.
package core

import (
	"bytes"
	"context"
	"encoding/json"
	"reflect"
	"strings"

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
			"żądanie bez wartości w polach wymaganych kontraktem: "+strings.Join(brakujace, ", "))
	}
	return niezgodnoscWyliczen(komenda, wzor, pola)
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

// ksztaltZadania oddaje kształt struktury treści żądania. Wzór przychodzi
// z obsługiwacza komendy, więc bywa wskaźnikiem; typ niebędący strukturą pól
// kontraktu nie niesie i sprawdzać w nim nie ma czego.
func ksztaltZadania(wzor any) reflect.Type {
	ksztalt := reflect.TypeOf(wzor)
	for ksztalt != nil && ksztalt.Kind() == reflect.Pointer {
		ksztalt = ksztalt.Elem()
	}
	if ksztalt == nil || ksztalt.Kind() != reflect.Struct {
		return nil
	}
	return ksztalt
}

// brakujacePolaWymagane zwraca nazwy pól, które kontrakt oznacza jako wymagane,
// a których treść żądania nie niesie albo niesie bez wartości. Nazwy idą
// w kolejności kontraktu, bo w tej kolejności czyta je człowiek w opisie
// komendy.
func brakujacePolaWymagane(wzor any, pola map[string]json.RawMessage) []string {
	ksztalt := ksztaltZadania(wzor)
	if ksztalt == nil {
		return nil
	}
	var brakujace []string
	for i := range ksztalt.NumField() {
		pole := ksztalt.Field(i)
		nazwa, wymagane := polePola(pole)
		if !wymagane {
			continue
		}
		surowa, jest := pola[nazwa]
		if !jest || !poleNiesieWartosc(pole.Type, surowa) {
			brakujace = append(brakujace, nazwa)
		}
	}
	return brakujace
}

// poleNiesieWartosc odróżnia wartość pola od samej obecności klucza. `null`
// nie jest wartością pola napisowego, a pusty napis nie jest wartością pola
// wyliczeniowego, bo żadne wyliczenie kontraktu pustej wartości nie zna —
// przepuszczone dochodzą do zapisu i zatrzymuje je dopiero CHECK sterownika
// bazy, wracając surowym tekstem sterownika zamiast odmową walidacji. Pusty
// napis w polu swobodnym zostaje wartością, bo kontrakt zna komendy, w których
// znaczy zdjęcie wskazania. Pole wielowartościowe rozstrzyga się samą
// obecnością, bo wycinek pusty jest wartością, a Go zapisuje go jako `null`.
func poleNiesieWartosc(typ reflect.Type, surowa json.RawMessage) bool {
	if typ.Kind() != reflect.String {
		return true
	}
	if bytes.Equal(bytes.TrimSpace(surowa), []byte("null")) {
		return false
	}
	var napis string
	if err := json.Unmarshal(surowa, &napis); err != nil {
		return true
	}
	return napis != "" || nazwaWyliczenia(typ) == ""
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

// niezgodnoscWyliczen odmawia żądaniu, którego pole wyliczeniowe niesie wartość spoza zakresu kontraktu; pola pozostałych typów zakresu nie mają i przechodzą.
func niezgodnoscWyliczen(komenda shared.MessageType, wzor any, pola map[string]json.RawMessage) error {
	ksztalt := ksztaltZadania(wzor)
	if ksztalt == nil {
		return nil
	}
	for i := range ksztalt.NumField() {
		pole := ksztalt.Field(i)
		nazwa, _ := polePola(pole)
		if nazwa == "" {
			continue
		}
		wyliczenie := nazwaWyliczenia(pole.Type)
		if wyliczenie == "" {
			continue
		}
		wartosc, jest := pola[nazwa]
		if !jest {
			continue
		}
		zakres, znany := shared.ZakresyWyliczen[wyliczenie]
		if !znany {
			return bladWyliczeniaBezZakresu(komenda, nazwa, wyliczenie)
		}
		if powod := wartoscPozaZakresem(nazwa, wartosc, zakres); powod != "" {
			return bladZgodnosciZKontraktem(komenda, powod)
		}
	}
	return nil
}

// nazwaWyliczenia oddaje nazwę wyliczenia kontraktu stojącego pod typem pola,
// albo pusty napis, gdy pole wyliczenia nie niesie. Generator kontraktu
// wystawia wyliczenie nazwanym typem napisowym, pole opcjonalne wskaźnikiem,
// a pole wielowartościowe wycinkiem — stąd zdejmowanie obu powłok. Typ
// napisowy bez własnego pakietu to zwykły `string`, czyli pole bez wyliczenia.
func nazwaWyliczenia(typ reflect.Type) string {
	for typ.Kind() == reflect.Pointer || typ.Kind() == reflect.Slice {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.String || typ.PkgPath() == "" {
		return ""
	}
	return typ.Name()
}

// wartoscPozaZakresem sprawdza jedną wartość, napis albo tablicę napisów, wobec zakresu wyliczenia kontraktu; wartość pusta oraz wartość `null` przechodzą bez sprawdzenia, bo dla pola opcjonalnego znaczą brak wskazania.
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

// bladZgodnosciZKontraktem nazywa niezgodność żądania z kontraktem. Komenda
// stoi w treści odmowy, bo odmowa idzie do Operatora oderwana od żądania.
func bladZgodnosciZKontraktem(komenda shared.MessageType, powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed,
		string(komenda)+": "+powod))
}

// bladWyliczeniaBezZakresu nazywa rozjazd typu pola z kontraktem: pole niesie
// nazwany typ napisowy, którego wykaz ZakresyWyliczen kontraktu nie zna, więc
// zakresu nie ma czym sprawdzić. Odmowa zamiast przepuszczenia, bo przepuszczone
// pole rozstrzyga dopiero CHECK sterownika bazy.
func bladWyliczeniaBezZakresu(komenda shared.MessageType, pole, wyliczenie string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeInternalError,
		string(komenda)+": pole "+pole+": wyliczenie "+wyliczenie+
			" bez zakresu w wykazie ZakresyWyliczen kontraktu"))
}
