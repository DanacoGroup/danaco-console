// Odpowiedzialność pliku: kolekcje zasobów Assets Panel —
// `design.collection.create`, `design.collection.assign`,
// `design.collection.list`. Metody stoją na `*adapterDesignu`
// (`adapter_modul_design.go`).
//
// Kolekcja jest bytem osobnym od etykiety. Etykieta jest słowem, kolekcja ma
// nazwę, opis i porządek — powód rozdziału stoi w nagłówku migracji 230.
//
// Przypisanie jest dokładką albo odjęciem, nigdy zastąpieniem. Zasoby wchodzące
// do kolekcji sprawdzamy PRZED zapisem: kolekcja pełna kodów, za którymi nie
// stoi żaden zasób, byłaby dokładnie tym rodzajem koperty, który ten moduł już
// raz oddał — wykazem bez treści. Zasób nieznany jest więc odmową całej
// komendy, nie cichym pominięciem, bo `changed` nie ma pola na „ten nie wszedł".
package core

import (
	"context"
	"fmt"
	"strings"

	"danacoconsole/server/internal/dane"
	"danacoconsole/shared"
)

// przedrostekKolekcjiDesign znakuje identyfikatory zewnętrzne kolekcji — obok
// `zasob-`, `plansza-` i `warstwa-`.
const przedrostekKolekcjiDesign = "kolekcja-"

// ZalozKolekcje zakłada kolekcję zasobów — obsługuje `design.collection.create`.
//
// Odpowiedź niesie kolekcję ODCZYTANĄ Z BAZY, nie echo żądania (kontrakt mówi
// to wprost): panel dostaje identyfikator, znacznik czasu i licznik zasobów
// takie, jakie naprawdę zostały zapisane.
func (a *adapterDesignu) ZalozKolekcje(ctx context.Context,
	z shared.DesignCollectionCreateRequest) (shared.DesignCollectionCreateResponse, error) {

	if z.WindowId == "" {
		return shared.DesignCollectionCreateResponse{}, bladWskazaniaDesignu(
			"komenda design.collection.create bez wskazania okna")
	}
	if strings.TrimSpace(z.Name) == "" {
		return shared.DesignCollectionCreateResponse{}, bladWskazaniaDesignu(
			"komenda design.collection.create bez nazwy kolekcji — kolekcja bez nazwy nie da się " +
				"odróżnić od żadnej innej w panelu")
	}

	zapisana, err := a.repozytorium.ZapiszKolekcjeDesignu(ctx, dane.KolekcjaDesignu{
		Kod:   nowyIdentyfikator(przedrostekKolekcjiDesign),
		Okno:  z.WindowId,
		Nazwa: strings.TrimSpace(z.Name),
		Opis:  z.Description,
	})
	if err != nil {
		return shared.DesignCollectionCreateResponse{}, bladDesignu(err)
	}
	return shared.DesignCollectionCreateResponse{Collection: kolekcjaKontraktuDesignu(zapisana)}, nil
}

