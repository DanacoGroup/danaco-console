// Odpowiedzialność pliku: rozpoznanie urządzenia bieżącego, czyli maszyny, na
// której działa ten rdzeń. Rdzeń zna swoją maszynę z systemu operacyjnego, a
// warstwa danych zna katalog urządzeń — spotykają się tutaj i nigdzie indziej.
//
// Punkt dostępu rodzaju `localDirectory` musi wskazać urządzenie (warunek CHECK
// tabeli `punkt_dostepu`), bo katalog lokalny istnieje na jednej maszynie.
// Rozpoznanie startowe zakłada wiersz maszyny, na której stoi rdzeń, i odświeża
// go przy każdym starcie — powtórzenie nie tworzy drugiej maszyny.
package core

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"danacoconsole/server/internal/dane"
)

// nazwaHostaNierozpoznana zastępuje nazwę hosta, gdy system jej nie poda. Brak
// nazwy nie może przerwać startu rdzenia — maszyna dostaje wtedy wiersz
// rozpoznawalny po samym systemie operacyjnym.
const nazwaHostaNierozpoznana = "nierozpoznany"

// znamionaMaszyny to fakty maszyny czytane z systemu operacyjnego.
//
// `Identyfikator` nie jest numerem sprzętu: składamy go z systemu i nazwy hosta,
// bo tyle da się ustalić bez sięgania po dane sprzętowe maszyny. Zmiana nazwy
// hosta daje więc nowe urządzenie w katalogu — z punktu widzenia katalogów
// lokalnych jest to inna maszyna.
type znamionaMaszyny struct {
	NazwaHosta       string
	SystemOperacyjny string
	Identyfikator    string
}

// odczytajZnamionaMaszyny ustala znamiona maszyny, na której działa proces.
func odczytajZnamionaMaszyny() znamionaMaszyny {
	host, err := os.Hostname()
	if err != nil {
		host = ""
	}
	host = strings.TrimSpace(host)
	if host == "" {
		host = nazwaHostaNierozpoznana
	}
	return znamionaMaszyny{
		NazwaHosta:       host,
		SystemOperacyjny: runtime.GOOS,
		Identyfikator:    runtime.GOOS + "/" + strings.ToLower(host),
	}
}

// urzadzenie przekłada znamiona na wiersz katalogu urządzeń. Maszyna, na której
// działa rdzeń, jest zaufana z założenia — kod rdzenia już się na niej wykonuje.
// Oznaczenie zaufania zapisuje się wyłącznie przy zakładaniu wiersza; odebranie
// go później jest decyzją Operatora i rozpoznanie startowe jej nie cofa.
func (z znamionaMaszyny) urzadzenie() dane.Urzadzenie {
	return dane.Urzadzenie{
		Nazwa:                  z.NazwaHosta,
		IdentyfikatorSprzetowy: z.Identyfikator,
		NazwaHosta:             z.NazwaHosta,
		SystemOperacyjny:       z.SystemOperacyjny,
		Zaufane:                true,
		Biezace:                true,
	}
}

// rozpoznajUrzadzenieBiezace zapewnia wiersz maszyny bieżącej i zwraca go wraz
// z kluczem nadanym przez bazę. Wywołanie powtórzone na tej samej maszynie
// zwraca ten sam wiersz — idempotencji pilnuje warstwa danych, bo tylko ona
// widzi jednocześnie identyfikator sprzętowy i oznaczenie maszyny bieżącej.
func rozpoznajUrzadzenieBiezace(kontekst context.Context,
	repozytorium dane.RepozytoriumUrzadzen) (dane.Urzadzenie, error) {

	if repozytorium == nil {
		return dane.Urzadzenie{}, fmt.Errorf("core: katalog urządzeń nie jest złożony")
	}
	return repozytorium.ZapewnijBiezace(kontekst, odczytajZnamionaMaszyny().urzadzenie())
}

// odnotujUrzadzenieBiezace wykonuje rozpoznanie przy montażu rdzenia. Nieudane
// rozpoznanie idzie do dziennika i nie przerywa startu: rdzeń bez wiersza swojej
// maszyny pracuje dalej, tylko katalog lokalny nie ma na czym stanąć, dopóki
// Operator nie wskaże urządzenia sam.
func odnotujUrzadzenieBiezace(kontekst context.Context, repozytorium dane.RepozytoriumUrzadzen,
	dziennik *log.Logger) {

	urzadzenie, err := rozpoznajUrzadzenieBiezace(kontekst, repozytorium)
	if dziennik == nil {
		return
	}
	if err != nil {
		dziennik.Printf("rozpoznanie urządzenia bieżącego nie powiodło się: %v", err)
		return
	}
	dziennik.Printf("urządzenie bieżące: #%d %q (host=%s system=%s identyfikator=%s)",
		urzadzenie.ID, urzadzenie.Nazwa, urzadzenie.NazwaHosta, urzadzenie.SystemOperacyjny,
		urzadzenie.IdentyfikatorSprzetowy)
}
