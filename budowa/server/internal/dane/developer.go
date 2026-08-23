// Odpowiedzialność pliku: obszar modułu Developer — wersja pliku edytora
// (tabela `developer_wersja_pliku`) i dziennik przebiegów budowania (tabela
// `developer_budowanie`) — byty obszaru i jego kontrakt. Odczyt leży
// w `developer_odczyt.go`, zapis w `developer_zapis.go`.
//
// Repozytorium nie dotyka plików na dysku. Treść pliku roboczego żyje
// w katalogu roboczym okna i czyta ją rdzeń; tutaj zapisuje się wyłącznie
// migawka zakładana na wyraźne żądanie (`createVersion`) — czyli to, czego na
// dysku po nadpisaniu już nie ma.
//
// Repozytorium nie prowadzi budowania. Uchwyt do biegnącego procesu ma rdzeń;
// tutaj zostaje ślad: zadanie, kod wyjścia, czasy i ogon logu.
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

// RepozytoriumDevelopera jest kontraktem obszaru Developer.
type RepozytoriumDevelopera interface {
	// ZapiszWersje zakłada migawkę treści pliku.
	ZapiszWersje(ctx context.Context, wersja WersjaPliku) error
	// OstatniaWersja zwraca najnowszą migawkę pliku okna. Brak migawki wraca
	// jako ErrBrakWiersza — plik bez wersji jest stanem zwykłym, nie usterką.
	OstatniaWersja(ctx context.Context, oknoKod, sciezka string) (WersjaPliku, error)
	// Wersje zwraca migawki pliku od najnowszej; `limit` niedodatni znaczy
	// wykaz pełny.
	Wersje(ctx context.Context, oknoKod, sciezka string, limit int) ([]WersjaPliku, error)
	// WersjaPoKodzie zwraca jedną migawkę. Czyta to
	// `developer.file.version.restore`, które zna wyłącznie identyfikator
	// wersji — ścieżka pliku wynika z migawki, a nie z żądania.
	WersjaPoKodzie(ctx context.Context, kod string) (WersjaPliku, error)

	// ZapiszPrzebieg zakłada albo odświeża wiersz przebiegu budowania.
	ZapiszPrzebieg(ctx context.Context, przebieg PrzebiegBudowania) error
	// ZakonczPrzebieg domyka przebieg wynikiem i ogonem logu. Wiersz już
	// domknięty zostaje bez zmiany — pierwszy prawdziwy kod wyjścia nie ma
	// prawa zostać nadpisany przez późniejsze przerwanie.
	ZakonczPrzebieg(ctx context.Context, kod string, stan shared.BuildStatus,
		kodWyjscia *int64, log string) error
	// OstatniPrzebieg zwraca najnowszy przebieg okna. Czyta to
	// `developer.build.run` z `stop=true` dla okna, w którym nic nie biegnie.
	OstatniPrzebieg(ctx context.Context, oknoKod string) (PrzebiegBudowania, error)
	// OsierocPrzebiegi przestawia przebiegi zostawione w stanie `running` przez
	// poprzedni bieg rdzenia na `stopped`. Rdzeń po restarcie nie ma do nich
	// uchwytu, więc wykazywanie ich jako czynnych byłoby nieprawdą.
	OsierocPrzebiegi(ctx context.Context) (int64, error)

	// Przebiegi zwraca dziennik budowań okna od najnowszego. Stan pusty znaczy
	// wykaz bez zawężenia, limit niedodatni — wykaz pełny.
	Przebiegi(ctx context.Context, oknoKod, stan string, limit int) ([]PrzebiegBudowania, error)
	// Przebieg zwraca jeden przebieg po jego identyfikatorze. Czyta to
	// `developer.build.log.get`, gdy przebieg zdążył się już domknąć.
	Przebieg(ctx context.Context, kod string) (PrzebiegBudowania, error)

	// ── Punkty przerwania (Run & Debug) ──────────────────────────────────────
	// ZapiszPunktPrzerwania zakłada albo odświeża punkt w pliku okna.
	ZapiszPunktPrzerwania(ctx context.Context, punkt PunktPrzerwania) error
	// UsunPunktPrzerwania zdejmuje punkt z wiersza pliku.
	UsunPunktPrzerwania(ctx context.Context, oknoKod, sciezka string, wiersz int64) error
	// PunktyPrzerwania zwraca punkty okna; ścieżka pusta znaczy wszystkie pliki.
	PunktyPrzerwania(ctx context.Context, oknoKod, sciezka string) ([]PunktPrzerwania, error)

	// ── Kolekcje zapytań (API Client) ────────────────────────────────────────
	// ZapiszKolekcjeApi zakłada albo nadpisuje kolekcję zapytań okna.
	ZapiszKolekcjeApi(ctx context.Context, kolekcja KolekcjaApi) error
	// KolekcjeApi zwraca kolekcje okna; kod niepusty zawęża do jednej.
	KolekcjeApi(ctx context.Context, oknoKod, kod string) ([]KolekcjaApi, error)

	// ── Połączenia bazodanowe (Data Console) ─────────────────────────────────
	// ZapiszPolaczenieDanych zakłada albo nadpisuje opis połączenia.
	ZapiszPolaczenieDanych(ctx context.Context, polaczenie PolaczenieDanych) error
	// PolaczeniaDanych zwraca połączenia okna w kolejności nazwy.
	PolaczeniaDanych(ctx context.Context, oknoKod string) ([]PolaczenieDanych, error)
	// PolaczenieDanychPoKodzie zwraca jedno połączenie. Brak wraca jako
	// ErrBrakWiersza — konsola SQL odróżnia „nie ma takiego połączenia” od
	// „odczyt zawiódł”.
	PolaczenieDanychPoKodzie(ctx context.Context, kod string) (PolaczenieDanych, error)

	// ── Skanowanie (bezpieczeństwo i jakość) ─────────────────────────────────
	// ZapiszSkan zakłada albo domyka przebieg skanowania.
	ZapiszSkan(ctx context.Context, skan PrzebiegSkanu) error
	// ZapiszZnaleziska dopisuje znaleziska przebiegu jedną transakcją.
	ZapiszZnaleziska(ctx context.Context, znaleziska []ZnaleziskoSkanu) error
	// Znaleziska zwraca spostrzeżenia wedle filtru, od najcięższych.
	Znaleziska(ctx context.Context, filtr FiltrZnalezisk) ([]ZnaleziskoSkanu, error)

	// ── Wyniki testów i pokrycie (Build Output) ──────────────────────────────
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
