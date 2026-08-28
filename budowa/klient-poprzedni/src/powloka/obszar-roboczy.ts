import './obszar-roboczy.css';

import { elementIkony } from '../ikony/ikony';
import type { PozycjaModulu, Srodowisko } from './srodowiska';

/**
 * Pas 4 powłoki — obszar roboczy: kontener treści wybranej pozycji nawigacji, ze sceną rozmowy
 * i uczciwym stanem pustym.
 */
export interface ObszarRoboczy {
  /** Kontener montowany w korpusie powłoki, obok nawigacji. */
  element: HTMLElement;
  /** Osadza scenę okien komunikacji — warstwę wspólną wszystkim modułom, montowaną raz na stałe. */
  osadzCzat(scena: HTMLElement): void;
  /** Odsłania widok modułu nad sceną czatu; `null` zostawia sam czat, widok osadzony raz nie znika. */
  pokaz(widok: HTMLElement | null): void;
  /** Chowa planszę w całości i odsłania stan pusty. */
  oproznij(): void;
  /** Ustawia treść stanu pustego dla wskazanej pozycji nawigacji. */
  zapowiedz(pozycja: PozycjaModulu, dane: Srodowisko, okna?: readonly string[]): void;
  /** Pasek uczciwości nad widokiem; `null` go chowa. */
  nota(tresc: string | null): void;
}

export function utworzObszarRoboczy(): ObszarRoboczy {
  /** Widoki modułów osadzone na stałe — dokładane, nigdy zdejmowane. */
  const widoki = new Set<HTMLElement>();

  const element = document.createElement('main');
  element.className = 'dn-obszar';
  element.setAttribute('aria-label', 'Obszar roboczy');

  const notaElement = document.createElement('p');
  notaElement.className = 'dn-obszar__nota';
  notaElement.setAttribute('role', 'status');
  notaElement.hidden = true;

  const zapowiedzElement = document.createElement('div');
  zapowiedzElement.className = 'dn-pusty-stan dn-obszar__zapowiedz';

  const znak = document.createElement('span');
  znak.className = 'dn-obszar__znak';

  const tytul = document.createElement('p');
  tytul.className = 'dn-pusty-stan-tytul';
  tytul.textContent = 'Obszar roboczy';

  const opis = document.createElement('p');
  opis.className = 'dn-pusty-stan-opis';
  opis.textContent = 'Wybierz moduł w bocznej nawigacji — jego treść wejdzie tutaj.';

  const wykaz = document.createElement('ul');
  wykaz.className = 'dn-obszar__braki';
  wykaz.hidden = true;

  zapowiedzElement.append(znak, tytul, opis, wykaz);

  // Plansza ma dwa wiersze: widok modułu na górze, scenę rozmowy — okno wiodące modułu — na dole.
  const plansza = document.createElement('div');
  plansza.className = 'dn-obszar__plansza';
  plansza.hidden = true;

  const strefaModulu = document.createElement('div');
  strefaModulu.className = 'dn-obszar__modul';
  strefaModulu.hidden = true;

  const strefaCzatu = document.createElement('div');
  strefaCzatu.className = 'dn-obszar__czat';

  // Zapowiedź stoi w strefie modułu, nie zamiast planszy, bo plansza niesie też scenę rozmowy.
  strefaModulu.append(zapowiedzElement);
  plansza.append(strefaModulu, strefaCzatu);
  element.append(notaElement, plansza);

  return {
    element,

    osadzCzat(scena) {
      if (scena.parentElement !== strefaCzatu) strefaCzatu.append(scena);
    },

    pokaz(widok) {
      if (widok !== null && !widoki.has(widok)) {
        widoki.add(widok);
        strefaModulu.append(widok);
      }
      // Widok modułu bywa jeden z wielu osadzonych — odsłaniamy wskazany, resztę chowamy; czat zostaje.
      for (const osadzony of widoki) osadzony.hidden = osadzony !== widok;
      strefaModulu.hidden = false;
      plansza.hidden = false;
      zapowiedzElement.hidden = true;
    },

    oproznij() {
      // Widok modułu schodzi, plansza zostaje: brak widoku mówi o tym w strefie, rozmowa pracuje dalej.
      for (const osadzony of widoki) osadzony.hidden = true;
      strefaModulu.hidden = false;
      plansza.hidden = false;
      zapowiedzElement.hidden = false;
    },

    zapowiedz(pozycja, dane, okna = []) {
      znak.replaceChildren(elementIkony(pozycja.ikona, { rozmiar: 24 }));
      tytul.textContent = `${dane.nazwa} · ${pozycja.nazwa}`;
      opis.textContent = zdanieStanu(pozycja, okna);
      wykaz.replaceChildren(...okna.map(wpisOkna));
      wykaz.hidden = okna.length === 0;
    },

    nota(tresc) {
      notaElement.textContent = tresc ?? '';
      notaElement.hidden = tresc === null || tresc === '';
    },
  };
}

/** Kod okna operacyjnego z katalogu rdzenia, wypisywany jako pojedynczy wiersz wykazu brakujących okien. */
function wpisOkna(kod: string): HTMLLIElement {
  const punkt = document.createElement('li');
  punkt.textContent = kod;
  return punkt;
}

/** Zdanie stanu pustego pozycji nawigacji, złożone z roli pozycji oraz prawdy o braku zbudowanych okien. */
function zdanieStanu(pozycja: PozycjaModulu, okna: readonly string[]): string {
  const rola = pozycja.opis === '' ? '' : `${wielkaLitera(pozycja.opis)}. `;
  if (okna.length === 0) {
    return `${rola}Ta pozycja nie ma jeszcze zbudowanego widoku — rdzeń nie przypisał jej okien operacyjnych.`;
  }
  return `${rola}Okna operacyjne tego modułu nie zostały jeszcze zbudowane. Rdzeń podaje ich katalog:`;
}

/** Zamienia pierwszą literę opisu modułu z rdzenia na wielką, bo stan pusty składa opis jako pełne zdanie. */
function wielkaLitera(tekst: string): string {
  return tekst.charAt(0).toUpperCase() + tekst.slice(1);
}
