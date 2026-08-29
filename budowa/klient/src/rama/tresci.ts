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
    /* Nazwy wprost z prototypu (`design/05-okna/moduly/studio.html`, pas
       `dn-narzedzia`): sześć rodzin czynności okna, po kolei od lewej. */
    zwinPanel: 'Zwiń szynę sesji',
    wstecz: 'Wstecz',
    naprzod: 'Do przodu',
    centrum: 'Centrum dowodzenia',
    odswiez: 'Odśwież widok',
    zrzut: 'Zrzut ekranu',
    schowek: 'Schowek',
    kolejka: 'Kolejka zadań',
    magistrala: 'Magistrala kontekstu',
    izolacja: 'Izolacja kontekstu',
    powiadomienia: 'Powiadomienia',
    skupienie: 'Tryb skupienia',
    pelnyEkran: 'Pełny ekran',
    wiecej: 'Więcej',
    dostosuj: 'Dostosuj wstążkę',
    szukaj: 'Szukaj',
    karty: 'Karty sesji',
    /* Czynność, której warstwa okien jeszcze nie prowadzi. Zasada zero blokad:
       niegotowość nazywa się zdaniem, przycisk nie znika. */
    zapowiedziane: 'Ta czynność wejdzie w kolejnym wydaniu.',
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
