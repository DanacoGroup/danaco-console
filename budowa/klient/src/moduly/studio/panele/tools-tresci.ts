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
    brakDokumentu: 'W tym oknie nie ma dokumentu.',
    brakZaznaczenia: 'Zakres zaznaczenia nie jest znany.',
  },

  operacje: {
    wlasna: 'własna',
    brakWykazu: 'Nie ma żadnej operacji do uruchomienia.',
  },

  uruchom: {
    przycisk: 'Uruchom operację',
    wBiegu: 'Uruchamianie…',
    brakWyniku: 'Wynik nie zawiera treści.',
  },

  odmowa: {
    brakOkna: 'Nie udało się otworzyć okna modułu.',
    lista: 'Nie udało się wczytać wykazu operacji',
    uruchomienie: 'Nie udało się uruchomić operacji',
    brakOpisu: 'Powód nie został podany.',
  },
} as const;
