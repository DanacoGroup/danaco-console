/**
 * Rama aplikacji — montaż. Składa trzy strefy z prototypu — szynę nawigacji,
 * belkę tytułową i pas stanu — i wstawia je w miejsce wskazane w dokumencie.
 * Ta sama zasada co w drodze wejścia: montaż jest jedynym miejscem, które
 * dotyka dokumentu.
 */

import type { Environment, Module, Session } from '../../../shared/contract.ts';
import { el, tekst } from './narzedzia.ts';
import { belka } from './skladniki/belka.ts';
import { szyna } from './skladniki/szyna.ts';
import { stan as pasStanu } from './skladniki/stan.ts';

export interface NastawyRamy {
  /** Miejsce w dokumencie, w które rama się wstawia. */
  miejsce: HTMLElement;
  /** Środowisko, do którego przebieg wszedł. */
  srodowisko: Environment;
  /** Moduły tego środowiska, w kolejności odebranej od rdzenia. */
  moduly: Module[];
  /** Karty sesji odtworzone przez rdzeń — ich liczba zasila pas stanu. */
  sesje: Session[];
}

/** Zamontowana rama wraz z drogą jej zdjęcia: odłączeniem słuchaczy zdarzeń. */
export interface ZamontowanaRama {
  zdejmij(): void;
}

/** Zapytanie o preferencję systemową motywu — bez odstępu po dwukropku, składnia CSS przyjmuje oba zapisy. */
const ZAPYTANIE_MOTYW_CIEMNY = '(prefers-color-scheme:dark)';

function motywCiemny(): boolean {
  const jawny = document.documentElement.dataset['theme'];
  if (jawny === 'dark') return true;
  if (jawny === 'light') return false;
  return globalThis.matchMedia?.(ZAPYTANIE_MOTYW_CIEMNY).matches === true;
}

export function zamontujRame(w: NastawyRamy): ZamontowanaRama {
  const glowna = el('main', { 'aria-label': tekst('glowna.etykieta') });
  const belkaWezel = belka({ srodowisko: w.srodowisko.name });
  let stanWezel = pasStanu({
    srodowisko: w.srodowisko.name,
    liczbaSesji: w.sesje.length,
    motywCiemny: motywCiemny(),
  });

  const powloka = el('div', { klasa: 'sta-powloka' }, [
    el('div', { klasa: 'dn-rama-korpus' }, [
      szyna({ moduly: w.moduly }),
      el('div', { klasa: 'dn-rama-prawa' }, [belkaWezel, glowna]),
    ]),
    stanWezel,
  ]);

  w.miejsce.replaceChildren(powloka);

  /* Wybór modułu w szynie jest sprawą samej ramy: modyfikuje jej własny
     tytuł, nie otwiera okna modułowego — to osobny teren. */
  function naKlikniecie(zdarzenie: Event): void {
    const przycisk = (zdarzenie.target as Element | null)?.closest(
      '.dn-szyna-poz--modul',
    ) as HTMLElement | null;
    if (przycisk === null) return;
    for (const inny of w.miejsce.querySelectorAll('.dn-szyna-poz--modul')) {
      inny.removeAttribute('aria-current');
    }
    przycisk.setAttribute('aria-current', 'true');
    const nazwa = przycisk.dataset['modulNazwa'] ?? w.srodowisko.name;
    belkaWezel.querySelector('[data-belka-tytul]')!.textContent =
      `${tekst('belka.marka')} ${tekst('belka.separator')} ${nazwa}`;
  }
  w.miejsce.addEventListener('click', naKlikniecie);

  const zapytanieMotywu = globalThis.matchMedia?.(ZAPYTANIE_MOTYW_CIEMNY);
  function naZmianeMotywu(): void {
    const nowy = pasStanu({
      srodowisko: w.srodowisko.name,
      liczbaSesji: w.sesje.length,
      motywCiemny: motywCiemny(),
    });
    stanWezel.replaceWith(nowy);
    stanWezel = nowy;
  }
  zapytanieMotywu?.addEventListener('change', naZmianeMotywu);

  return {
    zdejmij() {
      w.miejsce.removeEventListener('click', naKlikniecie);
      zapytanieMotywu?.removeEventListener('change', naZmianeMotywu);
    },
  };
}
