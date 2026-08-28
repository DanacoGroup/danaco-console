/*
 * Nastawa budowania pakietu klienta.
 *
 * Vite wchodzi jako narzędzie budowania, nie jako zależność pakietu: klient
 * nie ma i nie zyskuje tu ani jednej zależności produkcyjnej. Wytworem jest
 * dokument wraz z modułem TypeScriptu, arkuszami warstwy projektowej i
 * skryptami biblioteki ramy, bez śladu narzędzia.
 *
 * Rama wciąga bibliotekę `design/zasoby/` dwiema różnymi drogami, bo arkusze
 * i skrypty tej biblioteki mają różny kształt:
 *
 *   — Arkusze łączy `index.html` przez zwykły `<link rel="stylesheet">"
 *     wskazujący źródło w `design/zasoby/`. To zwykły arkusz CSS, więc Vite
 *     rozwija jego `@import`, przenosi kroje i obrazy wskazane przez `url()`
 *     i oddaje gotowy plik do `dist/assets/` — dokładnie tak samo, jak dla
 *     arkusza leżącego w katalogu klienta. Żadnej nastawy to nie wymaga.
 *
 *   — Skryptów Vite tą samą drogą przenieść nie może: to funkcje domknięte,
 *     nie moduły ECMAScript (żadnego `export`), więc znacznik `<script src>`
 *     bez `type="module"` Vite zostawia dosłownie, bez kopiowania pliku do
 *     `dist/` (ostrzeżenie budowania to potwierdza). Ścieżka źródłowa,
 *     zostawiona bez zmiany, po wydaniu wskazywałaby poza katalog serwowany
 *     przez rdzeń. Rozstrzygnięcie: wtyczka `closeBundle` niżej kopiuje
 *     dokładnie te skrypty biblioteki, których używa `design/05-okna/przeplyw/
 *     centrum-dowodzenia.html`, do `dist/zasoby/`, zachowując ich ścieżkę
 *     względną — `index.html` odwołuje się do nich przez `./zasoby/…`,
 *     ścieżkę już poprawną w wydanym pakiecie.
 */

import { copyFileSync, mkdirSync } from 'node:fs';
import { dirname, join } from 'node:path';
import { fileURLToPath } from 'node:url';

/** Korzeń projektu — katalog klienta, nie katalog roboczy powłoki. */
const korzen = new URL('.', import.meta.url).pathname;

/**
 * Korzeń repozytorium. Arkusze warstwy projektowej i wygenerowany kontrakt
 * leżą poza katalogiem klienta, więc serwer poglądowy musi mieć zgodę na ich
 * odczyt; budowanie czyta je bez tej zgody.
 */
const korzenRepozytorium = new URL('../../', import.meta.url).pathname;

/** Katalog biblioteki warstwy projektowej — jedyne źródło skryptów kopiowanych do pakietu. */
const KORZEN_ZASOBOW = new URL('../../design/zasoby/', import.meta.url);

/**
 * Skrypty biblioteki, w kolejności wczytania przez
 * `design/05-okna/przeplyw/centrum-dowodzenia.html` — osiemnaście plików,
 * ścieżka każdego względna wobec `KORZEN_ZASOBOW`.
 */
const SKRYPTY_BIBLIOTEKI = [
  'wspolne.js',
  'prototyp.js',
  'menu.js',
  'okna-modalne.js',
  'rama.js',
  'stanowisko.js',
  'stany.js',
  'pasek-stanu.js',
  'karty-okna.js',
  'panel-sesji.js',
  'pasek-okna.js',
  'okno-robocze.js',
  'okna/przelacznik-srodowisk.js',
  'okna/zakladki-paneli.js',
  'okna/centrum-dowodzenia.js',
  'okna/centrum-obszar.js',
  'okna/centrum-wejscie.js',
  'okna/centrum-dymki.js',
];

/** Kopiuje skrypty biblioteki do `dist/zasoby/` po zamknięciu paczki, zachowując ich ścieżkę względną. */
function kopiaSkryptowBiblioteki() {
  return {
    name: 'kopia-skryptow-biblioteki',
    apply: 'build',
    closeBundle() {
      for (const wzgledna of SKRYPTY_BIBLIOTEKI) {
        const zrodlo = fileURLToPath(new URL(wzgledna, KORZEN_ZASOBOW));
        const cel = join(korzen, 'dist', 'zasoby', wzgledna);
        mkdirSync(dirname(cel), { recursive: true });
        copyFileSync(zrodlo, cel);
      }
    },
  };
}

export default {
  root: korzen,

  /* Ścieżki względne, bo pakiet jest oddawany dwiema drogami: rdzeń serwuje go
     spod korzenia nasłuchu, a powłoka Tauri wczytuje z własnego protokołu.
     Ścieżka bezwzględna wiązałaby pakiet z jednym z tych dwóch miejsc. */
  base: './',

  plugins: [kopiaSkryptowBiblioteki()],

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
