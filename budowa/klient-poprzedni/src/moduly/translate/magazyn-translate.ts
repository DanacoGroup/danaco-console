import type { TranslationPanel, Window } from '../../../../shared/contract';

/**
 * Jedno źródło prawdy modułu Translate — zbiór danych bez wywołań, oddzielony od odczytu
 * prowadzonego przez stan-translate.ts, z trzema rozróżnialnymi fazami wczytywania danych z rdzenia.
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
  /** Ogłasza zmianę stanu wszystkim widokom modułu jednym wspólnym nasłuchem, bez duplikatu w oknach. */
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
