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

// Sekcja czynności sesji to druga sekcja menu nagłówka okna rozmowy, doklejana tylko gdy niepusta.

export interface OpcjeSekcjiCzynnosci {
  /** Droga do rdzenia; bez niej sekcja nie ma czym wykonać ani jednej czynności. */
  kanal: Kanal;
  /** Tryby widoku transkryptu z potoku rozmowy; pominięty znaczy, że gniazdo jeszcze nie ma okna. */
  widokZapisu?: PortWidokuZapisu;
  /** Dojście do liczby gniazd sceny pod pozycję otwarcia w nowym oknie; pominięte znaczy brak sceny. */
  noweOkno?: PortNowegoOkna;
}

export interface SekcjaCzynnosciSesji {
  /** Element doklejany do `sekcjeDalsze` menu paneli. */
  element: HTMLElement;
  /** Przerysowuje wiersze — po zmianie stanu sesji, trybu zapisu albo liczby gniazd. */
  odswiez(): void;
  /** Port trybów widoku transkryptu przychodzi później niż gniazdo; pustka zabiera wiersz z powrotem. */
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

    // Otwarcie w nowym oknie nie pyta rdzenia o stan sesji, tylko dostawia gniazdo na scenie.
    const noweOkno =
      opcje.noweOkno === undefined
        ? null
        : pozycjaNowegoOkna(opcje.noweOkno, () => przerysuj());

    // Nagłówek znika razem z pustym wykazem, bo napis nad pustką mówiłby, że coś tu jest, a nie ma.
    naglowek.hidden = pozycje.length === 0 && podmenu === null && noweOkno === null;

    // Kolejność wzorca ustala pozycje menu; podmenu wstawiamy po pierwszej czynności wpisu sesji.
    const wiersze: HTMLElement[] = pozycje.map((pozycja) => zbudujWierszCzynnosci(pozycja));
    if (podmenu !== null) {
      podmenu.odswiez();
      wiersze.splice(Math.min(1, wiersze.length), 0, podmenu.element);
    }
    if (noweOkno !== null) wiersze.unshift(zbudujWierszCzynnosci(noweOkno));
    wykaz.replaceChildren(...wiersze);
  }

  /** Skrót klawiaturowy sekcji działa bez modyfikatorów tylko wtedy, gdy ognisko siedzi w menu. */
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
      // Podmenu poprzedniego portu trzeba zwinąć przed porzuceniem, bo zwinięcie jest tanie i jawne.
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

/** Nagłówek sekcji czynności sesji — drobny, nie wersalikami, tak jak nagłówek sekcji paneli — pokazuje sekcję menu okna rozmowy. */
const NAGLOWEK = 'Sesja';
