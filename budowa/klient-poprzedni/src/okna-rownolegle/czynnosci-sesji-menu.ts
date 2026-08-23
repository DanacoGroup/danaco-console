import type { Kanal } from '../protokol/kanal';
import './czynnosci.css';
import {
  utworzPodmenuWidokuZapisu,
  type PodmenuWidokuZapisu,
  type PortWidokuZapisu,
} from './podmenu-widoku-zapisu';
import { pozycjaNowegoOkna, type PortNowegoOkna } from './pozycja-nowego-okna';
import { spisCzynnosciSesji } from './spis-czynnosci-sesji';
import { zbudujWierszCzynnosci, type PozycjaCzynnosciMenu } from './wiersz-czynnosci';
import { sledzWpisSesji, type ZrodloWpisuSesji } from './wpis-sesji-okna';

/**
 * Sekcja czynności sesji — druga sekcja menu `⋮` nagłówka okna rozmowy.
 *
 * Doklejana przez `OpcjeMenuPaneli.sekcjeDalsze`; kreska nad sekcją rysuje się
 * po stronie `menu-paneli.ts` i tylko wtedy, gdy sekcja niepusta. Kolejność
 * wierszy: `Otwórz w nowym oknie`, `Zmień nazwę`, `Widok transkryptu ›`,
 * `Archiwizuj`, `Usuń`. Pozycja bez pokrycia w rdzeniu nie powstaje w ogóle —
 * nie ma tu wiersza wygaszonego.
 *
 * Identyfikator sesji czytamy z kanału (`kanal.sesja().id()`) przy każdym
 * pytaniu, nie raz przy montażu: gniazdo powstaje przed uzgodnieniem z rdzeniem
 * i sesji wtedy jeszcze nie ma. Dopóki identyfikator jest pusty albo rdzeń nie
 * oddał wpisu, sekcja jest pusta.
 *
 * Skróty `R`, `A`, `D` łapie nasłuch tej sekcji i działają, dopóki menu jest
 * rozwinięte. Skrótu globalnego nie rejestrujemy: te same litery wpisane w polu
 * wypowiedzi mają pisać litery.
 */

export interface OpcjeSekcjiCzynnosci {
  /** Droga do rdzenia; bez niej sekcja nie ma czym wykonać ani jednej czynności. */
  kanal: Kanal;
  /**
   * Tryby widoku transkryptu podane przez potok `rozmowa/`.
   *
   * Pominięty znaczy „powłoka nie związała jeszcze gniazda z oknem" — wiersza
   * `Widok transkryptu ›` wtedy nie ma. Port dochodzi zwykle później, przez
   * `ustawWidokZapisu`: gniazdo powstaje przed uzgodnieniem, a rozmowę osadza
   * dopiero `aplikacja/wiazanie-gniazda.ts`.
   */
  widokZapisu?: PortWidokuZapisu;
  /**
   * Dojście do liczby gniazd sceny — pod pozycję `Otwórz w nowym oknie`.
   *
   * Pominięte znaczy „nie ma sceny, na którą dałoby się dostawić okno" —
   * pozycji wtedy nie ma.
   */
  noweOkno?: PortNowegoOkna;
}

export interface SekcjaCzynnosciSesji {
  /** Element doklejany do `sekcjeDalsze` menu paneli. */
  element: HTMLElement;
  /** Przerysowuje wiersze — po zmianie stanu sesji, trybu zapisu albo liczby gniazd. */
  odswiez(): void;
  /**
   * Podaje port trybów widoku transkryptu — albo go zabiera (`null`).
   *
   * Osobne wejście, bo port przychodzi później niż gniazdo: rozmowę osadza
   * `aplikacja/wiazanie-gniazda.ts` po uzgodnieniu z rdzeniem, a sekcja montuje
   * się razem z gniazdem. `null` przy rozłączeniu zabiera wiersz z powrotem,
   * żeby podmenu nie wołało trybów rozmowy, której już nie ma.
   */
  ustawWidokZapisu(port: PortWidokuZapisu | null): void;
  /** Zdejmuje subskrypcję rdzenia; wołane przy zejściu gniazda. */
  zamknij(): void;
}

