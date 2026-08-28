import { Cel, wskazKatalog as wskazNatywnie } from '../powloka/most-katalogow';

/**
 * Przejściówka do mostu powłoki natywnej `powloka/most-katalogow.ts`.
 *
 * Wskazanie katalogu jest jedną czynnością systemu operacyjnego wywoływaną
 * w dwóch sprawach: katalogu roboczego i punktu dostępu, więc implementacja
 * jest jedna.
 */
export { Cel, czyPowlokaNatywna, naWskazanieZZasobnika } from '../powloka/most-katalogow';

/**
 * Wskazanie katalogu przez powłokę natywną, z jednym zapisem wyniku dla widoków
 * dostępów: napis pusty znaczy rezygnację Operatora albo brak powłoki natywnej,
 * a zamiana `null` na napis pusty jest jedyną treścią tej funkcji.
 */
export async function wskazKatalog(cel: Cel = Cel.KatalogRoboczy): Promise<string> {
  return (await wskazNatywnie(cel)) ?? '';
}
