import { defineConfig } from 'vitest/config';

// Uruchamiacz sprawdzianów klienta, osobny od budowania: środowisko domyślne to node, a widoki dotykające dokumentu dostają jsdom przez dopasowanie ścieżek; katalog nadrzędny jest dostępny dla importu wspólnego kontraktu.
export default defineConfig({
  test: {
    environment: 'node',
    environmentMatchGlobs: [
      ['src/{motyw,komponenty,mission-control,widok-sterowania,strona-glowna}/**/*.test.ts', 'jsdom'],
      // Układ okien równoległych buduje elementy wprost w dokumencie, więc jego sprawdziany wymagają jsdom.
      ['src/okna-rownolegle/**/*.test.ts', 'jsdom'],
      // Okno Punktów Izolacji buduje panele wprost w dokumencie, więc jego sprawdziany wymagają jsdom.
      ['src/punkty-izolacji/**/*.test.ts', 'jsdom'],
      // Scena wejścia osadza się w dokumencie i sama się z niego zdejmuje, więc wymaga jsdom.
      ['src/ladowanie/**/*.test.ts', 'jsdom'],
      // Okna modułów budują widoki wprost w dokumencie, więc ich sprawdziany wymagają jsdom.
      ['src/moduly/**/*.test.ts', 'jsdom'],
      // Powierzchnia mobilna buduje ekran przeglądu wprost w dokumencie, więc wymaga jsdom.
      ['src/mobile/**/*.test.ts', 'jsdom'],
    ],
    include: ['src/**/*.test.ts'],
    // Uruchomienie bez ani jednego pliku sprawdzianu kończy się błędem, a nie
    // cichym powodzeniem.
    passWithNoTests: false,
    // Podglądy konsoli zakładane w sprawdzianach znikają po każdym z nich.
    restoreMocks: true,
    // Vitest domyślnie zaślepia importy arkuszy; sprawdziany widoków czytają je przez import surowy.
    css: true,
  },
  server: {
    fs: { allow: ['..'] },
  },
});
