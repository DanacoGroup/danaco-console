import type { ErrorInfo } from '../../../shared/contract';

/**
 * Komunikat o losie pojedynczej zmiany ustawienia.
 *
 * Niepowodzenie zmiany nie wyłącza sterowania ani nie zamyka okna — jest
 * wyłącznie informacją. Błąd techniczny dotyczy bieżącego wywołania, więc
 * kolejna próba idzie normalnie.
 */
export interface KomunikatZmiany {
  /** Treść pokazywana operatorowi. */
  tresc: string;
  /** Czy zmiana została potwierdzona przez rdzeń. */
  udany: boolean;
}

/**
 * Komunikat potwierdzenia zmiany ustawienia sterowania, wyświetlany w pasku komunikatów
 * tego okna aplikacji.
 */
export function potwierdzenie(nazwa: string): KomunikatZmiany {
  return { tresc: `${nazwa} — zapisane`, udany: true };
}

/**
 * Komunikat niepowodzenia zmiany.
 *
 * Kod błędu podajemy dosłownie; pole `retryable` kontraktu rozstrzyga, czy
 * ponowienie ma sens.
 */
export function niepowodzenie(nazwa: string, blad?: ErrorInfo): KomunikatZmiany {
  if (blad === undefined) {
    return { tresc: `${nazwa} — rdzeń nie potwierdził zmiany`, udany: false };
  }
  const ponowienie = blad.retryable ? ' · ponowienie ma sens' : '';
  return { tresc: `${nazwa} — ${blad.code}: ${blad.message}${ponowienie}`, udany: false };
}

/**
 * Odbiorca komunikatów kompletu sterowania tego okna, wywoływany przy każdej zmianie jego
 * ustawienia.
 */
export type OdbiorcaKomunikatu = (komunikat: KomunikatZmiany) => void;
