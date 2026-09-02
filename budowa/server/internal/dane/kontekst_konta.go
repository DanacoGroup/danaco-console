// Odpowiedzialność pliku: przeniesienie wskazania konta Operatora przez kontekst
// żądania, od chwili rozpoznania sesji bramki po zapytania sięgające pracy.
package dane

import (
	"context"
	"database/sql"
	"fmt"
)

// kluczKontaOperatora jest kluczem własnym pakietu, więc żaden inny pakiet nie
// nadpisze wskazania przypadkiem ani go nie odczyta bez tej funkcji.
type kluczKontaOperatora struct{}

/*
ZKontemOperatora zwraca kontekst niosący wskazanie konta, do którego należy
żądanie. Wskazanie stawia warstwa rozpoznająca sesję bramki — jedyne miejsce,
które wie, czyje jest połączenie.

Konto wędruje kontekstem, a nie osobnym argumentem każdej funkcji, bo dotyka
zapytań leżących głęboko pod warstwą, która je rozpoznaje, i przechodzi przez
ogniwa, których samo konto nie obchodzi.
*/
func ZKontemOperatora(ctx context.Context, kontoId int64) context.Context {
	return context.WithValue(ctx, kluczKontaOperatora{}, kontoId)
}

/*
KontoOperatora odczytuje wskazanie konta z kontekstu żądania. Zero znaczy
żądanie bez rozpoznanego konta — połączenie przed zalogowaniem albo czynność
rdzenia własną, nie wywołaną przez Operatora.

Zero nie jest usterką: zapytania czytają wtedy pracę konta najstarszego, bo tak
stała cała praca zapisana przed rozdzieleniem kont (migracja 407).
*/
func KontoOperatora(ctx context.Context) int64 {
	kontoId, jest := ctx.Value(kluczKontaOperatora{}).(int64)
	if !jest {
		return 0
	}
	return kontoId
}

/*
WarunekKonta zawęża wiersz do konta, do którego należy żądanie. Jest jeden na
wszystkie tabele niosące pracę Operatora, bo rozstrzygnięcie musi wypaść tak
samo w każdym zapytaniu — granica trzymająca w jednym, a puszczająca w sąsiednim
nie jest granicą.

Warunek bierze JEDEN argument: wynik KontoOperatora. Obie strony porównania
sprowadzają brak wskazania do konta najstarszego — wiersz bez `konto_id` powstał
przed rozdzieleniem kont i nie ma jak wskazać konta wstecz, a żądanie bez
rozpoznanego konta przychodzi z połączenia przed zalogowaniem. Porównanie idzie
przez IS, nie przez znak równości: na instalacji przed rejestracją obie strony
są puste, a pustka porównana znakiem równości nie jest prawdą i praca zastana
znikłaby z oczu.
*/
const WarunekKonta = `COALESCE(konto_id, (SELECT id FROM konto_wlasciciela ORDER BY id LIMIT 1))
	                  IS COALESCE(NULLIF(?, 0), (SELECT id FROM konto_wlasciciela ORDER BY id LIMIT 1))`

// WskazanieKonta zapisuje konto w wierszu zakładanym. Zero znaczyłoby konto
// o identyfikatorze zero, czyli żadne, i wiersz byłby niewidoczny dla własnego
// właściciela; wskazanie puste wchodzi jako NULL, tak jak wiersz zastany.
const WskazanieKonta = `NULLIF(?, 0)`

/*
sprawdzTrafienieZapisu odróżnia zapis wykonany od zapisu zatrzymanego przez
WarunekKonta przy ON CONFLICT DO UPDATE. Zero zmienionych wierszy znaczy wtedy,
że wiersz stoi i należy do innego konta; silnik nie zgłasza tego błędem, więc
odmowa musi powstać tutaj, zamiast cichego powodzenia albo błędu odczytu
po zapisie.
*/
func sprawdzTrafienieZapisu(wynik sql.Result, byt, wskazanie string) error {
	zmienione, err := wynik.RowsAffected()
	if err != nil {
		return fmt.Errorf("dane: nieznana liczba zapisanych wierszy (%s %q): %w", byt, wskazanie, err)
	}
	if zmienione == 0 {
		return fmt.Errorf("dane: %s %q należy do innego konta: %w", byt, wskazanie, ErrKolizjaWiersza)
	}
	return nil
}
