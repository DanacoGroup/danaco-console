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
    /* Trzy formy każdego rzeczownika: „1 znak”, „3 znaki”, „7 znaków”.
       Formę wybiera `liczebnik.ts`, słowa stoją tutaj — w katalogu. */
    jednostkaZnaki: { jedna: 'znak', kilka: 'znaki', wiele: 'znaków' },
    jednostkaSlowa: { jedna: 'słowo', kilka: 'słowa', wiele: 'słów' },
    brakDokumentu: 'W tym oknie nie ma dokumentu.',
    brakZaznaczenia: 'Nie zaznaczono fragmentu dokumentu.',
  },

  operacje: {
    wlasna: 'własna',
    brakWykazu: 'Brak dostępnych operacji.',
  },

  uruchom: {
    przycisk: 'Uruchom operację',
    wBiegu: 'Uruchamianie…',
    brakWyniku: 'Model nie zwrócił odpowiedzi.',
  },

  odmowa: {
    brakOkna: 'Nie udało się otworzyć okna modułu.',
    lista: 'Nie udało się wczytać wykazu operacji',
    uruchomienie: 'Nie udało się uruchomić operacji',
  },
} as const;
