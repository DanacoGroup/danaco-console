import type { SettingDefinition } from '../../../shared/contract';
import { napis } from './kontrolka';
import type { StanKonfiguracji } from './stan-konfiguracji';

/**
 * Reguła widoczności pola sterowana katalogiem (`visibleWhenKey`,
 * `visibleWhenValue`). Warunek żyje w katalogu, nie w kodzie interfejsu.
 * Warunek wskazujący klucz, którego katalog nie zna, nie chowa pola.
 */
export function polePozostajeWidoczne(
  definicja: SettingDefinition,
  stan: StanKonfiguracji,
): boolean {
  const klucz = definicja.visibleWhenKey;
  if (klucz === undefined || klucz === '') return true;

  const warunkujaca = stan.definicjaKlucza(klucz);
  if (warunkujaca === null) {
    console.warn('[konfiguracja] warunek widoczności wskazuje klucz spoza katalogu', klucz);
    return true;
  }

  const obowiazujaca = napis(stan.rozstrzygnij(warunkujaca).wartosc);
  return obowiazujaca === (definicja.visibleWhenValue ?? '');
}
