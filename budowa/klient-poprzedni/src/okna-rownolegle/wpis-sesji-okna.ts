import { ChangeKind, Command, EventType } from '../../../shared/contract';
import { utworzMagistrale, type Odsubskrybuj } from '../polaczenie/magistrala-zdarzen';
import type { Kanal } from '../protokol/kanal';
import { czyTablica, sprawdzKsztalt } from '../protokol/ksztalt-odpowiedzi';
import { wywolaj } from '../protokol/wywolanie';
import type { WpisSesji } from '../strona-glowna/zrodlo-sesji';

/**
 * Sesja tego okna rozmowy — jeden wpis, żywy.
 *
 * Sekcja czynności menu gniazda rozstrzyga o pozycjach na podstawie `WpisSesji`
 * (stan sesji, tytuł, projekt, liczba okien strumieniujących), a gniazdo zna
 * wyłącznie identyfikator sesji; to źródło prowadzi od identyfikatora do wpisu.
 *
 * `strona-glowna/zrodlo-sesji.ts` się tu nie nadaje: obsługuje sekcję sesji
 * w tle i odcina sesję bieżącego połączenia, czyli dokładnie tę, o którą pyta
 * menu. Komenda jest ta sama — `session.list` z żywym stanem — różni się filtr.
 *
 * Źródło nie zna DOM-u, nie wykonuje żadnej czynności i nie ma zdania o tym, co
 * z sesją wolno zrobić: oddaje wpis albo `null`. `null` znaczy brak wiedzy, nie
 * zakaz — dopóki rdzeń nie oddał wykazu, sekcja czynności jest krótsza, bo
 * czynności bez znanego stanu sesji nie da się uczciwie nazwać.
 */
export interface ZrodloWpisuSesji {
  /** Wpis sesji okna; `null`, dopóki rdzeń go nie oddał. */
  wpis(): WpisSesji | null;
  /** Zgłasza każdą zmianę wpisu — także pierwsze jego pojawienie się. */
  naZmiane(sluchacz: () => void): Odsubskrybuj;
  /**
   * Ponowne odpytanie rdzenia.
   *
   * Potrzebne, gdy identyfikator sesji dopiero się pojawił — gniazdo montuje
   * się przed uzgodnieniem — oraz po czynności, która sesję zmieniła.
   */
  odswiez(): void;
  /** Zdejmuje subskrypcję zdarzeń rdzenia; wołający musi to zrobić. */
  zamknij(): void;
}

/**
 * Śledzi jedną sesję po identyfikatorze.
 *
 * Identyfikator przychodzi funkcją, nie napisem: sesja powstaje po uzgodnieniu,
 * więc w chwili montażu gniazda bywa jeszcze pusta. Odpytanie rusza dopiero,
 * gdy identyfikator jest niepusty — żądanie o sesję „" nie miałoby o co pytać.
 */
export function sledzWpisSesji(kanal: Kanal, idSesji: () => string): ZrodloWpisuSesji {
  const zmiany = utworzMagistrale<void>();
  let biezacy: WpisSesji | null = null;
  let numerZapytania = 0;

  /**
   * Odpytanie rdzenia o wykaz i wyjęcie z niego naszej sesji.
   *
   * Odpowiedź przedawniona nie nadpisuje świeższej (licznik `numerZapytania`) —
   * dwa odpytania w locie potrafią wrócić w odwrotnej kolejności, a wtedy
   * starszy stan przykryłby nowszy.
   */
  function odpytaj(): void {
    const id = idSesji();
    if (id === '') return;

    const numer = (numerZapytania += 1);
    void wywolaj(kanal, Command.SessionList, { includePresence: true }).then((wynik) => {
      if (numer !== numerZapytania) return;
      const sprawdzony = sprawdzKsztalt(wynik, Command.SessionList, (tresc) =>
        czyTablica(tresc.sessions),
      );
      // Odmowa nie zeruje wpisu, który już mamy: menu ma wtedy pokazywać ostatnią
      // znaną prawdę, a nie zapadać się przy pierwszym potknięciu łączności.
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

  // Zmiana sesji przychodzi zdarzeniem, więc menu nie musi odpytywać w kółko.
  // Usunięcie zdejmuje wpis: sesji już nie ma i żadna czynność nie ma celu.
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
      // Licznik przesunięty zawczasu unieważnia odpowiedź, która jest jeszcze
      // w drodze — inaczej zapisałaby wpis po zejściu gniazda ze sceny.
      numerZapytania += 1;
    },
  };
}
