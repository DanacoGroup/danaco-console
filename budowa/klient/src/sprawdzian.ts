/**
 * Bieg sprawdzianów klienta.
 *
 * Klient nie ma zależności zewnętrznych, a wbudowany biegacz `node:test` nie da
 * się przy tym użyć: jego deklaracje typów mieszkają w pakiecie `@types/node`,
 * którego nowy klient nie zaciąga, więc `tsc --noEmit` odmówiłby każdemu
 * plikowi sprawdzianu. Bieg poniżej korzysta wyłącznie z tego, co daje samo
 * środowisko uruchomieniowe.
 *
 * Zbiór pusty jest tu niepowodzeniem, nie wynikiem: sprawdzian, który niczego
 * nie zmierzył, milczałby dokładnie tak samo jak sprawdzian zdany.
 */

/** Przerywa sprawdzian, gdy warunek nie zachodzi. */
export function sprawdz(warunek: boolean, opis: string): void {
  if (!warunek) throw new Error(opis);
}

/** Przerywa sprawdzian, gdy wartości różnią się zapisem JSON. */
export function rowne(otrzymane: unknown, oczekiwane: unknown, opis: string): void {
  const a = JSON.stringify(otrzymane);
  const b = JSON.stringify(oczekiwane);
  if (a !== b) throw new Error(`${opis}: otrzymano ${a}, oczekiwano ${b}`);
}

/** Czeka, aż zaplanowane wywołania zwrotne zdążą się wykonać. */
export function poOdstepie(ms = 0): Promise<void> {
  return new Promise((rozstrzygnij) => setTimeout(rozstrzygnij, ms));
}

/**
 * Wykonuje zbiór sprawdzianów i wypisuje wynik każdego z nich.
 *
 * Niepowodzenie choćby jednego kończy proces wyjątkiem — kod wyjścia różny od
 * zera jest jedynym sygnałem, którego wywołujący nie przeoczy.
 */
export async function bieg(
  zbior: string,
  sprawdziany: Record<string, () => void | Promise<void>>,
): Promise<void> {
  const nazwy = Object.keys(sprawdziany);
  if (nazwy.length === 0) {
    throw new Error(`${zbior}: zbiór sprawdzianów jest pusty — nie zmierzono niczego`);
  }
  console.log(`# ${zbior} — ${nazwy.length} sprawdzianów`);
  const niezdane: string[] = [];
  for (const nazwa of nazwy) {
    try {
      await sprawdziany[nazwa]?.();
      console.log(`  zdany:   ${nazwa}`);
    } catch (blad) {
      niezdane.push(nazwa);
      console.log(`  NIEZDANY: ${nazwa}`);
      console.log(`            ${blad instanceof Error ? blad.message : String(blad)}`);
    }
  }
  console.log(`# ${zbior} — zdane ${nazwy.length - niezdane.length} z ${nazwy.length}`);
  if (niezdane.length > 0) {
    throw new Error(`${zbior}: niezdane sprawdziany — ${niezdane.join(', ')}`);
  }
}
