import './przepelnienie.css';

import { elementIkony } from '../ikony/ikony';

/**
 * Przepełnienie pasa poziomego: przewijanie zamiast zwężania pozycji.
 *
 * Zwężanie odbiera czytelność wszystkim pozycjom naraz — przy kilkunastu
 * kartach każda ma tytuł ścięty do zera i żadnej nie da się rozpoznać, a wraz
 * z tytułem znika plakietka stanu, czyli informacja o awarii. Przewijanie
 * odbiera czytelność wyłącznie temu, co wyszło poza kadr. Warunkiem jest
 * podłoga szerokości pozycji w arkuszu pasa (`flex: 0 0 auto` z `min-width`);
 * bez niej pozycje kurczą się, zanim przewijak w ogóle ruszy.
 *
 * Wiązanie dokłada to, czego arkusz nie umie: utrzymanie pozycji czynnej
 * w kadrze i uchwyty krawędziowe. Uchwyty nie są bramą, tylko skrótem do tego,
 * co i tak zrobi kółko myszy albo klawiatura, i pokazują się wyłącznie wtedy,
 * gdy jest dokąd przewijać.
 */

/** Strona pasa, po której coś zostało poza kadrem. */
export type StronaPrzepelnienia = 'nie' | 'lewo' | 'prawo' | 'oba';

/** Wiązanie przepełnienia — trzyma nasłuchy, więc trzeba je odłączyć. */
export interface WiezPrzepelnienia {
  /**
   * Przelicza przepełnienie i przywraca pozycję czynną do kadru.
   *
   * Widok woła to po każdej migawce rdzenia, bo zmiana stanu potrafi uczynić
   * czynną pozycję stojącą poza kadrem.
   */
  odswiez(): void;
  /** Która strona jest dziś poza kadrem. */
  strona(): StronaPrzepelnienia;
  /** Odłącza nasłuchy i zdejmuje uchwyty. */
  rozlacz(): void;
}

/** Ustawienia wiązania; wszystkie poza `przewijak` mają wartość domyślną. */
export interface OpcjePrzepelnienia {
  /** Element z `overflow-x: auto` — pas, który przewija swoją zawartość. */
  przewijak: HTMLElement;
  /**
   * Element, który ma zostać w kadrze — zwykle karta czynna.
   *
   * `null` znaczy „nic nie trzeba pilnować"; pas pusty albo bez wskazania.
   */
  wKadrze?: () => HTMLElement | null;
  /**
   * Gdzie stawiać uchwyty krawędziowe. Domyślnie rodzic przewijaka, bo uchwyt
   * przyklejony do przewijanej treści jechałby razem z nią.
   */
  gospodarz?: HTMLElement;
  /** Nazwa pasa w etykietach uchwytów — „Karty sesji" daje „…kart sesji". */
  nazwaPasa?: string;
}

/** Ile procent kadru przesuwa jedno naciśnięcie uchwytu. */
const KROK_UCHWYTU = 0.8;

export function zwiazPrzepelnieniePasa(opcje: OpcjePrzepelnienia): WiezPrzepelnienia {
  const { przewijak } = opcje;
  const gospodarz = opcje.gospodarz ?? przewijak.parentElement ?? przewijak;
  const nazwa = opcje.nazwaPasa ?? 'pasa';
  const wKadrze = opcje.wKadrze ?? (() => null);

  // Klasa własna zamiast gołego atrybutu: reguły cienia krawędzi mają
  // obowiązywać pasy prowadzone tym wiązaniem, a nie każdy element z podobnym
  // atrybutem.
  przewijak.classList.add('dnp-przepelnienie__przewijak');

  const wLewo = utworzUchwyt('lewo', nazwa, () => przesun(-1));
  const wPrawo = utworzUchwyt('prawo', nazwa, () => przesun(1));

  // Uchwyty obejmują przewijak z obu stron, żeby kolejność czytania zgadzała
  // się z kierunkiem, który wskazują.
  gospodarz.insertBefore(wLewo, przewijak);
  if (przewijak.nextSibling === null) gospodarz.append(wPrawo);
  else gospodarz.insertBefore(wPrawo, przewijak.nextSibling);

  let biezaca: StronaPrzepelnienia = 'nie';

  function przesun(kierunek: 1 | -1): void {
    const krok = Math.max(1, Math.round(przewijak.clientWidth * KROK_UCHWYTU));
    // `scrollBy` bywa nieobecne w środowiskach bez rozkładu — wtedy przewijamy
    // wprost właściwością, zamiast wywracać widok.
    if (typeof przewijak.scrollBy === 'function') {
      przewijak.scrollBy({ left: kierunek * krok, behavior: 'smooth' });
    } else {
      przewijak.scrollLeft += kierunek * krok;
    }
    przelicz();
  }

  /** Która strona pasa została poza kadrem. */
  function policzStrone(): StronaPrzepelnienia {
    // Zaokrąglenia przeglądarki potrafią zostawić ułamek piksela; jeden piksel
    // luzu oszczędza uchwyt migoczący na pasie, który mieści się dokładnie.
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
    // `nearest` w obu osiach: pozycja wjeżdża do kadru najkrótszą drogą i nie
    // rusza przewijania pionowego całej powłoki.
    cel.scrollIntoView({ block: 'nearest', inline: 'nearest' });
  }

  const naPrzewiniecie = (): void => przelicz();
  przewijak.addEventListener('scroll', naPrzewiniecie, { passive: true });

  // Zmiana szerokości okna zmienia kadr, a nie treść — bez tego uchwyty
  // zostałyby widoczne na pasie, który już się mieści.
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

/** Uchwyt krawędziowy — skrót do przewinięcia, nie brama. */
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
  // Uchwyt jest skrótem dla myszy; wędrówka klawiaturowa pasa i tak dowozi
  // ognisko do pozycji poza kadrem, więc drugi przystanek w kolejności
  // tabulacji byłby przeszkodą, nie ułatwieniem.
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
