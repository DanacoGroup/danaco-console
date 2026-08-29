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
    /* Brzmienie z prototypu: szyna niesie i środowiska, i moduły. */
    etykieta: 'Nawigacja środowisk i modułów',
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
    brakSesji: 'Brak otwartych kart sesji.',
    brakProjektow: 'Wykaz projektów nie jest dostępny w tej wersji.',
  },

  glowna: {
    etykieta: 'Obszar roboczy',
    /* „Szyna nawigacji” to nasza nazwa strefy, nie Operatora — w oknie mówimy,
       gdzie ma spojrzeć, a nie jak ta część interfejsu nazywa się w budowie. */
    brakModulu: 'Wybierz moduł z paska po lewej stronie.',
  },

  stan: {
    etykieta: 'Pasek stanu',
    srodowisko: 'Środowisko',
    sesje: 'Sesje czynne',
    widok: 'Widok',
    motywJasny: 'jasny',
    motywCiemny: 'ciemny',
    operator: 'Operator',
    /* „Tożsamość” jest pojęciem z warstwy nakładania profilu na model —
       w pasku stanu Operator czyta wyłącznie, czy program zna jego konto. */
    brakTozsamosci: 'konto nierozpoznane',
  },
} as const;
