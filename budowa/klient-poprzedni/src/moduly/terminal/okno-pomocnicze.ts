import './terminal.css';

import { utworzRameOkna } from '../../komponenty/rama-okna';
import type { Kanal } from '../../protokol/kanal';
import { zamontujTerminal, type ZamontowanyTerminal } from './indeks';

/**
 * Terminal jako okno pomocnicze modułu gospodarza: nie jest samodzielnym
 * modułem ani niezależną sesją, tylko powierzchnią osadzoną w oknie gospodarza.
 */
export interface OknoTerminalaPomocnicze {
  /** Sekcja osadzana w złożeniu modułu gospodarza albo w pasie okien pomocniczych. */
  element: HTMLElement;
  /**
   * Podaje okno wykonania gospodarza; napis pusty znaczy brak okna od rdzenia.
   */
  ustawOkno(okno: string): void;
  /** Odczytuje okna terminala z rdzenia; bez okna gospodarza nie robi nic. */
  odswiez(): void;
  /** Zamyka nasłuch zdarzeń terminala założony przez to okno. */
  zamknij(): void;
}

/**
 * Zależności okna pomocniczego terminala, podawane przy jego złożeniu przez
 * moduł gospodarza, który to okno u siebie osadza.
 */
export interface OpcjeTerminalaPomocniczego {
  kanal: Kanal;
  /** Nazwa modułu gospodarza — wchodzi w etykiety i w zdania dla Operatora. */
  modul: string;
  /** Przedrostek klas gospodarza: `mdev`, `dg`, `mp`. Pominięty = same klasy biblioteki. */
  przedrostek?: string;
  /** Okno wykonania gospodarza, jeśli znane już przy tworzeniu. */
  okno?: string;
}

export function utworzOknoTerminalaPomocnicze(
  opcje: OpcjeTerminalaPomocniczego,
): OknoTerminalaPomocnicze {
  const rama = utworzRameOkna({
    tytul: 'Terminal',
    rola: 'pomocnicze',
    kod: 'terminal',
    modul: opcje.modul,
    przeznaczenie:
      `Karta powłoki wewnątrz modułu ${opcje.modul} — polecenia, procesy i wyjście na żywo. ` +
      'Terminal nie jest samodzielnym modułem: karta powstaje w OKNIE tego modułu, ' +
      'więc dziedziczy jego tryb uprawnień i katalog roboczy.',
    ...(opcje.przedrostek === undefined ? {} : { przedrostek: opcje.przedrostek }),
  });

  let zamontowany: ZamontowanyTerminal | null = null;
  let biezaceOkno = '';

  function pokazBrakOkna(): void {
    rama.cialo.replaceChildren(zdanieBezOkna(opcje.modul));
    rama.ustawZnacznik('bez okna gospodarza', 'ostrzezenie');
  }

  function ustawOkno(okno: string): void {
    const nowe = okno.trim();
    const stoi = zamontowany !== null;
    if (nowe === biezaceOkno && stoi === (nowe !== '')) return;

    zamontowany?.zamknij();
    zamontowany = null;
    biezaceOkno = nowe;

    if (nowe === '') {
      pokazBrakOkna();
      return;
    }
    // Montaż terminala sam podmienia zawartość gospodarza.
    zamontowany = zamontujTerminal(rama.cialo, opcje.kanal, { okno: nowe });
    rama.ustawZnacznik(`okno gospodarza: ${nowe}`, 'sukces');
  }

  ustawOkno(opcje.okno ?? '');

  return {
    element: rama.element,
    ustawOkno,
    odswiez: () => zamontowany?.odswiez(),
    zamknij() {
      zamontowany?.zamknij();
      zamontowany = null;
    },
  };
}

/**
 * Zdanie stanu opisujące gospodarza bez okna: nazywa brakujące okno i cytuje
 * odmowę, którą rdzeń odda przy próbie założenia karty powłoki bez wskazania.
 */
function zdanieBezOkna(nazwaModulu: string): HTMLElement {
  const element = document.createElement('p');
  element.className = 'dn-pole-opis dt-stan';
  element.dataset['stan'] = 'pusto';
  element.dataset['brak'] = 'okno-gospodarza';
  element.textContent =
    `Terminal stoi tu jako okno pomocnicze modułu ${nazwaModulu}, ale rdzeń nie dał temu ` +
    'modułowi okna wykonania. Karta powłoki bez okna nie ma ani trybu uprawnień, ani katalogu ' +
    'roboczego, więc rdzeń odmawia jej założenia wprost: „karta powłoki wymaga wskazania okna". ' +
    'To jest brak OKNA MODUŁU, a nie brak terminala i nie awaria — karty pojawią się tu ' +
    'razem z oknem sesji.';
  return element;
}
