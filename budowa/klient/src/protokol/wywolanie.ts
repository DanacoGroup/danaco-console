import type { Command, RequestOf, ResponseOf } from '../../../shared/contract.ts';
import type { Kanal, Wynik } from './kanal.ts';

/**
 * Wywołanie komendy kontraktu jako obietnica.
 *
 * Kanał rozdaje wynik przez wywołanie zwrotne, ponieważ tak wygodniej obsłużyć
 * strumień i ruch ciągły. Wołający pyta jednak punktowo — „powitaj rdzeń
 * i pokaż, co wróciło" — i dla niego obietnica jest formą naturalną. To nie
 * jest druga droga do rdzenia: każde wywołanie idzie tym samym
 * `kanal.wyslij`, z nazwą komendy wziętą wyłącznie ze stałych kontraktu.
 *
 * Obietnica nie jest odrzucana nigdy. Niepowodzenie wraca jako `Wynik` z polem
 * `blad`, więc wywołujący nie musi zakładać `try`. Gdy rdzeń nie odpowie
 * w ogóle — bo połączenie padło w trakcie — obietnica pozostaje
 * nierozstrzygnięta: kontrakt nie przewiduje limitu czasu, a rozłączenie
 * klienta nie kończy pracy rdzenia nad poleceniem.
 */
export function wywolaj<K extends Command>(
  kanal: Kanal,
  komenda: K,
  zadanie: RequestOf<K>,
): Promise<Wynik<ResponseOf<K>>> {
  return new Promise((rozstrzygnij) => {
    kanal.wyslij(komenda, zadanie, rozstrzygnij);
  });
}
