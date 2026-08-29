/**
 * Strefa 2 — szyna nawigacji. Pionowy pas z pozycją każdego środowiska, jakie
 * rdzeń wymienił w `environment.list`; środowisko bieżące stoi rozwinięte
 * i niesie pod sobą swoje moduły. Rdzeń podaje moduły wyłącznie środowiska,
 * do którego przebieg wszedł, więc pozostałe stoją zwinięte i puste — wejście
 * w nie woła `environment.enter` i składa ramę od nowa na tym, co wróci.
 *
 * Wybór modułu okna roboczego jest osobnym terenem, więc pozycja modułu
 * wyłącznie oznacza się jako bieżąca i zmienia tytuł w belce.
 */

import type { Environment, Module } from '../../../../shared/contract.ts';
import { ikony } from '../ikony.ts';
import { el, tekst, zeZnacznika } from '../narzedzia.ts';

export interface WlasciwosciSzyny {
  /** Środowisko, do którego przebieg wszedł. */
  srodowisko: Environment;
  /** Moduły środowiska bieżącego, w kolejności odebranej od rdzenia. */
  moduly: Module[];
  /** Wszystkie środowiska rdzenia; pusty wykaz zostawia w szynie samo bieżące. */
  srodowiska?: Environment[];
  /** Wejście w środowisko inne niż bieżące; brak — pozycje pozostałych stoją nieczynne. */
  wejdz?: (srodowisko: Environment) => void;
}

/** Rysunek środowiska po kodzie; kody spoza wykazu katalogu rdzenia dostają godło produktu. */
const ZNAK_SRODOWISKA: Record<string, keyof typeof ikony> = {
  talkin: 'srodowiskoTalkin',
  workspace: 'srodowiskoWorkspace',
  codestudio: 'srodowiskoCodestudio',
  multitaskingai: 'srodowiskoMultitaskingai',
};

/**
 * Rysunek modułu po kodzie — jeden klucz na moduł z katalogu rdzenia
 * (`internal/store/migracja_007_zaczyn_slownikow.sql`), żaden klucz się nie
 * powtarza. Moduł spoza tego wykazu dostaje rysunek zastępczy, odrębny od
 * wszystkich piętnastu.
 */
const ZNAK_MODULU: Record<string, keyof typeof ikony> = {
  studio: 'olowek',
  workspace: 'pulpit',
  automations: 'automatyzacja',
  browser: 'globus',
  research: 'badanie',
  library: 'biblioteka',
  translate: 'tlumacz',
  roundtable: 'debata',
  design: 'paleta',
  assistant: 'mikrofon',
  terminal: 'terminal',
  developer: 'kod',
  diagnostics: 'diagnostyka',
  apps: 'aplikacje',
  agents: 'agenci',
};

function znak(nazwa: keyof typeof ikony): SVGElement {
  const rysunek = zeZnacznika(ikony[nazwa]);
  rysunek.setAttribute('aria-hidden', 'true');
  return rysunek;
}

function pozycjaModulu(modul: Module, biezacy: boolean): HTMLElement {
  const rysunek = ZNAK_MODULU[modul.code] ?? 'ukladanka';
  return el(
    'button',
    {
      klasa: 'dn-szyna-poz dn-szyna-poz--modul dn-etykietka',
      type: 'button',
      'data-modul-nazwa': modul.name,
      'data-modul-kod': modul.code,
      'data-etykietka': modul.name,
      'aria-label': modul.name,
      'aria-current': biezacy ? 'true' : null,
    },
    [znak(rysunek)],
  );
}

export function szyna(w: WlasciwosciSzyny): HTMLElement {
  const znakMarki = zeZnacznika(ikony.godlo);
  znakMarki.setAttribute('aria-hidden', 'true');

  /* Wykaz pusty znaczy, że rdzeń nie wymienił środowisk — zostaje samo bieżące,
     bo do niego przebieg naprawdę wszedł. */
  const wykaz = w.srodowiska !== undefined && w.srodowiska.length > 0 ? w.srodowiska : [w.srodowisko];

  const pozycje: HTMLElement[] = [];
  for (const srodowisko of [...wykaz].sort((a, b) => a.order - b.order)) {
    const biezace = srodowisko.id === w.srodowisko.id;
    const idModuly = `dn-szyna-moduly-${srodowisko.code}`;
    const pozycja = el(
      'button',
      {
        klasa: 'dn-szyna-poz dn-szyna-poz--srodowisko dn-etykietka',
        type: 'button',
        'data-srodowisko': srodowisko.code,
        'data-biezace': biezace ? 'true' : null,
        'aria-current': biezace ? 'true' : null,
        'aria-expanded': biezace ? 'true' : 'false',
        'aria-controls': idModuly,
        'data-etykietka': srodowisko.name,
        'aria-label': srodowisko.name,
      },
      [znak(ZNAK_SRODOWISKA[srodowisko.code] ?? 'godlo')],
    );
    if (!biezace && w.wejdz !== undefined) {
      pozycja.addEventListener('click', () => w.wejdz?.(srodowisko));
    }
    pozycje.push(pozycja);
    pozycje.push(
      el(
        'div',
        { klasa: 'dn-szyna-moduly', id: idModuly, role: 'group', 'aria-label': srodowisko.name },
        biezace ? w.moduly.map((modul, indeks) => pozycjaModulu(modul, indeks === 0)) : [],
      ),
    );
  }

  const lista = el('div', { klasa: 'dn-szyna-nawigacji-lista', role: 'group', 'aria-label': tekst('szyna.etykieta') }, pozycje);

  return el('nav', { klasa: 'dn-szyna-nawigacji', 'aria-label': tekst('szyna.etykieta') }, [
    el('div', { klasa: 'dn-szyna-poz', 'aria-hidden': true }, [znakMarki]),
    el('div', { klasa: 'dn-szyna-tresc' }, [lista]),
  ]);
}
