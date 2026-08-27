// Ten plik utrzymuje granicę uprawnień modelu: poza wykazem narzędzi stoją
// komendy warstwy połączenia klienta i komendy zastrzeżone Operatorowi,
// których model nie rozszerza.
package narzedzia

import (
	"fmt"
	"strings"
	"sync"

	"danacoconsole/shared"
)

// komendyPozaWykazem odwzorowuje nazwę, jaką nosiłoby narzędzie komendy spoza
// wykazu, na samą komendę. Liczone raz na proces: kontrakt jest niezmienny
// w czasie życia procesu.
var komendyPozaWykazem = sync.OnceValue(zbierzKomendyPozaWykazem)

// odmowaNarzedzia buduje treść odmowy dla nazwy, której wykaz kontraktu nie
// niesie. Zwraca błąd, bo rozdzielnia oddaje go modelowi jako treść błędu
// narzędzia — połączenie żyje dalej.
func odmowaNarzedzia(nazwa string) error {
	komenda, zastrzezona := komendyPozaWykazem()[nazwa]
	if !zastrzezona {
		return fmt.Errorf("narzędzie %q nie występuje w wykazie kontraktu; komplet narzędzi zwraca tools/list", nazwa)
	}
	return fmt.Errorf("narzędzie %q nie jest udostępnione modelowi: komenda %s stoi poza wykazem narzędzi kontraktu. "+
		"Poza wykazem są warstwa połączenia klienta oraz punkty zastrzeżone Operatorowi — zakładanie "+
		"i zmiana punktów dostępu, nadań, kont oraz treści tożsamości. Odczyt tych rejestrów model ma, zapisu nie",
		nazwa, komenda)
}

// zbierzKomendyPozaWykazem odejmuje od kompletu komend kontraktu komendy objęte
// narzędziami i nadaje reszcie nazwy według wzorca nazewniczego kontraktu.
func zbierzKomendyPozaWykazem() map[string]shared.MessageType {
	przedrostek, separator, rozpoznany := wzorzecNazwy()
	if !rozpoznany {
		// Wzorca nie da się odczytać wyłącznie wtedy, gdy kontrakt nie niesie ani
		// jednego narzędzia.
		return map[string]shared.MessageType{}
	}
	objete := map[shared.MessageType]bool{}
	for _, komenda := range shared.KomendyNarzedzi {
		objete[komenda] = true
	}
	poza := map[string]shared.MessageType{}
	for _, komenda := range shared.WszystkieKomendy() {
		if objete[komenda] {
			continue
		}
		poza[nazwaNarzedzia(komenda, przedrostek, separator)] = komenda
	}
	return poza
}

// nazwaNarzedzia składa nazwę narzędzia z nazwy komendy: `access.point.add`
// przy przedrostku `danaco` i separatorze `_` daje `danaco_access_point_add`.
func nazwaNarzedzia(komenda shared.MessageType, przedrostek, separator string) string {
	czlony := append([]string{przedrostek}, strings.Split(string(komenda), ".")...)
	return strings.Join(czlony, separator)
}

// wzorzecNazwy odczytuje przedrostek i separator nazw narzędzi z samych
// deklaracji kontraktu, nie z literału zapisanego tutaj.
func wzorzecNazwy() (przedrostek, separator string, rozpoznany bool) {
	deklaracje := shared.NarzedziaModelu()
	if len(deklaracje) == 0 {
		return "", "", false
	}
	przedrostek, separator, rozpoznany = rozbierzNazwe(deklaracje[0])
	if !rozpoznany {
		return "", "", false
	}
	for _, pozycja := range deklaracje {
		if nazwaNarzedzia(pozycja.Command, przedrostek, separator) != pozycja.Name {
			return "", "", false
		}
	}
	return przedrostek, separator, true
}

// rozbierzNazwe wyprowadza przedrostek i separator z jednej deklaracji. Separator
// jest znakiem stojącym w nazwie narzędzia tuż przed ostatnim członem komendy.
func rozbierzNazwe(pozycja shared.ToolDeclaration) (przedrostek, separator string, rozpoznany bool) {
	czlony := strings.Split(string(pozycja.Command), ".")
	ostatni := czlony[len(czlony)-1]
	odciecie := len(pozycja.Name) - len(ostatni) - 1
	if ostatni == "" || odciecie <= 0 || !strings.HasSuffix(pozycja.Name, ostatni) {
		return "", "", false
	}
	separator = pozycja.Name[odciecie : odciecie+1]
	ogon := separator + strings.Join(czlony, separator)
	if !strings.HasSuffix(pozycja.Name, ogon) {
		return "", "", false
	}
	return strings.TrimSuffix(pozycja.Name, ogon), separator, true
}
