// Udostępnianie ikon zestawu i znaku marki po nazwie.
// Jedyny publiczny punkt dostępu — pozostałe warstwy interfejsu sięgają po
// ikony wyłącznie stąd, nie po pliki w `svg/` ani `zasoby-marki/`.

import {
  PROG_SYGNETU_UPROSZCZONEGO,
  zrodloGodla,
  zrodloGodlaMono,
  zrodloZnaku,
  type OdmianaZnaku,
  type PodlozeZnaku,
} from './marka';
import {
  IKONA_PELNOKOLOROWA,
  NAZWY_IKON,
  ZRODLA_IKON,
  type NazwaIkony,
  type NazwaZestawu,
} from './zrodla-ikon';

export {
  IKONA_PELNOKOLOROWA,
  NAZWY_IKON,
  PROG_SYGNETU_UPROSZCZONEGO,
  zrodloGodla,
  zrodloGodlaMono,
  zrodloZnaku,
  type NazwaIkony,
  type NazwaZestawu,
  type OdmianaZnaku,
  type PodlozeZnaku,
};

/**
 * Kanoniczne rozmiary renderowania ikony w pikselach. Wykaz nie wychodzi poza
 * moduł — nikt na zewnątrz nie wybiera rozmiaru z listy; wewnątrz wyznacza
 * rozmiar domyślny, więc największy rozmiar zapisany jest raz.
 */
const ROZMIARY_IKON = [14, 16, 20, 24] as const;

/** Rozmiar przyjmowany w pikselach, gdy wywołanie ikony go nie podaje — największy rozmiar z zestawu `ROZMIARY_IKON`. */
export const ROZMIAR_DOMYSLNY: number = ROZMIARY_IKON[ROZMIARY_IKON.length - 1];

/** Klasa CSS nadawana elementowi ikony, gdy wywołanie nie wskazuje klasy własnej w opcjach; domyślna to `dn-ikona`. */
export const KLASA_DOMYSLNA = 'dn-ikona';

export interface OpcjeIkony {
  /** Bok kwadratu ikony w pikselach. Domyślnie 24. */
  rozmiar?: number;
  /** Nazwa czytana przez technologie wspomagające; pominięta ukrywa ikonę ozdobną przed odczytem. */
  etykieta?: string;
  /** Klasa CSS ikony. Domyślnie `dn-ikona`. */
  klasa?: string;
}

export interface OpcjeZnaku extends OpcjeIkony {
  /** Podłoże osadzenia znaku. Domyślnie ciemne — pasek kokpitu. */
  podloze?: PodlozeZnaku;
}

/** Zwraca źródło SVG ikony o podanej nazwie jako czysty tekst, bez rozbioru i bez osadzenia w dokumencie. */
export function zrodloIkony(nazwa: NazwaIkony): string {
  return ZRODLA_IKON[nazwa];
}

/** Rozstrzyga, czy dowolny podany tekst jest nazwą ikony należącej do zestawu, zawężając typ do `NazwaIkony`. */
export function czyNazwaIkony(tekst: string): tekst is NazwaIkony {
  return Object.prototype.hasOwnProperty.call(ZRODLA_IKON, tekst);
}

// Czytnik powoływany przy pierwszym użyciu, nie przy wczytaniu modułu —
// dzięki temu sięgnięcie po samo źródło ikony nie wymaga środowiska przeglądarki.
let czytnik: DOMParser | undefined;

/**
 * Rozbiera źródło SVG i nadaje mu rozmiar, klasę oraz status dostępnościowy.
 * Wspólny krok dla ikon zestawu i znaku marki — jedno miejsce, w którym
 * powstaje element `<svg>` osadzany w dokumencie.
 */
function zbudujElement(zrodlo: string, opcje: OpcjeIkony): SVGSVGElement {
  const {
    rozmiar = ROZMIAR_DOMYSLNY,
    etykieta,
    klasa = KLASA_DOMYSLNA,
  } = opcje;

  czytnik ??= new DOMParser();
  const rozbior = czytnik.parseFromString(zrodlo, 'image/svg+xml');
  const element = rozbior.documentElement as unknown as SVGSVGElement;

  element.setAttribute('width', String(rozmiar));
  element.setAttribute('height', String(rozmiar));
  element.setAttribute('class', klasa);
  element.setAttribute('focusable', 'false');

  if (etykieta === undefined) {
    element.removeAttribute('aria-label');
    element.setAttribute('role', 'presentation');
    element.setAttribute('aria-hidden', 'true');
  } else {
    element.removeAttribute('aria-hidden');
    element.setAttribute('role', 'img');
    element.setAttribute('aria-label', etykieta);
  }

  return element;
}

/**
 * Buduje gotowy do osadzenia element ikony.
 * Ikona jest osadzana, a nie ładowana obrazem, bo tylko wtedy `currentColor`
 * przejmuje barwę tekstu kontekstu i ikona działa w obu motywach.
 */
export function elementIkony(nazwa: NazwaIkony, opcje: OpcjeIkony = {}): SVGSVGElement {
  return zbudujElement(zrodloIkony(nazwa), opcje);
}

/**
 * Buduje element godła marki. Odmianę dobiera podłoże, a nie motyw — znak
 * ma barwy własne. Od 16 px w dół użyta zostaje odmiana uproszczona.
 */
export function elementGodla(opcje: OpcjeZnaku = {}): SVGSVGElement {
  const { rozmiar = ROZMIAR_DOMYSLNY, podloze = 'ciemne' } = opcje;
  return zbudujElement(zrodloGodla(rozmiar, podloze), opcje);
}

/**
 * Szerokość znaku wyliczona z `viewBox` przy zadanej wysokości. Szerokość
 * musi paść jawnie, inaczej `<svg>` bierze całą szerokość rodzica i logotyp
 * wygląda na mniejszy. Brak `viewBox` zostawia bok kwadratowy, jak
 * w `zbudujElement`.
 */
function szerokoscZnaku(element: SVGSVGElement, wysokosc: number): number {
  const pola = (element.getAttribute('viewBox') ?? '').split(/[\s,]+/).map(Number);
  const [, , szerokoscBoku, wysokoscBoku] = pola;
  if (pola.length !== 4 || !Number.isFinite(szerokoscBoku) || !Number.isFinite(wysokoscBoku)) {
    return wysokosc;
  }
  if (wysokoscBoku === undefined || wysokoscBoku <= 0 || szerokoscBoku === undefined) {
    return wysokosc;
  }
  return Math.round((szerokoscBoku / wysokoscBoku) * wysokosc);
}

/**
 * Buduje element znaku marki w wybranej odmianie — sygnet, logotyp, poziomy
 * albo pionowy. Odmiany złożone mają inne proporcje niż kwadrat: tu `rozmiar`
 * rozstrzyga wysokość, a szerokość wylicza się osobno, przy samym znaku.
 */
export function elementZnaku(
  odmiana: OdmianaZnaku,
  opcje: OpcjeZnaku = {},
): SVGSVGElement {
  const { rozmiar = ROZMIAR_DOMYSLNY, podloze = 'ciemne' } = opcje;
  const element = zbudujElement(zrodloZnaku(odmiana, podloze), opcje);
  element.setAttribute('height', String(rozmiar));
  element.setAttribute('width', String(szerokoscZnaku(element, rozmiar)));
  return element;
}
