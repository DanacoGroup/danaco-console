/**
 * Bieg sprawdzianów klienta. Klient nie ma zależności zewnętrznych, a wbudowany
 * biegacz `node:test` nie da się użyć, więc bieg poniżej korzysta wyłącznie
 * z tego, co daje samo środowisko uruchomieniowe.
 */

/** Przerywa sprawdzian, gdy warunek nie zachodzi, zgłaszając wyjątek z podanym opisem niepowodzenia sprawdzianu. */
export function sprawdz(warunek: boolean, opis: string): void {
  if (!warunek) throw new Error(opis);
}

/** Przerywa sprawdzian, gdy wartości różnią się zapisem JSON, nazywając w komunikacie obie porównywane strony. */
export function rowne(otrzymane: unknown, oczekiwane: unknown, opis: string): void {
  const a = JSON.stringify(otrzymane);
  const b = JSON.stringify(oczekiwane);
  if (a !== b) throw new Error(`${opis}: otrzymano ${a}, oczekiwano ${b}`);
}

/** Czeka, aż zaplanowane wywołania zwrotne zdążą się wykonać, ustępując na wskazany czas w milisekundach. */
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
