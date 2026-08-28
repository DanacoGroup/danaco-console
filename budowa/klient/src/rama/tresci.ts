/**
 * Rama aplikacji — katalog treści, jedyne miejsce z tekstem widocznym dla
 * Operatora. Klucze idą hierarchicznie: strefa do elementu, jak ścieżka.
 */

/** Węzeł katalogu: łańcuch, liczba albo poddrzewo złożone z kolejnych węzłów tego samego rodzaju. */
export type WezelTresci = string | number | { [klucz: string]: WezelTresci };

export const tresci = {
  belka: {
    marka: 'Danaco Console',
    separator: '›',
  },

  szyna: {
    etykieta: 'Nawigacja środowiska',
  },

  glowna: {
    etykieta: 'Obszar roboczy',
  },

  stan: {
    etykieta: 'Pasek stanu',
    srodowisko: 'Środowisko',
    sesje: 'Sesje czynne',
    widok: 'Widok',
    motywJasny: 'jasny',
    motywCiemny: 'ciemny',
  },
} as const;
