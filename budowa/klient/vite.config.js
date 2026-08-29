// Nastawa budowania pakietu klienta.
//
// Skrypty biblioteki są funkcjami domkniętymi, nie modułami ECMAScript, więc
// Vite ich nie pakuje. Wtyczka `kopiaSkryptowBiblioteki` przenosi je do
// `dist/zasoby/` z zachowaniem ścieżki względnej.

import { readFileSync } from 'node:fs';
import { copyFileSync, mkdirSync } from 'node:fs';
import { dirname, join } from 'node:path';
import { fileURLToPath } from 'node:url';

const katalogKlienta = new URL('.', import.meta.url).pathname;

// Arkusze i wygenerowany kontrakt leżą poza katalogiem klienta; serwer poglądowy
// potrzebuje zgody na ich odczyt.
const katalogRepozytorium = new URL('../../', import.meta.url).pathname;

const KATALOG_ZASOBOW = new URL('../../design/zasoby/', import.meta.url);

// Pliki biblioteki kopiowane do wydania; ścieżki względne wobec KATALOG_ZASOBOW.
const PLIKI_BIBLIOTEKI = [
  'wspolne.js',
  'narzedzia-okien.js',
  'ekran-startowy.js',
  'powloka.js',
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
  'okna/studio.js',
  'okna/wejscie/tresci.js',
  'okna/wejscie/ikony.js',
  'okna/wejscie/skladniki/belka-okna.js',
  'okna/wejscie/skladniki/kolumna-tozsamosci.js',
  'okna/wejscie/skladniki/naglowek-ekranu.js',
  'okna/wejscie/skladniki/lista-etapow.js',
  'okna/wejscie/skladniki/baner.js',
  'okna/wejscie/skladniki/fraza-nawigacyjna.js',
  'okna/wejscie/skladniki/pas-dzialan.js',
  'okna/wejscie/skladniki/zakladki-pigulki.js',
  'okna/wejscie/skladniki/pole-tekstowe.js',
  'okna/wejscie/skladniki/pole-hasla.js',
  'okna/wejscie/skladniki/miernik-sily.js',
  'okna/wejscie/skladniki/pole-sesji.js',
  'okna/wejscie/skladniki/metody-logowania.js',
  'okna/wejscie/skladniki/pole-kodu.js',
  'okna/wejscie/skladniki/kroki-odzyskiwania.js',
  'okna/wejscie/skladniki/pasek-postepu.js',
  'okna/wejscie/ekrany/uruchomienie.js',
  'okna/wejscie/ekrany/dostep.js',
  'okna/wejscie/ekrany/przygotowanie.js',
  'okna/wejscie/montaz.js',
  'powloki.js',
  'bryla.js',
  'okna/instalator/tresci.js',
  'tresci/licencja.js',
  'okna/instalator/narzedzia.js',
  'okna/instalator/ikony.js',
  'okna/instalator/stany.js',
  'okna/instalator/skladniki/naglowek-bloku.js',
  'okna/instalator/skladniki/nawigacja-krokow.js',
  'okna/instalator/skladniki/naglowek-ekranu.js',
  'okna/instalator/skladniki/blok-danych.js',
  'okna/instalator/skladniki/pole-wyboru.js',
  'okna/instalator/skladniki/pole-sciezki.js',
  'okna/instalator/skladniki/karta-opcji.js',
  'okna/instalator/skladniki/baner.js',
  'okna/instalator/skladniki/lista-etapow.js',
  'okna/instalator/skladniki/pasek-postepu.js',
  'okna/instalator/skladniki/pasek-szczegolow.js',
  'okna/instalator/skladniki/pas-dzialan.js',
  'okna/instalator/skladniki/fraza-nawigacyjna.js',
  'okna/instalator/skladniki/tekst-ciagly.js',
  'okna/instalator/skladniki/dokument.js',
  'okna/instalator/skladniki/wiersz-miary.js',
  'okna/instalator/skladniki/blok-bledu.js',
  'okna/instalator/skladniki/odsylacz-pomocy.js',
  'okna/instalator/skladniki/okno-dialogowe.js',
  'okna/instalator/ekrany/1-wymagania.js',
  'okna/instalator/ekrany/2-licencja.js',
  'okna/instalator/ekrany/3-wersja.js',
  'okna/instalator/ekrany/4-lokalizacja.js',
  'okna/instalator/ekrany/5-instalacja.js',
  'okna/instalator/ekrany/6-podsumowanie.js',
  'okna/instalator/montaz.js',
  'okna/instalator.js',
  'kreator.css',
  'okna/przeplyw-wejscia.js',
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
      for (const wzgledna of PLIKI_BIBLIOTEKI) {
        const zrodlo = fileURLToPath(new URL(wzgledna, KATALOG_ZASOBOW));
        const cel = join(katalogKlienta, 'dist', 'zasoby', wzgledna);
        mkdirSync(dirname(cel), { recursive: true });
        copyFileSync(zrodlo, cel);
      }
    },
  };
}

