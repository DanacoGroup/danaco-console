/**
 * Strefa 2 — szyna nawigacji. Pionowy pas ze środowiskiem bieżącym u góry
 * listy i modułami tego środowiska pod nim; wybór modułu okna roboczego jest
 * osobnym terenem, więc pozycja tu wyłącznie oznacza się jako bieżąca
 * i zmienia tytuł w belce.
 */

import type { Environment, Module } from '../../../../shared/contract.ts';
import { ikony } from '../ikony.ts';
import { el, tekst, zeZnacznika } from '../narzedzia.ts';

export interface WlasciwosciSzyny {
  /** Środowisko, do którego przebieg wszedł. */
  srodowisko: Environment;
  /** Moduły środowiska bieżącego, w kolejności odebranej od rdzenia. */
  moduly: Module[];
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

  const rysunekSrodowiska = ZNAK_SRODOWISKA[w.srodowisko.code] ?? 'godlo';
  const idModuly = 'dn-szyna-moduly';

  const pozycjaSrodowiska = el(
    'button',
    {
      klasa: 'dn-szyna-poz dn-szyna-poz--srodowisko dn-etykietka',
      type: 'button',
      'data-srodowisko': w.srodowisko.code,
      'aria-current': 'true',
      'aria-expanded': 'true',
      'aria-controls': idModuly,
      'data-etykietka': w.srodowisko.name,
      'aria-label': w.srodowisko.name,
    },
    [znak(rysunekSrodowiska)],
  );

  const moduly = el(
    'div',
    { klasa: 'dn-szyna-moduly', id: idModuly, role: 'group', 'aria-label': w.srodowisko.name },
    w.moduly.map((modul, indeks) => pozycjaModulu(modul, indeks === 0)),
  );

  const lista = el('div', { klasa: 'dn-szyna-nawigacji-lista', role: 'group', 'aria-label': tekst('szyna.etykieta') }, [
    pozycjaSrodowiska,
    moduly,
  ]);

  return el('nav', { klasa: 'dn-szyna-nawigacji', 'aria-label': tekst('szyna.etykieta') }, [
    el('div', { klasa: 'dn-szyna-poz', 'aria-hidden': true }, [znakMarki]),
    el('div', { klasa: 'dn-szyna-tresc' }, [lista]),
  ]);
}
