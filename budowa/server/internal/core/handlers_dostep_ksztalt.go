// Plik sprawdza kształt punktu dostępu, żądanie wobec więzów schematu, i
// rozkłada adres mostu na kolumny wiersza.
package core

import (
	"errors"
	"strconv"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// sprawdzKsztaltPunktu odmawia założenia punktu, którego schemat i tak by nie
// przyjął, i podaje operatorowi powód odmowy.
func sprawdzKsztaltPunktu(z shared.AccessPointAddRequest) error {
	if strings.TrimSpace(z.Name) == "" {
		return odmowaKsztaltuPunktu("punkt dostępu wymaga nazwy")
	}
	if z.Kind == shared.AccessPointKindMcpBridge &&
		wartoscTekstu(z.Host) == "" && wartoscTekstu(z.Endpoint) == "" {
		return odmowaKsztaltuPunktu("punkt rodzaju mcpBridge wymaga nazwy maszyny albo adresu mostu")
	}
	return nil
}

// odmowaKsztaltuPunktu składa odmowę merytoryczną założenia punktu, z powodem
// właściwym niespełnionemu więzowi schematu.
func odmowaKsztaltuPunktu(powod string) error {
	return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeValidationFailed, powod))
}

// odmowaUrzadzeniaPunktu przekłada rozstrzygnięcia katalogu urządzeń na odmowę
// z powodem, nie na błąd wewnętrzny rdzenia.
func odmowaUrzadzeniaPunktu(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, dane.ErrBrakUrzadzeniaBiezacego):
		return odmowaKsztaltuPunktu("punkt rodzaju localDirectory wymaga urządzenia, " +
			"a serwer nie ma w katalogu wiersza swojej maszyny — rozpoznanie startowe " +
			"nie przebiegło; wskaż urządzenie polem deviceId albo powtórz po starcie serwera")
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
