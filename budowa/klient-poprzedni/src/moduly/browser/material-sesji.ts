/**
 * Materiał przechwycony w toku sesji przeglądania: zrzuty, archiwa stron
 * i monitory zmian. Pamięć tego, co zgromadziło okno przechwytywania, bez
 * odczytu z rdzenia i bez elementów widoku.
 */

import type { BrowserSnapshot } from '../../../../shared/contract';

/**
 * Czym pozycja materiału naprawdę jest według zawartości odpowiedzi rdzenia:
 * zrzut ekranu, archiwum strony albo sama jej treść. Rodzaju nie rozstrzyga
 * zamówienie okna, lecz pola oddanej migawki.
 */
export type RodzajPrzechwycenia = 'zrzut' | 'archiwum' | 'tresc';

export interface Przechwycenie {
  migawka: BrowserSnapshot;
  rodzaj: RodzajPrzechwycenia;
}

/**
 * Wynik ostatniego sprawdzenia monitora zmian: monitor nietknięty od
 * założenia, sprawdzony bez różnicy w treści albo sprawdzony ze stwierdzoną
 * zmianą względem treści odniesienia.
 */
export type WynikMonitora = 'nietkniety' | 'bez-zmian' | 'zmiana';

export interface MonitorZmian {
  /** Adres pilnowanej strony — klucz monitora. */
  adres: string;
  tytul: string;
  /** Treść strony z chwili założenia monitora — podstawa zestawienia. */
  trescOdniesienia: string;
  zalozonyO: number;
  /** Czas ostatniego sprawdzenia w milisekundach epoki; zero znaczy „ani razu". */
  sprawdzonyO: number;
  wynik: WynikMonitora;
  /** Różnica długości treści w znakach względem odniesienia. */
  roznicaZnakow: number;
}

export interface MaterialSesji {
  /** Przechwycenia od najnowszego. */
  przechwycenia(): readonly Przechwycenie[];
  /** Monitory w kolejności założenia. */
  monitory(): readonly MonitorZmian[];
  /** Dopisuje migawkę jako pozycję materiału i oddaje jej rozpoznany rodzaj. */
  dopisz(migawka: BrowserSnapshot): RodzajPrzechwycenia;
  // Zakłada monitor na stronie migawki; oddaje `false`, gdy monitor tego adresu już stoi.
  zalozMonitor(migawka: BrowserSnapshot): boolean;
  /** Zapisuje wynik sprawdzenia monitora; `null`, gdy monitora tego adresu nie ma. */
  zapiszSprawdzenie(adres: string, migawka: BrowserSnapshot): MonitorZmian | null;
  usunMonitor(adres: string): void;
}

export function utworzMaterialSesji(oglos: () => void): MaterialSesji {
  const przechwycenia: Przechwycenie[] = [];
  const monitory: MonitorZmian[] = [];

  return {
    przechwycenia: () => przechwycenia,
    monitory: () => monitory,

    dopisz(migawka) {
      const rodzaj = rozpoznajRodzaj(migawka);
      const pozycja = przechwycenia.findIndex((wpis) => wpis.migawka.id === migawka.id);
      if (pozycja === -1) przechwycenia.unshift({ migawka, rodzaj });
      else przechwycenia[pozycja] = { migawka, rodzaj };
      oglos();
      return rodzaj;
    },

    zalozMonitor(migawka) {
      if (monitory.some((wpis) => wpis.adres === migawka.url)) return false;
      monitory.push({
        adres: migawka.url,
        tytul: (migawka.title ?? '').trim(),
        trescOdniesienia: migawka.text ?? '',
        zalozonyO: migawka.capturedAt,
        sprawdzonyO: 0,
        wynik: 'nietkniety',
        roznicaZnakow: 0,
      });
      oglos();
      return true;
    },

    zapiszSprawdzenie(adres, migawka) {
      const monitor = monitory.find((wpis) => wpis.adres === adres);
      if (monitor === undefined) return null;
      const teraz = migawka.text ?? '';
      monitor.sprawdzonyO = migawka.capturedAt;
      monitor.wynik = teraz === monitor.trescOdniesienia ? 'bez-zmian' : 'zmiana';
      monitor.roznicaZnakow = teraz.length - monitor.trescOdniesienia.length;
      oglos();
      return monitor;
    },

    usunMonitor(adres) {
      const pozycja = monitory.findIndex((wpis) => wpis.adres === adres);
      if (pozycja === -1) return;
      monitory.splice(pozycja, 1);
      oglos();
    },
  };
}

/**
 * Rozpoznaje rodzaj pozycji z pól oddanej migawki: niepuste pole screenshotRef
 * daje zrzut, niepuste pole html daje archiwum, a pozostałe przypadki samą
 * treść strony.
 */
function rozpoznajRodzaj(migawka: BrowserSnapshot): RodzajPrzechwycenia {
  if ((migawka.screenshotRef ?? '').trim() !== '') return 'zrzut';
  if ((migawka.html ?? '').trim() !== '') return 'archiwum';
  return 'tresc';
}

/**
 * Oddaje nazwę rodzaju pozycji w brzmieniu pokazywanym w oknie: zrzut ekranu,
 * archiwum strony albo treść strony, zależnie od wartości wyliczenia
 * RodzajPrzechwycenia.
 */
export function nazwaRodzaju(rodzaj: RodzajPrzechwycenia): string {
  if (rodzaj === 'zrzut') return 'zrzut ekranu';
  if (rodzaj === 'archiwum') return 'archiwum strony';
  return 'treść strony';
}
