import type { ExecutionEnv } from '../../../shared/contract';
import { nazwaSrodowiska } from './etykiety-okna';
import type { ProfilModulu } from './profil-modulu';

/**
 * Kontekst okna widoczny w panelu — to, co okno wnosi do rozmowy poza samym
 * modułem. Wszystkie pola są opcjonalne: okno pokazuje to, co wie, i nie
 * zmyśla reszty.
 */
export interface KontekstOkna {
  /** Kanał modelu obsługujący okno. */
  kanalModelu?: string;
  /** Zasięg wykonania modelu. */
  srodowiskoWykonania?: ExecutionEnv;
  /** Katalogi robocze udostępnione oknu. */
  katalogiRobocze?: readonly string[];
  /** Liczba wpisów wątku; czytana przy każdym odświeżeniu panelu. */
  wpisyWatku?: () => number;
}

/** Panel kontekstu okna komunikacji, pokazujący moduł, kanał modelu, środowisko wykonania oraz wątek rozmowy w bieżącym oknie. */
export interface PanelKontekstu {
  /** Element montowany w oknie. */
  element: HTMLElement;
  /** Przestawia panel na moduł i odświeża dane okna. */
  ustaw(profil: ProfilModulu, kontekst: KontekstOkna): void;
}

/**
 * Panel kontekstu — czym okno dysponuje w bieżącym module.
 *
 * Wiersz „Wątek" pokazuje, że przy zmianie modułu historia rozmowy zostaje.
 * Licznik czyta tę samą listę wpisów, którą widzi operator, więc po przestawieniu
 * modułu liczba wpisów się nie zmienia.
 */
export function utworzPanelKontekstu(): PanelKontekstu {
  const element = document.createElement('section');
  element.className = 'dc-panel-kontekstu';
  element.setAttribute('aria-label', 'Kontekst okna komunikacji');

  const naglowek = document.createElement('h2');
  naglowek.className = 'dc-panel-kontekstu__naglowek';
  naglowek.textContent = 'Kontekst okna';

  const lista = document.createElement('dl');
  lista.className = 'dc-panel-kontekstu__lista';

  element.append(naglowek, lista);

  return {
    element,
    ustaw(profil, kontekst) {
      element.dataset['modul'] = profil.kod;
      lista.replaceChildren(...wiersze(profil, kontekst).map(([n, w]) => wiersz(n, w)));
    },
  };
}

/** Pozycje panelu w kolejności czytania: moduł, jego okna, środowisko wykonania modelu oraz licznik wpisów wątku. */
function wiersze(profil: ProfilModulu, kontekst: KontekstOkna): Array<[string, string]> {
  const pozycje: Array<[string, string]> = [['Przeznaczenie modułu', profil.przeznaczenie]];

  pozycje.push([
    'Okna kontekstu',
    profil.oknaKontekstu.length > 0
      ? profil.oknaKontekstu.join(' · ')
      : 'brak okien modułowych w rejestrze klienta',
  ]);

  if (kontekst.kanalModelu !== undefined && kontekst.kanalModelu.length > 0) {
    pozycje.push(['Kanał modelu', kontekst.kanalModelu]);
  }
  if (kontekst.srodowiskoWykonania !== undefined) {
    pozycje.push(['Zasięg wykonania', nazwaSrodowiska(kontekst.srodowiskoWykonania)]);
  }
  if (kontekst.katalogiRobocze !== undefined && kontekst.katalogiRobocze.length > 0) {
    pozycje.push(['Katalogi robocze', kontekst.katalogiRobocze.join(' · ')]);
  }
  if (kontekst.wpisyWatku !== undefined) {
    pozycje.push([
      'Wątek',
      `zachowany przy zmianie modułu — wpisów: ${kontekst.wpisyWatku()}`,
    ]);
  }
  return pozycje;
}

/** Jedna pozycja panelu kontekstu, niosąca nazwę pola oraz jego bieżącą wartość tekstową w oknie komunikacji. */
function wiersz(nazwa: string, wartosc: string): DocumentFragment {
  const fragment = document.createDocumentFragment();

  const etykieta = document.createElement('dt');
  etykieta.className = 'dc-panel-kontekstu__etykieta';
  etykieta.textContent = nazwa;

  const tresc = document.createElement('dd');
  tresc.className = 'dc-panel-kontekstu__wartosc';
  tresc.textContent = wartosc;

  fragment.append(etykieta, tresc);
  return fragment;
}
