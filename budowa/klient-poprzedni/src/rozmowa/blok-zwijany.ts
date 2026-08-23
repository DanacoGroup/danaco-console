import './blok.css';

import { elementIkony, type NazwaIkony } from '../ikony/ikony';

/** Nastawy bloku zwijanego. */
export interface NastawyBloku {
  /** Tytuł bloku — etykieta wersalikowa. */
  tytul: string;
  /** Ikona przy tytule. */
  ikona: NazwaIkony;
  /** Czy blok jest rozwinięty przy założeniu. */
  otwarty?: boolean;
  /** Odmiana wizualna: neutralna, akcentowana albo błędna. */
  odmiana?: 'neutralna' | 'akcent' | 'blad';
}

/** Blok zwijany — nagłówek klikalny i treść pod nim. */
export interface BlokZwijany {
  /** Element montowany we wpisie. */
  element: HTMLElement;
  /** Obszar treści do wypełnienia przez wywołującego. */
  tresc: HTMLElement;
  /** Ustawia podtytuł widoczny przy zwiniętym bloku. */
  ustawPodtytul(tekst: string): void;
  /** Pokazuje albo ukrywa cały blok. */
  pokaz(widoczny: boolean): void;
  /**
   * Rozwija albo zwija blok z zewnątrz — używa tego widok transkryptu.
   *
   * Wołane tylko przy zmianie trybu, nigdy przy odświeżeniu treści: blok
   * przestawiony ręcznie ma zostać w stanie, w jakim go zostawiono, także gdy
   * w środku tury przychodzą kolejne fragmenty.
   */
  ustawRozwiniecie(otwarty: boolean): void;
}

/**
 * Blok zwijany zbudowany na `details`/`summary`.
 *
 * Zwinięcie jest zachowaniem natywnym przeglądarki, więc blok działa
 * klawiaturą, ma poprawną semantykę dla technologii wspomagających i obywa się
 * bez obsługiwaczy zdarzeń. Brak treści oznacza ukrycie bloku, nie wyszarzenie.
 */
export function utworzBlokZwijany(nastawy: NastawyBloku): BlokZwijany {
  const element = document.createElement('details');
  element.className = `dc-blok dc-blok--${nastawy.odmiana ?? 'neutralna'}`;
  element.open = nastawy.otwarty === true;

  const naglowek = document.createElement('summary');
  naglowek.className = 'dc-blok__naglowek';

  const tytul = document.createElement('span');
  tytul.className = 'dc-blok__tytul dn-etykieta-wersalikowa';
  tytul.textContent = nastawy.tytul;

  const podtytul = document.createElement('span');
  podtytul.className = 'dc-blok__podtytul';

  naglowek.append(
    elementIkony(nastawy.ikona, { rozmiar: 16, klasa: 'dn-ikona dc-blok__ikona' }),
    tytul,
    podtytul,
  );

  const tresc = document.createElement('div');
  tresc.className = 'dc-blok__tresc';

  element.append(naglowek, tresc);

  return {
    element,
    tresc,
    ustawPodtytul(tekst) {
      podtytul.textContent = tekst;
    },
    pokaz(widoczny) {
      element.hidden = !widoczny;
    },
    ustawRozwiniecie(otwarty) {
      element.open = otwarty;
    },
  };
}

/** Wiersz pary „nazwa — wartość" w treści bloku. */
export function wierszDanych(nazwa: string, wartosc: string): HTMLElement {
  const wiersz = document.createElement('div');
  wiersz.className = 'dc-dana';

  const etykieta = document.createElement('span');
  etykieta.className = 'dc-dana__nazwa';
  etykieta.textContent = nazwa;

  const tresc = document.createElement('span');
  tresc.className = 'dc-dana__wartosc dn-kod dn-kod--wiersz';
  tresc.textContent = wartosc;

  wiersz.append(etykieta, tresc);
  return wiersz;
}

/** Blok tekstu w kroju monospacjowym — argv, prompt, wejście narzędzia. */
export function blokKodu(tresc: string): HTMLElement {
  const element = document.createElement('pre');
  element.className = 'dc-kod dn-kod';
  element.textContent = tresc;
  return element;
}
