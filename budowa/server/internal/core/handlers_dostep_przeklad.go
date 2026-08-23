// Odpowiedzialność pliku: przekład wierszy obszaru dostępów na struktury
// kontraktu i z powrotem. Jedno miejsce styku obu nazewnictw.
//
// IDENTYFIKATORY. Punkt wychodzi kontraktem pod swoim trwałym kodem, nie pod
// numerem wiersza — kod jest tym, czym Operator posługuje się w konfiguracji
// mostu i co przeżywa przeniesienie bazy. Nadanie wychodzi pod identyfikatorem
// zewnętrznym nadanym przez rdzeń, a gdy wiersz powstał wprost w bazie — pod
// numerem wiersza (identyfikatorWiersza, jak przy oknach i sesjach).
package core

import (
	"strconv"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// punktKontraktu przekłada wiersz katalogu na strukturę AccessPoint.
//
// pole credentialRef niesie ODWOŁANIE do wpisu w sejfie — nazwę wpisu
// albo ścieżkę profilu. Treści poświadczenia nie ma ani w kolumnie, ani tutaj.
func punktKontraktu(p dane.PunktDostepu) shared.AccessPoint {
	punkt := shared.AccessPoint{
		Id:            p.Kod,
		Name:          p.Nazwa,
		Description:   tekstOpcjonalny(p.Opis),
		Kind:          p.Rodzaj,
		Host:          tekstOpcjonalny(p.Host),
		Endpoint:      tekstOpcjonalny(p.AdresMostu()),
		BridgeName:    tekstOpcjonalny(p.NazwaMostu),
		Roots:         korzenieKontraktu(p.Korzenie),
		DefaultMode:   p.TrybDomyslny,
		CredentialRef: p.PoswiadczenieOdwolanie,
		Status:        p.Stan,
		CheckedAt:     chwilaOpcjonalna(p.Sprawdzono),
		Enabled:       p.Aktywny,
		CreatedAt:     chwilaBazy(p.Utworzono),
		UpdatedAt:     chwilaBazy(p.Zaktualizowano),
	}
	if p.UrzadzenieID != nil {
		punkt.DeviceId = tekstOpcjonalny(strconv.FormatInt(*p.UrzadzenieID, 10))
	}
	return punkt
}

// nadanieKontraktu przekłada wiersz nadania na strukturę AccessGrant. Kod punktu
// i identyfikator okna przychodzą z zewnątrz, bo wiersz niesie numery wierszy,
// a kontrakt — identyfikatory trwałe.
func nadanieKontraktu(n dane.Nadanie, kodPunktu, idOkna string) shared.AccessGrant {
	return shared.AccessGrant{
		Id:            identyfikatorWiersza(n.IdentyfikatorZewnetrzny, n.ID),
		WindowId:      idOkna,
		AccessPointId: kodPunktu,
		Mode:          n.Tryb,
		Roots:         korzenieKontraktu(n.Korzenie),
		Order:         n.Kolejnosc,
		Primary:       n.Glowne,
		Enabled:       n.Aktywne,
		CreatedAt:     chwilaBazy(n.Utworzono),
	}
}

// korzenieKontraktu oddaje listę korzeni tak, żeby pole kontraktu nigdy nie
// było wartością pustą typu nil — klient dostaje tablicę, także pustą.
func korzenieKontraktu(korzenie []string) []string {
	if korzenie == nil {
		return []string{}
	}
	return korzenie
}

// chwilaOpcjonalna przekłada znacznik czasu bazy na pole opcjonalne kontraktu.
// Brak znacznika i znacznik nieczytelny dają brak wartości, nie zero epoki —
// zero znaczyłoby „sprawdzono w 1970 roku”.
func chwilaOpcjonalna(znacznik *string) *int64 {
	if znacznik == nil || *znacznik == "" {
		return nil
	}
	chwila := chwilaBazy(*znacznik)
	if chwila == 0 {
		return nil
	}
	return &chwila
}

// urzadzenieWiersza przekłada identyfikator urządzenia z kontraktu na kolumnę.
//
// Identyfikatorem urządzenia w kontrakcie jest klucz wiersza katalogu — ten sam,
// który wychodzi z powrotem polem `AccessPoint.deviceId` (patrz punktKontraktu).
// Wartość nieliczbowa nie wskazuje żadnego urządzenia i wraca odmową. Ciche
// `nil` rozbiłoby się dopiero o więz schematu
// `CHECK(rodzaj <> 'localDirectory' OR urzadzenie_id IS NOT NULL)`, a Operator
// dostałby awarię zapisu zamiast powodu.
//
// Brak pola to nie to samo co pole niezrozumiałe: pominięcie jest dopuszczone
// kontraktem i znaczy „urządzenia nie wskazuję”, co dla katalogu lokalnego
// rozstrzyga warstwa trwałości maszyną, na której działa rdzeń.
func urzadzenieWiersza(urzadzenie *string) (*int64, error) {
	if urzadzenie == nil || *urzadzenie == "" {
		return nil, nil
	}
	numer, err := strconv.ParseInt(*urzadzenie, 10, 64)
	if err != nil || numer <= 0 {
		return nil, odmowaKsztaltuPunktu("identyfikator urządzenia " +
			strconv.Quote(*urzadzenie) + " nie jest identyfikatorem katalogu urządzeń")
	}
	return &numer, nil
}
