import type { SettingDefinition } from '../../../shared/contract';
import { napis } from './kontrolka';
import type { StanKonfiguracji } from './stan-konfiguracji';

/**
 * Reguła widoczności pola sterowana katalogiem (`visibleWhenKey`,
 * `visibleWhenValue`).
 *
 * Pozycja katalogu może zależeć od innej: pole „host wykonania" ma sens
 * dopiero przy zasięgu zdalnym, pole „ścieżka klucza" dopiero przy moście SSH.
 * Warunek żyje w katalogu, nie w kodzie interfejsu — zmiana zależności to nowy
 * wiersz, nie nowa gałąź w kliencie.
 *
 * Warunek wskazujący klucz, którego katalog nie zna, nie chowa pola. Ukrycie
 * pola z powodu braku metadanej byłoby cichym odebraniem Operatorowi
 * ustawienia.
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
