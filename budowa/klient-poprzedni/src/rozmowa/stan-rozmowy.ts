/**
 * Stan wysyłania rozmowy, jawny, żeby Operator widział biegnącą turę zanim przyjdzie
 * pierwszy znak odpowiedzi.
 */
export interface StanRozmowy {
  /** Czy tura biegnie — od wysłania komendy do znacznika końca strumienia. */
  wysyla: boolean;
  /** Wiadomość, której dotyczy bieżąca tura; pusta przed pierwszym fragmentem. */
  idWiadomosci: string;
  /** Liczba fragmentów odebranych w bieżącej turze. */
  fragmenty: number;
  // Ile pełnych sekund minęło od ostatniego znaku życia tury; poza turą równa się zeru.
  ciszaSekundy: number;
}

/**
 * Próg sekund milczenia strumienia, po którym okno mówi o ciszy wprost zamiast pokazywać
 * stały napis wysyłania.
 */
export const PROG_CISZY_SEKUND = 20;

/**
 * Opis stanu wysyłania rozmowy dla Operatora, przekazywany zawsze słowem, nigdy samą barwą
 * interfejsu.
 */
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
