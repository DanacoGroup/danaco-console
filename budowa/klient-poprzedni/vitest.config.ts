import { defineConfig } from 'vitest/config';

// Uruchamiacz sprawdzianów klienta.
//
// Osobny plik, a nie sekcja `test` dopisana do `vite.config.ts`: budowanie
// pakietu i sprawdzanie kodu to dwa różne byty. Vitest czyta ten plik
// z pierwszeństwem przed `vite.config.ts`, więc konfiguracja budowania zostaje
// nietknięta.
//
// Środowisko domyślne `node`, ponieważ warstwy protokołu i łączności nie
// dotykają dokumentu. Widoki dotykają go wprost (tworzą elementy, czytają
// arkusze, przełączają `data-theme`), więc ich sprawdziany dostają `jsdom`
// przez `environmentMatchGlobs`.
//
// `server.fs.allow` obejmuje katalog nadrzędny, bo warstwa protokołu importuje
// kontrakt z `budowa/shared/` — jedynego źródła prawdy nazw. Ta sama zgoda co
// w `vite.config.ts`.
export default defineConfig({
  test: {
    environment: 'node',
    environmentMatchGlobs: [
      ['src/{motyw,komponenty,mission-control,widok-sterowania,strona-glowna}/**/*.test.ts', 'jsdom'],
      // Układ okien równoległych buduje przełącznik, gniazda i pas relacji
      // wprost w dokumencie — bez `jsdom` jego sprawdziany nie miałyby czego
      // zbudować.
      ['src/okna-rownolegle/**/*.test.ts', 'jsdom'],
      // Okno Punktów Izolacji buduje trzy panele, macierz przełączników
      // i selektor zasięgu wprost w dokumencie — jego sprawdziany pytają
      // o zbudowane węzły, więc bez `jsdom` nie miałyby czego sprawdzić.
      ['src/punkty-izolacji/**/*.test.ts', 'jsdom'],
      // Scena wejścia stawia przesłonę z bryłą wprost w dokumencie i sama się
      // z niego zdejmuje — bez `jsdom` nie miałaby ani gdzie stanąć, ani skąd
      // zejść.
      ['src/ladowanie/**/*.test.ts', 'jsdom'],
      // Okna modułów budują swoje widoki wprost w dokumencie — ramy, paski
      // kontekstu, wykazy i panele wysuwane. Bez tego wpisu sprawdzian
      // któregokolwiek modułu nie miałby czego zbudować.
      ['src/moduly/**/*.test.ts', 'jsdom'],
      // Powierzchnia mobilna buduje ekran przeglądu zadań i procesów wprost
      // w dokumencie — kafle procesów wraz z przyciskami sterowania. Bez tego
      // wpisu sprawdzian nie miałby czego dotknąć.
      ['src/mobile/**/*.test.ts', 'jsdom'],
    ],
    include: ['src/**/*.test.ts'],
    // Uruchomienie bez ani jednego pliku sprawdzianu kończy się błędem, a nie
    // cichym powodzeniem.
    passWithNoTests: false,
    // Podglądy konsoli zakładane w sprawdzianach znikają po każdym z nich.
    restoreMocks: true,
    // Domyślnie Vitest zaślepia importy `.css` pustą treścią. Sprawdziany
    // widoków czytają arkusz przez import `?raw` i sprawdzają przełączenie
    // żetonu przez `getComputedStyle` po wstrzyknięciu reguł — bez tego
    // dostałyby pustkę zamiast arkusza produktu.
    css: true,
  },
  server: {
    fs: { allow: ['..'] },
  },
});
