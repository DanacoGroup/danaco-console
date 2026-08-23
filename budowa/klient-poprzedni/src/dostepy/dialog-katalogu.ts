import { Cel, wskazKatalog as wskazNatywnie } from '../powloka/most-katalogow';

/**
 * Przejściówka do mostu powłoki natywnej (`powloka/most-katalogow.ts`).
 *
 * Wskazanie katalogu jest jedną czynnością systemu operacyjnego wywoływaną
 * w dwóch sprawach: katalogu roboczego i punktu dostępu, więc implementacja
 * jest jedna. Most mieszka w `powloka/`, bo dotyczy powłoki natywnej, a nie
 * dostępów; tutaj zostaje wyłącznie zapis przyjęty w widokach dostępów:
 * napis pusty znaczy „nie wskazano”.
 *
 * Zamiana `null` na napis pusty jest jedyną treścią tego modułu.
 */

export { Cel, czyPowlokaNatywna, naWskazanieZZasobnika } from '../powloka/most-katalogow';

/** Wskazanie katalogu; napis pusty znaczy rezygnację albo brak powłoki. */
export async function wskazKatalog(cel: Cel = Cel.KatalogRoboczy): Promise<string> {
  return (await wskazNatywnie(cel)) ?? '';
}
