/**
 * Składnik — wykaz etapów. Wykaz kroków ze stanem i miarą po prawej; znak
 * kroku idzie za stanem, zależnie od tego, czy krok jest gotowy, w toku,
 * nieudany, czy czeka.
 */

import { ikony } from '../ikony.ts';
import { el, tekst, wykaz, zeZnacznika, type DanePodstawienia } from '../narzedzia.ts';

/** Stan jednego kroku wykazu, rozstrzygający, jaki znak stanie przy tym kroku na całej liście jego etapów. */
export type StanKroku = 'gotowy' | 'pracuje' | 'blad' | 'oczekuje';

/** Miara stojąca po prawej stronie kroku, opisana kluczem katalogu i danymi do podstawienia w jej treść. */
export interface MiaraKroku {
  /** Klucz miary w gałęzi `miary` katalogu. */
  klucz: string;
  /** Dane podstawiane w miarę. */
  dane?: DanePodstawienia;
}

export interface WlasciwosciWykazu {
  /** Gałąź katalogu z nazwami kroków, na przykład `uruchomienie.etapy`. */
  nazwy: string;
  /** Gałąź katalogu z miarami, na przykład `uruchomienie.stany`. */
  miary: string;
  /** Stan każdego kroku, po jednym na krok. */
  stany: StanKroku[];
  /** Miara każdego kroku, po jednej na krok. */
  wartosciMiar: MiaraKroku[];
  /** Klucz katalogu — nazwa obszaru dla czytnika ekranu. */
  obszar: string;
}

function znakKroku(stan: StanKroku, numer: number): Node {
  if (stan === 'pracuje') {
    return el('span', { klasa: 'dn-kropka dn-kropka--tetno', 'aria-hidden': 'true' });
  }
  if (stan === 'gotowy' || stan === 'blad') {
    const rysunek = zeZnacznika(stan === 'gotowy' ? ikony.ptaszek : ikony.krzyzyk);
    rysunek.setAttribute('aria-hidden', 'true');
    return rysunek;
  }
  return document.createTextNode(String(numer));
}

/** Odmiana kroku właściwa stanowi; stan oczekujący nie ma odmiany i zostaje na rysunku podstawowym. */
function odmianaKroku(stan: StanKroku): string {
  if (stan === 'pracuje') return ' dn-krok--pracuje';
  if (stan === 'gotowy') return ' dn-krok--poprawny';
  if (stan === 'blad') return ' dn-krok--wstrzymany';
  return '';
}

export function listaEtapow(w: WlasciwosciWykazu): HTMLElement {
  const kroki = wykaz(w.nazwy).map((nazwa, i) => {
    const stan = w.stany[i] ?? 'oczekuje';
    const miara = w.wartosciMiar[i] ?? { klucz: 'oczekuje' };
    return el('li', { klasa: `dn-krok dn-krok--pole${odmianaKroku(stan)}`, dane: { stan } }, [
      el('span', { klasa: 'dn-krok-znak', 'aria-hidden': 'true' }, [znakKroku(stan, i + 1)]),
      el('span', { tekst: nazwa }),
      el('span', {
        klasa: 'dn-krok-meta',
        tekst: tekst(`${w.miary}.${miara.klucz}`, miara.dane),
      }),
    ]);
  });

  return el(
    'ol',
    { klasa: 'dn-kolejka dn-kolejka--pola', 'aria-live': 'polite', 'aria-label': tekst(w.obszar) },
    kroki,
  );
}
