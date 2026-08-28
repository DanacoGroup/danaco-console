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

  dokument: {
    tytul: 'Dokument',
    nazwaNowego: 'Dokument bez tytułu',
    zakladanie: 'Zakładanie okna roboczego w rdzeniu…',
    etykietaTresci: 'Treść dokumentu',
    zapisz: 'Zapisz',
    zapisywanie: 'Zapisywanie…',
    zapisany: 'Zapisane w repozytorium sesji',
    wersja: 'wersja',
    bezWersji: 'bez wersji',
  },

  dokumentOdmowa: {
    sesja: 'Rdzeń nie założył sesji',
    kanal: 'Rdzeń nie zgłosił żadnego kanału modelu',
    okno: 'Rdzeń nie założył okna modułu',
    dokument: 'Rdzeń nie założył dokumentu',
    zapis: 'Rdzeń odmówił zapisu dokumentu',
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
