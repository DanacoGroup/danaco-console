/**
 * Stan wysyłania rozmowy.
 *
 * Stan jest jawny, ponieważ Operator ma widzieć, że tura biegnie, zanim
 * przyjdzie pierwszy znak odpowiedzi. Nie jest bramką: przycisk zatrzymania
 * pozostaje czynny w każdym stanie, a pole wpisywania nigdy nie jest
 * wyszarzane.
 */
export interface StanRozmowy {
  /** Czy tura biegnie — od wysłania komendy do znacznika końca strumienia. */
  wysyla: boolean;
  /** Wiadomość, której dotyczy bieżąca tura; pusta przed pierwszym fragmentem. */
  idWiadomosci: string;
  /** Liczba fragmentów odebranych w bieżącej turze. */
  fragmenty: number;
  /**
   * Ile pełnych sekund minęło od ostatniego znaku życia biegnącej tury —
   * wysłania komendy albo ostatniego fragmentu. Poza turą równa się zeru.
   */
  ciszaSekundy: number;
}

/**
 * Po ilu sekundach milczenia strumienia mówimy o nim wprost.
 *
 * Próg jest wysoki, bo model bywa cichy, kiedy rozumuje — alarm po kilku
 * sekundach byłby fałszywy. Po upływie progu przyczyna (zerwane łącze, martwy
 * proces kanału, długie rozumowanie) nie ma znaczenia: Operator dostaje opis
 * ciszy zamiast niezmiennego napisu „Wysyłanie…".
 */
export const PROG_CISZY_SEKUND = 20;

/** Opis stanu dla Operatora — zawsze słowem, nigdy samą barwą. */
export function opisStanu(stan: StanRozmowy): string {
  if (!stan.wysyla) return 'Gotowe';
  if (stan.ciszaSekundy >= PROG_CISZY_SEKUND) {
    return (
      `Cisza w strumieniu — od ${stan.ciszaSekundy} s nie przyszedł żaden fragment ` +
      `(odebrano: ${stan.fragmenty})`
    );
  }
  if (stan.fragmenty === 0) return 'Wysyłanie…';
  return `Odpowiedź w strumieniu — fragmentów: ${stan.fragmenty}`;
}
