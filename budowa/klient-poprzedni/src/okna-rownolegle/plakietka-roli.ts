import { WindowRole } from '../../../shared/contract';
import { elementIkony, type NazwaIkony } from '../ikony/ikony';
import { nazwaRoli } from '../okno-komunikacji/etykiety-okna';

/** Plakietka roli okna — widoczna na pierwszy rzut oka w nagłówku gniazda, niosąca ikonę i etykietę zamiast samej barwy. */
export interface PlakietkaRoli {
  element: HTMLElement;
  /** Przestawia plakietkę na inną rolę. */
  ustaw(rola: WindowRole): void;
}

/**
 * Plakietka roli okna w pętli koordynator-wykonawca nigdy nie sygnalizuje roli samą barwą: wariant sygnałowy dostaje wyłącznie koordynator, wyłącznie na powierzchni pigułki.
 */
export function utworzPlakietkeRoli(rola: WindowRole): PlakietkaRoli {
  const element = document.createElement('span');

  const nosnik = document.createElement('span');
  nosnik.className = 'dn-okna__rola-ikona';

  const napis = document.createElement('span');

  element.append(nosnik, napis);

  function ustaw(nowa: WindowRole): void {
    element.className = `dn-plakietka dn-plakietka--rola dn-okna__rola ${wariant(nowa)}`;
    element.dataset.rola = nowa;
    element.title = `Rola okna: ${nazwaRoli(nowa)}`;
    nosnik.replaceChildren(elementIkony(ikonaRoli(nowa), { rozmiar: 16 }));
    napis.textContent = nazwaRoli(nowa);
  }

  ustaw(rola);
  return { element, ustaw };
}

/**
 * Ikona roli.
 *
 * Waga — koordynator rozdziela i ocenia pracę; strzałka uruchomienia —
 * wykonawca pracę wykonuje; sylwetka — okno samodzielne pracuje wprost
 * z operatorem, poza pętlą.
 */
function ikonaRoli(rola: WindowRole): NazwaIkony {
  switch (rola) {
    case WindowRole.Coordinator:
      return 'waga';
    case WindowRole.Executor:
      return 'uruchom';
    case WindowRole.Standalone:
      return 'uzytkownik';
  }
}

/** Wariant barwny plakietki roli; okno samodzielne, bez pary, zostaje neutralne wobec sygnału koordynatora. */
function wariant(rola: WindowRole): string {
  switch (rola) {
    case WindowRole.Coordinator:
      return 'dn-plakietka--sygnal';
    case WindowRole.Executor:
      return 'dn-plakietka--informacja';
    case WindowRole.Standalone:
      return '';
  }
}
