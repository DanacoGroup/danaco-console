/*
 * Nastawa budowania pakietu klienta.
 *
 * Vite wchodzi jako narzędzie budowania, nie jako zależność pakietu: klient
 * nie ma i nie zyskuje tu ani jednej zależności produkcyjnej, a wytworem jest
 * dokument wraz z jednym modułem i jednym arkuszem, bez śladu narzędzia.
 * Wybór wynika z trzech rzeczy, których wytworzenie pakietu wymaga: rozwinięcia
 * `@import` w arkuszach warstwy projektowej wraz z przeniesieniem plików
 * krojów, przełożenia TypeScriptu ze specyfikatorem `.ts` na moduł
 * przeglądarki oraz wpięcia obu w dokument.
 *
 * Nastawa nie wnosi `import` z pakietu `vite` — narzędzie stoi na maszynie
 * globalnie, więc plik opisujący budowanie nie może zależeć od tego, czy
 * pakiet daje się rozwiązać z katalogu klienta.
 */

/** Korzeń projektu — katalog klienta, nie katalog roboczy powłoki. */
const korzen = new URL('.', import.meta.url).pathname;

/**
 * Korzeń repozytorium. Arkusze warstwy projektowej i wygenerowany kontrakt
 * leżą poza katalogiem klienta, więc serwer poglądowy musi mieć zgodę na ich
 * odczyt; budowanie czyta je bez tej zgody.
 */
const korzenRepozytorium = new URL('../../', import.meta.url).pathname;

export default {
  root: korzen,

  /* Ścieżki względne, bo pakiet jest oddawany dwiema drogami: rdzeń serwuje go
     spod korzenia nasłuchu, a powłoka Tauri wczytuje z własnego protokołu.
     Ścieżka bezwzględna wiązałaby pakiet z jednym z tych dwóch miejsc. */
  base: './',

  build: {
    outDir: 'dist',
    emptyOutDir: true,
    /* Kroje idą plikami, nie treścią wplecioną w arkusz: wpisany krój rośnie
       w każdym pakiecie o pełny rozmiar podzbioru, a przeglądarka nie może go
       pominąć na podstawie zakresu znaków. */
    assetsInlineLimit: 0,
  },

  server: {
    fs: { allow: [korzenRepozytorium] },
  },
};
