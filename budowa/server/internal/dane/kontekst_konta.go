// Odpowiedzialność pliku: przeniesienie wskazania konta Operatora przez kontekst
// żądania, od chwili rozpoznania sesji bramki po zapytania sięgające pracy.
package dane

import "context"

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
