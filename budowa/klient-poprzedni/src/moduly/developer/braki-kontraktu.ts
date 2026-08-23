import { KOMENDY } from '../../../../shared/contract';

/**
 * Powód nieczynnej kontrolki modułu Developer składany z kontraktu, nie wpisany
 * na stałe.
 *
 * `shared/contract.ts` niesie `KOMENDY` — wykaz komend kontraktu dostępny
 * w czasie działania, wytwarzany z `contract.json`. Zdanie powstaje z odczytu
 * tego wykazu przy składaniu okna, więc dopisanie komendy do kontraktu
 * przepisuje je samo; zdanie wpisane na stałe przestałoby być prawdziwe w dniu
 * takiej zmiany i nikt by tego nie zauważył.
 *
 * Zdanie nie orzeka, czy złożony rdzeń komendę rejestruje — to osobne pytanie
 * i odpowiada na nie `katalog-komend.ts`, który pyta rdzeń o jego rejestr.
 * Tutaj brak jest po stronie kontraktu i tylko o kontrakcie zdanie mówi.
 *
 * Bliźniaczy mechanizm stoi w module Diagnostics. Wspólnego bytu biblioteka
 * `komponenty/` dziś nie ma, a modułowi nie wolno sięgać do wnętrza sąsiada —
 * rozstrzygnięcie, czy taki byt ma powstać w bibliotece, należy do właściciela
 * projektu.
 */

/** Komendy obszaru odczytane z kontraktu w czasie działania. */
function komendyObszaru(obszar: string): readonly string[] {
  const przedrostek = `${obszar}.`;
  return [...(KOMENDY as readonly string[])]
    .filter((komenda) => komenda.startsWith(przedrostek))
    .sort();
}

/**
 * Zdanie powodu dla kontrolki bez pokrycia w kontrakcie.
 *
 * @param czegoByTrzeba czynność, której kontrolka miała dokonać, wraz ze
 *   wskazaniem komendy, która musiałaby powstać — to okno wie, po co ta pozycja
 *   stoi w inwentarzu opracowania.
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