/* Data składania pakietu. Katalog wydań stanowi, że wydania różni data, nie
   numer — sam numer w stopce nie nazywa więc, które wydanie Operator ma przed
   sobą. Wartość wchodzi przy budowaniu, bo w przeglądarce nie ma jej skąd wziąć. */
const DATA_SKLADANIA = new Date().toISOString().slice(0, 10);


/**
 * Wnętrze okna roboczego Studia pochodzi z prototypu, nie z kodu klienta.
 *
 * `zasoby/powloka.js` skleja powłokę wokół treści szablonu `#dn-tresc-okna`,
 * a treścią tą jest blok `main.st-okno-robocze` ze źródła kształtu. Wstrzyknięcie
 * przy budowaniu trzyma jedno źródło prawdy: zmiana prototypu wchodzi do
 * aplikacji przebudowaniem, bez przepisywania znacznika ręką.
 */
function trescOknaZPrototypu() {
  return {
    name: 'tresc-okna-z-prototypu',
    transformIndexHtml(html) {
      const zrodlo = new URL('../../design/05-okna/moduly/studio.html', import.meta.url);
      const prototyp = readFileSync(zrodlo, 'utf8');
      const poczatek = prototyp.indexOf('<main class="st-okno-robocze"');
      const koniec = prototyp.indexOf('</main>', poczatek);
      if (poczatek < 0 || koniec < 0) {
        throw new Error('prototyp Studia nie niesie bloku st-okno-robocze');
      }
      const blok = prototyp.slice(poczatek, koniec + '</main>'.length);
      return html.replace(
        '<template id="dn-tresc-okna"></template>',
        `<template id="dn-tresc-okna">${blok}</template>`,
      );
    },
  };
}

export default {
  root: katalogKlienta,

  define: {
    __DATA_SKLADANIA__: JSON.stringify(DATA_SKLADANIA),
  },

  /* Ścieżki względne, bo pakiet jest oddawany dwiema drogami: rdzeń serwuje go
     spod korzenia nasłuchu, a powłoka Tauri wczytuje z własnego protokołu.
     Ścieżka bezwzględna wiązałaby pakiet z jednym z tych dwóch miejsc. */
  base: './',

  plugins: [kopiaSkryptowBiblioteki(), trescOknaZPrototypu()],

  build: {
    rollupOptions: {
      input: {
        aplikacja: new URL('index.html', import.meta.url).pathname,
        instalator: new URL('instalator.html', import.meta.url).pathname,
      },
    },
    outDir: 'dist',
    emptyOutDir: true,
    /* Kroje idą plikami, nie treścią wplecioną w arkusz: wpisany krój rośnie
       w każdym pakiecie o pełny rozmiar podzbioru, a przeglądarka nie może go
       pominąć na podstawie zakresu znaków. */
    assetsInlineLimit: 0,
  },

  server: {
    fs: { allow: [katalogRepozytorium] },
  },
};
