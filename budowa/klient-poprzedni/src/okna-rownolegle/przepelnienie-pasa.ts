import './przepelnienie.css';

import { elementIkony } from '../ikony/ikony';

/**
 * Przepełnienie pasa: przewijanie zamiast zwężania.
 */

/** Strona pasa poziomego, po której coś zostało poza widocznym kadrem tego przewijanego paska z treścią. */
export type StronaPrzepelnienia = 'nie' | 'lewo' | 'prawo' | 'oba';

/**
 * Wiązanie przepełnienia tego pasa trzyma nasłuchy zdarzeń, więc trzeba je jawnie odłączyć po ich użyciu.
 */
export interface WiezPrzepelnienia {
  /**
   * Przelicza przepełnienie i przywraca pozycję czynną do kadru po migawce rdzenia.
   */
  odswiez(): void;
  /** Która strona jest dziś poza kadrem. */
  strona(): StronaPrzepelnienia;
  /** Odłącza nasłuchy i zdejmuje uchwyty. */
  rozlacz(): void;
}

/**
 * Ustawienia wiązania przepełnienia tego pasa; wszystkie poza elementem przewijaka mają wartość domyślną.
 */
export interface OpcjePrzepelnienia {
  /** Element z `overflow-x: auto` — pas, który przewija swoją zawartość. */
  przewijak: HTMLElement;
  /**
   * Element, który ma zostać w kadrze; null znaczy, że nic nie trzeba pilnować.
   */
  wKadrze?: () => HTMLElement | null;
  /**
   * Gdzie stawiać uchwyty krawędziowe; domyślnie rodzic przewijaka.
   */
  gospodarz?: HTMLElement;
  /** Nazwa pasa w etykietach uchwytów — „Karty sesji" daje „…kart sesji". */
  nazwaPasa?: string;
}

/**
 * Ile procent szerokości kadru przesuwa jedno naciśnięcie uchwytu krawędziowego tego pasa przewijania.
 */
const KROK_UCHWYTU = 0.8;

export function zwiazPrzepelnieniePasa(opcje: OpcjePrzepelnienia): WiezPrzepelnienia {
  const { przewijak } = opcje;
  const gospodarz = opcje.gospodarz ?? przewijak.parentElement ?? przewijak;
  const nazwa = opcje.nazwaPasa ?? 'pasa';
  const wKadrze = opcje.wKadrze ?? (() => null);

  // Klasa własna zamiast atrybutu: cień krawędzi ma dotyczyć tylko pasów tego wiązania.
  przewijak.classList.add('dnp-przepelnienie__przewijak');

  const wLewo = utworzUchwyt('lewo', nazwa, () => przesun(-1));
  const wPrawo = utworzUchwyt('prawo', nazwa, () => przesun(1));

  // Uchwyty obejmują przewijak z obu stron zgodnie z kierunkiem, który wskazują.
  gospodarz.insertBefore(wLewo, przewijak);
  if (przewijak.nextSibling === null) gospodarz.append(wPrawo);
  else gospodarz.insertBefore(wPrawo, przewijak.nextSibling);

  let biezaca: StronaPrzepelnienia = 'nie';

  function przesun(kierunek: 1 | -1): void {
    const krok = Math.max(1, Math.round(przewijak.clientWidth * KROK_UCHWYTU));
    // scrollBy bywa nieobecne bez rozkładu; wtedy przewijamy wprost właściwością.
    if (typeof przewijak.scrollBy === 'function') {
      przewijak.scrollBy({ left: kierunek * krok, behavior: 'smooth' });
    } else {
      przewijak.scrollLeft += kierunek * krok;
    }
    przelicz();
  }

  /** Która strona pasa została poza kadrem. */
  function policzStrone(): StronaPrzepelnienia {
    // Zaokrąglenia przeglądarki zostawiają ułamek piksela; luz chroni przed migotaniem.
    const zapas = przewijak.scrollWidth - przewijak.clientWidth;
    if (zapas <= 1) return 'nie';
    const lewo = przewijak.scrollLeft > 1;
    const prawo = przewijak.scrollLeft < zapas - 1;
    if (lewo && prawo) return 'oba';
    if (lewo) return 'lewo';
    if (prawo) return 'prawo';
    return 'nie';
  }

  function przelicz(): void {
    biezaca = policzStrone();
    przewijak.dataset.przepelnienie = biezaca;
    wLewo.hidden = biezaca !== 'lewo' && biezaca !== 'oba';
    wPrawo.hidden = biezaca !== 'prawo' && biezaca !== 'oba';
  }

  function przywrocDoKadru(): void {
    const cel = wKadrze();
    if (cel === null) return;
    // Środowisko bez rozkładu nie ma `scrollIntoView`; brak metody nie ma
    // wywracać widoku.
    if (typeof cel.scrollIntoView !== 'function') return;
    // Tryb nearest w obu osiach: pozycja wjeżdża do kadru, nie rusza przewijania pionowego.
    cel.scrollIntoView({ block: 'nearest', inline: 'nearest' });
  }

  const naPrzewiniecie = (): void => przelicz();
  przewijak.addEventListener('scroll', naPrzewiniecie, { passive: true });

  // Zmiana szerokości okna zmienia kadr, nie treść; bez tego uchwyty byłyby zbędne.
  const obserwator =
    typeof ResizeObserver === 'function' ? new ResizeObserver(() => przelicz()) : null;
  obserwator?.observe(przewijak);

  przelicz();

  return {
    odswiez() {
      przywrocDoKadru();
      przelicz();
    },
    strona: () => biezaca,
    rozlacz() {
      przewijak.removeEventListener('scroll', naPrzewiniecie);
      obserwator?.disconnect();
      przewijak.classList.remove('dnp-przepelnienie__przewijak');
      delete przewijak.dataset.przepelnienie;
      wLewo.remove();
      wPrawo.remove();
    },
  };
}

/**
 * Uchwyt krawędziowy pasa jest skrótem do przewinięcia treści, a nie osobną bramą sterującą tym pasem.
 */
function utworzUchwyt(
  strona: 'lewo' | 'prawo',
  nazwaPasa: string,
  przyNacisnieciu: () => void,
): HTMLButtonElement {
  const przycisk = document.createElement('button');
  przycisk.type = 'button';
  przycisk.className = `dn-btn-ikona dnp-przepelnienie__uchwyt dnp-przepelnienie__uchwyt--${strona}`;
  przycisk.dataset.uchwyt = strona;
  const opis = strona === 'lewo' ? `Przewiń ${nazwaPasa} w lewo` : `Przewiń ${nazwaPasa} w prawo`;
  przycisk.setAttribute('aria-label', opis);
  przycisk.title = opis;
  // Uchwyt jest skrótem dla myszy; wędrówka klawiaturowa i tak dowozi ognisko poza kadr.
  przycisk.tabIndex = -1;
  przycisk.append(elementIkony(strona === 'lewo' ? 'strzalka-lewo' : 'strzalka-prawo', {
    rozmiar: 16,
  }));
  przycisk.hidden = true;
  przycisk.addEventListener('click', (zdarzenie) => {
    zdarzenie.stopPropagation();
    przyNacisnieciu();
  });
  return przycisk;
}
