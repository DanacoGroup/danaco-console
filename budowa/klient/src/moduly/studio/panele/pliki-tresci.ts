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
    brak: 'Brak urządzeń wejściowych.',
    podajnik: 'z podajnikiem',
    skanuj: 'Skanuj',
    skanowanie: 'Skanowanie…',
  },

  kolejka: {
    tytul: 'Kolejka wczytywania',
    ladowanie: 'Wczytywanie kolejki…',
    brak: 'Kolejka wczytywania jest pusta.',
    /* Trzy formy: „1 strona”, „3 strony”, „7 stron”; wybór w `liczebnik.ts`. */
    stron: { jedna: 'strona', kilka: 'strony', wiele: 'stron' },
    pewnosc: 'pewność rozpoznania',
    zPisma: 'tekst z rozpoznania pisma',
    odrzucona: 'Tej pozycji nie udało się wczytać. Usuń ją z kolejki i spróbuj dodać plik ponownie.',
  },

  stan: {
    oczekuje: 'Oczekuje',
    przetwarzanie: 'Rozpoznawanie w toku',
    gotowa: 'Gotowa',
    ponowienie: 'Do ponowienia',
    odmowa: 'Odmowa',
  },

  dodaj: {
    etykietaSciezki: 'Położenie pliku',
    doloz: 'Dodaj plik',
    dokladanie: 'Dodawanie…',
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
    tytul: 'Rozpoznane słowa do poprawienia',
    etykietaSlowa: 'Słowo',
    zapisz: 'Zapisz',
    zapisywanie: 'Zapisywanie…',
    zapisano: 'Poprawka zapisana',
  },

  odmowa: {
    urzadzenia: 'Nie udało się wczytać wykazu urządzeń',
    kolejka: 'Nie udało się wczytać kolejki',
    dodanie: 'Nie udało się dodać pliku',
    rozpoznanie: 'Nie udało się rozpoznać pisma',
    poprawka: 'Nie udało się zapisać poprawki',
    przyjecie: 'Nie udało się przyjąć do edytora',
    adres: 'Nie udało się pobrać strony',
    skan: 'Nie udało się uruchomić skanowania',
  },
} as const;
