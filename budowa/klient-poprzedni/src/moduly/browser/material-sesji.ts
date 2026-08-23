import type { BrowserSnapshot } from '../../../../shared/contract';

/**
 * Materiał przechwycony w toku sesji przeglądania: zrzuty, archiwa stron
 * i monitory zmian.
 *
 * Jedna odpowiedzialność: pamięć tego, co Capture & Monitor Panel zgromadził.
 * Bez odczytu z rdzenia i bez ani jednego elementu widoku — wzorem
 * `zebrane-w-sesji.ts`.
 *
 * Zbiór nie jest odbiciem rdzenia, bo odbijać nie ma czego: zapis pozycji jako
 * wytworu sesji ma komendę `browser.artifact.add`, ale komendy odczytu wykazu
 * wytworów okna kontrakt nie niesie. Pozycją jest migawka, którą rdzeń
 * oddał na żądanie Operatora; po przeładowaniu karty wykaz zaczyna się od nowa,
 * a same migawki zostają w rdzeniu pod swoimi identyfikatorami. Panel mówi
 * o tym wprost, zamiast obiecywać trwałość, której nie ma.
 *
 * Rodzaj pozycji rozstrzyga zawartość odpowiedzi, nie zamówienie okna. Rdzeń
 * pobiera dziś stronę bez uruchamiania przeglądarki (`przegladarka_pobieranie.go`),
 * więc migawka zamówiona ze zrzutem potrafi przyjść bez niego — pozycja nazywa
 * się wtedy archiwum albo treścią, a nie zrzutem. Rozpoznanie po odpowiedzi
 * nazwie ją zrzutem sama, gdy obraz zacznie przychodzić.
 *
 * Monitor stoi na adresie, nie na migawce: sprawdzenie polega na ponownym
 * pobraniu tej samej strony i zestawieniu treści z zapamiętaną. Adres jest
 * więc jego jedynym rozsądnym kluczem — dwa monitory tej samej strony
 * pilnowałyby dokładnie tego samego.
 */

/** Co pozycja materiału naprawdę niesie — rozstrzygnięte po odpowiedzi rdzenia. */
export type RodzajPrzechwycenia = 'zrzut' | 'archiwum' | 'tresc';

export interface Przechwycenie {
  migawka: BrowserSnapshot;
  rodzaj: RodzajPrzechwycenia;
}

/** Wynik ostatniego sprawdzenia monitora. */
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
  /**
   * Zakłada monitor na stronie migawki i oddaje `false`, gdy monitor tego
   * adresu już stoi — okno ma o tym powiedzieć zamiast meldować założenie,
   * którego nie było.
   */
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

/** Czym pozycja jest według tego, co rdzeń w migawce oddał. */
function rozpoznajRodzaj(migawka: BrowserSnapshot): RodzajPrzechwycenia {
  if ((migawka.screenshotRef ?? '').trim() !== '') return 'zrzut';
  if ((migawka.html ?? '').trim() !== '') return 'archiwum';
  return 'tresc';
}

/** Nazwa rodzaju pozycji mówiona Operatorowi. */
export function nazwaRodzaju(rodzaj: RodzajPrzechwycenia): string {
  if (rodzaj === 'zrzut') return 'zrzut ekranu';
  if (rodzaj === 'archiwum') return 'archiwum strony';
  return 'treść strony';
}
