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


// Wnętrza okien biorą się z prototypów Właściciela w `design/05-okna/`.
// `zasoby/powloka.js` skleja powłokę wokół szablonu `#dn-tresc-okna`; szablon
// `#dn-tresc-studio` stoi obok i wchodzi na jego miejsce przy otwarciu modułu.
const WNETRZA_OKIEN = [
  {
    gniazdo: 'dn-tresc-okna',
    prototyp: '../../design/05-okna/przeplyw/centrum-dowodzenia.html',
    otwarcie: '<div class="dn-obszar">',
    zamkniecie: '<!-- /dn-obszar -->',
  },
  {
    gniazdo: 'dn-tresc-studio',
    prototyp: '../../design/05-okna/moduly/studio.html',
    otwarcie: '<main class="st-okno-robocze"',
    zamkniecie: '</main>',
  },
];

/** Wycina z prototypu blok między znacznikami i oddaje go wraz z zamknięciem. */
function blokPrototypu(tresc, otwarcie, zamkniecie) {
  const poczatek = tresc.indexOf(otwarcie);
  if (poczatek < 0) return null;
  const koniec = tresc.indexOf(zamkniecie, poczatek);
  if (koniec < 0) return null;
  return tresc.slice(poczatek, koniec + zamkniecie.length);
}

/** Wstawia wnętrza okien z prototypów w puste szablony dokumentu klienta. */
function wnetrzaOkienZPrototypow() {
  return {
    name: 'wnetrza-okien-z-prototypow',
    transformIndexHtml(html) {
      let wynik = html;
      for (const okno of WNETRZA_OKIEN) {
        const zrodlo = new URL(okno.prototyp, import.meta.url);
        const blok = blokPrototypu(readFileSync(zrodlo, 'utf8'), okno.otwarcie, okno.zamkniecie);
        if (blok === null) {
          throw new Error(`prototyp ${okno.prototyp} nie niesie bloku ${okno.otwarcie}`);
        }
        wynik = wynik.replace(
          `<template id="${okno.gniazdo}"></template>`,
          `<template id="${okno.gniazdo}">${blok}</template>`,
        );
      }
      return wynik;
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

  plugins: [kopiaSkryptowBiblioteki(), wnetrzaOkienZPrototypow()],

  build: {
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
