import { EventType, type SessionFocusChangedEvent } from '../../../shared/contract';
import { utworzMagistrale, type Odsubskrybuj } from './magistrala-zdarzen';
import type { ZrodloZdarzen } from './zrodlo-zdarzen';

/**
 * Obserwator zdarzenia session.focus.changed synchronizuje ognisko między urządzeniami konta i kieruje zmianę do właściwego odbiorcy albo do wszystkich, pomijając ładunek niezgodny z kontraktem.
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

  /** Okno bez wskazania zostaje dotychczasowe, dopóki ognisko nie przejdzie do innej karty sesji. */
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

/** Pola obowiązkowe ładunku zdarzenia — sprawdzane ręcznie, ponieważ sam kanał ładunku ich nie waliduje. */
function czyZmianaOgniska(wartosc: unknown): wartosc is SessionFocusChangedEvent {
  if (typeof wartosc !== 'object' || wartosc === null) return false;
  const kandydat = wartosc as Record<string, unknown>;
  return (
    typeof kandydat['sessionId'] === 'string' &&
    typeof kandydat['clientId'] === 'string' &&
    typeof kandydat['focusedAt'] === 'number'
  );
}
