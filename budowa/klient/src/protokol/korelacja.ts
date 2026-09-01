import {
  EnvelopeStatus,
  ErrorCode,
  type Envelope,
  type ErrorInfo,
  type MessageType,
} from '../../../shared/contract.ts';

/**
 * Wiązanie odpowiedzi z żądaniem po identyfikatorze koperty. Kontrakt
 * zapowiada dokładnie jedną odpowiedź na komendę, a odpowiedź powtarza
 * identyfikator żądania.
 */
export type OdbiorcaOdpowiedzi = (odpowiedz: Envelope) => void;

/*
Termin, po którym żądanie bez odpowiedzi zostaje rozstrzygnięte odmową.
Kontrakt nie zapowiada czasu odpowiedzi, a wywołanie zgubione w locie
zostawiłoby wołającemu obietnicę bez rozstrzygnięcia i okno czekające bez komunikatu.
*/
export const TERMIN_ODPOWIEDZI_MS = 30_000;

export interface Korelacja {
  /** Zapamiętuje odbiorcę odpowiedzi na żądanie o podanym identyfikatorze i typie. */
  zarejestruj(idZadania: string, typ: MessageType, odbiorca: OdbiorcaOdpowiedzi): void;
  /** Przekazuje odpowiedź jej odbiorcy; zwraca prawdę, gdy odbiorca istniał. */
  rozstrzygnij(odpowiedz: Envelope): boolean;
  /** Rozstrzyga oczekujące żądanie odmową; zwraca prawdę, gdy odbiorca istniał. */
  odmow(idZadania: string, blad: ErrorInfo): boolean;
  /** Rozstrzyga odmową wszystkie żądania oczekujące. */
  uniewaznijWszystkie(blad: ErrorInfo): void;
  /** Czy odpowiedź o tym identyfikatorze przyszła po terminie; wpis jest zdejmowany, bo odpowiedź jest jedna. */
  czySpozniona(idZadania: string): boolean;
  /** Liczba żądań oczekujących na odpowiedź. */
  oczekujace(): number;
}

/** Żądanie oczekujące na odpowiedź wraz z typem koperty i pilnującym go terminem. */
interface Wpis {
  typ: MessageType;
  odbiorca: OdbiorcaOdpowiedzi;
  termin: ReturnType<typeof setTimeout>;
}

export function utworzKorelacje(): Korelacja {
  const oczekujacy = new Map<string, Wpis>();
  /* Żądania rozstrzygnięte terminem, pamiętane przez jeszcze jeden termin:
     odpowiedź po nich jest spóźniona, nie obca, i kanał odrzuca ją bez doręczania. */
  const poTerminie = new Map<string, ReturnType<typeof setTimeout>>();

  /** Zdejmuje wpis wraz z jego terminem; zwraca go wołającemu albo pustkę, gdy już nie oczekiwał. */
  function zdejmij(idZadania: string): Wpis | undefined {
    const wpis = oczekujacy.get(idZadania);
    if (wpis === undefined) return undefined;
    clearTimeout(wpis.termin);
    oczekujacy.delete(idZadania);
    return wpis;
  }

  function zapamietajTermin(idZadania: string): void {
    poTerminie.set(
      idZadania,
      setTimeout(() => poTerminie.delete(idZadania), TERMIN_ODPOWIEDZI_MS),
    );
  }

  const korelacja: Korelacja = {
    zarejestruj(idZadania, typ, odbiorca) {
      if (idZadania.length === 0) return;
      const termin = setTimeout(() => {
        zapamietajTermin(idZadania);
        korelacja.odmow(idZadania, bladTerminu());
      }, TERMIN_ODPOWIEDZI_MS);
      oczekujacy.set(idZadania, { typ, odbiorca, termin });
    },

    rozstrzygnij(odpowiedz) {
      const wpis = zdejmij(odpowiedz.id);
      if (wpis === undefined) return false;
      wpis.odbiorca(odpowiedz);
      return true;
    },

    odmow(idZadania, blad) {
      const wpis = zdejmij(idZadania);
      if (wpis === undefined) return false;
      wpis.odbiorca(kopertaOdmowy(idZadania, wpis.typ, blad));
      return true;
    },

    uniewaznijWszystkie(blad) {
      for (const idZadania of [...oczekujacy.keys()]) {
        korelacja.odmow(idZadania, blad);
      }
    },

    czySpozniona(idZadania) {
      const zapomnienie = poTerminie.get(idZadania);
      if (zapomnienie === undefined) return false;
      clearTimeout(zapomnienie);
      poTerminie.delete(idZadania);
      return true;
    },

    oczekujace: () => oczekujacy.size,
  };

  return korelacja;
}

/** Koperta odpowiedzi składana przez klienta, gdy odpowiedzi rdzenia nie będzie; wywołujący czyta ją tak samo jak odmowę rdzenia. */
function kopertaOdmowy(idZadania: string, typ: MessageType, blad: ErrorInfo): Envelope {
  return {
    type: typ,
    id: idZadania,
    timestamp: Date.now(),
    status: EnvelopeStatus.Error,
    error: blad,
  };
}

/** Odmowa żądania, na które odpowiedź nie przyszła w terminie. */
function bladTerminu(): ErrorInfo {
  return {
    code: ErrorCode.ChannelUnavailable,
    message: 'Rdzeń nie odpowiedział w terminie',
    retryable: true,
  };
}
