import { elementIkony } from '../ikony/ikony';
import { ZrodloDanych } from './model-danych';

/**
 * Nagłówek pulpitu operacyjnego.
 *
 * Jedna odpowiedzialność: tytuł ekranu i jawne oznaczenie pochodzenia liczb.
 * Dopóki odczyt z rdzenia nie nadszedł (`ZrodloDanych.Oczekiwanie`), przy tytule
 * stoi plakietka „oczekiwanie na rdzeń" wraz ze zdaniem wyjaśniającym, żeby pusty
 * ekran nie został wzięty za pomiar. Po pierwszym odczycie plakietka mówi „dane
 * z rdzenia"; pulpit nie zna innych źródeł.
 */
export interface NaglowekPulpitu {
  element: HTMLElement;
  /** Ustawia oznaczenie pochodzenia liczb. */
  oznacz(zrodlo: ZrodloDanych): void;
}

/** Buduje nagłówek pulpitu. */
export function utworzNaglowekPulpitu(zrodlo: ZrodloDanych): NaglowekPulpitu {
  const element = document.createElement('header');
  element.className = 'mc-czolo';

  const blok = document.createElement('div');
  blok.className = 'mc-czolo__blok';

  const tytul = document.createElement('h1');
  tytul.className = 'mc-czolo__tytul';
  tytul.id = 'mc-tytul-pulpitu';
  tytul.textContent = 'Mission Control';

  const podtytul = document.createElement('p');
  podtytul.className = 'mc-czolo__podtytul';
  podtytul.textContent = 'Pulpit operacyjny cyfrowej organizacji — stan pracy, nie okno rozmowy.';

  blok.append(tytul, podtytul);

  const znacznik = document.createElement('span');

  const wyjasnienie = document.createElement('p');
  wyjasnienie.className = 'mc-czolo__wyjasnienie';

  element.append(blok, znacznik);

  const oznacz = (biezace: ZrodloDanych): void => {
    element.dataset.zrodlo = biezace;
    const oczekiwanie = biezace === ZrodloDanych.Oczekiwanie;

    znacznik.className = `dn-plakietka mc-czolo__znacznik ${
      oczekiwanie ? 'dn-plakietka--informacja' : 'dn-plakietka--sukces'
    }`;
    znacznik.replaceChildren();
    if (oczekiwanie) {
      znacznik.append(
        elementIkony('zegar', { rozmiar: 16 }),
        document.createTextNode('oczekiwanie na rdzeń'),
      );
      znacznik.title =
        'Odczyt z rdzenia jeszcze nie nadszedł. Pulpit niczego nie zmyśla — pokazuje stany puste.';
      wyjasnienie.textContent =
        'Odczyt z rdzenia jeszcze nie nadszedł — sekcje pokazują stany puste, nie wymyślone liczby.';
      blok.append(wyjasnienie);
    } else {
      znacznik.append(
        elementIkony('ptaszek-kolo', { rozmiar: 16 }),
        document.createTextNode('dane z rdzenia'),
      );
      znacznik.title = 'Liczby pochodzą z odczytu rdzenia przez kanał kontraktu.';
      wyjasnienie.remove();
    }
  };
  oznacz(zrodlo);

  return { element, oznacz };
}
