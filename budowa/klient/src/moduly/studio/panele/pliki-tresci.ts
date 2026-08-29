/**
 * Katalog treści panelu Pliki — jedyne miejsce w `panele/pliki.ts`, w którym
 * wolno zapisać zdanie widoczne dla Operatora. Ten sam układ co
 * `moduly/studio/tresci.ts`; osobny plik, żeby panel nie dotykał katalogu
 * współdzielonego przez sześć innych paneli budowanych równolegle.
 */

export const tresciPliki = {
  panel: {
    tytul: 'Pliki',
    brakOkna: 'Okno modułu jeszcze nie stanęło. Kolejka wczytywania pojawi się, gdy będzie gdzie ją osadzić.',
  },

  filtr: {
    etykieta: 'Filtruj',
    zacheta: 'Filtruj pliki…',
    brakTrafien: 'Żadna pozycja kolejki nie pasuje do filtra.',
  },

  wnoszenie: {
    tytul: 'Wnoszenie materiału',
    /* Wskazanie pliku okienkiem systemowym nie ma pokrycia w kontrakcie:
       komendy biorą położenie widziane przez rdzeń, a przeglądarka takiego
       położenia nie zna. Zdanie stoi w widoku zamiast przycisku-atrapy. */
    bezOkienka: 'Wskazanie pliku okienkiem systemowym jeszcze nie działa. Wpisz położenie pliku takie, jakie widzi je aplikacja.',
  },

  urzadzenia: {
    tytul: 'Urządzenia wejściowe',
    ladowanie: 'Sprawdzanie urządzeń wejściowych…',
    brak: 'Brak urządzeń wejściowych.',
    podajnik: 'z podajnikiem',
    skanuj: 'Skanuj',
    skanowanie: 'Skanowanie…',
    rodzaj: {
      skaner: 'Skaner',
      kamera: 'Kamera',
    },
  },

  kolejka: {
    tytul: 'Kolejka wczytywania',
    ladowanie: 'Wczytywanie kolejki…',
    brak: 'Kolejka wczytywania jest pusta.',
    /* Trzy formy: „1 strona”, „3 strony”, „7 stron”; wybór w `liczebnik.ts`. */
    stron: { jedna: 'strona', kilka: 'strony', wiele: 'stron' },
    pewnosc: 'pewność rozpoznania',
    zPisma: 'tekst z rozpoznania pisma',
    /* Usunięcia pozycji z kolejki nie ma dziś czym wykonać, więc zdanie nie
       może o nie prosić — kieruje na jedyną drogę, która działa naprawdę. */
    odrzucona: 'Tej pozycji nie udało się wczytać. Dodaj plik ponownie; usuwanie pozycji z kolejki jeszcze nie działa.',
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

  wprost: {
    etykieta: 'Położenie pliku wnoszonego do edytora',
    plik: 'Wnieś do edytora',
    pdf: 'Wnieś PDF',
    wnoszenie: 'Wnoszenie…',
    bilans: 'Wniesiono do edytora',
    stronyZTekstem: 'stron z tekstem',
    akapity: 'akapitów',
    tabele: 'tabel',
    tabeleNierozpoznane: 'układów tabelarycznych nierozpoznanych',
    obrazy: 'obrazów',
    obrazyPominiete: 'obrazów pominiętych',
    bezWarstwyTekstu: 'Plik nie niesie warstwy tekstowej. Pozycja stanęła w kolejce — rozpoznaj ją, zanim przyjmiesz wynik.',
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
    wniesienie: 'Nie udało się wnieść pliku do edytora',
  },
} as const;