export function utworzSekcjeCzynnosciSesji(
  opcje: OpcjeSekcjiCzynnosci,
): SekcjaCzynnosciSesji {
  const element = document.createElement('div');
  element.className = 'dn-czynnosci-sesji';
  element.setAttribute('role', 'group');
  element.setAttribute('aria-label', NAGLOWEK);

  const naglowek = document.createElement('p');
  naglowek.className = 'dn-czynnosci-sesji__naglowek';
  naglowek.textContent = NAGLOWEK;

  const wykaz = document.createElement('div');
  wykaz.className = 'dn-czynnosci-sesji__wykaz';
  element.append(naglowek, wykaz);

  const zrodlo: ZrodloWpisuSesji = sledzWpisSesji(opcje.kanal, () =>
    opcje.kanal.sesja().id(),
  );

  let podmenu: PodmenuWidokuZapisu | null =
    opcje.widokZapisu === undefined
      ? null
      : utworzPodmenuWidokuZapisu(opcje.widokZapisu);

  /** Pozycje ostatnio postawione — po nich idzie obsługa skrótów. */
  let pozycje: readonly PozycjaCzynnosciMenu[] = [];

  function przerysuj(): void {
    pozycje = spisCzynnosciSesji(zrodlo.wpis(), {
      kanal: opcje.kanal,
      odswiez: () => zrodlo.odswiez(),
    });

    // `Otwórz w nowym oknie` nie zależy od wpisu sesji, w odróżnieniu od
    // czterech pozostałych: nie pyta rdzenia o stan sesji, tylko dostawia
    // gniazdo na scenie. Pytanie o sufit idzie przy każdym rysowaniu, bo
    // figura modułu sceny przestawia się z rdzenia.
    const noweOkno =
      opcje.noweOkno === undefined
        ? null
        : pozycjaNowegoOkna(opcje.noweOkno, () => przerysuj());

    // Nagłówek znika razem z pustym wykazem: napis „Sesja" nad pustką mówiłby,
    // że coś tu jest, a nie ma. Kreski nad sekcją i tak wtedy nie ma, bo
    // `menu-paneli.ts` stawia ją wyłącznie przy niepustym `sekcjeDalsze`.
    naglowek.hidden = pozycje.length === 0 && podmenu === null && noweOkno === null;

    // Kolejność wzorca: `Otwórz w`, `Zmień nazwę`, `Widok transkryptu ›`,
    // `Archiwizuj`, `Usuń`. Podmenu wstawiamy po pierwszej czynności wpisu,
    // a przy jej braku na początek tamtego wykazu — dopiero potem na czoło
    // wchodzi „Otwórz w nowym oknie", którego wpis sesji nie dotyczy.
    const wiersze: HTMLElement[] = pozycje.map((pozycja) => zbudujWierszCzynnosci(pozycja));
    if (podmenu !== null) {
      podmenu.odswiez();
      wiersze.splice(Math.min(1, wiersze.length), 0, podmenu.element);
    }
    if (noweOkno !== null) wiersze.unshift(zbudujWierszCzynnosci(noweOkno));
    wykaz.replaceChildren(...wiersze);
  }

  /**
   * Skrót klawiaturowy sekcji.
   *
   * Litera bez modyfikatorów i tylko wtedy, gdy ognisko siedzi w menu — nasłuch
   * stoi na sekcji, więc zdarzenie musi do niej dojść. `Ctrl`/`Alt`/`Meta`
   * odpuszczamy, żeby nie odbierać przeglądarce jej własnych skrótów.
   */
  element.addEventListener('keydown', (zdarzenie) => {
    if (zdarzenie.ctrlKey || zdarzenie.altKey || zdarzenie.metaKey) return;
    if (zdarzenie.key.length !== 1) return;
    const litera = zdarzenie.key.toUpperCase();
    const pozycja = pozycje.find((wpis) => wpis.skrot === litera);
    if (pozycja === undefined) return;
    zdarzenie.preventDefault();
    zdarzenie.stopPropagation();
    pozycja.wykonaj();
  });

  const odsubskrybuj = zrodlo.naZmiane(przerysuj);
  przerysuj();

  return {
    element,
    odswiez: przerysuj,

    ustawWidokZapisu(port) {
      // Podmenu poprzedniego portu trzeba zwinąć przed porzuceniem: jego
      // wiersze trybów są ukrywane atrybutem `hidden`, po którym `menu-rozwijane.ts`
      // rozstrzyga o wędrówce ogniska. Element i tak wypada z drzewa przy
      // najbliższym `replaceChildren`, ale zwinięcie jest tu tanie i jawne.
      podmenu?.zwin();
      podmenu = port === null ? null : utworzPodmenuWidokuZapisu(port);
      przerysuj();
    },

    zamknij() {
      odsubskrybuj();
      zrodlo.zamknij();
      podmenu?.zwin();
    },
  };
}

/** Nagłówek sekcji — drobny, nie wersalikami, jak nagłówek sekcji paneli. */
const NAGLOWEK = 'Sesja';
