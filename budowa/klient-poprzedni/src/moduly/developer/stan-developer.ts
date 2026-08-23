import type { DeveloperFile } from '../../../../shared/contract';
import type { Odsubskrybuj } from '../../polaczenie/magistrala-zdarzen';
import type { ZrodloDeveloper } from './zrodlo-developer';

/**
 * Stan wspólny modułu Developer — jedna prawda dla czterech okien.
 *
 * Cztery okna pracują nad tym samym repozytorium: Project Tree wskazuje plik,
 * Code Editor go otwiera i zapisuje, Git Panel podstawia ścieżkę wskazaną do
 * pola czynności i sam wskazuje ścieżki ze swojego wyniku, a Build Output
 * wskazuje pliki ze zgłoszeń przebiegu. Gdyby każde okno trzymało własną
 * ścieżkę i własny korzeń, wskazanie pliku w drzewie nie dotarłoby do edytora.
 *
 * `windowId` stoi tutaj, bo kontrakt wymaga go we wszystkich pięciu komendach
 * obszaru, a okna mają go podać identycznie — inaczej rdzeń rozdzieliłby ich
 * pracę między dwa katalogi robocze.
 *
 * Stan nie wywołuje komend za okna: trzyma wybór i rozgłasza zmianę, a odczyt
 * i zapis należą do okien, bo tylko one mają stan ładowania i odmowy.
 */
export interface StanDevelopera {
  /** Okno modułu wymagane w każdej z pięciu komend obszaru. */
  okno(): string;
  /** Ścieżka wskazana; pusta znaczy „nie wskazano”. */
  sciezka(): string;
  /**
   * Ostatni plik oddany przez rdzeń — wraz z jego własną ścieżką.
   *
   * To nie to samo co `sciezka()`: wskazanie biegnie natychmiast po kliknięciu
   * w drzewie, a plik przychodzi dopiero odpowiedzią rdzenia. Zgodność obu
   * sprawdza się porównaniem `plik()?.path` ze `sciezka()` — i po to plik
   * niesie własną ścieżkę.
   *
   * Czytają to Code Editor (rozstrzyga, czy w polu leży już wskazany plik,
   * i bierze stąd `versionId` sprzed zapisu) oraz Project Tree (odróżnia węzeł
   * wskazany od węzła wczytanego do edytora).
   */
  plik(): DeveloperFile | null;
  /**
   * Treść pliku oddanego przez rdzeń — pusta, dopóki rdzeń jej nie podał.
   *
   * To nie jest treść pola edycji. Pole bywa zmienione i niezapisane, a rdzeń
   * po zapisie potrafi nie odesłać treści wcale; znacznik czystości pola należy
   * więc do Code Editora i stoi tam, nie tutaj.
   */
  tresc(): string;
  /** Wskazanie pliku bez jego treści — czynność Project Tree. */
  wskazPlik(sciezka: string): void;
  /** Plik po odczycie albo zapisie — czynność Code Editora. */
  ustawPlik(plik: DeveloperFile): void;
  /** Katalog korzenia drzewa; pusty znaczy „katalog roboczy okna”. */
  korzen(): string;
  /** Przestawia drzewo na inny korzeń i powiadamia okna. */
  ustawKorzen(sciezka: string): void;
  /** Subskrypcja zmiany wskazania, treści albo korzenia. */
  naZmiane(sluchacz: () => void): Odsubskrybuj;
  /** Zamyka nasłuch zdarzeń rdzenia. */
  zamknij(): void;
}

/** Zależności stanu: okno modułu oraz wskazania otwierane od razu. */
export interface OpcjeStanuDevelopera {
  okno: string;
  sciezka?: string;
  korzen?: string;
}

export function utworzStanDevelopera(
  zrodlo: ZrodloDeveloper,
  opcje: OpcjeStanuDevelopera,
): StanDevelopera {
  const sluchacze = new Set<() => void>();
  const idOkna = opcje.okno;
  let sciezkaBiezaca = opcje.sciezka ?? '';
  let korzenBiezacy = opcje.korzen ?? '';
  let plikBiezacy: DeveloperFile | null = null;

  function powiadom(): void {
    for (const sluchacz of [...sluchacze]) sluchacz();
  }

  // Przyrost budowania przychodzi także z pracy innego okna tego konta. Liczy
  // się wyłącznie przyrost dotyczący tego okna modułu: budowanie generuje pliki,
  // więc drzewo i plik w edytorze mogą być nieaktualne i okna mają się odczytać
  // ponownie.
  const odsubskrybujBudowanie = zrodlo.naZmianeBudowania((tresc) => {
    if (tresc.build.windowId !== idOkna) return;
    powiadom();
  });

  return {
    okno: () => idOkna,

    sciezka: () => sciezkaBiezaca,

    plik: () => plikBiezacy,

    tresc: () => plikBiezacy?.content ?? '',

    wskazPlik(sciezka) {
      const przyciety = sciezka.trim();
      if (przyciety === sciezkaBiezaca) return;
      sciezkaBiezaca = przyciety;
      // Plik zostaje, zmienia się samo wskazanie. Plik niesie własną `path`,
      // więc pomylić go ze wskazaniem nie sposób, a wyzerowanie kasowałoby
      // jedyną wiedzę o tym, co leży w polu edytora. Czytelnik pokazujący treść
      // porównuje `plik()?.path` ze `sciezka()`.
      powiadom();
    },

    ustawPlik(plik) {
      plikBiezacy = plik;
      sciezkaBiezaca = plik.path;
      powiadom();
    },

    korzen: () => korzenBiezacy,

    ustawKorzen(sciezka) {
      const przyciety = sciezka.trim();
      if (przyciety === korzenBiezacy) return;
      korzenBiezacy = przyciety;
      powiadom();
    },

    naZmiane(sluchacz) {
      sluchacze.add(sluchacz);
      return () => sluchacze.delete(sluchacz);
    },

    zamknij() {
      sluchacze.clear();
      odsubskrybujBudowanie();
    },
  };
}
