import type { TranslationPanel, Window } from '../../../../shared/contract';

/**
 * Jedno źródło prawdy modułu Translate — sam zbiór danych, bez wywołań.
 *
 * Zbiór jest oddzielony od odczytu: `stan-translate.ts` wie, jak zapytać rdzeń,
 * ten plik wie wyłącznie, co moduł już wie. Dzięki temu okna można sprawdzić na
 * samym zbiorze, a odczyt nie miesza się z pamięcią.
 *
 * Trzy pustki są rozróżnialne: faza `spoczynek` znaczy „jeszcze nie pytałem",
 * `odczyt` — „pytam", `gotowe` z pustym oknem — „rdzeń nie zna ani jednego okna
 * tej sesji". Zlanie ich w jedno kazałoby zgadywać, czy czekać, czy działać
 * (wzór: `dostepy/stany-odczytu.ts`).
 */
export type FazaKontekstu = 'spoczynek' | 'odczyt' | 'gotowe' | 'blad';

export interface MagazynTranslate {
  oknoModulu(): Window | null;
  ustawOkno(okno: Window | null): void;
  faza(): FazaKontekstu;
  powod(): string;
  ustawFaze(faza: FazaKontekstu, powod?: string): void;
  tekstZrodlowy(): string;
  jezykZrodlowy(): string;
  segmenty(): readonly string[];
  ustawZrodlo(tekst: string, jezyk: string, segmenty: readonly string[]): void;
  ustawSegmenty(segmenty: readonly string[]): void;
  panele(): readonly TranslationPanel[];
  wchlonPanel(panel: TranslationPanel): void;
  usunPanel(idPanelu: string): void;
  ustawPanele(panele: readonly TranslationPanel[]): void;
  obserwuj(sluchacz: () => void): () => void;
  /**
   * Ogłasza zmianę, której magazyn nie prowadzi u siebie.
   *
   * Cały moduł ma jeden nasłuch. Rejestr kanałów modelu leży poza magazynem
   * (`zrodlo-kanalow-translate.ts`) i jego danych nie dotyka, ale jego powrót ma
   * przerysować stery — a wszystkie widoki modułu przerysowuje jedno ogłoszenie
   * stanu. Drugi, równoległy nasłuch w oknach trzeba by odpinać w każdej
   * instancji panelu i pierwszy przeoczony byłby wyciekiem.
   */
  ogloszZmiane(): void;
  zapomnijSluchaczy(): void;
}

export function utworzMagazynTranslate(): MagazynTranslate {
  const sluchacze = new Set<() => void>();
  let okno: Window | null = null;
  let stanFazy: FazaKontekstu = 'spoczynek';
  let stanPowodu = '';
  let tekst = '';
  let jezyk = '';
  let podzial: readonly string[] = [];
  let wykaz: TranslationPanel[] = [];

  function oglos(): void {
    for (const sluchacz of [...sluchacze]) sluchacz();
  }

  return {
    oknoModulu: () => okno,
    faza: () => stanFazy,
    powod: () => stanPowodu,
    tekstZrodlowy: () => tekst,
    jezykZrodlowy: () => jezyk,
    segmenty: () => podzial,
    panele: () => wykaz,

    ustawOkno(nowe) {
      okno = nowe;
      oglos();
    },

    ustawFaze(faza, powod = '') {
      stanFazy = faza;
      stanPowodu = powod;
      oglos();
    },

    ustawZrodlo(nowyTekst, nowyJezyk, noweSegmenty) {
      tekst = nowyTekst;
      jezyk = nowyJezyk;
      podzial = [...noweSegmenty];
      oglos();
    },

    ustawSegmenty(noweSegmenty) {
      podzial = [...noweSegmenty];
      oglos();
    },

    wchlonPanel(panel) {
      const pozycja = wykaz.findIndex((wpis) => wpis.id === panel.id);
      if (pozycja === -1) wykaz = [...wykaz, panel];
      else wykaz = wykaz.map((wpis) => (wpis.id === panel.id ? panel : wpis));
      oglos();
    },

    usunPanel(idPanelu) {
      wykaz = wykaz.filter((wpis) => wpis.id !== idPanelu);
      oglos();
    },

    ustawPanele(nowe) {
      wykaz = [...nowe];
      oglos();
    },

    obserwuj(sluchacz) {
      sluchacze.add(sluchacz);
      return () => sluchacze.delete(sluchacz);
    },

    ogloszZmiane: oglos,

    zapomnijSluchaczy: () => sluchacze.clear(),
  };
}
