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
    minimalizuj: 'Minimalizuj okno',
    maksymalizuj: 'Maksymalizuj okno',
    zamknij: 'Zamknij okno',
  },

  szyna: {
    etykieta: 'Nawigacja środowiska',
    nowaSesja: 'Nowa sesja',
  },

  narzedzia: {
    etykieta: 'Pasek narzędzi',
    zwinPanel: 'Zwiń panel',
    wstecz: 'Wstecz',
    naprzod: 'Do przodu',
    odswiez: 'Odśwież widok',
    szukaj: 'Szukaj',
    karty: 'Karty sesji',
  },

  panel: {
    etykieta: 'Panel sesji',
    sekcje: 'Sekcje panelu',
    sesje: 'Sesje',
    projekty: 'Projekty',
    szukaj: 'Szukaj w panelu',
    nowaSesja: 'Nowa sesja',
    filtry: 'Filtry',
    brakSesji: 'Środowisko nie niesie jeszcze żadnej karty sesji.',
    brakProjektow: 'Rdzeń nie wystawia dziś osobnego wykazu projektów.',
  },

  glowna: {
    etykieta: 'Obszar roboczy',
    brakModulu: 'Wybierz moduł z szyny nawigacji po lewej.',
  },

  stan: {
    etykieta: 'Pasek stanu',
    srodowisko: 'Środowisko',
    sesje: 'Sesje czynne',
    widok: 'Widok',
    motywJasny: 'jasny',
    motywCiemny: 'ciemny',
    operator: 'Operator',
    brakTozsamosci: 'tożsamość nie przekazana',
  },
} as const;
