/**
 * Powód nieczynnej kontrolki składany z wykazu komend kontraktu, a nie wpisany
 * na stałe. Zdanie powstaje z odczytu wykazu przy składaniu okna, więc dopisanie
 * komendy do kontraktu przepisuje je samo.
 */

import { KOMENDY } from '../../../../shared/contract';

/**
 * Odczytuje z wykazu komend kontraktu komendy jednego obszaru: bierze te, które
 * noszą przedrostek obszaru, i zwraca je uporządkowane rosnąco.
 */
function komendyObszaru(obszar: string): readonly string[] {
  const przedrostek = `${obszar}.`;
  return [...(KOMENDY as readonly string[])].filter((komenda) => komenda.startsWith(przedrostek)).sort();
}

/**
 * Zdanie powodu dla kontrolki bez pokrycia w kontrakcie.
 *
 * @param czegoByTrzeba czynność, której kontrolka miała dokonać; zdanie podaje
 *   okno, bo ono zna przeznaczenie tej pozycji.
 * @param obszar przedrostek komend, w którym takiej komendy się szuka.
 */
export function powodBezKomendy(czegoByTrzeba: string, obszar = 'diagnostics'): string {
  const komendy = komendyObszaru(obszar);
  if (komendy.length === 0) {
    return `${czegoByTrzeba} Kontrakt nie niesie w obszarze ${obszar} ani jednej komendy — sprawdzone w wykazie kontraktu przy składaniu okna.`;
  }
  return `${czegoByTrzeba} Kontrakt niesie dziś w obszarze ${obszar} wyłącznie: ${komendy.join(', ')} — sprawdzone w wykazie kontraktu przy składaniu okna, nie wpisane na stałe.`;
}
