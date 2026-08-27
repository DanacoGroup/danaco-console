/**
 * Rozpoznanie kanału obrazowego po stronie klienta: czy danym kanałem da się
 * wygenerować obraz. Klucz adaptera bierze się z parametru `adapter`, a w jego
 * braku z rodzaju wiersza kanału; rozstrzygnięcie ostateczne należy do rdzenia.
 */
import type { Channel } from '../../../../shared/contract';

/**
 * Klucz adaptera kanału obrazowego. Napis rozstrzyga porównaniem dosłownym, więc
 * kanał podaje go w parametrze `adapter` albo w rodzaju wiersza dokładnie w tym
 * brzmieniu.
 */
export const ADAPTER_OBRAZY = 'obrazy';

/**
 * Klucz adaptera kanału — parametr `adapter`, a w jego braku rodzaj wiersza.
 *
 * Pusty wynik znaczy „kanał nie mówi, czym jest": ani parametru, ani rodzaju.
 * To inny stan niż „kanał tekstowy" i wykaz mówi o nim osobno.
 */
export function kluczAdaptera(kanal: Channel): string {
  const parametr = parametrKonfiguracji(kanal.config, 'adapter');
  if (parametr !== '') return parametr;
  return kanal.kind.trim();
}

/**
 * Czy kanał oddaje bajty obrazu — jedyny rodzaj, którym generowanie przejdzie.
 * Porównuje klucz adaptera kanału z kluczem obrazowym, więc kanał bez parametru
 * i bez rodzaju wypada poza wykaz.
 */
export function czyKanalObrazowy(kanal: Channel): boolean {
  return kluczAdaptera(kanal) === ADAPTER_OBRAZY;
}

/**
 * Kanały obrazowe rejestru w kolejności, w jakiej oddał je rdzeń. Przesiewa wykaz
 * wejściowy, nie kopiuje wierszy i nie zmienia ich porządku.
 */
export function kanalyObrazowe(kanaly: readonly Channel[]): readonly Channel[] {
  return kanaly.filter(czyKanalObrazowy);
}

/**
 * Dlaczego kanał nie nadaje się na silnik obrazów — zdanie dla wykazu.
 *
 * Puste znaczy „nadaje się". Zdanie powtarza powód, którym odmówiłby rdzeń, więc
 * powód jest znany bez zlecania generowania.
 */
export function powodNieprzydatnosci(kanal: Channel): string {
  if (czyKanalObrazowy(kanal)) return '';
  const klucz = kluczAdaptera(kanal);
  if (klucz === '') {
    return 'kanał nie podaje ani parametru adapter, ani rodzaju — rdzeń nie ma po czym poznać, czy odda obraz';
  }
  return `adapter „${klucz}" oddaje fragmenty tekstu, nie bajty obrazu`;
}

/**
 * Odczyt parametru kanału z pola `config` kontraktu.
 *
 * Pole jest w kontrakcie treścią nieokreśloną (`config?: unknown`), bo niesie
 * parametry dowolnego adaptera. Wszystko, co nie jest obiektem z napisem pod tym
 * kluczem, daje pustkę.
 */
function parametrKonfiguracji(config: unknown, klucz: string): string {
  if (config === null || typeof config !== 'object') return '';
  const wartosc = (config as Record<string, unknown>)[klucz];
  return typeof wartosc === 'string' ? wartosc.trim() : '';
}
