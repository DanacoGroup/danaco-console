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
    brakDokumentu: 'W tym oknie nie ma jeszcze dokumentu.',
    brakZaznaczenia:
      'Rozmiar zaznaczenia jeszcze się tu nie pokazuje — zaznaczenie z edytora nie dochodzi do tej karty. Operacja na zaznaczeniu działa mimo to.',
  },

  operacje: {
    wlasna: 'własna',
    /* Znaki wiersza wzięte z prototypu: „▸” przy operacji gotowej do wyboru,
       „✓” przy wybranej — ten sam zapis, co w wykazie zadań karty Plan. */
    znakWiersza: '▸',
    znakWyboru: '✓',
    brakWykazu: 'Wykaz operacji jest pusty.',
    warianty:
      'Wybór wariantu operacji — na przykład języka tłumaczenia — jeszcze nie działa: operacje przychodzą bez opisu swoich ustawień.',
  },

  uruchom: {
    przycisk: 'Uruchom operację',
    wBiegu: 'Uruchamianie…',
    brakWyniku: 'Model nie zwrócił odpowiedzi.',
    propozycja: 'Wynik odłożono jako propozycję zmiany — obejrzysz ją w karcie Diff/Grep Panel.',
  },

  usun: {
    przycisk: 'Usuń operację własną',
    wBiegu: 'Usuwanie…',
    fabryczna: 'Operacji fabrycznej nie da się usunąć — zostaje w wykazie.',
  },

  zapis: {
    naglowek: 'Operacja własna',
    etykietaNazwy: 'Nazwa operacji',
    zastepczaNazwa: 'Nazwa widoczna w wykazie',
    etykietaKategorii: 'Kategoria operacji',
    zastepczaKategoria: 'Kategoria, w której operacja stanie',
    etykietaTresci: 'Treść polecenia',
    zastepczaTresc: 'Polecenie, które model ma wykonać',
    przycisk: 'Zapisz operację',
    wBiegu: 'Zapisywanie…',
    zapisana: 'Operacja zapisana — stoi w wykazie powyżej.',
  },

  szukanie: {
    naglowek: 'Wyszukiwanie znaczeniowe',
    etykieta: 'Wyszukiwanie znaczeniowe w dokumencie',
    zastepczaTresc: 'Opisz, czego szukasz w dokumencie…',
    przycisk: 'Szukaj',
    wBiegu: 'Szukanie…',
    brakDopasowan: 'Żaden fragment dokumentu nie jest bliski temu zapytaniu.',
    bliskosc: 'bliskość',
    wiersz: 'wiersz',
  },

  odmowa: {
    brakOkna: 'Nie udało się otworzyć okna modułu.',
    lista: 'Nie udało się wczytać wykazu operacji',
    uruchomienie: 'Nie udało się uruchomić operacji',
    zapis: 'Nie udało się zapisać operacji',
    usuniecie: 'Nie udało się usunąć operacji',
    szukanie: 'Nie udało się przeszukać dokumentu',
  },
} as const;
