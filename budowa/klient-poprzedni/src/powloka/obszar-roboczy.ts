import './obszar-roboczy.css';

import { elementIkony } from '../ikony/ikony';
import type { PozycjaModulu, Srodowisko } from './srodowiska';

/**
 * Pas 4 powłoki — obszar roboczy.
 *
 * Jedna odpowiedzialność: kontener na treść wybranej pozycji nawigacji wraz
 * z uczciwym stanem pustym dla pozycji, która zbudowanego widoku nie ma.
 * Obszar nie zna żadnego widoku i żadnego rdzenia — przyjmuje gotowy element
 * od warstwy, która go składa.
 *
 * Plansza obszaru niesie dwie warstwy naraz: widok modułu w wierszu górnym
 * i scenę okien komunikacji pod nim. Okno rozmowy jest oknem wiodącym każdego
 * modułu, więc obszar pokazujący jedną warstwę naraz odbierałby rozmowę
 * każdemu modułowi, który doczekał się własnego widoku.
 *
 * Widok raz osadzony nie jest niszczony: `pokaz` dokłada go przy pierwszym
 * użyciu i dalej wyłącznie przestawia widoczność atrybutem `hidden`. W obszarze
 * stoi scena z żywym oknem rozmowy — gniazdem WebSocket, strumieniem odpowiedzi
 * modelu, historią wpisów i obserwatorami układu okien równoległych.
 * Odmontowanie sceny przy każdym przełączeniu modułu zrywałoby rozmowę
 * w połowie zdania; ukrycie zostawia ją nietkniętą.
 *
 * Stan pusty nie zapowiada modułu, którego nie ma — nazywa go i pisze wprost,
 * że jego okna operacyjne nie zostały zbudowane, wymieniając katalog kodów
 * wzięty z rdzenia.
 */
export interface ObszarRoboczy {
  /** Kontener montowany w korpusie powłoki, obok nawigacji. */
  element: HTMLElement;
  /**
   * Osadza scenę okien komunikacji — warstwę wspólną wszystkim modułom.
   * Wywoływana raz; scena zostaje w drzewie do końca życia powłoki.
   */
  osadzCzat(scena: HTMLElement): void;
  /**
   * Odsłania widok modułu nad sceną czatu; `null` zostawia sam czat.
   * Przy pierwszym użyciu dokłada widok i już go nie zdejmuje.
   */
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

  // Plansza — dwa wiersze. Widok modułu na górze, scena okien komunikacji pod
  // nim: rozmowa jest oknem wiodącym każdego modułu, więc bierze dolny wiersz
  // i nigdy nie schodzi z planszy przez wybór modułu.
  const plansza = document.createElement('div');
  plansza.className = 'dn-obszar__plansza';
  plansza.hidden = true;

  const strefaModulu = document.createElement('div');
  strefaModulu.className = 'dn-obszar__modul';
  strefaModulu.hidden = true;

  const strefaCzatu = document.createElement('div');
  strefaCzatu.className = 'dn-obszar__czat';

  // Zapowiedź stoi w strefie modułu, nie zamiast planszy: gdyby gasiła całą
  // planszę, gasiłaby razem z nią scenę rozmowy, a rozmowa jest oknem wiodącym
  // każdej pozycji i nie schodzi z planszy przez brak widoku modułu.
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
      // Widok modułu bywa jeden z wielu osadzonych — odsłaniamy wskazany,
      // resztę chowamy. Scena czatu nie jest jedną z tych warstw i zostaje.
      for (const osadzony of widoki) osadzony.hidden = osadzony !== widok;
      strefaModulu.hidden = false;
      plansza.hidden = false;
      zapowiedzElement.hidden = true;
    },

    oproznij() {
      // Widok modułu schodzi, plansza zostaje: pozycja bez widoku mówi o tym
      // wprost w strefie modułu, a rozmowa pracuje dalej obok niej.
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

/** Kod okna operacyjnego z katalogu rdzenia jako wiersz wykazu braków. */
function wpisOkna(kod: string): HTMLLIElement {
  const punkt = document.createElement('li');
  punkt.textContent = kod;
  return punkt;
}

/** Zdanie stanu pustego: rola pozycji plus prawda o braku okien. */
function zdanieStanu(pozycja: PozycjaModulu, okna: readonly string[]): string {
  const rola = pozycja.opis === '' ? '' : `${wielkaLitera(pozycja.opis)}. `;
  if (okna.length === 0) {
    return `${rola}Ta pozycja nie ma jeszcze zbudowanego widoku — rdzeń nie przypisał jej okien operacyjnych.`;
  }
  return `${rola}Okna operacyjne tego modułu nie zostały jeszcze zbudowane. Rdzeń podaje ich katalog:`;
}

/** Opis modułu z rdzenia zaczyna się małą literą — stan pusty jest zdaniem. */
function wielkaLitera(tekst: string): string {
  return tekst.charAt(0).toUpperCase() + tekst.slice(1);
}
