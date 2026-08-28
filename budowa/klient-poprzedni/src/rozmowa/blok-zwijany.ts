import './blok.css';

import { elementIkony, type NazwaIkony } from '../ikony/ikony';

/** Interfejs zawiera pełny zestaw nastaw przekazywanych przy tworzeniu bloku zwijanego wewnątrz wpisu rozmowy. */
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

/** Interfejs bloku zwijanego udostępnia element osadzany we wpisie, obszar treści oraz sterowanie widocznością i rozwinięciem. */
export interface BlokZwijany {
  /** Element montowany we wpisie. */
  element: HTMLElement;
  /** Obszar treści do wypełnienia przez wywołującego. */
  tresc: HTMLElement;
  /** Ustawia podtytuł widoczny przy zwiniętym bloku. */
  ustawPodtytul(tekst: string): void;
  /** Pokazuje albo ukrywa cały blok. */
  pokaz(widoczny: boolean): void;
  /** Rozwija albo zwija blok z zewnątrz, wywoływane wyłącznie przy zmianie trybu widoku. */
  ustawRozwiniecie(otwarty: boolean): void;
}

/** Funkcja tworzy blok zwijany oparty na natywnych elementach details i summary, z ikoną, tytułem i podtytułem w nagłówku. */
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

/** Funkcja buduje wiersz pary nazwa-wartość wyświetlany w treści bloku zwijanego, z wartością zapisaną krojem stałej szerokości. */
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

/** Funkcja buduje blok tekstu w kroju stałej szerokości, przeznaczony na treść wywołania, polecenia albo wynik narzędzia. */
export function blokKodu(tresc: string): HTMLElement {
  const element = document.createElement('pre');
  element.className = 'dc-kod dn-kod';
  element.textContent = tresc;
  return element;
}
