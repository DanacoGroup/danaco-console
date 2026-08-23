import { EventType, type SessionFocusChangedEvent } from '../../../shared/contract';
import { utworzMagistrale, type Odsubskrybuj } from './magistrala-zdarzen';
import type { ZrodloZdarzen } from './zrodlo-zdarzen';

/**
 * Obserwator zdarzenia `session.focus.changed`.
 *
 * Zdarzenie jest nośnikiem synchronizacji ogniska między urządzeniami konta:
 * ognisko należy do klienta, więc rdzeń rozgłasza je z `clientId`, a odbiorca
 * rozstrzyga, czy zmiana dotyczy jego karty. Obserwator trzyma ostatni znany
 * stan ogniska i kieruje zmianę do właściwego odbiorcy — albo do wszystkich
 * (`naZmiane`), albo do słuchacza jednej karty sesji (`naSesje`).
 *
 * Ładunek o kształcie niezgodnym z kontraktem trafia do dziennika i jest
 * pomijany; stan ogniska pozostaje wtedy taki, jaki był.
 */
export interface ObserwatorOgniska {
  /** Sesja ogniskowana ostatnio; pusta, dopóki nie padło pierwsze zdarzenie. */
  ogniskowanaSesja(): string;
  /** Okno ogniskowane wewnątrz sesji; puste, gdy zdarzenie go nie wskazało. */
  ogniskowaneOkno(): string;
  /** Subskrypcja każdej zmiany ogniska. */
  naZmiane(sluchacz: (zmiana: SessionFocusChangedEvent) => void): Odsubskrybuj;
  /** Subskrypcja zmian dotyczących jednej karty sesji — wejściem albo wyjściem ogniska. */
  naSesje(idSesji: string, sluchacz: (zmiana: SessionFocusChangedEvent) => void): Odsubskrybuj;
  /** Odłącza obserwatora od źródła zdarzeń. */
  odlacz(): void;
}

export function utworzObserwatorOgniska(zrodlo: ZrodloZdarzen): ObserwatorOgniska {
  const zmiany = utworzMagistrale<SessionFocusChangedEvent>();
  let sesja = '';
  let okno = '';

  const odsubskrybuj = zrodlo.naZdarzenie(EventType.SessionFocusChanged, (tresc) => {
    if (!czyZmianaOgniska(tresc)) {
      console.warn('[łączność] zmiana ogniska o niespodziewanym kształcie', tresc);
      return;
    }
    zapamietaj(tresc);
    zmiany.oglos(tresc);
  });

  /**
   * Zapis stanu ogniska.
   *
   * Okno bez wskazania zostaje dotychczasowe, dopóki ognisko nie przechodzi do
   * innej karty — wtedy poprzednie okno przestaje być prawdą i pole gaśnie.
   */
  function zapamietaj(zmiana: SessionFocusChangedEvent): void {
    const wskazane = zmiana.windowId;
    if (wskazane !== undefined && wskazane.length > 0) {
      okno = wskazane;
    } else if (zmiana.sessionId !== sesja) {
      okno = '';
    }
    sesja = zmiana.sessionId;
  }

  return {
    ogniskowanaSesja: () => sesja,
    ogniskowaneOkno: () => okno,

    naZmiane: (sluchacz) => zmiany.subskrybuj(sluchacz),

    naSesje(idSesji, sluchacz) {
      return zmiany.subskrybuj((zmiana) => {
        if (zmiana.sessionId !== idSesji && zmiana.previousSessionId !== idSesji) return;
        sluchacz(zmiana);
      });
    },

    odlacz: odsubskrybuj,
  };
}

/** Pola obowiązkowe ładunku zdarzenia — sprawdzane, bo kanał ładunku nie waliduje. */
function czyZmianaOgniska(wartosc: unknown): wartosc is SessionFocusChangedEvent {
  if (typeof wartosc !== 'object' || wartosc === null) return false;
  const kandydat = wartosc as Record<string, unknown>;
  return (
    typeof kandydat['sessionId'] === 'string' &&
    typeof kandydat['clientId'] === 'string' &&
    typeof kandydat['focusedAt'] === 'number'
  );
}
