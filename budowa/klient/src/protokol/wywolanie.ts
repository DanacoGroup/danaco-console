import type { Command, RequestOf, ResponseOf } from '../../../shared/contract.ts';
import type { Kanal, Wynik } from './kanal.ts';

/**
 * Wywołanie komendy kontraktu jako obietnica. Kanał rozdaje wynik przez
 * wywołanie zwrotne, ale wołający pyta punktowo, dla którego obietnica jest
 * formą naturalną. Obietnica nie jest odrzucana nigdy — niepowodzenie wraca
 * jako wynik z opisem błędu.
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
