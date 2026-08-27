import type { DeveloperFile } from '../../../../shared/contract';
import type { Odsubskrybuj } from '../../polaczenie/magistrala-zdarzen';
import type { ZrodloDeveloper } from './zrodlo-developer';

/**
 * Stan wspólny modułu Developer jest jedyną prawdą dla czterech okien pracujących
 * nad tym samym repozytorium. Przechowuje okno modułu, ścieżkę wskazaną, ostatni
 * plik oddany przez rdzeń oraz korzeń drzewa i rozgłasza każdą ich zmianę.
 */
export interface StanDevelopera {
  /** Okno modułu wymagane w każdej z pięciu komend obszaru. */
  okno(): string;
  /** Ścieżka wskazana; pusta znaczy „nie wskazano”. */
  sciezka(): string;
  /** Ostatni plik oddany przez rdzeń wraz z własną ścieżką; nie jest wskazaniem. */
  plik(): DeveloperFile | null;
  /** Treść pliku oddanego przez rdzeń; pusta do odpowiedzi, nie treść pola edycji. */
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

/**
 * Zależności stanu: okno modułu wymagane w każdej komendzie obszaru oraz
 * wskazania ścieżki i korzenia, którymi stan otwiera pracę zaraz po utworzeniu.
 */
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

  // Liczy się wyłącznie przyrost budowania dotyczący tego okna modułu.
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
      // Plik zostaje, zmienia się samo wskazanie; plik niesie własną ścieżkę.
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
