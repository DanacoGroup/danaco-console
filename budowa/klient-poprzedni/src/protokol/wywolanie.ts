import type { Command, RequestOf, ResponseOf } from '../../../shared/contract';
import type { Kanal, Wynik } from './kanal';

/** Wywołanie komendy kontraktu jako obietnica, rozstrzyganej zawsze wynikiem, nigdy odrzuceniem obietnicy. */
export function wywolaj<K extends Command>(
  kanal: Kanal,
  komenda: K,
  zadanie: RequestOf<K>,
): Promise<Wynik<ResponseOf<K>>> {
  return new Promise((rozstrzygnij) => {
    kanal.wyslij(komenda, zadanie, rozstrzygnij);
  });
}
