// Skrypty biblioteki nie są modułami ECMAScript, więc Vite ich nie pakuje.

import { readFileSync } from 'node:fs';
import { copyFileSync, mkdirSync } from 'node:fs';
import { dirname, join } from 'node:path';
import { fileURLToPath } from 'node:url';

const katalogKlienta = new URL('.', import.meta.url).pathname;

const katalogRepozytorium = new URL('../../', import.meta.url).pathname;

const KATALOG_ZASOBOW = new URL('../../design/zasoby/', import.meta.url);

const PLIKI_BIBLIOTEKI = [
  'wspolne.js',
  'narzedzia-okien.js',
  'ekran-startowy.js',
  'powloka.js',
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
  /* Poza wydaniem zostają `okna/centrum-wejscie.js`, `okna/centrum-dowodzenia.js`
     i `prototyp.js`: przełączały stan funkcji globalnych bez wywołania rdzenia,
     a zachowania okna niosą `src/wiazanie/centrum.ts` i `podstawa-okna.ts`. */
  'okna/centrum-obszar.js',
  'okna/centrum-dymki.js',
];

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

// Wydania różni data, nie numer; w przeglądarce nie ma jej skąd wziąć.
const DATA_SKLADANIA = new Date().toISOString().slice(0, 10);


const WNETRZA_OKIEN = [
  {
    gniazdo: 'dn-tresc-okna',
    prototyp: '../../design/05-okna/przeplyw/centrum-dowodzenia.html',
    otwarcie: '<div class="dn-obszar">',
    zamkniecie: '<!-- /dn-obszar -->',
  },
  // Osobny prototyp na środowisko, bo kafle niosą własne znaki modułów.
  {
    gniazdo: 'dn-tresc-przedsionek-talkin',
    prototyp: '../../design/05-okna/srodowiska/talkin-przedsionek.html',
    otwarcie: '<div class="pd-cialo">',
    zamkniecie: null,
  },
  {
    gniazdo: 'dn-tresc-przedsionek-workspace',
    prototyp: '../../design/05-okna/srodowiska/workspace-przedsionek.html',
    otwarcie: '<div class="pd-cialo">',
    zamkniecie: null,
  },
  {
    gniazdo: 'dn-tresc-przedsionek-codestudio',
    prototyp: '../../design/05-okna/srodowiska/codestudio-przedsionek.html',
    otwarcie: '<div class="pd-cialo">',
    zamkniecie: null,
  },
  {
    gniazdo: 'dn-tresc-przedsionek-multitaskingai',
    prototyp: '../../design/05-okna/srodowiska/multitaskingai-przedsionek.html',
    otwarcie: '<div class="pd-cialo">',
    zamkniecie: null,
  },
  {
    gniazdo: 'dn-tresc-studio',
    prototyp: '../../design/05-okna/moduly/studio.html',
    otwarcie: '<main class="st-okno-robocze"',
    zamkniecie: '</main>',
  },
];

function blokPrototypu(tresc, otwarcie, zamkniecie) {
  const poczatek = tresc.indexOf(otwarcie);
  if (poczatek < 0) return null;
  if (zamkniecie === null) return blokZbalansowany(tresc, poczatek);
  const koniec = tresc.indexOf(zamkniecie, poczatek);
  if (koniec < 0) return null;
  return tresc.slice(poczatek, koniec + zamkniecie.length);
}

// Prototyp bez komentarza zamykającego domyka się liczeniem otwarć elementu.
function blokZbalansowany(tresc, poczatek) {
  const znacznik = /<(\/?)div\b[^>]*>/g;
  znacznik.lastIndex = poczatek;
  let glebokosc = 0;
  let dopasowanie = znacznik.exec(tresc);
  while (dopasowanie !== null) {
    glebokosc += dopasowanie[1] === '/' ? -1 : 1;
    if (glebokosc === 0) return tresc.slice(poczatek, znacznik.lastIndex);
    dopasowanie = znacznik.exec(tresc);
  }
  return null;
}

// Reguły stoją w prototypie, nie w arkuszach; bez nich wnętrze traci układ.
function stylWpisany(tresc) {
  const dopasowanie = /<style[^>]*>([\s\S]*?)<\/style>/.exec(tresc);
  return dopasowanie === null ? '' : dopasowanie[1].trim();
}

function wnetrzaOkienZPrototypow() {
  return {
    name: 'wnetrza-okien-z-prototypow',
    transformIndexHtml(html) {
      let wynik = html;
      const reguly = new Map();
      for (const okno of WNETRZA_OKIEN) {
        const zrodlo = new URL(okno.prototyp, import.meta.url);
        const tresc = readFileSync(zrodlo, 'utf8');
        const blok = blokPrototypu(tresc, okno.otwarcie, okno.zamkniecie);
        if (blok === null) {
          throw new Error(`prototyp ${okno.prototyp} nie niesie bloku ${okno.otwarcie}`);
        }
        wynik = wynik.replace(
          `<template id="${okno.gniazdo}"></template>`,
          `<template id="${okno.gniazdo}">${blok}</template>`,
        );
        const styl = stylWpisany(tresc);
        if (styl !== '') reguly.set(okno.prototyp, styl);
      }
      const arkusz = [...reguly.values()].join('\n');
      return arkusz === ''
        ? wynik
        : wynik.replace('</head>', `<style>\n${arkusz}\n</style>\n</head>`);
    },
  };
}

export default {
  root: katalogKlienta,

  define: {
    __DATA_SKLADANIA__: JSON.stringify(DATA_SKLADANIA),
  },

  // Pakiet idzie dwiema drogami: nasłuchem rdzenia i protokołem powłoki Tauri.
  base: './',

  plugins: [kopiaSkryptowBiblioteki(), wnetrzaOkienZPrototypow()],

  build: {
    outDir: 'dist',
    emptyOutDir: true,
    // Krój wpisany w arkusz rośnie o pełny podzbiór i nie da się go pominąć.
    assetsInlineLimit: 0,
  },

  server: {
    fs: { allow: [katalogRepozytorium] },
  },
};
