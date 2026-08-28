import { ChangeKind, Command, EventType } from '../../../shared/contract';
import { utworzMagistrale, type Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import type { Kanal } from '../protokol/kanal';
import { czyTablica, sprawdzKsztalt } from '../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../protokol/wywolanie';
import type { WpisSesji } from '../strona-glowna/zrodlo-sesji';

/**
 * Źródło wpisu sesji bieżącego okna rozmowy, prowadzące od identyfikatora sesji do
 * jej pełnego wpisu, którego gniazdo samo nie zna.
 */
export interface ZrodloWpisuSesji {
  /** Wpis sesji okna; `null`, dopóki rdzeń go nie oddał. */
  wpis(): WpisSesji | null;
  /** Zgłasza każdą zmianę wpisu — także pierwsze jego pojawienie się. */
  naZmiane(sluchacz: () => void): Odsubskrybuj;
  /** Odpytanie po zmianie sesji albo gdy identyfikator dopiero się pojawił. */
  odswiez(): void;
  /** Zdejmuje subskrypcję zdarzeń rdzenia; wołający musi to zrobić. */
  zamknij(): void;
}

/**
 * Śledzi jedną sesję po identyfikatorze przekazanym funkcją, bo sesja powstaje dopiero
 * po uzgodnieniu i w chwili montażu gniazda bywa jeszcze pusta.
 */
export function sledzWpisSesji(kanal: Kanal, idSesji: () => string): ZrodloWpisuSesji {
  const zmiany = utworzMagistrale<void>();
  let biezacy: WpisSesji | null = null;
  let numerZapytania = 0;

  /** Odpytuje rdzeń o wykaz sesji i wyjmuje z niego naszą, odrzucając odpowiedzi przedawnione. */
  function odpytaj(): void {
    const id = idSesji();
    if (id === '') return;

    const numer = (numerZapytania += 1);
    void wywolaj(kanal, Command.SessionList, { includePresence: true }).then((wynik) => {
      if (numer !== numerZapytania) return;
      const sprawdzony = sprawdzKsztalt(wynik, Command.SessionList, (tresc) =>
        czyTablica(tresc.sessions),
      );
      // Odmowa nie zeruje wpisu; menu pokazuje ostatnią znaną prawdę, nie błąd łączności.
      if (!sprawdzony.udany) return;

      const sesja = (sprawdzony.wynik?.sessions ?? []).find((pozycja) => pozycja.id === id);
      if (sesja === undefined) return;
      const obecnosc = (sprawdzony.wynik?.presence ?? []).find(
        (odpis) => odpis.sessionId === id,
      );
      biezacy = obecnosc === undefined ? { sesja } : { sesja, obecnosc };
      zmiany.oglos();
    });
  }

  // Zmiana sesji przychodzi zdarzeniem; usunięcie zdejmuje wpis, bo czynność traci cel.
  const odsubskrybuj = kanal.naZdarzenie(EventType.SessionChanged, (zdarzenie) => {
    if (zdarzenie.session.id !== idSesji()) return;
    if (zdarzenie.change === ChangeKind.Deleted) {
      biezacy = null;
    } else {
      const obecnosc = zdarzenie.presence ?? biezacy?.obecnosc;
      biezacy = obecnosc === undefined
        ? { sesja: zdarzenie.session }
        : { sesja: zdarzenie.session, obecnosc };
    }
    zmiany.oglos();
  });

  odpytaj();

  return {
    wpis: () => biezacy,
    naZmiane: (sluchacz) => zmiany.subskrybuj(sluchacz),
    odswiez: odpytaj,

    zamknij() {
      odsubskrybuj();
      // Licznik przesunięty zawczasu unieważnia odpowiedź w drodze po zejściu gniazda ze sceny.
      numerZapytania += 1;
    },
  };
}