// PrzypiszDoKolekcji dokłada zasoby do kolekcji albo je z niej zdejmuje —
// obsługuje `design.collection.assign`.
//
// Sprawdzenie istnienia zasobów idzie tylko przy dokładaniu. Przy zdejmowaniu
// zasób z kolekcji bywa już usunięty z Assets Panelu, a odmowa zdjęcia go
// zamknęłaby Operatorowi jedyną drogę do posprzątania kolekcji.
func (a *adapterDesignu) PrzypiszDoKolekcji(ctx context.Context,
	z shared.DesignCollectionAssignRequest) (shared.DesignCollectionAssignResponse, error) {

	if strings.TrimSpace(z.CollectionId) == "" {
		return shared.DesignCollectionAssignResponse{}, bladWskazaniaDesignu(
			"komenda design.collection.assign bez wskazania kolekcji")
	}
	if len(z.AssetIds) == 0 {
		return shared.DesignCollectionAssignResponse{}, bladWskazaniaDesignu(
			"komenda design.collection.assign bez wskazania zasobów")
	}
	kolekcja, err := a.repozytorium.KolekcjaDesignuPoKodzie(ctx, strings.TrimSpace(z.CollectionId))
	if err != nil {
		return shared.DesignCollectionAssignResponse{}, bladNieznanejKolekcjiDesignu(z.CollectionId, err)
	}

	zdejmij := z.Remove != nil && *z.Remove
	if !zdejmij {
		for _, kod := range z.AssetIds {
			kod = strings.TrimSpace(kod)
			if kod == "" {
				continue
			}
			if _, err := a.repozytorium.Zasob(ctx, kod); err != nil {
				return shared.DesignCollectionAssignResponse{}, bladNieznanegoZasobuDesignu(kod, err)
			}
		}
	}

	zmienione, err := a.repozytorium.ZmienPrzypisaniaKolekcjiDesignu(ctx, kolekcja.ID,
		z.AssetIds, zdejmij)
	if err != nil {
		return shared.DesignCollectionAssignResponse{}, bladDesignu(err)
	}

	poZmianie, err := a.repozytorium.KolekcjaDesignuPoKodzie(ctx, kolekcja.Kod)
	if err != nil {
		return shared.DesignCollectionAssignResponse{}, bladDesignu(err)
	}
	return shared.DesignCollectionAssignResponse{
		Collection: kolekcjaKontraktuDesignu(poZmianie),
		Changed:    zmienione,
	}, nil
}

// Kolekcje zwraca kolekcje okna wraz z licznikiem zasobów — obsługuje
// `design.collection.list`.
func (a *adapterDesignu) Kolekcje(ctx context.Context,
	z shared.DesignCollectionListRequest) (shared.DesignCollectionListResponse, error) {

	if z.WindowId == "" {
		return shared.DesignCollectionListResponse{}, bladWskazaniaDesignu(
			"komenda design.collection.list bez wskazania okna")
	}
	wiersze, err := a.repozytorium.KolekcjeDesignu(ctx, z.WindowId, z.AssetId)
	if err != nil {
		return shared.DesignCollectionListResponse{}, bladDesignu(err)
	}
	kolekcje := make([]shared.DesignCollection, 0, len(wiersze))
	for _, wiersz := range wiersze {
		kolekcje = append(kolekcje, kolekcjaKontraktuDesignu(wiersz))
	}
	return shared.DesignCollectionListResponse{Collections: kolekcje, Total: len(kolekcje)}, nil
}

// kolekcjaKontraktuDesignu składa `DesignCollection` kontraktu z wiersza
// repozytorium. `AssetCount` bierze się z licznika policzonego przez bazę, nie
// z długości wykazu — obie liczby są tu równe, ale prawdą o kolekcji jest ta
// z bazy, a nie ta z tego, co akurat udało się wczytać.
func kolekcjaKontraktuDesignu(k dane.KolekcjaDesignu) shared.DesignCollection {
	return shared.DesignCollection{
		Id:          k.Kod,
		WindowId:    k.Okno,
		Name:        k.Nazwa,
		Description: k.Opis,
		AssetIds:    k.Zasoby,
		AssetCount:  k.Liczba,
		UpdatedAt:   chwilaBazy(k.Zaktualizowano),
	}
}

// bladNieznanejKolekcjiDesignu nazywa kolekcję, której rdzeń nie zna. Brak
// wiersza jest odmową `not_found` — panel ma po czym poznać, że jego wykaz
// kolekcji jest nieaktualny; usterka odczytu zostaje usterką rdzenia.
func bladNieznanejKolekcjiDesignu(kod string, err error) error {
	if czyBrakZasobuDesignu(err) {
		return bladNieznanegoBytuDesignu(fmt.Sprintf("kolekcji %s nie ma w module Design", kod))
	}
	return bladDesignu(err)
}
