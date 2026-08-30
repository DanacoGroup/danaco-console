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
  /* `okna/centrum-wejscie.js` do wydania nie wchodzi: przenosi przeglądarkę na
     osobny plik prototypu, a produkt zmienia wnętrze okna w miejscu. */
  'okna/centrum-dowodzenia.js',
  'okna/centrum-obszar.js',
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
    gniazdo: 'dn-tresc-przedsionek',
    prototyp: '../../design/05-okna/srodowiska/talkin-przedsionek.html',
    otwarcie: '<main class="pd-plotno">',
    zamkniecie: '</main>',
  },
  {
    gniazdo: 'dn-tresc-konfiguracja',
    prototyp: '../../design/05-okna/platformowe/konfiguracja.html',
    otwarcie: '<div class="sta-cialo">',
    zamkniecie: null,
  },
  {
    gniazdo: 'dn-tresc-pula-kont',
    prototyp: '../../design/05-okna/platformowe/konfiguracja.html',
    otwarcie: '<dialog class="dn-modal kf-modal-pula"',
    zamkniecie: '</dialog>',
  },
  {
    gniazdo: 'dn-tresc-studio',
    prototyp: '../../design/05-okna/moduly/studio.html',
    otwarcie: '<main class="st-okno-robocze"',
    zamkniecie: '</main>',
  },
  {
    gniazdo: 'dn-tresc-agents',
    prototyp: '../../design/05-okna/moduly/agents.html',
    otwarcie: '<div class="sta-cialo">',
    zamkniecie: '<!-- /sta-cialo -->',
  },
  {
    gniazdo: 'dn-tresc-apps',
    prototyp: '../../design/05-okna/moduly/apps.html',
    otwarcie: '<div class="sta-cialo">',
    zamkniecie: '<!-- /sta-cialo -->',
  },
  {
    gniazdo: 'dn-tresc-assistant',
    prototyp: '../../design/05-okna/moduly/assistant.html',
    otwarcie: '<div class="sta-cialo">',
    zamkniecie: '<!-- /sta-cialo -->',
  },
  {
    gniazdo: 'dn-tresc-automations',
    prototyp: '../../design/05-okna/moduly/automations.html',
    otwarcie: '<div class="sta-cialo">',
    zamkniecie: '<!-- /sta-cialo -->',
  },
  {
    gniazdo: 'dn-tresc-browser',
    prototyp: '../../design/05-okna/moduly/browser.html',
    otwarcie: '<div class="sta-cialo">',
    zamkniecie: '<!-- /sta-cialo -->',
  },
  {
    gniazdo: 'dn-tresc-design',
    prototyp: '../../design/05-okna/moduly/design.html',
    otwarcie: '<div class="sta-cialo">',
    zamkniecie: '<!-- /sta-cialo -->',
  },
  {
    gniazdo: 'dn-tresc-developer',
    prototyp: '../../design/05-okna/moduly/developer.html',
    otwarcie: '<div class="sta-cialo">',
    zamkniecie: '<!-- /sta-cialo -->',
  },
  {
    gniazdo: 'dn-tresc-diagnostics',
    prototyp: '../../design/05-okna/moduly/diagnostics.html',
    otwarcie: '<div class="sta-cialo">',
    zamkniecie: '<!-- /sta-cialo -->',
  },
  {
    gniazdo: 'dn-tresc-library',
    prototyp: '../../design/05-okna/moduly/library.html',
    otwarcie: '<div class="sta-cialo">',
    zamkniecie: '<!-- /sta-cialo -->',
  },
  {
    gniazdo: 'dn-tresc-research',
    prototyp: '../../design/05-okna/moduly/research.html',
    otwarcie: '<div class="sta-cialo">',
    zamkniecie: '<!-- /sta-cialo -->',
  },
  {
    gniazdo: 'dn-tresc-roundtable',
    prototyp: '../../design/05-okna/moduly/roundtable.html',
    otwarcie: '<div class="sta-cialo">',
    zamkniecie: '<!-- /sta-cialo -->',
  },
  {
    gniazdo: 'dn-tresc-terminal',
    prototyp: '../../design/05-okna/moduly/terminal.html',
    otwarcie: '<div class="sta-cialo">',
    zamkniecie: '<!-- /sta-cialo -->',
  },
  {
    gniazdo: 'dn-tresc-translate',
    prototyp: '../../design/05-okna/moduly/translate.html',
    otwarcie: '<div class="sta-cialo">',
    zamkniecie: '<!-- /sta-cialo -->',
  },
  {
    gniazdo: 'dn-tresc-workspace',
    prototyp: '../../design/05-okna/moduly/workspace.html',
    otwarcie: '<div class="sta-cialo">',
    zamkniecie: '<!-- /sta-cialo -->',
  },
];

/** Wycina z prototypu blok między znacznikami i oddaje go wraz z zamknięciem. */
function blokPrototypu(tresc, otwarcie, zamkniecie) {
  const poczatek = tresc.indexOf(otwarcie);
  if (poczatek < 0) return null;
  if (zamkniecie === null) return blokZbalansowany(tresc, poczatek);
  const koniec = tresc.indexOf(zamkniecie, poczatek);
  if (koniec < 0) return null;
  return tresc.slice(poczatek, koniec + zamkniecie.length);
}

/* Bloki bez znacznika zamykającego domyka się liczeniem otwarć i zamknięć
   elementu. Prototypy platformowe nie niosą komentarza zamykającego, a stały
   znacznik końca ramy zabrałby ze sobą zamknięcia elementów nadrzędnych. */
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

/* Reguły stojące w prototypie, a nie w arkuszach warstwy projektowej. Wnętrze
   okna bez nich składa się bez układu, więc jadą do dokumentu razem z treścią.
   Klasy prototypów są przedrostkowane osobno dla każdego okna, więc reguły
   zebrane z wielu plików nie nachodzą na siebie. */
function stylWpisany(tresc) {
  const dopasowanie = /<style[^>]*>([\s\S]*?)<\/style>/.exec(tresc);
  return dopasowanie === null ? '' : dopasowanie[1].trim();
}

/** Wstawia wnętrza okien z prototypów w puste szablony dokumentu klienta. */
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
