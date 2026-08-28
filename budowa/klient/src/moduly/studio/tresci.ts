/**
 * Moduł Studio — katalog treści, jedyne miejsce z tekstem widocznym dla
 * Operatora w tym module. Ten sam wzorzec co w `rama/tresci.ts`.
 */

export type WezelTresci = string | number | { [klucz: string]: WezelTresci };

export const tresci = {
  okno: {
    etykieta: 'Okno robocze — Studio',
  },

  panel: {
    tytul: 'Urządzenia wejściowe',
    ladowanie: 'Odpytywanie rdzenia…',
  },

  pusto: {
    tytul: 'Brak urządzeń wejściowych',
    opis: 'Rdzeń nie zgłosił żadnego skanera ani kamery podłączonych do maszyny.',
  },

  odmowa: {
    glowa: 'Rdzeń odmówił wykazu urządzeń',
    brakOpisu: 'Rdzeń nie podał powodu odmowy.',
  },

  urzadzenie: {
    podajnik: 'z podajnikiem',
    dpi: 'dpi',
  },
} as const;
