// Odpowiedzialność pliku: kształt punktu dostępu — sprawdzenie żądania wobec
// więzów schematu oraz rozłożenie adresu mostu na kolumny wiersza.
//
// Sprawdzenie jest tu, a nie w warstwie danych, bo dotyczy żądania: baza pilnuje
// więzu, ale zgłasza go jako awarię zapisu. Operator ma dostać odmowę
// merytoryczną z powodem, nie błąd wewnętrzny rdzenia.
package core

import (
	"errors"
	"strconv"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// sprawdzKsztaltPunktu odmawia założenia punktu, którego schemat i tak by nie
// przyjął, i podaje Operatorowi powód.
//
// Most MCP bez nazwy maszyny i bez adresu nie ma dokąd prowadzić.
//
// Katalog lokalny bez pola `deviceId` nie jest odrzucany. Pole jest w kontrakcie
// opcjonalne, a jego pominięcie znaczy „na tej maszynie”: katalog wskazany oknem
// powłoki leży na maszynie, na której działa rdzeń, i tę maszynę katalog
// urządzeń zna z rozpoznania startowego (kolumna `biezace`).
func sprawdzKsztaltPunktu(z shared.AccessPointAddRequest) error {
	if z.Kind == shared.AccessPointKindMcpBridge &&
		wartoscTekstu(z.Host) == "" && wartoscTekstu(z.Endpoint) == "" {
		return odmowaKsztaltuPunktu("punkt rodzaju mcpBridge wymaga nazwy maszyny albo adresu mostu")
	}
	return nil
}

// odmowaKsztaltuPunktu składa odmowę merytoryczną założenia punktu.
func odmowaKsztaltuPunktu(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed, powod))
}

// odmowaUrzadzeniaPunktu przekłada rozstrzygnięcia katalogu urządzeń na odmowę
// z powodem. Oba przypadki są brakiem w żądaniu albo w stanie platformy, nie
// awarią trwałości, więc Operator ma zobaczyć zdanie, a nie ciszę ani błąd
// wewnętrzny. Każdy inny błąd idzie dalej nietknięty — rdzeń nie zgaduje
// za bazę, czy zawiódł dysk, czy schemat.
func odmowaUrzadzeniaPunktu(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, dane.ErrBrakUrzadzeniaBiezacego):
		return odmowaKsztaltuPunktu("punkt rodzaju localDirectory wymaga urządzenia, " +
			"a rdzeń nie ma w katalogu wiersza swojej maszyny — rozpoznanie startowe " +
			"nie przebiegło; wskaż urządzenie polem deviceId albo powtórz po starcie rdzenia")
	case errors.Is(err, dane.ErrNieznaneUrzadzenie):
		return odmowaKsztaltuPunktu("wskazanego urządzenia nie ma w katalogu urządzeń")
	}
	return err
}

// rozlozAdresPunktu rozkłada adres mostu na kolumny wiersza. Katalog lokalny
// adresu nie ma — nie sięga się do niego po SSH.
func rozlozAdresPunktu(punkt *dane.PunktDostepu, host, endpoint *string) {
	if punkt.Rodzaj == shared.AccessPointKindLocalDirectory {
		punkt.Host = wartoscTekstu(host)
		return
	}
	adres := AdresMostu(shared.AccessPoint{Id: punkt.Kod, Host: host, Endpoint: endpoint,
		BridgeName: tekstOpcjonalny(punkt.NazwaMostu)})
	punkt.Host = adres.Host
	punkt.Uzytkownik = adres.Uzytkownik
	if port, err := strconv.Atoi(adres.Port); err == nil {
		punkt.Port = port
	}
}
