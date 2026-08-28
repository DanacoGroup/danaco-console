import type { BrowserNote, BrowserSource } from '../../../../shared/contract';

/**
 * Zbiór zebrany w toku sesji przeglądania: źródła, notatki i fragmenty
 * wyodrębnione z migawki. Jest pamięcią tego, co Operator zgromadził — bez
 * odczytu z rdzenia i bez ani jednego elementu widoku. Każda jego zmiana jest
 * ogłaszana oknom modułu.
 */
export interface ZebraneWSesji {
  /** Źródła okna, najnowsze na początku — w kolejności, w której daje je rdzeń. */
  zrodla(): readonly BrowserSource[];
  /** Notatki okna, najświeższe na początku; przypięte idą przed wszystkimi. */
  notatki(): readonly BrowserNote[];
  /** Fragmenty wyodrębnione z migawki paskiem zaznaczenia. */
  wyodrebnione(): readonly string[];
  dopiszZrodlo(zrodlo: BrowserSource): void;
  dopiszNotatke(notatka: BrowserNote): void;
  /** Podmienia wykaz źródeł tym, co oddał rdzeń — w całości, nie doszywając. */
  zastapZrodla(wykaz: readonly BrowserSource[]): void;
  /** Podmienia wykaz notatek tym, co oddał rdzeń — w całości, nie doszywając. */
  zastapNotatki(wykaz: readonly BrowserNote[]): void;
  /** Dopisuje fragment do wyodrębnionych i mówi, czy się dopisał. */
  dopiszWyodrebniony(fragment: string): boolean;
  /** Czy notatka jest przypięta na górze wykazu. */
  czyPrzypieta(idNotatki: string): boolean;
  /** Przestawia przypięcie notatki i zwraca stan po zmianie. */
  przelaczPrzypiecie(idNotatki: string): boolean;
  /** Klasyfikacja notatki nadana w tej karcie; pusty napis znaczy „bez klasyfikacji". */
  klasyfikacja(idNotatki: string): string;
  /** Nadaje albo zdejmuje klasyfikację notatki. */
  ustawKlasyfikacje(idNotatki: string, kod: string): void;
  /** Źródło o podanym identyfikatorze albo `null`. */
  zrodlo(idZrodla: string): BrowserSource | null;
}

export function utworzZebraneWSesji(oglos: () => void): ZebraneWSesji {
  const zrodla: BrowserSource[] = [];
  const notatki: BrowserNote[] = [];
  const wyodrebnione: string[] = [];
  const przypiete = new Set<string>();
  const klasyfikacje = new Map<string, string>();

  /** Przypięte na górze, reszta od najświeższej — porządek wykazu. */
  function uporzadkowane(): BrowserNote[] {
    const gora = notatki.filter((wpis) => przypiete.has(wpis.id));
    const reszta = notatki.filter((wpis) => !przypiete.has(wpis.id));
    return [...gora, ...reszta];
  }

  return {
    zrodla: () => zrodla,
    notatki: () => uporzadkowane(),
    wyodrebnione: () => wyodrebnione,

    dopiszZrodlo(zrodlo) {
      const pozycja = zrodla.findIndex((wpis) => wpis.id === zrodlo.id);
      if (pozycja === -1) zrodla.unshift(zrodlo);
      else zrodla[pozycja] = zrodlo;
      oglos();
    },

    dopiszNotatke(notatka) {
      const pozycja = notatki.findIndex((wpis) => wpis.id === notatka.id);
      if (pozycja === -1) notatki.unshift(notatka);
      else notatki[pozycja] = notatka;
      oglos();
    },

    zastapZrodla(wykaz) {
      zrodla.splice(0, zrodla.length, ...wykaz);
      oglos();
    },

    // Przypięcia zostają nietknięte: są własnością widoku, a nie rdzenia.
    zastapNotatki(wykaz) {
      notatki.splice(0, notatki.length, ...wykaz);
      oglos();
    },

    dopiszWyodrebniony(fragment) {
      const tresc = fragment.trim();
      if (tresc === '' || wyodrebnione.includes(tresc)) return false;
      wyodrebnione.push(tresc);
      oglos();
      return true;
    },

    czyPrzypieta: (idNotatki) => przypiete.has(idNotatki),

    przelaczPrzypiecie(idNotatki) {
      const bylo = przypiete.has(idNotatki);
      if (bylo) przypiete.delete(idNotatki);
      else przypiete.add(idNotatki);
      oglos();
      return !bylo;
    },

    klasyfikacja: (idNotatki) => klasyfikacje.get(idNotatki) ?? '',

    ustawKlasyfikacje(idNotatki, kod) {
      if (kod === '') klasyfikacje.delete(idNotatki);
      else klasyfikacje.set(idNotatki, kod);
      oglos();
    },

    zrodlo: (idZrodla) => zrodla.find((wpis) => wpis.id === idZrodla) ?? null,
  };
}
