// Odpowiedzialność pliku: rzeczy, których używa cała dobudowa modułu
// Automations — przekład chwili kontraktu na znacznik bazy, przedrostki
// identyfikatorów, wpis do dziennika audytu i dobór przebiegu.
package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// Przedrostki identyfikatorów bytów dobudowy. Byt nadany przez rdzeń wychodzi
// kontraktem pod własnym identyfikatorem, nie pod kluczem wiersza.
const (
	przedrostekSzablonuAutomatyki = "szablon-"
	przedrostekReguly             = "alarm-"
	przedrostekPunktuWznowienia   = "punkt-"
	przedrostekWpisuAudytu        = "audyt-"
	przedrostekZleceniaKolejki    = "zlec-"
	przedrostekPoswiadczenia      = "sekret-automatyki-"
)

// wykonawcaAudytu nazywa sprawcę czynności. Rdzeń pracuje na jednym koncie
// Operatora i kontrakt nie niesie tożsamości w żądaniach modułu, więc zapis
// niesie prawdę: czynność wykonał Operator tej instalacji.
const wykonawcaAudytu = "operator"

// znacznikChwiliAutomatyzacji przekłada milisekundy epoki kontraktu na znacznik bazy.
// Zero znaczy „bez wskazania” i daje znacznik pusty — zawężenia zapytań czytają
// pusty napis jako brak zawężenia, a nie jako początek epoki.
func znacznikChwiliAutomatyzacji(milisekundy int64) string {
	if milisekundy == 0 {
		return ""
	}
	return time.UnixMilli(milisekundy).UTC().Format(formatZnacznikaBazy)
}

// znacznikChwiliWskazanejAutomatyzacji zdejmuje wskaźnik z pola opcjonalnego
// chwili i przekłada na znacznik bazy, albo oddaje pustkę dla pola nieustawionego.
func znacznikChwiliWskazanejAutomatyzacji(milisekundy *int64) string {
	if milisekundy == nil {
		return ""
	}
	return znacznikChwiliAutomatyzacji(*milisekundy)
}

// zapisAudytu nanosi wpis dziennika audytu. Szczegóły są zapisem strukturalnym
// budowanym przez wołającego; brak szczegółów zostawia kolumnę pustą.
func (a *adapterAutomatyk) zapisAudytu(ctx context.Context, automatykaID *int64,
	czynnosc string, szczegoly any) {

	wpis := dane.WpisAudytuAutomatyki{
		Kod: nowyIdentyfikator(przedrostekWpisuAudytu), AutomatykaID: automatykaID,
		Wykonawca: wykonawcaAudytu, Czynnosc: czynnosc,
	}
	if szczegoly != nil {
		if tresc, err := json.Marshal(szczegoly); err == nil {
			zapis := string(tresc)
			wpis.Szczegoly = &zapis
		}
	}
	// Skutek zapisu celowo pominięty: audyt nie wywraca czynności, która się powiodła.
	_ = a.repozytorium.DopiszAudytAutomatyki(ctx, wpis)
}

// wierszPrzebiegu odnajduje przebieg po jego identyfikatorze kontraktu i nazywa
// brak wprost. Rodzina `automation.execution.*` wskazuje przebieg, a nie
// automatykę, więc ta droga powtarza się w niej sześć razy.
func (a *adapterAutomatyk) wierszPrzebiegu(ctx context.Context, kod string) (dane.Przebieg, error) {
	if kod == "" {
		return dane.Przebieg{}, bladWskazaniaAutomatyki("komenda bez wskazania przebiegu")
	}
	wiersz, err := a.repozytorium.Przebieg(ctx, kod)
	if err != nil {
		return dane.Przebieg{}, bladNieznanegoPrzebiegu(kod, err)
	}
	return wiersz, nil
}

// bladNieznanegoBytuAutomatyki odróżnia „bytu nie ma” od „odczyt się nie
// powiódł”. Jedna funkcja na całą dobudowę, bo powód jest zawsze ten sam.
func bladNieznanegoBytuAutomatyki(err error, powod string, wskazanie any) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			fmt.Sprintf("moduł Automations: %s%v", powod, wskazanie)))
	}
	return bladAutomatyki(err)
}

// zapisStrukturalny składa zapis strukturalny z wartości kontraktu. Wartość
// niezapisywalna daje brak zapisu, nie napis „null” — kolumna pusta znaczy
// „nie podano”, a napis „null” znaczyłby „podano wartość pustą”.
func zapisStrukturalny(wartosc any) *string {
	if wartosc == nil {
		return nil
	}
	tresc, err := json.Marshal(wartosc)
	if err != nil {
		return nil
	}
	zapis := string(tresc)
	return &zapis
}

// wartoscLiczbyLub zdejmuje wskaźnik z pola opcjonalnego, podstawiając wartość
// zastępczą przy braku wskazania.
func wartoscLiczbyLub(wskazanie *int, zastepcza int) int {
	if wskazanie == nil {
		return zastepcza
	}
	return *wskazanie
}
