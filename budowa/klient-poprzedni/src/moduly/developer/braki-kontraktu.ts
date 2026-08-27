/**
 * Powód nieczynnej kontrolki modułu Developer składany z kontraktu, nie wpisany
 * na stałe. Wykaz `KOMENDY` z `shared/contract.ts` powstaje z `contract.json`,
 * więc dopisanie komendy do kontraktu przepisuje zdanie powodu samo.
 */
import { KOMENDY } from '../../../../shared/contract';



/**
 * Komendy obszaru odczytane z wykazu `KOMENDY` w czasie działania, zawężone
 * przedrostkiem obszaru i uporządkowane rosnąco, gotowe do wypisania w zdaniu
 * powodu.
 */
function komendyObszaru(obszar: string): readonly string[] {
  const przedrostek = `${obszar}.`;
  return [...(KOMENDY as readonly string[])]
    .filter((komenda) => komenda.startsWith(przedrostek))
    .sort();
}

/**
 * Zdanie powodu dla kontrolki bez pokrycia w kontrakcie; obszar bez komend
 * i obszar z komendami dają dwa różne zdania.
 * @param czegoByTrzeba czynność, której kontrolka miała dokonać.
 * @param obszar przedrostek komend, w którym takiej komendy szukamy.
 */
export function powodBezKomendy(czegoByTrzeba: string, obszar = 'developer'): string {
  const komendy = komendyObszaru(obszar);
  if (komendy.length === 0) {
    return (
      `${czegoByTrzeba} Kontrakt nie niesie w obszarze ${obszar} ani jednej komendy — ` +
      'sprawdzone w wykazie kontraktu przy składaniu okna.'
    );
  }
  return (
    `${czegoByTrzeba} Kontrakt niesie dziś w obszarze ${obszar} wyłącznie: ${komendy.join(', ')} — ` +
    'sprawdzone w wykazie kontraktu przy składaniu okna, nie wpisane na stałe.'
  );
}
