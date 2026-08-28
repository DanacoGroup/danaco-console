// Warstwa danych obsługuje obszar modułu Developer: wersje pliku edytora
// w tabeli developer_wersja_pliku i dziennik przebiegów budowania w tabeli
// developer_budowanie.
package dane

import (
	"context"

	"danacoconsole/shared"
)

// WersjaPliku to wiersz tabeli `developer_wersja_pliku` — migawka treści pliku
// sprzed zapisu, założona przez `developer.file.save` z `createVersion`.
type WersjaPliku struct {
	Kod       string
	OknoKod   string
	Sciezka   string
	Tresc     string
	Rozmiar   int64
	Utworzono string
}

// PrzebiegBudowania to wiersz tabeli `developer_budowanie` — jeden przebieg
// zadania budowania wraz z jego wynikiem i ogonem logu.
type PrzebiegBudowania struct {
	Kod         string
	OknoKod     string
	Zadanie     string
	Argumenty   *string
	Stan        shared.BuildStatus
	KodWyjscia  *int64
	Log         *string
	Uruchomiono string
	Zakonczono  *string
}

// RepozytoriumDevelopera jest kontraktem obszaru Developer: wersje plików,
// przebiegi budowania, punkty przerwania, kolekcje zapytań, połączenia
// danych, skanowanie i wyniki testów.
type RepozytoriumDevelopera interface {
	// ZapiszWersje zakłada migawkę treści pliku.
	ZapiszWersje(ctx context.Context, wersja WersjaPliku) error
	// OstatniaWersja zwraca najnowszą migawkę pliku okna; brak migawki wraca jako ErrBrakWiersza.
	OstatniaWersja(ctx context.Context, oknoKod, sciezka string) (WersjaPliku, error)
	// Wersje zwraca migawki pliku od najnowszej; `limit` niedodatni znaczy
	// wykaz pełny.
	Wersje(ctx context.Context, oknoKod, sciezka string, limit int) ([]WersjaPliku, error)
	// WersjaPoKodzie zwraca jedną migawkę po jej identyfikatorze; ścieżka pliku wynika z migawki.
	WersjaPoKodzie(ctx context.Context, kod string) (WersjaPliku, error)

	// ZapiszPrzebieg zakłada albo odświeża wiersz przebiegu budowania.
	ZapiszPrzebieg(ctx context.Context, przebieg PrzebiegBudowania) error
	// ZakonczPrzebieg domyka przebieg wynikiem i ogonem logu; przebieg już domknięty zostaje bez zmiany.
	ZakonczPrzebieg(ctx context.Context, kod string, stan shared.BuildStatus,
		kodWyjscia *int64, log string) error
	// OstatniPrzebieg zwraca najnowszy przebieg budowania okna.
	OstatniPrzebieg(ctx context.Context, oknoKod string) (PrzebiegBudowania, error)
	// OsierocPrzebiegi przestawia przebiegi w stanie running po poprzednim biegu rdzenia na stopped.
	OsierocPrzebiegi(ctx context.Context) (int64, error)

	// Przebiegi zwraca dziennik budowań okna od najnowszego; pusty stan i niedodatni limit nie zawężają.
	Przebiegi(ctx context.Context, oknoKod, stan string, limit int) ([]PrzebiegBudowania, error)
	// Przebieg zwraca jeden przebieg budowania po jego identyfikatorze.
	Przebieg(ctx context.Context, kod string) (PrzebiegBudowania, error)

	// ZapiszPunktPrzerwania zakłada albo odświeża punkt przerwania w pliku okna.
	ZapiszPunktPrzerwania(ctx context.Context, punkt PunktPrzerwania) error
	// UsunPunktPrzerwania zdejmuje punkt z wiersza pliku.
	UsunPunktPrzerwania(ctx context.Context, oknoKod, sciezka string, wiersz int64) error
	// PunktyPrzerwania zwraca punkty okna; ścieżka pusta znaczy wszystkie pliki.
	PunktyPrzerwania(ctx context.Context, oknoKod, sciezka string) ([]PunktPrzerwania, error)

	// ZapiszKolekcjeApi zakłada albo nadpisuje kolekcję zapytań okna.
	ZapiszKolekcjeApi(ctx context.Context, kolekcja KolekcjaApi) error
	// KolekcjeApi zwraca kolekcje okna; kod niepusty zawęża do jednej.
	KolekcjeApi(ctx context.Context, oknoKod, kod string) ([]KolekcjaApi, error)

	// ZapiszPolaczenieDanych zakłada albo nadpisuje opis połączenia.
	ZapiszPolaczenieDanych(ctx context.Context, polaczenie PolaczenieDanych) error
	// PolaczeniaDanych zwraca połączenia okna w kolejności nazwy.
	PolaczeniaDanych(ctx context.Context, oknoKod string) ([]PolaczenieDanych, error)
	// PolaczenieDanychPoKodzie zwraca jedno połączenie; brak wraca jako ErrBrakWiersza.
	PolaczenieDanychPoKodzie(ctx context.Context, kod string) (PolaczenieDanych, error)

	// ZapiszSkan zakłada albo domyka przebieg skanowania.
	ZapiszSkan(ctx context.Context, skan PrzebiegSkanu) error
	// ZapiszZnaleziska dopisuje znaleziska przebiegu jedną transakcją.
	ZapiszZnaleziska(ctx context.Context, znaleziska []ZnaleziskoSkanu) error
	// Znaleziska zwraca spostrzeżenia wedle filtru, od najcięższych.
	Znaleziska(ctx context.Context, filtr FiltrZnalezisk) ([]ZnaleziskoSkanu, error)

	// ZapiszWynikiTestow zastępuje wyniki testów przebiegu.
	ZapiszWynikiTestow(ctx context.Context, budowanieKod string, wyniki []WynikTestu) error
	// WynikiTestow zwraca wyniki przebiegu; stan pusty znaczy wszystkie.
	WynikiTestow(ctx context.Context, budowanieKod, stan string) ([]WynikTestu, error)
	// ZapiszPokrycie zastępuje pomiar pokrycia przebiegu.
	ZapiszPokrycie(ctx context.Context, budowanieKod string, pokrycie []PokryciePliku) error
	// Pokrycie zwraca pomiar pokrycia przebiegu; ścieżka pusta znaczy wszystkie
	// pliki.
	Pokrycie(ctx context.Context, budowanieKod, sciezka string) ([]PokryciePliku, error)
}
