/**
 * Strefa 2 — szyna nawigacji. Pionowy pas z modułami środowiska, do którego
 * przebieg wszedł; wybór modułu okna roboczego jest osobnym terenem, więc
 * pozycja tu wyłącznie oznacza się jako bieżąca i zmienia tytuł w belce.
 */

import type { Module } from '../../../../shared/contract.ts';
import { ikony } from '../ikony.ts';
import { el, tekst, zeZnacznika } from '../narzedzia.ts';

export interface WlasciwosciSzyny {
  /** Moduły środowiska bieżącego, w kolejności odebranej od rdzenia. */
  moduly: Module[];
}

function znak(nazwa: keyof typeof ikony): SVGElement {
  const rysunek = zeZnacznika(ikony[nazwa]);
  rysunek.setAttribute('aria-hidden', 'true');
  return rysunek;
}

function pozycjaModulu(modul: Module, biezacy: boolean): HTMLElement {
  return el(
    'button',
    {
      klasa: 'dn-szyna-poz dn-szyna-poz--modul',
      type: 'button',
      'data-modul-nazwa': modul.name,
      'aria-label': modul.name,
      'aria-current': biezacy ? 'true' : null,
    },
    [znak('modul')],
  );
}

export function szyna(w: WlasciwosciSzyny): HTMLElement {
  const znakMarki = zeZnacznika(ikony.godlo);
  znakMarki.setAttribute('aria-hidden', 'true');

  const lista = el(
    'div',
    { klasa: 'dn-szyna-nawigacji-lista' },
    w.moduly.map((modul, indeks) => pozycjaModulu(modul, indeks === 0)),
  );

  return el('nav', { klasa: 'dn-szyna-nawigacji', 'aria-label': tekst('szyna.etykieta') }, [
    el('div', { klasa: 'dn-szyna-poz', 'aria-hidden': true }, [znakMarki]),
    lista,
  ]);
}
