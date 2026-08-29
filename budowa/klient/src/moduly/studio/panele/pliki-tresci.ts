/**
 * Katalog treści panelu Pliki — jedyne miejsce w `panele/pliki.ts`, w którym
 * wolno zapisać zdanie widoczne dla Operatora. Ten sam układ co
 * `moduly/studio/tresci.ts`; osobny plik, żeby panel nie dotykał katalogu
 * współdzielonego przez sześć innych paneli budowanych równolegle.
 */

export const tresciPliki = {
  urzadzenia: {
    tytul: 'Urządzenia wejściowe',
    ladowanie: 'Sprawdzanie urządzeń wejściowych…',
    brak: 'Rdzeń nie zgłosił żadnego urządzenia wejściowego.',
    podajnik: 'z podajnikiem',
    skanuj: 'Skanuj',
    skanowanie: 'Skanowanie…',
  },

  kolejka: {
    tytul: 'Kolejka wczytywania',
    ladowanie: 'Wczytywanie kolejki…',
    brak: 'Rdzeń nie zgłosił żadnej pozycji w kolejce wczytywania.',
    stron: 'stron',
    pewnosc: 'pewność rozpoznania',
    zPisma: 'tekst z rozpoznania pisma',
  },

  stan: {
    oczekuje: 'Oczekuje',
    przetwarzanie: 'Rozpoznawanie w toku',
    gotowa: 'Gotowa',
    ponowienie: 'Do ponowienia',
    odmowa: 'Odmowa',
  },

  dodaj: {
    etykietaSciezki: 'Ścieżka materiału widziana przez rdzeń',
    doloz: 'Dołóż plik',
    dokladanie: 'Dokładanie…',
  },

  adres: {
    etykieta: 'Adres strony',
    dolacz: 'Dołącz stronę',
    dolaczanie: 'Pobieranie strony…',
  },

  dzialania: {
    rozpoznaj: 'Rozpoznaj',
    rozpoznawanie: 'Rozpoznawanie…',
    przyjmij: 'Przyjmij do edytora',
    przyjmowanie: 'Przyjmowanie…',
  },

  tytulDokumentu: {
    etykieta: 'Tytuł zakładanego dokumentu',
  },

  poprawa: {
    tytul: 'Słowa rozpoznane — poprawa przed przyjęciem',
    etykietaSlowa: 'Słowo',
    zapisz: 'Zapisz',
    zapisywanie: 'Zapisywanie…',
    zapisano: 'Poprawka zapisana',
  },

  odmowa: {
    urzadzenia: 'Rdzeń odmówił wykazu urządzeń',
    kolejka: 'Rdzeń odmówił wykazu kolejki',
    dodanie: 'Rdzeń odmówił dołożenia materiału',
    rozpoznanie: 'Rdzeń odmówił rozpoznania',
    poprawka: 'Rdzeń odmówił zapisu poprawki',
    przyjecie: 'Rdzeń odmówił przyjęcia do edytora',
    adres: 'Rdzeń odmówił pobrania strony',
    skan: 'Rdzeń odmówił skanowania',
  },
} as const;
