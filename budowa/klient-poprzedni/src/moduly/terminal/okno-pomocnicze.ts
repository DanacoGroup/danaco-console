import './terminal.css';

import { utworzRameOkna } from '../../komponenty/rama-okna';
import type { Kanal } from '../../protokol/kanal';
import { zamontujTerminal, type ZamontowanyTerminal } from './indeks';

/**
 * Terminal jako okno pomocnicze modułu gospodarza.
 *
 * Terminal nie jest samodzielnym modułem ani niezależną sesją — stanowi
 * dodatkowe okno pomocnicze dostępne wewnątrz modułów Developer, Diagnostics
 * i Apps (patrz `okno-komunikacji/profile-inzynieria.ts`).
 *
 * Rdzeń wiąże kartę powłoki z oknem, a nie z modułem: `OtworzKarte`
 * w `server/internal/core/adapter_modul_terminal.go` żąda wyłącznie `windowId`
 * okna otwartego (`oknoWykonania` w `adapter_modul_terminal_okno.go`) i nie
 * sprawdza, czy okno należy do modułu `terminal`. Karta bierze z okna tryb
 * uprawnień i katalog roboczy, więc terminal w oknie Developera pracuje
 * w kontekście tego właśnie modułu.
 *
 * Plik nie powiela modułu: nie ma tu drugiego stanu terminala, drugiego źródła
 * komend ani drugiej konsoli. Okno woła `zamontujTerminal`
 * z `moduly/terminal/indeks.ts` — to samo złożenie, które stoi w module —
 * i podaje mu okno gospodarza.
 *
 * Gospodarz bywa zamontowany bez okna wykonania (moduł Diagnostics ma trzy
 * z czterech komend bez `windowId`). Okno pokazuje wtedy zdanie o braku okna
 * zamiast pustej konsoli, którą czytałoby się jako „nic nie biegnie".
 *
 * Plik nie zdejmuje Terminala z bocznej nawigacji: pozycja stoi w macierzy
 * widoczności rdzenia (`server/internal/store/migracja_007_zaczyn_slownikow.sql`),
 * a klient bierze wykaz modułów wyłącznie z `environment.enter` i `module.list`.
 *
 * Użycie:
 *
 *   const terminal = utworzOknoTerminalaPomocnicze({
 *     kanal, modul: 'Developer', przedrostek: 'mdev', okno: idOkna,
 *   });
 *   zlozenie.append(terminal.element);
 *   terminal.ustawOkno(idOkna);   // gdy okno przychodzi później, po odczycie
 *   // przy zamykaniu modułu: terminal.zamknij();
 */
export interface OknoTerminalaPomocnicze {
  /** Sekcja osadzana w złożeniu modułu gospodarza albo w pasie okien pomocniczych. */
  element: HTMLElement;
  /**
   * Podaje okno wykonania gospodarza. Pusty napis znaczy „rdzeń nie dał temu
   * modułowi okna" i wypełnia ciało zdaniem o braku. Wolno wołać wielokrotnie:
   * to samo okno nie przebudowuje niczego, inne przebudowuje złożenie wraz
   * z zamknięciem poprzednich subskrypcji.
   */
  ustawOkno(okno: string): void;
  /** Odczytuje okna terminala z rdzenia; bez okna gospodarza nie robi nic. */
  odswiez(): void;
  /** Zamyka nasłuch zdarzeń terminala założony przez to okno. */
  zamknij(): void;
}

/** Zależności okna pomocniczego. */
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
    // `zamontujTerminal` sam podmienia zawartość gospodarza, więc czyszczenie
    // ciała byłoby drugą drogą do tego samego skutku.
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
 * Zdanie stanu „gospodarz bez okna": nazywa brakujące okno i cytuje odmowę,
 * którą rdzeń odda przy próbie założenia karty — `bladZadaniaTerminala("karta
 * powłoki wymaga wskazania okna")` w `OtworzKarte`
 * (`server/internal/core/adapter_modul_terminal.go`).
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
