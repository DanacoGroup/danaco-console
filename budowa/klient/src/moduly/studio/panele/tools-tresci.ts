/**
 * Panel Tools Panel — katalog treści, jedyne miejsce tego panelu z tekstem
 * widocznym dla Operatora. Osobny plik od `moduly/studio/tresci.ts`, żeby
 * panel nie dotykał pliku dzielonego z pozostałymi sześcioma panelami terenu.
 */

export const tresciNarzedzi = {
  panel: {
    tytul: 'Tools Panel',
  },

  ladowanie: 'Wczytywanie operacji Tools Panel…',

  zakres: {
    etykietaPrzelacznika: 'Zakres działania',
    zaznaczenie: 'Zaznaczenie',
    dokument: 'Cały dokument',
    etykietaRozmiaru: 'Zakres:',
    jednostkaZnaki: 'znaków',
    jednostkaSlowa: 'słów',
    brakDokumentu: 'Rdzeń nie zgłosił dokumentu w tym oknie.',
    brakZaznaczenia: 'Rdzeń nie podał zakresu zaznaczenia.',
  },

  operacje: {
    wlasna: 'własna',
    brakWykazu: 'Rdzeń nie zgłosił żadnej operacji Tools Panel.',
  },

  uruchom: {
    przycisk: 'Uruchom operację',
    wBiegu: 'Uruchamianie…',
    brakWyniku: 'Rdzeń nie zwrócił treści wyniku.',
  },

  odmowa: {
    brakOkna: 'Rdzeń nie założył okna modułu.',
    lista: 'Rdzeń odmówił wykazu operacji',
    uruchomienie: 'Rdzeń odmówił uruchomienia operacji',
    brakOpisu: 'Rdzeń nie podał powodu odmowy.',
  },
} as const;
