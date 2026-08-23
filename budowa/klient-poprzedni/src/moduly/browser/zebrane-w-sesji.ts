import type { BrowserNote, BrowserSource } from '../../../../shared/contract';

/**
 * Zbiór zebrany w toku tej sesji przeglądania: źródła, notatki i fragmenty
 * wyodrębnione z migawki.
 *
 * Jedna odpowiedzialność: pamięć tego, co Operator zgromadził. Bez odczytu
 * z rdzenia i bez ani jednego elementu widoku.
 *
 * Zbiór jest odbiciem rdzenia, nie drugą prawdą. Wykaz przychodzi komendami
 * `browser.source.list` i `browser.note.list` i podmienia zawartość w całości
 * (`zastapZrodla`, `zastapNotatki`), a nie doszywa się do zastanej: rdzeń wie,
 * co w oknie jest, i to jego odpowiedź rozstrzyga. Dopisanie po
 * udanym `add` zostaje obok — nowa pozycja ma być widoczna od razu, bez
 * czekania na ponowny odczyt, ale to ten sam wiersz, który przyjdzie w wykazie.
 *
 * Przypięcie i klasyfikacja są własnością widoku, nie rdzenia. Kontrakt nie
 * niesie ani pola przypięcia, ani pola rodzaju notatki, więc jedno i drugie
 * zapamiętane tutaj nie udaje zapisu — Notes Panel mówi o tym w swoim opisie.
 *
 * Każda zmiana zbioru jest ogłaszana. Źródło dodane w Sources Panel ma się
 * pojawić w wyborze powiązania notatki w tej samej chwili; bez ogłoszenia
 * drugie okno zobaczyłoby je dopiero przy własnym odświeżeniu, a Operator
 * dostałby wybór bez pozycji, którą właśnie zapisał.
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
  /**
   * Dopisuje fragment do wyodrębnionych i mówi, czy się dopisał. Pusty i taki
   * sam jak zastany nie wchodzą — a okno ma o tym powiedzieć zamiast meldować
   * dopisanie, którego nie było.
   */
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

    // Przypięcia zostają nietknięte: są własnością widoku, nie rdzenia, więc
    // podmiana wykazu nie ma prawa ich zgubić. Przypięcie pozycji, której
    // rdzeń już nie oddaje, po prostu nic nie porządkuje.
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
