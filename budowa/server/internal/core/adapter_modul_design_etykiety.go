// Odpowiedzialność pliku: nadawanie etykiet zasobom Assets Panel —
// `design.asset.tag.set` na typie `*adapterDesignu` zadeklarowanym w
// `adapter_modul_design.go`. Jedno wypełnienie portu, plik osobny wedle
// odpowiedzialności, tak jak `adapter_modul_design_kompozycje.go`.
//
// Komenda domyka drogę zapisu etykiet: `design.asset.list` zawęża wykaz polem
// `tags`, a warstwa danych niesie `UstawEtykietyZasobu`, więc brakowało
// wyłącznie przejścia komenda → uchwyt → baza. Własnej migracji ta komenda nie
// potrzebuje.
//
// Pusty zestaw etykiet jest wartością, nie brakiem żądania: `tags: []` znaczy
// „zdejmij wszystkie etykiety” i jest drogą udaną. Odmowa należy się wyłącznie
// zasobowi, którego rdzeń nie zna — pomylenie tych przypadków odebrałoby
// Operatorowi jedyny sposób odetykietowania zasobu, więc pustka nie jest tu
// sprawdzana ani odrzucana.
package core

import (
	"context"
	"errors"

	"danacoconsole/server/internal/dane"
	"danacoconsole/server/internal/protocol"
	"danacoconsole/shared"
)

// UstawEtykietyZasobu podmienia komplet etykiet zasobu na nadesłany —
// obsługuje `design.asset.tag.set`. Semantyka wymiany (nie dokładania) jest
// zapisana w kontrakcie i wykonana w jednej transakcji warstwy danych, więc
// zasób nie zostaje przejściowo bez etykiet, gdy zapis padnie w połowie.
//
// Odpowiedź niesie etykiety odczytane z bazy, nie echo żądania: warstwa danych
// pomija wartości puste i zwraca zestaw uporządkowany, więc odesłanie `z.Tags`
// wprost mówiłoby o stanie, którego w bazie nie ma. Drugi odczyt jest ceną
// prawdy o skutku.
//
// PromptId zostaje pusty z tego samego powodu, co przy `Zasoby`: repozytorium
// nie ma przekładu klucza wiersza promptu na kod kontraktu — pole nie jest
// zgadywane.
func (a *adapterDesignu) UstawEtykietyZasobu(ctx context.Context,
	z shared.DesignAssetTagSetRequest) (shared.DesignAssetTagSetResponse, error) {

	if z.AssetId == "" {
		return shared.DesignAssetTagSetResponse{}, bladWskazaniaDesignu("komenda bez wskazania zasobu")
	}
	zasob, err := a.repozytorium.Zasob(ctx, z.AssetId)
	if err != nil {
		return shared.DesignAssetTagSetResponse{}, bladNieznanegoZasobuDesignu(z.AssetId, err)
	}
	if err := a.repozytorium.UstawEtykietyZasobu(ctx, zasob.ID, z.Tags); err != nil {
		return shared.DesignAssetTagSetResponse{}, bladDesignu(err)
	}
	etykiety, err := a.repozytorium.EtykietyZasobu(ctx, zasob.ID)
	if err != nil {
		return shared.DesignAssetTagSetResponse{}, bladDesignu(err)
	}
	return shared.DesignAssetTagSetResponse{Asset: zasobKontraktu(zasob, etykiety)}, nil
}

// czyBrakZasobuDesignu odróżnia „zasobu nie ma” od usterki odczytu bez
// zamieniania tego w błąd. Woła to `design.asset.remove`, dla którego brak
// zasobu jest odpowiedzią (`removed: false`), a nie odmową — reguła należy do
// warstwy danych (`ErrBrakWiersza`) i ma jedno miejsce w module.
func czyBrakZasobuDesignu(err error) bool {
	return errors.Is(err, dane.ErrBrakWiersza)
}

// bladNieznanegoZasobuDesignu odróżnia „zasobu nie ma” od „odczyt się nie
// powiódł”. Assets Panel ma po tym poznać, że wskazanie w wykazie jest
// nieaktualne, zamiast czytać usterkę rdzenia tam, gdzie jej nie było.
func bladNieznanegoZasobuDesignu(kod string, err error) error {
	if errors.Is(err, dane.ErrBrakWiersza) {
		return protocol.JakoError(protocol.NowyBlad(shared.ErrorCodeNotFound,
			"moduł Design: zasób nie istnieje: "+kod))
	}
	return bladDesignu(err)
}
