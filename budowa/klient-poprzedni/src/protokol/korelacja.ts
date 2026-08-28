import type { Envelope } from '../../../shared/contract';

/** Wiązanie odpowiedzi z żądaniem po identyfikatorze koperty, bez limitu czasu ani limitu żądań oczekujących. */
export type OdbiorcaOdpowiedzi = (odpowiedz: Envelope) => void;

export interface Korelacja {
  /** Zapamiętuje odbiorcę odpowiedzi na żądanie o podanym identyfikatorze. */
  zarejestruj(idZadania: string, odbiorca: OdbiorcaOdpowiedzi): void;
  /** Przekazuje odpowiedź jej odbiorcy; zwraca prawdę, gdy odbiorca istniał. */
  rozstrzygnij(odpowiedz: Envelope): boolean;
  /** Liczba żądań oczekujących na odpowiedź. */
  oczekujace(): number;
}

export function utworzKorelacje(): Korelacja {
  const odbiorcy = new Map<string, OdbiorcaOdpowiedzi>();

  return {
    zarejestruj(idZadania, odbiorca) {
      if (idZadania.length === 0) return;
      odbiorcy.set(idZadania, odbiorca);
    },

    rozstrzygnij(odpowiedz) {
      const odbiorca = odbiorcy.get(odpowiedz.id);
      if (odbiorca === undefined) return false;
      odbiorcy.delete(odpowiedz.id);
      odbiorca(odpowiedz);
      return true;
    },

    oczekujace: () => odbiorcy.size,
  };
}
